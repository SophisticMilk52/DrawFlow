package cmd

import (
	"net"
	"testing"
)

var path string = "../files/hola.txt"

func TestFileExists(t *testing.T) {

	_, err := fileExists("s")
	if err == nil {
		t.Log(err, nil)
	}
	val, _ := fileExists(path)
	if val == true {
		t.Logf("The path %s is a valid %t", path, val)
	}

}

func TestFileReader(t *testing.T) {
	bytes, _ := fileReader(path)
	if bytes != nil {
		t.Logf("The amount of bytes is %d", len(bytes))
	}
	_, err := fileReader("s")
	if err == nil {
		t.Errorf("path not found %s", err)
	}

}

func TestSender(t *testing.T) {
	//	servidor()
	addr := "localhost:8888"
	l, _ := net.Listen("tcp", addr)
	msg := Message{
		Type:    1,
		Msg:     "HI",
		Channel: Channel,
	}
	err := sender(msg)
	if err == nil {
		t.Log(err)
	}
	l.Close()
}
