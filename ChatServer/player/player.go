package player

import (
	"time"

	"github.com/pzqf/zEngine/zNet"
)

type Player struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Session   *zNet.Session
	RoomId    int32
	LoginTime time.Time
}

func (p *Player) SendData(protoId int32, msg interface{}) error {
	return p.Session.Send(protoId, msg)
}
