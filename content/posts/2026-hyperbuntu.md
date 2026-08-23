---
layout: post
title: hyprbuntu
tags: [tech]
date: 2026-07-07
---


So omarchy 4.0 dropped but maybe I'm not ready to go all in [see](https://マリウス.com/a-word-on-omarchy/#summary). Kind of grown used the the boringness of ubuntu but want [hyprland](https://wiki.hypr.land/)

This is close but maybe too much https://wiki.hypr.land/Hypr-Ecosystem/hyprlauncher/


So considering a meta deb package to just let me install the folliowing items and setup 

hyprland
hypridle  
waybar
fuzzel (hyperlauncher) 

maybe hyprpolkitagent and xdg-desktop-portal-hyprland

hyperlock seems not be multi user so using this for now to switch users.

```
!/bin/sh

loginctl lock-session 
sleep 0.5

gdbus call \
  --system \
  --dest org.gnome.DisplayManager \
  --object-path /org/gnome/DisplayManager/LocalDisplayFactory \
  --method org.gnome.DisplayManager.LocalDisplayFactory.CreateTransientDisplay
```

Fuzzel:
```
[main]
font=Ubuntu:size=14
width=42
lines=9

horizontal-pad=24
vertical-pad=14
inner-pad=10

icons-enabled=yes
icon-theme=Adwaita

prompt=> 
placeholder=Search applications...

[colors]
background=111111f5
text=e6e6e6ff
prompt=9ca3afff
placeholder=6b7280ff
input=f5f5f5ff

match=ffffffff

selection=292929ff
selection-text=ffffffff
selection-match=ffffffff

counter=737373ff
border=3f3f46ff

[border]
width=1
radius=14
selection-radius=8
```

hyperidle:
```
-> % more .config/hypr/hypridle.conf
general {
    # When logind asks the session to lock, start hyprlock.
    lock_cmd = pidof hyprlock || hyprlock

    # If GDM/logind successfully unlocks the session, dismiss hyprlock.
    unlock_cmd = pkill -USR1 hyprlock

    # Always lock before the machine itself suspends.
    before_sleep_cmd = loginctl lock-session

    # Make sure monitors come back on after resume.
    after_sleep_cmd = hyprctl dispatch dpms on
}

# After 10 minutes: turn monitors off.
listener {
    timeout = 600
    on-timeout = hyprctl dispatch dpms off
    on-resume = hyprctl dispatch dpms on
}

# After 15 minutes: lock the session.
listener {
    timeout = 900
    on-timeout = loginctl lock-session
}
```

hyperland itself:
```
# See https://wiki.hypr.land/Configuring/Monitors/
#monitor=,preferred,auto,auto
monitor = DP-1, preferred, 0x0, 1
monitor = HDMI-A-1, preferred, 3840x120, 1, transform, 1

###################
### MY PROGRAMS ###
###################

# See https://wiki.hypr.land/Configuring/Keywords/

# Set programs that you use
$terminal = ghostty
$fileManager = nautilus 
$menu = fuzzel 
$browser = google-chrome

bind = SUPER, RETURN, exec, $terminal
bind = SUPER, SPACE, exec, $menu
bind = SUPER, B, exec, $browser

# Omarchy-style close window
bind = SUPER, W, killactive,
#ghetto
bind = SUPER, L, exec,  ~/.local/bin/switch-user 

```