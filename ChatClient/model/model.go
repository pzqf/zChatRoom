package model

import (
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
	err := client.Send(proto.PlayerLogin, &proto.PlayerLoginReq{
		UserName: username,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
}

func SelectRoom(roomId int32) {
	err := client.Send(proto.PlayerEnterRoom, &proto.PlayerEnterRoomReq{
		RoomId: roomId,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
}

func Speak(content string) {
	err := client.Send(proto.PlayerSpeak, &proto.PlayerSpeakReq{
		Content: content,
	})
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
