package cui

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"zChatRoom/ChatClient/model"

	"github.com/jroimartin/gocui"
)

const (
	RoomListWidth = 50

	ViewRoomListTitle = "listRoom"
	ViewSelectTitle   = "selectRoom"
)

type RoomListUi struct {
	gui *gocui.Gui
}

func NewRoomList(g *gocui.Gui) *RoomListUi {
	return &RoomListUi{
		gui: g,
	}
}
func (u *RoomListUi) Show(roomList []string) error {
	tw, th := u.gui.Size()

	high := len(roomList)
	if high > th-8 {
		high = th - 8
	}
	if high < 5 {
		high = 5
	}

	top := (th-high-3)/2 - 1

	roomListViews, err := u.gui.SetView(ViewRoomListTitle, tw/2-RoomListWidth/2, top, tw/2+RoomListWidth/2, top+high+1)
	if err != nil && err != gocui.ErrUnknownView {
		log.Println("Failed to create tasks view:", err)
		return err
	}
	roomListViews.Title = " room list "
	roomListViews.FgColor = gocui.ColorCyan

	for _, v := range roomList {
		if roomListViews != nil {
			_, _ = fmt.Fprintln(roomListViews, v)
		}
	}

	inputViews, err := u.gui.SetView(ViewSelectTitle, tw/2-RoomListWidth/2, top+high+2, tw/2+RoomListWidth/2, top+high+4)
	if err != nil && err != gocui.ErrUnknownView {
		log.Println("Failed to create tasks view:", err)
		return err
	}
	inputViews.Title = "please input room id"
	inputViews.FgColor = gocui.ColorCyan
	inputViews.SelBgColor = gocui.ColorBlue
	inputViews.SelFgColor = gocui.ColorBlack

	inputViews.Editable = true
	u.gui.Cursor = true

	if err = u.gui.SetKeybinding(ViewSelectTitle, gocui.KeyEnter, gocui.ModNone, u.selectRoomInput); err != nil {
		return err
	}

	if _, err = u.gui.SetCurrentView(ViewSelectTitle); err != nil {
		return err
	}

	u.gui.Update(func(gui *gocui.Gui) error {
		return nil
	})

	return nil
}

func (u *RoomListUi) selectRoomInput(g *gocui.Gui, iv *gocui.View) error {
	iv.Rewind()

	roomIdStr := iv.Buffer()
	if roomIdStr == "" {
		return nil
	}
	roomIdStr = strings.Replace(roomIdStr, "\n", "", -1)
	roomId, err := strconv.Atoi(roomIdStr)
	if err != nil {
		return nil
	}

	if iv.Buffer() != "" {
		u.OnSelectRoom(int32(roomId))
	} else {
		return nil
	}

	iv.Clear()
	_ = iv.SetCursor(0, 0)

	return err
}

func (u *RoomListUi) OnSelectRoom(roomId int32) {
	model.SelectRoom(roomId)
}
