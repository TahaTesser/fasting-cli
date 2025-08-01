# Gemini Code Assistant Project: Fasting CLI

This document provides context for the Gemini Code Assistant to effectively assist with development of this project.

## 1. Project Overview

This project is a command-line interface (CLI) application for tracking intermittent fasting schedules. It is built with Go and utilizes the Bubble Tea and Lip Gloss libraries to create a rich, interactive terminal user interface.

The core features of the application include:
- Starting and stopping fasting timers.
- Viewing fasting history and statistics.
- Configuring different fasting protocols (e.g., 16:8, 18:6, 20:4, 23:1, 36:6).
- Displaying a real-time countdown timer and progress visualization.

## 2. Technologies and Key Libraries

- **Language:** Go
- **Core UI Framework:** [Bubble Tea](https://github.com/charmbracelet/bubbletea) - A powerful TUI (Terminal User Interface) framework based on The Elm Architecture.
- **Styling:** [Lip Gloss](https://github.com/charmbracelet/lipgloss) - A library for style-based rendering in the terminal, used for colors, layouts, and other visual flair.
- **Persistence:** (To be decided) - Data will likely be stored in a local file (e.g., JSON, SQLite).

## 3. Project Structure

The project follows a standard Go application layout:

```
.
├── main.go           # Application entry point
├── model/            # Data structures (e.g., FastingSession, UserProfile)
├── view/             # UI components and rendering logic (Bubble Tea views)
├── state/            # Application state management
├── cmd/              # Command-line commands and flags
└── go.mod            # Go module definition
```

## 4. Development Workflow

### 4.1. Building the Application

To build the application, use the standard `go build` command:

```shell
go build -o fasting-cli .
```

### 4.2. Running the Application

Execute the compiled binary to run the application:

```shell
./fasting-cli
```

Alternatively, run directly using `go run`:

```shell
go run main.go
```

### 4.3. Running Tests

Tests will be located in `_test.go` files alongside the code they are testing. Run all tests with:

```shell
go test ./...
```

## 5. Coding Style and Conventions

- Follow standard Go formatting (`gofmt`).
- Adhere to the principles of The Elm Architecture as implemented by Bubble Tea: Model, View, Update.
- Keep UI (View) logic separate from application state (Model) and business logic.
- Use Lip Gloss for all terminal styling to ensure a consistent look and feel.
