package main

import (
	"log"
	"os"

	"github.com/CameronXie/endpoint-monitor/internal/cmd"
	"github.com/CameronXie/endpoint-monitor/internal/status"
	"github.com/CameronXie/endpoint-monitor/internal/storage"
	"github.com/CameronXie/endpoint-monitor/pkg/monitor"
	"github.com/CameronXie/endpoint-monitor/pkg/scheduler"
	"github.com/sirupsen/logrus"
)

const (
	dbURL    = "INFLUXDB_URL"
	dbToken  = "INFLUXDB_TOKEN"
	dbOrg    = "INFLUXDB_ORG"
	dbBucket = "INFLUXDB_BUCKET"
)

func main() {
	l := logrus.New()
	err := cmd.Execute(
		status.New(
			scheduler.New(),
			monitor.New(),
			storage.NewInfluxDBStorage(&storage.InfluxdbStorageConfig{
				URL:    os.Getenv(dbURL),
				Token:  os.Getenv(dbToken),
				Org:    os.Getenv(dbOrg),
				Bucket: os.Getenv(dbBucket),
			}),
			l,
		),
		l,
	)

	if err != nil {
		log.Fatal(err)
	}
}
