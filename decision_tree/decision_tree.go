package decision_tree

import (
	"github.com/GeorgeAPaul/Random_forest/dt_node"
)

type BinaryTree struct {
	Root *dt_node.BinaryNode
}

// func (b BinaryTree) String() string {
// 	return fmt.Sprintf("[%v,%v]", b.Root, "plop")
// }
