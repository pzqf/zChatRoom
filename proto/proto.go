package proto

import "time"

const (
	PlayerLogin     = 1001
	PlayerLogout    = 1002
	PlayerEnterRoom = 1003
	PlayerLeaveRoom = 1004
	PlayerSpeak     = 1005
	SpeakBroadcast  = 1010

	RoomList       = 2001
	RoomPlayerList = 2002

	TestPing = 3001
)

type PlayerLoginReq struct {
	UserName string `json:"user_name"`
}

type PlayerLoginRes struct {
	Code    int32  `json:"err_code"`
	Message string `json:"message"`
}

type PlayerLogoutRes struct {
	Code    int32  `json:"err_code"`
	Message string `json:"message"`
}

type PlayerEnterRoomReq struct {
	RoomId int32 `json:"room_id"`
}

type PlayerEnterRoomRes struct {
	Code            int32         `json:"err_code"`
	Message         string        `json:"message"`
	ChatHistoryList []ChatMessage `json:"chat_history_list"`
}

type PlayerLeaveRoomReq struct {
}

type PlayerLeaveRoomRes struct {
	Code    int32  `json:"err_code"`
	Message string `json:"message"`
}

type ChatMessage struct {
	Uid     string `json:"uid"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Time    int64  `json:"time"`
}

type PlayerSpeakReq struct {
	//Id      string `json:"id"`
	Content string `json:"content"`
}

type PlayerSpeakRes struct {
	Code    int32  `json:"err_code"`
	Message string `json:"message"`
}

type RoomInfo struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
}

type RoomListRes struct {
	Code     int32      `json:"err_code"`
	Message  string     `json:"message"`
	RoomList []RoomInfo `json:"room_list"`
}

type RoomPlayerInfo struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type RoomPlayerListRes struct {
	Code           int32            `json:"err_code"`
	Message        string           `json:"message"`
	RoomPlayerList []RoomPlayerInfo `json:"room_player_list"`
}

type TestPingReq struct {
	Id   int32     `json:"id"`
	Name string    `json:"name"`
	Time time.Time `json:"time"`
}

type TestPingRes struct {
	Id   int32     `json:"id"`
	Name string    `json:"name"`
	Time time.Time `json:"time"`
}
