---
date: "2018-11-18T00:00:00Z"
tags:
- beer
title: Brewie Repair
---

So follow up on the [Brewie Review]({{ site.baseurl }}{% post_url 2018-6-2-Brewie %}). Was having a blast. Just orded a double recipe from [morebeer](https://morebeer.com). But while my sculpin clone was mashing, the whole brewie died. Changing outlets didn't help but unlike during ordering customer support was pretty helpful. Their first guess was a fuse so I grabbed a 15 amp ceramic microwave fuse (should be good for higher temps) at home depot and cracked open the case. (about 12 short screws and 4 feet have to be removed, not too bad)

![brewie under the hood](https://images.weserv.nl/?w=650&url=https://pk70ha.dm.files.1drv.com/y4mDUSCLRR9-Wj1mi539ASJ8bDkeaPM9WUKqANL_IF7-bA3vO65d4IfEU7xVd7ZJVNOJSD0zhjKC8k6NTrgIa1T2a5hdMDTlCfA36XnlCNBqU4id9tUMed2OpQqfFj7xNnsK7t58__bvfVopwNEWAkp2C9blC672we7hyUMCP5N6eLoCzOYnB6UJ8KtDGvsn9O59OywtfkRVsr6ST69iHYiOA?width=4032&height=3024&cropmode=none)

Some thoughts:
* Circuit board is fairly accesible.
* Tubing already had some residue (black stuff). Might be replacable.
* can see the chiller path pretty easily. If this is a complete disaster least I could steal the chiller. 
* The heaters and pumps are harder to get to. (hard to tell pumps from plumbing servos)

Sure enough fuse was blown here's a blown one next to the other ones.

![fuse](https://images.weserv.nl/?w=650&url=https://nwev7w.dm.files.1drv.com/y4mPBo8Q8BdILjtsKlMlVuV-J69cNppliqR1Be9ucOraJmorP603lUNoMmoMR9eM0Zx6xd4zGltB_91tn-cyWcvBDkKRy_ga3P1FeeWy0B30uaqgfr7UE3REV72jB6HViBu1n4O53lhDV5nVX0OaG1kgt2pcRW_OXWdCYULpQm_okP_rTIFmxNfzuGHgv1qjkYttX_CkbLBhDilg30yLfaNsg?width=3024&height=4032&cropmode=none)

Booted right back up but shortly after smelled a burning smell. Boil heater seemed on. Restarting and running a stanizing cleaning cycle didn't help. But again customer support was helpful (takes 24 hours they are in hungary). Said likely they SSR (solid state relay) was at fault. These are basically the high voltage electronic switches for the heaters. They shipped a new one but also when I asked linked to [where I could get one on amazon](https://www.amazon.com/Omron-G3NA-220B-DC5-24-Indicator-Phototriac-Isolation/dp/B003B2Z0N6/ref=sr_1_fkmr0_1?ie=UTF8&amp;qid=1542305128&amp;sr=8-1-fkmr0&amp;keywords=20+amp+omron+ssr). That's great because I like fixing things myself. I even found [another version](https://www.amazon.com/gp/product/B074FT4VXB/ref=oh_aui_detailpage_o00_s00?ie=UTF8&psc=1) that was 10 dollars instead of 60 (buyer beware.)

You can see the relay in the brewie on the lower right here:

![relay](https://images.weserv.nl/?w=650&url=https://o9lleg.dm.files.1drv.com/y4mRTSNZLRYP-82IKXuUj8CbQZL9Fo3LXbm0JNg2mhYqs3kC2rZiBMaWsIrM2bWlXYTacUVuPT9UwbMw37Ps7N-3KKIVIrKUlv6qZltLbm-z7wLyQEu3lvp7JIGLCJib_uRoO4Jk-D2s1F1v1gbwAOjJl_g3ljkPOuyHpYDlTiA1gD_f7jdp63wV62cG-VKlB10S3USuBsp_Wxn2vu4SmhVAw?width=4032&height=3024&cropmode=none)


Relays showed up and did fix the problem (amazon one also installed fine but holding for future). Removing and reinstalling the relay to the mount was pain since the electroics box is super small.
Also customer serive told me the wrong one to replace so deduced it with a multimeter. Really want to draw up a circuit diagram and replace this with a rasberry pi and external relays, pumps and tubing.
[This looks](https://www.raspberrypi.org/blog/homebrew-beer-brewing-pi/) [pretty](https://opensource.com/article/17/7/brewing-beer-python-and-raspberry-pi) [intersting](https://web.craftbeerpi.com/) 

Either way trying to brew soon.




