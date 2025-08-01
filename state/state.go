package state

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/tahatesser/fasting-cli/model"
)

const (
	configDirName  = "fasting-cli"
	sessionFileName = "session.json"
)

type State struct {
	FastingSession model.FastingSession
	Spinner        spinner.Model
}

func (s *State) Init() tea.Cmd {
	// Load the fasting session on startup
	session, err := LoadFastingSession()
	if err != nil {
		// Handle error, maybe log it or show a message to the user
		// For now, we'll just proceed with a default empty session
		_ = err // Suppress unused error warning
	}
	s.FastingSession = session

	return tea.Batch(s.Spinner.Tick, tick())
}

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (s *State) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		if s.FastingSession.IsActive {
			elapsed := time.Since(s.FastingSession.StartTime)
			if elapsed >= s.FastingSession.Duration {
				s.StopFasting() // Fasting session completed
			} else {
				// Save session periodically to ensure progress is saved
				_ = SaveFastingSession(s.FastingSession)
			}
			return s, tick()
		}
		return s, nil
	case spinner.TickMsg:
		var cmd tea.Cmd
		s.Spinner, cmd = s.Spinner.Update(msg)
		return s, cmd

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			// Save the current session before quitting
			_ = SaveFastingSession(s.FastingSession)
			return s, tea.Quit
		case "s": // Start fasting (handled by CLI command, not TUI)
			// This case is primarily for testing or if a default start is desired without CLI args
			// If you want to enable starting from TUI, you'd need to add input for duration/protocol
			// For now, we'll just prevent starting if already active.
			if !s.FastingSession.IsActive {
				// Optionally, you could add a default start here for quick testing:
				// s.StartFasting(time.Now().UTC(), 16*time.Hour, "16-8")
				// return s, tick()
			}
		case "p": // Stop fasting
			if s.FastingSession.IsActive {
				s.StopFasting()
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return s, nil
}

func (s *State) StartFasting(startTime time.Time, duration time.Duration, protocolName string) {
	s.FastingSession = model.FastingSession{
		StartTime:    startTime,
		Duration:     duration,
		ProtocolName: protocolName,
		IsActive:     true,
	}
	_ = SaveFastingSession(s.FastingSession)
}

func (s *State) StopFasting() {
	s.FastingSession.IsActive = false
	_ = SaveFastingSession(s.FastingSession)
}

func (s *State) View() string {
	return ""
}

// GetConfigDirPath returns the path to the application's configuration directory.
func GetConfigDirPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".config", configDirName), nil
}

// SaveFastingSession saves the current fasting session to a file.
func SaveFastingSession(session model.FastingSession) error {
	configPath, err := GetConfigDirPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(configPath, 0755); err != nil {
		return err
	}

	filePath := filepath.Join(configPath, sessionFileName)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // For pretty printing
	return encoder.Encode(session)
}

// LoadFastingSession loads the fasting session from a file.
func LoadFastingSession() (model.FastingSession, error) {
	configPath, err := GetConfigDirPath()
	if err != nil {
		return model.FastingSession{}, err
	}

	filePath := filepath.Join(configPath, sessionFileName)
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// No session file found, return a default empty session
			return model.FastingSession{
				StartTime: time.Time{},
				Duration:  0,
				IsActive:  false,
			}, nil
		}
		return model.FastingSession{}, err
	}
	defer file.Close()

	var session model.FastingSession
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&session)
	if err != nil {
		return model.FastingSession{}, err
	}

	return session, nil
}
