# kubectl-hpe_oneview
## [kubectl extension plugin](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/) for [HPE OneView](https://www.hpe.com/us/en/integrated-systems/software.html) example.

This is a demonstration of writing Kubernetes kubectl extenstion plugin to view some of the physical infrastructure parameters of [HPE Synergy](https://www.hpe.com/pl/en/integrated-systems/synergy.html) hardware platform via HPE OneView API and [Go Language bindings](https://github.com/HewlettPackard/oneview-golang) from with some additions.

Here's the instructions to test this in one's environment:

### Prerequsites:
1. Install Go and git
2. Access to HPE Synergy with Server Profile Template that is used to craete profiles that will run Kubernetes nodes.
2. Set the following environment variables:
```
OV_ENDPOINT=<e.g. https://....>
OV_USERNAME=<your username>
OV_PASSWORD=<your password>
OV_AUTHLOGINDOMAIN=<domain> //optional, can be empty
OV_PROFILETEMPLATE=<template_name>
```

### Instructions:
1. clone this repo.
2. install modules:
  ```
  go get github.com/mpolski/oneview-golang-temp
  go get github.com/HewlettPackard/oneview-golang
  ```
3. build executable:
  ```
  go build kubectl-hpe_oneview.go
  ```
4. move the exectuable to a folder in PATH, allow for it to be executed
  ```
  chmod +x kubectl-hpe_oneview.go
  ```
5. Verfiy kubectl can see the plugin
  ```
  kubectl plugin list
  ```
6. If the plugin is listed, it's ready to use. Sample command and output:

```
mpolski@fpd-jumphost:~$ kubectl hpe-oneview -compute
NAME                     vCPU                RAM[GB]             STATUS              POWER STATE         MODEL                LOCATION [ENCLOSURE, BAY]
fpd-node10.cec.hpe.com    40                  96                  OK                  On                  Synergy 480 Gen10    Rack3, bay 10
fpd-node2.cec.hpe.com    40                  96                  OK                  On                  Synergy 480 Gen10    Rack3, bay 2
fpd-node1.cec.hpe.com    40                  96                  OK                  On                  Synergy 480 Gen10    Rack3, bay 1
fpd-node7.cec.hpe.com    40                  96                  OK                  On                  Synergy 480 Gen10    Rack3, bay 7
fpd-node3.cec.hpe.com    40                  96                  OK                  On                  Synergy 480 Gen10    Rack3, bay 3
fpd-node5.cec.hpe.com    40                  96                  OK                  On                  Synergy 480 Gen10    Rack3, bay 5
fpd-node9.cec.hpe.com    40                  96                  OK                  On                  Synergy 480 Gen10    Rack3, bay 9
fpd-node4.cec.hpe.com    40                  96                  OK                  On                  Synergy 480 Gen10    Rack3, bay 4
fpd-node6.cec.hpe.com    40                  96                  OK                  On                  Synergy 480 Gen10    Rack3, bay 6
fpd-node8.cec.hpe.com    40                  96                  OK                  On                  Synergy 480 Gen10    Rack3, bay 8

TOTAL NODES:     10

Servers deployed using template: MP_Konvoy_worker_v1.2

mpolski@fpd-jumphost:~$
```

### TO DO
```
kubectl hpe-oneview -addnode 1
```
Would deploy a node using a server profile template specified and then join the cluster.
