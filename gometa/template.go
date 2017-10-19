package gometa

import (
	"html/template"
)

var tmpl = template.Must(template.New("index").Parse(`<!doctype html>
<html>
<head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
	<meta name="go-import" content="{{.ImportRoot}} {{.VCS}} {{.VCSRoot}}">
</head>
<body>{{.ImportRoot}}</body>
</html>
`))

type data struct {
	ImportRoot string
	VCS        string
	VCSRoot    string
	Suffix     string
}
