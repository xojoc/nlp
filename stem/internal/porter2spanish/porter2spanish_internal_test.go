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

package porter2spanish

import "testing"

func test(t *testing.T, fn func([]byte) []byte, pairs []string) {
	for i := 0; i < len(pairs); i += 2 {
		actual := string(fn([]byte(pairs[i])))
		if actual != pairs[i+1] {
			t.Errorf("fn(%q) = %v; want %v", pairs[i], actual, pairs[i+1])
		}
	}

}

func TestR1(t *testing.T) {
	pairs := []string{"beautiful", "iful", "beauty", "y", "beau", "", "animadversion", "imadversion", "sprinkled", "kled", "eucharist", "harist"}
	test(t, r1, pairs)
}

func TestR2(t *testing.T) {
	pairs := []string{"beautiful", "ul", "beauty", "", "beau", "", "animadversion", "adversion", "sprinkled", "", "eucharist", "ist"}
	test(t, r2, pairs)
}

func TestRv(t *testing.T) {
	pairs := []string{"macho", "ho", "oliva", "va", "trabajo", "bajo", "àureo", "eo"}
	test(t, rv, pairs)
}

func TestStep0(t *testing.T) {
	pairs := []string{"haciéndola", "haciendo"}
	test(t, step0, pairs)
}
