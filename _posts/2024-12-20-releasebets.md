---
layout: post
title: Release Bets
tags: [ tech ]
---

## The Problem 
So you're a big wig at some tech company and you want to make a big jump in quality. A common and reasonable sounding thing to woud be to bring in some bar raisers. You don't have enough people you trust to do it at the Pull Request level so you decide to have them review releases yet. (Continuous Deployment isn't happening for .... reasons....). Now you hae an interesting incentive problem. 

Bar raisers don't have much incentive to get the code out but they have a big inentive not to let a bad hange out. [Type I vs Type II errors](https://en.wikipedia.org/wiki/Type_I_and_type_II_errors). The bar raisers don't have much signal about what releases/changes are dangerous or not so they're gonna want to talk about each change. You'll end up spending a constant amount of time on low risk releases eating up alot of the review time budget. Higher risk release might get LESS disussion just beause low risk items ate up available time. So how do you get a better signal?

## Enter the bet

If you want to make a release you should give the bar raisers odds their release will cause an outage (defining outages is hard but lets assume we have one). If I'm really save I migh say this release only has a 1% chance of causing an outage (anyone that says 0% is a fool). Then the bar raiser now has some information and if they want to talk about the issue they need to ante a minimum amount of money to the bet (lets say $10). If after disucssion the odds don't change the bet is locked. If the odds change (a good thing new information being exposed!) then releaser forfeits the ante and the bar raiser can still take whatever odds that come out if they want. This both focuses the discussions using time most effetively and aligns incentives to bring out truth rather than block or override. 

This works best for single change releases. If you do bulk you need PR reviewers to all give preditions and aggregate them which is harder but still doable. This is basically a prediction market. It works better the more people that can participate and the more liquid the market is. Ideally I'd be able to go bet against the release of some division far from me if I have reason to doubt. Insider Trading is a goodness here!
