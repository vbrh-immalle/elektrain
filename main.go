package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

var db ElekTrainDb

type tickMsg struct{}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

type TerminalSize struct {
	Width  int
	Height int
}

type model struct {
	CurrentTerminalSize TerminalSize
	Textinput           textinput.Model
	HuidigeOpgave       Opgave
	State               string
	Ticks               int
}

func initialModel() model {
	ti := textinput.NewModel()
	ti.Placeholder = "Antwoord"
	ti.Focus()
	ti.CharLimit = 15
	ti.Width = 15

	opg, hasNext := db.RandomNogJuistTeBeantwoordenOpgave()
	if !hasNext {
		log.Fatalf("We hebben geen vragen!")
	}

	m := model{
		Textinput:     ti,
		HuidigeOpgave: opg,
		State:         "entering answer",
		Ticks:         0,
	}

	return m
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, tea.EnterAltScreen, tick())
}

func (m model) handleWindowSizeMsg(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.CurrentTerminalSize.Width = msg.Width
	m.CurrentTerminalSize.Height = msg.Height

	return m, nil
}

func (m model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		return m, tea.Quit
	case tea.KeyEnter:
		switch m.State {
		case "entering answer":
			if m.HuidigeOpgave.Antwoord == m.Textinput.Value() {
				m.State = "answer correct"
				m.HuidigeOpgave.AantalJuisteAntwoorden += 1
			} else {
				m.State = "answer wrong"
				m.HuidigeOpgave.AantalFouteAntwoorden += 1
			}
			db.UpdateOpgave(m.HuidigeOpgave)
			return m, nil
		case "answer correct":
			fallthrough
		case "answer wrong":
			var hasNext bool
			m.HuidigeOpgave, hasNext = db.RandomNogJuistTeBeantwoordenOpgave()
			if hasNext {
				m.State = "entering answer"
			} else {

				m.State = "finished"
			}
			m.Textinput.Reset()
			return m, nil
		case "finished":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.Textinput, cmd = m.Textinput.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleWindowSizeMsg(msg)
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)
	case tickMsg:
		m.Ticks += 1
		if m.State != "finished" {
			return m, tick()
		}
		return m, nil
	}

	return m, cmd
}

func (m model) View() string {
	s := ""

	if m.CurrentTerminalSize.Width < 80 || m.CurrentTerminalSize.Height < 20 {
		s += wordwrap.String("Maak je terminal groter! (ESC om te stoppen)", m.CurrentTerminalSize.Width)
		return s
	}

	title := "ElekTrain"
	headerStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder(), true, true).
		Align(lipgloss.Center).
		Height(1).
		//Padding(1).
		Width(len(title) + 3*2).
		MarginLeft(m.CurrentTerminalSize.Width/2 - len(title) + 3).
		Foreground(lipgloss.Color("#00CC00"))
	s += headerStyle.Render(title)

	variableStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00CC00")).
		PaddingLeft(3).
		PaddingRight(1).
		Width(6)

	vraagBoxMarginLeft := m.CurrentTerminalSize.Width / 8
	vraagBoxPaddingHor := 5
	vraagBoxWidth := m.CurrentTerminalSize.Width - vraagBoxMarginLeft*2
	vraagBoxTextWidth := vraagBoxWidth - 2*vraagBoxPaddingHor

	vraagStyle := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder(), true, true).
		PaddingLeft(vraagBoxPaddingHor).
		PaddingRight(vraagBoxPaddingHor).
		PaddingTop(1).
		PaddingBottom(1).
		Align(lipgloss.Center).
		MarginLeft(vraagBoxMarginLeft).
		Foreground(lipgloss.Color("#999999"))

	antwoordBoxMarginLeft := vraagBoxMarginLeft

	antwoordBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder(), true, true).
		PaddingLeft(5).
		PaddingRight(5).
		Width(25).
		MarginLeft(antwoordBoxMarginLeft).
		Foreground(lipgloss.Color("#FFFFFF"))

	foutStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#CC0000"))

	s += "\n\n"
	s += "                  Seconden bezig: ["
	s += variableStyle.Render(strconv.Itoa(m.Ticks))
	s += "]\n"

	s += "Totaal aantal gegeven antwoorden: ["
	s += variableStyle.Render(strconv.Itoa(db.CountTotaalAantalAntwoorden()))
	s += "]\n"

	s += "                  Aantal opgaven: ["
	s += variableStyle.Render(strconv.Itoa(db.Count()))
	s += "]\n"

	s += "       Nog juist te beantwoorden: ["
	s += variableStyle.Render(strconv.Itoa(db.CountOpgavesTeGaan()))
	s += "]\n"

	s += "        Aantal juiste antwoorden: ["
	s += variableStyle.Render(strconv.Itoa(db.CountJuistBeantwoord()))
	s += "]\n"

	if m.State != "finished" {
		s += "\n"
		s += fmt.Sprintln(vraagStyle.Render(wordwrap.String(m.HuidigeOpgave.Vraag, vraagBoxTextWidth)))
	}

	switch m.State {
	case "entering answer":
		s += antwoordBoxStyle.Render(m.Textinput.View())
		s += "\n\nDruk"
		s += variableStyle.Render("ESCAPE")
		s += " om te stoppen."
	case "answer correct":
		s += "\n"
		s += variableStyle.Render("JUIST!")
		s += "\n"
		s += fmt.Sprintf("\nDruk [%s] om verder te gaan.\n", variableStyle.Render("ENTER"))
	case "answer wrong":
		s += "\n"
		s += foutStyle.Render("FOUT!")
		s += "\n"
		s += "Het juiste antwoord was ["
		s += variableStyle.Render(m.HuidigeOpgave.Antwoord)
		s += "].\n"
		s += fmt.Sprintf("\nDruk [%s] om verder te gaan.\n", variableStyle.Render("ENTER"))
	case "finished":
		s += "\nJe kan het! Je kan voor de zekerheid natuurlijk altijd nog eens proberen ;-)\n"
		s += "\nDruk "
		s += variableStyle.Render("ENTER")
		s += " om terug te keren naar de shell."
	}

	return s
}

func main() {
	db = &HardCodedDb{}
	db.Init()

	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, error: %v\n", err)
		os.Exit(1)
	}
}
