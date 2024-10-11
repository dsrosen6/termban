package main

import (
	_ "github.com/mattn/go-sqlite3"
)

type task struct {
	id          int
	title       string
	description string
	status
}
