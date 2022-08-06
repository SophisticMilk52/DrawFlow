/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/gob"
	"errors"
	"fmt"
	"net"

	"github.com/spf13/cobra"
)

var channel string

// receiveCmd represents the receive command
var receiveCmd = &cobra.Command{
	Use:   "receive <channel>",
	Short: "Program a client to receive documents",
	Long:  `Especify witch channel do you want your client to run`,

	Args: func(cmd *cobra.Command, args []string) error {
		if channel == "" && len(args) < 1 {
			return errors.New("accepts 1 arg(s)")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		addr := "localhost:"
		s, err := net.Listen("tcp", addr)
		initSend(s.Addr().String())
		initReceive(s, err)
	},
}

func ClientSend(address string) {
	addr := "localhost:8888"
	message := Message{
		Type:    0,
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
	if err != nil {
		fmt.Println(err)
	}
	c.Close()
}
func ClientReceive(s net.Listener, err error) {
	if err != nil {
		fmt.Print(err)
		return
	}
	for {
		c, err := s.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleClientReceive(c)
	}
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
	receiveCmd.Flags().StringVarP(&channel, "channel", "c", "", "Especify the channel of the client")
	receiveCmd.MarkFlagRequired("channel")
}
func initSend(addr string) {
	go ClientSend(addr)
}
func initReceive(s net.Listener, err error) {
	ClientReceive(s, err)
}

type Message struct {
	Type    int
	Msg     string
	Channel string
	Addres  string
}
