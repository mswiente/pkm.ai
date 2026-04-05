package skills

import "embed"

// FS contains all skill definition files embedded in the binary.
//
//go:embed *.md
var FS embed.FS
