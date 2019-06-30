---
layout: post
title: A CI/CD we can believe in
tags: [ tech ]
---

## Here to Help not to Hinder. 

Should be Faster than local and  Encourage rapid iteration. 

### Autocomplete and hope considered harmful.
* Changes goto master in the middle of the night. 
* Higher chance of resolve conflicts.
* PR Rot / comment angst / change averersion.

## Awesomely awful agents

### BuildTracker / Jenkins / Vsts agents
* I have some machines.
* Give me a script to run
* I’ll try and clean up. 
* If I’m awesome I’ll do some caching

### I can do anything! Oh no so can anyone else!

* base image bloat  / install 
* the sleep 100 task.
* What size should these agents be?
* How much parallelism is allowed?

## A Common Tragedy 

* Searchgold hilarity.
* Artificial limits
* Best of intentions. 
* Large organization fun.

## A private pool
Paul doesn't really want to run a build lab. But here we ares. Vm's are so 2011 lets use k8/docker

* [Base image](https://msdata.visualstudio.com/Vienna/_git/kubernetes-util?path=%2Fazureml-infra%2Fsrc%2Fdockeragent%2Fdockerfile&version=GBmaster)
* [Yaml](azureml-infra/deploy/templates/k8s/dockeragent.yaml) 
* Stateful vs stateless. (liveness probe)
* Autoscaling / Sizing. 
* Pat Secret
* Indrection

### A yaml aside 
* pipelines look pretty neat.
* Nothing is hidden
* Can change build/release atomically with change and put in pr
* More secure/auditable / One branch overlap.
* Waitin on approval tasks?

## Agentless Tasks.
* Been here before.
* Total control
* Limitless parallelism. 
* Agentless K8's? Agentless docker build? 
