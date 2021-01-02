package treeprinter

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"text/template/parse"

	"git.k3s.test/nonlogical/golang-block-templates/internal/utils"
)

type NodeTypeStringer parse.NodeType

func (t NodeTypeStringer) String() string {
	switch parse.NodeType(t) {
	case parse.NodeText:
		return "NodeText"
	case parse.NodeList:
		return "NodeList"
	case parse.NodeRange:
		return "NodeRange"
	case parse.NodeIf:
		return "NodeIf"
	case parse.NodeWith:
		return "NodeWith"
	case parse.NodeAction:
		return "NodeAction"
	default:
		return fmt.Sprintf("? %v ?", t)
	}
}

func PrintNode(w io.Writer, indent string, node parse.Node) {
	if reflect.ValueOf(node).IsNil() {
		return
	}

	subIndent := indent+indent
	nodeType := NodeTypeStringer(node.Type())

	// Print Header:
	_, _ = fmt.Fprintf(w, "%s[%#v]:", indent, nodeType)

	if n, ok := node.(*parse.ListNode); ok {
		_, _ = fmt.Fprintf(w, "\n")
		for _, cn := range n.Nodes {
			PrintNode(w, subIndent, cn)
		}
	} else if n, ok := node.(*parse.ActionNode); ok {
		_, _ = fmt.Fprintf(w, " %s\n", repr(n.String()))
	} else if n, ok := node.(*parse.TextNode); ok {
		_, _ = fmt.Fprintf(w, " %s\n", repr(n.Text))
	} else if n, ok := utils.AsBranchNode(node); ok {
		_, _ = fmt.Fprintf(w, "\n |- (List):\n")
		PrintNode(w, subIndent, n.List)
		_, _ = fmt.Fprintf(w, "\n |- (ElseList):\n")
		PrintNode(w, subIndent, n.ElseList)
	} else {
		_, _ = fmt.Fprintf(w, " %T", node)
	}
}

func repr(i interface{}) string {
	out, _ := json.Marshal(i)
	return string(out)
}
