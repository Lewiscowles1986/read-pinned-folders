package main

import (
	"fmt"
	"os"
	"path"

	"github.com/richardlehane/mscfb"
)

func main() {
	homeDir, _ := os.UserHomeDir()
	roamingDir := path.Join(homeDir, "AppData", "Roaming")
	windowsRoamingDir := path.Join(roamingDir, "Microsoft", "Windows")
	automaticDestinations := path.Join(windowsRoamingDir, "Recent", "AutomaticDestinations")
	filePath := path.Join(automaticDestinations, "f01b4d95cf55d32a.automaticDestinations-ms")

	file, _ := os.Open(filePath)
	defer file.Close()

	doc, err := mscfb.New(file)
	if err != nil {
		panic(err)
	}

	paths, err := parseAutomaticDestinationFile(doc, file)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%d Paths found\n", len(paths))
	for _, path := range paths {
		fmt.Println(path)
	}
}
