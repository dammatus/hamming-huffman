package main

import (
	"fmt"
	"net/http"

	"github.com/dammatus/hamming-huffman/hamming"
)

func main() {
	// Handlers
	http.HandleFunc("hamming/files", hamming.ControlarArchivo)
	http.HandleFunc("hamming/", hamming.ControlarHamming)
	http.HandleFunc("hamming/resultados", hamming.MostrarResultados)
	http.HandleFunc("/", principal)

	// Iniciar el servidor
	fmt.Println("Servidor escuchando en http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func principal(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Si es GET, mostramos el formulario
		http.ServeFile(w, r, "/index.html")
	} else if r.Method == "POST" {
		// Si es POST, enviamos el archivo a la función controlarArchivo

	}
}
