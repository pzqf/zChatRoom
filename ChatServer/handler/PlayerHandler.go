package handler

import (
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

func PlayerLogin(session *zNet.Session, packet *zNet.NetPacket) {
	var reqData proto.PlayerLoginReq
	resData := proto.PlayerLoginRes{
		Code:    proto.ErrNil,
		Message: "success",
	}
	err := packet.DecodeData(&reqData)
	if err != nil {
		resData.Code = proto.RrrLogin
		resData.Message = err.Error()
		_ = session.Send(proto.PlayerLogin, resData)
		return
	}

	err = playerMgr.CheckPlayerName(reqData.UserName)
	if err != nil {
		resData.Code = proto.RrrLogin
		resData.Message = err.Error()
		_ = session.Send(proto.PlayerLogin, resData)
		return
	}

	newPlayer := player.Player{
		Id:        uuid.NewV4().String(),
		Session:   session,
		Name:      reqData.UserName,
		LoginTime: time.Now(),
	}
	playerMgr.AddPlayer(&newPlayer)

	zLog.Info("Add new Player", zap.String("id", newPlayer.Id))

	_ = session.Send(proto.PlayerLogin, resData)

	RoomList(session, nil)

}

func PlayerLogout(session *zNet.Session, packet *zNet.NetPacket) {
	resData := proto.PlayerLogoutRes{
		Code:    proto.ErrNil,
		Message: "success",
	}

	p, err := playerMgr.GetPlayerBySid(session.GetSid())
	if err != nil {
		resData.Code = proto.RrrLogout
		resData.Message = err.Error()
		_ = session.Send(proto.PlayerEnterRoom, resData)
		return
	}

	if p.RoomId != 0 {
		r, err := room.GetRoom(p.RoomId)
		if err != nil {
			_ = r.DelPlayer(p.Id)
		}
		p.RoomId = 0
	}
	_ = session.Send(proto.PlayerEnterRoom, resData)
}

func PlayerEnterRoom(session *zNet.Session, packet *zNet.NetPacket) {
	var reqData proto.PlayerEnterRoomReq
	resData := proto.PlayerEnterRoomRes{
		Code:    proto.ErrNil,
		Message: "success",
	}

	err := packet.DecodeData(&reqData)
	if err != nil {
		resData.Code = proto.ErrEnterRoom
		resData.Message = err.Error()
		_ = session.Send(proto.PlayerEnterRoom, resData)
		return
	}

	p, err := playerMgr.GetPlayerBySid(session.GetSid())
	if err != nil {
		resData.Code = proto.ErrEnterRoom
		resData.Message = err.Error()
		_ = session.Send(proto.PlayerEnterRoom, resData)
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
		_ = session.Send(proto.PlayerEnterRoom, resData)
		return
	}

	err = r.AddPlayer(p)
	if err != nil {
		resData.Code = proto.ErrEnterRoom
		resData.Message = err.Error()
		_ = session.Send(proto.PlayerEnterRoom, resData)
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

	_ = session.Send(proto.PlayerEnterRoom, resData)
	//fmt.Println("玩家", p.Name, "进入房间", reqData.RoomId, "成功")
	zLog.Info("Player enter room", zap.String("id", p.Id), zap.String("name", p.Name), zap.Int32("room_id", reqData.RoomId))

	time.Sleep(10 * time.Millisecond)
	//r.UpdateRoomPlayerList()
}

func PlayerLeaveRoom(session *zNet.Session, packet *zNet.NetPacket) {
	resData := proto.PlayerLeaveRoomRes{
		Code:    proto.ErrNil,
		Message: "success",
	}

	p, err := playerMgr.GetPlayerBySid(session.GetSid())
	if err != nil {
		resData.Code = proto.ErrLeaveRoom
		resData.Message = err.Error()
		_ = session.Send(proto.PlayerLeaveRoom, resData)
		return
	}

	if p.RoomId == 0 {
		resData.Code = proto.ErrLeaveRoom
		resData.Message = "not in any room"
		_ = session.Send(proto.PlayerLeaveRoom, resData)
		return
	}

	r, err := room.GetRoom(p.RoomId)
	if err != nil {
		resData.Code = proto.ErrLeaveRoom
		resData.Message = err.Error()
		_ = session.Send(proto.PlayerLeaveRoom, resData)
		return
	}
	err = r.DelPlayer(p.Id)
	if err != nil {
		resData.Code = proto.ErrLeaveRoom
		resData.Message = err.Error()
		_ = session.Send(proto.PlayerLeaveRoom, resData)
		return
	}
	p.RoomId = 0

	_ = session.Send(proto.PlayerLeaveRoom, resData)

	//r.UpdateRoomPlayerList()
}

func PlayerSpeak(session *zNet.Session, packet *zNet.NetPacket) {
	var reqData proto.PlayerSpeakReq
	resData := proto.PlayerSpeakRes{
		Code:    proto.ErrNil,
		Message: "success",
	}
	err := packet.DecodeData(&reqData)
	if err != nil {
		log.Println("receives speak", reqData.Content, time.Now())
		resData.Code = proto.ErrSpeak
		resData.Message = err.Error()
		_ = session.Send(proto.PlayerSpeak, resData)
		return
	}

	p, err := playerMgr.GetPlayerBySid(session.GetSid())
	if err != nil {
		resData.Code = proto.ErrSpeak
		resData.Message = err.Error()
		_ = session.Send(proto.PlayerSpeak, resData)
		return
	}

	if p.RoomId == 0 {
		resData.Code = proto.ErrSpeak
		resData.Message = "not in any chat room"
		_ = session.Send(proto.PlayerSpeak, resData)
		return
	}

	r, err := room.GetRoom(p.RoomId)
	if err != nil {
		resData.Code = proto.ErrSpeak
		resData.Message = err.Error()
		_ = session.Send(proto.PlayerSpeak, resData)
		return
	}

	if reqData.Content == "" {
		_ = session.Send(proto.PlayerSpeak, resData)
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
