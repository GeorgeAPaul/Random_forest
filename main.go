package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

func main() {
	fmt.Println("The project begins...")
	fmt.Println(open_csv_column("Reduced Features for TAI project.csv", 4))
	//open_csv_column("test.csv", 1)

}

func open_csv_column(path string, column int) (a []float64) {

	file, _ := os.Open(path)
	reader := csv.NewReader(file)

	var data []float64

	//deal with header
	_, _ = reader.Read()

	for {

		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		record_float, _ := strconv.ParseFloat(record[column], 64)

		data = append(data, record_float)
	}

	return data
}
