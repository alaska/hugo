package helpers

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/sourcegraph/syntaxhighlight"
	jww "github.com/spf13/jwalterweatherman"
)

var tmplCodeTable, _ = template.New("Table Parser").Funcs(template.FuncMap{"add": func(a, b int) int { return a + b }}).Parse(`<table class='codetable' id='codetable-{{ .TableNumber }}'>
{{ range $idx, $line := .Lines}}<tr><td class='codetable-line-number'>{{ add $idx 1 }}</td><td id='codetable-{{ $.TableNumber }}-line-{{ add $idx 1 }}' class='codetable-line'>{{ $line }}</tr>
{{ end }}</table>`)

func CodeHighlight(code string, useTable string, tableNum int64) string {
	code = strings.Trim(code, "\n")
	highlighted, err := syntaxhighlight.AsHTML([]byte(code))
	if err != nil {
		jww.ERROR.Print(err.Error())
		return code
	}
	highlightedStr := string(highlighted)
	if useTable == "true" {
		splitCode := strings.Split(highlightedStr, "\n")
		for i, val := range splitCode {
			if val == "" {
				//replace empty lines with zero-width space to maintain height
				splitCode[i] = "&#8203;"
			}
		}
		var tableOut bytes.Buffer
		err = tmplCodeTable.Execute(&tableOut, struct {
			TableNumber int64
			Lines       []string
		}{
			tableNum,
			splitCode,
		})
		if err != nil {
			jww.ERROR.Print(err.Error())
			return code
		}
		return (tableOut.String())
	}
	return (string(highlighted))
}
