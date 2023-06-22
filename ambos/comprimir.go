package ambos

import (
	"fmt"
	"hamming-huffman/comprimir/huffman"
	"net/http"
)

func Comprimir(w http.ResponseWriter, contenido []byte) string {

	text := string(contenido)

	fmt.Println(text)

	freqs := make(map[rune]int)
	for _, ch := range text {
		freqs[ch]++
	}
	raiz := huffman.ConstruirArbol(freqs)

	fmt.Println("Codigos Huffman:")
	huffman.PrintCodes(raiz, []byte{})

	fmt.Println(freqs)
	fmt.Printf("Tamaño: %d\n", len(text))

	compacted := huffman.Compacted(text, raiz)

	//fmt.Println("Compactado: " + binary)

	fmt.Println("Codigo Huffman: ", compacted)

	err := huffman.SaveCompacted(compacted, raiz)
	if err != nil {
		fmt.Println("Error al guardar el archivo: ", err)
	} else {
		fmt.Println("Datos comprimidos exitosamente!")
	}

	return compacted
}
