package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
)

func main() {
	fmt.Println("The project begins...")
	//col := open_csv_column("Reduced Features for TAI project.csv", 4)
	//label := open_csv_column("Reduced Features for TAI project.csv", 151)
	data := open_csv("test.csv")
	fmt.Println(data)
	fmt.Print(gini_index(data, 0))

}

func open_csv(path string) (data [][]float64) {

	file, _ := os.Open(path)
	reader := csv.NewReader(file)

	//deal with header
	_, _ = reader.Read()

	for i := 0; i >= 0; i++ {

		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		data = append(data, make([]float64, 0))

		for column := range record {
			record_float, _ := strconv.ParseFloat(record[column], 64)
			data[i] = append(data[i], record_float)
		}
	}
	return data
}

func gini_index(data [][]float64, column int) (gini_index float64) {

	threshold := 5.5
	//total := float64(len(data))

	var rows_above [][]float64
	var rows_below [][]float64

	above_label0 := 0.
	below_label0 := 0.
	above_label1 := 0.
	below_label1 := 0.

	for _, row := range data {

		if row[column] < threshold {
			rows_below = append(rows_below, row)

			if row[len(row)-1] == 0 {
				below_label0++
			} else {
				below_label1++
			}
		} else {
			rows_above = append(rows_above, row)

			if row[len(row)-1] == 0 {
				above_label0++
			} else {
				above_label1++
			}
		}
	}

	total_above := above_label0 + above_label1
	total_below := below_label0 + below_label1

	fmt.Printf("Above %v\n", rows_above)
	fmt.Printf("Below %v\n", rows_below)

	gini_above := 1 - math.Pow(above_label0/total_above, 2) - math.Pow(above_label1/total_above, 2)
	gini_below := 1 - math.Pow(below_label0/total_below, 2) - math.Pow(below_label1/total_below, 2)

	average_gini := (gini_above + gini_below) / 2
	return average_gini //gini
}
