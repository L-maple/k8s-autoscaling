# dynamic-provisioner
提供了csi-provisioner的本地方式，只要在本目录运行 kubectl apply -f . 即可。

然后运行example/下里的部署文件，kubectl get pv 可看到自动配置的pv了。
