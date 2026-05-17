package twmerge

// MergeConfigs returns base after applying the override/extend operations from ext.
// The base config is mutated in place.
func MergeConfigs(base *Config, ext *ConfigExtension) *Config {
	if ext == nil {
		return base
	}

	if ext.CacheSize != nil {
		base.CacheSize = *ext.CacheSize
	}
	if ext.Prefix != nil {
		base.Prefix = *ext.Prefix
	}
	if ext.ExperimentalParseClassName != nil {
		base.ExperimentalParseClassName = ext.ExperimentalParseClassName
	}

	if ov := ext.Override; ov != nil {
		overrideMap(base.Theme, ov.Theme)
		overrideMap(base.ClassGroups, ov.ClassGroups)
		overrideMap(base.ConflictingClassGroups, ov.ConflictingClassGroups)
		overrideMap(base.ConflictingClassGroupModifiers, ov.ConflictingClassGroupModifiers)
		if ov.PostfixLookupClassGroups != nil {
			base.PostfixLookupClassGroups = ov.PostfixLookupClassGroups
		}
		if ov.OrderSensitiveModifiers != nil {
			base.OrderSensitiveModifiers = ov.OrderSensitiveModifiers
		}
	}

	if ex := ext.Extend; ex != nil {
		mergeGroupMap(base.Theme, ex.Theme)
		mergeGroupMap(base.ClassGroups, ex.ClassGroups)
		mergeStringSliceMap(base.ConflictingClassGroups, ex.ConflictingClassGroups)
		mergeStringSliceMap(base.ConflictingClassGroupModifiers, ex.ConflictingClassGroupModifiers)
		if len(ex.PostfixLookupClassGroups) > 0 {
			base.PostfixLookupClassGroups = append(base.PostfixLookupClassGroups, ex.PostfixLookupClassGroups...)
		}
		if len(ex.OrderSensitiveModifiers) > 0 {
			base.OrderSensitiveModifiers = append(base.OrderSensitiveModifiers, ex.OrderSensitiveModifiers...)
		}
	}

	return base
}

func overrideMap[V any](base map[string]V, override map[string]V) {
	for k, v := range override {
		base[k] = v
	}
}

func mergeGroupMap(base map[string]ClassGroup, extend map[string]ClassGroup) {
	for k, v := range extend {
		if existing, ok := base[k]; ok {
			combined := make(ClassGroup, 0, len(existing)+len(v))
			combined = append(combined, existing...)
			combined = append(combined, v...)
			base[k] = combined
		} else {
			base[k] = v
		}
	}
}

func mergeStringSliceMap(base map[string][]string, extend map[string][]string) {
	for k, v := range extend {
		if existing, ok := base[k]; ok {
			combined := make([]string, 0, len(existing)+len(v))
			combined = append(combined, existing...)
			combined = append(combined, v...)
			base[k] = combined
		} else {
			base[k] = v
		}
	}
}
