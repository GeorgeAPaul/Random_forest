package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"
)

func main() {

	rand.Seed(time.Now().UnixNano())
	//data := open_csv("test.csv")
	//training_data := data[:10]
	//data := open_csv("Reduced Features for TAI project.csv")
	data := open_csv("Raisin_Dataset.csv")

	training_data := data[100:801]

	//forest := plant_forest(training_data, 100, 6)

	// for i := 0; i < len(forest); i++ {
	// 	fmt.Println(forest[i].Root)
	// }

	f, _ := os.Create("prediction.txt")
	defer f.Close()

	average_correct := 0.
	reps := 1

	for i := 0; i < reps; i++ {
		correct := 0.
		forest := plant_forest(training_data, 100, 6)
		for j := 0; j < len(data); j++ {
			if j < 100 || j >= 800 {
				class := classify_forest(data[j], forest)
				if class == int(data[j][len(data[j])-1]) {
					correct++
				}
				_, err := f.WriteString(strconv.Itoa(class) + "\n")
				if err != nil {
					log.Fatal(err)
				}
			}
		}
		fmt.Printf("Correct: %v\n", correct)
		average_correct += correct
	}

	average_correct = average_correct / float64(reps)

	fmt.Printf("Average correct: %v\n", average_correct)

	// f1, _ := os.Create("random.txt")
	// defer f1.Close()

	// for i := 0; i < len(data); i++ {
	// 	if i <= 100 || i > 800 {
	// 		_, err := f1.WriteString(strconv.Itoa(rand.Intn(2)) + "\n")
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}
	// 	}
	// }

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

func plant_forest(data [][]float64, ntrees int, mtry int) (forest []*DecisionTree) {

	for i := 0; i <= ntrees; i++ {

		//bagging_size := len(data)
		var bag [][]float64

		for j := 0; j < len(data); j++ {
			bag = append(bag, data[rand.Intn(len(data))])

		}

		//fmt.Printf("Bag: %v\n", i)
		//fmt.Println(bag)

		root := &Node{Data: nil, Left: nil, Right: nil}
		tree := &DecisionTree{Root: root}
		forest = append(forest, tree)

		populate_dt_node(bag, 0, tree.Root, mtry)

		//fmt.Printf("Tree: %v\n", tree.Root.Data)
		//fmt.Printf("TreeLeft: %v\n", tree.Root.Left.Data)
		//fmt.Printf("TreeRight: %v\n", tree.Root.Right.Data)

		fmt.Printf("Tree done %v\n", i)
	}

	return forest
}

func classify_forest(row []float64, forest []*DecisionTree) int {

	var yay int
	var nay int

	for i := 0; i < len(forest); i++ {

		//fmt.Printf("Tree: %v\n", i)
		//fmt.Printf("Tree data: %v\n", forest[i].Root.Data)
		//fmt.Printf("Tree Left: %v\n", forest[i].Root.Left)
		//fmt.Printf("Tree Right: %v\n", forest[i].Root.Right)

		result := classify_tree(row, forest[i])

		if result == 1 {
			yay++
		} else {
			nay++
		}
	}

	fmt.Printf("yay: %v\n", yay)
	fmt.Printf("nay: %v\n", nay)
	if yay > nay {
		return 1
	}

	return 0
}

func classify_tree(row []float64, decision_tree *DecisionTree) int {

	node := decision_tree.Root

	for i := 0; i >= 0; i++ {

		//fmt.Println(node.Data)

		column := node.Data[i][0]
		split := node.Data[i][1]
		direction := node.Data[i][2]

		//fmt.Printf("Column %v\n", column)
		//fmt.Printf("Split %v\n", split)

		if row[int(column)] >= split {
			node = node.Right
			//fmt.Println("Right")
			if node == nil && direction == 1 {
				return 1
			} else if node == nil && direction == 0 {
				return 0
			}
		} else if row[int(column)] < split {
			node = node.Left
			//fmt.Println("Right")
			if node == nil && direction == 0 {
				return 1
			} else if node == nil && direction == 1 {
				return 0
			}
		}
	}

	return 99
}

func populate_dt_node(data [][]float64, current_depth int, node *Node, mtry int) {

	//fmt.Println("Adding node")

	max_depth := len(data[0]) - 1

	var column int

	var best_column int
	best_gini := 1.
	var best_split float64
	var best_direction int
	var best_rows_above [][]float64
	var best_rows_below [][]float64

	for i := 0; i < mtry; i++ {

		for {
			//rand.Seed(time.Now().UnixNano())
			column = rand.Intn(len(data[0]) - 1)
			if !contains(node.Data, column) {
				break
			}
		}

		gini, split, direction, rows_above, rows_below := find_best_split(data, column)

		if gini < best_gini {
			//fmt.Printf("Tries %v\n", i)
			//fmt.Printf("Gini %v\n", gini)
			//fmt.Printf("Column %v\n", split)
			//fmt.Printf("Direction %v\n", direction)

			best_column = column
			best_gini = gini
			best_split = split
			best_direction = direction
			best_rows_above = rows_above
			best_rows_below = rows_below
		}

	}

	//fmt.Println(best_column)

	node.Add_data(best_column, best_split, best_direction)

	if current_depth+1 == max_depth {
		return
	}
	if len(best_rows_above) == 0 || len(best_rows_below) == 0 {
		return
	}

	node.Add_nodes(best_column, best_split, best_direction)

	populate_dt_node(best_rows_below, current_depth+1, node.Left, mtry)
	populate_dt_node(best_rows_above, current_depth+1, node.Right, mtry)

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

		//fmt.Println(gini, split, direction)

		if gini < best_gini {
			best_gini = gini
			best_split = split
			best_direction = direction
			best_rows_above = rows_above
			best_rows_below = rows_below
		}

		lower = row[column]

	}

	return best_gini, best_split, best_direction, best_rows_above, best_rows_below
}

func gini_impurity(data [][]float64, column int, threshold float64) (gini_index float64, direction int, rows_above [][]float64, rows_below [][]float64) {

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
	total := float64(len(data))

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
		//rand.Seed(time.Now().UnixNano())
		direction = rand.Intn(2)
	}

	weighted_gini := (total_above/total)*gini_above + (total_below/total)*gini_below

	return weighted_gini, direction, rows_above, rows_below
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

func (n *Node) Add_data(col int, split float64, direction int) {

	n.Data = append(n.Data, []float64{float64(col), split, float64(direction)})

}

func (n *Node) Add_nodes(col int, split float64, direction int) {

	n.Left = &Node{Data: n.Data, Left: nil, Right: nil}
	n.Right = &Node{Data: n.Data, Left: nil, Right: nil}

}
