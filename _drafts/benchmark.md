pmiller@millerbuntu:~/aks-rp/resourceprovider/server$ time go build ./...

real    2m40.371s
user    13m53.231s
sys     2m5.590s



pmiller@millerbuntu:~/aks-rp/resourceprovider/server$ time go build ./...

real    2m40.956s
user    20m7.778s
sys     4m51.363s


millerarmcloud
-> % time go build ./...
go build ./...  
629.53s user 
68.39s system
918% cpu 
1:16.00 total


pmiller@millercloud4 [06:03:30] [~/aks-rp6/resourceprovider/server] [pmiller/cvebreakdown]
-> % time go build ./...
go build ./... 
291.98s user 
31.69s system 
826% cpu 
39.171 total

## With go clean --cache and go build -a 
pmiller@millercloud4 [06:06:53] [~/aks-rp6/resourceprovider/server] [pmiller/cvebreakdown]
-> % time go build -a ./...
go build -a ./...  
958.86s user
76.35s system 
1062% cpu 
1:37.42 total

pmiller@millerarmcloud [06:07:10] [~/aks-rp/resourceprovider/server] [master]
-> % time go build -a ./...
go build -a ./...  
630.92s user 
70.19s system 
940% cpu 
1:14.52 total

##with 32 arm cores

go build -a ./...  660.47s user 81.14s system 734% cpu 1:41.03 total
Trying again to see if cahing was in play
go build -a ./...  662.80s user 80.33s system 1233% cpu 1:00.25 total


