// Package resources provides access to embedded data files
package resources

import "embed"

//go:embed data/placeholder.json data/fedramp-high.json data/fedramp-moderate.json
var Data embed.FS
