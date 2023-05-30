package huffman

import (
	"bufio"
	"bytes"
	"encoding/binary"
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

	dataBytes := make([]byte, (originalLength+7)/8)
	_, err = reader.Read(dataBytes)
	if err != nil {
		return "", nil, err
	}

	data := ""

	//Lee los bytes del archivo y los convierte en una cadena binaria
	for _, byteVal := range dataBytes {
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
	caracter, err := reader.ReadByte()
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

	return &arbol{0, rune(caracter), hijoIzq, hijoDer}, nil
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
			resultado.WriteByte(byte(nodoActual.c))
			nodoActual = raiz
		}
	}
	return resultado.String()
}
