package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	prompt       string
	promptArray  []rune
	userInput    []rune
	duration     int
	startedAt    time.Time
	finishedAt   time.Time
	elapsed      time.Duration
	cursorPos    int
	currentCount int
	done         bool
	typing       bool
	err          error
	wpm          float64
	width        int
	height       int
}

func initialModel() model {
	data := "Typing quickly is not just about moving your fingers fast; it's about rhythm, accuracy, and focus. The more you practice, the more your hands learn the flow of the keyboard, almost like a musician mastering an instrument. Stay relaxed, keep your eyes on the screen, and let your thoughts guide your movements. Every mistake is a small lesson, reminding you to slow down, breathe, and improve with each word you type."

	return model{
		prompt:       data,
		promptArray:  []rune(data),
		userInput:    []rune{},
		duration:     30,
		startedAt:    time.Time{},
		finishedAt:   time.Time{},
		elapsed:      0,
		currentCount: 0,
		done:         false,
		typing:       false,
		err:          nil,
	}
}

type tickMsg struct{}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m = initialModel()
			return m, tea.Quit
		case "backspace":
			if len(m.userInput) > 0 {
				m.userInput = m.userInput[:len(m.userInput)-1]
				if m.cursorPos > 0 {
					m.cursorPos--
				}
			}
		case "ctrl+r":
			m = m.restart_func()
		default:
			if !m.done && !m.typing {
				m.startedAt = time.Now()
				m.typing = true
				if len(msg.String()) == 1 {
					m.userInput = append(m.userInput, []rune(msg.String())...)
					m.cursorPos++
					if m.promptArray[len(m.userInput)-1] == m.userInput[len(m.userInput)-1] {
						m.currentCount++
					}
				}

				return m, tick()
			} else if m.typing && !m.done {
				if len(msg.String()) == 1 {
					m.userInput = append(m.userInput, []rune(msg.String())...)
					m.cursorPos++
					if m.promptArray[len(m.userInput)-1] == m.userInput[len(m.userInput)-1] {
						m.currentCount++
					}
				}
			} else {
				return m, nil
			}
		}

	case tickMsg:
		if m.done {
			return m, nil
		}
		m.duration--
		if m.duration <= 0 {
			m.done = true
			m.finishedAt = time.Now()
			m.elapsed = time.Since(m.startedAt)
			m.wpm = calculateWPM(m.currentCount, m.duration, m.elapsed)
			return m, nil
		}
		return m, tick()
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m model) restart_func() model {
	width := m.width
	height := m.height
	rm := initialModel()
	rm.width = width
	rm.height = height
	return rm
}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m model) View() string {

	var (
		correctStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
		wrongStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555"))
		pendingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#777777"))
	)

	var style_header = lipgloss.NewStyle().
		Width(m.width-4).
		BorderStyle(lipgloss.NormalBorder()).
		Bold(true).
		Align(lipgloss.Center, lipgloss.Center).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#331886ff"))

	var align_center = lipgloss.NewStyle().
		Width(m.width-4).
		Align(lipgloss.Center, lipgloss.Center)

	fullscreenStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center)

	text_warp := lipgloss.NewStyle().Width(m.width - 4)

	// if m.width == 0 || m.height == 0 {
	// 	return "Loading…"
	// }

	var quote strings.Builder
	for i, r := range m.promptArray {
		if !m.done && i == m.cursorPos {
			quote.WriteString(cursor_style.Render("|"))
		}

		if i < len(m.userInput) {

			if m.userInput[i] == r {
				quote.WriteString(correctStyle.Render(string(r)))
			} else {
				quote.WriteString(wrongStyle.Render(string(r)))
			}
		} else {
			quote.WriteString(pendingStyle.Render(string(r)))
		}
	}

	s := style_header.Render("Typing Test\nMessauere your WPM\n")
	if !m.typing && !m.done {
		s += "\n" + align_center.Render("To start typing press start") + "\n"
		s += text_warp.Render(quote.String())
	} else if !m.done && m.typing {
		s += "\n"
		s += text_warp.Render(quote.String())
		s += fmt.Sprintf("\n current user input count %v", m.currentCount)
		s += fmt.Sprintf("\n %v time left", m.duration)
	} else if m.done {
		s += fmt.Sprintf("\nTest Finished")
		s += fmt.Sprintf("\n Your Result is %v WPM", m.wpm)
		s += fmt.Sprintf("\n To restart press \"ctrl+r\"")
	}
	return fullscreenStyle.Render(style_border.Render(s))
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func calculateWPM(currentCount int, duration int, elapsed time.Duration) float64 {
	minutesPassed := elapsed.Minutes()
	if minutesPassed == 0 && duration == 0 {
		return 0
	}
	return float64(currentCount) / 5.0 / minutesPassed
}

// styles

var style_border = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder())

var cursor_style = lipgloss.NewStyle().Bold(true)
