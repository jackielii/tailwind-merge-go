package twmerge

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// loadBenchmarkData reads tests/tw-merge-benchmark-data.json from the sibling
// tailwind-merge repo. Returns nil and skips the benchmark when the file is
// absent so the suite is portable.
func loadBenchmarkData(b *testing.B) [][]ClassNameValue {
	path := filepath.Join("..", "..", "..", "tailwind-merge", "tests", "tw-merge-benchmark-data.json")
	data, err := os.ReadFile(path)
	if err != nil {
		b.Skipf("benchmark data not found at %s (%v) — clone github.com/dcastil/tailwind-merge next to this repo to enable", path, err)
		return nil
	}
	var raw [][]any
	if err := json.Unmarshal(data, &raw); err != nil {
		b.Fatalf("parse benchmark data: %v", err)
	}
	out := make([][]ClassNameValue, len(raw))
	for i, row := range raw {
		args := make([]ClassNameValue, len(row))
		copy(args, row)
		out[i] = args
	}
	return out
}

func BenchmarkInit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m := CreateTwMerge(nil)
		_ = m("flex")
	}
}

func BenchmarkSimple(b *testing.B) {
	m := CreateTwMerge(nil)
	// Prime cache miss by running once outside the timing loop.
	_ = m("flex mx-10 px-10", "mr-5 pr-5")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m("flex mx-10 px-10", "mr-5 pr-5")
	}
}

func BenchmarkHeavy(b *testing.B) {
	m := CreateTwMerge(nil)
	args := []ClassNameValue{
		"font-medium text-sm leading-16",
		"group/button relative isolate items-center justify-center overflow-hidden rounded-md outline-none transition [-webkit-app-region:no-drag] focus-visible:ring focus-visible:ring-primary",
		"inline-flex",
		"bg-primary-50 ring ring-primary-200",
		"text-primary dark:text-primary-900 hover:bg-primary-100",
		false,
		"font-medium text-sm leading-16 gap-4 px-6 py-4",
		nil,
		"p-0 size-24",
		nil,
	}
	_ = m(args...) // warm cache
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m(args...)
	}
}

func BenchmarkCollectionWithCache(b *testing.B) {
	data := loadBenchmarkData(b)
	if data == nil {
		return
	}
	m := CreateTwMerge(nil)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, row := range data {
			_ = m(row...)
		}
	}
}

func BenchmarkCollectionWithoutCache(b *testing.B) {
	data := loadBenchmarkData(b)
	if data == nil {
		return
	}
	cfg := GetDefaultConfig()
	cfg.CacheSize = 0
	m := CreateTwMerge(cfg)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, row := range data {
			_ = m(row...)
		}
	}
}

// ultraLongClassList constructs the same 200-iteration synthetic class list the
// JS benchmark uses (padding/margin/width/height conflicts + occasional modifiers).
func ultraLongClassList() []ClassNameValue {
	out := make([]ClassNameValue, 0, 2200)
	for i := 0; i < 200; i++ {
		out = append(out,
			fmt.Sprintf("p-%d", i%20),
			fmt.Sprintf("px-%d", i%20),
			fmt.Sprintf("py-%d", i%20),
			fmt.Sprintf("m-%d", i%20),
			fmt.Sprintf("mx-%d", i%20),
			fmt.Sprintf("my-%d", i%20),
			fmt.Sprintf("w-%d", i%20),
			fmt.Sprintf("h-%d", i%20),
			fmt.Sprintf("text-%d", i%10),
			fmt.Sprintf("bg-%d", i%10),
		)
		if i%10 == 0 {
			out = append(out,
				fmt.Sprintf("hover:p-%d", i%20),
				fmt.Sprintf("focus:m-%d", i%20),
			)
		}
	}
	return out
}

func BenchmarkUltraLongWithoutCache(b *testing.B) {
	cfg := GetDefaultConfig()
	cfg.CacheSize = 0
	m := CreateTwMerge(cfg)
	args := ultraLongClassList()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m(args...)
	}
}

func BenchmarkUltraLongWithCache(b *testing.B) {
	m := CreateTwMerge(nil)
	args := ultraLongClassList()
	_ = m(args...) // warm cache
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m(args...)
	}
}

func BenchmarkCacheHit(b *testing.B) {
	// Best case: a previously-merged input should be a one-hop cache lookup.
	m := CreateTwMerge(nil)
	const in = "flex mx-10 px-10 mr-5 pr-5 hover:bg-red-500"
	_ = m(in)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m(in)
	}
}

// BenchmarkMemoryFootprint reports the heap allocation increase from building
// a merger and running the heavy test once. Use `go test -bench Memory -benchmem`.
func BenchmarkMemoryFootprint(b *testing.B) {
	runtime.GC()
	var before runtime.MemStats
	runtime.ReadMemStats(&before)
	for i := 0; i < b.N; i++ {
		m := CreateTwMerge(nil)
		_ = m("flex mx-10 px-10", "mr-5 pr-5")
	}
	runtime.GC()
	var after runtime.MemStats
	runtime.ReadMemStats(&after)
	b.ReportMetric(float64(after.HeapAlloc-before.HeapAlloc)/float64(b.N), "heap-bytes/op-delta")
}
