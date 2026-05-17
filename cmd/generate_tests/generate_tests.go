// Command generate_tests converts the JavaScript tailwind-merge test files in
// ../tailwind-merge/tests/*.test.ts into a Go test file at
// pkg/twmerge/generated_test.go.
//
// It only translates tests of the simple form:
//
//	expect(twMerge('class string')).toBe('expected output')
//	expect(twMerge('a', 'b', 'c')).toBe('expected')
//	expect(twMerge(`...`)).toBe(`...`)
//
// Tests using `createTailwindMerge`, `extendTailwindMerge`, `twJoin`, custom
// configs, or extension hooks are skipped — those need hand-written Go tests
// because they depend on JS-specific features.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type testCase struct {
	in  string
	out string
}

type testGroup struct {
	source      string
	description string
	cases       []testCase
}

var (
	// Locate `test('name', ...)` or `test("name", ...)`.
	testNameRE = regexp.MustCompile(`(?m)^test\(['"](.+?)['"]\s*,`)
	// Match `expect(twMerge(<args>)).toBe(<expected>)` capturing the args inside.
	expectRE = regexp.MustCompile(`(?s)expect\(\s*twMerge\((.*?)\)\s*\)\.toBe\((.*?)\)`)
)

func main() {
	var inDir, outFile string
	flag.StringVar(&inDir, "in", "../tailwind-merge/tests", "directory containing .test.ts files")
	flag.StringVar(&outFile, "out", "pkg/twmerge/generated_test.go", "output Go test file path")
	flag.Parse()

	entries, err := os.ReadDir(inDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var groups []testGroup
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".test.ts") {
			continue
		}
		path := filepath.Join(inDir, e.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		src := string(data)
		// Skip whole files that need a custom config — we don't translate those.
		if containsAny(src, "createTailwindMerge", "extendTailwindMerge") {
			// Some of these files also have plain twMerge tests; allow only those
			// that don't reference the unsupported helpers.
			// We still try block-by-block extraction below; per-block skipping handles it.
		}
		fileGroups := parseFile(e.Name(), src)
		groups = append(groups, fileGroups...)
	}

	sort.SliceStable(groups, func(i, j int) bool {
		if groups[i].source != groups[j].source {
			return groups[i].source < groups[j].source
		}
		return groups[i].description < groups[j].description
	})

	out := renderGo(groups)
	if err := os.WriteFile(outFile, []byte(out), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("wrote %s with %d groups\n", outFile, len(groups))
}

func containsAny(src string, needles ...string) bool {
	for _, n := range needles {
		if strings.Contains(src, n) {
			return true
		}
	}
	return false
}

// parseFile extracts test blocks from a TS test file. We assume blocks start at
// "test(" at column 0 and end at "})" at column 0.
func parseFile(name, src string) []testGroup {
	var groups []testGroup
	lines := strings.Split(src, "\n")
	var i int
	for i < len(lines) {
		line := lines[i]
		if !strings.HasPrefix(line, "test(") {
			i++
			continue
		}
		// gather lines until matching `})` at column 0.
		start := i
		end := i
		depth := 0
		for j := i; j < len(lines); j++ {
			for _, c := range lines[j] {
				if c == '(' || c == '{' || c == '[' {
					depth++
				} else if c == ')' || c == '}' || c == ']' {
					depth--
				}
			}
			if depth <= 0 {
				end = j
				break
			}
		}
		block := strings.Join(lines[start:end+1], "\n")
		i = end + 1

		// skip blocks that use unsupported helpers.
		if containsAny(block, "createTailwindMerge", "extendTailwindMerge", "twJoin", "parseClassName") {
			continue
		}
		desc := ""
		if m := testNameRE.FindStringSubmatch(block); m != nil {
			desc = m[1]
		}
		cases := extractCases(block)
		if len(cases) == 0 {
			continue
		}
		groups = append(groups, testGroup{
			source:      name,
			description: desc,
			cases:       cases,
		})
	}
	return groups
}

func extractCases(block string) []testCase {
	var cases []testCase
	matches := expectRE.FindAllStringSubmatch(block, -1)
	for _, m := range matches {
		args := strings.TrimSpace(m[1])
		expected := strings.TrimSpace(m[2])
		in, ok := flattenArgs(args)
		if !ok {
			continue
		}
		out, ok := decodeStringLiteral(expected)
		if !ok {
			continue
		}
		cases = append(cases, testCase{in: in, out: out})
	}
	return cases
}

// flattenArgs turns the args of `twMerge(...)` into a single Go string
// argument. We support multi-arg calls by joining with a space — twMerge's
// own behaviour for plain string args is equivalent to "a b c".
// Returns ok=false if any argument is not a plain string literal.
func flattenArgs(args string) (string, bool) {
	parts := splitTopLevel(args, ',')
	var out []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		s, ok := decodeStringLiteral(p)
		if !ok {
			return "", false
		}
		if s != "" {
			out = append(out, s)
		}
	}
	return strings.Join(out, " "), true
}

func splitTopLevel(s string, sep rune) []string {
	var parts []string
	depth := 0
	start := 0
	var inStr rune
	for i, c := range s {
		switch c {
		case '\'', '"', '`':
			if inStr == 0 {
				inStr = c
			} else if inStr == c {
				inStr = 0
			}
		case '(', '[', '{':
			if inStr == 0 {
				depth++
			}
		case ')', ']', '}':
			if inStr == 0 {
				depth--
			}
		}
		if inStr == 0 && depth == 0 && c == sep {
			parts = append(parts, s[start:i])
			start = i + 1
		}
	}
	parts = append(parts, s[start:])
	return parts
}

// decodeStringLiteral parses a single JS string literal ('...', "...", or `...`)
// into its string value. Returns ok=false if the input is not a single literal
// (e.g. it's an expression, template with interpolations, etc.).
func decodeStringLiteral(s string) (string, bool) {
	s = strings.TrimSpace(s)
	// Drop trailing comma if any
	s = strings.TrimRight(s, ",")
	s = strings.TrimSpace(s)
	if len(s) < 2 {
		return "", false
	}
	q := s[0]
	if q != '\'' && q != '"' && q != '`' {
		return "", false
	}
	if s[len(s)-1] != q {
		return "", false
	}
	body := s[1 : len(s)-1]
	if q == '`' && strings.Contains(body, "${") {
		return "", false
	}
	// Unescape only the quote, backslash, and standard escapes.
	var b strings.Builder
	for i := 0; i < len(body); i++ {
		c := body[i]
		if c == '\\' && i+1 < len(body) {
			n := body[i+1]
			switch n {
			case 'n':
				b.WriteByte('\n')
			case 't':
				b.WriteByte('\t')
			case 'r':
				b.WriteByte('\r')
			case '\\':
				b.WriteByte('\\')
			case '\'', '"', '`':
				b.WriteByte(n)
			default:
				b.WriteByte('\\')
				b.WriteByte(n)
			}
			i++
			continue
		}
		// Inside `template literals` newlines are literal.
		if (c == '\n' || c == '\r') && q != '`' {
			return "", false
		}
		b.WriteByte(c)
	}
	return b.String(), true
}

func renderGo(groups []testGroup) string {
	var b strings.Builder
	b.WriteString(`// Code generated by cmd/generate_tests. DO NOT EDIT.
//
// Source: github.com/dcastil/tailwind-merge tests/*.test.ts

package twmerge

import "testing"

`)
	// Group tests by source file so failures point back to the JS source.
	bySource := map[string][]testGroup{}
	var sources []string
	for _, g := range groups {
		if _, ok := bySource[g.source]; !ok {
			sources = append(sources, g.source)
		}
		bySource[g.source] = append(bySource[g.source], g)
	}
	sort.Strings(sources)

	for _, src := range sources {
		fnName := "TestGenerated_" + sanitize(strings.TrimSuffix(src, ".test.ts"))
		fmt.Fprintf(&b, "// From %s\n", src)
		fmt.Fprintf(&b, "func %s(t *testing.T) {\n", fnName)
		b.WriteString("\ttests := []struct {\n\t\tdesc string\n\t\tin   string\n\t\tout  string\n\t}{\n")
		for _, g := range bySource[src] {
			for _, c := range g.cases {
				fmt.Fprintf(&b, "\t\t{desc: %q, in: %q, out: %q},\n", g.description, c.in, c.out)
			}
		}
		b.WriteString("\t}\n\tfor _, tc := range tests {\n\t\tt.Run(tc.desc+\"/\"+tc.in, func(t *testing.T) {\n\t\t\tgot := Merge(tc.in)\n\t\t\tif got != tc.out {\n\t\t\t\tt.Errorf(\"Merge(%q) = %q; want %q\", tc.in, got, tc.out)\n\t\t\t}\n\t\t})\n\t}\n}\n\n")
	}
	return b.String()
}

func sanitize(s string) string {
	var b strings.Builder
	upperNext := true
	for _, c := range s {
		switch {
		case c >= 'a' && c <= 'z':
			if upperNext {
				b.WriteRune(c - 32)
				upperNext = false
			} else {
				b.WriteRune(c)
			}
		case c >= 'A' && c <= 'Z':
			b.WriteRune(c)
			upperNext = false
		case c >= '0' && c <= '9':
			b.WriteRune(c)
			upperNext = false
		default:
			upperNext = true
		}
	}
	return b.String()
}
