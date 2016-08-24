#!/bin/sh

# Linux
KERNEL="./vmlinuz-3.10.0-327.el7.x86_64"
INITRD="./initramfs-3.10.0-327.el7.x86_64.img"
CMDLINE="earlyprintk=serial quiet console=ttyS0 acpi=off root=/dev/mapper/centos-root rd.lvm.lv=centos/root rd.lvm.lv=centos/swap rw"

MEM="-m 5G"
SMP="-c 4"
NET="-s 2:0,virtio-net"
#IMG_CD="-s 3,ahci-cd,/Users/will/Downloads/CentOS-7-x86_64-Minimal-1511.iso"
IMG_HDD="-s 4,virtio-blk,./hdd.img"
PCI_DEV="-s 0:0,hostbridge -s 31,lpc"
LPC_DEV="-l com1,stdio"
#ACPI="-A"
UUID="-U deadbeef-dead-dead-dead-deaddeafbeef"

# Linux
sudo xhyve $ACPI $MEM $SMP $PCI_DEV $LPC_DEV $NET $IMG_CD $IMG_HDD $UUID \
  -f kexec,$KERNEL,$INITRD,"$CMDLINE" &> vm.log &
