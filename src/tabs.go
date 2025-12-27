// GoEdit v2.0
// Copyright © Prof. Dr. Michael Stal, 2025
// All rights reserved.

package main

import (
    "path/filepath"
)

type TabManager struct {
    tabs      []*Tab
    activeTab int
}

type Tab struct {
    buffer    *Buffer
    cursor    *Cursor
    offsetRow int
    offsetCol int
}

func NewTabManager() *TabManager {
    return &TabManager{
        tabs:      make([]*Tab, 0),
        activeTab: -1,
    }
}

func (tm *TabManager) AddTab(filename string) error {
    buffer, err := NewBuffer(filename)
    if err != nil {
        return err
    }
    
    tab := &Tab{
        buffer:    buffer,
        cursor:    NewCursor(),
        offsetRow: 0,
        offsetCol: 0,
    }
    
    tm.tabs = append(tm.tabs, tab)
    tm.activeTab = len(tm.tabs) - 1
    
    return nil
}

func (tm *TabManager) GetActiveTab() *Tab {
    if tm.activeTab < 0 || tm.activeTab >= len(tm.tabs) {
        return nil
    }
    return tm.tabs[tm.activeTab]
}

func (tm *TabManager) NextTab() {
    if len(tm.tabs) == 0 {
        return
    }
    tm.activeTab = (tm.activeTab + 1) % len(tm.tabs)
}

func (tm *TabManager) PrevTab() {
    if len(tm.tabs) == 0 {
        return
    }
    tm.activeTab--
    if tm.activeTab < 0 {
        tm.activeTab = len(tm.tabs) - 1
    }
}

func (tm *TabManager) CloseTab() bool {
    if len(tm.tabs) == 0 {
        return false
    }
    
    if len(tm.tabs) == 1 {
        return true
    }
    
    tm.tabs = append(tm.tabs[:tm.activeTab], tm.tabs[tm.activeTab+1:]...)
    
    if tm.activeTab >= len(tm.tabs) {
        tm.activeTab = len(tm.tabs) - 1
    }
    
    return false
}

func (tm *TabManager) GetTabCount() int {
    return len(tm.tabs)
}

func (tm *TabManager) GetTabName(index int) string {
    if index < 0 || index >= len(tm.tabs) {
        return ""
    }
    
    tab := tm.tabs[index]
    if tab == nil || tab.buffer == nil {
        return "[Error]"
    }
    
    filename := tab.buffer.filename
    if filename == "" {
        filename = "[No Name]"
    } else {
        filename = filepath.Base(filename)
    }
    
    if tab.buffer.modified {
        return filename + " [+]"
    }
    return filename
}

func (tm *TabManager) IsActiveTab(index int) bool {
    return index == tm.activeTab
} 

