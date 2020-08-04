package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/golang/glog"
	// config "github.com/prebid/prebid-server/stored_requests/config"
	"github.com/prebid/prebid-server/pbsmetrics"
	"github.com/prebid/prebid-server/stored_requests/events"
)

type StoredRequestType string

const (
	Amp      StoredRequestType = "amp"
	Auction  StoredRequestType = "auction"
	Category StoredRequestType = "category"
	Video    StoredRequestType = "video"
)

// This function queries the database to get all the data, and is guaranteed to return
// an EventProducer with a single "events.Save" object already in the channel before returning.
//
// The string query should return Rows with the following columns and types:
//
//   1. id: string
//   2. data: JSON
//   3. type: string ("request" or "imp")
//
func LoadAll(srType StoredRequestType, ctx context.Context, db *sql.DB, metricsEngine pbsmetrics.MetricsEngine, query string) (eventProducer *PostgresLoader, err error) {
	if db == nil {
		glog.Fatal("The Stored Request Postgres Startup needs a database connection to work.")
	}
	eventProducer = &PostgresLoader{
		metricsEngine: metricsEngine,
		srType:        srType,
		saves:         make(chan events.Save, 1),
	}
	err = eventProducer.doFetch(ctx, db, query)
	return
}

type PostgresLoader struct {
	metricsEngine pbsmetrics.MetricsEngine
	srType        StoredRequestType
	saves         chan events.Save
}

func (loader *PostgresLoader) doFetch(ctx context.Context, db *sql.DB, query string) error {
	startTime := time.Now()
	glog.Infof("Loading all %s Stored Requests from Postgres with: %s", loader.srType, query)
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		glog.Warningf("Failed to fetch %s Stored Requests from Postgres on startup. The app might be a bit slow to start. Error was: %v", loader.srType, err)
		loader.saves <- events.Save{}
		elapsedTime := time.Since(startTime)
		loader.metricsEngine.RecordStoredRequestLoadTime(
			pbsmetrics.StoredRequestLabels{
				StoredRequestType: "Auction", //TODO(bfs): !!!!!!
				PullType:          pbsmetrics.PeriodicTaskInitial,
			}, elapsedTime)
		return err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			glog.Warningf("Failed to close DB connection after fetching %s Stored Requests from Postgres: %v", loader.srType, err)
		}
		elapsedTime := time.Since(startTime)
		loader.metricsEngine.RecordStoredRequestLoadTime(
			pbsmetrics.StoredRequestLabels{
				StoredRequestType: "Auction", //TODO(bfs): !!!!!!
				PullType:          pbsmetrics.PeriodicTaskInitial,
			}, elapsedTime)
	}()

	if err := sendEvents(rows, loader.saves, nil); err != nil {
		glog.Warningf("Failed to fetch %s Stored Requests from Postgres on startup. Things might be a bit slow to start: %v", loader.srType, err)
		loader.saves <- events.Save{}
		return err
	}

	return nil
}

func (e *PostgresLoader) Saves() <-chan events.Save {
	return e.saves
}

func (e *PostgresLoader) Invalidations() <-chan events.Invalidation {
	return nil
}
