---
layout: post
title: slog is neat
tags: [ tech ]
---

Why did golang get a new logging framework after ~10 years? 

## existing ["log"]([url](https://pkg.go.dev/log)) package. 
* Its like printf but theres a timestamp!
* Well I can Fatal and Panic I guess thats nice in main.
```
log.Println("Hello, 世界")
log.Printf("I will se you in %s\n", time.Hour)
log.Fatal("Oh no its all gone wrong")

2009/11/10 23:00:00 Hello, 世界
2009/11/10 23:00:00 I will se you in 1h0m0s
2009/11/10 23:00:00 Oh now its all gone wrong
```
  
* There  is [Setoutput](https://cs.opensource.google/go/go/+/refs/tags/go1.25.3:src/log/log.go;l=364) Lets come back to this later
* And I can get a Default logger but I can't SET a default logger. 

## What should libraries do
Above is pretty great if your in main.go but what if you're writing a package for reuse?

* Should libraries log at all?
  * Well I can pass up errors. (wrap them if you're extra cool!)
  * What about that thing that seems wrong but I'm suppressing I don't won't to break compatbilities but pretty sure this is wront
  * https://github.com/Azure/eno/pull/532

* How do get the logger?
  *   Package level variable/Function?
  *   Pass it everywhere?
  *   Shove it on context
  *   Have seen all of these 

* How do log respectfully with out overlogging?
    
* What if I want to customize my log format?
    * Json / Binary
    * Coloring / Emoji✅ 

## What about X
* zap and logrus are pretty neat. They introduced level, structure, sinks and alot of the stuff I am going to talk about below. 
* But if I'm writing a library how do I know what interface to get/log to?
* have seen several instances where logging adapters are regularly rewritten. Usually losing somethign each time. 
* Sometimes its noice to just have the standard library 


## Enter Slog
[Gopkg](https://pkg.go.dev/log/slog@go1.25.4) and [Blog](https://go.dev/blog/slog)

* A high levl interface you can impment behnd
* A text and json implementation
* Set Default

## Levels
Debug/Info/Warn/Error
* Throw all your junk into debug!
* indecisive? Warn
* make up your own log levels and go crazy. 
* but don't forget that blog post

## Structure
* String parsing stinks! We have kusto and lots of other fancy databases!
* Request ids and other tracing items.
* Latencies as durations! 
* Customizing and 

### Putting things on context?
Slog context

https://betterstack.com/community/guides/logging/golang-contextual-logging/

Not everything on context is logge handler has to pull it off. So libraries don't have a great way to add things?

## Handlers!


## Sinks
Sometimes you do just need a writer.


## awesome slog.
https://github.com/go-slog/awesome-slog
