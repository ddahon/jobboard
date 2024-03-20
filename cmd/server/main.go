package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ddahon/jobboard/internal/pkg/models"
	_ "github.com/lib/pq"
)

func main() {
	var err error
	connStr := "postgresql://postgres:password@localhost:5432/jobs?sslmode=disable"
	err = models.InitDB(connStr)
	if err != nil {
		log.Fatalln(err)
	}

	jobs, err := models.GetAllJobs()
	if err != nil {
		log.Fatalf("Failed to retrieve jobs from DB: %v", err)
	}

	vueAppPath := os.Getenv("VUE_APP_PATH")
	if vueAppPath == "" {
		panic("VUE_APP_PATH environment variable is not set")
	}
	log.Printf("Found %v jobs", len(jobs))
	server := http.NewServeMux()
	server.Handle("/", http.FileServer(http.Dir(vueAppPath)))

	log.Println("Starting to listen on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", server))
}
