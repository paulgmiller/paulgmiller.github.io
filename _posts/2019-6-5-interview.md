---
layout: post
title: My Interview question
tags: [ tech ]
---

Find all pass codes for the old android unlock screen. Can't pass over the same node. 


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

Interviewee has to write this part. Whatever langage is fine.
If they're done we talk about using compute (parallism) or memory (memoizing) to speed it up. 

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

or golang

```go
package main

import (
	"fmt"
	"strings"
)

var graph = map[string][]string{
	"a": {"b", "d", "e"},
	"b": {"a", "c", "e"},
	"c": {"b", "e", "f"},
	"d": {"a", "e", "g"},
	"e": {"a", "b", "c", "d", "f", "g", "h", "i"},
	"f": {"c", "e", "i"},
	"g": {"d", "e", "h"},
	"h": {"g", "e", "i"},
	"i": {"h", "e", "f"},
}

func dfs(node, visited string, test func(string)) {
	newvisted := visited + node
	test(newvisted)
	for _, neighbor := range graph[node] {
		if !strings.Contains(newvisted, neighbor) {
			dfs(neighbor, newvisted, test)
		}
	}
}

func main() {
	for node := range graph {
		dfs(node, "", func (t string) {
			fmt.Printf(t + "\n")
		})
	}
}
```

