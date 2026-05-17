package twmerge

import "strings"

const classPartSeparator = "-"
const arbitraryPropertyPrefix = "arbitrary.."

type classValidatorObject struct {
	classGroupID string
	validator    ClassValidator
	// catchall is true for validators that accept arbitrary inputs (IsAny,
	// IsAnyNonArbitrary, etc.). When several class groups share a trie node,
	// non-catchall validators must be tested first so specific ones (IsNumber,
	// IsArbitraryLength, …) win over generic ones — JS preserves this via
	// object insertion order, which Go maps do not.
	catchall bool
}

// classPartObject is a node in the class-name trie.
type classPartObject struct {
	nextPart     map[string]*classPartObject
	validators   []classValidatorObject
	classGroupID string
}

func newClassPartObject() *classPartObject {
	return &classPartObject{nextPart: map[string]*classPartObject{}}
}

// validatorCanary is an input no real Tailwind class would have — used to
// probe whether a validator is a catchall (returns true for anything).
const validatorCanary = "\x00twmerge-catchall-probe\x00"

func isCatchallValidator(v ClassValidator) bool {
	// Recover from panics in case a validator assumes non-empty input.
	defer func() { _ = recover() }()
	return v(validatorCanary)
}

// classGroupUtils bundles class-group lookup helpers.
type classGroupUtils struct {
	getClassGroupID            func(className string) (string, bool)
	getConflictingClassGroupIDs func(classGroupID string, hasPostfixModifier bool) []string
}

func createClassGroupUtils(config *Config) classGroupUtils {
	classMap := createClassMap(config)
	conflictingClassGroups := config.ConflictingClassGroups
	conflictingClassGroupModifiers := config.ConflictingClassGroupModifiers

	getClassGroupID := func(className string) (string, bool) {
		if len(className) > 1 && className[0] == '[' && className[len(className)-1] == ']' {
			id, ok := getGroupIDForArbitraryProperty(className)
			return id, ok
		}

		classParts := strings.Split(className, classPartSeparator)
		// Classes like "-inset-1" produce an empty string as first part.
		startIndex := 0
		if len(classParts) > 1 && classParts[0] == "" {
			startIndex = 1
		}
		return getGroupRecursive(classParts, startIndex, classMap)
	}

	getConflictingClassGroupIDs := func(classGroupID string, hasPostfixModifier bool) []string {
		if hasPostfixModifier {
			modifierConflicts := conflictingClassGroupModifiers[classGroupID]
			baseConflicts := conflictingClassGroups[classGroupID]
			switch {
			case modifierConflicts != nil && baseConflicts != nil:
				combined := make([]string, 0, len(baseConflicts)+len(modifierConflicts))
				combined = append(combined, baseConflicts...)
				combined = append(combined, modifierConflicts...)
				return combined
			case modifierConflicts != nil:
				return modifierConflicts
			default:
				return baseConflicts
			}
		}
		return conflictingClassGroups[classGroupID]
	}

	return classGroupUtils{
		getClassGroupID:             getClassGroupID,
		getConflictingClassGroupIDs: getConflictingClassGroupIDs,
	}
}

func getGroupRecursive(classParts []string, startIndex int, node *classPartObject) (string, bool) {
	if startIndex >= len(classParts) {
		if node.classGroupID != "" {
			return node.classGroupID, true
		}
		return "", false
	}

	currentPart := classParts[startIndex]
	if next, ok := node.nextPart[currentPart]; ok {
		if id, found := getGroupRecursive(classParts, startIndex+1, next); found {
			return id, true
		}
	}

	if len(node.validators) == 0 {
		return "", false
	}

	var classRest string
	if startIndex == 0 {
		classRest = strings.Join(classParts, classPartSeparator)
	} else {
		classRest = strings.Join(classParts[startIndex:], classPartSeparator)
	}
	// Non-catchall validators tested first; see classValidatorObject.catchall.
	for _, v := range node.validators {
		if v.catchall {
			continue
		}
		if v.validator(classRest) {
			return v.classGroupID, true
		}
	}
	for _, v := range node.validators {
		if !v.catchall {
			continue
		}
		if v.validator(classRest) {
			return v.classGroupID, true
		}
	}
	return "", false
}

func getGroupIDForArbitraryProperty(className string) (string, bool) {
	inner := className[1 : len(className)-1]
	colon := strings.IndexByte(inner, ':')
	if colon == -1 {
		return "", false
	}
	property := inner[:colon]
	if property == "" {
		return "", false
	}
	return arbitraryPropertyPrefix + property, true
}

func createClassMap(config *Config) *classPartObject {
	root := newClassPartObject()
	// If an explicit order is provided, follow it so that validators in
	// shared trie nodes resolve in the JS source order.
	if len(config.ClassGroupOrder) > 0 {
		seen := make(map[string]struct{}, len(config.ClassGroups))
		for _, id := range config.ClassGroupOrder {
			group, ok := config.ClassGroups[id]
			if !ok {
				continue
			}
			seen[id] = struct{}{}
			processClassesRecursively(group, root, id, config.Theme)
		}
		// Any class groups not mentioned in the order get processed afterwards
		// (e.g. user-added groups via ExtendTailwindMerge that didn't update
		// ClassGroupOrder). Their relative order is non-deterministic but they
		// come strictly after the ordered ones.
		for id, group := range config.ClassGroups {
			if _, done := seen[id]; done {
				continue
			}
			processClassesRecursively(group, root, id, config.Theme)
		}
		return root
	}
	for id, group := range config.ClassGroups {
		processClassesRecursively(group, root, id, config.Theme)
	}
	return root
}

func processClassesRecursively(group ClassGroup, node *classPartObject, classGroupID string, theme Theme) {
	for _, def := range group {
		processClassDefinition(def, node, classGroupID, theme)
	}
}

func processClassDefinition(def any, node *classPartObject, classGroupID string, theme Theme) {
	switch d := def.(type) {
	case string:
		if d == "" {
			node.classGroupID = classGroupID
			return
		}
		target := getPart(node, d)
		target.classGroupID = classGroupID
	case ThemeGetter:
		sub := theme[d.key]
		processClassesRecursively(sub, node, classGroupID, theme)
	case ClassValidator:
		node.validators = append(node.validators, classValidatorObject{
			classGroupID: classGroupID,
			validator:    d,
			catchall:     isCatchallValidator(d),
		})
	case func(string) bool:
		// allow plain function literals to be used as validators
		v := ClassValidator(d)
		node.validators = append(node.validators, classValidatorObject{
			classGroupID: classGroupID,
			validator:    v,
			catchall:     isCatchallValidator(v),
		})
	case map[string]ClassGroup:
		// ClassObject (it's an alias for map[string][]any).
		for prefix, subGroup := range d {
			processClassesRecursively(subGroup, getPart(node, prefix), classGroupID, theme)
		}
	case map[string]any:
		// permissive form: values may also be []any
		for prefix, subGroup := range d {
			if g, ok := subGroup.(ClassGroup); ok {
				processClassesRecursively(g, getPart(node, prefix), classGroupID, theme)
			} else if g, ok := subGroup.([]any); ok {
				processClassesRecursively(g, getPart(node, prefix), classGroupID, theme)
			}
		}
	default:
		// unsupported definition type — ignore so the rest of the config still loads
	}
}

func getPart(node *classPartObject, path string) *classPartObject {
	current := node
	for _, p := range strings.Split(path, classPartSeparator) {
		next, ok := current.nextPart[p]
		if !ok {
			next = newClassPartObject()
			current.nextPart[p] = next
		}
		current = next
	}
	return current
}
