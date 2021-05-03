# Kupenstack Installation



Prerequisites:

* kubernetes cluster
* Docker
* Kubeconfig file(lets say `config`)

To install kupenstack run command:

```
docker container run -d --name kupenstack_installer -v $(pwd)/config:/root/.kube/config -v $(pwd)/kupenstack.yaml:/etc/kupenstack/kupenstack.yaml  parthyadav/kupenstack:latest 
```



### Quick-Start

* Install docker, if not present

* Install kind, if not present

* Run bash script:

  ```
  #!/bin/bash
  
  # create cluster
  kind create cluster
  
  # prepare kubeconfig file
  cp $HOME/.kube/config $(pwd)/config && sed -i -e 's/https:\/\/127.0.0.1:40525/https:\/\/kind-control-plane:6443/g' $(pwd)/config
  
  # prepare configuration for kupenstack installation
  tee $(pwd)/kupenstack.yaml << EOF
  spec:
    controlNodes:
      - kind-control-plane
  EOF
  
  # run installer
  docker container run --rm --name kupenstack_installer --network kind -v $(pwd)/config:/root/.kube/config -v$(pwd)/kupenstack.yaml:/etc/kupenstack/kupenstack.yaml  parthyadav/kupenstack:latest
  ```

