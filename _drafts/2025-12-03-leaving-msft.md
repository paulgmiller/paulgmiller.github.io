---
layout: post
title: why I'm leaving microsoft
tags: [ tech, family]
---


It's been 20 years. So probably inevitable but thought I'd write up whats changed.

## Why I am NOT  going.
* Product. I love giving devs a managed k8s they can build things on. I use this product for my own hobbies and its pretty great
* Money. While microsoft is getting outgunned in pay I never felt treated unfairly. 
* To join FANG/Mag 7. I worry they are all the same and have same issues below. 

## The year of Pain

Microsoft has always had a set of bad internal tools that never had devs as a customer. ADO was a breath of fresh air > 10 years ago but we've failed to move to github.
Instead in the face of quality and security there has a been a resurgence of geneva, and ev2, one branch tooling that is widely reviled and painful for engineering. 
At the same time there have been seveal initiatives that put alot of unfunded burden on engineering. 

### Release approvals/tooling. 
Despite having to get approvals on the last iteration of every PR and every release we decided to add R2D which might have been the most beauracratic move I've seen at miscrosoft. 
https://msazure.visualstudio.com/CloudNativeCompute/_wiki/wikis/PersonalPlayground/875657/R2D-after-1.5-weeks
I sat in a ScottGu meeting  where several vp's exposed how this is a bad idea (Brendan qoute: R2D team has no incentive to engineer it self out of existance). Our dev infra team eventually hid this from us with leases. 
(though at the cost of half baked, dubious compatibility tests) but I still get losts souls asking me to do an R2D approval once a month on something I have zero context on. 

Code of Coduct is a farse. It just liability coverage that makes the devs feel compeletely untrusted and forced to routinely make and normalize small violations. 

Ev2 has been a usability nightmare for years but for us was basically a script runner under ADO that we'd try and ignore. Hoever this year we're basically forced to onboard to managed sdp despite numerous usability issues. 
When we talk wit the team they are at least upfront about "Scott Gu told us to focus on Red tape and onboarding and user experience is somethign we're just now getting to". 
ev2 (and srm) aren't just usabity problems they also make it  hard to know what is successfully deployed where and violate sdp or delay releases with no notification
https://msazure.visualstudio.com/CloudNativeCompute/_wiki/wikis/PersonalPlayground/893656/Watermelon-logic



### Security Mandates.
Security is good! Msft should care more about security. 
But after midnight blizzard we wen didnt really fund product teams we put in mandates and hired PMs to chase them down. 
Some make sense. (MSI only should have been policy since before 2020) but some of the most painful and always reoccuring like fixing EVERY CVE is a big burden with minimal benefit. 

The return of VPN requirements was particular egregious to me. We moved away from VPN in 2010 to zero trust with the idea that corpnet (which has maybe 10x in size since) was fundamentally not trustable. 
But we went back to that and broke many things (notably kusto with really bad error messages) along the way.  Nobody seems to own the VPN or can explain why I have to reconnect all the time. 
I found the tailscale ban on personal corp attatched machines poorly defended. 



### Identity

I login dozens of times a day. I am so numbed to this I will press accept on any dialog box you give me which it own security hole. The groups and accounts are still split across idweb and one identity. 
Native Passkeys don't work. For three years intern new account setup has failed in a different way. 

JIT is a joke. Everyone approves anything. There are never audits. If r2d or security was serious they would start with jit audits and every jit would come with a work item. 

Device control flow was also just "hey we shut if off no we haven't really thought about how we do ssh" 
Cloud screen just disables your account before on call because you missed a 5am friday email with no way to get an exception. 

### Building open souce

Dalec forces you to fork open soruces build projects for fips and azure linux.  Mostly listening this because it was yet another unfunded initiative where things were deprecated abd we were just told to deal. 

### Sovriegn Clouds
The US and china clouds have always been a disaster but dev mostly ignored them and assumed they were broken but without heads up we were told 
3 new clouds were coming with new nuances and we were already "behind on first one" second one started with no postmortem of problme with first instead ICM was abused to get things fixed. 
https://msazure.visualstudio.com/CloudNativeCompute/_wiki/wikis/PersonalPlayground/893602/Cloud-Buildout-Bad-Engineering-Culture
This completely broke some people who should be standing up for qualty an security by normalizing one offs, manual touches and missing depedndencies/testing/instrumentation. 
When you raise this issue you just get "well we signed a contract" and complete defeatism. This is a again a Scott Gu level problem. 

This short sightness immediately fell on on calls like a hammer right before the holidays in a way that was all to predictable. 


## Attrition  / Hiring. 
I've lost one report because "AKS is just too big and entering a different phase" and  coworker to a non-product team because they "wanted to do real engineering work"

But we've also lost a chunk of SDE II's because we were pretty flippant about it being a bad hiring market even though AI startups are hiring like crazy especially for anyone with k8s experience.

The actual hiring experience is horrendous as our recruiters dribble us resumes, hiring manager takes the brunt of tech screening and after we work through candidates we often lose on salary or just slowness. 
Again there seems to be a sense that people should be desperate for jobs but that is NOT the case for those with k8s experience. 

Another core set of AKS people are going to try and do baremetal kubernetes and I think that will get them around some of the Pain I mention above but if they are successful they'll be right back in.

Hiring retaining was already very hard before we entered this year of pain. 

## Morale/expenses/travel/desks/devices.

Most of these have been ongoing since 2020 but overlap with year of pain as constant baseline. 

You'd think msft was a poor company from the 1950's. Being constantly told theres no morale budge for 5 years just feels like a fraud. When I did pay for half my own ticket to kubecon I didn't bother expensing any meals becasue the tooling is so paulful. We've gotten two teams up since the pandemic despite encourageing remote work foir 5 years and now we've flipped to ROT. I don't actually mind RTO but the flippancy with which it was done was astonishing. We litterly announced rto 4 months after we changed all intern offers to remote only and killed relocaiton budget. 

Interaction with space and device teams also reflects the on going beuaracratic sclerosis. Trying to get a intern a desk is treated as "why would you bother the space team!" When you ask for a cheaper better laptop you're told, "just be happy with what you have". Like internal tooling we are not the customer we are just cost.  