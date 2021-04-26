---
layout: post
title: Distributed Twitter.
tags: [ tech ]
---

##Update:
Wrote this over here: https://github.com/paulgmiller/zebu

## Why 
* Should we be concerned with facebook/twitter massively booting people after the 1/6 resurrection, or aws booting Parler?  Maybe not?
* Should we worry about [Myanmar](https://www.engadget.com/myanmar-internet-instagram-twitter-shutdown-175341979.html) or the [great china firewall](https://www.technologyreview.com/2020/07/08/1004876/the-internet-is-changing-drastically-for-hong-kongs-citizens/#:~:text=The%20fallout:%20Effectively,%20this%20brings%20Hong%20Kong%20into,are%20permitted%20to%20operate%20only%20if%20they%20comply.)? Maybe yes? 
* I think [futurama](https://theinfosphere.org/Old_Man_Waterfall) covered some of the tradeoffs pretty well.

I'm going to ignore this though and plow through to a solution because like 3d printed guns I think they might be inevitable.

In the US you can obviously set up a blog and serve it up from your home machine but it's not trivial and you can get DOS's pretty trivially. So lets solve that problem too. 

## Content.
[IPFS](https://ipfs.io/) mostly has this solved. I can take any content and publish an immutable hash and anyone else on the planet as nodes share around tables of where to find content. If nobody else looks at my content it dissapears pretty easily but as soon as other people look a it 

So your tweets would essentiall be a linked list of immutable hash documetns.
```
paul@millernuc:~$ ipfs add
Welcome to distributed twitter
added QmZruVCZVLKgMokmgr45CTMsYs8u2v1FgWraR32S1dHbLW QmZruVCZVLKgMokmgr45CTMsYs8u2v1FgWraR32S1dHbLW

paul@millernuc:~$ ipfs add
prev: QmZruVCZVLKgMokmgr45CTMsYs8u2v1FgWraR32S1dHbLW
----
Was this a mistake what have I done?

added QmbX8JVs8jWdLFroZ72y9ydK1aRimJoB3dYYTq55SmVAHR
```

You can go edit past tweets by creating new versions of every tweet since but it won't remove them. Here I've used a header to point to the previous post. This could obviously be json/yaml to hold other metadata. The content below the ---- I would just store markdown.

To update new content you'd use [ipns](https://docs.ipfs.io/concepts/ipns/) to point the hash of the key you own to your head post

```
paul@millernuc:~$ ipfs name publish --key=wilkes /ipfs/QmbX8JVs8jWdLFroZ72y9ydK1aRimJoB3dYYTq55SmVAHR
Published to k51qzi5uqu5dl59p3uol4d71qpeapvvtc98napcdfehkavsquau1svu8xia8rs: /ipfs/QmbX8JVs8jWdLFroZ72y9ydK1aRimJoB3dYYTq55SmVAHR
```

Here's the above 
* https://cloudflare-ipfs.com/ipns/k51qzi5uqu5dl59p3uol4d71qpeapvvtc98napcdfehkavsquau1svu8xia8rs
* https://cloudflare-ipfs.com/ipfs/QmbX8JVs8jWdLFroZ72y9ydK1aRimJoB3dYYTq55SmVAHR


## Finding People. 
Okay obviously people aren't great at finding/remembering large hashes so this all still seems dumb. DNS solves this in the real world and you could [use dns to solve this](https://docs.ipfs.io/concepts/dnslink/#dnslink) like so https://cloudflare-ipfs.com/ipns/dtwitter.swiperank.com

But china is going to block your dns server and a legitimate dns registrar may not want to have anything to do with your rantings (actually curious if they can legally deny registrations in US).

Enter [Ethereum Naming Service](https://docs.ens.domains/). The Blockchain is good at establishing global consensus on who controls said digitial asset. So it can tell us who owns the right to point a name at a ipns hash. 

Now I can just tell people to find me at https://cloudflare-ipfs.com/ipns/northbriton.eth. Further more I control subdomains so I can give away "handles" to randos or verified real meat people depending on what I want to do. 

But I should probably tell who I follow so instead of just having my ipns point to my head post I can actually point to a document that is the list of people I follow and my head post. Couple advantages of making this public.
* I can browse a "feed" of those I follow even if I don't have access to my private data. 
* my client can mirror the content of those I follow automatically. 
* it could also mirror the content of my follow's follows creating a decent amount of redundancy. 

I can always create a seperate account that follows people I don't want to make public if you want to be a coward. But you could also have seperate accounts for differnt interest groups. 

But now I have a way to find my and others posts so we have the start of somethng

## Who will curate the feed!

Well nobody, the default view is just a chronological posts from people I follow. But all the info is public so any enterprising go-getter can generate a feed they think will be useful for you. They can crawl all the other handles in ENS and tell you who is @'ing your or retweeeting you. They can use fancy machine learning or they can have susan curate something. They will give you ads or charge a fee but if you don't like it there should be healthy competition. If you are just posting boring stuff you can pay a service to "pin" your content so old posts aren't lost. 

There's no direct messsaging in the current conception though giving people you follow access to your signal handle seems possible. 

Anonymous accounts with few followers should pretty much never show up in your feed. I am not sure why twitter does this except people actually like to argue with trolls.

## Problems.
### Identity. 
People are not good managing key pairs. My naive approache would be to recommend/demand you sync it between two machines. If you lose it you have to point at your last post and tell people about your new handle. 

### App stores / phones
Turns out people use phones more than computers. App stores could just block this.
Then you're stuck with the web which is more open but a tighter sandbox. 
On old boring PC's it's harder to block. Your ipfs node could even download new upadtes to its client. 
There is already ipfs nodes that run on phones and a java script ipfs node infact. 
For ENS to work you'd need some sort of ethereum light client to go find ipns/ipfs addresses. Seems possible but haven't verified it exists. 

### Blocking all content from "outside"
I'm not actually sure how national governments shut down the internet. I'm guessing they mostly block dns, but they could shut down incoming pipes. 
I am really hoping satallite internet becomes cheap enough that when governments shutdown the pipes people just open connections to the stars and share. 

