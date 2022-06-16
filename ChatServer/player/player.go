package player

import (
	"time"

	"github.com/pzqf/zEngine/zObject"

	"github.com/pzqf/zEngine/zNet"
)

type Player struct {
	//Id        string `json:"id"`
	zObject.Object
	Name      string `json:"name"`
	Session   *zNet.TcpServerSession
	RoomId    int32
	LoginTime time.Time
}
