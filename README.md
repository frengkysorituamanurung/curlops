# Curl Keeper

A TUI (Text User Interface) application for SRE to manage and organize curl commands by category.

## Features

- 📝 Save curl commands with names and categories
- 🗂️ Organize commands by category
- 🔍 Search and filter commands
- 📋 Copy commands to clipboard
- 💾 Persistent storage in `~/.curlops/`
- 🎨 Beautiful terminal interface

## Installation

### From Source

```bash
git clone <repository-url>
cd curlops
go build -o curlops
sudo mv curlops /usr/local/bin/
```

### Using Go Install

```bash
go install github.com/yourusername/curlops@latest
```

## Usage

Run the application:

```bash
curlops
```

### Keyboard Shortcuts

**List View:**
- `a` - Add new curl command
- `e` - Edit selected command
- `c` - Copy selected command to clipboard
- `d` - Delete selected command
- `Enter` - View command details
- `q` - Quit application

**Add/Edit View:**
- `Tab` - Move to next field
- `Shift+Tab` - Move to previous field
- `Enter` - Save command
- `Esc` - Cancel and back to list

**Detail View:**
- `c` - Copy command to clipboard
- `e` - Edit command
- `Esc` - Back to list

## Storage

Commands are stored in JSON format at:
- macOS/Linux: `~/.curlops/commands.json`

## Example Commands

```json
{
  "name": "Check API Health",
  "curl": "curl -X GET https://api.example.com/health",
  "service": "monitoring",
  "description": "Health check endpoint for production API"
}
```

## License

MIT License
