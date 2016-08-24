#!/bin/bash

echo "--> cleanup up staled kubernetes mount points ..."
mount | grep kubernetes | cut -d ' ' -f 3 | xargs sudo umount -f &>/dev/null

set -e

source_dir="/data/src/github.com/openshift/origin"
osdir="/var/lib/openshift"
etcdir="/tmp/etcd"

echo "--> cleanup up previous origin data ..."
ip=$(/sbin/ip -o -4 addr list eth0 | awk '{print $4}' | cut -d/ -f1)
rm -rf ${osdir} && mkdir -p /var/lib/openshift
rm -rf ${etcdir} && mkdir -p /tmp/etcd

sudo /home/mfojtik/bin/openshift-start.sh & ospid=$?

echo "--> waiting for origin api server ..."
set +e
while true; do
    curl --max-time 2 -kfs https://${ip}:8443/healthz &>/dev/null
    if [ $? -eq 0 ]; then
        break
    fi
    sleep 1
done
set -e

pushd ${osdir} >/dev/null
master="openshift.local.config/master"

chmod a+rwX ${master}/admin.kubeconfig
chmod +r ${master}/openshift-registry.kubeconfig \
         ${master}/openshift-router.kubeconfig
export CURL_CA_BUNDLE=${master}/ca.crt
sudo chmod a+rwX ${master}/admin.kubeconfig

echo "--> installing docker registry ..."
oadm registry \
    --latest-images \
    --namespace=default \
    --config=${master}/admin.kubeconfig

echo "--> creating haproxy router ..."
oadm policy add-scc-to-user hostnetwork \
    --serviceaccount=router \
    --config=${master}/admin.kubeconfig
oadm router \
    --latest-images \
    --config=${master}/admin.kubeconfig \
    --service-account=router

echo "--> waiting for openshift namespace to be ready ..."
set +e
while true; do
    oc get namespace/openshift --config=${master}/admin.kubeconfig &>/dev/null
    if [ $? -eq 0 ]; then
        break
    fi
    sleep 1
done

set -e
echo "--> importing image streams ..."
oc create \
    -f ${source_dir}/examples/image-streams/image-streams-centos7.json \
    -n openshift  \
    --config=${master}/admin.kubeconfig

echo "--> granting view access to test-admin user ..."
oadm policy add-role-to-user view test-admin --config=${master}/admin.kubeconfig

oc login ${ip}:8443 --certificate-authority=${master}/ca.crt -u test-admin -p test &>/dev/nul
oc new-project test --display-name="Dev sample" --description="Dev project" &>/dev/null
oc logout
popd >/dev/null
echo
echo "To login:"
echo "oc login ${ip}:8443 --certificate-authority=/var/lib/openshift/${master}/ca.crt -u test-admin -p test"
echo
echo "Waiting for openshift to finish bootstraping ..."
sleep 5 && kill ${ospid}
echo
echo
