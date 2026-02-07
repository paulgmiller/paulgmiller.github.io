---
layout: post
title: "First week no msft"
author: "Paul Miller"
tags: [family,tech]
comments: true
date: 2026-02-06
---

I am no lomnger employed by msft after 20 years of work and all school before hand.

## Monday

5k treadmill and
Kids drop off. Wife and me breakfast.
30 minutes of Careme (more e2e tests.)

9:30-4PM funeral

1 hour of omarchy/hyperland and minimal careme.

Pickup nate cooked dinner for nate.
Folded laundry.
Cooked dinner for wintana while she had wintana at soccer.
Beer pickles.

Healthcare. Family Calendar. Careme e2e test.

His/Him netfli

## Tuesday
Kids denitist and drop off.
Kids Bus and Othodonteist

Clerks - 6 hours!
First two hours codex completely failed.
Tried manually in afternoon and realized I needed frontend code and Clerks had no way to get a session cookie on dev ( :( ). Also clerks jwts might just be too short and I might need to make my own?
So really a design issue that codex/claude would happily just throw more shit at the wall with none of it sticking  https://github.com/paulgmiller/careme/pull/166

Kids pickup dinner (Mussels!) dominoes. House cleaning

Is Clerks the right choice? Shoudl I manage my own session rather than use javascript on every page?

Moved cloudflare over to LLC account

Watch stock market go down

Look at other auth sides  options in evening. (Stych, Authzero and and Entra)


## Wedneday

Beef stew prep with normal breakfast/dishes stuff. 

About 2 hours refactoring auth https://github.com/paulgmiller/careme/pull/166 so its not as clerk specific. Got happier

New lunch place Keg drop off. Birthday beer purchase

Ski bus. 6 hours

Eat stew. Watch pats 2015 superbowl run and read to kids. 

Watch stock market go down. 

Get beat by wife in shuffle puck and random tv.

## Thursday
Kids breafast/dishes/lunches

Careme
Dockerfile based tailwind - 30 minuts. https://github.com/paulgmiller/careme/pull/171
Deployed auth to production!!! (fixed signout and some other things along the way) https://github.com/paulgmiller/careme/pull/166 
Did dark mode then reverted 
Better readiness probe https://github.com/paulgmiller/careme/pull/183/changes
Wrote a tiny cli to mopve users and list ingredients

Talked about exchange funs with fidelity
Ate leftovers

After dinner did a night time run. 
Read to kids
Folded laundry and drank a beer. 
No computer activity in night time

# Friday

Kid breakfast and drop off. 
Samena Weight work out and Swim
QFC shopping bought to meals and lots of otehr kitchen stuff.
Set up parent teach confernce 

Careme. 
Came to conclusion Instacart does not want me intermediating stores and customers. You can search for ingredients but they give a html shopping list back for customer
Signed up for some sites that scrape albertsons and whole foods.
Lunch
Realized Walmart might have public api. Learned about affliate vs developer. Had codex write bad code on legacy api then got it to new one.
Figured out their strange auth. Tried to debug get stores on zip returning no results. Proved auth was right with taxonomy call but still not sure why no stores. 
So have a partially  working walmart client. https://github.com/paulgmiller/careme/pull/187

## GOals

Careme auth
Benji loft (nope)
Shed shelves (nope).
Figure out bus. (done)
Workout every day?
Lots of kids appointments 

