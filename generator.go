// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at <http://mozilla.org/MPL/2.0/>.

package short

import (
	"errors"
	"fmt"
	"io"
	"reflect"
)

func ErrFilteredOut() error { return errFilteredOut }

var errFilteredOut = errors.New("filtered out")

type Generator[Of any] interface {
	Generate(src io.Reader) (Of, error)
}

func Always[Of any](v Of) Generator[Of] {
	return constGenerator[Of]{v}
}

type constGenerator[Of any] struct {
	v Of
}

func (c constGenerator[Of]) Generate(_ io.Reader) (Of, error) {
	return c.v, nil
}

func Int() Generator[int] { return intGen{} }

type intGen struct{}

var intBytes = reflect.TypeOf(int(0)).Size()

func (intGen) Generate(src io.Reader) (int, error) {
	// FIXME: should generate negatives too
	p := make([]byte, intBytes)
	n, err := src.Read(p)
	if err != nil {
		return 0, err
	}
	if n != len(p) {
		err := fmt.Errorf(
			"too few bytes read (got %d, want %d): %w",
			n, len(p), io.ErrNoProgress)
		return 0, err
	}

	out := uint64(0)
	for _, b := range p {
		out = out*255 + uint64(b)
	}
	return int(out), nil
}

func Bool() Generator[bool] { return boolGen{} }

type boolGen struct{}

func (boolGen) Generate(src io.Reader) (bool, error) {
	var buf [1]byte
	_, err := src.Read(buf[:])
	if err != nil {
		return false, err
	}
	return buf[0]%2 == 0, nil
}

func Filter[Of any](
	gen Generator[Of],
	filter func(Of) (cause string, ok bool),
) Generator[Of] {

	return filterGenerator[Of]{gen, filter}
}

type filterGenerator[Of any] struct {
	gen    Generator[Of]
	filter func(Of) (cause string, ok bool)
}

func (fg filterGenerator[Of]) Generate(in io.Reader) (Of, error) {
	var zero Of

	kase, err := fg.gen.Generate(in)
	if err != nil {
		return zero, err
	}

	cause, ok := fg.filter(kase)
	if !ok {
		return zero, fmt.Errorf("%s: %w", cause, errFilteredOut)
	}

	return kase, nil
}
