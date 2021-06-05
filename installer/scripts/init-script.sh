#!/bin/bash


# Add node labels

# label openstack-control-plane=enabled
controlnodes=$(yq -r .spec.controlNodes /etc/kupenstack/config/kupenstack.yaml)
for node in "${controlnodes[@]}"; do nodename=$(echo $node | sed 's/[]"/[]//g'); kubectl label nodes $nodename openstack-control-plane=enabled; done

# label openstack-compute-node=enabled
computeNodes=$(yq -r .spec.computeNodes /etc/kupenstack/config/kupenstack.yaml)
for node in "${computeNodes[@]}"; do nodename=$(echo $node | sed 's/[]"/[]//g'); kubectl label nodes $nodename openstack-compute-node=enabled; done

# label linuxbridge=enabled
for node in "${controlnodes[@]}"; do nodename=$(echo $node | sed 's/[]"/[]//g'); kubectl label nodes $nodename linuxbridge=enabled; done
for node in "${computeNodes[@]}"; do nodename=$(echo $node | sed 's/[]"/[]//g'); kubectl label nodes $nodename linuxbridge=enabled; done


mkdir /etc/kupenstack/auth/openstack

authUrl="http://keystone.kupenstack.svc.cluster.local/v3"
username="admin"
password="password"
domain="Default"
tenant="admin"

echo $(echo $authUrl | base64 ) > /etc/kupenstack/auth/openstack/authUrl
echo $(echo $username | base64 ) > /etc/kupenstack/auth/openstack/username
echo $(echo $password | base64 ) > /etc/kupenstack/auth/openstack/password
echo $(echo $domain | base64 ) > /etc/kupenstack/auth/openstack/domain
echo $(echo $tenant | base64 ) > /etc/kupenstack/auth/openstack/tenant


# Init Helm
/etc/kupenstack/helm.sh

# Deploy Ingress
/etc/kupenstack/ingress.sh

# Deploy MariaDB
/etc/kupenstack/mariadb.sh

# Deploy Rabbitmq
/etc/kupenstack/rabbitmq.sh

# Deploy Memcached
/etc/kupenstack/memcached.sh

# Deploy Keystone
/etc/kupenstack/keystone.sh

# Deploy Horizon
/etc/kupenstack/horizon.sh

# Deploy Glance
/etc/kupenstack/glance.sh

# Deploy Compute-Kit(Libvirt, Placement, Nova, Neutron, Linux-Bridge)
/etc/kupenstack/compute-kit.sh

echo "Deployment Completed"

