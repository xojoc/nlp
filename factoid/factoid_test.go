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

package factoid_test

import (
	"testing"

	"gitlab.com/gonlp/factoid"
)

var questions = []string{
	"richard stallman",
	"when was richard stallman born?",
	"how old is richard stallman",
	"the color of the sun",
	"sun's color",
	"human eye's shape",
	"sound of christabel",
	"what continent is scotland in?",
	"what are the 7 wonders of the world",
}

var meanings = []factoid.Meaning{
	{Subject: "richard stallman"},
	{Subject: "richard stallman", Attribute: "birth date"},
	{Subject: "richard stallman", Attribute: "age"},
	{Subject: "sun", Attribute: "color"},
	{Subject: "sun", Attribute: "color"},
	{Subject: "human eye", Attribute: "shape"},
	{Subject: "christabel", Attribute: "sound"},
	{Subject: "scotland", Attribute: "continent"},
	{Subject: "7 wonders of the world", Action: "list"},
	{Subject: "", Attribute: ""},
}

func TestParse(t *testing.T) {
	for i, q := range questions {
		m := factoid.Parse(q)
		if m == nil {
			t.Errorf("no meaning for %q\n", q)
			continue
		}
		if *m != meanings[i] {
			t.Errorf("Parse(%q) = %+v; want %+v", q, m, meanings[i])
		}
	}
}
