package dt_node

type BinaryNode struct {
	Left  *BinaryNode
	Right *BinaryNode
	Data  [][]float64
}

func (n *BinaryNode) Add_nodes(col int, split float64) {

	n.Left = &BinaryNode{Data: append(n.Data, []float64{float64(col), split, 0.}), Left: nil, Right: nil}
	n.Right = &BinaryNode{Data: append(n.Data, []float64{float64(col), split, 1.}), Left: nil, Right: nil}

}
