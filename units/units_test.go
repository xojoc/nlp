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

package units_test

import (
	"fmt"
	"testing"

	"xojoc.pw/must"
	"xojoc.pw/nlp/units"
)

func ExampleConvert() {
	tnum, err := units.Convert(10, "m", "in")
	must.OK(err)
	fmt.Printf("%.3f in", tnum)
	//Output: 393.701 in
}

func ExampleEnglish() {
	fnum, fu, tu, err := units.English("10 cm to km")
	must.OK(err)
	tnum, err := units.Convert(fnum, fu, tu)
	must.OK(err)
	fmt.Println(tnum, tu)
	// Output: 0.0001 km
}

type entry struct {
	fnum  float64
	funit string
	tunit string
}

var english map[string]*entry = map[string]*entry{
	"10 cm to km":                     {10, "cm", "km"},
	"how many centimeters in a meter": {1, "m", "cm"},
	"£1 to euro":                      {1, "£", "€"},
	"10 kilograms to grams":           {10, "kg", "g"},
	"10 nonsense to nope":             {0, "", ""},
}

func TestEnglish(t *testing.T) {
	for i, e := range english {
		fnum, funit, tunit, _ := units.English(i)
		if fnum != e.fnum || funit != e.funit || tunit != e.tunit {
			t.Logf("%v: got: %v %q %q -- want: %v %q %q\n", i, fnum, funit, tunit, e.fnum, e.funit, e.tunit)
			t.Fail()
		}
	}
}

type conversion struct {
	fnum  float64
	funit string
	tunit string
	tnum  float64
}

var conversions []conversion = []conversion{
	{10, "m", "cm", 1000},
	{10, "cm", "in", 3.9370078740157477},
}

func TestConvert(t *testing.T) {
	for _, v := range conversions {
		tnum, _ := units.Convert(v.fnum, v.funit, v.tunit)
		if tnum != v.tnum {
			t.Fatalf("%v %v %v: got: %v -- want: %v\n", v.fnum, v.funit, v.tunit, tnum, v.tnum)
		}
	}
}
