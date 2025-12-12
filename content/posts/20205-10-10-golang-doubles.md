---
layout: post
title: Test doubles in golang
date: 2025-10-10
tags: [tech]
---

[Test doubles in golang](https://en.wikipedia.org/wiki/Test_double)

```
type Doer  interface {
     Do() (int,error)
}

type Dummy  struct {}

func (d Dummy) (int,error) { panic("im a dummy") }


type Stub struct {}
func (s Stub) (int, error) { return 5, nil }

type Spy struct {
    count int
}
func (s Spy) Do (int, err)  { s.count++ return 5, nil }
func (s Spy) Count() int { return s.count}


type Fake struct {
      Value  int
      Errror error
}
func (f Fake) Do() (int, err) { return f.Value, f.Error }
```

Mocks are evil. Don't use mocks.
