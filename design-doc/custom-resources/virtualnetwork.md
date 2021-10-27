# VirtualNetwork

* [Summary](#Summary)
* [Motivation](#Motivation)
* [Design Details](#Design-Details)
  * [API](#API)
  * [Overview](#Overview)
    * [Cidr](#Cidr)
    * [Allocation Pools](#Allocation-Pools)
    * [Gateway IP](#Gateway-IP)
    * [Disable Dhcp](#Disable-Dhcp)
    * [MTU](#MTU)
  * [Creation Considerations](#Creation-Considerations)
  * [Deletion Considerations](#Deletion-Considerations)
  * [Updation Considerations](#Updation-Considerations)
* [Other implementation considerations](#Other-implementation-considerations)
  * [Deletion](#Deletion)
  * [Updation](#Updation)
* [Functioning at OpenStack](#Functioning-at-OpenStack)

### Summary

This document covers design specification, functionality, details for **VirtualNetwork** custom resource(CR) in KupenStack. A Virtual Network CR manages the network, subnet resources and related networking in the associated OpenStack cluster.

### Motivation

Developers want to architect and manage network topology for their workloads in a flexible way. Declarative virtual networks are a desired feature in a Kubernetes-native environment.  OpenStack has multiple networking concepts like Networks, Subnets, DHCP, IP pools, DNS, Host routes, etc. for allowing users to create complete virtual network infrastructure. The motivation is to support the initial idea of allowing the user to build and architect complete virtual network infrastructure for their applications in a KupenStack cluster while simplifying the concepts and usage as compared to the previous implementation of these in OpenStack.

The ability to easily build and manage a virtual network gives various benefits and allows developers to isolate and secure their workloads in separate networks.

### Design Details

#### API

```yaml
apiVersion: kupenstack.io/v1alpha1
kind: VirtualNetwork
# shortName=vn

metadata:
  # scope=Namespaced
  name: sample-network
  namespace: default

spec:

  # CIDR block to use for this private network in IPv4 or IPv6.
  # required=false, type=string, mutable=flase
  cidr: "10.10.10.0/24"
  
  # AllocationPools are a way of telling DHCP agent what IP to use from the cidr block.
  # required=flase, type=array
  allocationPools:
    - startIP: "10.10.10.100"
      endIP: "10.10.10.150"
  
  # Statically defined the IP that gateway(router) should be given in this network.
  # required=false, type=string
  gatewayIP: "10.10.10.5"
  
  # Does not creates DHCP agent in this network when true.
  # required=false, type=boolean, default=false
  disableDhcp: false
  
  # The maximum transmission unit(MTU) value to address fragementation.
  # Minimum value is 64 for IPv4 and 1280 for IPv6.
  # required=false, type=integer
  mtu: 100

status:
  
  # Id of network at openstack cloud
  # type=string
  id: a29d8-1d73n-2dw45-h4hr2
  
  # Whether virtual network is ready for use or not
  # type=boolean
  ready: true
  
  # Name of the virtualrouter that this virtualnetwork is connected to.
  # type=string
  gatewayName: "none"
  
  # String representing ratio of used IPs to total no. of IPs in available virtualnetwork.
  # type=string
  ipUsed: "2/50"

```

**Output on `kubectl get virtualnetworks` or `kubectl get vn`**

```
NAME             IP-USAGE   GATEWAY   AGE
sample-network   2/50       none      47m
```

**Output on `kubectl get vn -o wide`**

```
NAME             CIDR            IP-USAGE   GATEWAY   gatewayIP    AGE
sample-network   10.10.10.0/24   2/50       none      10.10.10.5   47m
```

#### Overview

Every Virtual Network CR has a CIDR block to define the IP block to use for it. Any VirtualMachine that attaches to this Virtual Network(VN) is automatically assigned an IP from this VN. 

###### Cidr

CIDR can be either IPv4 or IPv6 and are immutable. Once the network is created with a CIDR block then it cannot be changed for it. All other fields in VN follows same IPv4 or IPv6 conventions as used by Cidr. It is a string data type.

CIDR is an optional field. If CIDR is not provided on creation then any `10.*.*.0/24` block is automatically assigned if it does not overlap with any existing VN, where `*` can be any value between 0-254.  If CIDR is provided on creation then the VN is created with that CIDR even if it overlaps with any of the existing CIDR.

Also, note that VN are namespace scoped resources therefore when CIDR is not provided then during automatic allocation of CIDR if no non-overlapping CIDR are available globally then non-overlapping CIDR is chose locally in the same namespace.

Also, from above please note that by default IPv4 CIDR is used if not provided.

###### Allocation Pools

Allocation pools can control what IP to use from a CIDR block. Allocation pools are optional hence, If no allocation pools are configured then all IP in CIDR are available for use. Allocation pools can be configured with startIP, endIP or both.

* When StartIP is provided then allocatable IP start with startIP till the end of the CIDR block.
* When EndIP is provided then allocatable IP starts with the first IP of CIDR till the endIP in block.
* When StartIP, EndIP both are provided then a specific pool starting with StartIP in CIDR till endIP in CIDR is used.

Allocation pools take a list of such StartIP, endIP pairs while configuring. Multiple pairs of StartIP, EndIP can be defined in allocation pools to defined multiple pool ranges to use for IP allocation.

StartIP, EndIP can define exact IP from CIDR like `10.10.10.100` or can use wildcard notation for the network part in the IP address like `*.*.*.100`. Note: `10.*.*.100` is also valid, but `10.10.10.*` is not valid as the startIP/endIP are not clear in this case.

###### Gateway IP

Gateway IP is an optional field. When a VN connects to a virtual router then it needs to assign an IP to it. By default, VN automatically assigns the first available IP from CIDR to the virtual router, this behaviour can be overridden by statically providing what gateway IP to use for the virtual router with this field.

###### Disable DHCP

This is an optional field. By default, DHCP is enabled for a VN. DHCP can be disabled for this network by setting this field to true.

###### MTU

The maximum transmission unit(MTU) value to use for this VN.

#### Creation Considerations

VN is not dependent on any field from other CRs in the cluster. As CIDR is an optional field, VN can be created with minimal configuration as follow:

```yaml
# Minimal Virtual Network example
apiVersion: kupenstack.io/v1alpha1
kind: VirtualNetwork
metadata:
  name: sample-network
---
```

Although, if VN is created with a specific CIDR block defined in its definition then considerations have to be taken regarding overlapping CIDR. As two VN can be created with the same CIDR but cannot to be attached to the same router if they have so.

Also, it is to be noted that when you define allocation pools for your VN without defining CIDR block then consider using wildcard notations to prevent errors.

```yaml
# Allocation pools Virtual Network example
apiVersion: kupenstack.io/v1alpha1
kind: VirtualNetwork
metadata:
  name: sample-network
spec:
  allocationPools:
    - StartIP: "*.*.*.100"
---
```

#### Deletion Considerations

VN does not have any pre-clean-up requirements. Deleting VN deletes it. But it is to be noted that if a VN is deleted while any virtual machine is attached to it or router then they will loose connectivity to the VN as it gets deleted.

#### Updation Considerations

All fields of VN are mutable except CIDR. When a VN is updated then it will try to reconcile itself to the new state.

### Other implementation considerations

#### Deletion

Before deleting VN, all resources in it have to be deleted. These are OpenStack ports, subnet, network resources internally. VN controller first cleans all ports, then the subnet, and then the network. Then deletes the VN.

Meanwhile, VM controller in every reconcile loop checks the desired network connections of vm to actual network connection, if they do not match then it creates required ports and connect to those networks.

During deletion of VN, first, the ready state of the VN has to be set to be false then the VN controller should start deleting all ports in it. Setting status ready as false indicates VM that the virtual network is not available to make changes.

#### Updation

Note VN controllers only reconcile for VNs desired configuration. VN controller is not responsible for ensuring connection to VM or virtual router. Those are handled individually by VM and virtual router controllers respectively.

### Functioning at OpenStack

*Note: This section describes how Virtual Networks are implemented internally using OpenStack. This section serve as an extra documentation to explain what is happening behind at the OpenStack. Although as a KupenStack user who is working with custom resources, this knowledge may not be required. Feel free to skip this section.*

* When a VN is created in KupenStack then a Network and Subnet are created in OpenStack for it. VN stores the reference of the Network ID at OpenStack. 
* The Subnet created uses same name as the VN.
* VN creates Network with Admin State True.
* VN creates Network with shared as True.
* VN creates Subnet with Gateway enabled.
* VN does not configure any DNS Name Server or Host Routes.

