package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type CurlCommand struct {
	Name string `json:"name"`
	Curl string `json:"curl"`
	Service string `json:"service"`
	Desc string `json:"description"`
}

func (c CurlCommand) Title() string { 
	return c.Name 
}

func (c CurlCommand) Description() string { 
	var result strings.Builder
	
	// First line: service tag
	if c.Service != "" {
		result.WriteString("[" + c.Service + "]")
		result.WriteString("\n")
	}
	
	// Curl command with word wrapping
	result.WriteString("  ")
	result.WriteString(wrapText(c.Curl, 100, "  "))
	
	// Description on separate line if available
	if c.Desc != "" {
		result.WriteString("\n  ")
		result.WriteString(c.Desc)
	}
	
	return result.String()
}

// wrapText wraps text to specified width with indent for continuation lines
func wrapText(text string, width int, indent string) string {
	if len(text) <= width {
		return text
	}
	
	var result strings.Builder
	words := strings.Fields(text)
	lineLen := 0
	firstLine := true
	
	for i, word := range words {
		wordLen := len(word)
		
		if lineLen+wordLen+1 > width && lineLen > 0 {
			result.WriteString("\n")
			if !firstLine {
				result.WriteString(indent)
			}
			firstLine = false
			result.WriteString(indent)
			result.WriteString(word)
			lineLen = len(indent) + wordLen
		} else {
			if i > 0 {
				result.WriteString(" ")
				lineLen++
			}
			result.WriteString(word)
			lineLen += wordLen
		}
	}
	
	return result.String()
}

func (c CurlCommand) FilterValue() string { 
	// Include name, service, description, and curl for comprehensive search
	// Convert to lowercase for case-insensitive matching
	value := c.Name + " " + c.Service + " " + c.Desc + " " + c.Curl
	return strings.ToLower(value)
}

func getStoragePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".curlkeeper", "commands.json"), nil
}

func initStorage() error {
	path, err := getStoragePath()
	if err != nil {
		return err
	}
	
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return saveCommands([]CurlCommand{})
	}
	
	return nil
}

func loadCommands() ([]CurlCommand, error) {
	path, err := getStoragePath()
	if err != nil {
		return nil, err
	}
	
	data, err := os.ReadFile(path)
	if err != nil {
		return []CurlCommand{}, nil
	}
	
	var commands []CurlCommand
	if err := json.Unmarshal(data, &commands); err != nil {
		return nil, err
	}
	
	return commands, nil
}

func saveCommands(commands []CurlCommand) error {
	path, err := getStoragePath()
	if err != nil {
		return err
	}
	
	data, err := json.MarshalIndent(commands, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(path, data, 0644)
}
