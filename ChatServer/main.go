package main

import (
	"log"
	_ "net/http/pprof"
	"zChatRoom/ChatServer/handler"
	"zChatRoom/ChatServer/playerMgr"
	"zChatRoom/ChatServer/room"

	"github.com/pzqf/zEngine/zNet"
	"github.com/pzqf/zEngine/zSignal"
	"github.com/pzqf/zUtil/zKeyWordFilter"
)

func main() {
	zKeyWordFilter.InitDefaultFilter()
	err := zKeyWordFilter.ParseFromFile(`keyword.txt`)
	if err != nil {
		log.Println("KeyWordFilter.ParseFromFile error ", err)
		return
	}

	playerMgr.InitDefaultPlayerMgr()

	room.InitDefaultRoomMgr(2)

	zNet.InitDefaultTcpServer(":9106", 10000)
	zNet.InitPacket(zNet.PacketCodeJson, zNet.MaxNetPacketDataSize)
	zNet.GetDefaultTcpServer().SetRemoveSessionCallBack(playerMgr.OnSessionClose)

	err = handler.Init()
	if err != nil {
		log.Printf("RegisterHandler error %d", 1)
		return
	}
	err = zNet.GetDefaultTcpServer().Start()
	if err != nil {
		log.Printf(err.Error())
		return
	}
	log.Printf("Tcp server listing on  %d", 9106)

	zSignal.GracefulExit()
	log.Printf("server will be shut off")
	zNet.GetDefaultTcpServer().Close()

	log.Printf("====>>> FBI warning , server exit <<<=====")
}
