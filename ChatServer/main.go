package main

import (
	"fmt"
	"log"
	_ "net/http/pprof"
	"zChatRoom/ChatServer/handler"
	"zChatRoom/ChatServer/playerMgr"
	"zChatRoom/ChatServer/room"

	"github.com/pkg/profile"

	"go.uber.org/zap"

	"github.com/pzqf/zEngine/zLog"
	"github.com/pzqf/zEngine/zNet"
	"github.com/pzqf/zEngine/zSignal"
	"github.com/pzqf/zUtil/zKeyWordFilter"
)

func main() {
	stopper := profile.Start(profile.CPUProfile, profile.ProfilePath("."), profile.NoShutdownHook)
	defer stopper.Stop()
	// go tool pprof -http=:9999 cpu.pprof

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

	maxRoomCount := int32(50)
	room.InitDefaultRoomMgr(maxRoomCount)
	zLog.Info("room manager init success", zap.Int32("maxRoomCount", maxRoomCount))

	err = handler.Init()
	if err != nil {
		zLog.Error("RegisterHandler error %d", zap.Error(err))
		return
	}
	zLog.Info("handler init success")

	netCfg := zNet.Config{
		MaxPacketDataSize: zNet.DefaultPacketDataSize,
		ListenAddress:     fmt.Sprintf(":%d", 9160),
	}
	zNet.InitTcpServerDefault(&netCfg,
		zNet.WithMaxClientCount(10000),
		zNet.WithSidInitio(10000),
		zNet.WithHeartbeat(30),
	)
	zNet.GetTcpServerDefault().SetRemoveSessionCallBack(playerMgr.OnSessionClose)
	zLog.Info("tcp server init success", zap.Int("maxClientCount", 10000),
		zap.Int32("MaxNetPacketDataSize", zNet.DefaultPacketDataSize),
	)

	err = zNet.GetTcpServerDefault().Start()
	if err != nil {
		zLog.Error("start tcp server error", zap.Error(err))
		return
	}

	zSignal.GracefulExit()
	zLog.Info("server will be shutdown")
	zNet.GetTcpServerDefault().Close()

	log.Printf("====>>> FBI warning , server exit <<<=====")
}
