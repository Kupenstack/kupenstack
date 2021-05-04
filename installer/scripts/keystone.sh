#!/bin/bash

OPENSTACK_HELM_ROOT_PATH=/opt/openstack-helm

make -C ${OPENSTACK_HELM_ROOT_PATH} keystone

helm upgrade --install keystone ${OPENSTACK_HELM_ROOT_PATH}/keystone \
    --namespace=openstack \
    --set pod.replicas.api=1


# wait for keystone
cd /opt/openstack-helm && ./tools/deployment/common/wait-for-pods.sh openstack

