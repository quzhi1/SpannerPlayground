package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"

	"cloud.google.com/go/spanner"
	"github.com/quzhi1/spanner-playground/util"
	"google.golang.org/api/iterator"
)

const pageSize = 5

func main() {
	ctx := context.Background()

	// Point to local spanner
	err := os.Setenv("SPANNER_EMULATOR_HOST", "localhost:9010")
	if err != nil {
		panic(err)
	}

	// Create client
	client, err := spanner.NewClient(
		ctx,
		fmt.Sprintf(
			"projects/%s/instances/%s/databases/%s",
			util.ProjectID,
			util.InstanceID,
			util.DbName,
		),
	)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// Query first page
	nextPageToken := ""
	for i := 0; ; i++ {
		sql := PaginationSql(pageSize, nextPageToken)
		fmt.Printf("SQL for first page: %s\n", sql)
		iter := client.Single().Query(ctx, spanner.Statement{SQL: sql})
		defer iter.Stop()
		var timestamp int64
		itemCount := 0
		for {
			row, err := iter.Next()
			if err == iterator.Done {
				break
			} else if err != nil {
				panic(err)
			}
			var publicApplicationId, name string
			if err := row.Columns(&publicApplicationId, &name, &timestamp); err != nil {
				panic(err)
			}
			itemCount++
			fmt.Printf("%s %s %d\n", publicApplicationId, name, timestamp)
		}
		if itemCount < pageSize {
			fmt.Printf("Pagination is done.")
			break
		}
		nextPageToken = GeneratePageToken(timestamp)
		fmt.Printf("Query is done. Next page token: %s\n", nextPageToken)
	}

}

func PaginationSql(pageSize int, pageToken string) string {
	if pageToken == "" {
		return fmt.Sprintf(
			"SELECT PublicApplicationID, Name, Time FROM Application ORDER BY Time DESC LIMIT %d",
			pageSize,
		)
	} else {
		timestamp := DecodePageToken(pageToken)
		return fmt.Sprintf(
			"SELECT PublicApplicationID, Name, Time FROM Application WHERE Time < %d ORDER BY Time DESC LIMIT %d",
			timestamp,
			pageSize,
		)
	}
}

func GeneratePageToken(timestamp int64) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("time_%d", timestamp)))
}

func DecodePageToken(pageToken string) int64 {
	decoded, err := base64.StdEncoding.DecodeString(pageToken)
	if err != nil {
		panic(err)
	}

	timestampStr := strings.Replace(string(decoded), "time_", "", 1)
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		panic(err)
	}
	return timestamp
}
