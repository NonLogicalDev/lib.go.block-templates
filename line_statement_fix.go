package blocktemplates

import (
	"fmt"
	"reflect"
	"text/template/parse"

	"git.k3s.test/nonlogical/golang-block-templates/internal/utils"
)

func adjustNode(prevText *parse.TextNode, node parse.Node, nextText *parse.TextNode) {
	if reflect.ValueOf(node).IsNil() {
		return
	}
	if n, ok := node.(*parse.ListNode); ok {
		adjustListNode(n)
	} else if bn, ok := utils.AsBranchNode(node); ok {
		adjustBranchNode(prevText, bn, nextText)
	} else if _, ok := node.(*parse.TextNode); ok {
		// pass
	} else if _, ok := node.(*parse.ActionNode); ok {
		// pass
	} else {
		panic(fmt.Sprintf("%T", node))
	}
}

func adjustListNode(node *parse.ListNode) {
	for i, cn := range node.Nodes {
		var (
			prevTextIdx = i - 1
			prevText *parse.TextNode

			nextTextIdx = i + 1
			nextText *parse.TextNode
		)

		if prevTextIdx >= 0 {
			prevText, _ = node.Nodes[prevTextIdx].(*parse.TextNode)
		}
		if nextTextIdx <= len(node.Nodes) - 1 {
			nextText, _ = node.Nodes[nextTextIdx].(*parse.TextNode)
		}

		adjustNode(prevText, cn, nextText)
	}
}

func adjustBranchNode(prevText *parse.TextNode, node *parse.BranchNode, nextText *parse.TextNode) {
	// Handling Branch Nodes:
	//
	// [A]{{branch}}[B]
	//   TB1
	// [C]{{else}}[D]
	//   TB2
	// [E]{{end}}[F]
	//
	// if
	//    !validEnd(A) || !validStart(B) ||
	//    !validEnd(C) || !validStart(D) ||
	//    !validEnd(E) || !validStart(F)
	//  --> !bail!
	//
	// if validStart(B) && validEnd(A) -> A = trimLastLine(A)
	// if validStart(D) && validEnd(C) -> C = trimLastLine(C)
	// if validStart(F) && validEnd(E) -> E = trimLastLine(E)

	var (
		aText *parse.TextNode
		bText *parse.TextNode
		cText *parse.TextNode
		dText *parse.TextNode
		eText *parse.TextNode
		fText *parse.TextNode
	)

	// Find all text nodes:
	aText = prevText
	if node.List != nil && len(node.List.Nodes) > 0 {
		bText, _ = node.List.Nodes[0].(*parse.TextNode)
		cText, _ = node.List.Nodes[len(node.List.Nodes)-1].(*parse.TextNode)
		adjustNode(nil, node.List, nil)
	}
	if node.ElseList != nil && len(node.ElseList.Nodes) > 0 {
		dText, _ = node.ElseList.Nodes[0].(*parse.TextNode)
		eText, _ = node.ElseList.Nodes[len(node.ElseList.Nodes)-1].(*parse.TextNode)
		adjustNode(nil, node.ElseList, nil)
	}
	fText = nextText

	var completionHooks []func()

	// Make sure all text nodes are valid for block modify.
	if aText != nil {
		modifiedText, ok := trimEndWhitespace(aText.Text)
		if !ok {
			return
		}
		completionHooks = append(completionHooks, func() {
			aText.Text = modifiedText
		})
	}

	if bText != nil {
		modifiedText, ok := trimStartWhitespace(bText.Text)
		if !ok {
			return
		}
		completionHooks = append(completionHooks, func() {
			bText.Text = modifiedText
		})
	}

	if cText != nil {
		modifiedText, ok := trimEndWhitespace(cText.Text)
		if !ok {
			return
		}
		completionHooks = append(completionHooks, func() {
			cText.Text = modifiedText
		})
	}

	if dText != nil {
		modifiedText, ok := trimStartWhitespace(dText.Text)
		if !ok {
			return
		}
		completionHooks = append(completionHooks, func() {
			dText.Text = modifiedText
		})
	}

	if eText != nil {
		modifiedText, ok := trimEndWhitespace(eText.Text)
		if !ok {
			return
		}
		completionHooks = append(completionHooks, func() {
			eText.Text = modifiedText
		})
	}

	if fText != nil {
		modifiedText, ok := trimStartWhitespace(fText.Text)
		if !ok {
			return
		}
		completionHooks = append(completionHooks, func() {
			fText.Text = modifiedText
		})
	}


	// Run Completion Hooks
	for _, hook := range completionHooks {
		hook()
	}
}

