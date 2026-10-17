package ui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
)

// Run starts the interactive program. root is the product directory.
func Run(root string) error {
	app, err := newApp(root, 0)
	if err != nil {
		return err
	}
	_, err = tea.NewProgram(Model{app: app}, tea.WithFilter(ignoreInterrupt)).Run()
	return err
}

// Model is the Bubble Tea program. It draws only the current frame.
type Model struct {
	app *app
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(func() tea.Msg { return tea.RequestWindowSize() }, m.app.pending)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if m.app.width == 0 && msg.Width > 0 {
			if err := m.app.relayout(msg.Width); err != nil {
				m.app.status = err.Error()
			}
		}
		return m, nil
	case tea.KeyboardEnhancementsMsg:
		m.app.keyMode = keysEnhanced
		if m.app.status == noShiftMsg {
			m.app.status = "Prêt"
		}
		return m, nil
	case tea.PasteMsg:
		if m.app.running {
			return m, nil
		}
		m.app.paste(msg.Content)
		return m, nil
	case spinMsg:
		if !m.app.running || msg.gen != m.app.runGen {
			return m, nil
		}
		m.app.spin = (m.app.spin + 1) % len(spinFrames)
		m.app.status = spinFrames[m.app.spin] + " En cours : " + m.app.cmdLabel
		return m, spinCmd(msg.gen)
	case cmdDone:
		if msg.gen != m.app.runGen {
			return m, nil
		}
		m.app.running = false
		m.app.cancel = nil
		if msg.cancelled {
			m.app.status = "Annulé"
			return m, nil
		}
		if msg.zone != "" {
			k := key(m.app.winID(), msg.zone)
			m.app.val[k] = strings.TrimRight(msg.out, "\n")
			m.app.scroll[k] = 0
		}
		m.app.status = "Terminé " + itoa(msg.code)
		return m, nil
	case listSpin:
		if msg.gen != m.app.listGen || m.app.listCancel == nil {
			return m, nil
		}
		m.app.spin = (m.app.spin + 1) % len(spinFrames)
		m.app.status = spinFrames[m.app.spin] + " En cours : " + m.app.listLabel
		return m, listSpinCmd(msg.gen)
	case listDone:
		m.app.onListDone(msg)
		return m, nil
	case tea.KeyPressMsg:
		return m, m.app.onKey(msg)
	default:
		return m, nil
	}
}

func (m Model) View() tea.View {
	v := tea.NewView(m.app.frame())
	v.AltScreen = true
	if x, y, ok := m.app.cursor(); ok {
		v.Cursor = tea.NewCursor(x, y)
	}
	return v
}

func (a *app) onKey(msg tea.KeyPressMsg) tea.Cmd {
	if a.edit != nil && cancelEditKey(msg) {
		a.cancelEdit()
		return nil
	}
	if quitKey(msg) {
		if a.running {
			a.stopRun()
		}
		return func() tea.Msg { return tea.Quit() }
	}
	name := msg.String()
	if shiftEnter(msg) {
		name = "shift+enter"
		a.keyMode = keysEnhanced
		if a.status == noShiftMsg {
			a.status = "Prêt"
		}
	} else if saveEditKey(msg) {
		name = "ctrl+s"
	} else if a.keyMode == keysUnknown {
		a.keyMode = keysPlain
	}
	if a.running {
		if name == "esc" {
			a.stopRun()
			a.status = "Annulé"
		}
		return nil
	}
	if a.edit != nil {
		return a.editKey(name)
	}
	switch name {
	case "left":
		if err := a.moveFocus(-1); err != nil {
			a.status = err.Error()
		}
	case "right":
		if err := a.moveFocus(1); err != nil {
			a.status = err.Error()
		}
	case "up":
		if err := a.moveVertical(-1); err != nil {
			a.status = err.Error()
		}
	case "down":
		if err := a.moveVertical(1); err != nil {
			a.status = err.Error()
		}
	case "enter":
		act, _ := a.action("enter")
		cmd, err := a.perform(act)
		if err != nil {
			a.status = err.Error()
		}
		return cmd
	case "esc":
		if a.listCancel != nil {
			a.stopList()
			a.status = "Annulé"
			return nil
		}
		if err := a.goBack(); err != nil {
			a.status = err.Error()
		}
		return a.takePending()
	}
	return nil
}

func ignoreInterrupt(_ tea.Model, msg tea.Msg) tea.Msg {
	if _, ok := msg.(tea.InterruptMsg); ok {
		return nil
	}
	return msg
}

func quitKey(msg tea.KeyPressMsg) bool {
	if msg.String() == "ctrl+q" || msg.Keystroke() == "ctrl+q" {
		return true
	}
	k := msg.Key()
	return k.Mod.Contains(tea.ModCtrl) && (k.Code == 'q' || k.Code == 'Q')
}

func cancelEditKey(msg tea.KeyPressMsg) bool {
	if msg.String() == "ctrl+c" || msg.String() == "\x03" || msg.Keystroke() == "ctrl+c" {
		return true
	}
	k := msg.Key()
	return k.Mod.Contains(tea.ModCtrl) && (k.Code == 'c' || k.Code == 'C')
}

func saveEditKey(msg tea.KeyPressMsg) bool {
	if msg.String() == "ctrl+s" || msg.Keystroke() == "ctrl+s" {
		return true
	}
	k := msg.Key()
	return k.Mod.Contains(tea.ModCtrl) && (k.Code == 's' || k.Code == 'S')
}

func shiftEnter(msg tea.KeyPressMsg) bool {
	if msg.String() == "shift+enter" || msg.Keystroke() == "shift+enter" {
		return true
	}
	k := msg.Key()
	if !k.Mod.Contains(tea.ModShift) {
		return false
	}
	switch k.Code {
	case tea.KeyEnter, tea.KeyKpEnter:
		return true
	default:
		return false
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b [16]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		b[i] = '-'
	}
	return string(b[i:])
}
