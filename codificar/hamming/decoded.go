package hamming

import "fmt"

func DecodeHamming(encoded []byte, blockSize int, infoSize int, error bool, parityBits int) []byte {
	decoded := make([]byte, 0)
	var decodedBlock = make([]byte, infoSize)
	for k := 0; k < len(encoded); k += blockSize {
		blockEncoded := encoded[k : k+blockSize]
		var j = 0
		//Corregir error aca
		if error {
			sindrome := checkHamming(blockEncoded, parityBits)
			if sindrome != 0 {
				blockEncoded = correctBlock(blockEncoded, sindrome)
			}
		}
		/*
			Se queda en bucle aca
		*/
		for i := 0; j < len(blockEncoded); i++ {
			if !isPowerOfTwO(i + 1) {
				decodedBlock[j] = blockEncoded[i]
				j++
			}
		}
		decoded = append(decoded, decodedBlock...)

	}
	fmt.Println("DEcoded:", decoded)
	return decoded
}

func isPowerOfTwO(n int) bool {
	return (n != 0) && ((n & (n - 1)) == 0)
}
