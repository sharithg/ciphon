package api

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/sharithg/siphon/internal/remote"
	"github.com/sharithg/siphon/internal/repository"
)

type Node struct {
	Id       string `json:"id"`
	Host     string `json:"host"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Port     int    `json:"port"`
	Status   string `json:"status"`
}

func (app *Application) installTools(nodeId uuid.UUID, token string, sshConn *remote.SshConn) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	err := sshConn.InstallTools(token)
	if err != nil {
		log.Printf("Failed to install tools for node %s: %v", nodeId, err)
		updateErr := app.Repository.UpdateNodeStatus(ctx, repository.UpdateNodeStatusParams{
			Status: "error",
			ID:     nodeId,
		})
		if updateErr != nil {
			return updateErr
		}
		return err
	}
	log.Printf("Successfully installed tools for node %s", nodeId)
	err = app.Repository.UpdateNodeStatus(ctx, repository.UpdateNodeStatusParams{
		Status: "healthy",
		ID:     nodeId,
	})
	if err != nil {
		return err
	}
	return nil
}

func (app *Application) createNodeHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Printf("Error parsing multipart form: %v", err)
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("pem")
	if err != nil {
		log.Printf("Error retrieving pem file from form: %v", err)
		http.Error(w, "PEM file not present or badly formatted", http.StatusBadRequest)
		return
	}
	defer file.Close()

	name := r.FormValue("name")
	host := r.FormValue("host")
	user := r.FormValue("user")
	portStr := r.FormValue("port")

	port, err := strconv.ParseInt(portStr, 10, 64)
	if err != nil {
		log.Printf("Port must be a number: %v", err)
		http.Error(w, "Port must be a number", http.StatusInternalServerError)
		return
	}

	tempFile, err := os.CreateTemp("", "key-*.pem")
	if err != nil {
		log.Printf("Error creating temporary PEM file: %v", err)
		http.Error(w, "Error creating temporary PEM file", http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Error reading PEM file: %v", err)
		http.Error(w, "Error reading PEM file", http.StatusInternalServerError)
		return
	}

	sshConn, err := remote.New(host, user, fileBytes, true)
	if err != nil {
		log.Printf("Error establishing ssh conn: %v", err)
		http.Error(w, "Error establishing ssh conn", http.StatusInternalServerError)
		return
	}

	if err = sshConn.Ping(); err != nil {
		log.Printf("Error connecting to server: %v", err)
		http.Error(w, "Error connecting to server", http.StatusInternalServerError)
		return
	}

	pemFileEncoded := base64.StdEncoding.EncodeToString(fileBytes)

	token := uuid.New().String()

	n := repository.CreateNodeParams{
		Host:       host,
		Name:       name,
		PemFile:    pemFileEncoded,
		Username:   user,
		Port:       int32(port),
		AgentToken: token,
	}

	id, err := app.Repository.CreateNode(r.Context(), n)
	if err != nil {
		log.Printf("Error adding node to database: %v", err)
		http.Error(w, "Error creating new node", http.StatusBadRequest)
		return
	}

	err = app.installTools(id, token, sshConn)

	if err != nil {
		app.badRequestResponse(w, r, errors.New("error instaling tools in node"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Node added successfully"))

}

func (app *Application) installToolsForNode(w http.ResponseWriter, r *http.Request) {

	id, ok := app.parseUUIDParam(w, r, "nodeId")

	if !ok {
		app.badRequestResponse(w, r, errors.New("invalid id"))
		return
	}

	node, err := app.Repository.GetNodeById(r.Context(), id)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	pemBytes, err := base64.StdEncoding.DecodeString(node.PemFile)
	if err != nil {
		app.internalServerError(w, r, fmt.Errorf("failed to decode PEM file: %w", err))
		return
	}

	sshConn, err := remote.New(node.Host, node.Username, pemBytes, true)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	go app.installTools(node.ID, node.AgentToken, sshConn)

}

func (app *Application) getNodesHandler(w http.ResponseWriter, r *http.Request) {
	nodes, err := app.Repository.GetAllNodes(r.Context())

	var nodesList []Node
	for _, node := range nodes {
		nodesList = append(nodesList, Node{
			Id:       node.ID.String(),
			Host:     node.Host,
			Name:     node.Name,
			Username: node.Username,
			Status:   node.Status,
		})
	}

	if err != nil {
		log.Printf("Error fetching nodes: %v", err)
		http.Error(w, "Error fetching nodes", http.StatusBadRequest)
		return
	}

	app.jsonResponse(w, http.StatusOK, nodesList)
}
