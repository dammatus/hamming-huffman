package main

import (
	"fmt"
	"hamming-huffman/codificar/hamming"
	"hamming-huffman/comprimir/huffman"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type Hamm struct {
	Contenido    string
	Codificado   string
	Decodificado string
}

type Huff struct {
	Contenido    string
	Comprimido   string
	Decomprimido string
}

const (
	bitsParity32    = 5
	bitsParity2048  = 11
	bitsParity65536 = 16
	bitsInfo32      = 26
	bitsInfo2048    = 2036
	bitsInfo65536   = 65519
)

var blockSize int

func main() {
	// Configurar los manejadores de las rutas
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/codificar/", codificarHandler)
	http.HandleFunc("/codificar/files", archivosCodHandler)
	http.HandleFunc("/codificar/resultados", mostrarResultadosCod)
	http.HandleFunc("/comprimir/", comprimirHandler)
	http.HandleFunc("/comprimir/files", archivosCompHandler)
	http.HandleFunc("/comprimir/resultados", mostrarResultadosComp)
	http.HandleFunc("/hamming-huffman/", hammingHuffmanHandler)

	// Iniciar el servidor
	fmt.Println("Servidor escuchando en http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.ServeFile(w, r, "index.html")
	}
}

func codificarHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Si es GET, mostramos el formulario
		http.ServeFile(w, r, "codificar/hamming.html")
	} else if r.Method == "POST" {
		// Si es POST, enviamos el archivo a la función archivosHandler
		archivosCodHandler(w, r)
	}
}

func comprimirHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Si es GET, mostramos el formulario
		http.ServeFile(w, r, "comprimir/huffman.html")
	} else if r.Method == "POST" {
		// Si es POST, enviamos el archivo a la función archivosHandler
		archivosCompHandler(w, r)
	}
}

func hammingHuffmanHandler(w http.ResponseWriter, r *http.Request) {
	// Agregar funcionalidad después
}

/*
Funciones para la codificacion de un archivo
*/
func archivosCodHandler(w http.ResponseWriter, r *http.Request) {

	// Parsea la petición y extrae el archivo subido
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Error al procesar el archivo subido: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Extrae el archivo subido
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error al procesar el archivo subido: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Extrae el tipo de codificación
	blockSizeStr := r.FormValue("tipoCod")
	blockSize, err = strconv.Atoi(blockSizeStr)
	if err != nil {
		http.Error(w, "Error al convertir el tamaño de bloque: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Extrae el error
	errorStr := r.FormValue("error")
	hasError := false
	if errorStr == "Si" {
		hasError = true
	}

	// Crea la carpeta "codificar/files" si no existe
	err = os.MkdirAll("codificar/files", os.ModePerm)
	if err != nil {
		http.Error(w, "No se pudo crear la carpeta en el servidor", http.StatusInternalServerError)
		return
	}

	// Crea el archivo en el servidor
	f, err := os.Create(filepath.Join("codificar/files", "archivo.txt"))
	if err != nil {
		http.Error(w, "No se pudo crear el archivo en el servidor", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// Copia el contenido del archivo subido al archivo en el servidor
	_, err = io.Copy(f, file)
	if err != nil {
		http.Error(w, "No se pudo guardar el archivo en el servidor", http.StatusInternalServerError)
		return
	}

	// Aplicar Hamming
	var parityBits, infoBits int

	switch blockSize {
	case 32:
		parityBits = bitsParity32
		infoBits = bitsInfo32
	case 2048:
		parityBits = bitsParity2048
		infoBits = bitsInfo2048
	case 65536:
		parityBits = bitsParity65536
		infoBits = bitsInfo65536
	default:
		http.Error(w, "El tamaño de bloque es inválido", http.StatusBadRequest)
		return
	}

	// Leer el contenido del archivo
	contenido, err := ioutil.ReadFile(filepath.Join("codificar/files", "archivo.txt"))
	if err != nil {
		http.Error(w, "No se pudo leer el archivo subido", http.StatusInternalServerError)
		return
	}

	// Convertir el contenido a bits y aplicar Hamming
	bits := hamming.ByteToBits(contenido, blockSize)
	encode := hamming.AplicandoHamming(bits, blockSize, parityBits, infoBits, hasError)

	// Convertir el resultado a texto y escribirlo en un archivo
	ascii := hamming.BinToASCII(encode)

	codificadoFileName := ""
	switch blockSize {
	case 32:
		codificadoFileName = "codificado.HA1"
	case 2048:
		codificadoFileName = "codificado.HA2"
	case 65536:
		codificadoFileName = "codificado.HA3"
	}
	//Este es el que se mostrara en la pagina
	if err := ioutil.WriteFile(filepath.Join("codificar/files", "codificado.txt"), []byte(ascii), 0644); err != nil {
		http.Error(w, "No se pudo guardar el archivo codificado", http.StatusInternalServerError)
		return
	}
	//Este es el que cumple con la consigna, se guarda en la carpeta resultados
	if err := ioutil.WriteFile(filepath.Join("codificar/resultados", codificadoFileName), []byte(ascii), 0644); err != nil {
		http.Error(w, "No se pudo guardar el archivo codificado", http.StatusInternalServerError)
		return
	}

	// Decodificar el contenido y escribirlo en un archivo (Sin corregir)
	decode := hamming.DecodeHamming(encode, blockSize, infoBits, false, parityBits)
	asciiDeco := hamming.BitsToByte(decode)
	decoded := string(asciiDeco)

	decodificadoFileName := ""
	switch blockSize {
	case 32:
		decodificadoFileName = "decodificado.DE1"
	case 2048:
		decodificadoFileName = "decodificado.DE2"
	case 65536:
		decodificadoFileName = "decodificado.DE3"
	}

	//Este es el que se mostrara en la pagina
	if err := ioutil.WriteFile(filepath.Join("codificar/files", "decodificado.txt"), []byte(decoded[:len(contenido)]), 0644); err != nil {
		http.Error(w, "No se pudo guardar el archivo decodificado", http.StatusInternalServerError)
		return
	}
	//Este es el que cumple con la consigna, se guarda en la carpeta resultados
	if err := ioutil.WriteFile(filepath.Join("codificar/resultados", decodificadoFileName), []byte(decoded[:len(contenido)]), 0644); err != nil {
		http.Error(w, "No se pudo guardar el archivo decodificado", http.StatusInternalServerError)
		return
	}

	if hasError {
		// Decodificar el contenido y escribirlo en un archivo (Corregido)
		decode = hamming.DecodeHamming(encode, blockSize, infoBits, hasError, parityBits)
		asciiDeco = hamming.BitsToByte(decode)
		decoded = string(asciiDeco)

		corregidoFileName := ""
		switch blockSize {
		case 32:
			corregidoFileName = "decodificado.DC1"
		case 2048:
			corregidoFileName = "decodificado.DC2"
		case 65536:
			corregidoFileName = "decodificado.DC3"
		}
		//Este es el que se mostrara en la pagina
		if err := ioutil.WriteFile(filepath.Join("codificar/files", "decodificado.txt"), []byte(decoded[:len(contenido)]), 0644); err != nil {
			http.Error(w, "No se pudo guardar el archivo decodificado", http.StatusInternalServerError)
			return
		}
		//Este es el que cumple con la consigna, se guarda en la carpeta resultados
		if err := ioutil.WriteFile(filepath.Join("codificar/resultados", corregidoFileName), []byte(decoded[:len(contenido)]), 0644); err != nil {
			http.Error(w, "No se pudo guardar el archivo decodificado", http.StatusInternalServerError)
			return
		}

	}

	// Servir el archivo HTML con los resultados
	mostrarResultadosCod(w, r)
}

func mostrarResultadosCod(w http.ResponseWriter, _ *http.Request) {
	// Establecer el tipo de contenido para que se muestre en utf-8 (igual los acentos no los muestra bien)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Leer el archivo original
	contenido, err := ioutil.ReadFile(filepath.Join("codificar/files", "archivo.txt"))
	if err != nil {
		http.Error(w, "No se pudo leer el archivo", http.StatusInternalServerError)
		return
	}

	// Leer el archivo codificado
	codificado, err := ioutil.ReadFile(filepath.Join("codificar/files", "codificado.txt"))
	if err != nil {
		http.Error(w, "No se pudo leer el archivo codificado", http.StatusInternalServerError)
		return
	}

	// Leer el archivo decodificado
	decodificado, err := ioutil.ReadFile(filepath.Join("codificar/files", "decodificado.txt"))
	if err != nil {
		http.Error(w, "No se pudo leer el archivo decodificado", http.StatusInternalServerError)
		return
	}

	// Crear un mapa de datos para la plantilla HTML
	data := Hamm{
		Contenido:    string(contenido),
		Codificado:   string(codificado),
		Decodificado: string(decodificado),
	}

	// Leer la plantilla HTML
	tmpl, err := template.ParseFiles("codificar/resultados/resultados.html")
	if err != nil {
		http.Error(w, "No se pudo leer la plantilla HTML", http.StatusInternalServerError)
		return
	}

	// Pasar los datos a la plantilla HTML
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "No se pudo procesar la plantilla HTML", http.StatusInternalServerError)
		return
	}
}

/*
Fin funciones para la codificacion de un archivo
*/

/*
Funciones para la compactacion de un archivo
*/
func archivosCompHandler(w http.ResponseWriter, r *http.Request) {
	// Parsea la petición y extrae el archivo subido
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Error al procesar el archivo subido: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Extrae el archivo subido
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error al procesar el archivo subido: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Crea la carpeta "comprimir/files" si no existe
	err = os.MkdirAll("comprimir/files", os.ModePerm)
	if err != nil {
		http.Error(w, "No se pudo crear la carpeta en el servidor", http.StatusInternalServerError)
		return
	}

	// Crea el archivo en el servidor
	f, err := os.Create(filepath.Join("comprimir/files", "archivo.txt"))
	if err != nil {
		http.Error(w, "No se pudo crear el archivo en el servidor", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// Copia el contenido del archivo subido al archivo en el servidor
	_, err = io.Copy(f, file)
	if err != nil {
		http.Error(w, "No se pudo guardar el archivo en el servidor", http.StatusInternalServerError)
		return
	}

	// Aplicar Huffman
	// Leer el contenido del archivo
	contenido, err := ioutil.ReadFile(filepath.Join("comprimir/files", "archivo.txt"))
	if err != nil {
		http.Error(w, "No se pudo leer el archivo subido", http.StatusInternalServerError)
		return
	}
	text := string(contenido)

	fmt.Println(text)

	freqs := make(map[rune]int)
	for _, ch := range text {
		freqs[ch]++
	}
	raiz := huffman.ConstruirArbol(freqs)

	fmt.Println("Codigos Huffman:")
	huffman.PrintCodes(raiz, []byte{})

	fmt.Println(freqs)
	fmt.Printf("Tamaño: %d\n", len(text))

	//binary := huffman.Compacted(text, raiz)
	compacted := huffman.Compacted(text, raiz)

	//fmt.Println("Compactado: " + binary)

	//compacted := huffman.BinaryToBytes(binary)

	fmt.Println(compacted)

	//Este es el que se mostrara en la pagina
	/* if err := ioutil.WriteFile(filepath.Join("comprimir/files", "comprimido.txt"), compacted, 0644); err != nil {
		http.Error(w, "No se pudo guardar el archivo codificado", http.StatusInternalServerError)
		return
	} */

	//Este es el que cumple con la consigna, se guarda en la carpeta resultados
	/* if err := ioutil.WriteFile(filepath.Join("comprimir/resultados", "comprimido.huf"), compacted, 0644); err != nil {
		http.Error(w, "No se pudo guardar el archivo codificado", http.StatusInternalServerError)
		return
	} */
	err = huffman.SaveCompacted(compacted)
	if err != nil {
		fmt.Println("Error al guardar el archivo: ", err)
	} else {
		fmt.Println("Datos comprimidos exitosamente!")
	}

	mostrarResultadosComp(w, r)
}

func mostrarResultadosComp(w http.ResponseWriter, _ *http.Request) {
	// Establecer el tipo de contenido para que se muestre en utf-8 (igual los acentos no los muestra bien)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Leer el archivo original
	contenido, err := ioutil.ReadFile(filepath.Join("comprimir/files", "archivo.txt"))
	if err != nil {
		http.Error(w, "No se pudo leer el archivo", http.StatusInternalServerError)
		return
	}

	// Leer el archivo comprimido
	comprimido, err := ioutil.ReadFile(filepath.Join("comprimir/files", "comprimido.txt"))
	if err != nil {
		http.Error(w, "No se pudo leer el archivo codificado", http.StatusInternalServerError)
		return
	}

	// Leer el archivo descomprimido
	decomprimido, err := ioutil.ReadFile(filepath.Join("comprimir/files", "archivo.txt")) //Todavia no se descomprime
	if err != nil {
		http.Error(w, "No se pudo leer el archivo decodificado", http.StatusInternalServerError)
		return
	}

	// Crear un mapa de datos para la plantilla HTML
	data := Huff{
		Contenido:    string(contenido),
		Comprimido:   string(comprimido),
		Decomprimido: string(decomprimido),
	}

	// Leer la plantilla HTML
	tmpl, err := template.ParseFiles("comprimir/resultados/resultados.html")
	if err != nil {
		http.Error(w, "No se pudo leer la plantilla HTML", http.StatusInternalServerError)
		return
	}

	// Pasar los datos a la plantilla HTML
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "No se pudo procesar la plantilla HTML", http.StatusInternalServerError)
		return
	}
}