package cui

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
	"zChatRoom/ChatClient/model"

	"github.com/jroimartin/gocui"
)

const (
	TitleHigh       = 3
	ListPlayerWidth = 20
	ShortcutWidth   = 30
	InputHigh       = 3

	ViewTitleTitle      = "title"
	ViewTitleListPlayer = "list players"
	ViewTitleHistory    = "history"
	ViewTitleShortcut   = "shortcut"
	ViewTitleInput      = "input"
)

type ChatUi struct {
	gui       *gocui.Gui
	lastMsg   []string
	lastIndex int
	mu        sync.Mutex
}

func NewChat(g *gocui.Gui) *ChatUi {
	return &ChatUi{
		gui:       g,
		lastIndex: 0,
		mu:        sync.Mutex{},
	}
}

func (u *ChatUi) Show() error {
	tw, th := u.gui.Size()

	////////////////////////
	titleView, err := u.gui.SetView(ViewTitleTitle, 0, 0, tw-1, TitleHigh)
	if err != nil && err != gocui.ErrUnknownView {
		log.Println("Failed to create title view:", err)
		return err
	}
	//titleView.Title = ViewTitleListPlayer
	titleView.FgColor = gocui.ColorYellow
	titleView.Highlight = true
	titleView.SelFgColor = gocui.ColorGreen
	str := ""
	b := (tw-1-len("❤❤❤ chat online ❤❤❤ "))/2 - 12
	if b <= 10 {
		b = 10
	}
	for i := 0; i < b; i++ {
		str += " "
	}
	_, _ = fmt.Fprintln(titleView, str+"❤❤❤ chat online ❤❤❤ ")
	_, _ = fmt.Fprintln(titleView, str+"       -- by zqf▄︻┻┳═一…… ☆")

	////////////////////////
	listPlayerView, err := u.gui.SetView(ViewTitleListPlayer, 0, TitleHigh+1, ListPlayerWidth, th-InputHigh-1)
	if err != nil && err != gocui.ErrUnknownView {
		log.Println("Failed to create projects view:", err)
		return err
	}
	listPlayerView.Title = " player list "
	listPlayerView.FgColor = gocui.ColorCyan
	listPlayerView.Highlight = true

	////////////////////////
	historyView, err := u.gui.SetView(ViewTitleHistory, ListPlayerWidth+1, TitleHigh+1, tw-ShortcutWidth-1, th-InputHigh-1)
	if err != nil && err != gocui.ErrUnknownView {
		log.Println("Failed to create tasks view:", err)
		return err
	}
	historyView.Title = " chat history "
	historyView.FgColor = gocui.ColorCyan
	historyView.Autoscroll = true
	historyView.Wrap = true

	////////////////////////
	shortcutView, err := u.gui.SetView(ViewTitleShortcut, tw-ShortcutWidth, TitleHigh+1, tw-1, th-InputHigh-1)
	if err != nil && err != gocui.ErrUnknownView {
		log.Println("Failed to create output view:", err)
		return err
	}
	shortcutView.Title = "shortcut and gm command"
	shortcutView.FgColor = gocui.ColorWhite
	shortcutView.SelBgColor = gocui.ColorBlue
	shortcutView.SelFgColor = gocui.ColorBlack

	_, _ = fmt.Fprintln(shortcutView, "ctrl+c, quit")
	_, _ = fmt.Fprintln(shortcutView, "ctrl+q, quit room")
	_, _ = fmt.Fprintln(shortcutView, "-------------")
	_, _ = fmt.Fprintln(shortcutView, "GM: /stats player name")
	_, _ = fmt.Fprintln(shortcutView, "    like: /stats aa")
	_, _ = fmt.Fprintln(shortcutView, "GM: /popular room id")
	_, _ = fmt.Fprintln(shortcutView, "    like: /popular 1")

	////////////////////////
	inputViews, err := u.gui.SetView(ViewTitleInput, 0, th-InputHigh, tw-1, th-1)
	if err != nil && err != gocui.ErrUnknownView {
		log.Println("Failed to create tasks view:", err)
		return err
	}
	inputViews.Title = " please input your speech "
	inputViews.FgColor = gocui.ColorCyan
	inputViews.SelBgColor = gocui.ColorBlue
	inputViews.SelFgColor = gocui.ColorBlack

	inputViews.Editable = true
	u.gui.Cursor = true
	if _, err = u.gui.SetCurrentView(ViewTitleInput); err != nil {
		return err
	}
	if err = u.gui.SetKeybinding(ViewTitleInput, gocui.KeyEnter, gocui.ModNone, u.copyInput); err != nil {
		return err
	}

	if err = u.gui.SetKeybinding(ViewTitleInput, gocui.KeyArrowUp, gocui.ModNone, u.arrowUp); err != nil {
		return err
	}

	if err = u.gui.SetKeybinding(ViewTitleInput, gocui.KeyArrowDown, gocui.ModNone, u.arrowDown); err != nil {
		return err
	}

	return nil
}

func (u *ChatUi) copyInput(g *gocui.Gui, iv *gocui.View) error {
	var err error
	iv.Rewind()

	if iv.Buffer() != "" {
		content := strings.Replace(iv.Buffer(), "\n", "", -1)
		u.lastMsg = append(u.lastMsg, content)
		u.OnSpeak(content)
		u.lastIndex = 0
	} else {
		return nil
	}

	iv.Clear()
	_ = iv.SetCursor(0, 0)

	return err
}

func (u *ChatUi) arrowUp(g *gocui.Gui, iv *gocui.View) error {
	iv.Clear()
	if u.lastIndex < 0 {
		u.lastIndex = 0
	}
	if u.lastIndex >= 0 && u.lastIndex <= len(u.lastMsg)-1 {
		_, _ = fmt.Fprintln(iv, u.lastMsg[u.lastIndex])
		_ = iv.SetCursor(len(u.lastMsg[u.lastIndex]), 0)
	}
	u.lastIndex++
	if u.lastIndex >= len(u.lastMsg) {
		u.lastIndex = len(u.lastMsg) - 1
	}
	return nil
}

func (u *ChatUi) arrowDown(g *gocui.Gui, iv *gocui.View) error {
	iv.Clear()

	if u.lastIndex >= len(u.lastMsg) {
		u.lastIndex = len(u.lastMsg) - 1
	}

	if u.lastIndex >= 0 && u.lastIndex <= len(u.lastMsg)-1 {
		_, _ = fmt.Fprintln(iv, u.lastMsg[u.lastIndex])
		_ = iv.SetCursor(len(u.lastMsg[u.lastIndex]), 0)
	}
	u.lastIndex--
	if u.lastIndex < 0 {
		u.lastIndex = 0
	}
	return nil
}

func (u *ChatUi) SpeakBroadcast(content string) {
	u.mu.Lock()
	defer u.mu.Unlock()
	historyView, err := u.gui.View(ViewTitleHistory)
	if err != nil {
		return
	}
	if historyView != nil {
		_, _ = fmt.Fprintln(historyView, content)
	}
	G.Update(func(gui *gocui.Gui) error {
		return nil
	})
}

func (u *ChatUi) OnSpeak(content string) {
	model.Speak(content)
}

func (u *ChatUi) ShowPlayerList(list []string) {
	u.mu.Lock()
	defer u.mu.Unlock()
	time.Sleep(time.Second)
	listPlayerView, err := u.gui.View(ViewTitleListPlayer)
	if err != nil {
		return
	}
	if listPlayerView != nil {
		listPlayerView.Clear()
		for _, v := range list {
			_, _ = fmt.Fprintln(listPlayerView, v)
		}
	}
	G.Update(func(gui *gocui.Gui) error {
		return nil
	})
}

func destroyChatUi(g *gocui.Gui) error {
	g.DeleteKeybindings(ViewTitleTitle)
	if err := g.DeleteView(ViewTitleTitle); err != nil {
		fmt.Println(err)
		return err
	}
	g.DeleteKeybindings(ViewTitleListPlayer)
	if err := g.DeleteView(ViewTitleListPlayer); err != nil {
		fmt.Println(err)
		return err
	}
	g.DeleteKeybindings(ViewTitleHistory)
	if err := g.DeleteView(ViewTitleHistory); err != nil {
		fmt.Println(err)
		return err
	}
	g.DeleteKeybindings(ViewTitleShortcut)
	if err := g.DeleteView(ViewTitleShortcut); err != nil {
		fmt.Println(err)
		return err
	}

	g.DeleteKeybindings(ViewTitleInput)
	if err := g.DeleteView(ViewTitleInput); err != nil {
		fmt.Println(err)
		return err
	}

	_, _ = g.SetCurrentView(ViewSelectTitle)

	return nil
}
