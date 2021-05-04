FROM ubuntu:bionic-20210416

RUN apt-get update

# install pip3
RUN apt-get -y install python3-pip

# install helm
RUN apt-get -y install wget
RUN wget  https://get.helm.sh/helm-v2.16.9-linux-amd64.tar.gz
RUN tar -zxvf helm-v2.16.9-linux-amd64.tar.gz
RUN mv linux-amd64/helm /bin/helm


# install kubectl
RUN apt-get -y install curl
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
RUN install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
RUN rm kubectl

# clone openstack-helm
RUN apt -y install git
WORKDIR /opt
RUN git clone https://github.com/openstack/openstack-helm
RUN git clone https://github.com/openstack/openstack-helm-infra
WORKDIR openstack-helm

# Setup Clients
RUN pip3 install \
  -c${UPPER_CONSTRAINTS_FILE:=https://releases.openstack.org/constraints/upper/${OPENSTACK_RELEASE:-stein}} \
  cmd2 python-openstackclient python-heatclient

ENV OSH_INFRA_PATH=/opt/openstack-helm-infra
ENV HELM_CHART_ROOT_PATH=/opt/openstack-helm-infra

COPY configs/clouds.yaml /etc/openstack/clouds.yaml
RUN chown -R $(id -un): /etc/openstack


# The ingress controller
COPY configs/ingress-kube-system.yaml /tmp/ingress-kube-system.yaml
COPY configs/ingress-component.yaml /tmp/ingress-component.yaml
COPY scripts/ingress.sh /etc/kupenstack/ingress.sh

# mariaDB
COPY configs/mariadb.yaml /tmp/mariadb.yaml

# Setup
RUN pip3 install yq
RUN apt -y install jq

COPY scripts/init-script.sh /etc/kupenstack/init-script.sh
COPY scripts/helm.sh /etc/kupenstack/helm.sh
COPY scripts/mariadb.sh /etc/kupenstack/mariadb.sh
COPY scripts/rabbitmq.sh /etc/kupenstack/rabbitmq.sh
COPY scripts/memcached.sh /etc/kupenstack/memcached.sh
COPY scripts/keystone.sh /etc/kupenstack/keystone.sh
COPY scripts/horizon.sh /etc/kupenstack/horizon.sh

RUN chmod +x /etc/kupenstack/init-script.sh
RUN chmod +x /etc/kupenstack/helm.sh
RUN chmod +x /etc/kupenstack/ingress.sh
RUN chmod +x /etc/kupenstack/mariadb.sh
RUN chmod +x /etc/kupenstack/rabbitmq.sh
RUN chmod +x /etc/kupenstack/memcached.sh
RUN chmod +x /etc/kupenstack/keystone.sh
RUN chmod +x /etc/kupenstack/horizon.sh

CMD ["/etc/kupenstack/init-script.sh"]
