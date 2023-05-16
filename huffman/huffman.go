package main

import (
	"fmt"
)

func main() {
	text := "fffffffffffffffffffffffffffffffffffffffffffffeeeeeeeeeeeeeeeedddddddddddddccccccccccccbbbbbbbbbaaaaa"
	freqs := make(map[rune]int)
	for _, ch := range text {
		freqs[ch]++
	}
	raiz := huffman.ConstruirArbol(freqs)
	fmt.Println("Codigos Huffman:")
	PrintCodes(raiz, []byte{})

	fmt.Println(freqs)
	fmt.Printf("Tamaño: %d", len(text))
}
