package instances

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

type CLI struct {
	db             *Database
	cloudProviders map[string]CloudProvider
}

func NewCLI(db *Database, cloudProviders map[string]CloudProvider) *CLI {
	return &CLI{db, cloudProviders}
}

func (c *CLI) Run(args []string) error {
	if len(args) == 0 {
		return errors.New("use subcommand")
	}

	switch args[0] {
	case "add":
		return c.addInstance(args[1:])
	case "rm":
		return c.removeInstance(args[1:])
	case "status":
		return c.getInstanceStatus(args[1:])
	case "start":
		return c.startInstance(args[1:])
	case "stop":
		return c.stopInstance(args[1:])
	case "list":
		return c.listInstances(args[1:])
	default:
		return errors.New("unknown subcommand")
	}

}

func (c *CLI) addInstance(args []string) error {
	var cloudName, instanceName string
	addCmd := flag.NewFlagSet("add", flag.ContinueOnError)
	addCmd.Usage = func() {
		fmt.Print(
			"Usage: instances add [OPTIONS] INSTANCE_ID\n\n",
			"Add the instance INSTANCE_ID to the tracked instances\n\n",
		)
		addCmd.PrintDefaults()
	}
	addCmd.StringVar(&cloudName, "cloud", "", "the cloud provider (one of AWS, Azure, GCP)")
	addCmd.StringVar(&instanceName, "name", "", "the name under which to store the instance (by default, the instance name in the cloud provider)")

	err := addCmd.Parse(args)
	if err != nil {
		return err
	}

	if addCmd.NArg() == 0 {
		addCmd.Usage()
		return errors.New("missing instance ID")
	}

	if addCmd.NArg() > 1 {
		addCmd.Usage()
		return errors.New("only one instance ID can be provided")
	}

	instanceId := addCmd.Arg(0)

	cloudProvider, exists := c.cloudProviders[strings.ToLower(cloudName)]
	if !exists {
		return fmt.Errorf("unsupported cloud provider %q", cloudName)
	}

	err = c.db.AddInstance(instanceId, instanceName, cloudProvider)
	if err != nil {
		return err
	}

	return nil
}

func (c *CLI) removeInstance(args []string) error {
	removeCmd := flag.NewFlagSet("rm", flag.ContinueOnError)
	removeCmd.Usage = func() {
		fmt.Print(
			"Usage: instances rm INSTANCE_NAME\n\n",
			"Remove the instance INSTANCE_NAME from the list of tracked instances\n\n",
		)
		removeCmd.PrintDefaults()
	}

	name, err := parseInstanceName(removeCmd, args)
	if err != nil {
		return err
	}

	err = c.db.RemoveInstance(name)
	if err != nil {
		return err
	}

	return nil
}

func (c *CLI) getInstanceStatus(args []string) error {
	statusCmd := flag.NewFlagSet("status", flag.ContinueOnError)
	statusCmd.Usage = func() {
		fmt.Print(
			"Usage: instances status INSTANCE_NAME\n\n",
			"Print the status of the instance INSTANCE_NAME\n\n",
		)
		statusCmd.PrintDefaults()
	}

	name, err := parseInstanceName(statusCmd, args)
	if err != nil {
		return err
	}

	instance, err := c.db.GetInstance(name)
	if err != nil {
		return err
	}

	cloudProvider, err := instance.GetCloudProvider(c.cloudProviders)
	if err != nil {
		return fmt.Errorf("could not get instance status: %v", err)
	}

	status, err := cloudProvider.GetInstanceStatus(instance.Id)
	if err != nil {
		return err
	}
	fmt.Println(status)

	return nil
}

func (c *CLI) startInstance(args []string) error {
	startCmd := flag.NewFlagSet("start", flag.ContinueOnError)
	startCmd.Usage = func() {
		fmt.Print(
			"Usage: instances start INSTANCE_NAME\n\n",
			"Start the instance INSTANCE_NAME\n\n",
		)
		startCmd.PrintDefaults()
	}

	name, err := parseInstanceName(startCmd, args)
	if err != nil {
		return err
	}

	instance, err := c.db.GetInstance(name)
	if err != nil {
		return err
	}

	cloudProvider, err := instance.GetCloudProvider(c.cloudProviders)
	if err != nil {
		return err
	}

	err = cloudProvider.StartInstance(instance.Id)
	if err != nil {
		return err
	}

	return nil
}

func (c *CLI) stopInstance(args []string) error {
	stopCmd := flag.NewFlagSet("stop", flag.ContinueOnError)
	stopCmd.Usage = func() {
		fmt.Print(
			"Usage: instances stop INSTANCE_NAME\n\n",
			"Stop the instance INSTANCE_NAME\n\n",
		)
		stopCmd.PrintDefaults()
	}

	name, err := parseInstanceName(stopCmd, args)
	if err != nil {
		return err
	}

	instance, err := c.db.GetInstance(name)
	if err != nil {
		return err
	}

	cloudProvider, err := instance.GetCloudProvider(c.cloudProviders)
	if err != nil {
		return err
	}

	err = cloudProvider.StopInstance(instance.Id)
	if err != nil {
		return err
	}

	return nil
}

func (c *CLI) listInstances(args []string) error {
	var cloudName string
	listCmd := flag.NewFlagSet("list", flag.ContinueOnError)
	listCmd.StringVar(&cloudName, "cloud", "", "the cloud provider to list instances from")
	listCmd.Usage = func() {
		fmt.Print(
			"Usage: instances list [OPTIONS]\n\n",
			"List the instances\n\n",
		)
	}

	err := listCmd.Parse(args)
	if err != nil {
		return err
	}

	if len(listCmd.Args()) > 0 {
		return errors.New("list doesn't take positional arguments")
	}

	for name, instance := range c.db.Instances {
		if cloudName != "" && !strings.EqualFold(instance.CloudProviderName, cloudName) {
			continue
		}

		fmt.Printf("name: %s\tid: %s\tcloud provider: %s\n", name, instance.Id, instance.CloudProviderName)
	}

	return nil
}

func parseInstanceName(cmd *flag.FlagSet, args []string) (string, error) {
	err := cmd.Parse(args)
	if err != nil {
		return "", err
	}

	if len(cmd.Args()) == 0 {
		cmd.Usage()
		return "", errors.New("missing instance name")
	}

	if len(cmd.Args()) > 1 {
		cmd.Usage()
		return "", errors.New("only one instance name can be provided")
	}

	return cmd.Arg(0), nil
}
