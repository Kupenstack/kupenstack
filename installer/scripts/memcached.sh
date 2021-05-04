#!/bin/bash

make -C ${HELM_CHART_ROOT_PATH} memcached

helm upgrade --install memcached ${HELM_CHART_ROOT_PATH}/memcached --namespace=openstack

# wait for memcached
cd /opt/openstack-helm && ./tools/deployment/common/wait-for-pods.sh openstack
