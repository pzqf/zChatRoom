package room

import (
	"container/list"
	"errors"
	"fmt"
	"log"
	"time"
	"zChatRoom/ChatServer/player"
	"zChatRoom/ChatServer/segmenter"
	"zChatRoom/proto"

	"github.com/pzqf/zUtil/zMap"
	"github.com/pzqf/zUtil/zQueue"
)

const maxHistoryChatCount = 50

type Room struct {
	Id               int32  `json:"id"`
	Name             string `json:"name"`
	PlayerList       zMap.Map
	HistoryChatQueue zQueue.Queue
	wordCount        zMap.Map
}

func NewRoom(id int32) *Room {
	return &Room{
		Id:   id,
		Name: fmt.Sprintf("room_%d", id),
	}
}

func (r *Room) AddPlayer(p *player.Player) error {
	r.PlayerList.Store(p.Id, p)
	r.UpdateRoomPlayerList()

	return nil
}

func (r *Room) DelPlayer(uid string) error {
	p, exist := r.PlayerList.Get(uid)
	if !exist {
		return errors.New("player not in room")
	}
	name := p.(*player.Player).Name
	r.PlayerList.Delete(uid)
	r.UpdateRoomPlayerList()
	log.Println("player", name, "left", r.Name, r.Id)
	chatMsg := proto.ChatMessage{
		Content: name + " left room",
		Time:    time.Now().Unix(),
	}
	r.BroadcastChatMsg(chatMsg)
	return nil
}

func (r *Room) GetPlayerList() []*player.Player {
	var lp []*player.Player
	r.PlayerList.Range(func(key, value interface{}) bool {
		lp = append(lp, value.(*player.Player))
		return true
	})
	return lp
}

func (r *Room) UpdateRoomPlayerList() {
	resPlayerData := proto.RoomPlayerListRes{
		Code:    0,
		Message: "success",
	}

	r.PlayerList.Range(func(key, value interface{}) bool {
		p := value.(*player.Player)
		resPlayerData.RoomPlayerList = append(resPlayerData.RoomPlayerList, proto.RoomPlayerInfo{
			Id:   p.Id,
			Name: p.Name,
		})
		return true
	})

	r.PlayerList.Range(func(key, value interface{}) bool {
		p := value.(*player.Player)
		_ = p.Session.Send(proto.RoomPlayerList, resPlayerData)
		return true
	})
}

func (r *Room) NewSpeak(chatMsg proto.ChatMessage) error {
	//log.Println("玩家", chatMsg.Name, "在房间", r.Id, "发言:", chatMsg.Content)

	r.HistoryChatQueue.Enqueue(chatMsg)
	if r.HistoryChatQueue.Length() > maxHistoryChatCount {
		r.HistoryChatQueue.Dequeue()
	}

	r.BroadcastChatMsg(chatMsg)

	wordSlice := segmenter.Segment(chatMsg.Content)
	for _, word := range wordSlice {
		if _, ok := r.wordCount.Get(word); !ok {
			//r.wordCount[word] = &list.List{}
			r.wordCount.Store(word, &list.List{})
		}
		l, _ := r.wordCount.Get(word)
		l.(*list.List).PushBack(time.Now())
	}

	r.wordCount.Range(func(key, value interface{}) bool {
		wordList := value.(*list.List)
		var n *list.Element
		for e := wordList.Front(); e != nil; e = n {
			n = e.Next()
			if time.Now().Sub(e.Value.(time.Time)).Seconds() > 10*60 {
				wordList.Remove(e)
			}
		}
		return true
	})

	return nil
}

func (r *Room) BroadcastChatMsg(chatMsg proto.ChatMessage) {
	r.PlayerList.Range(func(key, value interface{}) bool {
		p := value.(*player.Player)
		if p.Session != nil {
			p.Session.Send(proto.SpeakBroadcast, chatMsg)
		}
		return true
	})
}

func (r *Room) GetHighWord() string {
	r.wordCount.Range(func(key, value interface{}) bool {
		wordList := value.(*list.List)
		var n *list.Element
		for e := wordList.Front(); e != nil; e = n {
			n = e.Next()
			if time.Now().Sub(e.Value.(time.Time)).Seconds() > 10*60 {
				wordList.Remove(e)
			}
		}
		return true
	})

	highWord := ""
	MaxCount := 0

	r.wordCount.Range(func(key, value interface{}) bool {
		wordList := value.(*list.List)
		if wordList.Len() > MaxCount {
			highWord = key.(string)
			MaxCount = wordList.Len()
		}
		return true
	})

	return highWord
}
