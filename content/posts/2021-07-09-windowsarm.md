---
layout: post
title: Windows Arm
tags: [ tech ]
date: 2021-07-09
---

## WHy am I doing this?
* Battery
* Cellular connectivity
* Cost
* dev works happens on azue vm through vs code remote

## What's broke
* System32\openssh is there in cmd.exe but not "windows powershell". POwershell core is fine. [issue](https://github.com/PowerShell/Win32-OpenSSH/issues/1566)
* Chocolatey install infiniteloops.
* Git for windows is emulated (not a big deal?) [issue](https://github.com/git-for-windows/git/issues/2346)
* WSL from windows store. Trying manual instructions
  

## What's fine
* edge
* windows terminal
* vscode

## Haven't tried.
* Windows 11