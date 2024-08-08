---
layout: post
title: Evictions KEP 
tags: [ tech ]
---
Update: Had an intern try this: https://github.com/javier090/k8s-pdb-autoscaler

Need to actually fit thisin this format
https://github.com/kubernetes/enhancements?tab=readme-ov-file#is-my-thing-an-enhancement

## Problem
THere are often situtions where users are unwilling to overprovision deployments. But they are willing to surge pods for deployments to avoid unncessary disruptions. However there is no surge during maitnance events so that either creates downtime or 

The simplest example of this is a replica count of 1 and max surge 1 on a deployment but also a Pod Distription Budge with a MinAvailable:1. I can now never do mainence on the node in question. But this app is different from one who says maxunavailble : 0 indicating it needs special handling its max surge already shows intention it can be disrupted as long as a new pod is surged in first. 

Why would someone do this? Why not just overprovision?  Well they may be offering a free sample product where keepign COGS low is acutally more improtant than 100% up time but that doesn't mean they want downtime for every maintence event. A temporary surge is still cheap. 

You could also have cases where you have 2 replicas but one of them is unready for reasons specific to it (was hit by a ramdom cosmic ray) and in that case the its actually better to surge in new pods to replace the unready one. 

## Proposal

Eviction attempts/failures actually recorded to pods as a status field.  (they are already events see below)
Then deployment controller has leeway to surge a new pods for each pod it sees with failed evictions.
The pdb could then be kept in sync with deployments desired replicas and only allow evictions after deployer had surged new new ready pods. 

### Do we surge on first eviction failure? 
Not clear maybe we should wait till there are a couple to make sure an unready pod resolves it self? API server could know if tehre are ohter unready replicas in pdb or not and give more context but that seems more complicated. 

### Why not events.
Events are much noisier. I theory deployment controller could watch events but I think that is going to cause a much bigger increase. Its already watching pods for 

### Why not watch pdbs themselves
Since its based on labels the knowing does this pdb just apply to this deployment is hard. 


### This may also help with not dropping traffic 
[Endpoints controller tries to remove terminating pods](https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/)
but again if you're current replicas are handling all the traffic they can handle (aggressive hpa settings or jsut currently running hot) then you may need surged nodes to be ready to  handle existing traffic. 


