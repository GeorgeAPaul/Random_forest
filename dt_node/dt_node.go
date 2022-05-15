package dt_node

type BinaryNode struct {
	Left  *BinaryNode
	Right *BinaryNode
	Data  [][]float64
}

func (n *BinaryNode) Add_nodes(col int, split float64) {

	n.Data = append(n.Data, []float64{float64(col), split})

	n.Left = &BinaryNode{Data: n.Data, Left: nil, Right: nil}
	n.Right = &BinaryNode{Data: n.Data, Left: nil, Right: nil}

}
