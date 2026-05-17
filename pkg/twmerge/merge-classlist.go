package twmerge

import (
	"regexp"
	"strings"
)

var splitClassesRegex = regexp.MustCompile(`\s+`)

type mergeClassListFn func(classList string) string

func createMergeClassList(utils *configUtils) mergeClassListFn {
	return func(classList string) string {
		trimmed := strings.TrimSpace(classList)
		if trimmed == "" {
			return ""
		}
		classNames := splitClassesRegex.Split(trimmed, -1)
		// classGroupsInConflict acts as a set; preserves insertion order.
		conflictSet := make(map[string]struct{}, len(classNames))
		// resultParts collected from the right so the original order is restored at the end.
		var resultParts []string

		for i := len(classNames) - 1; i >= 0; i-- {
			originalClassName := classNames[i]
			parsed := utils.parseClassName(originalClassName)

			if parsed.IsExternal {
				resultParts = append(resultParts, originalClassName)
				continue
			}

			hasPostfixModifier := parsed.MaybePostfixModifierPosition != -1
			var classGroupID string
			var found bool

			if hasPostfixModifier {
				baseWithoutPostfix := parsed.BaseClassName[:parsed.MaybePostfixModifierPosition]
				classGroupID, found = utils.getClassGroupID(baseWithoutPostfix)

				if found && utils.postfixLookupClassGroupIDs[classGroupID] {
					if newID, ok := utils.getClassGroupID(parsed.BaseClassName); ok && newID != classGroupID {
						classGroupID = newID
						hasPostfixModifier = false
					}
				}
			} else {
				classGroupID, found = utils.getClassGroupID(parsed.BaseClassName)
			}

			if !found {
				if !hasPostfixModifier {
					resultParts = append(resultParts, originalClassName)
					continue
				}
				classGroupID, found = utils.getClassGroupID(parsed.BaseClassName)
				if !found {
					resultParts = append(resultParts, originalClassName)
					continue
				}
				hasPostfixModifier = false
			}

			// Build modifier prefix.
			var variantModifier string
			switch len(parsed.Modifiers) {
			case 0:
				variantModifier = ""
			case 1:
				variantModifier = parsed.Modifiers[0]
			default:
				variantModifier = strings.Join(utils.sortModifiers(parsed.Modifiers), ":")
			}

			modifierID := variantModifier
			if parsed.HasImportantModifier {
				modifierID += importantModifier
			}

			classID := modifierID + classGroupID
			if _, exists := conflictSet[classID]; exists {
				continue
			}
			conflictSet[classID] = struct{}{}

			for _, group := range utils.getConflictingClassGroupIDs(classGroupID, hasPostfixModifier) {
				conflictSet[modifierID+group] = struct{}{}
			}

			resultParts = append(resultParts, originalClassName)
		}

		// resultParts is in reverse insertion order; reverse to restore original-left-to-right.
		var b strings.Builder
		for i := len(resultParts) - 1; i >= 0; i-- {
			if b.Len() > 0 {
				b.WriteByte(' ')
			}
			b.WriteString(resultParts[i])
		}
		return b.String()
	}
}
