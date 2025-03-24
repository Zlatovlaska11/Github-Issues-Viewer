package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

type CreateIssueFormPage struct {
	titleInput       textinput.Model
	descriptionInput textinput.Model
	focusIndex       int
	done             bool
	result           CreateIssueResult
}

type CreateIssueResult struct {
	Title       string
	Description string
	Submitted   bool
}

func NewCreateIssueFormPage() CreateIssueFormPage {
	title := textinput.New()
	title.Placeholder = "Issue title"
	title.Focus()
	title.CharLimit = 100
	title.Width = 30

	description := textinput.New()
	description.Placeholder = "Issue description"
	description.CharLimit = 200
	description.Width = 30

	return CreateIssueFormPage{
		titleInput:       title,
		descriptionInput: description,
		focusIndex:       0,
	}
}

func (p CreateIssueFormPage) Update(msg tea.Msg) (CreateIssueFormPage, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab", "enter", "up", "down":
			if msg.String() == "enter" && p.focusIndex == 1 {
				// Submit
				p.done = true
				p.result = CreateIssueResult{
					Title:       p.titleInput.Value(),
					Description: p.descriptionInput.Value(),
					Submitted:   true,
				}
				return p, nil
			}

			// Change focus
			if msg.String() == "up" || msg.String() == "shift+tab" {
				p.focusIndex--
			} else {
				p.focusIndex++
			}

			if p.focusIndex > 1 {
				p.focusIndex = 0
			} else if p.focusIndex < 0 {
				p.focusIndex = 1
			}

			if p.focusIndex == 0 {
				p.titleInput.Focus()
				p.descriptionInput.Blur()
			} else {
				p.titleInput.Blur()
				p.descriptionInput.Focus()
			}
		}
	}

	var cmd1, cmd2 tea.Cmd
	p.titleInput, cmd1 = p.titleInput.Update(msg)
	p.descriptionInput, cmd2 = p.descriptionInput.Update(msg)
	return p, tea.Batch(cmd1, cmd2)
}


func (p CreateIssueFormPage) View() string {
	var b strings.Builder
	fmt.Fprintln(&b, "📝 Create New Issue\n")
	fmt.Fprintln(&b, "Title:")
	fmt.Fprintln(&b, p.titleInput.View())
	fmt.Fprintln(&b, "\nDescription:")
	fmt.Fprintln(&b, p.descriptionInput.View())

	button := "[ Submit ]"
	if p.focusIndex == 1 {
		button = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render(button)
	}
	fmt.Fprintf(&b, "\n\n%s", button)
	fmt.Fprintln(&b, "\n\n[Tab to switch, Enter to submit, Backspace to cancel]")

	return b.String()
}

