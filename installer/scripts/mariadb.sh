#!/bin/bash

make -C ${HELM_CHART_ROOT_PATH} mariadb

helm upgrade --install mariadb ${HELM_CHART_ROOT_PATH}/mariadb \
    --namespace=openstack \
    --set volume.use_local_path_for_single_pod_cluster.enabled=true \
    --set volume.enabled=false \
    --values=/tmp/mariadb.yaml


# wait for mariadb
cd /opt/openstack-helm && ./tools/deployment/common/wait-for-pods.sh openstack
