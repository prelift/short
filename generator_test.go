// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at <http://mozilla.org/MPL/2.0/>.

package short_test

import (
	"bytes"
	"io"
	"math/rand"
	"strconv"
	"testing"

	"github.com/szabba/assert/v2"
	"github.com/szabba/assert/v2/assertions/theerr"
	"github.com/szabba/assert/v2/assertions/theval"

	"github.com/prelift/short"
)

func TestAlways(t *testing.T) {
	// given
	g := short.Always(123)

	// when
	v, err := g.Generate(nil)

	// then
	assert.Using(t.Errorf).
		That(theval.Equal(v, 123)).
		That(theerr.IsNil(err))
}

func TestBool(t *testing.T) {
	SucceedsGivenUnboundedInput(t, short.Bool)

	tests := map[string]struct {
		Src   []byte
		Value bool
		Err   error
	}{
		"EOF":      {Value: false, Err: io.EOF},
		"EvenByte": {Src: []byte{8}, Value: true},
		"OddByte":  {Src: []byte{7}, Value: false},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// given
			src := bytes.NewBuffer(tt.Src)

			// when
			v, err := short.Bool().Generate(src)

			// then
			assert.Using(t.Errorf).
				That(theval.Equal(v, tt.Value)).
				That(theerr.Is(err, tt.Err))
		})
	}
}

func TestInt(t *testing.T) {

	SucceedsGivenUnboundedInput(t, short.Int)

	t.Run("Unbiased", func(t *testing.T) {

		t.Run("Sign", func(t *testing.T) {
			// given
			const sampleSize = 1000
			gen := short.Int()
			rng := rand.New(rand.NewSource(rand.Int63()))

			neg, total := 0, 0

			// where
			for i := 0; i < 10_000; i++ {
				n, err := gen.Generate(rng)
				if err == nil {
					total++
					if n < 0 {
						neg++
					}
				}
			}

			// then
			negRatio := float64(neg) / float64(total)
			assert.
				Using(t.Errorf).
				That(total >= 100, "got low sample size (%d)", total).
				That(
					0.45 <= negRatio && negRatio <= 0.55,
					"negative ratio %f outside 0.5±0.05 range",
					negRatio)
		})

		t.Run("Oddity", func(t *testing.T) {
			// given
			const sampleSize = 1000
			gen := short.Int()
			rng := rand.New(rand.NewSource(rand.Int63()))

			odd, total := 0, 0

			// where
			for i := 0; i < 10_000; i++ {
				n, err := gen.Generate(rng)
				if err == nil {
					total++
					if n%2 == 1 {
						odd++
					}
				}
			}

			// then
			oddRatio := float64(odd) / float64(total)
			assert.
				Using(t.Errorf).
				That(total >= 100, "got low sample size (%d)", total).
				That(
					0.45 <= oddRatio && oddRatio <= 0.55,
					"odd number ratio %f outside the  0.5±0.05 range",
					oddRatio)
		})
	})
}

func SucceedsGivenUnboundedInput[Out any](t *testing.T, newGen func() short.Generator[Out]) {
	t.Run("SucceedsGivenUnboundedInput", func(t *testing.T) {
		// given
		gen := newGen()
		rng := rand.New(rand.NewSource(rand.Int63()))

		for i := 0; i < 1000; i++ {
			t.Run(strconv.Itoa(i), func(t *testing.T) {

				// when
				_, err := gen.Generate(rng)

				// then
				assert.Using(t.Errorf).That(theerr.IsNil(err))
			})
		}
	})
}
