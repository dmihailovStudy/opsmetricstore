package models

import "time"

type Gauge struct {
	Timestamp time.Time
	Name      string
	Value     float64
}

type Gauges []Gauge
