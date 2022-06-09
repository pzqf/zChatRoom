package handler

import (
	"encoding/json"
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
	if err := zNet.RegisterHandler(proto.TestPing, TestPingRes); err != nil {
		log.Printf("RegisterHandler error %d", err)
		return
	}
}

func PlayerLoginRes(session zNet.Session, protoId int32, data []byte) {
	var d proto.PlayerLoginRes
	err := json.Unmarshal(data, &d)
	//err := packet.DecodeData(&data)
	if err != nil {
		return
	}

	if d.Code != proto.ErrNil {
		cui.ShowLoginUi()
		cui.ShowDialog(d.Message, cui.DialogTypeError)
		return
	}
}

func PlayerLogoutRes(session zNet.Session, protoId int32, data []byte) {}

func PlayerEnterRoomRes(session zNet.Session, protoId int32, data []byte) {
	var d proto.PlayerEnterRoomRes
	err := json.Unmarshal(data, &d)
	if err != nil {
		return
	}

	if d.Code != proto.ErrNil {
		cui.ShowDialog(d.Message, cui.DialogTypeError)
		return
	}

	cui.ShowChatUi()
	for _, v := range d.ChatHistoryList {
		cui.GetChatUi().SpeakBroadcast(formatSpeakContent(v))
	}
	cui.G.Update(func(gui *gocui.Gui) error {
		return nil
	})

	go func() {
		for true {
			sendData, err := json.Marshal(proto.TestPingReq{
				Time: time.Now(),
			})
			if err != nil {
				return
			}
			err = session.Send(proto.TestPing, sendData)
			if err != nil {
				return
			}
			time.Sleep(3 * time.Second)
		}
	}()
}

func PlayerLeaveRoomRes(session zNet.Session, protoId int32, data []byte) {

	log.Println("离开房间成功")
}

func PlayerSpeakRes(session zNet.Session, protoId int32, data []byte) {
	var d proto.PlayerSpeakRes
	err := json.Unmarshal(data, &d)
	if err != nil {
		return
	}
	if d.Code != proto.ErrNil {
		cui.ShowDialog(d.Message, cui.DialogTypeError)
		return
	}
}

func SpeakBroadcast(session zNet.Session, protoId int32, data []byte) {
	//var data proto.ChatMessage
	var d proto.ChatMessage
	err := json.Unmarshal(data, &d)
	//err := packet.DecodeData(&data)
	if err != nil {
		//log.Println(err)
		cui.ShowDialog(err.Error(), cui.DialogTypeError)
		return
	}

	cui.GetChatUi().SpeakBroadcast(formatSpeakContent(d))
}

func RoomListRes(session zNet.Session, protoId int32, data []byte) {
	//var data proto.RoomListRes
	var d proto.RoomListRes
	err := json.Unmarshal(data, &d)
	if err != nil {
		return
	}

	if d.Code != proto.ErrNil {
		cui.ShowDialog(d.Message, cui.DialogTypeError)
		return
	}

	var list []string
	list = append(list, "\t\tid\t\t\tname")
	for _, v := range d.RoomList {
		list = append(list, fmt.Sprintf("\t\t%d\t\t\t%s", v.Id, v.Name))
	}
	cui.ShowRoomUi(list)
}

func RoomPlayerListRes(session zNet.Session, protoId int32, data []byte) {
	//var data proto.RoomPlayerListRes
	var d proto.RoomPlayerListRes
	err := json.Unmarshal(data, &d)
	//err := packet.DecodeData(&data)
	if err != nil {
		return
	}

	if d.Code != proto.ErrNil {
		cui.ShowDialog(d.Message, cui.DialogTypeError)
		return
	}

	var list []string
	for _, v := range d.RoomPlayerList {
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

func TestPingRes(session zNet.Session, protoId int32, data []byte) {
	var d proto.TestPingRes
	err := json.Unmarshal(data, &d)
	if err != nil {
		return
	}

	cui.GetChatUi().SpeakBroadcast(fmt.Sprintf("==========ping====%s", time.Now().Sub(d.Time).String()))
}
