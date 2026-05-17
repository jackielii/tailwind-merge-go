<div align="center">
    <br />
    <a href="https://github.com/jackielii/tailwind-merge-go">
        <img src="https://raw.githubusercontent.com/jackielii/tailwind-merge-go/master/assets/logo.svg" alt="tailwind-merge-go" height="150px" />
    </a>
</div>

# tailwind-merge-go - Tailwind Merge For Golang

<a href="https://pkg.go.dev/github.com/jackielii/tailwind-merge-go"><img src="https://pkg.go.dev/badge/github.com//github.com/jackielii/tailwind-merge-go.svg" alt="Go Reference" /></a>
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/jackielii/tailwind-merge-go)](https://goreportcard.com/report/github.com/jackielii/tailwind-merge-go)
[![Coverage Status](https://coveralls.io/repos/github/jackielii/tailwind-merge-go/badge.svg?branch=master)](https://coveralls.io/github/jackielii/tailwind-merge-go?branch=master)

Utility function to efficiently merge Tailwind CSS classes in Golang without style conflicts. This library aims to be as close as possible to a 1:1 copy of the original [dcastil/tailwind-merge](https://github.com/dcastil/tailwind-merge/) library written in javascript.

```go
import (
	"fmt"

	twmerge "github.com/jackielii/tailwind-merge-go"
)

func main() {

	// example usage
	c := twmerge.Merge("px-4 px-10 p-1")
	fmt.Println(c) // "p-1"
}
```

- Supports Tailwind CSS v4 (and the v3 syntax kept around for compatibility — leading `!` important, etc.)
- Mirrors the JS `tailwind-merge` v3.6.0 architecture (theme groups, `fromTheme` getters, postfix-modifier conflicts, order-sensitive modifiers, postfix lookup groups, prefix-as-external classes)
- Support for extending the default configuration
- Support for providing your own caching solution
- Its in 0.1.0, can I use it? Sure! I will personally be deploying this to prod. It's only in pre 1.0 because there some extra features I want to add before the 1.0 release (see roadmap)

## [Why use it?](https://github.com/dcastil/tailwind-merge/blob/v2.2.1/docs/what-is-it-for.md)

- See [tailwind-merge](https://github.com/dcastil/tailwind-merge/blob/v2.2.1/docs/what-is-it-for.md)
- Or Watch this amazing video on it

[Watch this introduction video from Simon Vrachliotis (@simonswiss) ↓ ![The "why" behind tailwind-merge](https://img.youtube.com/vi/tfgLd5ZSNPc/maxresdefault.jpg)](https://www.youtube.com/watch?v=tfgLd5ZSNPc (Watch YouTube video "Tailwind-Merge Is Incredibly Useful — And Here's Why!"))

## [Limitations](https://github.com/dcastil/tailwind-merge/blob/v2.2.1/docs/limitations.md)

- See [tailwind-merge](https://github.com/dcastil/tailwind-merge/blob/v2.2.1/docs/limitations.md)

## Advanced Examples

You might also want to check out the advanced example at `/cmd/examples/advanced`

### Provide Your Own or Extend Default Config

```go
import (
	// Use the pkg/twmerge sub-package directly to access the types, helpers, and validators.
	twmerge "github.com/jackielii/tailwind-merge-go/pkg/twmerge"
)

var twMerger twmerge.TwMergeFn

func main() {
	// Start from the default config and mutate it freely; GetDefaultConfig returns a fresh copy.
	config := twmerge.GetDefaultConfig()

	// e.g. add a custom class group
	config.ClassGroups["custom-shadow"] = twmerge.ClassGroup{"custom-shadow-sm", "custom-shadow-md"}

	twMerger = twmerge.CreateTwMerge(config)
	fmt.Println(twMerger("px-4 px-10", "p-20")) // "p-20"
}
```

Or use `ExtendTailwindMerge` to declaratively extend the default:

```go
twMerger := twmerge.ExtendTailwindMerge(&twmerge.ConfigExtension{
	Extend: &twmerge.ConfigGroups{
		ClassGroups: twmerge.ClassGroupsMap{
			"custom-shadow": twmerge.ClassGroup{"custom-shadow-sm", "custom-shadow-md"},
		},
	},
})
```

### Provide your own Cache

The default cache is an LRU. To plug in your own, satisfy the `cache.ICache` interface and pass it via `WithCache`:

```go
type ICache interface {
	Get(string) string
	Set(string, string)
}
```

```go
import (
	twmerge "github.com/jackielii/tailwind-merge-go/pkg/twmerge"
	"github.com/jackielii/tailwind-merge-go/pkg/lru"
)

func main() {
	twMerger := twmerge.CreateTwMerge(nil, twmerge.WithCache(lru.Make(10000)))
	fmt.Println(twMerger("px-4 px-10", "p-20")) // "p-20"
}
```

## Contributing

Checkout the [contributing docs](./CONTRIBUTING.md)

## Roadmap

- Improve cache concurrent performance by locking on a per key basis -> https://github.com/EagleChen/mapmutex
- Build the class map on initialization and have a simple config style
- replace regex with more performant solution
- Move arbitrary value delimeters '[' & ']' to config somehow?
- Plugins & easy plugin api.

## Acknowledgments

- Credit for all the hard work goes to [dcastil/tailwind-merge](https://github.com/dcastil/tailwind-merge/).
  - For the tests I used
  - For the approach and the code. I mostly translated from js to go
  - For the logo
- Big thank you to [tylantz/go-tailwind-merge/](https://github.com/tylantz/go-tailwind-merge/tree/main) for pushing me to finally do this by writing a very interesting version of this same idea (I encourage you to check it out) and for the code to generate a go test file based on tailwind-merge's tests
