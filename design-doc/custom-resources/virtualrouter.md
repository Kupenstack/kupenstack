# VirtualRouter

* [Summary](#Summary)
* [Motivation](#Motivation)
* [Design Details](#Design-Details)
  * [API](#API)
  * [Overview](#Overview)
    * [ExternalNetwork](#ExternalNetwork)
    * [VirtualNetworks](#VirtualNetworks)
    * [ExtendScope](#ExtendScope)
    * [Routes](#Routes)
  * [Creation Considerations](#Creation-Considerations)
  * [Deletion Considerations](#Deletion-Considerations)
  * [Updation Considerations](#Updation-Considerations)
* [Other implementation considerations](#Other-implementation-considerations)
* [Functioning at OpenStack](#Functioning-at-OpenStack)

### Summary

This document covers design specification, functionality, details for **VirtualRouter** custom resource(CR) in KupenStack. A Virtual Router CR manages router resources in the associated OpenStack cluster.

### Motivation

In networking, routers are used to provide connectivity between hosts on different networks. They decide how a packet should be routed. In KupenStack there are various VirtualNetworks(VN) separated by namespaces. VirtualMachines(VM) connected to these VN may require to communicate with VMs in another VN. A Virtual Router can solve this problem. 

Also, a VM may want to connect to an external network, a Virtual Router can solve this problem as well.

Therefore, it is required to have a Virtual Router CR that can dynamically connect to these VNs and an external network to provide required connectivity to workloads connected to them.

### Design Details

#### API

```yaml
apiVersion: kupenstack.io/v1alpha1
kind: VirtualRouter
# shortName=vr

metadata:
  # scope=Namespaced
  name: sample-router
  namespace: default

spec:

  # Name of the External Network that the router uses as its gateway.
  # required=false, type=string
  externalNetwork: provider-net
  
  # Virtual Networks that should connect to this router.
  # required=false, type=object
  virtualNetworks:
    selector:
      # MatchLabels takes multiple key:value labels. These labels are used to
      # search the required virtual networks to connect.
      matchLabels:
    	key: value
    	key2: value
  
  # ExtendScope takes a list of the names of additional namespaces inside which router
  # should search vitual networks, in addition to its own namespace.
  # required=false, type=array
  extendScope:
    - "backend-ns"
    - "database-ns"
    - "dev-ns"
  
  # Routes defines the static routes.
  # required=false, type=array
  routes:
    - destinationCidr: "10.10.100.0/24"
      nextHopIP: "172.168.10.5"

status:
  
  # Id of router at openstack cloud.
  # type=string
  id: a29d8-1d73n-2dw45-h4hr2
  
  # Whether virtual router is ready for use or not.
  # type=boolean
  ready: true
  
  # Contains list of all the virtual networks connected to this virtual router.
  # type=array
  virtualNetworkList: [ "Namespace/Name" ]
  
  # IP address of virtual router from external network
  # type=string
  exteranlIP: "172.29.249.149"

```

**Output on `kubectl get virtualrouters` or `kubectl get vr`**

```
NAME            EXTERNAL-NETWORK   EXTERNAL-IP      AGE
sample-router   provider-net       172.29.249.149   47m
```

**Output on `kubectl get vr -o wide`**

```
NAME            EXTERNAL-NETWORK   EXTERNAL-IP      AGE
sample-router   provider-net       172.29.249.149   47m
```

#### Overview

Virtual Router CR is a namespaced resource. It can connect to many virtual networks and only one external network. Virtual Router(VR) uses the external network as its default gateway. By default the scope of VR is limited to the namespace it is created in.

###### ExternalNetwork

External networks are special networks in KupenStack. These are cluster scoped. A VR can connect to anyone from the available external networks. This external network is then used as the default gateway by this VR.

###### VirtualNetworks

VR can connect to multiple virtual networks. This field uses the approach of labels and selectors to dynamically connect to all the available virtual networks with matching labels in its scope. By default the scope of VR is limited to the namespace it is created in.

###### ExtendScope

This field takes a list of namespaces and extends the scope of VR from the current namespace to the new namespaces specified in the list. This may be required to provide connectivity between VN from different namespaces.

###### Routes

This field takes a list of static routes to use for this virtual router. Every route is defined by two fields `destinationCidr` and `nextHopIP`. The `destinationCidr` defines the IP range that the packet may be destined to and `nextHopIP` tells the IP to which this packet should be routed to.

#### Creation Considerations

All fields in a VR are optional, hence VR can be created with minimal configuration as follow:

```yaml
# Minimal Virtual Network example
apiVersion: kupenstack.io/v1alpha1
kind: VirtualRouter
metadata:
  name: sample-router
---
```

Suppose the goal is to route traffic between workloads of two VNs namely, `my-vn1` and `my-vn2` in `ns1` then we need to create a VR that connects these two VN. First, make sure that both VNs have some common labels that can be used to select them let's say `router`=`my-vr1`. Now create a VR in the same namespace `ns1` that tries to connect to these two VN as follow:

```yaml
# VR that connects to all VN with labels router=my-vr1 (VN in same namespace)
apiVersion: kupenstack.io/v1alpha1
kind: VirtualRouter
metadata:
  name: sample-router
  namespace: ns1
spec:
  virtualNetworks:
    selector:
      matchLabels:
        router: my-vr1
---
```

Now, our two VNs have a virtual router connecting them.

Let us say that `my-vn2` is in a different namespace called `ns2` instead of `ns1` then the above VR custom resource has to be slightly modified. We need to add `ns2` in the scope of our VR as follow:

```yaml
# VR that connects to all VN with labels router=my-vr1 (VN in different namespace)
apiVersion: kupenstack.io/v1alpha1
kind: VirtualRouter
metadata:
  name: sample-router
  namespace: ns1
spec:
  extendScope:
    - "ns2"
  virtualNetworks:
    selector:
      matchLabels:
        router: my-vr1
---
```

#### Deletion Considerations

VR does not have any pre-clean-up requirements. Deleting a VR deletes it. But it is to be noted that if a VR is deleted then all VN connected to it will loose connection to it.

#### Updation Considerations

All fields of VR are mutable and can be changed.

### Other implementation considerations

Nil.

### Functioning at OpenStack

*Note: This section describes how Virtual Routers are implemented internally using OpenStack. This section serve as an extra documentation to explain what is happening behind at the OpenStack. Although as a KupenStack user who is working with custom resources, this knowledge may not be required. Feel free to skip this section.*

* When a VR is created in KupenStack then a Router resource is created in OpenStack for it. VR stores the reference of the Router ID from OpenStack. 
* VR creates a Router with Admin State True.
