package main

import (
	"fmt"
	"hamming-huffman/ambos"
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

type Resultados struct {
	Contenido     string
	Resultado     string
	OriginalPeso  string
	ResultadoPeso string
	Promedio      string
	Error         string
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
var resultado Resultados

func main() {
	// Configurar los manejadores de las rutas
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/codificar/", codificarHandler)
	http.HandleFunc("/codificar/files", archivosCodHandler)
	http.HandleFunc("/codificar/resultados", mostrarResultadosCod)
	http.HandleFunc("/comprimir/", comprimirHandler)
	http.HandleFunc("/comprimir/files", archivosCompHandler)
	http.HandleFunc("/comprimir/resultados", mostrarResultadosComp)
	http.HandleFunc("/ambos/", ambosHandler)
	http.HandleFunc("/ambos/files", archivosAmbosHandler)
	http.HandleFunc("/ambos/resultados", mostrarAmbosResultados)

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

func ambosHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Si es GET, mostramos el formulario
		http.ServeFile(w, r, "ambos/ambos.html")
	} else if r.Method == "POST" {
		// Si es POST, enviamos el archivo a la función archivosHandler
		archivosAmbosHandler(w, r)
	}
}

/*
Funciones para la codificacion y compresion de un archivo
*/
func archivosAmbosHandler(w http.ResponseWriter, r *http.Request) {
	// Extrae el archivo subido
	file, header, err := r.FormFile("archivo")
	if err != nil {
		http.Error(w, "Error al procesar el archivo subido: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Crea la carpeta "ambos/files" si no existe
	err = os.MkdirAll("ambos/files", os.ModePerm)
	if err != nil {
		http.Error(w, "No se pudo crear la carpeta en el servidor", http.StatusInternalServerError)
		return
	}

	// Crea el archivo en el servidor
	f, err := os.Create(filepath.Join("ambos/files", "archivo.txt"))
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

	// Opcion seleccionada
	opcion := r.FormValue("options")
	if err != nil {
		http.Error(w, "Error al convertir el tamaño de bloque: "+err.Error(), http.StatusBadRequest)
		return
	}

	/*
		Opciones
	*/
	fmt.Println(opcion)
	switch opcion {
	case "codificar":
		{
			// Extrae el tipo de codificación
			blockSizeStr := r.FormValue("block")
			blockSize, err := strconv.Atoi(blockSizeStr)
			if err != nil {
				http.Error(w, "Error al convertir el tamaño de bloque: "+err.Error(), http.StatusBadRequest)
				return
			}

			// Extrae el error
			errorStr := r.FormValue("Error")
			hasError := false
			dosErrores := false
			if errorStr == "error1" || errorStr == "error2" {
				hasError = true
				if errorStr == "error2" {
					dosErrores = true
				}
			}
			// Leer el contenido del archivo
			contenido, err := ioutil.ReadFile(filepath.Join("ambos/files", "archivo.txt"))
			if err != nil {
				http.Error(w, "No se pudo leer el archivo subido", http.StatusInternalServerError)
				return
			}
			// Codifica el archivo
			ambos.Codificar(w, blockSize, contenido, hasError, dosErrores)
			// Leer el contenido del archivo
			codificado, err := ioutil.ReadFile(filepath.Join("ambos/files", "codificado.txt"))
			if err != nil {
				http.Error(w, "No se pudo leer el archivo subido", http.StatusInternalServerError)
				return
			}
			archivoOriginal := "ambos/files/archivo.txt"
			// Obtener información del archivo
			info, err := os.Stat(archivoOriginal)
			if err != nil {
				fmt.Println("Error al obtener información del archivo:", err)
				return
			}
			archivoComprimido := "ambos/files/codificado.txt"
			// Obtener información del archivo
			cod, err := os.Stat(archivoComprimido)
			if err != nil {
				fmt.Println("Error al obtener información del archivo:", err)
				return
			}
			var errores string
			if hasError {
				errores = "Codificado con Error"
			} else {
				errores = "Codificado sin Error"
			}
			// Obtener el tamaño del archivo en bytes
			tamañoOriginal := info.Size()
			tamañoCodificado := cod.Size()
			codificacion := (tamañoCodificado * 100 / tamañoOriginal) - 100
			// se mostrara en el HTML
			resultado = Resultados{
				Contenido:     string(contenido),
				Resultado:     string(codificado),
				OriginalPeso:  fmt.Sprintf("Tamaño del archivo original: %d bytes", tamañoOriginal),
				ResultadoPeso: fmt.Sprintf("Tamaño del archivo Codificado: %d bytes", tamañoCodificado),
				Promedio:      fmt.Sprintf("El Archivo Creció un %d %%", codificacion),
				Error:         errores,
			}

			mostrarAmbosResultados(w, r)
		}
	case "decodificar":
		{
			// Ver extensión
			extension := filepath.Ext(header.Filename)
			hasError := false
			var parityBits, infoBits int
			switch extension {
			case ".HA1":
				parityBits = bitsParity32
				infoBits = bitsInfo32
				blockSize = 32
			case ".HE1":
				hasError = true
				parityBits = bitsParity32
				infoBits = bitsInfo32
				blockSize = 32
			case ".HA2":
				parityBits = bitsParity2048
				infoBits = bitsInfo2048
				blockSize = 2048
			case ".HE2":
				hasError = true
				parityBits = bitsParity2048
				infoBits = bitsInfo2048
				blockSize = 2048
			case ".HA3":
				parityBits = bitsParity65536
				infoBits = bitsInfo65536
				blockSize = 65536
			case ".HE3":
				hasError = true
				parityBits = bitsParity65536
				infoBits = bitsInfo65536
				blockSize = 65536
			}
			// Leer el contenido del archivo
			contenido, err := ioutil.ReadFile(filepath.Join("ambos/files", "archivo.txt"))
			if err != nil {
				http.Error(w, "No se pudo leer el archivo subido", http.StatusInternalServerError)
				return
			}

			// DEcodifica
			count, Filename := ambos.Decodificar(w, contenido, blockSize, infoBits, hasError, parityBits)
			// Leer el contenido del archivo
			decodificado, err := ioutil.ReadFile(filepath.Join("ambos/resultados", Filename))
			if err != nil {
				http.Error(w, "No se pudo leer el archivo subido", http.StatusInternalServerError)
				return
			}
			archivoOriginal := "ambos/files/archivo.txt"
			// Obtener información del archivo
			info, err := os.Stat(archivoOriginal)
			if err != nil {
				fmt.Println("Error al obtener información del archivo:", err)
				return
			}
			archivoDecodificado := "ambos/files/decodificado.txt"
			// Obtener información del archivo
			dec, err := os.Stat(archivoDecodificado)
			if err != nil {
				fmt.Println("Error al obtener información del archivo:", err)
				return
			}
			// Obtener el tamaño del archivo en bytes
			tamañoOriginal := info.Size()
			tamañoDecodificado := dec.Size()
			decodificacion := 100 - (tamañoDecodificado * 100 / tamañoOriginal)
			var errores string
			if hasError {
				errores = fmt.Sprintf("Hubo %d errores", count)
			} else {
				errores = "Decodificado sin Error"
			}
			// se mostrara en el HTML
			resultado = Resultados{
				Contenido:     string(contenido),
				Resultado:     string(decodificado),
				OriginalPeso:  fmt.Sprintf("Tamaño del archivo original: %d bytes", tamañoOriginal),
				ResultadoPeso: fmt.Sprintf("Tamaño del archivo Decodificado: %d bytes", tamañoDecodificado),
				Promedio:      fmt.Sprintf("El Archivo Disminuyó un %d %%", decodificacion),
				Error:         errores,
			}

			mostrarAmbosResultados(w, r)
		}
	case "comprimir":
		{
			// Leer el contenido del archivo
			contenido, err := ioutil.ReadFile(filepath.Join("ambos/files", "archivo.txt"))
			if err != nil {
				http.Error(w, "No se pudo leer el archivo subido", http.StatusInternalServerError)
				return
			}
			// Comprime
			ambos.Comprimir(w, contenido)
			// Leer el contenido del archivo
			zip, err := ioutil.ReadFile(filepath.Join("ambos/resultados", "comprimido.huf"))
			if err != nil {
				http.Error(w, "No se pudo leer el archivo subido", http.StatusInternalServerError)
				return
			}
			archivoOriginal := "ambos/files/archivo.txt"
			// Obtener información del archivo
			info, err := os.Stat(archivoOriginal)
			if err != nil {
				fmt.Println("Error al obtener información del archivo:", err)
				return
			}
			archivoComprimido := "ambos/resultados/comprimido.huf"
			// Obtener información del archivo
			comp, err := os.Stat(archivoComprimido)
			if err != nil {
				fmt.Println("Error al obtener información del archivo:", err)
				return
			}
			// Obtener el tamaño del archivo en bytes
			tamañoOriginal := info.Size()
			tamañoComprimido := comp.Size()
			compresion := 100 - (tamañoComprimido * 100 / tamañoOriginal)

			// se mostrara en el HTML
			resultado = Resultados{
				Contenido:     string(contenido),
				Resultado:     string(zip),
				OriginalPeso:  fmt.Sprintf("Tamaño del archivo original: %d bytes", tamañoOriginal),
				ResultadoPeso: fmt.Sprintf("Tamaño del archivo comprimido: %d bytes", tamañoComprimido),
				Promedio:      fmt.Sprintf("El Archivo se Compactó un %d %%", compresion),
			}

			mostrarAmbosResultados(w, r)
		}
	case "descomprimir":
		{
			// Leer el contenido del archivo
			contenido, err := ioutil.ReadFile(filepath.Join("ambos/files", "archivo.txt"))
			if err != nil {
				http.Error(w, "No se pudo leer el archivo subido", http.StatusInternalServerError)
				return
			}
			ambos.Descomprimir(w, string(contenido))
			// Leer el contenido del archivo
			unzip, err := ioutil.ReadFile(filepath.Join("ambos/resultados", "descomprimido.dhu"))
			if err != nil {
				http.Error(w, "No se pudo leer el archivo subido", http.StatusInternalServerError)
				return
			}
			archivoOriginal := "ambos/files/archivo.txt"
			// Obtener información del archivo
			info, err := os.Stat(archivoOriginal)
			if err != nil {
				fmt.Println("Error al obtener información del archivo:", err)
				return
			}
			archivoDescomprimido := "ambos/resultados/descomprimido.dhu"
			// Obtener información del archivo
			desc, err := os.Stat(archivoDescomprimido)
			if err != nil {
				fmt.Println("Error al obtener información del archivo:", err)
				return
			}
			// Obtener el tamaño del archivo en bytes
			tamañoOriginal := info.Size()
			tamañoDescomprimido := desc.Size()

			descompresion := 100 - (tamañoDescomprimido * 100 / tamañoOriginal)
			// se mostrara en el HTML
			resultado = Resultados{
				Contenido:     string(contenido),
				Resultado:     string(unzip),
				OriginalPeso:  fmt.Sprintf("Tamaño del archivo Comprimido: %d bytes", tamañoOriginal),
				ResultadoPeso: fmt.Sprintf("Tamaño del archivo Descomprimido: %d bytes", tamañoDescomprimido),
				Promedio:      fmt.Sprintf("El Archivo Descompactó un %d %%", descompresion),
			}
			mostrarAmbosResultados(w, r)
		}
	case "codificarComprimir":
		{
			// Se coprimira y luego se codificara el archivo codificado
			// Extrae el tipo de codificación
			blockSizeStr := r.FormValue("block")
			blockSize, err := strconv.Atoi(blockSizeStr)
			if err != nil {
				http.Error(w, "Error al convertir el tamaño de bloque: "+err.Error(), http.StatusBadRequest)
				return
			}

			// Extrae el error
			errorStr := r.FormValue("Error")
			hasError := false
			dosErrores := false
			if errorStr == "error1" || errorStr == "error2" {
				hasError = true
				if errorStr == "error2" {
					dosErrores = true
				}
			}
			// Leer el contenido del archivo
			contenido, err := ioutil.ReadFile(filepath.Join("ambos/files", "archivo.txt"))
			if err != nil {
				http.Error(w, "No se pudo leer el archivo subido", http.StatusInternalServerError)
				return
			}
			// Comprime
			ambos.Comprimir(w, contenido)
			// Leer el contenido del archivo
			zip, err := ioutil.ReadFile(filepath.Join("ambos/resultados", "comprimido.huf"))
			if err != nil {
				http.Error(w, "No se pudo leer el archivo subido", http.StatusInternalServerError)
				return
			}
			// Codifica el archivo
			ambos.Codificar(w, blockSize, zip, hasError, dosErrores)
			// Leer el contenido del archivo
			codificado, err := ioutil.ReadFile(filepath.Join("ambos/files", "codificado.txt"))
			if err != nil {
				http.Error(w, "No se pudo leer el archivo subido", http.StatusInternalServerError)
				return
			}
			archivoOriginal := "ambos/files/archivo.txt"
			// Obtener información del archivo
			info, err := os.Stat(archivoOriginal)
			if err != nil {
				fmt.Println("Error al obtener información del archivo:", err)
				return
			}
			archivoComprimido := "ambos/files/codificado.txt"
			// Obtener información del archivo
			cod, err := os.Stat(archivoComprimido)
			if err != nil {
				fmt.Println("Error al obtener información del archivo:", err)
				return
			}
			var errores string
			if hasError {
				errores = "Codificado con Error"
			} else {
				errores = "Codificado sin Error"
			}
			// Obtener el tamaño del archivo en bytes
			tamañoOriginal := info.Size()
			tamañoCodificado := cod.Size()
			codificacion := (tamañoCodificado * 100 / tamañoOriginal) - 100
			// se mostrara en el HTML
			resultado = Resultados{
				Contenido:     string(contenido),
				Resultado:     string(codificado),
				OriginalPeso:  fmt.Sprintf("Tamaño del archivo original: %d bytes", tamañoOriginal),
				ResultadoPeso: fmt.Sprintf("Tamaño del archivo original: %d bytes", tamañoCodificado),
				Promedio:      fmt.Sprintf("El Archivo cambió un %d %%", codificacion),
				Error:         errores,
			}
			mostrarAmbosResultados(w, r)
		}
	case "descomprimirDecodificar":
		{
			// Ver extensión
			extension := filepath.Ext(header.Filename)
			hasError := false
			var parityBits, infoBits int
			switch extension {
			case ".HA1":
				parityBits = bitsParity32
				infoBits = bitsInfo32
				blockSize = 32
			case ".HE1":
				hasError = true
				parityBits = bitsParity32
				infoBits = bitsInfo32
				blockSize = 32
			case ".HA2":
				parityBits = bitsParity2048
				infoBits = bitsInfo2048
				blockSize = 2048
			case ".HE2":
				hasError = true
				parityBits = bitsParity2048
				infoBits = bitsInfo2048
				blockSize = 2048
			case ".HA3":
				parityBits = bitsParity65536
				infoBits = bitsInfo65536
				blockSize = 65536
			case ".HE3":
				hasError = true
				parityBits = bitsParity65536
				infoBits = bitsInfo65536
				blockSize = 65536
			}
			// Leer el contenido del archivo
			contenido, err := ioutil.ReadFile(filepath.Join("ambos/files", "archivo.txt"))
			if err != nil {
				http.Error(w, "No se pudo leer el archivo subido", http.StatusInternalServerError)
				return
			}

			// DEcodifica
			count, _ := ambos.Decodificar(w, contenido, blockSize, infoBits, hasError, parityBits)
			// Leer el contenido del archivo
			decodificado, err := ioutil.ReadFile(filepath.Join("ambos/files", "decodificado.txt"))
			if err != nil {
				http.Error(w, "No se pudo leer el archivo subido", http.StatusInternalServerError)
				return
			}
			//
			if err := ioutil.WriteFile(filepath.Join("ambos/files", "archivo.txt"), decodificado, 0644); err != nil {
				http.Error(w, "No se pudo guardar el archivo descomprimido.txt", http.StatusInternalServerError)
				return
			}
			ambos.Descomprimir(w, string(decodificado))
			// Leer el contenido del archivo
			unzip, err := ioutil.ReadFile(filepath.Join("ambos/resultados", "descomprimido.dhu"))
			if err != nil {
				http.Error(w, "No se pudo leer el archivo subido desc", http.StatusInternalServerError)
				return
			}
			archivoOriginal := "ambos/files/archivo.txt"
			// Obtener información del archivo
			info, err := os.Stat(archivoOriginal)
			if err != nil {
				fmt.Println("Error al obtener información del archivo:", err)
				return
			}
			archivoDescomprimido := "ambos/resultados/descomprimido.dhu"
			// Obtener información del archivo
			desc, err := os.Stat(archivoDescomprimido)
			if err != nil {
				fmt.Println("Error al obtener información del archivo:", err)
				return
			}
			// Obtener el tamaño del archivo en bytes
			tamañoOriginal := info.Size()
			tamañoDescomprimido := desc.Size()

			descompresion := 100 - (tamañoDescomprimido * 100 / tamañoOriginal)
			// se mostrara en el HTML
			resultado = Resultados{
				Contenido:     string(contenido),
				Resultado:     string(unzip),
				OriginalPeso:  fmt.Sprintf("Tamaño del archivo Comprimido: %d bytes", tamañoOriginal),
				ResultadoPeso: fmt.Sprintf("Tamaño del archivo Descomprimido: %d bytes", tamañoDescomprimido),
				Promedio:      fmt.Sprintf("El Archivo cambió un %d %%", descompresion),
				Error:         fmt.Sprintf("Hubo %d errores", count),
			}
			mostrarAmbosResultados(w, r)
		}
	default:
		{
			fmt.Println("Error de opcion")
		}

	}

}

func mostrarAmbosResultados(w http.ResponseWriter, _ *http.Request) {
	// Leer la plantilla HTML
	tmpl, err := template.ParseFiles("ambos/resultados/resultados.html")
	if err != nil {
		http.Error(w, "No se pudo leer la plantilla HTML", http.StatusInternalServerError)
		return
	}

	// Pasar los datos a la plantilla HTML
	err = tmpl.Execute(w, resultado)
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
	fmt.Println(bits)
	encode := hamming.AplicandoHamming(bits, blockSize, parityBits, infoBits, hasError, false)
	fmt.Println(encode)
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
	if err := ioutil.WriteFile(filepath.Join("codificar/files", "codificado.txt"), ascii, 0644); err != nil {
		http.Error(w, "No se pudo guardar el archivo codificado", http.StatusInternalServerError)
		return
	}
	//Este es el que cumple con la consigna, se guarda en la carpeta resultados
	if err := ioutil.WriteFile(filepath.Join("codificar/resultados", codificadoFileName), ascii, 0644); err != nil {
		http.Error(w, "No se pudo guardar el archivo codificado", http.StatusInternalServerError)
		return
	}

	// Decodificar el contenido y escribirlo en un archivo (Sin corregir)
	decode, _ := hamming.DecodeHamming(encode, blockSize, infoBits, false, parityBits)
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
		decode, _ = hamming.DecodeHamming(encode, blockSize, infoBits, hasError, parityBits)
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

	compacted := huffman.Compacted(text, raiz)

	//fmt.Println("Compactado: " + binary)

	fmt.Println("Codigo Huffman: ", compacted)

	err = huffman.SaveCompacted(compacted, raiz)
	if err != nil {
		fmt.Println("Error al guardar el archivo: ", err)
	} else {
		fmt.Println("Datos comprimidos exitosamente!")
	}

	unziped, raizRecuperada, error := huffman.GetFromCompacted()

	if error == nil {
		fmt.Println("Recuperados del archivo: ", unziped)
	}

	fmt.Println("Resultado: ", huffman.DecodeData(raizRecuperada, unziped))

	//Este es el comprimido que se mostrara en la pagina
	compact := huffman.BinaryToBytes(compacted)

	if err := ioutil.WriteFile(filepath.Join("comprimir/files", "comprimido.txt"), compact, 0644); err != nil {
		http.Error(w, "No se pudo guardar el archivo comprimido", http.StatusInternalServerError)
		return
	}

	//Este es el descomprimido que se mostrara en la pagina
	unzip := huffman.DecodeData(raizRecuperada, unziped)

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
	comprimido, err := ioutil.ReadFile(filepath.Join("comprimir/files", "comprimido.txt"))
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
