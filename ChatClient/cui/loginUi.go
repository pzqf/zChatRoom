package cui

import (
	"fmt"
	"log"
	"strings"
	"zChatRoom/ChatClient/model"

	"github.com/jroimartin/gocui"
)

const (
	LoginWidth      = 30
	LoginPanelWidth = 50
	LoginPanelHigh  = 12

	ViewLoginPanelTitle = "Login_panel"
	ViewLoginInputTitle = "player_name_input"
)

type LoginUi struct {
	gui *gocui.Gui
}

func NewLogin(g *gocui.Gui) *LoginUi {
	return &LoginUi{
		gui: g,
	}
}

func (u *LoginUi) Show() error {
	tw, th := u.gui.Size()

	loginPanelView, err := u.gui.SetView(ViewLoginPanelTitle, tw/2-LoginPanelWidth/2, th/2-LoginPanelHigh/2, tw/2+LoginPanelWidth/2, th/2+LoginPanelHigh/2)
	if err != nil && err != gocui.ErrUnknownView {
		log.Println("Failed to create tasks view:", err)
		return err
	}
	_, _ = fmt.Fprintln(loginPanelView, "")
	_, _ = fmt.Fprintln(loginPanelView, "")
	_, _ = fmt.Fprintln(loginPanelView, "")
	_, _ = fmt.Fprintln(loginPanelView, "               ❤online chat room❤")

	inputViews, err := u.gui.SetView(ViewLoginInputTitle, tw/2-LoginWidth/2, th/2+1, tw/2+LoginWidth/2, th/2+3)
	if err != nil && err != gocui.ErrUnknownView {
		log.Println("Failed to create tasks view:", err)
		return err
	}
	inputViews.Title = "please input you name"
	inputViews.FgColor = gocui.ColorCyan
	inputViews.SelBgColor = gocui.ColorBlue
	inputViews.SelFgColor = gocui.ColorBlack

	inputViews.Editable = true
	u.gui.Cursor = true

	if err = u.gui.SetKeybinding(ViewLoginInputTitle, gocui.KeyEnter, gocui.ModNone, u.playerNameInput); err != nil {
		return err
	}

	if _, err = u.gui.SetCurrentView(ViewLoginInputTitle); err != nil {
		return err
	}

	u.gui.Update(func(gui *gocui.Gui) error {
		return nil
	})

	return nil
}

func (u *LoginUi) playerNameInput(g *gocui.Gui, iv *gocui.View) error {
	var err error
	iv.Rewind()

	if iv.Buffer() != "" {
		name := iv.Buffer()
		name = strings.Replace(name, "\n", "", -1)
		u.OnInput(name)
	} else {
		return nil
	}

	//iv.Editable = false
	//iv.Clear()
	//_ = iv.SetCursor(0, 0)

	g.DeleteKeybindings(ViewLoginInputTitle)
	if err = g.DeleteView(ViewLoginInputTitle); err != nil {
		return err
	}
	if err = g.DeleteView(ViewLoginPanelTitle); err != nil {
		return err
	}

	return nil
}

func (u *LoginUi) OnInput(str string) {
	model.Login(str)
}
