# OpenStack Cluster Life Cycle Management

* [Summary](#Summary)
* [Approach](#Approach)
* [KupenStack config file](#KupenStack-config-file)

### Summary

KupenStack manages a self-provisioned OpenStack cluster. The doc describes the approach for OpenStack cluster LCM.

## Approach

KupenStack deploys and manages all OpenStack components as containers through Kubernetes APIs. This makes OpenStack highly customizable and pluggable. For deploying these containers KupenStack relies on OpenStack-Helm projects. KupenStack takes [OpenStack-Helm Project](https://github.com/openstack/openstack-helm) charts as a standard interface for deploying containerized OpenStack cloud. While OpenStack-Helm Projects targets building charts for OpenStack, KupenStack focuses on taking these charts as a baseline and providing automated operator operations for lcm of OpenStack cluster. So, that users do not have to understand and manually manage these container deployments. The user should describe the desired state of the cluster and KupenStack takes it as its responsibility to bring up the OpenStack cluster to the desired state.

​            The desired configuration of the OpenStack cluster comes from the OpenStackNode custom resource(cr). Each OpenStack Node (osknode) has `status.desiredClusterConfiguration` which are generated from OpenStackCloudConfigurationProfile(OCCP) and are values passed to OpenStack-Helm charts during the management of OpenStack-Helm deployments by KupenStack. 

​            KupenStack heavily inspires by the design and architecture of Kubernetes. For every OpenStack component(Nova, Neutron, etc..), KupenStack runs a separate reconciliation loop that manages all containers deployed by it. Hence, the nova reconciliation loop only manages nova pods and so on.

​             Since KupenStack has principles of not modifying OpenStack and takes OpenStack-Helm as a standard for describing OpenStack deployments, therefore, KupenStack can deploy any OpenStack container images if they are compatible with OpenStack-Helm Project.

## KupenStack config file

### API

```yaml
apiVersion: kupenstack.io/v1alpha1
kind: KupenstackConfiguration
metadata:
  name: configfile
spec:
  # Name of default profile to apply on each node.
  # required=true, type=object
  defaultProfile:
    name: profile-sample
    namespace: default
  
  # List of osk nodes from k8s cluster.
  nodes:
    - name: node12
      disabled: true
    - name: kind-control-plane
      type: control,compute
      disable: false
```

