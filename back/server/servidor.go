package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"time"
)

var mapa map[string][]string = make(map[string][]string)

type Message struct {
	Type    int
	Msg     string
	Channel string
	Addres  string
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
	channel := make(chan Message)
	aux := host + ":3000"
	fmt.Printf("Listening on host: %s, port: %s\n", host, port)
	for {

		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("Llego una solicitud")

		go handleClient(c, aux, addr, channel)
		v := <-channel
		if v.Type == 1 {
			time.Sleep(2 * time.Second)
			go sendData(v)
		}
	}
}

//
func handleClient(c net.Conn, host, addr string, canal chan Message) {
	var receiver Message

	err := gob.NewDecoder(c).Decode(&receiver)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		canal <- receiver
		if receiver.Type == 0 {
			fmt.Println("CHANNEL:", receiver.Channel)
			fmt.Println("Local adress:", receiver.Addres)
			mapa[receiver.Channel] = append(mapa[receiver.Channel], receiver.Addres)
		}

	}

}
func sendData(msg Message) {
	value := mapa[msg.Channel]
	for i := 0; i < len(value); i++ {

		c, err := net.Dial("tcp", value[i])
		if err != nil {
			fmt.Println(err)
			return
		}
		err = gob.NewEncoder(c).Encode(msg)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Mensaje enviado con exito")
		c.Close()
	}
	return

}
func main() {
	go servidor()
	fmt.Println("Esperando clientes.....")
	var input string
	fmt.Scanln(&input)
}
