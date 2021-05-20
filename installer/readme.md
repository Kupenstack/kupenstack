# Quick-Start



Prerequisites:

* Kubernetes cluster.

### Installation

```shell
kubectl apply -f https://raw.githubusercontent.com/Kupenstack/kupenstack/main/config/demo/kupenstack-controller-manager.yaml
```

### Usage

Let try something simple like keypair.

While creating keypair we can specify our own public key or let kupenstack generate one for us. In this example we will use automatic generation. Create following yaml file:

```yaml
apiVersion: kupenstack.io/v1alpha1
kind: KeyPair
metadata:
  name: keypair-sample-2
  namespace: default
# we can also give our own public if we want
# spec:
#	publicKey: ssh-rsa AAAAB3NzaC1yc2EAsjadaskj
```

Now, lets apply it.

```
kubectl apply -f https://raw.githubusercontent.com/Kupenstack/kupenstack/main/config/samples/keypair-without-public-key.yaml
```

This should create our keypair at openstack.

```shell
kubectl get keypairs -o wide
```

```
NAME               IN-USE   AGE   PRIVATE-KEY
keypair-sample-2   false    4s   keypair-sample-2-b5wqw
```



