package gm

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
	"zChatRoom/ChatServer/playerMgr"
	"zChatRoom/ChatServer/room"
	"zChatRoom/proto"

	"github.com/pzqf/zEngine/zNet"
	"github.com/pzqf/zUtil/zTime"
)

func Process(session zNet.Session, content string) {
	resData := proto.PlayerSpeakRes{
		Code:    0,
		Message: "success",
	}
	p, err := playerMgr.GetPlayerBySid(session.GetSid())
	if err != nil {
		resData.Code = 2
		resData.Message = err.Error()
		d, _ := json.Marshal(resData)
		_ = session.Send(proto.PlayerSpeak, d)
		//_ = session.Send(proto.PlayerSpeak, resData)
		return
	}

	wordList := strings.Split(content, " ")
	if len(wordList) <= 0 {
		return
	}

	command := wordList[0]

	msg := proto.ChatMessage{
		Time: time.Now().Unix(),
	}

	switch command {
	case "/stats":
		if len(wordList) < 2 {
			msg.Content = "please input player name"
			d, _ := json.Marshal(msg)
			_ = session.Send(proto.SpeakBroadcast, d)
			//_ = session.Send(proto.SpeakBroadcast, msg)
			return
		}

		playerName := wordList[1]

		otherPlayer := playerMgr.GetPlayerByName(playerName)
		if p == nil {
			msg.Content = "player not online"
			d, _ := json.Marshal(msg)
			_ = session.Send(proto.SpeakBroadcast, d)
			return
		}

		msg.Content = fmt.Sprintf("GM: player:%s, login time:%s, online:%s, room id:%d", playerName,
			zTime.Time2String(otherPlayer.LoginTime), time.Now().Sub(otherPlayer.LoginTime).String(), otherPlayer.RoomId)

		d, _ := json.Marshal(msg)
		_ = session.Send(proto.SpeakBroadcast, d)

	case "/popular":
		if len(wordList) < 2 {
			msg.Content = "please input room id"
			d, _ := json.Marshal(msg)
			_ = session.Send(proto.SpeakBroadcast, d)
			return
		}

		roomIdStr := wordList[1]

		roomId, err := strconv.Atoi(roomIdStr)
		if err != nil {
			msg.Content = "room id wrong"
			d, _ := json.Marshal(msg)
			_ = session.Send(proto.SpeakBroadcast, d)
			return
		}

		roomInfo, err := room.GetRoom(int32(roomId))
		if err != nil {
			msg.Content = "can't find the room"
			d, _ := json.Marshal(msg)
			_ = session.Send(proto.SpeakBroadcast, d)
			return
		}

		msg.Content = fmt.Sprintf("GM: %s ", roomInfo.GetHighWord())
		d, _ := json.Marshal(msg)
		_ = session.Send(proto.SpeakBroadcast, d)
	default:
		msg.Content = fmt.Sprintf("GM: can't find command %s ", command)
		d, _ := json.Marshal(msg)
		_ = session.Send(proto.SpeakBroadcast, d)
	}
}
