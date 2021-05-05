#!/bin/bash

OPENSTACK_HELM_ROOT_PATH=/opt/openstack-helm


########
# Libvirt
########

make -C ${HELM_CHART_ROOT_PATH} libvirt

tee /tmp/libvirt.yaml << EOF
network:
  backend: 
    - linuxbridge
EOF

helm upgrade --install libvirt ${HELM_CHART_ROOT_PATH}/libvirt \
  --namespace=openstack \
  --set conf.ceph.enabled=false \
  --values=/tmp/libvirt.yaml



########
# PLACEMENT
########

make -C ${OPENSTACK_HELM_ROOT_PATH} placement

helm upgrade --install placement ${OPENSTACK_HELM_ROOT_PATH}/placement \
    --namespace=openstack \
    --set pod.replicas.api=1



########
# NOVA
########

case "${OPENSTACK_RELEASE}" in
  "queens")
    DEPLOY_SEPARATE_PLACEMENT="no"
    ;;
  "rocky")
    DEPLOY_SEPARATE_PLACEMENT="no"
    ;;
  "stein")
    DEPLOY_SEPARATE_PLACEMENT="yes"
    ;;
  *)
    DEPLOY_SEPARATE_PLACEMENT="yes"
    ;;
esac


if [[ "${DEPLOY_SEPARATE_PLACEMENT}" == "yes" ]]; then
  OSH_EXTRA_HELM_ARGS_NOVA="--values=${OPENSTACK_HELM_ROOT_PATH}/nova/values_overrides/train-disable-nova-placement.yaml"
fi


tee /tmp/nova.yaml << EOF
network:
  backend: 
    - linuxbridge
pod:
  replicas:
    osapi: 1
    conductor: 1
    consoleauth: 1
EOF

make -C ${OPENSTACK_HELM_ROOT_PATH} nova

helm upgrade --install nova ${OPENSTACK_HELM_ROOT_PATH}/nova --namespace=openstack \
      --values=/tmp/nova.yaml \
      --set bootstrap.wait_for_computes.enabled=true \
      --set conf.ceph.enabled=false \
      --set conf.nova.libvirt.virt_type=qemu \
      --set conf.nova.libvirt.cpu_mode=none \
      ${OSH_EXTRA_HELM_ARGS_NOVA}


########
# NEUTRON
########

make -C ${OPENSTACK_HELM_ROOT_PATH} neutron

tee /tmp/neutron.yaml << EOF
network:
  backend: 
    - linuxbridge
dependencies:
  dynamic:
    targeted:
      linuxbridge:
        dhcp:
          pod:
            - requireSameNode: true
              labels:
                application: neutron
                component: neutron-lb-agent
        l3:
          pod:
            - requireSameNode: true
              labels:
                application: neutron
                component: neutron-lb-agent
        metadata:
          pod:
            - requireSameNode: true
              labels:
                application: neutron
                component: neutron-lb-agent
        lb_agent:
          pod: null
conf:
  neutron:
    DEFAULT:
      interface_driver: linuxbridge
  dhcp_agent:
    DEFAULT:
      interface_driver: linuxbridge
  l3_agent:
    DEFAULT:
      interface_driver: linuxbridge
EOF


helm upgrade --install neutron ${OPENSTACK_HELM_ROOT_PATH}/neutron \
    --namespace=openstack \
    --values=/tmp/neutron.yaml

# wait for everything
cd /opt/openstack-helm && ./tools/deployment/common/wait-for-pods.sh openstack

