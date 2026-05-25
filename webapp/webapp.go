// Package webapp contains the embedded SPA frontend.
package webapp

import "embed"

// FS holds the embedded webapp build output.
//
//go:embed build
var FS embed.FS
