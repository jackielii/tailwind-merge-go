package twmerge

import (
	"regexp"
	"strconv"
	"strings"
)

// Regexes mirror validators.ts.
var (
	arbitraryValueRegex    = regexp.MustCompile(`(?i)^\[(?:(\w[\w-]*):)?(.+)\]$`)
	arbitraryVariableRegex = regexp.MustCompile(`(?i)^\((?:(\w[\w-]*):)?(.+)\)$`)
	fractionRegex          = regexp.MustCompile(`^\d+(?:\.\d+)?\/\d+(?:\.\d+)?$`)
	tshirtUnitRegex        = regexp.MustCompile(`^(\d+(\.\d+)?)?(xs|sm|md|lg|xl)$`)
	lengthUnitRegex        = regexp.MustCompile(
		`\d+(%|px|r?em|[sdl]?v([hwib]|min|max)|pt|pc|in|cm|mm|cap|ch|ex|r?lh|cq(w|h|i|b|min|max))|\b(calc|min|max|clamp)\(.+\)|^0$`,
	)
	colorFunctionRegex = regexp.MustCompile(`^(rgba?|hsla?|hwb|(ok)?(lab|lch)|color-mix)\(.+\)$`)
	shadowRegex        = regexp.MustCompile(`^(inset_)?-?((\d+)?\.?(\d+)[a-z]+|0)_-?((\d+)?\.?(\d+)[a-z]+|0)`)
	imageRegex         = regexp.MustCompile(
		`^(url|image|image-set|cross-fade|element|(repeating-)?(linear|radial|conic)-gradient)\(.+\)$`,
	)
)

// IsFraction reports whether v matches a fractional value like 1/2 or 1.5/2.5.
func IsFraction(v string) bool { return fractionRegex.MatchString(v) }

// IsNumber reports whether v parses as a finite number.
func IsNumber(v string) bool {
	if v == "" {
		return false
	}
	_, err := strconv.ParseFloat(v, 64)
	return err == nil
}

// IsInteger reports whether v parses as an integer.
func IsInteger(v string) bool {
	if v == "" {
		return false
	}
	_, err := strconv.Atoi(v)
	return err == nil
}

// IsPercent reports whether v ends with % and the prefix is a number.
func IsPercent(v string) bool {
	if !strings.HasSuffix(v, "%") {
		return false
	}
	return IsNumber(v[:len(v)-1])
}

// IsTshirtSize reports whether v matches a t-shirt size like sm, md, 2xl.
func IsTshirtSize(v string) bool { return tshirtUnitRegex.MatchString(v) }

// IsAny always returns true.
func IsAny(string) bool { return true }

// IsNever always returns false.
func IsNever(string) bool { return false }

func isLengthOnly(v string) bool {
	// colorFunctionRegex check avoids matching e.g. hsl(0 0% 0%) as a length.
	return lengthUnitRegex.MatchString(v) && !colorFunctionRegex.MatchString(v)
}

func isShadow(v string) bool { return shadowRegex.MatchString(v) }
func isImage(v string) bool  { return imageRegex.MatchString(v) }

// IsAnyNonArbitrary reports whether v is neither an arbitrary value nor an arbitrary variable.
func IsAnyNonArbitrary(v string) bool {
	return !IsArbitraryValue(v) && !IsArbitraryVariable(v)
}

// IsNamedContainerQuery matches @container/<name>, @container-size/<name>, @container-normal/<name>.
func IsNamedContainerQuery(v string) bool {
	if !strings.HasPrefix(v, "@container") {
		return false
	}
	if len(v) > 11 && v[10] == '/' {
		return true
	}
	if strings.HasPrefix(v[10:], "-size/") && len(v) > 16 {
		return true
	}
	if strings.HasPrefix(v[10:], "-normal/") && len(v) > 18 {
		return true
	}
	return false
}

// IsArbitrarySize matches [size:...] etc with an appropriate label.
func IsArbitrarySize(v string) bool { return getIsArbitraryValue(v, isLabelSize, IsNever) }

// IsArbitraryValue reports whether v is wrapped in [].
func IsArbitraryValue(v string) bool { return arbitraryValueRegex.MatchString(v) }

// IsArbitraryLength reports whether v is [length:...] or [<length-only-value>].
func IsArbitraryLength(v string) bool { return getIsArbitraryValue(v, isLabelLength, isLengthOnly) }

// IsArbitraryNumber reports whether v is [number:...] or [<number>].
func IsArbitraryNumber(v string) bool { return getIsArbitraryValue(v, isLabelNumber, IsNumber) }

// IsArbitraryWeight reports whether v is [number:...] / [weight:...].
func IsArbitraryWeight(v string) bool { return getIsArbitraryValue(v, isLabelWeight, IsAny) }

// IsArbitraryFamilyName reports whether v is [family-name:...].
func IsArbitraryFamilyName(v string) bool {
	return getIsArbitraryValue(v, isLabelFamilyName, IsNever)
}

// IsArbitraryPosition reports whether v is [position:...] / [percentage:...].
func IsArbitraryPosition(v string) bool { return getIsArbitraryValue(v, isLabelPosition, IsNever) }

// IsArbitraryImage reports whether v is [image:...] / [url:...] / [<image-fn>].
func IsArbitraryImage(v string) bool { return getIsArbitraryValue(v, isLabelImage, isImage) }

// IsArbitraryShadow reports whether v is [shadow:...] or [<shadow-value>].
func IsArbitraryShadow(v string) bool { return getIsArbitraryValue(v, isLabelShadow, isShadow) }

// IsArbitraryVariable reports whether v is wrapped in ().
func IsArbitraryVariable(v string) bool { return arbitraryVariableRegex.MatchString(v) }

// IsArbitraryVariableLength reports whether v is (length:...).
func IsArbitraryVariableLength(v string) bool {
	return getIsArbitraryVariable(v, isLabelLength, false)
}

// IsArbitraryVariableFamilyName reports whether v is (family-name:...).
func IsArbitraryVariableFamilyName(v string) bool {
	return getIsArbitraryVariable(v, isLabelFamilyName, false)
}

// IsArbitraryVariablePosition reports whether v is (position:...) or (percentage:...).
func IsArbitraryVariablePosition(v string) bool {
	return getIsArbitraryVariable(v, isLabelPosition, false)
}

// IsArbitraryVariableSize reports whether v is (length:...) / (size:...) / (bg-size:...).
func IsArbitraryVariableSize(v string) bool {
	return getIsArbitraryVariable(v, isLabelSize, false)
}

// IsArbitraryVariableImage reports whether v is (image:...) / (url:...).
func IsArbitraryVariableImage(v string) bool {
	return getIsArbitraryVariable(v, isLabelImage, false)
}

// IsArbitraryVariableShadow matches (shadow:...) or (...) without label.
func IsArbitraryVariableShadow(v string) bool {
	return getIsArbitraryVariable(v, isLabelShadow, true)
}

// IsArbitraryVariableWeight matches (number:...) / (weight:...) or (...) without label.
func IsArbitraryVariableWeight(v string) bool {
	return getIsArbitraryVariable(v, isLabelWeight, true)
}

func getIsArbitraryValue(v string, testLabel func(string) bool, testValue func(string) bool) bool {
	m := arbitraryValueRegex.FindStringSubmatch(v)
	if m == nil {
		return false
	}
	if m[1] != "" {
		return testLabel(m[1])
	}
	return testValue(m[2])
}

func getIsArbitraryVariable(v string, testLabel func(string) bool, shouldMatchNoLabel bool) bool {
	m := arbitraryVariableRegex.FindStringSubmatch(v)
	if m == nil {
		return false
	}
	if m[1] != "" {
		return testLabel(m[1])
	}
	return shouldMatchNoLabel
}

func isLabelPosition(l string) bool   { return l == "position" || l == "percentage" }
func isLabelImage(l string) bool      { return l == "image" || l == "url" }
func isLabelSize(l string) bool       { return l == "length" || l == "size" || l == "bg-size" }
func isLabelLength(l string) bool     { return l == "length" }
func isLabelNumber(l string) bool     { return l == "number" }
func isLabelFamilyName(l string) bool { return l == "family-name" }
func isLabelWeight(l string) bool     { return l == "number" || l == "weight" }
func isLabelShadow(l string) bool     { return l == "shadow" }
