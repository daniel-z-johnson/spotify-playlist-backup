package templates

import "embed"

//go:embed *.gohtml
var TemplatesFS embed.FS

//go:embed js/*.js
var JSFS embed.FS
