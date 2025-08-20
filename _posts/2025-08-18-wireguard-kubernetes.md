---
layout: post
title: wireguard kubernetes
tags: [ tech ]
---

What if you built your kubernetes cluster with wireguard/tailscale from  the ground up. 

* Api servers are always private except for wireguard port (use a side car)
* Head scale lets you join the cluster. (can we use certs for this?) 
* Nodes can are locked down except for wireguard port (use the top of rack)
* node wireguard nic get a /24 and we use a bridge to hand out pod ips.
* Do you need boostrapping tokens? do you need tls? Alot of tls becomes uncessary
* Can tailscale provide the load balancer service? 


