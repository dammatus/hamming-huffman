package ambos

import (
	"hamming-huffman/codificar/hamming"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

const (
	bitsParity32    = 5
	bitsParity2048  = 11
	bitsParity65536 = 16
	bitsInfo32      = 26
	bitsInfo2048    = 2036
	bitsInfo65536   = 65519
)

func Codificar(w http.ResponseWriter, blockSize int, contenido []byte, hasError bool) {
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

	// Convertir el contenido a bits y aplicar Hamming
	bits := hamming.ByteToBits(contenido, blockSize)
	encode := hamming.AplicandoHamming(bits, blockSize, parityBits, infoBits, hasError)

	// Convertir el resultado a texto y escribirlo en un archivo
	ascii := hamming.BinToASCII(encode)

	//Este es el que se mostrara en la pagina
	if err := ioutil.WriteFile(filepath.Join("ambos/files", "codificado.txt"), []byte(ascii), 0644); err != nil {
		http.Error(w, "No se pudo guardar el archivo codificado", http.StatusInternalServerError)
		return
	}
	switch blockSize {
	case 32:
		if hasError {
			if err := ioutil.WriteFile(filepath.Join("ambos/resultados", "codificado.HE1"), []byte(ascii), 0644); err != nil {
				http.Error(w, "No se pudo guardar el archivo codificado", http.StatusInternalServerError)
				return
			}
		} else {
			if err := ioutil.WriteFile(filepath.Join("ambos/resultados", "codificado.HA1"), []byte(ascii), 0644); err != nil {
				http.Error(w, "No se pudo guardar el archivo codificado", http.StatusInternalServerError)
				return
			}
		}
	case 2048:
		if hasError {
			if err := ioutil.WriteFile(filepath.Join("ambos/resultados", "codificado.HE2"), []byte(ascii), 0644); err != nil {
				http.Error(w, "No se pudo guardar el archivo codificado", http.StatusInternalServerError)
				return
			}
		} else {
			if err := ioutil.WriteFile(filepath.Join("ambos/resultados", "codificado.HA2"), []byte(ascii), 0644); err != nil {
				http.Error(w, "No se pudo guardar el archivo codificado", http.StatusInternalServerError)
				return
			}
		}
	case 65536:
		if hasError {
			if err := ioutil.WriteFile(filepath.Join("ambos/resultados", "codificado.HE3"), []byte(ascii), 0644); err != nil {
				http.Error(w, "No se pudo guardar el archivo codificado", http.StatusInternalServerError)
				return
			}
		} else {
			if err := ioutil.WriteFile(filepath.Join("ambos/resultados", "codificado.HA3"), []byte(ascii), 0644); err != nil {
				http.Error(w, "No se pudo guardar el archivo codificado", http.StatusInternalServerError)
				return
			}
		}
	default:
		http.Error(w, "El tamaño de bloque es inválido", http.StatusBadRequest)
		return
	}

}

func Decodificar(w http.ResponseWriter, encode string, blockSize int, infoBits int, hasError bool, parityBits int) {
	// Decodificar el contenido y escribirlo en un archivo (Sin corregir)
	bin := hamming.ASCIIToBin(encode)
	decode := hamming.DecodeHamming(bin, blockSize, infoBits, false, parityBits)
	asciiDeco := hamming.BitsToByte(decode)
	decoded := string(asciiDeco)
	//Este es el que se mostrara en la pagina
	if err := ioutil.WriteFile(filepath.Join("ambos/files", "decodificado.txt"), []byte(decoded[:len(bin)]), 0644); err != nil {
		http.Error(w, "No se pudo guardar el archivo decodificado", http.StatusInternalServerError)
		return
	}

}
