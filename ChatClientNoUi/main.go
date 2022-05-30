package main

import (
	"flag"
	"fmt"
	"time"
	"zChatRoom/ChatClientNoUi/handler"
	"zChatRoom/proto"

	"github.com/pzqf/zEngine/zSignal"

	//"zChatRoom/ChatClientNoUi/model"

	"github.com/pzqf/zEngine/zNet"
)

func main() {
	address := flag.String("a", "127.0.0.1", "server address")
	count := flag.Int("n", 100, "client count")
	flag.Parse()
	handler.Init()
	clientCount := *count
	zNet.InitPacket(zNet.PacketCodeJson, zNet.MaxNetPacketDataSize)

	for i := 0; i < clientCount; i++ {
		go func(n int) {
			var cli = zNet.TcpClient{}
			fmt.Println("connect to server", *address)
			err := cli.ConnectToServer(*address, 9106)
			if err != nil {
				fmt.Printf("Connect:, err:%s \n", err.Error())
				return
			}

			err = cli.Send(proto.PlayerLogin, &proto.PlayerLoginReq{
				UserName: fmt.Sprintf("player_%d", n),
			})
			if err != nil {
				fmt.Println(err)
				return
			}

			defer cli.Close()
			fmt.Println("Connect success :", n)
			select {}
		}(i)

		time.Sleep(time.Millisecond * 100)
	}
	zSignal.GracefulExit()
}
