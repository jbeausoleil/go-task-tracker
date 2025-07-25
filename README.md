# Task Tracker (Bubble Tea)

A modern, interactive task tracker built with [Bubble Tea](https://github.com/charmbracelet/bubbletea).
This version supports:
- **Natural language due dates** (`tomorrow`, `next week`, `next Monday`, `3 days from now`).
- **Blue color scheme** for a clean modern look.
- **Persistent storage** in `tasks.json`.
- **Commands**: Add, List, Mark as Done, Delete, and filter by date.

## Build and Run
```bash
go mod tidy
go run main.go
```

## Cross Compile
```bash
GOOS=linux GOARCH=amd64 go build -o dist/task-cli main.go
```

## Key Bindings
- `↑/↓`: Navigate
- `Enter`: Select
- `ESC`: Go back
- `q`: Quit