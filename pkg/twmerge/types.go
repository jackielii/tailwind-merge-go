package twmerge

// ClassValidator is a function that tests whether a class part matches a group.
type ClassValidator func(classPart string) bool

// ThemeGetter retrieves a class group from the configured theme.
// Use FromTheme to construct one.
type ThemeGetter struct {
	key string
}

// Key returns the theme key the getter looks up.
func (g ThemeGetter) Key() string { return g.key }

// ClassGroup is a list of class definitions. Each item is one of:
//   - string: a literal class part (e.g. "auto", "")
//   - ClassValidator: a function that validates an arbitrary value
//   - ThemeGetter: a reference to a theme key
//   - ClassObject: a nested record producing prefixed class parts
type ClassGroup = []any

// ClassObject is a nested record mapping a class-part prefix to its sub-groups.
// e.g. {"object": ["contain", "cover", ...]} produces "object-contain", "object-cover".
type ClassObject = map[string]ClassGroup

// Theme maps theme group IDs to class groups.
type Theme = map[string]ClassGroup

// ClassGroupsMap maps class group IDs to class groups.
type ClassGroupsMap = map[string]ClassGroup

// ConflictingClassGroupsMap maps a class group ID to the IDs of class groups
// it conflicts with (i.e. removes when the class is present).
type ConflictingClassGroupsMap = map[string][]string

// ParsedClassName is the output of the class-name parser.
type ParsedClassName struct {
	// IsExternal indicates the class is not a Tailwind class and should pass through.
	IsExternal bool
	// Modifiers in the order they appear (e.g. ["hover", "dark"] for "hover:dark:bg-white").
	Modifiers []string
	// HasImportantModifier reports whether the class has a `!` important marker.
	HasImportantModifier bool
	// BaseClassName is the base class without preceding modifiers or important marker.
	BaseClassName string
	// MaybePostfixModifierPosition is the index of a possible postfix modifier ('/'),
	// or -1 if none was found.
	MaybePostfixModifierPosition int
}

// ParseClassNameFn parses a class name into its components.
type ParseClassNameFn func(className string) ParsedClassName

// ExperimentalParseClassNameParam is passed to ExperimentalParseClassName hooks.
type ExperimentalParseClassNameParam struct {
	ClassName      string
	ParseClassName ParseClassNameFn
}

// Config is the tailwind-merge configuration.
type Config struct {
	// CacheSize is the size of the LRU cache used for memoizing results.
	// No cache is used for values <= 0.
	CacheSize int
	// Prefix added to Tailwind-generated classes.
	Prefix string
	// ExperimentalParseClassName allows customising parsing of individual classes.
	// May introduce breaking changes in minor versions.
	ExperimentalParseClassName func(p ExperimentalParseClassNameParam) ParsedClassName

	// Theme scales used in ClassGroups via ThemeGetters.
	Theme Theme
	// ClassGroups maps class group IDs to their class definitions.
	ClassGroups ClassGroupsMap
	// ClassGroupOrder optionally lists ClassGroups keys in the order they
	// should be processed when building the lookup trie. This matters when
	// two groups share a trie node and their validators overlap — the order
	// determines which group wins on a tie (JS preserves source insertion
	// order; Go maps don't, so the order is provided explicitly).
	// If empty, ClassGroups is iterated in map order (non-deterministic).
	ClassGroupOrder []string
	// ConflictingClassGroups maps a class group ID to the IDs of conflicting groups.
	// If a class from the key ID is present, preceding classes from the values are removed.
	ConflictingClassGroups ConflictingClassGroupsMap
	// ConflictingClassGroupModifiers maps a class group ID to IDs of class groups
	// conflicting only when a postfix modifier is present (e.g. text-lg/loose vs leading-*).
	ConflictingClassGroupModifiers ConflictingClassGroupsMap
	// PostfixLookupClassGroups lists class group IDs that should be resolved again
	// with their postfix modifier attached.
	PostfixLookupClassGroups []string
	// OrderSensitiveModifiers lists modifiers whose order among other modifiers
	// changes the target element and so must be preserved.
	OrderSensitiveModifiers []string
}

// ConfigExtension is a partial config used to extend or override a base Config.
type ConfigExtension struct {
	CacheSize                  *int
	Prefix                     *string
	ExperimentalParseClassName func(p ExperimentalParseClassNameParam) ParsedClassName

	Override *ConfigGroups
	Extend   *ConfigGroups
}

// ConfigGroups is the dynamic part of a Config used in ConfigExtension.
type ConfigGroups struct {
	Theme                          Theme
	ClassGroups                    ClassGroupsMap
	ConflictingClassGroups         ConflictingClassGroupsMap
	ConflictingClassGroupModifiers ConflictingClassGroupsMap
	PostfixLookupClassGroups       []string
	OrderSensitiveModifiers        []string
}
