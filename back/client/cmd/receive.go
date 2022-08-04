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
		initReceive()
		fmt.Println("--------------")
	},
}

func ClientReceive() {
	s, err := net.Listen("tcp", ":9999")
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
	receiver := Message{
		Type:    0,
		Channel: channel,
		Msg:     "",
	}

	err := gob.NewDecoder(c).Decode(&receiver)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("Mensaje:", receiver.Msg)
	}

}

func init() {
	rootCmd.AddCommand(receiveCmd)
	receiveCmd.Flags().StringVarP(&channel, "channel", "c", "", "Especify the channel of the client")
	receiveCmd.MarkFlagRequired("channel")
}

func initReceive() {
	go ClientReceive()
	var input string
	fmt.Scanln(&input)
}

type Message struct {
	Type    int
	Msg     string
	Channel string
}
