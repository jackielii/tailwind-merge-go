package twmerge

// configUtils bundles per-config helpers used by the merge function.
type configUtils struct {
	parseClassName              ParseClassNameFn
	sortModifiers               func([]string) []string
	getClassGroupID             func(className string) (string, bool)
	getConflictingClassGroupIDs func(classGroupID string, hasPostfixModifier bool) []string
	postfixLookupClassGroupIDs  map[string]bool
}

func createConfigUtils(config *Config) *configUtils {
	cg := createClassGroupUtils(config)
	return &configUtils{
		parseClassName:              createParseClassName(config),
		sortModifiers:               createSortModifiers(config),
		getClassGroupID:             cg.getClassGroupID,
		getConflictingClassGroupIDs: cg.getConflictingClassGroupIDs,
		postfixLookupClassGroupIDs:  createPostfixLookupClassGroupIDs(config),
	}
}

func createPostfixLookupClassGroupIDs(config *Config) map[string]bool {
	lookup := make(map[string]bool, len(config.PostfixLookupClassGroups))
	for _, id := range config.PostfixLookupClassGroups {
		lookup[id] = true
	}
	return lookup
}
