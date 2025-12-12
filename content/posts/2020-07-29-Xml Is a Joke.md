---
layout: post
title: Xml is a Joke
tags: [ tech ]
date: 2020-07-29
---
Way back in the ancient times of 2006 I had a boss who wrote a xml sax parser in C++ because libraries were hard i guess.
I got to debug bugs and use it and evntually wrote a object representation on top of it which used more memory and was 10% less awful to use.
But everything was terrible so xml probably just didn't have good libraries right.

Today I had to edit someone's app.config and append to xml attributes that were really maps[string][]string / Dictionary<string,List<string>> that were hand serialized as"foo:bar,bunk;moo:brown,cow"

Yuck. But this was just some awful dev right? So I aked to friends how they would serialize a map dict in xml.

I got 

```
<MyAwesomeMap>
  <Entry>
    <Key>my key</Key>
    <Value>
      <List>
        <Entry>First one</Entry>
        <Entry> Second one </Entry>
      </List
    </Value>
  </Entry>
</MyAwesomeMap>
```
and 
```
<maptype>
  <entry key="key">value</entry>
</maptype>
```

Okay those are better but certainly not consisteny. What does my old buddy powershell do


```
> @{foo="paul"} | ConvertTo-Xml -as string
<?xml version="1.0" encoding="utf-8"?>
<Objects>
  <Object Type="System.Collections.Hashtable">
    <Property Name="Key" Type="System.String">foo</Property>
    <Property Name="Value" Type="System.String">paul</Property>
  </Object>
</Objects>
```

Woof that's not great here's json for comparsion

```
> @{foo="paul"} | ConvertTo-Json
{
  "foo": "paul"
}
```

What aboutu go and python? Turns out their standard library marshallers won't even take maps/dicts. 
    
    
https://stackoverflow.com/questions/30928770/marshall-map-to-xml-in-go
    
https://pypi.org/project/dicttoxml/

Maps/dicts seem like pretty fundamental data types. Maybe XML was just never good all along.
