package Dt_node

type BinaryNode struct {
	Left  *BinaryNode
	Right *BinaryNode
	Data  []float64
}

func (n *BinaryNode) add_left_node(col int, split float64) {

	n.Left = &BinaryNode{Data: []float64{float64(col), split}, Left: nil, Right: nil}

}

func (n *BinaryNode) add_right_node(col int, split float64) {

	n.Right = &BinaryNode{Data: []float64{float64(col), split}, Left: nil, Right: nil}
}
