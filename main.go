package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"./server"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("api key: ")
	apiKey, _ := reader.ReadString('\n')
	apiKey = strings.TrimRight(apiKey, "\r\n")
	server.StartServer(apiKey)
}
