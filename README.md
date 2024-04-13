# lxi

Go-based implementation of the LAN eXtensions for Instrumentation (LXI)
standard.

[![GoDoc][godoc badge]][godoc link]
[![Go Report Card][report badge]][report card]
[![License Badge][license badge]][LICENSE.txt]

## Overview

This packages enables controlling LXI compatible test equipment (e.g.,
oscilloscopes, function generators, multimeters, etc.) over Ethernet. While this
package can be used by itself to send Standard Commands for Programmable
Instruments ([SCPI][]) commands to a piece of test equipment, it also serves to
provide an Instrument interface for both the [ivi][] and [visa][] packages. The
[ivi][] package provides standardized APIs for programming test instruments
following the [Interchangeable Virtual Instrument (IVI) standard][ivi-specs].

## Installation

```bash
$ go get github.com/gotmc/lxi
```

## Documentation

Documentation can be found at either:

- <https://godoc.org/github.com/gotmc/lxi>
- <http://localhost:6060/pkg/github.com/gotmc/lxi/> after running `$
godoc -http=:6060`

## Contributing

Contributions are welcome! To contribute please:

1. Fork the repository
2. Create a feature branch
3. Code
4. Submit a [pull request][]

### Testing

Prior to submitting a [pull request][], please run the tests using either [GNU
Make][make]:

```bash
$ make check
$ make lint
```

or you can use [Just][]:

```bash
$ just check
$ just lint
```

To update and view the test coverage report using [Make][] run:

```bash
$ make cover
```

or you can use [Just][]:

```bash
$ just cover
```

## License

[lxi][] is released under the MIT license. Please see the
[LICENSE.txt][] file for more information.

[godoc badge]: https://godoc.org/github.com/gotmc/lxi?status.svg
[godoc link]: https://godoc.org/github.com/gotmc/lxi
[ivi]: https://github.com/gotmc/ivi
[ivi-foundation]: http://www.ivifoundation.org/
[ivi-specs]: http://www.ivifoundation.org/specifications/
[just]: https://just.systems/man/en/
[LICENSE.txt]: https://github.com/gotmc/lxi/blob/master/LICENSE.txt
[license badge]: https://img.shields.io/badge/license-MIT-blue.svg
[lxi]: https://github.com/gotmc/lxi
[make]: https://www.gnu.org/software/make/
[pull request]: https://help.github.com/articles/using-pull-requests
[report badge]: https://goreportcard.com/badge/github.com/gotmc/lxi
[report card]: https://goreportcard.com/report/github.com/gotmc/lxi
[scpi]: https://www.ivifoundation.org/About-IVI/scpi.html
[visa]: https://github.com/gotmc/visa
