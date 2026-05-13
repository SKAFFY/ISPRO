package model

import "time"

type Entry struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type EntryInput struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Link struct {
	ID        int64     `json:"id"`
	SourceID  int64     `json:"source_id"`
	TargetID  int64     `json:"target_id"`
	CreatedAt time.Time `json:"created_at"`
}

type LinkInput struct {
	SourceID int64 `json:"source_id"`
	TargetID int64 `json:"target_id"`
}