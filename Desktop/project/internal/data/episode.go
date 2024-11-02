package data

import "time"

type Episode struct {
	ID         int64     `json:"id"`
	CreatedAt  time.Time `json:"-"`
	Title      string    `json:"title"`
	Year       int32     `json:"year,omitempty"`
	Runtime    Runtime   `json:"runtime,omitempty"`
	Characters []string  `json:"characters,omitempty"`
	Version    int32     `json:"version"`
}
