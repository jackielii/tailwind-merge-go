package twmerge

// GetDefaultConfig returns a fresh copy of the default tailwind-merge configuration.
// Each call returns new maps/slices so mutation by callers doesn't leak.
func GetDefaultConfig() *Config {
	// Theme getters for theme variable namespaces.
	// @see https://tailwindcss.com/docs/theme#theme-variable-namespaces
	themeColor := FromTheme("color")
	themeFont := FromTheme("font")
	themeText := FromTheme("text")
	themeFontWeight := FromTheme("font-weight")
	themeTracking := FromTheme("tracking")
	themeLeading := FromTheme("leading")
	themeBreakpoint := FromTheme("breakpoint")
	themeContainer := FromTheme("container")
	themeSpacing := FromTheme("spacing")
	themeRadius := FromTheme("radius")
	themeShadow := FromTheme("shadow")
	themeInsetShadow := FromTheme("inset-shadow")
	themeTextShadow := FromTheme("text-shadow")
	themeDropShadow := FromTheme("drop-shadow")
	themeBlur := FromTheme("blur")
	themePerspective := FromTheme("perspective")
	themeAspect := FromTheme("aspect")
	themeEase := FromTheme("ease")
	themeAnimate := FromTheme("animate")

	// concat builds a new ClassGroup by appending all provided groups in order.
	// This mirrors TS spread syntax (e.g. [...scaleX(), 'foo']).
	concat := func(groups ...ClassGroup) ClassGroup {
		var out ClassGroup
		for _, g := range groups {
			out = append(out, g...)
		}
		return out
	}

	// Helpers to avoid repeating the same scales.
	// Each helper returns a freshly allocated slice so mutation by callers
	// doesn't leak to other groups.
	scaleBreak := func() ClassGroup {
		return ClassGroup{"auto", "avoid", "all", "avoid-page", "page", "left", "right", "column"}
	}
	scalePosition := func() ClassGroup {
		return ClassGroup{
			"center",
			"top",
			"bottom",
			"left",
			"right",
			"top-left",
			// Deprecated since Tailwind CSS v4.1.0, see https://github.com/tailwindlabs/tailwindcss/pull/17378
			"left-top",
			"top-right",
			// Deprecated since Tailwind CSS v4.1.0, see https://github.com/tailwindlabs/tailwindcss/pull/17378
			"right-top",
			"bottom-right",
			// Deprecated since Tailwind CSS v4.1.0, see https://github.com/tailwindlabs/tailwindcss/pull/17378
			"right-bottom",
			"bottom-left",
			// Deprecated since Tailwind CSS v4.1.0, see https://github.com/tailwindlabs/tailwindcss/pull/17378
			"left-bottom",
		}
	}
	scalePositionWithArbitrary := func() ClassGroup {
		return concat(scalePosition(), ClassGroup{IsArbitraryVariable, IsArbitraryValue})
	}
	scaleOverflow := func() ClassGroup {
		return ClassGroup{"auto", "hidden", "clip", "visible", "scroll"}
	}
	scaleOverscroll := func() ClassGroup {
		return ClassGroup{"auto", "contain", "none"}
	}
	scaleUnambiguousSpacing := func() ClassGroup {
		return ClassGroup{IsArbitraryVariable, IsArbitraryValue, themeSpacing}
	}
	scaleInset := func() ClassGroup {
		return concat(ClassGroup{IsFraction, "full", "auto"}, scaleUnambiguousSpacing())
	}
	scaleGridTemplateColsRows := func() ClassGroup {
		return ClassGroup{IsInteger, "none", "subgrid", IsArbitraryVariable, IsArbitraryValue}
	}
	scaleGridColRowStartAndEnd := func() ClassGroup {
		return ClassGroup{
			"auto",
			ClassObject{"span": ClassGroup{"full", IsInteger, IsArbitraryVariable, IsArbitraryValue}},
			IsInteger,
			IsArbitraryVariable,
			IsArbitraryValue,
		}
	}
	scaleGridColRowStartOrEnd := func() ClassGroup {
		return ClassGroup{IsInteger, "auto", IsArbitraryVariable, IsArbitraryValue}
	}
	scaleGridAutoColsRows := func() ClassGroup {
		return ClassGroup{"auto", "min", "max", "fr", IsArbitraryVariable, IsArbitraryValue}
	}
	scaleAlignPrimaryAxis := func() ClassGroup {
		return ClassGroup{
			"start",
			"end",
			"center",
			"between",
			"around",
			"evenly",
			"stretch",
			"baseline",
			"center-safe",
			"end-safe",
		}
	}
	scaleAlignSecondaryAxis := func() ClassGroup {
		return ClassGroup{"start", "end", "center", "stretch", "center-safe", "end-safe"}
	}
	scaleMargin := func() ClassGroup {
		return concat(ClassGroup{"auto"}, scaleUnambiguousSpacing())
	}
	scaleSizing := func() ClassGroup {
		return concat(ClassGroup{
			IsFraction,
			"auto",
			"full",
			"dvw",
			"dvh",
			"lvw",
			"lvh",
			"svw",
			"svh",
			"min",
			"max",
			"fit",
		}, scaleUnambiguousSpacing())
	}
	scaleSizingInline := func() ClassGroup {
		return concat(ClassGroup{
			IsFraction,
			"screen",
			"full",
			"dvw",
			"lvw",
			"svw",
			"min",
			"max",
			"fit",
		}, scaleUnambiguousSpacing())
	}
	scaleSizingBlock := func() ClassGroup {
		return concat(ClassGroup{
			IsFraction,
			"screen",
			"full",
			"lh",
			"dvh",
			"lvh",
			"svh",
			"min",
			"max",
			"fit",
		}, scaleUnambiguousSpacing())
	}
	scaleColor := func() ClassGroup {
		return ClassGroup{themeColor, IsArbitraryVariable, IsArbitraryValue}
	}
	scaleBgPosition := func() ClassGroup {
		return concat(scalePosition(), ClassGroup{
			IsArbitraryVariablePosition,
			IsArbitraryPosition,
			ClassObject{"position": ClassGroup{IsArbitraryVariable, IsArbitraryValue}},
		})
	}
	scaleBgRepeat := func() ClassGroup {
		return ClassGroup{
			"no-repeat",
			ClassObject{"repeat": ClassGroup{"", "x", "y", "space", "round"}},
		}
	}
	scaleBgSize := func() ClassGroup {
		return ClassGroup{
			"auto",
			"cover",
			"contain",
			IsArbitraryVariableSize,
			IsArbitrarySize,
			ClassObject{"size": ClassGroup{IsArbitraryVariable, IsArbitraryValue}},
		}
	}
	scaleGradientStopPosition := func() ClassGroup {
		return ClassGroup{IsPercent, IsArbitraryVariableLength, IsArbitraryLength}
	}
	scaleRadius := func() ClassGroup {
		return ClassGroup{
			// Deprecated since Tailwind CSS v4.0.0
			"",
			"none",
			"full",
			themeRadius,
			IsArbitraryVariable,
			IsArbitraryValue,
		}
	}
	scaleBorderWidth := func() ClassGroup {
		return ClassGroup{"", IsNumber, IsArbitraryVariableLength, IsArbitraryLength}
	}
	scaleLineStyle := func() ClassGroup {
		return ClassGroup{"solid", "dashed", "dotted", "double"}
	}
	scaleBlendMode := func() ClassGroup {
		return ClassGroup{
			"normal",
			"multiply",
			"screen",
			"overlay",
			"darken",
			"lighten",
			"color-dodge",
			"color-burn",
			"hard-light",
			"soft-light",
			"difference",
			"exclusion",
			"hue",
			"saturation",
			"color",
			"luminosity",
		}
	}
	scaleMaskImagePosition := func() ClassGroup {
		return ClassGroup{IsNumber, IsPercent, IsArbitraryVariablePosition, IsArbitraryPosition}
	}
	scaleBlur := func() ClassGroup {
		return ClassGroup{
			// Deprecated since Tailwind CSS v4.0.0
			"",
			"none",
			themeBlur,
			IsArbitraryVariable,
			IsArbitraryValue,
		}
	}
	scaleRotate := func() ClassGroup {
		return ClassGroup{"none", IsNumber, IsArbitraryVariable, IsArbitraryValue}
	}
	scaleScale := func() ClassGroup {
		return ClassGroup{"none", IsNumber, IsArbitraryVariable, IsArbitraryValue}
	}
	scaleSkew := func() ClassGroup {
		return ClassGroup{IsNumber, IsArbitraryVariable, IsArbitraryValue}
	}
	scaleTranslate := func() ClassGroup {
		return concat(ClassGroup{IsFraction, "full"}, scaleUnambiguousSpacing())
	}

	return &Config{
		CacheSize: 500,
		Theme: Theme{
			"animate":      ClassGroup{"spin", "ping", "pulse", "bounce"},
			"aspect":       ClassGroup{"video"},
			"blur":         ClassGroup{IsTshirtSize},
			"breakpoint":   ClassGroup{IsTshirtSize},
			"color":        ClassGroup{IsAny},
			"container":    ClassGroup{IsTshirtSize},
			"drop-shadow": ClassGroup{IsTshirtSize},
			"ease":         ClassGroup{"in", "out", "in-out"},
			"font":         ClassGroup{IsAnyNonArbitrary},
			"font-weight": ClassGroup{
				"thin",
				"extralight",
				"light",
				"normal",
				"medium",
				"semibold",
				"bold",
				"extrabold",
				"black",
			},
			"inset-shadow": ClassGroup{IsTshirtSize},
			"leading":      ClassGroup{"none", "tight", "snug", "normal", "relaxed", "loose"},
			"perspective":  ClassGroup{"dramatic", "near", "normal", "midrange", "distant", "none"},
			"radius":       ClassGroup{IsTshirtSize},
			"shadow":       ClassGroup{IsTshirtSize},
			"spacing":      ClassGroup{"px", IsNumber},
			"text":         ClassGroup{IsTshirtSize},
			"text-shadow":  ClassGroup{IsTshirtSize},
			"tracking":     ClassGroup{"tighter", "tight", "normal", "wide", "wider", "widest"},
		},
		ClassGroups: ClassGroupsMap{
			// --------------
			// --- Layout ---
			// --------------

			// Aspect Ratio
			// @see https://tailwindcss.com/docs/aspect-ratio
			"aspect": ClassGroup{
				ClassObject{
					"aspect": ClassGroup{
						"auto",
						"square",
						IsFraction,
						IsArbitraryValue,
						IsArbitraryVariable,
						themeAspect,
					},
				},
			},
			// Container
			// @see https://tailwindcss.com/docs/container
			// @deprecated since Tailwind CSS v4.0.0
			"container": ClassGroup{"container"},
			// Container Type
			// @see https://tailwindcss.com/docs/responsive-design#container-queries
			"container-type": ClassGroup{
				ClassObject{"@container": ClassGroup{"", "normal", "size", IsArbitraryVariable, IsArbitraryValue}},
			},
			// Container Name
			// @see https://tailwindcss.com/docs/responsive-design#named-containers
			"container-named": ClassGroup{IsNamedContainerQuery},
			// Columns
			// @see https://tailwindcss.com/docs/columns
			"columns": ClassGroup{
				ClassObject{"columns": ClassGroup{IsNumber, IsArbitraryValue, IsArbitraryVariable, themeContainer}},
			},
			// Break After
			// @see https://tailwindcss.com/docs/break-after
			"break-after": ClassGroup{ClassObject{"break-after": scaleBreak()}},
			// Break Before
			// @see https://tailwindcss.com/docs/break-before
			"break-before": ClassGroup{ClassObject{"break-before": scaleBreak()}},
			// Break Inside
			// @see https://tailwindcss.com/docs/break-inside
			"break-inside": ClassGroup{
				ClassObject{"break-inside": ClassGroup{"auto", "avoid", "avoid-page", "avoid-column"}},
			},
			// Box Decoration Break
			// @see https://tailwindcss.com/docs/box-decoration-break
			"box-decoration": ClassGroup{ClassObject{"box-decoration": ClassGroup{"slice", "clone"}}},
			// Box Sizing
			// @see https://tailwindcss.com/docs/box-sizing
			"box": ClassGroup{ClassObject{"box": ClassGroup{"border", "content"}}},
			// Display
			// @see https://tailwindcss.com/docs/display
			"display": ClassGroup{
				"block",
				"inline-block",
				"inline",
				"flex",
				"inline-flex",
				"table",
				"inline-table",
				"table-caption",
				"table-cell",
				"table-column",
				"table-column-group",
				"table-footer-group",
				"table-header-group",
				"table-row-group",
				"table-row",
				"flow-root",
				"grid",
				"inline-grid",
				"contents",
				"list-item",
				"hidden",
			},
			// Screen Reader Only
			// @see https://tailwindcss.com/docs/display#screen-reader-only
			"sr": ClassGroup{"sr-only", "not-sr-only"},
			// Floats
			// @see https://tailwindcss.com/docs/float
			"float": ClassGroup{ClassObject{"float": ClassGroup{"right", "left", "none", "start", "end"}}},
			// Clear
			// @see https://tailwindcss.com/docs/clear
			"clear": ClassGroup{ClassObject{"clear": ClassGroup{"left", "right", "both", "none", "start", "end"}}},
			// Isolation
			// @see https://tailwindcss.com/docs/isolation
			"isolation": ClassGroup{"isolate", "isolation-auto"},
			// Object Fit
			// @see https://tailwindcss.com/docs/object-fit
			"object-fit": ClassGroup{ClassObject{"object": ClassGroup{"contain", "cover", "fill", "none", "scale-down"}}},
			// Object Position
			// @see https://tailwindcss.com/docs/object-position
			"object-position": ClassGroup{ClassObject{"object": scalePositionWithArbitrary()}},
			// Overflow
			// @see https://tailwindcss.com/docs/overflow
			"overflow": ClassGroup{ClassObject{"overflow": scaleOverflow()}},
			// Overflow X
			// @see https://tailwindcss.com/docs/overflow
			"overflow-x": ClassGroup{ClassObject{"overflow-x": scaleOverflow()}},
			// Overflow Y
			// @see https://tailwindcss.com/docs/overflow
			"overflow-y": ClassGroup{ClassObject{"overflow-y": scaleOverflow()}},
			// Overscroll Behavior
			// @see https://tailwindcss.com/docs/overscroll-behavior
			"overscroll": ClassGroup{ClassObject{"overscroll": scaleOverscroll()}},
			// Overscroll Behavior X
			// @see https://tailwindcss.com/docs/overscroll-behavior
			"overscroll-x": ClassGroup{ClassObject{"overscroll-x": scaleOverscroll()}},
			// Overscroll Behavior Y
			// @see https://tailwindcss.com/docs/overscroll-behavior
			"overscroll-y": ClassGroup{ClassObject{"overscroll-y": scaleOverscroll()}},
			// Position
			// @see https://tailwindcss.com/docs/position
			"position": ClassGroup{"static", "fixed", "absolute", "relative", "sticky"},
			// Inset
			// @see https://tailwindcss.com/docs/top-right-bottom-left
			"inset": ClassGroup{ClassObject{"inset": scaleInset()}},
			// Inset Inline
			// @see https://tailwindcss.com/docs/top-right-bottom-left
			"inset-x": ClassGroup{ClassObject{"inset-x": scaleInset()}},
			// Inset Block
			// @see https://tailwindcss.com/docs/top-right-bottom-left
			"inset-y": ClassGroup{ClassObject{"inset-y": scaleInset()}},
			// Inset Inline Start
			// @see https://tailwindcss.com/docs/top-right-bottom-left
			// @todo class group will be renamed to `inset-s` in next major release
			"start": ClassGroup{
				ClassObject{
					"inset-s": scaleInset(),
					// @deprecated since Tailwind CSS v4.2.0 in favor of `inset-s-*` utilities.
					// @see https://github.com/tailwindlabs/tailwindcss/pull/19613
					"start": scaleInset(),
				},
			},
			// Inset Inline End
			// @see https://tailwindcss.com/docs/top-right-bottom-left
			// @todo class group will be renamed to `inset-e` in next major release
			"end": ClassGroup{
				ClassObject{
					"inset-e": scaleInset(),
					// @deprecated since Tailwind CSS v4.2.0 in favor of `inset-e-*` utilities.
					// @see https://github.com/tailwindlabs/tailwindcss/pull/19613
					"end": scaleInset(),
				},
			},
			// Inset Block Start
			// @see https://tailwindcss.com/docs/top-right-bottom-left
			"inset-bs": ClassGroup{ClassObject{"inset-bs": scaleInset()}},
			// Inset Block End
			// @see https://tailwindcss.com/docs/top-right-bottom-left
			"inset-be": ClassGroup{ClassObject{"inset-be": scaleInset()}},
			// Top
			// @see https://tailwindcss.com/docs/top-right-bottom-left
			"top": ClassGroup{ClassObject{"top": scaleInset()}},
			// Right
			// @see https://tailwindcss.com/docs/top-right-bottom-left
			"right": ClassGroup{ClassObject{"right": scaleInset()}},
			// Bottom
			// @see https://tailwindcss.com/docs/top-right-bottom-left
			"bottom": ClassGroup{ClassObject{"bottom": scaleInset()}},
			// Left
			// @see https://tailwindcss.com/docs/top-right-bottom-left
			"left": ClassGroup{ClassObject{"left": scaleInset()}},
			// Visibility
			// @see https://tailwindcss.com/docs/visibility
			"visibility": ClassGroup{"visible", "invisible", "collapse"},
			// Z-Index
			// @see https://tailwindcss.com/docs/z-index
			"z": ClassGroup{ClassObject{"z": ClassGroup{IsInteger, "auto", IsArbitraryVariable, IsArbitraryValue}}},

			// ------------------------
			// --- Flexbox and Grid ---
			// ------------------------

			// Flex Basis
			// @see https://tailwindcss.com/docs/flex-basis
			"basis": ClassGroup{
				ClassObject{
					"basis": concat(ClassGroup{IsFraction, "full", "auto", themeContainer}, scaleUnambiguousSpacing()),
				},
			},
			// Flex Direction
			// @see https://tailwindcss.com/docs/flex-direction
			"flex-direction": ClassGroup{ClassObject{"flex": ClassGroup{"row", "row-reverse", "col", "col-reverse"}}},
			// Flex Wrap
			// @see https://tailwindcss.com/docs/flex-wrap
			"flex-wrap": ClassGroup{ClassObject{"flex": ClassGroup{"nowrap", "wrap", "wrap-reverse"}}},
			// Flex
			// @see https://tailwindcss.com/docs/flex
			"flex": ClassGroup{ClassObject{"flex": ClassGroup{IsNumber, IsFraction, "auto", "initial", "none", IsArbitraryValue}}},
			// Flex Grow
			// @see https://tailwindcss.com/docs/flex-grow
			"grow": ClassGroup{ClassObject{"grow": ClassGroup{"", IsNumber, IsArbitraryVariable, IsArbitraryValue}}},
			// Flex Shrink
			// @see https://tailwindcss.com/docs/flex-shrink
			"shrink": ClassGroup{ClassObject{"shrink": ClassGroup{"", IsNumber, IsArbitraryVariable, IsArbitraryValue}}},
			// Order
			// @see https://tailwindcss.com/docs/order
			"order": ClassGroup{
				ClassObject{
					"order": ClassGroup{
						IsInteger,
						"first",
						"last",
						"none",
						IsArbitraryVariable,
						IsArbitraryValue,
					},
				},
			},
			// Grid Template Columns
			// @see https://tailwindcss.com/docs/grid-template-columns
			"grid-cols": ClassGroup{ClassObject{"grid-cols": scaleGridTemplateColsRows()}},
			// Grid Column Start / End
			// @see https://tailwindcss.com/docs/grid-column
			"col-start-end": ClassGroup{ClassObject{"col": scaleGridColRowStartAndEnd()}},
			// Grid Column Start
			// @see https://tailwindcss.com/docs/grid-column
			"col-start": ClassGroup{ClassObject{"col-start": scaleGridColRowStartOrEnd()}},
			// Grid Column End
			// @see https://tailwindcss.com/docs/grid-column
			"col-end": ClassGroup{ClassObject{"col-end": scaleGridColRowStartOrEnd()}},
			// Grid Template Rows
			// @see https://tailwindcss.com/docs/grid-template-rows
			"grid-rows": ClassGroup{ClassObject{"grid-rows": scaleGridTemplateColsRows()}},
			// Grid Row Start / End
			// @see https://tailwindcss.com/docs/grid-row
			"row-start-end": ClassGroup{ClassObject{"row": scaleGridColRowStartAndEnd()}},
			// Grid Row Start
			// @see https://tailwindcss.com/docs/grid-row
			"row-start": ClassGroup{ClassObject{"row-start": scaleGridColRowStartOrEnd()}},
			// Grid Row End
			// @see https://tailwindcss.com/docs/grid-row
			"row-end": ClassGroup{ClassObject{"row-end": scaleGridColRowStartOrEnd()}},
			// Grid Auto Flow
			// @see https://tailwindcss.com/docs/grid-auto-flow
			"grid-flow": ClassGroup{ClassObject{"grid-flow": ClassGroup{"row", "col", "dense", "row-dense", "col-dense"}}},
			// Grid Auto Columns
			// @see https://tailwindcss.com/docs/grid-auto-columns
			"auto-cols": ClassGroup{ClassObject{"auto-cols": scaleGridAutoColsRows()}},
			// Grid Auto Rows
			// @see https://tailwindcss.com/docs/grid-auto-rows
			"auto-rows": ClassGroup{ClassObject{"auto-rows": scaleGridAutoColsRows()}},
			// Gap
			// @see https://tailwindcss.com/docs/gap
			"gap": ClassGroup{ClassObject{"gap": scaleUnambiguousSpacing()}},
			// Gap X
			// @see https://tailwindcss.com/docs/gap
			"gap-x": ClassGroup{ClassObject{"gap-x": scaleUnambiguousSpacing()}},
			// Gap Y
			// @see https://tailwindcss.com/docs/gap
			"gap-y": ClassGroup{ClassObject{"gap-y": scaleUnambiguousSpacing()}},
			// Justify Content
			// @see https://tailwindcss.com/docs/justify-content
			"justify-content": ClassGroup{ClassObject{"justify": concat(scaleAlignPrimaryAxis(), ClassGroup{"normal"})}},
			// Justify Items
			// @see https://tailwindcss.com/docs/justify-items
			"justify-items": ClassGroup{ClassObject{"justify-items": concat(scaleAlignSecondaryAxis(), ClassGroup{"normal"})}},
			// Justify Self
			// @see https://tailwindcss.com/docs/justify-self
			"justify-self": ClassGroup{ClassObject{"justify-self": concat(ClassGroup{"auto"}, scaleAlignSecondaryAxis())}},
			// Align Content
			// @see https://tailwindcss.com/docs/align-content
			"align-content": ClassGroup{ClassObject{"content": concat(ClassGroup{"normal"}, scaleAlignPrimaryAxis())}},
			// Align Items
			// @see https://tailwindcss.com/docs/align-items
			"align-items": ClassGroup{ClassObject{"items": concat(scaleAlignSecondaryAxis(), ClassGroup{ClassObject{"baseline": ClassGroup{"", "last"}}})}},
			// Align Self
			// @see https://tailwindcss.com/docs/align-self
			"align-self": ClassGroup{
				ClassObject{"self": concat(ClassGroup{"auto"}, scaleAlignSecondaryAxis(), ClassGroup{ClassObject{"baseline": ClassGroup{"", "last"}}})},
			},
			// Place Content
			// @see https://tailwindcss.com/docs/place-content
			"place-content": ClassGroup{ClassObject{"place-content": scaleAlignPrimaryAxis()}},
			// Place Items
			// @see https://tailwindcss.com/docs/place-items
			"place-items": ClassGroup{ClassObject{"place-items": concat(scaleAlignSecondaryAxis(), ClassGroup{"baseline"})}},
			// Place Self
			// @see https://tailwindcss.com/docs/place-self
			"place-self": ClassGroup{ClassObject{"place-self": concat(ClassGroup{"auto"}, scaleAlignSecondaryAxis())}},
			// Spacing
			// Padding
			// @see https://tailwindcss.com/docs/padding
			"p": ClassGroup{ClassObject{"p": scaleUnambiguousSpacing()}},
			// Padding Inline
			// @see https://tailwindcss.com/docs/padding
			"px": ClassGroup{ClassObject{"px": scaleUnambiguousSpacing()}},
			// Padding Block
			// @see https://tailwindcss.com/docs/padding
			"py": ClassGroup{ClassObject{"py": scaleUnambiguousSpacing()}},
			// Padding Inline Start
			// @see https://tailwindcss.com/docs/padding
			"ps": ClassGroup{ClassObject{"ps": scaleUnambiguousSpacing()}},
			// Padding Inline End
			// @see https://tailwindcss.com/docs/padding
			"pe": ClassGroup{ClassObject{"pe": scaleUnambiguousSpacing()}},
			// Padding Block Start
			// @see https://tailwindcss.com/docs/padding
			"pbs": ClassGroup{ClassObject{"pbs": scaleUnambiguousSpacing()}},
			// Padding Block End
			// @see https://tailwindcss.com/docs/padding
			"pbe": ClassGroup{ClassObject{"pbe": scaleUnambiguousSpacing()}},
			// Padding Top
			// @see https://tailwindcss.com/docs/padding
			"pt": ClassGroup{ClassObject{"pt": scaleUnambiguousSpacing()}},
			// Padding Right
			// @see https://tailwindcss.com/docs/padding
			"pr": ClassGroup{ClassObject{"pr": scaleUnambiguousSpacing()}},
			// Padding Bottom
			// @see https://tailwindcss.com/docs/padding
			"pb": ClassGroup{ClassObject{"pb": scaleUnambiguousSpacing()}},
			// Padding Left
			// @see https://tailwindcss.com/docs/padding
			"pl": ClassGroup{ClassObject{"pl": scaleUnambiguousSpacing()}},
			// Margin
			// @see https://tailwindcss.com/docs/margin
			"m": ClassGroup{ClassObject{"m": scaleMargin()}},
			// Margin Inline
			// @see https://tailwindcss.com/docs/margin
			"mx": ClassGroup{ClassObject{"mx": scaleMargin()}},
			// Margin Block
			// @see https://tailwindcss.com/docs/margin
			"my": ClassGroup{ClassObject{"my": scaleMargin()}},
			// Margin Inline Start
			// @see https://tailwindcss.com/docs/margin
			"ms": ClassGroup{ClassObject{"ms": scaleMargin()}},
			// Margin Inline End
			// @see https://tailwindcss.com/docs/margin
			"me": ClassGroup{ClassObject{"me": scaleMargin()}},
			// Margin Block Start
			// @see https://tailwindcss.com/docs/margin
			"mbs": ClassGroup{ClassObject{"mbs": scaleMargin()}},
			// Margin Block End
			// @see https://tailwindcss.com/docs/margin
			"mbe": ClassGroup{ClassObject{"mbe": scaleMargin()}},
			// Margin Top
			// @see https://tailwindcss.com/docs/margin
			"mt": ClassGroup{ClassObject{"mt": scaleMargin()}},
			// Margin Right
			// @see https://tailwindcss.com/docs/margin
			"mr": ClassGroup{ClassObject{"mr": scaleMargin()}},
			// Margin Bottom
			// @see https://tailwindcss.com/docs/margin
			"mb": ClassGroup{ClassObject{"mb": scaleMargin()}},
			// Margin Left
			// @see https://tailwindcss.com/docs/margin
			"ml": ClassGroup{ClassObject{"ml": scaleMargin()}},
			// Space Between X
			// @see https://tailwindcss.com/docs/margin#adding-space-between-children
			"space-x": ClassGroup{ClassObject{"space-x": scaleUnambiguousSpacing()}},
			// Space Between X Reverse
			// @see https://tailwindcss.com/docs/margin#adding-space-between-children
			"space-x-reverse": ClassGroup{"space-x-reverse"},
			// Space Between Y
			// @see https://tailwindcss.com/docs/margin#adding-space-between-children
			"space-y": ClassGroup{ClassObject{"space-y": scaleUnambiguousSpacing()}},
			// Space Between Y Reverse
			// @see https://tailwindcss.com/docs/margin#adding-space-between-children
			"space-y-reverse": ClassGroup{"space-y-reverse"},

			// --------------
			// --- Sizing ---
			// --------------

			// Size
			// @see https://tailwindcss.com/docs/width#setting-both-width-and-height
			"size": ClassGroup{ClassObject{"size": scaleSizing()}},
			// Inline Size
			// @see https://tailwindcss.com/docs/width
			"inline-size": ClassGroup{ClassObject{"inline": concat(ClassGroup{"auto"}, scaleSizingInline())}},
			// Min-Inline Size
			// @see https://tailwindcss.com/docs/min-width
			"min-inline-size": ClassGroup{ClassObject{"min-inline": concat(ClassGroup{"auto"}, scaleSizingInline())}},
			// Max-Inline Size
			// @see https://tailwindcss.com/docs/max-width
			"max-inline-size": ClassGroup{ClassObject{"max-inline": concat(ClassGroup{"none"}, scaleSizingInline())}},
			// Block Size
			// @see https://tailwindcss.com/docs/height
			"block-size": ClassGroup{ClassObject{"block": concat(ClassGroup{"auto"}, scaleSizingBlock())}},
			// Min-Block Size
			// @see https://tailwindcss.com/docs/min-height
			"min-block-size": ClassGroup{ClassObject{"min-block": concat(ClassGroup{"auto"}, scaleSizingBlock())}},
			// Max-Block Size
			// @see https://tailwindcss.com/docs/max-height
			"max-block-size": ClassGroup{ClassObject{"max-block": concat(ClassGroup{"none"}, scaleSizingBlock())}},
			// Width
			// @see https://tailwindcss.com/docs/width
			"w": ClassGroup{ClassObject{"w": concat(ClassGroup{themeContainer, "screen"}, scaleSizing())}},
			// Min-Width
			// @see https://tailwindcss.com/docs/min-width
			"min-w": ClassGroup{
				ClassObject{
					"min-w": concat(ClassGroup{
						themeContainer,
						"screen",
						// Deprecated. @see https://github.com/tailwindlabs/tailwindcss.com/issues/2027#issuecomment-2620152757
						"none",
					}, scaleSizing()),
				},
			},
			// Max-Width
			// @see https://tailwindcss.com/docs/max-width
			"max-w": ClassGroup{
				ClassObject{
					"max-w": concat(ClassGroup{
						themeContainer,
						"screen",
						"none",
						// Deprecated since Tailwind CSS v4.0.0. @see https://github.com/tailwindlabs/tailwindcss.com/issues/2027#issuecomment-2620152757
						"prose",
						// Deprecated since Tailwind CSS v4.0.0. @see https://github.com/tailwindlabs/tailwindcss.com/issues/2027#issuecomment-2620152757
						ClassObject{"screen": ClassGroup{themeBreakpoint}},
					}, scaleSizing()),
				},
			},
			// Height
			// @see https://tailwindcss.com/docs/height
			"h": ClassGroup{ClassObject{"h": concat(ClassGroup{"screen", "lh"}, scaleSizing())}},
			// Min-Height
			// @see https://tailwindcss.com/docs/min-height
			"min-h": ClassGroup{ClassObject{"min-h": concat(ClassGroup{"screen", "lh", "none"}, scaleSizing())}},
			// Max-Height
			// @see https://tailwindcss.com/docs/max-height
			"max-h": ClassGroup{ClassObject{"max-h": concat(ClassGroup{"screen", "lh"}, scaleSizing())}},

			// ------------------
			// --- Typography ---
			// ------------------

			// Font Size
			// @see https://tailwindcss.com/docs/font-size
			"font-size": ClassGroup{
				ClassObject{"text": ClassGroup{"base", themeText, IsArbitraryVariableLength, IsArbitraryLength}},
			},
			// Font Smoothing
			// @see https://tailwindcss.com/docs/font-smoothing
			"font-smoothing": ClassGroup{"antialiased", "subpixel-antialiased"},
			// Font Style
			// @see https://tailwindcss.com/docs/font-style
			"font-style": ClassGroup{"italic", "not-italic"},
			// Font Weight
			// @see https://tailwindcss.com/docs/font-weight
			"font-weight": ClassGroup{
				ClassObject{
					"font": ClassGroup{themeFontWeight, IsArbitraryVariableWeight, IsArbitraryWeight},
				},
			},
			// Font Stretch
			// @see https://tailwindcss.com/docs/font-stretch
			"font-stretch": ClassGroup{
				ClassObject{
					"font-stretch": ClassGroup{
						"ultra-condensed",
						"extra-condensed",
						"condensed",
						"semi-condensed",
						"normal",
						"semi-expanded",
						"expanded",
						"extra-expanded",
						"ultra-expanded",
						IsPercent,
						IsArbitraryValue,
					},
				},
			},
			// Font Family
			// @see https://tailwindcss.com/docs/font-family
			"font-family": ClassGroup{
				ClassObject{"font": ClassGroup{IsArbitraryVariableFamilyName, IsArbitraryFamilyName, themeFont}},
			},
			// Font Feature Settings
			// @see https://tailwindcss.com/docs/font-feature-settings
			"font-features": ClassGroup{ClassObject{"font-features": ClassGroup{IsArbitraryValue}}},
			// Font Variant Numeric
			// @see https://tailwindcss.com/docs/font-variant-numeric
			"fvn-normal": ClassGroup{"normal-nums"},
			// Font Variant Numeric
			// @see https://tailwindcss.com/docs/font-variant-numeric
			"fvn-ordinal": ClassGroup{"ordinal"},
			// Font Variant Numeric
			// @see https://tailwindcss.com/docs/font-variant-numeric
			"fvn-slashed-zero": ClassGroup{"slashed-zero"},
			// Font Variant Numeric
			// @see https://tailwindcss.com/docs/font-variant-numeric
			"fvn-figure": ClassGroup{"lining-nums", "oldstyle-nums"},
			// Font Variant Numeric
			// @see https://tailwindcss.com/docs/font-variant-numeric
			"fvn-spacing": ClassGroup{"proportional-nums", "tabular-nums"},
			// Font Variant Numeric
			// @see https://tailwindcss.com/docs/font-variant-numeric
			"fvn-fraction": ClassGroup{"diagonal-fractions", "stacked-fractions"},
			// Letter Spacing
			// @see https://tailwindcss.com/docs/letter-spacing
			"tracking": ClassGroup{ClassObject{"tracking": ClassGroup{themeTracking, IsArbitraryVariable, IsArbitraryValue}}},
			// Line Clamp
			// @see https://tailwindcss.com/docs/line-clamp
			"line-clamp": ClassGroup{
				ClassObject{"line-clamp": ClassGroup{IsNumber, "none", IsArbitraryVariable, IsArbitraryNumber}},
			},
			// Line Height
			// @see https://tailwindcss.com/docs/line-height
			"leading": ClassGroup{
				ClassObject{
					"leading": concat(ClassGroup{
						// Deprecated since Tailwind CSS v4.0.0. @see https://github.com/tailwindlabs/tailwindcss.com/issues/2027#issuecomment-2620152757
						themeLeading,
					}, scaleUnambiguousSpacing()),
				},
			},
			// List Style Image
			// @see https://tailwindcss.com/docs/list-style-image
			"list-image": ClassGroup{ClassObject{"list-image": ClassGroup{"none", IsArbitraryVariable, IsArbitraryValue}}},
			// List Style Position
			// @see https://tailwindcss.com/docs/list-style-position
			"list-style-position": ClassGroup{ClassObject{"list": ClassGroup{"inside", "outside"}}},
			// List Style Type
			// @see https://tailwindcss.com/docs/list-style-type
			"list-style-type": ClassGroup{
				ClassObject{"list": ClassGroup{"disc", "decimal", "none", IsArbitraryVariable, IsArbitraryValue}},
			},
			// Text Alignment
			// @see https://tailwindcss.com/docs/text-align
			"text-alignment": ClassGroup{ClassObject{"text": ClassGroup{"left", "center", "right", "justify", "start", "end"}}},
			// Placeholder Color
			// @deprecated since Tailwind CSS v3.0.0
			// @see https://v3.tailwindcss.com/docs/placeholder-color
			"placeholder-color": ClassGroup{ClassObject{"placeholder": scaleColor()}},
			// Text Color
			// @see https://tailwindcss.com/docs/text-color
			"text-color": ClassGroup{ClassObject{"text": scaleColor()}},
			// Text Decoration
			// @see https://tailwindcss.com/docs/text-decoration
			"text-decoration": ClassGroup{"underline", "overline", "line-through", "no-underline"},
			// Text Decoration Style
			// @see https://tailwindcss.com/docs/text-decoration-style
			"text-decoration-style": ClassGroup{ClassObject{"decoration": concat(scaleLineStyle(), ClassGroup{"wavy"})}},
			// Text Decoration Thickness
			// @see https://tailwindcss.com/docs/text-decoration-thickness
			"text-decoration-thickness": ClassGroup{
				ClassObject{
					"decoration": ClassGroup{
						IsNumber,
						"from-font",
						"auto",
						IsArbitraryVariable,
						IsArbitraryLength,
					},
				},
			},
			// Text Decoration Color
			// @see https://tailwindcss.com/docs/text-decoration-color
			"text-decoration-color": ClassGroup{ClassObject{"decoration": scaleColor()}},
			// Text Underline Offset
			// @see https://tailwindcss.com/docs/text-underline-offset
			"underline-offset": ClassGroup{
				ClassObject{"underline-offset": ClassGroup{IsNumber, "auto", IsArbitraryVariable, IsArbitraryValue}},
			},
			// Text Transform
			// @see https://tailwindcss.com/docs/text-transform
			"text-transform": ClassGroup{"uppercase", "lowercase", "capitalize", "normal-case"},
			// Text Overflow
			// @see https://tailwindcss.com/docs/text-overflow
			"text-overflow": ClassGroup{"truncate", "text-ellipsis", "text-clip"},
			// Text Wrap
			// @see https://tailwindcss.com/docs/text-wrap
			"text-wrap": ClassGroup{ClassObject{"text": ClassGroup{"wrap", "nowrap", "balance", "pretty"}}},
			// Text Indent
			// @see https://tailwindcss.com/docs/text-indent
			"indent": ClassGroup{ClassObject{"indent": scaleUnambiguousSpacing()}},
			// Tab Size
			// @see https://tailwindcss.com/docs/tab-size
			"tab-size": ClassGroup{ClassObject{"tab": ClassGroup{IsInteger, IsArbitraryVariable, IsArbitraryValue}}},
			// Vertical Alignment
			// @see https://tailwindcss.com/docs/vertical-align
			"vertical-align": ClassGroup{
				ClassObject{
					"align": ClassGroup{
						"baseline",
						"top",
						"middle",
						"bottom",
						"text-top",
						"text-bottom",
						"sub",
						"super",
						IsArbitraryVariable,
						IsArbitraryValue,
					},
				},
			},
			// Whitespace
			// @see https://tailwindcss.com/docs/whitespace
			"whitespace": ClassGroup{
				ClassObject{"whitespace": ClassGroup{"normal", "nowrap", "pre", "pre-line", "pre-wrap", "break-spaces"}},
			},
			// Word Break
			// @see https://tailwindcss.com/docs/word-break
			"break": ClassGroup{ClassObject{"break": ClassGroup{"normal", "words", "all", "keep"}}},
			// Overflow Wrap
			// @see https://tailwindcss.com/docs/overflow-wrap
			"wrap": ClassGroup{ClassObject{"wrap": ClassGroup{"break-word", "anywhere", "normal"}}},
			// Hyphens
			// @see https://tailwindcss.com/docs/hyphens
			"hyphens": ClassGroup{ClassObject{"hyphens": ClassGroup{"none", "manual", "auto"}}},
			// Content
			// @see https://tailwindcss.com/docs/content
			"content": ClassGroup{ClassObject{"content": ClassGroup{"none", IsArbitraryVariable, IsArbitraryValue}}},

			// -------------------
			// --- Backgrounds ---
			// -------------------

			// Background Attachment
			// @see https://tailwindcss.com/docs/background-attachment
			"bg-attachment": ClassGroup{ClassObject{"bg": ClassGroup{"fixed", "local", "scroll"}}},
			// Background Clip
			// @see https://tailwindcss.com/docs/background-clip
			"bg-clip": ClassGroup{ClassObject{"bg-clip": ClassGroup{"border", "padding", "content", "text"}}},
			// Background Origin
			// @see https://tailwindcss.com/docs/background-origin
			"bg-origin": ClassGroup{ClassObject{"bg-origin": ClassGroup{"border", "padding", "content"}}},
			// Background Position
			// @see https://tailwindcss.com/docs/background-position
			"bg-position": ClassGroup{ClassObject{"bg": scaleBgPosition()}},
			// Background Repeat
			// @see https://tailwindcss.com/docs/background-repeat
			"bg-repeat": ClassGroup{ClassObject{"bg": scaleBgRepeat()}},
			// Background Size
			// @see https://tailwindcss.com/docs/background-size
			"bg-size": ClassGroup{ClassObject{"bg": scaleBgSize()}},
			// Background Image
			// @see https://tailwindcss.com/docs/background-image
			"bg-image": ClassGroup{
				ClassObject{
					"bg": ClassGroup{
						"none",
						ClassObject{
							"linear": ClassGroup{
								ClassObject{"to": ClassGroup{"t", "tr", "r", "br", "b", "bl", "l", "tl"}},
								IsInteger,
								IsArbitraryVariable,
								IsArbitraryValue,
							},
							"radial": ClassGroup{"", IsArbitraryVariable, IsArbitraryValue},
							"conic":  ClassGroup{IsInteger, IsArbitraryVariable, IsArbitraryValue},
						},
						IsArbitraryVariableImage,
						IsArbitraryImage,
					},
				},
			},
			// Background Color
			// @see https://tailwindcss.com/docs/background-color
			"bg-color": ClassGroup{ClassObject{"bg": scaleColor()}},
			// Gradient Color Stops From Position
			// @see https://tailwindcss.com/docs/gradient-color-stops
			"gradient-from-pos": ClassGroup{ClassObject{"from": scaleGradientStopPosition()}},
			// Gradient Color Stops Via Position
			// @see https://tailwindcss.com/docs/gradient-color-stops
			"gradient-via-pos": ClassGroup{ClassObject{"via": scaleGradientStopPosition()}},
			// Gradient Color Stops To Position
			// @see https://tailwindcss.com/docs/gradient-color-stops
			"gradient-to-pos": ClassGroup{ClassObject{"to": scaleGradientStopPosition()}},
			// Gradient Color Stops From
			// @see https://tailwindcss.com/docs/gradient-color-stops
			"gradient-from": ClassGroup{ClassObject{"from": scaleColor()}},
			// Gradient Color Stops Via
			// @see https://tailwindcss.com/docs/gradient-color-stops
			"gradient-via": ClassGroup{ClassObject{"via": scaleColor()}},
			// Gradient Color Stops To
			// @see https://tailwindcss.com/docs/gradient-color-stops
			"gradient-to": ClassGroup{ClassObject{"to": scaleColor()}},

			// ---------------
			// --- Borders ---
			// ---------------

			// Border Radius
			// @see https://tailwindcss.com/docs/border-radius
			"rounded": ClassGroup{ClassObject{"rounded": scaleRadius()}},
			// Border Radius Start
			// @see https://tailwindcss.com/docs/border-radius
			"rounded-s": ClassGroup{ClassObject{"rounded-s": scaleRadius()}},
			// Border Radius End
			// @see https://tailwindcss.com/docs/border-radius
			"rounded-e": ClassGroup{ClassObject{"rounded-e": scaleRadius()}},
			// Border Radius Top
			// @see https://tailwindcss.com/docs/border-radius
			"rounded-t": ClassGroup{ClassObject{"rounded-t": scaleRadius()}},
			// Border Radius Right
			// @see https://tailwindcss.com/docs/border-radius
			"rounded-r": ClassGroup{ClassObject{"rounded-r": scaleRadius()}},
			// Border Radius Bottom
			// @see https://tailwindcss.com/docs/border-radius
			"rounded-b": ClassGroup{ClassObject{"rounded-b": scaleRadius()}},
			// Border Radius Left
			// @see https://tailwindcss.com/docs/border-radius
			"rounded-l": ClassGroup{ClassObject{"rounded-l": scaleRadius()}},
			// Border Radius Start Start
			// @see https://tailwindcss.com/docs/border-radius
			"rounded-ss": ClassGroup{ClassObject{"rounded-ss": scaleRadius()}},
			// Border Radius Start End
			// @see https://tailwindcss.com/docs/border-radius
			"rounded-se": ClassGroup{ClassObject{"rounded-se": scaleRadius()}},
			// Border Radius End End
			// @see https://tailwindcss.com/docs/border-radius
			"rounded-ee": ClassGroup{ClassObject{"rounded-ee": scaleRadius()}},
			// Border Radius End Start
			// @see https://tailwindcss.com/docs/border-radius
			"rounded-es": ClassGroup{ClassObject{"rounded-es": scaleRadius()}},
			// Border Radius Top Left
			// @see https://tailwindcss.com/docs/border-radius
			"rounded-tl": ClassGroup{ClassObject{"rounded-tl": scaleRadius()}},
			// Border Radius Top Right
			// @see https://tailwindcss.com/docs/border-radius
			"rounded-tr": ClassGroup{ClassObject{"rounded-tr": scaleRadius()}},
			// Border Radius Bottom Right
			// @see https://tailwindcss.com/docs/border-radius
			"rounded-br": ClassGroup{ClassObject{"rounded-br": scaleRadius()}},
			// Border Radius Bottom Left
			// @see https://tailwindcss.com/docs/border-radius
			"rounded-bl": ClassGroup{ClassObject{"rounded-bl": scaleRadius()}},
			// Border Width
			// @see https://tailwindcss.com/docs/border-width
			"border-w": ClassGroup{ClassObject{"border": scaleBorderWidth()}},
			// Border Width Inline
			// @see https://tailwindcss.com/docs/border-width
			"border-w-x": ClassGroup{ClassObject{"border-x": scaleBorderWidth()}},
			// Border Width Block
			// @see https://tailwindcss.com/docs/border-width
			"border-w-y": ClassGroup{ClassObject{"border-y": scaleBorderWidth()}},
			// Border Width Inline Start
			// @see https://tailwindcss.com/docs/border-width
			"border-w-s": ClassGroup{ClassObject{"border-s": scaleBorderWidth()}},
			// Border Width Inline End
			// @see https://tailwindcss.com/docs/border-width
			"border-w-e": ClassGroup{ClassObject{"border-e": scaleBorderWidth()}},
			// Border Width Block Start
			// @see https://tailwindcss.com/docs/border-width
			"border-w-bs": ClassGroup{ClassObject{"border-bs": scaleBorderWidth()}},
			// Border Width Block End
			// @see https://tailwindcss.com/docs/border-width
			"border-w-be": ClassGroup{ClassObject{"border-be": scaleBorderWidth()}},
			// Border Width Top
			// @see https://tailwindcss.com/docs/border-width
			"border-w-t": ClassGroup{ClassObject{"border-t": scaleBorderWidth()}},
			// Border Width Right
			// @see https://tailwindcss.com/docs/border-width
			"border-w-r": ClassGroup{ClassObject{"border-r": scaleBorderWidth()}},
			// Border Width Bottom
			// @see https://tailwindcss.com/docs/border-width
			"border-w-b": ClassGroup{ClassObject{"border-b": scaleBorderWidth()}},
			// Border Width Left
			// @see https://tailwindcss.com/docs/border-width
			"border-w-l": ClassGroup{ClassObject{"border-l": scaleBorderWidth()}},
			// Divide Width X
			// @see https://tailwindcss.com/docs/border-width#between-children
			"divide-x": ClassGroup{ClassObject{"divide-x": scaleBorderWidth()}},
			// Divide Width X Reverse
			// @see https://tailwindcss.com/docs/border-width#between-children
			"divide-x-reverse": ClassGroup{"divide-x-reverse"},
			// Divide Width Y
			// @see https://tailwindcss.com/docs/border-width#between-children
			"divide-y": ClassGroup{ClassObject{"divide-y": scaleBorderWidth()}},
			// Divide Width Y Reverse
			// @see https://tailwindcss.com/docs/border-width#between-children
			"divide-y-reverse": ClassGroup{"divide-y-reverse"},
			// Border Style
			// @see https://tailwindcss.com/docs/border-style
			"border-style": ClassGroup{ClassObject{"border": concat(scaleLineStyle(), ClassGroup{"hidden", "none"})}},
			// Divide Style
			// @see https://tailwindcss.com/docs/border-style#setting-the-divider-style
			"divide-style": ClassGroup{ClassObject{"divide": concat(scaleLineStyle(), ClassGroup{"hidden", "none"})}},
			// Border Color
			// @see https://tailwindcss.com/docs/border-color
			"border-color": ClassGroup{ClassObject{"border": scaleColor()}},
			// Border Color Inline
			// @see https://tailwindcss.com/docs/border-color
			"border-color-x": ClassGroup{ClassObject{"border-x": scaleColor()}},
			// Border Color Block
			// @see https://tailwindcss.com/docs/border-color
			"border-color-y": ClassGroup{ClassObject{"border-y": scaleColor()}},
			// Border Color Inline Start
			// @see https://tailwindcss.com/docs/border-color
			"border-color-s": ClassGroup{ClassObject{"border-s": scaleColor()}},
			// Border Color Inline End
			// @see https://tailwindcss.com/docs/border-color
			"border-color-e": ClassGroup{ClassObject{"border-e": scaleColor()}},
			// Border Color Block Start
			// @see https://tailwindcss.com/docs/border-color
			"border-color-bs": ClassGroup{ClassObject{"border-bs": scaleColor()}},
			// Border Color Block End
			// @see https://tailwindcss.com/docs/border-color
			"border-color-be": ClassGroup{ClassObject{"border-be": scaleColor()}},
			// Border Color Top
			// @see https://tailwindcss.com/docs/border-color
			"border-color-t": ClassGroup{ClassObject{"border-t": scaleColor()}},
			// Border Color Right
			// @see https://tailwindcss.com/docs/border-color
			"border-color-r": ClassGroup{ClassObject{"border-r": scaleColor()}},
			// Border Color Bottom
			// @see https://tailwindcss.com/docs/border-color
			"border-color-b": ClassGroup{ClassObject{"border-b": scaleColor()}},
			// Border Color Left
			// @see https://tailwindcss.com/docs/border-color
			"border-color-l": ClassGroup{ClassObject{"border-l": scaleColor()}},
			// Divide Color
			// @see https://tailwindcss.com/docs/divide-color
			"divide-color": ClassGroup{ClassObject{"divide": scaleColor()}},
			// Outline Style
			// @see https://tailwindcss.com/docs/outline-style
			"outline-style": ClassGroup{ClassObject{"outline": concat(scaleLineStyle(), ClassGroup{"none", "hidden"})}},
			// Outline Offset
			// @see https://tailwindcss.com/docs/outline-offset
			"outline-offset": ClassGroup{
				ClassObject{"outline-offset": ClassGroup{IsNumber, IsArbitraryVariable, IsArbitraryValue}},
			},
			// Outline Width
			// @see https://tailwindcss.com/docs/outline-width
			"outline-w": ClassGroup{
				ClassObject{"outline": ClassGroup{"", IsNumber, IsArbitraryVariableLength, IsArbitraryLength}},
			},
			// Outline Color
			// @see https://tailwindcss.com/docs/outline-color
			"outline-color": ClassGroup{ClassObject{"outline": scaleColor()}},

			// ---------------
			// --- Effects ---
			// ---------------

			// Box Shadow
			// @see https://tailwindcss.com/docs/box-shadow
			"shadow": ClassGroup{
				ClassObject{
					"shadow": ClassGroup{
						// Deprecated since Tailwind CSS v4.0.0
						"",
						"none",
						themeShadow,
						IsArbitraryVariableShadow,
						IsArbitraryShadow,
					},
				},
			},
			// Box Shadow Color
			// @see https://tailwindcss.com/docs/box-shadow#setting-the-shadow-color
			"shadow-color": ClassGroup{ClassObject{"shadow": scaleColor()}},
			// Inset Box Shadow
			// @see https://tailwindcss.com/docs/box-shadow#adding-an-inset-shadow
			"inset-shadow": ClassGroup{
				ClassObject{
					"inset-shadow": ClassGroup{
						"none",
						themeInsetShadow,
						IsArbitraryVariableShadow,
						IsArbitraryShadow,
					},
				},
			},
			// Inset Box Shadow Color
			// @see https://tailwindcss.com/docs/box-shadow#setting-the-inset-shadow-color
			"inset-shadow-color": ClassGroup{ClassObject{"inset-shadow": scaleColor()}},
			// Ring Width
			// @see https://tailwindcss.com/docs/box-shadow#adding-a-ring
			"ring-w": ClassGroup{ClassObject{"ring": scaleBorderWidth()}},
			// Ring Width Inset
			// @see https://v3.tailwindcss.com/docs/ring-width#inset-rings
			// @deprecated since Tailwind CSS v4.0.0
			// @see https://github.com/tailwindlabs/tailwindcss/blob/v4.0.0/packages/tailwindcss/src/utilities.ts#L4158
			"ring-w-inset": ClassGroup{"ring-inset"},
			// Ring Color
			// @see https://tailwindcss.com/docs/box-shadow#setting-the-ring-color
			"ring-color": ClassGroup{ClassObject{"ring": scaleColor()}},
			// Ring Offset Width
			// @see https://v3.tailwindcss.com/docs/ring-offset-width
			// @deprecated since Tailwind CSS v4.0.0
			// @see https://github.com/tailwindlabs/tailwindcss/blob/v4.0.0/packages/tailwindcss/src/utilities.ts#L4158
			"ring-offset-w": ClassGroup{ClassObject{"ring-offset": ClassGroup{IsNumber, IsArbitraryLength}}},
			// Ring Offset Color
			// @see https://v3.tailwindcss.com/docs/ring-offset-color
			// @deprecated since Tailwind CSS v4.0.0
			// @see https://github.com/tailwindlabs/tailwindcss/blob/v4.0.0/packages/tailwindcss/src/utilities.ts#L4158
			"ring-offset-color": ClassGroup{ClassObject{"ring-offset": scaleColor()}},
			// Inset Ring Width
			// @see https://tailwindcss.com/docs/box-shadow#adding-an-inset-ring
			"inset-ring-w": ClassGroup{ClassObject{"inset-ring": scaleBorderWidth()}},
			// Inset Ring Color
			// @see https://tailwindcss.com/docs/box-shadow#setting-the-inset-ring-color
			"inset-ring-color": ClassGroup{ClassObject{"inset-ring": scaleColor()}},
			// Text Shadow
			// @see https://tailwindcss.com/docs/text-shadow
			"text-shadow": ClassGroup{
				ClassObject{
					"text-shadow": ClassGroup{
						"none",
						themeTextShadow,
						IsArbitraryVariableShadow,
						IsArbitraryShadow,
					},
				},
			},
			// Text Shadow Color
			// @see https://tailwindcss.com/docs/text-shadow#setting-the-shadow-color
			"text-shadow-color": ClassGroup{ClassObject{"text-shadow": scaleColor()}},
			// Opacity
			// @see https://tailwindcss.com/docs/opacity
			"opacity": ClassGroup{ClassObject{"opacity": ClassGroup{IsNumber, IsArbitraryVariable, IsArbitraryValue}}},
			// Mix Blend Mode
			// @see https://tailwindcss.com/docs/mix-blend-mode
			"mix-blend": ClassGroup{ClassObject{"mix-blend": concat(scaleBlendMode(), ClassGroup{"plus-darker", "plus-lighter"})}},
			// Background Blend Mode
			// @see https://tailwindcss.com/docs/background-blend-mode
			"bg-blend": ClassGroup{ClassObject{"bg-blend": scaleBlendMode()}},
			// Mask Clip
			// @see https://tailwindcss.com/docs/mask-clip
			"mask-clip": ClassGroup{
				ClassObject{"mask-clip": ClassGroup{"border", "padding", "content", "fill", "stroke", "view"}},
				"mask-no-clip",
			},
			// Mask Composite
			// @see https://tailwindcss.com/docs/mask-composite
			"mask-composite": ClassGroup{ClassObject{"mask": ClassGroup{"add", "subtract", "intersect", "exclude"}}},
			// Mask Image
			// @see https://tailwindcss.com/docs/mask-image
			"mask-image-linear-pos":         ClassGroup{ClassObject{"mask-linear": ClassGroup{IsNumber}}},
			"mask-image-linear-from-pos":    ClassGroup{ClassObject{"mask-linear-from": scaleMaskImagePosition()}},
			"mask-image-linear-to-pos":      ClassGroup{ClassObject{"mask-linear-to": scaleMaskImagePosition()}},
			"mask-image-linear-from-color":  ClassGroup{ClassObject{"mask-linear-from": scaleColor()}},
			"mask-image-linear-to-color":    ClassGroup{ClassObject{"mask-linear-to": scaleColor()}},
			"mask-image-t-from-pos":         ClassGroup{ClassObject{"mask-t-from": scaleMaskImagePosition()}},
			"mask-image-t-to-pos":           ClassGroup{ClassObject{"mask-t-to": scaleMaskImagePosition()}},
			"mask-image-t-from-color":       ClassGroup{ClassObject{"mask-t-from": scaleColor()}},
			"mask-image-t-to-color":         ClassGroup{ClassObject{"mask-t-to": scaleColor()}},
			"mask-image-r-from-pos":         ClassGroup{ClassObject{"mask-r-from": scaleMaskImagePosition()}},
			"mask-image-r-to-pos":           ClassGroup{ClassObject{"mask-r-to": scaleMaskImagePosition()}},
			"mask-image-r-from-color":       ClassGroup{ClassObject{"mask-r-from": scaleColor()}},
			"mask-image-r-to-color":         ClassGroup{ClassObject{"mask-r-to": scaleColor()}},
			"mask-image-b-from-pos":         ClassGroup{ClassObject{"mask-b-from": scaleMaskImagePosition()}},
			"mask-image-b-to-pos":           ClassGroup{ClassObject{"mask-b-to": scaleMaskImagePosition()}},
			"mask-image-b-from-color":       ClassGroup{ClassObject{"mask-b-from": scaleColor()}},
			"mask-image-b-to-color":         ClassGroup{ClassObject{"mask-b-to": scaleColor()}},
			"mask-image-l-from-pos":         ClassGroup{ClassObject{"mask-l-from": scaleMaskImagePosition()}},
			"mask-image-l-to-pos":           ClassGroup{ClassObject{"mask-l-to": scaleMaskImagePosition()}},
			"mask-image-l-from-color":       ClassGroup{ClassObject{"mask-l-from": scaleColor()}},
			"mask-image-l-to-color":         ClassGroup{ClassObject{"mask-l-to": scaleColor()}},
			"mask-image-x-from-pos":         ClassGroup{ClassObject{"mask-x-from": scaleMaskImagePosition()}},
			"mask-image-x-to-pos":           ClassGroup{ClassObject{"mask-x-to": scaleMaskImagePosition()}},
			"mask-image-x-from-color":       ClassGroup{ClassObject{"mask-x-from": scaleColor()}},
			"mask-image-x-to-color":         ClassGroup{ClassObject{"mask-x-to": scaleColor()}},
			"mask-image-y-from-pos":         ClassGroup{ClassObject{"mask-y-from": scaleMaskImagePosition()}},
			"mask-image-y-to-pos":           ClassGroup{ClassObject{"mask-y-to": scaleMaskImagePosition()}},
			"mask-image-y-from-color":       ClassGroup{ClassObject{"mask-y-from": scaleColor()}},
			"mask-image-y-to-color":         ClassGroup{ClassObject{"mask-y-to": scaleColor()}},
			"mask-image-radial":             ClassGroup{ClassObject{"mask-radial": ClassGroup{IsArbitraryVariable, IsArbitraryValue}}},
			"mask-image-radial-from-pos":    ClassGroup{ClassObject{"mask-radial-from": scaleMaskImagePosition()}},
			"mask-image-radial-to-pos":      ClassGroup{ClassObject{"mask-radial-to": scaleMaskImagePosition()}},
			"mask-image-radial-from-color":  ClassGroup{ClassObject{"mask-radial-from": scaleColor()}},
			"mask-image-radial-to-color":    ClassGroup{ClassObject{"mask-radial-to": scaleColor()}},
			"mask-image-radial-shape":       ClassGroup{ClassObject{"mask-radial": ClassGroup{"circle", "ellipse"}}},
			"mask-image-radial-size": ClassGroup{
				ClassObject{
					"mask-radial": ClassGroup{
						ClassObject{
							"closest":  ClassGroup{"side", "corner"},
							"farthest": ClassGroup{"side", "corner"},
						},
					},
				},
			},
			"mask-image-radial-pos":        ClassGroup{ClassObject{"mask-radial-at": scalePosition()}},
			"mask-image-conic-pos":         ClassGroup{ClassObject{"mask-conic": ClassGroup{IsNumber}}},
			"mask-image-conic-from-pos":    ClassGroup{ClassObject{"mask-conic-from": scaleMaskImagePosition()}},
			"mask-image-conic-to-pos":      ClassGroup{ClassObject{"mask-conic-to": scaleMaskImagePosition()}},
			"mask-image-conic-from-color":  ClassGroup{ClassObject{"mask-conic-from": scaleColor()}},
			"mask-image-conic-to-color":    ClassGroup{ClassObject{"mask-conic-to": scaleColor()}},
			// Mask Mode
			// @see https://tailwindcss.com/docs/mask-mode
			"mask-mode": ClassGroup{ClassObject{"mask": ClassGroup{"alpha", "luminance", "match"}}},
			// Mask Origin
			// @see https://tailwindcss.com/docs/mask-origin
			"mask-origin": ClassGroup{
				ClassObject{"mask-origin": ClassGroup{"border", "padding", "content", "fill", "stroke", "view"}},
			},
			// Mask Position
			// @see https://tailwindcss.com/docs/mask-position
			"mask-position": ClassGroup{ClassObject{"mask": scaleBgPosition()}},
			// Mask Repeat
			// @see https://tailwindcss.com/docs/mask-repeat
			"mask-repeat": ClassGroup{ClassObject{"mask": scaleBgRepeat()}},
			// Mask Size
			// @see https://tailwindcss.com/docs/mask-size
			"mask-size": ClassGroup{ClassObject{"mask": scaleBgSize()}},
			// Mask Type
			// @see https://tailwindcss.com/docs/mask-type
			"mask-type": ClassGroup{ClassObject{"mask-type": ClassGroup{"alpha", "luminance"}}},
			// Mask Image
			// @see https://tailwindcss.com/docs/mask-image
			"mask-image": ClassGroup{ClassObject{"mask": ClassGroup{"none", IsArbitraryVariable, IsArbitraryValue}}},

			// ---------------
			// --- Filters ---
			// ---------------

			// Filter
			// @see https://tailwindcss.com/docs/filter
			"filter": ClassGroup{
				ClassObject{
					"filter": ClassGroup{
						// Deprecated since Tailwind CSS v3.0.0
						"",
						"none",
						IsArbitraryVariable,
						IsArbitraryValue,
					},
				},
			},
			// Blur
			// @see https://tailwindcss.com/docs/blur
			"blur": ClassGroup{ClassObject{"blur": scaleBlur()}},
			// Brightness
			// @see https://tailwindcss.com/docs/brightness
			"brightness": ClassGroup{ClassObject{"brightness": ClassGroup{IsNumber, IsArbitraryVariable, IsArbitraryValue}}},
			// Contrast
			// @see https://tailwindcss.com/docs/contrast
			"contrast": ClassGroup{ClassObject{"contrast": ClassGroup{IsNumber, IsArbitraryVariable, IsArbitraryValue}}},
			// Drop Shadow
			// @see https://tailwindcss.com/docs/drop-shadow
			"drop-shadow": ClassGroup{
				ClassObject{
					"drop-shadow": ClassGroup{
						// Deprecated since Tailwind CSS v4.0.0
						"",
						"none",
						themeDropShadow,
						IsArbitraryVariableShadow,
						IsArbitraryShadow,
					},
				},
			},
			// Drop Shadow Color
			// @see https://tailwindcss.com/docs/filter-drop-shadow#setting-the-shadow-color
			"drop-shadow-color": ClassGroup{ClassObject{"drop-shadow": scaleColor()}},
			// Grayscale
			// @see https://tailwindcss.com/docs/grayscale
			"grayscale": ClassGroup{ClassObject{"grayscale": ClassGroup{"", IsNumber, IsArbitraryVariable, IsArbitraryValue}}},
			// Hue Rotate
			// @see https://tailwindcss.com/docs/hue-rotate
			"hue-rotate": ClassGroup{ClassObject{"hue-rotate": ClassGroup{IsNumber, IsArbitraryVariable, IsArbitraryValue}}},
			// Invert
			// @see https://tailwindcss.com/docs/invert
			"invert": ClassGroup{ClassObject{"invert": ClassGroup{"", IsNumber, IsArbitraryVariable, IsArbitraryValue}}},
			// Saturate
			// @see https://tailwindcss.com/docs/saturate
			"saturate": ClassGroup{ClassObject{"saturate": ClassGroup{IsNumber, IsArbitraryVariable, IsArbitraryValue}}},
			// Sepia
			// @see https://tailwindcss.com/docs/sepia
			"sepia": ClassGroup{ClassObject{"sepia": ClassGroup{"", IsNumber, IsArbitraryVariable, IsArbitraryValue}}},
			// Backdrop Filter
			// @see https://tailwindcss.com/docs/backdrop-filter
			"backdrop-filter": ClassGroup{
				ClassObject{
					"backdrop-filter": ClassGroup{
						// Deprecated since Tailwind CSS v3.0.0
						"",
						"none",
						IsArbitraryVariable,
						IsArbitraryValue,
					},
				},
			},
			// Backdrop Blur
			// @see https://tailwindcss.com/docs/backdrop-blur
			"backdrop-blur": ClassGroup{ClassObject{"backdrop-blur": scaleBlur()}},
			// Backdrop Brightness
			// @see https://tailwindcss.com/docs/backdrop-brightness
			"backdrop-brightness": ClassGroup{
				ClassObject{"backdrop-brightness": ClassGroup{IsNumber, IsArbitraryVariable, IsArbitraryValue}},
			},
			// Backdrop Contrast
			// @see https://tailwindcss.com/docs/backdrop-contrast
			"backdrop-contrast": ClassGroup{
				ClassObject{"backdrop-contrast": ClassGroup{IsNumber, IsArbitraryVariable, IsArbitraryValue}},
			},
			// Backdrop Grayscale
			// @see https://tailwindcss.com/docs/backdrop-grayscale
			"backdrop-grayscale": ClassGroup{
				ClassObject{"backdrop-grayscale": ClassGroup{"", IsNumber, IsArbitraryVariable, IsArbitraryValue}},
			},
			// Backdrop Hue Rotate
			// @see https://tailwindcss.com/docs/backdrop-hue-rotate
			"backdrop-hue-rotate": ClassGroup{
				ClassObject{"backdrop-hue-rotate": ClassGroup{IsNumber, IsArbitraryVariable, IsArbitraryValue}},
			},
			// Backdrop Invert
			// @see https://tailwindcss.com/docs/backdrop-invert
			"backdrop-invert": ClassGroup{
				ClassObject{"backdrop-invert": ClassGroup{"", IsNumber, IsArbitraryVariable, IsArbitraryValue}},
			},
			// Backdrop Opacity
			// @see https://tailwindcss.com/docs/backdrop-opacity
			"backdrop-opacity": ClassGroup{
				ClassObject{"backdrop-opacity": ClassGroup{IsNumber, IsArbitraryVariable, IsArbitraryValue}},
			},
			// Backdrop Saturate
			// @see https://tailwindcss.com/docs/backdrop-saturate
			"backdrop-saturate": ClassGroup{
				ClassObject{"backdrop-saturate": ClassGroup{IsNumber, IsArbitraryVariable, IsArbitraryValue}},
			},
			// Backdrop Sepia
			// @see https://tailwindcss.com/docs/backdrop-sepia
			"backdrop-sepia": ClassGroup{
				ClassObject{"backdrop-sepia": ClassGroup{"", IsNumber, IsArbitraryVariable, IsArbitraryValue}},
			},

			// --------------
			// --- Tables ---
			// --------------

			// Border Collapse
			// @see https://tailwindcss.com/docs/border-collapse
			"border-collapse": ClassGroup{ClassObject{"border": ClassGroup{"collapse", "separate"}}},
			// Border Spacing
			// @see https://tailwindcss.com/docs/border-spacing
			"border-spacing": ClassGroup{ClassObject{"border-spacing": scaleUnambiguousSpacing()}},
			// Border Spacing X
			// @see https://tailwindcss.com/docs/border-spacing
			"border-spacing-x": ClassGroup{ClassObject{"border-spacing-x": scaleUnambiguousSpacing()}},
			// Border Spacing Y
			// @see https://tailwindcss.com/docs/border-spacing
			"border-spacing-y": ClassGroup{ClassObject{"border-spacing-y": scaleUnambiguousSpacing()}},
			// Table Layout
			// @see https://tailwindcss.com/docs/table-layout
			"table-layout": ClassGroup{ClassObject{"table": ClassGroup{"auto", "fixed"}}},
			// Caption Side
			// @see https://tailwindcss.com/docs/caption-side
			"caption": ClassGroup{ClassObject{"caption": ClassGroup{"top", "bottom"}}},

			// ---------------------------------
			// --- Transitions and Animation ---
			// ---------------------------------

			// Transition Property
			// @see https://tailwindcss.com/docs/transition-property
			"transition": ClassGroup{
				ClassObject{
					"transition": ClassGroup{
						"",
						"all",
						"colors",
						"opacity",
						"shadow",
						"transform",
						"none",
						IsArbitraryVariable,
						IsArbitraryValue,
					},
				},
			},
			// Transition Behavior
			// @see https://tailwindcss.com/docs/transition-behavior
			"transition-behavior": ClassGroup{ClassObject{"transition": ClassGroup{"normal", "discrete"}}},
			// Transition Duration
			// @see https://tailwindcss.com/docs/transition-duration
			"duration": ClassGroup{ClassObject{"duration": ClassGroup{IsNumber, "initial", IsArbitraryVariable, IsArbitraryValue}}},
			// Transition Timing Function
			// @see https://tailwindcss.com/docs/transition-timing-function
			"ease": ClassGroup{
				ClassObject{"ease": ClassGroup{"linear", "initial", themeEase, IsArbitraryVariable, IsArbitraryValue}},
			},
			// Transition Delay
			// @see https://tailwindcss.com/docs/transition-delay
			"delay": ClassGroup{ClassObject{"delay": ClassGroup{IsNumber, IsArbitraryVariable, IsArbitraryValue}}},
			// Animation
			// @see https://tailwindcss.com/docs/animation
			"animate": ClassGroup{ClassObject{"animate": ClassGroup{"none", themeAnimate, IsArbitraryVariable, IsArbitraryValue}}},

			// ------------------
			// --- Transforms ---
			// ------------------

			// Backface Visibility
			// @see https://tailwindcss.com/docs/backface-visibility
			"backface": ClassGroup{ClassObject{"backface": ClassGroup{"hidden", "visible"}}},
			// Perspective
			// @see https://tailwindcss.com/docs/perspective
			"perspective": ClassGroup{
				ClassObject{"perspective": ClassGroup{themePerspective, IsArbitraryVariable, IsArbitraryValue}},
			},
			// Perspective Origin
			// @see https://tailwindcss.com/docs/perspective-origin
			"perspective-origin": ClassGroup{ClassObject{"perspective-origin": scalePositionWithArbitrary()}},
			// Rotate
			// @see https://tailwindcss.com/docs/rotate
			"rotate": ClassGroup{ClassObject{"rotate": scaleRotate()}},
			// Rotate X
			// @see https://tailwindcss.com/docs/rotate
			"rotate-x": ClassGroup{ClassObject{"rotate-x": scaleRotate()}},
			// Rotate Y
			// @see https://tailwindcss.com/docs/rotate
			"rotate-y": ClassGroup{ClassObject{"rotate-y": scaleRotate()}},
			// Rotate Z
			// @see https://tailwindcss.com/docs/rotate
			"rotate-z": ClassGroup{ClassObject{"rotate-z": scaleRotate()}},
			// Scale
			// @see https://tailwindcss.com/docs/scale
			"scale": ClassGroup{ClassObject{"scale": scaleScale()}},
			// Scale X
			// @see https://tailwindcss.com/docs/scale
			"scale-x": ClassGroup{ClassObject{"scale-x": scaleScale()}},
			// Scale Y
			// @see https://tailwindcss.com/docs/scale
			"scale-y": ClassGroup{ClassObject{"scale-y": scaleScale()}},
			// Scale Z
			// @see https://tailwindcss.com/docs/scale
			"scale-z": ClassGroup{ClassObject{"scale-z": scaleScale()}},
			// Scale 3D
			// @see https://tailwindcss.com/docs/scale
			"scale-3d": ClassGroup{"scale-3d"},
			// Skew
			// @see https://tailwindcss.com/docs/skew
			"skew": ClassGroup{ClassObject{"skew": scaleSkew()}},
			// Skew X
			// @see https://tailwindcss.com/docs/skew
			"skew-x": ClassGroup{ClassObject{"skew-x": scaleSkew()}},
			// Skew Y
			// @see https://tailwindcss.com/docs/skew
			"skew-y": ClassGroup{ClassObject{"skew-y": scaleSkew()}},
			// Transform
			// @see https://tailwindcss.com/docs/transform
			"transform": ClassGroup{
				ClassObject{"transform": ClassGroup{IsArbitraryVariable, IsArbitraryValue, "", "none", "gpu", "cpu"}},
			},
			// Transform Origin
			// @see https://tailwindcss.com/docs/transform-origin
			"transform-origin": ClassGroup{ClassObject{"origin": scalePositionWithArbitrary()}},
			// Transform Style
			// @see https://tailwindcss.com/docs/transform-style
			"transform-style": ClassGroup{ClassObject{"transform": ClassGroup{"3d", "flat"}}},
			// Translate
			// @see https://tailwindcss.com/docs/translate
			"translate": ClassGroup{ClassObject{"translate": scaleTranslate()}},
			// Translate X
			// @see https://tailwindcss.com/docs/translate
			"translate-x": ClassGroup{ClassObject{"translate-x": scaleTranslate()}},
			// Translate Y
			// @see https://tailwindcss.com/docs/translate
			"translate-y": ClassGroup{ClassObject{"translate-y": scaleTranslate()}},
			// Translate Z
			// @see https://tailwindcss.com/docs/translate
			"translate-z": ClassGroup{ClassObject{"translate-z": scaleTranslate()}},
			// Translate None
			// @see https://tailwindcss.com/docs/translate
			"translate-none": ClassGroup{"translate-none"},
			// Zoom
			// @see https://tailwindcss.com/docs/zoom
			"zoom": ClassGroup{ClassObject{"zoom": ClassGroup{IsInteger, IsArbitraryVariable, IsArbitraryValue}}},

			// ---------------------
			// --- Interactivity ---
			// ---------------------

			// Accent Color
			// @see https://tailwindcss.com/docs/accent-color
			"accent": ClassGroup{ClassObject{"accent": scaleColor()}},
			// Appearance
			// @see https://tailwindcss.com/docs/appearance
			"appearance": ClassGroup{ClassObject{"appearance": ClassGroup{"none", "auto"}}},
			// Caret Color
			// @see https://tailwindcss.com/docs/just-in-time-mode#caret-color-utilities
			"caret-color": ClassGroup{ClassObject{"caret": scaleColor()}},
			// Color Scheme
			// @see https://tailwindcss.com/docs/color-scheme
			"color-scheme": ClassGroup{
				ClassObject{"scheme": ClassGroup{"normal", "dark", "light", "light-dark", "only-dark", "only-light"}},
			},
			// Cursor
			// @see https://tailwindcss.com/docs/cursor
			"cursor": ClassGroup{
				ClassObject{
					"cursor": ClassGroup{
						"auto",
						"default",
						"pointer",
						"wait",
						"text",
						"move",
						"help",
						"not-allowed",
						"none",
						"context-menu",
						"progress",
						"cell",
						"crosshair",
						"vertical-text",
						"alias",
						"copy",
						"no-drop",
						"grab",
						"grabbing",
						"all-scroll",
						"col-resize",
						"row-resize",
						"n-resize",
						"e-resize",
						"s-resize",
						"w-resize",
						"ne-resize",
						"nw-resize",
						"se-resize",
						"sw-resize",
						"ew-resize",
						"ns-resize",
						"nesw-resize",
						"nwse-resize",
						"zoom-in",
						"zoom-out",
						IsArbitraryVariable,
						IsArbitraryValue,
					},
				},
			},
			// Field Sizing
			// @see https://tailwindcss.com/docs/field-sizing
			"field-sizing": ClassGroup{ClassObject{"field-sizing": ClassGroup{"fixed", "content"}}},
			// Pointer Events
			// @see https://tailwindcss.com/docs/pointer-events
			"pointer-events": ClassGroup{ClassObject{"pointer-events": ClassGroup{"auto", "none"}}},
			// Resize
			// @see https://tailwindcss.com/docs/resize
			"resize": ClassGroup{ClassObject{"resize": ClassGroup{"none", "", "y", "x"}}},
			// Scroll Behavior
			// @see https://tailwindcss.com/docs/scroll-behavior
			"scroll-behavior": ClassGroup{ClassObject{"scroll": ClassGroup{"auto", "smooth"}}},
			// Scrollbar Thumb Color
			// @see https://tailwindcss.com/docs/scrollbar-color
			"scrollbar-thumb-color": ClassGroup{ClassObject{"scrollbar-thumb": scaleColor()}},
			// Scrollbar Track Color
			// @see https://tailwindcss.com/docs/scrollbar-color
			"scrollbar-track-color": ClassGroup{ClassObject{"scrollbar-track": scaleColor()}},
			// Scrollbar Gutter
			// @see https://tailwindcss.com/docs/scrollbar-gutter
			"scrollbar-gutter": ClassGroup{ClassObject{"scrollbar-gutter": ClassGroup{"auto", "stable", "both"}}},
			// Scrollbar Width
			// @see https://tailwindcss.com/docs/scrollbar-width
			"scrollbar-w": ClassGroup{ClassObject{"scrollbar": ClassGroup{"auto", "thin", "none"}}},
			// Scroll Margin
			// @see https://tailwindcss.com/docs/scroll-margin
			"scroll-m": ClassGroup{ClassObject{"scroll-m": scaleUnambiguousSpacing()}},
			// Scroll Margin Inline
			// @see https://tailwindcss.com/docs/scroll-margin
			"scroll-mx": ClassGroup{ClassObject{"scroll-mx": scaleUnambiguousSpacing()}},
			// Scroll Margin Block
			// @see https://tailwindcss.com/docs/scroll-margin
			"scroll-my": ClassGroup{ClassObject{"scroll-my": scaleUnambiguousSpacing()}},
			// Scroll Margin Inline Start
			// @see https://tailwindcss.com/docs/scroll-margin
			"scroll-ms": ClassGroup{ClassObject{"scroll-ms": scaleUnambiguousSpacing()}},
			// Scroll Margin Inline End
			// @see https://tailwindcss.com/docs/scroll-margin
			"scroll-me": ClassGroup{ClassObject{"scroll-me": scaleUnambiguousSpacing()}},
			// Scroll Margin Block Start
			// @see https://tailwindcss.com/docs/scroll-margin
			"scroll-mbs": ClassGroup{ClassObject{"scroll-mbs": scaleUnambiguousSpacing()}},
			// Scroll Margin Block End
			// @see https://tailwindcss.com/docs/scroll-margin
			"scroll-mbe": ClassGroup{ClassObject{"scroll-mbe": scaleUnambiguousSpacing()}},
			// Scroll Margin Top
			// @see https://tailwindcss.com/docs/scroll-margin
			"scroll-mt": ClassGroup{ClassObject{"scroll-mt": scaleUnambiguousSpacing()}},
			// Scroll Margin Right
			// @see https://tailwindcss.com/docs/scroll-margin
			"scroll-mr": ClassGroup{ClassObject{"scroll-mr": scaleUnambiguousSpacing()}},
			// Scroll Margin Bottom
			// @see https://tailwindcss.com/docs/scroll-margin
			"scroll-mb": ClassGroup{ClassObject{"scroll-mb": scaleUnambiguousSpacing()}},
			// Scroll Margin Left
			// @see https://tailwindcss.com/docs/scroll-margin
			"scroll-ml": ClassGroup{ClassObject{"scroll-ml": scaleUnambiguousSpacing()}},
			// Scroll Padding
			// @see https://tailwindcss.com/docs/scroll-padding
			"scroll-p": ClassGroup{ClassObject{"scroll-p": scaleUnambiguousSpacing()}},
			// Scroll Padding Inline
			// @see https://tailwindcss.com/docs/scroll-padding
			"scroll-px": ClassGroup{ClassObject{"scroll-px": scaleUnambiguousSpacing()}},
			// Scroll Padding Block
			// @see https://tailwindcss.com/docs/scroll-padding
			"scroll-py": ClassGroup{ClassObject{"scroll-py": scaleUnambiguousSpacing()}},
			// Scroll Padding Inline Start
			// @see https://tailwindcss.com/docs/scroll-padding
			"scroll-ps": ClassGroup{ClassObject{"scroll-ps": scaleUnambiguousSpacing()}},
			// Scroll Padding Inline End
			// @see https://tailwindcss.com/docs/scroll-padding
			"scroll-pe": ClassGroup{ClassObject{"scroll-pe": scaleUnambiguousSpacing()}},
			// Scroll Padding Block Start
			// @see https://tailwindcss.com/docs/scroll-padding
			"scroll-pbs": ClassGroup{ClassObject{"scroll-pbs": scaleUnambiguousSpacing()}},
			// Scroll Padding Block End
			// @see https://tailwindcss.com/docs/scroll-padding
			"scroll-pbe": ClassGroup{ClassObject{"scroll-pbe": scaleUnambiguousSpacing()}},
			// Scroll Padding Top
			// @see https://tailwindcss.com/docs/scroll-padding
			"scroll-pt": ClassGroup{ClassObject{"scroll-pt": scaleUnambiguousSpacing()}},
			// Scroll Padding Right
			// @see https://tailwindcss.com/docs/scroll-padding
			"scroll-pr": ClassGroup{ClassObject{"scroll-pr": scaleUnambiguousSpacing()}},
			// Scroll Padding Bottom
			// @see https://tailwindcss.com/docs/scroll-padding
			"scroll-pb": ClassGroup{ClassObject{"scroll-pb": scaleUnambiguousSpacing()}},
			// Scroll Padding Left
			// @see https://tailwindcss.com/docs/scroll-padding
			"scroll-pl": ClassGroup{ClassObject{"scroll-pl": scaleUnambiguousSpacing()}},
			// Scroll Snap Align
			// @see https://tailwindcss.com/docs/scroll-snap-align
			"snap-align": ClassGroup{ClassObject{"snap": ClassGroup{"start", "end", "center", "align-none"}}},
			// Scroll Snap Stop
			// @see https://tailwindcss.com/docs/scroll-snap-stop
			"snap-stop": ClassGroup{ClassObject{"snap": ClassGroup{"normal", "always"}}},
			// Scroll Snap Type
			// @see https://tailwindcss.com/docs/scroll-snap-type
			"snap-type": ClassGroup{ClassObject{"snap": ClassGroup{"none", "x", "y", "both"}}},
			// Scroll Snap Type Strictness
			// @see https://tailwindcss.com/docs/scroll-snap-type
			"snap-strictness": ClassGroup{ClassObject{"snap": ClassGroup{"mandatory", "proximity"}}},
			// Touch Action
			// @see https://tailwindcss.com/docs/touch-action
			"touch": ClassGroup{ClassObject{"touch": ClassGroup{"auto", "none", "manipulation"}}},
			// Touch Action X
			// @see https://tailwindcss.com/docs/touch-action
			"touch-x": ClassGroup{ClassObject{"touch-pan": ClassGroup{"x", "left", "right"}}},
			// Touch Action Y
			// @see https://tailwindcss.com/docs/touch-action
			"touch-y": ClassGroup{ClassObject{"touch-pan": ClassGroup{"y", "up", "down"}}},
			// Touch Action Pinch Zoom
			// @see https://tailwindcss.com/docs/touch-action
			"touch-pz": ClassGroup{"touch-pinch-zoom"},
			// User Select
			// @see https://tailwindcss.com/docs/user-select
			"select": ClassGroup{ClassObject{"select": ClassGroup{"none", "text", "all", "auto"}}},
			// Will Change
			// @see https://tailwindcss.com/docs/will-change
			"will-change": ClassGroup{
				ClassObject{
					"will-change": ClassGroup{
						"auto",
						"scroll",
						"contents",
						"transform",
						IsArbitraryVariable,
						IsArbitraryValue,
					},
				},
			},

			// -----------
			// --- SVG ---
			// -----------

			// Fill
			// @see https://tailwindcss.com/docs/fill
			"fill": ClassGroup{ClassObject{"fill": concat(ClassGroup{"none"}, scaleColor())}},
			// Stroke Width
			// @see https://tailwindcss.com/docs/stroke-width
			"stroke-w": ClassGroup{
				ClassObject{
					"stroke": ClassGroup{
						IsNumber,
						IsArbitraryVariableLength,
						IsArbitraryLength,
						IsArbitraryNumber,
					},
				},
			},
			// Stroke
			// @see https://tailwindcss.com/docs/stroke
			"stroke": ClassGroup{ClassObject{"stroke": concat(ClassGroup{"none"}, scaleColor())}},

			// ---------------------
			// --- Accessibility ---
			// ---------------------

			// Forced Color Adjust
			// @see https://tailwindcss.com/docs/forced-color-adjust
			"forced-color-adjust": ClassGroup{ClassObject{"forced-color-adjust": ClassGroup{"auto", "none"}}},
		},
		ConflictingClassGroups: ConflictingClassGroupsMap{
			"container-named": {"container-type"},
			"overflow":        {"overflow-x", "overflow-y"},
			"overscroll":      {"overscroll-x", "overscroll-y"},
			"inset": {
				"inset-x",
				"inset-y",
				"inset-bs",
				"inset-be",
				"start",
				"end",
				"top",
				"right",
				"bottom",
				"left",
			},
			"inset-x":   {"right", "left"},
			"inset-y":   {"top", "bottom"},
			"flex":      {"basis", "grow", "shrink"},
			"gap":       {"gap-x", "gap-y"},
			"p":         {"px", "py", "ps", "pe", "pbs", "pbe", "pt", "pr", "pb", "pl"},
			"px":        {"pr", "pl"},
			"py":        {"pt", "pb"},
			"m":         {"mx", "my", "ms", "me", "mbs", "mbe", "mt", "mr", "mb", "ml"},
			"mx":        {"mr", "ml"},
			"my":        {"mt", "mb"},
			"size":      {"w", "h"},
			"font-size": {"leading"},
			"fvn-normal": {
				"fvn-ordinal",
				"fvn-slashed-zero",
				"fvn-figure",
				"fvn-spacing",
				"fvn-fraction",
			},
			"fvn-ordinal":      {"fvn-normal"},
			"fvn-slashed-zero": {"fvn-normal"},
			"fvn-figure":       {"fvn-normal"},
			"fvn-spacing":      {"fvn-normal"},
			"fvn-fraction":     {"fvn-normal"},
			"line-clamp":       {"display", "overflow"},
			"rounded": {
				"rounded-s",
				"rounded-e",
				"rounded-t",
				"rounded-r",
				"rounded-b",
				"rounded-l",
				"rounded-ss",
				"rounded-se",
				"rounded-ee",
				"rounded-es",
				"rounded-tl",
				"rounded-tr",
				"rounded-br",
				"rounded-bl",
			},
			"rounded-s":      {"rounded-ss", "rounded-es"},
			"rounded-e":      {"rounded-se", "rounded-ee"},
			"rounded-t":      {"rounded-tl", "rounded-tr"},
			"rounded-r":      {"rounded-tr", "rounded-br"},
			"rounded-b":      {"rounded-br", "rounded-bl"},
			"rounded-l":      {"rounded-tl", "rounded-bl"},
			"border-spacing": {"border-spacing-x", "border-spacing-y"},
			"border-w": {
				"border-w-x",
				"border-w-y",
				"border-w-s",
				"border-w-e",
				"border-w-bs",
				"border-w-be",
				"border-w-t",
				"border-w-r",
				"border-w-b",
				"border-w-l",
			},
			"border-w-x": {"border-w-r", "border-w-l"},
			"border-w-y": {"border-w-t", "border-w-b"},
			"border-color": {
				"border-color-x",
				"border-color-y",
				"border-color-s",
				"border-color-e",
				"border-color-bs",
				"border-color-be",
				"border-color-t",
				"border-color-r",
				"border-color-b",
				"border-color-l",
			},
			"border-color-x": {"border-color-r", "border-color-l"},
			"border-color-y": {"border-color-t", "border-color-b"},
			"translate":      {"translate-x", "translate-y", "translate-none"},
			"translate-none": {"translate", "translate-x", "translate-y", "translate-z"},
			"scroll-m": {
				"scroll-mx",
				"scroll-my",
				"scroll-ms",
				"scroll-me",
				"scroll-mbs",
				"scroll-mbe",
				"scroll-mt",
				"scroll-mr",
				"scroll-mb",
				"scroll-ml",
			},
			"scroll-mx": {"scroll-mr", "scroll-ml"},
			"scroll-my": {"scroll-mt", "scroll-mb"},
			"scroll-p": {
				"scroll-px",
				"scroll-py",
				"scroll-ps",
				"scroll-pe",
				"scroll-pbs",
				"scroll-pbe",
				"scroll-pt",
				"scroll-pr",
				"scroll-pb",
				"scroll-pl",
			},
			"scroll-px": {"scroll-pr", "scroll-pl"},
			"scroll-py": {"scroll-pt", "scroll-pb"},
			"touch":     {"touch-x", "touch-y", "touch-pz"},
			"touch-x":   {"touch"},
			"touch-y":   {"touch"},
			"touch-pz":  {"touch"},
		},
		ConflictingClassGroupModifiers: ConflictingClassGroupsMap{
			"font-size": {"leading"},
		},
		// ClassGroupOrder lists ClassGroups keys in JS source order so that
		// shared trie-node validators resolve consistently. See
		// Config.ClassGroupOrder for the rationale.
		ClassGroupOrder: []string{
			"aspect", "container", "container-type", "container-named", "columns",
			"break-after", "break-before", "break-inside", "box-decoration", "box",
			"display", "sr", "float", "clear", "isolation", "object-fit", "object-position",
			"overflow", "overflow-x", "overflow-y", "overscroll", "overscroll-x", "overscroll-y",
			"position", "inset", "inset-x", "inset-y", "start", "end", "inset-bs", "inset-be",
			"top", "right", "bottom", "left", "visibility", "z",
			"basis", "flex-direction", "flex-wrap", "flex", "grow", "shrink", "order",
			"grid-cols", "col-start-end", "col-start", "col-end",
			"grid-rows", "row-start-end", "row-start", "row-end",
			"grid-flow", "auto-cols", "auto-rows",
			"gap", "gap-x", "gap-y",
			"justify-content", "justify-items", "justify-self",
			"align-content", "align-items", "align-self",
			"place-content", "place-items", "place-self",
			"p", "px", "py", "ps", "pe", "pbs", "pbe", "pt", "pr", "pb", "pl",
			"m", "mx", "my", "ms", "me", "mbs", "mbe", "mt", "mr", "mb", "ml",
			"space-x", "space-x-reverse", "space-y", "space-y-reverse",
			"size", "inline-size", "min-inline-size", "max-inline-size",
			"block-size", "min-block-size", "max-block-size",
			"w", "min-w", "max-w", "h", "min-h", "max-h",
			"font-size", "font-smoothing", "font-style", "font-weight", "font-stretch",
			"font-family", "font-features",
			"fvn-normal", "fvn-ordinal", "fvn-slashed-zero", "fvn-figure", "fvn-spacing", "fvn-fraction",
			"tracking", "line-clamp", "leading",
			"list-image", "list-style-position", "list-style-type",
			"text-alignment", "placeholder-color", "text-color",
			"text-decoration", "text-decoration-style", "text-decoration-thickness", "text-decoration-color",
			"underline-offset", "text-transform", "text-overflow", "text-wrap",
			"indent", "tab-size", "vertical-align", "whitespace", "break", "wrap", "hyphens", "content",
			"bg-attachment", "bg-clip", "bg-origin", "bg-position", "bg-repeat", "bg-size", "bg-image", "bg-color",
			"gradient-from-pos", "gradient-via-pos", "gradient-to-pos",
			"gradient-from", "gradient-via", "gradient-to",
			"rounded", "rounded-s", "rounded-e", "rounded-t", "rounded-r", "rounded-b", "rounded-l",
			"rounded-ss", "rounded-se", "rounded-ee", "rounded-es",
			"rounded-tl", "rounded-tr", "rounded-br", "rounded-bl",
			"border-w", "border-w-x", "border-w-y", "border-w-s", "border-w-e",
			"border-w-bs", "border-w-be", "border-w-t", "border-w-r", "border-w-b", "border-w-l",
			"divide-x", "divide-x-reverse", "divide-y", "divide-y-reverse",
			"border-style", "divide-style",
			"border-color", "border-color-x", "border-color-y",
			"border-color-s", "border-color-e", "border-color-bs", "border-color-be",
			"border-color-t", "border-color-r", "border-color-b", "border-color-l", "divide-color",
			"outline-style", "outline-offset", "outline-w", "outline-color",
			"shadow", "shadow-color", "inset-shadow", "inset-shadow-color",
			"ring-w", "ring-w-inset", "ring-color", "ring-offset-w", "ring-offset-color",
			"inset-ring-w", "inset-ring-color", "text-shadow", "text-shadow-color",
			"opacity", "mix-blend", "bg-blend",
			"mask-clip", "mask-composite",
			"mask-image-linear-pos", "mask-image-linear-from-pos", "mask-image-linear-to-pos",
			"mask-image-linear-from-color", "mask-image-linear-to-color",
			"mask-image-t-from-pos", "mask-image-t-to-pos", "mask-image-t-from-color", "mask-image-t-to-color",
			"mask-image-r-from-pos", "mask-image-r-to-pos", "mask-image-r-from-color", "mask-image-r-to-color",
			"mask-image-b-from-pos", "mask-image-b-to-pos", "mask-image-b-from-color", "mask-image-b-to-color",
			"mask-image-l-from-pos", "mask-image-l-to-pos", "mask-image-l-from-color", "mask-image-l-to-color",
			"mask-image-x-from-pos", "mask-image-x-to-pos", "mask-image-x-from-color", "mask-image-x-to-color",
			"mask-image-y-from-pos", "mask-image-y-to-pos", "mask-image-y-from-color", "mask-image-y-to-color",
			"mask-image-radial", "mask-image-radial-from-pos", "mask-image-radial-to-pos",
			"mask-image-radial-from-color", "mask-image-radial-to-color",
			"mask-image-radial-shape", "mask-image-radial-size", "mask-image-radial-pos",
			"mask-image-conic-pos", "mask-image-conic-from-pos", "mask-image-conic-to-pos",
			"mask-image-conic-from-color", "mask-image-conic-to-color",
			"mask-mode", "mask-origin", "mask-position", "mask-repeat", "mask-size", "mask-type", "mask-image",
			"filter", "blur", "brightness", "contrast",
			"drop-shadow", "drop-shadow-color",
			"grayscale", "hue-rotate", "invert", "saturate", "sepia",
			"backdrop-filter", "backdrop-blur", "backdrop-brightness", "backdrop-contrast",
			"backdrop-grayscale", "backdrop-hue-rotate", "backdrop-invert",
			"backdrop-opacity", "backdrop-saturate", "backdrop-sepia",
			"border-collapse", "border-spacing", "border-spacing-x", "border-spacing-y",
			"table-layout", "caption",
			"transition", "transition-behavior", "duration", "ease", "delay",
			"animate", "backface", "perspective", "perspective-origin",
			"rotate", "rotate-x", "rotate-y", "rotate-z",
			"scale", "scale-x", "scale-y", "scale-z", "scale-3d",
			"skew", "skew-x", "skew-y",
			"transform", "transform-origin", "transform-style",
			"translate", "translate-x", "translate-y", "translate-z", "translate-none",
			"zoom",
			"accent", "appearance", "caret-color", "color-scheme", "cursor",
			"field-sizing", "pointer-events", "resize", "scroll-behavior",
			"scrollbar-thumb-color", "scrollbar-track-color", "scrollbar-gutter", "scrollbar-w",
			"scroll-m", "scroll-mx", "scroll-my", "scroll-ms", "scroll-me",
			"scroll-mbs", "scroll-mbe", "scroll-mt", "scroll-mr", "scroll-mb", "scroll-ml",
			"scroll-p", "scroll-px", "scroll-py", "scroll-ps", "scroll-pe",
			"scroll-pbs", "scroll-pbe", "scroll-pt", "scroll-pr", "scroll-pb", "scroll-pl",
			"snap-align", "snap-stop", "snap-type", "snap-strictness",
			"touch", "touch-x", "touch-y", "touch-pz",
			"select", "will-change",
			"fill", "stroke-w", "stroke",
			"forced-color-adjust",
		},
		PostfixLookupClassGroups: []string{"container-type"},
		OrderSensitiveModifiers: []string{
			"*",
			"**",
			"after",
			"backdrop",
			"before",
			"details-content",
			"file",
			"first-letter",
			"first-line",
			"marker",
			"placeholder",
			"selection",
		},
	}
}
