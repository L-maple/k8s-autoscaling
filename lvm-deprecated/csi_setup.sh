
sudo kubectl apply -f lvm-plugin.yaml
sudo kubectl apply -f lvm-provisioner.yaml
sudo kubectl apply -f rbac.yaml
sudo kubectl apply -f resizer/csi-resizer.yaml
