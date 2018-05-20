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

// Parse simple factoid questions.
package factoid

import "strings"

// Examples of factoid questions:
//   https://raw.githubusercontent.com/brmson/dataset-factoid-curated/master/curated-full.tsv

// strip quotes
// remove stop words? (of the a an)

/*
 Type of factoid questions:
   - Generic information about an object: "Richard Stallman", "Gandhi", "Free Software Foundation"
   - Attribute of an object: "When did _ die", "When was _ born", "How hot is the Sun", "What continent is Scotland in?",
   - "Creator": "Who painted "Sunflowers"?", "Who makes viagra?"

 When is Fashion week in NYC?
 What site did Lindbergh begin his flight from in 1927?
 What country is the holy city of Mecca located in?
 What war is connected with the book "Charge of the Light Brigade"?
 Where was the first atomic bomb detonated?
 What continent is Scotland in?
 How deep is Crater Lake?

*/

// who won the nobel peace prize? (nobel peace prize winners)
// attribute of attribute
// time range
// capital of germany

/*
1. Who is the author of the book, "The Iron Lady:  A Biography of Margaret Thatcher"?
2. What was the monetary value of the Nobel Peace Prize in 1989?
3. What does the Peugeot company manufacture?
4. How much did Mercury spend on advertising in 1993?
5. What is the name of the managing director of Apricot Computer?
7. What debts did Qintex group leave?
8. What is the name of the rare neurological disease with symptoms such as:
  involuntary movements (tics), swearing,
    and incoherent vocalizations (grunts, shouts, etc.)
*/

/*
who invented surf music?
which english translation of the bible is used in official catholic liturgies?
how tall is the sears tower?
*/

//  Name the first private citizen to fly in space
//  When was the internal combustion engine invented ?
//  How hot does the inside of an active volcano get?

/*
	{"how old is X|age of X|X's age|when was X born", birthDeathOutput},
	// video
	{"X's pronunciation|how to pronunce X", pronunciationOutput},
	// gender
*/

/*
* who made this
* Who discovered the first star?
* who invented the television
* all X films / lists

* who started the world war 2
* why is the sky blue

## first order
* tell me about (get description)

## second order

* who is /  (first get name, then make sure it's human)
* what is (first get name, then make sure it's *not* human)
* is __ dead?
* is __ alive?
* how old is
* where is (openstreetmap)

## third order
* [10,5] largest cities in Europe (get cities, in europe, order by larginess)

* who killed
* who gave birth (mother father)
* who are the parents of
* who are the children of
* who is the oldest [known] human being
* why the sky is blue or how Napoleon died?
*/

type Type int

const (
	Date Type = iota
	Location
	Measure
	Property // color
)

type Meaning struct {
	Subject    string
	Attribute  string
	Location   string // Answer valid in this location: When is fashion week in NYC? // TODO: normalize?
	Period     string // Answer valid in this period (August, 2016, II century) // TODO: normalize?
	AnswerType Type
}

// attributes: age, birth date, death date, born, color, shape, sound, audio, video, pronunciation

var attributeSynonyms = map[string]string{
	"born":  "birth date",
	"old":   "age",
	"audio": "sound",
	"hot":   "temperature",
}

func normalize(m *Meaning) {
	a := attributeSynonyms[m.Attribute]
	if a != "" {
		m.Attribute = a
	}
	s := m.Subject
	if len(s) > 3 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}
	m.Subject = s
}

type factoid func(string) *Meaning

var factoids = []factoid{
	simplePattern,
}

func Parse(s string) *Meaning {
	s = strings.TrimRight(s, "?")
	for _, f := range factoids {
		m := f(s)
		if m != nil {
			normalize(m)
			return m
		}
	}
	return nil
}
