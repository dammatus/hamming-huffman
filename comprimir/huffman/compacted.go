package huffman

import (
	"bufio"
	"os"
	"strconv"
)

func Compacted(texto string, arbol *arbol) string {
	codigos := make(map[rune]string)
	obtenerCodigos(arbol, "", codigos) // Generar los códigos Huffman a partir del árbol

	var compactado string

	for _, ch := range texto {
		compactado += codigos[ch]
	}

	return compactado
}

func obtenerCodigos(nodo *arbol, prefijo string, codigos map[rune]string) {
	if nodo == nil {
		return
	}

	if nodo.izq == nil && nodo.der == nil {
		// Es un nodo hoja, asignar el prefijo como código Huffman para el carácter
		codigos[nodo.c] = prefijo
		return
	}
	// Recursivamente obtener los códigos Huffman de los subárboles izquierdo y derecho
	obtenerCodigos(nodo.izq, prefijo+"0", codigos)
	obtenerCodigos(nodo.der, prefijo+"1", codigos)
}

func BinaryToBytes(binaryString string) []byte {
	// Convierte el string binario en un número entero sin signo
	number, _ := strconv.ParseUint(binaryString, 2, len(binaryString))

	// Convierte el número entero en un slice de bytes
	bytes := make([]byte, len(binaryString)/8)
	for i := range bytes {
		bytes[i] = byte(number >> (8 * (len(bytes) - 1 - i)))
	}

	return bytes
}

func SaveCompacted(compacted string) error {
	file, err := os.Create("./comprimir/resultados/comprimido.huf")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for i := 0; i < len(compacted); i += 8 {
		if i+8 <= len(compacted) {
			byteStr := compacted[i : i+8]
			byteVal := byte(0)
			for j := 0; j < 8; j++ {
				if byteStr[j] == '1' {
					byteVal |= 1 << (7 - j)
				}
			}
			writer.WriteByte(byteVal)
		}
	}
	err = writer.Flush()
	if err != nil {
		return err
	}
	return nil
}