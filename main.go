package main

import (
	"fmt"
	"net/http"
	"os"
	"io"
	"github.com/google/uuid"
	"strings"
	"regexp"
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

func sanitizeFilename(filename string) string {
    filename = strings.ReplaceAll(filename, " ", "_")
    reg := regexp.MustCompile(`[^a-zA-Z0-9._-]`)
    filename = reg.ReplaceAllString(filename, "")
    
    return filename
}

func uploadHandler(w http.ResponseWriter, r *http.Request){
	
	// handle file size limit to 100MB
	r.Body = http.MaxBytesReader(w, r.Body, 100<<20)

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}	
	defer file.Close()

	id := uuid.New().String()
	safeFilename := id + "_" + sanitizeFilename(header.Filename)

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