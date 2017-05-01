// Written by https://xojoc.pw. Apache 2.0 license. No warranty.

// Package units contains functions to convert between units.
package units // import "xojoc.pw/nlp/units"

import (
	"bufio"
	"fmt"
	"log"
	"strconv"
	"strings"

	"xojoc.pw/must"
)

// TODO: use bignum

// https://github.com/ryantenney/gnu-units/blob/master/units.dat
var unitsPrefixes = map[string]float64{
	"yotta": 1e24,
	"Y":     1e24,
	"zetta": 1e21,
	"Z":     1e21,
	"exa":   1e18,
	"E":     1e18,
	"peta":  1e15,
	"P":     1e15,
	"tera":  1e12,
	"T":     1e12,
	"giga":  1e9,
	"G":     1e9,
	"mega":  1e6,
	"M":     1e6,
	"kilo":  1e3,
	"k":     1e3,
	"hecto": 1e2,
	"h":     1e2,
	"deka":  1e1,
	"da":    1e1,
	"":      1,
	"deci":  1e-1,
	"d":     1e-1,
	"centi": 1e-2,
	"c":     1e-2,
	"milli": 1e-3,
	"m":     1e-3,
	"micro": 1e-6,
	"µ":     1e-6,
	"nano":  1e-9,
	"n":     1e-9,
	"pico":  1e-12,
	"p":     1e-12,
	"femto": 1e-15,
	"f":     1e-15,
	"atto":  1e-18,
	"a":     1e-18,
	"zepto": 1e-21,
	"z":     1e-21,
	"yocto": 1e-24,
	"y":     1e-24,
}

type unit struct {
	value     float64
	primitive string
}

const units = `
kg		!
m		!
meter		m
meters		m
inch		2.54 cm
inches		inch
in		inch
foot		12 inch
feet		foot
ft		foot
yard		3 ft
yd		yard
mile		5280 ft
`

var allUnits = map[string]unit{}

func init() {
	scanner := bufio.NewScanner(strings.NewReader(units))
	for scanner.Scan() {
		fs := strings.Fields(scanner.Text())
		if len(fs) == 0 || fs[0][0] == '#' {
			continue
		}
		u := unit{1, fs[0]}
		if fs[1] != "!" {
			if len(fs) == 3 {
				v, err := strconv.ParseFloat(fs[1], 64)
				must.OK(err)
				u = allUnits[fs[2]]
				u.value *= v
			} else {
				u = allUnits[fs[1]]
			}
		}
		for pk, pv := range unitsPrefixes {
			allUnits[pk+fs[0]] = unit{pv * u.value, u.primitive}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func toKelvin(fnum float64, f string) float64 {
	switch f {
	case "K":
		return fnum
	case "°C":
		return fnum + 273.15
	case "°F":
		return (fnum + 459.67) * 5 / 9
	default:
		panic("unknown temperature unit: " + f)
	}
}

func fromKelvin(fnum float64, t string) float64 {
	switch t {
	case "K":
		return fnum
	case "°C":
		return fnum - 273.15
	case "°F":
		return fnum*9/5 - 459.67
	default:
		panic("unknown temperature unit: " + t)
	}
}

func convertTemperature(fnum float64, f string, t string) (float64, error) {
	aliases := map[string]string{
		"°C":         "°C",
		"K":          "K",
		"°F":         "°F",
		"celsius":    "°C",
		"kelvin":     "K",
		"fahrenheit": "°F",
	}
	var ok bool
	f, ok = aliases[f]
	if !ok {
		return 0, fmt.Errorf("unit %q not known", f)
	}
	t, ok = aliases[t]
	if !ok {
		return 0, fmt.Errorf("unit %q not known", t)
	}
	return fromKelvin(toKelvin(fnum, f), t), nil
}

// Convert converts one unit to another.
// Returns an error if units are not compatible.
func Convert(fnum float64, f string, t string) (float64, error) {
	fu, ok := allUnits[f]
	if !ok {
		return convertTemperature(fnum, f, t)
	}
	fu.value *= fnum
	tu, ok := allUnits[t]
	if !ok {
		return 0, fmt.Errorf("unit %q not known", t)
	}
	if fu.primitive != tu.primitive {
		return 0, fmt.Errorf("cannot convert %q to %q", fu.primitive, tu.primitive)
	}
	return fu.value / tu.value, nil
}
