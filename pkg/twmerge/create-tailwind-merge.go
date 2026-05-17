package twmerge

import (
	"sync"

	"github.com/jackielii/tailwind-merge-go/pkg/cache"
	"github.com/jackielii/tailwind-merge-go/pkg/lru"
)

// CreateConfigFn returns a base Config. Used by CreateTwMerge / ExtendTailwindMerge.
type CreateConfigFn func() *Config

// CreateConfigSubsequentFn transforms an existing Config (e.g. to apply an extension).
type CreateConfigSubsequentFn func(*Config) *Config

// TwMergeFn merges Tailwind CSS class strings, removing later-conflicting classes.
type TwMergeFn func(args ...ClassNameValue) string

// CreateTailwindMerge builds a Tailwind class merger from a sequence of config-creators.
// Initialisation is deferred until the merger is first invoked (matching the JS version).
//
// You can pass a custom cache via WithCache.
func CreateTailwindMerge(first CreateConfigFn, rest ...CreateConfigSubsequentFn) TwMergeFn {
	return createTailwindMerge(first, rest, nil)
}

// CreateTwMergeOption configures CreateTwMerge.
type CreateTwMergeOption func(*createTwMergeOptions)

type createTwMergeOptions struct {
	cache cache.ICache
}

// WithCache supplies a custom cache implementation. Pass nil to disable caching
// (the default LRU cache from pkg/lru is used otherwise).
func WithCache(c cache.ICache) CreateTwMergeOption {
	return func(o *createTwMergeOptions) { o.cache = c }
}

// CreateTwMerge is a convenience constructor that mirrors the previous Go API.
// Pass options to override the default LRU cache. `config` may be nil to use the
// default config.
func CreateTwMerge(config *Config, opts ...CreateTwMergeOption) TwMergeFn {
	o := createTwMergeOptions{}
	for _, opt := range opts {
		opt(&o)
	}
	createConfig := func() *Config {
		if config != nil {
			return config
		}
		return GetDefaultConfig()
	}
	return createTailwindMerge(createConfig, nil, o.cache)
}

func createTailwindMerge(first CreateConfigFn, rest []CreateConfigSubsequentFn, customCache cache.ICache) TwMergeFn {
	var (
		initOnce       sync.Once
		utils          *configUtils
		c              cache.ICache
		mergeClassList mergeClassListFn
	)

	doInit := func() {
		config := first()
		for _, fn := range rest {
			config = fn(config)
		}
		utils = createConfigUtils(config)
		mergeClassList = createMergeClassList(utils)

		if customCache != nil {
			c = customCache
		} else if config.CacheSize > 0 {
			c = lru.Make(config.CacheSize)
		} else {
			c = noopCache{}
		}
	}

	return func(args ...ClassNameValue) string {
		classList := TwJoin(args...)
		if classList == "" {
			return ""
		}
		initOnce.Do(doInit)
		if cached := c.Get(classList); cached != "" {
			return cached
		}
		result := mergeClassList(classList)
		c.Set(classList, result)
		return result
	}
}

// ExtendTailwindMerge returns a merger that uses the default config extended by ext
// followed by additional config-creators.
func ExtendTailwindMerge(ext *ConfigExtension, rest ...CreateConfigSubsequentFn) TwMergeFn {
	first := func() *Config {
		return MergeConfigs(GetDefaultConfig(), ext)
	}
	return createTailwindMerge(first, rest, nil)
}

// noopCache is used when caching is disabled.
type noopCache struct{}

func (noopCache) Get(string) string { return "" }
func (noopCache) Set(string, string) {}

// Merge is the default top-level merge function using the default config.
var Merge TwMergeFn = CreateTwMerge(nil)
