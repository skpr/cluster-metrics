package main

import "gopkg.in/alecthomas/kingpin.v2"

var cliFrequency = kingpin.Flag("frequency", "How often to poll for metrics data").Default("60s").Duration()

func main() {
	kingpin.Parse()
}
