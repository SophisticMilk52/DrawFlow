/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/gob"
	"fmt"
	"net"
	"time"

	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "The server that manages the clients",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		listen()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

var mapa map[string][]string = make(map[string][]string)

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
func listen() {
	go servidor()
	fmt.Println("Esperando clientes.....")
	var input string
	fmt.Scanln(&input)
}
