package hamming

import (
	"math"
	"math/rand"
	"time"
)

/*
*
La idea es que haga el xor contando los primeros 2^n
Despues de un salto de 2^n, porque esos no se cuentan
y siga contando 2^n. Empezando desde la posición 2^n
*
*/

/*
para i desde 0 hasta len(bloque) con paso i++

	contador = contador xor encoded[i]
	j++
	if j == 4
		i += j
		j = 0
*/
func aplicaBitsDeControl(n int, encoded []byte, blockSize int) []byte {
	var contador byte
	var j = 0
	for i := (int(math.Pow(2, float64(n))) - 1); i < blockSize; i++ {

		contador ^= encoded[i]
		j++
		if j == int(math.Pow(2, float64(n))) {
			i += j
			j = 0
		}
	}
	encoded[int(math.Pow(2, float64(n)))-1] = contador
	return encoded
}

// Aplica el bit de paridad en el final
func aplicarBitDePariedad(encoded []byte) []byte {
	var contador byte
	for _, i := range encoded {
		contador ^= encoded[i] //XOR para ver paridad
	}
	encoded[len(encoded)-1] = contador
	return encoded
}

// Aplica, a un bloque, los bits de control en sus debidas posiciones
func encode(info []byte, parityBits int, blockSize int) []byte {

	var encoded = make([]byte, blockSize)
	var j = 0
	for i := 0; i < blockSize; i++ {
		if !isPowerOfTwo(i + 1) {
			encoded[i] = info[j]
			j++
		}
	}
	for i := 0; i < parityBits; i++ {
		encoded = aplicaBitsDeControl(i, encoded, blockSize)
	}
	encoded = aplicarBitDePariedad(encoded)

	return encoded
}

// Realiza todo el proceso de Hamming para todos los bloques encesarios
func AplicandoHamming(info []byte, blockSize int, parityBits int, infoBits int, error bool, dosErrores bool) []byte {
	/*
		Esta funcion tomara todos los bloques de informacion, y le aplicara hamming
		para luego concatenarlos en un slice "encoded", que contendra toda la cadena de info
		ya hammingnizada (si esa palabra existe).
	*/

	encoded := make([]byte, 0)
	/*
		si el slice info es más chico que blockSize, se crea un nuevo slice llamado completado de tamaño blockSize
		se copian los elementos de info al inicio del slice y se rellena el resto del slice con ceros
	*/
	if len(info) < blockSize {

		completado := make([]byte, blockSize) // crear slice de blockSize con ceros

		copy(completado, info) // copiar elementos de info al inicio del slice

		for i := len(info); i < blockSize; i++ { // rellenar el resto del slice con ceros
			completado[i] = 0
		}
		info = completado
	}
	for i := 0; i < len(info); i += infoBits {
		// Recuperar el siguiente bloque de info
		var temp []byte
		if i+infoBits > len(info) { //verificando que los bits de info alcancen para llenar un bloque
			temp = make([]byte, infoBits)
			for i := range temp { //sino, se rellena con ceros temp y se agregan los datos que quedan en info
				temp[i] = 0
			}
			copy(temp, info[i:])
		} else {
			temp = info[i : i+infoBits]
		}
		cod := make([]byte, infoBits)
		// Vemos si tienen que generarse con error o no
		if error {
			rand.Seed(time.Now().UnixNano())
			prob := rand.Intn(100)
			if prob < 10 {
				cod = GenerarErrorEnbloque(encode(temp, parityBits, blockSize), dosErrores)
			} else {
				cod = encode(temp, parityBits, blockSize)
			}
		} else {
			// Codificar el bloque y agregarlo a la salida
			cod = encode(temp, parityBits, blockSize)
		}
		encoded = append(encoded, cod...)
	}

	return encoded
}

// Potencia de 2
func isPowerOfTwo(n int) bool {
	return (n != 0) && ((n & (n - 1)) == 0)
}
