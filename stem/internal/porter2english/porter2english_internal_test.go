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

func TestStep1a(t *testing.T) {
	var pairs = []string{
		"ties", "tie", "cries", "cri", "gas", "gas", "this", "this", "gaps", "gap", "kiwis", "kiwi"}
	test(t, step1a, pairs)
}

func TestStep1b(t *testing.T) {
	var pairs = []string{"luxuriatingly", "luxuriate", "hopped", "hop", "hoping", "hope", "writing", "write"}
	test(t, step1b, pairs)
}
func TestStep1c(t *testing.T) {
	var pairs = []string{"cry", "cri", "by", "by", "say", "say"}
	test(t, step1c, pairs)
}

func TestStep2(t *testing.T) {
	var pairs = []string{"fepotional", "fepotion", "belogi", "belog"}
	test(t, step2, pairs)
}
