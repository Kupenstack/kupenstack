# Build the manager binary
FROM golang:1.16 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY apis/ apis/
COPY controllers/ controllers/
COPY pkg/ pkg/
COPY oskops/ oskops/
# COPY ook-operator/ ook-operator/


# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o manager main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM ubuntu:bionic-20210827

RUN apt-get update

# install pip3
# RUN apt-get -y install python3-pip

# install helm
RUN apt-get -y install wget
RUN wget  https://get.helm.sh/helm-v3.7.0-linux-amd64.tar.gz
RUN tar -zxvf helm-v3.7.0-linux-amd64.tar.gz
RUN mv linux-amd64/helm /bin/helm

# # install kubectl
# RUN apt-get -y install curl
# RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
# RUN install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
# RUN rm kubectl

# # clone openstack-helm
# RUN apt -y install git
# WORKDIR /opt
# RUN git clone https://github.com/openstack/openstack-helm
# RUN git clone https://github.com/openstack/openstack-helm-infra
# WORKDIR openstack-helm

# # Setup Clients
# RUN pip3 install \
#   -c${UPPER_CONSTRAINTS_FILE:=https://releases.openstack.org/constraints/upper/${OPENSTACK_RELEASE:-stein}} \
#   cmd2 python-openstackclient python-heatclient

# RUN pip3 install yq
# RUN apt -y install jq


# ENV OSH_PATH=/opt/openstack-helm
# ENV OSH_INFRA_PATH=/opt/openstack-helm-infra



# WORKDIR /workspace
# COPY ook-operator/settings/ ook-operator/settings/
# COPY ook-operator/pkg/actions/ ook-operator/pkg/actions/

# RUN chmod +x ook-operator/pkg/actions/cluster/apply
# RUN chmod +x ook-operator/pkg/actions/glance/apply
# RUN chmod +x ook-operator/pkg/actions/helm/initCreds
# RUN chmod +x ook-operator/pkg/actions/helm/initHelm
# RUN chmod +x ook-operator/pkg/actions/horizon/apply
# RUN chmod +x ook-operator/pkg/actions/ingress/apply
# RUN chmod +x ook-operator/pkg/actions/keystone/apply
# RUN chmod +x ook-operator/pkg/actions/libvirt/apply
# RUN chmod +x ook-operator/pkg/actions/mariadb/apply
# RUN chmod +x ook-operator/pkg/actions/memcached/apply
# RUN chmod +x ook-operator/pkg/actions/neutron/apply
# RUN chmod +x ook-operator/pkg/actions/nova/apply
# RUN chmod +x ook-operator/pkg/actions/placement/apply
# RUN chmod +x ook-operator/pkg/actions/rabbitmq/apply


WORKDIR /
# COPY ook-operator/settings/ /etc/kupenstack/
COPY --from=builder /workspace/manager .

ENTRYPOINT ["/manager"]


