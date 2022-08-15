/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

var channel string

// receiveCmd represents the receive command
var receiveCmd = &cobra.Command{
	Use:   "receive <channel>",
	Short: "Program a client to receive documents",
	Long:  `Specify which channel do you want your client to run`,

	Args: func(_ *cobra.Command, args []string) error {
		if channel == "" && len(args) < 1 {
			return errors.New("accepts 1 arg(s)")
		}
		return nil
	},
	Run: func(_ *cobra.Command, _ []string) {
		addr := "localhost:"
		s, err := net.Listen("tcp", addr)
		if err != nil {
			log.Fatalf("Could establish TCP connection: %s", err.Error())
		}
		initSend(s.Addr().String())
		go checkOut(s.Addr().String())
		initReceive(s)
	},
}

func clientSend(address string, t int) {
	addr := "localhost:8888"
	message := Message{
		Type:    t,
		Channel: channel,
		Msg:     "",
		Addres:  address,
	}
	c, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = gob.NewEncoder(c).Encode(message)
	fmt.Println("Sucessfull connection!")
	if err != nil {
		fmt.Println(err)
	}
	c.Close()
}

func clientReceive(s net.Listener) {
	for {
		c, err := s.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleClientReceive(c)
	}
}

func checkOut(addr string) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Close session")
		clientSend(addr, 3)
		os.Exit(1)
	}()
}

func handleClientReceive(c net.Conn) {
	var receiver Message
	err := gob.NewDecoder(c).Decode(&receiver)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("Mensaje:", receiver.Msg)
	}
	c.Close()
}

func init() {
	rootCmd.AddCommand(receiveCmd)
	// TODO: Revisar para que sirve ese parametro en blanco
	receiveCmd.Flags().StringVarP(&channel, "channel", "c", "", "Especify the channel of the client")
	receiveCmd.MarkFlagRequired("channel")
}

func initSend(addr string) {
	go clientSend(addr, 0)
}

func initReceive(s net.Listener) {
	clientReceive(s)
}

type Message struct {
	Type    int
	Msg     string
	Channel string
	Addres  string
}
