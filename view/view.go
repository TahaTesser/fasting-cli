package view

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/tahatesser/fasting-cli/model"
)

// Styles
var (
	// Main container style
	containerStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 2).
		Width(60)

	// Header style
	headerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("6")).
		Background(lipgloss.Color("236")).
		Padding(0, 1).
		Bold(true).
		Align(lipgloss.Center)

	// Status style (e.g., "Fasting" or "Not fasting")
	statusStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("10")). // Green for active, red for inactive
		Bold(true)

	// Time display style
	timeStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("12")). // Blue
		Bold(true)

	// Progress bar styles
	progressBarWidth = 40
	progressBarDone  = lipgloss.NewStyle().Background(lipgloss.Color("6")).Foreground(lipgloss.Color("6")).Render(" ") // Green
	progressBarLeft  = lipgloss.NewStyle().Background(lipgloss.Color("8")).Foreground(lipgloss.Color("8")).Render(" ") // Dark gray

	// Percentage style
	percentageStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("7")). // Light gray
		Align(lipgloss.Right)
)

func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
}

func View(session model.FastingSession, spinner spinner.Model) string {
	if session.IsActive {
		elapsed := time.Since(session.StartTime)
		remaining := session.Duration - elapsed
		endTime := session.StartTime.Add(session.Duration)

		if remaining < 0 {
			remaining = 0
		}

		progress := float64(elapsed) / float64(session.Duration)
		if progress > 1.0 { // Cap progress at 100%
			progress = 1.0
		}

		doneBars := int(float64(progressBarWidth) * progress)
		leftBars := progressBarWidth - doneBars

		progressBar := lipgloss.JoinHorizontal(
			lipgloss.Top,
			lipgloss.NewStyle().Width(doneBars).SetString(progressBarDone).String(),
			lipgloss.NewStyle().Width(leftBars).SetString(progressBarLeft).String(),
		)

		// Build the content
		statusText := statusStyle.Foreground(lipgloss.Color("10")).Render("Fasting")
		remainingText := timeStyle.Render(formatDuration(remaining) + " remaining")
		endTimeText := timeStyle.Render("Ends: " + endTime.Format("Mon, 02 Jan 2006 15:04:05 MST"))
		percentageText := percentageStyle.Render(fmt.Sprintf("%.1f%%", progress*100))

		// Combine elements
		statusLine := lipgloss.JoinHorizontal(lipgloss.Left, spinner.View(), " ", statusText)
		timeInfo := lipgloss.JoinVertical(lipgloss.Left, remainingText, endTimeText)

		content := lipgloss.JoinVertical(
			lipgloss.Left,
			statusLine,
			"", // Spacer
			timeInfo,
			"", // Spacer
			lipgloss.JoinHorizontal(lipgloss.Left, progressBar, " ", percentageText),
		)

		return containerStyle.Render(lipgloss.JoinVertical(lipgloss.Center, content))
	}

	// Not fasting view
	return containerStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Center,
			statusStyle.Foreground(lipgloss.Color("9")).Render("Not fasting"), // Red for inactive
			"",
			lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Render("Start a new session with 'fasting-cli start <duration>'"),
		),
	)
}

