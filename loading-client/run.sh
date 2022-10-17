#!/bin/bash
n=0
waitseconds=${1:-1}
echo $seconds
while :
do
	n=$(expr $n + 1)
	printf "[%06d] Trying to load for $IP" $n
	hey -c 1 -n 100 http://$IP/foo/bar
	sleep $waitseconds
done
