package main

import (
	"encoding/json"
	"fmt"
	"github.com/cmu440/bitcoin"
	"github.com/cmu440/lsp"
	"os"
	"strconv"
)

func main() {
	const numArgs = 4
	if len(os.Args) != numArgs {
		fmt.Printf("Usage: ./%s <hostport> <message> <maxNonce>", os.Args[0])
		return
	}
	hostport := os.Args[1]
	message := os.Args[2]
	maxNonce, err := strconv.ParseUint(os.Args[3], 10, 64)
	if err != nil {
		fmt.Printf("%s is not a number.\n", os.Args[3])
		return
	}

	client, err := lsp.NewClient(hostport, lsp.NewParams())
	if err != nil {
		//fmt.Println("Failed to connect to server:", err)
		printDisconnected()
		return
	}

	//defer client.Close()

	// generate a new request
	request := bitcoin.NewRequest(message, 0, maxNonce)
	rawMsg, _ := json.Marshal(request)

	//if err != nil {
	//fmt.Println("Failed to marshall the request")
	//fmt.Println(err.Error())
	//	return
	//}

	err = client.Write(rawMsg)
	if err != nil {
		printDisconnected()
		return
	}
	data, readerr := client.Read()
	if readerr != nil {
		printDisconnected()
		return
	}

	result := new(bitcoin.Message)
	err = json.Unmarshal(data, result)
	if err != nil {
		//fmt.Println("Failed to marshall the request")
		//fmt.Println(err.Error())
		return
	}

	if result.Type == bitcoin.Result {
		printResult(result.Hash, result.Nonce)
	}

}

// printResult prints the final result to stdout.
func printResult(hash, nonce uint64) {
	fmt.Println("Result", hash, nonce)
}

// printDisconnected prints a disconnected message to stdout.
func printDisconnected() {
	fmt.Println("Disconnected")
}
