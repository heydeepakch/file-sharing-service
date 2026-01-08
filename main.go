package main

import (
	"fmt"
	"net/http"
	"os"
	"io"
	"github.com/google/uuid"
)

func main(){
	err := os.MkdirAll("uploads", 0755)
    if err != nil {
        fmt.Println("Error creating uploads directory:", err)
        return
    }
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/upload", uploadHandler)
	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("uploads"))))
	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func uploadHandler(w http.ResponseWriter, r *http.Request){
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}	
	defer file.Close()

	id := uuid.New().String()
	safeFilename := id + "_" + header.Filename

	dst, err := os.Create("uploads/" + safeFilename)
	if err != nil{
		http.Error(w, "Unable to create file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil{
		http.Error(w, "Unable to copy file", http.StatusInternalServerError)
		return
	}

	link := "http://localhost:8080/files/" + safeFilename
	fmt.Fprintln(w, "File uploaded successfully.\n\n Download Link:\n", link)

	fmt.Println ("Uploaded File: ", header.Filename)
}