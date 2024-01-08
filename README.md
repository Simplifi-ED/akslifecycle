<div align="center">
<h1 align="center">akslifecycle</h1>
<br />
<img alt="License: GPL v3" src="https://img.shields.io/badge/License-GPLv3-blue.svg"/><br>
<br>
This repository provides a command-line interface (CLI) to manage the lifecycle of Azure Kubernetes Service (AKS) node pools. The node pools are defined in a config.yml file. The CLI tool can start and stop these node pools according to a defined schedule, which is also specified in the config.yml file.
</div>

***

## Getting Started

To get started with this CLI tool, you need to have the following prerequisites:

- Azure CLI: You need to have the Azure CLI installed and configured on your machine. You can check your Azure CLI version by running az --version in your terminal. If you need to install or upgrade, refer to the Azure CLI installation guide.
- Go: This project is written in Go, so you need to have Go installed on your machine. You can download Go from the official Go website.
- You need to export these env variables:
  - `AZURE_SUBSCRIPTION_ID`: This is the ID of your Azure subscription.
  - `AZURE_CLIENT_ID`: This is the ID of your Azure Active Directory (AD) app registration.
  - `AZURE_TENANT_ID`: This is the ID of your Azure AD tenant.
  - `AZURE_CLIENT_SECRET`: This is the secret of your Azure AD app registration.

### Installation

```sh
git clone https://github.com/muandane/akslifecycle.git
cd akslifecycle
go build
```

### Usage

Before using the CLI tool, you need to define your node pools and the schedule for starting and stopping them in the config.yml file. Here is an example of what the config.yml file might look like:

```yaml
resources:
- ResourceGroupName: rg-1
  ClusterName: cluster-1
  NodePools:
  - nodepool1
  startSchedule: "0 8 * * 1-5"
  stopSchedule: "0 18 * * 1-5"
- ResourceGroupName: rg-2
  ClusterName: cluster-2
  NodePools:
  - nodepool2
  startSchedule: "0 9 * * 1-5"
  stopSchedule: "0 19 * * 1-5"

```

In this example, nodepool1 is scheduled to start at 8:00 and stop at 18:00, while nodepool2 is scheduled to start at 9:00 and stop at 19:00. The schedules are defined in cron format.

To start the CLI tool, run the akslifecycle executable:

```sh
./akslifecycle
```

This will start the CLI tool, which will then start and stop the node pools according to the schedule defined in the config.yml file.

### License

This project is licensed under the Apache license

### Show your support

Leave a ⭐ if you like this project

***
