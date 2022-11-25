#!/bin/bash
set +x

# install dependencies
sudo apt install -y \
    git \
    make \
    build-essential \
    flex \
    bison \
    libelf-dev \
    libssl-dev \
    qemu

# download and build stable linux kernel
git clone git://git.kernel.org/pub/scm/linux/kernel/git/stable/linux-stable.git
cd linux-stable
make defconfig
make -j $(nproc)

cd -

# create rootfs for the kernel image
mkdir -p rootfs
cd rootfs
mkdir -p bin dev etc lib mnt proc sbin sys tmp var
cd -

# download busybox to be able to run basic commands
curl -L 'https://www.busybox.net/downloads/binaries/1.26.2-defconfig-multiarch/busybox-x86_64' > rootfs/bin/busybox
chmod +x rootfs/bin/busybox

# create init bin to run after kernel boot
cat >> rootfs/init << EOF
#!/bin/busybox sh

/bin/busybox mount -t devtmpfs  devtmpfs  /dev
/bin/busybox mount -t proc      proc      /proc
/bin/busybox mount -t sysfs     sysfs     /sys
/bin/busybox mount -t tmpfs     tmpfs     /tmp

/bin/busybox echo "HELLO WORLD!"
/bin/busybox sh
EOF
chmod +x rootfs/init

# create cpio image from the rootfs
cd rootfs
find . | cpio -ov --format=newc | gzip -9 >../initram.img
cd -

# Emulate qemu machine that uses the built kernel image and boots from the initram.img image
qemu-system-x86_64 -kernel linux-stable/arch/x86/boot/bzImage \
                   -initrd initram.img \
                   -append "console=ttyS0" \
                   -nographic
