/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var File string
var Channel string

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send <path> <ch>",
	Short: "Client that sends files",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("File", File)
		fmt.Println("Channel", Channel)
		var filename string
		var err error
		var argument string

		if File != "" {
			argument = File
		} else {
			argument = args[0]
		}
		fileExists, err := FileExists(argument)
		if err != nil {
			fmt.Println(err)
		}
		if fileExists {
			filename, err = filepath.Abs(argument)
			if err != nil {
				fmt.Println(err.Error())

			}
		} else {
			fmt.Printf("File %v doest not Exists", argument)
			return
		}

		msg := Message{
			Type:    1,
			Msg:     string(fileReader(filename)),
			Channel: Channel,
		}
		sender(filename, msg)

		if err != nil {
			log.Fatal("VS Code executable file not found in %PATH%")
		}
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)
	sendCmd.Flags().StringVarP(&File, "file", "f", "", "Especify the route of the file")
	sendCmd.Flags().StringVarP(&Channel, "channel", "c", "", "Especify the channel of the client")
	sendCmd.MarkFlagRequired("file")
	sendCmd.MarkFlagRequired("channel")
}
func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func fileReader(path string) []byte {
	Myfile, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file!!!")
	}

	byteBuff := make([]byte, 1073741824)
	totalLen, err := Myfile.Read(byteBuff)
	if err != nil {
		fmt.Println(err)
	}
	bytes := byteBuff[:totalLen]
	fmt.Printf("File Data: \n%s\n", string(bytes))

	Myfile.Close()
	return bytes
}
func sender(path string, message Message) {
	c, err := net.Dial("tcp", ":8888")
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
