package model

import "time"

type Task struct {
	Status   string
	Priority int
	Duration time.Duration
	Project  string
	Tags     []string
	Mentions []string
	Title    string
	Note     string
	URLs     []string
	Line     int
}
