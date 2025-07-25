# Task Tracker (Bubble Tea Edition)

An interactive, terminal-based task tracker built with [Charmbracelet Bubble Tea](https://github.com/charmbracelet/bubbletea).

---

## Features

### Core Features
- **Add tasks** with a description and due date.
- **Natural language due date parsing**:
  - Examples: `tomorrow`, `next week`, `3 days from now`, `next monday`.
  - Also supports `YYYY-MM-DD` format.
- **Color-coded statuses**:
  - `todo` → Green.
  - `in-progress` → Blue.
  - `done` → Grey with strikethrough.
- **Mark tasks as done** and record the completion date.
- **Delete tasks** interactively.
- **Persistent storage** using `tasks.json`.

### Listing
- **List all tasks**, sorted by creation time.
- **Filter tasks**:
  - **Completed tasks today**.
  - **Tasks due this week**.

---

## Installation

### Prerequisites
- Go 1.21 or higher
- Git

### Clone and Build
```bash
git clone https://github.com/yourusername/go-task-tracker.git
cd go-task-tracker
go mod tidy
make build
```

---

## Usage

### Start the Program
```bash
./task-cli
```

### Key Bindings
- `↑/↓`: Navigate the menu
- `Enter`: Select option
- `ESC`: Back or cancel current input
- `q`: Quit the program

---

## Adding a Task
1. Select **Add Task** from the menu.
2. Enter a description.
3. Enter a due date (supports natural language like `tomorrow`, `next week`, or `YYYY-MM-DD`).
4. The task will be added with a `todo` status.

---

## Listing Tasks
- **List All Tasks** shows every task with STATUS, DESCRIPTION, DUE date, and CREATED date.
- **List Completed Tasks Today** shows tasks that were marked done today.
- **List Tasks Due This Week** shows tasks with due dates within the current week.

---

## Marking and Deleting Tasks
- **Mark Task as Done**: Choose a task from the list and mark it as completed.
- **Delete Task**: Choose a task to remove it permanently.

---

## Colors and Styles
- `STATUS` is color-coded:
  - **Green** for `todo`.
  - **Blue** for `in-progress`.
  - **Gray with strikethrough** for `done`.
- Dates are styled in **yellow and cyan** for due and created dates.

---

## Makefile
The included `Makefile` supports:
```bash
make build          # Build for current OS
make build-mac      # Build for macOS
make build-linux    # Build for Linux
make cross-compile  # Build for both macOS and Linux
make clean          # Clean binaries
```

---

## Roadmap
- Add tagging for tasks.
- Add priority levels.
- Export tasks to CSV.
- Integration with calendar APIs.

---

## License
This project is licensed under the MIT License.