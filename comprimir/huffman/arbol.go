package huffman

import (
	"container/heap"
	"fmt"
)

type arbol struct {
	freq int    // Frecuencia total del subárbol
	c    rune   // Símbolo del nodo (solo se usa en hojas)
	izq  *arbol // Hijo izquierdo
	der  *arbol // Hijo derecho
}

type parvaArboles []*arbol

// Inicio de funciones de interface
func (aParva parvaArboles) Len() int {
	return len(aParva)
}

func (aParva parvaArboles) Less(i, j int) bool {
	return aParva[i].freq < aParva[j].freq
}

func (aParva parvaArboles) Swap(i, j int) {
	aParva[i], aParva[j] = aParva[j], aParva[i]
}

// Agrega un árbol al monticulo de árboles
func (aParva *parvaArboles) Push(ele interface{}) {
	*aParva = append(*aParva, ele.(*arbol))
}

func (aParva *parvaArboles) Pop() interface{} {
	anterior := *aParva
	n := len(anterior)
	elemento := anterior[n-1]
	*aParva = anterior[0 : n-1]
	return elemento
}

//Fin de funciones de interface

/*
Toma un mapa de frecuencias donde las claves son símbolos y los valores son sus frecuencias en el mensaje,
y devuelve el árbol Huffman correspondiente.
*/
func ConstruirArbol(frecuencias map[rune]int) *arbol {
	var arboles parvaArboles
	for c, frec := range frecuencias {
		arboles = append(arboles, &arbol{frec, c, nil, nil})
	}
	heap.Init(&arboles)
	for arboles.Len() > 1 {
		//Toma los dos subarboles con menor frecuencia
		a := heap.Pop(&arboles).(*arbol)
		b := heap.Pop(&arboles).(*arbol)

		//Crea un nuevo nodo con la suma de las frecuencias
		newNode := &arbol{a.freq + b.freq, 0, a, b}

		//Agrega el nuevo nodo de nuevo a la parva
		heap.Push(&arboles, newNode)
	}
	// Cuando solo queda un árbol en la heap, este es el árbol Huffman, y se devuelve su raíz como resultado de la función
	return heap.Pop(&arboles).(*arbol) //Se hace así porque heap devuelve una interfaz vacia y asi se castea a tipo arbol
}

func PrintCodes(raiz *arbol, prefix []byte) {
	if raiz == nil {
		return
	}
	if raiz.c != 0 {
		fmt.Printf("%c: %s\n", raiz.c, string(prefix))
	} else {
		PrintCodes(raiz.izq, append(prefix, '0'))
		PrintCodes(raiz.der, append(prefix, '1'))
	}
}
