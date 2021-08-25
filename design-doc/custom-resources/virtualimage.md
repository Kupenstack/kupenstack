# VirtualImage

* [Summary](#Summary)
* [Motivation](#Motivation)
* [Design Details](#Design-Details)
  * [API](#API)
  * [Overview](#Overview)
  * [Creation Considerations](#Creation-Considerations)
  * [Updation Considerations](#Updation-Considerations)
  * [Deletion Considerations](#Deletion-Considerations)
* [Other implementation considerations](#Other-implementation-considerations)
  * [Updation](#Updation)
  * [Deletion](#Deletion)
* [Functioning at OpenStack](#Functioning-at-OpenStack)

### Summary

This document covers design specification, functionality details of **VirtualImage** custom resource(CR) in KupenStack. A Virtual Image CR manages glance image in the associated OpenStack cluster.

### Motivation

Images resources are required for creating Virtual Machine in OpenStack. Before creating Virtual Machines users need to create and upload image in glance service of OpenStack. Therefore it is desired to have custom resources in KupenStack to manage these images.

### Design Details

#### API

```yaml
apiVersion: kupenstack.io/v1alpha1
kind: VirtualImage

metadata:
  # scope=Namespaced
  name: image-sample
  namespace: default

spec:

  # Source contains url to pull image from.
  # required=true, type=string, mutable=flase
  src: https://some-url.com
  
  # Disk format of the image.
  # required=true, type=enum(raw;qcow2;iso;vdi;ami;ari;aki;vhd;vmdk), mutable=flase
  format: raw
  
  # ContainerFormat is the format of the container.
  # required=false, type=enum(bare;ami;ari;aki;ovf); default=bare, mutable=flase
  containerFormat: bare
  
  # Minimum disk size in GB required to boot this image.
  # required=false, type=integer, mutable=flase
  minDisk: 0
  
  # Minimum ram size in MB required to boot this image.
  # required=false, type=integer, mutable=flase
  minRam: 0

status:
  
  # Id of image at openstack cloud
  # type=string
  id: a29d8-1d73n-2dw45-h4hr2
  
  # Whether image is ready for use or not
  # type=boolean
  ready: true
  
  # Size of image stored in openstack cloud.
  # type=string
  size: "56.78mb"
  
  usage:
  
    # Contains list of all the virtual machines using this image.
    # type=array
    virtualMachines: [ "Namespace/Name" ]
  
    # InUse is true if the image is being used by one or more virtual machine.
    # type=boolean
    inUse: false
```

**Output on `kubectl get virtualimages`**

```
NAME           READY   IN-USE   SIZE      AGE
image-sample   true    true     56.78mb   47m
```

#### Overview

Virtual Image CR have a source url from where it pulls image and stores locally. Virtual Machines CR uses these virtual images to create virtual machines.

###### src, format, containerFormat, minDisk, minRam

`src` is any valid url that virtual machine image is available at. `format` describes the format of the image disk. `containerFormat` specifies the format of the image container. `minDisk` specifies the amount of disk space in GB that is required to boot the image. `minRam` specifies the amount of RAM in MB that is required to boot the image.

#### Creation Considerations

`src` and `format` are required fields therefore they need to be known prior to image creation. `src` and `format` are immutable fields therefore cannot be changed after creation.

#### Updation Considerations

All fields in Virtual Image CR are immutable, therefore cannot be updated. 

#### Deletion Considerations

If a Virtual Image is deleted while it is in use by any Virtual Machines then its deletion will get in pending state. Once all Virtual Machine using it will be deleted then Virtual Image will delete itself. Note: As long as Virtual Image is in pending state for deletion and have not been deleted, new Virtual Machines can be created from it.

### Other implementation considerations

#### Updation

The list of virtual machines in the status is not updated by Virtual Image controller. Whenever a Virtual Machine is created that uses this Virtual Image then the Virtual Machine controller updates status of Virtual Image by adding virtual machine's `Namespace/Name` to the usage list. Similarly it removes that `Namespace/Name` when Virtual Machine is deleted.

#### Deletion

During deletion of Virtual Image, the controller first waits for all Virtual Machines using it to get deleted. Once no Virtual Machine is using this virtual image then first the `ready` field in the status of the virtual image is set to false then actual deletion of virtual image starts. This ensures no new Virtual Machines can be created while the image is being deleted.

### Functioning at OpenStack

*Note: This section describes how Virtual Images are implemented internally using OpenStack. This section serve as an extra documentation to explain what is happening behind at the OpenStack. Although as a KupenStack user who is working with custom resources, this knowledge may not be required. Feel free to skip this section.*

* OpenStack image is created with visibility as public.
* OpenStack image is created with protected as false.

