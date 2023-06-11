# `short.Check` - a

> :warning: Experimental software.

`short.Check` is a property testing library for Go.

Essentially:

* You write a property - a function that should never fail for _any_ input it gets.
* You say how to generate random values.
* It generates many random inputs looking for an example where the property always holds.
* It attempts to find _smaller_ inputs.
    *
* It reports all the checks it's run.

## Examples

TBW: write some.

## What do you mean "smaller" inputs?

TBW:

## vs `"testing/quick"`

TBW: why `"testing/quick"` is subpar.

## vs fuzzing

TBW: when to reach for property testing vs fuzzing.

## Licensing

Copyright Karol Marcjan 2023-

`short.Check` is released under the Mozilla Public License 2.0.
For more details see the [`LICENSE`] file.

[`LICENSE`]: ./LICENSE
