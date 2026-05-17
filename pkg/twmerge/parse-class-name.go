package twmerge

import "strings"

const importantModifier = "!"
const modifierSeparator = ':'

func createParseClassName(config *Config) ParseClassNameFn {
	parse := makeBaseParser()

	if config.Prefix != "" {
		fullPrefix := config.Prefix + string(modifierSeparator)
		base := parse
		parse = func(className string) ParsedClassName {
			if strings.HasPrefix(className, fullPrefix) {
				return base(className[len(fullPrefix):])
			}
			return ParsedClassName{
				IsExternal:                   true,
				BaseClassName:                className,
				MaybePostfixModifierPosition: -1,
			}
		}
	}

	if config.ExperimentalParseClassName != nil {
		base := parse
		hook := config.ExperimentalParseClassName
		parse = func(className string) ParsedClassName {
			return hook(ExperimentalParseClassNameParam{
				ClassName:      className,
				ParseClassName: base,
			})
		}
	}

	return parse
}

// makeBaseParser produces the parser without prefix / experimental wrapping.
// Mirrors createParseClassName's inner closure in parse-class-name.ts.
func makeBaseParser() ParseClassNameFn {
	return func(className string) ParsedClassName {
		modifiers := []string{}
		bracketDepth := 0
		parenDepth := 0
		modifierStart := 0
		postfixModifierPosition := -1

		for i := 0; i < len(className); i++ {
			c := className[i]

			if bracketDepth == 0 && parenDepth == 0 {
				if c == byte(modifierSeparator) {
					modifiers = append(modifiers, className[modifierStart:i])
					modifierStart = i + 1
					continue
				}
				if c == '/' {
					postfixModifierPosition = i
					continue
				}
			}

			switch c {
			case '[':
				bracketDepth++
			case ']':
				bracketDepth--
			case '(':
				parenDepth++
			case ')':
				parenDepth--
			}
		}

		var baseWithImportant string
		if len(modifiers) == 0 {
			baseWithImportant = className
		} else {
			baseWithImportant = className[modifierStart:]
		}

		baseClassName := baseWithImportant
		hasImportant := false

		if strings.HasSuffix(baseWithImportant, importantModifier) {
			baseClassName = baseWithImportant[:len(baseWithImportant)-1]
			hasImportant = true
		} else if strings.HasPrefix(baseWithImportant, importantModifier) {
			// Tailwind CSS v3 legacy syntax: leading `!`.
			baseClassName = baseWithImportant[1:]
			hasImportant = true
		}

		maybePostfix := -1
		if postfixModifierPosition != -1 && postfixModifierPosition > modifierStart {
			maybePostfix = postfixModifierPosition - modifierStart
		}

		return ParsedClassName{
			Modifiers:                    modifiers,
			HasImportantModifier:         hasImportant,
			BaseClassName:                baseClassName,
			MaybePostfixModifierPosition: maybePostfix,
		}
	}
}
