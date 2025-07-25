package main

import (
    "fmt"
    "os"
    "sort"
    "strings"
    "time"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/bubbles/textinput"
    "github.com/charmbracelet/lipgloss"
    "github.com/jbeausoleil/task-tracker-bubbletea-blue-natural/internal/task"
    "github.com/jinzhu/now"
)

type state int
type addStep int

const (
    menu state = iota
    list
    add
    delete
    mark
    completedToday
    dueThisWeek
)

const (
    addDesc addStep = iota
    addDue
)

// Blue theme styles
var (
    titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00BFFF"))
    selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#1E90FF"))
    normalStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#87CEFA"))
    doneStyle     = lipgloss.NewStyle().Strikethrough(true).Foreground(lipgloss.Color("#5F9EA0"))
    helpStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#5F9EA0"))
)

type model struct {
    cursor     int
    choices    []string
    state      state
    tasks      []task.Task
    textInput  textinput.Model
    dueInput   textinput.Model
    addStep    addStep
    store      *task.Store
    tempDesc   string
}

func initialModel() model {
    descInput := textinput.New()
    descInput.Placeholder = "Type task description"
    descInput.Focus()
    descInput.CharLimit = 128
    descInput.Width = 40

    dueInput := textinput.New()
    dueInput.Placeholder = "YYYY-MM-DD or 'tomorrow'"
    dueInput.CharLimit = 30
    dueInput.Width = 35

    store := task.NewStore("tasks.json")
    tasks, _ := store.Load()

    return model{
        choices: []string{
            "Add Task",
            "List All Tasks",
            "List Completed Tasks Today",
            "List Tasks Due This Week",
            "Mark Task as Done",
            "Delete Task",
            "Quit",
        },
        state:     menu,
        tasks:     tasks,
        textInput: descInput,
        dueInput:  dueInput,
        addStep:   addDesc,
        store:     store,
    }
}

func (m model) Init() tea.Cmd {
    return nil
}

// parseDueDate supports natural language using jinzhu/now
func parseDueDate(input string) time.Time {
    nowRef := time.Now()
    input = strings.ToLower(strings.TrimSpace(input))

    switch input {
    case "", "default":
        return nowRef.Add(48 * time.Hour)
    case "today":
        return nowRef
    case "tomorrow":
        return nowRef.Add(24 * time.Hour)
    case "next week":
        return nowRef.AddDate(0, 0, 7)
    }

    // Fuzzy parsing with jinzhu/now
    if t, err := now.Parse(input); err == nil {
        return t
    }

    // Fallback to YYYY-MM-DD
    if t, err := time.Parse("2006-01-02", input); err == nil {
        return t
    }

    // Default fallback
    return nowRef.Add(48 * time.Hour)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd
    switch m.state {
    case menu:
        switch msg := msg.(type) {
        case tea.KeyMsg:
            switch msg.String() {
            case "ctrl+c", "q":
                return m, tea.Quit
            case "up":
                if m.cursor > 0 {
                    m.cursor--
                }
            case "down":
                if m.cursor < len(m.choices)-1 {
                    m.cursor++
                }
            case "enter":
                switch m.cursor {
                case 0:
                    m.state = add
                    m.addStep = addDesc
                case 1:
                    m.state = list
                case 2:
                    m.state = completedToday
                case 3:
                    m.state = dueThisWeek
                case 4:
                    m.state = mark
                    m.cursor = 0
                case 5:
                    m.state = delete
                    m.cursor = 0
                case 6:
                    return m, tea.Quit
                }
            }
        }

    case add:
        if m.addStep == addDesc {
            m.textInput, cmd = m.textInput.Update(msg)
            switch msg := msg.(type) {
            case tea.KeyMsg:
                switch msg.String() {
                case "esc":
                    m.state = menu
                case "enter":
                    if m.textInput.Value() != "" {
                        m.tempDesc = m.textInput.Value()
                        m.textInput.SetValue("")
                        m.addStep = addDue
                        m.dueInput.Focus()
                    } else {
                        m.state = menu
                    }
                }
            }
        } else if m.addStep == addDue {
            m.dueInput, cmd = m.dueInput.Update(msg)
            switch msg := msg.(type) {
            case tea.KeyMsg:
                switch msg.String() {
                case "esc":
                    m.state = menu
                case "enter":
                    dueDate := parseDueDate(m.dueInput.Value())
                    newTask := task.Task{
                        Description: m.tempDesc,
                        Status:      "todo",
                        CreatedAt:   time.Now(),
                        DueDate:     dueDate,
                    }
                    m.tasks = append(m.tasks, newTask)
                    m.store.Save(m.tasks)
                    m.tempDesc = ""
                    m.dueInput.SetValue("")
                    m.state = menu
                }
            }
        }

    case list, completedToday, dueThisWeek:
        switch msg := msg.(type) {
        case tea.KeyMsg:
            if msg.String() == "esc" {
                m.state = menu
            }
        }

    case mark:
        switch msg := msg.(type) {
        case tea.KeyMsg:
            switch msg.String() {
            case "up":
                if m.cursor > 0 {
                    m.cursor--
                }
            case "down":
                if m.cursor < len(m.tasks)-1 {
                    m.cursor++
                }
            case "enter":
                if len(m.tasks) > 0 {
                    m.tasks[m.cursor].Status = "done"
                    m.tasks[m.cursor].CompletedAt = time.Now()
                    m.store.Save(m.tasks)
                }
            case "esc":
                m.state = menu
            }
        }

    case delete:
        switch msg := msg.(type) {
        case tea.KeyMsg:
            switch msg.String() {
            case "up":
                if m.cursor > 0 {
                    m.cursor--
                }
            case "down":
                if m.cursor < len(m.tasks)-1 {
                    m.cursor++
                }
            case "enter":
                if len(m.tasks) > 0 {
                    m.tasks = append(m.tasks[:m.cursor], m.tasks[m.cursor+1:]...)
                    if m.cursor >= len(m.tasks) && m.cursor > 0 {
                        m.cursor--
                    }
                    m.store.Save(m.tasks)
                }
            case "esc":
                m.state = menu
            }
        }
    }
    return m, cmd
}

func (m model) View() string {
    switch m.state {
    case list:
        return m.renderTasks(m.tasks, "All Tasks") + helpStyle.Render("\n[↑/↓]: Scroll  [ESC]: Back  [q]: Quit")

    case completedToday:
        todayTasks := task.CompletedToday(m.tasks)
        return m.renderTasks(todayTasks, "Completed Tasks Today") + helpStyle.Render("\n[↑/↓]: Scroll  [ESC]: Back  [q]: Quit")

    case dueThisWeek:
        weekTasks := task.DueThisWeek(m.tasks)
        return m.renderTasks(weekTasks, "Tasks Due This Week") + helpStyle.Render("\n[↑/↓]: Scroll  [ESC]: Back  [q]: Quit")

    case add:
        if m.addStep == addDesc {
            return fmt.Sprintf(
                "%s\n\n%s\n\n(Press Enter to continue or ESC to cancel)\n",
                titleStyle.Render("Add a new task (Description):"),
                m.textInput.View(),
            )
        } else {
            return fmt.Sprintf(
                "%s\n\n%s\n\n(Enter YYYY-MM-DD or natural text: 'tomorrow', 'next week', '3 days from now')\n(Press Enter to save or ESC to cancel)\n",
                titleStyle.Render("Add a new task (Due Date):"),
                m.dueInput.View(),
            )
        }

    case mark:
        return m.renderSelection("Mark Task as Done") + helpStyle.Render("\n[↑/↓]: Navigate  [Enter]: Mark Done  [ESC]: Back")

    case delete:
        return m.renderSelection("Delete Task") + helpStyle.Render("\n[↑/↓]: Navigate  [Enter]: Delete  [ESC]: Back")

    default:
        s := titleStyle.Render("Task Tracker (Bubble Tea)") + "\n\n"
        for i, choice := range m.choices {
            style := normalStyle
            if m.cursor == i {
                style = selectedStyle
            }
            s += style.Render(fmt.Sprintf("%s", choice)) + "\n"
        }
        s += "\n" + helpStyle.Render("[↑/↓]: Navigate  [Enter]: Select  [q]: Quit") + "\n"
        return s
    }
}

func (m model) renderTasks(tasks []task.Task, title string) string {
    if len(tasks) == 0 {
        return fmt.Sprintf("No tasks found. Press ESC to return.")
    }
    sort.Slice(tasks, func(i, j int) bool {
        return tasks[i].CreatedAt.Before(tasks[j].CreatedAt)
    })
    header := fmt.Sprintf("%-3s %-8s %-30s %-12s %-16s", "#", "STATUS", "DESCRIPTION", "DUE DATE", "CREATED")
    s := titleStyle.Render(title) + "\n\n" + lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00BFFF")).Render(header) + "\n"
    s += strings.Repeat("-", len(header)) + "\n"
    for i, t := range tasks {
        due := t.DueDate.Format("2006-01-02")
        created := t.CreatedAt.Format("2006-01-02 15:04")
        line := fmt.Sprintf("%-3d %-8s %-30s %-12s %-16s", i+1, t.Status, t.Description, due, created)
        if t.Status == "done" {
            line = doneStyle.Render(line)
        }
        s += line + "\n"
    }
    return s
}

func (m model) renderSelection(title string) string {
    if len(m.tasks) == 0 {
        return fmt.Sprintf("No tasks available. Press ESC to return.")
    }
    s := titleStyle.Render(title) + "\n\n"
    for i, t := range m.tasks {
        cursor := " "
        if m.cursor == i {
            cursor = ">"
        }
        s += fmt.Sprintf("%s [%s] %s\n", cursor, t.Status, t.Description)
    }
    return s
}

func main() {
    p := tea.NewProgram(initialModel())
    if err := p.Start(); err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }
}