---
layout: post
title: Evictions KEP 
tags: [ tech ]
---



The KEP I want to write is that eviction attempts actually get recorded to pods than then deployment controller has leeway to surge in new pods when it sees pods are trying to evict. 
The pdb could then be kept in sync with deployments desired replicas and only allow evictions after deployer had surged ne new ready pods. 
Right now you basically have to overprovision to survive maintence events (upgrades, node freezes) and over provissionign is probably still the best idea because not all down events are predictable.
But I think there is a subnet of customer that don't want to overprovision for COGS reason but still want downtime only int he unpredictable cases.



### Why not events.
Events are much noisier. I theory deployment controller could watch events but I think that is going to cause a much bigger increase. Its already watching pods for 

### Doesn't solve draining
Endpoints controller does something with terminating 
https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/

https://kubernetes.io/docs/tutorials/services/pods-and-endpoint-termination-flow/

