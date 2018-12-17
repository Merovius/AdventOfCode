package main

import (
	"fmt"
	"log"
)

func main() {
	n := readTree()
	fmt.Println(n.Sum())
}

type Node struct {
	Children []Node
	Metadata []int
}

func (n Node) Sum() int {
	var total int
	for _, c := range n.Children {
		total += c.Sum()
	}
	for _, v := range n.Metadata {
		total += v
	}
	return total
}

func readTree() Node {
	var c, md int
	_, err := fmt.Scanf("%d %d", &c, &md)
	if err != nil {
		log.Fatal(err)
	}
	var n Node
	for i := 0; i < c; i++ {
		n.Children = append(n.Children, readTree())
	}
	for i := 0; i < md; i++ {
		var x int
		_, err := fmt.Scanf("%d", &x)
		if err != nil {
			log.Fatal(err)
		}
		n.Metadata = append(n.Metadata, x)
	}
	return n
}
