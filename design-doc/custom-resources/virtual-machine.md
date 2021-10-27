# VirtualMachine

* [Summary](#Summary)
* [Motivation](#Motivation)
* [Design Details](#Design-Details)
  * [API](#API)
  * [Overview](#Overview)
    * [Running](#running)
    * [Vcpu](#vcpu)
    * [Memory](#memory)
    * [Swap](#swap)
    * [Public Key](#publickey)
    * [Virtual Networks](#virtualNetworks)
    * [Disk Attachments](#diskAttachments)
  * [Deletion Considerations](#Deletion-Considerations)
* [Other implementation considerations](#Other-implementation-considerations)
* [Functioning at OpenStack](#Functioning-at-OpenStack)

### Summary

This document covers design specification, functionality, details for **VirtualMachine** custom resource(CR) in KupenStack. A Virtual Machine CR manages virtual machine resources and related configurations in the associated OpenStack cluster.

### Motivation

Virtual Machines are the primitive workload type in OpenStack. They have various use-cases, and it is required to have higher-level operations automated on them to ensure features like high-availability, self-healing, scalability, etc.

### Design Details

#### API

```yaml
apiVersion: kupenstack.io/v1alpha1
kind: VirtualMachine
# shortName=vm

metadata:
  # scope=Namespaced
  name: sample-vm
  namespace: default

spec:

  # Desired state of vm is set to running when true and shutdown when false.
  # required=false, type=boolean, default=true
  running: true
  
  # Vcpu count for the machine.
  # required=true, type=integer 
  vcpu: 3
  
  # Memory size in megabytes for the machine.
  # required=true, type=resource.Quantity
  memory: 500Mi
  
  # Swap storage size in gigabytes for the machine.
  # required=false, type=resource.Quantity, default=0
  swap: 1Gi
  
  # Public ssh key to add to the server.
  # required=false, type=string
  publicKey: "ssh-rsa AAAAB3Nza....oJaUMfHvJs= parth@kupenstack.io"

  # List of virtual networks to connect this virtual machine to.
  # required=false, type=array
  virtualNetworks:
    - "net-db"
    - "net-backend"
    - "net-front-link"
  
  # List of volumes to attach to the machine as disk.
  # required=true, type=array
  diskAttachments:
    - name: "myvol"
    - name: "vol-db"
    - name: "vol-users-creds"

status:
  
  # Id of virtual machine at openstack cloud.
  # type=string
  id: a29d8-1d73n-2dw45-h4hr2
  
  # Status of virtual network.
  # type=string
  phase: "Running"
  
  # Hostname on which machine is running.
  # type=string
  node: "node1"
  
  # Contains the status on comma separated list of network(IP) that virtual machine has.
  # type=string
  ip: "net-db(10.10.120.21), net-backend(10.10.150.12), net-front-link(10.10.90.100)"
  
  # number of times the VM is restarted.
  # type=integer
  restartCount: 0
```

**Output on `kubectl get virtualmachines` or `kubectl get vm`**

```
NAME        STATUS    RESTARTS   AGE
sample-vm   Running   0          47m
```

**Output on `kubectl get vm -o wide`**

```
NAME        STATUS    RESTARTS   AGE   NODE    NETWORKS(IP)
sample-vm   Running   0          47m   node1   net-db(10.10.120.21), net-backend(10.10.150.12), net-front-link(10.10.90.100)
```

#### Overview

A Virtual Machine CR creates and manages a virtual machine(VM). VMs are treated as ephemeral workloads. A VM can be attached/detached to multiple virtual networks and volumes during its lifetime. Volumes attached to a VM are treated as disk devices unlike `volumeMounts` in Pods.

###### running

Set state of the virtual machine to either running or shutdown(when running is false).

###### vcpu

This field specifies the number of virtual cpu cores to allocate for the virtual machine.

###### memory

This field specifies the memory size required to allocate for the virtual machine.

###### swap

This field specifies the swap memory size required to allocate for the virtual machine. It defaults to 0 when not given.

###### publicKey

This field specifies the public ssh key to invoke into the virtual machine. 

###### virtualNetworks

`virtualNetworks` takes a list of names of all the virtual networks to connect this virtual machine to. These virtual networks should exists in the same namespace in which virtual machine lives. 

###### diskAttachments

This field specifies the list of volumes to attach to the virtual machine as disk. The volumes must exist in the same namespace as that of virtual machine. When a VM is booted it starts with the first available bootable disk.

#### Deletion Considerations

Deletion of a VM will not delete the volumes or virtual networks linked to it.

### Other implementation considerations

### Functioning at OpenStack

*Note: This section describes how Virtual Machines are implemented internally using OpenStack. This section serve as an extra documentation to explain what is happening behind at the OpenStack. Although as a KupenStack user who is working with custom resources, this knowledge may not be required. Feel free to skip this section.*

* For a VM, flavors, keypairs are automatically managed in the Openstack cluster.
* A VM is not created directly using the image and root disk, instead booted from attached volumes.
* 

