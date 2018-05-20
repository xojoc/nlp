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

// Package stem is a collection of stemming algorithms.
package stem

import (
	"xojoc.pw/nlp/stem/internal/porter2english"
	"xojoc.pw/nlp/stem/internal/porter2italian"
	"xojoc.pw/nlp/stem/internal/porter2spanish"
)

type Interface interface {
	// StemBytes returns the stem of b.
	// NOTE: StemBytes may modify b or return a sublice of it.
	StemBytes(b []byte) []byte
	// StemString returns the stem of s.
	StemString(s string) string
	// StemBytes assumes the input is already normalized.
	// You can use NormalizeBytes to normalize it.
	// NOTE: NormalizeBytes may modify b or return a sublice of it.
	NormalizeBytes(b []byte) []byte
	// StemString assumes the input is already normalized.
	// You can use NormalizeString to normalize it.
	NormalizeString(s string) string
}

type Porter2English struct{}

var _ Interface = Porter2English{}

func (Porter2English) StemBytes(b []byte) []byte {
	return porter2english.StemBytes(b)
}
func (Porter2English) StemString(s string) string {
	return porter2english.StemString(s)
}
func (Porter2English) NormalizeBytes(b []byte) []byte {
	return porter2english.NormalizeBytes(b)
}
func (Porter2English) NormalizeString(s string) string {
	return porter2english.NormalizeString(s)
}

type Porter2Italian struct{}

var _ Interface = Porter2Italian{}

func (Porter2Italian) StemBytes(b []byte) []byte {
	return porter2italian.StemBytes(b)
}
func (Porter2Italian) StemString(s string) string {
	return porter2italian.StemString(s)
}
func (Porter2Italian) NormalizeBytes(b []byte) []byte {
	return porter2italian.NormalizeBytes(b)
}
func (Porter2Italian) NormalizeString(s string) string {
	return porter2italian.NormalizeString(s)
}

type Porter2Spanish struct{}

var _ Interface = Porter2Spanish{}

func (Porter2Spanish) StemBytes(b []byte) []byte {
	return porter2spanish.StemBytes(b)
}
func (Porter2Spanish) StemString(s string) string {
	return porter2spanish.StemString(s)
}
func (Porter2Spanish) NormalizeBytes(b []byte) []byte {
	return porter2spanish.NormalizeBytes(b)
}
func (Porter2Spanish) NormalizeString(s string) string {
	return porter2spanish.NormalizeString(s)
}
