package cui

import (
	"log"
	"zChatRoom/ChatClient/model"

	"github.com/jroimartin/gocui"
)

//

var G *gocui.Gui

var LastView = ViewLoginInputTitle

func InitUi() {
	var err error
	G, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Println("Failed to create a GUI:", err)
		return
	}

	G.Highlight = true
	G.SelFgColor = gocui.ColorBlue
	G.BgColor = gocui.ColorBlack
	G.FgColor = gocui.ColorWhite

	err = G.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		return gocui.ErrQuit
	})

	err = G.SetKeybinding("", gocui.KeyCtrlQ, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		if LastView == ViewTitleInput {
			_ = destroyChatUi(G)
			model.GetRoomList()
		}
		return err
	})

	if err != nil {
		log.Println("Could not set key binding:", err)
		return
	}
}

func StartUi() error {
	if G != nil {
		err := G.MainLoop()
		log.Println("Main loop has finished:", err)
		if err != nil {
			return err
		}
	}

	return nil
}

func StopUi() {
	if G != nil {
		G.Close()
	}
}

func ShowLoginUi() {
	l := NewLogin(G)
	_ = l.Show()
	LastView = ViewLoginInputTitle
}

func ShowRoomUi(roomList []string) {
	r := NewRoomList(G)
	_ = r.Show(roomList)
	LastView = ViewSelectTitle
}

var c *ChatUi

func ShowChatUi() {
	c = NewChat(G)
	_ = c.Show()
	LastView = ViewTitleInput
}

func ShowDialog(str string, dt DialogType) {
	d := NewDialog(G, dt)
	_ = d.Show(str)
}

func GetChatUi() *ChatUi {
	return c
}
