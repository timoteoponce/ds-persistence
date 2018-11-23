package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"encoding/json"

	"github.com/gorilla/mux"
)

const StoragePath = "./files/"

type Document struct {
	Id   string
	Name string
	Size int64
	Path string `json:"-"`
}

func getDocuments(w http.ResponseWriter, r *http.Request) {
	docs := findDocuments("")
	w.Header().Set("Content-Type", "application/json")
	log.Printf("Data is %v\n", docs)
	json.NewEncoder(w).Encode(docs)
}

func getDocumentsById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	docs := findDocuments(vars["id"])
	if len(docs) == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Not found")
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(docs[0])
	}
}

func deleteDocumentsById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	docs := findDocuments(vars["id"])
	if len(docs) == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Not found")
	} else {
		handleError(os.Remove(docs[0].Path))
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "DELETED")
	}
}

func findDocuments(id string) []Document {
	if _, err := os.Stat(StoragePath); os.IsNotExist(err) {
		os.Mkdir(StoragePath, os.ModePerm)
	}
	var files []string
	err := filepath.Walk(StoragePath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	handleError(err)
	docs := make([]Document, 0)
	for _, f := range files {
		doc := fileToDocument(f)
		if id == "" || id == doc.Id {
			docs = append(docs, doc)
		}
	}
	return docs
}

func fileToDocument(f string) Document {
	file, err := os.Open(f)
	handleError(err)
	defer file.Close()
	info, err := file.Stat()
	handleError(err)
	return Document{Id: checksum(file), Name: file.Name(), Size: info.Size(), Path: f}
}

func checksum(f *os.File) string {
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/documents", getDocuments).Methods("GET")
	router.HandleFunc("/documents/{id}", getDocumentsById).Methods("GET")
	router.HandleFunc("/documents/{id}", deleteDocumentsById).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":9000", router))
}
