package main

import (
	"fmt"
	"os"

	bubbles_list "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type status int

const (
	todo status = iota
	inProgress
	done
)

/* CUSTOM ITEM */
type Task struct {
	status      status
	title       string
	description string
}

// Implement the list.item interface
func (t Task) FilterValue() string {
	return t.title
}

func (t Task) Title() string {
	return t.title
}

func (t Task) Description() string {
	return t.description
}

// main model
type Model struct {
	list bubbles_list.Model
	err  error
}

func NewModel() *Model {
	return &Model{}
}

func (m *Model) initList(width int, height int) {
	m.list = bubbles_list.New([]bubbles_list.Item{}, bubbles_list.NewDefaultDelegate(), width, height)
	m.list.Title = "To Do"
	m.list.SetItems([]bubbles_list.Item{
		Task{status: todo, title: "Buy Milk from the store!", description: "Tesco is running a sale..."},
		Task{status: todo, title: "Look at clouds", description: "They are so fluffy!"},
		Task{status: todo, title: "Call Mom", description: "It's her birthday!"},
		Task{status: todo, title: "Go for a run", description: "It's been a while..."},
		Task{status: todo, title: "Read a book", description: "I have a few in mind..."},
	})
}

func (m Model) Init() tea.Cmd {
	m.initList(40, 10)
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.initList(msg.Width, msg.Height)
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.list.View()
}

func main() {
	m := NewModel()
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
