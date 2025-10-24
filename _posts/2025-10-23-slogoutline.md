---
layout: post
title: slog is neat
tags: [ tech ]
---



Slog talk out line.

## existing log package.
its like printf but theres a timestamp!
Well there is [Setoutput](https://cs.opensource.google/go/go/+/refs/tags/go1.25.3:src/log/log.go;l=364)

## what should libraries do
* Well I can pass up errors
* What about that thing that seems wrong but I'm suppressing I don't won't to break compatbilities but pretty sure this is wront
* https://github.com/Azure/eno/pull/532
* How do I associate these things with the current request.
* How do log respectfully with out overlogging

## What about X
zap and logrus are pretty neat

But something nice about just being in

## Levels
* Throw all your
* but don't forget that blog post

## Structure
* Request ids! Traces

### Putting things on context?
Slog context

https://betterstack.com/community/guides/logging/golang-contextual-logging/

Not everything on context is logge handler has to pull it off. So libraries don't have a great way to add things?

## Handlers!


## Sinks
Sometimes you do just need a writer.


## awesome slog.
https://github.com/go-slog/awesome-slog