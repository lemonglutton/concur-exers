package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	initializeClient()
}

func initializeClient() {
	conn, err := net.Dial("tcp", HOST+":"+PORT)
	if err != nil {
		log.Fatalf("Dial failed: %v", err.Error())
	}
	defer conn.Close()

	serverReader := bufio.NewReader(conn)
	inputReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Provide me with Your username. Username needs to have at least 8 characters\n")
		fmt.Print("-> ")
		input, err := inputReader.ReadString('\n')
		if err != nil {
			log.Printf("Stdin reading err: %v", err.Error())
		}
		_, err = conn.Write([]byte(input))
		if err != nil {
			log.Printf("Message sending error: %v", err.Error())
		}

		response, err := serverReader.ReadString('\n')
		if err != nil {
			log.Printf("Stdin reading err: %v", err.Error())
		}
		log.Printf("%v\n", response)

		if response != errToShort && response != errDuplicated {
			return
		}
	}

}
