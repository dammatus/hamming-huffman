package codificarcomprimir

import (
	"hamming-huffman/codificar/hamming"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

func Codificar(contenido []byte, parityBits int, infoBits int, blockSize int, hasError bool, w http.ResponseWriter) {
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
	if err := ioutil.WriteFile(filepath.Join("codificar-comprimir/files", "codificado.txt"), []byte(ascii), 0644); err != nil {
		http.Error(w, "No se pudo guardar el archivo codificado", http.StatusInternalServerError)
		return
	}
	//Este es el que cumple con la consigna, se guarda en la carpeta resultados
	if err := ioutil.WriteFile(filepath.Join("codificar-comprimir/resultados", codificadoFileName), []byte(ascii), 0644); err != nil {
		http.Error(w, "No se pudo guardar el archivo codificado", http.StatusInternalServerError)
		return
	}

}

func Decodificar(contenido []byte, hasError bool, w http.ResponseWriter, tam int, extension string, parityBits int, infoBits int) {
	var blockSize int //determinar dependiendo de la extension

	// Decodificar el contenido y escribirlo en un archivo (Sin corregir)
	decode := hamming.DecodeHamming(contenido, blockSize, infoBits, false, parityBits)
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
	if err := ioutil.WriteFile(filepath.Join("codificar-comprimir/files", "decodificado.txt"), []byte(decoded[:tam]), 0644); err != nil {
		http.Error(w, "No se pudo guardar el archivo decodificado", http.StatusInternalServerError)
		return
	}
	//Este es el que cumple con la consigna, se guarda en la carpeta resultados
	if err := ioutil.WriteFile(filepath.Join("codificar-comprimir/resultados", decodificadoFileName), []byte(decoded[:tam]), 0644); err != nil {
		http.Error(w, "No se pudo guardar el archivo decodificado", http.StatusInternalServerError)
		return
	}

	if hasError {
		// Decodificar el contenido y escribirlo en un archivo (Corregido)
		decode = hamming.DecodeHamming(contenido, blockSize, infoBits, hasError, parityBits)
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
}
