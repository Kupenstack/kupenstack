# OpenstackNodes

* [Summary](#Summary)
* [Motivation](#Motivation)
* [Design Details](#Design-Details)
  * [API](#API)
  * [Overview](#Overview)

### Summary

This document covers design specification, functionality, details for **OpenstackNode** custom resource(CR) in KupenStack. An OpenStack Node CR manages the details and status of the desired OpenStack cluster self-provisioned by KupenStack.

### Motivation

KupenStack aims to automate OpenStack deployments with containers and operators. These automated deployments have node level lcm of OpenStack cluster. Therefore, it is required to track the state of the cluster for each node individually. Having OpenStack Nodes similar to Kubernetes Nodes can simplify debugging and monitoring the OpenStack cluster.

An OpenStack node can store the desired OpenStack configurations for that node and the show status of all components of OpenStack for that node.

### Design Details

#### API

```yaml
apiVersion: kupenstack.io/v1alpha1
kind: OpenstackNode
# shortName=osknode

metadata:
  # scope=Cluster
  name: node1

spec:
  
  # OpenStack Cloud Configuration Profile(OCCP) used by this node.
  # required=true, type=object
  openstackCloudConfigurationProfileRef:
    name: ""
    namespace: ""
  
status:
  
  # Generated configration from OCCP.
  # type=object
  desiredNodeConfiguration: {}
  
  # Status of OpenStack cluster components for this osknode.
  # type=string
  status: Ready 
```

**Output on `kubectl get openstacknodes` or `kubectl get osknodes`**

```
NAME             STATUS   ROLES           PROFILE          AGE
sample-network   Ready    control-plane   sample-profile   47m
```

#### Overview

OpenStack Nodes are automatically created by KupenStack. For every Kubernetes node, we have an OpenStack Node with the same name. The purpose of OpenStack Nodes is to keep track of OpenStack components and their configuration for that node. OpenStack Nodes drives the desired OpenStack configurations from the occp profile used by them.
