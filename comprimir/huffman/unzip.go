package huffman

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"strings"
)

func GetFromCompacted() (string, *arbol, error) {
	//Abre el archivo para lectura
	var builder strings.Builder
	file, err := os.Open("./comprimir/resultados/comprimido.huf")

	if err != nil {
		return "", nil, err
	}
	defer file.Close()

	//Crea un lector de búfer
	reader := bufio.NewReader(file)

	//Lee la longitud original de los datos del archivo
	lengthBytes := make([]byte, 4)
	_, err = reader.Read(lengthBytes)
	if err != nil {
		return "", nil, err
	}

	originalLength := int(binary.LittleEndian.Uint32(lengthBytes))

	dataBytes := make([]byte, 1024)
	var acumulador []byte
	/* _, err = reader.Read(dataBytes)
	if err != nil {
		return "", err
	} */
	for {
		n, err := reader.Read(dataBytes)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println("Error al leer el archivo: ", err)
				return "", nil, err
			}
		}
		if n == 0 {
			break
		}
		acumulador = append(acumulador, dataBytes...)
		if len(acumulador) >= originalLength-1024 {
			break
		}
	}

	leeEstoTambien := make([]byte, originalLength-len(acumulador))
	reader.Read(leeEstoTambien)
	acumulador = append(acumulador, leeEstoTambien...)

	data := ""

	//Lee los bytes del archivo y los convierte en una cadena binaria
	for _, byteVal := range acumulador {
		for j := 7; j >= 0; j-- {
			if (byteVal >> j & 1) == 1 {
				builder.WriteByte('1')
			} else {
				builder.WriteByte('0')
			}
		}
	}
	data = builder.String()
	data = data[:originalLength]
	raiz, err := cargarArbol(reader)

	if err != nil {
		if err == io.EOF {
			err = nil
		} else {
			return "", nil, err
		}
	}
	return data, raiz, nil
}

func cargarArbol(reader *bufio.Reader) (*arbol, error) {
	flag, err := reader.ReadByte()
	if err != nil {
		if err == io.EOF {
			err = nil
		} else {
			return nil, err
		}
	}

	if flag == 0 { //nodo nulo
		return nil, nil
	}
	freq, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	caracter, _, err := reader.ReadRune()
	if err != nil {
		return nil, err
	}

	hijoIzq, err := cargarArbol(reader)
	if err != nil {
		return nil, err
	}

	hijoDer, err := cargarArbol(reader)
	if err != nil {
		return nil, err
	}

	return &arbol{int(freq), caracter, hijoIzq, hijoDer}, nil
}

func DecodeData(raiz *arbol, data string) string {
	var resultado bytes.Buffer
	nodoActual := raiz

	dataLength := len(data)

	for i := 0; i < dataLength; i++ {
		if data[i] == '0' {
			nodoActual = nodoActual.izq
		} else {
			nodoActual = nodoActual.der
		}
		if nodoActual.izq == nil && nodoActual.der == nil {
			resultado.WriteRune(nodoActual.c)
			nodoActual = raiz
		}
	}
	return resultado.String()
}

func LoadMap() (map[rune]int, error) {
	file, err := os.Open("./comprimir/resultados/freq.dat")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)

	var freq map[rune]int
	err = decoder.Decode(&freq)
	if err != nil {
		return nil, err
	}
	return freq, nil
}
