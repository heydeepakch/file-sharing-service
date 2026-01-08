package main

import (
	"fmt"
	"net/http"
)

func main(){
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/upload", uploadHandler)
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
	fmt.Println ("Uploaded File: ", header.Filename)
	fmt.Fprintln(w, "Upload Successful")
}