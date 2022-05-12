package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"
)

func main() {
	data := open_csv("test.csv")
	//data := open_csv("Reduced Features for TAI project.csv")
	//fmt.Println(data)

	//fmt.Print(gini_index(data, 0, 5.5))

	//gini, split := find_best_split(data, 0)
	//fmt.Println(gini)
	//fmt.Println(split)

	fmt.Println(create_decision_tree(data))

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

func create_decision_tree(data [][]float64) (tree [][]float64) {

	var column_index []int

	for i := 0; i < len(data[0])-1; i++ {

		column_index = append(column_index, i)
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(column_index), func(i, j int) { column_index[i], column_index[j] = column_index[j], column_index[i] })

	for _, column := range column_index {
		_, split := find_best_split(data, column)
		tree = append(tree, []float64{float64(column), split})
	}

	return tree
}

func find_best_split(data [][]float64, column int) (best_gini float64, best_split float64) {

	sort.Slice(data, func(i, j int) bool { return data[i][column] < data[j][column] })

	lower := 0.
	best_gini = 1

	for _, row := range data {

		split := (lower + row[column]) / 2

		gini := gini_index(data, column, split)

		if gini < best_gini {
			best_gini = gini
			best_split = split
		}
		//fmt.Println(lower)
		//fmt.Println(row[column])
		//fmt.Println(gini)
		//fmt.Println(best_gini)

		lower = row[column]

	}

	return best_gini, best_split
}

func gini_index(data [][]float64, column int, threshold float64) (gini_index float64) {

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

	//fmt.Printf("Above %v\n", rows_above)
	//fmt.Printf("Below %v\n", rows_below)

	gini_above := 0.
	gini_below := 0.

	if total_above != 0 {
		gini_above = 1 - math.Pow(above_label0/total_above, 2) - math.Pow(above_label1/total_above, 2)
	}

	if total_below != 0 {
		gini_below = 1 - math.Pow(below_label0/total_below, 2) - math.Pow(below_label1/total_below, 2)
	}

	//fmt.Printf("Gini Above %v\n", gini_above)
	//fmt.Printf("Gini Below %v\n", gini_below)

	average_gini := (gini_above + gini_below) / 2

	return average_gini //gini
}
