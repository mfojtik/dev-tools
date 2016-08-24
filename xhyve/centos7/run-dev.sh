#!/bin/bash

echo "---> Powering up dev box ..."

pushd ${HOME}/xhyve/centos7 >/dev/null
${HOME}/xhyve/centos7/run.sh
popd >/dev/null

echo -n "---> Waiting for dev box to become ready ..."
while true; do
  ssh -o ConnectTimeout=1 -o ConnectionAttempts=1 dev uptime &>/dev/null; retval=$?
  [[ "$retval" == "0" ]] && break
  printf "." && sleep 1
done
echo

setup_dir=${HOME}/go/src/github.com/mfojtik/dev-tools/xhyve/centos7
ansible-playbook -i ${setup_dir}/hosts -s ${setup_dir}/dev.yaml
