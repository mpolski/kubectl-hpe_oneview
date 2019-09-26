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
  $ kubectl plugin list
The following compatible plugins are available:

/usr/local/bin/kubectl-hpe_oneview
$

  ```
6. If the plugin is listed as above, it's ready to use. Sample commands and outputs (it works with local storage only so far, requires API 1000 since /localStorage is supported from this version onwards):

```
$ kubectl hpe-oneview -compute
NAME                    vCPU                RAM[GB]             STATUS              POWER STATE         MODEL                LOCATION [ENCLOSURE, BAY]
fpd-node3.cec.hpe.com    40                  96                  OK                  On                  Synergy 480 Gen10    Rack3, bay 3
fpd-node10.cec.hpe.com    40                  96                  OK                  On                  Synergy 480 Gen10    Rack3, bay 10
fpd-node1.cec.hpe.com    40                  96                  OK                  On                  Synergy 480 Gen10    Rack3, bay 1
fpd-node5.cec.hpe.com    40                  96                  OK                  On                  Synergy 480 Gen10    Rack3, bay 5
fpd-node9.cec.hpe.com    40                  96                  OK                  On                  Synergy 480 Gen10    Rack3, bay 9
fpd-node4.cec.hpe.com    40                  96                  OK                  On                  Synergy 480 Gen10    Rack3, bay 4
fpd-node6.cec.hpe.com    40                  96                  OK                  On                  Synergy 480 Gen10    Rack3, bay 6
fpd-node8.cec.hpe.com    40                  96                  OK                  On                  Synergy 480 Gen10    Rack3, bay 8
fpd-node2.cec.hpe.com    40                  96                  OK                  On                  Synergy 480 Gen10    Rack3, bay 2
fpd-node7.cec.hpe.com    40                  96                  OK                  On                  Synergy 480 Gen10    Rack3, bay 7

TOTAL NODES:     10

Servers deployed using template: MP_Konvoy_worker_v1.2

$
$ kubectl hpe-oneview -storage

LOCAL STORAGE - DISKS IN COMPUTE NODES OR DY3940 (P416):

  NAME              CONTROLLER / STATUS                      DISK NO.            CAPACITY[GB]        INTERFACE           MEDIA               DISK HEALTH         SERIAL NO.          MODEL
  fpd-node5.cec.hpe.com
                     HPE Smart Array P204i-c SR Gen10 / OK
                                                             0                   1000                SAS                 HDD                 OK                  W471EHMX            MM1000JFJTH
                     HPE Smart Array P416ie-m SR G10 / OK
                                                             0                   960                 SATA                SSD                 OK                  184320F2F539        VK000960GWSXH
                                                             1                   960                 SATA                SSD                 OK                  1843211378E7        VK000960GWSXH
  fpd-node10.cec.hpe.com
                     HPE Smart Array P204i-c SR Gen10 / OK
                                                             0                   480                 SAS                 SSD                 OK                  S3GSNX0K100048      VO000480JWDAR
                     HPE Smart Array P416ie-m SR G10 / OK
                                                             0                   960                 SATA                SSD                 OK                  184321137992        VK000960GWSXH
                                                             1                   960                 SATA                SSD                 OK                  184321137C94        VK000960GWSXH
  fpd-node9.cec.hpe.com
                     HPE Smart Array P204i-c SR Gen10 / OK
                                                             0                   480                 SAS                 SSD                 OK                  S3GSNX0K100226      VO000480JWDAR
                     HPE Smart Array P416ie-m SR G10 / OK
                                                             0                   960                 SATA                SSD                 OK                  1843211379EC        VK000960GWSXH
                                                             1                   960                 SATA                SSD                 OK                  184321137D16        VK000960GWSXH
  fpd-node4.cec.hpe.com
                     HPE Smart Array P204i-c SR Gen10 / OK
                                                             0                   1000                SAS                 HDD                 OK                  W471EHN8            MM1000JFJTH
                     HPE Smart Array P416ie-m SR G10 / OK
                                                             0                   960                 SATA                SSD                 OK                  184321137C55        VK000960GWSXH
                                                             1                   960                 SATA                SSD                 OK                  184321137863        VK000960GWSXH
  fpd-node7.cec.hpe.com
                     HPE Smart Array P204i-c SR Gen10 / OK
                                                             0                   1000                SAS                 HDD                 OK                  W471BAHJ            MM1000JFJTH
                     HPE Smart Array P416ie-m SR G10 / OK
                                                             0                   960                 SATA                SSD                 OK                  184320F2F298        VK000960GWSXH
                                                             1                   960                 SATA                SSD                 OK                  184321137780        VK000960GWSXH
  fpd-node3.cec.hpe.com
                     HPE Smart Array P204i-c SR Gen10 / OK
                                                             0                   1000                SAS                 HDD                 OK                  W471EJ4Y0000E824FWZL    MM1000JFJTH
                     HPE Smart Array P416ie-m SR G10 / OK
                                                             0                   960                 SATA                SSD                 OK                  184321137D17        VK000960GWSXH
                                                             1                   960                 SATA                SSD                 OK                  18432113854A        VK000960GWSXH
  fpd-node6.cec.hpe.com
                     HPE Smart Array P204i-c SR Gen10 / OK
                                                             0                   1000                SAS                 HDD                 OK                  W471EJ8N            MM1000JFJTH
                     HPE Smart Array P416ie-m SR G10 / OK
                                                             0                   960                 SATA                SSD                 OK                  184321137D10        VK000960GWSXH
                                                             1                   960                 SATA                SSD                 OK                  184321137C7C        VK000960GWSXH
  fpd-node8.cec.hpe.com
                     HPE Smart Array P204i-c SR Gen10 / OK
                                                             0                   480                 SAS                 SSD                 OK                  S3GSNX0K100046      VO000480JWDAR
                     HPE Smart Array P416ie-m SR G10 / OK
                                                             0                   960                 SATA                SSD                 OK                  1843211379F4        VK000960GWSXH
                                                             1                   960                 SATA                SSD                 OK                  1843211379E0        VK000960GWSXH
  fpd-node2.cec.hpe.com
                     HPE Smart Array P204i-c SR Gen10 / OK
                                                             0                   1000                SAS                 HDD                 OK                  W471EJL9            MM1000JFJTH
                     HPE Smart Array P416ie-m SR G10 / OK
                                                             0                   960                 SATA                SSD                 OK                  18432113788F        VK000960GWSXH
                                                             1                   960                 SATA                SSD                 OK                  184321137965        VK000960GWSXH
  fpd-node1.cec.hpe.com
                     HPE Smart Array P204i-c SR Gen10 / OK
                                                             0                   480                 SAS                 SSD                 OK                  S3GSNX0K100051      VO000480JWDAR
                     HPE Smart Array P416ie-m SR G10 / OK
                                                             0                   960                 SATA                SSD                 OK                  18432113795D        VK000960GWSXH
                                                             1                   960                 SATA                SSD                 OK                  184321137C42        VK000960GWSXH
$





$ kubectl hpe-oneview -compute
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

$
```

### TO DO
```
kubectl hpe-oneview -addnode 1
```
Would deploy a node using a server profile template specified and then join the cluster.
