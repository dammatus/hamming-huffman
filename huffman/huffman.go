package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dammatus/hamming-huffman/huffman/modulesHuffman"
)

func main() {

	// // Handlers
	// http.HandleFunc("/files", ControlarArchivo)
	// http.HandleFunc("/", ControlarHuffman)
	// //http.HandleFunc("/resultados", MostrarResultados)

	// // Iniciar el servidor
	// fmt.Println("Servidor escuchando en http://localhost:8080")
	// http.ListenAndServe(":8080", nil)

	text := "fffffffffffffffffffffffffffffffffffffffffffffeeeeeeeeeeeeeeeedddddddddddddccccccccccccbbbbbbbbbaaaaa"
	freqs := make(map[rune]int)
	for _, ch := range text {
		freqs[ch]++
	}
	raiz := modulesHuffman.ConstruirArbol(freqs)
	fmt.Println("Codigos Huffman:")
	modulesHuffman.PrintCodes(raiz, []byte{})

	fmt.Println(freqs)
	fmt.Printf("Tamaño: %d", len(text))
}

func ControlarHuffman(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Si es GET, mostramos el formulario
		http.ServeFile(w, r, "/huffman.html")
	} else if r.Method == "POST" {
		// Si es POST, enviamos el archivo a la función controlarArchivo
		ControlarArchivo(w, r)
	}
}
func ControlarArchivo(w http.ResponseWriter, r *http.Request) {

	// Parsea la petición y extrae el archivo subido
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Error al procesar el archivo subido: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Extrae el archivo subido
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error al procesar el archivo subido: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Crea la carpeta "files" si no existe
	os.Mkdir("files", os.ModePerm)

	// Crea el archivo en el servidor
	//Aqui deberiamos cambiar el tipo de archivo dependiendo el tipo de codificacion
	f, err := os.Create(filepath.Join("files", "archivo.txt"))
	if err != nil {
		http.Error(w, "No se pudo crear el archivo en el servidor", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// Copia el contenido del archivo subido al archivo en el servidor
	if _, err := io.Copy(f, file); err != nil {
		http.Error(w, "No se pudo guardar el archivo en el servidor", http.StatusInternalServerError)
		return
	}

	// Aplicar Huffman

	// Leer el contenido del archivo
	contenido, err := ioutil.ReadFile(filepath.Join("files", handler.Filename)) //handler.Filename tendrá "archivo.txt"
	if err != nil {
		http.Error(w, "No se pudo leer el archivo subido", http.StatusInternalServerError)
		return
	}

	// Convertir el contenido a texto y aplicar Huffman
	text := string(contenido)
	freqs := make(map[rune]int)
	for _, ch := range text {
		freqs[ch]++
	}
	raiz := modulesHuffman.ConstruirArbol(freqs)
	fmt.Println("Codigos Huffman:")
	modulesHuffman.PrintCodes(raiz, []byte{})

	fmt.Println(freqs)
	fmt.Printf("Tamaño: %d", len(text))

	// // Convertir el resultado a texto y escribirlo en un archivo
	// if err := ioutil.WriteFile(filepath.Join("files", "comprimido.txt"), []byte(text), 0644); err != nil {
	// 	http.Error(w, "No se pudo guardar el archivo comprimido", http.StatusInternalServerError)
	// 	return
	// }

	// // Descomprimir archivo y copiarlo en archivo nuevo
	// descomprimido := make([]byte, 0)
	// if err := ioutil.WriteFile(filepath.Join("files", "descomprimido.txt"), []byte(descomprimido[:len(contenido)]), 0644); err != nil {
	// 	http.Error(w, "No se pudo guardar el archivo descomprimido", http.StatusInternalServerError)
	// 	return
	// }

	// Servir el archivo HTML con los resultados
	//http.HandlerFunc(MostrarResultados).ServeHTTP(w, r)

}
