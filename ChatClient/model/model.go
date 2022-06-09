package model

import (
	"encoding/json"
	"fmt"
	"zChatRoom/proto"

	"github.com/pzqf/zEngine/zNet"
)

var client *zNet.TcpClient
var isLogin bool

func Init(tcpClient *zNet.TcpClient) {
	client = tcpClient
	isLogin = false
}

func Login(username string) {
	data, err := json.Marshal(proto.PlayerLoginReq{
		UserName: username,
	})
	if err != nil {
		return
	}
	err = client.Send(proto.PlayerLogin, data)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func SelectRoom(roomId int32) {
	data, err := json.Marshal(proto.PlayerEnterRoomReq{
		RoomId: roomId,
	})
	err = client.Send(proto.PlayerEnterRoom, data)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func Speak(content string) {
	data, err := json.Marshal(proto.PlayerSpeakReq{
		Content: content,
	})
	err = client.Send(proto.PlayerSpeak, data)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func GetRoomList() {
	if GetLogin() {
		err := client.Send(proto.RoomList, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func SetLogin(login bool) {
	isLogin = login
}

func GetLogin() bool {
	return isLogin
}

func HeartBeat() {
	for {
		if client != nil {
			client.Send(0, nil)
		}
	}

}
