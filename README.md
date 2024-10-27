## stalker
Stalk and Hunt Flake Testgrid Jobs 

Summarizes failures in the testgrid board for CI signal enumeration, currently fetching `sig-release-master-blocking` and `sig-release-master-informing`


```
./stalker abstract

* [sig-release-master-blocking] gce-cos-master-reboot    https://testgrid.k8s.io/sig-release-master-blocking#gce-cos-master-reboot
Tab stats: 8 of 9 (88.9%) recent columns passed (304 of 306 or 99.3% cells)

- ci-kubernetes-e2e-gci-gce-reboot.Overall                                                                                                                                                                                                                   --
 R 2024-09-30 14:22:22 -0300 -03 Build still running...

- kubetest.diffResources                                                                                                                                                                                                                                     --
 F 2024-09-30 09:44:22 -0300 -03 Error: 2 leaked resources
+NAME                     MACHINE_TYPE   PRE...on-template  e2-standard-2               2024-09-30T05:11:15.602-07:00
S 2024-09-30 05:06:22 -0300 -03 Error: 53 leaked resources
+NAME               ADDRESS/RANGE  TYPE    ...S    1000      tcp,udp,icmp,esp,ah,sctp                          False

- kubetest.TearDown                                                                                                                                                                                                                                          --
 F 2024-09-30 05:06:22 -0300 -03 error during ./hack/e2e-internal/e2e-down.sh (interrupted): signal: interrupt

- kubetest.Timeout                                                                                                                                                                                                                                           --
 F 2024-09-30 05:06:22 -0300 -03 kubetest --timeout triggered

- Kubernetes e2e suite.[It] [sig-cloud-provider-gcp] Reboot [Disruptive] [Feature:Reboot] each node by ordering unclean reboot and ensure they function upon restart                                                                                         --
 F 2024-09-23 13:18:20 -0300 -03 [FAILED] Test failed; at least one node failed to reboot in the time g...detected after the initial failure. These are visible in the timeline


- kubetest.Test                                                                                                                                                                                                                                              --
 F 2024-09-23 13:18:20 -0300 -03 error during ./hack/ginkgo-e2e.sh --ginkgo.focus=\[Feature:Reboot\] --...=8 --report-dir=/logs/artifacts --disable-log-dump=true: exit status 1

- kubetest.Up                                                                                                                                                                                                                                                --
 F 2024-09-16 08:54:32 -0300 -03 error during ./hack/e2e-internal/e2e-up.sh: exit status 2
F 2024-09-16 08:24:33 -0300 -03 error during ./hack/e2e-internal/e2e-up.sh: exit status 2
F 2024-09-16 07:54:32 -0300 -03 error during ./hack/e2e-internal/e2e-up.sh: exit status 2
F 2024-09-16 07:24:32 -0300 -03 error during ./hack/e2e-internal/e2e-up.sh: exit status 2

```
