package main

import (
	"log"
	"net/http"

	"github.com/ddahon/jobboard/cmd/server/views"
	"github.com/ddahon/jobboard/internal/models"
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := views.Index(jobs).Render(r.Context(), w); err != nil {
			log.Printf("Failed to respond to request: %v", err)
		}

	})
	log.Println("Starting to listen on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
