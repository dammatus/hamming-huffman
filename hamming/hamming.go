package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"text/template"
)

type Resultados struct {
	Contenido    string
	Codificado   string
	Decodificado string
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

// Abre el archivo a codificar
func readFile(file string) string {
	// Lee el contenido del archivo
	datos, err := ioutil.ReadFile(file)
	if err != nil {
		// Si ocurre un error, devuelve una cadena vacía
		return ""
	}
	// Convierte el slice de bytes a un string y lo devuelve
	return string(datos)
}

// Escribe en un archivo la codificacion
func writeFile(file string, datos string) error {
	// Escribe el contenido en el archivo
	err := ioutil.WriteFile(file, []byte(datos), 0644)
	if err != nil {
		// Si ocurre un error, devuelve el error
		return err
	}
	// Si no hay errores, devuelve nil
	return nil
}

func main() {

	// Handlers
	http.HandleFunc("/files", controlarArchivo)
	http.HandleFunc("/", controlarHamming)
	http.HandleFunc("/resultados", mostrarResultados)

	// Iniciar el servidor
	fmt.Println("Servidor escuchando en http://localhost:8080")
	http.ListenAndServe(":8080", nil)

}

func controlarHamming(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Si es GET, mostramos el formulario
		http.ServeFile(w, r, "hamming.html")
	} else if r.Method == "POST" {
		// Si es POST, enviamos el archivo a la función controlarArchivo
		controlarArchivo(w, r)
	}
}

func controlarArchivo(w http.ResponseWriter, r *http.Request) {

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

	// Extrae el tipo de codificación
	blockSizeStr := r.FormValue("tipoCod")
	blockSize, err := strconv.Atoi(blockSizeStr)
	if err != nil {
		http.Error(w, "Error al convertir el tamaño de bloque: "+err.Error(), http.StatusBadRequest)
		return
	}
	// Extrae el error
	errorStr := r.FormValue("error")
	error := false
	if errorStr == "Si" {
		error = true
	}

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

	// Aplicar Hamming
	var parityBits int
	var infoBits int
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
	contenido, err := ioutil.ReadFile(filepath.Join("files", handler.Filename)) //handler.Filename tendrá "archivo.txt"
	if err != nil {
		http.Error(w, "No se pudo leer el archivo subido", http.StatusInternalServerError)
		return
	}

	// Convertir el contenido a bits y aplicar Hamming
	bits := modulesHamming.ByteToBits(contenido, blockSize)
	encode := modulesHamming.AplicandoHamming(bits, blockSize, parityBits, infoBits, error)

	// Convertir el resultado a texto y escribirlo en un archivo
	ascii := modulesHamming.BinToASCII(encode)
	if err := ioutil.WriteFile(filepath.Join("files", "codificado.txt"), []byte(ascii), 0644); err != nil {
		http.Error(w, "No se pudo guardar el archivo codificado", http.StatusInternalServerError)
		return
	}
	switch blockSize {
	case 32:
		if err := ioutil.WriteFile(filepath.Join("resultados", "codificado.HA1"), []byte(ascii), 0644); err != nil {
			http.Error(w, "No se pudo guardar el archivo codificado", http.StatusInternalServerError)
			return
		}
	case 2048:
		if err := ioutil.WriteFile(filepath.Join("resultados", "codificado.HA2"), []byte(ascii), 0644); err != nil {
			http.Error(w, "No se pudo guardar el archivo codificado", http.StatusInternalServerError)
			return
		}
	case 65536:
		if err := ioutil.WriteFile(filepath.Join("resultados", "codificado.HA3"), []byte(ascii), 0644); err != nil {
			http.Error(w, "No se pudo guardar el archivo codificado", http.StatusInternalServerError)
			return
		}
	}
	// Decodificar el contenido y escribirlo en un archivo (Sin corregir)
	decode := modulesHamming.DecodeHamming(encode, blockSize, infoBits, false, parityBits)
	asciiDeco := modulesHamming.BitsToByte(decode)
	decoded := string(asciiDeco)
	if err := ioutil.WriteFile(filepath.Join("files", "decodificadoSC.txt"), []byte(decoded[:len(contenido)]), 0644); err != nil {
		http.Error(w, "No se pudo guardar el archivo decodificado", http.StatusInternalServerError)
		return
	}
	switch blockSize {
	case 32:
		if err := ioutil.WriteFile(filepath.Join("resultados", "decodificado.DE1"), []byte(decoded[:len(contenido)]), 0644); err != nil {
			http.Error(w, "No se pudo guardar el archivo decodificado", http.StatusInternalServerError)
			return
		}
	case 2048:
		if err := ioutil.WriteFile(filepath.Join("resultados", "decodificado.DE2"), []byte(decoded[:len(contenido)]), 0644); err != nil {
			http.Error(w, "No se pudo guardar el archivo decodificado", http.StatusInternalServerError)
			return
		}
	case 65536:
		if err := ioutil.WriteFile(filepath.Join("resultados", "decodificado.DE3"), []byte(decoded[:len(contenido)]), 0644); err != nil {
			http.Error(w, "No se pudo guardar el archivo decodificado", http.StatusInternalServerError)
			return
		}
	}
	if error {
		// Decodificar el contenido y escribirlo en un archivo (Corregido)
		decode = modulesHamming.DecodeHamming(encode, blockSize, infoBits, error, parityBits)
		asciiDeco = modulesHamming.BitsToByte(decode)
		decoded = string(asciiDeco)
		if err := ioutil.WriteFile(filepath.Join("files", "decodificado.txt"), []byte(decoded[:len(contenido)]), 0644); err != nil {
			http.Error(w, "No se pudo guardar el archivo decodificado", http.StatusInternalServerError)
			return
		}
		switch blockSize {
		case 32:
			if err := ioutil.WriteFile(filepath.Join("resultados", "decodificado.DC1"), []byte(decoded[:len(contenido)]), 0644); err != nil {
				http.Error(w, "No se pudo guardar el archivo decodificado", http.StatusInternalServerError)
				return
			}
		case 2048:
			if err := ioutil.WriteFile(filepath.Join("resultados", "decodificado.DC2"), []byte(decoded[:len(contenido)]), 0644); err != nil {
				http.Error(w, "No se pudo guardar el archivo decodificado", http.StatusInternalServerError)
				return
			}
		case 65536:
			if err := ioutil.WriteFile(filepath.Join("resultados", "decodificado.DC3"), []byte(decoded[:len(contenido)]), 0644); err != nil {
				http.Error(w, "No se pudo guardar el archivo decodificado", http.StatusInternalServerError)
				return
			}
		}
	}

	fmt.Println("***************************************************")
	fmt.Println(bits[:infoBits]) // Texto a codificar
	fmt.Println("***************************************************")
	fmt.Println(encode[:blockSize]) // Codificacion
	fmt.Println("***************************************************")
	fmt.Println(decode[:blockSize]) // Error en codificacion
	fmt.Println("***************************************************")

	// Servir el archivo HTML con los resultados
	http.HandlerFunc(mostrarResultados).ServeHTTP(w, r)

}

func mostrarResultados(w http.ResponseWriter, _ *http.Request) {

	// Establecer el tipo de contenido para que se muestre en utf-8 (igual loas acentos no los muestra bien)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Leer el archivo original
	contenido, err := ioutil.ReadFile(filepath.Join("files", "archivo.txt"))
	if err != nil {
		http.Error(w, "No se pudo leer el archivo ", http.StatusInternalServerError)
		return
	}

	// Leer el archivo codificado
	codificado, err := ioutil.ReadFile(filepath.Join("files", "codificado.txt"))
	if err != nil {
		http.Error(w, "No se pudo leer el archivo codificado", http.StatusInternalServerError)
		return
	}

	// Leer el archivo decodificado
	decodificado, err := ioutil.ReadFile(filepath.Join("files", "decodificado.txt"))
	if err != nil {
		http.Error(w, "No se pudo leer el archivo decodificado", http.StatusInternalServerError)
		return
	}

	// Crear un mapa de datos para la plantilla HTML
	data := map[string]string{
		"Contenido":    string(contenido),
		"Codificado":   string(codificado),
		"Decodificado": string(decodificado),
	}

	// Leer la plantilla HTML
	tmpl, err := template.ParseFiles("resultados/resultados.html")
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
