---
layout: post
title: hyprbuntu
tags: [tech]
date: 2026-08-22
---


So omarchy 4.0 dropped but maybe I'm not ready to go all in. 
[Some interesting points here ](https://マリウス.com/a-word-on-omarchy/#summary). K
Kind of grown used the the boringness of ubuntu but want [hyprland](https://wiki.hypr.land/). 
But also want my kid to be able to use bone stock gnome on same machine. 

Seems like I can make ubuntu's login screen (gdm) launch hyperland for me and keep rest of the ssytem stick ubuntu. Lets see how this breaks 

https://github.com/paulgmiller/hybuntu

### omarchy 4 quatgro noites

tailscale (through omarchy)
golang (ui pacman got 1.27 where did it puit it)
docker (installed need to make it non sudo)
vscode (through omarchy )
jellyfin ( thorugh pacman -S jellyfin-server)
codex (though omarchy )
spotify (through omarchy, login failed )
ghosttty (ui pacman?)+
gh cli (sudo pacman -S github-cli)
zsh (ui pacman or  pacman -S zsh)
oh-my-zhs


changed monitors.lua

changed lookandfeel.lua

-- https://wiki.hypr.land/Configuring/Basics/Variables/#general
 hl.config({
   general = {
     -- No gaps between windows or borders.
     gaps_in = 0,
     gaps_out = 0,
     border_size = 0,

--     -- Change to niri-like side-scrolling layout.
--     layout = "scrolling",
   },
 })

 local omarchy_gdk_scale = 1
local omarchy_monitor_scale = 1

hl.env("GDK_SCALE", tostring(omarchy_gdk_scale))

hl.monitor({
  output = "",
  mode = "preferred",
  position = "auto",
  scale = omarchy_monitor_scale
})

hl.monitor({
  output = "DP-1",
  mode = "1920x1200@59.95",
  position = "0x0",
  scale = 1,
  transform = 1
})

