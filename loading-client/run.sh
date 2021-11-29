#!/bin/bash

while :
do
	hey -c 2 -n 200 http://$IP
	sleep 1
done
