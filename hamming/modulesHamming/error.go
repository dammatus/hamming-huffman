package modules

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"
)

func GenerarErrorEnbloque(encoded []byte) []byte {

	// Semilla del generador de numeros aleatorios
	rand.Seed(time.Now().UnixNano())

	// Generar un número aleatorio entre 0 y blockSize -> [0,blockSize)
	num := rand.Intn(len(encoded))

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
func checkHamming(bloque []byte, parityBits int) (sindrome int64) {
	sindrome = 0
	comparador := make([]byte, len(bloque))
	// bitsBloque := getControlBits(bloque, parityBits)
	for i := 0; i < len(bloque); i++ {
		if isPowerOfTwo(i + 1) {
			comparador[i] = 0
		} else {
			comparador[i] = bloque[i]
		}
	}
	fmt.Println("********************-------------*****************")
	fmt.Println(comparador)
	fmt.Println(bloque)
	fmt.Println("********************-------------*****************")
	//vuelvo a aplicar los bits de control pero en un comparador,
	//para así despues comparar que ambos sean iguales
	for i := 0; i < parityBits; i++ {
		comparador = aplicaBitsDeControl(i, comparador, len(bloque))
	}

	// bitsComparador := getControlBits(comparador, parityBits)

	// j := 0
	// for i := parityBits - 1; i >= 0; i-- {
	// 	if bitsBloque[i] != bitsComparador[i] {
	// 		sindrome += int(math.Pow(2, float64(j)))
	// 	}

	// 	j++
	// }
	var control byte
	var aux []byte
	for i := 0; i < parityBits; i++ {

		for j := int(math.Pow(2, float64(i))); j < len(bloque); j += int(math.Pow(2, float64(i))) {
			control = bloque[j] ^ comparador[j]
			aux = append(aux, control)
		}
	}
	var sindStr string
	sindStr = string(aux)

	sindrome, _ = strconv.ParseInt(sindStr, 2, 64)

	return
}

/*
correctBlock:
Recibe como argumento un bloque y su sindrome (conseguido de checkHamming)
Devuelve el bloque corregido
ALERTA: solo se debe llamar correctBlock si hubo algún error, es decir si sindrome != 0
*/
func correctBlock(bloque []byte, sindrome int64) []byte {
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
	}
	return retorno
}
