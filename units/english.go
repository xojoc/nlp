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

package units

import (
	"fmt"
	"strconv"
	"strings"
)

const englishAliases = `
kilogram	kg

`

// FIXME: aliases:
//  centimeters
// TODO: handle "how many ounces in a pound"

// TODO: stemming: meters -> meter, inches -> inch
var unitsAliases = map[string]string{
	"meter": "m",
}

// English parses an english phrase and prepares the arguments for Convert. TODO
func English(s string) (fnum float64, funit string, tunit string, err error) {
	words := strings.Fields(s)
	if len(words) == 4 {
		if words[2] == "to" || words[2] == "in" {
			words = append(words[:2], words[3:]...)
		} else {
			return fnum, funit, tunit, fmt.Errorf("cannot parse phrase: %q", s)
		}
	}
	if len(words) != 3 {
		return fnum, funit, tunit, fmt.Errorf("cannot parse phrase: %q", s)
	}
	fnum, err = strconv.ParseFloat(words[0], 64)
	if err != nil {
		return fnum, funit, tunit, err
	}

	funit = words[1]
	tunit = words[2]
	return
}
