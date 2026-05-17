package twmerge

// FromTheme returns a ThemeGetter referencing a theme key. When used inside a
// ClassGroup definition, the class-map builder substitutes the contents of
// `theme[key]` at that position.
func FromTheme(key string) ThemeGetter {
	return ThemeGetter{key: key}
}
