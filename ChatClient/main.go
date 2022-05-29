package main

import (
	"flag"
	"fmt"
	"zChatRoom/ChatClient/cui"
	"zChatRoom/ChatClient/handler"
	"zChatRoom/ChatClient/model"

	"github.com/pzqf/zEngine/zNet"
)

func main() {
	address := flag.String("a", "127.0.0.1", "server address")
	flag.Parse()
	handler.Init()
	var cli = zNet.TcpClient{}
	fmt.Println("connect to server", *address)
	err := cli.ConnectToServer(*address, 9106)
	if err != nil {
		fmt.Printf("Connect:, err:%s \n", err.Error())
		return
	}
	model.Init(&cli)

	defer cli.Close()
	fmt.Println("Connect success :")

	cui.InitUi()

	cui.ShowLoginUi()

	cui.StartUi()
	cui.StopUi()
	cli.Close()
}
