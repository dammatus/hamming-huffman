package hamming

import (
	"math/rand"
	"time"
)

func GenerarErrorEnbloque(encoded []byte, dosErrores bool) []byte {

	// Semilla del generador de numeros aleatorios
	rand.Seed(time.Now().UnixNano())

	// Generar un número aleatorio entre 0 y blockSize -> [0,blockSize)
	num := rand.Intn(len(encoded) - 1)
	if dosErrores {
		// Generar un número aleatorio entre 0 y blockSize -> [0,blockSize)
		e := rand.Intn(len(encoded) - 1)
		encoded[e] = encoded[e] ^ 1
	}

	// Si en num hay un 0 -> lo cambia a 1. Si hay un 1 -> Lo cambia a 0
	encoded[num] = encoded[num] ^ 1

	return encoded
}

/*
checkHamming:
Recibe como argumento:
  - Un bloque de información (el de 32 bits que ya tiene aplicado el hamming)
  - La cantidad de bits de control (Es probable que este valor se pueda calcular solo con el len de bloque)

Devuelve el sindrome de ese bloque (Creo que es util que devuelva directamente el entero en base 10,
habria que restarle 1 para que se corresponda con la posición del arreglo donde esta el bit a cambiar)
*/
func checkHamming(bloque []byte, parityBits int) (sindrome int) {
	sindrome = 0
	sindromeBits := make([]byte, parityBits)
	comparador := make([]byte, len(bloque))
	bitsBloque := getControlBits(bloque, parityBits)
	for i := 0; i < len(bloque); i++ {
		if !isPowerOfTwo(i + 1) {
			comparador[i] = bloque[i]
		}
	}

	//vuelvo a aplicar los bits de control pero en un comparador,
	//para así despues comparar que ambos sean iguales
	for i := 0; i < parityBits; i++ {
		comparador = aplicaBitsDeControl(i, comparador, len(bloque))
	}

	bitsComparador := getControlBits(comparador, parityBits)

	j := 0
	for i := parityBits - 1; i >= 0; i-- {
		if bitsBloque[i] != bitsComparador[i] {
			sindromeBits[j] = 1
		} else {
			sindromeBits[j] = 0
		}

		j++
	}

	sindrome = binaryToDecimal(sindromeBits)

	return
}

/*
correctBlock:
Recibe como argumento un bloque y su sindrome (conseguido de checkHamming)
Devuelve el bloque corregido
ALERTA: solo se debe llamar correctBlock si hubo algún error, es decir si sindrome != 0
*/
func correctBlock(bloque []byte, sindrome int) []byte {
	bloque[sindrome-1] ^= 1 //Si no anda puede que sea porque no se permite el xor con 1 directamente
	return bloque
}

/*
getControlBits:
Función que toma como argumento un bloque de bits y retorna sus bits de control
Fue creada para usarse en checkHamming
*/
func getControlBits(bloque []byte, parityBits int) []byte {
	retorno := make([]byte, parityBits)
	j := 0
	for i := 1; i < len(bloque)-1; i = i * 2 {
		retorno[j] = bloque[i-1]
		j++
	}
	return retorno
}
