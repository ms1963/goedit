// GoEdit v2.0
// Copyright © Prof. Dr. Michael Stal, 2025
// All rights reserved.

package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

const maxUndoLevels = 50

type Buffer struct {
    lines     []string
    filename  string
    modified  bool
    undoStack []BufferState
    redoStack []BufferState
}

type BufferState struct {
    lines     []string
    cursorRow int
    cursorCol int
}

func NewBuffer(filename string) (*Buffer, error) {
    b := &Buffer{
        lines:     []string{""},
        filename:  filename,
        modified:  false,
        undoStack: make([]BufferState, 0, maxUndoLevels),
        redoStack: make([]BufferState, 0, maxUndoLevels),
    }

    if filename != "" {
        if err := b.Load(); err != nil {
            if !os.IsNotExist(err) {
                return nil, err
            }
        }
    }

    return b, nil
}

func (b *Buffer) Load() error {
    file, err := os.Open(b.filename)
    if err != nil {
        return err
    }
    defer file.Close()

    b.lines = []string{}
    scanner := bufio.NewScanner(file)
    
    const maxCapacity = 1024 * 1024
    buf := make([]byte, maxCapacity)
    scanner.Buffer(buf, maxCapacity)

    for scanner.Scan() {
        b.lines = append(b.lines, scanner.Text())
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

func (b *Buffer) Save() error {
    if b.filename == "" {
        return fmt.Errorf("no filename specified")
    }

    tempFile := b.filename + ".tmp"
    file, err := os.Create(tempFile)
    if err != nil {
        return err
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

    if err := os.Rename(tempFile, b.filename); err != nil {
        os.Remove(tempFile)
        return err
    }

    b.modified = false
    return nil
}

func (b *Buffer) GetLine(row int) string {
    if row < 0 || row >= len(b.lines) {
        return ""
    }
    return b.lines[row]
}

func (b *Buffer) LineCount() int {
    return len(b.lines)
}

func (b *Buffer) GetText() string {
    return strings.Join(b.lines, "\n")
}

func (b *Buffer) InsertChar(row, col int, ch rune) {
    if row < 0 || row >= len(b.lines) {
        return
    }

    line := b.lines[row]
    if col < 0 {
        col = 0
    }
    if col > len(line) {
        col = len(line)
    }

    newLine := line[:col] + string(ch) + line[col:]
    b.lines[row] = newLine
    b.modified = true
}

func (b *Buffer) DeleteChar(row, col int) {
    if row < 0 || row >= len(b.lines) {
        return
    }

    if col > 0 {
        line := b.lines[row]
        if col > len(line) {
            col = len(line)
        }
        b.lines[row] = line[:col-1] + line[col:]
        b.modified = true
    } else if row > 0 {
        prevLine := b.lines[row-1]
        currentLine := b.lines[row]
        b.lines[row-1] = prevLine + currentLine
        b.lines = append(b.lines[:row], b.lines[row+1:]...)
        b.modified = true
    }
}

func (b *Buffer) DeleteCharForward(row, col int) {
    if row < 0 || row >= len(b.lines) {
        return
    }

    line := b.lines[row]
    if col >= 0 && col < len(line) {
        b.lines[row] = line[:col] + line[col+1:]
        b.modified = true
    }
}

func (b *Buffer) InsertNewline(row, col int) {
    if row < 0 || row >= len(b.lines) {
        return
    }

    line := b.lines[row]
    if col < 0 {
        col = 0
    }
    if col > len(line) {
        col = len(line)
    }

    before := line[:col]
    after := line[col:]

    b.lines[row] = before
    newLines := make([]string, len(b.lines)+1)
    copy(newLines, b.lines[:row+1])
    newLines[row+1] = after
    copy(newLines[row+2:], b.lines[row+1:])
    b.lines = newLines

    b.modified = true
}

func (b *Buffer) DeleteLine(row int) {
    if row < 0 || row >= len(b.lines) {
        return
    }

    if len(b.lines) == 1 {
        b.lines[0] = ""
    } else {
        b.lines = append(b.lines[:row], b.lines[row+1:]...)
    }
    b.modified = true
}

func (b *Buffer) AppendToLine(row int, text string) {
    if row < 0 || row >= len(b.lines) {
        return
    }
    b.lines[row] += text
    b.modified = true
}

func (b *Buffer) InsertText(row, col int, text string) {
    if row < 0 || row >= len(b.lines) {
        return
    }

    line := b.lines[row]
    if col < 0 {
        col = 0
    }
    if col > len(line) {
        col = len(line)
    }

    lines := strings.Split(text, "\n")
    if len(lines) == 1 {
        b.lines[row] = line[:col] + text + line[col:]
    } else {
        before := line[:col]
        after := line[col:]

        b.lines[row] = before + lines[0]

        newLines := make([]string, len(b.lines)+len(lines)-1)
        copy(newLines, b.lines[:row+1])
        
        for i := 1; i < len(lines)-1; i++ {
            newLines[row+i] = lines[i]
        }
        
        newLines[row+len(lines)-1] = lines[len(lines)-1] + after
        copy(newLines[row+len(lines):], b.lines[row+1:])
        b.lines = newLines
    }

    b.modified = true
}

func (b *Buffer) SaveState(cursorRow, cursorCol int) {
    linesCopy := make([]string, len(b.lines))
    copy(linesCopy, b.lines)

    state := BufferState{
        lines:     linesCopy,
        cursorRow: cursorRow,
        cursorCol: cursorCol,
    }

    b.undoStack = append(b.undoStack, state)
    if len(b.undoStack) > maxUndoLevels {
        b.undoStack = b.undoStack[1:]
    }

    b.redoStack = b.redoStack[:0]
}

func (b *Buffer) Undo() (int, int, bool) {
    if len(b.undoStack) == 0 {
        return 0, 0, false
    }

    currentState := BufferState{
        lines: make([]string, len(b.lines)),
    }
    copy(currentState.lines, b.lines)
    b.redoStack = append(b.redoStack, currentState)

    state := b.undoStack[len(b.undoStack)-1]
    b.undoStack = b.undoStack[:len(b.undoStack)-1]

    b.lines = make([]string, len(state.lines))
    copy(b.lines, state.lines)
    b.modified = true

    return state.cursorRow, state.cursorCol, true
}

func (b *Buffer) Redo() (int, int, bool) {
    if len(b.redoStack) == 0 {
        return 0, 0, false
    }

    currentState := BufferState{
        lines: make([]string, len(b.lines)),
    }
    copy(currentState.lines, b.lines)
    b.undoStack = append(b.undoStack, currentState)

    state := b.redoStack[len(b.redoStack)-1]
    b.redoStack = b.redoStack[:len(b.redoStack)-1]

    b.lines = make([]string, len(state.lines))
    copy(b.lines, state.lines)
    b.modified = true

    return state.cursorRow, state.cursorCol, true
}
