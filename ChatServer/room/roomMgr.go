package room

import (
	"errors"
	"sync/atomic"

	"github.com/pzqf/zUtil/zMap"
)

type Mgr struct {
	RoomList     zMap.Map
	MaxRoomCount int32
	roomIdIndex  int32
}

func NewRoomMgr(maxRoomCount int32) *Mgr {
	m := Mgr{
		MaxRoomCount: maxRoomCount,
		roomIdIndex:  0,
	}
	for i := int32(0); i < maxRoomCount; i++ {
		_, _ = m.AddRoom()
	}

	return &m
}

var DefaultMgr *Mgr

func InitDefaultRoomMgr(maxRoomCount int32) {
	DefaultMgr = NewRoomMgr(maxRoomCount)
}

func (m *Mgr) AddRoom() (*Room, error) {
	if m.RoomList.Len() >= m.MaxRoomCount {
		return nil, errors.New("room count over max")
	}

	id := atomic.AddInt32(&m.roomIdIndex, 1)
	newRoom := NewRoom(id)
	m.RoomList.Store(newRoom.Id, newRoom)

	return newRoom, nil
}

func GetRoomList() ([]*Room, error) {
	var list []*Room
	DefaultMgr.RoomList.Range(func(key, value interface{}) bool {
		list = append(list, value.(*Room))
		return true
	})

	return list, nil
}

func GetRoom(roomId int32) (*Room, error) {
	var r *Room
	DefaultMgr.RoomList.Range(func(key, value interface{}) bool {
		if value.(*Room).Id == roomId {
			r = value.(*Room)
			return false
		}
		return true
	})
	if r == nil {
		return nil, errors.New("can't find room")
	}

	return r, nil
}
