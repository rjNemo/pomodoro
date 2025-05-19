package main

import (
	"log"
	"os/exec"
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

			cmd := exec.Command("terminal-notifier",
				"-title", "Timer Complete",
				"-message", m.name,
				"-ignoreDnD",
			)
			if err := cmd.Run(); err != nil {
				log.Printf("error sending notification: %v", err)
			}
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
