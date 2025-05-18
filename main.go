package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	padding  = 2
	maxWidth = 80
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

type (
	model struct {
		name        string
		totalTime   time.Duration
		elapsedTime time.Duration
		progress    progress.Model
	}

	TickMsg time.Time
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("error running the Pomodoro timer %v", err)
	}
}

func run() error {
	name := flag.String("name", "Pomodoro", "Name of the timer")
	durationString := flag.String("duration", "25m", "Duration of the timer")
	flag.Parse()

	duration, err := time.ParseDuration(*durationString)
	if err != nil {
		return fmt.Errorf("could not parse duration %s: %w", *durationString, err)
	}

	p := tea.NewProgram(initialModel(*name, duration))
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("could not start the timer app %w", err)
	}
	return nil
}

func initialModel(name string, duration time.Duration) model {
	return model{
		name:      name,
		totalTime: duration,
		progress:  progress.New(progress.WithDefaultGradient()),
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.progress.Width = min(msg.Width-padding*2-4, maxWidth)
		return m, nil
	case TickMsg:
		if m.progress.Percent() == 1.0 {
			return m, tea.Quit
		}
		m.elapsedTime += time.Second
		cmd := m.progress.SetPercent(float64(m.elapsedTime) / float64(m.totalTime))

		return m, tea.Batch(tickCmd(), cmd)

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	}
	return m, nil
}

func (m model) View() string {
	pad := strings.Repeat(" ", padding)
	return "\n" +
		pad + m.name + " for " + m.totalTime.String() + "\n\n" +
		pad + m.elapsedTime.String() + "\n\n" +
		pad + m.progress.View() + "\n\n" +
		pad + helpStyle("Press Q to quit")
}
