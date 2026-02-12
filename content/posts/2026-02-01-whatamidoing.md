---
layout: post
title: "First Week: No MSFT"
author: "Paul Miller"
tags: [family, tech]
comments: true
date: 2026-02-06
---

I am no longer employed by MSFT after 20 years of work, plus all the school beforehand.
## Week 1

### Monday

- 5k treadmill  
- Kids drop-off  
- Breakfast with my wife  
- 30 minutes of Careme (more end-to-end tests)

**9:30–4:00 PM:** Funeral

- 1 hour of Omarchy/Hyprland and minimal Careme work  
- Picked up Nate; cooked dinner for Nate  
- Folded laundry  
- Cooked dinner for Wintana while she was at soccer  
- Beer pickles

Healthcare admin, family calendar, Careme e2e tests.

Netflix (His/Him).

### Tuesday

- Kids dentist and drop-off  
- Kids bus and orthodontist

**Clerk — 6 hours.**  
The first two hours Codex completely failed. Tried manually in the afternoon and realized I needed frontend code, and Clerk has no way to get a session cookie in dev 😞. Also, Clerk JWTs might just be too short-lived, and I may need to mint my own.

So really a design issue that Codex/Claude would happily throw more stuff at without any of it sticking:  
https://github.com/paulgmiller/careme/pull/166

- Kids pickup, dinner (mussels!), Dominoes, house cleaning

Is Clerk the right choice? Should I manage my own session rather than use JavaScript on every page?

- Moved Cloudflare over to the LLC account  
- Watched the stock market go down  
- Looked at other auth options in the evening (Stytch, Auth0, Entra)

### Wednesday

- Beef stew prep plus normal breakfast/dishes stuff

- About 2 hours refactoring auth so it’s less Clerk-specific — felt better afterward:  
  https://github.com/paulgmiller/careme/pull/166

- New lunch place, keg drop-off, birthday beer purchase  
- Ski bus — 6 hours

- Ate stew  
- Watched the Pats’ 2015 Super Bowl run and read to the kids

- Watched the stock market go down  
- Got beaten by my wife at shuffle puck and random TV

### Thursday

- Kids breakfast, dishes, lunches

**Careme work:**
- Dockerfile-based Tailwind — 30 minutes  
  https://github.com/paulgmiller/careme/pull/171
- Deployed auth to production!!! (fixed sign-out and other issues along the way)  
  https://github.com/paulgmiller/careme/pull/166
- Did dark mode, then reverted  
- Better readiness probe  
  https://github.com/paulgmiller/careme/pull/183/changes
- Wrote a tiny CLI to move users and list ingredients

- Talked about exchange funds with Fidelity  
- Ate leftovers

After dinner:
- Nighttime run  
- Read to kids  
- Folded laundry and drank a beer  
- No computer activity at night

### Friday

- Kids breakfast and drop-off  
- Sammamish weight workout and swim  
- QFC shopping: bought two meals and lots of other kitchen stuff  
- Set up parent–teacher conference

 **Careme:**
- Concluded Instacart does not want me intermediating stores and customers. You can search for ingredients, but they give back an HTML shopping list for the customer.
- Signed up for sites that scrape Albertsons and Whole Foods.
- Poke lunch and a walk.
- Realized Walmart might have a public API. Learned about affiliate vs. d

## Week 2 
Over weekend cooked an excellent careme dinner for friends. 
https://careme.cooking/recipe/le76qpHPhb-EK2KB7TqGZA==
https://photos.app.goo.gl/83E1cXchei89Qsg77

Applied to Y combinator and started locations ui update

### Monday. 
* Samena workout. Squat and bike.
* Finished auth bug
* locations ui to prod.
* Started better wine pairing work.
* two dinners and lanundry.

### Tuesday
* Pancake Breakfast for kids
* Cache talk on excahnge funds (20% in realestate due to irs rules)
* Got questions on recipes working though not enough context and bad reload behavioru https://github.com/paulgmiller/careme/pull/203
* 4 mile run and met ruben for lunch. Totally failed to find good food rubens gave us bad seattle dogs and nachos.
* Fraield to trademark, Created many issues most non tehnica, Minor updates, SEO research, sitemap started. 
* Colins pub
* Missed bus. Careme steelhead (slaw great, trout meh)
* India scraping compay did not schedule a 9pm meeting as promised



