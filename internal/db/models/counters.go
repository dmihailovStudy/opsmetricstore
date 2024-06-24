package models

import "time"

type Counter struct {
	Timestamp time.Time
	Name      string
	Value     int64
}

type Counters []Counter
