package main

import "fmt"

var tokens []string

func main() {
	start()
	send_info()
	// go spred()
	// go block_dc()

}

func send_info() {
	tokens = append(tokens, "hi")
	fmt.Println(tokens)
}

func start() {
}
