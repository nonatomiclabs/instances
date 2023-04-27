package instances

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
)

type Database struct {
	Instances map[string]Instance `json:"instances"`
	support   io.ReadWriter
}

// NewDatabase creates a new Database populated with the content read from the given
// io.ReadWriter.
func NewDatabase(support io.ReadWriter) (*Database, error) {
	database := Database{
		support: support,
	}
	err := json.NewDecoder(support).Decode(&database)
	if err != nil {
		return &database, fmt.Errorf("open database: %s", err)
	}

	return &database, nil
}

// Save saves the database to the provided io.Writer.
func (d *Database) Save() error {
	b, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return fmt.Errorf("serialize database: %s", err)
	}

	_, err = d.support.Write(b)
	if err != nil {
		return fmt.Errorf("save database: %s", err)
	}
	return nil
}

// AddInstance adds an instance to the database.
func (d *Database) AddInstance(id string, name string, cloudProvider CloudProvider) error {
	log.Printf("adding instance %s", id)

	if _, instanceExists := d.Instances[name]; instanceExists {
		return fmt.Errorf("instance %q exists already", name)
	}

	for instanceName, instance := range d.Instances {
		if instance.Id == id {
			return fmt.Errorf("instance id %q already referenced by instance %q", id, instanceName)
		}
	}

	_, err := cloudProvider.GetInstanceStatus(id)
	if err != nil {
		return err
	}

	d.Instances[name] = Instance{Id: id, CloudProviderName: cloudProvider.GetName()}

	return nil
}

// GetInstance gets an instance from the database
func (d *Database) GetInstance(name string) (Instance, error) {
	instance, instanceExists := d.Instances[name]
	if !instanceExists {
		return Instance{}, fmt.Errorf("no instance named %s", name)
	}
	return instance, nil
}

// RemoveInstance removes an instance from the database
func (d *Database) RemoveInstance(name string) error {
	_, instanceExists := d.Instances[name]
	if !instanceExists {
		return fmt.Errorf("no instance named %s", name)
	}
	delete(d.Instances, name)
	return nil
}
