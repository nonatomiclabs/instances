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

func main() {
	userDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dbPath := filepath.Join(userDir, ".instances.db.json")

	if _, err := os.Stat(dbPath); errors.Is(err, os.ErrNotExist) {
		// The database doesn't exist yet, create an empty one
		emptyDatabase := instances.Database{
			Instances: map[string]instances.Instance{},
		}
		f, err := os.OpenFile(dbPath, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			fmt.Printf("create default database file: %s\n", err)
			os.Exit(1)
		}
		err = json.NewEncoder(f).Encode(emptyDatabase)
		if err != nil {
			fmt.Printf("write default databse: %s\n", err)
			os.Exit(1)
		}
	}

	f, err := os.OpenFile(dbPath, os.O_RDWR, 0644)
	if err != nil {
		fmt.Printf("open database file: %s\n", err)
		os.Exit(1)
	}
	defer f.Close()

	db, err := instances.NewDatabase(f)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer func() {
		err = f.Truncate(0)   // TODO: handle error
		_, err = f.Seek(0, 0) // TODO: handle error
		db.Save()             // TODO: handle error
	}()

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
