package utils

import "text/template/parse"

func AsBranchNode(node parse.Node) (*parse.BranchNode, bool) {
	switch n := node.(type) {
	case *parse.RangeNode:
		return &n.BranchNode, true
	case *parse.IfNode:
		return &n.BranchNode, true
	case *parse.WithNode:
		return &n.BranchNode, true
	}
	return nil, false
}
