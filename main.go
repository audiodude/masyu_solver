package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type color uint8

const (
	EMPTY color = 0
	WHITE color = 1
	BLACK color = 2
)

type node struct {
	north *node
	south *node
	east *node
	west *node
	color color
}

func init_puzzle(f *os.File) *node {
	fmt.Println("In init_puzzle")
	scanner := bufio.NewScanner(f)

	var root_node *node
	var last_row []*node
	for scanner.Scan() {
		var (
			last_node *node
			next_node *node
		)
		spots := strings.Split(scanner.Text(), "")
		for i, spot := range spots {
			var color color
			switch spot {
			case ".":
				color = EMPTY
			case "b":
				color = BLACK
			case "w":
				color = WHITE
			}
			if root_node == nil {
				root_node = &node{color: color}
				last_node = root_node
			} else {
				next_node = &node{west: last_node, color: color}
				if last_node != nil {
					last_node.east = next_node
				}
				if len(last_row) >= i + 1 {
					last_row[i].south = next_node
					next_node.north = last_row[i]
				}
				last_node = next_node
			}
			last_row = append(last_row, last_node)
		}
	}
	return root_node
}

func main() {
	var puzzle *node
	if len(os.Args) != 2 {
		panic("First program argument should be the file name of a puzzle")
	}
	for _, file_name := range os.Args[1:] {
		fmt.Println(file_name)
		f, err := os.Open(file_name)
		if err != nil { panic(err) }

		defer func() {
			if err := f.Close(); err != nil {
				panic(err)
			}
    }()

		puzzle = init_puzzle(f)
	}
	fmt.Printf("%v\n", puzzle)
}