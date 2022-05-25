package gm

import (
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

func Process(session *zNet.Session, content string) {
	resData := proto.PlayerSpeakRes{
		Code:    0,
		Message: "success",
	}
	p, err := playerMgr.GetPlayerBySid(session.GetSid())
	if err != nil {
		resData.Code = 2
		resData.Message = err.Error()
		_ = session.Send(proto.PlayerSpeak, resData)
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
			_ = session.Send(proto.SpeakBroadcast, msg)
			return
		}

		playerName := wordList[1]

		otherPlayer := playerMgr.GetPlayerByName(playerName)
		if p == nil {
			msg.Content = "player not online"
			_ = session.Send(proto.SpeakBroadcast, msg)
			return
		}

		msg.Content = fmt.Sprintf("GM: player:%s, login time:%s, online:%s, room id:%d", playerName,
			zTime.Time2String(otherPlayer.LoginTime), time.Now().Sub(otherPlayer.LoginTime).String(), otherPlayer.RoomId)

		_ = session.Send(proto.SpeakBroadcast, msg)

	case "/popular":
		if len(wordList) < 2 {
			msg.Content = "please input room id"
			_ = session.Send(proto.SpeakBroadcast, msg)
			return
		}

		roomIdStr := wordList[1]

		roomId, err := strconv.Atoi(roomIdStr)
		if err != nil {
			msg.Content = "room id wrong"
			_ = session.Send(proto.SpeakBroadcast, msg)
			return
		}

		roomInfo, err := room.GetRoom(int32(roomId))
		if err != nil {
			msg.Content = "can't find the room"
			_ = session.Send(proto.SpeakBroadcast, msg)
			return
		}

		msg.Content = fmt.Sprintf("GM: %s ", roomInfo.GetHighWord())
		_ = session.Send(proto.SpeakBroadcast, msg)
	default:
		msg.Content = fmt.Sprintf("GM: can't find command %s ", command)
		_ = session.Send(proto.SpeakBroadcast, msg)
	}
}
