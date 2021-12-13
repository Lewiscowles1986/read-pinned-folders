package main

import (
	"fmt"
	"os"
	"path"
)

func main() {
	homeDir, _ := os.UserHomeDir()
	roamingDir := path.Join(homeDir, "AppData", "Roaming")
	windowsRoamingDir := path.Join(roamingDir, "Microsoft", "Windows")
	automaticDestinations := path.Join(windowsRoamingDir, "Recent", "AutomaticDestinations")
	paths, err := parseAutomaticDestinationFile(
		path.Join(automaticDestinations, "f01b4d95cf55d32a.automaticDestinations-ms"),
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d Paths found\n", len(paths))
	for _, path := range paths {
		fmt.Println(path)
	}
}
