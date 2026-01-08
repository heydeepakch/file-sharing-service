package main

import (
	"fmt"
	"net/http"
	"os"
	"io"
	"github.com/google/uuid"
	"strings"
	"regexp"
	"path/filepath"
	"time"
)

var fileDB = make(map[string]FileMeta)

func main(){

	loadDB()

	err := os.MkdirAll("uploads", 0755)
    if err != nil {
        fmt.Println("Error creating uploads directory:", err)
        return
    }
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/download/", downloadHandler)
	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("uploads"))))
	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/download/")
	meta, ok := fileDB[id]
	if !ok {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, "uploads/" + meta.Path)
	fmt.Fprintln(w, "File downloaded successfully.")
}

func sanitizeFilename(filename string) string {
    filename = strings.ReplaceAll(filename, " ", "_")
    reg := regexp.MustCompile(`[^a-zA-Z0-9._-]`)
    filename = reg.ReplaceAllString(filename, "")
    
    return filename
}

func isAllowedFileType(filename string) bool {
    allowedExtensions := []string{".pdf", ".jpg", ".jpeg", ".png", ".gif", ".zip"}
    
    ext := strings.ToLower(filepath.Ext(filename))
    for _, allowed := range allowedExtensions {
        if ext == allowed {
            return true
        }
    }
    return false
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

	if !isAllowedFileType(header.Filename) {
		http.Error(w, "Unsupported file type", http.StatusUnsupportedMediaType)
		return
	}

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

	publicID := uuid.New().String()
	fileDB[publicID] = FileMeta{
		ID: publicID,
		Name: header.Filename,
		Path: safeFilename,
		Size: header.Size,
		UploadedAt: time.Now(),
	}

	saveDB()

	link := "http://localhost:8080/download/" + publicID
	fmt.Fprintln(w, "File uploaded successfully.\n\n Download Link:\n", link)

	fmt.Println ("Uploaded File: ", header.Filename)
}