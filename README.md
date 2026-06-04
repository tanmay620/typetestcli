# Terminal Typing Speed Test

A terminal-based typing speed test built with Go, Bubble Tea, and Lip Gloss. The application provides a simple environment for practicing typing, tracking accuracy, and measuring words per minute (WPM) directly from the command line.

## Overview

This project presents a short typing challenge where users type a predefined passage within a fixed time limit. As characters are entered, the application highlights correct and incorrect input in real time and calculates typing speed when the test ends.

## Features

* 30-second typing challenge
* Real-time character validation
* Countdown timer
* WPM calculation based on correctly typed characters
* Interactive terminal interface built with Bubble Tea
* Styled text rendering with Lip Gloss
* Restart functionality without restarting the application

## Preview

```text
┌───────────────────────────────────────────────┐
│                 Typing Test                   │
│               Measure Your WPM                │
└───────────────────────────────────────────────┘

Typing quickly is not just about moving your fingers...
```

Character states are displayed as follows:

* Correct characters appear normally
* Incorrect characters are highlighted
* Remaining characters are shown in a muted style
* A cursor indicates the current typing position

## Installation

### Prerequisites

* Go 1.22 or later

### Clone the Repository

```bash
git clone https://github.com/yourusername/typing-test.git
cd typing-test
```

### Install Dependencies

```bash
go mod tidy
```

### Run the Application

```bash
go run .
```

Or build a standalone executable:

```bash
go build -o typing-test
./typing-test
```

## Controls

| Key           | Action                    |
| ------------- | ------------------------- |
| Any character | Start or continue typing  |
| Backspace     | Remove the last character |
| Ctrl + R      | Restart the test          |
| Esc           | Exit the application      |
| Ctrl + C      | Exit the application      |

## Technologies Used

* Go
* Bubble Tea
* Lip Gloss

### Dependencies

```go
github.com/charmbracelet/bubbletea
github.com/charmbracelet/lipgloss
```

## Possible Improvements

Some ideas for future development include:

* Color Scheme
* Randomized typing passages
* Custom test durations

## License

This project is licensed under the MIT License.

## Acknowledgements

This project uses the excellent Bubble Tea and Lip Gloss libraries from the Charm ecosystem to create an interactive terminal experience.

