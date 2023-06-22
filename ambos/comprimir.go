package ambos

import (
	"fmt"
	"hamming-huffman/comprimir/huffman"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

func Comprimir(w http.ResponseWriter, contenido []byte) {

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

	err := huffman.SaveCompactedAmbos(compacted, raiz)
	if err != nil {
		fmt.Println("Error al guardar el archivo: ", err)
	} else {
		fmt.Println("Datos comprimidos exitosamente!")
	}

}

func Descomprimir(w http.ResponseWriter, compacted string) {
	unziped, raizRecuperada, error := huffman.GetFromUnzip()

	if error == nil {
		fmt.Println("Recuperados del archivo: ", unziped)
	}

	fmt.Println("Resultado: ", huffman.DecodeData(raizRecuperada, unziped))

	//Este es el descomprimido que se mostrara en la pagina
	unzip := huffman.DecodeData(raizRecuperada, unziped)

	if err := ioutil.WriteFile(filepath.Join("ambos/files", "descomprimido.txt"), []byte(unzip), 0644); err != nil {
		http.Error(w, "No se pudo guardar el archivo descomprimido.txt", http.StatusInternalServerError)
		return
	}
	//Este es el descomprimido huf
	if err := ioutil.WriteFile(filepath.Join("ambos/resultados", "descomprimido.dhu"), []byte(unzip), 0644); err != nil {
		http.Error(w, "No se pudo guardar el archivo descomprimido.dhu", http.StatusInternalServerError)
		return
	}
}
