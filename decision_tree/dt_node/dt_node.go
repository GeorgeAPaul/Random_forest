package dt_node

type BinaryNode struct {
	left  *BinaryNode
	right *BinaryNode
	data  []float64
}

func (n *BinaryNode) add_left_node(col int, split float64) {

	n.left = &BinaryNode{data: []float64{float64(col), split}, left: nil, right: nil}

}

func (n *BinaryNode) add_right_node(col int, split float64) {

	n.right = &BinaryNode{data: []float64{float64(col), split}, left: nil, right: nil}
}
