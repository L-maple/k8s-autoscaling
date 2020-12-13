# format && mount
MOUNT_POINT=/root/keti1-ruc/disks
DISK=/dev/mapper/centos-home

umount $DISK
mkfs.ext4 $DISK
DISK_UUID=$(blkid -s UUID -o value $DISK)
mkdir -p $MOUNT_POINT/$DISK_UUID
mount -t ext4 $DISK $MOUNT_POINT/$DISK_UUID
echo UUID=`sudo blkid -s UUID -o value /dev/mapper/centos-home` $MOUNT_POINT/$DISK_UUID ext4 defaults 0 2 | sudo tee -a /etc/fstab

# Create multiple directories and bind mount them into discovery directory
#for i in $(seq 1 5); do
#  sudo mkdir -p ${MOUNT_POINT}/${DISK_UUID}/vol${i} ${MOUNT_POINT}/disks/${DISK_UUID}_vol${i}
#  sudo mount --bind ${MOUNT_POINT}/${DISK_UUID}/vol${i} ${MOUNT_POINT}/disks/${DISK_UUID}_vol${i}
#done

# Persistent bind mount entries into /etc/fstab
#for i in $(seq 1 5); do
#  echo ${MOUNT_POINT}/${DISK_UUID}/vol${i} ${MOUNT_POINT}/disks/${DISK_UUID}_vol${i} none bind 0 0 | sudo tee -a /etc/fstab
#done
 
