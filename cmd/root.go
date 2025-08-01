package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	
	"github.com/tahatesser/fasting-cli/state"
	"github.com/tahatesser/fasting-cli/view"
)

type appModel struct {
	state *state.State
}

func initialModel() *appModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return &appModel{
		state: &state.State{
			Spinner: s,
		},
	}
}

func (m *appModel) Init() tea.Cmd {
	return tea.Batch(m.state.Init(), m.state.Spinner.Tick)
}

func (m *appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	newModel, cmd := m.state.Update(msg)
	if s, ok := newModel.(*state.State); ok {
		m.state = s
	}
	return m, cmd
}

func (m *appModel) View() string {
	return view.View(m.state.FastingSession, m.state.Spinner)
}

var rootCmd = &cobra.Command{
	Use:   "fasting-cli",
	Short: "A CLI for tracking intermittent fasting.",
	Long:  `A simple CLI application to help you track your intermittent fasting schedule.`,
	Run: func(cmd *cobra.Command, args []string) {
		// If no subcommand is given, run the TUI to show current status
		p := tea.NewProgram(initialModel())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	},
}

var startCmd = &cobra.Command{
	Use:   "start [duration]",
	Short: "Start a new fasting session",
	Long:  `Start a new fasting session with a specified duration (e.g., 16h, 18h30m). You can optionally specify how long ago the fast started using --ago, and a protocol name.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		durationStr := args[0]
		duration, err := time.ParseDuration(durationStr)
		if err != nil {
			fmt.Printf("Error parsing duration: %v\n", err)
			os.Exit(1)
		}

		agoStr, _ := cmd.Flags().GetString("ago")
		protocolName, _ := cmd.Flags().GetString("protocol")

		var startTime time.Time
		if agoStr != "" {
			agoDuration, err := time.ParseDuration(agoStr)
			if err != nil {
				fmt.Printf("Error parsing --ago duration: %v. Please use format like 2h, 30m, 1h30m\n", err)
				os.Exit(1)
			}
			startTime = time.Now().UTC().Add(-agoDuration)
		} else {
			startTime = time.Now().UTC()
		}

		// Safeguards
		maxDuration := 48 * time.Hour // Maximum allowed fasting duration
		if duration > maxDuration {
			fmt.Printf("Error: Fasting duration cannot exceed %s.\n", maxDuration)
			os.Exit(1)
		}

		// Calculate end time based on provided start time and duration
		endTime := startTime.Add(duration)
		if endTime.Before(time.Now().UTC()) {
			fmt.Println("Error: The calculated end time is in the past. Please adjust start time or duration.")
			os.Exit(1)
		}

		if protocolName == "" {
			protocolName = "Custom"
		}

		app := initialModel()
		app.state.StartFasting(startTime, duration, protocolName)
		fmt.Printf("Fasting session started for %s (Protocol: %s)!\n", duration.String(), protocolName)
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the current fasting session",
	Long:  `Stop the currently active fasting session.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load the current fasting session directly
		session, err := state.LoadFastingSession()
		if err != nil {
			fmt.Printf("Error loading session: %v\n", err)
			os.Exit(1)
		}

		if session.IsActive {
			session.IsActive = false
			_ = state.SaveFastingSession(session)
			fmt.Println("Fasting session stopped.")
		} else {
			fmt.Println("No active fasting session to stop.")
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().String("ago", "", "Optional: Specify how long ago the fast started (e.g., 2h, 30m, 1h30m).")
	startCmd.Flags().String("protocol", "", "Optional: Specify a fasting protocol name (e.g., 16-8, 18-6)")
	rootCmd.AddCommand(stopCmd)
}

// Execute adds all child commands to the root command and sets flags appropriately.

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
