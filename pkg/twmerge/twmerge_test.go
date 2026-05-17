package twmerge

import "testing"

func TestPostfixModifierConflicts(t *testing.T) {
	cacheSize := 10
	customMerge := CreateTwMerge(&Config{
		CacheSize: cacheSize,
		Theme:     Theme{},
		ClassGroups: ClassGroupsMap{
			"foo": ClassGroup{"foo-1/2", "foo-2/3"},
			"bar": ClassGroup{"bar-1", "bar-2"},
			"baz": ClassGroup{"baz-1", "baz-2"},
		},
		ConflictingClassGroups: ConflictingClassGroupsMap{},
		ConflictingClassGroupModifiers: ConflictingClassGroupsMap{
			"baz": {"bar"},
		},
		OrderSensitiveModifiers: []string{},
	})

	cases := []struct {
		in, out string
	}{
		{"foo-1/2 foo-2/3", "foo-2/3"},
		{"bar-1 bar-2", "bar-2"},
		{"bar-1 baz-1", "bar-1 baz-1"},
		{"bar-1/2 bar-2", "bar-2"},
		{"bar-2 bar-1/2", "bar-1/2"},
		{"bar-1 baz-1/2", "baz-1/2"},
	}
	for _, c := range cases {
		if got := customMerge(c.in); got != c.out {
			t.Errorf("Merge(%q) = %q; want %q", c.in, got, c.out)
		}
	}
}

func TestPostfixLookupClassGroups(t *testing.T) {
	customMerge := CreateTwMerge(&Config{
		CacheSize: 10,
		Theme:     Theme{},
		ClassGroups: ClassGroupsMap{
			"base": ClassGroup{ClassObject{"foo": ClassGroup{"bar"}}},
			"named": ClassGroup{ClassValidator(func(v string) bool {
				return v == "foo-bar/baz"
			})},
		},
		ConflictingClassGroups: ConflictingClassGroupsMap{
			"named": {"base"},
		},
		ConflictingClassGroupModifiers: ConflictingClassGroupsMap{},
		PostfixLookupClassGroups:       []string{"base"},
		OrderSensitiveModifiers:        []string{},
	})

	cases := []struct {
		in, out string
	}{
		{"foo-bar foo-bar/baz", "foo-bar/baz"},
		{"foo-bar/baz foo-bar", "foo-bar/baz foo-bar"},
	}
	for _, c := range cases {
		if got := customMerge(c.in); got != c.out {
			t.Errorf("Merge(%q) = %q; want %q", c.in, got, c.out)
		}
	}
}

func TestOrderSensitiveModifiers(t *testing.T) {
	customMerge := CreateTwMerge(&Config{
		CacheSize: 10,
		Theme:     Theme{},
		ClassGroups: ClassGroupsMap{
			"foo": ClassGroup{"foo-1", "foo-2"},
		},
		ConflictingClassGroups:         ConflictingClassGroupsMap{},
		ConflictingClassGroupModifiers: ConflictingClassGroupsMap{},
		OrderSensitiveModifiers:        []string{"a", "b"},
	})

	cases := []struct {
		in, out string
	}{
		{"c:d:e:foo-1 d:c:e:foo-2", "d:c:e:foo-2"},
		{"a:b:foo-1 a:b:foo-2", "a:b:foo-2"},
		{"a:b:foo-1 b:a:foo-2", "a:b:foo-1 b:a:foo-2"},
		{"x:y:a:z:foo-1 y:x:a:z:foo-2", "y:x:a:z:foo-2"},
	}
	for _, c := range cases {
		if got := customMerge(c.in); got != c.out {
			t.Errorf("Merge(%q) = %q; want %q", c.in, got, c.out)
		}
	}
}

func TestPrefix(t *testing.T) {
	prefix := "tw"
	customMerge := CreateTwMerge(&Config{
		CacheSize: 10,
		Prefix:    prefix,
		Theme:     Theme{},
		ClassGroups: ClassGroupsMap{
			"display": ClassGroup{"block", "inline", "flex"},
		},
		ConflictingClassGroups:         ConflictingClassGroupsMap{},
		ConflictingClassGroupModifiers: ConflictingClassGroupsMap{},
		OrderSensitiveModifiers:        []string{},
	})

	// Prefixed Tailwind classes should be merged.
	if got := customMerge("tw:block tw:inline"); got != "tw:inline" {
		t.Errorf("prefixed merge: got %q", got)
	}
	// Non-prefixed classes are external — pass through, no merging.
	if got := customMerge("block inline"); got != "block inline" {
		t.Errorf("external pass-through: got %q", got)
	}
}

func TestExtendTailwindMerge(t *testing.T) {
	custom := ExtendTailwindMerge(&ConfigExtension{
		Extend: &ConfigGroups{
			ClassGroups: ClassGroupsMap{
				"custom-shadow": ClassGroup{"custom-shadow-sm", "custom-shadow-md"},
			},
		},
	})
	if got := custom("custom-shadow-sm custom-shadow-md"); got != "custom-shadow-md" {
		t.Errorf("got %q", got)
	}
	// Default behaviour still works.
	if got := custom("px-4 px-2"); got != "px-2" {
		t.Errorf("default still works: got %q", got)
	}
}

func TestExternalClassesPassThrough(t *testing.T) {
	cases := []struct {
		in, out string
	}{
		{"", ""},
		{"   ", ""},
		{"some-non-tw-class", "some-non-tw-class"},
		{"px-4 some-non-tw-class px-8", "some-non-tw-class px-8"},
	}
	for _, c := range cases {
		if got := Merge(c.in); got != c.out {
			t.Errorf("Merge(%q) = %q; want %q", c.in, got, c.out)
		}
	}
}

func TestTwJoin(t *testing.T) {
	cases := []struct {
		args []ClassNameValue
		out  string
	}{
		{[]ClassNameValue{"a", "b", "c"}, "a b c"},
		{[]ClassNameValue{"a", nil, "b"}, "a b"},
		{[]ClassNameValue{"a", false, "b"}, "a b"},
		{[]ClassNameValue{"a", []ClassNameValue{"b", []ClassNameValue{"c"}}}, "a b c"},
		{[]ClassNameValue{[]string{"a", "b"}, "c"}, "a b c"},
	}
	for _, c := range cases {
		if got := TwJoin(c.args...); got != c.out {
			t.Errorf("TwJoin(%v) = %q; want %q", c.args, got, c.out)
		}
	}
}
