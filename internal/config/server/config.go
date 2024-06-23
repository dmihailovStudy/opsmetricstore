package server

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog/log"
)

// run flags
const AFlag = "a"
const ADefault = "localhost:8080"
const AUsage = "specify the url"

const IFlag = "i"
const IDefault = 300
const IUsage = "interval to save storage"

const FFlag = "f"
const FDefault = "tmp/metrics-db.json"
const FUsage = "path to save snapshot"

const RFlag = "r"
const RDefault = true
const RUsage = "restore start snapshot?"

// routing paths
const MainPath = "/"
const GetMetricByURLPath = "/value/:metricType/:metricName"
const GetMetricByJSONPath = "/value"
const UpdateByURLPath = "/update/:metricType/:metricName/:metricValue"
const UpdateByJSONPath = "/update"

type Config struct {
	Address       string `env:"ADDRESS"`
	StoreInterval int    `env:"STORE_INTERVAL"`
	Path          string `env:"FILE_STORAGE_PATH"`
	Restore       bool   `env:"RESTORE"`
}

func (c *Config) Load() {
	var endpoint string
	var interval int
	var path string
	var restore bool
	if err := env.Parse(c); err != nil {
		errMsg := "Load(): parse env config"
		log.Warn().Err(err).Msg(errMsg)
	}

	// read flags
	flag.StringVar(&endpoint, AFlag, ADefault, AUsage)
	flag.IntVar(&interval, IFlag, IDefault, IUsage)
	flag.StringVar(&path, FFlag, FDefault, FUsage)
	flag.BoolVar(&restore, RFlag, RDefault, RUsage)
	flag.Parse()

	if c.Address == "" {
		c.Address = endpoint
	}
	if c.StoreInterval == 0 {
		c.StoreInterval = interval
	}
	if c.Path == "" {
		c.Path = path
	}
	if !c.Restore {
		c.Restore = restore
	}
}
