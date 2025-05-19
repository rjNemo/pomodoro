package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

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
