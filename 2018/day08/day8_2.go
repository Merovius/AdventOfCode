package main

import (
	"fmt"
	"log"
)

func main() {
	n := readTree()
	fmt.Println(n.Value())
}

type Node struct {
	Children []Node
	Metadata []int
}

func (n Node) Value() int {
	if len(n.Children) == 0 {
		val := 0
		for _, v := range n.Metadata {
			val += v
		}
		return val
	}
	val := 0
	for _, md := range n.Metadata {
		md -= 1
		if md < 0 || md >= len(n.Children) {
			continue
		}
		val += n.Children[md].Value()
	}
	return val
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
