#!/bin/bash
waitseconds=${1:-1}
n=0
while :
do
	n=$(expr $n + 1)
	printf "[%06d] Trying to load for $IP" $n
	hey -c 10 -n 100 http://$IP/foo/bar
	sleep $waitseconds
done
