package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

import (
	"net/http"
	"html/template"
)

type color uint8
type line uint8
const (
	EMPTY color = 0
	WHITE color = 1
	BLACK color = 2
	VERT line = 3
	HORZ line = 4
	NE line = 5
	SE line = 6
	SW line = 7
	NW line = 8
)

type conn struct {
	node *node
	path bool
}

type node struct {
	north *conn
	south *conn
	east *conn
	west *conn
	color color
	lines []line
}

func init_puzzle(f *os.File) *node {
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
				next_node = &node{west: &conn{node: last_node}, color: color}
				if last_node != nil {
					last_node.east = &conn{node: next_node}
				}
				if len(last_row) >= i + 1 {
					last_row[i].south = &conn{node: next_node}
					next_node.north = &conn{node: last_row[i]}
				}
				last_node = next_node
			}
			if len(last_row) >= i + 1 {
				last_row[i] = last_node
			} else {
				last_row = append(last_row, last_node)
			}
		}
	}
	return root_node
}

func print_puzzle(root_node *node) {
	if root_node == nil {
		return
	}
	node := root_node
	for node.south != nil && node.east != nil {
		row_start := node
		for node != nil {
			switch node.color {
			case EMPTY:
				fmt.Print(".")
			case WHITE:
				fmt.Print("○")
			case BLACK:
				fmt.Print("●")
			}

			if node.east != nil {
				node = node.east.node
			} else {
				node = nil
			}
		}
		fmt.Print("\n")
		node = row_start.south.node
	}
}

func fmt_puzzle_html(root_node *node) string {
	if root_node == nil {
		return ""
	}
	var buffer bytes.Buffer
	buffer.WriteString("<div class=\"puzzle\">")

	node := root_node
	for node.south != nil && node.east != nil {
		row_start := node
		buffer.WriteString("\n  <div class=\"row\">")
		for node != nil {
			switch node.color {
			case EMPTY:
				buffer.WriteString("\n    <div class=\"square\"><div class=\"em\"></div></div>")
			case WHITE:
				buffer.WriteString("\n    <div class=\"square\"><div class=\"wh\"></div></div>")
			case BLACK:
				buffer.WriteString("\n    <div class=\"square\"><div class=\"bl\"></div></div>")
			}

			if node.east != nil {
				node = node.east.node
			} else {
				node = nil
			}
		}
		buffer.WriteString("\n  </div>")
		node = row_start.south.node
	}

	buffer.WriteString("\n</div>")
	return buffer.String()
}

func main() {
	var puzzle *node
	if len(os.Args) != 4 {
		panic("usage: masyu_solver puzzle.txt template_dir/ static_dir/")
	}

	file_name := os.Args[1]
	f, err := os.Open(file_name)
	if err != nil { panic(err) }

	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()
	puzzle = init_puzzle(f)

	tmpl_path := filepath.Join(os.Args[2], "index.go.html")
	tmpl := template.Must(template.ParseFiles(tmpl_path))

	print_puzzle(puzzle)
	puzzle_html := fmt_puzzle_html(puzzle)
	var data = map[string]template.HTML {
		"puzzle": template.HTML(puzzle_html),
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(
		http.Dir("./" + os.Args[3]))))
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	  tmpl.Execute(w, data)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
