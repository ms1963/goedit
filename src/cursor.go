package main

// Cursor represents the cursor position
type Cursor struct {
        Row int
        Col int
}

// NewCursor creates a new cursor
func NewCursor() *Cursor {
        return &Cursor{Row: 0, Col: 0}
}
