package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

func main() {
	fmt.Println("The project begins...")
	open_csv("Reduced Features for TAI project.csv")

}

func open_csv(path string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file", err)
	}
	reader := csv.NewReader(file)
	records, _ := reader.ReadAll()
	fmt.Println(records)
}
