*A successful open-source project is about better coordination in a team of hundreds of developers working on the same codebase.*

# Contributing Guides

Welcome to the KupenStack project and thank you for considering contributing to KupenStack. We recognise every interaction with the project as a contribution. You can help us and be part of the community by creating issues, improving documentation, fixing bugs, adding new features or interacting with the community.

If you are interested in contributing code then please start by reading this documentation. If you have questions then please reach us out on slack.

## Contributing Code

* [Communicate your intent](#Communicate-your-intent).
* [Development and making changes](#Developer-guide).

### Communicate your intent

We use GitHub for project management as most developers are familiar and comfortable with it. As a new contributor, Github issues are a good place to look at what are open issues to work on. Choose an issue that is not assigned to anyone, understand it and if you are interested then comment on it indicating that you want it to be assigned to you. This ensures better collaboration and lets others know that you are already working on it. If you want to add an enhancement or new feature then you can open an issue regarding it. Explain it thoroughly. We will be very happy to have you. You can reuse existing issue templates to open such an issue.  If you are new then good-first-issue is a good place to start. Next is [the development guide](#Developer-guide).

### Developer guide

As a  developer, a good knowledge on the following will be very helpful:

* Programming in Golang.
* Basics of Kubernetes.
* Basics of OpenStack.
* Development experience with "extending Kubernetes".

Although these won't be an extreme blocker on every topic, and we will also try to help you onboard the project wherever possible.

KupenStack project provides a very easy to use makefile to work with the project. Here are few commands:

* `make generate` generates the code used for new crds in KupenStack.
* `make manifest` generates ClusterRole, CustomResourceDefinition and WebhookConfiguration manifest files.
* `make install` applies the above-generated files into the Kubernetes cluster.
* `make run` locally runs the KupenStack on the developer machine. Make sure you have golang 1.16 or more.

Run `make help` for more commands.

