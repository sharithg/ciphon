package ws

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"sort"

	"github.com/gorilla/websocket"
	"github.com/sharithg/siphon/internal/utils"
)

type BaseEvent struct {
	Type string `json:"type"`
}

type Command struct {
	Id      string `json:"id"`
	Cmd     string `json:"cmd"`
	Order   int    `json:"order"`
	WorkDir string `json:"workDir"`
}

type Commands struct {
	BaseEvent
	Image    string    `json:"image"`
	Commands []Command `json:"commands"`
}

type CommandOutput struct {
	OutputType string `json:"outputType"`
	Output     string `json:"output"`
	CmdType    string `json:"cmdType"`
	Id         string `json:"id"`
}

func (app *Application) serveWs(w http.ResponseWriter, r *http.Request) {

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade to WebSocket: %v", err)
		return
	}
	defer ws.Close()

	log.Printf("Client connected from %s", ws.RemoteAddr())

	for {
		var event Commands
		err := ws.ReadJSON(&event)
		if err != nil {
			log.Printf("Receive failed: %s", err.Error())
			break
		}
		switch event.Type {
		case "run_command":
			go func() {
				if err := app.handleWebsocketEvent1(ws, event); err != nil {
					log.Println("Error writing to stdin:", err)
				}
			}()
		default:
			log.Printf("Unknown event type: %s", event.Type)
		}
	}

}

func (a *Application) handleWebsocketEvent1(conn *websocket.Conn, e Commands) error {
	fmt.Println("Starting Docker operation")

	ctx := context.Background()

	runId := utils.RandStringBytes(10)

	err := a.Store.Containers.Set(runId, "running")

	if err != nil {
		slog.Error("error saving container state", "err", err)
		sendError(conn, err)
		return nil
	}

	outputChan := make(chan CommandOutput)

	// setup
	stdoutSetupFunc := stdoutHandler(outputChan, "setup", "")
	stderrSetupFunc := stderrHandler(outputChan, "setup", "")

	// terdown
	stdoutTeardownFunc := stdoutHandler(outputChan, "teardown", "")
	stderrTeardownFunc := stderrHandler(outputChan, "teardown", "")

	terdownFunc := func() {
		running := CommandOutput{
			CmdType: "running",
		}
		if err := sendOutput(conn, running); err != nil {
			slog.Error("error sending output over websocket", "err", err)
			return
		}

		if err := a.Docker.StopAndRemoveContainer(ctx, runId, stdoutTeardownFunc, stderrTeardownFunc); err != nil {
			slog.Error("error stopping and deleting containers", "err", err)
			sendError(conn, err)
			return
		}

		doneCmd := CommandOutput{
			CmdType: "doneCmd",
		}
		if err := sendOutput(conn, doneCmd); err != nil {
			slog.Error("error sending output over websocket", "err", err)
			return
		}

		output := CommandOutput{
			CmdType: "done",
		}

		if err := sendOutput(conn, output); err != nil {
			slog.Error("error sending output over websocket", "err", err)
		}
	}

	go func() {
		defer close(outputChan)

		running := CommandOutput{
			CmdType: "running",
		}

		if err := sendOutput(conn, running); err != nil {
			slog.Error("error sending output over websocket", "err", err)
			return
		}

		if err := a.Docker.PullImageAndStreamOutput(ctx, e.Image, stdoutSetupFunc, stderrSetupFunc); err != nil {
			slog.Error("error pulling docker image", "err", err)
			sendError(conn, err)
			return
		}

		if err := sendOutput(conn, running); err != nil {
			slog.Error("error sending output over websocket", "err", err)
			return
		}

		if err := a.Docker.RunBackgroundContainer(runId, e.Image, stdoutSetupFunc, stderrSetupFunc); err != nil {
			slog.Error("error running background container", "err", err)
			sendError(conn, err)
			return
		}

		doneCmd := CommandOutput{
			CmdType: "doneCmd",
		}
		if err := sendOutput(conn, doneCmd); err != nil {
			slog.Error("error sending output over websocket", "err", err)
			return
		}

		defer terdownFunc()

		commands := e.Commands

		sort.Slice(commands, func(i, j int) bool {
			return commands[i].Order < commands[j].Order
		})

		for _, cmd := range e.Commands {

			running := CommandOutput{
				CmdType: "running",
				Id:      cmd.Id,
			}

			stdoutCmdFunc := stdoutHandler(outputChan, "cmd", cmd.Id)
			stderrCmdFunc := stderrHandler(outputChan, "cmd", cmd.Id)

			workDir := "/"

			if cmd.WorkDir != "" {
				workDir = cmd.WorkDir
			}

			if err := sendOutput(conn, running); err != nil {
				slog.Error("error sending output over websocket", "err", err)
				return
			}

			if err := a.Docker.ExecAndStreamLogs(runId, workDir, cmd.Cmd, stdoutCmdFunc, stderrCmdFunc); err != nil {
				slog.Error("error running command", "err", err)
				sendError(conn, err)
				return
			}
			doneCmd := CommandOutput{
				CmdType: "doneCmd",
				Id:      cmd.Id,
			}
			if err := sendOutput(conn, doneCmd); err != nil {
				slog.Error("error sending output over websocket", "err", err)
				return
			}
		}

	}()

	for output := range outputChan {
		if err := sendOutput(conn, output); err != nil {
			slog.Error("error sending output over websocket", "err", err)
			return err
		}
	}

	err = a.Store.Containers.Set(runId, "complete")

	if err != nil {
		slog.Error("error saving container state", "err", err)
		sendError(conn, err)
		return nil
	}

	return nil
}

func stdoutHandler(outputChan chan CommandOutput, cmdType string, cmdId string) func(string) {
	return func(message string) {
		outputChan <- CommandOutput{
			OutputType: "stdout",
			Output:     message,
			CmdType:    cmdType,
			Id:         cmdId,
		}
	}
}

func stderrHandler(outputChan chan CommandOutput, cmdType string, cmdId string) func(string) {
	return func(message string) {
		outputChan <- CommandOutput{
			OutputType: "stderr",
			Output:     message,
			CmdType:    cmdType,
			Id:         cmdId,
		}
	}
}

func sendOutput(conn *websocket.Conn, output CommandOutput) error {
	if err := conn.WriteJSON(output); err != nil {
		return fmt.Errorf("failed to send %s message: %v", output.OutputType, err)
	}
	return nil
}

func sendError(conn *websocket.Conn, errMsg error) error {
	err := conn.WriteJSON(CommandOutput{
		CmdType: "error",
		Output:  errMsg.Error(),
	})
	if err != nil {
		return fmt.Errorf("can't send: %s", err.Error())
	}

	output := CommandOutput{
		CmdType: "done",
	}

	if err := sendOutput(conn, output); err != nil {
		slog.Error("error sending output over websocket", "err", err)
	}

	return nil
}
