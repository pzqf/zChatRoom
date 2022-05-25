package cui

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

type DialogUi struct {
	gui        *gocui.Gui
	dialogType DialogType
}

const (
	MaxDialogWidth = 140
	MinDialogWidth = 40

	ViewDialogTitle = "dialog"
)

type DialogType int

const (
	DialogTypeInfo  = DialogType(1)
	DialogTypeWarn  = DialogType(2)
	DialogTypeError = DialogType(3)
)

func NewDialog(g *gocui.Gui, dt DialogType) *DialogUi {
	if dt < 0 || dt > 3 {
		dt = DialogTypeInfo
	}

	return &DialogUi{
		gui:        g,
		dialogType: dt,
	}
}

func (u *DialogUi) Show(content string) error {
	tw, th := u.gui.Size()

	dialogWidth := len(content) + 4
	if dialogWidth > MaxDialogWidth/2 {
		dialogWidth = MaxDialogWidth / 2
	}
	if dialogWidth < MinDialogWidth {
		dialogWidth = MinDialogWidth
	}

	v, err := u.gui.SetView(ViewDialogTitle, tw/2-dialogWidth/2, th/2-3, tw/2+dialogWidth/2, th/2+2)
	if err != nil && err != gocui.ErrUnknownView {
		log.Println("Failed to create tasks view:", err)
		return err
	}

	v.FgColor = gocui.ColorCyan
	switch u.dialogType {
	case DialogTypeInfo:
		v.Title = "info"
	case DialogTypeWarn:
		v.Title = "warning"
		v.FgColor = gocui.ColorYellow
	case DialogTypeError:
		v.Title = "error"
		v.FgColor = gocui.ColorRed
	default:
		v.Title = "info"
	}

	//v.SelBgColor = gocui.ColorBlue
	//v.SelFgColor = gocui.ColorBlack

	//v.Editable = true
	u.gui.Cursor = false

	if len(content) >= MaxDialogWidth/2-4 {
		content = content[0:MaxDialogWidth/2-8] + "..."
	}

	//show content
	c := (dialogWidth - len(content)) / 2
	str := "\n\n"
	for i := 0; i < c; i++ {
		str += " "
	}
	_, _ = fmt.Fprintln(v, str+content)

	//show ok
	str = "\n"
	for i := 0; i < dialogWidth/2-5; i++ {
		str += " "
	}
	_, _ = fmt.Fprintln(v, str+"[enter]")

	if _, err = u.gui.SetCurrentView(ViewDialogTitle); err != nil {
		return err
	}
	if err = u.gui.SetKeybinding(ViewDialogTitle, gocui.KeyEnter, gocui.ModNone, u.close); err != nil {
		return err
	}

	u.gui.Update(func(gui *gocui.Gui) error {
		return nil
	})

	return nil
}

func (u *DialogUi) Close() error {
	_ = u.close(u.gui, nil)
	return nil
}

func (u *DialogUi) close(g *gocui.Gui, iv *gocui.View) error {
	g.DeleteKeybindings(ViewDialogTitle)
	if err := g.DeleteView(ViewDialogTitle); err != nil {
		fmt.Println(err)
		return err
	}

	_, err := g.SetCurrentView(LastView)
	if err != nil {
		return err
	}

	return nil
}
