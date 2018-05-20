/*  Copyright (C) 2018 Alexandru Cojocaru

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>. */

package factoid

import (
	"regexp"
	"strings"
)

// S - subject
// A - attribute
// L - location

var simplePatterns = []string{
	"S's A",                // human eye's color
	"(the )?A of (the )?S", // the temperature of the sun
	"how A is S",           // how old is Richard Stallman (TODO: isAdjective(Y))
	"when is S in L",       // when is fashion week in NYC

}

var simplePatternsRegexp []*regexp.Regexp

func init() {
	for _, p := range simplePatterns {
		p = strings.Replace(p, "S", "(?P<subject>.+)", -1)
		p = strings.Replace(p, "A", "(?P<attribute>[[:alpha:]]+)", -1)
		p = "^" + p + "$"
		simplePatternsRegexp = append(simplePatternsRegexp, regexp.MustCompile(p))
	}
}

func extract(s string, r *regexp.Regexp, n string) string {
	m := r.FindStringSubmatchIndex(s)
	if m == nil {
		return ""
	}
	return string(r.ExpandString(nil, "$"+n, s, m))
}

func simplePattern(s string) *Meaning {
	for _, r := range simplePatternsRegexp {
		// TODO: check attribute and subject
		a := extract(s, r, "attribute")
		j := extract(s, r, "subject")
		if j != "" && a != "" {
			return &Meaning{Subject: j, Attribute: a}
		}
	}
	return nil
}
