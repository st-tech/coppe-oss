package services

import (
	"context"
	"testing"
	"time"
)

func TestTimeout(t *testing.T) {
	sql := `SELECT SUM(4 * POW(-1, n) / (2 * n + 1)) as pi FROM
		(
			SELECT  GENERATE_ARRAY(1000 * l - 1000, 1000 * l - 1) AS m
			FROM UNNEST(GENERATE_ARRAY(1, 1000000)) AS l
		),
		UNNEST(m) AS n` // compute pi using Google's resource cf. https://qiita.com/shiozaki/items/97648fee7849cc45e8ab

	tctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := queryToBigQuery(tctx, sql, true)
	if err != nil {
		t.Logf("Timeout detected as: %v", err)
		return
	}
	t.Error("ERROR: timeout was not detected")
}
