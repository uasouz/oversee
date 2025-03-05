package main

import "oversee/collector"

func main() {
	collectorApi := collector.NewCollectorAPI()

	collectorApi.Serve()
}
