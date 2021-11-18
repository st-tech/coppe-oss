package services

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

/*
	Reference from: https://github.com/GoogleCloudPlatform/golang-samples/blob/bdc987b4624a0939603bb9f0a74eb2b815aa6577/bigquery/snippets/querying/bigquery_query.go
*/

func queryToBigQuery(ctx context.Context, query string, disableCache bool) ([]BigQueryRow, error) {
	projectId := os.Getenv("GCP_PROJECT_ID")

	client, err := bigquery.NewClient(ctx, projectId)
	if err != nil {
		return nil, fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	q := client.Query(query)
	q.DisableQueryCache = disableCache
	q.Location = "US"

	job, err := q.Run(ctx)
	if err != nil {
		return nil, err
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return nil, err
	}
	if err := status.Err(); err != nil {
		return nil, err
	}
	it, err := job.Read(ctx)
	if err != nil {
		return nil, err
	}

	queryRows := 1000 // reasonable number for human to read & not crushing on memory
	if it.TotalRows < 1000 {
		queryRows = int(it.TotalRows)
	}

	rows := make([]BigQueryRow, queryRows)

	for i := 0; i < queryRows; i++ {
		err := it.Next(&rows[i])
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
	}
	return rows, nil
}
