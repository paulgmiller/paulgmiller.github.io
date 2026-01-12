---
layout: post
title: why I'm leaving microsoft
tags: [ tech, family]
---


It's time for Paul to go. It's been 20 years, so probably inevitable, but I thought I'd write up what's changed.

## Why I am NOT going.
* **Product.** I love giving devs a managed k8s they can build things on. I use this product for my own hobbies and it's pretty great.
* **Money.** While Microsoft is getting outgunned in pay, I never felt treated unfairly.
* **To join FAANG/Mag 7.** I worry they are all the same and have the same issues below.

## The Year of Pain

Microsoft has always had a set of bad internal tools that never had external devs as a customer. ADO was a breath of fresh air >10 years ago, but we've failed to move to GitHub.  
Instead, in the face of quality and security, there has been a resurgence of Geneva, EV2, and one-branch tooling that is widely reviled and painful for engineering and no external customer would touch.   
At the same time, there have been several initiatives that put a lot of unfunded burden on engineering.
Where as before AKS biggest pain points had been internal/self inflicted and thus fixable ( adxmon replacing brittle prom->geneva, eno breaking up monoilith, validation spagheeti code refactor) the year of pain seems very much outside our control. 

### Release approvals/tooling
Despite having to get approvals on the last iteration of every PR and every release, we decided to add R2D, which might have been the most bureaucratic move I've seen at Microsoft.  
https://msazure.visualstudio.com/CloudNativeCompute/_wiki/wikis/PersonalPlayground/875657/R2D-after-1.5-weeks  

I dialed into a ScottGu meeting where several VPs explained why this is a bad idea (Brendan quote: R2D team has no incentive to engineer itself out of existence). But it went nowhere. 
Our dev infra team eventually hid this from us with leases (though at the cost of half-baked, dubious compatibility tests), but I still get lost souls asking me to do an R2D approval once a month for something I have zero context on.

Code of Conduct is a farce. It's just liability ass covering that makes the devs feel completely untrusted and forced to routinely make and normalize small violations. Will you do an honest post mortem if you violated this? 

EV2 has been a usability nightmare for years, but for us was basically a script runner under ADO that we'd try to ignore. However, this year we're basically forced to onboard to managed SDP despite numerous usability issues.  
When we talk with the team, they are at least upfront about, “ScottGu told us to focus on red tape and onboarding, and user experience is something we're just now getting to.”  
EV2 (and SRM) aren't just usability problems—they also make it hard to know what is successfully deployed where and cause SDP violations or delay releases with no notification.  
https://msazure.visualstudio.com/CloudNativeCompute/_wiki/wikis/PersonalPlayground/893656/Watermelon-logic
Our devinfra team tries to do a lot to cover here but now we just have TWO internal toolsets people have to learn. 

### Security Mandates
Security is good! MSFT should care more about security.  
But after Midnight Blizzard, we didn't really fund product teams—we put in mandates and hired PMs to chase them down.

Some make sense (MSI-only should have been policy since before 2020), but some of the most painful and always recurring—like fixing *every* CVE—are a big burden with minimal benefit.

The return of VPN requirements was particularly egregious to me. We moved away from VPN in 2010 to zero trust with the idea that corpnet (which has maybe 10× in size since) was fundamentally not trustworthy.  
But without justification we went back to that and broke many things (notably Kusto, with really confusing error messages) along the way. Nobody seems to own the VPN or can explain why I have to reconnect all the time.  
I found the Tailscale ban on personal corp-attached machines poorly defended.

### Identity

I log in dozens of times a day. I am so numbed to this that I will press accept on any dialog box you give me, which is its own security hole. The groups and accounts are still split across IDWeb and OneIdentity.  
Native passkeys don't work. For three years, intern new-account setup has failed in a different way.

JIT is a joke. Everyone approves anything. There are never audits. If R2D or security were serious, they would start with JIT audits, and every JIT would come with a work item.

Device control flow was also just “hey, we shut it off—no, we haven't really thought about how we do SSH.”  
Cloud Screen just disables your account before on-call because you missed a 5am Friday email, with no way to get an exception.

### Building open source

Dalec forces you to fork open-source build projects for FIPS and Azure Linux. Mostly listing this because it was yet another unfunded initiative where things were deprecated and we were just told to deal.

### Sovereign Clouds
The US and China clouds have always been a disaster, but dev mostly ignored them and assumed they were broken. Without heads-up, we were told three new clouds were coming with new nuances, and we were already “behind on the first one.”  
The second one started with no postmortem of problems with the first; instead, ICM was abused to get things fixed.  
https://msazure.visualstudio.com/CloudNativeCompute/_wiki/wikis/PersonalPlayground/893602/Cloud-Buildout-Bad-Engineering-Culture  

This completely broke some people who should be standing up for quality and security by normalizing one-offs, manual touches, and missing dependencies/testing/instrumentation.  
When you raise this issue, you just get “well, we signed a contract” and complete defeatism. This is again a ScottGu-level problem.

This short-sightedness immediately fell on on-calls like a hammer right before the holidays in a way that was all too predictable. 
Often this just doubled up existing weaknesses above. A new identity was created wh


## Attrition / Hiring
I've lost one report because “AKS is just too big and entering a different phase,” and a coworker to a non-product team because they “wanted to do real engineering work.”

But we've also lost a chunk of SDE IIs because we were pretty flippant about it being a bad hiring market, even though AI startups are hiring like crazy—especially for anyone with k8s experience.

The actual hiring experience is horrendous: recruiters dribble us resumes, hiring managers take the brunt of tech screening, and after we work through candidates we often lose on salary or just slowness.  
Again, there seems to be a sense that people should be desperate for jobs, but that is NOT the case for those with k8s experience.

Another core set of AKS people are going to try to do bare-metal Kubernetes, and I think that will get them around some of the pain I mention above—but if they are successful, they'll be right back in.

Hiring and retaining was already very hard before we entered this year of pain.

## Morale / expenses / travel / desks / devices

Most of these have been ongoing since 2020, but overlap with the year of pain as a constant baseline.

You'd think MSFT was a poor company from the 1950s. Being constantly told there's no morale budget for 5 years just feels like a fraud.  
When I paid for half my own ticket to KubeCon, I didn't bother expensing any meals because the tooling is so painful.

We've gotten two teams up since the pandemic despite encouraging remote work for 5 years, and now we've flipped to RTO. I don't actually mind RTO, but the flippancy with which it was done was astonishing.  
We literally announced RTO four months after we changed all intern offers to remote-only and killed relocation budgets.

Interaction with space and device teams also reflects the ongoing bureaucratic sclerosis. Trying to get an intern a desk is treated as “why would you bother the space team?”  
When you ask for a cheaper, better laptop you're told, “just be happy with what you have”  
Like internal tooling, we are not the customer—we are just cost.



