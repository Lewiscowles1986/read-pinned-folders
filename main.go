package main

import "fmt"

func main() {
	paths, err := parseAutomaticDestinationFile("C:\\Users\\lewis\\AppData\\Roaming\\Microsoft\\Windows\\Recent\\AutomaticDestinations\\f01b4d95cf55d32a.automaticDestinations-ms")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d Paths found\n", len(paths))
	for _, path := range paths {
		fmt.Println(path)
	}
}
