# forge

[![Build Status](https://travis-ci.org/256dpi/forge.svg?branch=master)](https://travis-ci.org/256dpi/forge)
[![Coverage Status](https://coveralls.io/repos/github/256dpi/forge/badge.svg?branch=master)](https://coveralls.io/github/256dpi/forge?branch=master)
[![GoDoc](https://godoc.org/github.com/256dpi/forge?status.svg)](http://godoc.org/github.com/256dpi/forge)
[![Release](https://img.shields.io/github/release/256dpi/forge.svg)](https://github.com/256dpi/forge/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/256dpi/forge)](https://goreportcard.com/report/github.com/256dpi/forge)

**A toolkit for managing long-running tasks in Go.**

## Example

```go
// prepare service
service := &forge.Service{}

// set reporter
var errs []error
service.Report(func(err error) {
    errs = append(errs, err)
})

// run task
i := 0
service.Run(1, func() error {
    i++

    if i == 5 {
        return forge.ErrDone
    }
    if i%2 == 0 {
        return errors.New("foo")
    }

    return nil
}, func() {
    fmt.Println("finalize")
})

// wait for exit
<-service.Done()

// print output
fmt.Println(i)
fmt.Println(errs)

// Output:
// finalize
// 5
// [foo foo]
```
