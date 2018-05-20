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

package porter2english

import (
	"bytes"

	. "xojoc.pw/nlp/stem/internal/porter2"
)

// http://snowballstem.org/algorithms/english/stemmer.html

const vowels = "aeiouy"

var vowelsWXY = vowels + "wxY"

var r1Bare = R1(vowels)

/*
func r1_bare(s []byte) []byte {
	for len(s) > 0 && !isVowel(s[0]) {
		s = s[1:]
	}
	for len(s) > 0 && isVowel(s[0]) {
		s = s[1:]
	}
	if len(s) > 0 {
		s = s[1:]
	}
	return s
}
*/

var r1Exceptions = [][]byte{[]byte("gener"), []byte("commun"), []byte("arsen")}

func r1(s []byte) []byte {
	for _, p := range r1Exceptions {
		if bytes.HasPrefix(s, p) {
			return s[len(p):]
		}
	}
	return r1Bare(s)

}

func r2(s []byte) []byte {
	return r1Bare(r1(s))
}

func hasSuffix(s []byte, suffix string) bool {
	if len(s) < len(suffix) {
		return false
	}
	d := len(s) - len(suffix)
	for i := 0; i < len(suffix); i++ {
		if s[d+i] != suffix[i] {
			return false
		}
	}
	return true
}

func anySuffix(s []byte, suffixes ...string) bool {
	for _, x := range suffixes {
		if hasSuffix(s, x) {
			return true
		}
	}
	return false
}

func containsAny(s []byte, els string) bool {
	for i := 0; i < len(els); i++ {
		for j := range s {
			if els[i] == s[j] {
				return true
			}
		}
	}
	return false
}

func isAny(b byte, s string) bool {
	for i := 0; i < len(s); i++ {
		if b == s[i] {
			return true
		}
	}
	return false
}

func isVowel(b byte) bool {
	return isAny(b, vowels)
}

func hasDouble(s []byte) bool {
	return anySuffix(s, "bb", "dd", "ff", "gg", "mm", "nn", "pp", "rr", "tt")
}

func endsShortSyllable(s []byte) bool {
	if len(s) < 2 {
		return false
	} else if len(s) == 2 && isVowel(s[0]) && !isVowel(s[1]) {
		return true
	}
	i := len(s) - 2
	return isVowel(s[i]) && !isAny(s[i+1], vowelsWXY) && !isVowel(s[i-1])
}

func isShort(s []byte) bool {
	return len(r1(s)) == 0 && endsShortSyllable(s)
}

var step0Step = NewStep([]Suffix{
	{"'s' 's '", Delete, nil, nil, nil},
})

func step0(s []byte) []byte {
	return step0Step.Apply(s)
}

func step1aIed(s []byte, suffix []byte) []byte {
	s = s[:len(s)-len(suffix)]
	if len(s) > 1 {
		s = append(s, 'i')
	} else {
		s = append(s, 'i', 'e')
	}
	return s
}

func step1aNothing(s []byte, suffix []byte) []byte {
	return s
}

func step1aScb(s []byte, suffix []byte) []byte {
	if len(s) > 2 &&
		containsAny(s[:len(s)-2], vowels) {
		s = s[:len(s)-1]
	}
	return s
}

var step1aStep = NewStep([]Suffix{
	{"sses", Truncate("ss"), nil, nil, nil},
	{"ied ies", step1aIed, nil, nil, nil},
	{"s", step1aScb, nil, nil, nil},
	{"us ss", step1aNothing, nil, nil, nil},
})

func step1a(s []byte) []byte {
	return step1aStep.Apply(s)
}

func step1bReplace(ss []byte, suffix []byte) []byte {
	if len(ss) > len(suffix) &&
		containsAny(ss[:len(ss)-len(suffix)], vowels) {
		ss = ss[:len(ss)-len(suffix)]
		if anySuffix(ss, "at", "bl", "iz") {
			ss = append(ss, 'e')
		} else if hasDouble(ss) {
			ss = ss[:len(ss)-1]
		} else if isShort(ss) {
			ss = append(ss, 'e')
		}
	}
	return ss
}

var step1bStep = NewStep([]Suffix{
	{"eed eedly", Truncate("ee"), nil, r1, nil},
	{"ed edly ing ingly", step1bReplace, nil, nil, nil},
})

func step1b(s []byte) []byte {
	return step1bStep.Apply(s)
}

func step1c(s []byte) []byte {
	if len(s) >= 3 && (s[len(s)-1] == 'y' || s[len(s)-1] == 'Y') && !isVowel(s[len(s)-2]) {
		s[len(s)-1] = 'i'
	}
	return s
}

func step2Ogi(s []byte, suffix []byte) []byte {
	if len(s) > len(suffix) && s[len(s)-len(suffix)-1] == 'l' {
		return s[:len(s)-1]
	}
	return s
}

func step2Li(s []byte, suffix []byte) []byte {
	const validLiEnding = "cdeghkmnrt"
	if len(s) > len(suffix) && isAny(s[len(s)-len(suffix)-1], validLiEnding) {
		return s[:len(s)-len(suffix)]
	}
	return s
}

var step2Step = NewStep([]Suffix{
	{"tional", Truncate("tion"), nil, r1, nil},
	{"enci", Replace("ence"), nil, r1, nil},
	{"anci", Replace("ance"), nil, r1, nil},
	{"abli", Replace("able"), nil, r1, nil},
	{"entli", Truncate("ent"), nil, r1, nil},
	{"izer ization", Replace("ize"), nil, r1, nil},
	{"ational ation ator", Replace("ate"), nil, r1, nil},
	{"alism aliti alli", Truncate("al"), nil, r1, nil},
	{"fulness", Truncate("ful"), nil, r1, nil},
	{"ousli ousness", Truncate("ous"), nil, r1, nil},
	{"iveness iviti", Replace("ive"), nil, r1, nil},
	{"biliti bli", Replace("ble"), nil, r1, nil},
	{"ogi", step2Ogi, nil, r1, nil},
	{"fulli", Truncate("ful"), nil, r1, nil},
	{"lessli", Truncate("less"), nil, r1, nil},
	{"li", step2Li, nil, r1, nil},
})

func step2(s []byte) []byte {
	return step2Step.Apply(s)
}

var step3Step = NewStep([]Suffix{
	{"tional", Replace("tion"), nil, r1, nil},
	{"ational", Replace("ate"), nil, r1, nil},
	{"alize", Replace("al"), nil, r1, nil},
	{"icate iciti ical", Replace("ic"), nil, r1, nil},
	{"ful ness", Delete, nil, r1, nil},
	{"ative", Delete, nil, r2, nil},
})

func step3(s []byte) []byte {
	return step3Step.Apply(s)
}

func step4Ion(s []byte, suffix []byte) []byte {
	if len(s) > len(suffix) && isAny(s[len(s)-len(suffix)-1], "st") {
		s = s[:len(s)-len(suffix)]
	}
	return s
}

var step4Step = NewStep([]Suffix{
	{"al ance ence er ic able ible ant ement ment ent ism ate iti ous ive ize", Delete, nil, r2, nil},
	{"ion", step4Ion, nil, r2, nil},
})

func step4(s []byte) []byte {
	return step4Step.Apply(s)
}

func lastIs(s []byte, b byte) bool {
	return len(s) > 0 && s[len(s)-1] == b
}

func step5(s []byte) []byte {
	switch {
	case lastIs(s, 'e'):
		if lastIs(r2(s), 'e') {
			s = s[:len(s)-1]
		} else if lastIs(r1(s), 'e') {
			if !endsShortSyllable(s[:len(s)-1]) {
				s = s[:len(s)-1]
			}
		}
	case lastIs(s, 'l'):
		if lastIs(r2(s), 'l') && len(s) > 1 && s[len(s)-2] == 'l' {
			s = s[:len(s)-1]
		}
	}
	return s
}

var specialWords = map[string]string{
	"skis":   "ski",
	"skies":  "sky",
	"dying":  "die",
	"lying":  "lie",
	"tying":  "tie",
	"idly":   "idl",
	"gently": "gentl",
	"ugly":   "ugli",
	"early":  "earli",
	"only":   "onli",
	"singly": "singl",

	// invariant
	"sky":    "sky",
	"news":   "news",
	"howe":   "howe",
	"atlas":  "atlas",
	"cosmos": "cosmos",
	"bias":   "bias",
	"andes":  "andes",
}

var step1aInvariants = map[string]struct{}{
	"inning":  struct{}{},
	"outing":  struct{}{},
	"canning": struct{}{},
	"herring": struct{}{},
	"earring": struct{}{},
	"proceed": struct{}{},
	"exceed":  struct{}{},
	"succeed": struct{}{},
}

func StemBytes(s []byte) []byte {
	if len(s) <= 2 {
		return s
	}
	if s[0] == '\'' {
		s = s[1:]
	}
	s = step0(s)

	if v, ok := specialWords[string(s)]; ok {
		return []byte(v)
	}

	for i := range s {
		if s[i] == 'y' && (i == 0 || isVowel(s[i-1])) {
			s[i] = 'Y'
			break
		}
	}

	s = step1a(s)
	if _, ok := step1aInvariants[string(s)]; ok {
		return s
	}
	s = step1b(s)
	s = step1c(s)
	s = step2(s)
	s = step3(s)
	s = step4(s)
	s = step5(s)
	for i := range s {
		if s[i] == 'Y' {
			s[i] = 'y'
		}
	}
	return s
}

func StemString(s string) string {
	return string(StemBytes([]byte(s)))
}

func NormalizeBytes(b []byte) []byte {
	return bytes.ToLower(bytes.TrimSpace(b))
}

func NormalizeString(s string) string {
	return string(NormalizeBytes([]byte(s)))
}
