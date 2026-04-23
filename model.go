package main

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type view int

const (
	listView view = iota
	addView
	detailView
	editView
	searchView
)

type model struct {
	list       list.Model
	currentView view
	commands   []CurlCommand
	cursor     int
	editIndex  int
	
	// Add/Edit form inputs
	nameInput        textinput.Model
	curlInput        textinput.Model
	serviceInput     textinput.Model
	descriptionInput textinput.Model
	focusIndex       int
	
	// Search
	searchInput   textinput.Model
	searchResults []CurlCommand
	searchCursor  int
	
	// Status message
	statusMsg string
	
	width  int
	height int
	err    error
}

func initialModel() model {
	commands, err := loadCommands()
	if err != nil {
		commands = []CurlCommand{}
	}
	
	items := make([]list.Item, len(commands))
	for i, cmd := range commands {
		items[i] = cmd
	}
	
	// Create delegate with better styling for multi-line items
	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = true
	delegate.SetHeight(5) // Increased height for wrapped text
	delegate.SetSpacing(1) // Add spacing between items
	
	l := list.New(items, delegate, 0, 0)
	l.Title = "Curl Keeper - Your Curl Commands Manager"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false) // We have custom help
	
	// Use default rank function which uses fuzzy matching on FilterValue
	// The default should work, but let's make sure items are set correctly
	
	nameInput := textinput.New()
	nameInput.Placeholder = "Command name"
	nameInput.Focus()
	nameInput.Width = 60
	
	curlInput := textinput.New()
	curlInput.Placeholder = "curl command (e.g., curl -X GET https://api.example.com)"
	curlInput.Width = 60
	
	serviceInput := textinput.New()
	serviceInput.Placeholder = "Service name (e.g., api, monitoring, database)"
	serviceInput.Width = 60
	
	descriptionInput := textinput.New()
	descriptionInput.Placeholder = "Description (optional)"
	descriptionInput.Width = 60
	
	searchInput := textinput.New()
	searchInput.Placeholder = "Type to search..."
	searchInput.Width = 60
	
	return model{
		list:             l,
		currentView:      listView,
		commands:         commands,
		nameInput:        nameInput,
		curlInput:        curlInput,
		serviceInput:     serviceInput,
		descriptionInput: descriptionInput,
		searchInput:      searchInput,
		searchResults:    []CurlCommand{},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}
