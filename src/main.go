package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/image", imageHandler)
	fmt.Println("Servidor rodando em http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Teste imagem!")
}
