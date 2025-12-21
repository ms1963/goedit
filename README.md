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
git clone https://github.com/ms1963/goedit.git

- Configuration
- AI Integration Setup
- Examples
- Troubleshooting
- Building from Source
- Contributing
- License


## ✨ Features




### Core Editing

- 🎯 Intuitive Interface - Familiar keyboard shortcuts
- 📝 Full Text Editing - Insert, delete, copy, cut, paste
- ↩️ Undo/Redo - 50 levels of history
- 🔍 Search - Case-insensitive with wraparound
- 📍 Go to Line - Quick navigation
- 📋 Clipboard - Internal copy/paste support




### Advanced Features

- 🤖 AI Integration - Built-in Ollama LLM support
- 💾 Atomic Saves - Safe file writing
- 🌐 Cross-Platform - Windows, Linux, macOS
- 📏 Status Bar - Real-time file info
- 🎨 Smart Indentation - 4-space tabs
- 🔒 Quit Protection - Unsaved change warnings





### Performance

- ⚡ Fast - Efficient rendering and minimal memory footprint
- 📄 Large Files - Handles files up to 1MB line length
- 🚀 Responsive - Smooth scrolling and instant feedback


### 📦 Prerequisites
Required


Go 1.21+
Download Go
go version


Terminal
Any modern terminal emulator

Windows Terminal
iTerm2
GNOME Terminal


Git (optional)
For cloning repository
git --version


Optional (for AI features)

Ollama - Install Ollama
LLM Model - Any Ollama-compatible model (llama2, codellama, etc.)


## 🚀 Installation
From Source
Step 1: Clone or Download
### Option A: Using git
git clone https://github.com/ms1963/goedit.git
cd goedit

### Option B: Download ZIP and extract
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
git clone https://github.com/ms1963/goedit.git &amp;&amp; cd goedit &amp;&amp; go mod init goedit &amp;&amp; go get github.com/gdamore/tcell/v2 &amp;&amp; go build -o goedit


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
def calc(x,y,z):
    return x+y*z-x/y

#Ask AI: "Refactor this function with proper naming and documentation"
#Ctrl+L → Enter → Ctrl+K

#Result:
def calculate_expression(x: float, y: float, z: float) -> float:
    """
    Calculate the expression: x + (y * z) - (x / y)
    
    Args:
        x: First operand
        y: Second operand (divisor, must not be zero)
        z: Third operand
        
    Returns:
        Result of the calculation
        
    Raises:
        ZeroDivisionError: If y is zero
    """
    if y == 0:
        raise ZeroDivisionError("Cannot divide by zero")
    
    multiplication = y * z
    division = x / y
    result = x + multiplication - division
    
    return result
```


## 🔧 Troubleshooting
Common Issues and Solutions

### ❌ "Failed to initialize screen"

Problem: Terminal doesn't support required features
Solutions:
#### Set proper TERM variable
export TERM=xterm-256color

#### Or try
export TERM=screen-256color

#### On Windows, use Windows Terminal or PowerShell
- Avoid old CMD.exe

Verify terminal:
```
echo $TERM
#Should show: xterm-256color or similar
```



### ❌ "Failed to connect to Ollama"

Problem: Ollama is not running or wrong URL

Solutions:
#### 1. Check if Ollama is running
curl http://localhost:11434/api/tags

#### 2. Start Ollama if not running
ollama serve

#### 3. Check firewall settings
# Ensure port 11434 is not blocked

#### 4. Verify URL and port
goedit -ollama http://localhost:11434 file.txt

####5. Check Ollama logs
ollama logs

Test Ollama directly:
curl http://localhost:11434/api/generate -d '{
  "model": "llama2",
  "prompt": "Hello",
  "stream": false
}'




### ❌ "Model not found"

Problem: Requested model not installed

Solutions:
#### 1. List installed models
ollama list

#### 2. Pull the model
ollama pull llama2

#### 3. Use installed model
goedit -model llama2 file.txt

#### 4. Check model name spelling
#Correct: llama2, codellama, mistral
#Incorrect: llama-2, code-llama




### ❌ Terminal too small

Problem: Terminal window is too small

Solutions:

Resize terminal to at least 80 columns × 24 rows
GoEdit requires minimum 3 lines height
Use fullscreen mode: F11 (most terminals)

Check terminal size:
echo "Columns: $COLUMNS, Rows: $LINES"




### ❌ Characters not displaying correctly

Problem: Encoding issues

Solutions:
#Set UTF-8 encoding
export LANG=en_US.UTF-8
export LC_ALL=en_US.UTF-8

#### On Windows, ensure terminal supports UTF-8
#### Use Windows Terminal (recommended)

Verify encoding:
locale
#Should show UTF-8 encoding




### ❌ Can't save file

Problem: Permission denied or directory doesn't exist

Solutions:
#### 1. Check permissions
ls -la /path/to/file

#### 2. Create directory if needed
mkdir -p /path/to/directory

#### 3. Check write permissions
touch /path/to/test.txt
rm /path/to/test.txt

#### 4. Use correct path
- Absolute: /home/user/file.txt
- Relative: ./file.txt

#### 5. For system files, use sudo (not recommended for regular editing)
sudo goedit /etc/config




### ❌ Slow AI responses

Problem: Model is large or system is slow

Solutions:
#### 1. Use smaller, faster model
goedit -model mistral file.txt

#### 2. Use quantized model (smaller, faster)
```
ollama pull llama2:7b-q4_0
goedit -model llama2:7b-q4_0 file.txt
```

#### 3. Check system resources
#Ensure enough RAM (8GB+ recommended for 7B models)

#### 4. Close other applications

####5. Use GPU if available
#Ollama automatically uses GPU when available

Model performance comparison:

- mistral - Fastest
- llama2 - Fast
- llama2:13b - Medium
- llama2:70b - Slow (requires powerful hardware)




### ❌ Build errors

Problem: Compilation fails

Solutions:
#### 1. Ensure Go version is 1.21+
go version

#### 2.Clean and rebuild
```
go clean
rm go.sum
go mod tidy
go build
```

#### 3.Update dependencies
```
go get -u ./...
go mod tidy

#4. Check for syntax errors
go fmt ./...
go vet ./...

#5. Verify all files are present
#Required: main.go, buffer.go, cursor.go, ollama.go, go.mod
```


### Debug Mode
```
#Test Ollama connection
curl -X POST http://localhost:11434/api/generate \
  -H "Content-Type: application/json" \
  -d '{
    "model": "llama2",
    "prompt": "Say hello",
    "stream": false
  }'

#Check Ollama logs
ollama logs

#Verify Go installation
go env

#Test terminal capabilities
tput colors  # Should show 256 or more

Getting Help
#Show help
goedit -help

#Show version
goedit -version

#Check Go version
go version

#Check Ollama status
ollama list
```

### Performance Tips

For large files:

- Use Ctrl+G to jump to specific lines
- Use Ctrl+F to find content
- Save frequently


For AI features:

- Use smaller models for faster responses
- Be specific in prompts
- Use codellama for code, llama2 for text


For slow terminals:

- Reduce terminal font size
- Use hardware acceleration
- Close unnecessary applications




🛠️ Building from Source
Development Setup
#Clone repository
git clone https://github.com/ms1963/goedit.git
cd goedit


#Install dependencies
go mod download


#Run without building
go run . test.txt


#Build for development
go build -o goedit


#Build with debug info
go build -gcflags="all=-N -l" -o goedit-debug


Build Optimizations
#Optimized build (smaller binary)
go build -ldflags="-s -w" -o goedit


#With version info
VERSION="1.0.0"
go build -ldflags="-X main.version=$VERSION -s -w" -o goedit


#Static binary (Linux - no external dependencies)
CGO_ENABLED=0 go build -ldflags="-s -w" -o goedit


#Verify binary size
ls -lh goedit

Testing
#Run tests (if implemented)
go test ./...

#Verbose output
go test -v ./...

#Test specific package
go test -v ./buffer

#Run benchmarks
go test -bench=. ./...

#Test coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

### Code Quality

#Format code
go fmt ./...


#Lint code
go vet ./...


#Static analysis (install golangci-lint first)
golangci-lint run


#Check for common mistakes
staticcheck ./...


Project Structure

```
goedit/
├── main.go          # Main entry point and editor logic
├── buffer.go        # Text buffer management
├── cursor.go        # Cursor position handling
├── ollama.go        # Ollama API client
├── go.mod           # Go module definition
├── go.sum           # Dependency checksums
├── README.md        # This file
├── LICENSE          # MIT License
└── .gitignore       # Git ignore rules

File Descriptions
```


File
Lines
Purpose



main.go
~800
Editor UI, keyboard handling, rendering


buffer.go
~400
Text buffer, undo/redo, file I/O


cursor.go
~15
Cursor position management


ollama.go
~150
Ollama API client, LLM integration


Dependencies
// Direct dependencies
github.com/gdamore/tcell/v2  // Terminal handling

// Indirect dependencies
github.com/gdamore/encoding
github.com/lucasb-eyer/go-colorful
github.com/mattn/go-runewidth
github.com/rivo/uniseg
golang.org/x/sys
golang.org/x/term
golang.org/x/text


🤝 Contributing
We welcome contributions! Here's how you can help:
Ways to Contribute

🐛 Report bugs - Open an issue with details




💡 Suggest features - Share your ideas


📖 Improve docs - Fix typos, add examples


🔧 Submit PRs - Fix bugs or add features


⭐ Star the repo - Show your support

### Development Process

Fork the repository
#Click "Fork" on GitHub
git clone https://github.com/ms1963/goedit.git
cd goedit


Create a feature branch
git checkout -b feature/amazing-feature


Make your changes
#Edit files
#Test thoroughly
go test ./...
go build


Commit your changes
git add .
git commit -m "feat: Add amazing feature"

#### Commit message format:

- feat: New feature
- fix: Bug fix
- docs: Documentation
- style: Formatting
- refactor: Code restructuring
- test: Tests
- chore: Maintenance


#### Push to your fork
git push origin feature/amazing-feature


#### Open a Pull Request

Go to GitHub
Click "New Pull Request"
Describe your changes
Link related issues



### Development Guidelines

- Follow Go best practices


- Add comments for complex logic


- Test on Windows, Linux, and macOS


- Update README for new features


✅ Keep dependencies minimal


✅ Write clear commit messages


✅ Add tests for new features



Code Style
#Format code
go fmt ./...

#Check for issues
go vet ./...

#Run linter
golangci-lint run

Testing Checklist
Before submitting PR:

 Code builds without errors
 All tests pass
 Tested on target platforms
 Documentation updated
 No breaking changes (or documented)
 Commit messages are clear


### 📄 License
MIT License
Copyright (c) 2025 Prof. Dr. Michael Stal


Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

🙏 Acknowledgments

tcell - Excellent terminal handling library


Ollama - Making local LLMs accessible


Inspired by: nano, vim, emacs, and modern editors


Community: All contributors and users



📞 Support
Get Help

📚 Documentation: This README


🐛 Bug Reports: GitHub Issues


💬 Discussions: GitHub Discussions


📧 Email: michael.stal@gmail.com


### Useful Links

Go Documentation


Ollama Documentation


tcell Documentation


🗺️ Roadmap
Planned Features

 Syntax highlighting - Language-specific coloring

 
 Multiple file tabs - Edit multiple files simultaneously

 
 System clipboard - Integration with OS clipboard

 
 Mouse support - Click to position cursor

 
 Configuration file - Persistent settings

 
 Plugin system - Extensibility

 
 Line numbers - Optional line number display

 
 Code folding - Collapse/expand code blocks

 
 Git integration - Show git status, diff

 
 Themes - Customizable color schemes

 
 Auto-completion - Context-aware suggestions

 
 Bracket matching - Highlight matching brackets

 
 Multiple cursors - Edit multiple locations

 
 Regex search - Advanced search patterns
 

Version History

- v1.0.0 (2025-12) - Initial release

  
- Basic text editing


- Undo/redo support


- Search functionality


- Ollama integration


- Cross-platform support






Made with ❤️ by the GoEdit Team



Happy Editing! 🚀



⬆ Back to Top
