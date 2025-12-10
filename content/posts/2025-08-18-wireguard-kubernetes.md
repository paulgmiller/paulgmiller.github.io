---
layout: post
title: wireguard kubernetes
tags: [ tech ]
---

What if you built your kubernetes cluster with wireguard/tailscale from  the ground up. 

* Api servers are always private except for wireguard port
     * Use a side car for each api server pod https://tailscale.com/kb/1185/kubernetes#sample-sidecar. Drop a tailnet secret in ccp.
     * like api server vnet integration pods show up in your tailnet/overlay even if they have no nodes. Can be hit directly with endpoint slices 
     * No konnectivity proxy ingress and egress goes direct to api server pods
     * Could still have a internal load balancer for legacy nodes but ideally want none of that.
     * Rest of ccp does not need to be in tailnet but could be
* Head scale lets you join the cluster. (can we use certs to verify this public key is from a baremetal node? We can't be the ones giving them a secret?)
* Nodes  are locked down except for wireguard port (use the top of rack)
  * node wireguard nic get a /24 and we use a bridge to hand out pod ips. (could use cilium or soemthign else if you want net pols but don't need it for ipam)
  * node ip should be just the first ip. 
* Do you need boostrapping tokens? If you're on the network is that enough? Alot of tls becomes uncessary
* Can tailscale provide the load balancer service? https://tailscale.com/kb/1223/funnel to get things to kubenet.

Can we prove this out by 
1) sidecarring api servers on a standalone
2) Creating a BYO CNI cluster
3) Adding a Demonset to VM nodes that boostraps tailscale and gets a /24 + brige plugin
4) Create a script that adds us to tailscale then gets a bootstrap token and joins api server from some random home machine? 

