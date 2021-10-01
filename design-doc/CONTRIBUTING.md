# Design Docs

* [Contributing Guide: Custom Resource Design Documents](#Custom-Resource-Design-Document-Contributing-Guide)
  * [How to document CR APIs](#How-to-document-cr-apis)





## Custom Resource Design Document Contributing Guide

[custom-resources/](./custom-resources/) folder contains design documents for each custom resource defined in KupenStack. These are abstract design documents explaining API and the behaviour of a CR. Any design changes have to first go through changes in these documents after which they are implemented in the actual implementation of KupenStack. Each documentation has a table of contents. These docs can follow any structure as needed. A commonly recommended structure is as follow:

```markdown
### Summary
This section contains an intro about what custom resource this document focuses on.
### Motivation
This section contains why this custom resource is needed? what problems does it solves? what are its use-cases?
### Design Details
	#### API
	This section contains an abstract, formatted YAML explaining the custom resource API. As explained in the #How-to-document-cr-api section below.
	#### Overview
	This section contains a detailed explanation of each field in the API and how they can be used.
	#### Creation/Deletion/Updation Considerations
	These sections contain extra usage considerations for users while creation/updation/deletion which may not have been covered in previous sections.
### Other implementation considerations
This section is extra notes by KupenStack developers to point out few considerations that were kept in mind while implementing the custom resource. This section can be further divided into ####Creation, ####Updation ####Deletions, etc.
### Functioning at OpenStack
This section is extra documentation to give an overview of how we implement the custom resource using OpenStack types internally. These documents may not explain everything, instead are abstract for curious readers.
```

These custom resource design docs do not cover a complete explanation of custom resources instead  serves as a reference blueprint for their implementations.

 

### How to document CR APIs

We follow [api-conventions](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md) set by the k8s-sig-architecture while defining custom resource APIs.

For explaining the custom resource API in our design doc take any example YAML for the CR.

```yaml
apiVersion: kupenstack.io/v1alpha1
kind: FooBar
metadata:
  name: sample-foobar
  namespace: default
spec:
  foo: bar
```

Then for every field, we add a short explanation by adding comments on them.

```yaml
apiVersion: kupenstack.io/v1alpha1
kind: FooBar
metadata:
  name: sample-foobar
  namespace: default
spec:
  # Tells about the bar
  foo: example-bar
```

To explain datatype and other properties for these field we use key=value conventions in comment, for example:

```yaml
apiVersion: kupenstack.io/v1alpha1
kind: FooBar
metadata:
  name: sample-foobar
  namespace: default
spec:
  # Tells about the bar
  # required=false, type=string, mutable=flase
  foo: example-bar
```

The above example can be interpreted as: the name of the field is `foo` and is an optional field. The datatype of foo is a string and it is an immutable field. The description of `foo` explains that it tells about the bar. More on these in [Reference](#reference).

In API doc we also add example output of `kubectl get ` and `kubectl get -o wide `  for the custom resource.

Please see existing design-docs as a reference. Design-docs are loosely written documents working as reference blueprints for implementation. A complete explanation of custom resources will be documented separately for each type in API reference documentations. 

#### Reference

* `required`: 

  ​	values: true, false

  ​	example: `required=false`

  ​	description: Whether this field is required or not.

* `type`: 

  ​	values: integer, string, boolean, array, object, resource.Quantity, enum

  ​	example: `type=string` or `type=enum(bare;ami;ari;aki;ovf)`

  ​	description: Describes the data type for this field.

* `mutable`: 

  ​	values: true, false

  ​	example: `mutable=false`

  ​	description: Whether is field is mutable or not.

* `default`:

  ​	values: based on datatype

  ​	example: `default=0` 

  ​	description: Describes default value of this field.

* `scope`:

  ​	values: Namespaced, Cluster

  ​	example: `scope=Namespaced`

  ​	description: Describes whether the CR is cluster scoped or namespaced resource. Add this as comment under `metadata` .

* `shortName`;

  ​	values: any string

  ​	example: `shortName=pvc`

  ​	description: Describes short name for CR Kind if required. Add this as comment below kind field of CR.

