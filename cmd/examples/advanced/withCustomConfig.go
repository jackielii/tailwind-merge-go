package main

import (
	"fmt"

	"github.com/jackielii/tailwind-merge-go/pkg/lru"
	twmerge "github.com/jackielii/tailwind-merge-go/pkg/twmerge"
)

var twMerger twmerge.TwMergeFn

func main() {
	// Start from the default config and tweak it.
	config := twmerge.GetDefaultConfig()

	// Bring your own cache implementation if the default LRU doesn't fit.
	customCache := lru.Make(10000)

	// Build the merger.
	twMerger = twmerge.CreateTwMerge(config, twmerge.WithCache(customCache))

	fmt.Println(twMerger("px-4 px-10", "p-20")) // output: "p-20"
}
