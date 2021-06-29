package main

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/grafana/loki/pkg/sizing"
	"github.com/grafana/loki/pkg/util/flagext"
)

type Config struct {
	BytesPerSecond flagext.ByteSize
}

func (c *Config) RegisterFlags(f *flag.FlagSet) {
	f.Var(&c.BytesPerSecond, "bytes-per-second", "[human readable] How many bytes per second the cluster should receive, i.e. (50MB)")
}

func (c *Config) Validate() error {
	if c.BytesPerSecond <= 0 {
		return errors.New("must specify bytes-per-second")
	}
	return nil
}

func main() {
	var cfg Config
	cfg.RegisterFlags(flag.CommandLine)
	flag.Parse()
	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}

	cluster := sizing.SizeCluster(cfg.BytesPerSecond.Val())

	printClusterArchitecture(&cluster)
}

// TODO: Add verbose flag to include the "request" (min resources) in addition to "limit" (max resources)
func printClusterArchitecture(c *sizing.ClusterResources) {

	// loop through all components, and print out how many replicas of each component we're recommending.
	/*
		Format will look like
		"""
		Overall Requirements for a Loki cluster than can handle X volume of ingest
		Number of Nodes: 2
		Memory Required: 1000 MB
		CPUs Required: 34
		Disk Required: 100 GB

		List of all components in the Loki cluster, the number of replicas of each, and the resources required per replica

		Ingester: 5 replicas, each with:
			2000 MB RAM
			10 GB Disk
			5 CPU

		Distributor: 2 replicas, each with:
			1000 MB RAM
			1 GB Disk
			2 CPU
		"""
	*/

	totals := c.Totals()

	// TODO: Actually populate the value of X volume of ingest
	fmt.Println("Overall Requirements for a Loki cluster than can handle X volume of ingest")
	fmt.Printf("\tNumber of Nodes: %d\n", c.NumNodes())
	fmt.Printf("\tMemory Required: %v\n", totals.MemoryLimits)
	fmt.Printf("\tCPUs Required: %v\n", totals.CPULimits)
	fmt.Printf("\tDisk Required: %d GB\n", totals.DiskGB)

	fmt.Printf("\n")

	fmt.Printf("List of all components in the Loki cluster, the number of replicas of each, and the resources required per replica\n")

	for _, component := range c.Components() {
		if component != nil {
			fmt.Printf("%v: %d replicas, each of which requires\n", component.Name, component.Replicas)
			fmt.Printf("\t%v MB of memory\n", component.Resources.MemoryLimits)
			fmt.Printf("\t%v CPU\n", component.Resources.CPULimits)
			fmt.Printf("\t%d GB of disk\n", component.Resources.DiskGB)
		}
	}

}