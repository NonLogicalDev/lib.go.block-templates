package blocktemplates

import (
	"bytes"
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/require"
)

type dict = map[string]interface{}
type list = []interface{}

var data = dict{
	"List": list{
		dict{"Name": "1", "SubList": list{"a", "b", "c"}},
		dict{"Name": "2", "SubList": nil},
	},
}

func tplDef(s string) string {
	return strings.TrimSuffix(strings.TrimPrefix(s, "\n"), "\n")
}

func TestModify(t *testing.T) {
	table := []struct {
	name           string
	template       string
	templateManual string
}{
		{
			name: "ranges get transformed correctly",
			template: tplDef(`
:BEGIN A:
{{ range .List }}
{{.Name}}:
{{end}}
:END A:
`),
			templateManual: tplDef(`
:BEGIN A:{{ range .List }}
{{.Name}}:{{end}}
:END A:
`),
		},
		{
			name: "whitespaces get trimmed from both sides of control blocks",
			template: tplDef(`
:BEGIN A:
   {{ range .List }}   `+ /* Note extra 3 spaces */ `
{{.Name}}:
{{end}}
:END A:
`),
			templateManual: tplDef(`
:BEGIN A:{{ range .List }}`+ /* Note no extra spaces */`
{{.Name}}:{{end}}
:END A:
`),
		},
		{
			name: "nested expresstions get transformed correctly",
			template: tplDef(`
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
`),
			templateManual: tplDef(`
:BEGIN B:{{ range .List }}
:ITEM: {{ .Name }}{{ range .SubList }}{{ if (eq . "a") }}
  :SI:IF CLAUSE: {{ . }}{{ else if (eq . "b") }}
  :SI:ELSE IF CLAUSE: {{ . }}{{ else }}
  :SI:ELSE CLAUSE: {{ . }}{{ end }}{{ else }}
  :NOTHING:{{ end }}{{ end }}
:END B:
`),
		},
		{
			name: "chars on control block lines block whitespace trim",
			template: tplDef(`
:BEGIN A1:
a{{ range .List }}b
{{.Name}}:
{{else}}
NOTHING:
e{{end}}f
:END A1:
`),
			templateManual: tplDef(`
:BEGIN A1:
a{{ range .List }}b
{{.Name}}:
{{else}}
NOTHING:
e{{end}}f
:END A1:
`),
		},
	}
	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			tpl := template.New("(root)")
			tplOriginal, err := tpl.New("original").Parse(tt.template)
			require.NoError(t, err)

			adjustNode(nil, tplOriginal.Root, nil)

			tplOriginalOut := bytes.NewBuffer(nil)
			err = tplOriginal.Execute(tplOriginalOut, data)
			require.NoError(t, err)

			tplManual, err := tpl.New("manual").Parse(tt.templateManual)
			require.NoError(t, err)

			tplManualOut := bytes.NewBuffer(nil)
			err = tplManual.Execute(tplManualOut, data)
			require.NoError(t, err)

			require.Equal(t, tplManualOut.String(), tplOriginalOut.String())
		})
	}
}
