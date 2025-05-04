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
	dataCenterIDStr := os.Getenv("DATA_CENTER_ID")
	machineIDStr := os.Getenv("MACHINE_ID")

	dataCenterID, err := strconv.Atoi(dataCenterIDStr)
	if err != nil {
		log.Fatalf("Invalid DATA_CENTER_ID: %v", err)
	}

	machineID, err := strconv.Atoi(machineIDStr)
	if err != nil {
		log.Fatalf("Invalid MACHINE_ID: %v", err)
	}

	s := snowid.NewSnowIDGenerator(dataCenterID, machineID, snowid.DefaultEpoch)
	if err != nil {
		log.Fatalf("Failed to Initiate SnowID generator: %v", err)
	}

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello Unique ID server"))
	})

	http.HandleFunc("POST /ids", func(w http.ResponseWriter, r *http.Request) {
		id, err := s.GenerateID()
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
		log.Println("ID", id.String())
		w.Header().Set("Content-Type", "application/json")
		res := map[string]any{
			"status": 200,
			"id":     id.String(),
		}
		json.NewEncoder(w).Encode(res)
	})

	http.HandleFunc("GET /ids/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		idStr := r.PathValue("id")

		// Parse binary id
		id, err := snowid.ParseID(idStr, snowid.DefaultEpoch)
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
			"status":    200,
			"data":      id,
			"datetime":  id.Datetime().UTC().String(),
			"id_binary": id.StringBinary(),
		}
		json.NewEncoder(w).Encode(res)
	})

	log.Println("Listening at 8000")

	// Reset records every 10 seconds
	// g.AutoResetRecords(10 * time.Second)
	http.ListenAndServe(":8000", nil)
}
