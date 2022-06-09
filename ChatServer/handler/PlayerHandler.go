package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"
	"zChatRoom/ChatServer/gm"
	"zChatRoom/ChatServer/player"
	"zChatRoom/ChatServer/playerMgr"
	"zChatRoom/ChatServer/room"
	"zChatRoom/proto"

	"github.com/pzqf/zEngine/zLog"
	"go.uber.org/zap"

	"github.com/pzqf/zUtil/zKeyWordFilter"

	"github.com/pzqf/zEngine/zNet"
	uuid "github.com/satori/go.uuid"
)

func RegisterPlayerHandler() error {
	err := zNet.RegisterHandler(proto.PlayerLogin, PlayerLogin)
	if err != nil {
		return err
	}
	err = zNet.RegisterHandler(proto.PlayerLogout, PlayerLogout)
	if err != nil {
		return err
	}

	err = zNet.RegisterHandler(proto.PlayerEnterRoom, PlayerEnterRoom)
	if err != nil {
		return err
	}

	err = zNet.RegisterHandler(proto.PlayerLeaveRoom, PlayerLeaveRoom)
	if err != nil {
		return err
	}

	err = zNet.RegisterHandler(proto.PlayerSpeak, PlayerSpeak)
	if err != nil {
		return err
	}
	return nil
}

func PlayerLogin(session zNet.Session, protoId int32, data []byte) {
	var reqData proto.PlayerLoginReq
	resData := proto.PlayerLoginRes{
		Code:    proto.ErrNil,
		Message: "success",
	}
	fmt.Println(string(data))
	err := json.Unmarshal(data, &reqData)
	//err := packet.DecodeData(&reqData)
	if err != nil {
		resData.Code = proto.RrrLogin
		resData.Message = err.Error() + "ddddddd"
		d, _ := json.Marshal(resData)
		_ = session.Send(proto.PlayerLogin, d)
		return
	}

	err = playerMgr.CheckPlayerName(reqData.UserName)
	if err != nil {
		resData.Code = proto.RrrLogin
		resData.Message = err.Error()
		d, _ := json.Marshal(resData)
		_ = session.Send(proto.PlayerLogin, d)
		return
	}

	newPlayer := player.Player{
		Id:        uuid.NewV4().String(),
		Session:   session.(*zNet.TcpServerSession),
		Name:      reqData.UserName,
		LoginTime: time.Now(),
	}
	playerMgr.AddPlayer(&newPlayer)

	zLog.Info("Add new Player", zap.String("id", newPlayer.Id))

	d, _ := json.Marshal(resData)
	_ = session.Send(proto.PlayerLogin, d)

	RoomList(session, proto.RoomList, nil)

}

func PlayerLogout(session zNet.Session, protoId int32, data []byte) {
	resData := proto.PlayerLogoutRes{
		Code:    proto.ErrNil,
		Message: "success",
	}
	//si := session.(*zNet.TcpServerSession)

	p, err := playerMgr.GetPlayerBySid(session.GetSid())
	if err != nil {
		resData.Code = proto.RrrLogout
		resData.Message = err.Error()
		d, _ := json.Marshal(resData)
		_ = session.Send(proto.PlayerLogout, d)
		return
	}

	if p.RoomId != 0 {
		r, err := room.GetRoom(p.RoomId)
		if err != nil {
			_ = r.DelPlayer(p.Id)
		}
		p.RoomId = 0
	}
	d, _ := json.Marshal(resData)
	_ = session.Send(proto.PlayerLogout, d)
}

func PlayerEnterRoom(session zNet.Session, protoId int32, data []byte) {
	var reqData proto.PlayerEnterRoomReq
	resData := proto.PlayerEnterRoomRes{
		Code:    proto.ErrNil,
		Message: "success",
	}
	defer func() {
		d, _ := json.Marshal(resData)
		_ = session.Send(proto.PlayerEnterRoom, d)
	}()

	//err := packet.DecodeData(&reqData)
	err := json.Unmarshal(data, &reqData)
	if err != nil {
		resData.Code = proto.ErrEnterRoom
		resData.Message = err.Error()
		return
	}

	p, err := playerMgr.GetPlayerBySid(session.GetSid())
	if err != nil {
		resData.Code = proto.ErrEnterRoom
		resData.Message = err.Error()
		return
	}

	if p.RoomId != 0 {
		r, err := room.GetRoom(reqData.RoomId)
		if err != nil {
			_ = r.DelPlayer(p.Id)
		}
		p.RoomId = 0
	}

	r, err := room.GetRoom(reqData.RoomId)
	if err != nil {
		resData.Code = proto.ErrEnterRoom
		resData.Message = err.Error()
		return
	}

	err = r.AddPlayer(p)
	if err != nil {
		resData.Code = proto.ErrEnterRoom
		resData.Message = err.Error()
		return
	}

	p.RoomId = reqData.RoomId
	list := r.HistoryChatQueue.Get()
	for _, v := range list {
		resData.ChatHistoryList = append(resData.ChatHistoryList, v.(proto.ChatMessage))
	}
	sort.Slice(resData.ChatHistoryList, func(i, j int) bool {
		return resData.ChatHistoryList[i].Time < resData.ChatHistoryList[j].Time
	})

	//fmt.Println("玩家", p.Name, "进入房间", reqData.RoomId, "成功")
	zLog.Info("Player enter room", zap.String("id", p.Id), zap.String("name", p.Name), zap.Int32("room_id", reqData.RoomId))

	time.Sleep(10 * time.Millisecond)
	//r.UpdateRoomPlayerList()
}

func PlayerLeaveRoom(session zNet.Session, protoId int32, data []byte) {
	resData := proto.PlayerLeaveRoomRes{
		Code:    proto.ErrNil,
		Message: "success",
	}
	defer func() {
		d, _ := json.Marshal(resData)
		_ = session.Send(proto.PlayerLeaveRoom, d)
	}()

	p, err := playerMgr.GetPlayerBySid(session.GetSid())
	if err != nil {
		resData.Code = proto.ErrLeaveRoom
		resData.Message = err.Error()
		return
	}

	if p.RoomId == 0 {
		resData.Code = proto.ErrLeaveRoom
		resData.Message = "not in any room"
		return
	}

	r, err := room.GetRoom(p.RoomId)
	if err != nil {
		resData.Code = proto.ErrLeaveRoom
		resData.Message = err.Error()
		return
	}
	err = r.DelPlayer(p.Id)
	if err != nil {
		resData.Code = proto.ErrLeaveRoom
		resData.Message = err.Error()
		return
	}
	p.RoomId = 0

	//r.UpdateRoomPlayerList()
}

func PlayerSpeak(session zNet.Session, protoId int32, data []byte) {
	var reqData proto.PlayerSpeakReq
	resData := proto.PlayerSpeakRes{
		Code:    proto.ErrNil,
		Message: "success",
	}
	defer func() {
		d, _ := json.Marshal(resData)
		_ = session.Send(proto.PlayerSpeak, d)
	}()
	//err := packet.DecodeData(&reqData)
	err := json.Unmarshal(data, &reqData)
	if err != nil {
		log.Println("receives speak", reqData.Content, time.Now())
		resData.Code = proto.ErrSpeak
		resData.Message = err.Error()
		return
	}

	p, err := playerMgr.GetPlayerBySid(session.GetSid())
	if err != nil {
		resData.Code = proto.ErrSpeak
		resData.Message = err.Error()
		return
	}

	if p.RoomId == 0 {
		resData.Code = proto.ErrSpeak
		resData.Message = "not in any chat room"
		return
	}

	r, err := room.GetRoom(p.RoomId)
	if err != nil {
		resData.Code = proto.ErrSpeak
		resData.Message = err.Error()
		return
	}

	if reqData.Content == "" {
		return
	}

	content := strings.TrimSpace(reqData.Content)
	if content[0] == '/' {
		gm.Process(session, content)
		return
	}

	content = zKeyWordFilter.Filter(reqData.Content)

	msg := proto.ChatMessage{
		Uid:     p.Id,
		Name:    p.Name,
		Time:    time.Now().Unix(),
		Content: content,
	}
	_ = r.NewSpeak(msg)
}
