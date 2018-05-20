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

package stem_test

import (
	"fmt"

	"gitlab.com/gonlp/stem"
)

func Example() {
	st := stem.Porter2English{}
	s := "  NoRmAlIzEd  "
	s = st.NormalizeString(s)
	fmt.Printf("%s\n", s)
	s = st.StemString(s)
	fmt.Printf("%s\n", s)
	//Output:
	// normalized
	// normal
}

func ExamplePorter2English_StemBytes() {
	st := stem.Porter2English{}
	b := []byte("running")
	// StemBytes may modify b or return a subslice of it.
	b = st.StemBytes(b)
	fmt.Printf("%s", b)
	//Output: run
}
func ExamplePorter2English_StemString() {
	st := stem.Porter2English{}
	fmt.Println(st.StemString("enjoying"))
	//Output: enjoy
}
func ExamplePorter2English_NormalizeBytes() {
	st := stem.Porter2English{}
	b := []byte("  NoRmAlIzEd  ")
	// NormalizeBytes may modify b or return a subslice of it.
	b = st.NormalizeBytes(b)
	fmt.Printf("%s\n", b)
	// StemBytes may modify b or return a subslice of it.
	b = st.StemBytes(b)
	fmt.Printf("%s", b)
	//Output: normalized
	// normal
}
