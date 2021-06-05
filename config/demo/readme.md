# Quick-Start

**Note: This demo is not an complete implementation.**

### Before you begin

To follow this demo, you need:

* A Kubernetes cluster with each node having minimum 8gb ram.

For this demo, we will use KinD tool to create cluster, although any Kubernetes cluster should work.

<details>
    <summary>Click here for KinD cluster creation instructions.</summary>
    <p>
    <b>Download Kind on linux:</b>
        <br>
        <code>curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.11.0/kind-linux-amd64</code><br>
        <code>chmod +x ./kind</code><br>
        <code>mv ./kind /bin/kind</code><br>
    </p><br>
    <p>
        <b>Now create cluster:</b><br>
        <code>kind create cluster</code>
    </p>
</details>

Note: You can try this demo in a VM on any public cloud like AWS, GCP, If you do not want to try on your system.



Note: If you are not using single node KinD k8s-cluster, then you will need to update Confimap with names of controller nodes and compute nodes in the below deployment. Names will be same as node name in kubernetes cluster. `kubectl get nodes`

### Installing KupenStack

```shell
kubectl apply -f https://raw.githubusercontent.com/Kupenstack/kupenstack/main/config/demo/kupenstack-controller-manager.yaml
```

Wait for some time, till our deployment gets completed and our kupenstack-controller-manager becomes ready. It should take about 11 minutes on good internet connection.

**Verify:**

```
kubectl get pods -n kupenstack-control
NAME                                             READY   STATUS    RESTARTS   AGE
kupenstack-controller-manager-85ff8ddb64-rf6qx   1/1     Running   0          12m
```

Now, once deployment is complete, you can try `kubectl get pods -n kupenstack`. The output will be full containerized OpenStack deployment running on Kubernetes. 

```
NAME                                           READY   STATUS      RESTARTS   AGE
glance-api-6c9c8bbcdc-t4q78                    1/1     Running     0          53m
glance-bootstrap-89qpd                         0/1     Completed   0          53m
glance-db-init-c7tb2                           0/1     Completed   0          53m
glance-db-sync-7p5f5                           0/1     Completed   0          53m
glance-ks-endpoints-l68hx                      0/3     Completed   0          53m
glance-ks-service-rft5d                        0/1     Completed   0          53m
glance-ks-user-pklqf                           0/1     Completed   0          53m
glance-metadefs-load-dqjsl                     0/1     Completed   0          53m
glance-rabbit-init-7lwtx                       0/1     Completed   0          53m
glance-storage-init-t5sdp                      0/1     Completed   0          53m
horizon-7b8497bf88-g5qgl                       1/1     Running     0          55m
horizon-db-init-sdpmp                          0/1     Completed   0          55m
horizon-db-sync-k6m4f                          0/1     Completed   0          55m
ingress-6cbc96b8fd-s2c4n                       1/1     Running     0          60m
ingress-error-pages-755bf859fd-rl6nr           1/1     Running     0          60m
keystone-api-879d7748f-5vh4v                   1/1     Running     0          57m
keystone-bootstrap-wtdsm                       0/1     Completed   0          56m
keystone-credential-setup-xnk8j                0/1     Completed   0          57m
keystone-db-init-bstc6                         0/1     Completed   0          57m
keystone-db-sync-pr6kg                         0/1     Completed   0          57m
keystone-domain-manage-bk6pr                   0/1     Completed   0          56m
keystone-fernet-setup-9dtkj                    0/1     Completed   0          57m
keystone-rabbit-init-zvdkj                     0/1     Completed   0          56m
libvirt-libvirt-default-kd4dn                  1/1     Running     0          52m
mariadb-ingress-65d4fd8d6f-6mm7j               1/1     Running     0          60m
mariadb-ingress-error-pages-854c4d5469-r8ls6   1/1     Running     0          60m
mariadb-server-0                               1/1     Running     0          60m
memcached-memcached-654f7f6956-vqhdg           1/1     Running     0          58m
neutron-db-init-9m799                          0/1     Completed   0          51m
neutron-db-sync-zk7n7                          0/1     Completed   0          51m
neutron-dhcp-agent-default-txvw6               1/1     Running     0          51m
neutron-ks-endpoints-ggl4h                     0/3     Completed   0          50m
neutron-ks-service-2lsqd                       0/1     Completed   0          50m
neutron-ks-user-55b2z                          0/1     Completed   0          50m
neutron-l3-agent-default-st898                 1/1     Running     0          51m
neutron-lb-agent-default-txm7m                 1/1     Running     0          51m
neutron-metadata-agent-default-tnrxg           1/1     Running     0          51m
neutron-netns-cleanup-cron-default-k4h4l       1/1     Running     0          51m
neutron-rabbit-init-657sn                      0/1     Completed   0          51m
neutron-server-f6d997db6-mhdk9                 1/1     Running     0          51m
nova-api-metadata-6b994f4f69-k68gg             1/1     Running     1          52m
nova-api-osapi-588d44bc9b-qkgs6                1/1     Running     0          52m
nova-bootstrap-kslqp                           0/1     Completed   0          52m
nova-cell-setup-bdflf                          0/1     Completed   0          52m
nova-compute-default-zfgx8                     1/1     Running     0          52m
nova-conductor-64fb679b4b-s5gkd                1/1     Running     0          52m
nova-consoleauth-7ffd74c78c-qr424              1/1     Running     0          52m
nova-db-init-9qh98                             0/3     Completed   0          52m
nova-db-sync-bg29p                             0/1     Completed   0          52m
nova-ks-endpoints-6zg5z                        0/3     Completed   1          52m
nova-ks-service-vnf55                          0/1     Completed   0          52m
nova-ks-user-kfbst                             0/1     Completed   0          52m
nova-novncproxy-6df88bdcf7-s8g7t               1/1     Running     0          52m
nova-rabbit-init-kzl7v                         0/1     Completed   0          52m
nova-scheduler-5d55b8d878-9v5t2                1/1     Running     0          52m
placement-api-7ff8b4f57c-t598n                 1/1     Running     0          52m
placement-db-init-krx6z                        0/1     Completed   0          52m
placement-db-sync-rh84c                        0/1     Completed   0          52m
placement-ks-endpoints-t7z9q                   0/3     Completed   0          52m
placement-ks-service-q22n2                     0/1     Completed   0          52m
placement-ks-user-m98th                        0/1     Completed   0          52m
rabbitmq-cluster-wait-md2gb                    0/1     Completed   0          58m
rabbitmq-rabbitmq-0                            1/1     Running     0          58m
```



### Create some resources

Let us create a Virtual Machine using KupenStack.

Apply following yaml file to Kubernetes:

```yaml
apiVersion: kupenstack.io/v1alpha1
kind: Image
metadata:
  name: image-sample
spec:
  src: http://download.cirros-cloud.net/0.5.1/cirros-0.5.1-x86_64-disk.img
  format: raw

---
apiVersion: kupenstack.io/v1alpha1
kind: Flavor
metadata:
  name: flavor-sample
spec:
  vcpu: 2
  ram: 500
  disk: 1
  rxtx: "1.0"

---
apiVersion: kupenstack.io/v1alpha1
kind: VirtualMachine
metadata:
  name: virtualmachine-sample
spec:
  image: image-sample
  flavor: flavor-sample
```

Apply:

```bash
kubectl apply -f https://raw.githubusercontent.com/Kupenstack/kupenstack/main/config/demo/vm-sample.yaml
```

### Know your deployment

Now, let us explore what we did.



* `kubectl get virtualmachines` this will show

```
NAME                    NODE                 STATE   AGE
virtualmachine-sample   kind-control-plane   BUILD   12s
```

Our Virtual machine is in build state, let us wait for some time and check again.



* `kubectl get vm -o wide`

```
NAME                    NODE                 STATE     NETWORKS(IP)            AGE
virtualmachine-sample   kind-control-plane   Running   default(10.10.1.190)    62s
```

Now, our Virtual machine is in running state. We can also see that VM is in `default` network and has IP `10.10.1.190` . Remember we have not deployed any networks. KupenStack by default keeps all VM in `default` network if network is not specified in VM Manifest. Since, we had not created any network called `default ` KupenStack automatically created it for us. We can check it by:

* `kubectl get networks`

```
NAME      IN-USE   AGE
default   true     56s
```



* `kubectl get flavors` gives us following output:

```
NAME            IN-USE   AGE
flavor-sample   true     38s
```



* `kubectl get images` give us output:

```
NAME           IN-USE   READY   AGE
image-sample   true     true    43s
```

Note: our Image is in Ready state. It means it has been downloaded and uploaded to glance successfully.



* `kubectl get keypairs -o wide`

```
NAME             IN-USE   AGE     PRIVATE-KEY
keypair-sample   true     50s     keypair-sample-542zr
```

Note: we didn't specified any public key while creating key-pair so KupenStack automatically generated one and is exposing it through a secret. We can copy our private-key from this secret and then delete this secret.



* If we look at our Namespaces, then we can see there are annotations for corresponding Project in OpenStack. For example: `kubectl describe ns default` give us:

```
Name:         default
Labels:       <none>
Annotations:  kupenstack.io/external-project-id: 4c7c184c14314cb08244a1fcc47d0bf5
              kupenstack.io/external-project-name: default-5fg8k
Status:       Active

No resource quota.

No LimitRange resource.
```

In the annotations we can see a Project named `default-5fg8k` has been created in OpenStack which maps to this Namespace in Kubernetes. Every resource that is created in this Namespace will be created in corresponding Project at OpenStack.



We can cross-verify our above deployments in Horizon also, Go to port `32020 ` on the node's ip in browser:

Now, you can login and check, (Credentials are: Domain: "default", user: "admin", password: "password")

![img](horizon-screenshot.png)

