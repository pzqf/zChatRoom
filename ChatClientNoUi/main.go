package main

import (
	"encoding/json"
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
	address := flag.String("a", "192.168.50.16", "server address")
	count := flag.Int("n", 1000, "client count")
	begin := flag.Int("b", 1, "begin count")
	flag.Parse()
	handler.Init()
	clientCount := *count
	zNet.InitPacket(zNet.DefaultPacketDataSize)

	for i := *begin; i < *begin+clientCount; i++ {
		go func(n int) {
			var cli = zNet.TcpClient{}
			fmt.Println("connect to server", *address)
			err := cli.ConnectToServer(*address, 9160, "", 30)
			if err != nil {
				fmt.Printf("Connect:, err:%s \n", err.Error())
				return
			}

			data, err := json.Marshal(proto.PlayerLoginReq{
				UserName: fmt.Sprintf("player_%d", n),
			})
			if err != nil {
				return
			}
			err = cli.Send(proto.PlayerLogin, data)
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
