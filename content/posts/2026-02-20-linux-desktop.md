---
layout: post
title: "Linux desktop"
author: "Paul Miller"
tags: [ tech]
comments: true
date: 2026-02-06
---

On the first day after my last day at MSFT, I reimaged all my machines. Windows itself is not bad (SSH, PowerShell, and WSL are fine), but the junk our IT department installs definitely is.
So while I was reimaging, I decided it was time for Linux desktops on my three machines.

## [Bee / Nuc](https://www.amazon.com/dp/B0CTQBGPH9?ref_=ppx_hzsearch_conn_dt_b_fed_asin_title_2) + [Omakub](https://omakub.org/)

This has worked pretty well. Camera and audio are fine. I had to fight a Bluetooth mouse disconnect, but ChatGPT seemed to figure it out.
Chrome, VSCode, Spotify all seem fine.
Still using GNOME Terminal over Alacritty some of the time.
Omakub certainly has options.


##  [Omarchy](https://omarchy.org/) on my [pc part picker build (ryzen)](https://blog.northbriton.net/posts/2022-01-11-new-pc-build/)
This has been fun. Tiling window manager, no overlap. It is better at forcing you to organize things.
Browser tabs still get out of control though.
I have not dug into Arch/Pacman yet.
But it does seem to go down when not in use overnight, so I need to figure out uptime.



## Ubuntu 25 on Lenovo Snapdragon
[FAQ](https://discourse.ubuntu.com/t/faq-ubuntu-25-04-25-10-on-snapdragon-x-elite/61016)
Oof, this is pretty painful. Writing this because it just hard-rebooted twice, and ChatGPT and I could not find much in journalctl.
Suspend does not really work and drains your battery overnight.
Dealing with that by powering off on laptop lid close, but it sucks to lose anything you did not save.
Audio and camera don't work.
Even when not suspended, battery life does not seem quite as good.
No Chrome, which is mostly a problem since it is my password manager. I exported and imported, but it does not seem sustainable.

## Still in queue.
Bluefin/Fedora and ground-up NixOS. Also, do I take the Neovim pill?
