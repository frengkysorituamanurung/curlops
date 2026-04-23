package main

import (
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width, msg.Height-2)
		return m, nil
		
	case tea.KeyMsg:
		switch m.currentView {
		case listView:
			return m.updateListView(msg)
		case addView:
			return m.updateAddView(msg)
		case detailView:
			return m.updateDetailView(msg)
		case editView:
			return m.updateEditView(msg)
		case searchView:
			return m.updateSearchView(msg)
		}
	}
	
	return m, nil
}

func (m model) updateListView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Check if list is in filtering mode
	isFiltering := m.list.FilterState() == list.Filtering
	
	// If filtering, let the list handle ALL keys except quit
	if isFiltering {
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		default:
			// Pass everything to list when filtering
			var cmd tea.Cmd
			m.list, cmd = m.list.Update(msg)
			return m, cmd
		}
	}
	
	// Clear status message on any key press except 'c'
	if msg.String() != "c" {
		m.statusMsg = ""
	}
	
	// Normal mode (not filtering) - handle our custom keys
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "/":
		// Enter custom search mode
		m.currentView = searchView
		m.searchInput.SetValue("")
		m.searchInput.Focus()
		m.searchResults = []CurlCommand{}
		m.searchCursor = 0
		return m, nil
	case "a":
		m.currentView = addView
		m.focusIndex = 0
		m.nameInput.Focus()
		m.nameInput.SetValue("")
		m.curlInput.SetValue("")
		m.serviceInput.SetValue("")
		m.descriptionInput.SetValue("")
		return m, nil
	case "enter":
		if len(m.commands) > 0 {
			m.currentView = detailView
			m.cursor = m.list.Index()
		}
		return m, nil
	case "e":
		// Edit directly from list view
		if len(m.commands) > 0 {
			idx := m.list.Index()
			if idx >= 0 && idx < len(m.commands) {
				m.currentView = editView
				m.editIndex = idx
				m.cursor = idx
				m.focusIndex = 0
				
				// Load current values into inputs
				cmd := m.commands[idx]
				m.nameInput.SetValue(cmd.Name)
				m.curlInput.SetValue(cmd.Curl)
				m.serviceInput.SetValue(cmd.Service)
				m.descriptionInput.SetValue(cmd.Desc)
				m.nameInput.Focus()
			}
		}
		return m, nil
	case "c":
		// Copy directly from list view
		if len(m.commands) > 0 {
			idx := m.list.Index()
			if idx >= 0 && idx < len(m.commands) {
				err := clipboard.WriteAll(m.commands[idx].Curl)
				if err != nil {
					m.statusMsg = "❌ Failed to copy to clipboard"
				} else {
					m.statusMsg = "✓ Copied to clipboard!"
				}
			}
		}
		return m, nil
	case "d":
		if len(m.commands) > 0 {
			idx := m.list.Index()
			// Safe delete - handle edge cases
			if idx >= 0 && idx < len(m.commands) {
				// Remove the item at index
				newCommands := make([]CurlCommand, 0, len(m.commands)-1)
				newCommands = append(newCommands, m.commands[:idx]...)
				newCommands = append(newCommands, m.commands[idx+1:]...)
				m.commands = newCommands
				
				saveCommands(m.commands)
				
				// Update list items
				items := make([]list.Item, len(m.commands))
				for i, cmd := range m.commands {
					items[i] = cmd
				}
				m.list.SetItems(items)
			}
		}
		return m, nil
	}
	
	// Pass all other keys to list
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) updateAddView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.currentView = listView
		m.nameInput.SetValue("")
		m.curlInput.SetValue("")
		m.serviceInput.SetValue("")
		m.descriptionInput.SetValue("")
		return m, nil
	case "tab", "shift+tab":
		if msg.String() == "tab" {
			m.focusIndex = (m.focusIndex + 1) % 4
		} else {
			m.focusIndex = (m.focusIndex - 1 + 4) % 4
		}
		m.nameInput.Blur()
		m.curlInput.Blur()
		m.serviceInput.Blur()
		m.descriptionInput.Blur()
		
		switch m.focusIndex {
		case 0:
			m.nameInput.Focus()
		case 1:
			m.curlInput.Focus()
		case 2:
			m.serviceInput.Focus()
		case 3:
			m.descriptionInput.Focus()
		}
		return m, nil
	case "enter":
		if m.nameInput.Value() != "" && m.curlInput.Value() != "" {
			cmd := CurlCommand{
				Name:        m.nameInput.Value(),
				Curl:        m.curlInput.Value(),
				Service:     m.serviceInput.Value(),
				Desc: m.descriptionInput.Value(),
			}
			m.commands = append(m.commands, cmd)
			
			// Save to disk immediately
			if err := saveCommands(m.commands); err != nil {
				// Handle error but continue
				m.err = err
			}
			
			// Update list items
			items := make([]list.Item, len(m.commands))
			for i, c := range m.commands {
				items[i] = c
			}
			m.list.SetItems(items)
			
			m.nameInput.SetValue("")
			m.curlInput.SetValue("")
			m.serviceInput.SetValue("")
			m.descriptionInput.SetValue("")
			m.currentView = listView
		}
		return m, nil
	}
	
	var cmd tea.Cmd
	switch m.focusIndex {
	case 0:
		m.nameInput, cmd = m.nameInput.Update(msg)
	case 1:
		m.curlInput, cmd = m.curlInput.Update(msg)
	case 2:
		m.serviceInput, cmd = m.serviceInput.Update(msg)
	case 3:
		m.descriptionInput, cmd = m.descriptionInput.Update(msg)
	}
	return m, cmd
}

func (m model) updateDetailView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.currentView = listView
		m.statusMsg = ""
		return m, nil
	case "c":
		if m.cursor < len(m.commands) {
			err := clipboard.WriteAll(m.commands[m.cursor].Curl)
			if err != nil {
				m.statusMsg = "❌ Failed to copy to clipboard"
			} else {
				m.statusMsg = "✓ Copied to clipboard!"
			}
		}
		return m, nil
	case "e":
		if m.cursor < len(m.commands) {
			m.currentView = editView
			m.editIndex = m.cursor
			m.focusIndex = 0
			
			// Load current values into inputs
			cmd := m.commands[m.cursor]
			m.nameInput.SetValue(cmd.Name)
			m.curlInput.SetValue(cmd.Curl)
			m.serviceInput.SetValue(cmd.Service)
			m.descriptionInput.SetValue(cmd.Desc)
			m.nameInput.Focus()
		}
		return m, nil
	}
	return m, nil
}

func (m model) updateEditView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.currentView = listView
		m.nameInput.SetValue("")
		m.curlInput.SetValue("")
		m.serviceInput.SetValue("")
		m.descriptionInput.SetValue("")
		return m, nil
	case "tab", "shift+tab":
		if msg.String() == "tab" {
			m.focusIndex = (m.focusIndex + 1) % 4
		} else {
			m.focusIndex = (m.focusIndex - 1 + 4) % 4
		}
		m.nameInput.Blur()
		m.curlInput.Blur()
		m.serviceInput.Blur()
		m.descriptionInput.Blur()
		
		switch m.focusIndex {
		case 0:
			m.nameInput.Focus()
		case 1:
			m.curlInput.Focus()
		case 2:
			m.serviceInput.Focus()
		case 3:
			m.descriptionInput.Focus()
		}
		return m, nil
	case "enter":
		if m.nameInput.Value() != "" && m.curlInput.Value() != "" {
			// Update the command
			m.commands[m.editIndex] = CurlCommand{
				Name:        m.nameInput.Value(),
				Curl:        m.curlInput.Value(),
				Service:     m.serviceInput.Value(),
				Desc: m.descriptionInput.Value(),
			}
			
			// Save to disk immediately
			if err := saveCommands(m.commands); err != nil {
				m.err = err
			}
			
			// Update list items
			items := make([]list.Item, len(m.commands))
			for i, c := range m.commands {
				items[i] = c
			}
			m.list.SetItems(items)
			
			m.nameInput.SetValue("")
			m.curlInput.SetValue("")
			m.serviceInput.SetValue("")
			m.descriptionInput.SetValue("")
			m.currentView = listView
		}
		return m, nil
	}
	
	var cmd tea.Cmd
	switch m.focusIndex {
	case 0:
		m.nameInput, cmd = m.nameInput.Update(msg)
	case 1:
		m.curlInput, cmd = m.curlInput.Update(msg)
	case 2:
		m.serviceInput, cmd = m.serviceInput.Update(msg)
	case 3:
		m.descriptionInput, cmd = m.descriptionInput.Update(msg)
	}
	return m, cmd
}

func (m model) updateSearchView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.currentView = listView
		m.searchInput.SetValue("")
		m.searchResults = []CurlCommand{}
		return m, nil
	case "up":
		if m.searchCursor > 0 {
			m.searchCursor--
		}
		return m, nil
	case "down":
		if m.searchCursor < len(m.searchResults)-1 {
			m.searchCursor++
		}
		return m, nil
	case "enter":
		if len(m.searchResults) > 0 && m.searchCursor < len(m.searchResults) {
			// Find the index in original commands
			selectedCmd := m.searchResults[m.searchCursor]
			for i, cmd := range m.commands {
				if cmd.Name == selectedCmd.Name && cmd.Curl == selectedCmd.Curl {
					m.cursor = i
					m.currentView = detailView
					break
				}
			}
		}
		return m, nil
	case "ctrl+c":
		// Only Ctrl+C for copy (works on both Linux and macOS terminal)
		if len(m.searchResults) > 0 && m.searchCursor < len(m.searchResults) {
			err := clipboard.WriteAll(m.searchResults[m.searchCursor].Curl)
			if err != nil {
				m.statusMsg = "❌ Failed to copy to clipboard"
			} else {
				m.statusMsg = "✓ Copied to clipboard!"
			}
		}
		return m, nil
	}
	
	// Update search input - all other keys including 'c' go to input
	var cmd tea.Cmd
	m.searchInput, cmd = m.searchInput.Update(msg)
	
	// Perform search
	query := strings.ToLower(strings.TrimSpace(m.searchInput.Value()))
	if query == "" {
		m.searchResults = []CurlCommand{}
	} else {
		m.searchResults = []CurlCommand{}
		for _, command := range m.commands {
			searchText := strings.ToLower(command.Name + " " + command.Service + " " + command.Desc + " " + command.Curl)
			if strings.Contains(searchText, query) {
				m.searchResults = append(m.searchResults, command)
			}
		}
	}
	
	// Reset cursor if results changed
	if m.searchCursor >= len(m.searchResults) {
		m.searchCursor = 0
	}
	
	return m, cmd
}
