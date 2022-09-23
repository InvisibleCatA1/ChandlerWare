package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	fmt.Println("Starting...")
	fmt.Println("Loading files...")

	content, err := ioutil.ReadFile("stub/stub.go")
	if err != nil {
		fmt.Println(err)
	}

	_, err2 := os.Stat("tmp/new_stub.go")
	if err2 != nil {
		err3 := ioutil.WriteFile("/tmp/new_stub.go", content, 0644)
		if err3 != nil {
			fmt.Println(err)
		}
	}

}
