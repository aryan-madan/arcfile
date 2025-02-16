package public

import "embed"

//go:embed css js *.html img
var Templates embed.FS
