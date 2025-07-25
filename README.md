# Task Tracker (Bubble Tea Extended)

This version includes:
- Add tasks with description and due date.
- Prompt to add a note immediately after adding a task (optional).
- View, add, and delete notes for tasks.
- New "View Notes" menu option:
  - Select task by arrow keys or by typing its number (e.g., `3`).
- Mark tasks as done and list them in aligned columns.
- Natural language due dates (e.g., `tomorrow`, `next week`).

## Controls
- **Main Menu**: [↑/↓] Navigate • [Enter] Select • [q] Quit
- **Task List**: [↑/↓] Scroll • [Enter] View Notes • [d] Mark Done • [ESC] Back
- **Notes View**: [a] Add Note • [x] Delete Last Note • [ESC] Back
- **Number Keys**: Directly jump to a task in "View Notes".

## Build
```bash
go mod tidy
go build -o task-cli main.go
```