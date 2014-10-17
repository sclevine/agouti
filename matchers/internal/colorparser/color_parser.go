package colorparser

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Color struct {
	R, G, B uint8
	A       float64
}

func (c Color) String() string {
	return fmt.Sprintf("Color{R:%d, G:%d, B:%d, A:%.2f}", c.R, c.G, c.B, c.A)
}

var (
	shortRGBHexRE    = regexp.MustCompile(`^#([0-9a-fA-F])([0-9a-fA-F])([0-9a-fA-F])$`)
	longRGBHexRE     = regexp.MustCompile(`^#([0-9a-fA-F]{2})([0-9a-fA-F]{2})([0-9a-fA-F]{2})$`)
	rgbIntegerRE     = regexp.MustCompile(`^rgb\(\s*([-0-9]+)\s*,\s*([-0-9]+)\s*,\s*([-0-9]+)\s*\)$`)
	rgbPercentageRE  = regexp.MustCompile(`^rgb\(\s*([-0-9.]+)%\s*,\s*([-0-9.]+)%\s*,\s*([-0-9.]+)%\s*\)$`)
	rgbaIntegerRE    = regexp.MustCompile(`^rgba\(\s*(-?[0-9]+)\s*,\s*(-?[0-9]+)\s*,\s*(-?[0-9]+)\s*,\s*(-?[0-9.]+)\s*\)$`)
	rgbaPercentageRE = regexp.MustCompile(`^rgba\(\s*(-?[0-9.]+)%\s*,\s*(-?[0-9.]+)%\s*,\s*(-?[0-9.]+)%\s*,\s*(-?[0-9.]+)\s*\)$`)
	hslRE            = regexp.MustCompile(`^hsl\(\s*(-?[0-9]+)\s*,\s*(-?[0-9.]+)%\s*,\s*(-?[0-9.]+)%\s*\)$`)
	hslaRE           = regexp.MustCompile(`^hsla\(\s*(-?[0-9]+)\s*,\s*(-?[0-9.]+)%\s*,\s*(-?[0-9.]+)%\s*,\s*(-?[0-9.]+)\s*\)$`)
)

func ParseCSSColor(color string) (Color, error) {
	color = strings.Trim(color, " ")
	rgba, ok := colorLookup[color]
	if ok {
		return rgba, nil
	}
	switch {
	case shortRGBHexRE.MatchString(color):
		return parseShortRGBHex(color)
	case longRGBHexRE.MatchString(color):
		return parseLongRGBHex(color)
	case rgbIntegerRE.MatchString(color):
		return parseRGBInteger(color)
	case rgbPercentageRE.MatchString(color):
		return parseRGBPercentage(color)
	case rgbaIntegerRE.MatchString(color):
		return parseRGBAInteger(color)
	case rgbaPercentageRE.MatchString(color):
		return parseRGBAPercentage(color)
	case hslRE.MatchString(color):
		return parseHSL(color)
	case hslaRE.MatchString(color):
		return parseHSLA(color)
	default:
		return Color{}, errors.New("unparseable color")
	}
}

func parseShortRGBHex(color string) (Color, error) {
	components := shortRGBHexRE.FindStringSubmatch(color)
	if len(components) != 4 {
		return Color{}, errors.New("invalid rgb hex")
	}
	r, err := strconv.ParseUint(components[1]+components[1], 16, 32)
	if err != nil {
		return Color{}, err
	}
	g, err := strconv.ParseUint(components[2]+components[2], 16, 32)
	if err != nil {
		return Color{}, err
	}
	b, err := strconv.ParseUint(components[3]+components[3], 16, 32)
	if err != nil {
		return Color{}, err
	}
	return Color{uint8(r), uint8(g), uint8(b), 1.0}, nil
}

func parseLongRGBHex(color string) (Color, error) {
	components := longRGBHexRE.FindStringSubmatch(color)
	if len(components) != 4 {
		return Color{}, errors.New("invalid rgb hex")
	}
	r, err := strconv.ParseUint(components[1], 16, 32)
	if err != nil {
		return Color{}, err
	}
	g, err := strconv.ParseUint(components[2], 16, 32)
	if err != nil {
		return Color{}, err
	}
	b, err := strconv.ParseUint(components[3], 16, 32)
	if err != nil {
		return Color{}, err
	}
	return Color{uint8(r), uint8(g), uint8(b), 1.0}, nil
}

func parseRGBInteger(color string) (Color, error) {
	components := rgbIntegerRE.FindStringSubmatch(color)
	if len(components) != 4 {
		return Color{}, errors.New("invalid rgb")
	}
	r, err := strconv.ParseInt(components[1], 10, 64)
	if err != nil {
		return Color{}, err
	}
	g, err := strconv.ParseInt(components[2], 10, 64)
	if err != nil {
		return Color{}, err
	}
	b, err := strconv.ParseInt(components[3], 10, 64)
	if err != nil {
		return Color{}, err
	}
	return Color{
		clamp255(r),
		clamp255(g),
		clamp255(b),
		1.0,
	}, nil
}

func parseRGBPercentage(color string) (Color, error) {
	components := rgbPercentageRE.FindStringSubmatch(color)
	if len(components) != 4 {
		return Color{}, errors.New("invalid rgb percentage")
	}
	r, err := strconv.ParseFloat(components[1], 64)
	if err != nil {
		return Color{}, err
	}
	g, err := strconv.ParseFloat(components[2], 64)
	if err != nil {
		return Color{}, err
	}
	b, err := strconv.ParseFloat(components[3], 64)
	if err != nil {
		return Color{}, err
	}
	return Color{
		round255(r / 100.0 * 255),
		round255(g / 100.0 * 255),
		round255(b / 100.0 * 255),
		1.0,
	}, nil
}

func parseRGBAInteger(color string) (Color, error) {
	components := rgbaIntegerRE.FindStringSubmatch(color)
	if len(components) != 5 {
		return Color{}, errors.New("invalid rgba")
	}
	r, err := strconv.ParseInt(components[1], 10, 64)
	if err != nil {
		return Color{}, err
	}
	g, err := strconv.ParseInt(components[2], 10, 64)
	if err != nil {
		return Color{}, err
	}
	b, err := strconv.ParseInt(components[3], 10, 64)
	if err != nil {
		return Color{}, err
	}
	a, err := strconv.ParseFloat(components[4], 64)
	if err != nil {
		return Color{}, err
	}
	return Color{
		clamp255(r),
		clamp255(g),
		clamp255(b),
		clamp1(a),
	}, nil
}

func parseRGBAPercentage(color string) (Color, error) {
	components := rgbaPercentageRE.FindStringSubmatch(color)
	if len(components) != 5 {
		return Color{}, errors.New("invalid rgb percentage")
	}
	r, err := strconv.ParseFloat(components[1], 64)
	if err != nil {
		return Color{}, err
	}
	g, err := strconv.ParseFloat(components[2], 64)
	if err != nil {
		return Color{}, err
	}
	b, err := strconv.ParseFloat(components[3], 64)
	if err != nil {
		return Color{}, err
	}
	a, err := strconv.ParseFloat(components[4], 64)
	if err != nil {
		return Color{}, err
	}
	return Color{
		round255(r / 100.0 * 255),
		round255(g / 100.0 * 255),
		round255(b / 100.0 * 255),
		clamp1(a),
	}, nil
}

func parseHSL(color string) (Color, error) {
	components := hslRE.FindStringSubmatch(color)
	if len(components) != 4 {
		return Color{}, errors.New("invalid hsl percentage")
	}
	h, err := strconv.ParseInt(components[1], 10, 64)
	if err != nil {
		return Color{}, err
	}
	s, err := strconv.ParseFloat(components[2], 64)
	if err != nil {
		return Color{}, err
	}
	l, err := strconv.ParseFloat(components[3], 64)
	if err != nil {
		return Color{}, err
	}

	return colorFromHSL(h, s, l, 1.0), nil
}

func parseHSLA(color string) (Color, error) {
	components := hslaRE.FindStringSubmatch(color)
	if len(components) != 5 {
		return Color{}, errors.New("invalid hsl percentage")
	}
	h, err := strconv.ParseInt(components[1], 10, 64)
	if err != nil {
		return Color{}, err
	}
	s, err := strconv.ParseFloat(components[2], 64)
	if err != nil {
		return Color{}, err
	}
	l, err := strconv.ParseFloat(components[3], 64)
	if err != nil {
		return Color{}, err
	}
	a, err := strconv.ParseFloat(components[4], 64)
	if err != nil {
		return Color{}, err
	}

	return colorFromHSL(h, s, l, clamp1(a)), nil
}

func clamp255(value int64) uint8 {
	if value < 0 {
		return 0
	}
	if value > 255 {
		return 255
	}

	return uint8(value)
}

func round255(value float64) uint8 {
	value = value + 0.5 //round!
	if value < 0 {
		return 0
	}
	if value > 255 {
		return 255
	}

	return uint8(value)
}

func clamp1(value float64) float64 {
	if value < 0 {
		return 0
	}
	if value > 1 {
		return 1
	}
	return value
}

func colorFromHSL(hDegrees int64, s float64, l float64, a float64) Color {
	//see http://www.w3.org/TR/css3-color
	h := float64((((hDegrees % 360) + 360) % 360)) / 360.0
	s = clamp1(s / 100.0)
	l = clamp1(l / 100.0)

	var m2 float64
	if l <= 0.5 {
		m2 = l * (s + 1)
	} else {
		m2 = l + s - l*s
	}
	m1 := l*2 - m2

	return Color{
		R: hueToRGB(m1, m2, h+1.0/3.0),
		G: hueToRGB(m1, m2, h),
		B: hueToRGB(m1, m2, h-1.0/3.0),
		A: a,
	}
}

func hueToRGB(m1 float64, m2 float64, h float64) uint8 {
	if h < 0 {
		h = h + 1
	}
	if h > 1 {
		h = h - 1
	}
	if h*6 < 1 {
		return uint8(255 * (m1 + (m2-m1)*h*6.0))
	}
	if h*2 < 1 {
		return uint8(255 * m2)
	}
	if h*3 < 2 {
		return uint8(255 * (m1 + (m2-m1)*(2.0/3.0-h)*6.0))
	}
	return uint8(255 * m1)
}

var colorLookup = map[string]Color{
	"aliceblue":            Color{240, 248, 255, 1.0},
	"antiquewhite":         Color{250, 235, 215, 1.0},
	"aqua":                 Color{0, 255, 255, 1.0},
	"aquamarine":           Color{127, 255, 212, 1.0},
	"azure":                Color{240, 255, 255, 1.0},
	"beige":                Color{245, 245, 220, 1.0},
	"bisque":               Color{255, 228, 196, 1.0},
	"black":                Color{0, 0, 0, 1.0},
	"blanchedalmond":       Color{255, 235, 205, 1.0},
	"blue":                 Color{0, 0, 255, 1.0},
	"blueviolet":           Color{138, 43, 226, 1.0},
	"brown":                Color{165, 42, 42, 1.0},
	"burlywood":            Color{222, 184, 135, 1.0},
	"cadetblue":            Color{95, 158, 160, 1.0},
	"chartreuse":           Color{127, 255, 0, 1.0},
	"chocolate":            Color{210, 105, 30, 1.0},
	"coral":                Color{255, 127, 80, 1.0},
	"cornflowerblue":       Color{100, 149, 237, 1.0},
	"cornsilk":             Color{255, 248, 220, 1.0},
	"crimson":              Color{220, 20, 60, 1.0},
	"cyan":                 Color{0, 255, 255, 1.0},
	"darkblue":             Color{0, 0, 139, 1.0},
	"darkcyan":             Color{0, 139, 139, 1.0},
	"darkgoldenrod":        Color{184, 134, 11, 1.0},
	"darkgray":             Color{169, 169, 169, 1.0},
	"darkgreen":            Color{0, 100, 0, 1.0},
	"darkgrey":             Color{169, 169, 169, 1.0},
	"darkkhaki":            Color{189, 183, 107, 1.0},
	"darkmagenta":          Color{139, 0, 139, 1.0},
	"darkolivegreen":       Color{85, 107, 47, 1.0},
	"darkorange":           Color{255, 140, 0, 1.0},
	"darkorchid":           Color{153, 50, 204, 1.0},
	"darkred":              Color{139, 0, 0, 1.0},
	"darksalmon":           Color{233, 150, 122, 1.0},
	"darkseagreen":         Color{143, 188, 143, 1.0},
	"darkslateblue":        Color{72, 61, 139, 1.0},
	"darkslategray":        Color{47, 79, 79, 1.0},
	"darkslategrey":        Color{47, 79, 79, 1.0},
	"darkturquoise":        Color{0, 206, 209, 1.0},
	"darkviolet":           Color{148, 0, 211, 1.0},
	"deeppink":             Color{255, 20, 147, 1.0},
	"deepskyblue":          Color{0, 191, 255, 1.0},
	"dimgray":              Color{105, 105, 105, 1.0},
	"dimgrey":              Color{105, 105, 105, 1.0},
	"dodgerblue":           Color{30, 144, 255, 1.0},
	"firebrick":            Color{178, 34, 34, 1.0},
	"floralwhite":          Color{255, 250, 240, 1.0},
	"forestgreen":          Color{34, 139, 34, 1.0},
	"fuchsia":              Color{255, 0, 255, 1.0},
	"gainsboro":            Color{220, 220, 220, 1.0},
	"ghostwhite":           Color{248, 248, 255, 1.0},
	"gold":                 Color{255, 215, 0, 1.0},
	"goldenrod":            Color{218, 165, 32, 1.0},
	"gray":                 Color{128, 128, 128, 1.0},
	"green":                Color{0, 128, 0, 1.0},
	"greenyellow":          Color{173, 255, 47, 1.0},
	"grey":                 Color{128, 128, 128, 1.0},
	"honeydew":             Color{240, 255, 240, 1.0},
	"hotpink":              Color{255, 105, 180, 1.0},
	"indianred":            Color{205, 92, 92, 1.0},
	"indigo":               Color{75, 0, 130, 1.0},
	"ivory":                Color{255, 255, 240, 1.0},
	"khaki":                Color{240, 230, 140, 1.0},
	"lavender":             Color{230, 230, 250, 1.0},
	"lavenderblush":        Color{255, 240, 245, 1.0},
	"lawngreen":            Color{124, 252, 0, 1.0},
	"lemonchiffon":         Color{255, 250, 205, 1.0},
	"lightblue":            Color{173, 216, 230, 1.0},
	"lightcoral":           Color{240, 128, 128, 1.0},
	"lightcyan":            Color{224, 255, 255, 1.0},
	"lightgoldenrodyellow": Color{250, 250, 210, 1.0},
	"lightgray":            Color{211, 211, 211, 1.0},
	"lightgreen":           Color{144, 238, 144, 1.0},
	"lightgrey":            Color{211, 211, 211, 1.0},
	"lightpink":            Color{255, 182, 193, 1.0},
	"lightsalmon":          Color{255, 160, 122, 1.0},
	"lightseagreen":        Color{32, 178, 170, 1.0},
	"lightskyblue":         Color{135, 206, 250, 1.0},
	"lightslategray":       Color{119, 136, 153, 1.0},
	"lightslategrey":       Color{119, 136, 153, 1.0},
	"lightsteelblue":       Color{176, 196, 222, 1.0},
	"lightyellow":          Color{255, 255, 224, 1.0},
	"lime":                 Color{0, 255, 0, 1.0},
	"limegreen":            Color{50, 205, 50, 1.0},
	"linen":                Color{250, 240, 230, 1.0},
	"magenta":              Color{255, 0, 255, 1.0},
	"maroon":               Color{128, 0, 0, 1.0},
	"mediumaquamarine":     Color{102, 205, 170, 1.0},
	"mediumblue":           Color{0, 0, 205, 1.0},
	"mediumorchid":         Color{186, 85, 211, 1.0},
	"mediumpurple":         Color{147, 112, 219, 1.0},
	"mediumseagreen":       Color{60, 179, 113, 1.0},
	"mediumslateblue":      Color{123, 104, 238, 1.0},
	"mediumspringgreen":    Color{0, 250, 154, 1.0},
	"mediumturquoise":      Color{72, 209, 204, 1.0},
	"mediumvioletred":      Color{199, 21, 133, 1.0},
	"midnightblue":         Color{25, 25, 112, 1.0},
	"mintcream":            Color{245, 255, 250, 1.0},
	"mistyrose":            Color{255, 228, 225, 1.0},
	"moccasin":             Color{255, 228, 181, 1.0},
	"navajowhite":          Color{255, 222, 173, 1.0},
	"navy":                 Color{0, 0, 128, 1.0},
	"oldlace":              Color{253, 245, 230, 1.0},
	"olive":                Color{128, 128, 0, 1.0},
	"olivedrab":            Color{107, 142, 35, 1.0},
	"orange":               Color{255, 165, 0, 1.0},
	"orangered":            Color{255, 69, 0, 1.0},
	"orchid":               Color{218, 112, 214, 1.0},
	"palegoldenrod":        Color{238, 232, 170, 1.0},
	"palegreen":            Color{152, 251, 152, 1.0},
	"paleturquoise":        Color{175, 238, 238, 1.0},
	"palevioletred":        Color{219, 112, 147, 1.0},
	"papayawhip":           Color{255, 239, 213, 1.0},
	"peachpuff":            Color{255, 218, 185, 1.0},
	"peru":                 Color{205, 133, 63, 1.0},
	"pink":                 Color{255, 192, 203, 1.0},
	"plum":                 Color{221, 160, 221, 1.0},
	"powderblue":           Color{176, 224, 230, 1.0},
	"purple":               Color{128, 0, 128, 1.0},
	"red":                  Color{255, 0, 0, 1.0},
	"rosybrown":            Color{188, 143, 143, 1.0},
	"royalblue":            Color{65, 105, 225, 1.0},
	"saddlebrown":          Color{139, 69, 19, 1.0},
	"salmon":               Color{250, 128, 114, 1.0},
	"sandybrown":           Color{244, 164, 96, 1.0},
	"seagreen":             Color{46, 139, 87, 1.0},
	"seashell":             Color{255, 245, 238, 1.0},
	"sienna":               Color{160, 82, 45, 1.0},
	"silver":               Color{192, 192, 192, 1.0},
	"skyblue":              Color{135, 206, 235, 1.0},
	"slateblue":            Color{106, 90, 205, 1.0},
	"slategray":            Color{112, 128, 144, 1.0},
	"slategrey":            Color{112, 128, 144, 1.0},
	"snow":                 Color{255, 250, 250, 1.0},
	"springgreen":          Color{0, 255, 127, 1.0},
	"steelblue":            Color{70, 130, 180, 1.0},
	"tan":                  Color{210, 180, 140, 1.0},
	"teal":                 Color{0, 128, 128, 1.0},
	"thistle":              Color{216, 191, 216, 1.0},
	"tomato":               Color{255, 99, 71, 1.0},
	"turquoise":            Color{64, 224, 208, 1.0},
	"violet":               Color{238, 130, 238, 1.0},
	"wheat":                Color{245, 222, 179, 1.0},
	"white":                Color{255, 255, 255, 1.0},
	"whitesmoke":           Color{245, 245, 245, 1.0},
	"yellow":               Color{255, 255, 0, 1.0},
	"yellowgreen":          Color{154, 205, 50, 1.0},
}
