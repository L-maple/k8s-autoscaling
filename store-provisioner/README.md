# autoScalingHdfs

## Step1. Run the pv_setup.sh to prepare volume
sh pv_setup.sh

## Step2. Deploy the local-volume-provisioner.yaml
kubectl apply -f local-volume-provisioner.yaml
