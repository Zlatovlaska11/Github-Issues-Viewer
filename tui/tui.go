package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
	"zlatolas/projectManager/dataSchemes"
)

type model struct {
	page         int
	titleInput   textinput.Model
	descInput    textinput.Model
	labelInput   textinput.Model
	issueTitle   string
	issueDesc    string
	issueLabels  string
	currentInput int

	table table.Model
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

const (
	page1 = iota
	page2
	page3
)

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, tea.Quit
		case "up":
			if m.currentInput > 0 {
				m.currentInput--
				m.updateFocus() // Update focus when navigating
			}
		case "down":
			if m.currentInput < 2 { // 0-based index for 3 inputs
				m.currentInput++
				m.updateFocus() // Update focus when navigating
			}
		case "tab":
			// Switch between pages
			if m.page == page1 {
				m.page = page2
			} else {
				m.page = page1
			}
    case "q":
      if m.page == page3 {
        m.page = page1
      }
		case "enter":
			if m.page != 0 {
				m.page = page3
			} else {
				m.titleInput.Reset()
				m.descInput.Reset()
				m.labelInput.Reset()
			}
			// TODO add some sending logic
		}
	}

	// Update the input fields if on page 1
	var cmd tea.Cmd
	if m.page == page1 {
		switch m.currentInput {
		case 0:
			m.titleInput, cmd = m.titleInput.Update(msg)
		case 1:
			m.descInput, cmd = m.descInput.Update(msg)
		case 2:
			m.labelInput, cmd = m.labelInput.Update(msg)
		}
		return m, cmd
	}
	m.table, cmd = m.table.Update(msg)

	return m, nil
}

// Update the focus based on current input field
func (m *model) updateFocus() {
	switch m.currentInput {
	case 0:
		m.titleInput.Focus()
		m.descInput.Blur()
		m.labelInput.Blur()
	case 1:
		m.titleInput.Blur()
		m.descInput.Focus()
		m.labelInput.Blur()
	case 2:
		m.titleInput.Blur()
		m.descInput.Blur()
		m.labelInput.Focus()
	}
}

func (m model) View() string {
	var view string

	// Tab-like view at the top
	view += "== Tabs ==\n"
	if m.page == page1 {
		view += "-> Page 1\nPage 2\n"
	} else {
		view += "Page 1\n-> Page 2\n"
	}

	// Main view content based on current page
	switch m.page {
	case page1:
		// Page 1 with the GitHub Issue form
		view += "\nPage 1: Create a GitHub Issue\n\n"
		view += fmt.Sprintf("Title: %s\n", m.titleInput.View())
		view += fmt.Sprintf("Description: %s\n", m.descInput.View())
		view += fmt.Sprintf("Labels: %s\n", m.labelInput.View())
		view += fmt.Sprintf("[Submit]")
		view += "\nPress 'Tab' to go to the summary page\nPress 'esc' to quit."
	case page2:
		view += m.table.View() + "\nPress q to quit."
	case page3:
		view += fmt.Sprintf("issue number: %s\n", m.table.SelectedRow()[0])
		view += fmt.Sprintf("description: %s\n", m.table.SelectedRow()[1])
		view += fmt.Sprintf("status: %s\n", m.table.SelectedRow()[2])
		view += fmt.Sprintf("assignee: %s\n", m.table.SelectedRow()[3])
	}

	return view
}

func InitTui() {
	// Initialize form inputs for GitHub issue creation
	titleInput := textinput.New()
	titleInput.Placeholder = "Issue Title"
	titleInput.Focus()

	descInput := textinput.New()
	descInput.Placeholder = "Issue Description"

	labelInput := textinput.New()
	labelInput.Placeholder = "Issue Labels (comma separated)"

	columns := []table.Column{
		{Title: "Number", Width: 4},
		{Title: "Title", Width: 10},
		{Title: "State", Width: 10},
		{Title: "Asignee", Width: 10},
	}

	issues := dataschemes.GetIssues("Strnadi-Mobile-App")

	rows := dataschemes.ParseIssues(issues)

	rws := []table.Row{}
	for _, issue := range rows.Issues {
		row := table.Row{
			fmt.Sprintf("#%d", issue.Number),
			issue.Title,
			issue.State,
			issue.Asignee,
		}
		rws = append(rws, row)
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rws),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	// Initialize the program with a model containing the form inputs
	p := tea.NewProgram(model{
		page:         page1,
		table:        t,
		titleInput:   titleInput,
		descInput:    descInput,
		labelInput:   labelInput,
		currentInput: 0,
	})

	if err := p.Start(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
