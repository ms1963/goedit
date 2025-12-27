// GoEdit v2.0
// Copyright © Prof. Dr. Michael Stal, 2025
// All rights reserved.

package main

type Cursor struct {
    Row int
    Col int
}

func NewCursor() *Cursor {
    return &Cursor{Row: 0, Col: 0}
}
