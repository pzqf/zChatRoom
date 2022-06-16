package handler

import (
	"encoding/json"
	"zChatRoom/ChatServer/room"
	"zChatRoom/proto"

	"github.com/pzqf/zEngine/zNet"
)

func RegisterRoomHandler() error {
	err := zNet.RegisterHandler(proto.RoomList, RoomList)
	if err != nil {
		return err
	}

	return nil
}

func RoomList(session zNet.Session, protoId int32, data []byte) {
	resData := proto.RoomListRes{
		Code:    0,
		Message: "success",
	}
	roomList, err := room.GetRoomList()
	if err != nil {
		resData.Code = proto.ErrRoomList
		resData.Message = err.Error()
		d, _ := json.Marshal(resData)
		_ = session.Send(proto.RoomList, d)
	}
	for _, v := range roomList {
		roomInfo := proto.RoomInfo{
			Id:   v.GetId().(room.RoomIdType),
			Name: v.Name,
		}
		resData.RoomList = append(resData.RoomList, roomInfo)
	}
	d, _ := json.Marshal(resData)
	_ = session.Send(proto.RoomList, d)
}
