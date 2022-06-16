package handler

import (
	"encoding/json"
	"zChatRoom/proto"

	"github.com/pzqf/zEngine/zNet"
)

func Init() error {

	zNet.InitDispatcherWorkerPool(10000)

	err := RegisterPlayerHandler()
	if err != nil {
		return err
	}
	err = RegisterRoomHandler()
	if err != nil {
		return err
	}
	err = zNet.RegisterHandler(proto.TestPing, TestPing)
	if err != nil {
		return err
	}

	return nil
}

func TestPing(session zNet.Session, protoId int32, data []byte) {
	var reqData proto.TestPingReq

	_ = json.Unmarshal(data, &reqData)

	resData := proto.TestPingRes{
		Id:   reqData.Id,
		Name: reqData.Name,
		Time: reqData.Time,
	}
	d, _ := json.Marshal(resData)
	_ = session.Send(proto.TestPing, d)
}
