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
	//data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	//data := open_csv("Reduced Features for TAI project.csv")

	//fmt.Println(find_best_split(data, 0))
	//fmt.Println(find_best_split(data, 1))
	//fmt.Println(find_best_split(data, 2))

	fmt.Println(data)

	root := &Node{Data: nil, Left: nil, Right: nil}
	tree := &DecisionTree{Root: root}

	populate_dt_node(data, 0, tree.Root)

	//rand.Seed(time.Now().UnixNano())
	//test_row := 2 //rand.Intn(len(data) - 1)

	//fmt.Printf("Classified %v\n", classify(data[test_row], tree))

	fmt.Printf("%+v\n", tree)
	fmt.Printf("%+v\n", *tree.Root)
	fmt.Printf("%+v\n", *tree.Root.Left)
	//fmt.Printf("%+v\n", *tree.Root.Left.Left)
	//fmt.Printf("%+v\n", *tree.Root.Left.Left.Left)
	// fmt.Printf("%+v\n", *tree.Root.Left)
	// fmt.Printf("%+v\n", *tree.Root.Right)
	// fmt.Printf("%+v\n", *tree.Root.Left.Left)
	// fmt.Printf("%+v\n", *tree.Root.Left.Right)
	// fmt.Printf("%+v\n", *tree.Root.Right.Right)
	// fmt.Printf("%+v\n", *tree.Root.Right.Left)

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

func classify(row []float64, decision_tree *DecisionTree) int {

	fmt.Printf("Row %v\n", row)

	node := decision_tree.Root

	for i := 0; i < 3; i++ {

		column := node.Data[i][0]
		split := node.Data[i][1]
		direction := node.Data[i][2]

		fmt.Printf("Column %v\n", column)
		fmt.Printf("Split %v\n", split)

		if row[int(column)] > split && direction == 1 {
			node = node.Right
			fmt.Println("Right")
			if node.Left == nil {
				return 1
			}
		} else if row[int(column)] < split && direction == 0 {
			node = node.Right
			fmt.Println("Right")
			if node.Left == nil {
				return 1
			}
		} else {
			node = node.Left
			fmt.Println("Left")
			if node.Left == nil {
				return 0
			}
		}
	}

	return 99
}

func populate_dt_node(data [][]float64, current_depth int, node *Node) {

	max_depth := len(data[0]) - 1

	//fmt.Println(max_depth)

	//if current_depth%1 == 0 {
	//fmt.Println(current_depth)
	//}

	if current_depth == max_depth {
		return
	} else {

		var column int

		for {
			rand.Seed(time.Now().UnixNano())
			column = rand.Intn(len(data[0]) - 1)
			if !contains(node.Data, column) {
				break
			}
		}

		_, split, direction, rows_above, rows_below := find_best_split(data, column)

		if len(rows_above) == 0 || len(rows_below) == 0 {
			fmt.Println("HEY")
			return
		}

		node.Add_nodes(column, split, direction)

		populate_dt_node(rows_below, current_depth+1, node.Left)
		populate_dt_node(rows_above, current_depth+1, node.Right)

	}
}

func find_best_split(data [][]float64, column int) (best_gini float64, best_split float64, best_direction int, best_rows_above [][]float64, best_rows_below [][]float64) {

	sorted_data := make([][]float64, len(data))

	// https://stackoverflow.com/questions/45465368/golang-multidimensional-slice-copy
	for i := range data {
		sorted_data[i] = make([]float64, len(data[i]))
		copy(sorted_data[i], data[i])
	}

	sort.Slice(sorted_data, func(i, j int) bool { return sorted_data[i][column] < sorted_data[j][column] })

	lower := 0.
	best_gini = 1

	for _, row := range sorted_data {

		split := (lower + row[column]) / 2

		gini, direction, rows_above, rows_below := gini_impurity(sorted_data, column, split)

		if gini < best_gini {
			best_gini = gini
			best_split = split
			best_direction = direction
			best_rows_above = rows_above
			best_rows_below = rows_below
		}
		//fmt.Println(lower)
		//fmt.Println(row[column])
		//fmt.Println(gini)
		//fmt.Println(best_split)

		lower = row[column]

	}

	return best_gini, best_split, best_direction, best_rows_above, best_rows_below
}

func gini_impurity(data [][]float64, column int, threshold float64) (gini_index float64, direction int, rows_above [][]float64, rows_below [][]float64) {

	//var rows_above [][]float64
	//var rows_below [][]float64

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

	total_above := float64(len(rows_above))
	total_below := float64(len(rows_below))

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

	ratio_above := above_label1 / total_above
	ratio_below := below_label1 / total_below

	if ratio_above > ratio_below {
		direction = 1
	} else if ratio_above < ratio_below {
		direction = 0
	} else {
		rand.Seed(time.Now().UnixNano())
		direction = rand.Intn(2)
	}

	//fmt.Printf("Gini Above %v\n", gini_above)
	//fmt.Printf("Gini Below %v\n", gini_below)

	average_gini := (gini_above + gini_below) / 2

	return average_gini, direction, rows_above, rows_below //gini
}

func contains(s [][]float64, e int) bool {
	for _, a := range s {
		if int(a[0]) == e {
			return true
		}
	}
	return false
}

type DecisionTree struct {
	Root *Node
}

type Node struct {
	Left  *Node
	Right *Node
	Data  [][]float64
}

func (n *Node) Add_nodes(col int, split float64, direction int) {

	n.Data = append(n.Data, []float64{float64(col), split, float64(direction)})

	n.Left = &Node{Data: n.Data, Left: nil, Right: nil}
	n.Right = &Node{Data: n.Data, Left: nil, Right: nil}

}
