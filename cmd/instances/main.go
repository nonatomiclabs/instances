package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/nonatomiclabs/instances"
)

func openDatabase(filePath string) (*instances.Database, error) {
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		// The database doesn't exist yet, create an empty one
		emptyDatabase := instances.Database{
			Instances: map[string]instances.Instance{},
		}
		f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return &instances.Database{}, fmt.Errorf("create default database file: %s", err)
		}
		err = json.NewEncoder(f).Encode(emptyDatabase)
		if err != nil {
			return &instances.Database{}, fmt.Errorf("write default databse: %s", err)
		}
	}

	f, err := os.Open(filePath)
	if err != nil {
		return &instances.Database{}, fmt.Errorf("open database file: %s", err)
	}
	defer f.Close()

	db, err := instances.NewDatabase(f)
	if err != nil {
		return &instances.Database{}, err
	}

	return db, nil
}

func main() {
	userDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dbPath := filepath.Join(userDir, ".instances.db.json")
	db, err := openDatabase(dbPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}
	ec2Client := ec2.NewFromConfig(cfg)

	cloudProviders := map[string]instances.CloudProvider{
		"aws": instances.AWSCloud{Ec2Client: ec2Client},
	}

	CLI := instances.NewCLI(db, cloudProviders)

	if err = CLI.Run(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
