package blocktemplates

import (
	htmlTemplate "html/template"
	textTemplate "text/template"
	"text/template/parse"
)

func FixTemplateTree(tree *parse.Tree) *parse.Tree {
	tree = tree.Copy()
	adjustNode(nil, tree.Root, nil)
	return tree
}

func FixLineStatementsText(t *textTemplate.Template) *textTemplate.Template {
	tClone, _ := t.Clone()
	tClone.Tree = FixTemplateTree(tClone.Tree)
	return tClone
}

func FixLineStatementsHtml(t *htmlTemplate.Template) *htmlTemplate.Template {
	tClone, _ := t.Clone()
	tClone.Tree = FixTemplateTree(tClone.Tree)
	return tClone
}
