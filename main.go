package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Theme struct {
	Name        string
	Correct     string
	Wrong       string
	Pending     string
	HeaderBg    string
	HeaderFg    string
	BorderColor string
}

var themes = []Theme{
	{
		Name:        "Tokyo Night (Dark & Moodish)",
		Correct:     "#9ece6a",
		Wrong:       "#f7768e",
		Pending:     "#565f89",
		HeaderBg:    "#24283b",
		HeaderFg:    "#c0caf5",
		BorderColor: "#414868",
	},
	{
		Name:        "Cyberpunk 2077 (Neon High-Contrast)",
		Correct:     "#00f0ff",
		Wrong:       "#ff0055",
		Pending:     "#5a5a00",
		HeaderBg:    "#fcee0a",
		HeaderFg:    "#000000",
		BorderColor: "#fcee0a",
	},
	{
		Name:        "Matrix (Digital Terminal Minimalist)",
		Correct:     "#00ff00",
		Wrong:       "#ff0000",
		Pending:     "#003300",
		HeaderBg:    "#000000",
		HeaderFg:    "#00ff00",
		BorderColor: "#00ff00",
	},
	{
		Name:        "Gruvbox Retro (Warm, Earthy & Sepia)",
		Correct:     "#b8bb26",
		Wrong:       "#fb4934",
		Pending:     "#928374",
		HeaderBg:    "#3c3836",
		HeaderFg:    "#ebdbb2",
		BorderColor: "#7c6f64",
	},
	{
		Name:        "Paperback (Crisp Light Mode)",
		Correct:     "#22863a",
		Wrong:       "#cb2431",
		Pending:     "#6a737d",
		HeaderBg:    "#e1e4e8",
		HeaderFg:    "#24292e",
		BorderColor: "#d1d5da",
	},
	{
		Name:        "Synthwave '84 (Rad Outrun Sunset)",
		Correct:     "#39ff14",
		Wrong:       "#ff007f",
		Pending:     "#711c91",
		HeaderBg:    "#2b013a",
		HeaderFg:    "#00ffff",
		BorderColor: "#ff007f",
	},
	{
		Name:        "Just Black",
		Correct:     "#111",
		Wrong:       "#111",
		Pending:     "#111",
		HeaderBg:    "#111",
		HeaderFg:    "#fff",
		BorderColor: "#111",
	},
}

type screen int

const (
	screenTyping screen = iota
	screenTheme
	ScreenMenu
)

type model struct {
	prompt        string
	promptArray   []rune
	userInput     []rune
	duration      int
	startedAt     time.Time
	finishedAt    time.Time
	elapsed       time.Duration
	cursorPos     int
	currentCount  int
	done          bool
	typing        bool
	err           error
	wpm           float64
	width         int
	height        int
	currentScreen screen
	selectedTheme int
}

func initialModel() model {
	data := "Typing quickly is not just about moving your fingers fast; it's about rhythm, accuracy, and focus. The more you practice, the more your hands learn the flow of the keyboard, almost like a musician mastering an instrument. Stay relaxed, keep your eyes on the screen, and let your thoughts guide your movements. Every mistake is a small lesson, reminding you to slow down, breathe, and improve with each word you type."

	return model{
		prompt:        data,
		promptArray:   []rune(data),
		userInput:     []rune{},
		duration:      30,
		startedAt:     time.Time{},
		finishedAt:    time.Time{},
		elapsed:       0,
		currentCount:  0,
		done:          false,
		typing:        false,
		err:           nil,
		currentScreen: screenTyping,
		selectedTheme: 0,
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
			if m.currentScreen != screenTyping {
				m.currentScreen = screenTyping
				return m, nil
			}
			m = initialModel()
			return m, tea.Quit
		case "ctrl+r":
			m = m.restart_func()
			return m, nil
		case "ctrl+t":
			if !m.typing {
				m.currentScreen = screenTheme
				return m, nil
			}
		}
		// theme screen control
		if m.currentScreen == screenTheme {
			switch msg.String() {
			case "up", "k":
				if m.selectedTheme > 0 {
					m.selectedTheme--
				} else if m.selectedTheme == 0 {
					m.selectedTheme = len(themes) - 1
				}
			case "down", "j":
				if m.selectedTheme < len(themes)-1 {
					m.selectedTheme++
				} else if m.selectedTheme == len(themes)-1 {
					m.selectedTheme = 0
				}
			case "enter":
				m.currentScreen = screenTyping
			}
			return m, nil
		}
		// typing screen control
		switch msg.String() {
		case "backspace":
			if len(m.userInput) > 0 {
				m.userInput = m.userInput[:len(m.userInput)-1]
				if m.cursorPos > 0 {
					m.cursorPos--
				}
			}
		default:
			if len(msg.Runes) == 1 {
				r := msg.Runes[0]
				if r < 32 {
					return m, nil
				}
			}
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
	rm.selectedTheme = m.selectedTheme
	return rm
}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m model) View() string {

	currentTheme := themes[m.selectedTheme]
	var (
		correctStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(currentTheme.Correct))
		wrongStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color(currentTheme.Wrong))
		pendingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(currentTheme.Pending))
	)

	var style_header = lipgloss.NewStyle().
		Width(m.width-4).
		BorderStyle(lipgloss.NormalBorder()).
		Bold(true).
		Align(lipgloss.Center, lipgloss.Center).
		Foreground(lipgloss.Color(currentTheme.HeaderFg)).
		Background(lipgloss.Color(currentTheme.HeaderBg))

	var style_border = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(currentTheme.BorderColor))

	var cursor_style = lipgloss.NewStyle().Bold(true)

	var wpmStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(currentTheme.Correct)).Bold(true)

	var style_button = lipgloss.NewStyle().
		Background(lipgloss.Color(currentTheme.BorderColor)).
		Foreground(lipgloss.Color(currentTheme.HeaderFg)).
		Padding(0, 1).
		Bold(true)

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

	// Render theme screen
	if m.currentScreen == screenTheme {
		var b strings.Builder
		b.WriteString(style_header.Render("Color Themes"))
		b.WriteString("\n\n")
		b.WriteString("Use ↑/↓ or j/k to navigate, Enter to select theme\n\n")
		for i, t := range themes {
			cursor := " "
			if i == m.selectedTheme {
				cursor = ">"
				b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(t.HeaderFg)).Bold(true).Render(fmt.Sprintf("%s %s", cursor, t.Name)))
				b.WriteString("\n")
			} else {
				fmt.Fprintf(&b, "%s %s\n", cursor, t.Name)
			}
		}
		return fullscreenStyle.Render(style_border.Render(b.String()))
	}

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
	var themeBtn string
	if !m.typing && !m.done {
		themeBtn = style_button.Render("[Ctrl+T] Themes")
		s += "\n" + lipgloss.NewStyle().Width(m.width-4).Align(lipgloss.Right).Render(themeBtn) + "\n"
		s += "\n" + align_center.Render("To start typing press start") + "\n"
		s += text_warp.Render(quote.String())
	} else if !m.done && m.typing {
		s += "\n"
		s += text_warp.Render(quote.String())
		s += fmt.Sprintf("\n current user input count %v", m.currentCount)
		s += fmt.Sprintf("\n %v time left", m.duration)
	} else if m.done {
		s += "\nTest Finished"
		s += wpmStyle.Render(fmt.Sprintf("\n Your Result is %v WPM", m.wpm))
		s += "\n To restart press \"ctrl+r\""
	}
	return fullscreenStyle.Render(style_border.Render(s))
}

func main() {
	p := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Printf("error: %v", err)
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
