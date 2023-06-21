package main

import (
	"fmt"
	codificarcomprimir "hamming-huffman/codificar-comprimir"
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

type HammHuff struct {
	Contenido string
	Resultado string
}
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
var result HammHuff

func main() {
	// Configurar los manejadores de las rutas
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/codificar/", codificarHandler)
	http.HandleFunc("/codificar/files", archivosCodHandler)
	http.HandleFunc("/codificar/resultados", mostrarResultadosCod)
	http.HandleFunc("/comprimir/", comprimirHandler)
	http.HandleFunc("/comprimir/files", archivosCompHandler)
	http.HandleFunc("/comprimir/resultados", mostrarResultadosComp)
	http.HandleFunc("/codificar-comprimir/", codificarComprimirHandler)
	http.HandleFunc("/codificar-comprimir/resultados", mostrarResultadosCodComp)

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

func codificarComprimirHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Si es GET, mostramos el formulario
		http.ServeFile(w, r, "codificar-comprimir/codificar-comprimir.html")
	} else if r.Method == "POST" {
		// Agregar funcionalidad despues
		archivosCodCompHandler(w, r)
	}
}

/*
Funciones para la codificacion y compresion de un archivo
*/
func archivosCodCompHandler(w http.ResponseWriter, r *http.Request) {
	// Parsea la petición y extrae el archivo subido
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Error al procesar el archivo subido: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Extrae el archivo subido
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error al procesar el archivo subido: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Obtiene la extensión del archivo subido
	extension := filepath.Ext(fileHeader.Filename)
	fmt.Println("Extensión del archivo:", extension)

	// Extrae el tipo de codificación
	blockSizeStr := r.FormValue("block")
	blockSize, err = strconv.Atoi(blockSizeStr)
	if err != nil {
		http.Error(w, "Error al convertir el tamaño de bloque: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Extrae el error
	errorStr := r.FormValue("Error")
	hasError := false
	var errores int
	if errorStr == "error1" || errorStr == "error2" {
		hasError = true
		if errorStr == "error1" {
			errores = 1
		} else {
			errores = 2
		}
	}
	fmt.Print(errores)

	// Crea la carpeta "codificar-comprimir/files" si no existe
	err = os.MkdirAll("codificar-comprimir/files", os.ModePerm)
	if err != nil {
		http.Error(w, "No se pudo crear la carpeta en el servidor", http.StatusInternalServerError)
		return
	}

	// Crea el archivo en el servidor
	f, err := os.Create(filepath.Join("codificar-comprimir/files", "archivo.txt"))
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
	// Leer el contenido del archivo
	contenido, err := ioutil.ReadFile(filepath.Join("codificar-comprimir/files", "archivo.txt"))
	if err != nil {
		http.Error(w, "No se pudo leer el archivo subido", http.StatusInternalServerError)
		return
	}

	// Extrae la operacion a realizar
	options := r.FormValue("options")
	switch options {
	// Codificar
	case "option1":
		/*
			Luego controlar que si las extensiones no son del tipo txt habra errores
		*/
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

		codificarcomprimir.Codificar(contenido, parityBits, infoBits, blockSize, hasError, w)

		// Leer el archivo original
		contenido, err := ioutil.ReadFile(filepath.Join("codificar-comprimir/files", "archivo.txt"))
		if err != nil {
			http.Error(w, "No se pudo leer el archivo", http.StatusInternalServerError)
			return
		}
		// Leer el archivo codificado
		codificado, err := ioutil.ReadFile(filepath.Join("codificar-comprimir/files", "codificado.txt"))
		if err != nil {
			http.Error(w, "No se pudo leer el archivo codificado", http.StatusInternalServerError)
			return
		}

		// Crear un mapa de datos para la plantilla HTML

		result := HammHuff{
			Contenido: string(contenido),
			Resultado: string(codificado),
		}
		fmt.Print(result)
		mostrarResultadosCodComp(w, r)

	// Decodificar
	case "option2":
		/*
			Luego controlar que si las extensiones no son del tipo txt habra errores
		*/
		// Aplicar Hamming
		var parityBits, infoBits int

		switch extension {
		case ".DE1":
			parityBits = bitsParity32
			infoBits = bitsInfo32
		case ".DC1":
			parityBits = bitsParity32
			infoBits = bitsInfo32
		case ".DE2":
			parityBits = bitsParity2048
			infoBits = bitsInfo2048
		case ".DC2":
			parityBits = bitsParity2048
			infoBits = bitsInfo2048
		case ".DE3":
			parityBits = bitsParity65536
			infoBits = bitsInfo65536
		case ".DC3":
			parityBits = bitsParity65536
			infoBits = bitsInfo65536
		default:
			http.Error(w, "El tamaño de bloque es inválido", http.StatusBadRequest)
			return
		}
		if errores == 2 {
			fmt.Println("Hay dos errores solo se corregira 1")
		}
		codificarcomprimir.Decodificar(contenido, hasError, w, len(contenido), extension, parityBits, infoBits)

	// Comprimir
	case "option3":
	// Descomprimir
	case "option4":
	// Compactar y Codificar
	case "option5":
	// Descompactar y Decodificar
	case "option6":
	// Caso de error - Deja el archivo tal cual está
	case "":
	}

}

func mostrarResultadosCodComp(w http.ResponseWriter, _ *http.Request) {
	// Establecer el tipo de contenido para que se muestre en utf-8 (igual los acentos no los muestra bien)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Leer la plantilla HTML
	tmpl, err := template.ParseFiles("codificar-comprimir/resultados/resultados.html")
	if err != nil {
		http.Error(w, "No se pudo leer la plantilla HTML", http.StatusInternalServerError)
		return
	}
	// Pasar los datos a la plantilla HTML
	err = tmpl.Execute(w, result)
	if err != nil {
		http.Error(w, "No se pudo procesar la plantilla HTML", http.StatusInternalServerError)
		return
	}
}

/*
Fin de Funciones para la codificacion y compresion de un archivo
*/

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

	freqs := make(map[rune]int)
	for _, ch := range text {
		freqs[ch]++
	}
	raiz := huffman.ConstruirArbol(freqs)
	fmt.Println(text)
	compacted := huffman.Compacted(text, raiz)

	err = huffman.SaveCompacted(compacted, raiz)
	if err != nil {
		http.Error(w, "No se pudo comprimir el archivo1", http.StatusInternalServerError)
	}

	unziped, _, error := huffman.GetFromCompacted() //La raiz que se recupera aca es la que va en la linea 383 donde se define unzip

	if error != nil {
		http.Error(w, "No se pudo comprimir el archivo2", http.StatusInternalServerError)
		fmt.Println(error)
	}

	//Este es el comprimido que se mostrara en la pagina
	compact := huffman.BinaryToBytes(compacted)

	if err := ioutil.WriteFile(filepath.Join("comprimir/files", "comprimido.txt"), compact, 0644); err != nil {
		http.Error(w, "No se pudo guardar el archivo comprimido", http.StatusInternalServerError)
		return
	}

	//Este es el descomprimido que se mostrara en la pagina
	unzip := huffman.DecodeData(raiz, unziped) //Aca raiz deberia ser raizRecuperada para que sea fiel... pero raizRecuperada no esta funcionando correctamente... ya lo wa arreglar
	//fmt.Println(unzip)
	if err := ioutil.WriteFile(filepath.Join("comprimir/files", "descomprimido.txt"), []byte(unzip), 0644); err != nil {
		http.Error(w, "No se pudo guardar el archivo descomprimido.txt", http.StatusInternalServerError)
		return
	}
	//Este es el descomprimido huf
	if err := ioutil.WriteFile(filepath.Join("comprimir/resultados", "descomprimido.dhu"), []byte(unzip), 0644); err != nil {
		http.Error(w, "No se pudo guardar el archivo descomprimido.dhu", http.StatusInternalServerError)
		return
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
	comprimido, err := ioutil.ReadFile(filepath.Join("comprimir/resultados", "comprimido.huf"))
	if err != nil {
		http.Error(w, "No se pudo leer el archivo codificado", http.StatusInternalServerError)
		return
	}

	// Leer el archivo descomprimido
	decomprimido, err := ioutil.ReadFile(filepath.Join("comprimir/files", "descomprimido.txt"))
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

/*
Fin funciones para la compresion de un archivo
*/
