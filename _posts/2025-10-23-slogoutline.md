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
```
[*](https://dave.cheney.net/2015/11/05/lets-talk-about-logging)
```
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
  *   Pass it everywhere? Save it as a member?
  *   Shove it on context
  *   Have seen all of these 

* How do log respectfully with out overlogging?
    
* What if I want to customize my log format?
    * Json / Binary
    * [Coloring](https://gitlab.com/greyxor/slogor) / Emoji✅ 

## What about X
* zap and logrus are pretty neat. They introduced level, structure, sinks and alot of the stuff I am going to talk about below. 
* But if I'm writing a library how do I know what interface to get/log to?
* have seen several instances where logging adapters are regularly rewritten. Usually losing something each time. 
* Sometimes its nice to just have the standard library 


## Enter Slog
[Gopkg](https://pkg.go.dev/log/slog@go1.25.4) and [Blog](https://go.dev/blog/slog)

* Package and Type level log functions.
   *   Most call into [func Log(ctx context.Context, level Level, msg string, args ...any)](https://pkg.go.dev/log/slog@go1.25.4#LogAttrs)
* A high level [interface](https://pkg.go.dev/log/slog@go1.25.4#Handler) you can impment behind. (Handlers) 
* A text and json implementation
* Set Default !! (Claim: just use the slog.Info/Error methods instead of passing loggers around most of the time)

## Levels
[Debug/Info/Warn/Error](https://pkg.go.dev/log/slog@go1.25.4#Level)
* Throw all your junk into debug! Main.go can [set level](https://pkg.go.dev/log/slog@go1.25.4#SetLogLoggerLevel) and toss it. 
* indecisive? Warn away!
* make up your own log levels and go crazy. 

```
	slog.Debug("Hello, 世界")
	slog.Info("I will see you in ", "time", time.Hour)
	slog.Warn("It's getting hot in here")
	slog.Error("Nuts", "error", fmt.Errorf("next time gadget"))
	slog.Log(context.TODO(), slog.Level(10), "more imporant than an error"
```
[*](https://go.dev/play/p/fy06nhHWPpl)
```
  2009/11/10 23:00:00 INFO I will se you in  time=1h0m0s
  2009/11/10 23:00:00 WARN It's getting hot in here
  2009/11/10 23:00:00 ERROR Nuts error="next time gadget"
  2009/11/10 23:00:00 ERROR+2 more imporant than an error
```

## Structure
* String parsing stinks! We have kusto and lots of other fancy databases!
* Request ids and other tracing items.  Latencies as durations! Metadata!
* Attrs has lots of nice built in types and use and there's groups. (I have not used groups)
* Globally setting attributes on default or a particular log instance. With/WithAttrs.

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
* Okay what if I want to add an attributes that all things below me in the call stack can use? [Chatgpt](https://chatgpt.com/share/691d5e6c-8bc0-8009-b6b6-de2b28c6ea7e)
* Well I can stick it on the context but then the code I am calling has to pull them out.
* There is a InfoContext and Equivalent for each helper. 
* I didn't find a great solution for this unless you control the handler too.
* Then you can use https://github.com/PumpkinSeed/slog-context. [Example blog](https://betterstack.com/community/guides/logging/golang-contextual-logging/)
* Written by someone named PumpkinSeed so it must be good. I kind of wanted a standard lib way to do it though. 

## Handlers!
* The fools gave us an [interface](https://pkg.go.dev/log/slog@go1.25.4#Handler): 	Handle(context.Context, Record) error
* Maybe you don't want json maybe you want bson or thrift or whatever compact nonense.
* Maybe to forward records to kusto, blog, kafka, service bus.
* You can do whatever you want! 
* Sometimes you might just want to do an io.Writer. But then you lose per record control. 
* My Azure append blog example.
    * [Started as a handler](https://github.com/paulgmiller/careme/blob/8621cce4cd5f240d881a0bf4e22cb26c281c8e8e/internal/logsink/appendblob.go
    * HIt an interesting deadlock
    * [Went to a writer](https://github.com/paulgmiller/careme/blob/master/internal/logsink/appendblob.go) much simpler  but regret it because now I might split log records

## awesome slog.
* Lots more stuff here https://github.com/go-slog/awesome-slog
* Want to log test to console but put json to soemwhere else? [slog-multi!](https://github.com/samber/slog-multi) 
* Want to log directly to team? (WHY!) [Here's a forwarder](https://github.com/samber/slog-microsoft-teams). Maybe datadog or kafka make more sense.
* Miss Logrus? [There's an adapter](https://github.com/go-slog/awesome-slog?tab=readme-ov-file#adapters)
