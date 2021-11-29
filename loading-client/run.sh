#!/bin/bash
n=0
waitseconds=${1:-1}
echo $seconds
while :
do
	n=$(expr $n + 1)
	printf "[%06d] Trying to load for $IP" $n
	hey -c 2 -n 200 http://$IP
	sleep $waitseconds
done
