package repository

import (
	"cloud.google.com/go/bigquery"
)

type BigQueryRepo struct {
	client    *bigquery.Client
	projectID string
}

func NewBigQueryRepo(client *bigquery.Client, projectID string) *BigQueryRepo {
	return &BigQueryRepo{
		client:    client,
		projectID: projectID,
	}
}
