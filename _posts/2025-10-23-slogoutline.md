---
layout: post
title: slog is neat
tags: [ tech ]
---

Why did golang get a new logging framework after ~10 years? 

## Existing ["log"]([url](https://pkg.go.dev/log)) package
* It's like printf but there's a timestamp!
* I can call Fatal and Panic; I guess that's nice in main.

```
log.Println("Hello, 世界")
log.Printf("I will see you in %s\n", time.Hour)
log.Fatal("Oh no it's all gone wrong")
```
[*](https://dave.cheney.net/2015/11/05/lets-talk-about-logging)
```
2009/11/10 23:00:00 Hello, 世界
2009/11/10 23:00:00 I will see you in 1h0m0s
2009/11/10 23:00:00 Oh no it's all gone wrong
```
  
* There is [SetOutput](https://cs.opensource.google/go/go/+/refs/tags/go1.25.3:src/log/log.go;l=364); let's come back to this later.
* I can get a default logger, but I can't set a default logger.

## What should libraries do?
This approach is pretty great if you're in main.go, but what if you're writing a package for reuse?

* Should libraries log at all?
  * Well, I can pass up errors (wrap them if you're extra cool!).
  * What about that thing that seems wrong but I'm suppressing? I don't want to break compatibilities, but I'm pretty sure this is wrong.
  * https://github.com/Azure/eno/pull/532

* How do I get the logger?
  *   Use a package-level variable or function?
  *   Pass it everywhere and save it as a member?
  *   Shove it onto context?
  *   I have seen all of these approaches.

* How do I log respectfully without over-logging?
    
* What if I want to customize my log format?
    * JSON or binary
    * [Coloring](https://gitlab.com/greyxor/slogor) or emoji ✅

## What about X
* zap and logrus are pretty neat. They introduced level, structure, sinks, and a lot of the stuff I am going to talk about below.
* But if I'm writing a library, how do I know what interface to get/log to?
* I have seen several instances where logging adapters are regularly rewritten, usually losing something each time.
* Sometimes it's nice to just have the standard library.


## Enter Slog
[Gopkg](https://pkg.go.dev/log/slog@go1.25.4) and [Blog](https://go.dev/blog/slog)

* There are package- and type-level log functions.
   *   Most call into [func Log(ctx context.Context, level Level, msg string, args ...any)](https://pkg.go.dev/log/slog@go1.25.4#LogAttrs)
* There is a high-level [interface](https://pkg.go.dev/log/slog@go1.25.4#Handler) you can implement behind (handlers).
* It includes text and JSON implementations.
* You can set the default! (Claim: just use the slog.Info/Error methods instead of passing loggers around most of the time.)

## Levels
[Debug/Info/Warn/Error](https://pkg.go.dev/log/slog@go1.25.4#Level)
* Throw all your junk into debug! Main.go can [set level](https://pkg.go.dev/log/slog@go1.25.4#SetLogLoggerLevel) and toss it.
* Indecisive? Warn away!
* Make up your own log levels and go crazy.

```
	slog.Debug("Hello, 世界")
	slog.Info("I will see you in ", "time", time.Hour)
	slog.Warn("It's getting hot in here")
	slog.Error("Nuts", "error", fmt.Errorf("next time gadget"))
	slog.Log(context.TODO(), slog.Level(10), "more important than an error"
```
[*](https://go.dev/play/p/fy06nhHWPpl)
```
  2009/11/10 23:00:00 INFO I will see you in  time=1h0m0s
  2009/11/10 23:00:00 WARN It's getting hot in here
  2009/11/10 23:00:00 ERROR Nuts error="next time gadget"
  2009/11/10 23:00:00 ERROR+2 more important than an error
```

## Structure
* String parsing stinks because we have Kusto and lots of other fancy databases!
* Request IDs, other tracing items, latencies as durations, and general metadata all live happily here.
* Attrs have lots of nice built-in types, and there are groups (I have not used groups).
* You can set attributes globally on the default or a particular log instance with With/WithAttrs.

```
slog.With("author", "pmiller")
	slog.Info("I will see you in ", "time", time.Hour, slog.Duration("strongtypedtime", time.Minute), slog.Float64("pi", math.Pi))
	attrs := []slog.Attr{
		slog.String("author", "jsonpaul"),
		slog.String("program", "playground"),
	}
	handler := slog.NewJSONHandler(os.Stdout, nil).WithAttrs(attrs)

	slog.SetDefault(slog.New(handler))
	slog.Info("I will see you in ", "time", time.Hour, slog.Duration("strongtypedtime", time.Minute), slog.Float64("pi", math.Pi))
```
[*](https://go.dev/play/p/6UbejghCS5a)
```
  2009/11/10 23:00:00 INFO I will see you in  time=1h0m0s strongtypedtime=1m0s pi=3.141592653589793
  {"time":"2009-11-10T23:00:00Z","level":"INFO","msg":"I will see you in ","author":"jsonpaul","program":"playground","time":3600000000000,"strongtypedtime":60000000000,"pi":3.141592653589793}
```

### Putting things on context?
* Okay, what if I want to add attributes that everything below me in the call stack can use? [ChatGPT](https://chatgpt.com/share/691d5e6c-8bc0-8009-b6b6-de2b28c6ea7e)
* I can stick them on the context, but then the code I am calling has to pull them out.
* There is an InfoContext and equivalent for each helper.
* I didn't find a great solution for this unless you control the handler too.
* Then you can use https://github.com/PumpkinSeed/slog-context. [Example blog](https://betterstack.com/community/guides/logging/golang-contextual-logging/)
* Written by someone named PumpkinSeed, so it must be good. I kind of wanted a standard library way to do it, though.

## Handlers!
* The fools gave us an [interface](https://pkg.go.dev/log/slog@go1.25.4#Handler): 	Handle(context.Context, Record) error
* Maybe you don't want JSON; maybe you want BSON or Thrift or whatever compact nonsense.
* Maybe you want to forward records to Kusto, Blob storage, Kafka, or Service Bus.
* You can do whatever you want! 
* Sometimes you might just want to do an io.Writer, but then you lose per-record control.
* My Azure append blob example.
    * [Started as a handler](https://github.com/paulgmiller/careme/blob/8621cce4cd5f240d881a0bf4e22cb26c281c8e8e/internal/logsink/appendblob.go
    * Hit an interesting deadlock
    * [Went to a writer](https://github.com/paulgmiller/careme/blob/master/internal/logsink/appendblob.go) much simpler but regret it because now I might split log records

## Awesome slog
* Lots more stuff here: https://github.com/go-slog/awesome-slog
* Want to log tests to console but put JSON somewhere else? [slog-multi!](https://github.com/samber/slog-multi)
* Want to log directly to Teams? (WHY!) [Here's a forwarder](https://github.com/samber/slog-microsoft-teams). Maybe Datadog or Kafka make more sense.
* Miss Logrus? [There's an adapter](https://github.com/go-slog/awesome-slog?tab=readme-ov-file#adapters)
