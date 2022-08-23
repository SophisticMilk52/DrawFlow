/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

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
var clients []Receiver
var senders []Sender
var mutex sync.Mutex
var wg sync.WaitGroup

func server() {
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
	go runApi()
	fmt.Printf("Listening on host: %s, port: %s\n", host, port)

	for {
		wg.Add(3)
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleClient(c, channel, &mutex, &wg)

		v := <-channel
		if v.Type == 1 {
			t := Sender{
				Addres:  c.LocalAddr().String(),
				Channel: v.Channel,
				Size:    len([]byte(v.Msg)),
				Time:    time.Now().Format(time.ANSIC),
			}
			go sendData(v, t, &mutex, &wg)

		} else if v.Type == 3 {
			go checkConection(v.Addres, &mutex, &wg)
		}

	}

}

func checkConection(addr string, m *sync.Mutex, wg *sync.WaitGroup) {

	for {
		for i := 0; i < len(clients); i++ {
			if clients[i].Addres == addr {
				m.Lock()
				clients = append(clients[:i], clients[i+1:]...)
				fmt.Printf("Logout client %s", addr)
				fmt.Println("")
				m.Unlock()
				wg.Done()
				continue
			}
		}
		break
	}
}

func handleClient(c net.Conn, canal chan Message, m *sync.Mutex, wg *sync.WaitGroup) {

	var receiver Message

	err := gob.NewDecoder(c).Decode(&receiver)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		canal <- receiver
		if receiver.Type == 0 {

			fmt.Printf("New request arrived in channel %s and adress %s", receiver.Channel, receiver.Addres)
			fmt.Println("")
			mapa[receiver.Channel] = append(mapa[receiver.Channel], receiver.Addres)

			t := Receiver{
				Addres:  receiver.Addres,
				Channel: receiver.Channel,
				Time:    time.Now().Format(time.ANSIC),
			}
			m.Lock()
			clients = append(clients, t)
			m.Unlock()
			wg.Done()
		} else {
			c.Close()
		}

	}

}

func sendData(msg Message, t Sender, m *sync.Mutex, wg *sync.WaitGroup) error {

	value := mapa[msg.Channel]
	for i := 0; i < len(value); i++ {
		c, err := net.Dial("tcp", value[i])
		if err != nil {
			fmt.Println(err)
			return err
		}

		err = gob.NewEncoder(c).Encode(msg)
		if err != nil {
			fmt.Println(err)
			return err
		}
		m.Lock()
		senders = append(senders, t)
		fmt.Println("Message successfully send")
		m.Unlock()
		wg.Done()
		c.Close()
	}
	return nil
}
func listen() {
	defer wg.Wait()
	go server()
	fmt.Println("Waiting for clients.....")
	var input string
	fmt.Scanln(&input)
}

type Receiver struct {
	Addres  string `json:"addres"`
	Channel string `json:"channel"`
	Time    string `json:"timestamp"`
}
type Sender struct {
	Addres  string `json:"addres"`
	Channel string `json:"channel"`
	Size    int    `json:"size"`
	Time    string `json:"timestamp"`
}

func runApi() {
	mux := http.NewServeMux()
	mux.HandleFunc("/subscribers", getSubscribers)
	mux.HandleFunc("/sender", getSenders)
	server := &http.Server{
		Addr:    "localhost:9090",
		Handler: mux,
	}
	log.Fatal(server.ListenAndServe())
}

func getSenders(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(senders)
}

func getSubscribers(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(clients)
}
