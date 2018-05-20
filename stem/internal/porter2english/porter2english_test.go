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

package porter2english_test

import (
	"bufio"
	"bytes"
	"strings"
	"testing"

	"gitlab.com/gonlp/stem/internal/porter2english"
	"gitlab.com/xojoc/util"
)

func TestStemBytes(t *testing.T) {
	f := util.MustOpen("testfiles/vocabulary.txt")
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		words := strings.Fields(scanner.Text())
		s := porter2english.StemString(words[0])
		if s != words[1] {
			t.Errorf("StemString(%q): expected %q got %q\n", words[0], words[1], s)
		}
	}

	util.Fatal(scanner.Err())
}

var words [][]byte

func loadWords() {
	if words != nil {
		return
	}
	f := util.MustOpen("testfiles/vocabulary.txt")
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		ws := bytes.Fields(scanner.Bytes())
		words = append(words, ws[0])
	}
	util.Fatal(scanner.Err())

}

func BenchmarkStemBytes(b *testing.B) {
	loadWords()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, w := range words {
			_ = porter2english.StemBytes(w)
		}
	}
}
