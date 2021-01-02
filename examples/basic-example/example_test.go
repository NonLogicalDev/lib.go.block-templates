package main

import (
	"fmt"
	"os"
	"text/template"

	blocktemplates "git.k3s.test/nonlogical/golang-block-templates"
)

type dict = map[string]interface{}
type list = []interface{}

var data = dict{
	"List": list{
		dict{"Name": "1", "SubList": list{"a", "b", "c"}},
		dict{"Name": "2", "SubList": nil},
	},
}

const _tplLineBlock = `
:BEGIN A:
	{{ range .List }}
{{.Name}}:
	{{end}}
:END A:

:BEGIN B:
	{{ range .List }}
:ITEM: {{ .Name }}
		{{ range .SubList }}
			{{ if (eq . "a") }}
  :SI:IF CLAUSE: {{ . }}
			{{ else if (eq . "b") }}
  :SI:ELSE IF CLAUSE: {{ . }}
			{{ else }}
  :SI:ELSE CLAUSE: {{ . }}
			{{ end }}
		{{ else }}
  :NOTHING:
		{{ end }}
	{{ end }}
:END B:
`

const _tplManualFix = `
:BEGIN A:{{ range .List }}
{{.Name}}:{{end}}
:END A:

:BEGIN B:{{ range .List }}
:ITEM: {{ .Name }}{{ range .SubList }}{{ if (eq . "a") }}
  :SI:IF CLAUSE: {{ . }}{{ else if (eq . "b") }}
  :SI:ELSE IF CLAUSE: {{ . }}{{ else }}
  :SI:ELSE CLAUSE: {{ . }}{{ end }}{{ else }}
  :NOTHING:{{ end }}{{ end }}
:END B:
`

func ExampleBasic() {
	//fmt.Println("============================================================")
	//fmt.Println("Before Fix:")
	//fmt.Println("============================================================")
	//{
	//	tc, _ := template.New("(test)").Parse(_tpl)
	//	_ = tc.Execute(os.Stdout, data)
	//}

	fmt.Println("============================================================")
	fmt.Println("After Fix:")
	fmt.Println("============================================================")
	{
		tc, _ := template.New("(test)").Parse(_tplLineBlock)
		tc = blocktemplates.FixLineStatementsText(tc)
		_ = tc.Execute(os.Stdout, data)
	}

	// Output:
	// ============================================================
	// After Fix:
	// ============================================================
	//
	// :BEGIN A:
	// 1:
	// 2:
	// :END A:
	//
	// :BEGIN B:
	// :ITEM: 1
	//   :SI:IF CLAUSE: a
	//   :SI:ELSE IF CLAUSE: b
	//   :SI:ELSE CLAUSE: c
	// :ITEM: 2
	//   :NOTHING:
	// :END B:
}

