package main

import (
	"encoding/gob"
	"fmt"
	"net"
)

var mapa map[string]string = make(map[string]string)

type Message struct {
	Type    int
	Msg     string
	Channel string
}

func servidor() {
	addr := "localhost:8888"
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	host, port, err := net.SplitHostPort(l.Addr().String())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Listening on host: %s, port: %s\n", host, port)
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleClient(c, host, addr)
	}
}

//
func handleClient(c net.Conn, host, addr string) {
	var receiver Message

	err := gob.NewDecoder(c).Decode(&receiver)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		if receiver.Type == 0 {

			mapa[receiver.Channel] = host
		} else {
			fmt.Println("Mensaje:", receiver.Msg)
			fmt.Println("Canal", receiver.Channel)
			sendData(addr, receiver)
		}

	}

}
func sendData(addr string, msg Message) {
	c, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = gob.NewEncoder(c).Encode(msg)
	if err != nil {
		fmt.Println(err)
	}
	c.Close()

}
func main() {
	go servidor()
	fmt.Println("Esperando clientes.....")
	var input string
	fmt.Scanln(&input)
}
