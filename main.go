package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var fileDB = make(map[string]FileMeta)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	loadDB()

	// err := os.MkdirAll("uploads", 0755)
	// if err != nil {
	//     fmt.Println("Error creating uploads directory:", err)
	//     return
	// }
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/download/", downloadHandler)
	// http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("uploads"))))
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server is running on port", port)
	http.ListenAndServe(":"+port, nil)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/download/")

	mu.Lock()
	meta, ok := fileDB[id]
	mu.Unlock()

	if !ok {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Create presigned URL for R2
	storage := newStorage()
	presign := s3.NewPresignClient(storage)
	bucketName := os.Getenv("R2_BUCKET_NAME")

	res, err := presign.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(meta.Path),
	}, s3.WithPresignExpires(15*time.Minute))

	if err != nil {
		http.Error(w, "Unable to generate download link", http.StatusInternalServerError)
		return
	}

	// Redirect to the presigned URL
	http.Redirect(w, r, res.URL, http.StatusFound)
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

func uploadHandler(w http.ResponseWriter, r *http.Request) {

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

	// Upload to R2 instead of disk
	storage := newStorage()
	bucketName := os.Getenv("R2_BUCKET_NAME")

	_, err = storage.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(safeFilename),
		Body:   file,
	})
	if err != nil {
		http.Error(w, "Unable to upload file to R2: "+err.Error(), http.StatusInternalServerError)
		return
	}

	publicID := uuid.New().String()
	fileDB[publicID] = FileMeta{
		ID:         publicID,
		Name:       header.Filename,
		Path:       safeFilename,
		Size:       header.Size,
		UploadedAt: time.Now(),
	}

	saveDB()

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	link := baseURL + "/download/" + publicID
	fmt.Fprintln(w, "File uploaded successfully.\n\n Download Link:\n", link)

	fmt.Println("Uploaded File: ", header.Filename)
}
