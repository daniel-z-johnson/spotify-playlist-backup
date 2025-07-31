package templates

import (
	"embed"
)

//go:embed *.gohtml
var TemplatesFS embed.FS
