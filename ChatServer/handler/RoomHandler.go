package handler

import (
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

func RoomList(session *zNet.Session, packet *zNet.NetPacket) {
	resData := proto.RoomListRes{
		Code:    0,
		Message: "success",
	}
	roomList, err := room.GetRoomList()
	if err != nil {
		resData.Code = proto.ErrRoomList
		resData.Message = err.Error()
		_ = session.Send(proto.PlayerLogin, resData)
	}
	for _, v := range roomList {
		roomInfo := proto.RoomInfo{
			Id:   v.Id,
			Name: v.Name,
		}
		resData.RoomList = append(resData.RoomList, roomInfo)
	}

	_ = session.Send(proto.RoomList, resData)
}
