package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ddahon/jobboard/cmd/server/views"
	"github.com/ddahon/jobboard/internal/pkg/models"
	"github.com/spf13/viper"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please specify the config file path in the arguments")
	}
	getConfig(os.Args[1])
	var err error
	dbPath := viper.GetString("dbPath")
	port := viper.GetString("port")
	sslEnabled := viper.GetBool("sslEnabled")
	err = models.InitDB(dbPath)
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
	if sslEnabled {
		certFile := viper.GetString("sslCertFile")
		keyFile := viper.GetString("sslKeyFile")
		log.Fatal(http.ListenAndServeTLS(":443", certFile, keyFile, nil))
	} else {
		log.Fatal(http.ListenAndServe(":"+port, nil))
	}
}

func getConfig(path string) {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(fmt.Errorf("config file not found in %v: %w", path, err))
		} else {
			panic(fmt.Errorf("error while reading config file: %w", err))
		}
	}
}
