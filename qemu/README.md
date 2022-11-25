# QEMU script

This script is built to run on an UBUNTU 20.04 LTS machine. The script creates and runs a Linux filesystem image using QEMU:
- Downloads dependency tools
- Download stable Linux kernel and builds it with default config
- Create basic rootfs for the kernel image to boot from
- Create a basic `HELLO WORLD` init process on the rootfs
- Runs a QEMU machine using the built kernel and the rootfs
