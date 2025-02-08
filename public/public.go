package public

import "embed"

//go:embed css js *.html
var Templates embed.FS
