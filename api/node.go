package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sharithg/siphon/internal/remote"
	"github.com/sharithg/siphon/internal/storage"
)

type Node struct {
	Id     string `json:"id"`
	Host   string `json:"host"`
	Name   string `json:"name"`
	User   string `json:"user"`
	Port   int    `json:"port"`
	Status string `json:"status"`
}

func (app *Application) installTools(ctx context.Context, nodeId, token string, sshConn *remote.SshConn) error {
	err := sshConn.InstallTools(token)
	if err != nil {
		log.Printf("Failed to install tools for node %s: %v", nodeId, err)
		app.Store.Nodes.UpdateStatus(ctx, nodeId, "error")
		return err
	}
	log.Printf("Successfully installed tools for node %s", nodeId)
	app.Store.Nodes.UpdateStatus(ctx, nodeId, "healthy")
	return nil
}

func (app *Application) createNodeHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("File Upload Endpoint Hit")

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

	n := storage.Node{
		Host:       host,
		Name:       name,
		PemFile:    pemFileEncoded,
		User:       user,
		Port:       port,
		AgentToken: token,
	}

	id, err := app.Store.Nodes.Create(r.Context(), n)
	if err != nil {
		log.Printf("Error adding node to database: %v", err)
		http.Error(w, "Error creating new node", http.StatusBadRequest)
		return
	}

	err = app.installTools(r.Context(), id, token, sshConn)

	if err != nil {
		app.badRequestResponse(w, r, errors.New("error instaling tools in node"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Node added successfully"))

}

func (app *Application) installToolsForNode(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "nodeId")

	node, err := app.Store.Nodes.GetById(r.Context(), idParam)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	pemBytes, err := base64.StdEncoding.DecodeString(node.PemFile)
	if err != nil {
		app.internalServerError(w, r, fmt.Errorf("failed to decode PEM file: %w", err))
		return
	}

	sshConn, err := remote.New(node.Host, node.User, pemBytes, true)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if node == nil {
		app.notFoundResponse(w, r, errors.New("node not found"))
		return
	}

	go app.installTools(r.Context(), node.Id, node.AgentToken, sshConn)

}

func (app *Application) getNodesHandler(w http.ResponseWriter, r *http.Request) {
	nodes, err := app.Store.Nodes.All(r.Context())

	var nodesList []Node
	for _, node := range nodes {
		nodesList = append(nodesList, Node{
			Id:     node.Id,
			Host:   node.Host,
			Name:   node.Name,
			User:   node.User,
			Status: node.Status,
		})
	}

	if err != nil {
		log.Printf("Error fetching nodes: %v", err)
		http.Error(w, "Error fetching nodes", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(nodesList)
}
