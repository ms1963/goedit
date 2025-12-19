
package main

import (
        "bufio"
        "fmt"
        "os"
        "path/filepath"
        "strings"
)

// Buffer represents the text buffer
type Buffer struct {
        lines        []string
        filename     string
        modified     bool
        history      []BufferState
        historyIndex int
}

// BufferState represents a snapshot of the buffer for undo/redo
type BufferState struct {
        lines     []string
        cursorRow int
        cursorCol int
}

// NewBuffer creates a new buffer
func NewBuffer(filename string) (*Buffer, error) {
        b := &Buffer{
                lines:        []string{""},
                filename:     filename,
                modified:     false,
                history:      make([]BufferState, 0, 50),
                historyIndex: -1,
        }

        if filename != "" {
                if err := b.Load(); err != nil && !os.IsNotExist(err) {
                        return nil, err
                }
        }

        b.SaveState(0, 0)
        return b, nil
}

// Load loads the file into the buffer
func (b *Buffer) Load() error {
        file, err := os.Open(b.filename)
        if err != nil {
                return err
        }
        defer file.Close()

        b.lines = make([]string, 0)
        scanner := bufio.NewScanner(file)

        // Increase buffer size for long lines
        const maxCapacity = 1024 * 1024
        buf := make([]byte, maxCapacity)
        scanner.Buffer(buf, maxCapacity)

        for scanner.Scan() {
                line := scanner.Text()
                // Remove carriage returns for cross-platform compatibility
                line = strings.TrimRight(line, "\r")
                b.lines = append(b.lines, line)
        }

        if err := scanner.Err(); err != nil {
                return err
        }

        if len(b.lines) == 0 {
                b.lines = []string{""}
        }

        b.modified = false
        return nil
}

// Save saves the buffer to file
func (b *Buffer) Save() error {
        if b.filename == "" {
                return fmt.Errorf("no filename specified")
        }

        // Create directory if needed
        dir := filepath.Dir(b.filename)
        if dir != "" && dir != "." {
                if err := os.MkdirAll(dir, 0755); err != nil {
                        return fmt.Errorf("failed to create directory: %w", err)
                }
        }

        // Write to temp file first for atomic save
        tempFile := b.filename + ".tmp"
        file, err := os.Create(tempFile)
        if err != nil {
                return fmt.Errorf("failed to create temp file: %w", err)
        }

        writer := bufio.NewWriter(file)
        for i, line := range b.lines {
                if i > 0 {
                        if _, err := writer.WriteString("\n"); err != nil {
                                file.Close()
                                os.Remove(tempFile)
                                return err
                        }
                }
                if _, err := writer.WriteString(line); err != nil {
                        file.Close()
                        os.Remove(tempFile)
                        return err
                }
        }

        if err := writer.Flush(); err != nil {
                file.Close()
                os.Remove(tempFile)
                return err
        }

        if err := file.Close(); err != nil {
                os.Remove(tempFile)
                return err
        }

        // Atomic rename
        if err := os.Rename(tempFile, b.filename); err != nil {
                os.Remove(tempFile)
                return fmt.Errorf("failed to save file: %w", err)
        }

        b.modified = false
        return nil
}

// InsertChar inserts a character at position
func (b *Buffer) InsertChar(row, col int, ch rune) {
        if !b.isValidRow(row) {
                return
        }

        line := b.lines[row]
        col = b.clampCol(col, len(line))

        newLine := line[:col] + string(ch) + line[col:]
        b.lines[row] = newLine
        b.modified = true
}

// DeleteChar deletes a character at position (backspace)
func (b *Buffer) DeleteChar(row, col int) {
        if !b.isValidRow(row) {
                return
        }

        line := b.lines[row]
        if col > 0 && col <= len(line) {
                b.lines[row] = line[:col-1] + line[col:]
                b.modified = true
        } else if col == 0 && row > 0 {
                // Merge with previous line
                prevLine := b.lines[row-1]
                b.lines[row-1] = prevLine + line
                b.lines = append(b.lines[:row], b.lines[row+1:]...)
                if len(b.lines) == 0 {
                        b.lines = []string{""}
                }
                b.modified = true
        }
}

// DeleteCharForward deletes character at cursor (delete key)
func (b *Buffer) DeleteCharForward(row, col int) {
        if !b.isValidRow(row) {
                return
        }

        line := b.lines[row]
        if col >= 0 && col < len(line) {
                b.lines[row] = line[:col] + line[col+1:]
                b.modified = true
        }
}

// InsertNewline inserts a newline at position
func (b *Buffer) InsertNewline(row, col int) {
        if !b.isValidRow(row) {
                return
        }

        line := b.lines[row]
        col = b.clampCol(col, len(line))

        // Split line at cursor
        newLines := make([]string, len(b.lines)+1)
        copy(newLines, b.lines[:row])
        newLines[row] = line[:col]
        newLines[row+1] = line[col:]
        copy(newLines[row+2:], b.lines[row+1:])

        b.lines = newLines
        b.modified = true
}

// GetLine returns a line
func (b *Buffer) GetLine(row int) string {
        if !b.isValidRow(row) {
                return ""
        }
        return b.lines[row]
}

// LineCount returns the number of lines
func (b *Buffer) LineCount() int {
        if len(b.lines) == 0 {
                return 1
        }
        return len(b.lines)
}

// GetText returns all text
func (b *Buffer) GetText() string {
        return strings.Join(b.lines, "\n")
}

// SetText replaces all text
func (b *Buffer) SetText(text string) {
        b.lines = strings.Split(text, "\n")
        if len(b.lines) == 0 {
                b.lines = []string{""}
        }
        b.modified = true
}

// InsertText inserts text at cursor position
func (b *Buffer) InsertText(row, col int, text string) {
        if !b.isValidRow(row) {
                return
        }

        if text == "" {
                return
        }

        col = b.clampCol(col, len(b.lines[row]))
        lines := strings.Split(text, "\n")

        if len(lines) == 0 {
                return
        }

        if len(lines) == 1 {
                // Single line insert
                line := b.lines[row]
                b.lines[row] = line[:col] + text + line[col:]
        } else {
                // Multi-line insert
                line := b.lines[row]
                firstPart := line[:col] + lines[0]
                lastPart := lines[len(lines)-1] + line[col:]

                newLines := make([]string, 0, len(b.lines)+len(lines)-1)
                newLines = append(newLines, b.lines[:row]...)
                newLines = append(newLines, firstPart)
                if len(lines) > 2 {
                        newLines = append(newLines, lines[1:len(lines)-1]...)
                }
                newLines = append(newLines, lastPart)
                newLines = append(newLines, b.lines[row+1:]...)

                b.lines = newLines
        }
        b.modified = true
}

// DeleteLine deletes a line
func (b *Buffer) DeleteLine(row int) {
        if !b.isValidRow(row) {
                return
        }

        b.lines = append(b.lines[:row], b.lines[row+1:]...)
        if len(b.lines) == 0 {
                b.lines = []string{""}
        }
        b.modified = true
}

// AppendToLine appends text to a line
func (b *Buffer) AppendToLine(row int, text string) {
        if !b.isValidRow(row) {
                return
        }
        b.lines[row] += text
        b.modified = true
}

// SaveState saves state for undo/redo
func (b *Buffer) SaveState(cursorRow, cursorCol int) {
        // Validate cursor position
        if cursorRow < 0 {
                cursorRow = 0
        }
        if cursorCol < 0 {
                cursorCol = 0
        }

        // Remove any states after current index
        if b.historyIndex >= 0 && b.historyIndex < len(b.history)-1 {
                b.history = b.history[:b.historyIndex+1]
        }

        // Create snapshot
        linesCopy := make([]string, len(b.lines))
        copy(linesCopy, b.lines)

        state := BufferState{
                lines:     linesCopy,
                cursorRow: cursorRow,
                cursorCol: cursorCol,
        }

        b.history = append(b.history, state)
        b.historyIndex++

        // Limit history size to prevent memory issues
        const maxHistory = 50
        if len(b.history) > maxHistory {
                overflow := len(b.history) - maxHistory
                b.history = b.history[overflow:]
                b.historyIndex = len(b.history) - 1
        }
}

// Undo reverts to previous state
func (b *Buffer) Undo() (int, int, bool) {
        if b.historyIndex <= 0 {
                return 0, 0, false
        }

        b.historyIndex--
        state := b.history[b.historyIndex]
        b.lines = make([]string, len(state.lines))
        copy(b.lines, state.lines)
        
        if len(b.lines) == 0 {
                b.lines = []string{""}
        }
        
        b.modified = true

        return state.cursorRow, state.cursorCol, true
}

// Redo moves forward in history
func (b *Buffer) Redo() (int, int, bool) {
        if b.historyIndex >= len(b.history)-1 {
                return 0, 0, false
        }

        b.historyIndex++
        state := b.history[b.historyIndex]
        b.lines = make([]string, len(state.lines))
        copy(b.lines, state.lines)
        
        if len(b.lines) == 0 {
                b.lines = []string{""}
        }
        
        b.modified = true

        return state.cursorRow, state.cursorCol, true
}

// Helper methods

func (b *Buffer) isValidRow(row int) bool {
        return row >= 0 && row < len(b.lines) && len(b.lines) > 0
}

func (b *Buffer) clampCol(col, maxCol int) int {
        if col < 0 {
                return 0
        }
        if col > maxCol {
                return maxCol
        }
        return col
}
