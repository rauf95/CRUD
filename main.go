package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strconv"
)

var db *sql.DB

type UpdateReq struct {
	Type   string `json:"type"`
	Region string `json:"region"`
}

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

	connStr := "user=postgres password=postgres sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Start server 2000\n")
	log.Fatal(http.ListenAndServe(":2000", router))
}


func getForest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params, _ := strconv.Atoi(mux.Vars(r)["id"])

	f := Forest{}
	err := db.QueryRowContext(r.Context(), `SELECT * FROM forests WHERE id = $1`, params).Scan(&f.ID, &f.Type, &f.Region)
	if err != nil {
		panic(err)
	}
	json.NewEncoder(w).Encode(f)
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

	err := db.QueryRowContext(r.Context(), "insert into forests (region, kind) values ($1,$2) returning id", newForest.Region, newForest.Type).
		Scan(&newForest.ID)
	if err != nil {
		log.Println("ERROR", "msg", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newForest)
}

func updateForest(w http.ResponseWriter, r *http.Request) {

	idParams := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParams)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("ID couldn't be converted to integer"))
		return
	}

	updateReq := UpdateReq{}
	err = json.NewDecoder(r.Body).Decode(&updateReq)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("couldn't unmarshal request"))
		return
	}

	for i := range forests {

		if id == forests[i].ID {
			forests[i].Type = updateReq.Type
			forests[i].Region = updateReq.Region
		}
	}
	_, err = db.ExecContext(r.Context(), "update forests set region = $1, kind = $2 where id = $3   ", updateReq.Region, updateReq.Type, id)
	if err != nil {
		log.Println("ERROR", "msg", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func deleteForest(w http.ResponseWriter, r *http.Request) {
	params, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("couldn't parse id"))
		return
	}

	_, err = db.ExecContext(r.Context(), "delete from forests where id = $1", params)
	if err != nil {
		log.Println("ERROR", "msg", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
