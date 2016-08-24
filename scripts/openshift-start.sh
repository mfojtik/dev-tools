#!/bin/bash
ip=$(/sbin/ip -o -4 addr list eth0 | awk '{print $4}' | cut -d/ -f1)

osdir="/var/lib/openshift"
mkdir -p ${osdir}
echo "--> starting origin server on (\"https://${ip}:8443\") ..."
echo
cd ${osdir}
exec openshift start --public-master=https://${ip}:8443 \
  --etcd-dir=/tmp/etcd  \
  --latest-images \
  --loglevel=8 \
  --volume-dir=volumes "$@" &> server.log
