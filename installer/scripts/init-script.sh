#!/bin/bash


# Add node labels

# label openstack-control-plane=enabled
controlnodes=$(yq -r .spec.controlNodes /etc/kupenstack/kupenstack.yaml)
for node in "${controlnodes[@]}"; do nodename=$(echo $node | sed 's/[]"/[]//g'); kubectl label nodes $nodename openstack-control-plane=enabled; done


# Init Helm
/etc/kupenstack/helm.sh

# Deploy Ingress
/etc/kupenstack/ingress.sh
