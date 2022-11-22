package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"cloud.google.com/go/spanner"
	database "cloud.google.com/go/spanner/admin/database/apiv1"
	instance "cloud.google.com/go/spanner/admin/instance/apiv1"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/option"
	adminpb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
	instancepb "google.golang.org/genproto/googleapis/spanner/admin/instance/v1"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

var projectID = "uas-project"
var instanceID = "uas-instance"
var dbName = "uas-db"

func main() {
	ctx := context.Background()

	// Point to local spanner
	err := os.Setenv("SPANNER_EMULATOR_HOST", "localhost:9010")
	if err != nil {
		panic(err)
	}

	// Create instance
	err = createInstance(ctx)
	if err != nil {
		panic(err)
	}

	// Create database
	err = createDatabase(ctx)
	if err != nil {
		panic(err)
	}

	// Insert record
	err = insertRecords(ctx)
	if err != nil {
		panic(err)
	}
}

func createDatabase(ctx context.Context) error {
	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return err
	}
	defer adminClient.Close()

	b, _ := os.ReadFile("spanner/schema/spanner-schema.ddl")
	statements := strings.Split(string(b), ";")
	fmt.Println(statements)
	op, err := adminClient.CreateDatabase(ctx, &adminpb.CreateDatabaseRequest{
		Parent:          "projects/uas-project/instances/uas-instance",
		CreateStatement: "CREATE DATABASE `uas-db`",
		ExtraStatements: statements,
	},
	)
	if err != nil {
		return err
	}
	if _, err := op.Wait(ctx); err != nil {
		return err
	}
	fmt.Println("database created")
	return nil
}

func createInstance(ctx context.Context) error {
	instanceAdmin, err := instance.NewInstanceAdminClient(ctx)
	if err != nil {
		return err
	}
	defer instanceAdmin.Close()
	instance, _ := instanceAdmin.GetInstance(ctx, &instancepb.GetInstanceRequest{
		Name:      fmt.Sprintf("projects/%s/instances/%s", projectID, instanceID),
		FieldMask: &fieldmaskpb.FieldMask{},
	})

	if instance != nil {
		return nil
	}

	op, err := instanceAdmin.CreateInstance(ctx, &instancepb.CreateInstanceRequest{
		Parent:     fmt.Sprintf("projects/%s", projectID),
		InstanceId: instanceID,
		Instance: &instancepb.Instance{
			Config:      fmt.Sprintf("projects/%s/instanceConfigs/%s", projectID, "regional-us-central1"),
			DisplayName: instanceID,
			NodeCount:   1,
			Labels:      map[string]string{"cloud_spanner_samples": "true"},
		},
	})
	if err != nil {
		return fmt.Errorf("could not create instance %s: %v", fmt.Sprintf("projects/%s/instances/%s", projectID, instanceID), err)
	}
	// Wait for the instance creation to finish.
	i, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("waiting for instance creation to finish failed: %v", err)
	}
	// The instance may not be ready to serve yet.
	if i.State != instancepb.Instance_READY {
		fmt.Println("instance state is not READY yet")
	}
	fmt.Println("instance created")
	return nil
}

func insertRecords(ctx context.Context) error {
	client, err := NewSpannerClient(
		ctx,
		fmt.Sprintf(
			"projects/%s/instances/%s/databases/%s",
			projectID,
			instanceID,
			dbName,
		),
	)
	if err != nil {
		return err
	}

	for i := 0; i < 100; i++ {

		m := spanner.Insert(
			"Application",
			[]string{
				"PublicApplicationID",
				"Name",
				"IconURL",
			},
			[]interface{}{
				uuid.New().String(),
				"Zhi Qu " + fmt.Sprint(i),
				"https://www.google.com",
			},
		)

		_, err = client.Apply(ctx, []*spanner.Mutation{m})
		if err != nil {
			return err
		}
	}

	fmt.Println("data inserted")

	return nil
}

func NewSpannerClient(
	ctx context.Context,
	db string,
	opts ...option.ClientOption) (*spanner.Client, error) {
	spannerConfig := spanner.ClientConfig{
		SessionPoolConfig: spanner.DefaultSessionPoolConfig,
	}

	client, err := spanner.NewClientWithConfig(ctx, db, spannerConfig, opts...)
	if err != nil {
		log.Err(err).Msg("Cannot create spanner client")
		return nil, err
	}

	return client, nil
}
