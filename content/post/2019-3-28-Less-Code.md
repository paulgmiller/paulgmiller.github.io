---
author: Paul Miller
comments: true
date: "2019-03-28T00:00:00Z"
tags:
- tech
title: Less Code
---
## Proposition
Value of a developer is features / (lines of code + number of dependenceis)
## Lessons
* Unittests can be bad if mocking balloons code size (also are you really testing your code or the mocks). Acceptance tests that are mostly setup are bad. Run end to end tests against a real environment.
* Dependencies on services are great. Libraries good. Not invented here is anathema. 
* Config flags/flighting while normally great can baloon code size if not cleaned aggressively.
* Yagni is importnant. Arictect Astronaughts not needed. Refactor to be general when you have > 2 users of code.
* Beware the young competent developer who goes and creates a whole new system from scratch
* microservices really need to have little boiler plate. This is part of the appeal of lambda/functions. 

