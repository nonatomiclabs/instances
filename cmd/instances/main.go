package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

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
	const dbPath = "/Users/jean/Code/instances/sample_config_2.json"
	db, err := openDatabase(dbPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cloudProviders := map[string]instances.CloudProvider{
		"aws": instances.MockAWSCloud{},
	}

	CLI := instances.NewCLI(db, cloudProviders)

	if err = CLI.Run(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
