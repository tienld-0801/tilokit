package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	padding  = 2
	maxWidth = 80
)

var (
	titleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		MarginBottom(1)
)

type ProgressModel struct {
	progress progress.Model
	title    string
	status   string
	steps    []string
	current  int
	done     bool
}

type ProgressMsg struct {
	Step     string
	Progress float64
	Done     bool
}

func NewProgressModel(title string, steps []string) ProgressModel {
	return ProgressModel{
		progress: progress.New(progress.WithDefaultGradient()),
		title:    title,
		steps:    steps,
		current:  0,
		done:     false,
	}
}

func (m ProgressModel) Init() tea.Cmd {
	return nil
}

func (m ProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		w := msg.Width - padding*2 - 4
		if w < 10 {
			w = 10
		}
		if w > maxWidth {
			w = maxWidth
		}
		m.progress.Width = w
		return m, nil

	case ProgressMsg:
		m.status = msg.Step
		// Advance current step based on the incoming label
		for i, s := range m.steps {
			if s == msg.Step {
				if i >= m.current {
					m.current = i
				}
				break
			}
		}
		if msg.Done {
			m.done = true
			m.current = len(m.steps) // mark all as completed in the view
			m.progress.SetPercent(1.0)
			return m, tea.Quit
		}
		// Keep below 100% until Done to avoid "stuck" full bar
		percent := msg.Progress
		if percent > 0.99 {
			percent = 0.99
		}
		cmd := m.progress.SetPercent(percent)
		return m, cmd

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	default:
		return m, nil
	}
	return m, nil
}

func (m ProgressModel) View() string {
	if m.done {
		return ""
	}

	pad := strings.Repeat(" ", padding)

	var stepsList strings.Builder
	for i, step := range m.steps {
		if i < m.current {
			stepsList.WriteString(fmt.Sprintf("%s✅ %s\n", pad, step))
		} else if i == m.current {
			stepsList.WriteString(fmt.Sprintf("%s🔄 %s\n", pad, step))
		} else {
			stepsList.WriteString(fmt.Sprintf("%s⏳ %s\n", pad, step))
		}
	}

	return fmt.Sprintf("\n%s\n\n%s%s\n\n%s\n",
		titleStyle.Render(m.title),
		pad, m.progress.View(),
		stepsList.String())
}

// RunProjectProgress starts the animated progress UI and returns a done channel
// that is closed when the TUI program exits. Callers can wait on it instead of sleeping.
func RunProjectProgress(title string, steps []string, progressChan <-chan ProgressMsg) (<-chan struct{}, error) {
	model := NewProgressModel(title, steps)
	
	p := tea.NewProgram(model)
	done := make(chan struct{})
	
	go func() {
		defer close(done)
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error running progress: %v\n", err)
		}
	}()

	// Listen for progress updates
	go func() {
		for msg := range progressChan {
			p.Send(msg)
			if msg.Done {
				p.Quit()
				break
			}
		}
	}()

	return done, nil
}
