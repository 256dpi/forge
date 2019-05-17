# forge

[![Build Status](https://travis-ci.org/256dpi/forge.svg?branch=master)](https://travis-ci.org/256dpi/forge)
[![Coverage Status](https://coveralls.io/repos/github/256dpi/forge/badge.svg?branch=master)](https://coveralls.io/github/256dpi/forge?branch=master)
[![GoDoc](https://godoc.org/github.com/256dpi/forge?status.svg)](http://godoc.org/github.com/256dpi/forge)
[![Release](https://img.shields.io/github/release/256dpi/forge.svg)](https://github.com/256dpi/forge/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/256dpi/forge)](https://goreportcard.com/report/github.com/256dpi/forge)

**A toolkit for building task pipelines in Go.**

## Example

```go
service := &Service{}

service.Report(func(err error) {
    panic(err.Error())
})

i := 0
service.Source(1, func(values chan<- Value) error {
    i++
    values <- i
    if i == 5 {
        return ErrDone
    }
    return nil
}, 1)

service.Batch(1, 3, time.Millisecond, 1)

service.FilterFunc(1, func(v Value, out chan<- Value) error {
    out <- v
    return nil
}, 1)

var out []Value
service.SinkFunc(1, func(v Value) error {
    out = append(out, v)
    return nil
})

<-service.Done()

fmt.Printf("%+v\n", out)

// Output:
// [[1 2 3] [4 5]]
```
