package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/sharithg/siphon/models"
	"github.com/sharithg/siphon/ssh"
)

type Node struct {
	Id   string `json:"id"`
	Host string `json:"host"`
	Name string `json:"name"`
	User string `json:"user"`
}

func (env *Env) AddNode(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

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

	conn, err := ssh.New(host, user, fileBytes, true)
	if err != nil {
		log.Printf("Error establishing ssh conn: %v", err)
		http.Error(w, "Error establishing ssh conn", http.StatusInternalServerError)
		return
	}

	if err = conn.Ping(); err != nil {
		log.Printf("Error connecting to server: %v", err)
		http.Error(w, "Error connecting to server", http.StatusInternalServerError)
		return
	}

	_, err = tempFile.Write(fileBytes)
	if err != nil {
		log.Printf("Error writing to temporary PEM file: %v", err)
		http.Error(w, "Error writing to temporary PEM file", http.StatusInternalServerError)
		return
	}

	objectName := fmt.Sprintf("assets/%s/key.pem", name)
	info, err := env.Storage.Upload(ctx, "node-pem-files", objectName, tempFile.Name(), "application/x-x509-ca-ce")
	if err != nil {
		log.Printf("Error uploading PEM file to storage: %v", err)
		http.Error(w, "Error uploading PEM file to storage", http.StatusInternalServerError)
		return
	}

	n := models.Node{
		Host:    host,
		Name:    name,
		PemFile: info.Key,
		User:    user,
	}

	if err = env.Nodes.AddNode(n); err != nil {
		log.Printf("Error adding node to database: %v", err)
		http.Error(w, "Error creating new node", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Node added successfully"))
}

func (env *Env) GetNodes(w http.ResponseWriter, r *http.Request) {
	nodes, err := env.Nodes.All()

	var nodesList []Node
	for _, node := range nodes {
		nodesList = append(nodesList, Node{
			Id:   node.Id,
			Host: node.Host,
			Name: node.Name,
			User: node.User,
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
