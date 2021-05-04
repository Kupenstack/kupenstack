#!/bin/bash

OPENSTACK_HELM_ROOT_PATH=/opt/openstack-helm

make -C ${OPENSTACK_HELM_ROOT_PATH} horizon

helm upgrade --install horizon ${OPENSTACK_HELM_ROOT_PATH}/horizon \
    --namespace=openstack \
    --set pod.replicas.server=1


# wait for horizon
cd /opt/openstack-helm && ./tools/deployment/common/wait-for-pods.sh openstack

