package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/ddahon/jobboard/cmd/server/views"
	"github.com/ddahon/jobboard/internal/models"
	_ "github.com/lib/pq"
)

func main() {
	connStr := "postgresql://postgres:password@localhost:5432/jobs?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Could not open DB: %v: ", err)
	}
	res, err := db.Query("SELECT * FROM jobs")
	if err != nil {
		log.Fatalf("Failed to retrieve jobs from DB: %v", err)
	}
	var jobs []models.Job

	for res.Next() {
		var job models.Job
		res.Scan(&job.Id, &job.Description, &job.Title, &job.Link)
		jobs = append(jobs, job)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := views.Index(jobs).Render(r.Context(), w); err != nil {
			log.Printf("Failed to respond to request: %v", err)
		}

	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
