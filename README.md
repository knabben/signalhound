# SignalHunter

Hunts Flake And Failing Testgrid Jobs

SignalHunter monitors TestGrid dashboards to identify and summarize test failures and flaking patterns in Kubernetes
CI/CD pipelines. It provides actionable insights for CI signal enumeration, currently focusing on
`sig-release-master-blocking` and `sig-release-master-informing` dashboards.

## Features

### 📊 Test Monitoring Dashboard

Run SignalHunter with the `abstract` command to launch an interactive text user interface (TUI) that displays:

* Board#Tabs combinations in the first panel for easy navigation
* Test listings when selecting specific board combinations
*  Dual information panels:
** Left panel: Slack summary from #release-ci-signal channel (Markdown formatted)
** Right panel: GitHub issue template with Kubernetes defaults (Markdown formatted)

### 📋 Draft issues automatically in the CI Signal Board
Access drafts in the DRAFTING section after selecting a panel and pressing Ctrl-B
Configure with a Personal Access Token (PAT) with appropriate repository permissions

* Clipboard Integration
 
Press Ctrl-Space on any panel to copy content to clipboard
Currently optimized for WSL2 environments

## Installation and Build

Prerequisites

* Go 1.24 or later
* Git
* Github Token PAT
* Kubernetes cluster (kind)

### Running at runtime

```bash
git clone https://github.com/knabben/signalhunter.git
cd signalhunter
make run  # for abstract and standalone
make run-controller # for running the controller outside the cluster
```

### To Deploy on the cluster

**Build and push your image to the location specified by `IMG`:**

```sh
make docker-build docker-push IMG=<some-registry>/signalhound:tag
```

**Install the CRDs into the cluster:**

```sh
make install
```

**Deploy the Manager to the cluster with the image specified by `IMG`:**

```sh
make deploy IMG=<some-registry>/signalhound:tag
```

**Create instances of your solution**

You can apply the samples (examples) from the config/sample:

```sh
kubectl apply -k config/samples/
```

### To Uninstall

**Delete the instances (CRs) from the cluster:**

```sh
kubectl delete -k config/samples/
```

**Delete the APIs(CRDs) from the cluster:**

```sh
make uninstall
```

**UnDeploy the controller from the cluster:**

```sh
make undeploy
```

## Support

Create an issue on GitHub for bug reports and feature requests 
* Join the #sig-release channel in Kubernetes Slack for community support
* Check the documentation for detailed usage guides

## License

Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

--- 

SignalHunter is a community project focused on improving Kubernetes CI reliability through better test signal
monitoring.
