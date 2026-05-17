// Package twmerge is a top-level convenience that re-exports the default merger
// from pkg/twmerge.
//
// For custom configurations, caches, or extension, import
// github.com/jackielii/tailwind-merge-go/pkg/twmerge directly.
package twmerge

import (
	pkg "github.com/jackielii/tailwind-merge-go/pkg/twmerge"
)

// Merge merges Tailwind CSS class strings using the default configuration.
func Merge(args ...pkg.ClassNameValue) string {
	return pkg.Merge(args...)
}
