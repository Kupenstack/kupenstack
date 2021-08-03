# OpenstackCloudConfigurationProfile

* [Summary](#Summary)
* [Motivation](#Motivation)
* [Design Details](#Design-Details)
  * [API](#API)
  * [Overview](#Overview)

### Summary

This document covers design specification, functionality, details for **OpenStackCloudConfigurationProfile** custom resource(CR) in KupenStack. A OpenStackCloudConfigurationProfile(occp) defines a profile with configurations for various OpenStack components to be deployed by KupenStack controllers.

### Motivation

KupenStack aims to automate OpenStack deployments with containers and operators. These automated deployments should have simplified OpenStack configuration management. A common, reusable OpenStack cloud configuration file is needed in KupenStack that can be shared as profiles.

Defining all OpenStack configurations in a single artifacts will further make it easy to validate and debug configuration errors.

### Design Details

#### API

```yaml
apiVersion: cluster.kupenstack.io/v1alpha1
kind: OpenStackCloudConfigurationProfile
# shortName=occp

metadata:
  # scope=Namespaced
  name: sample-profile
  namespace: default

spec:

  # The parent profile to inherit and override in this definition.
  # required=false, type=string
  from: "prod-profile.mynamespace"
  
  # Keystone related confs
  # required=false, type=object
  keystone:
  
    # Configures number of replicas for each pods.
    # requried=false, type=object
    replicas:
      
      # Number of keystone-api pods.
      # requried=false, type=integer, default=1
      api: 1
    
    # Reference: Values.conf in openstack-helm keystone chart.
    # required=false, type=object
    conf: {}
   
   
  # Glance related confs
  # required=false, type=object
  glance:
    
    # Whether to disable this component
    # required=false, type=boolean, default=false
    disable: false
  
    # Configures number of replicas for each pods.
    # requried=false, type=object
    replicas:
      
      # Number of glance-api pods.
      # requried=false, type=integer, default=1
      api: 1
      
      # Number of glance-registry pods.
      # requried=false, type=integer, default=1
      registry: 1
    
    # Reference: Values.conf in openstack-helm glance chart.
    # required=false, type=object
    conf: {}
    
    
  # Horizon related confs
  # required=false, type=object
  horizon:
    
    # Whether to disable this component
    # required=false, type=boolean, default=false
    disable: false
  
    # Configures number of replicas for each pods.
    # requried=false, type=object
    replicas:
      
      # Number of horizon-server pods.
      # requried=false, type=integer, default=1
      server: 1
    
    # Reference: Values.conf in openstack-helm horizon chart.
    # required=false, type=object
    conf: {}
  
    
  # Nova related confs
  # required=false, type=object
  nova:
    
    # Whether to disable this component
    # required=false, type=boolean, default=false
    disable: false
  
    # Configures number of replicas for each pods.
    # requried=false, type=object
    replicas:
      
      # Number of Nova api-metadata pods.
      # requried=false, type=integer, default=1
      metadata: 1
      
      # Number of Nova ironic pods.
      # requried=false, type=integer, default=1
      ironic: 1
      
      # Number of Nova placement pods.
      # requried=false, type=integer, default=1
      placement: 1
      
      # Number of nova-api-osapi pods.
      # requried=false, type=integer, default=1
      osapi: 1
      
      # Number of Nova conductor pods.
      # requried=false, type=integer, default=1
      conductor: 1
    
    # Reference: Values.conf in openstack-helm nova chart.
    # required=false, type=object
    conf: {}
  
  
  # Neutron related confs
  # required=false, type=object
  neutron:
    
    # Whether to disable this component
    # required=false, type=boolean, default=false
    disable: false
  
    # Configures number of replicas for each pods.
    # requried=false, type=object
    replicas:
      
      # Number of neutron-server pods.
      # requried=false, type=integer, default=1
      server: 1
      
      # Number of neutron-ironic-agent pods.
      # requried=false, type=integer, default=1
      ironicAgent: 1
    
    # Reference: Values.conf in openstack-helm neutron chart.
    # required=false, type=object
    conf: {}
  
  
  # Placement related confs
  # required=false, type=object
  placement:
    
    # Whether to disable this component
    # required=false, type=boolean, default=false
    disable: false
  
    # Configures number of replicas for each pods.
    # requried=false, type=object
    replicas:
      
      # Number of placement-api pods.
      # requried=false, type=integer, default=1
      api: 1
    
    # Reference: Values.conf in openstack-helm placement chart.
    # required=false, type=object
    conf: {}  

```

**Output on `kubectl get openstackcloudconfigurationprofiles` or `kubectl get occp`**

```
NAME             AGE
sample-profile   47m
```

#### Overview

KupenStack deploys containerized OpenStack using OpenStack-Helm project and assumes all conf values defaulting to charts in OpenStack-Helm. OpenStack Cloud Configuration Profiles declare what values to override on default OpenStack-Helm charts values. If a field is not provided in the OCCP CR then it means the deployment should use default values for that field from OpenStack-Helm charts.

Currently the OCCP doc covers:

* Keystone
* Glance
* Horizon
* Nova
* Neutron
* Placement

An OCCP profile can reuse any existing profile deployed in the cluster or on the internet with valid url. For example:

**Case 1**

`sample-profile` inherits from `prod-profile` in same namespace.

```yaml
apiVersion: cluster.kupenstack.io/v1alpha1
kind: OpenStackCloudConfigurationProfile
metadata:
  name: sample-profile
spec:
  from: prod-profile
```

**Case 2**

`sample-profile` inherits from `prod-profile` in `default` namespace.

```yaml
apiVersion: cluster.kupenstack.io/v1alpha1
kind: OpenStackCloudConfigurationProfile
metadata:
  name: sample-profile
spec:
  from: prod-profile.default
```

**Case 3**

`sample-profile` inherits from `prod-profile` from github url.

```yaml
apiVersion: cluster.kupenstack.io/v1alpha1
kind: OpenStackCloudConfigurationProfile
metadata:
  name: sample-profile
spec:
  from: https://raw.githubusercontent.com/kupenstack/example/prod-profile.yaml
```

The new OCCP profile overrides the values of the parent profile.

