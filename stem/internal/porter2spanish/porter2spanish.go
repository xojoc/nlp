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

import (
	"bytes"

	. "xojoc.pw/nlp/stem/internal/porter2"
)

// http://snowballstem.org/algorithms/spanish/stemmer.html

const vowels = "aeiouáéíóúü"

var (
	r1 = R1(vowels)
	r2 = R2(vowels)
	rv = SpanishRV(vowels)
)

func stripAfter(after string) func(s []byte, suffix []byte) []byte {
	return func(s []byte, suffix []byte) []byte {
		if bytes.HasSuffix(s, []byte(after+string(suffix))) {
			return s[:len(s)-len(suffix)+len(after)]
		}
		return s
	}
}

func step0(s []byte) []byte {
	return s
}

var step1Step = NewStep([]Suffix{
	{"anza anzas ico ica icos icas ismo ismos able ables ible ibles ista istas oso osa osos osas amiento amientos imiento imientos", Delete, nil, r2, nil},
	{"adora ador ación adoras adores aciones ante antes ancia ancias", Delete, nil, r2, []Suffix{{"ic", Delete, nil, r2, nil}}},
	{"logía logías", Truncate("log"), nil, r2, nil},
	{"ución uciones", Truncate("u"), nil, r2, nil},
	{"encia encias", Replace("ente"), nil, r2, nil},
	{"amente", Delete, nil, r1,
		[]Suffix{
			{"iv", Delete, nil, r2, []Suffix{{"at", Delete, nil, r2, nil}}},
			{"os ic ad", Delete, nil, r2, nil},
		},
	},
	{"mente", Delete, nil, r2, []Suffix{{"ante able ible", Delete, nil, r2, nil}}},
	{"idad idades", Delete, nil, r2, []Suffix{{"abil ic iv", Delete, nil, r2, nil}}},
	{"iva ivo ivas ivos", Delete, nil, r2, []Suffix{{"at", Delete, nil, r2, nil}}},
})

func step1(s []byte) []byte {
	return step1Step.Apply(s)
}

var step2aStep = NewStep([]Suffix{
	{"ya ye yan yen yeron yendo yo yó yas yes yais yamos", stripAfter("u"), rv, nil, nil},
})

func step2a(s []byte) []byte {
	return step2aStep.Apply(s)
}

func step2bGu(s []byte, suffix []byte) []byte {
	s = s[:len(s)-len(suffix)]
	if bytes.HasSuffix(s, []byte("gu")) {
		s = s[:len(s)-1]
	}
	return s
}

var step2bStep = NewStep([]Suffix{
	{"en es éis emos", step2bGu, rv, nil, nil},
	{"arían arías arán arás aríais aría aréis aríamos aremos ará aré erían erías erán erás eríais ería eréis eríamos eremos erá eré irían irías irán irás iríais iría iréis iríamos iremos irá iré aba ada ida ía ara iera ad ed id ase iese aste iste an aban ían aran ieran asen iesen aron ieron ado ido ando iendo ió ar er ir as abas adas idas ías aras ieras ases ieses ís áis abais íais arais ierais aseis ieseis asteis isteis ados idos amos ábamos íamos imos áramos iéramos iésemos ásemos", Delete, rv, nil, nil},
})

func step2b(s []byte) []byte {
	return step2bStep.Apply(s)
}

var step3Step = NewStep([]Suffix{
	{"os a o á í ó", Delete, rv, nil, nil},
	{"e é", Delete, rv, nil, nil},
})

func step3(s []byte) []byte {
	return step3Step.Apply(s)
}

func stripAccents(s []byte) []byte {
	return bytes.Map(func(r rune) rune {
		switch r {
		case 'á':
			return 'a'
		case 'é':
			return 'e'
		case 'í':
			return 'i'
		case 'ó':
			return 'o'
		case 'ú':
			return 'u'
		case 'ü':
			return 'u'
		default:
			return r
		}
	}, s)
}

func StemBytes(s []byte) []byte {
	s = step0(s)
	s1 := step1(s)
	if bytes.Equal(s1, s) {
		sa := step2a(s)
		if bytes.Equal(sa, s) {
			s = step2b(s)
		} else {
			s = sa
		}
	} else {
		s = s1
	}
	s = step3(s)
	for i, b := range s {
		if b == 'I' {
			s[i] = 'i'
		} else if b == 'U' {
			s[i] = 'u'
		}
	}
	s = stripAccents(s)
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
