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

// Code common to all porter2 algorithms.
package porter2

import (
	"bytes"
	"strings"
)

type action struct {
	Suffix       []byte
	Callback     func([]byte, []byte) []byte
	MatchRegion  func([]byte) []byte
	ActionRegion func([]byte) []byte
	Step         *Step
}

type node struct {
	leafs          [255]*node
	suffixCallback *action
}

func insertSuffix(n *node, x []byte, cb *action) *node {
	if n == nil {
		n = &node{}
	}
	if len(x) == 0 {
		n.suffixCallback = cb
	} else {
		n.leafs[x[0]] = insertSuffix(n.leafs[x[0]], x[1:], cb)
	}
	return n
}

func reverse(a []byte) []byte {
	b := make([]byte, len(a))
	for i, x := range a {
		b[len(b)-i-1] = x
	}
	return b
}

func newTrie(suffixes []*action) *node {
	var root *node
	for _, cb := range suffixes {
		root = insertSuffix(root, reverse(cb.Suffix), cb)
	}
	return root
}

type Suffix struct {
	Suffixes        string
	Callback        func([]byte, []byte) []byte
	MatchRegion     func([]byte) []byte
	ActionRegion    func([]byte) []byte
	SuffixCallbacks []Suffix
}

type Step struct {
	suffixes *node
}

func NewStep(suffixCallbacks []Suffix) *Step {
	if len(suffixCallbacks) == 0 {
		return nil
	}
	s := &Step{}
	var callbacks []*action
	for _, x := range suffixCallbacks {
		suffixes := strings.Fields(x.Suffixes)
		for _, f := range suffixes {
			c := action{[]byte(f), x.Callback, x.MatchRegion, x.ActionRegion, NewStep(x.SuffixCallbacks)}
			callbacks = append(callbacks, &c)

		}
	}
	s.suffixes = newTrie(callbacks)
	return s
}

func (s *Step) Apply(str []byte) []byte { // return true if changed?
	if s == nil {
		return str
	}
	var suffix *action
	p := s.suffixes
	for i := len(str) - 1; i >= 0; i-- {
		p = p.leafs[str[i]]
		if p == nil {
			break
		}
		sx := p.suffixCallback
		if sx != nil {
			if sx.MatchRegion == nil || bytes.HasSuffix(sx.MatchRegion(str), sx.Suffix) {
				suffix = sx
			}
		}
	}

	if suffix == nil {
		return str
	}
	if suffix.ActionRegion == nil {
		str = suffix.Callback(str, suffix.Suffix)
	} else if bytes.HasSuffix(suffix.ActionRegion(str), suffix.Suffix) {
		str = suffix.Callback(str, suffix.Suffix)
	}
	return suffix.Step.Apply(str)
}

// Common callbacks.

func Delete(s []byte, suffix []byte) []byte {
	return s[:len(s)-len(suffix)]
}

func Truncate(to string) func([]byte, []byte) []byte {
	return func(s []byte, suffix []byte) []byte {
		return s[:len(s)-len(suffix)+len(to)]
	}
}

func Replace(by string) func(s []byte, suffix []byte) []byte {
	return func(s []byte, suffix []byte) []byte {
		s = s[:len(s)-len(suffix)]
		return append(s, by...)
	}
}
