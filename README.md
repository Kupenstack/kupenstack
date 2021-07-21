# kupenStack
KupenStack provides easy to use Cloud-Native OpenStack.

---

KupenStack is a Kubernetes add-on that provides OpenStack to Kubernetes users. KupenStack lets users declaratively provision an OpenStack cluster and resources(like VM, Subnet, Router, etc) from this OpenStack cluster using CRs(Custom Resources).

#### Design principles

* KupenStack does not modify OpenStack projects in anyways. Instead, it forms a layer of cloud-native operations on OpenStack resources and OOK deployments. Abstracts everything in the form of easy to use CRs.
* KupenStack manages only one OpenStack deployment. KupenStack when deployed to a Kubernetes cluster, it allows us to configure, scale, upgrade individual OpenStack components(like nova, neutron, glance, etc) using CRs on that Kubernetes cluster.

These design principles let KupenStack be used with:

* Any Kubernetes cluster
* Any Kubernetes tools(like Helm, Kustomize, etc).
* Any Kubernetes top-level orchestration/automation/multi-cluster tools like Crossplane, Airship, etc.

#### Project Goals/Focus

* For admins, life-saving OpenStack management. Admins should not have to ssh to every node to fix/troubleshoot OpenStack deployments.
* For developers, ability to use utilities/stacks like Helm, Prometheus, Grafana, GitOps, etc with OpenStack VMs and other resources.
* Cloud-Native operations and abstractions on OpenStack components, resources, concepts.

------



#### Demo

Self-try [demo with instructions](config/demo/readme.md)



https://user-images.githubusercontent.com/28928589/121054295-d7f63680-c7d9-11eb-9c25-f80ffa4cad4d.mp4

------



#### Extra Resources:

* [Research Paper: KupenStack: Kubernetes based Cloud-Native OpenStack](https://arxiv.org/pdf/2106.02956.pdf)
* [Slides: LFN Virtual Developer & Testing Forums](https://wiki.lfnetworking.org/display/LN/2021-06-08+-+Anuket%3A+Cloud-Native+Openstack)



#### Community:

* [Slack invite](https://join.slack.com/t/kupenstack/shared_invite/zt-rpkca4zk-HKF1ewJifKcEvHlrdMBVrQ)
* **Meetings:** Every month 1st Tuesday (16:00 UTC)
  * Next meeting: 3rd August 2021
  * Meeting link and minutes: https://docs.google.com/document/d/1jTwZkWtA6fevh3oDSuTrXKg6Ty56yCTwAWGwk5vlSgk/edit?usp=sharing

------



#### FAQ

**Q. Does KupenStack also installs the Kubernetes cluster?**

No, Kubernetes cluster is a prerequisite, KupenStack does not have any responsibility for provisioning a Kubernetes cluster. Instead, KupenStack takes benefits of this k8s cluster to simplify its OpenStack cluster management operations. So, now OpenStack users do not have to ssh to every node to fix/troubleshoot OpenStack deployments.

**Q. How to use KupenStack. Does it provide any -ctl ?**

KupenStack leverages k8s CRs. `kubectl` tool is being used to manage KupenStack. This also means other k8s utilities like `helm`,  `kustomize`, etc. can be used on KupenStack as is.

**Q. What about multi-cluster operations? Does KupenStack deploys and manages multiple OpenStack clusters?**

No, single KupenStack deployment focuses on managing a single OpenStack cluster. It keeps the design simple by having 1 k8s cluster == 1 OpenStack cluster. For multiple OpenStack clusters, KupenStack can be utilised by other projects like Airship, KubeFed, Crossplane, etc that focuses on solving this problem statement.

**Q. How does KupenStack builds images for OpenStack components?**

KupenStack project does not aim at maintaining OpenStack container images. Instead, it leverages the OpenStack-Helm project to keep deployments image agnostic. Any open-source or commercial images compatible with OpenStack-Helm charts will work with KupenStack as is.

**Q. Can we use KupenStack for a hybrid cloud of OpenStack?**

Kubernetes provides a hybrid cloud by adding on-prem nodes and public cloud nodes to the same cluster as long as those nodes are reachable to each other. KupenStack leverage Kubernetes nodes to deploy OpenStack cluster.

