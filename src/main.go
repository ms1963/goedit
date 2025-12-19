package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "path/filepath"
    "strings"
    "sync"
    "time"

    "github.com/gdamore/tcell/v2"
)

const version = "1.0.0"

type Editor struct {
    screen       tcell.Screen
    buffer       *Buffer
    cursor       *Cursor
    offsetRow    int
    offsetCol    int
    width        int
    height       int
    statusMsg    string
    mode         EditorMode
    findQuery    string
    clipboard    string
    llmClient    *OllamaClient
    llmPrompt    string
    llmResponse  string
    inputBuffer  string
    quitAttempts int
    aiMutex      sync.Mutex
    aiInProgress bool
    aiCancel     chan bool
}

type EditorMode int

const (
    ModeNormal EditorMode = iota
    ModeFind
    ModeGoto
    ModeLLM
    ModeFilename
)

func NewEditor(filename, ollamaURL, model string) (*Editor, error) {
    screen, err := tcell.NewScreen()
    if err != nil {
        return nil, fmt.Errorf("failed to create screen: %w", err)
    }

    if err := screen.Init(); err != nil {
        return nil, fmt.Errorf("failed to initialize screen: %w", err)
    }

    screen.EnableMouse()
    screen.EnablePaste()
    screen.Clear()

    buffer, err := NewBuffer(filename)
    if err != nil {
        screen.Fini()
        return nil, fmt.Errorf("failed to create buffer: %w", err)
    }

    width, height := screen.Size()
    if height < 3 {
        screen.Fini()
        return nil, fmt.Errorf("terminal too small (need at least 3 lines)")
    }

    return &Editor{
        screen:       screen,
        buffer:       buffer,
        cursor:       NewCursor(),
        width:        width,
        height:       height - 2,
        statusMsg:    "Ctrl+Q: Quit | Ctrl+S: Save | Ctrl+L: AI | Ctrl+F: Find",
        mode:         ModeNormal,
        llmClient:    NewOllamaClient(ollamaURL, model),
        aiInProgress: false,
        aiCancel:     make(chan bool, 1),
    }, nil
}

func (e *Editor) checkOllamaSetup() error {
    if !e.llmClient.IsAvailable() {
        return fmt.Errorf("Ollama not running. Start with: ollama serve")
    }
    
    if err := e.llmClient.CheckModel(); err != nil {
        return err
    }
    
    return nil
}

func (e *Editor) Run() error {
    if e.screen == nil {
        return fmt.Errorf("screen not initialized")
    }

    defer e.screen.Fini()

    e.render()

    for {
        ev := e.screen.PollEvent()
        if ev == nil {
            continue
        }

        if !e.handleEvent(ev) {
            return nil
        }

        e.render()
    }
}

func (e *Editor) handleEvent(ev tcell.Event) bool {
    switch ev := ev.(type) {
    case *tcell.EventResize:
        e.width, e.height = ev.Size()
        if e.height < 3 {
            e.height = 3
        }
        e.height -= 2
        e.screen.Sync()
        return true

    case *tcell.EventKey:
        if ev.Key() == tcell.KeyEscape {
            e.aiMutex.Lock()
            if e.aiInProgress {
                select {
                case e.aiCancel <- true:
                default:
                }
                e.aiInProgress = false
                e.statusMsg = "AI request cancelled"
                e.aiMutex.Unlock()
                if e.mode == ModeLLM {
                    e.mode = ModeNormal
                }
                return true
            }
            e.aiMutex.Unlock()
        }

        return e.handleKey(ev)
    }

    return true
}

func (e *Editor) handleKey(ev *tcell.EventKey) bool {
    if ev == nil {
        return true
    }

    switch e.mode {
    case ModeFind:
        return e.handleFindMode(ev)
    case ModeGoto:
        return e.handleGotoMode(ev)
    case ModeLLM:
        return e.handleLLMMode(ev)
    case ModeFilename:
        return e.handleFilenameMode(ev)
    default:
        return e.handleNormalMode(ev)
    }
}

func (e *Editor) handleNormalMode(ev *tcell.EventKey) bool {
    mod := ev.Modifiers()

    switch ev.Key() {
    case tcell.KeyCtrlQ:
        e.aiMutex.Lock()
        inProgress := e.aiInProgress
        e.aiMutex.Unlock()

        if inProgress {
            e.statusMsg = "AI in progress. Press Esc to cancel, then Ctrl+Q to quit"
            return true
        }

        if e.buffer.modified && e.quitAttempts == 0 {
            e.statusMsg = "File modified! Press Ctrl+Q again to quit, Ctrl+S to save"
            e.quitAttempts++
            return true
        }
        return false

    case tcell.KeyCtrlS:
        e.quitAttempts = 0
        e.saveFile()

    case tcell.KeyCtrlF:
        e.mode = ModeFind
        e.inputBuffer = ""
        e.statusMsg = "Find: "

    case tcell.KeyCtrlG:
        e.mode = ModeGoto
        e.inputBuffer = ""
        e.statusMsg = "Go to line: "

    case tcell.KeyCtrlL:
        if err := e.checkOllamaSetup(); err != nil {
            e.statusMsg = fmt.Sprintf("AI unavailable: %v", err)
            return true
        }
        e.mode = ModeLLM
        e.inputBuffer = ""
        e.statusMsg = "Ask AI: "

    case tcell.KeyCtrlK:
        if e.llmResponse != "" {
            oldRow := e.cursor.Row
            oldCol := e.cursor.Col
            e.buffer.InsertText(e.cursor.Row, e.cursor.Col, e.llmResponse)

            lines := strings.Split(e.llmResponse, "\n")
            if len(lines) > 1 {
                e.cursor.Row = oldRow + len(lines) - 1
                if e.cursor.Row < 0 {
                    e.cursor.Row = 0
                }
                lastLine := lines[len(lines)-1]
                e.cursor.Col = len(lastLine)
            } else {
                e.cursor.Col = oldCol + len(e.llmResponse)
            }

            e.ensureCursorValid()
            e.buffer.SaveState(e.cursor.Row, e.cursor.Col)
            e.statusMsg = "AI response inserted at cursor"
        } else {
            e.statusMsg = "No AI response available. Use Ctrl+L to ask AI first"
        }

    case tcell.KeyCtrlA:
        e.clipboard = e.buffer.GetText()
        e.statusMsg = "All text copied to clipboard"

    case tcell.KeyCtrlC:
        if e.cursor.Row >= 0 && e.cursor.Row < e.buffer.LineCount() {
            line := e.buffer.GetLine(e.cursor.Row)
            e.clipboard = line
            e.statusMsg = "Current line copied to clipboard"
        }

    case tcell.KeyCtrlX:
        if e.cursor.Row >= 0 && e.cursor.Row < e.buffer.LineCount() {
            line := e.buffer.GetLine(e.cursor.Row)
            e.clipboard = line
            e.buffer.DeleteLine(e.cursor.Row)
            if e.cursor.Row >= e.buffer.LineCount() && e.cursor.Row > 0 {
                e.cursor.Row--
            }
            e.cursor.Col = 0
            e.ensureCursorValid()
            e.buffer.SaveState(e.cursor.Row, e.cursor.Col)
            e.statusMsg = "Current line cut to clipboard"
        }

    case tcell.KeyCtrlV:
        if e.clipboard != "" {
            oldRow := e.cursor.Row
            oldCol := e.cursor.Col
            e.buffer.InsertText(e.cursor.Row, e.cursor.Col, e.clipboard)

            lines := strings.Split(e.clipboard, "\n")
            if len(lines) > 1 {
                e.cursor.Row = oldRow + len(lines) - 1
                if e.cursor.Row < 0 {
                    e.cursor.Row = 0
                }
                lastLine := lines[len(lines)-1]
                e.cursor.Col = len(lastLine)
            } else {
                e.cursor.Col = oldCol + len(e.clipboard)
            }

            e.ensureCursorValid()
            e.buffer.SaveState(e.cursor.Row, e.cursor.Col)
            e.statusMsg = "Clipboard content pasted"
        } else {
            e.statusMsg = "Clipboard is empty"
        }

    case tcell.KeyCtrlZ:
        if row, col, ok := e.buffer.Undo(); ok {
            e.cursor.Row = row
            e.cursor.Col = col
            e.ensureCursorValid()
            e.statusMsg = "Undo successful"
        } else {
            e.statusMsg = "Nothing to undo"
        }

    case tcell.KeyCtrlY:
        if row, col, ok := e.buffer.Redo(); ok {
            e.cursor.Row = row
            e.cursor.Col = col
            e.ensureCursorValid()
            e.statusMsg = "Redo successful"
        } else {
            e.statusMsg = "Nothing to redo"
        }

    case tcell.KeyUp:
        if e.cursor.Row > 0 {
            e.cursor.Row--
            e.ensureCursorValid()
        }

    case tcell.KeyDown:
        if e.cursor.Row < e.buffer.LineCount()-1 {
            e.cursor.Row++
            e.ensureCursorValid()
        }

    case tcell.KeyLeft:
        if e.cursor.Col > 0 {
            e.cursor.Col--
        } else if e.cursor.Row > 0 {
            e.cursor.Row--
            e.cursor.Col = len(e.buffer.GetLine(e.cursor.Row))
        }

    case tcell.KeyRight:
        lineLen := len(e.buffer.GetLine(e.cursor.Row))
        if e.cursor.Col < lineLen {
            e.cursor.Col++
        } else if e.cursor.Row < e.buffer.LineCount()-1 {
            e.cursor.Row++
            e.cursor.Col = 0
        }

    case tcell.KeyHome:
        if mod&tcell.ModCtrl != 0 {
            e.cursor.Row = 0
            e.cursor.Col = 0
            e.ensureCursorValid()
            e.statusMsg = "Moved to start of file"
        } else {
            e.cursor.Col = 0
        }

    case tcell.KeyEnd:
        if mod&tcell.ModCtrl != 0 {
            e.cursor.Row = e.buffer.LineCount() - 1
            if e.cursor.Row < 0 {
                e.cursor.Row = 0
            }
            e.cursor.Col = len(e.buffer.GetLine(e.cursor.Row))
            e.ensureCursorValid()
            e.statusMsg = "Moved to end of file"
        } else {
            e.cursor.Col = len(e.buffer.GetLine(e.cursor.Row))
        }

    case tcell.KeyPgUp:
        e.cursor.Row -= e.height
        if e.cursor.Row < 0 {
            e.cursor.Row = 0
        }
        e.ensureCursorValid()

    case tcell.KeyPgDn:
        e.cursor.Row += e.height
        if e.cursor.Row >= e.buffer.LineCount() {
            e.cursor.Row = e.buffer.LineCount() - 1
        }
        if e.cursor.Row < 0 {
            e.cursor.Row = 0
        }
        e.ensureCursorValid()

    case tcell.KeyEnter:
        e.buffer.InsertNewline(e.cursor.Row, e.cursor.Col)
        e.cursor.Row++
        if e.cursor.Row < 0 {
            e.cursor.Row = 0
        }
        e.cursor.Col = 0
        e.ensureCursorValid()
        e.buffer.SaveState(e.cursor.Row, e.cursor.Col)
        e.quitAttempts = 0

    case tcell.KeyBackspace, tcell.KeyBackspace2:
        if e.cursor.Col > 0 {
            e.buffer.DeleteChar(e.cursor.Row, e.cursor.Col)
            e.cursor.Col--
        } else if e.cursor.Row > 0 {
            prevLineLen := len(e.buffer.GetLine(e.cursor.Row - 1))
            e.buffer.DeleteChar(e.cursor.Row, e.cursor.Col)
            e.cursor.Row--
            e.cursor.Col = prevLineLen
        }
        e.ensureCursorValid()
        e.buffer.SaveState(e.cursor.Row, e.cursor.Col)
        e.quitAttempts = 0

    case tcell.KeyDelete:
        lineLen := len(e.buffer.GetLine(e.cursor.Row))
        if e.cursor.Col < lineLen {
            e.buffer.DeleteCharForward(e.cursor.Row, e.cursor.Col)
        } else if e.cursor.Row < e.buffer.LineCount()-1 {
            nextLine := e.buffer.GetLine(e.cursor.Row + 1)
            e.buffer.AppendToLine(e.cursor.Row, nextLine)
            e.buffer.DeleteLine(e.cursor.Row + 1)
        }
        e.ensureCursorValid()
        e.buffer.SaveState(e.cursor.Row, e.cursor.Col)
        e.quitAttempts = 0

    case tcell.KeyTab:
        for i := 0; i < 4; i++ {
            e.buffer.InsertChar(e.cursor.Row, e.cursor.Col, ' ')
            e.cursor.Col++
        }
        e.ensureCursorValid()
        e.buffer.SaveState(e.cursor.Row, e.cursor.Col)
        e.quitAttempts = 0

    case tcell.KeyRune:
        e.buffer.InsertChar(e.cursor.Row, e.cursor.Col, ev.Rune())
        e.cursor.Col++
        e.ensureCursorValid()
        e.buffer.SaveState(e.cursor.Row, e.cursor.Col)
        e.quitAttempts = 0
    }

    return true
}

func (e *Editor) handleFindMode(ev *tcell.EventKey) bool {
    switch ev.Key() {
    case tcell.KeyEscape:
        e.mode = ModeNormal
        e.statusMsg = "Search cancelled"
    case tcell.KeyEnter:
        e.findQuery = e.inputBuffer
        e.findText()
        e.mode = ModeNormal
    case tcell.KeyBackspace, tcell.KeyBackspace2:
        if len(e.inputBuffer) > 0 {
            e.inputBuffer = e.inputBuffer[:len(e.inputBuffer)-1]
        }
        e.statusMsg = "Find: " + e.inputBuffer
    case tcell.KeyRune:
        e.inputBuffer += string(ev.Rune())
        e.statusMsg = "Find: " + e.inputBuffer
    }
    return true
}

func (e *Editor) handleGotoMode(ev *tcell.EventKey) bool {
    switch ev.Key() {
    case tcell.KeyEscape:
        e.mode = ModeNormal
        e.statusMsg = "Go to line cancelled"
    case tcell.KeyEnter:
        var lineNum int
        _, err := fmt.Sscanf(e.inputBuffer, "%d", &lineNum)
        if err == nil && lineNum > 0 && lineNum <= e.buffer.LineCount() {
            e.cursor.Row = lineNum - 1
            e.cursor.Col = 0
            e.ensureCursorValid()
            e.statusMsg = fmt.Sprintf("Jumped to line %d", lineNum)
        } else {
            e.statusMsg = "Invalid line number"
        }
        e.mode = ModeNormal
    case tcell.KeyBackspace, tcell.KeyBackspace2:
        if len(e.inputBuffer) > 0 {
            e.inputBuffer = e.inputBuffer[:len(e.inputBuffer)-1]
        }
        e.statusMsg = "Go to line: " + e.inputBuffer
    case tcell.KeyRune:
        if ev.Rune() >= '0' && ev.Rune() <= '9' {
            e.inputBuffer += string(ev.Rune())
            e.statusMsg = "Go to line: " + e.inputBuffer
        }
    }
    return true
}

func (e *Editor) handleLLMMode(ev *tcell.EventKey) bool {
    switch ev.Key() {
    case tcell.KeyEscape:
        e.mode = ModeNormal
        e.statusMsg = "AI prompt cancelled"
    case tcell.KeyEnter:
        e.llmPrompt = e.inputBuffer
        e.askLLMAsync()
        e.mode = ModeNormal
    case tcell.KeyBackspace, tcell.KeyBackspace2:
        if len(e.inputBuffer) > 0 {
            e.inputBuffer = e.inputBuffer[:len(e.inputBuffer)-1]
        }
        e.statusMsg = "Ask AI: " + e.inputBuffer
    case tcell.KeyRune:
        e.inputBuffer += string(ev.Rune())
        e.statusMsg = "Ask AI: " + e.inputBuffer
    }
    return true
}

func (e *Editor) handleFilenameMode(ev *tcell.EventKey) bool {
    switch ev.Key() {
    case tcell.KeyEscape:
        e.mode = ModeNormal
        e.statusMsg = "Save cancelled"
    case tcell.KeyEnter:
        e.buffer.filename = e.inputBuffer
        e.saveFile()
        e.mode = ModeNormal
    case tcell.KeyBackspace, tcell.KeyBackspace2:
        if len(e.inputBuffer) > 0 {
            e.inputBuffer = e.inputBuffer[:len(e.inputBuffer)-1]
        }
        e.statusMsg = "Filename: " + e.inputBuffer
    case tcell.KeyRune:
        e.inputBuffer += string(ev.Rune())
        e.statusMsg = "Filename: " + e.inputBuffer
    }
    return true
}

func (e *Editor) saveFile() {
    if e.buffer.filename == "" {
        e.mode = ModeFilename
        e.inputBuffer = ""
        e.statusMsg = "Enter filename: "
        return
    }

    if err := e.buffer.Save(); err != nil {
        e.statusMsg = fmt.Sprintf("Save failed: %v", err)
    } else {
        basename := filepath.Base(e.buffer.filename)
        e.statusMsg = fmt.Sprintf("Saved '%s' (%d lines)", basename, e.buffer.LineCount())
    }
}

func (e *Editor) findText() {
    if e.findQuery == "" {
        e.statusMsg = "No search query entered"
        return
    }

    totalLines := e.buffer.LineCount()
    if totalLines == 0 {
        e.statusMsg = "Buffer is empty"
        return
    }

    startRow := e.cursor.Row
    startCol := e.cursor.Col + 1

    if startRow < 0 {
        startRow = 0
    }
    if startRow >= totalLines {
        startRow = totalLines - 1
    }

    for i := 0; i < totalLines; i++ {
        row := (startRow + i) % totalLines
        line := e.buffer.GetLine(row)

        searchFrom := 0
        if row == startRow && i == 0 {
            searchFrom = startCol
        }

        if searchFrom >= len(line) {
            continue
        }

        lowerLine := strings.ToLower(line[searchFrom:])
        lowerQuery := strings.ToLower(e.findQuery)
        idx := strings.Index(lowerLine, lowerQuery)

        if idx != -1 {
            e.cursor.Row = row
            e.cursor.Col = searchFrom + idx
            e.ensureCursorValid()
            e.statusMsg = fmt.Sprintf("Found '%s' at line %d, column %d", e.findQuery, row+1, e.cursor.Col+1)
            return
        }
    }

    e.statusMsg = fmt.Sprintf("'%s' not found in document", e.findQuery)
}

func (e *Editor) askLLMAsync() {
    if e.llmPrompt == "" {
        e.statusMsg = "No prompt entered"
        return
    }

    e.aiMutex.Lock()
    if e.aiInProgress {
        e.aiMutex.Unlock()
        e.statusMsg = "AI request already in progress. Press Esc to cancel current request"
        return
    }
    e.aiInProgress = true
    e.aiMutex.Unlock()

    if !e.llmClient.IsAvailable() {
        e.aiMutex.Lock()
        e.aiInProgress = false
        e.statusMsg = "Cannot connect to Ollama. Is it running? Try: ollama serve"
        e.aiMutex.Unlock()
        return
    }

    if err := e.llmClient.CheckModel(); err != nil {
        e.aiMutex.Lock()
        e.aiInProgress = false
        errMsg := err.Error()
        if len(errMsg) > 80 {
            errMsg = errMsg[:77] + "..."
        }
        e.statusMsg = fmt.Sprintf("Model error: %s", errMsg)
        e.aiMutex.Unlock()
        return
    }

    e.statusMsg = "Processing AI request... (Press Esc to cancel)"

    go func() {
        prompt := e.llmPrompt
        
        done := make(chan bool, 1)
        var response string
        var err error
        
        go func() {
            response, err = e.llmClient.GenerateWithCancel(prompt, e.aiCancel)
            done <- true
        }()
        
        select {
        case <-done:
        case <-time.After(90 * time.Second):
            select {
            case e.aiCancel <- true:
            default:
            }
            err = fmt.Errorf("request timeout (90s)")
        }

        e.aiMutex.Lock()
        defer e.aiMutex.Unlock()

        e.aiInProgress = false

        if err != nil {
            errMsg := err.Error()
            if errMsg == "cancelled" {
                e.statusMsg = "AI request cancelled by user"
            } else if strings.Contains(errMsg, "cannot connect") {
                e.statusMsg = "Cannot connect to Ollama. Run: ollama serve"
            } else if strings.Contains(errMsg, "model") && strings.Contains(errMsg, "not found") {
                modelName := e.llmClient.model
                e.statusMsg = fmt.Sprintf("Model '%s' not found. Run: ollama pull %s", modelName, modelName)
            } else if strings.Contains(errMsg, "timeout") {
                e.statusMsg = "AI request timeout. Try a simpler prompt or check Ollama"
            } else {
                if len(errMsg) > 70 {
                    errMsg = errMsg[:67] + "..."
                }
                e.statusMsg = fmt.Sprintf("AI error: %s", errMsg)
            }
            e.llmResponse = ""
            return
        }

        if response == "" {
            e.statusMsg = "AI returned empty response. Try rephrasing your prompt"
            e.llmResponse = ""
            return
        }

        e.llmResponse = response
        
        preview := response
        preview = strings.ReplaceAll(preview, "\n", " ")
        preview = strings.ReplaceAll(preview, "\r", "")
        preview = strings.TrimSpace(preview)
        
        if len(preview) > 60 {
            preview = preview[:57] + "..."
        }
        
        if preview == "" {
            preview = "[Response ready]"
        }
        
        responseLines := strings.Count(response, "\n") + 1
        e.statusMsg = fmt.Sprintf("AI response ready (%d lines). Preview: %s | Press Ctrl+K to insert", responseLines, preview)
    }()
}

func (e *Editor) ensureCursorValid() {
    if e.cursor.Row < 0 {
        e.cursor.Row = 0
    }
    maxRow := e.buffer.LineCount() - 1
    if maxRow < 0 {
        maxRow = 0
    }
    if e.cursor.Row > maxRow {
        e.cursor.Row = maxRow
    }

    lineLen := len(e.buffer.GetLine(e.cursor.Row))
    if e.cursor.Col > lineLen {
        e.cursor.Col = lineLen
    }
    if e.cursor.Col < 0 {
        e.cursor.Col = 0
    }
}

func (e *Editor) render() {
    if e.screen == nil {
        return
    }

    e.screen.Clear()

    e.ensureCursorValid()

    if e.cursor.Row < e.offsetRow {
        e.offsetRow = e.cursor.Row
    }
    if e.cursor.Row >= e.offsetRow+e.height && e.height > 0 {
        e.offsetRow = e.cursor.Row - e.height + 1
    }
    if e.offsetRow < 0 {
        e.offsetRow = 0
    }

    for y := 0; y < e.height; y++ {
        row := y + e.offsetRow
        if row >= e.buffer.LineCount() {
            e.drawString(0, y, "~", tcell.StyleDefault.Foreground(tcell.ColorBlue))
            continue
        }

        line := e.buffer.GetLine(row)
        e.drawString(0, y, line, tcell.StyleDefault)
    }

    e.renderStatusBar()

    screenY := e.cursor.Row - e.offsetRow
    screenX := e.cursor.Col

    if screenX >= e.width {
        screenX = e.width - 1
    }
    if screenX < 0 {
        screenX = 0
    }
    if screenY >= e.height {
        screenY = e.height - 1
    }
    if screenY < 0 {
        screenY = 0
    }

    e.screen.ShowCursor(screenX, screenY)
    e.screen.Show()
}

func (e *Editor) renderStatusBar() {
    y := e.height
    if y < 0 {
        return
    }

    style := tcell.StyleDefault.
        Background(tcell.ColorWhite).
        Foreground(tcell.ColorBlack)

    e.aiMutex.Lock()
    inProgress := e.aiInProgress
    e.aiMutex.Unlock()

    if inProgress {
        style = tcell.StyleDefault.
            Background(tcell.ColorYellow).
            Foreground(tcell.ColorBlack)
    }

    for x := 0; x < e.width; x++ {
        if y >= 0 && y < e.height+2 {
            e.screen.SetContent(x, y, ' ', nil, style)
        }
        if y+1 >= 0 && y+1 < e.height+2 {
            e.screen.SetContent(x, y+1, ' ', nil, style)
        }
    }

    statusMsg := e.statusMsg
    if len(statusMsg) > e.width && e.width > 3 {
        statusMsg = statusMsg[:e.width-3] + "..."
    }
    e.drawString(0, y, statusMsg, style)

    modMark := ""
    if e.buffer.modified {
        modMark = " [+]"
    }
    filename := e.buffer.filename
    if filename == "" {
        filename = "[No Name]"
    } else {
        filename = filepath.Base(filename)
        if len(filename) > 20 {
            filename = filename[:17] + "..."
        }
    }

    info := fmt.Sprintf("%s%s | Ln %d/%d | Col %d",
        filename, modMark, e.cursor.Row+1, e.buffer.LineCount(), e.cursor.Col+1)

    if len(info) > e.width && e.width > 0 {
        info = info[:e.width]
    }

    if y+1 >= 0 && y+1 < e.height+2 {
        e.drawString(0, y+1, info, style)
    }
}

func (e *Editor) drawString(x, y int, s string, style tcell.Style) {
    if e.screen == nil || y < 0 || y >= e.height+2 || x < 0 {
        return
    }

    for i, r := range s {
        posX := x + i
        if posX >= e.width || posX < 0 {
            break
        }
        e.screen.SetContent(posX, y, r, nil, style)
    }
}

func main() {
    ollamaURL := flag.String("ollama", "http://localhost:11434", "Ollama API URL")
    model := flag.String("model", "llama2", "LLM model to use")
    showVersion := flag.Bool("version", false, "Show version")
    showHelp := flag.Bool("help", false, "Show help")

    flag.Parse()

    if *showVersion {
        fmt.Printf("GoEdit v%s\n", version)
        os.Exit(0)
    }

    if *showHelp {
        printHelp()
        os.Exit(0)
    }

    var filename string
    if flag.NArg() > 0 {
        filename = flag.Arg(0)
        absPath, err := filepath.Abs(filename)
        if err == nil {
            filename = absPath
        }
    }

    ed, err := NewEditor(filename, *ollamaURL, *model)
    if err != nil {
        log.Fatalf("Failed to create editor: %v", err)
    }

    if err := ed.Run(); err != nil {
        log.Fatalf("Editor error: %v", err)
    }
}

func printHelp() {
    fmt.Println("GoEdit - A minimal yet powerful terminal text editor")
    fmt.Println("\nUsage:")
    fmt.Println("  goedit [options] [filename]")
    fmt.Println("\nOptions:")
    fmt.Println("  -ollama string    Ollama API URL (default: http://localhost:11434)")
    fmt.Println("  -model string     LLM model to use (default: llama2)")
    fmt.Println("  -version          Show version")
    fmt.Println("  -help             Show this help")
    fmt.Println("\nKeyboard Shortcuts:")
    fmt.Println("  Ctrl+S           Save file")
    fmt.Println("  Ctrl+Q           Quit")
    fmt.Println("  Ctrl+F           Find text")
    fmt.Println("  Ctrl+G           Go to line")
    fmt.Println("  Ctrl+A           Select all")
    fmt.Println("  Ctrl+C           Copy line")
    fmt.Println("  Ctrl+X           Cut line")
    fmt.Println("  Ctrl+V           Paste")
    fmt.Println("  Ctrl+Z           Undo")
    fmt.Println("  Ctrl+Y           Redo")
    fmt.Println("  Ctrl+L           Ask LLM (AI Assistant)")
    fmt.Println("  Ctrl+K           Insert LLM response at cursor")
    fmt.Println("  Esc              Cancel AI request")
    fmt.Println("  Tab              Insert 4 spaces")
    fmt.Println("  Home/End         Line start/end")
    fmt.Println("  Page Up/Down     Scroll page")
}
