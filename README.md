# Fasting CLI

A simple CLI application to help you track your intermittent fasting schedule.

## Features

- Start and stop fasting timers with custom durations.
- Persist fasting sessions across application restarts.
- Display real-time fasting progress with a visual bar.

## Usage

### Build and Run

To build the application:

```shell
go build -o fasting-cli .
```

To run the application:

```shell
./fasting-cli
```

### Commands

- **Start a fasting session:**

  ```shell
  ./fasting-cli start <total_fast_duration> [--ago <time_elapsed_since_start>] [--protocol <name>]
  ```
  `<total_fast_duration>`: The total duration of the fasting period (e.g., `16h`, `18h30m`, `20h`). This argument is always required.

  **Options:**
  - `--ago <time_elapsed_since_start>`: Optional. Specify how long ago the fast *actually* started (e.g., `2h`, `30m`, `1h30m`). If not provided, the current time will be used as the start time.
  - `--protocol <name>`: Optional. Specify a fasting protocol name (e.g., `16-8`, `18-6`, `20-4`). If not provided, "Custom" will be used.

  **Example:** To start a 16-hour fast that began 5 hours and 43 minutes ago:
  ```bash
  ./fasting-cli start 16h --ago 5h43m
  ```

  **Safeguards:**
  - The fasting duration cannot exceed 48 hours.
  - The calculated end time (start time + total_fast_duration) cannot be in the past.

- **Stop the current fasting session:**

  ```shell
  ./fasting-cli stop
  ```

- **View current fasting progress:**

  Simply run the application without any arguments:
  ```shell
  ./fasting-cli
  ```
  If a session is active, it will display the remaining time and a progress bar.

### Keyboard Shortcuts (when viewing progress)

- Press `q` or `Ctrl+c` to quit the application.

