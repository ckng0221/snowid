package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/ckng0221/snowid"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()

	if err != nil {
		fmt.Println("Error loading .env file")
	}
}

func main() {

	dataCenterIdStr := os.Getenv("DATA_CENTER_ID")
	machineIdStr := os.Getenv("MACHINE_ID")

	dataCenterId, err := strconv.Atoi(dataCenterIdStr)
	if err != nil {
		log.Fatalf("Invalid DATA_CENTER_ID: %v", err)
	}

	machineId, err := strconv.Atoi(machineIdStr)
	if err != nil {
		log.Fatalf("Invalid MACHINE_ID: %v", err)
	}

	g, err := snowid.NewSnowIdGenerator(dataCenterId, machineId, snowid.DefaultEpoch)
	if err != nil {
		log.Fatalf("Failed to create ID generator: %v", err)
	}

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello ID generator"))
	})

	http.HandleFunc("POST /snowids", func(w http.ResponseWriter, r *http.Request) {
		id := g.GenerateId()
		log.Println("id", id.String())
		w.Header().Set("Content-Type", "application/json")
		res := map[string]any{
			"status": 200,
			"id":     id.String(),
		}
		json.NewEncoder(w).Encode(res)
	})

	http.HandleFunc("GET /snowids/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		id := r.PathValue("id")

		// Parse binary id
		idObj, err := snowid.ParseId(id, snowid.DefaultEpoch)
		if err != nil {
			log.Print(err.Error())
			errMsg := "Internal server error"
			res := map[string]any{
				"status": 500,
				"error":  errMsg,
			}
			json.NewEncoder(w).Encode(res)
			return
		}

		res := map[string]any{
			"status":   200,
			"data":     idObj,
			"datetime": idObj.Datetime().UTC().String(),
		}
		json.NewEncoder(w).Encode(res)
	})

	log.Println("Listening at 8000")

	// Reset records every 10 seconds
	// g.AutoResetRecords(10 * time.Second)
	http.ListenAndServe(":8000", nil)
}
