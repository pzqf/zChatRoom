package playerMgr

import (
	"errors"
	"fmt"
	"time"
	"zChatRoom/ChatServer/player"
	"zChatRoom/ChatServer/room"

	"github.com/pzqf/zEngine/zNet"

	"github.com/pzqf/zUtil/zMap"
)

type Mgr struct {
	PlayerList zMap.Map
}

var mgr *Mgr

func InitDefaultPlayerMgr() {
	mgr = &Mgr{}
	go func() {
		for true {
			fmt.Println(time.Now(), "online player count:", mgr.PlayerList.Len())
			time.Sleep(time.Second * 5)
		}
	}()
}

func CheckPlayerName(name string) error {
	find := false
	mgr.PlayerList.Range(func(key, value interface{}) bool {
		playerInfo := value.(*player.Player)
		if playerInfo.Name == name {
			find = true
			return false
		}
		return true
	})

	if find {
		return errors.New("find same name player")
	}

	return nil
}
func GetDefaultMgr() *Mgr {
	return mgr
}

func AddPlayer(p *player.Player) {
	mgr.PlayerList.Store(p.Id, p)

	fmt.Println("添加玩家", p.Name, "成功，当前玩家总数:", mgr.PlayerList.Len())
}

func GetPlayerBySid(sid zNet.SessionIdType) (*player.Player, error) {
	var p *player.Player
	mgr.PlayerList.Range(func(key, value interface{}) bool {
		playerInfo := value.(*player.Player)
		if playerInfo.Session.GetSid() == sid {
			p = playerInfo
			return false
		}
		return true
	})
	if p == nil {
		return nil, errors.New("can't find player")
	}

	return p, nil
}

func GetPlayerByName(playerName string) *player.Player {
	var p *player.Player
	mgr.PlayerList.Range(func(key, value interface{}) bool {
		playerInfo := value.(*player.Player)
		if playerInfo.Name == playerName {
			p = playerInfo
			return false
		}
		return true
	})
	return p
}

func OnSessionClose(sid zNet.SessionIdType) {
	mgr.PlayerList.Range(func(key, value interface{}) bool {
		playerInfo := value.(*player.Player)
		if playerInfo.Session.GetSid() != sid {
			return true
		}

		if playerInfo.RoomId != 0 {
			r, err := room.GetRoom(playerInfo.RoomId)
			if err != nil {
				return false
			}
			err = r.DelPlayer(playerInfo.Id)
			if err != nil {
				return false
			}
		}
		mgr.PlayerList.Delete(playerInfo.Id)
		return false
	})
}
