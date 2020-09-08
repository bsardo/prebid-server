package postgres

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/golang/glog"
	"github.com/prebid/prebid-server/config"
	"github.com/prebid/prebid-server/pbsmetrics"
	"github.com/prebid/prebid-server/stored_requests/events"
	"github.com/prebid/prebid-server/util/timeutil"
)

type PostgresEventProducerConfig struct {
	Db                 *sql.DB
	RequestType        config.DataType
	CacheInitQuery     string
	CacheInitTimeout   time.Duration
	CacheUpdateQuery   string
	CacheUpdateTimeout time.Duration
	MetricsEngine      pbsmetrics.MetricsEngine
}

type PostgresEventProducer struct {
	cfg           PostgresEventProducerConfig
	lastUpdate    time.Time
	invalidations chan events.Invalidation
	saves         chan events.Save
	time          timeutil.Time
}

func NewPostgresEventProducer(cfg PostgresEventProducerConfig) (eventProducer *PostgresEventProducer) {
	if cfg.Db == nil {
		glog.Fatalf("The Postgres Stored %s Loader needs a database connection to work.", cfg.RequestType)
	}

	return &PostgresEventProducer{
		cfg:           cfg,
		lastUpdate:    time.Time{},
		saves:         make(chan events.Save, 1),
		invalidations: make(chan events.Invalidation, 1),
		time:          &timeutil.RealTime{},
	}
}

func (e *PostgresEventProducer) Run() error {
	if e.lastUpdate.IsZero() {
		return e.fetchAll()
	} else {
		return e.fetchDelta()
	}
}

func (e *PostgresEventProducer) Saves() <-chan events.Save {
	return e.saves
}

func (e *PostgresEventProducer) Invalidations() <-chan events.Invalidation {
	return e.invalidations
}

func (e *PostgresEventProducer) fetchAll() error {
	timeout := e.cfg.CacheInitTimeout * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	startTime := e.time.Now().UTC()
	rows, err := e.cfg.Db.QueryContext(ctx, e.cfg.CacheInitQuery)
	elapsedTime := time.Since(startTime)
	e.recordFetchTime(elapsedTime, pbsmetrics.FetchAll)

	if err != nil {
		glog.Warningf("Failed to fetch all Stored %s data from the DB: %v", e.cfg.RequestType, err)
		// e.saves <- events.Save{} //TODO: do we need to do this?
		return err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			glog.Warningf("Failed to close the Stored %s DB connection: %v", e.cfg.RequestType, err)
		}
	}()
	if err := e.sendEvents(rows); err != nil {
		glog.Warningf("Failed to load all Stored %s data from the DB: %v", e.cfg.RequestType, err)
		// e.saves <- events.Save{} //TODO: do we need to do this?
		return err
	} else {
		e.lastUpdate = startTime
	}
	return nil
}

func (e *PostgresEventProducer) fetchDelta() error {
	timeout := e.cfg.CacheUpdateTimeout * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	startTime := e.time.Now().UTC()
	rows, err := e.cfg.Db.QueryContext(ctx, e.cfg.CacheUpdateQuery, e.lastUpdate)
	elapsedTime := time.Since(startTime)
	e.recordFetchTime(elapsedTime, pbsmetrics.FetchDelta) // record fetch and load time or just fetch time? prob just fetch time since experienced db timeouts

	if err != nil {
		glog.Warningf("Failed to fetch updated Stored %s data from the DB: %v", e.cfg.RequestType, err)
		return err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			glog.Warningf("Failed to close the Stored %s DB connection: %v", e.cfg.RequestType, err)
		}
	}()
	if err := e.sendEvents(rows); err != nil {
		glog.Warningf("Failed to load updated Stored %s data from the DB: %v", e.cfg.RequestType, err)
		return err
	} else {
		e.lastUpdate = startTime
	}
	return nil
}

func (e *PostgresEventProducer) recordFetchTime(elapsedTime time.Duration, fetchType pbsmetrics.StoredDataFetchType) {
	e.cfg.MetricsEngine.RecordStoredDataFetchTime(
		pbsmetrics.StoredDataTypeLabels{
			DataType:      e.metricsStoredDataType(),
			DataFetchType: fetchType,
		}, elapsedTime)
}

func (e *PostgresEventProducer) metricsStoredDataType() pbsmetrics.StoredDataType {
	return map[config.DataType]pbsmetrics.StoredDataType{
		config.RequestDataType:    pbsmetrics.RequestDataType,
		config.CategoryDataType:   pbsmetrics.CategoryDataType,
		config.VideoDataType:      pbsmetrics.VideoDataType,
		config.AMPRequestDataType: pbsmetrics.AMPRequestDataType,
		config.AccountDataType:    pbsmetrics.AccountDataType,
	}[e.cfg.RequestType]
}

// sendEvents reads the rows and sends notifications into the channel for any updates.
// If it returns an error, then callers can be certain that no events were sent to the channels.
func (e *PostgresEventProducer) sendEvents(rows *sql.Rows) (err error) {
	storedRequestData := make(map[string]json.RawMessage)
	storedImpData := make(map[string]json.RawMessage)

	var requestInvalidations []string
	var impInvalidations []string

	for rows.Next() {
		var id string
		var data []byte
		var dataType string

		// Beware #338... we really don't want to save corrupt data
		if err := rows.Scan(&id, &data, &dataType); err != nil {
			return err
		}

		switch dataType {
		case "request":
			if len(data) == 0 || bytes.Equal(data, []byte("null")) {
				requestInvalidations = append(requestInvalidations, id)
			} else {
				storedRequestData[id] = data
			}
		case "imp":
			if len(data) == 0 || bytes.Equal(data, []byte("null")) {
				impInvalidations = append(impInvalidations, id)
			} else {
				storedImpData[id] = data
			}
		default:
			glog.Warningf("Stored Data with id=%s has invalid type: %s. This will be ignored.", id, dataType)
		}
	}

	// Beware #338... we really don't want to save corrupt data
	if rows.Err() != nil {
		return rows.Err()
	}

	if len(storedRequestData) > 0 || len(storedImpData) > 0 {
		e.saves <- events.Save{
			Requests: storedRequestData,
			Imps:     storedImpData,
		}
	}

	//TODO: do we really need to check that this isn't a fetch all here?
	if (len(requestInvalidations) > 0 || len(impInvalidations) > 0) && !e.lastUpdate.IsZero() {
		e.invalidations <- events.Invalidation{
			Requests: requestInvalidations,
			Imps:     impInvalidations,
		}
	}

	return
}