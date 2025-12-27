// GoEdit v2.0
// Copyright © Prof. Dr. Michael Stal, 2025
// All rights reserved.

package main

import (
    "context"
    "os/exec"
    "runtime"
    "strings"
    "time"
)

type ClipboardManager struct {
    fallback string
}

func NewClipboardManager() *ClipboardManager {
    return &ClipboardManager{
        fallback: "",
    }
}

func (cm *ClipboardManager) Copy(text string) error {
    if text == "" {
        cm.fallback = ""
        return nil
    }

    var cmdPath string
    var cmdArgs []string
    
    switch runtime.GOOS {
    case "darwin":
        cmdPath = "pbcopy"
        cmdArgs = []string{}
    case "linux":
        if _, err := exec.LookPath("xclip"); err == nil {
            cmdPath = "xclip"
            cmdArgs = []string{"-selection", "clipboard"}
        } else if _, err := exec.LookPath("xsel"); err == nil {
            cmdPath = "xsel"
            cmdArgs = []string{"--clipboard", "--input"}
        } else {
            cm.fallback = text
            return nil
        }
    case "windows":
        cmdPath = "powershell.exe"
        cmdArgs = []string{"-NoProfile", "-NonInteractive", "-Command", "$input | Set-Clipboard"}
    default:
        cm.fallback = text
        return nil
    }
    
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()
    
    cmd := exec.CommandContext(ctx, cmdPath, cmdArgs...)
    
    stdin, err := cmd.StdinPipe()
    if err != nil {
        cm.fallback = text
        return nil
    }
    
    if err := cmd.Start(); err != nil {
        cm.fallback = text
        return nil
    }
    
    _, writeErr := stdin.Write([]byte(text))
    closeErr := stdin.Close()
    
    if writeErr != nil || closeErr != nil {
        cm.fallback = text
        cmd.Process.Kill()
        return nil
    }
    
    waitErr := cmd.Wait()
    if waitErr != nil && ctx.Err() != context.DeadlineExceeded {
        cm.fallback = text
        return nil
    }
    
    cm.fallback = text
    return nil
}

func (cm *ClipboardManager) Paste() (string, error) {
    var cmdPath string
    var cmdArgs []string
    
    switch runtime.GOOS {
    case "darwin":
        cmdPath = "pbpaste"
        cmdArgs = []string{}
    case "linux":
        if _, err := exec.LookPath("xclip"); err == nil {
            cmdPath = "xclip"
            cmdArgs = []string{"-selection", "clipboard", "-o"}
        } else if _, err := exec.LookPath("xsel"); err == nil {
            cmdPath = "xsel"
            cmdArgs = []string{"--clipboard", "--output"}
        } else {
            return cm.fallback, nil
        }
    case "windows":
        cmdPath = "powershell.exe"
        cmdArgs = []string{"-NoProfile", "-NonInteractive", "-Command", "Get-Clipboard"}
    default:
        return cm.fallback, nil
    }
    
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()
    
    cmd := exec.CommandContext(ctx, cmdPath, cmdArgs...)
    
    output, err := cmd.Output()
    if err == nil && len(output) > 0 {
        result := string(output)
        result = strings.TrimRight(result, "\r\n")
        if result != "" {
            cm.fallback = result
            return result, nil
        }
    }
    
    return cm.fallback, nil
}
