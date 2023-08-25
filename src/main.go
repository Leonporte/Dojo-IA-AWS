package main

import (
	"fmt"
	"net/http"
	// "os"
	// "github.com/aws/aws-sdk-go/aws/session"
	// "github.com/aws/aws-sdk-go/service/textract"
)

func main() {
	http.HandleFunc("/image", imageHandler)
	fmt.Println("Servidor rodando em http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Teste imagem!")
	removeBackgroundHandler(w, r)
}

func removeBackgroundHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Teste imagem!")
}
