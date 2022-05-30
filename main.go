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

	//TUNABLE PARAMETERES IN CAPS

	//Seed for randomness
	rand.Seed(time.Now().UnixNano())

	//CHOOSE YOUR DATA
	//data := open_csv("heart_failure_clinical_records_dataset.csv")
	//data := open_csv("Algerian_forest_fires_dataset_UPDATE.csv")
	//data := open_csv("SouthGermanCredit.csv")
	data := open_csv("Raisin_Dataset.csv")
	//data := open_csv("test.csv")

	//Which rows to use as training data?
	begin_train := 100
	end_train := 800

	training_data := data[begin_train:end_train]

	//Write predictions to file
	f, _ := os.Create("prediction.txt")
	defer f.Close()

	average_correct := 0.

	//HOW MANY FORESTS TO CREATE
	reps := 10

	for i := 0; i < reps; i++ {
		total := 0.
		correct := 0.

		//FOREST ARGUMENTS plant_forest(number of trees, max tree depth, features to randomly sample)
		//Build forest
		forest := plant_forest(training_data, 50, 3, 3)
		for j := 0; j < len(data); j++ {

			//Exclude training data
			if j < begin_train || j > end_train {
				class := classify_forest(data[j], forest)

				//Check if correct
				if class == int(data[j][len(data[j])-1]) {
					correct++
					//fmt.Println("Correct!")
				} else {
					//fmt.Println("Not correct...")
				}

				//Write to file
				_, err := f.WriteString(strconv.Itoa(class) + "\n")
				if err != nil {
					log.Fatal(err)
				}

				total++
			}
		}

		//Prints correct predictions for each forest as percentage
		fmt.Printf("Correct: %v%%\n", correct/total*100)
		average_correct += correct / total * 100
	}

	//Prints average correct predictions for all the forests
	average_correct = average_correct / float64(reps)

	fmt.Printf("Average correct: %v%%\n", average_correct)

	//All below was used to generate random predictions
	f1, _ := os.Create("random.txt")
	defer f1.Close()

	for i := 0; i < len(data); i++ {
		if i <= begin_train || i > end_train {
			_, err := f1.WriteString(strconv.Itoa(rand.Intn(2)) + "\n")
			if err != nil {
				log.Fatal(err)
			}
		}
	}

}

func open_csv(path string) (data [][]float64) {

	//Open file
	file, _ := os.Open(path)
	reader := csv.NewReader(file)

	//Ignore header
	_, _ = reader.Read()

	//Endless loop until end of file that reads lines and stores them in a 2d slice
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

	fmt.Println("CSV opened!")

	return data
}

func plant_forest(data [][]float64, ntrees int, max_depth int, mtry int) (forest []*DecisionTree) {

	//Loop for as many trees as defined by user
	for i := 0; i < ntrees; i++ {

		var bag [][]float64

		//Sample training data using bootstrag aggregation
		for j := 0; j < len(data); j++ {
			bag = append(bag, data[rand.Intn(len(data))])

		}

		//Initialise new tree and add to forest
		root := &Node{Data: nil, Left: nil, Right: nil}
		tree := &DecisionTree{Root: root}
		forest = append(forest, tree)

		//Build tree using populate_dt_node
		populate_dt_node(bag, 0, tree.Root, mtry, max_depth)

		fmt.Printf("%v trees constructed\n", i)
	}

	return forest
}

func classify_forest(row []float64, forest []*DecisionTree) int {

	//If no trees in forest somethin ain't right
	if len(forest) == 0 {
		fmt.Println("No trees in forest!")
		return 99
	}

	//Vars for counting votes
	var yay int
	var nay int

	//Loop over forest collecting votes
	for i := 0; i < len(forest); i++ {

		result := classify_tree(row, forest[i])

		if result == 1 {
			yay++
		} else {
			nay++
		}
	}

	//fmt.Printf("Votes for 1: %v\n", yay)
	//fmt.Printf("Votes for 0: %v\n", nay)

	//If more yays than nays classify as 1
	if yay > nay {
		return 1
	}

	//Classify as 0
	return 0
}

func classify_tree(row []float64, decision_tree *DecisionTree) int {

	//Set curren node to tree root
	node := decision_tree.Root

	//Loop over tree until leaf node reached
	for i := 0; i >= 0; i++ {

		//Update splt values to values from current node
		column := node.Data[i][0]
		split := node.Data[i][1]
		direction := node.Data[i][2]

		//If data >= split go right, if less go left. If lead node reached return classification
		if row[int(column)] >= split {
			node = node.Right
			if node == nil && direction == 1 {
				return 1
			} else if node == nil && direction == 0 {
				return 0
			}
		} else if row[int(column)] < split {
			node = node.Left
			if node == nil && direction == 0 {
				return 1
			} else if node == nil && direction == 1 {
				return 0
			}
		}
	}

	//Something went wrong as classification was not reached!
	return 99
}

func populate_dt_node(data [][]float64, current_depth int, node *Node, mtry int, max_depth int) {

	var column int

	//Variables to sore the best values when sampling different features
	var best_column int
	best_gini := 1.
	var best_split float64
	var best_direction int
	var best_rows_above [][]float64
	var best_rows_below [][]float64

	//Sample mtry number of random features
	for i := 0; i < mtry; i++ {

		//Generate random column number that has not been used already
		for {
			column = rand.Intn(len(data[0]) - 1)
			if !contains(node.Data, column) {
				break
			}
		}

		gini, split, direction, rows_above, rows_below := find_best_split(data, column)

		//Only retain if it's the best gini
		if gini < best_gini {

			best_column = column
			best_gini = gini
			best_split = split
			best_direction = direction
			best_rows_above = rows_above
			best_rows_below = rows_below
		}

	}

	//Add best split data to current node
	node.Add_data(best_column, best_split, best_direction)

	//If reached current depth return
	if current_depth+1 == max_depth {
		return
	}

	//If data is perfectly split return
	if best_gini == 0 {
		return
	}

	//If no rows are above or below the best split then training data is well split already
	if len(best_rows_above) == 0 || len(best_rows_below) == 0 {
		return
	}

	// Update current node pointers to new nodes
	node.Add_nodes(best_column, best_split, best_direction)

	//Recursively call itself to build the tree
	populate_dt_node(best_rows_below, current_depth+1, node.Left, mtry, max_depth)
	populate_dt_node(best_rows_above, current_depth+1, node.Right, mtry, max_depth)

}

func find_best_split(data [][]float64, column int) (best_gini float64, best_split float64, best_direction int, best_rows_above [][]float64, best_rows_below [][]float64) {

	//Slice to hold sorted training data
	sorted_data := make([][]float64, len(data))

	//https://stackoverflow.com/questions/45465368/golang-multidimensional-slice-copy
	for i := range data {
		sorted_data[i] = make([]float64, len(data[i]))
		copy(sorted_data[i], data[i])
	}

	//Sort slice
	sort.Slice(sorted_data, func(i, j int) bool { return sorted_data[i][column] < sorted_data[j][column] })

	lower := 0.
	best_gini = 1

	//Loop over training data and calculate weighted gini impurity for each split point
	for _, row := range sorted_data {

		//Split points are mid points between rows
		split := (lower + row[column]) / 2

		gini, direction, rows_above, rows_below := gini_impurity(sorted_data, column, split)

		//Only retain if we have the lowest weighted gini impurity
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

	//Initialise variables for counting rows above and below the split.
	above_label0 := 0.
	below_label0 := 0.
	above_label1 := 0.
	below_label1 := 0.

	//Loop over training data counting rows for the above variables
	for _, row := range data {

		if row[column] < threshold {
			rows_below = append(rows_below, row) //Add row to be returned to calling func

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

	gini_above := 0.
	gini_below := 0.

	//Calculating gini impurity for above and below the split
	if total_above != 0 {
		gini_above = 1 - math.Pow(above_label0/total_above, 2) - math.Pow(above_label1/total_above, 2)
	}

	if total_below != 0 {
		gini_below = 1 - math.Pow(below_label0/total_below, 2) - math.Pow(below_label1/total_below, 2)
	}

	//Counting labels above and below the split to work out the direction
	ratio_above := above_label1 / total_above
	ratio_below := below_label1 / total_below

	//Deciding direction, random if equal ratios
	if ratio_above > ratio_below {
		direction = 1
	} else if ratio_above < ratio_below {
		direction = 0
	} else {
		direction = rand.Intn(2)
	}

	//Calculate weighted gini impurity
	weighted_gini := (total_above/total)*gini_above + (total_below/total)*gini_below

	return weighted_gini, direction, rows_above, rows_below
}

//https://stackoverflow.com/questions/10485743/contains-method-for-a-slice
func contains(s [][]float64, e int) bool {
	for _, a := range s {
		if int(a[0]) == e {
			return true
		}
	}
	return false
}

//Decsion tree
type DecisionTree struct {
	Root *Node
}

//Decsion tree node
type Node struct {
	Left  *Node
	Right *Node
	Data  [][]float64
}

//For appending the current nodes split
func (n *Node) Add_data(col int, split float64, direction int) {

	n.Data = append(n.Data, []float64{float64(col), split, float64(direction)})

}

//Populating the Left and Right pointers
func (n *Node) Add_nodes(col int, split float64, direction int) {

	n.Left = &Node{Data: n.Data, Left: nil, Right: nil}
	n.Right = &Node{Data: n.Data, Left: nil, Right: nil}

}
