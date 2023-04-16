package instances

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
)

type Database struct {
	Instances map[string]Instance `json:"instances"`
}

// NewDatabase creates a new Database populated with the content read from the given
// io.Reader.
func NewDatabase(r io.Reader) (*Database, error) {

	var database *Database
	err := json.NewDecoder(r).Decode(&database)
	if err != nil {
		return database, fmt.Errorf("open database: %s", err)
	}

	return database, nil
}

// Save saves the database to the provided io.Writer.
func (d *Database) Save(w io.Writer) error {
	b, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return fmt.Errorf("serialize database: %s", err)
	}
	_, err = w.Write(b)
	if err != nil {
		return fmt.Errorf("save database: %s", err)
	}
	return nil
}

func (d *Database) AddInstance(id string, name string, cloudProvider CloudProvider) error {
	log.Printf("adding instance %s", id)

	if _, instanceExists := d.Instances[name]; instanceExists {
		return fmt.Errorf("instance %s exists already", name)
	}

	_, err := cloudProvider.GetInstanceStatus(id)
	if errors.Is(err, ErrInstanceNotFound) {
		return err
	}

	d.Instances[name] = Instance{Id: id, CloudProviderName: cloudProvider.GetName()}

	return nil
}

func (d *Database) GetInstance(name string) (Instance, error) {
	instance, instanceExists := d.Instances[name]
	if !instanceExists {
		return Instance{}, fmt.Errorf("no instance named %s", name)
	}
	return instance, nil
}

func (d *Database) RemoveInstance(name string) error {
	_, instanceExists := d.Instances[name]
	if !instanceExists {
		return fmt.Errorf("no instance named %s", name)
	}
	delete(d.Instances, name)
	return nil
}
