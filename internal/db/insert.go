package db

import (
	"fmt"
	"github.com/dmihailovStudy/opsmetricstore/internal/db/models"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

func InsertGauges(gauges models.Gauges) error {
	tableName := "gauges"

	if len(gauges) == 0 {
		log.Info().Msg("InsertGauges(): gauges empty")
		return nil
	}

	query := &strings.Builder{}
	query.WriteString(fmt.Sprintf(`INSERT INTO %s 
    (
		timestamp,  
		name,     
		value
     ) VALUES `, tableName))
	for i, gauge := range gauges {
		val := fmt.Sprintf(
			`('%s','%s','%v')`,
			gauge.Timestamp.Format(time.DateTime),
			gauge.Name,
			gauge.Value,
		)
		if i > 0 {
			query.WriteRune(',')
		}
		query.WriteString(val)
	}
	_, err := conn.Exec(query.String())
	if err != nil {
		log.Warn().Err(err).Msg("InsertGauges(): insert error")
		return err
	}
	return nil
}

func InsertCounters(counters models.Counters) error {
	tableName := "counters"

	if len(counters) == 0 {
		log.Info().Msg("InsertCounters(): counters empty")
		return nil
	}

	query := &strings.Builder{}
	query.WriteString(fmt.Sprintf(`INSERT INTO %s 
    (
		timestamp,  
		name,     
		value
     ) VALUES `, tableName))
	for i, counter := range counters {
		val := fmt.Sprintf(
			`('%s','%s','%v')`,
			counter.Timestamp.Format(time.DateTime),
			counter.Name,
			counter.Value,
		)
		if i > 0 {
			query.WriteRune(',')
		}
		query.WriteString(val)
	}
	_, err := conn.Exec(query.String())
	if err != nil {
		log.Warn().Err(err).Msg("InsertCounters(): insert error")
		return err
	}
	return nil
}
