#!/bin/bash

echo "---> Backuping tools and configuration ..."
echo -n "---> Powering down ..."
ssh -t dev -- sudo /usr/sbin/poweroff &> /dev/null
while true; do
  [[ -z "$(pidof xhyve)" ]] && break
  printf "." && sleep 1
done
echo
