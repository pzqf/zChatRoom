package room

import (
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
	"zChatRoom/ChatServer/player"
	"zChatRoom/ChatServer/segmenter"
	"zChatRoom/proto"

	"github.com/pzqf/zEngine/zObject"

	"github.com/pzqf/zUtil/zList"

	"github.com/pzqf/zUtil/zMap"
	"github.com/pzqf/zUtil/zQueue"
)

const maxHistoryChatCount = 50

type Room struct {
	zObject.BaseObject
	Name             string `json:"name"`
	PlayerList       zMap.Map
	HistoryChatQueue zQueue.Queue
	wordCount        zMap.Map
}

type RoomIdType = int32

func NewRoom(id RoomIdType) *Room {
	r := &Room{
		Name: fmt.Sprintf("room_%d", id),
	}
	r.SetId(id)
	ticker := time.NewTicker(60 * time.Second)
	go func() {
		for range ticker.C {
			r.wordCount.Range(func(key, value interface{}) bool {
				wordList := value.(*zList.List)
				wordList.Range(func(e *list.Element, value any) bool {
					if time.Now().Sub(e.Value.(time.Time)).Seconds() > 10*60 {
						wordList.Remove(e)
					}
					return true
				})
				return true
			})

		}
	}()

	ticker2 := time.NewTicker(time.Second)
	go func() {
		for range ticker2.C {
			r.UpdateRoomPlayerList()
		}
	}()

	return r
}

func (r *Room) AddPlayer(p *player.Player) error {
	r.PlayerList.Store(p.GetId(), p)
	return nil
}

func (r *Room) DelPlayer(uid string) error {
	p, exist := r.PlayerList.Get(uid)
	if !exist {
		return errors.New("player not in room")
	}
	name := p.(*player.Player).Name
	r.PlayerList.Delete(uid)
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
			Id:   p.Id.(string),
			Name: p.Name,
		})
		return true
	})

	r.PlayerList.Range(func(key, value interface{}) bool {
		p := value.(*player.Player)
		d, _ := json.Marshal(resPlayerData)
		_ = p.Session.Send(proto.RoomPlayerList, d)
		return true
	})
}

func (r *Room) NewSpeak(chatMsg proto.ChatMessage) error {
	//log.Println("??????", chatMsg.Name, "?????????", r.Id, "??????:", chatMsg.Content)
	//zLog.Info("Player speak", zap.String("name", chatMsg.Name),
	//	zap.Int32("room_id", r.Id), zap.String("content", chatMsg.Content))

	r.HistoryChatQueue.Enqueue(chatMsg)
	if r.HistoryChatQueue.Length() > maxHistoryChatCount {
		r.HistoryChatQueue.Dequeue()
	}

	r.BroadcastChatMsg(chatMsg)

	wordSlice := segmenter.Segment(chatMsg.Content)
	for _, word := range wordSlice {
		if _, ok := r.wordCount.Get(word); !ok {
			r.wordCount.Store(word, zList.New())
		}
		l, _ := r.wordCount.Get(word)
		l.(*zList.List).PushBack(time.Now())
	}

	return nil
}

func (r *Room) BroadcastChatMsg(chatMsg proto.ChatMessage) {
	r.PlayerList.Range(func(key, value interface{}) bool {
		p := value.(*player.Player)
		if p.Session != nil {
			d, _ := json.Marshal(chatMsg)
			_ = p.Session.Send(proto.SpeakBroadcast, d)
		}
		return true
	})
}

func (r *Room) GetHighWord() string {
	fmt.Println("GetHighWord", r.wordCount.Len())
	highWord := ""
	MaxCount := 0

	r.wordCount.Range(func(key, value interface{}) bool {
		wordList := value.(*zList.List)
		fmt.Println(key, wordList.Len())
		if wordList.Len() > MaxCount {
			highWord = key.(string)
			MaxCount = wordList.Len()
		}
		return true
	})

	return highWord
}
