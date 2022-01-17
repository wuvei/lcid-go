#!/bin/bash

for pid in $(pgrep -f go);
do
{ kill $pid && wait $pid; } 2>/dev/null
done

cd /root/lcid-go

nohup /usr/local/go/bin/go run ./main.go > /dev/null 2>&1 &