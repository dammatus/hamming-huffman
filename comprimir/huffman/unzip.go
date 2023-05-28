package huffman

import (
	"bufio"
	"encoding/binary"
	"os"
)

func GetFromCompacted() (string, *arbol, error) {
	//Abre el archivo para lectura
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
				data += "1"
			} else {
				data += "0"
			}
		}
	}

	data = data[:originalLength]

	raiz, err := cargarArbol(reader)

	if err != nil {
		return "", nil, err
	}

	return data, raiz, nil
}

func cargarArbol(reader *bufio.Reader) (*arbol, error) {
	flag, err := reader.ReadByte()
	if err != nil {
		return nil, err
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
	resultado := ""
	nodoActual := raiz

	dataLength := len(data)

	for i := 0; i < dataLength; i++ {
		if data[i] == '0' {
			nodoActual = nodoActual.izq
		} else {
			nodoActual = nodoActual.der
		}
		if nodoActual.izq == nil && nodoActual.der == nil {
			resultado += string(nodoActual.c)
			nodoActual = raiz
		}
	}
	return resultado
}
