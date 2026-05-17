package twmerge

import "sort"

// createSortModifiers returns a function that sorts modifiers, preserving the
// relative order of arbitrary `[...]` and order-sensitive modifiers while
// alphabetising the segments between them.
func createSortModifiers(config *Config) func([]string) []string {
	orderSensitive := make(map[string]struct{}, len(config.OrderSensitiveModifiers))
	for _, m := range config.OrderSensitiveModifiers {
		orderSensitive[m] = struct{}{}
	}

	return func(modifiers []string) []string {
		if len(modifiers) < 2 {
			return modifiers
		}
		result := make([]string, 0, len(modifiers))
		segment := make([]string, 0)

		flush := func() {
			if len(segment) == 0 {
				return
			}
			sort.Strings(segment)
			result = append(result, segment...)
			segment = segment[:0]
		}

		for _, m := range modifiers {
			isArbitrary := len(m) > 0 && m[0] == '['
			_, isOrderSensitive := orderSensitive[m]
			if isArbitrary || isOrderSensitive {
				flush()
				result = append(result, m)
				continue
			}
			segment = append(segment, m)
		}
		flush()

		return result
	}
}
