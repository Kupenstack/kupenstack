# Volume

* [Summary](#Summary)
* [Motivation](#Motivation)
* [Design Details](#Design-Details)
  * [API](#API)
  * [Overview](#Overview)
  * [Updation Considerations](#Updation-Considerations)
  * [Deletion Considerations](#Deletion-Considerations)
* [Other implementation considerations](#Other-implementation-considerations)
* [Functioning at OpenStack](#Functioning-at-OpenStack)

### Summary

This document covers design specification, functionality details of **Volume** custom resource(CR) in KupenStack. A Volume CR manages volumes in the associated OpenStack cluster.

### Motivation

Volumes are required to store and move data for workloads in a cloud-agnostic way as these volumes can come from various providers. Virtual Machines can be booted using bootable volumes or store application data in them.

### Design Details

#### API

```yaml
apiVersion: kupenstack.io/v1alpha1
kind: Volume

metadata:
  # scope=Namespaced
  name: volume-sample
  namespace: default

spec:

  # Size of the volume
  # required=true, type=resource.Quantity
  size: 2Gi

  # Data source to initialize volume with.
  # required=false, type=object
  source:
    
    # Intialize volume with virutal machine image data from given Url.
    # required=false, type=string, mutable=flase
    image: "https://some-url.com"

status:
  
  # Id of volume at openstack cloud
  # type=string
  id: a29d8-1d73n-2dw45-h4hr2
  
  # Whether volume is ready for use or not
  # type=boolean
  ready: true

  # Contains list of all the virtual machines using this volume.
  # type=array
  virtualMachines: [ "Namespace/Name" ]
  
  # InUse is true if the volume is being attached to one or more virtual machines.
  # type=boolean
  inUse: false
```

**Output on `kubectl get virtualimages`**

```
NAME            READY   IN-USE   SIZE      AGE
volume-sample   true    true     2GB       47m
```

#### Overview

Volume CR manages a volume. This volume can come from various types of storage backends. Can be of various types, have different access permissions, sizes and capabilities. A Volume CR is focused on 3 types of declarations:

* Storage: Storage related details like size, type, provider, class, permissions, capabilities, etc.
* Data: Data to use for initialization of this volume
* Mounts: Mounts can include Volume mount, Configmap mount, Secret mount, Persistent Volumes, etc 

Volume CR is managed by developers in their applications and is designed to be independent of cluster implementations, making them portable.

#### Updation Considerations

Once a Image is used to initialize a Volume then the image field in the source becomes immutable and cannot be changed.

#### Deletion Considerations

Volumes cannot be deleted while attached to Virtual Machines.

### Other implementation considerations

### Functioning at OpenStack

*Note: This section describes how Volumes are implemented internally using OpenStack. This section serve as an extra documentation to explain what is happening behind at the OpenStack. Although as a KupenStack user who is working with custom resources, this knowledge may not be required. Feel free to skip this section.*

* 

