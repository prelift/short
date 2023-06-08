// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at <http://mozilla.org/MPL/2.0/>.

package short_test

import (
	"fmt"
	"testing"

	"github.com/prelift/short"
	"github.com/szabba/assert/v2"
	"github.com/szabba/assert/v2/assertions/theslice"
	"github.com/szabba/assert/v2/assertions/theval"
)

func TestCheckTrivialProperty(t *testing.T) {
	// given
	config := short.Config[int]{
		Generator: short.Int(),
		Property:  func(_ int) error { return nil },
	}

	// when
	res := config.Check()

	// then
	tried := len(res.Cases.Passed) + len(res.Cases.Failed)

	assert.Using(t.Errorf).
		That(res.Passed(), "check did not pass").
		That(!res.Failed(), "check failed").
		That(theval.Equal(tried, 10_000))
}

func TestCheckAllIntsAreEven(t *testing.T) {
	// given
	config := short.Config[int]{
		Generator: short.Int(),
		Property: func(n int) error {
			if n%2 != 0 {
				return fmt.Errorf("%d is odd", n)
			}
			return nil
		},
	}

	// when
	res := config.Check()

	// then
	assert.Using(t.Errorf).
		That(!res.Passed(), "check passed").
		That(res.Failed(), "check did not fail").
		That(theslice.NotEmpty(res.Cases.Failed))
}
