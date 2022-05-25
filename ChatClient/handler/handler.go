package handler

import (
	"fmt"
	"log"
	"time"
	"zChatRoom/ChatClient/cui"
	"zChatRoom/proto"

	"github.com/jroimartin/gocui"

	"github.com/pzqf/zUtil/zTime"

	"github.com/pzqf/zEngine/zNet"
)

func Init() {
	if err := zNet.RegisterHandler(proto.PlayerLogin, PlayerLoginRes); err != nil {
		log.Printf("RegisterHandler error %d", err)
		return
	}
	if err := zNet.RegisterHandler(proto.PlayerLogout, PlayerLogoutRes); err != nil {
		log.Printf("RegisterHandler error %d", err)
		return
	}
	if err := zNet.RegisterHandler(proto.PlayerEnterRoom, PlayerEnterRoomRes); err != nil {
		log.Printf("RegisterHandler error %d", err)
		return
	}
	if err := zNet.RegisterHandler(proto.PlayerLeaveRoom, PlayerLeaveRoomRes); err != nil {
		log.Printf("RegisterHandler error %d", err)
		return
	}
	if err := zNet.RegisterHandler(proto.PlayerSpeak, PlayerSpeakRes); err != nil {
		log.Printf("RegisterHandler error %d", err)
		return
	}
	if err := zNet.RegisterHandler(proto.SpeakBroadcast, SpeakBroadcast); err != nil {
		log.Printf("RegisterHandler error %d", err)
		return
	}
	if err := zNet.RegisterHandler(proto.RoomList, RoomListRes); err != nil {
		log.Printf("RegisterHandler error %d", err)
		return
	}
	if err := zNet.RegisterHandler(proto.RoomPlayerList, RoomPlayerListRes); err != nil {
		log.Printf("RegisterHandler error %d", err)
		return
	}
}

func PlayerLoginRes(session *zNet.Session, packet *zNet.NetPacket) {
	var data proto.PlayerLoginRes
	err := packet.DecodeData(&data)
	if err != nil {
		return
	}

	if data.Code != proto.ErrNil {
		cui.ShowLoginUi()
		cui.ShowDialog(data.Message, cui.DialogTypeError)
		return
	}
}

func PlayerLogoutRes(session *zNet.Session, packet *zNet.NetPacket) {}

func PlayerEnterRoomRes(session *zNet.Session, packet *zNet.NetPacket) {
	var data proto.PlayerEnterRoomRes
	err := packet.DecodeData(&data)
	if err != nil {
		return
	}

	if data.Code != proto.ErrNil {
		cui.ShowDialog(data.Message, cui.DialogTypeError)
		return
	}

	cui.ShowChatUi()
	for _, v := range data.ChatHistoryList {
		cui.GetChatUi().SpeakBroadcast(formatSpeakContent(v))
	}
	cui.G.Update(func(gui *gocui.Gui) error {
		return nil
	})
}

func PlayerLeaveRoomRes(session *zNet.Session, packet *zNet.NetPacket) {

	log.Println("离开房间成功")
}

func PlayerSpeakRes(session *zNet.Session, packet *zNet.NetPacket) {
	var data proto.PlayerSpeakRes
	_ = packet.DecodeData(&data)
	if data.Code != proto.ErrNil {
		cui.ShowDialog(data.Message, cui.DialogTypeError)
		return
	}
	log.Println("离开房间成功")
}

func SpeakBroadcast(session *zNet.Session, packet *zNet.NetPacket) {
	var data proto.ChatMessage
	err := packet.DecodeData(&data)
	if err != nil {
		//log.Println(err)
		cui.ShowDialog(err.Error(), cui.DialogTypeError)
		return
	}

	cui.GetChatUi().SpeakBroadcast(formatSpeakContent(data))
}

func RoomListRes(session *zNet.Session, packet *zNet.NetPacket) {
	var data proto.RoomListRes
	err := packet.DecodeData(&data)
	if err != nil {
		return
	}

	if data.Code != proto.ErrNil {
		cui.ShowDialog(data.Message, cui.DialogTypeError)
		return
	}

	var list []string
	list = append(list, "\t\tid\t\t\tname")
	for _, v := range data.RoomList {
		list = append(list, fmt.Sprintf("\t\t%d\t\t\t%s", v.Id, v.Name))
	}
	cui.ShowRoomUi(list)
}

func RoomPlayerListRes(session *zNet.Session, packet *zNet.NetPacket) {
	var data proto.RoomPlayerListRes
	err := packet.DecodeData(&data)
	if err != nil {
		return
	}

	if data.Code != proto.ErrNil {
		cui.ShowDialog(data.Message, cui.DialogTypeError)
		return
	}

	var list []string
	for _, v := range data.RoomPlayerList {
		list = append(list, v.Name)
	}

	time.Sleep(10 * time.Millisecond)
	cui.GetChatUi().ShowPlayerList(list)
}

func formatSpeakContent(data proto.ChatMessage) string {
	str := ""
	if data.Name != "" {
		str = fmt.Sprintf("%s [ %s ] say: %s", zTime.Seconds2String(data.Time), data.Name, data.Content)
	} else {
		str = fmt.Sprintf("%s   %s", zTime.Seconds2String(data.Time), data.Content)
	}

	return str
}
