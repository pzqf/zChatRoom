package main

import (
	"fmt"
	"log"
	_ "net/http/pprof"
	"zChatRoom/ChatServer/handler"
	"zChatRoom/ChatServer/playerMgr"
	"zChatRoom/ChatServer/room"

	"go.uber.org/zap"

	"github.com/pzqf/zEngine/zLog"
	"github.com/pzqf/zEngine/zNet"
	"github.com/pzqf/zEngine/zSignal"
	"github.com/pzqf/zUtil/zKeyWordFilter"
)

func main() {
	cfg := zLog.Config{
		Level:    zLog.InfoLevel,
		Console:  true,
		Filename: "./logs/server.log",
		MaxSize:  128,
	}
	err := zLog.InitLogger(&cfg)
	if err != nil {
		fmt.Println(err)
		return
	}

	zLog.Info("server start....")

	zKeyWordFilter.InitDefaultFilter()
	err = zKeyWordFilter.ParseFromFile(`keyword.txt`)
	if err != nil {
		zLog.Error("KeyWordFilter.ParseFromFile error ", zap.Error(err))
		return
	}

	zLog.Info("keyword filter init success")

	playerMgr.InitDefaultPlayerMgr()
	zLog.Info("player manager init success")

	maxRoomCount := int32(10)
	room.InitDefaultRoomMgr(maxRoomCount)
	zLog.Info("room manager init success", zap.Int32("maxRoomCount", maxRoomCount))

	zNet.InitDefaultTcpServer(fmt.Sprintf(":%d", 9106),
		zNet.WithMaxClientCount(10000),
		zNet.WithSidInitio(10000),
		zNet.WithPacketCodeType(zNet.PacketCodeJson),
		zNet.WithMaxPacketDataSize(zNet.MaxNetPacketDataSize),
		zNet.WithDispatcherPoolSize(3000),
	)
	zNet.GetDefaultTcpServer().SetRemoveSessionCallBack(playerMgr.OnSessionClose)
	zLog.Info("tcp server init success", zap.Int("maxClientCount", 10000),
		zap.Int32("MaxNetPacketDataSize", zNet.MaxNetPacketDataSize),
		zap.Int("PacketCodeType", int(zNet.PacketCodeJson)),
	)

	err = handler.Init()
	if err != nil {
		zLog.Error("RegisterHandler error %d", zap.Error(err))
		return
	}
	zLog.Info("handler init success")

	err = zNet.GetDefaultTcpServer().Start()
	if err != nil {
		zLog.Error("start tcp server error", zap.Error(err))
		return
	}
	zLog.Info("Tcp server listing on", zap.Int("port", 9106))

	zSignal.GracefulExit()
	zLog.Info("server will be shutdown")
	zNet.GetDefaultTcpServer().Close()

	log.Printf("====>>> FBI warning , server exit <<<=====")
}
