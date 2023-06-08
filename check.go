// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at <http://mozilla.org/MPL/2.0/>.

package short

import (
	"bytes"
	cryptorand "crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"math/big"
	"math/rand"
)

type Config[Input any] struct {
	Generator Generator[Input]
	Property  func(Input) error

	Seed   *int64
	Source rand.Source

	Limit struct {
		InitialSamples int
	}
}

func (c Config[Input]) Check() Result[Input] {
	r := Result[Input]{Config: c}
	src := r.source()
	r.sampleUntilFailure(src)
	r.seekSimpler(src)
	r.reverseFailures()
	return r
}

type Result[Input any] struct {
	Config Config[Input]

	GeneratorErrors []error
	Cases           struct {
		Passed []Input
		Failed []Failure[Input]
	}
}

type Failure[Input any] struct {
	Case      Input
	Err       error
	BytesRead []byte
}

func (res Result[Input]) Passed() bool { return !res.Failed() }

func (res Result[Input]) Failed() bool { return len(res.Cases.Failed) > 0 }

func (res Result[Input]) source() rand.Source {

	seed := res.seed()

	if res.Config.Source != nil {
		res.Config.Source.Seed(seed)
		return res.Config.Source
	}

	return rand.NewSource(seed)
}

func (res Result[Input]) seed() int64 {
	if res.Config.Seed != nil {
		return *res.Config.Seed
	}

	var buf [8]byte
	_, err := cryptorand.Read(buf[:])
	if err != nil {
		panic("cannot seed random source")
	}

	return int64(binary.LittleEndian.Uint64(buf[:]))
}

func (res *Result[Input]) sampleUntilFailure(src rand.Source) {
	const maxTries = 10_000
	for i := 0; i < maxTries; i++ {

		kase, bytesRead, err := res.generate(src)
		if err != nil {
			continue
		}

		err = res.check(kase)
		if err != nil {
			res.recordFailure(kase, err, bytesRead)
			return
		}

		res.recordSuccess(kase)
	}
}

func (res *Result[Input]) seekSimpler(src rand.Source) {
	const maxTries = 10_000
	if len(res.Cases.Failed) == 0 {
		return
	}

	for i := 0; i < maxTries; i++ {
		res.seekOneSimpler(src)
	}
}

func (res *Result[Input]) reverseFailures() {
	fails := make([]Failure[Input], 0, len(res.Cases.Failed))

	for _, f := range res.Cases.Failed {
		fails = append(fails, f)
	}
	res.Cases.Failed = fails
}

func (res *Result[Input]) seekOneSimpler(src rand.Source) {
	minFail := res.Cases.Failed[len(res.Cases.Failed)-1]
	bs := minFail.BytesRead
	var err error

	bs, err = res.shorter(minFail.BytesRead, src)
	if err != nil {
		err = fmt.Errorf("failed to sample input smaller than %x: %w", bs, err)
		res.GeneratorErrors = append(res.GeneratorErrors, err)
		return
	}

	buf := bytes.NewBuffer(bs)
	in, err := res.generateFromReader(buf)
	if err != nil {
		return
	}

	err = res.check(in)
	if err == nil {
		res.recordSuccess(in)
		return
	}

	res.recordFailure(in, err, bs)
}

func (res *Result[Input]) shorter(bs []byte, src rand.Source) ([]byte, error) {
	max := new(big.Int).SetBytes(bs)
	rng := rand.New(src)
	s, err := cryptorand.Int(rng, max)
	if err != nil {
		return nil, err
	}
	return s.Bytes(), nil
}

func (res *Result[Input]) generate(src rand.Source) (Input, []byte, error) {
	rng := rand.New(src)
	buf := new(bytes.Buffer)
	r := io.TeeReader(rng, buf)
	kase, err := res.generateFromReader(r)
	return kase, buf.Bytes(), err
}

func (res *Result[Input]) generateFromReader(r io.Reader) (Input, error) {
	kase, err := res.Config.Generator.Generate(r)
	if err != nil {
		res.GeneratorErrors = append(res.GeneratorErrors, err)
	}
	return kase, err
}

func (res Result[Input]) check(kase Input) error {
	return res.Config.Property(kase)
}

func (res *Result[Input]) recordFailure(kase Input, err error, bytesRead []byte) {
	fail := Failure[Input]{kase, err, bytesRead}
	res.Cases.Failed = append(res.Cases.Failed, fail)
}

func (res *Result[Input]) recordSuccess(kase Input) {
	res.Cases.Passed = append(res.Cases.Passed, kase)
}
