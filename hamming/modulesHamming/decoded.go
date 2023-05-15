package modules

import (
	"fmt"
)

func DecodeHamming(encoded []byte, blockSize int, infoSize int, error bool, parityBits int) []byte {
	decoded := make([]byte, 0)
	var decodedBlock = make([]byte, infoSize)
	for k := 0; k < len(encoded); k += blockSize {
		blockEncoded := encoded[k : k+blockSize]
		var j = 0
		//Corregir error aca
		if error {
			sindrome := checkHamming(blockEncoded, parityBits)
			fmt.Println(sindrome)
			// // Convertir el slice de bytes a una cadena
			// bitStr := string(sindrome)
			// fmt.Println(bitStr)
			// // Convertir la cadena de bits a un entero
			// num, _ := strconv.ParseInt(bitStr, 2, 0)
			// fmt.Println(num)
			// // Convertir el entero a un tipo int
			// pos := int(num)
			// fmt.Println(pos)
			fmt.Println("********************************************")
			if sindrome != 0 {
				blockEncoded = correctBlock(blockEncoded, sindrome)
			}
		}

		for i := 0; i < len(blockEncoded); i++ {
			if !isPowerOfTwO(i + 1) {
				decodedBlock[j] = blockEncoded[i]
				j++
			}
		}
		decoded = append(decoded, decodedBlock...)
	}
	return decoded
}

func isPowerOfTwO(n int) bool {
	return (n != 0) && ((n & (n - 1)) == 0)
}
