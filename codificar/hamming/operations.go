package hamming

import (
	"math"
)

// Convierte de binario a texto
// func BinToASCII(bin []byte) string {
// 	/*
// 	   devuelve los mismo que la funcion de abajo, solo que de tipo string
// 	*/
// 	ascii := ""
// 	for i := 0; i < len(bin); i += 8 { // se recorre el slice de 8 bits en 8
// 		end := i + 8
// 		if end > len(bin) {
// 			end = len(bin)
// 		}
// 		bits := bin[i:end]
// 		n := byte(0)                     // inicializo un byte en cero que sera usado para construir el byte ASCII correspondiente
// 		for j := 0; j < len(bits); j++ { // para cada uno, se desplaza n un bit hacia la izquierda y se verifica si el bit es 1 o 0
// 			n <<= 1
// 			if bits[j] == 0x01 { // Si el bit es 1, se utiliza el operador "|" para poner el último bit de "n" en 1.
// 				n |= 1
// 			}
// 		}
// 		ascii += string(n) // se agrega el byte n al string ascii
// 	}

//		return ascii
//	}

// Convierte de binario a texto
func BinToASCII(bin []byte) []byte {
	ascii := make([]byte, 0)
	for i := 0; i < len(bin); i += 8 { // se recorre el slice de 8 bits en 8
		end := i + 8
		if end > len(bin) {
			end = len(bin)
		}
		bits := bin[i:end]
		n := byte(0)                     // inicializo un byte en cero que sera usado para construir el byte ASCII correspondiente
		for j := 0; j < len(bits); j++ { // para cada uno, se desplaza n un bit hacia la izquierda y se verifica si el bit es 1 o 0
			n <<= 1
			if bits[j] == 0x01 { // Si el bit es 1, se utiliza el operador "|" para poner el último bit de "n" en 1.
				n |= 1
			}
		}
		ascii = append(ascii, n) // se agrega el byte n al slice ascii
	}

	return ascii
}

// Convierte de bytes a bits
func ByteToBits(slice []byte, blockSize int) []byte {
	/*
		Convierte un slice de bytes en una secuencia de bits.
		Para cada byte del slice, toma sus bits de izquierda a derecha
		y los agrega a un slice de bytes llamado "bits"
	*/

	bits := make([]byte, 0, len(slice)*8)

	for _, b := range slice {
		for i := 7; i >= 0; i-- {
			bit := (b >> uint(i)) & 1
			bits = append(bits, bit)
		}
	}
	return bits
}

func BitsToByte(bits []byte) []byte {
	/*
		Convierte una secuencia de bits en un slice de bytes.
		Agrupa los bits de 8 en 8 y los convierte en un byte.
	*/

	// Aseguramos que la longitud de bits sea múltiplo de 8
	numBits := len(bits)
	if numBits%8 != 0 {
		pad := make([]byte, 8-(numBits%8))
		bits = append(bits, pad...)
	}

	// Convertimos los bits en bytes
	numBytes := len(bits) / 8
	bytes := make([]byte, numBytes)
	for i := 0; i < numBytes; i++ {
		for j := 0; j < 8; j++ {
			bit := bits[(i*8)+j]
			bytes[i] = (bytes[i] << 1) | bit
		}
	}
	return bytes
}

func binaryToDecimal(binary []byte) int {
	var decimal int
	for i := len(binary) - 1; i >= 0; i-- {
		decimal += int(binary[i]) * int(math.Pow(2, float64(len(binary)-i-1)))
	}
	return decimal
}
