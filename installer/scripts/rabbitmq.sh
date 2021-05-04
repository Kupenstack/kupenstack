#!/bin/bash

make -C ${HELM_CHART_ROOT_PATH} rabbitmq

helm upgrade --install rabbitmq ${HELM_CHART_ROOT_PATH}/rabbitmq \
    --namespace=openstack \
    --set volume.enabled=false \
    --set pod.replicas.server=1


# wait for rabbitmq
cd /opt/openstack-helm && ./tools/deployment/common/wait-for-pods.sh openstack
