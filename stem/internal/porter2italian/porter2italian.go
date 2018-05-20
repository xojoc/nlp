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

package porter2italian

import (
	"bytes"
	"strings"
	"unicode/utf8"

	. "xojoc.pw/nlp/stem/internal/porter2"
)

// http://snowballstem.org/algorithms/italian/stemmer.html

var vowels = "aeiouàèìòù"

func isVowel(r rune) bool {
	return strings.ContainsRune(vowels, r)
}

var (
	r1 = R1(vowels)
	r2 = R2(vowels)
	rv = SpanishRV(vowels)
)

func normalize(s []byte) []byte {
	p := s
	pr := utf8.RuneError
	for {
		r, l := utf8.DecodeRune(p)
		if l == 0 {
			break
		}
		switch r {
		case 'á':
			utf8.EncodeRune(p, 'à')
		case 'é':
			utf8.EncodeRune(p, 'è')
		case 'í':
			utf8.EncodeRune(p, 'ì')
		case 'ó':
			utf8.EncodeRune(p, 'ò')
		case 'ú':
			utf8.EncodeRune(p, 'ù')
		case 'i', 'u':
			if r == 'u' && pr == 'q' {
				utf8.EncodeRune(p, 'U')
			} else {
				nr, _ := utf8.DecodeRune(p[l:])
				if isVowel(pr) && isVowel(nr) {
					if r == 'i' {
						utf8.EncodeRune(p, 'I')
					} else if r == 'u' {
						utf8.EncodeRune(p, 'U')
					} else {
						panic("unreachable")
					}
				}
			}

		}
		pr, _ = utf8.DecodeRune(p)
		//	pr = r
		p = p[l:]
	}
	return s
}

func step0CB(s []byte, suffix []byte) []byte {
	m := s[:len(s)-len(suffix)]
	v := rv(m)
	if bytes.HasSuffix(v, []byte("ando")) || bytes.HasSuffix(v, []byte("endo")) {
		return m
	} else if bytes.HasSuffix(v, []byte("ar")) || bytes.HasSuffix(v, []byte("er")) || bytes.HasSuffix(v, []byte("ir")) {
		return append(m, 'e')
	}
	return s
}

var step0Step = NewStep([]Suffix{
	{"ci gli la le li lo mi ne si ti vi sene gliela gliele glieli glielo gliene mela mele meli melo mene tela tele teli telo tene cela cele celi celo cene vela vele veli velo vene", step0CB, rv, nil, nil},
})

func step0(s []byte) []byte {
	return step0Step.Apply(s)
}

var step1Step = NewStep([]Suffix{
	{"anza anze ico ici ica ice iche ichi ismo ismi abile abili ibile ibili ista iste isti istà istè istì oso osi osa ose mente atrice atrici ante anti", Delete, nil, r2, nil},
	{"azione azioni atore atori", Delete, nil, r2, []Suffix{{"ic", Delete, nil, r2, nil}}},
	{"logia logie", Truncate("log"), nil, r2, nil},
	{"uzione uzioni usione usioni", Truncate("u"), nil, r2, nil},
	{"enza enze", Replace("ente"), nil, r2, nil},
	{"amento amenti imento imenti", Delete, nil, rv, nil},
	{"amente", Delete, nil, r1,
		[]Suffix{
			{"iv", Delete, nil, r2, []Suffix{{"at", Delete, nil, r2, nil}}},
			{"os ic abil", Delete, nil, r2, nil},
		},
	},
	{"ità", Delete, nil, r2, []Suffix{{"abil ic iv", Delete, nil, r2, nil}}},
	{"ivo ivi iva ive", Delete, nil, r2, []Suffix{{"at", Delete, nil, r2, []Suffix{{"ic", Delete, nil, r2, nil}}}}},
})

func step1(s []byte) []byte {
	return step1Step.Apply(s)

}

var step2Step = NewStep([]Suffix{
	{"ammo ando ano are arono asse assero assi assimo ata ate ati ato ava avamo avano avate avi avo emmo enda ende endi endo erà erai eranno ere erebbe erebbero erei eremmo eremo ereste eresti erete erò erono essero ete eva evamo evano evate evi evo Yamo iamo immo irà irai iranno ire irebbe irebbero irei iremmo iremo ireste iresti irete irò irono isca iscano isce isci isco iscono issero ita ite iti ito iva ivamo ivano ivate ivi ivo ono uta ute uti uto ar ir", Delete, rv, nil, nil},
})

func step2(s []byte) []byte {
	return step2Step.Apply(s)
}

var step3aStep = NewStep([]Suffix{
	{"a e i o à è ì ò", Delete, nil, rv, []Suffix{{"i", Delete, rv, nil, nil}}},
})

func step3a(s []byte) []byte {
	return step3aStep.Apply(s)
}

var step3bStep = NewStep([]Suffix{
	{"ch", Truncate("c"), nil, rv, nil},
	{"gh", Truncate("g"), nil, rv, nil},
})

func step3b(s []byte) []byte {
	return step3bStep.Apply(s)
}

func StemBytes(s []byte) []byte {
	s = normalize(s)
	s = step0(s)
	s1 := step1(s)
	if bytes.Equal(s1, s) {
		s = step2(s)
	} else {
		s = s1
	}
	s = step3a(s)
	s = step3b(s)
	for i, b := range s {
		if b == 'I' {
			s[i] = 'i'
		} else if b == 'U' {
			s[i] = 'u'
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
