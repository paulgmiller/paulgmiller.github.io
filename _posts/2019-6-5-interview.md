﻿---
layout: post
title: My Interview quesition
tags: [ tech ]
---

Find all pass codes for the old android unlock scrren.
More to come 

```
# a - b - c
# | \ | / |
# d - e - f
# | / | \ |
# g - h - i
# 
graph = {}
graph['a'] = ['b', 'd', 'e']
graph['b'] = ['a', 'c', 'e']
graph['c'] = ['b', 'e', 'f']
graph['d'] = ['a', 'e', 'g']
graph['e'] = ['a', 'b', 'c', 'd', 'f', 'g', 'h', 'i']
graph['f'] = ['c', 'e', 'i']
graph['g'] = ['d', 'e', 'h']
graph['h'] = ['g', 'e', 'i']
graph['i'] = ['h', 'e', 'f']
```

Interviewee has to right this part. Whatever langage is fine.
```python
def find_all(a_graph, test):
  for g in a_graph:
      dfs(g, '', test)
      
def dfs(node, visited, test):
  newvisted = visited + node
  test(newvisted)
  for neighbor in graph[node]:
      if neighbor not in visited:
          dfs(neighbor, newvisted, test)
          
find_all(graph, print)  
```          

Here's a coworker's [java solution](https://gist.github.com/maxormo/84fdbecbbf328544eaa5c3aa464b6674)

Here's an iterative solution 
```python          
def dfs_iter(node,  test):
  stack = [node]
  while stack:
    path = stack.pop()
    test(path)
    for neighbor in graph[path[-1]]:
      if not neighbor in path:
        stack.append(path + neighbor)
```
