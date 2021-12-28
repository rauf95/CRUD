package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type Forest struct {
	ID     int    `json:"id"`
	Type   string `json:"type"`
	Region string `json:"region"`
}

var forests []Forest = []Forest{
	{
		ID:     1,
		Type:   "rain,tropic",
		Region: "Africa,Asia",
	},
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/forest/{id}", getForest).Methods("GET")
	router.HandleFunc("/forest", getAllForest).Methods("GET")
	router.HandleFunc("/forest", createForest).Methods("POST")
	router.HandleFunc("/forest/{id}", updateForest).Methods("PUT")
	router.HandleFunc("/forest/{id}", deleteForest).Methods("DELETE")

	fmt.Printf("Start server 2000\n")
	log.Fatal(http.ListenAndServe(":2000", router))
}

func getForest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params, _ := strconv.Atoi(mux.Vars(r)["id"])

	for _, item := range forests {
		if item.ID == params {
			err := json.NewEncoder(w).Encode(item)
			if err != nil {
				return
			}
		}
	}
}

func getAllForest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(forests)

	if err != nil {
		return
	}
}

func createForest(w http.ResponseWriter, r *http.Request) {
	var newForest Forest
	json.NewDecoder(r.Body).Decode(&newForest)
	forests = append(forests, newForest)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(forests)

}

func updateForest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params, _ := strconv.Atoi(mux.Vars(r)["id"])

	for index, item := range forests {
		if item.ID == params {
			forests = append(forests[:index], forests[index+1])
			var forest Forest
			_ = json.NewDecoder(r.Body).Decode(&forest)
			forest.ID = params
			forests = append(forests, forest)
			err := json.NewEncoder(w).Encode(forest)
			if err != nil {
				return
			}
		}
		return
	}

}

func deleteForest(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params, _ := strconv.Atoi(mux.Vars(r)["id"])
	for i, item := range forests {
		if item.ID == params {
			forests = append(forests[:i], forests[i+i])
			fmt.Printf("forest with ID %v is deleted", params)

		}
	}
}
