package main

import (
	"log"
	"net/http"

	"encoding/json"

	"github.com/gorilla/mux"
)

type Document struct {
	Id   string
	Name string
	Size int
}

func getDocuments(w http.ResponseWriter, r *http.Request) {
	var docs []Document
	docs = append(docs,
		Document{Id: "doc-1", Name: "Report.docx", Size: 1500})
	docs = append(docs,
		Document{Id: "doc-2", Name: "Sheet.xlsx", Size: 5000})
	docs = append(docs,
		Document{Id: "doc-3", Name: "Container.tar", Size: 50000})
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(docs)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/documents", getDocuments).Methods("GET")
	log.Fatal(http.ListenAndServe(":9000", router))
}
