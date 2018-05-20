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

package porter2

import (
	"strings"
	"unicode/utf8"
)

// http://snowballstem.org/texts/r1r2.html

func isVowel(r rune, vowels string) bool {
	return strings.ContainsRune(vowels, r)
}

func isConsonant(r rune, vowels string) bool {
	return !isVowel(r, vowels)
}

func R1(vowels string) func([]byte) []byte {
	return func(s []byte) []byte {
		for len(s) > 0 && isConsonant(rune(s[0]), vowels) {
			s = s[1:]
		}
		for {
			r, l := utf8.DecodeRune(s)
			if r == utf8.RuneError {
				return s
			}
			if isVowel(r, vowels) {
				s = s[l:]
			} else {
				break
			}
		}
		if len(s) > 0 {
			s = s[1:]
		}
		return s
	}
}

func R2(vowels string) func([]byte) []byte {
	r1 := R1(vowels)
	return func(s []byte) []byte {
		return r1(r1(s))
	}
}

// http://snowballstem.org/algorithms/spanish/stemmer.html

func SpanishRV(vowels string) func([]byte) []byte {
	return func(s []byte) []byte {
		r1, l := utf8.DecodeRune(s)
		s = s[l:]
		r2, l := utf8.DecodeRune(s)
		if r2 == utf8.RuneError {
			return s[len(s):]
		}
		s = s[l:]

		if isConsonant(r2, vowels) {
			for {
				r, l := utf8.DecodeRune(s)
				if l == 0 {
					break
				}
				s = s[l:]
				if isVowel(r, vowels) {
					break
				}
			}
			return s
		}

		if isVowel(r1, vowels) && isVowel(r2, vowels) {
			for {
				r, l := utf8.DecodeRune(s)
				if l == 0 {
					break
				}
				s = s[l:]
				if isConsonant(r, vowels) {
					break
				}
			}
			return s
		}

		_, l = utf8.DecodeRune(s)
		return s[l:]
	}
}
