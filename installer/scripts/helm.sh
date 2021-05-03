#!/bin/bash

kubectl create serviceaccount --namespace kube-system tiller

kubectl create clusterrolebinding tiller-cluster-rule --clusterrole=cluster-admin --serviceaccount=kube-system:tiller

kubectl patch deploy --namespace kube-system tiller-deploy -p '{"spec":{"template":{"spec":{"serviceAccount":"tiller"}}}}'

helm init --stable-repo-url https://charts.helm.sh/stable --service-account tiller

# run helm serve in background
helm serve &

# wait for tiller
cd /opt/openstack-helm && ./tools/deployment/common/wait-for-pods.sh kube-system


make -C ${HELM_CHART_ROOT_PATH} helm-toolkit

