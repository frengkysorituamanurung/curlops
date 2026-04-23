package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")).
			Padding(1, 0)
	
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Padding(1, 0)
	
	inputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF75B5"))
	
	detailStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1, 2).
			Width(80)
)

func (m model) View() string {
	switch m.currentView {
	case listView:
		return m.viewList()
	case addView:
		return m.viewAdd()
	case detailView:
		return m.viewDetail()
	case editView:
		return m.viewEdit()
	case searchView:
		return m.viewSearch()
	}
	return ""
}

func (m model) viewList() string {
	help := helpStyle.Render("/: search • a: add • e: edit • c: copy • d: delete • enter: view • q: quit")
	status := ""
	if m.statusMsg != "" {
		status = "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Render(m.statusMsg)
	}
	return fmt.Sprintf("%s%s\n%s", m.list.View(), status, help)
}

func (m model) viewAdd() string {
	var b strings.Builder
	
	b.WriteString(titleStyle.Render("Add New Curl Command"))
	b.WriteString("\n\n")
	
	b.WriteString(inputStyle.Render("Name:"))
	b.WriteString("\n")
	b.WriteString(m.nameInput.View())
	b.WriteString("\n\n")
	
	b.WriteString(inputStyle.Render("Curl Command:"))
	b.WriteString("\n")
	b.WriteString(m.curlInput.View())
	b.WriteString("\n\n")
	
	b.WriteString(inputStyle.Render("Service:"))
	b.WriteString("\n")
	b.WriteString(m.serviceInput.View())
	b.WriteString("\n\n")
	
	b.WriteString(inputStyle.Render("Description (optional):"))
	b.WriteString("\n")
	b.WriteString(m.descriptionInput.View())
	b.WriteString("\n\n")
	
	b.WriteString(helpStyle.Render("tab: next field • enter: save • esc: cancel"))
	
	return b.String()
}

func (m model) viewDetail() string {
	if m.cursor >= len(m.commands) {
		return "No command selected"
	}
	
	cmd := m.commands[m.cursor]
	
	var b strings.Builder
	b.WriteString(titleStyle.Render(cmd.Name))
	b.WriteString("\n\n")
	
	if cmd.Service != "" {
		b.WriteString(fmt.Sprintf("Service: %s\n\n", cmd.Service))
	}
	
	if cmd.Desc != "" {
		b.WriteString(fmt.Sprintf("Description: %s\n\n", cmd.Desc))
	}
	
	b.WriteString(fmt.Sprintf("Curl Command:\n%s\n\n", cmd.Curl))
	
	if m.statusMsg != "" {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Render(m.statusMsg))
		b.WriteString("\n\n")
	}
	
	b.WriteString(helpStyle.Render("c: copy to clipboard • e: edit • esc: back"))
	
	return detailStyle.Render(b.String())
}

func (m model) viewEdit() string {
	var b strings.Builder
	
	b.WriteString(titleStyle.Render("Edit Curl Command"))
	b.WriteString("\n\n")
	
	b.WriteString(inputStyle.Render("Name:"))
	b.WriteString("\n")
	b.WriteString(m.nameInput.View())
	b.WriteString("\n\n")
	
	b.WriteString(inputStyle.Render("Curl Command:"))
	b.WriteString("\n")
	b.WriteString(m.curlInput.View())
	b.WriteString("\n\n")
	
	b.WriteString(inputStyle.Render("Service:"))
	b.WriteString("\n")
	b.WriteString(m.serviceInput.View())
	b.WriteString("\n\n")
	
	b.WriteString(inputStyle.Render("Description (optional):"))
	b.WriteString("\n")
	b.WriteString(m.descriptionInput.View())
	b.WriteString("\n\n")
	
	b.WriteString(helpStyle.Render("tab: next field • enter: save • esc: cancel"))
	
	return b.String()
}

func (m model) viewSearch() string {
	var b strings.Builder
	
	b.WriteString(titleStyle.Render("Search Curl Commands"))
	b.WriteString("\n\n")
	
	b.WriteString(inputStyle.Render("Search: "))
	b.WriteString(m.searchInput.View())
	b.WriteString("\n\n")
	
	// Show search results
	if m.searchInput.Value() == "" {
		b.WriteString(helpStyle.Render("Type to search by name, category, or command..."))
	} else if len(m.searchResults) == 0 {
		b.WriteString(helpStyle.Render("No results found"))
	} else {
		b.WriteString(lipgloss.NewStyle().Bold(true).Render(fmt.Sprintf("Found %d result(s):", len(m.searchResults))))
		b.WriteString("\n\n")
		
		// Display results
		for i, cmd := range m.searchResults {
			style := lipgloss.NewStyle()
			if i == m.searchCursor {
				style = style.Foreground(lipgloss.Color("#7D56F4")).Bold(true)
				b.WriteString(style.Render("▶ "))
			} else {
				b.WriteString("  ")
			}
			
			b.WriteString(style.Render(cmd.Name))
			if cmd.Service != "" {
				b.WriteString(style.Render(fmt.Sprintf(" [%s]", cmd.Service)))
			}
			b.WriteString("\n")
			
			// Show description if available
			if cmd.Desc != "" {
				descPreview := cmd.Desc
				if len(descPreview) > 60 {
					descPreview = descPreview[:60] + "..."
				}
				b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Render("  " + descPreview))
				b.WriteString("\n")
			}
			
			// Show truncated curl command
			curlPreview := cmd.Curl
			if len(curlPreview) > 60 {
				curlPreview = curlPreview[:60] + "..."
			}
			b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render("  " + curlPreview))
			b.WriteString("\n\n")
		}
	}
	
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("↑/↓: navigate • enter: view detail • ctrl+c: copy • esc: back"))
	
	return b.String()
}
