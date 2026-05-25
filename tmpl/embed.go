// Package tmpl contains embedded email templates.
package tmpl

import "embed"

// FS holds the embedded template files.
//
//go:embed *.tmpl
var FS embed.FS
