# GoEdit v2.0

**A Modern Terminal Text Editor with AI Integration**

Copyright © Prof. Dr. Michael Stal, 2025  
All Rights Reserved

---

## 🚀 Overview

GoEdit is a powerful, lightweight terminal-based text editor written in Go, featuring advanced AI assistance through Ollama integration. It combines the simplicity of classic terminal editors with modern features like multi-file tabs, system clipboard integration, and real-time AI code generation.

### ✨ Key Features

- 🎯 **Multi-File Tabs** - Edit multiple files simultaneously with easy tab switching
- 🤖 **AI Assistant** - Integrated Ollama support with streaming responses
- 📋 **System Clipboard** - Full integration with macOS, Linux, and Windows clipboards
- ↩️ **Undo/Redo** - 50 levels of undo/redo support
- 🔍 **Search** - Fast text search across your documents
- ⚡ **Lightweight** - Minimal dependencies, fast startup
- 🎨 **Clean UI** - Intuitive tab bar and status indicators
- 💾 **Safe Saves** - Atomic file writes with temporary file protection

---

## 📦 Installation

### Prerequisites

- **Go 1.18+** - [Download Go](https://golang.org/dl/)
- **Ollama** (optional, for AI features) - [Download Ollama](https://ollama.ai/)

### Build from Source

```bash
# Clone the repository
git clone https://github.com/yourusername/goedit.git
cd goedit

# Initialize Go module
go mod init goedit

# Get dependencies
go get github.com/gdamore/tcell/v2

# Build
go build -o goedit

# Install (optional)
sudo mv goedit /usr/local/bin/

Quick Install (One-liner)
git clone https://github.com/yourusername/goedit.git &amp;&amp; cd goedit &amp;&amp; go mod init goedit &amp;&amp; go get github.com/gdamore/tcell/v2 &amp;&amp; go build -o goedit


🎮 Usage
Basic Usage
# Open a new file
goedit

# Edit existing file
goedit myfile.txt

# Edit multiple files in tabs
goedit file1.go file2.go file3.go

# Use specific AI model
goedit -model codellama main.go

# Enable streaming AI responses
goedit -stream -model llama2 document.md

# Custom Ollama server
goedit -ollama http://192.168.1.100:11434 file.txt

Command Line Options



Option
Default
Description



-ollama
http://localhost:11434
Ollama API URL


-model
llama2
LLM model to use


-stream
false
Enable streaming AI responses


-version
-
Show version information


-help
-
Display help message



⌨️ Keyboard Shortcuts
File Operations



Shortcut
Action



Ctrl+S
Save current file


Ctrl+Q
Quit editor (with unsaved changes warning)


Ctrl+T
Create new tab


Ctrl+W
Close current tab


Tab
Switch to next tab


Shift+Tab
Switch to previous tab


Editing



Shortcut
Action



Ctrl+A
Copy all text to system clipboard


Ctrl+C
Copy current line to system clipboard


Ctrl+X
Cut current line to system clipboard


Ctrl+V
Paste from system clipboard


Ctrl+Z
Undo (50 levels)


Ctrl+Y
Redo


Tab (in text)
Insert 4 spaces


Backspace
Delete character before cursor


Delete
Delete character at cursor


Navigation



Shortcut
Action



Ctrl+F
Find text


Ctrl+G
Go to line number


Home
Move to line start


End
Move to line end


Ctrl+Home
Move to file start


Ctrl+End
Move to file end


Page Up
Scroll page up


Page Down
Scroll page down


Arrow Keys
Move cursor


AI Assistant



Shortcut
Action



Ctrl+L
Ask AI (opens prompt)


Ctrl+K
Insert AI response at cursor


Esc
Cancel AI request



🤖 AI Integration
GoEdit integrates seamlessly with Ollama for AI-powered code assistance.
Setup Ollama

Install Ollama:
# macOS
brew install ollama

# Linux
curl https://ollama.ai/install.sh | sh

# Windows
# Download from https://ollama.ai/download


Start Ollama Server:
ollama serve


Pull a Model:
# General purpose
ollama pull llama2

# Code-focused
ollama pull codellama

# Other options
ollama pull mistral
ollama pull phi



Using AI Features
Standard Mode (Default)
goedit -model codellama mycode.py


Press Ctrl+L to open AI prompt
Type your question (e.g., "Add error handling to this function")
Press Enter to send
Wait for complete response
Press Ctrl+K to insert at cursor

Streaming Mode
goedit -stream -model llama2 document.md


Press Ctrl+L to open AI prompt
Type your question
Press Enter
Watch response appear in real-time
Press Ctrl+K to insert

Example AI Prompts

"Explain this code"
"Add error handling"
"Write unit tests for this function"
"Optimize this algorithm"
"Add documentation comments"
"Convert this to Python"
"Find potential bugs"


📁 File Structure
goedit/
├── cursor.go       # Cursor position management
├── buffer.go       # Text buffer with undo/redo
├── clipboard.go    # OS clipboard integration
├── tabs.go         # Multi-file tab management
├── ollama.go       # AI integration with streaming
├── main.go         # Main editor logic and UI
├── go.mod          # Go module definition
└── README.md       # This file


🎨 User Interface
Tab Bar
┌─────────────────────────────────────────────────────┐
│ 1:main.go [+] │ 2:utils.go │ 3:README.md           │
└─────────────────────────────────────────────────────┘


Active tab is highlighted
[+] indicates unsaved changes
Tab number for quick reference

Status Bar
┌─────────────────────────────────────────────────────┐
│ Ctrl+Q: Quit | Ctrl+S: Save | Ctrl+L: AI           │
│ main.go [+] | Ln 42/150 | Col 12 | Tab 1/3         │
└─────────────────────────────────────────────────────┘


Top line: Current status message
Bottom line: File info and position
Yellow background during AI processing


🔧 Configuration
Environment Variables
GoEdit respects the following environment variables:
# Set default Ollama URL
export OLLAMA_HOST=http://localhost:11434

# Terminal settings (handled by tcell)
export TERM=xterm-256color

Clipboard Support



OS
Primary Method
Fallback



macOS
pbcopy/pbpaste
Internal buffer


Linux
xclip or xsel
Internal buffer


Windows
PowerShell clipboard
Internal buffer


Note: On Linux, install xclip or xsel for system clipboard support:
# Debian/Ubuntu
sudo apt-get install xclip

# Fedora/RHEL
sudo dnf install xclip

# Arch
sudo pacman -S xclip


🐛 Troubleshooting
AI Features Not Working
Problem: "Cannot connect to Ollama"
Solution:
# Check if Ollama is running
curl http://localhost:11434/api/tags

# Start Ollama
ollama serve

# Verify model is installed
ollama list
ollama pull llama2

Clipboard Not Working
Problem: Copy/paste doesn't work with system clipboard
Solution:
# Linux - Install clipboard utilities
sudo apt-get install xclip

# macOS - Should work out of the box
# Windows - Ensure PowerShell is available

# Test clipboard manually
echo "test" | xclip -selection clipboard
xclip -selection clipboard -o

Terminal Display Issues
Problem: Strange characters or rendering issues
Solution:
# Set proper terminal type
export TERM=xterm-256color

# Or try
export TERM=screen-256color

# Resize terminal
# Press Ctrl+L to refresh in GoEdit

File Save Errors
Problem: "Permission denied" when saving
Solution:
# Check file permissions
ls -la yourfile.txt

# Make file writable
chmod u+w yourfile.txt

# Check directory permissions
ls -la $(dirname yourfile.txt)


🚦 System Requirements
Minimum Requirements

OS: Linux, macOS, Windows, BSD
RAM: 10 MB
Terminal: Any terminal with UTF-8 support
Go: 1.18+ (for building)

Recommended

Terminal: iTerm2 (macOS), Alacritty, or Windows Terminal
Font: Monospace font with Unicode support
Ollama: For AI features
Clipboard: xclip/xsel (Linux)


📊 Performance

Startup Time: < 50ms
Memory Usage: ~15 MB (without AI)
File Size Limit: 1 GB per file
Undo Levels: 50 (configurable in source)
Max Tabs: Limited by available memory


🔒 Security
Safe File Operations

Atomic Writes: Files are written to temporary files first
Backup on Error: Original file preserved if save fails
Permission Preservation: File permissions maintained

AI Privacy

Local Processing: All AI runs through your local Ollama instance
No Cloud: No data sent to external servers (unless you configure remote Ollama)
Model Control: You control which models are used


🤝 Contributing
Contributions are welcome! Please follow these guidelines:

Fork the Repository
Create a Feature Branchgit checkout -b feature/amazing-feature


Commit Your Changesgit commit -m 'Add amazing feature'


Push to Branchgit push origin feature/amazing-feature


Open a Pull Request

Code Style

Follow standard Go conventions
Run go fmt before committing
Add comments for complex logic
Write tests for new features


📝 License
Copyright © Prof. Dr. Michael Stal, 2025All Rights Reserved
This software is proprietary and confidential. Unauthorized copying, distribution, or use of this software, via any medium, is strictly prohibited.

🙏 Acknowledgments

tcell - Terminal cell library for Go
Ollama - Local LLM runtime
Go Team - For the amazing Go language


📞 Support
Getting Help

Issues: GitHub Issues
Discussions: GitHub Discussions
Email: your.email@example.com

Reporting Bugs
Please include:

GoEdit version (goedit -version)
Operating system and version
Terminal emulator
Steps to reproduce
Expected vs actual behavior


🗺️ Roadmap
Planned Features

 Syntax highlighting
 Line numbers (toggle)
 Split panes
 Macro recording
 Plugin system
 Configuration file support
 Multiple cursor support
 Git integration
 File tree sidebar
 Themes support

Under Consideration

 Mouse support enhancement
 Remote file editing (SSH)
 Collaborative editing
 LSP (Language Server Protocol) support
 Fuzzy file finder
 Project-wide search and replace


📚 Examples
Example 1: Quick Note Taking
# Start editor
goedit notes.txt

# Type your notes
# Press Ctrl+S to save
# Press Ctrl+Q to quit

Example 2: Multi-File Editing
# Open multiple files
goedit src/main.go src/utils.go README.md

# Use Tab to switch between files
# Edit each file
# Ctrl+S saves current file
# Ctrl+W closes current tab

Example 3: AI-Assisted Coding
# Start with AI model
goedit -stream -model codellama main.go

# Write some code
# Press Ctrl+L
# Ask: "Add error handling to this function"
# Watch response stream in
# Press Ctrl+K to insert

Example 4: Clipboard Workflow
# Open file
goedit document.txt

# Copy entire file: Ctrl+A
# Copy single line: Ctrl+C
# Cut line: Ctrl+X
# Paste: Ctrl+V

# Clipboard works with other applications!


🔍 FAQ
Q: How do I change the AI model?A: Use the -model flag: goedit -model codellama file.go
Q: Can I use GoEdit without Ollama?A: Yes! All features except AI work without Ollama.
Q: How do I enable streaming AI?A: Use the -stream flag: goedit -stream file.txt
Q: What's the difference between streaming and non-streaming AI?A: Streaming shows responses in real-time; non-streaming waits for complete response.
Q: How many files can I open at once?A: Limited only by available memory. Tested with 50+ tabs.
Q: Does it support Unicode?A: Yes! Full UTF-8 support for all languages.
Q: Can I customize keyboard shortcuts?A: Not yet, but it's on the roadmap!
Q: How do I save without a filename?A: Press Ctrl+S, and you'll be prompted to enter a filename.
Q: What happens if I try to quit with unsaved changes?A: GoEdit warns you and requires a second Ctrl+Q to confirm.
Q: Can I use this over SSH?A: Yes! Works perfectly in SSH sessions.

📈 Changelog
Version 2.0.0 (2025-01-XX)
New Features:

✨ Multi-file tab support
🤖 AI integration with Ollama
📋 System clipboard integration (macOS, Linux, Windows)
🌊 Streaming AI responses
↩️ 50 levels of undo/redo
🎨 Enhanced UI with tab bar

Improvements:

⚡ Faster rendering
🔒 Atomic file saves
🛡️ Better error handling
🧵 Thread-safe operations
📊 Improved status messages

Bug Fixes:

Fixed race conditions in AI processing
Fixed clipboard timeout issues
Fixed cursor positioning edge cases
Fixed tab overflow rendering


🌟 Star History
If you find GoEdit useful, please consider giving it a star on GitHub! ⭐

Made with ❤️ and Go
GoEdit v2.0 - The Modern Terminal EditorCopyright © Prof. Dr. Michael Stal, 2025
```
