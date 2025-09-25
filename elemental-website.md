### AI Context

This is the official Elemental documentation. Elemental is a software stack enabling centralized, full cloud-native OS management with Kubernetes.

# Merged Knowledge Base for AI

## Article: airgap.md

---
sidebar_label: Air-Gapped Installation
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/airgap"/>
</head>

## Install Elemental in an Air-Gapped Environment

### Assumptions
A Rancher air-gapped installation should be already configured as per the [official Rancher documentation](https://ranchermanager.docs.rancher.com/pages-for-subheaders/air-gapped-helm-cli-install).
In particular, a private registry should be available in the air-gapped infrastructure.

### Overview
In order to run Elemental in an air-gapped environment the following artifacts are needed:
- the Elemental Operator charts
- the container images referenced in the charts (the *elemental-operator* and *seedimage-builder* images)
- the containerized OS images

Moreover, it could be handy to create a *channel image* referencing the containerized OS images available.
The official channel image (the *elemental-channel* one) references absolute URLs of the OS images on the official suse registry, so it cannot be used as-is in an air-gapped scenario.

### Elemental Air-Gapped installation from the command line
All the required steps can be accomplished by executing the
[`elemental-airgap.sh` script](https://raw.githubusercontent.com/rancher/elemental-operator/main/scripts/elemental-airgap.sh)
from a host with Internet access.

The Elemental charts are a required parameter to the script and can be provided as downloaded archives, URLs or as one of
the `stable`, `staging` and `dev` keywords, to let the script retrieve the correct chart version for you.

`elemental-airgap.sh` inspects the Elemental Operator chart, identifies all the required container images, downloads and saves them in a single docker archive.
It also builds a new OS channel image with the OS image URLs pointing to the private registry passed as argument
(which is a mandatory argument too).

The latest version of the elemental script can be easily downloaded from the official github repo:
```shell showLineNumbers
wget https://raw.githubusercontent.com/rancher/elemental-operator/main/scripts/elemental-airgap.sh
chmod 755 elemental-airgap.sh
```

Let's now download all the artifacts and build a custom channel from the latest stable release of Elemental:

<Tabs>
<TabItem value="dockerArchive" label="Create a Docker archive" default>

```shell showLineNumbers
./elemental-airgap.sh stable -r <REGISTRY.YOURDOMAIN.COM:PORT>
```

once completed (the script may take a while) the following files will be available in the current dir:
- `elemental-operator-crds-chart-<*VERSION*>.tgz`
- `elemental-operator-chart-<*VERSION*>.tgz`
- `elemental-images.txt`
- `elemental-images.tar.gz`

</TabItem>
<TabItem value="haulerArchive" label="Create a Hauler archive" default>

```shell showLineNumbers
./elemental-airgap.sh -ha stable -r <REGISTRY.YOURDOMAIN.COM:PORT>
```
once completed (the script may take a while) both the charts and the container images will be packed in the
hauler archive named `elemental-haul.tar.zst` .

</TabItem>
</Tabs>

#### Elemental installation
The files and archives created by the script should be copied to a host which:
- Has access to the private registry.
- Has the kubectl binary installed and configured to access the air-gapped Rancher cluster.
- Has the helm binary installed.


<Tabs>
<TabItem value="dockerArchive" label="Install from a Docker archive" default>

If the private registry requires authentication you need to log with docker into it:
```shellnocolor showLineNumbers
docker login <REGISTRY.YOURDOMAIN.COM:PORT>
```
Two steps are needed to perform the Elemental installation:
1. load the archive with all the required container images on the private registry:
this could be done using the `rancher-load-images.sh` script distributed with the Rancher release and already used for the Rancher air-gapped deployment:
```shellnocolor showLineNumbers
rancher-load-images.sh \
   --image-list elemental-images.txt \
   --images elemental-images.tar.gz \
   --registry <REGISTRY.YOURDOMAIN.COM:PORT>
```
2. install the downloaded elemental charts configuring the local registry and the newly created channel:
```shellnocolor showLineNumbers
helm upgrade --create-namespace -n cattle-elemental-system \
  --install elemental-operator-crds elemental-operator-crds-chart-<VERSION>.tgz

helm upgrade --create-namespace -n cattle-elemental-system \
  --install elemental-operator elemental-operator-chart-<VERSION>.tgz \
  --set registryUrl=<REGISTRY.YOURDOMAIN.COM:PORT> \
  --set channel.repository=rancher/elemental-channel-<REGISTRY.YOURDOMAIN.COM>
```
</TabItem>
<TabItem value="haulerArchive" label="Install from a Hauler archive" default>
To install from a [Hauler](https://rancherfederal.github.io/hauler-docs/) archive (`-ha` option in `elemental-airgap.sh`)
Hauler installation is also a requirement on the host from where the installation is performed.

If the private registry requires authentication you need to log with Hauler into it:
```shellnocolor showLineNumbers
hauler login <REGISTRY.YOURDOMAIN.COM:PORT> -u $USERNAME -p $PASSWORD
```

Three steps are needed to perform the Elemental installation:
1. Load the 'elemental-haul.tar.zst' Haul archive in the Hauler instance in the airgapped infrastructure:
```shellnocolor
hauler store load 'elemental-haul.tar.zst'
```
2. If the local registry in the air-gapped environment is not server by Hauler,
load the Haul archive in the local registry:
```shellnocolor
hauler store copy registry://<REGISTRY.YOURDOMAIN.COM:PORT>
```
:::info Hauler can also serve as a registry
In case the air-gapped local registry is served by an Hauler instance, just load the Haul archive directly there
(as shown in step (1)) and skip step (2).
:::

3. Extract the elemental charts from the Hauler store and install them:
```shellnocolor
hauler store extract elemental-operator-crds-chart-<ELEMENTAL-VERSION>.tgz

hauler store extract elemental-operator-chart-<ELEMENTAL-VERSION>.tgz

helm upgrade --create-namespace -n cattle-elemental-system \
  --install elemental-operator-crds elemental-operator-crds-chart-<ELEMENTAL-VERSION>.tgz

helm upgrade --create-namespace -n cattle-elemental-system \
  --install elemental-operator elemental-operator-chart-<ELEMENTAL-VERSION>.tgz \
  --set registryUrl=<REGISTRY.YOURDOMAIN.COM:PORT> \
  --set channel.repository=rancher/elemental-channel-<REGISTRY.YOURDOMAIN.COM:PORT>
```
</TabItem>
</Tabs>

:::info The elemental airgap script outputs the required commands
The `elemental-airgap.sh` scripts prints out the required commands shown above but using the actual chart version and the provided registry URL to allow to easily copy and paste the exact commands.
:::

### Elemental Air-Gapped installation from the Rancher Marketplace
A Rancher air-gapped installation includes also the Elemental Operator charts and the operator and seedimage container
images.

To collect the missing OS images and to build an OS channel image for your private registry execute the
[`elemental-airgap.sh` script](https://raw.githubusercontent.com/rancher/elemental-operator/main/scripts/elemental-airgap.sh)
from an host with Internet access, using the `-co` option.

As an example, let's target the `elemental-channel` image from the latest stable release of Elemental.
The script will take care of downloading the Elemental operator chart (if needed), extract the OS channel image URL,
download it, inspect all the OS images referenced, download all of them and create a new OS channel with links to the
private registry of the air-gapped scenario.
<Tabs>
<TabItem value="dockerArchive" label="Create a Docker archive" default>
```shell showLineNumbers
wget https://raw.githubusercontent.com/rancher/elemental-operator/main/scripts/elemental-airgap.sh
chmod 755 elemental-airgap.sh
./elemental-airgap.sh stable -co -r <REGISTRY.YOURDOMAIN.COM:PORT>
```
once completed (the script may take a while) the following files will be available in the current dir:
- `elemental-operator-crds-chart-<*VERSION*>.tgz`
- `elemental-operator-chart-<*VERSION*>.tgz`
- `elemental-images.txt`
- `elemental-images.tar.gz`
</TabItem>
<TabItem value="haulerArchive" label="Create a Hauler archive" default>
```shell showLineNumbers
./elemental-airgap.sh -ha -co stable -r <REGISTRY.YOURDOMAIN.COM:PORT>
```
once completed (the script may take a while) the container images will be packed in the
hauler archive named `elemental-haul.tar.zst`.
</TabItem>
</Tabs>

#### Elemental installation
The generated archive should be loaded to the air-gapped private registry.

<Tabs>
<TabItem value="dockerArchive" label="Install from a Docker archive" default>

If the private registry requires authentication you need to log with docker into it:
```shellnocolor showLineNumbers
docker login <REGISTRY.YOURDOMAIN.COM:PORT>

The script will print out the commands required to load the images via the Rancher `rancher-load-images.sh` tool, used
for the Rancher air-gapped installations. It should be something like:

```shell showLineNumbers
NEXT STEPS:

1) Load the 'elemental-images.tar.gz' to the local registry (<REGISTRY.YOURDOMAIN.COM:PORT>)
   available in the airgapped infrastructure:

./rancher-load-images.sh \
   --image-list elemental-images.txt \
   --images elemental-images.tar.gz \
   --registry <REGISTRY.YOURDOMAIN.COM:PORT>
```
Once the OS and channel images are loaded, you should skip the point (2) from the script output
(which will install the Elemental charts from the downloaded archives)
and instead perform the Elemental Operator installation from the Rancher UI.
</TabItem>
<TabItem value="haulerArchive" label="Install from a Hauler archive" default>
If the private registry requires authentication you need to log with Hauler into it:
```shellnocolor showLineNumbers
hauler login <REGISTRY.YOURDOMAIN.COM:PORT> -u $USERNAME -p $PASSWORD
```

The script will print out the commands required to load the images. It should be something like:

```shell showLineNumbers
NEXT STEPS:

1. Load the 'elemental-haul.tar.zst' Haul archive in the Hauler instance in the airgapped infrastructure:

hauler store load 'elemental-haul.tar.zst'

2. If the local registry in the air-gapped environment is not server by Hauler,
load the Haul archive in the local registry:

hauler store copy registry://<REGISTRY.YOURDOMAIN.COM:PORT>
```
Once the OS and channel images are loaded, you should skip the point (3) from the script output
(which will install the Elemental charts from the downloaded archives)
and instead perform the Elemental Operator installation from the Rancher UI.
</TabItem>
</Tabs>

When requested, put the full path of the OS channel image just uploaded in your private registry:
![Elemental OS Channel](images/airgap-os-channel-image.png)

### Elemental UI Extension
Rancher 2.7.x doesn't support UI extensions plugin in air-gapped environments, and so the Elemental UI is not available in Rancher 2.7.x.

The Elemental UI plugin will be present in the available UI extensions in Rancher 2.8.0.


---

## Article: architecture-clusterdeployment.md

---
sidebar_label: Kubernetes cluster provisioning
title: ''
---

## Kubernetes cluster provisioning
The goal of the Kubernetes Cluster deployment phase is to create a new RKE2 or K3s cluster using the available [MachineInventories](machineinventory-reference.md), i.e., the hosts that have successfully completed the [Machine onboarding](architecture-machineonboarding.md) phase.

The Elemental Kubernetes cluster deployment involves the following steps:
* The user creates a [MachineInventorySelectorTemplate](machineinventoryselectortemplate-reference.md) resource: it allows to define a _selector_ to identify a subset of the available [MachineInventories](machineinventory-reference.md) based on the value of their labels.
* The user defines a [Rancher cluster](cluster-reference.md) and adds to the `machinePools` definition a reference to the [MachineInventorySelectorTemplate](machineinventoryselectortemplate-reference.md) created in the step before.
* The [Rancher RKE2/K3s Cluster provisioning](https://ranchermanager.docs.rancher.com/how-to-guides/new-user-uuides/launch-kubernetes-with-rancher#rke2) reacts to the Rancher cluster resource creation by generating a number of [MachineInventorySelectors](machineinventoryselector-reference.md) resources equal to the _quantity_ specified in the _machinePools_.
* The Elemental Operator pairs each generated [MachineInventorySelector](machineinventoryselector-reference.md) resource with an available [MachineInventory](machineinventory-reference.md) and installs the [rancher-system-agent](https://github.com/rancher/system-agent) daemon on the host tracked by the [MachineInventory](machineinventory-reference.md).
The [Rancher RKE2/K3s Cluster provisioning](https://ranchermanager.docs.rancher.com/how-to-guides/new-user-uuides/launch-kubernetes-with-rancher#rke2) takes over the K3s/RKE provisioning using [rancher-system-agent](https://github.com/rancher/system-agent) _plans_: it installs the required components (e.g., containerd, K3s, ...) and creates the configuration files till the successful deployment of the new Kubernetes cluster.


---

## Article: architecture-components.md

---
sidebar_label: Elemental components
title: ''
---

# Elemental components

The components required to provide the [Elemental services](architecture-services.md) are:
* the [``elemental``](#elemental-command-line-tool) command line tool
* the [``elemental-operator``](#elemental-operator-daemon) daemon
* the [``elemental-register``](#elemental-register-command-line-tool) command line tool
* the [``elemental-system-agent``](#elemental-system-agent-daemon) daemon
* one or more [Elemental OS container images](#elemental-os-container-image)

### ``elemental`` command line tool
The ``elemental`` tool is part of the <Vars name="elemental_toolkit_name" link="elemental_toolkit_url" /> project.
It performs the actual OS installation and upgrade operations on the host and is used to execute the [cloud-config](cloud-config-reference.md) directives added in the [Elemental CRDs](custom-resources.md).

The ``elemental`` binary is included in all the base OS images distributed with Elemental.

### ``elemental-operator`` daemon
The **elemental-operator** daemon performs two main tasks:
1. embeds the Elemental Kubernetes controllers to manage all the [Elemental CRDs](custom-resources.md)
2. exposes the _registration endpoints_ to allow the host to register and download the OS installation configuration during the [machine onboarding](architecture-machineonboarding.md)

The `elemental-operator` daemon is deployed on the Rancher cluster as a `Deployment` via the [Elemental Operator Helm Chart](elementaloperatorchart-reference.md).

### ``elemental-register`` command line tool
The **elemental-register** binary is the client used to register the host against the _registration endpoints_ exposed by the [elemental-operator](#elemental-operator-daemon). It collects and forwards the host data to allow the [elemental-operator](#elemental-operator-daemon) to fill the [SMBIOS](smbios.md) and the [Hardware](hardwarelabels.md) label templates.

If the registraton phase is performed successfully, the **elemental-register** gets the full configuration stored in the [MachineRegistration](machineregistration-reference.md) from the **elemental-operator**.
As the last step, the **elemental-register** client calls the [elemental](#elemental-command-line-tool) binary passing the retrieved configuration to kick off the OS installation.

### ``elemental-system-agent`` daemon
The ``elemental-system-agent`` is built from the [Rancher System Agent project](https://github.com/rancher/system-agent) and allows Elemental to deploy _plans_ to assist with the host provisioning.
Notably, the ``rancher-system-agent`` installation and configuration required for the [Kubernetes Cluster provisioning service](architecture-clusterdeployment.md) is performed through an ``elemental-system-agent`` _plan_.

### Elemental OS container image
An <Vars name="elemental_toolkit_name" link="elemental_toolkit_url" /> OS image is an OCI container image containing all the files that will make up the OS of the target host. It contains not only all the binaries and libraries, but also the kernel and the boot files required by a linux system.

The Elemental OS image is an opinionated <Vars name="elemental_toolkit_name" link="elemental_toolkit_url" /> OS image which is based on [SLE Micro](https://www.suse.com/products/micro/) and contains specific Elemental configurations and binaries (the [elemental](#elemental-command-line-tool) and the [elemental-register](#elemental-register-command-line-tool) ones).

The Elemental OS images are tracked in the [ManagedOSVersions](managedosversion-reference.md) resources. The [ManagedOSVersions](managedosversion-reference.md) resources are dynamically created from [ManagedOSVersionChannel](managedosversionchannel-reference.md) resources. A default[ManagedOSVersionChannel](managedosversionchannel-reference.md) resource is deployed with each Elemental Operator installation.

---

## Article: architecture-machineonboarding.md

---
sidebar_label: Machine onboarding
title: ''
---

## Machine onboarding
The goal of the Machine Onboarding service is to register a new machine in the Elemental Operator catalog ([MachineInventory](machineinventory-reference.md) resources) and provision the host with the OS, applying the desired configuration.

The Machine Onboarding involves the following steps:
* The user creates a [MachineRegistration](machineregistration-reference.md) resource: the MachineRegistration includes the configuration required for the OS installation.
As soon as the MachineRegistration is created, the Elemental Operator exposes an HTTP _registration endpoint_: it is the entrypoint required by the onboarding hosts to register to the operator.
* The user requests a self-installing image via a [SeedImage](seedimage-reference.md) resource: it only requires a reference to a MachineRegistration resource and to an OS container image compatible with the Elemental Toolkit.
As soon as the SeedImage is created, the Elemental Operator triggers the build process of the self-installing image: once completed, the URL is exposed in the `status.downloadURL` field of the [SeedImage](seedimage-reference.md) resource.
* The user downloads the self-installing image and uses it to boot an unprovisioned host:
the host [authenticates](https://elemental.docs.rancher.com/authentication) to the _registration endpoint_ of the Elemental Operator, gets the full configuration stored in the MachineRegistration and installs the OS on the local storage device. As soon as the host has completed the registration process, the Elemental Operator creates a unique [MachineInventory](machineinventory-reference.md) resource tracking the host.
The self-installing image can be used to onboard any number of hosts.

---

## Article: architecture-services.md

# Services

Elemental functionality is provided through the following list of services:
* [Machine Onboarding](architecture-machineonboarding.md): allows hosts to be OS provisioned and registered in the Elemental Operator catalog.
* [Kubernetes Cluster Deployment](architecture-clusterdeployment.md): triggers the creation of a Kubernetes Cluster on the available host registered in the Elemental Operator catalog. 

---

## Article: architecture.md

---
sidebar_label: Overview
title: ''
---

import useBaseUrl from '@docusaurus/useBaseUrl';
import ThemedImage from '@theme/ThemedImage';

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/architecture"/>
</head>

# Architecture overview
Elemental is the combination of two main projects: the <Vars name="elemental_toolkit_name" link="elemental_toolkit_url" /> and the <Vars name="elemental_operator_name" link="elemental_operator_url" />.

The <Vars name="elemental_toolkit_name" link="elemental_toolkit_url" /> enables OS installation and updates from OCI container images, so the OS can be distributed and retrieved via container registries.

The <Vars name="elemental_operator_name" link="elemental_operator_url" /> extends Rancher with OS provisioning and OS management functionalities, leveraging the Elemental Toolkit.
It bridges the gap between the <Vars name="elemental_toolkit_name" link="elemental_toolkit_url"/> and the
[Rancher RKE2/K3s Cluster provisioning](https://ranchermanager.docs.rancher.com/how-to-guides/new-user-guides/launch-kubernetes-with-rancher#launching-kubernetes-on-new-nodes-in-an-infrastructure-provider-1),
to provide a seamless experience going from hosts without OS to fully configured Kubernetes Clusters.

This is achieved through a [set of services](architecture-services.md) offered by Elemental that are made possible thanks to the [components](architecture-components.md) borrowed from the Elemental Toolkit and the Elemental Operator projects.

<ThemedImage
  alt="Elemental Architecture"
  sources={{
    light: useBaseUrl('img/elemental-architecture-v1.5.png'),
    dark: useBaseUrl('img/elemental-architecture-v1.5-dark.png'),
  }}
/>

---

## Article: authentication.md

---
sidebar_label: Authentication
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/authentication"/>
</head>

import RegistrationEmulateTpm from "!!raw-loader!@site/examples/authentication/registration-emulate-tpm.yaml"
import RegistrationMac from "!!raw-loader!@site/examples/authentication/registration-mac.yaml"
import RegistrationSysUUID from "!!raw-loader!@site/examples/authentication/registration-sys-uuid.yaml"

# Authentication

Authentication happens during Machine onboarding, when the Machine _registers_ to the
<Vars name="elemental_operator_name" />.

Elemental by default authenticates hosts through the _Trusted Platform Module_ (TPM):
the machine is authenticated through _attestation_, i.e., the machine proofs its identity through its
TPM device.

In order for _attestation_ to work, each onboarding machine
**must have a TPM 2.0 device**,
otherwise would not be able to register using secure TPM authentication.

## TPM alternatives
**The only officially supported registration method is based on TPM attestation and requires devices TPM 2.0 enabled.**

If you want to enroll devices without a TPM 2.0 chip bypassing secure authentication, there are multiple ways to uniquely _identify_ those machines and allow registration:
* emulating a TPM device via a simple software implementation
* identifying themselves through their network MAC address
* identifying themselves using their SMBIOS UUID

The authentication/identification method can be specified in the
[config:elemental:registration:auth field](machineregistration-reference.md#configelementalregistration) of the MachineRegistration resource.

:::warning
The only secure and officially supported **authentication** method in Elemental is the default one, based on TPM attestation.
The TPM alternatives can be used for demo purposes or local deployments but are not recommended for production use as the onboarding machines identity is not securely verified.
:::

### TPM emulation
TPM emulation performs authentication using a software that mimics TPM behavior and which is embedded in the <Vars name="elemental_register_name" />.
The keys of the emulated TPM device are all generated by a single seed in a deterministic way: same seed results in the same TPM keys, so a different seed should be picked up in each enrolling host.

TPM emulation is enabled configuring the `emulate-tpm` and `emulated-tpm-seed` fields in the `MachineRegistration` configuration (see the [config:elemental:registration section in the MachineRegistration reference](machineregistration-reference.md#configelementalregistration) for more details).

<CodeBlock language="yaml" title="example MachineRegistration using TPM emulation" showLineNumbers>{RegistrationEmulateTpm}</CodeBlock>

TPM emulation configuration is detailed in the [TPM emulation configuration section](tpm.md#add-tpm-emulation-to-bare-metal-machine).

### MAC address identification
When using MAC address identification, the host registers to the <Vars name="elemental_operator_name" /> using the MAC address from its Network Interface Card (NIC) as an identifier.
In case the machine has more than one network interface, the MAC addresses are sorted lexicographically and the first one is selected.

To replace TPM authentication with MAC address identification, it is enough to set the `mac` value to the `auth` field in the [config:elemental:registration section in the MachineRegistration reference](machineregistration-reference.md#configelementalregistration).

<CodeBlock language="yaml" title="example MachineRegistration using the MAC address as machine identifier" showLineNumbers>{RegistrationMac}</CodeBlock>

:::warning
The MAC address is considered unique by the <Vars name="elemental_operator_name" />.
This is true for physical devices, while if using VirtualMachines from different hypervisors and different network segments it is up to the administrator to ensure that the registering VMs have a unique MAC address.
:::

### SMBIOS UUID identification
The System Management BIOS (SMBIOS) specification defines data structures that can be used to read management information produced by the BIOS of a host.

When using the `sys-uuid` value as the `auth` field of the [config:elemental:registration section in the MachineRegistration](machineregistration-reference.md#configelementalregistration), the host registers to the <Vars name="elemental_operator_name" /> using the `UUID` value from the `System Information` table of the host SMBIOS data.

<CodeBlock language="yaml" title="example MachineRegistration using the UUID from the SMBIOS System Information table as machine identifier" showLineNumbers>{RegistrationSysUUID}</CodeBlock>

:::warning 
The SMBIOS `System information/UUID` value should be filled by the hardware vendor as a unique UUID for the host.

The SMBIOS data is not always reliable. This depends on the manufacturer. You may experience the UUID being missing, or the same UUID being applied to multiple devices within the same batch.

It is up to the administrator to ensure that the machines have unique `System information/UUID` SMBIOS values (the  `dmidecode` tool could be of help), otherwise the machines will keep overwriting the same `MachineInventory` resource and the Elemental provisioning will fail.
:::


---

## Article: backup.md

---
sidebar_label: Backup
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/backup"/>
</head>

# Backup

Since Elemental runs as part of Rancher, the Elemental resources are bundled in the Rancher backup.  
For more details about Rancher backups, restore, and disaster recovery options, please follow the official [Rancher documentation](https://ranchermanager.docs.rancher.com/pages-for-subheaders/backup-restore-configuration).

## Install rancher-backup operator for Rancher

Follow the [Rancher backup guide](https://docs.ranchermanager.rancher.io/how-to-guides/new-user-guides/backup-restore-and-disaster-recovery/back-up-rancher) to learn how to install and configure the Rancher backup-operator.  

Note that for single node Rancher installations the backup workflow is different.  
You may follow the official [documentation](https://ranchermanager.docs.rancher.com/v2.6/how-to-guides/new-user-guides/backup-restore-and-disaster-recovery/back-up-docker-installed-rancher) to learn more.

## Backup Elemental with rancher-backup operator

Create a `backup object` (adapted to your needs) to backup Rancher running on a Kubernetes cluster.  

```yaml showLineNumbers
apiVersion: resources.cattle.io/v1
kind: Backup
metadata:
  name: rancher-backup
spec:
  resourceSetName: rancher-resource-set
  schedule: "10 3 * * *"
  retentionCount: 10
```

The rancher-backup operator offers several options for schedule, encryption, and storage classes.  
You can explore all options by reading the [official documentation](https://ranchermanager.docs.rancher.com/reference-guides/backup-restore-configuration/backup-configuration).

Check logs from rancher-backup operator.

```shell showLineNumbers
kubectl logs -n cattle-resources-system -l app.kubernetes.io/name=rancher-backup -f
```

Verify if backup file was created on Persistent Volume.

```shell showLineNumbers
...
INFO[2022/10/17 07:45:04] Finding files starting with /var/lib/backups/rancher-backup-430169aa-edde-4a61-85e8-858f625a755b*.tar.gz 
INFO[2022/10/17 07:45:04] File rancher-backup-430169aa-edde-4a61-85e8-858f625a755b-2022-10-17T05-15-00Z.tar.gz was created at 2022-10-17 0
...
```


---

## Article: certificate-authority.md

---
sidebar_label: Certificate Authority Verification
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/certificate-authority"/>
</head>

# Certificate Authority Verification

The `elemental-register` and `elemental-system-agent` rely on the [Rancher's Certificate Authority configuration](https://ranchermanager.docs.rancher.com/getting-started/installation-and-upgrade/resources/update-rancher-certificate) to verify the [MachineRegistration](https://elemental.docs.rancher.com/machineregistration-reference#configelementalregistration) URL, and to remotely watch plans.  

Depending on whether the Certificate Authority is private or public, you may want to instruct the agent to enforce `strict` CA verification, or to use the system trust store instead.  

From Rancher `2.9`, the [agent-tls-mode](https://ranchermanager.docs.rancher.com/getting-started/installation-and-upgrade/installation-references/tls-settings#agent-tls-enforcement) global setting will also apply to the installation of Elemental agents.  
Note that if the `agent-tls-mode` setting changes, Elemental machines will need to be [reset](./reset.md) in order for the setting to apply.  

## Private CA certificate lifecycle

When using a private CA, the recommendation is to always make sure that the same CA is also used for Rancher.  
Elemental will make use of the `cacerts`, when including the CA cert to be trusted by agents. This is the same value as it appears on the `https://my.rancher.example/cacerts` URL.  

Note however that it will not be possible to update this value after an Elemental machine has been installed.  
Replacing the CA certificate on Rancher may lead to Elemental machines not being able to re-connect to Rancher and operating normally, when the `agent-tls-mode` is set to `strict`.  

For this reason the recommendation is to use the `agent-tls-mode: system-store` setting instead and manage the lifecycle of CA certs on Elemental machines directly, when a private Certificate Authority is in use.  

The CA cert can be installed in [custom OS images](./custom-images.md) directly, or passed as a cloud-init configuration in Elemental resources.  
For example the initial CA certificate can be included in the [MachineRegistration](./machineregistration-reference.md#configcloud-config) `cloud-config`:  

```yaml
apiVersion: elemental.cattle.io/v1beta1
kind: MachineRegistration
metadata:
  name: fire-nodes
  namespace: fleet-default
spec:
  config:
    cloud-config:
      write_files:
        - path: /etc/pki/trust/anchors/rancher-ca.pem
          permission: 0444
          content: |-
            -----BEGIN CERTIFICATE-----
            MIIDETCCAfmgAwIBAgIRAK0J3NrgPllXUiGYrA9sTlUwDQYJKoZIhvcNAQELBQAw
            IjEgMB4GA1UEAxMXZWxlbWVudGFsLXNlbGZzaWduZWQtY2EwHhcNMjQxMDA3MTEw
            ODM5WhcNMzUwODAxMTEwODM5WjAiMSAwHgYDVQQDExdlbGVtZW50YWwtc2VsZnNp
            Z25lZC1jYTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAJm2esaQL82b
            rWMpnurmyiutruWvdWUT0Dci2+I7vI1CRs7Gqq93in3+HOoEuaJhS4eZQT9AFyaq
            msijMa3cTYUDhTbOAvPs27E/mSBeQyKd/hJuQ0B8vl47Z1ixOpUHdMOBsZDI0XF5
            yjVTj4nTZXW5n0zZpnmEs4DhLJLJc6icjQLdHDsSj/LeTy8alyTtkOaWcPjFppNI
            6M5a1BWJPhNKGlFpezqfjtJogxbOEAohpN4DUKvqebRWnC+4MjhqUcEW5sXatFTH
            F7MGbVSqQk/f7lzIuke4nvWd0FGPyk/sD31rXT2/2eHkcTJEanRq3bwWQNXQynQ1
            wdqIH1TtfMUCAwEAAaNCMEAwDgYDVR0PAQH/BAQDAgKkMA8GA1UdEwEB/wQFMAMB
            Af8wHQYDVR0OBBYEFIv2OZVFAhh8HzoEwjlf5GivNf6IMA0GCSqGSIb3DQEBCwUA
            A4IBAQAfNUNQKZ02oTo9q+/nbS8kIuhwzSTtNzKflQU5oibpDSAxYlx2gsYqppb/
            w7voj+GiONQR22PrCFh+Kr7aGr/GZh6oXg47dK4Es2dVeE8qdqW3WtZ8oj/OJxmP
            7TqWZdGf7TAxfgNzIpGjWFw/coJ7dcYbDrcZFWG5oQpTbLHK/ECMPWytGVRjrqE6
            baLJ85AVqF9rcCb0giXzvzS6/IpyAe7+Q4WvdzY1uaLQSwkBtpt9OM/O35GmeFUR
            OUkPxQ15e+3tUnDLUDnkTk3xMVRvJehnk/I75auqlUra55KLqfd6SUEbGP3MU9ZI
            12xVJHQTSN8XWh0++9jNG0eSMe75
            -----END CERTIFICATE-----
      runcmd:
        - update-ca-certificates
```

Before the CA cert is replaced on Rancher, the new CA cert can be included on Elemental machines by upgrading them.  
The **new** CA cert can be configured in the [ManagedOSImage](./managedosimage-reference.md#cloudconfig) `cloudConfig`:  

```yaml
apiVersion: elemental.cattle.io/v1beta1
kind: ManagedOSImage
metadata:
  name: ca-cert-upgrade
  namespace: fleet-default
spec:
  # The cloudConfig will be applied after node reboot
  cloudConfig:
    write_files:
      - path: /etc/pki/trust/anchors/rancher-ca-new.pem
        permission: 0444
        content: |-
          -----BEGIN CERTIFICATE-----
          MIIDEDCCAfigAwIBAgIQVgcMnY4HFB5+bZ9yhLaFkTANBgkqhkiG9w0BAQsFADAi
          MSAwHgYDVQQDExdlbGVtZW50YWwtc2VsZnNpZ25lZC1jYTAeFw0yNDEwMDcxMjUx
          MzZaFw0zNTA4MDExMjUxMzZaMCIxIDAeBgNVBAMTF2VsZW1lbnRhbC1zZWxmc2ln
          bmVkLWNhMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvwokj48hQFF3
          +v/ObqPOOyYszL/Nyv8/BomPgBia/GXGe8mkQHEWUXFS4P6KdMGQQU3X6Pm071qG
          QEWEIy95szy1H/q1DgZQCM5fjYPcfFJMopQ28vJEk58/9PePr/GZRWAeAhmMKZeg
          HP/wpuUMEdEh7vGYjKjVuIJiFgT2lVDKqrtRIon+L1iIP3IRmVa49UzmdW2wM79W
          a1nv52+EZaw3UDSLPonvs29AZG8M+NuENlefHWEwYVpDEwF9lXinfL3wMw36gIo4
          X4LmStP9WU4mvglrR8Zwj1M9COMrYbBYQ86jUGM0L0eNG52Uflsn+0ttLRhgkpba
          wAl8jAZWdQIDAQABo0IwQDAOBgNVHQ8BAf8EBAMCAqQwDwYDVR0TAQH/BAUwAwEB
          /zAdBgNVHQ4EFgQUmOGv0AwumUlwQDdULL2dLik/6FUwDQYJKoZIhvcNAQELBQAD
          ggEBACMmDLbKOgz5Zo1pSLTYc08Nb5sRTK/bW24IZ67cfdPstvTQBDAH5+obAjus
          N2Linl/IAsN8K2cnoBq1gM3sST+YDVOBdItZXwe8jybk3IoJPdzE63l//ReTyTSg
          OamwUR6qHcLZ9XNwS4z8WYNy3mDLO6dgq7udb2DHm/0mvyi3Q0oRvsrI+9JCCrgz
          YTFWEWhbpfUzH+dheISMYJx3l/iIFJajaASWKtGBMnp9G+RC2HhDcDwBnW/4JT1h
          wqvat7kdRIxcWHtW482JKRyfa58QidqA7nIBblZJuWqpo4etAVZTCV/caFKbn/Ek
          FrT88MNiy5xsimgQSdt9vptOvJc=
          -----END CERTIFICATE-----
    runcmd:
      - update-ca-certificates
  osImage: "registry.suse.com/suse/sl-micro/6.0/baremetal-os-container:2.1.2-3.59"
  clusterTargets:
    - clusterName: volcano
  upgradeContainer:
    envs:
      # Use FORCE to force an upgrade. 
      # This is convenient when the `osImage` is the same, and only the `cloudConfig` changed.
      - name: FORCE
        value: "false"
```


---

## Article: channels.md

---
sidebar_label: Channels
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/channels"/>
</head>

import MangedOSVersionChannelJson from "!!raw-loader!@site/examples/upgrade/managed-os-version-channel-json.yaml"
import ManagedOSVersionChannelCustom from "!!raw-loader!@site/examples/upgrade/managed-os-version-channel-custom.yaml"
import Versions from "../examples/upgrade/versions.raw!=!raw-loader!@site/examples/upgrade/versions.json"

## Channels

The <Vars name="elemental_operator_name"/> allows subscription to one or more [ManagedOSVersionChannels](./managedosversionchannel-reference.md), to automatically populate a list of [ManagedOSVersions](./managedosversion-reference.md) ready to be consumed to build new ISOs using a [SeedImage](./seedimage-reference.md), or to upgrade existing Elemental nodes to new OS versions using the [ManagedOSImage](./managedosimage-reference.md).  

A channel is normally distributed as an OCI container image, but it is also possible to reference the URI of a JSON file directly containing a list of `ManagedOSVersion`. Note that the best practice is to distribute channels using images, so that distribution is consistent with all other images needed by the <Vars name="elemental_operator_name"/>. This can be beneficial for example when deploying in an Airgapped environment.

<Tabs>
<TabItem value="jsonSyncer" label="Json syncer">

This syncer will fetch a json from url and parse it into valid `ManagedOSVersion` resources.

<CodeBlock language="yaml" title="managed-os-version-channel-json.yaml" showLineNumbers>{MangedOSVersionChannelJson}</CodeBlock>

</TabItem>
<TabItem value="customSyncer" label="Custom syncer">

A custom syncer allows more flexibility on how to gather `ManagedOSVersion` by allowing custom commands with custom images.

This type of syncer allows to run a given command with arguments and env vars in a custom image and output a json file to `/data/output`.
The generated data is then automounted by the syncer and then parsed so it can gather create the proper versions.

Elemental project provides channels to list all `ManagedOSVersions` released as a custom syncer.  
See the channel resource definition below:

<CodeBlock language="yaml" title="managed-os-version-channel.yaml" showLineNumbers>{ManagedOSVersionChannelCustom}</CodeBlock>

</TabItem>
</Tabs>

### Available Channels

Elemental maintains a list of channels that can be used out of the box.  

| Base OS          | BaseOS Version | Flavor     | Channel URI                                                        |
|------------------|----------------|------------|--------------------------------------------------------------------|
| SL Micro         | 6.1            | Base       | registry.suse.com/rancher/elemental-channel/sl-micro:6.1-base      |
| SL Micro         | 6.1            | Bare-metal | registry.suse.com/rancher/elemental-channel/sl-micro:6.1-baremetal |
| SL Micro         | 6.1            | KVM        | registry.suse.com/rancher/elemental-channel/sl-micro:6.1-kvm       |
| SL Micro         | 6.1            | RT         | registry.suse.com/rancher/elemental-channel/sl-micro:6.1-rt        |

#### Finding Elemental channels

Using crane we can find the channels maintained:

```
$ crane ls -O registry.suse.com/rancher/elemental-channel/sl-micro
6.0-baremetal
6.0-base
6.0-kvm
6.0-rt
6.1-baremetal
6.1-base
6.1-kvm
6.1-rt
<snip>
```

### Flavors

Elemental distributes different OS flavors that can better fit specific use cases.

| Flavor     | Description                                                                       | Reference                                                                                                |
|------------|-----------------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------|
| Base       | A minimal image that can be used as base to build custom images.                  | [Source](https://github.com/rancher/elemental/blob/v2.1.x/.obs/dockerfile/micro-base-os/Dockerfile)      |
| Bare-metal | Contains bare-metal and usability packages. Can be used for any generic workload. | [Source](https://github.com/rancher/elemental/blob/v2.1.x/.obs/dockerfile/micro-baremetal-os/Dockerfile) |
| KVM        | Ready to be used with KVM. Contains QEMU Guest agent by default.                  | [Source](https://github.com/rancher/elemental/blob/v2.1.x/.obs/dockerfile/micro-kvm-os/Dockerfile)       |
| RT         | Like bare-metal images, but includes a Real-Time kernel.                          | [Source](https://github.com/rancher/elemental/blob/v2.1.x/.obs/dockerfile/micro-rt-os/Dockerfile)        |

### Channels lifecycle and best practices

Once a new `ManagedOSVersionChannel` is created, the <Vars name="elemental_operator_name"/> will periodically sync the channel provided JSON list, and convert it to new `ManagedOSVersions`.  
All synced `ManagedOSVersions` will be owned by the `ManagedOSVersionChannel`. Deleting the `ManagedOSVersionChannel` will lead to the deletion of all `ManagedOSVersions` on cascade.  

Note that the `ManagedOSVersionChannel` supports automatic clean up of no longer in sync `ManagedOSVersions`, when the `ManagedOSVersionChannel.spec.deleteNoLongerInSyncVersions` option is enabled.  

When a `ManagedOSVersion` is scheduled for deletion, a finalizer will make sure that there is no active reference on any `ManagedOSImage`.  

If a `ManagedOSVersion` can not be deleted, you can find out by which resources it is referenced:  

```bash
kubectl -n fleet-default get managedosimages -l elemental.cattle.io/managed-os-version-name=my-deleted-os-version
```

When using multiple channels it's important to keep a proper naming strategy to always have a quick, human readable reference on the owned `ManagedOSVersions`.  
It is recommended to name any channel as: `{BaseOS}-{BaseOSVersion}-{Flavor}`.  

This should allow the user to use the `ManagedOSVersion` name as the specific Elemental build version of the image, while keeping a reference on the Base OS and Base OS version from the parent channel.  
On the Rancher UI this will look something like the following image:  

![Channel naming](images/channel-naming.png)

### Making your own Channels

The only requirement to make your own custom syncer is to make it output a JSON file to `/data/output` and keep the correct JSON structure.  

The file is a JSON array containing ISO and Container entries.  
Each entry in the array is mapped 1:1 with a [ManagedOSVersion](./managedosversion-reference.md) object.  

`"type": "iso"` entries must contain a bootable Elemental ISO and are used by [SeedImages](./seedimage-reference.md), while `"type": "container"` entries are used by [ManagedOSImage](./managedosimage-reference.md) for Elemental upgrades.  

If in doubt, the [elemental-channels](https://github.com/rancher-sandbox/elemental-channels) project can be used as a reference implementation on how to build and maintain your own channels.

When creating new entries, be mindful of the naming strategy you choose, in order to avoid collisions with other channels, since they may end up syncing different `ManagedOSVersion` with the same name.  
A best practice is to use the convention: `{Flavor}-{Version}-{Type}`  
A sample of the JSON format is as follows:  

<CodeBlock language="json" title="versions.json" showLineNumbers>{Versions}</CodeBlock>


---

## Article: cloud-config-reference.md

---
sidebar_label: Cloud-config reference
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/cloud-config-reference"/>
</head>

# Cloud-config Reference

Node OS images build using the [Elemental Toolkit](https://github.com/rancher/elemental-toolkit) are expected
to be initialized and configured by using [yip](https://github.com/rancher/yip). Yip is a small utility to
apply a set of actions and configurations to the system described with yaml files. Yip is integrated and consumed
as a library within the elemental client binary (see `elemental run-stage --help`). Yip groups the configurations
and actions to apply in arbitrary _stages_, for instance the `elemental run-stage network` command call would only
apply defined actions and configuration for the stage named `network`. Note from Yip perspective stages can be run at
any time as they are simply grouping a set of actions under an arbitrary name.

Elemental Toolkit integrates five predefined stages into the OS boot process.

1. **`pre-rootfs`**: this stage runs on early boot inside the init ram disk, just *before* mounting the root device (typically at `/sysroot`).
   This stage can be used to define first-boot steps that are required before mounting the root device. A good example is expanding the
   root device partition. Executed as part of the `initrd-root-device.target`.

2. **`rootfs`**: this stage runs on early boot inside the init ram disk, just *after* mounting the root device (typically at `/sysroot`).
   This stage can be used to define first-boot steps like creating new partitions. Ephemeral and  persistent paths are typically defined
   at this stage. Executed as part of the `initrd-root-fs.target`.

3. **`initramfs`**: this stage runs inside the init ram disk too, but on a later stage just before switching root. This stage runs in a chrooted
   environment to the actual root device. This stage is handy to set some system parameters that might be relevant to systemd
   after switching root. For instance, additional systemd unit files could be added here before the systemd from the actual root is executed.
   Executed as part of the `initrd.target`.

4. **`fs`**: this stage is the first one executed after switching root and it is executed as part of the `sysinit.target` which runs once all
   all local filesystems and mountpoints defined in fstab are mounted and ready.

5. **`network`**: this stage is executed as part of the `multi-user.target` and after the `network-online.target` is reached. This stage can be used
   to run actions and configurations that require network connectivity. For instance this stage is used to run the very first node registration and
   and installation from a live ISO.

6. **`boot`**: this stage is executed as part of the `multi-user.target` and before the `getty.target`. This is the default stage to run cloud-config
   data provided using the supported cloud-init syntax. See [cloud-init compatibility](cloud-config-reference.md#compatibility-with-cloud-init-format) section.

By default, `elemental` reads the yaml configuration files from the following paths in order: `/system/oem`, `/oem` and `/usr/local/cloud-config`.

In Elemental Operator, all kubernetes resources including a `cloud-config` field can be expressed in either [yip](#configuration-syntax) or
[cloud-init compatible](#compatibility-with-cloud-init-format) syntax. This includes resources such as `MachineRegistration`, `SeedImage`, and `ManagedOSImage`.

:::note
In contrast to similar projects such as _Cloud Init_, Yip does not keep records or caches of executed stages and steps,
all stages and its associated configuration is executed at every boot.
:::


## Elemental client cloud-config hooks

In addition to the defined cloud-config stages at boot described in the previous section the Elemental client also
honors some specific stages, referenced as hooks, to customize the behavior of these subcommands: `install`, 
`upgrade`, `reset` and `build-disk`. Each of these subcommands has it's own set of four different cloud-config stages executed at
analog phases of the specific subcommand execution.

Hooks are essetially a way to provide permanent changes to system that can't be easily expressed as part of an OCI container or that
are not easily achievable with the `elemental` client configuration options. A good example could be handling the firmware in EFI
partition for Raspberry Pi devices.

:::warning
Note most hooks are executed in the host environment with privileges, so they are potentially destructive operations. In most
cases regular cloud-config operations at boot time are sufficient to setup the system. Also to include additional software in an image
the preferred option is to build a [derivative image](custom-images.md) and not abuse of hooks to install additional software.
:::


### Hook stages

* Before stages: `before-install`, `before-upgrade`, `before-reset`, `before-disk`
  These stages are executed once the working directories and environment are prepared but before starting the actual action. In
  `install`, `upgrade` and `reset` steps this happens once all the associated partitions are created and mounted, but before stating
  the deployment of any image.

* After chrooted stages: `after-install-chroot`, `after-upgrade-chroot`, `after-reset-chroot`, `after-disk-chroot`
  These stages are executed after deploying the target system into the working area into a chroot environment rooted to the
  actual deployed image. Since this happens in a chroot env the elemental client analyses the hooks present in the deployed image, not in the host.
  Only `/oem` is shared with the host if available.

* After stages: `after-install`, `after-upgrade`, `after-reset`, `after-disk`
  These stages are executed after deploying the target system into the working area from the host environment. At this stages all
  partitions are still mounted and available in RW mode. This particular set of stages analyses the hooks present in the host
  and into the equivalent set of paths chrooted to the deployed image. The hook however is not executed in a chroot environment.
  This is helpful to provide `after-*` hooks within the deployed image.

* Post stages: `post-install`, `post-upgrade`, `post-reset`, `post-disk`
  These stages are executed at end before exiting the command and running a cleanup process. At this stage the image is already deployed
  and locked in a read-only subvolume or filesystem. Partitions are still mounted at this stage.

:::note
Note installation hooks are not applied as part of the [MachineRegistration.config.cloud-config](machineregistration-reference.md#configcloud-config).
In order to provide installation hooks they can be included as part of the [SeedImage.cloud-config](seedimage-reference.md#seedimagespec-reference),
as they need to be present in the installation media.
:::

## Configuration syntax

Yip has its own syntax, it essentially requires to define _stages_ and a list of steps for each stage. Steps are executed in
order and each step can be a combination different action types (e.g run commands, create files, set hostname, etc.).

Consider the following example:

```yaml
stages:
  initramfs:
  - name: "Setup users"
    ensure_entities:
    - path: /etc/shadow
      entity: |
          kind: "shadow"
          username: "root"
          password: "root"
  boot:
  - files:
    - path: /tmp/script.sh
      content: |
        #!/bin/sh
        echo "test"
      permissions: 0777
      owner: 1000
      group: 100
  - commands:
    - /tmp/script.sh
```

In the above example there are two stages: `initramfs` and `boot`.
The `initramfs` stage initializes a sample user.
The `boot` stage includes two steps, one to create an executable script file and a second one
that actually runs the script.

Yip also supports `*.before` and `*.after` suffix modifiers to any given stage. For instance, running the `network` stage
results into running first `network.before` stages found in config files and then `network` and finally `network.after`.

See the full reference of applicable keys in steps documented in
[yip project](https://github.com/rancher/yip?tab=readme-ov-file#configuration-reference) itself.

Below is an example of the above configuration embedded in a MachineRegistration resource.

<details>
  <summary>MachineRegistration example</summary>

```yaml showLineNumbers
apiVersion: elemental.cattle.io/v1beta1
kind: MachineRegistration
metadata:
  name: my-nodes
  namespace: fleet-default
spec:
  config:
    cloud-config:
      name: "A registration driven config"
      stages:
        after-install-chroot:
        - name: "Set serial console"
          commands:
          - grub2-editenv /oem/grubenv set extra_cmdline="console=ttyS0"
        initramfs:
        - name: "Setup users"
          ensure_entities:
          - path: /etc/shadow
            entity: |
                kind: "shadow"
                username: "root"
                password: "root"
        boot:
        - files:
          - path: /tmp/script.sh
            content: |
              #!/bin/sh
              echo "test"
            permissions: 0777
            owner: 1000
            group: 100
        - commands:
          - /tmp/script.sh
    elemental:
      install:
        reboot: true
        device: /dev/sda
        debug: true
  machineName: my-machine
  machineInventoryLabels:
    element: fire
```

</details>

## Compatibility with Cloud Init format

A subset of the official [cloud-config spec](http://cloudinit.readthedocs.org/en/latest/topics/format.html#cloud-config-data) is implemented by yip.
More specific the supported cloud-init keys are: `users`, `ssh_authorized_keys`, `runcmd`, `hostname` and `write_files` are implemented.

If a yaml file starts with `#cloud-config` it is parsed as a standard cloud-init, associated it to the yip `boot` stage.
For example:

```yaml
#cloud-config

# Note groups are delivered as list, not as comma separated values
users:
- name: "bar"
  passwd: "foo"
  groups:
  - "users"
  homedir: "/home/foo"
  shell: "/bin/bash"
  ssh_authorized_keys:
  - faaapploo

# Assigns these keys to the first user in users or root if there
# is none
ssh_authorized_keys:
- asddadfafefa

# Run these commands once the system has fully booted
# Each command is expressed as a sinlge string, no nested lists
runcmd:
- echo hello world

# Hostname
hostname: myserver

# Write arbitrary files
write_files:
- encoding: b64
  content: CiMgVGhpcyBmaWxlIGNvbnRyb2xzIHRoZSBzdGF0ZSBvZiBTRUxpbnV4
  path: /foo/bar
  permissions: "0644"
  owner: "bar"
```

Below is an example of the above configuration embedded in a MachineRegistration resource.

<details>
  <summary>MachineRegistration example</summary>

```yaml showLineNumbers
apiVersion: elemental.cattle.io/v1beta1
kind: MachineRegistration
metadata:
  name: my-nodes
  namespace: fleet-default
spec:
  config:
    cloud-config:
      users:
      - name: "bar"
        passwd: "foo"
        groups:
        - "users"
        homedir: "/home/foo"
        shell: "/bin/bash"
        ssh_authorized_keys:
        - faaapploo
      ssh_authorized_keys:
      - asdd
      runcmd:
      - foo
      write_files:
      - encoding: b64
        content: CiMgVGhpcyBmaWxlIGNvbnRyb2xzIHRoZSBzdGF0ZSBvZiBTRUxpbnV4
        path: /foo/bar
        permissions: "0644"
        owner: "bar"
    elemental:
      install:
        reboot: true
        device: /dev/sda
        debug: true
  machineName: my-machine
  machineInventoryLabels:
    element: fire
```

</details>


---

## Article: cluster-reference.md

---
sidebar_label: Cluster reference
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/cluster-reference"/>
</head>

import Machinepools from "!!raw-loader!@site/examples/clusters/clusters-several-machinepools.yml"

# Cluster reference

A `Cluster` definition includes a `kubernetesVersion` and a list of `machinePools` to deploy the cluster to.

For how to select a `kubernetesVersion` please check our [Kubernetes Versions](kubernetesversions.md) page.

A `machinePool` is a bundle of configuration with a `ObjectReference` so the cluster is deployed to those `machinePools`
with the proper roles (etcd, control-plane, worker) with a quantity (how many nodes to deploy from this pool) and some extra configurations (rolling update config, max unhealthy nodes, etc...).

<details>
  <summary>Example</summary>

  ```yaml showLineNumbers
  kind: Cluster
  apiVersion: provisioning.cattle.io/v1
  metadata:
    name: ...
    namespace: ...
  spec:
    rkeConfig:
      machinePools:
        - name: ...
          controlPlaneRole: ...
          etcdRole: ...
          workerRole: ...
          quantity: ...
          machineConfigRef:
            apiVersion: elemental.cattle.io/v1beta1
            kind: MachineInventorySelectorTemplate
            name: ...
        - name: ...
          controlPlaneRole: ...
          etcdRole: ...
          workerRole: ...
          quantity: ...
          machineConfigRef:
            apiVersion: elemental.cattle.io/v1beta1
            kind: MachineInventorySelectorTemplate
            name: ...
  ```

</details>

It's also possible to disable cluster components via the `Cluster` object in `spec.rkeConfig.machineGlobalConfig`, for example:

<details>
  <summary>Service Disabling Example</summary>

  ```yaml showLineNumbers
  kind: Cluster
  apiVersion: provisioning.cattle.io/v1
  metadata:
    name: ...
    namespace: ...
  spec:
    rkeConfig:
      machinePools:
        - name: ...
          controlPlaneRole: ...
          etcdRole: ...
          workerRole: ...
          quantity: ...
          machineConfigRef:
            apiVersion: elemental.cattle.io/v1beta1
            kind: MachineInventorySelectorTemplate
            name: ...
      machineGlobalConfig:
        disable:
          - servicelb
          - ...
  ```

</details>

## rkeConfig.machinePools

A list of `machinePools`. A minimum of 1 `machinePools` is required for the cluster to be deployed to.

## machinePools Spec Reference

| Key                  | Type   | Default value   | Description                                                          |
|----------------------|--------|-----------------|----------------------------------------------------------------------|
| controlPlaneRole     | bool   | false           | Set machines in this pool as control-plane                           |
| etcdRole             | bool   | false           | Set machines in this pool as etcd                                    |
| workerRole           | bool   | false           | Set machines in this pool as worker                                  |
| name                 | string | nil             | Name for this pool                                                   |
| quantity             | int    | nil             | Number of machines to deploy from this pool                          |
| unhealthyNodeTimeout | int    | nil             | Timeout for unhealthy node health checks                             |
| machineConfigRef     | int    | ObjectReference | Reference to an object used to know what nodes are part of this pool |

A minimum of `quantity` set to one is required for this pool to be used.
Basically translates to how many nodes from this pool are going to be setup for this cluster.

<details>
  <summary>Example</summary>

  ```yaml showLineNumbers
  kind: Cluster
  apiVersion: provisioning.cattle.io/v1
  metadata:
    name: cluster-example
    namespace: example-default
  spec:
    rkeConfig:
      machinePools:
        - name: examplePool 
          controlPlaneRole: true
          etcdRole: true
          workerRole: false
          quantity: 3
          unhealthyNodeTimeout: 0s
          machineConfigRef:
            apiVersion: elemental.cattle.io/v1beta1
            kind: MachineInventorySelectorTemplate
            name: exampleSelector
  ```

</details>

## machineConfigRef Spec Reference

A `machineConfigRef` is a generic k8s `ObjectReference` which usually contain a
`kind` `name` and `apiVersion` to point to a different object.

In Elemental, we set this to a [MachineInventorySelectorTemplate](machineinventoryselectortemplate-reference.md).
This allows us to point to more than one object by using the selector.

### Example

The example below creates a cluster that uses 2 different `machinePool`'s to set different nodes to control-plane and workers nodes,
based on 2 different `MachineInventorySelectorTemplate` that select their nodes based on a `MachineInventory` label (location):

:::warning warning
The labels for the example are manual set labels, they are not set by Elemental automatically..

For automatic labels generated by Elemental please check the [SMBIOS](smbios.md) page.
:::

<CodeBlock language="yaml" title="Example of a cluster with more than one machinePool" showLineNumbers>{Machinepools}</CodeBlock>


---

## Article: custom-certificate.md

---
sidebar_label: Add a custom certificate
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/custom-certificate"/>
</head>


### How to add a custom certificate

Prerequisite: A certificate in `.pem` format

Goal: Make a custom certificate available system-wide

:::note This is for certificates used by system-level services.
Kubernetes workloads should bring their certificates within the
container image instead.
:::

In order to install a custom certificate we need to

* copy the `.pem` file to `/etc/pki/trust/anchors/`
* run `update-ca-certificates`

The respective `cloud-config` snippet looks like this:

```yaml
write_files:
  - path: /etc/pki/trust/anchors/my-custom-certificate.pem
    permission: 0444
    content: |-
      -----BEGIN CERTIFICATE-----
      ...
      -----END CERTIFICATE-----
runcmd:
  - update-ca-certificates
```

(actual certificate content omitted for brevity reasons)


---

## Article: custom-images.md

---
sidebar_label: Build Custom OS Images
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/custom-images"/>
</head>

# How to build and use custom OS images

## Remastering an OS image with a custom Dockerfile

Since OS images provided by Elemental are container images, they can also be used as a base image
in a Dockerfile in order to create a new container image.

The Elemental project publishes several flavors for images:
* baremetal: An image containing firmware and drivers suitable for baremetal deployment.
* rt: Based on the baremetal image, but contains the real-time kernel.
* kvm: A slimmer image suitable for VMs.
* base: The base system needed for Elemental used by the other flavors.

Imagine some additional packages from an extra repository is required, the following example
showcases how this could be added:

```docker showLineNumbers
# The version of Elemental to modify.
FROM registry.suse.com/suse/sl-micro/6.1/baremetal-os-container:latest

# Custom commands
RUN rpm --import <repo-signing-key-url> && \
    zypper addrepo --refresh <repo_url> extra_repo && \
    zypper install -y <extra_package>

# IMPORTANT: /etc/os-release is used for versioning/upgrade. The
# values here should reflect the tag of the image currently being built
ARG IMAGE_REPO=norepo
ARG IMAGE_TAG=latest
RUN \
    sed -i -e "s|^IMAGE_REPO=.*|IMAGE_REPO=\"${IMAGE_REPO}\"|g" /etc/os-release && \
    sed -i -e "s|^IMAGE_TAG=.*|IMAGE_TAG=\"${IMAGE_TAG}\"|g" /etc/os-release && \
    sed -i -e "s|^IMAGE=.*|IMAGE=\"${IMAGE_REPO}:${IMAGE_TAG}\"|g" /etc/os-release

# IMPORTANT: it is good practice to recreate the initrd and re-apply `elemental-init`
# command that was used in the base image. This ensures that any eventual change that should
# be synced in initrd included binaries is also applied there and consistent.
RUN elemental init --force elemental-rootfs,grub-config,dracut-config,cloud-config-essentials,elemental-setup
```

Where `latest` is the base version we want to customize.

And then the following commands

```bash showLineNumbers
docker build --build-arg IMAGE_REPO=myrepo/custom-build \
             --build-arg IMAGE_TAG=v1.1.1 \
             -t myrepo/custom-build:v1.1.1 .
docker push myrepo/custom-build:v1.1.1
```

The new customized OS is available as the Docker image `myrepo/custom-build:v1.1.1` and it can
be run and verified using docker with

```bash showLineNumbers
docker run -it myrepo/custom-build:v1.1.1 bash
```

## Create a custom bootable installation ISO

Elemental leverages container images to build its root filesystems; therefore, it is possible
to use it in a multi-stage environment to create custom bootable media that bundles a custom container image.

```docker showLineNumbers
FROM registry.suse.com/suse/sl-micro/6.1/baremetal-os-container:latest AS os

# Check the previous section on building custom images

# The released OS already includes the toolchain for building ISOs
FROM registry.suse.com/suse/sl-micro/6.1/baremetal-os-container:latest AS builder

ARG TARGETARCH
WORKDIR /iso
COPY --from=os / rootfs

# work around buildah issue: https://github.com/containers/buildah/issues/4242
RUN rm -f rootfs/etc/resolv.conf

RUN elemental build-iso \
        dir:rootfs \
        --bootloader-in-rootfs \
        --squash-no-compression \
        -o /output -n "elemental-${TARGETARCH}"

FROM busybox
COPY --from=builder /output /elemental-iso

ENTRYPOINT ["busybox", "sh", "-c"]
```

Build it with regular `docker build` command:

```bash showLineNumbers
docker build -t myrepo/custom-build:v1.1.1 \
              --build-arg IMAGE_REPO=myrepo/custom-build-iso \
              --build-arg IMAGE_TAG=v1.1.1 \
              .
```

The resulting container image is actually a container image including the ISO,
this container image can be pushed to an OCI registry too. The ISO image can be
extracted from the container to the current folder by executing the container as:

```bash showLineNumbers
docker run --rm -v $(pwd):/host mytest-image "busybox cp /elemental-iso/*.iso /host"
```

The new customized installation media can be found in `elemental-<arch>.iso`.

The above container run is equivalent to what *elemental-operator* does to extract
the ISO from a container to build a new one including the registration URL,
hence this is also a good check mark to verify the container can be pushed to a
registry and used by the *elemental-operator* as a `baseImage` for a
[SeedImage](seedimage-reference) resource.

## List custom images as a ManagedOSVersion resource

In Elemental listing OS container images and ISO container images as ManagedOSVersion
resources is not mandatory but handy. Specially from a UI perspective this makes
the custom images visible and easy to use from the Elemental UI extension.

Continuing the example from the previous section a custom OS container referenced as
`myrepo/custom-build:v1.1.1` was built and eventually pushed to a registry. Then this
image is ready to be added as a ManagedOSVersion resource with:

```yaml showLineNumbers
apiVersion: elemental.cattle.io/v1beta1
kind: ManagedOSVersion
metadata:
  name: v1.1.1-custom-build
  namespace: fleet-default
spec:
  metadata:
    displayName: Custom build image
    upgradeImage: myrepo/custom-build:v1.1.1
  type: container
  version: v1.1.1
```

Note the `type: container` states this is a container OS. This makes the image `myrepo/custom-build:v1.1.1`
eligible for OS upgrades from the UI.

Finally, the custom container for the ISO `myrepo/custom-build-iso:v1.1.1` can also be included
as a ManagedOSVersion resource with:

```yaml showLineNumbers
apiVersion: elemental.cattle.io/v1beta1
kind: ManagedOSVersion
metadata:
  name: v1.1.1-custom-build-iso
  namespace: fleet-default
spec:
  metadata:
    displayName: Custom build ISO image
    uri: myrepo/custom-build-iso:v1.1.1
  type: iso
  version: v1.1.1
```

Note the  `type: iso` states this is an ISO. This makes the image `myrepo/custom-build-iso:v1.1.1`
eligible for SeedImages generation from UI.

## Custom partition size

When building custom images, it's important to take in account disk partition sizes, to ensure the image and the upgrade snapshots can fit correctly over time.  
A partitions configuration can be included in your custom image, or alternatively it can be conveniently applied to the [SeedImage](./seedimage-reference.md) used to generate the install media.  
Note that all `size` values are expressed in megabytes, and a value of `0` will take the rest of the disk. This is the default behavior of the `persistent` partition if no `size` has been defined for it. For more information, see the full [configuration sample](https://github.com/rancher/elemental-toolkit/blob/main/config.yaml.example).  

```yaml
apiVersion: elemental.cattle.io/v1beta1
kind: SeedImage
metadata:
  name: custom-partitions-iso
  namespace: fleet-default
spec:
  cloud-config:
    write_files:
    - path: /etc/elemental/config.d/partitions.yaml
      content: |
        install:
          partitions:
            recovery:
              size: 8192
            state:
              size: 16384
    - path: /etc/elemental/config.d/snapshotter.yaml
      content: |
        snapshotter:
          max-snaps: 2
  baseImage: myrepo/custom-build-iso:v1.1.1
  registrationRef:
    name: my-machine-registration
    namespace: fleet-default
```

The `state` partition will hold all system snapshots. Therefore when sizing this partition, the following formula can be considered: `$image_size * ($max_number_of_snapshots + 1 + 1)`.  
The `$max_number_of_snapshots` can be similarly configured with a custom configuration file as shown in the sample above.  
Note that by default it's `4` for the `btrfs` snapshotter type, and `2` for the `loopdevice` type.  
You can configure the snapshotter type in use editing the [MachineRegistration](./machineregistration-reference.md#configelementalinstallsnapshotter).  

Since the state partition is also used for the <Vars name="elemental_toolkit_name" link="elemental_toolkit_url"/> work directory, it's best to leave an additional `$image_size` worth of free space, so that the image can be unpacked correctly for example when running upgrades.  

Lastly, an extra `$image_size` free space can be used as a safe margin to keep. This is especially important when using the `loopdevice` snapshotter type, in case newer images will grow in size from the originally installed one.  
On the contrary, the `btrfs` snapshotter can be used instead to save space on the `state` partition, or to use the same space to keep more snapshots.  

## Finding Elemental base images

Using crane we can find the following SL-Micro images suitable for extending:

```
$ crane catalog registry.suse.com | grep -i "suse/sl-micro"
suse/sl-micro/6.0/baremetal-iso-image
suse/sl-micro/6.0/baremetal-os-container
suse/sl-micro/6.0/base-iso-image
suse/sl-micro/6.0/base-os-container
suse/sl-micro/6.0/kvm-iso-image
suse/sl-micro/6.0/kvm-os-container
suse/sl-micro/6.0/rt-iso-image
suse/sl-micro/6.0/rt-os-container
suse/sl-micro/6.1/baremetal-iso-image
suse/sl-micro/6.1/baremetal-os-container
suse/sl-micro/6.1/base-iso-image
suse/sl-micro/6.1/base-os-container
suse/sl-micro/6.1/kvm-iso-image
suse/sl-micro/6.1/kvm-os-container
suse/sl-micro/6.1/rt-iso-image
suse/sl-micro/6.1/rt-os-container
```

The images with `-iso-image` suffix contains a pre-built ISO image and a busybox system to be able to copy the contents to a volume.
Images with a `-os-container` suffix contains a root filesystem that can be used as the base for custom images.


---

## Article: custom-install.md

---
sidebar_label: Customize Elemental Installation
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/custom-install"/>
</head>

# Customize Elemental Installation

Elemental installed OS images can be customized in different ways.

One option is to remaster container OS images by simply using a docker build.
SL Micro images are regular container images, so it is absolutely possible to create
a new image using a Dockerfile based on SLE Micro. See [Build Custom OS Images](/custom-images.md)
section for further details, this is the preferred option.

Alternatively, it is also possible to provide additional resources and configuration
within the installation media so that during installation, or eventually at boot time,
additional binaries such as drivers or extra configuration files can be included.

This section focuses on how to customize the installation process from a given OS image. 

## Customization options

Elemental installation can be customized in three different non-exclusive ways. First,
including additional cloud-config files into the installed system, second, by including
additional cloud-config file into the installation media, and finally, by adding a custom
Elemental client configuration file (`/etc/elemental/config.yaml`).

1. The additional cloud-config files included into the installed system are useful to run
   custom operations at boot time. See the [Cloud Config Reference](cloud-config-reference.md).

2. The additional cloud-config files into the installation media are useful to run custom operations
   at install time or to customize the installation environment to match certain specific needs.

3. A custom Elemental client configuration file is, by default, located at `/etc/elemental/config.yaml`
   and it can also be split into multiple yaml files under `/etc/elemental/config.d` directory.
   See the configuration file
   [reference](https://rancher.github.io/elemental-toolkit/docs/customizing/general_configuration/).

A common pattern is to combine the three ways described above depending on the specific needs.

### Adding additional cloud-config files within the installed OS

In order to include additional cloud-config files during the installation they need
to be added to the installation data into the MachineRegistration resource. The most simple
way to achieve it is by adding the desired cloud-config directly as part of the MachineRegistration [devoted
section](machineregistration-reference.md#configcloud-config). See the example below


```yaml showLineNumbers
apiVersion: elemental.cattle.io/v1beta1
kind: MachineRegistration
metadata:
  name: my-nodes
  namespace: fleet-default
spec:
  ...
  config:
    ...
    cloud-config:
      stages:
        boot:
        - name: "Adding 'admin' user"
          users:
            admin:
              passwd: mysecretpasswd
```

Alternatively, cloud-config files path can also be referenced explicitly so the configuration is not
strictly needed to live within the MachineRegistration resource itself. The `config-urls` section of the
MachineRegistration is used for this exact purpose. See
[MachineRegistration reference](/machineregistration-reference) page.

`config-urls` is a list of string literals where each item is an HTTP URL or a local path pointing to a
cloud-config file. The local path is evaluated during
the installation, hence it must exists within the installation media, commonly an ISO image.

By default, Elemental live systems mount the ISO root at `/run/initramfs/live` which is also the default path set for `config-url` in `MachineRegistrations`:
See the example below:

```yaml showLineNumbers
apiVersion: elemental.cattle.io/v1beta1
kind: MachineRegistration
metadata:
  name: my-nodes
  namespace: fleet-default
spec:
  ...
  config:
    ...
    elemental:
      ...
      install:
        ...
        config-urls:
        - "/run/initramfs/live/oem/custom_config.yaml"
```
Elemental live ISOs, when booted, have the ISO root mounted at `/run/initramfs/live`.
According to that, the ISO for the example above is expected to include the `/oem/custom_config.yaml` file.

:::note
`/run/initramfs/live` is a readonly mountpoint and it's not an appropriate path for dynamically generated content at ISO boot.
:::

### Adding additional cloud-config files within the installation media

Adding additional cloud-config files within the installation media might be required to configure the
installation environment (e.g. setting the network connectivity to register) or to provide some
[installation hooks](cloud-config-reference.md#elemental-client-cloud-config-hooks) to run custom
logic during the installation itself.

In Elemental the [SeedImage](seedimage-reference.md) resources are the responsibles of handling the installation
media. So the simplest place to include additional cloud-config data is within the installation media is by
adding it to the [`cloud-config` section](seedimage-reference.md#seedimagespec-reference). By doing so the
given cloud-config will be already evaluated form the very first boot of the installation media and also
honored during the installation phase in case some install hook is provided. Consider the following example:

```yaml showLineNumbers
apiVersion: elemental.cattle.io/v1beta1
kind: SeedImage
metadata:
  name: custom-seed
  namespace: fleet-default
spec:
  ...
  cloud-config:
    stages:
      post-install:
      - name: "Run custom script after installation"
        commands:
        - |
          echo "This is a custom script"
          echo "For instance, this could be used to handle extra drives for an LVM group"
      boot:
      - name: "Add proxy setup for the installation media"
        files:
        - path: /etc/sysconfig/proxy
          permissions: 0664
          content: |
            PROXY_ENABLED="yes"
            HTTP_PROXY=http://<MY_PROXY>:<MY_PORT>
            HTTPS_PROXY=https://<MY_PROXY>:<MY_PORT>
            NO_PROXY="localhost, 127.0.0.1"
```

### Custom Elemental client configuration file

[Elemental client](https://github.com/rancher/elemental-toolkit/blob/main/docs/elemental.md) `install`, `upgrade` and 
`reset` commands can be configured with a
[custom configuration file](https://rancher.github.io/elemental-toolkit/docs/customizing/general_configuration/) located
by default in `/etc/elemental/config.yaml`. If you have multiple yaml files, you need to add them in the
`/etc/elemental/config.d` directory.

The following example sets an additional extra partition during the installation:

```yaml showLineNumbers
install:
  extra-partitions:
  - size: 10240
    fs: ext4
    label: EXTRA_PARTITION
```

In order to make it available at installation time this could be done my adding the extra file as part
of the SeedImage resource cloud-config as it is described in the previous section of this page. Consider
the example:

```yaml showLineNumbers
apiVersion: elemental.cattle.io/v1beta1
kind: SeedImage
metadata:
  name: custom-seed
  namespace: fleet-default
spec:
  ...
  cloud-config:
    stages:
      boot:
      - name: "Add Elemental client configuration file"
        files:
        - path: /etc/elemental/config.d/extra-partition.yaml
          permissions: 0664
          content: |
            install:
              extra-partitions:
              - size: 10240
                fs: ext4
                label: EXTRA_PARTITION
```


---

## Article: custom-resources.md

---
sidebar_label: Custom Resources
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/custom-resources"/>
</head>

# Custom Resources

The <Vars name="elemental_operator_name" /> allows control of the Elemental Nodes by extending the Kubernetes APIs with a set of _elemental.cattle.io_ [Kubernetes CRDs](https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/):

- [MachineRegistration](machineregistration-reference.md)
- [MachineInventory](machineinventory-reference.md)
- [MachineInventorySelector](machineinventoryselector-reference.md)
- [MachineInventorySelectorTemplate](machineinventoryselectortemplate-reference.md)
- [ManagedOSImage](managedosimage-reference.md)
- [ManagedOSVersion](managedosversion-reference.md)
- [ManagedOSVersionChannel](managedosversionchannel-reference.md)
- [SeedImage](seedimage-reference.md)

---

## Article: elemental-plans.md

---
sidebar_label: Elemental plans
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/elemental-plans"/>
</head>

## Introduction

Elemental uses the [Rancher System Agent](https://github.com/rancher/system-agent), renamed to Elemental System Agent, to initially bootstrap the node with a simple plan.

The plan will apply the following configurations:

- Set some labels for the node
- Set the proper hostname according to the `MachineInventory` value
- Install the default Rancher System Agent from Rancher Server, and install the proper Kubernetes components

The bootstrap service also accepts local plans stored under `/var/lib/elemental/agent/plans`. Any plan written
in there will also be applied during the initial node start after the installation is completed.

:::tip
The local plans run only during the initial Elemental bootstrap **before** Kubernetes is installed on the node.
:::

## Types of Plans

The type of plans that Elemental can use are:

- One time instructions: Only run once
- Periodic instructions: They run periodically
- Files: Creates files
- Probes: http probes

:::tip
Both one time instructions and periodic instructions can run either a direct command or a docker image.
:::

## Adding local plans on Elemental

You can add local plans to Elemental as part of the `MachineRegistration` CRD, in the `cloud-config` section as follows:

```yaml showLineNumbers
apiVersion: elemental.cattle.io/v1beta1
kind: MachineRegistration
metadata:
  name: my-nodes
  namespace: fleet-default
spec:
  config:
    cloud-config:
      users:
        - name: root
          passwd: root
      write_files:
        - path: /var/lib/elemental/agent/plans/mycustomplan.plan
          permissions: "0600"
          content: |
            {"instructions":
                [
                  {
                    "name":"set hostname",
                    "command":"hostnamectl",
                    "args": ["set-hostname", "myHostname"]
                  },
                  {
                    "name":"stop sshd service",
                    "command":"systemctl",
                    "args": ["stop", "sshd"]
                  }
                ]
            }
    elemental:
      install:
        reboot: true
        device: /dev/sda
        debug: true
  machineName: my-machine
  machineInventoryLabels:
    element: fire
```

## Plan examples

The following plans are provided as a quick reference and are not guaranteed to work in your environment. To learn more about plans please check [Rancher System Agent](https://github.com/rancher/system-agent).

<Tabs>
<TabItem value="example1" label="Example 1: one time instructions" default>

```json showLineNumbers
{"instructions":
    [
        {
            "name":"set hostname",
            "command":"hostnamectl",
            "args": ["set-hostname", "myHostname"]
        },
        {
            "name":"stop sshd service",
            "command":"systemctl",
            "args": ["stop", "sshd"]
        }
    ]
}
```

</TabItem>
<TabItem value="example2" label="Example 2: periodic instructions">

```json showLineNumbers
{"periodicInstructions":
    [
        {
            "name":"set hostname",
            "image":"ghcr.io/rancher-sandbox/elemental-example-plan:main"
            "command": "run.sh"
        }
    ]
}
```

</TabItem>
<TabItem value="example3" label="Example 3: files">

```json showLineNumbers
{"files":
    [
        {
            "content":"Welcome to the system",
            "path":"/etc/motd",
            "permissions": "0644"
        }
    ]
}
```

</TabItem>
<TabItem value="example4" label="Example 4: probes">

```json showLineNumbers
{"probes":
    "probe1": {
        "name": "Service Up",
        "httpGet": {
            "url": "http://10.0.0.1/healthz",
            "insecure": "false",
            "clientCert": "....",
            "clientKey": "....",
            "caCert": "....."
        }   
    }
}
```

</TabItem>
</Tabs>


---

## Article: elemental_behind_proxy.md

---
sidebar_label: Elemental behind proxy
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/elemental_behind_proxy"/>
</head>

import RegistrationProxy from "!!raw-loader!@site/examples/proxy/registration-proxy.yaml"
import SeedimageProxy from "!!raw-loader!@site/examples/proxy/seedimage-proxy.yaml"
import ClusterProxy from "!!raw-loader!@site/examples/proxy/cluster-proxy.yaml"

## Introduction

In a lot of enterprise environments, servers or VMs running on premises do not have direct Internet access. Instead, the connection to external services is done through a HTTP(S) proxy for security reasons. This tutorial shows you how to set up an Elemental deployment in such an environment.

:::caution important note
This guide will not cover the Rancher installation behind a proxy. It's a different use case and you can find the detailed documentation [here](https://ranchermanager.docs.rancher.com/pages-for-subheaders/rancher-behind-an-http-proxy).
:::

:::info info
For this documentation, we assume you are using a SUSE family system (like SLE Micro), so proxy settings have to be written in `/etc/sysconfig/proxy`.
:::

Proxy settings must be configured in the following locations:

- Machine Registration Endpoint
- SeedImage resource
- Elemental cluster configuration

The `elemental-system-agent` needs proxy settings to reach the Rancher Manager.
To achieve that, you need to fill the cloud-init section of the Machine Registration Endpoint.

You can do it either with [UI](quickstart-ui#add-a-machine-registration-endpoint) or [CLI](quickstart-cli#prepare-your-kubernetes-resources).

<Tabs>
<TabItem value="cliRegistration" label="CLI" default>
<CodeBlock language="yaml" title="registration.yaml" showLineNumbers>{RegistrationProxy}</CodeBlock>
</TabItem>
<TabItem value="uiRegistration" label="UI" default>

![Add proxy settings in Machine Registration](images/proxy-settings-machine-registration-ui.png)
</TabItem>
</Tabs>

## Elemental-register

[Elemental-register](architecture-components#elemental-register-command-line-tool) is the first communication endpoint between the new host and Rancher Manager, this is the first place where proxy settings need to be set.

:::warning warning
At the time of writing, it's only possible to configure proxy settings for the ISO with the CLI. The proxy settings aren't implemented in the UI.
:::

The process happens when you boot your Elemental ISO for the first time, in order to configure the proxy settings you have to include a `cloud-init` definition in the ISO.
To do that, you have to create a [SeedImage](seedimage-reference) definition.

<CodeBlock language="yaml" title="seedimage.yaml" showLineNumbers>{SeedimageProxy}</CodeBlock>

Apply the YAML with `kubectl` and then, print your SeedImage definition to get the URL to download it:

```bash showLineNumbers
kubectl apply -f <my_seedimage_yaml_file>
kubectl get seedimage <seed_image_name> -n <namespace> -o yaml
```

Boot the ISO and you should see your new system appears in [Machine inventory](machineinventory-reference.md).

## Create Elemental cluster

For this step, you can use either the UI or CLI.

<Tabs>
<TabItem value="cliCluster" label="CLI" default>
<CodeBlock language="yaml" title="cluster.yaml" showLineNumbers>{ClusterProxy}</CodeBlock>
You can see that proxy settings are added below `agentEnvVars`.
</TabItem>
<TabItem value="uiCluster" label="UI" default>

![Add proxy settings for Elemental cluster](images/proxy-settings-cluster-ui.png)
</TabItem>
</Tabs>


---

## Article: elementaloperatorchart-reference.md

---
sidebar_label: Elemental Operator Helm Chart
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/elementaloperatorchart-reference"/>
</head>

# Elemental Operator Helm Chart

The <Vars name="elemental_operator_name" link="elemental_operator_url" /> is responsible for managing the Elemental versions and maintaining a machine inventory to assist with edge or bare metal installations.

The associated chart bootstraps an elemental-operator deployment on the [Rancher Manager v2.6](https://rancher.com/docs/rancher/v2.6/) cluster using the [Helm](https://helm.sh) package manager.

## Prerequisites

- Rancher Manager version v2.6
- Helm client version v3.8.0+

## Get Helm chart info

```console showLineNumbers
helm pull oci://registry.suse.com/rancher/elemental-operator-chart
helm show all oci://registry.suse.com/rancher/elemental-operator-chart
```

## Install Chart

```console showLineNumbers
helm install --create-namespace -n cattle-elemental-system elemental-operator-crds \
             oci://registry.suse.com/rancher/elemental-operator-crds-chart
helm install --create-namespace -n cattle-elemental-system elemental-operator \
             oci://registry.suse.com/rancher/elemental-operator-chart
```

The command deploys elemental-operator on the Kubernetes cluster in the default configuration.

_See [configuration](#configuration) below._

_See [helm install](https://helm.sh/docs/helm/helm_install/) for command documentation._

## Uninstall Chart

```console showLineNumbers
helm uninstall -n cattle-elemental-system elemental-operator
```

This removes all the Kubernetes components associated with the chart and deletes the release.

_See [helm uninstall](https://helm.sh/docs/helm/helm_uninstall/) for command documentation._

## Upgrading Chart

```console showLineNumbers
helm upgrade -n cattle-elemental-system \
             --install elemental-operator \
             oci://registry.suse.com/rancher/elemental-operator-chart
```

_See [helm upgrade](https://helm.sh/docs/helm/helm_upgrade/) for command documentation._

## Configuration

See [Customizing the Chart Before Installing](https://helm.sh/docs/intro/using_helm/#customizing-the-chart-before-installing). To see all configurable options with detailed comments, visit the chart's [values](#values), or run these configuration commands:

```console showLineNumbers
helm show values oci://registry.suse.com/rancher/elemental-operator-chart
```

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| image.empty | string | `rancher/pause:3.1` |  |
| image.repository | string | `registry.suse.com/rancher/elemental-operator-chart` | Source image for elemental-operator with repository name  |
| image.tag | tag | `""` |  |
| image.imagePullPolicy | string | `IfNotPresent` |  |
| noProxy | string | `127.0.0.0/8,10.0.0.0/8,172.16.0.0/12,192.168.0.0/16,.svc,.cluster.local" | Comma separated list of domains or ip addresses that will not use the proxy |
| global.cattle.systemDefaultRegistry | string | `""` | Default container registry name  |
| sync_interval | string | `"60m"` | Default sync interval for upgrade channel |
| sync_namespaces | list | `[]` | Namespace the operator will watch for, leave empty for all |
| debug | bool | `false` | Enable debug output for operator |
| nodeSelector.kubernetes.io/os | string | `linux` |  |
| tolerations | object | `{}` |  |
| tolerations.key | string | `cattle.io/os` |  |
| tolerations.operator | string | `"Equal"` |  |
| tolerations.value | string | `"linux"` |  |
| tolerations.effect | string | `NoSchedule` |  |


---

## Article: extra-rpms.md

---
sidebar_label: Add third party RPMs
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/custom-install/extra-rpms"/>
</head>

## How to add third party RPMs at install time

This example is covering the case in which extra RPMS (e.g. specific hardware drivers) are included
into the ISO and during the installation they are installed over the OS image.

For that use case the following files are required:

* additional RPMs to install
* additional cloud-config file defining hooks to copy and install binaries as installation hooks

We can handle it all into a single `SeedImage.spec.cloud-config` section assuming the additional
RPM is available to download it from a remote server at installation time.

Consider the following cloud-config data which could be used as the content of the `cloud-config`
section in a [SeedImage resource](seedimage-reference.md#seedimagespec-reference).

```yaml showLineNumbers
name: "Install extra drivers"
stages:
  before-install:
  # Preload data to the persistent storage
  # During installation persistent partition is mounted at /run/elemental/persistent
  - downloads:
    - path: /tmp/some_package.rpm
      url: "<REMOTE_PACKAGE_URL>"
      permissions: 0777
  - commands:
    - mkdir -o /run/elemental/persistent/extra-pkgs
    - cp -p /tmp/some_package.rpm /run/elemental/persistent/extra-pkgs

  after-install-chroot:
  # Install the package at install time
  - commands:
    - rpm -iv /run/elemental/extra-pkgs/some_package.rpm

  # Include to the install system analog upgrade and reset hooks
  - files:
    - path: /oem/extra-pkg.yaml
      permissions: 0664
      content: |
        name: "Install extra drivers"
        stages:
          after-upgrade-chroot:
          # Install the package after upgrading to a new image
          - commands:
            - rpm -iv /run/elemental/extra-pkgs/some_package.rpm
          after-reset-chroot:
          # Install the package on reset
          - commands:
            - rpm -iv /run/elemental/extra-pkgs/some_package.rpm
```

Note the installation hooks only cover installation procedures, so that additional cloud-config data should
be also part of the installed system in order to keep installing the package as part of the upgrade or reset
processes.

### Repacking a generated ISO image with extra files

Alternatively, if having the dynamic download of content at install time is not a desired behavior
and an ISO already including all the extra binaries is the actual goal this can is also possible
but requires manually repacking the ISO.

Using `xorriso`, the linux utility to create ISOs, this turns to be a relatively easy process.

Let's create an `overlay` directory to create the directory tree that needs to be
added into the ISO root. In that case the `overlay` directory could contain:

```yaml showLineNumbers
overlay/
  data/
    extra-pkgs/
      some_package.rpm
  iso-config/
    install_hooks.yaml
```

We are assuming the `install_hooks.yaml` is the content the actual cloud-config exposed in the previous
section which is manually included into image instead of being embedded in a SeedImage resource. The `data` folder is eventually including the binaries we want to append into the ISO.

Assuming we already downloaded an Elemental ISO tied to a specific MachineRegistration with the following
`xorriso` call all the `overlay` folder contents would be included into a new ISO:

```bash
xorriso -indev elemental.x86_64.iso -outdev elemental.custom.x86_64.iso -map overlay / -boot_image any replay
```

Requires `xorriso` equal or higher than version 1.5.

The contents of `install_hooks.yaml` could eventually be the same as the previous section but omitting the RPM
package download and also adapting the RPM path to `/run/initramfs/live/data/extra-pkgs/some_package.rpm` as
the ISO root folder is mounted at `/run/initramfs/live`.


---

## Article: hardwarelabels.md

---
sidebar_label: Hardware
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/hardwarelabels"/>
</head>

import Registration from "!!raw-loader!@site/examples/quickstart/registration-hardware-dhcphostname.yaml"

:::warning
Hardware Template Variables have been deprecated: please use the new
[Label Templates' Variables](label-templates#label-templates-variables) when possible.

Check the [deprecated variables page](label-templates-deprecated) and the
[conversion table](label-templates-deprecated#hardware-labels-to-new-label-templates-variables-table)
for a smooth transition.
:::

## Hardware Template Variables

When a node is registered, hardware data is collected and made available to the MachineRegistration in a way similar to [SMBIOS variables](smbios.md).

This data can be used for easy identification and selection via a [MachineSelector](machineinventoryselectortemplate-reference.md).

The following are available for templating:

| Variable                                                      | Description                                                           |
| ------------------------------------------------------------- | --------------------------------------------------------------------- |
| `${System Data/Runtime/Hostname}`                             | The hostname of the node (at registration time)                       |
| `${System Data/Memory/Total Physical Bytes}`                  | The total RAM memory in the node, expressed in bytes                  |
| `${System Data/CPU/Total Cores}`                              | Total CPU cores                                                       |
| `${System Data/CPU/Total Threads}`                            | Total CPU threads                                                     |
| `${System Data/CPU/Vendor}`                                   | CPU vendor                                                            |
| `${System Data/CPU/Model}`                                    | CPU model                                                             |
| `${System Data/GPU/Vendor}`                                   | GPU vendor (Only available if the node has an identifiable GPU)       |
| `${System Data/GPU/Model}`                                    | GPU model (Only available if the node has an identifiable GPU)        |
| `${System Data/Network/Number Interfaces}`                    | Number of network interfaces in the system                            |
| `${System Data/Network/{Iface name}/Name}`                    | Network interface name                                                |
| `${System Data/Network/{Iface name}/IsVirtual}`               | Boolean indicating virtual network interface                          |
| `${System Data/Block Devices/Number Devices}`                 | Number of block devices in the system (includes DVD and USB drives)   |
| `${System Data/Block Devices/{Disk name}/Name}`               | Device name of the block device (i.e. sda, sr0, vda, etc...)          |
| `${System Data/Block Devices/{Disk name}/Removable}`          | Whether this block device is removable (i.e. DVD)                     |
| `${System Data/Block Devices/{Disk name}/Size}`               | Total space in this block device, expressed in bytes                  |
| `${System Data/Block Devices/{Disk name}/Drive Type}`         | Drive type of this block device, see table below                      |
| `${System Data/Block Devices/{Disk name}/Storage Controller}` | Controller type for this block device connection, see table below     |

:::info info
On both `Block Devices` and `Network` the device name is used as a sub-block, as there could be more than one device.
:::

### Block device drive types

| Type    | Description                     |
|---------|---------------------------------|
| HDD     | Hard disk drive                 |
| FDD     | Floppy disk drive               |
| ODD     | Optical disk drive              |
| SSD     | Solid-state drive               |
| virtual | virtual drive i.e. loop devices |
| Unknown | unknown drive type              |

### Block device controller types

| Type    | Description                                                    |
|---------|----------------------------------------------------------------|
| IDE     | Integrated Drive Electronics                                   |
| SCSI    | Small computer system interface                                |
| NVMe    | Non-volatile Memory Express                                    |
| MMC     | Multi-media controller (used for mobile phone storage devices) |
| virtio  | Virtualized storage controller/driver                          |
| loop    | loop device                                                    |
| Unknown | unknown controller type                                        |



---

## Article: hostname.md

---
sidebar_label: Customize hostname
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/hostname"/>
</head>

## Customize hostname
### Elemental hostname management overview

When a host boots from the Elemental ISO, the hostname is temporarily set (*transient hostname*) to the one provided from the DHCP server.
If the DHCP server doesn't provide a hostname, the Elemental ISO provides a transient hostname
in the form: `rancher-${RANDOM}`.

As soon as the boot process is finished, the registration phase kicks in: the host connects to
the Elemental Operator, which creates a `MachineInventory` for the host.
Each host registered with the Elemental Operator is tracked by a `MachineInventory` resource.

**The `name` of the `MachineInventory` resource associated with the node is the permanent (static) hostname**
**eventually set to the host**.
This permanent hostname **is set on the node during the K8s cluster provisioning phase only**.
Before the K8s provisioning phase, the node hostname is either the DHCP assigned one or `rancher-${RANDOM}`.

For the remainder of this section we will refer to the `hostname` meaning the *permanent* hostname,
i.e., the hostname that is set after the host has been provisioned as part of a K8s cluster.

### Default hostname

The default name assigned to each newly created MachineInventory is in the form `m-{$UUID}`.
When the host is provisioned as part of a Cluster, that `m-{UUID}` name is set as the hostname of
the corresponding host, overriding the previous assigned hostname (`rancher-{$RANDOM}` or the DHCP assigned one).

### Set a custom hostname

The hostname can be specified setting the [`machineName`](machineregistration-reference.md#machinename) field in the
['MachineRegistration'](machineregistration-reference.md) resource.

The hostname set in the `machineName` field is expected to be in a template form, in order to be uniquely generated
for each registering node, using the [Random](label-templates-random.md), [SMBIOS](smbios.md) and [Hardware](hardwarelabels.md)
variables from the [Label Templates](label-templates.md) feature.

:::caution important note
The `machineName` field in the `MachineRegistration` resource is used as the blueprint not
only for the hostname of the registering host, but also for the name of the `MachineInventory` resource
created to track the host.

This means that if you don't use a templated `machineName` such to generate a unique name for each
host that will boot using the same `MachineRegistration` data (i.e., the same ISO), only the first
registering host will be successful while the others will fail: the `MachineInventory` name must be
unique.
:::

import Registration from "!!raw-loader!@site/examples/quickstart/registration-hardware-dhcphostname.yaml"

### Keep the hostname assigned from DHCP
In order to keep the hostname assigned from the DHCP server before the host registers to the operator,
the `MachineRegistration` [`machineName field`](machineregistration-reference.md#machinename) should be set
to the `${System Data/Runtime/Hostname}` [Hardware Label](hardwarelabels.md).

This way Elemental will use the current hostname as the `MachineInventory` name during
the registration phase, which will be later set as the static hostname of the host during the
provisioning phase.

<CodeBlock language="yaml" title="registration example with hostname and MachineInventory name set on the hostname got by the DHCP server" showLineNumbers>{Registration}</CodeBlock>


---

## Article: index.md

---
slug: /
sidebar_label: Overview
title: ''
---

import ButtonGroup from '@mui/material/ButtonGroup';
import Button from '@mui/material/Button';

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com"/>
</head>

# OS Management for Rancher

Manage your OS appliance for Kubernetes nodes as simple OCI images and turn them into self installing
ISO or disk images.

<ButtonGroup size="small">
  <Button variant="contained" href="https://elemental.docs.rancher.com/release-notes">Latest Release:</Button>
  <Button variant="outlined" href="https://elemental.docs.rancher.com/release-notes">v1.7.3</Button>
</ButtonGroup>

## Welcome to Elemental

Elemental is a software stack enabling centralized, full cloud-native OS management with Kubernetes.

The Elemental Stack consists of a handful of packages on top of SLE Micro:

- **elemental-toolkit** - Includes a set of OS utilities to enable OS management via containers. Includes dracut modules, bootloader configuration, cloud-init style configuration services, etc.
- **elemental-operator** - Connects to Rancher Manager and handles MachineRegistration and MachineInventory CRDs.
- **elemental-register** - Registers machines via machineRegistrations and installs them via elemental-cli.
- **elemental-cli** - Installs any elemental-toolkit based derivative. Basically an installer based on our A/B install and upgrade system.
- **rancher-system-agent** - Runs on the installed system and gets instructions ("Plans") from Rancher Manager what to install and run on the system.

Cluster Node OSes are built and maintained via container images through the <Vars name="elemental_cli_name" /> and they can be installed on new hosts using the <Vars link="elemental_ui_url" name="elemental_ui_name" /> for [Rancher Manager](https://www.rancher.com/products/rancher) or the <Vars link="elemental_cli_url" name="elemental_cli_name" />.

The <Vars link="elemental_operator_url" name="elemental_operator_name" /> and the <Vars link="ranchersystemagent_url" name="ranchersystemagent_name" /> enable Rancher Manager to fully control Elemental clusters, from the installation and management of the OS on the Nodes to the provisioning of new K3s or RKE2 clusters in a centralized way.

### Elemental on x86-64 hardware

Elemental is production ready and fully supported on x86-64 starting with Rancher v2.7.0.

### Elemental on ARM hardware

ARM (aarch64) is functional in the development stage. ARM is currently only tested on Raspberry Pi 4 Model B with k3s 1.24.8 (or later). Feedback is welcome.

### Elemental on other hardware

Elemental is currently targeting 'edge' scenarios and does therefore not support other hardware. We will re-assess this as the market evolves.

## Ready to give it a try?

Get an Elemental Cluster up and running with your preferred method

* With Rancher manager [Elemental plugin](quickstart-ui)
* With the [Elemental CLI](quickstart-cli)

:::note What's next?
Want more details? Take a look at the [Architecture](architecture) section or reach out to the <Vars link="elemental_slack_url" name="elemental_slack_name" /> Slack channel.
:::


---

## Article: installation.md

---
sidebar_label: Installation
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/installation"/>
</head>

# Installation

## Overview

Elemental stack provides OS management using OCI containers and Kubernetes. The Elemental
stack installation encompasses the installation of the <Vars name="elemental_operator_name" /> into the
management cluster and the creation and use of installation media to provide
the OS into the Cluster Nodes. See [Architecture](architecture.md) section to
read about the interaction of the components.

The installation configuration is mostly applied and set as part of the registration process.
The registration process is done by the `elemental-register` (the <Vars name="elemental_operator_name" /> client part)
who is the responsible to register nodes in a Rancher management cluster and fetch the installation configuration.

Please refer to the [Quick Start](quickstart-cli.md) guide for simple step by step deployment instructions.

## Elemental Operator Installation

The <Vars name="elemental_operator_name" /> is responsible for managing the Elemental versions and
maintaining a machine inventory to assist with edge or bare metal installations. <Vars name="elemental_operator_name" />
requires a cluster including the Rancher Manager and it can be installed with a helm chart.

See <Vars name="elemental_operator_name" /> [helm chart reference](elementaloperatorchart-reference.md) for install,
uninstall, upgrade and configuration details.

## Prepare Kubernetes Resources

Once the <Vars name="elemental_operator_name" /> is up and running within the management cluster a couple of kubernetes
resources are required in order to prepare an Elemental based cluster deployment.

* [MachineInventorySelectorTemplate](machineinventoryselectortemplate-reference.md):
  This resource identifies the criteria to match registered boxes (listed as part of the MachineInventory)
  against available Rancher 2.6 Clusters. As soon as there is a match the selected kubernetes cluster takes
  ownership of the registered box.
  
* [MachineRegistration](machineregistration-reference.md):
  This resource defines OS deployment details for any machine attempting to register. The machine
  registration is the entrance for Elemental nodes as it handles the authentication (based on TPM),
  the OS deployment and the node inclusion into to the MachineInventory so it can be added
  to a cluster when there is a match based on a MachineInventorySelectorTemplate. The MachineRegistration
  object includes the machine registration URL that nodes use to register against it.

A Rancher Cluster resource is also required to deploy Elemental, it can be manually created as exemplified in
the [Quick Start](quickstart-cli.md) guide or created from the Rancher 2.6 UI.

## Prepare Installation Media

The installation media is the media that will be used to kick start an OS deployment. Currently
the supported media is a live ISO. The live ISO must include the registration configuration yaml hence it must
crafted once the MachineRegistration is created. The installation media is generated by creating [Seed Image](/seedimage-reference.md)
resources (see [quick start](quickstart-cli#preparing-the-installation-seed-image) and [custom images](/custom-images.md#create-a-custom-bootable-installation-iso)).

The live ISO supports PXE booting for direct integration with [SUSE Manager](https://documentation.suse.com/suma/4.3/en/suse-manager/client-configuration/autoinst-distributions.html#based-on-iso-image).

Within MachineRegistration only a subset of OS installation parameters can be configured, all available parameters are listed
at [MachineRegistration](machineregistration-reference.md) reference page.

In order to configure the installation beyond the common options provided within the
[`elemental.install`](machineregistration-reference.md#configelementalinstall) section a `config.yaml`
configuration file can be included into the ISO (see [Custom Images](/custom-install.md#custom-elemental-client-configuration-file)).
Note any configuration applied as part of `elemental.install` section of the MachineRegistration will be
applied on top of the settings included in any custom `config.yaml` file.

Most likely the cloud-init configuration is enough to configure and set the deployed node at boot, however
if for some reason firstboot actions or scripts are required it is possible to also include
Rancher System Agent plans into the installation media. Refer to the [Elemental Plans](elemental-plans.md) section for details and
some example plans. The plans could be included into the squashed rootfs at `/var/lib/elemental/agent/plans`
folder and they would be seen by the system agent at firstboot.

## Start Installation Process

The installation starts by booting the installation media on a node. Once the installation media has booted it will
attempt to contact the management cluster and register to it by calling `elemental-register` command.
As the registration yaml configuration is already included into the ISO `elemental-register` knows the registration URL and
any other required data for the registration.

On a succeeded registration the installation media will start the installation into the host based
on the configuration already included in the media and the MachineRegistration parameters. As soon as the installation
is done the node is ready to reboot. The deployed OS includes a system agent plan to
kick start a regular rancher provisioning process to install the selected kubernetes version, once booted, after
some minutes the node installation is finalized and the node is included into the cluster and visible through
the Rancher UI.

## Deployed Partition Table

Once the operating system is installed the OS partition table, according to default values, will look like

| Label          | Default Size    | Contains                                                    |
|----------------|-----------------|-------------------------------------------------------------|
| COS_GRUB       | 64 MiB          | UEFI Boot partition                                         |
| COS_STATE      | 8 GiB           | A/B bootable file system images constructed from OCI images |
| COS_OEM        | 64 MiB          | OEM cloud-config files and other data                       |
| COS_RECOVERY   | 4 GiB           | Recovery file system image if COS_STATE is destroyed        |
| COS_PERSISTENT | Remaining space | All contents of the persistent folders                      |

Note this is the basic structure of any OS built by the <Vars name="elemental_toolkit_name" link="elemental_toolkit_url" />

## Elemental Immutable Root

One of the characteristics of Elemental OSes is the setup of an immutable root
filesystem where some ephemeral or persistent locations are applied on top of
it. The default folders structure is listed in the matrix below.

| Path                    | Read-Only | Ephemeral | Persistent |
|-------------------------|:---------:|:---------:|:----------:|
| /                       |     x     |           |            |
| /etc                    |           |     x     |            |
| /etc/cni                |           |           |     x      |
| /etc/iscsi              |           |           |     x      |
| /etc/rancher            |           |           |     x      |
| /etc/ssh                |           |           |     x      |
| /etc/systemd            |           |           |     x      |
| /srv                    |           |     x     |            |
| /home                   |           |           |     x      |
| /opt                    |           |           |     x      |
| /root                   |           |           |     x      |
| /var                    |           |     x     |            |
| /usr/libexec            |           |           |     x      |
| /var/lib/cni            |           |           |     x      |
| /var/lib/kubelet        |           |           |     x      |
| /var/lib/rancher        |           |           |     x      |
| /var/lib/elemental      |           |           |     x      |
| /var/lib/NetworkManager |           |           |     x      |
| /var/lib/calico         |           |           |     x      |
| /var/log                |           |           |     x      |


---

## Article: inventory-management.md

---
sidebar_label: Inventory Management
title: ''
version_badge: '1.3.0'
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/inventory-management"/>
</head>

## Inventory Management

The Elemental operator can hold an inventory of machines and
the mapping of the machine to it's configuration and assigned cluster.

### MachineInventory

The `MachineInventory` holds all the relevant information for a registered machine.  

Upon successful registration, the `MachineInventory` inherits all the `machineInventoryLabels`
and the `machineInventoryAnnotations` defined in the associated `MachineRegistration`.

The registering host sends also a bunch of [`system annotations`](#system-annotations) tracking information regarding the authentication
method used, the running OS version and the current IP address.

Those annotations are added to the associated [MachineInventory](machineinventory-reference.md).

Elemental machines attempt a registration update every 30 minutes to update labels and annotations.

#### System Annotations
| Key                                         | Description                                                                                |
|---------------------------------------------|--------------------------------------------------------------------------------------------|
| elemental.cattle.io/auth                    | Authentication used during registration (one of 'tpm', 'emulated-tpm', 'mac', 'sys-uuid')  |
| elemental.cattle.io/registration-ip         | IP address used during last registration                                                   |
| elemental.cattle.io/os.unmanaged            | Only present when set to 'true', disables OS management functionality on the tracked host  |
| elemental.cattle.io/name                    | 'NAME' from /etc/os-release                                                                |
| elemental.cattle.io/version                 | 'VERSION' from /etc/os-release                                                             |
| elemental.cattle.io/version-id              | 'VERSION_ID' from /etc/os-release                                                          |
| elemental.cattle.io/id                      | 'ID' from /etc/os-release                                                                  |
| elemental.cattle.io/pretty-name             | 'PRETTY_NAME' from /etc/os-release                                                         |
| elemental.cattle.io/image                   | 'IMAGE' from /etc/os-release                                                               |
| elemental.cattle.io/cpe-name                | 'CPE_NAME' from /etc/os-release                                                            |

#### Reference

```yaml
apiVersion: elemental.cattle.io/v1beta1
kind: MachineInventory
metadata:
  # Machine annotations can be useful to identify hosts
  annotations:
    elemental.cattle.io/auth: tpm
    elemental.cattle.io/registration-ip: 192.168.122.152
  labels:
    # A label inherited from the MachineRegistration definition
    element: fire
    # Generic SMBIOS labels that are typically populated with
    # the MachineRegister approach
    machineUUID: f266c64b-3972-40e7-9937-3dc4a311436c
    manufacturer: QEMU
    productName: Standard-PC-Q35-ICH9-2009
    serialNumber: Not-Specified
    # Custom labels can be applied to each MachineInventory
    myCustomLabel: foo 
  name: m-479ab68e-00ff-4081-a731-5b1a76610289
  # The namespace must match the namespace of the cluster
  # assigned to the clusters.provisioning.cattle.io resource
  namespace: fleet-default
  # A reference to the MachineInventorySelector that links the 
  # machine to a Cluster definition
  ownerReferences:
  - apiVersion: elemental.cattle.io/v1beta1
    controller: true
    kind: MachineInventorySelector
    name: fire-machine-selector-qcn7d
    uid: 0a1f751e-4ca9-4a0d-919a-97ba1f434d12
spec:
  # The hash of the TPM EK public key. This is used if you are
  # using TPM2 to identify nodes. Nodes can report their TPM
  # hash by using the MachineRegistration.
  tpmHash: d68795c6192af9922692f050b...
```

### MachineRegistration

`MachineRegistration` holds information on how to install, reset, and configure all connected Elemental machines.  

The `spec.machineInventoryLabels` and `spec.machineInventoryAnnotations` fields hold label and annotation templates
rendered to actual labels and annotations applied to the [MachineInventories](machineinventory-reference.md) tracking
the registered machines.

Elemental machines attempt a registration update every 30 minutes to update labels and annotations.

While it's possible to modify the `spec.config` definition, updates to the `spec.config` will be ignored by machines that already completed installation.
Machines that couldn't complete the installation will try again every 30 minutes by reloading the remote `MachineRegistration` definition.
This can be useful to correct `spec.config` mistakes that prevent successful installation (for ex. `spec.config.elemental.install.device`), without having to create a new `MachineRegistration` and a new ISO.

#### Reference

```yaml
apiVersion: elemental.cattle.io/v1beta1
kind: MachineRegistration
metadata:
  name: fire-nodes
  # The namespace must match the namespace of the cluster
  # assigned to the clusters.provisioning.cattle.io resource
  namespace: fleet-default
spec:
  # The cloud config that will be used to provision the node
  config:
    cloud-config:
      users:
        - name: root
          passwd: root
    elemental:
      install:
        reboot: true
        device: /dev/sda
        debug: true
      reset:
        enabled: true
        debug: true
        reset-persistent: true
        reset-oem: true
        reboot: true
  # Labels to be added to the created MachineInventory object
  machineInventoryLabels:
    element: fire
    manufacturer: "${System Information/Manufacturer}"
    productName: "${System Information/Product Name}"
    serialNumber: "${System Information/Serial Number}"
    machineUUID: "${System Information/UUID}"
  # Annotations to be added to the created MachineInventory object
  machineInventoryAnnotations: {}
```


---

## Article: kubernetesversions.md

---
sidebar_label: Kubernetes versions
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/kubernetesversions"/>
</head>

## Valid Versions

The list of valid versions for the `kubernetesVersion` field can be determined
from the Rancher metadata using the following commands.

__k3s:__
```bash showLineNumbers
curl -sL https://raw.githubusercontent.com/rancher/kontainer-driver-metadata/release-v2.6/data/data.json | jq -r '.k3s.releases[].version'
```
__rke2:__
```bash showLineNumbers
curl -sL https://raw.githubusercontent.com/rancher/kontainer-driver-metadata/release-v2.6/data/data.json | jq -r '.rke2.releases[].version'
```


---

## Article: label-templates-baseboard.md

---
sidebar_label: BaseBoard
title: ''
---

## BaseBoard Template Variables

The information collected in the BaseBoard variables defines attributes of a system baseboard (for
example, a motherboard, planar, server blade, or other standard system module).

Subset of SMBIOS Baseboard (or Module) information data (type 2).

For more information on SMBIOS see the
[SMSBIOS specification](https://www.dmtf.org/sites/default/files/standards/documents/DSP0134_3.7.1.pdf).

| Variable                    | Description               | from  |
| --------------------------- | ------------------------- | ----- |
| `${BaseBoard/AssetTag}`     | Motherboard asset tag     | 1.7.0 |
| `${BaseBoard/Product}`      | Motherboard product name  | 1.7.0 |
| `${BaseBoard/SerialNumber}` | Motherboard serial number | 1.7.0 |
| `${BaseBoard/Vendor}`       | Motherboard vendor        | 1.7.0 |
| `${BaseBoard/Version}`      | Motherboard revision      | 1.7.0 |

---

## Article: label-templates-bios.md

---
sidebar_label: BIOS
title: ''
---

## BIOS Template Variables

BIOS information.

Subset of SMBIOS BIOS Information data (type 0).

For more information on SMBIOS see the
[SMSBIOS specification](https://www.dmtf.org/sites/default/files/standards/documents/DSP0134_3.7.1.pdf).

| Variable             | Description           | from  |
| -------------------- | --------------------- | ----- |
| `${BIOS/Date}`       | BIOS release date     | 1.7.0 |
| `${BIOS/Vendor}`     | BIOS vendor           | 1.7.0 |
| `${BIOS/Version}`    | BIOS version number   | 1.7.0 |

---

## Article: label-templates-chassis.md

---
sidebar_label: Chassis
title: ''
---

## Chassis Template Variables

The information in the Chassis variables defines attributes of the system’s mechanical enclosure(s).

Subset of SMBIOS System Enclosure or Chassis (type 3).

For more information on SMBIOS see the
[SMSBIOS specification](https://www.dmtf.org/sites/default/files/standards/documents/DSP0134_3.7.1.pdf).

| Variable                     | Description                | from  |
| ---------------------------- | -------------------------- | ----- |
| `${Chassis/AssetTag}`        | Chassis asset tag          | 1.7.0 |
| `${Chassis/SerialNumber}`    | Chassis serial number      | 1.7.0 |
| `${Chassis/TypeDescription}` | Chassis type (description) | 1.7.0 |
| `${Chassis/Type}`            | Chassis type (number/code) | 1.7.0 |
| `${Chassis/Vendor}`          | Chassis vendor             | 1.7.0 |
| `${Chassis/Version}`         | Chassis version            | 1.7.0 |

---

## Article: label-templates-cpu.md

---
sidebar_label: CPU
title: ''
---

## CPU Template Variables

The information collected in the CPU template variables defines attributes of the processors
installed on the system.

The data is exposed by processor number (0, 1, 2, ...).

To reference the first processor the processor number (_`<Number>`_) can be omitted
(e.g., the variable `${CPU/Processor/Model}` is the same of
`${CPU/Processor/0/Model}`).


| Variable                                   | Description                                    | from  |
| ------------------------------------------ | ---------------------------------------------- | ----- |
| `${CPU/Processor/Capabilities}`            | Comma separated list of processor capabilities | 1.7.0 |
| `${CPU/Processor/ID}`                      | Processor ID                                   | 1.7.0 |
| `${CPU/Processor/Model}`                   | Processor model                                | 1.7.0 |
| `${CPU/Processor/NumCores}`                | Number of cores                                | 1.7.0 |
| `${CPU/Processor/NumThreads}`              | Number of threads                              | 1.7.0 | 
| `${CPU/Processor/Vendor}`                  | Processor vendor                               | 1.7.0 |
| `${CPU/Processor/_<Number>_/Capabilities}` | Comma separated list of processor _number_     | 1.7.0 |
| `${CPU/Processor/_<Number>_/ID}`           | Processor _number_ ID                          | 1.7.0 |
| `${CPU/Processor/_<Number>_/Model}`        | Processor _number_ model                       | 1.7.0 |
| `${CPU/Processor/_<Number>_/NumCores}`     | Number of cores of processor _number_          | 1.7.0 |
| `${CPU/Processor/_<Number>_/NumThreads}`   | Number of threads of processor _number_        | 1.7.0 |
| `${CPU/Processor/_<Number>_/Vendor}`       | Processor _number_ vendor                      | 1.7.0 |
| `${CPU/TotalCores}`                        | Total number of cores in the system            | 1.7.0 |
| `${CPU/TotalProcessors}`                   | Number of processors in the system             | 1.7.0 |
| `${CPU/TotalThreads}`                      | Total number of threads in the system          | 1.7.0 |

---

## Article: label-templates-deprecated.md

---
sidebar_label: Deprecated Variables
title: ''
---

## Label Templates' deprecated variables

[**SMBIOS**](smbios) and [**Hardware**](hardwarelabels) variable families have been deprecated in the Elemental Operator v1.7.0.

While these variable families are still correctly rendered in current Elemental Operator versions,
their support will be removed in future Elemental Operator releases.

If you are still using these old Label Templates' variable families, please switch to
[the newer available ones](label-templates#label-templates-variables) if possible.

To simplify the adoption of the new variables, check the tables mapping the deprecated variables to the new
supported ones:
* [Hardware Labels to new Label Templates' variables](#hardware-labels-to-new-label-templates-variables-table)
* [SMBIOS labels to new Label Templates' variables](#smbios-labels-to-new-label-templates-variable-table)

#### Hardware Labels to new Label Templates' variables table

| Deprecated Variable                                               | Equivalent new Variable |
| ----------------------------------------------------------------- | ----------------------- |
| `${System Data/Block Devices/Number Devices}`                     | `${Storage/TotalDisks}` |
| `${System Data/Block Devices/_<DeviceName>_/Drive Type} `         | `${Storage/Disks/_<DeviceName>_/DriveType}` |
| `${System Data/Block Devices/_<DeviceName>_/Name} `               | `${Storage/Disks/_<DeviceName>_/Name}` |
| `${System Data/Block Devices/_<DeviceName>_/Removable} `          | `${Storage/Disks/_<DeviceName>_/Removable}` |
| `${System Data/Block Devices/_<DeviceName>_/Size} `               | `${Storage/Disks/_<DeviceName>_/Size}` |
| `${System Data/Block Devices/_<DeviceName>_/Storage Controller} ` | `${Storage/Disks/_<DeviceName>_/StorageController}` |
| `${System Data/CPU/Model}`                                        | `${CPU/Model}` |
| `${System Data/CPU/Total Cores}`                                  | `${CPU/TotalCores}` |
| `${System Data/CPU/Total Threads}`                                | `${CPU/TotalThreads}` |
| `${System Data/CPU/Vendor}`                                       | `${CPU/Processor/Vendor}` |
| `${System Data/GPU/Model}`                                        | `${GPU/GraphicsCard/ProductName}` |
| `${System Data/GPU/Vendor}`                                       | `${GPU/GraphicsCard/VendorName}` |
| `${System Data/Memory/Total Physical Bytes}`                      | `${Memory/TotalPhysicalBytes}` |
| `${System Data/Memory/Total Usable Bytes}`                        | `${Memory/TotalUsableBytes}` |
| `${System Data/Network/Number Interfaces}`                        | `${Network/TotalNICs}` |
| `${System Data/Network/_<IfaceName>_/IsVirtual}`                  | `${Network/NICs/_<IfaceName>_/IsVirtual}` |
| `${System Data/Network/_<IfaceName>_/MacAddress}`                 | `${Network/NICs/_<IfaceName>_/MacAddress}` |
| `${System Data/Network/_<IfaceName>_/Name}`                       | `${Network/NICs/_<IfaceName>_/Name}` |
| `${System Data/Runtime/Hostname}`                                 | `${Runtime/Hostname}` |

#### SMBIOS Labels to new Label Templates' variable table
| Deprecated Variable                                               | Equivalent new Variable |
| ----------------------------------------------------------------- | ----------------------- |
| `${BIOS Information/Address}`                        | _not available_ |
| `${BIOS Information/BIOS Revision}`                  | _not available_ |
| `${BIOS Information/ROM Size}`                       | _not available_ |
| `${BIOS Information/Release Date}`                   | `${BIOS/Date}`   |
| `${BIOS Information/Runtime Size}`                   | _not available_  |
| `${BIOS Information/Vendor}`                         | `${BIOS/Vendor}` |
| `${BIOS Information/Version}`                        | `${BIOS/Version}` |
| `${Base Board Information/Asset Tag}`                | `${BaseBoard/AssetTag}` |
| `${Base Board Information/Chassis Handle}`           | _not available_ |
| `${Base Board Information/Contained Object Handles}` | _not available_ |
| `${Base Board Information/Location In Chassis}`      | _not available_ |
| `${Base Board Information/Manufacturer}`             | `${BaseBoard/Vendor}` |
| `${Base Board Information/Product Name}`             | `${BaseBoard/Product}` |
| `${Base Board Information/Serial Number}`            | `${BaseBoard/SerialNumber}` |
| `${Base Board Information/Type}`                     | _not available_ |
| `${Base Board Information/Version}`                  | `${BaseBoard/Version}` |
| `${Chassis Information/Asset Tag}`                   | `${Chassis/AssetTag}` |
| `${Chassis Information/Boot-up State}`               | _not available_ |
| `${Chassis Information/Contained Elements}`          | _not available_ |
| `${Chassis Information/Height}`                      | _not available_ |
| `${Chassis Information/Lock}`                        | _not available_ |
| `${Chassis Information/Manufacturer}`                | `${Chassis/Vendor}` |
| `${Chassis Information/Number Of Power Cords}`       | _not available_ |
| `${Chassis Information/OEM Information}`             | _not available_ |
| `${Chassis Information/Power Supply State}`          | _not available_ |
| `${Chassis Information/SKU Number}`                  | _not available_ |
| `${Chassis Information/Security Status}`             | _not available_ |
| `${Chassis Information/Serial Number}`               | `${Chassis/SerialNumber}` |
| `${Chassis Information/Thermal State}`               | _not available_ |
| `${Chassis Information/Type}`                        | `${Chassis/TypeDescription}` |
| `${Chassis Information/Version}`                     | `${Chassis/Version}` |
| `${System Information/Family}`                       | `${Product/Family}` |
| `${System Information/Manufacturer}`                 | `${Product/Vendor}` |
| `${System Information/Product Name}`                 | `${Product/Name}`   |
| `${System Information/SKU Number}`                   | `${Product/SKU}`    |
| `${System Information/Serial Number}`                | `${Product/SerialNumber}` |
| `${System Information/UUID}`                         | `${Product/UUID}` |
| `${System Information/Version}`                      | `${Product/Version}` |
| `${System Information/Wake-up Type}`                 | _not available_ |


---

## Article: label-templates-gpu.md

---
sidebar_label: GPU
title: ''
---

## GPU Template Variables

The information collected in the GPU template variables defines the attributes of the graphics
cards installed in the system.

The data is exposed by graphic card number (0, 1, 2, ...)

To reference the first graphic card, the graphic card number (_`<Number>`_) can be omitted
(e.g., the variable `${GPU/GraphicsCards/Driver}` is the same of `${GPU/GraphicsCards/0/Driver}`).

| Variable                                      | Description                       | from  |
| --------------------------------------------- | --------------------------------- | ----- |
| `${GPU/GraphicsCards/Driver}`                 | GPU card driver                   | 1.7.0 |
| `${GPU/GraphicsCards/ProductName}`            | GPU card model                    | 1.7.0 |
| `${GPU/GraphicsCards/VendorName}`             | GPU card vendor                   | 1.7.0 |
| `${GPU/GraphicsCards/_<Number>_/Driver}`      | GPU card _number_ driver          | 1.7.0 |
| `${GPU/GraphicsCards/_<Number>_/ProductName}` | GPU card _number_ model           | 1.7.0 |
| `${GPU/GraphicsCards/_<Number>_/VendorName}`  | GPU card _number_ vendor          | 1.7.0 |
| `${GPU/TotalCards}`                           | Number of GPU cards in the system | 1.7.0 |

---

## Article: label-templates-legacy.md



---

## Article: label-templates-memory.md

---
sidebar_label: Memory
title: ''
---

## Memory Template Variables

The information collected in the Memory template variables, defines attributes of the 
(volatile) memory available on the system.

| Variable                        | Description               | from  |
| ------------------------------- | ------------------------- | ----- |
| `${Memory/Modules}`             | Memory modules            | 1.7.0 |
| `${Memory/SupportedPageSizes}`  | Allowed memory page sizes | 1.7.0 |
| `${Memory/TotalPhysicalBytes}`  | Total memory              | 1.7.0 |
| `${Memory/TotalUsableBytes}`    | Total memory available    | 1.7.0 |

---

## Article: label-templates-network.md

---
sidebar_label: Network
title: ''
---

## Network Template Variables

The information collected in the Network template variables defines attributes of the network
interfaces available in the system.

The NICs are tracked both by number (0, 1, 2, ...) and by interface name (_eth0_, _virbr0_, ...).

The IPv4Addresses and IPv6Addresses are tracked by number (0, 1, ...) and the first IPv4 and IPv6
address entries (i.e., IPv4Addresses.0 and IPv6Addresses.0) are tracked in the IPv4Address and the
IPv6Address variables.

For example, if `${Network/NICs/2/Name}` is equal to `eth0`, then the network interface
Mac Address can be retrieved via both the `${Network/NICs/2/MacAddress}` and the
`${Network/NICs/eth0/MacAddress}` variables.

`eth0` (first) IPv4 address can be retrieved with both the `${Network/NICs/eth0/IPv4Address}`
and the `${Network/NICs/eth0/IPv4Addresses/0}` variables.
The secondary `eth0` IPv4 address can be retrieved with `${Network/NICs/eth0/IPv4Addresses/1}` instead.

| Variable                                                  | Description                                       |  from |
| --------------------------------------------------------- | ------------------------------------------------- | ----- |
| `${Network/NICs/_<IfaceName>_/AdvertisedLinkModes}`       | Ethernet Link Modes advertised as available       | 1.7.0 |
| `${Network/NICs/_<IfaceName>_/Duplex}`                    | Ethernet Duplex mode                              | 1.7.0 |
| `${Network/NICs/_<IfaceName>_/IPv4Address}`               | First IPv4 address assigned to the interface      | 1.7.2 |
| `${Network/NICs/_<IfaceName>_/IPv4Addresses/_<Number>_}`  | \<Number\>ith assigned IPv4 address               | 1.7.2 |
| `${Network/NICs/_<IfaceName>_/IPv6Address}`               | First IPv6 address assigned to the interface      | 1.7.2 |
| `${Network/NICs/_<IfaceName>_/IPv6Addresses/_<Number>_}`  | \<Number\>ith assigned IPv6 address               | 1.7.2 |
| `${Network/NICs/_<IfaceName>_/IsVirtual}`                 | Whether the NIC is real ("false") or not ("true") | 1.7.0 |
| `${Network/NICs/_<IfaceName>_/MacAddress}`                | MAC of the interface                              | 1.7.0 |
| `${Network/NICs/_<IfaceName>_/Name}`                      | Interface name                                    | 1.7.0 |
| `${Network/NICs/_<IfaceName>_/Speed}`                     | Speed of the link                                 | 1.7.0 |
| `${Network/NICs/_<IfaceName>_/SupportedLinkModes}`        | Ethernet Link Modes supported by the NIC          | 1.7.0 |
| `${Network/NICs/_<IfaceName>_/SupportedPorts}`            | NIC's available ports                             | 1.7.0 |
| `${Network/NICs/_<Number>_/AdvertisedLinkModes}`          | Ethernet Link Modes advertised as available       | 1.7.0 |
| `${Network/NICs/_<Number>_/Duplex}`                       | Ethernet Duplex mode                              | 1.7.0 |
| `${Network/NICs/_<Number>_/IPv4Address}`                  | First IPv4 address assigned to the interface      | 1.7.2 |
| `${Network/NICs/_<Number>_/IPv4Addresses/_<Number>_}`     | \<Number\>ith assigned IPv4 address               | 1.7.2 |
| `${Network/NICs/_<Number>_/IPv6Address}`                  | First IPv6 address assigned to the interface      | 1.7.2 |
| `${Network/NICs/_<Number>_/IPv6Addresses/_<Number>_}`     | \<Number\>ith assigned IPv6 address               | 1.7.2 |
| `${Network/NICs/_<Number>_/IsVirtual}`                    | Whether the NIC is real ("false") or not ("true") | 1.7.0 |
| `${Network/NICs/_<Number>_/MacAddress}`                   | MAC of the interface                              | 1.7.0 |
| `${Network/NICs/_<Number>_/Name}`                         | Interface name                                    | 1.7.0 |
| `${Network/NICs/_<Number>_/Speed}`                        | Speed of the link                                 | 1.7.0 |
| `${Network/NICs/_<Number>_/SupportedLinkModes}`           | Ethernet Link Modes supported by the NIC          | 1.7.0 |
| `${Network/NICs/_<Number>_/SupportedPorts}`               | NIC's available ports                             | 1.7.0 |
| `${Network/TotalNICs/}`                                   | Number of NICs present in the system              | 1.7.0 |


---

## Article: label-templates-product.md

---
sidebar_label: Product
title: ''
---

## Product Template Variables

The information in this variables refers to attributes of the overall system.

Subset of SMBIOS System Information data (type 1).

For more information on SMBIOS see the
[SMSBIOS specification](https://www.dmtf.org/sites/default/files/standards/documents/DSP0134_3.7.1.pdf).

| Variable                  | Description               | from  |
| ------------------------- | ------------------------- | ----- |
| `${Product/Family}`       | System Family             | 1.7.0 |
| `${Product/Name}`         | Product Name              | 1.7.0 |
| `${Product/SKU}`          | System SKU Number         | 1.7.0 |
| `${Product/SerialNumber}` | System serial number      | 1.7.0 |
| `${Product/UUID}`         | System UUID               | 1.7.0 |
| `${Product/Vendor}`       | System vendor             | 1.7.0 |
| `${Product/Version}`      | System version            | 1.7.0 |

---

## Article: label-templates-random.md

---
sidebar_label: Random
title: ''
---

import Registration from "!!raw-loader!@site/examples/quickstart/registration-random-hostname.yaml"

## Random Template Variables

Random template variables are built-in in the Elemental Operator.

They allow to include random `Int`, `Hex` or `UUID` values in custom [label templates](label-templates).

The values are computed on the fly during the `label template variables` rendering.

:::info Random label templates are rendered only once

A label template containing a Random variable is rendered only if the
[MachineInventory](machineinventory-reference) of the registering
host doesn't have a value for that label yet (a label with the same key is missing or its value is empty).

So, the three cases in which a label template with a Random variable is rendered are:
1. the host is registering for the first time and the [MachineInventory](machineinventory-reference) is created anew
2. the label template has been added to the MachineRegistration after the host (re-)registered last time
3. the [MachineInventory](machineinventory-reference) label matching the label template (same label key) has been manually removed
or its value has been cleared out
:::

| Variable                 | Description                                                        | from  |
| ------------------------ | ------------------------------------------------------------------ | ----- |
| `${Random/UUID}`         | random UUID (e.g., fd95324a-c26b-4e28-8727-1dcec293a0ec)           | 1.7.0 |
| `${Random/Hex/[1-32]}`   | random hexadecimal string of the specified length (min 1, max 32)  | 1.7.0 |
| `${Random/Int/[MAXINT]`  | random integer (min 0, max MAXINT-1)                               | 1.7.0 |


:::note Rendering Examples
| template value        | rendered value example               |
|-----------------------|--------------------------------------|
| `${Random/UUID}`      | fd95324a-c26b-4e28-8727-1dcec293a0ec |
| `${Random/Hex/12}`    | acd231f222b8                         |
| `${Random/Int/10000}` | 9432                                 |
:::

The Random Template Variables can be handy for generating custom hostnames to be assigned to the registering host.

Since the hostname must be unique and is assigned through the
[MachineRegistration](machineregistration-reference) `spec.machineName` field, Random variables can be used
to ensure uniqueness of a group of host sharing the same custom prefix and/or suffix.

Check the [HowTo/Customize hostname](hostname) section for more information.

<CodeBlock language="yaml" title="registration example Random template variables" showLineNumbers>{Registration}</CodeBlock>


---

## Article: label-templates-runtime.md

---
sidebar_label: Runtime
title: ''
---

## Runtime Template Variables

Runtime template variables collect information about the system configuration.

| Variable              | Description      | from  |
| --------------------- | ---------------- | ----- |
| `${Runtime/Hostname}` | System hostname  | 1.7.0 |

---

## Article: label-templates-storage.md

---
sidebar_label: Storage
title: ''
---

## Storage Template Variables

The information collected in the Storage template variables defines attributes of the disks available
in the system.

The Disks are tracked both by number (0, 1, 2, ...) and by device name (_sda_, _vda_, _nvme0n1_, ...).

For example, if `${Storage/Disks/2/Name}` is equal to `sda`, the the disk model can be retrieved via
both the `${Storage/Disks/2/Model}` and the `${Storage/Disks/sda/Model}` variables.

:::note
The tracked disks include also optical disk drives and removable devices (e.g., USB sticks).
:::

| Variable                                             | Description                           | from  |
| ---------------------------------------------------- | ------------------------------------- | ----- |
| `${Storage/Disks/_<Number>_/DriveType}`              | Disk type                             | 1.7.0 |
| `${Storage/Disks/_<Number>_/Model}`                  | Disk model                            | 1.7.0 |
| `${Storage/Disks/_<Number>_/Name}`                   | Device name                           | 1.7.0 |
| `${Storage/Disks/_<Number>_/Removable}`              | is removable media?                   | 1.7.0 |
| `${Storage/Disks/_<Number>_/Size}`                   | Size ( bytes)                         | 1.7.0 |
| `${Storage/Disks/_<Number>_/StorageController}`      | Controller type                       | 1.7.0 |
| `${Storage/Disks/_<DeviceName>_/DriveType}`          | Disk type                             | 1.7.0 |
| `${Storage/Disks/_<DeviceName>_/Model}`              | Disk model                            | 1.7.0 |
| `${Storage/Disks/_<DeviceName>_/Name}`               | Device name                           | 1.7.0 |
| `${Storage/Disks/_<DeviceName>_/Removable}`          | is removable media?                   | 1.7.0 |
| `${Storage/Disks/_<DeviceName>_/Size}`               | Size ( bytes)                         | 1.7.0 |
| `${Storage/Disks/_<DeviceName>_/StorageController}`  | Controller type                       | 1.7.0 |
| `${Storage/TotalDisks}`                              | Number of disks present in the system | 1.7.0 |


### Disks' DriveTypes

| Type    | Description                     |
|---------|---------------------------------|
| HDD     | Hard disk drive                 |
| FDD     | Floppy disk drive               |
| ODD     | Optical disk drive              |
| SSD     | Solid-state drive               |
| virtual | virtual drive i.e. loop devices |
| Unknown | unknown drive type              |

### Disks' StorageController types

| Type    | Description                                                    |
|---------|----------------------------------------------------------------|
| IDE     | Integrated Drive Electronics                                   |
| SCSI    | Small computer system interface                                |
| NVMe    | Non-volatile Memory Express                                    |
| MMC     | Multi-media controller (used for mobile phone storage devices) |
| virtio  | Virtualized storage controller/driver                          |
| loop    | loop device                                                    |
| Unknown | unknown controller type                                        |


---

## Article: label-templates.md

---
sidebar_label: Label Templates
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/label-templates"/>
</head>

import Registration from "!!raw-loader!@site/examples/quickstart/registration-hardware-dhcphostname.yaml"

# Label Templates Overview
Elemental allows to specify *label templates* in the `spec.machineInventoryLabels` and `spec.machineInventoryAnnotations` sections of the
[MachineRegistration](machineregistration-reference) resources.

Their format is the canonical `key`:`value` used in Kubernetes labels and annotations.

These label templates are converted to actual labels and annotations attached to each
[MachineInventory](machineinventory-reference) resource created during the
[machine onboarding](architecture-machineonboarding) phase.

The resulting labels and annotations have the same `key` of the label template.

The associated `value` is generated:
* **rendering the [`label template variables`](#label-template-variables)** (if present)
* **`sanitizing`** the resulting value (in the case of the labels only)


:::info
The Elemental templating functionality covers also the [MachineRegistration](machineregistration-reference) `spec.machineName` field,
which defines the resulting hostname of the registering machine and the `name` of the associated
[MachineInventory](machineinventory-reference) resource.

See the [Machine Name](#custom-machine-names) section for more details.
:::

## Label Templates' Variables
Elemental Label Templating includes a set of predefined variables that could be used inside the `value` of
the *label templates* specified in the [MachineRegistration](machineregistration-reference).

The syntax used to specify label templates' variables is:

**\$\{ *VARFAMILY* \/ *VARPATH* \}**

where _VARFAMILY_ defines a group (family) of supported variables and _VARPATH_ defines the actual variable
name inside the belonging family group.

Elemental currently supports the following _template variable families_:

* [**BIOS**](label-templates-bios): **\$\{ BIOS \/** _VARPATH_ **\}**
* [**BaseBoard**](label-templates-baseboard): **\$\{ BaseBoard \/** _VARPATH_ **\}**
* [**CPU**](label-templates-cpu): **\$\{ CPU \/** _VARPATH_ **\}**
* [**Chassis**](label-templates-chassis): **\$\{ Chassis \/** _VARPATH_ **\}**
* [**GPU**](label-templates-gpu): **\$\{ GPU \/** _VARPATH_ **\}**
* [**Memory**](label-templates-memory): **\$\{ Memory \/** _VARPATH_ **\}**
* [**Network**](label-templates-network): **\$\{ Network \/** _VARPATH_ **\}**
* [**Product**](label-templates-product): **\$\{ Product \/** _VARPATH_ **\}**
* [**Runtime**](label-templates-runtime): **\$\{ Runtime \/** _VARPATH_ **\}**
* [**Storage**](label-templates-storage): **\$\{ Storage \/** _VARPATH_ **\}**

:::warning
All the _template variable families_ (but _`Random`_) are enabled only if [MachineRegistration](machineregistration-reference.md)'s
`elemental.registration.no-smbios` field is set to `false` (default).

When the `elemental.registration.no-smbios` field is set to `true`, the registering machines do not
send any data required for rendering the template variables, so no variables will be available, but
the **Random** variables, which are the only notable exception.

**Random** variables are always available since they are built-in on the operator side.
They are also special since they are computed only once: see the
[Random Template variables](label-templates-random) section for more details.
:::

Template variables can be mixed with static text to form the actual labels assigned to
([MachineInventories](machineinventory-reference)).

:::note Rendering Examples
* Label Template tracking the number of CPU cores of the registering host (assume host has 4 cores):
  * original label: **cpu: $\{CPU\/TotalCores\}-cores**
  * rendered label: **cpu: 4-cores**
* Label Template tracking the SMBIOS UUID of the registering host:
  * original label: **sbios-UUID: \$\{Product\/UUID\}**
  * rendered label: **sbios-UUID: fd95324a-c26b-4e28-8727-1dcec293a0ec**
:::

## Sanitization
Once the label template value has been rendered accordingly to the included [label template variables](#label-template-variables), the resulting value is `sanitized` before being assigned to the resulting label.

**The `sanitization` enforces the label value to only contain letters (capitalized or not), numbers and the hyphen (`-`), point (`.`) and underscore (`_`) characters**:
all the characters not included are substituted with an hyphen.

Any character at the beginning and at the end of the label value must be a letter or a number.
If it is not, it is dropped.

Two consecutive hyphens are replaced with one.

:::note Rendering Example
* Label Template with sanitization of prohibited chars:
  * original label: **sanitized: this:needs--sanitizing!**
  * rendered label: **sanitized: this-needs-sanitizing**
:::

## Usage of Label Templates
Label Templates allow to automatically attach and update labels and annotations to each host's
[MachineInventory](machineinventory-reference) every time an host registers to the Elemental Operator.

:::info
Registration happens not only during the [onboarding phase](architecture-machineonboarding): each host
re-registers every 30 minutes (and every time it reboots).
During the re-registration, the Label Templates in the associated
[MachineRegistration](machineregistration-reference) are re-evaluated and added/updated in the
[MachineInventory](machineinventory-reference).
:::

There are basically three main cases where the label templates can be of use:
* to attach hardware data to the Elemental Catalog
* to add selectors to pick up hosts for Cluster Provisioning
* to define a custom template for the Machine Names

### Hardware data for the Elemental catalog
The Label Templates' variables can be used to attach hardware data to each
[MachineInventory](machineinventory-reference) resource (which form the Elemental Catalog).

In this case, annotations may be a better choice since their values are not sanitized.

### Selectors for Cluster Provisioning
The `Label Templates` can be used to generate labels used to identify and select machines
with special hardware properties to form new Kubernetes Clusters.

The labels attached to each [MachineInventory](machineinventory-reference) are eligible to _selector_
for the [MachineInventorySelectorTemplate](machineinventoryselectortemplate-reference) resource
(see the [Kubernetes Cluster provisioning](architecture-clusterdeployment#kubernetes-cluster-provisioning)
section for more details).

### Custom Machine Names
The hostname of the onboarding machine can be specified using the
[MachineRegistration](machineregistration-reference) `spec.machineName` field.

`spec.machineName` value undergoes the same `Label Templates' variables` and `sanitization` processes
reserved to the `spec.machineInventoryLabels` field.

:::warning
There is one notable difference between the [MachineRegistration](machineregistration-reference) `spec.machineName` and `spec.machineInventorylabels` fields: during the [sanitization](#sanitization) process
the underscore (`_`) is substituted as the other forbidden characters (i.e., it is substituted by an
hyphen: `-`).
This is required as the underscore is not allowed in linux hostnames.
:::

For more information on how to define the hostname for Elemental hosts, see the
[HowTo/Customize hostname](hostname) section.

:::note Rendering Example
* Define an hostname template like `SLE-Micro-[random string of 6 hexadecimal values]`:
  * MachineRegistration spec: **machineName: SLE-Micro-\$\{Random\/Hex\/6\}**
  * MachineInventory name: **SLE-Micro-32ad41**
:::

## Label Templates in action

<CodeBlock language="yaml" title="registration example with Label Templates' variables" showLineNumbers>{Registration}</CodeBlock>


---

## Article: lvm-drives-example.md

---
sidebar_label: Extra LVM volume group
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/custom-install/lvm-drives-example"/>
</head>

## How to add extra LVM volume group disks during the installation

This example is covering the setup of a host with multiple disks and some of them used
as part of a LVM setup.

As an example, we have an host with three disks (`/dev/sda`, `/dev/sdb`
and `/dev/sdc`). 

The first disk is used for a regular Elemental installation
and the other remaining two are used as part of a LVM group where arbitrary logical volumes
are created, formatted and mounted at boot time via an extended `fstab` file.

For this example cloud-config steps are required in two different stages. First, some
[installation hooks](cloud-config-reference.md#elemental-client-cloud-config-hooks) are
needed to prepare and handle LVM volumes during the installation. Second, a cloud-config
is required at boot time to ensure the created LVM volumes are included in `/etc/fstab`
and consequently mounted.

Installation hooks can be included to the `SeedImage.spec.cloud-config` section with something
like:

```yaml showLineNumbers
apiVersion: elemental.cattle.io/v1beta1
kind: SeedImage
metadata:
  name: custom-seed
  namespace: fleet-default
spec:
  ...
  cloud-config:
    name: "Create LVM logic volumes over some physical disks"
    stages:
      post-install:
      - name: "Create physical volume, volume group and logical volumes"
        if: '[ -e "/dev/sdb" ] && [ -e "/dev/sdc" ]'
        commands:
        - | 
          # Create the physical volume, volume group and logical volumes
          pvcreate /dev/sdb /dev/sdc
          vgcreate elementalLVM /dev/sdb /dev/sdc
          lvcreate -L 8G -n elementalVol1 elementalLVM
          lvcreate -l 100%FREE -n elementalVol2 elementalLVM
          # Trigger udev detection
          if [ ! -e "/dev/elementalLVM/elementalVol1" ] || [ ! -e "/dev/elementalLVM/elementalVol2" ]; then
            sleep 10
            udevadm settle
          fi
          # Ensure devices are already available
          [ -e "/dev/elementalLVM/elementalVol1" ] || exit 1
          [ -e "/dev/elementalLVM/elementalVol2" ] || exit 1
          # Format logical volumes with a known label for later use in fstab
          mkfs.xfs -L eVol1 /dev/elementalLVM/elementalVol1
          mkfs.xfs -L eVol2 /dev/elementalLVM/elementalVol2
```

The LVM devices are created and formatted as desired. This is a good
example of an installation hook, as this setup is only needed once, at installation
time. As an alternative, the same action could be done on first boot, however it would
require a more sophisticated logic to ensure it's only applied once at first boot.

Finally, the boot time cloud-config data contains the mount point settings to trigger the
mounts. The Elemental OS `fstab` file is ephemeral and it's dynamically created at boot time.
That's why it doesn't exist during the installation and can't be used in an installation hook.

Consider the following example to customize the `/etc/fstab` file using the `cloud-config` section
of a MachineRegistration resource:

```yaml showLineNumbers
apiVersion: elemental.cattle.io/v1beta1
kind: MachineRegistration
metadata:
  name: my-nodes
  namespace: fleet-default
spec:
  ...
  config:
    ...
    cloud-config:
      stages:
        initramfs:
        - name: "Extend fstab to mount LVM logical volumes at boot"
          commands:
          - |
            echo "LABEL=eVol1 /run/elemental/eVol1  xfs defaults  0 0" >> /etc/fstab
            echo "LABEL=eVol2 /run/elemental/eVol2  xfs defaults  0 0" >> /etc/fstab
```

:::note
The `initramfs` [stage](cloud-config-reference.md) is the last stage before switching to the actual
root tree. At this stage, the `/etc/fstab` file already exists and can be adapted before
switching root. Once running in the final root tree, SystemD will handle the rest of the
initialization and apply it.
:::


---

## Article: machineinventory-reference.md

---
sidebar_label: MachineInventory reference
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/machineinventory-reference"/>
</head>

# MachineInventory
When a new host registers successfully, the <Vars name="elemental_operator_name" /> creates a MachineInventory resource representing that particular host.
The MachineInventory stores the TPM hash of the tracked host, retrieved during the registration process, and allows to execute arbitrary commands (plans) on the machine.

A MachineInventory has two conditions:

- `AdoptionReady`, which indicates the machine has been adopted by a selector to be part of a cluster.
- `Ready`, which indicates that the machine has been registered and provisioned with an Elemental OS.


---

## Article: machineinventoryselector-reference.md

---
sidebar_label: MachineInventorySelector reference
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/machineinventoryselector-reference"/>
</head>

# MachineInventorySelector reference

A MachineInventorySelector selects MachineInventories based on applied selectors (usually pattern matching on MachineInventory label values).

MachineInventorySelectors have two conditions:

- `InventoryReady`, turns to true if the MachineInventorySelector has found a matching MachineInventory and has successfully set itself as the MachineInventory owner.
- `Ready`, tracks if the selector already adopted a machine and started the kubernetes provisioning process (node bootstrap).

---

## Article: machineinventoryselectortemplate-reference.md

---
sidebar_label: MachineInventorySelectorTemplate reference
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/machineinventoryselectortemplate-reference"/>
</head>

# MachineInventorySelectorTemplate reference

The MachineInventorySelectorTemplate is a user defined resource that will be used as the blueprint to create the required [MachineInventorySelectors](machineinventoryselector-reference.md).

It is the resource responsible of defining the matching criteria to pair an inventoried machine with a Cluster resource.

The relevant key is the `selector` which includes label selector expressions.

```yaml title="MachineInventorySelectorTemplate" showLineNumbers
apiVersion: elemental.cattle.io/v1beta1
kind: MachineInventorySelectorTemplate
metadata:
  name: my-machine-selector
  namespace: fleet-default
spec:
  template:
    spec:
      selector:
        ...
```

`template.spec.selector` can include `matchLabels` and or `matchExpressions` keys.

#### template.spec.selector.matchLabels

It is a map of `{key,value}` pairs `(map[string]string)`. When multiple labels are provided all labels must match.

<details>
  <summary>Example</summary>

  ```yaml showLineNumbers
  ...
  spec:
    template:
      spec:
        selector:
          matchLabels:
            element: fire
            manufacturer: somevalue
  ```

</details>

A Cluster defined with the above selector will only attempt to provision nodes inventoried including these two labels.

#### template.spec.selector.matchExpressions

It is a list of label selectors, each label selectors can be defined as:

| Key               | Type              | Description                                     |
|-------------------|-------------------|-------------------------------------------------|
| key               | string            | This is the label key the selector applies on   |
| operator          | string            | Represents the relationship of the key to a set of values. Valid operators are 'In', 'NotIn', 'Exists' and 'DoesNotExist' |
| values            | []string          | Values is an array of string values. If the operator is 'In' or 'NotIn', the values array must be non-empty. If the operator is 'Exists' or 'DoesNotExist', the values array must be empty |

<details>
  <summary>Example</summary>

  ```yaml showLineNumbers
  ...
  spec:
    template:
      spec:
        selector:
          matchExpressions:
            - key: element
              operator: In
              values: [ 'fire' ]
            - key: manufacturer
              operator: Exists
  ```

</details>

A Cluster defined with the above selector will only attempt to provision nodes inventoried with the `element=fire` label and including a `manufacturer` label defined with any value.


---

## Article: machineregistration-reference.md

---
sidebar_label: MachineRegistration
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/machineregistration"/>
</head>

# MachineRegistration reference

The MachineRegistration resource allows:
* to configure the registration process
* to provide OS installation parameters
* to define the [Elemental services](architecture-services.md) enabled for the registering machine
Once created it generates the registration URL used by nodes to register and start the [machine onboarding](architecture-machineonboarding.md) process.

The MachineRegistration has a `Ready` condition which turns to true when the <Vars name="elemental_operator_name" /> has successfully generated the registration URL and an associated `ServiceAccount`. From this point on the target host can connect to the registration URL to kick off the provisioning process.

An HTTP GET request against the registration URL returns the _registration file_: a .yaml file containing the registration data (i.e., the _spec:config:elemental:registration_ section from the just created MachineRegistration).
The registration file contains all the required data to allow the target host to perform self registration and start the Elemental provisioning.

There are several keys that can be configured under a `MachineRegistration` resource spec.

```yaml title="MachineRegistration" showLineNumbers
apiVersion: elemental.cattle.io/v1beta1
kind: MachineRegistration
metadata:
  name: my-nodes
  namespace: fleet-default
spec:
  machineName: name
  machineInventoryLabels:
    label: value
  machineInventoryAnnotations:
    annotation: value
  config:
    cloud-config:
        ...
    elemental:
        registration:
            ...
        install:
            ... 
```

#### config.cloud-config

Contains the cloud-configuration to be injected in the node.  
Both yip and cloud-init syntax are supported. See the [Cloud Config Reference](cloud-config-reference.md) for full information.

The cloud-configuration provided in this field is not evaluated during the installation, it is
just added to the node so it gets evaluated on reboot.

#### config.network

Contains the Declarative Networking configuration, supporting integration with [CAPI IPAM Providers](https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20220125-ipam-integration.md#ipam-provider).  
See the [Declarative Networking Reference](networking.md) for full information.  
Any `configurator` value different than `none` will denote that the Network is managed by Elemental.  

| Key               | Type      | Default value | Description                                                                                                   |
|-------------------|-----------|---------------|---------------------------------------------------------------------------------------------------------------|
| configurator      | string    | none          | The network configurator type to use (`none`, `nmc`, `nmstate`, or `nmconnections`)                           |
| ipAddresses       | objRefMap | empty         | A map of `IPPool` references. Map keys can be used for IPAddress substitution in the network config template. |
| config            | obj       | empty         | The network config template. Syntax varies depending on the `configurator` in use.                            |

#### config.elemental.registration
Contains the configuration used for the connection and the initial registration to the <Vars name="elemental_operator_name" />.

Supports the following values:

| Key               | Type   | Default value | Description                                                                                                                                    |
|-------------------|--------|---------------|------------------------------------------------------------------------------------------------------------------------------------------------|
| url               | string | empty         | URL to connect to the <Vars name="elemental_operator_name" />                                                                                  |
| ca-cert           | string | empty         | CA to validate the certificate provided by the server at 'url' (required if the certificate is not signed by a public CA)                      |
| no-smbios         | bool   | false         | Whether SMBIOS data should be sent to the <Vars name="elemental_operator_name" /> (see the [SMBIOS reference](smbios.md) for more information) |
| no-toolkit        | bool   | false         | Disables the <Vars name="elemental_toolkit_name" link="elemental_toolkit_url" /> support and allows registration of an [unmanaged OS](unmanaged-os.md)                                           |

:::warning
The following values are for development purposes only.

| Key               | Type   | Default value | Description                                                                                                                                       |
|-------------------|--------|---------------|---------------------------------------------------------------------------------------------------------------------------------------------------|
| auth              | string | tpm           | Authentication method to use during registration, one of `tpm`, `mac` or `sys-uuid`. See [Authentication](authentication.md) for more information |
| emulate-tpm       | bool   | false         | This will use software emulation of the TPM (required for hosts without TPM hardware)                                                             |
| emulated-tpm-seed | int64  | 1             | Fixed seed to use with 'emulate-tpm'. Set to -1 to get a random seed. See [TPM](tpm.md) for more information                                      |

:::


#### config.elemental.install

Contains the installation configuration that would be applied via `elemental-register --install` when booted from an ISO and passed to [`elemental install`](https://github.com/rancher/elemental-toolkit/blob/main/docs/elemental_install.md)

Supports the following values:

| Key             | Type   | Default value | Description                                                                                                                                |
|-----------------|--------|---------------|--------------------------------------------------------------------------------------------------------------------------------------------|
| firmware        | string | efi           | Firmware to install ('efi' or 'bios')                                                                                                      |
| device          | string | empty         | Device to install the system to                                                                                                            |
| device-selector | string | empty         | Rules for picking device to install the system to                                                                                          |
| no-format       | bool   | false         | Don’t format disks. It is implied that COS_STATE, COS_RECOVERY, COS_PERSISTENT, COS_OEM partitions are already existing on the target disk |
| config-urls     | list   | empty         | Cloud-init config files locations                                                                                                          |
| iso             | string | empty         | Performs an installation from the ISO url instead of the running ISO                                                                       |
| system-uri      | string | empty         | Sets the system image source and its type (e.g. 'docker:registry.org/image:tag') instead of using the running ISO                          |
| debug           | bool   | false         | Enable debug output                                                                                                                        |
| tty             | string | empty         | Add named tty to grub                                                                                                                      |
| poweroff        | bool   | false         | Shutdown the system after install                                                                                                          |
| reboot          | bool   | false         | Reboot the system after install                                                                                                            |
| snapshotter     | obj    | empty         | Snapshotter configuration. See [reference](#configelementalinstallsnapshotter)                                                             |
| eject-cd        | bool   | false         | Try to eject the cd on reboot                                                                                                              |

:::warning warning
In case of using both `iso` and `system-uri` the `iso` value takes precedence
:::

It is only required to specify either the `device` or `device-selector` fields for a successful install, the rest of the parameters are all optional.

If both `device` and `device-selector` is specified the value of `device` is used and `device-selector` is ignored.

<details>
<summary>Example</summary>

  ```yaml showLineNumbers
  apiVersion: elemental.cattle.io/v1beta1
  kind: MachineRegistration
  metadata:
    name: my-nodes
    namespace: fleet-default
  spec:
    config:
      elemental:
        install:
          device: /dev/sda
          debug: true
          reboot: true
          eject-cd: true
          system-uri: registry.suse.com/rancher/sle-micro/5.5:latest
  ```
</details>

#### config.elemental.install.device-selector

The `device-selector` field can be used to dynamically pick device during installation. The field contains a list of rules that looks like the following:

<details>
<summary>Example device-selector based on device name</summary>
  ```yaml showLineNumbers
  device-selector:
  - key: Name
    operator: In
    values:
    - /dev/sda
    - /dev/vda
    - /dev/nvme0
  ```
</details>

<details>
<summary>Example device-selector based on device size</summary>
  ```yaml showLineNumbers
  device-selector:
  - key: Size
    operator: Lt
    values:
    - 100Gi
  - key: Size
    operator: Gt
    values:
    - 30Gi
  ```
</details>

The currently supported operators are:

| Operator            | Description                                       |
| ------------------- | ------------------------------------------------- |
| In                  | The key matches one of the provided values        |
| NotIn               | The key does not match any of the provided values |
| Gt                  | The key is greater than a single provided value   |
| Lt                  | The key is lesser than  a single provided value   |

The currently supported keys are:

| Key                 | Description                                                                    |
| ------------------- | ------------------------------------------------------------------------------ |
| Name                | The device name (eg. /dev/sda)                                                 |
| Size                | The device size (values can be specified using kubernetes resources, eg 100Gi) |

The rules are AND:ed together, which means all rules must match the targeted device.

#### config.elemental.install.snapshotter

You can configure how Elemental manages snapshots on the installed machine.  
New snapshots are created for example when [upgrading](./upgrade) the machine with a new OS image.  
The `loopdevice` snapshotter will unpack new images on a `ext4` filesystem, while the `btrfs` snapshotter will make use of the underlying [btrfs snapshots](https://archive.kernel.org/oldwiki/btrfs.wiki.kernel.org/index.php/SysadminGuide.html#Snapshots) functionality, greatly reducing the amount of disk space needed to store multiple snapshots.  

| Key      | Type   | Default value | Description                                                                     |
|----------|--------|---------------|---------------------------------------------------------------------------------|
| type     | string | loopdevice    | Type of device used to manage snapshots in OS images ('loopdevice' or 'btrfs'). |
| maxSnaps | int    | 2             | Maximum amount of snapshots to keep.                                            |

#### config.elemental.reset

Contains the reset configuration that would be applied via `elemental-register --reset`, when booted from the recovery partition and passed to [`elemental reset`](https://github.com/rancher/elemental-toolkit/blob/main/docs/elemental_reset.md)

Supports the following values:

| Key               | Type   | Default value | Description                                                                                                                                |
|-------------------|--------|---------------|-------------------------------------------------------------------------------------------------------------------|
| enabled           | bool   | false         | MachineInventories created from this MachineRegistration will have reset functionality enabled                    |
| reset-persistent  | bool   | true          | Format the COS_PERSISTENT partition                                                                               |
| reset-oem         | bool   | true          | Format the COS_OEM partition                                                                                      |
| config-urls       | list   | empty         | Cloud-init config files                                                                                           |
| system-uri        | string | empty         | Sets the system image source and its type (e.g. 'docker:registry.org/image:tag') instead of using the running ISO |
| debug             | bool   | false         | Enable debug output                                                                                               |
| poweroff          | bool   | false         | Shutdown the system after reset                                                                                   |
| reboot            | bool   | true          | Reboot the system after reset                                                                                     |

<details>
<summary>Example</summary>

  ```yaml showLineNumbers
  apiVersion: elemental.cattle.io/v1beta1
  kind: MachineRegistration
  metadata:
    name: my-nodes
    namespace: fleet-default
  spec:
    config:
      elemental:
        reset:
          enabled: true
          debug: true
          reset-persistent: true
          reset-oem: true
          reboot: true
          system-uri: registry.suse.com/rancher/sle-micro/5.5:latest
  ```
</details>

#### machineName

Template used to derive the hostname to be set to the node and as the name of the associated [MachineInventory](machineinventory-reference.md) kubernetes resource.

The value is interpolated using [Label Templates](label-templates.md).

:::info
If no `machineName` is specified, a default one in the form `m-$UUID` will be set.

See the [Customize Hostname section](hostname.md#customize-hostname) for further details.
:::

<details>
<summary>Example</summary>

  ```yaml showLineNumbers
  apiVersion: elemental.cattle.io/v1beta1
  kind: MachineRegistration
  metadata:
    name: my-nodes
    namespace: fleet-default
  spec:
    machineName: hostname-test-4
  ```

</details>

#### machineInventoryLabels

Labels to be set to the `MachineInventory` created from this `MachineRegistration`.

The label values are interpolated using [Label Templates](label-templates.md).

These labels could be used to establish a selection criteria in [MachineInventorySelectorTemplate](machineinventoryselectortemplate-reference.md).

Elemental nodes will run `elemental-register` every 30 minutes.

It is possible to update the `machineInventoryLabels` so that all registered nodes apply the new labels on the next successful registration update.

<details>
<summary>Example</summary>

  ```yaml showLineNumbers
  apiVersion: elemental.cattle.io/v1beta1
  kind: MachineRegistration
  metadata:
    name: my-nodes
    namespace: fleet-default
  spec:
    machineInventoryLabels:
      my.prefix.io/element: fire
      my.prefix.io/cpus: 32
      my.prefix.io/manufacturer: "${System Information/Manufacturer}"
      my.prefix.io/productName: "${System Information/Product Name}"
      my.prefix.io/serialNumber: "${System Information/Serial Number}"
      my.prefix.io/machineUUID: "${System Information/UUID}"
  ```

</details>

#### machineInventoryAnnotations

Annotations that will be set to the `MachineInventory` that is created from this `MachineRegistration`
`Key: value` type

<details>
<summary>Example</summary>

  ```yaml showLineNumbers
  apiVersion: elemental.cattle.io/v1beta1
  kind: MachineRegistration
  metadata:
    name: my-nodes
    namespace: fleet-default
  spec:
    machineInventoryAnnotations:
      owner: bob
      version: 1.0.0
  ```

</details>


---

## Article: managedosimage-reference.md

---
sidebar_label: ManagedOSImage reference
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/managedosimage-reference"/>
</head>

import ClusterTarget from "!!raw-loader!@site/examples/upgrade/upgrade-cluster-target.yaml"

# ManagedOSImage reference

The ManagedOSImage resource is responsible of defining an OS image or image version that needs to be applied to each node in a set of targeted Clusters.  
Once created, the ManagedOSImage resource can be updated with a new `osImage` or `managedOSVersionName` to trigger a new upgrade.  
Similarly, an existing ManagedOSImage can be updated to target new Clusters.  

There are several keys that can be configured under a `ManagedOSImage` resource spec.

<CodeBlock language="yaml" title="upgrade-cluster-target.yaml" showLineNumbers>{ClusterTarget}</CodeBlock>

#### ManagedOSImageSpec reference

| Key                    | Type   | Default value | Description                                                                                                                                |
|------------------------|--------|---------------|--------------------------------------------------------------------------------------------------------------------------------------------|
| osImage                | string | empty         | The fully qualified image to upgrade nodes to. This value has priority over `managedOSVersionName` if both are configured.                 |
| managedOSVersionName   | string | empty         | The name of a `ManagedOSVersion` to upgrade nodes to.                                                                                      |
| cloudConfig            | object | null          | A cloud-init or yip config to apply to the nodes during upgrades. See [reference](#cloudconfig).                                                  |
| nodeSelector           | object | null          | This selector can be used to target specific nodes within the `clusterTargets`. See [reference](#nodeselector).                            |
| concurrency            | int    | 1             | How many nodes within the same cluster should be upgraded at the same time.                                                                |
| cordon                 | bool   | true          | Set this to true if the nodes should be cordoned before applying the upgrade. Ineffective when `drain` is also configured.                 |
| drain                  | object | See ref       | Configure if and how nodes should be drained before applying the upgrade. See [reference](#drain).                                         |
| prepare                | object | null          | The prepare init container, if specified, is run before cordon/drain which is run before the upgrade container. See [reference](#prepare). |
| upgradeContainer       | object | null          | The upgrade container that will run the upgrade on the nodes. See [reference](#upgradecontainer).                                          |
| clusterRolloutStrategy | object | null          | RolloverStrategy controls the rollout of the upgrade bundle across clusters. See [reference](#clusterrolloutstrategy).                     |
| clusterTargets         | list   | null          | Declares clusters to deploy the upgrade plan to. See [reference](#clustertargets).                                                         |

#### cloudConfig

This describes a cloud-init or yip config that will be copied to each upgraded node to the `/oem/90_operator.yaml` path.  
This config will be applied by the system after reboot.  
For more information and examples, see the `MachineRegistration` `spec.config.cloud-config` [reference](./cloud-config-reference.md).  

#### nodeSelector

This [Label Selector](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#label-selectors) can be used to restrict the upgrades to only a certain set of nodes within the targeted Clusters.

<details>
  <summary>Example</summary>

  ```yaml showLineNumbers
  nodeSelector:
    matchExpressions:
      - {key: my-node/label, operator: Exists}
  ```
  
</details>

#### drain

Configure if and how nodes should be drained.  
To disable drain during upgrades you can configure this field to `null`.  
Drain is enabled by default.  

The drain settings directly translates to the [kubectl drain](https://kubernetes.io/docs/reference/kubectl/generated/kubectl_drain/) command being executed on the node before upgrade.  

| Key                      | Type             | Default value | Description                                                                                                                                |
|--------------------------|------------------|---------------|--------------------------------------------------------------------------------------------------------------------------------------------|
| timeout                  | time.Duration    | null          | The length of time to wait before giving up draining a node, zero means infinite.                                                          |
| gracePeriod              | int              | null          | Period of time in seconds given to each pod to terminate gracefully. If negative, the default value specified in the pod will be used.     |
| deleteEmptydirData       | bool             | true          | Continue even if there are pods using emptyDir (local data that will be deleted when the node is drained).                                 |
| ignoreDaemonSets         | bool             | true          | Ignore DaemonSet-managed pods.                                                                                                             |
| force                    | bool             | true          | Continue even if there are pods that do not declare a controller.                                                                          |
| disableEviction          | bool             | false         | Force drain to use delete, even if eviction is supported. This will bypass checking PodDisruptionBudgets, use with caution.                |
| skipWaitForDeleteTimeout | int              | 60            | If pod DeletionTimestamp older than N seconds, skip waiting for the pod. Seconds must be greater than 0 to skip.                           |
| podSelector              | label selector   | null          | Label selector to filter pods on the node. Only selected pods will be evicted.                                                             |

#### prepare

Defines a `prepare` Init container that is ran before the `upgrade` container executing the upgrade job on a node.  
The keys directly translate to the [container](https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/pod-v1/#Container) specification.  
Note that the node filesystem is mounted at `/host` inside the container.  

| Key                      | Type   | Default value | Description                                                          |
|--------------------------|--------|---------------|----------------------------------------------------------------------|
| image                    | string | empty         | Container image name.                                                |
| command                  | list   | empty         | Entrypoint array.                                                    |
| args                     | list   | empty         | Arguments to the entrypoint.                                         |
| env                      | list   | empty         | List of environment variables to set in the container.               |
| envFrom                  | list   | empty         | List of sources to populate environment variables in the container.  |
| volumes                  | list   | empty         | List of `hostPath` volumes. See [reference](#preparevolumes).        |
| securityContext          | object | null          | The security options the ephemeral container should be run with.     |

##### prepare.volumes

Each volume definition will translate to a [hostPath](https://kubernetes.io/docs/concepts/storage/volumes/#hostpath-volume-types) volume (`source`) which will be mounted in the container (`destination`).  
Note that by default the host root filesystem `/` will always be mounted at `/host`.  

| Key          | Type   | Default value | Description                 |
|--------------|--------|---------------|-----------------------------|
| name         | string | empty         | Volume name.                |
| source       | string | empty         | HostPath volume path.       |
| destination  | string | empty         | HostPath volume mount path. |

<details>
  <summary>Example</summary>

  ```yaml showLineNumbers
  volumes:
    - name: my-custom-volume
      source: /foo
      destination: /foo
  ```
  
</details>

#### upgradeContainer

Defines the `upgrade` container executing the upgrade job on a node.  
The keys directly translate to the [container](https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/pod-v1/#Container) specification.  
Note that the node filesystem is mounted at `/host` inside the container.  

:::warning warning
When using any Elemental or [Elemental based image](./custom-images.md) you are expected to only edit the `env` key to optionally set the `FORCE`, `UPGRADE_RECOVERY`, or `UPGRADE_RECOVERY_ONLY` variables.  
For more info you can read the [upgrade](./upgrade.md#upgrade-via-command-line-interface) documentation.  
Any other change to the `upgradeContainer` may result in issues during upgrades.  
:::

#### clusterRolloutStrategy

This controls the rollout of the bundle across clusters.  
For more information you can read the [reference documentation](https://fleet.rancher.io/0.9/ref-crds#rolloutstrategy).  

#### clusterTargets

Select Clusters to be targeted for the OS image upgrade.  
For more information you can read the [reference documentation](https://fleet.rancher.io/0.9/ref-crds#bundletarget).  

<details>
  <summary>Example</summary>

  ```yaml showLineNumbers
  clusterTargets:
    - clusterName: volcano
  ```
  
</details>


---

## Article: managedosversion-reference.md

---
sidebar_label: ManagedOSVersion reference
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/managedosversion-reference"/>
</head>

# ManagedOSVersion reference

The `ManagedOSVersion` resource is responsible of defining an OS image version that can be used with `SeedImage` or `ManagedOSImage` resources.

There are several keys that can be configured under a `ManagedOSVersion` resource spec.

```yaml title="managedosversion-example.yaml" showLineNumbers
apiVersion: elemental.cattle.io/v1beta1
kind: ManagedOSVersion
metadata:
  labels:
    elemental.cattle.io/channel: elemental-channel
  name: v2.0.2
  namespace: fleet-default
spec:
  metadata:
    displayName: SLE Micro
    upgradeImage: registry.suse.com/suse/sle-micro/5.5:2.0.2
  type: container
  version: v2.0.2
```

#### ManagedOSVersionSpec reference

| Key              | Type   | Default value | Description                                                                  |
|------------------|--------|---------------|------------------------------------------------------------------------------|
| metadata         | object | null          | This defines some data about the OS image. See [reference](#metadata)        |
| minVersion       | string | null          | Not used                                                                     |
| type             | string | empty         | Defines the OS image type, could be `container` or `iso`                     |
| version          | string | empty         | OS image version                                                             |
| upgradeContainer | object | null          | An upgrade container that can be defined. See [reference](#upgradecontainer) |

<details>
  <summary>ISO image example</summary>

  ```yaml showLineNumbers
  metadata:
    displayName: SLE Micro ISO x86_64
    uri: registry.suse.com/suse/sl-micro/6.0/baremetal-iso-image:2.2.0
  type: iso
  version: v2.2.0
  ```
  
</details>

<details>
  <summary>Container image example</summary>

  ```yaml showLineNumbers
  metadata:
    displayName: SLE Micro
    upgradeImage: registry.suse.com/suse/sl-micro/6.0/baremetal-os-container:2.2.0
  type: container
  version: v2.2.0
  ```
  
</details>

#### metadata

This describes the needed information to define an OS image in Elemental.

If `type` is set to `container`:
| Key          | Type   | Default value | Description                                                         |
|--------------|--------|---------------|---------------------------------------------------------------------|
| displayName  | string | empty         | OS image name as seen in Rancher UI                                 |
| upgradeImage | string | empty         | Fully qualified Container image (OCI reference or HTTP URI)        |
| platforms    | list   | empty         | The supported platforms (`linux/x86_64`, `linux/aarch64`, or both). |

If `type` is set to `iso`:
| Key          | Type   | Default value | Description                                                         |
|--------------|--------|---------------|---------------------------------------------------------------------|
| displayName  | string | empty         | OS image name as seen in Rancher UI                                 |
| uri          | string | empty         | Fully qualified ISO image                                           |
| platforms    | list   | empty         | The supported platforms (`linux/x86_64`, `linux/aarch64`, or both). |

#### upgradeContainer

This allows to overwrite the default `upgrade` field of System Upgrade Controller plans (see [upgrade components](/upgrade-lifecycle.md#components)) based on this ManagedOSVersion.
These keys are translated by the System Upgrade Controller to a Kubernetes [container](https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/pod-v1/#Container) specification.
This is the container responsible of running an OS upgrade.


---

## Article: managedosversionchannel-reference.md

---
sidebar_label: ManagedOSVersionChannel reference
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/managedosversionchannel-reference"/>
</head>

# ManagedOSVersionChannel reference

The `ManagedOSVersionChannel` resource is responsible of defining OS image channels which embed
information about ready to use OS and ISO images.

The `ManagedOSVersionChannel` is digested by the [elemental-operator](architecture-components.md#elemental-operator-daemon)
to produce a set of [ManagedOSVersion](managedosversion-reference.md) resources, which are used
for building installation images and to perform OS upgrades.

See [Operator/Channels](channels.md) for more details.

There are several keys that can be configured under a `ManagedOSVersionChannel` resource spec.

```yaml title="managedosversionchannel-example.yaml" showLineNumbers
apiVersion: elemental.cattle.io/v1beta1
kind: ManagedOSVersionChannel
metadata:
  name: elemental-channel
  namespace: fleet-default
spec:
  options:
    image: registry.suse.com/rancher/elemental-channel:1.4.2
  registry: priv-registry.local
  syncInterval: 1h
  type: custom
```

#### ManagedOSVersionChannelSpec reference

| Key                          | Type   | Default value | Description                                                                             |
|------------------------------|--------|---------------|-----------------------------------------------------------------------------------------|
| deleteNoLongerInSyncVersions | bool   | false         | Automatically delete deprecated OS versions that are no longer included in the channel  |
| enabled                      | bool   | true          | Enables this channel. Allowing syncing of OS versions.                                  |
| options                      | object | null          | Defines the optional information that can be added in an OS channel                    |
| registry                     | string | empty         | Registry prepended to all the embedded URIs during ManagedOSVersion resource generation |
| syncInterval                 | string | 1h            | Defines when to sync the OS channel                                                     |
| type                         | string | empty         | Defines the channel type, only `custom` is supported now                                |
| upgradeContainer             | object | null          | An upgrade container that can be defined. See [reference](#upgradecontainer)            |

#### registry
When set, the value is prepended to all the embedded OS and ISO URLs found in the channel to form the actual [ManagedOSVersion](managedosversion-reference.md)
resources.
This makes the airgap setup easier, as the channel image can be used as is (while the actual OS and ISO images should still be extracted and loaded
on the private registry).

#### upgradeContainer

This allows to overwrite the default `upgrade` field of System Upgrade Controller plans (see [upgrade components](/upgrade-lifecycle.md#components)) based on this ManagedOSVersion.
These keys are translated by the System Upgrade Controller to a Kubernetes [container](https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/pod-v1/#Container) specification.
This is the container responsible of running an OS upgrade.


---

## Article: networking-static.md

---
sidebar_label: Static Configuration
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/networking"/>
</head>

import YipNmcStaticConfig from "!!raw-loader!@site/examples/network/yip-nmc-static-config.yaml"

:::warning warning
Declarative Networking is in **technology preview** state with limited support.  
As such it is not recommended for production environments.
:::

## Static Network with nm-configurator

The `nm-configurator` [per node configuration](https://github.com/suse-edge/nm-configurator?tab=readme-ov-file#per-node-configurations) can be used to statically assign IP addresses to individual machines, based on the NIC's MAC addresses.  

This solution does not require a remote IPAM provider, but requires the user to maintain mapping between known MAC addresses and IP Addresses.  

In this example, we are going to customize an Elemental image, and include a [yip config](./cloud-config-reference.md#configuration-syntax) that will apply the static network config early at boot.

The following configuration covers 2 nodes with a very basic network setup.  

<CodeBlock language="yaml" title="99_static_network_config.yaml" showLineNumbers>{YipNmcStaticConfig}</CodeBlock>

## Include the static network configuration in a custom OS image

We can extend an Elemental image to include the static network configuration in `/system/oem`.  
Any Elemental powered OS, where [Elemental Toolkit](https://github.com/rancher/elemental-toolkit) is running, will evaluate any config in this directory when executing any stage.  
Additionally we are going to customize the image to install the required `nmc` binary.  

```docker showLineNumbers
# The version of Elemental to modify
FROM registry.suse.com/suse/sl-micro/6.0/baremetal-os-container:latest

# Install the static network config
COPY 99_static_network_config.yaml /system/oem/99_static_network_config.yaml

# Install nmc
RUN curl -LO https://github.com/suse-edge/nm-configurator/releases/download/v0.3.1/nmc-linux-x86_64 && \
    install -o root -g root -m 0755 nmc-linux-x86_64 /usr/sbin/nmc

# IMPORTANT: /etc/os-release is used for versioning/upgrade.
ARG IMAGE_REPO=norepo
ARG IMAGE_TAG=latest
RUN \
    sed -i -e "s|^IMAGE_REPO=.*|IMAGE_REPO=\"${IMAGE_REPO}\"|g" /etc/os-release && \
    sed -i -e "s|^IMAGE_TAG=.*|IMAGE_TAG=\"${IMAGE_TAG}\"|g" /etc/os-release && \
    sed -i -e "s|^IMAGE=.*|IMAGE=\"${IMAGE_REPO}:${IMAGE_TAG}\"|g" /etc/os-release

# IMPORTANT: it is good practice to recreate the initrd and re-apply `elemental-init`
RUN elemental init --force elemental-rootfs,grub-config,dracut-config,cloud-config-essentials,elemental-setup
```

The OS container can now be built and pushed to your registry:  

```bash showLineNumbers
docker build --build-arg IMAGE_REPO=myrepo/static-network-os \
             --build-arg IMAGE_TAG=v1.1.1 \
             -t myrepo/static-network-os:v1.1.1 .
docker push myrepo/static-network-os:v1.1.1
```

Note that since the static network config is included in the image, this requires you to build and maintain different images with different configurations, most likely one per nodepool or one per cluster for example.  
Custom OS images need to be maintained to create the initial ISO bootable image (via a [SeedImage](./seedimage-reference.md)), but also when [upgrading](./upgrade.md) the machines.  
The desired static network config **must** be present on the OS images used with the [ManagedOSImage](./managedosimage-reference.md), otherwise the config will be missing when booting from the upgraded system.  

The custom OS image can also be used as it is to build a bootable raw disk image:  

```yaml showLineNumbers
apiVersion: elemental.cattle.io/v1beta1
kind: SeedImage
metadata:
  name: my-raw-image
  namespace: fleet-default
spec:
  type: raw
  baseImage: myrepo/static-network-os:v1.1.1
  registrationRef:
    apiVersion: elemental.cattle.io/v1beta1
    kind: MachineRegistration
    name: my-registration
    namespace: fleet-default
```

## Create a bootable ISO

You can now [build an ISO container](./custom-images.md#create-a-custom-bootable-installation-iso) from this OS container image. For more information on how to customize Elemental images, please refer to the [documentation](./custom-images.md).  

```docker showLineNumbers
FROM myrepo/static-network-os:v1.1.1 AS os
FROM myrepo/static-network-os:v1.1.1 AS builder

WORKDIR /iso
COPY --from=os / rootfs

# work around buildah issue: https://github.com/containers/buildah/issues/4242
RUN rm -f rootfs/etc/resolv.conf

RUN elemental build-iso \
        dir:rootfs \
        --bootloader-in-rootfs \
        --squash-no-compression \
        -o /output -n "elemental"

FROM busybox
COPY --from=builder /output /elemental-iso

ENTRYPOINT ["busybox", "sh", "-c"]
```

```bash showLineNumbers
docker build -t myrepo/static-network-iso:v1.1.1 .
docker push myrepo/static-network-iso:v1.1.1
```

Once the ISO container is published on your registry, you can refer to it in the [SeedImage](./seedimage-reference.md) like any other Elemental distributed ISO image.  

```yaml
apiVersion: elemental.cattle.io/v1beta1
kind: SeedImage
metadata:
  name: my-iso
  namespace: fleet-default
spec:
  type: iso
  baseImage: myrepo/static-network-iso:v1.1.1
  registrationRef:
    apiVersion: elemental.cattle.io/v1beta1
    kind: MachineRegistration
    name: my-registration
    namespace: fleet-default
```

Note that the static network config will now be evaluated when the installation media boots, then it will be installed on the system as part of the base image.  


---

## Article: networking.md

---
sidebar_label: IPAM Driven Networking
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/networking"/>
</head>

import RegistrationWithNetwork from "!!raw-loader!@site/examples/network/machineregistration.yaml"
import RegistrationWithNetworkNmc from "!!raw-loader!@site/examples/network/machineregistration-nmc.yaml"
import RegistrationWithNetworkNmstate from "!!raw-loader!@site/examples/network/machineregistration-nmstate.yaml"
import RegistrationWithNetworkNmconnections from "!!raw-loader!@site/examples/network/machineregistration-nmconnections.yaml"
import YipNmcStaticConfig from "!!raw-loader!@site/examples/network/yip-nmc-static-config.yaml"

:::warning warning
Declarative Networking is in **technology preview** state with limited support.  
As such it is not recommended for production environments.
:::

## Network configuration with Elemental

The [MachineRegistration](machineregistration-reference) supports Declarative Networking and integration with [CAPI IPAM Providers](https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20220125-ipam-integration.md#ipam-provider).  

### Prerequisites

- A DHCP server is still required for the first boot registration and reset of machines. For this reason Lease Time can be kept minimal, as for the entire lifecycle of the machine, the IPAM driven IP Addresses will be used.  

- An IPAM Provider of your choice is installed on the Rancher management cluster.  
  For example the [InCluster IPAM Provider](https://github.com/kubernetes-sigs/cluster-api-ipam-provider-in-cluster).

- [NetworkManager](https://networkmanager.dev) needs to be installed on OS images and it can be directly configured using the `nmconnections` network [configurator](#configurators).  
  Already included in Elemental provided images.  
  
- (optionally) [nmc](https://github.com/suse-edge/nm-configurator/releases) can be used with `nmc` network [configurator](#configurators).  

- (optionally) [nmstatectl](https://github.com/nmstate/nmstate/releases) can be used with `nmstate` network [configurator](#configurators).

#### Installing nmc or nmstatectl on OS images

When using the `nmc` or `nmstate` configurators, the [nmc](https://github.com/suse-edge/nm-configurator/releases) or [nmstatectl](https://github.com/nmstate/nmstate/releases) tools need to be installed on the machine.  

Currently this can be achieved by [customizing an Elemental OS image](./custom-images.md#remastering-an-os-image-with-a-custom-dockerfile) with a custom command:

<Tabs>

<TabItem value="nmc" label="nmc" default>

```yaml
# Install nmc
RUN curl -LO https://github.com/suse-edge/nm-configurator/releases/download/v0.3.3/nmc-linux-x86_64 && \
    install -o root -g root -m 0755 nmc-linux-x86_64 /usr/sbin/nmc
```

</TabItem>

<TabItem value="nmstatectl" label="nmstatectl" default>

```yaml
# Install nmstatectl
RUN curl -LO https://github.com/nmstate/nmstate/releases/download/v2.2.40/nmstatectl-linux-x64.zip && \
    unzip nmstatectl-linux-x64.zip && \
    chmod +x nmstatectl && \
    mv ./nmstatectl /usr/sbin/nmstatectl && \
    rm nmstatectl-linux-x64.zip
```

</TabItem>

</Tabs>

### How to install the CAPI IPAM Provider

The recommended way to install any CAPI Provider into Rancher is to use [Rancher Turtles](https://turtles.docs.rancher.com).  
Rancher Turtles will allow the user to install and manage the lifecycle of any CAPI Provider.  
To install it on your system please follow the [documentation](https://turtles.docs.rancher.com/getting-started/install-rancher-turtles/using_helm).  

Once Rancher Turtles is installed, installing an IPAM CAPI Provider, for example the [InCluster IPAM Provider](https://github.com/kubernetes-sigs/cluster-api-ipam-provider-in-cluster), can be accomplished applying the following resource:

```yaml
kind: CAPIProvider
metadata:
  name: in-cluster
  namespace: default
spec:
  name: in-cluster
  type: ipam
  fetchConfig:
    url: "https://github.com/kubernetes-sigs/cluster-api-ipam-provider-in-cluster/releases"
  version: v0.1.0
```

#### Without Rancher Turtles

An alternative option to install a CAPI IPAM Provider is to directly apply the manifest in the Rancher cluster.  
Note that this solution may eventually lead to conflicts with the applied CRDs and resources, as they need to be applied and maintained manually.  

1. The `ipaddresses.ipam.cluster.x-k8s.io` and `ipaddressclaims.ipam.cluster.x-k8s.io` CRDs must be installed on the Rancher management cluster:  

    ```bash
    kubectl apply -f https://raw.githubusercontent.com/kubernetes-sigs/cluster-api/main/config/crd/bases/ipam.cluster.x-k8s.io_ipaddressclaims.yaml
    kubectl apply -f https://raw.githubusercontent.com/kubernetes-sigs/cluster-api/main/config/crd/bases/ipam.cluster.x-k8s.io_ipaddresses.yaml
    ```

    :::info info
    These CRDs are expected to eventually be part of Rancher, not requiring manual installation.  
    See: https://github.com/rancher/rancher/issues/46385
    :::

1. Install the [InCluster IPAM Provider](https://github.com/kubernetes-sigs/cluster-api-ipam-provider-in-cluster) from the released manifest:  

    ```bash
    kubectl apply -f https://github.com/kubernetes-sigs/cluster-api-ipam-provider-in-cluster/releases/download/v0.1.0/ipam-components.yaml
    ```

### Configuring Network

The `network` section of the `MachineRegistration` allows users to define:

1. A map of IPPool references.
1. A network config template (in this case `nmc` configurator is in use).  

For example:

<CodeBlock language="yaml" title="example MachineRegistration using Declarative Networking" showLineNumbers>{RegistrationWithNetwork}</CodeBlock>

Here we can observe that one `InClusterIPPool` has been defined, since we are using the [InCluster IPAM Provider](https://github.com/kubernetes-sigs/cluster-api-ipam-provider-in-cluster) for this example.  

Next we are going to reference this IPPool in the `MachineRegistration`. The key for this reference is `inventory-ip`, and we are only going to need one IP per registered Machine. If your machine has more than one NIC, you can define more references, and use different IPPools as well, for example:  

```yaml
ipAddresses:
  main-nic-ip:
    apiGroup: ipam.cluster.x-k8s.io
    kind: InClusterIPPool
    name: elemental-inventory-pool
  secondary-nic-ip:
    apiGroup: ipam.cluster.x-k8s.io
    kind: InClusterIPPool
    name: elemental-inventory-pool
  private-nic-ip:
    apiGroup: ipam.cluster.x-k8s.io
    kind: InClusterIPPool
    name: elemental-private-pool
```

Each defined IPPool reference key can be used for the network config template:

```yaml
config:
  dns-resolver:
    config:
      server:
      - 192.168.122.1
      search: []
  routes:
    config:
    - destination: 0.0.0.0/0
      next-hop-interface: eth0
      next-hop-address: 192.168.122.1
      metric: 150
      table-id: 254
  interfaces:
    - name: eth0
      type: ethernet
      description: Main-NIC
      state: up
      ipv4:
        enabled: true
        dhcp: false
        address:
        - ip: "{inventory-ip}"
          prefix-length: 24
      ipv6:
        enabled: false
```

The snippet above is almost 1:1 [nm-configurator syntax](https://github.com/suse-edge/nm-configurator?tab=readme-ov-file#unified-configurations), with the only exception of the `{inventory-ip}` placeholder.  
During the installation or reset phases of Elemental machines, the `elemental-operator` will claim one IP Address from the referenced IP Pool, and substitute the `{inventory-ip}` placeholder with a real IP Address.  

### Claimed IPAddresses

The `IPAddressClaim` will follow the entire lifecycle of the `MachineInventory`, ensuring that each registered machine will be assigned unique IPs.  
Each claim is named after the `MachineInventory` that uses it, as `$MachineInventoryName-$IPPoolRefKey`, for example:  

```yaml
apiVersion: ipam.cluster.x-k8s.io/v1beta1
kind: IPAddressClaim
metadata:
  finalizers:
    - ipam.cluster.x-k8s.io/ReleaseAddress
  name: m-e5331e3b-1e1b-4ce7-b080-235ed9a6d07c-inventory-ip
  namespace: fleet-default
  ownerReferences:
    - apiVersion: elemental.cattle.io/v1beta1
      kind: MachineInventory
      name: m-e5331e3b-1e1b-4ce7-b080-235ed9a6d07c
spec:
  poolRef:
    apiGroup: ipam.cluster.x-k8s.io
    kind: InClusterIPPool
    name: elemental-inventory-pool
status:
  addressRef:
    name: m-e5331e3b-1e1b-4ce7-b080-235ed9a6d07c-inventory-ip
```

Whenever a `MachineInventory` is deleted, the default (DHCP) network configuration will be restored, deleting any network profile on the machine and restarting the network stack. Finally the assigned IPs will be released.  

For more information and details on how troubleshoot issues, please consult the [documentation](./troubleshooting-network.md).

### Configurators

On the Elemental machine, `elemental-register` can configure the `NetworkManager` in different ways.  
The configurator in use is defined in the [MachineRegistration.spec.network](./machineregistration-reference.md#confignetwork):  

- [nmc](https://github.com/suse-edge/nm-configurator)
- [nmstate](https://nmstate.io/)
- [nmconnections](https://networkmanager.pages.freedesktop.org/NetworkManager/NetworkManager/nm-settings-keyfile.html)

<Tabs>

<TabItem value="nmc" label="nmc" default>

The `nmc` configurator uses the [nm-configurator unified syntax](https://github.com/suse-edge/nm-configurator?tab=readme-ov-file#unified-configurations) to generate NetworkManager's connection files.  

<details>
  <summary>example MachineRegistration using nmc configurator</summary>
<CodeBlock language="yaml" showLineNumbers>{RegistrationWithNetworkNmc}</CodeBlock>
</details>

</TabItem>

<TabItem value="nmstate" label="nmstate" default>

The `nmstate` configurator uses [nmstate syntax](https://nmstate.io/examples.html) to generate NetworkManager's connection files.  
Note that [nmstatectl](https://github.com/nmstate/nmstate/releases) needs to be installed on the Elemental system to use this configurator. This is not included by default in Elemental images, but can be installed when building a [custom image](./custom-images.md).

<details>
  <summary>example MachineRegistration using nmstate configurator</summary>
<CodeBlock language="yaml" showLineNumbers>{RegistrationWithNetworkNmstate}</CodeBlock>
</details>

</TabItem>

<TabItem value="nmconnections" label="nmconnections" default>

The `nmconnections` configurator is the simplest option available and allows the user to directly write `nmconnection` files.  
Defining these files for complex network setups may be challenging, but it's always possible to use [nmcli](https://networkmanager.dev/docs/api/latest/nmcli.html), or even [nmstate](https://nmstate.io), or [nm-configurator](https://github.com/suse-edge/nm-configurator), and use the generated `nmconnection` files as a template.  
This configurator only needs `NetworkManager`, without any extra dependency.  

<details>
  <summary>example MachineRegistration using nmconnections configurator</summary>
<CodeBlock language="yaml" showLineNumbers>{RegistrationWithNetworkNmconnections}</CodeBlock>
</details>

</TabItem>

</Tabs>


---

## Article: ntp.md

---
sidebar_label: Configure NTP
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/ntp"/>
</head>


## NTP configuration with Elemental

### Overview
The default OS channel shipped with Elemental provides NTP support via `systemd-timesyncd`.

This page covers configuring `systemd-timesyncd` with the provided SLE Micro images, which is
pre-configured with some default (fallback) NTP servers: (_[0-3].suse.pool.ntp.org_).

The easiest way to configure a specific NTP server, is to drop a configuration file in the
`/etc/systemd/timesyncd.conf.d` directory.
The directory and the configuration file should be accessible by the `systemd-timesync` user.

### Configure a static NTP server
The NTP configuration can be provided via a cloud-config snippet added to the MachineRegistration
configuration.

We will need to:
* ensure the `timesyncd.conf.d` directory can be read by the `systemd-timesync` user
* write the custom config file in the `timesyncd.conf.d` directory
* restart the `systemd-timesyncd` service to use the new configuration

As an example, let's see how to configure `ntp.ripe.net` as the primary NTP server (lines 6-14):
```yaml showLineNumbers
config:
  cloud-config:
    users:
      - name: root
        passwd: root
    write_files:
      - content: |
          [Time]
          NTP=ntp.ripe.net
        path: /etc/systemd/timesyncd.conf.d/custom-ntp.conf
        permissions: 644
    runcmd:
      - chmod 755 /etc/systemd/timesyncd.conf.d
      - systemctl restart systemd-timesyncd
  elemental:
    install:
      device: /dev/vda
      reboot: true
machineInventoryLabels:
  element: fire

```

### Configure NTP from DHCP
In order to get the NTP server from the network via the NTP DHCP option, we need
a NetworkManager dispatcher script to reconfigure dynamically the `systemd-timesync` service when
needed.

We will have both to:
* provide the dispatcher script which creates and deletes the systemd-timesyncd config files
* enable the NetworkManager-dispatcher service

See lines 6-34 in the following MachineRegistration configuration example:

```yaml showLineNumbers
config:
  cloud-config:
    users:
      - name: root
        passwd: root
    write_files:
      - content: |
          #! /usr/bin/bash

          [ -n "$CONNECTION_UUID" ] || exit

          INTERFACE=$1
          ACTION=$2

          case $ACTION in
              up | dhcp4-change | dhcp6-change)
                  [ -n "$DHCP4_NTP_SERVERS" ] || exit
                  mkdir -p /etc/systemd/timesyncd.conf.d/
                  cat<<EOF > /etc/systemd/timesyncd.conf.d/$CONNECTION_UUID.conf
          [Time]
          NTP=$DHCP4_NTP_SERVERS
          RootDistanceMaxSec=15
          EOF
                  systemctl restart systemd-timesyncd
                  ;;
              down)
                  rm -f /etc/systemd/timesyncd.conf.d/$CONNECTION_UUID.conf
                  systemctl restart systemd-timesyncd
                  ;;
          esac
        path: /etc/NetworkManager/dispatcher.d/10-update-timesyncd
        permissions: 700
    runcmd:
      - systemctl enable NetworkManager-dispatcher
  elemental:
    install:
      device: /dev/vda
      reboot: true
machineInventoryLabels:
  element: fire

```


---

## Article: quickstart-cli.md

---
sidebar_label: Elemental the command line way
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/quickstart-cli"/>
</head>

import Cluster from "!!raw-loader!@site/examples/quickstart/cluster.yaml"
import Registration from "!!raw-loader!@site/examples/quickstart/registration.yaml"
import RegistrationRPi from "!!raw-loader!@site/examples/quickstart/rpi-registration.yaml"
import Selector from "!!raw-loader!@site/examples/quickstart/selector.yaml"
import Prereqs from './partials/_quickstart-prereqs.md'
import Operator from './partials/_elemental-operator-install.md'
import SeedImage from "!!raw-loader!@site/examples/quickstart/seedimage.yaml"

# Elemental the command line way

Follow this guide to have an auto-deployed cluster via rke2/k3s and managed by Rancher 
with the only help of an Elemental ISO.

<Prereqs />

<Operator />

## Prepare your kubernetes resources

Node deployment starts with a `MachineRegistration`, identifying a set of machines sharing the same configuration (disk drives, network, etc.).

The `MachineRegistration` is needed to perform the deployment of the Elemental OS on the target hosts. When booting up, each host registers to the Elemental Operator which tracks the new host with a `MachineInventory` resource.

Then it continues with having a Cluster resource that uses a `MachineInventorySelectorTemplate` to know which machines are for that cluster.

This selector is a simple matcher based on labels set in the `MachineInventory`, so if your selector is matching on the label `cluster-id` with a value `cluster-id-val`
and your `MachineInventory` has that same `cluster-id`:`cluster-id-val` label, it will match and be bootstrapped as part of the cluster.

In this quickstart we are going to deploy the resources to provision a cluster named *volcano* that will match on `MachineInventory`s with the label *element*:*fire*.

<Tabs>
<TabItem value="manualYaml" label="Manually creating the resource yamls" default>

You will need to create the following files:

<CodeBlock language="yaml" title="selector.yaml" showLineNumbers>{Selector}</CodeBlock>

As you can see this is a very simple selector that looks for `MachineInventory`s having a label with the key `element` and the value `fire`.

<CodeBlock language="yaml" title="cluster.yaml" showLineNumbers>{Cluster}</CodeBlock>

As you can see the `machineConfigRef` is of kind `MachineInventorySelectorTemplate` with the name `fire-machine-selector`: it matches the selector we created.

You can get more information about cluster options like [`machineGlobalConfig`](https://ranchermanager.docs.rancher.com/reference-guides/cluster-configuration/rancher-server-configuration/rke2-cluster-configuration#machineglobalconfig) or [`machineSelectorConfig`](https://ranchermanager.docs.rancher.com/reference-guides/cluster-configuration/rancher-server-configuration/rke2-cluster-configuration#machineselectorconfig) directly in the [Rancher Manager documentation](https://ranchermanager.docs.rancher.com).

<Tabs>
<TabItem value="normalRegistration" label="Registration" default>
<CodeBlock language="yaml" title="registration.yaml" showLineNumbers>{Registration}</CodeBlock>
</TabItem>
<TabItem value="rpiRegistration" label="Registration for Raspberry Pi" default>
<CodeBlock language="yaml" title="rpi-registration.yaml" showLineNumbers>{RegistrationRPi}</CodeBlock>

For deployment on Raspberry Pi, you need to enable emulated TPM 
(except you have [a hardware TPM for Raspberry Pi](https://thepihut.com/products/letstrust-tpm-for-raspberry-pi)).
You also need to disable writing to the EFI store (since Raspberry Pi doesn't have one) via `disable-boot-entry: true`.

</TabItem>
</Tabs>

The `MachineRegistration` defines the registration and installation configuration. Once created, the Elemental operator exposes a unique URL to be used with the `elemental-register` binary to reach out to the management cluster and register the machine during installation: if the registration is successful, the operator creates a `MachineInventory` tracking the machine, which can be used to provision the machine as a node of our cluster.
We define the label matching our selector here, although it can also be added later to the created `MachineInventory`s.



:::warning warning
Make sure to modify the registration.yaml above to set the proper install device to point to a valid device based on your node configuration (i.e. /dev/sda, /dev/vda, /dev/nvme0, etc...).

The SD-card on a Raspberry Pi is usually `/dev/mmcblk0`.
:::

<Tabs>
<TabItem value="seedImagex86" label="Seed Image (x86_64)" default>
<CodeBlock language="yaml" title="seedimage.yaml" showLineNumbers>{SeedImage}</CodeBlock>

The `SeedImage` is required to generate the *seed image* (like a bootable ISO) that will boot and start the Elemental provisioning on the target machines.

Now that we have defined all the configuration files let's apply them to create the proper resources in Kubernetes:

```shell showLineNumbers
kubectl apply -f selector.yaml 
kubectl apply -f cluster.yaml 
kubectl apply -f registration.yaml
kubectl apply -f seedimage.yaml
```

</TabItem>
<TabItem value="seedImagerpi" label="Seed Image for Raspberry Pi" default>

The `SeedImage` resource, which automates the creation of an Elemental bootable image (the *seed image*), does not support Raspberry Pi ISOs yet (click [here](raspi-disk.md) for a guide to build a raw disk image).

We will generate a *seed image* manually in the [next section](quickstart-cli.md#preparing-the-installation-seed-image).

Now that we have defined all the configuration files let's apply them to create the proper resources in Kubernetes:

```shell showLineNumbers
kubectl apply -f selector.yaml 
kubectl apply -f cluster.yaml 
kubectl apply -f registration.yaml
```

</TabItem>
</Tabs>

</TabItem>
<TabItem value="repofiles" label="Using quickstart files from Elemental docs repo directly">

You can directly apply the quickstart example resource files from the [Elemental docs repository](https://github.com/rancher/elemental-docs).

:::warning warning
The quickstart example resource files assume the default storage of the target host to be mapped to the `/dev/sda`.
If your host storage device file is different, you have to change the registration.yaml file before applying it, changing the `config.elemental.install.device` accordingly.
:::

```bash showLineNumbers
kubectl apply -f https://raw.githubusercontent.com/rancher/elemental-docs/main/examples/quickstart/selector.yaml
kubectl apply -f https://raw.githubusercontent.com/rancher/elemental-docs/main/examples/quickstart/cluster.yaml
kubectl apply -f https://raw.githubusercontent.com/rancher/elemental-docs/main/examples/quickstart/registration.yaml
kubectl apply -f https://raw.githubusercontent.com/rancher/elemental-docs/main/examples/quickstart/seedimage.yaml (not for aarch64 yet)
```

</TabItem>
</Tabs>

## Preparing the installation (seed) image


This is the last step: you need an Elemental seed image that includes the initial registration config, so it can be auto registered, installed and fully deployed as part of your cluster.

:::note note
The initial registration config file is generated when you create a `Machine Registration`.

You can download it with:

```shell
wget --no-check-certificate `kubectl get machineregistration -n fleet-default fire-nodes -o jsonpath="{.status.registrationURL}"` -O initial-registration.yaml
```
:::

The contents of the registration config file are nothing more than the registration URL that the node needs to register, the proper server certificate and few options for the registration process.

Once generated, a seed image can be used to provision any number of machines.

<Tabs>
<TabItem value="download" label="Downloading the quickstart ISO">

The seed image created by the `SeedImage` resource above can be downloaded as an ISO via the following script:

```shell showLineNumbers
kubectl wait --for=condition=ready pod -n fleet-default fire-img
wget --no-check-certificate `kubectl get seedimage -n fleet-default fire-img -o jsonpath="{.status.downloadURL}"` -O elemental.x86_64.iso
```

The first command waits for the ISO to be built and ready, the second one downloads it in the current directory with the name `elemental-x86_64.iso`.

</TabItem>
<TabItem value="manual_raw" label="Preparing the seed image (aarch64) manually">

Elemental's support for Raspberry Pi is primarily for demonstration purposes at this point. Therefore the installation process is modelled similar to x86-64. You boot from a seed image (an USB stick in this case) and install to a storage medium (SD-card for Raspberry Pi).

#### Retrieving the prebuilt seed image

```shell showLineNumbers
wget -q https://download.opensuse.org/repositories/isv:/Rancher:/Elemental:/Stable/containers/rpi.raw
```

##### Verifying the download

In order to verify the integrity of the downloaded artifacts, you
should do a checksum verification:


```shell showLineNumbers
wget -q https://download.opensuse.org/repositories/isv:/Rancher:/Elemental:/Stable/containers/rpi.raw.sha256
sha256sum -c rpi.raw.sha256
```

This should print `rpi.raw: OK` as output.

#### Injecting the registration information

Adding the `initial-registration.yaml` isn't scripted yet. This is still a manual process:

The written USB stick will have two partitions. `RPI_BOOT` contains the boot loader files and `COS_LIVE` the Elemental files.
Mount the `COS_LIVE` partition and write `initial-registration.yaml` as `livecd-cloud-config.yaml` to this partition.

If you've mounted the USB stick with a file manager, this command should work to copy the registration information:

```shell showLineNumbers
sudo cp initial-registration.yaml /run/media/$USER/COS_LIVE/livecd-cloud-config.yaml
```

If you prefer using some CLI tools:

```shell showLineNumbers
IMAGE=rpi.raw
DEST=$(mktemp -d)

SECTORSIZE=$(sfdisk -J ${IMAGE} | jq '.partitiontable.sectorsize')
DATAPARTITIONSTART=$(sfdisk -J ${IMAGE} | jq '.partitiontable.partitions[1].start')
sudo mount -o rw,loop,offset=$((${SECTORSIZE}*${DATAPARTITIONSTART})) ${IMAGE} ${DEST}
sudo cp initial-registration.yaml ${DEST}/livecd-cloud-config.yaml
sudo umount ${DEST}
rmdir ${DEST}
```

#### Writing the seed image to a USB stick

The `.raw` image needs to be written to a USB stick to boot from. This can be done with `dd` on the Linux command line if you're comfortable with this command.
[openSUSE](https://www.opensuse.org) has nice instructions on how to write an image to a storage medium for [Linux](https://en.opensuse.org/SDB:Live_USB_stick),
[Windows](https://en.opensuse.org/SDB:Create_a_Live_USB_stick_using_Windows), and [OS X](https://en.opensuse.org/SDB:Create_a_Live_USB_stick_using_macOS).

#### Booting the Raspberry Pi

Now unmount the USB stick and plug it into your Raspberry Pi.

Plug a large (32 GB or more) and **fast** (!!) micro SD-card into the respective slot.

Connect the system to ethernet.

A powercycle will reboot the Pi. Everything else is identical to x86-64.

:::warning warning
Make sure the micro SD-card is unpartitioned. Otherwise the Pi bootloader will try to boot from it and fail.
:::

</TabItem>
</Tabs>

You can now boot your nodes with this image and they will:

- Register with the registrationURL given and create a per-machine `MachineInventory`
- Install SLE Micro to the given device
- Reboot

### Selecting the right machines to join a cluster

The `MachineInventorySelectorTemplate` selects the machines needed to provision the cluster from the `MachineInventory`s having the *element:fire* label.
We have added the *element*:*fire* label in the `MachineRegistration` `machineInventoryLabels` map, so all the `MachineInventory`s originated from it already have the label.
One could anyway skip the label from the `MachineRegistration` and add it later:

```shell showLineNumbers
kubectl -n fleet-default label machineinventory $(kubectl get machineinventory -n fleet-default --no-headers -o custom-columns=":metadata.name") element=fire
```

As soon as `MachineInventory`s with the *element*:*fire* are present, the corresponding machines auto-deploy the cluster via the chosen provider (k3s/rke).

After a few minutes your new cluster will be fully provisioned!!

## How can I choose the kubernetes version and deployer for the cluster?

In your cluster.yaml file there is a key in the `Spec` called `kubernetesVersion`. That sets the version and deployer that will be used for the cluster,
for example Kubernetes`v1.24.8` for rke2 would be `v1.24.8+rke2r1` and for k3s `v1.24.8+k3s1`.

To see all compatible versions check the [Rancher Support Matrix](https://www.suse.com/suse-rancher/support-matrix/all-supported-versions/) PDF for rke/rke2/k3s versions and their components.

You can also check our [Version doc](kubernetesversions.md) to know how to obtain those versions.

Check our [Cluster Spec](cluster-reference.md) page for more info about the `Cluster` resource.

## How can I follow what is going on behind the scenes?

You should be able to follow along what the machine is doing via:

- During ISO boot:
  - ssh into the machine (user/pass: root/ros):
    - running `journalctl -f -t elemental` shows you the progress of the registration (*elemental-register*) and the installation of Elemental (*elemental install*).
- Once the system is installed:
  - On the Rancher UI -> `Cluster Management` allows you to see your new cluster and the `Provisioning Log` in the cluster details
  - ssh into the machine (user/pass: Whatever your configured on the registration.yaml under `Spec.config.cloud-config.users`):
    - running `journalctl -f -u elemental-system-agent` shows the output of the initial elemental config and the installation of the `rancher-system-agent`
    - running `journalctl -f -u rancher-system-agent` shows the output of the bootstrap of cluster components like k3s
    - running `journalctl -f -u k3s` shows the logs of the k3s deployment

## Optional: Install the Elemental UI extension via CLI

Create and apply a new ClusterRepo to add the official Rancher Extensions repository.
```
apiVersion: catalog.cattle.io/v1
kind: ClusterRepo
metadata:
  name: rancher-ui-charts
spec:
  gitBranch: main
  gitRepo: https://github.com/rancher/ui-plugin-charts
```

Add a helm repo for the Rancher UI extensions charts.
```
helm repo add rancher-ui-plugins https://raw.githubusercontent.com/rancher/ui-plugin-charts/main/
```

Install the Elemental UI extension.
```
helm install elemental rancher-ui-plugins/elemental -n cattle-ui-plugin-system
```


---

## Article: quickstart-ui.md

---
sidebar_label: Elemental the visual way
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/quickstart-ui"/>
</head>

import Cluster from "!!raw-loader!@site/examples/quickstart/cluster.yaml"
import Registration from "!!raw-loader!@site/examples/quickstart/registration.yaml"
import RegistrationRPi from "!!raw-loader!@site/examples/quickstart/rpi-registration.yaml"
import Selector from "!!raw-loader!@site/examples/quickstart/selector.yaml"

# Elemental the visual way

:::note
The following instructions need Rancher 2.9.x at least.  
:::

This quickstart will show you how to deploy the Elemental plugin and operator into an existing Rancher Manager instance.

Once installed, you'll be able to provision a new Elemental cluster based on RKE2 or K3s.

However, if you want to install staging or dev operator, you can only do it in [CLI mode](quickstart-cli#non-stable-installations).

## Add the Official Rancher Extensions Repository

If the Elemental extension is not available, you need to add the `Official Rancher Extensions Repository`:

![Add Rancher Manager Extensions repository](images/quickstart-ui-extension-repository.png)

If this repository can not be added through the Extensions UI settings, it can be manually managed by [adding](https://ranchermanager.docs.rancher.com/how-to-guides/new-user-guides/helm-charts-in-rancher#manage-repositories) the `https://github.com/rancher/ui-plugin-charts` `git` repository, or alternatively by applying the following resource:  

```yaml
apiVersion: catalog.cattle.io/v1
kind: ClusterRepo
metadata:
  name: rancher-ui-charts
spec:
  gitBranch: main
  gitRepo: https://github.com/rancher/ui-plugin-charts
```

## Install the elemental plugin

After the Rancher Manager Extensions Support is enabled, you can install the `elemental` plugin as follow:

* Under the `Available` tab you will see `elemental` plugin available

![Rancher Manager Available plugins](images/quickstart-ui-extensions-available.png)

:::note
If the `Available` tab shows no entries, refresh the page. The `elemental` plugin will then appear.
:::

* Click on the `Install` button, a popup will appear and click on `Install` again to continue.

![Elemental plugin install](images/quickstart-ui-elemental-plugin-install.png)

* On the `Installed` tab, the `elemental` plugin is now listed.

:::note
If the `elemental` plugin is listed and the status stays at `Installing...`, refresh the page. The `elemental` plugin will display correctly.
:::

Once the `elemental` plugin installed, you can see the `OS Management` option in the Rancher Manager menu, refresh the page if you do not see it.

![Rancher Manager OS Management menu](images/quickstart-ui-elemental-plugin-menu.png)

## Install the elemental operator

:::note
The following guide will show you how to install the operator through the Elemental UI. But you can also install it directly from the Marketplace.
:::

Click on the OS Management button in the navigation menu.

If the operator is not already installed, the elemental ui will let you deploy it by clicking on the `Install Elemental Operator` button:
![Button to deploy elemental operator](images/quickstart-ui-extension-operator-button.png)

It will redirect you to the Rancher Marketplace to install the operator.

Click on the `Next` button:
![Install Elemental operator screenshot 1](images/quickstart-ui-extension-operator-install-1.png)

In this screen, you can customize or use the default values, click on `Install` to continue:
![Install Elemental operator screenshot 2](images/quickstart-ui-extension-operator-install-2.png)

You should see `elemental-operator-crds`and `elemental-operator` deployed in the `cattle-elemental-system` namespace:
![Install Elemental operator screenshot 3](images/quickstart-ui-extension-operator-install-3.png)

:::warning
If you do not see them, make sure to select the correct namespace at the top of the page.
:::

## Add a Machine Registration Endpoint

In the OS Management dashboard, click the `Create Registration Endpoint` button.

![OS Management registration endpoints](images/quickstart-ui-registration-endpoint-create.png)

Now here either you can enter each detail in its respective places or you can edit this as YAML and create the endpoint in one go. Here we'll edit every fields.

![Create a Registration Endpoint with UI](images/quickstart-ui-registration-endpoint-create-details.png)

:::info main options
`name: elemental-cluster1`: change this as per your need

`device-selector`: The [device-selector](machineregistration-reference#configelementalinstalldevice-selector) field can be used to dynamically pick device during installation. The field contains a list of rules to select the device you want.

`snapshotter`: Type of device used to manage snapshots in OS images.
:::

Once you create the machine registration end point it should show up as active.

![Machine registered in Registration Endpoints](images/quickstart-ui-registration-endpoint-complete.png)

## Preparing the installation (seed) image

Now this is the last step, you need to prepare a seed image that includes the initial registration config, so
it can be auto registered, installed and fully deployed as part of your cluster. The contents of the file are nothing
more than the registration URL that the node needs to register and the proper server certificate, so it can connect securely.

This seed image can then be used to provision an infinite number of machines.

The seed image is created as a Kubernetes resource above and can be built using the `Build Media` button, but first, you have to select ISO or RAW image.

In opposite to ISO where it needs two devices (device with ISO and another disk where to install Elemental), RAW image allows to boot from a single device and directly install the operating system in the device.
RAW image only contains a boot and a recovery partition and it boots first into recovery mode to install Elemental (for information, the process is similar to the [reset](reset#reset-workflow) one).


![Build Media in Registration Endpoints](images/quickstart-ui-registration-endpoint-build-media.png)

Once the build is done, media can be downloaded using the `Download Media` button:

![Download Media in Registration Endpoints](images/quickstart-ui-registration-endpoint-download-media.png)

You can now boot your nodes with this image and they will:

- Register with the registrationURL given and create a per-machine `MachineInventory`
- Install SLE Micro to the given device
- Reboot

## Machine Inventory

When nodes are booting up for the first time, they connect to Rancher Manager and a [`Machine Inventory`](machineinventory-reference.md) is created for each node.

![Machine Inventory menu](images/quickstart-ui-machine-inventory-menu.png)

Custom columns are based on `Machine Inventory Labels` which you can add when you create your `Machine Registration Endpoint`:

![Machine Registration Endpoint Hardware Labels](images/quickstart-ui-registration-endpoint-hardware-labels.png)

On the following screenshot, [`Hardware Labels`](hardwarelabels#hardware-labels) are used as custom columns:

You can also add custom columns by clicking on the three dots menu.

![Machine Inventory custom columns](images/quickstart-ui-machine-inventory-custom-columns.png)

Finally, you can also filter your `Machine Inventory` using those labels.

For instance if you only want to see your AMD machines, you can filter on `CPUModel` like below:

![Machine Inventory filtering](images/quickstart-ui-machine-inventory-filtering.png)

## Create your first Elemental Cluster

Now let's use those `Machine Inventory` to create a cluster by clicking on `Create Elemental Cluster` :

![Create Elemental Cluster button](images/quickstart-ui-create-cluster-button.png)

For your Elemental cluster, you can either choose K3s or RKE2 for Kubernetes.

![Elemental Cluster Creation Screen](images/quickstart-ui-create-cluster-standard-screen-.png)

Most of the options are coming from Rancher, that's why we will not detail all the possibilities.
Feel free to check the [Rancher Manager documentation](https://ranchermanager.docs.rancher.com/pages-for-subheaders/rancher-server-configuration) if you want to know more.

However, it is important to highlight the `Inventory of Machines Selector Template` section.

It lets you choose which `Machine Inventory` you want to use to create your Elemental cluster using the previously defined `Machine Inventory Labels` :

![Use Machine Inventory Selector Template](images/quickstart-ui-create-cluster-machine-selector-template.png)

As our three Machine Inventories contain the label `CPUVendor` with the key `AuthenticAMD`, the three machines will be used to create the Elemental cluster.


---

## Article: rancher-ip.md

---
sidebar_label: Rancher IP address options
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/rancher-ip"/>
</head>

import Registration from "!!raw-loader!@site/examples/labeltemplates/registration.yaml"

# Configure K3s and RKE2 internal and external IP addresses

K3s and RKE2 allow to specify the internal and external IP addresses of a node via the
`--node-ip` and `--node-external-ip` parameters.

Rancher Provisioning allows to fill those parameters collecting the IP address from each node's network
interface.

Anyway, the methods to collect those ip options for _Custom_ Rancher Clusters and
_Elemental_ Rancher Clusters differ.

## K3s and RKE2 internal and external IP addresses configuration in Rancher Cluster provisioning
Rancher provisioning allows to specify in the Rancher Agent section the network interfaces that should
be bound to the internal and external IP addresses of each provisioned node.

This is performed by adding the network interface names from which the IP addresses should be extracted
in the `CATTLE_INTERNAL_ADDRESS` and the `CATTLE_ADDRESS` Agent Environment Variables
([see Rancher docs](https://ranchermanager.docs.rancher.com/reference-guides/cluster-configuration/rancher-server-configuration/use-existing-nodes/rancher-agent-options#ip-address-options)).

:::note
The `CATTLE_INTERNAL_ADDRESS` and the `CATTLE_ADDRESS` Agent Environment Variables can be directly filled
with the desired IP addresses. Anyway, since the internal and external IP addresses are *per node*,
this would not work but for single node clusters.
:::

:::warning
While during the creation of an Elemental Cluster it is possible to add Agent Environment Variables,
the `CATTLE_ADDRESS` and the `CATTLE_INTERNAL_ADDRESS` ones are ignored and would not result in the
configuration of the internal and external IP addresses of the provisioned nodes.
:::

## Configure K3s and RKE2 internal and external IP addresses in Elemental Clusters
Elemental allows to configure the internal and external IP addresses of the Cluster Nodes attaching
the `elemental.cattle.io/InternalIP` and `elemental.cattle.io/ExternalIP` labels to the
[MachineInventory](machineinventory-reference.md) resources tracking the target nodes.

These labels when attached to a [MachineInventory](machineinventory-reference.md) resource are used
to fill the internal and external IP addresses of the associated nodes.

The labels can be added to the [MachineRegistration](machineregistration-reference.md)
`machineInventoryLabels` fields, using the
[Network Label Template](label-templates-network.md) IP address variables as values, in order
to allow to collect the IP addresses of each node.

Example: [MachineRegistration](machineregistration-reference.md) where nodes will have the internal IP
address set from interface eth0 and the external IP address from eth1.

<CodeBlock language="yaml" title="registration.yaml" showLineNumbers>{Registration}</CodeBlock>


---

## Article: rancher-vmware.md

---
sidebar_label: How to use Elemental with Rancher and VMware
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/rancher-vmware"/>
</head>

import Registration from "!!raw-loader!@site/examples/quickstart/registration.yaml"
import SeedImage from "!!raw-loader!@site/examples/quickstart/seedimage.yaml"

# How to use Elemental with Rancher and VMware

## Excerpt

In this document we will see how to set a virtual machine in VMware workstation to boot Elemental nodes.

## Prerequisites

1. Rancher 2.8 or higher installed and running. See the quick start guides.

## Step 1: Create the registration end point

We need to create an ISO image to bootstrap nodes and register against the Rancher instance. For that
a registration end point ([MachineRegistration](machineregistration-reference.md) resource) is required.

See this example of MachineRegistration:

<CodeBlock language="yaml" title="registration.yaml" showLineNumbers>{Registration}</CodeBlock>

The above MachineRegistration assumes the nodes include TPM 2.0. In case the virtualized target machine does
not include a virtual TPM device a software emulation can be configured in the
`config.elemental.registration` section.

Consider the following example:

```yaml showLineNumbers
apiVersion: elemental.cattle.io/v1beta1
kind: MachineRegistration
metadata:
  ...
spec:
  config:
    cloud-config:
      ...
    elemental:
      install:
        ...
      registration:
        emulate-tpm: true
        emulated-tpm-seed: -1
  machineInventoryLabels:
    ...
```

`emulated-tpm-seed: -1` sets the client to use a random seed to compute TPM hash, this is useful to be capable
to reuse the same registration end point definition for multiple machines. See further [TPM documentation](tpm.md).

## Step 2: Create the installation ISO

The installation media needs to be tied to a specific registration end point. This is created and handled
with the [SeedImage](seedimage-reference.md) resource.

Consider the following example:

<CodeBlock language="yaml" title="seedimage.yaml" showLineNumbers>{SeedImage}</CodeBlock>

Once the SeedImage resource is created it starts building an ISO with the provided OS image and linking it to
the given registration end point. Once done a download URL will be available in the SeedImage resource status.

You can download it with:

```shell
wget --no-check-certificate --content-disposition $(kubectl get seedimages.elemental.cattle.io -n fleet-default fire-img -o jsonpath="{.status.downloadURL}")
```

## Step 3: Boot the target device

Now ideally you would just burn the iso to a usb drive and boot your edge device using the usb device and once it boots and become active in Rancher under machine inventory you can select and create a cluster from it, however here we will use a vm to mimic an edge device for testing.

### 3.1 Prepare the VM to emulate TPM

In VMware workstation create a vm the way you would do normally, make sure to give the HDD size at least 40 GB.

Now edit the machine settings and go to the "Options" tab. The very last option would be "Advanced".

Click on "advanced" and on the right window pane change the firmware type from "BIOS" to "UEFI" and check the "Enable secure boot" option as follow:

* Default settings with BIOS selected

![VM boot options with BIOS](images/rancher-vmware-vm-boot-bios.png)

* Updated settings with UEFI selected and secure boot enabled

![VM boot options with UEFI](images/rancher-vmware-vm-boot-uefi.png)

Now on the same "Options" tab click on the "Access Control" option and click on "Encrypt" on the right side.

![Access control menu](images/rancher-vmware-access-control-menu.png)

This will ask you to enter a password to encrypt the machine. Enter a password and click on "Encrypt"

![Access control encryption credentials](images/rancher-vmware-access-control-encrypt.png)  

This is important to add the TPM Hardware. Next go back to the Hardware options and click on "Add"

And add the TPM (Trusted Platform Module) hardware and click on "Finish"

Now with the completion of this step our VM is ready.

### 3.2 Boot the VM with the elemental ISO

Next add the ISO that we created earlier in the VM and boot it up.

It should boot up with the ISO and start installing Elemental:

![Elemental OS install grub menu](images/rancher-vmware-elemental-install-grub.png)

![Elemental OS install logs](images/rancher-vmware-elemental-install-logs.png)

And once it is complete it will reboot the VM and it will show up as active under the machine inventory in Rancher.



---

## Article: raspi-disk.md

---
sidebar_label: Building raw disk images for Raspberry Pi
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/raspi-disk"/>
</head>


### How to build raw disk images for Raspberry Pi

This guide will show how we can build a raw disk image that can be written to an SD-card and booted without any other installation media.

:::caution
Any data on the SD-card will be erased, please only use a SD-card without anything important on it.

The SD-card must be reasonably large (32 GB or more) and **fast** (!!).
:::

```yaml title="SeedImage resource" showLineNumbers
apiVersion: elemental.cattle.io/v1beta1
kind: SeedImage
metadata:
  name: fire-img
  namespace: fleet-default
spec:
  type: raw
  baseImage: registry.opensuse.org/isv/rancher/elemental/staging/containers/suse/sl-micro/6.0/baremetal-os-container:latest
  targetPlatform: linux/arm64
  registrationRef:
    apiVersion: elemental.cattle.io/v1beta1
    kind: MachineRegistration
    name: fire-nodes
    namespace: fleet-default
```

Check the logs for the build pod using:

```shell
kubectl logs -n fleet-default fire-img -f -c build
```

When the build is finished we can download the image file using wget:

```shell
wget --no-check-certificate $(kubectl get seedimage -n fleet-default fire-img -o jsonpath="{.status.downloadURL}") -O sle-micro.arm64.raw
```

Now we can write the `.raw` image to the SD-card. This can be done with `dd` on the Linux command line if you're comfortable with this command.
[openSUSE](https://www.opensuse.org) has nice instructions on how to write an image to a storage medium for [Linux](https://en.opensuse.org/SDB:Live_USB_stick),
[Windows](https://en.opensuse.org/SDB:Create_a_Live_USB_stick_using_Windows), and [OS X](https://en.opensuse.org/SDB:Create_a_Live_USB_stick_using_macOS).

### Starting the machine

The raw disk image will only include the EFI partition, OEM partition and
recovery partition. On first boot the system will boot into the recovery system
to expand and add missing partitions. After expansion it will register with
rancher and reboot.

If an error occurs during registration phase the journal can be found using
`journalctl -u elemental-register-reset`.


---

## Article: release-notes.md

---
sidebar_label: Release Notes
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/release-notes"/>
</head>

# Release Notes

The Elemental project stack is made of various components such as the `Operator` and `UI` for example.

Each of these components have an independent development lifecycle with its own versioning. Once a new version is ready, meaning it fully integrates with the others components of the Elemental project stack, a release is made.

Here's the different components, their latest version and a link to the respective release notes on GitHub:

| Name                                                                 | Version | Release Notes                                                                |
|----------------------------------------------------------------------|---------|------------------------------------------------------------------------------|
| [Elemental Operator](https://github.com/rancher/elemental-operator/) | v1.7.3  | [Link](https://github.com/rancher/elemental-operator/releases/tag/v1.7.3)    |
| [Elemental Toolkit](https://github.com/rancher/elemental-toolkit/)   | v2.2.2  | [Link](https://github.com/rancher/elemental-toolkit/releases/tag/v2.2.2)     |
| [Elemental Linux](https://github.com/rancher/elemental)              | v2.2.0  | [Link](https://github.com/rancher/elemental/releases/tag/v2.2.0)             |
| [Elemental UI](https://github.com/rancher/elemental-ui)              | v3.0.1  | [Link](https://github.com/rancher/elemental-ui/releases/tag/elemental-3.0.1) |

:::note Information on docs versioning

The docs versioning is based on the `Elemental Operator` component as it's the user "entrypoint" to the Elemental project stack.

:::

## Install or Upgrade to latest release

In order to install this release of the Elemental Operator check the project documentation.

For already existing deployments use the following Helm commands to upgrade:

```
# Install/upgrade the CRDS chart
helm upgrade \
    --install -n cattle-elemental-system --create-namespace elemental-operator-crds \
    oci://registry.suse.com/rancher/elemental-operator-crds-chart

# Install/upgrade the operator chart
helm upgrade \
    --install -n cattle-elemental-system --create-namespace elemental-operator \
    oci://registry.suse.com/rancher/elemental-operator-chart
```

To install or upgrade from the helm chart repository use:

```
helm repo add elemental-stable https://rancher.github.io/elemental-operator/stable/
```

and installed or upgraded with

```
# Install/upgrade the CRDS chart
helm upgrade --install -n cattle-elemental-system --create-namespace \
    elemental-operator-crds elemental-stable/elemental-operator-crds

# Install/upgrade the operator chart
helm upgrade --install -n cattle-elemental-system --create-namespace \
    elemental-operator elemental-stable/elemental-operator
```

## Known issues

### Selinux in permissive mode

Setting selinux in enforcing mode is not supported as of today with Elemental.

### Install hooks not applicable in MachineRegistration resources

The cloud-config defined in `MachineRegistrations` is not applying `after-install-chroot` stage. Since
SL Micro 6.1 in order to apply `after-install-chroot` [yip stages](cloud-config-reference#elemental-client-cloud-config-hooks)
they should be defined as part of the `SeedImage` cloud-config. This stage is executed at install time and
so that it needs to be present in the installation media.

### ManagedOSVersion of type ISO may report a wrong version number

The `ManagedOSVersions` used for OS installation and upgrades come from the OS Channel (`ManagedOSVersionChannel`)
shipped with the Elemental Operator. The Channel contains two *types* of `ManagedOSVersions`: `container` and `iso`,
where the former is used for OS upgrades and the latter for new installations.
The `iso` types are sometimes labelled with a OS version lower than the actual one. This can be easily spotted by
checking if the latest version of the available `ManagedOSVersions` of type `container` lacks a matching version of a
`ManagedOSVersion` of type `iso`.

Example: the latest OS version actually present in the `registry.suse.com/rancher/elemental-channel/sl-micro:6.1-baremetal`
OS channel is `v2.2.0-4.4`. The ManagedOSVersion of type `container` is correctly labelled `v2.2.0-4.4`, while the latest
version of the ManagedOSVersion of type `iso` is `v2.2.0-4.3`: the `iso` type contains instead the OS version `v2.2.0-4.4`,
as would result by checking the `/etc/os-release` file of the installed machine.

### Predictable Network Interface Names

The SLE Micro OS images with versions v2.1.1 and v2.1.2 (released in the default
[ManagedOSVersionChannel](managedosversionchannel-reference))
adopt predictable network interface names by default.

This is a change from SLE Micro OS images previously released, so you should expect your
Elemental hosts to switch the network interface names from the `ethX` template to the `enpXsY` one.

You can disable the predictable network interface names by passing the `net.ifnames=0` argument
to the kernel command line. To make it permanent:

```sh
grub2-editenv /oem/grubenv set extra_cmdline=net.ifnames=0
```

:::warning
The adoption of the predictable network interface names feature was not a planned one:
it will be reverted in the next SLE Micro OS images starting from version v2.1.3.
These OS images will include the `net.ifnames=0` kernel command line argument by default.  
The v2.1.3 OS images will be released via the default Elemental 1.6 channel.
:::

### SSH root access

The SLE Micro OS images released in the current Elemental version (through the default
[ManagedOSVersionChannel](managedosversionchannel-reference)) do not allow ssh root access
via password anymore. Easyest workaround is to either configure ssh root access via an ssh
key or add a new user to the system.

### Kernel Panic on hypervisors

OS Images based on SL Micro 6.0 can fail to boot with a kernel panic on virtual machines using an unsupported CPU type.  
The `x86-64-v2` instruction set is required. For best compatibility CPU host passthrough is recommended.


---

## Article: removable-device-cloudconfig.md

---
sidebar_label: Include cloud-config from removable devices
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/removable-device-cloudconfig"/>
</head>


### How to include cloud-config files from removable devices

Elemental nodes supports loading [cloud-config](cloud-config-reference.md) files from specific block devices.
In particular supports loading cloud-config files from an ISO having `CIDATA` as the volume ID or any vFAT formatted
device labeled with `CIDATA`. If a device matching this criteria is found on early boot the Elemental client will
read it and look for a `user-data` file in its root.

As an example an ISO including a cloud-config file can be created on a Linux host with the procedure below.

Create a `user-data` file with the cloud-config data in it. In the example below we just set a
proxy:

```yaml title="user-data" showLineNumbers
#cloud-config
write_files:
- path: /etc/sysconfig/proxy
  append: true
  content: |
    PROXY_ENABLED="yes"
    HTTP_PROXY=http://some.domain.org:8080
    HTTPS_PROXY=https://some.domain.org:8080
    NO_PROXY="localhost, 127.0.0.1"
```

Once the `user-data` file exists create an ISO including only this file by using the `mkisofs` Linux utility:

```bash
mkisof -o cidata.iso -V CIDATA -J -r user-data
```

The result is an ISO labeled with `CIDATA` including the `user-data` file.

At boot the `user-data` file will be copied as is to `/oem/user-data` and in case it contains cloud-config data
an extra copy will be added as `/oem/user-data.yaml`. The file `/oem/user-data.yaml` will be parsed
on any later cloud-init stage.

Since the data is copied to `/oem` it will be persistent, hence on follow up reboots the removable device is
not required to be present any more. If still present on follow up reboots, it just overwrites any
already pre-existing data.

#### Include non cloud-config data

If the `user-data` is not containing cloud-config data the Elemental client will just copy it as
is to `/oem/user-data`. Only `*.yaml` files are parsed when executing cloud-init stages, so in that
case the file will be ignored by cloud-init services.

If the `user-data` contains a script the Elemental client will, in addition, try to execute it. The way
Elemental client determines if `user-data` is a script or not is by the presence of a _Shebang_ in the
first line. For example, the previous `user-data` file could be rewritten as:


```bash title="user-data" showLineNumbers
#!/bin/bash

cat <<EOF >> /etc/sysconfig/proxy
PROXY_ENABLED="yes"
HTTP_PROXY=http://some.domain.org:8080
HTTPS_PROXY=https://some.domain.org:8080
NO_PROXY="localhost, 127.0.0.1"
EOF
```


---

## Article: reset.md

---
sidebar_label: Machine Reset
title: ''
version_badge: '1.3.0'
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/reset"/>
</head>

## Machine Reset

There are two ways to reset Elemental machines to their original state or decommission them:

1. When deleting a Cluster, all associated machines will be reset  

![Delete a Cluster to reset all machines](images/reset-cluster-deletion.png)

2. When managing a Cluster, simply delete the Node that needs to be reset  

![Delete a single node to reset it](images/reset-single-node-deletion.png)

### Reset workflow

Once the related `MachineInventory` is flagged for deletion, a reset plan will be executed by the `elemental-system-agent` running on the machine.  

If the machine is still running, this plan will:

1. Reboot the machine in recovery mode.
2. Execute `systemctl start elemental-register-reset`.  
   This will fetch the remote `MachineRegistration` and apply the `spec.config.elemental.reset` options to reset the machine.  
   A new `MachineInventory` will be created and the `spec.config.cloud-config` defined in the `MachineRegistration` will be applied again.  

Note that the `MachineRegistration` reference will **not** change, the machine will **not** be reinstalled, the `COS_PERSISTENT` and `COS_OEM` partition will be cleared by default if reset is `enabled`. For more information, you can consult the [Partition Table](installation#deployed-partition-table).  

Since the `cloud-config` is re-applied during the reset workflow, you can reset a machine to apply updates from the `MachineRegistration` definition, for example to rotate `users` credentials and authorized keys. It is strongly recommended to enable the `reset-oem` option, to avoid conflicts with previously configured cloud-configs.  

If you need to bind a machine to a different `MachineRegistration` and trigger a new full installation, you need to reprovision it again using a new image.  

### Enable machine reset

In order to allow machines to be reset automatically, the `spec.config.elemental.reset.enabled` flag of the `MachineRegistration` should be toggled.  
This is off by default, but once activated, all newly created `MachineInventory` will inherit this setting automatically.  
For example:

```yaml
apiVersion: elemental.cattle.io/v1beta1
kind: MachineRegistration
metadata:
  name: fire-nodes
  namespace: fleet-default
spec:
  config:
    elemental:
      reset:
        enabled: true
        reset-persistent: true
        reset-oem: true
        # These cloud-init configs will be created during reset and will persist on the system after
        config-urls: 
          - "https://my.cloud.init/reset-plan-1.yaml"
          - "https://my.cloud.init/reset-plan-2.yaml"
        # You can select a different image to run the reset.  
        # Note that this image will not be installed on the system.
        system-uri: "my.oci.registry/reset-image:latest"
        power-off: false
        reboot: true
```

It is also possible to enable reset at a `MachineInventory` level, for example in scenarios where some machines are physical and will benefit from an automatic reset, and some others are virtual and can simply be destroyed and reprovisioned as needed.  
In order to flag a single `MachineInventory` to allow reset, you can use the `elemental.cattle.io/resettable: true` annotation.  
For example:  

```yaml
apiVersion: elemental.cattle.io/v1beta1
kind: MachineInventory
metadata:
  annotations:
    elemental.cattle.io/resettable: "true"
```


---

## Article: restore.md

---
sidebar_label: Restore
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/restore"/>
</head>

# Restore

Follow this guide to restore an Elemental configuration from a backup with Rancher.

## Prepare rancher-backup operator and backup files for restoring

Go to official [Rancher documentation](https://docs.ranchermanager.rancher.io/how-to-guides/new-user-guides/backup-restore-and-disaster-recovery/restore-rancher) and make sure that `rancher-backup operator` is installed and has access to backup files.

You first have to follow this [documentation](backup.md) to do a backup of Elemental resources.

## Restore the Elemental configuration with rancher-backup operator

Create a `restore object` (adapted to your needs) to restore the backup tarball.

```yaml showLineNumbers
apiVersion: resources.cattle.io/v1
kind: Restore
metadata:
  name: elemental-restore
spec:
  backupFilename: rancher-backup-430169aa-edde-4a61-85e8-858f625a755b-2022-10-17T05-15-00Z.tar.gz
```

Apply manifest on Kubernetes.

```shell showLineNumbers
kubectl apply -f elemental-restore.yaml
```

Check logs from rancher-backup operator.

```shell showLineNumbers
kubectl logs -n cattle-resources-system -l app.kubernetes.io/name=rancher-backup -f
```

Verify if backup file was restore successfully.

```shell showLineNumbers
...
INFO[2022/10/31 06:34:50] Processing controllerRef apps/v1/deployments/rancher 
INFO[2022/10/31 06:34:50] Done restoring
...
```

Continue with procedure from [Rancher documentation](https://docs.ranchermanager.rancher.io/how-to-guides/new-user-guides/backup-restore-and-disaster-recovery/migrate-rancher-to-new-cluster)


---

## Article: seedimage-reference.md

---
sidebar_label: SeedImage reference
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/seedimage-reference"/>
</head>

# SeedImage reference

A `SeedImage` resource allows to build an installation media that can be used to install Elemental onto a node.
It requires a `baseImage`, i.e., a URL to an Elemental installation ISO or node container image, and a `registrationRef` reference to a `MachineRegistration` resource, from which the registration part of the Elemental configuration is extracted and injected in the media to produce the final *seed image*.
It is also possible to inject customizations in the `cloud-config` field. Both yip and cloud-init syntax are supported. See the [Cloud Config Reference](cloud-config-reference.md) for full information.

Once the seed image is ready, the download URL is shared in the `.status.downloadURL` field.
It stays available for download for `cleanupAfterMinutes` minutes (default is `60`, 1 hour), after which it is deleted.
Setting `retriggerBuild` to `true` retriggers the seed image build process while setting `cleanupAfterMinutes` to `0` keeps the seed image around till the `SeedImage` resource is deleted.

The `SeedImage` resource also has a `type` field which can be set to either `iso`, to build an ISO, or `raw` to build a raw disk image. Raw disk images can be copied directly to the target drive and on first boot will automatically boot into a recovery partition to expand the drive to use the available disk space and register the node, after which it will reboot the same way as for the ISO installation.

If no `BuildContainer` is specified for the seed-image it will be automatically filled in based on default values and `type`.

Building a SeedImage for a different platform is accomplished using the `targetPlatform` field. The platform is specified using `os/arch`, for example (`linux/x86_64` or `linux/aarch64`). By default the image will be built for the same platform that the operator is hosted on.

:::warning seed images may fill up local storage
The seed images are kept on the node's local storage: pay attention to the number of `SeedImage` resources you start concurrently and to the ones you may leave around with the auto-cleanup feature disabled (`cleanupAfterMinutes` = `0`) as you may exhaust the storage on your cluster nodes.
:::

## SeedImageSpec reference

| Key                 | Type        | Default value | Description                                                                                                                                    |
|---------------------|-------------|---------------|------------------------------------------------------------------------------------------------------------------------------------------------|
| baseImage           | string      | empty         | The base Elemental image used to build the seed image.                                                                                         |
| registrationRef     | object ref. | null          | A reference to a MachineRegistration that will be used for all installed machines to register.                                                 |
| buildContainer      | object      | null          | Settings for a custom container used to generate the downloadable image. (See [documentation](#buildcontainer)).                               |
| cleanupAfterMinutes | int         | 60            | The time after which the built seed image will be cleaned up. Active downloads will finish before the image is removed.                        |
| retriggerBuild      | bool        | false         | Trigger to build again a cleaned up seed image.                                                                                                |
| size                | string      | 6Gi           | Specifies the size of the volume used to store the image.                                                                                      |
| type                | string      | iso           | Specifies the type of seed image to built. `iso` or `raw`. (See [documentation](#iso-and-raw-images))                                          |
| targetPlatform      | string      | empty         | Specifies the target platform for the built image. Example: `linux/x86_64` or `linux/aarch64`. (See [documentation](#multi-platform-support)). |
| cloud-config        | object      | null          | Contains cloud-config data to be included in the generated image. (See [documentation](./cloud-config-reference.md)).                          |

### BuildContainer

The `buildContainer` settings can be used to customize the `build` init container within the `SeedImage`'s pod.  
This could be the case for example when building custom Elemental images.  

```yaml
buildContainer:
  name: "custom-build"
  image: my.registry.com/elemental-custom-builder:1.2.3
  command:
    - build-image
  args:
    - foo
    - bar
  imagePullPolicy: Always
```

Note that the container will additionally have two volumes mounted at `/iso` and `/overlay`.  
The SeedImage build process expects the build container to place the build artifact in `/iso/$(ELEMENTAL_OUTPUT_NAME)`.  

Configuration files are available in:

- `/overlay/reg/livecd-cloud-config.yaml`: A configuration file that can be used by `elemental-register` to register the machine.

- `/overlay/iso-config/cloud-config.yaml`: The cloud-config defined in `SeedImage.spec.cloud-config`

The following list of environment variables can also be used within the custom build container:

- `ELEMENTAL_DEVICE`: The `MachineRegistration.spec.config.elemental.install.device` value.
- `ELEMENTAL_REGISTRATION_URL`: The unique URL of the MachineRegistration.
- `ELEMENTAL_BASE_IMAGE`: The base image defined in the `SeedImage`.
- `ELEMENTAL_OUTPUT_NAME`: The expected file name of the build artifact.

### ISO and Raw images

The `SeedImage` is able to build `iso` or `raw` image types.  
Note that Elemental ships two different flavors of images, `iso` or `container` types. See [ManagedOSversion's type](./managedosversion-reference.md#managedosversionspec-reference).

When building a `iso` `SeedImage`, you can use an `iso` Elemental image.  
`iso` images contain a pre-built `.iso` artifact. This is the default Elemental way of shipping official ISOs, so that they don't need to be rebuilt every time you define a `SeedImage`.

<details>
  <summary>ISO SeedImage example</summary>

  ```yaml showLineNumbers
  apiVersion: elemental.cattle.io/v1beta1
  kind: SeedImage
  metadata:
    name: fire-iso
    namespace: fleet-default
  spec:
    type: iso
    baseImage: registry.suse.com/suse/sl-micro/6.0/baremetal-iso-image:2.1.1-3.36
    registrationRef:
      apiVersion: elemental.cattle.io/v1beta1
      kind: MachineRegistration
      name: fire-nodes
      namespace: fleet-default
  ```

</details>

Alternatively, when building a `raw` `SeedImage`, you should use `container` Elemental images. These images are also used during the upgrade process (See: [ManagedOSImage](./managedosimage-reference.md)), but can be used to build `raw` `SeedImages` as well.  

<details>
  <summary>Raw SeedImage example</summary>

  ```yaml showLineNumbers
  apiVersion: elemental.cattle.io/v1beta1
  kind: SeedImage
  metadata:
    name: fire-raw
    namespace: fleet-default
  spec:
    type: raw
    baseImage: registry.suse.com/suse/sl-micro/6.0/baremetal-os-container:2.1.1-3.29
    registrationRef:
      apiVersion: elemental.cattle.io/v1beta1
      kind: MachineRegistration
      name: fire-nodes
      namespace: fleet-default
  ```

</details>

### Multi-Platform support

Elemental ships `linux/x86_64` and `linux/aarch64` images for most flavors.  
In order to determine whether a `ManagedOSVersion` image supports both platforms, you can verify the `ManagedOSVersion.spec.metadata.platform` values. (See [documentation](./managedosversion-reference.md#metadata)).

When defining a `SeedImage`, you can then use this value for the image's `targetPlatform`.  
Leaving the `targetPlatform` empty, will default to the platform where the `elemental-operator` is running.  

<details>
  <summary>Raw aarch64 SeedImage example</summary>

  ```yaml showLineNumbers
  apiVersion: elemental.cattle.io/v1beta1
  kind: SeedImage
  metadata:
    name: fire-raw-aarch64
    namespace: fleet-default
  spec:
    targetPlatform: linux/aarch64
    type: raw
    baseImage: registry.suse.com/suse/sl-micro/6.0/baremetal-os-container:2.1.1-3.29
    registrationRef:
      apiVersion: elemental.cattle.io/v1beta1
      kind: MachineRegistration
      name: fire-nodes
      namespace: fleet-default
  ```

</details>

## Downloadable URLs

The `SeedImage` resource tracks the seed image build process through two status conditions:

- **Ready**: tracks the creation of all the required child resources that perform the actual build process.
- **SeedImageReady**: tracks the status of the build process in the child resources.

Alternatively it is also possible to wait for the `SeedImage` pod to be ready:

```bash
kubectl wait --for=condition=ready pod -n fleet-default fire-img
```

Waiting on Ready conditions is a best practice before downloading any artifact.

Once a `SeedImage` is ready, the `.status.downloadURL` will contain the downloadable URL.  
Note that the URL will use the same endpoint as Rancher, so beware of HTTPS validation when using self signed certificates.  

```bash
kubectl get seedimage -n fleet-default fire-img -o jsonpath="{.status.downloadURL}"
```

The checksum of the image is also available to verify the download was correct:

```bash
kubectl get seedimage -n fleet-default fire-img -o jsonpath="{.status.checksumURL}"
```


---

## Article: smbios.md

---
sidebar_label: SMBIOS
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/smbios"/>
</head>

import Registration from "!!raw-loader!@site/examples/quickstart/registration.yaml"

:::warning
SMBIOS Template Variables have been deprecated: please use the new
[Label Templates' Variables](label-templates#label-templates-variables) when possible.

Check the [deprecated variables page](label-templates-deprecated) and the
[conversion table](label-templates-deprecated#smbios-labels-to-new-label-templates-variable-table)
for a smooth transition.
:::

## SMBIOS Template Variables

The System Management BIOS (SMBIOS) specification defines data structures (and access methods) that can be used to read management information produced by the BIOS of a computer.

This allows us to gather hardware information about the running system and use that as part of our labels.

The [Elemental Register Client](architecture-components#elemental-register-command-line-tool) gathers SMBIOS data running the `dmidecode` binary during the initial registration of the machine and
sends that data to the [Elemental Operator](architecture-components#elemental-operator-daemon).

That data is used to render the [template label variables](label-templates#label-template-variables) in the [MachineInventory](machineinventory-reference) associated to that machine.

:::note Example
Having the following SMBIOS data:

```console showLineNumbers
System Information
	Manufacturer: My manufacturer
	Product Name: Awesome PC
	Version: Not Specified
	Serial Number: THX1138
	Family: Toretto
```

And setting the `machineName` to `serial-${System Information/Serial Number}` would result in the final value of `serial-THX1138`
:::

A good use of SMBIOS data is to set up different labels for all your machines and get those values from the hardware directly.

Having your `machineInventoryLabels` on the [machineRegistration](machineregistration-reference.md) set to SMBIOS data would allow 
you to use selectors down the line to select similar machines.

For example using the following label `cpuFamily: "${Processor Information/Family}` would allow you to use a selector to search for i7 cpus in your machine fleet.

<CodeBlock language="yaml" title="registration example with SMBIOS template variables" showLineNumbers>{Registration}</CodeBlock>


---

## Article: tpm.md

---
sidebar_label: Trusted Platform Module (TPM)
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/tpm"/>
</head>

import RegistrationTpm from "!!raw-loader!@site/examples/quickstart/registration-tpm.yaml"

# Trusted Platform Module 2.0 (TPM)

Trusted Platform Module (TPM, also known as ISO/IEC 11889) is an international standard for a secure cryptoprocessor, a dedicated microcontroller designed to secure hardware through integrated cryptographic keys. The term can also refer to a chip conforming to the standard.

## Add TPM module to virtual machine

Easy way to add TPM to virtual machine is to use Libvirt with Virt-manager

### Create Virtual Machine

After starting virt-manager create new virtual machine

![Create new VM](images/tpm1.png)

### Verify and edit hardware module list

On the hardware configuration screen, verify list of modules and click ***Add Hardware*** button

![Devices list](images/tpm2.png)

### Add TPM module to VM

From the list of emulated devices choose TPM module and add it to VM

![Add TPM module](images/tpm3.png)

### Finish VM configuration

On the last screen verify once again if TPM module was added properly

![Verify TPM](images/tpm4.png)

## Add TPM emulation to bare metal machine

During applying `#!yaml MachineRegistration` add following key to the yaml `config:elemental:registration:emulate-tpm: true`

:::info
If you plan to deploy more than 1 machine with TPM emulation, make sure to set `config:elemental:registration:emulated-tpm-seed: -1`
so the seed used for the TPM emulation is randomized per machine. Otherwise, you will get the same TPM Hash for all deployed machines and only the last
one to be registered will be valid.
:::

<CodeBlock language="yaml" title="registration-tpm.yaml" showLineNumbers>{RegistrationTpm}</CodeBlock>


---

## Article: troubleshooting-label-templates.md

---
sidebar_label: Label Templates
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/troubleshooting-label-templates"/>
</head>

# Troubleshooting Label Templates

In order to check/debug the rendering of the Label Templates in each of the registering host,
starting from version 1.7.0 the `elemental-register` binary can now dump to console the available
Label Templates and their values using the command *dumpdata*.

Example:
```shell showLineNumbers
rancher-6368:~ # elemental-register dumpdata
${BIOS/Date}                                        : "02/02/2022"
${BIOS/Vendor}                                      : "EDK II"
${BIOS/Version}                                     : "unknown"
${BaseBoard/AssetTag}                               : "unknown"
${BaseBoard/Product}                                : "unknown"
${BaseBoard/SerialNumber}                           : "unknown"
${BaseBoard/Vendor}                                 : "unknown"
${BaseBoard/Version}                                : "unknown"
${CPU/Processor/Capabilities}                       : "fpu,vme,de,pse,tsc,msr,pae,mce,cx8,apic,sep,mtrr,pge,mca,cmov,pat,pse36,clflush,mmx,fxsr,sse,sse2,ss,syscall,nx,pdpe1gb,rdtscp,lm,constant_tsc,arch_perfmon,rep_good,nopl,xtopology,cpuid,tsc_known_freq,pni,pclmulqdq,vmx,ssse3,fma,cx16,pdcm,pcid,sse4_1,sse4_2,x2apic,movbe,popcnt,tsc_deadline_timer,aes,xsave,avx,f16c,rdrand,hypervisor,lahf_lm,abm,cpuid_fault,pti,ssbd,ibrs,ibpb,stibp,tpr_shadow,flexpriority,ept,vpid,ept_ad,fsgsbase,tsc_adjust,bmi1,avx2,smep,bmi2,erms,invpcid,xsaveopt,arat,vnmi,umip,md_clear,flush_l1d,arch_capabilities"
${CPU/Processor/ID}                                 : "0"
${CPU/Processor/Model}                              : "Intel(R) Xeon(R) CPU E5-2630 v3 @ 2.40GHz"
${CPU/Processor/NumCores}                           : "1"
${CPU/Processor/NumThreads}                         : "1"
${CPU/Processor/Vendor}                             : "GenuineIntel"
${CPU/Processors/0/Capabilities}                    : "fpu,vme,de,pse,tsc,msr,pae,mce,cx8,apic,sep,mtrr,pge,mca,cmov,pat,pse36,clflush,mmx,fxsr,sse,sse2,ss,syscall,nx,pdpe1gb,rdtscp,lm,constant_tsc,arch_perfmon,rep_good,nopl,xtopology,cpuid,tsc_known_freq,pni,pclmulqdq,vmx,ssse3,fma,cx16,pdcm,pcid,sse4_1,sse4_2,x2apic,movbe,popcnt,tsc_deadline_timer,aes,xsave,avx,f16c,rdrand,hypervisor,lahf_lm,abm,cpuid_fault,pti,ssbd,ibrs,ibpb,stibp,tpr_shadow,flexpriority,ept,vpid,ept_ad,fsgsbase,tsc_adjust,bmi1,avx2,smep,bmi2,erms,invpcid,xsaveopt,arat,vnmi,umip,md_clear,flush_l1d,arch_capabilities"
${CPU/Processors/0/ID}                              : "0"
${CPU/Processors/0/Model}                           : "Intel(R) Xeon(R) CPU E5-2630 v3 @ 2.40GHz"
${CPU/Processors/0/NumCores}                        : "1"
${CPU/Processors/0/NumThreads}                      : "1"
${CPU/Processors/0/Vendor}                          : "GenuineIntel"
${CPU/Processors/1/Capabilities}                    : "fpu,vme,de,pse,tsc,msr,pae,mce,cx8,apic,sep,mtrr,pge,mca,cmov,pat,pse36,clflush,mmx,fxsr,sse,sse2,ss,syscall,nx,pdpe1gb,rdtscp,lm,constant_tsc,arch_perfmon,rep_good,nopl,xtopology,cpuid,tsc_known_freq,pni,pclmulqdq,vmx,ssse3,fma,cx16,pdcm,pcid,sse4_1,sse4_2,x2apic,movbe,popcnt,tsc_deadline_timer,aes,xsave,avx,f16c,rdrand,hypervisor,lahf_lm,abm,cpuid_fault,pti,ssbd,ibrs,ibpb,stibp,tpr_shadow,flexpriority,ept,vpid,ept_ad,fsgsbase,tsc_adjust,bmi1,avx2,smep,bmi2,erms,invpcid,xsaveopt,arat,vnmi,umip,md_clear,flush_l1d,arch_capabilities"
${CPU/Processors/1/ID}                              : "1"
${CPU/Processors/1/Model}                           : "Intel(R) Xeon(R) CPU E5-2630 v3 @ 2.40GHz"
${CPU/Processors/1/NumCores}                        : "1"
${CPU/Processors/1/NumThreads}                      : "1"
${CPU/Processors/1/Vendor}                          : "GenuineIntel"
${CPU/Processors/2/Capabilities}                    : "fpu,vme,de,pse,tsc,msr,pae,mce,cx8,apic,sep,mtrr,pge,mca,cmov,pat,pse36,clflush,mmx,fxsr,sse,sse2,ss,syscall,nx,pdpe1gb,rdtscp,lm,constant_tsc,arch_perfmon,rep_good,nopl,xtopology,cpuid,tsc_known_freq,pni,pclmulqdq,vmx,ssse3,fma,cx16,pdcm,pcid,sse4_1,sse4_2,x2apic,movbe,popcnt,tsc_deadline_timer,aes,xsave,avx,f16c,rdrand,hypervisor,lahf_lm,abm,cpuid_fault,pti,ssbd,ibrs,ibpb,stibp,tpr_shadow,flexpriority,ept,vpid,ept_ad,fsgsbase,tsc_adjust,bmi1,avx2,smep,bmi2,erms,invpcid,xsaveopt,arat,vnmi,umip,md_clear,flush_l1d,arch_capabilities"
${CPU/Processors/2/ID}                              : "2"
${CPU/Processors/2/Model}                           : "Intel(R) Xeon(R) CPU E5-2630 v3 @ 2.40GHz"
${CPU/Processors/2/NumCores}                        : "1"
${CPU/Processors/2/NumThreads}                      : "1"
${CPU/Processors/2/Vendor}                          : "GenuineIntel"
${CPU/Processors/3/Capabilities}                    : "fpu,vme,de,pse,tsc,msr,pae,mce,cx8,apic,sep,mtrr,pge,mca,cmov,pat,pse36,clflush,mmx,fxsr,sse,sse2,ss,syscall,nx,pdpe1gb,rdtscp,lm,constant_tsc,arch_perfmon,rep_good,nopl,xtopology,cpuid,tsc_known_freq,pni,pclmulqdq,vmx,ssse3,fma,cx16,pdcm,pcid,sse4_1,sse4_2,x2apic,movbe,popcnt,tsc_deadline_timer,aes,xsave,avx,f16c,rdrand,hypervisor,lahf_lm,abm,cpuid_fault,pti,ssbd,ibrs,ibpb,stibp,tpr_shadow,flexpriority,ept,vpid,ept_ad,fsgsbase,tsc_adjust,bmi1,avx2,smep,bmi2,erms,invpcid,xsaveopt,arat,vnmi,umip,md_clear,flush_l1d,arch_capabilities"
${CPU/Processors/3/ID}                              : "3"
${CPU/Processors/3/Model}                           : "Intel(R) Xeon(R) CPU E5-2630 v3 @ 2.40GHz"
${CPU/Processors/3/NumCores}                        : "1"
${CPU/Processors/3/NumThreads}                      : "1"
${CPU/Processors/3/Vendor}                          : "GenuineIntel"
${CPU/TotalCores}                                   : "4"
${CPU/TotalProcessors}                              : "4"
${CPU/TotalThreads}                                 : "4"
${Chassis/AssetTag}                                 : ""
${Chassis/SerialNumber}                             : ""
${Chassis/TypeDescription}                          : "Other"
${Chassis/Type}                                     : "1"
${Chassis/Vendor}                                   : "QEMU"
${Chassis/Version}                                  : "pc-q35-8.2"
${GPU/GraphicsCards/0/Driver}                       : "virtio-pci"
${GPU/GraphicsCards/0/ProductName}                  : "Virtio GPU"
${GPU/GraphicsCards/0/VendorName}                   : "Red Hat, Inc."
${GPU/GraphicsCards/Driver}                         : "virtio-pci"
${GPU/GraphicsCards/ProductName}                    : "Virtio GPU"
${GPU/GraphicsCards/VendorName}                     : "Red Hat, Inc."
${GPU/TotalCards}                                   : "1"
${Memory/Modules}                                   : ""
${Memory/SupportedPageSizes}                        : "1073741824, 2097152"
${Memory/TotalPhysicalBytes}                        : "8589934592"
${Memory/TotalUsableBytes}                          : "8312393728"
${Network/NICs/0/AdvertisedLinkModes}               : ""
${Network/NICs/0/Duplex}                            : ""
${Network/NICs/0/IPv4Addresses/0}                   : "172.21.1.195"
${Network/NICs/0/IPv4Address}                       : "172.21.1.195"
${Network/NICs/0/IPv6Addresses/0}                   : "fe80::5f8e:31ac:bd6e:52df"
${Network/NICs/0/IPv6Address}                       : "fe80::5f8e:31ac:bd6e:52df"
${Network/NICs/0/IsVirtual}                         : "false"
${Network/NICs/0/MacAddress}                        : "52:54:00:8e:e6:39"
${Network/NICs/0/Name}                              : "eth0"
${Network/NICs/0/Speed}                             : ""
${Network/NICs/0/SupportedLinkModes}                : ""
${Network/NICs/0/SupportedPorts}                    : ""
${Network/NICs/eth0/AdvertisedLinkModes}            : ""
${Network/NICs/eth0/Duplex}                         : ""
${Network/NICs/eth0/IPv4Addresses/0}                : "172.21.1.195"
${Network/NICs/eth0/IPv4Address}                    : "172.21.1.195"
${Network/NICs/eth0/IPv6Addresses/0}                : "fe80::5f8e:31ac:bd6e:52df"
${Network/NICs/eth0/IPv6Address}                    : "fe80::5f8e:31ac:bd6e:52df"
${Network/NICs/eth0/IsVirtual}                      : "false"
${Network/NICs/eth0/MacAddress}                     : "52:54:00:8e:e6:39"
${Network/NICs/eth0/Name}                           : "eth0"
${Network/NICs/eth0/Speed}                          : ""
${Network/NICs/eth0/SupportedLinkModes}             : ""
${Network/NICs/eth0/SupportedPorts}                 : ""
${Network/TotalNICs}                                : "1"
${Product/Family}                                   : ""
${Product/Name}                                     : "Standard PC (Q35 + ICH9, 2009)"
${Product/SKU}                                      : ""
${Product/SerialNumber}                             : ""
${Product/UUID}                                     : "d4f972c7-edcc-40a5-98e1-67753bf84add"
${Product/Vendor}                                   : "QEMU"
${Product/Version}                                  : "pc-q35-8.2"
${Runtime/Hostname}                                 : "rancher-6368"
${Storage/Disks/0/DriveType}                        : "ODD"
${Storage/Disks/0/Model}                            : "QEMU_DVD-ROM"
${Storage/Disks/0/Name}                             : "sr0"
${Storage/Disks/0/Removable}                        : "true"
${Storage/Disks/0/Size}                             : "1073741312"
${Storage/Disks/0/StorageController}                : "SCSI"
${Storage/Disks/1/DriveType}                        : "HDD"
${Storage/Disks/1/Model}                            : "unknown"
${Storage/Disks/1/Name}                             : "vda"
${Storage/Disks/1/Removable}                        : "false"
${Storage/Disks/1/Size}                             : "32212254720"
${Storage/Disks/1/StorageController}                : "virtio"
${Storage/Disks/sr0/DriveType}                      : "ODD"
${Storage/Disks/sr0/Model}                          : "QEMU_DVD-ROM"
${Storage/Disks/sr0/Name}                           : "sr0"
${Storage/Disks/sr0/Removable}                      : "true"
${Storage/Disks/sr0/Size}                           : "1073741312"
${Storage/Disks/sr0/StorageController}              : "SCSI"
${Storage/Disks/vda/DriveType}                      : "HDD"
${Storage/Disks/vda/Model}                          : "unknown"
${Storage/Disks/vda/Name}                           : "vda"
${Storage/Disks/vda/Removable}                      : "false"
${Storage/Disks/vda/Size}                           : "32212254720"
${Storage/Disks/vda/StorageController}              : "virtio"
${Storage/TotalDisks}                               : "2"

```


---

## Article: troubleshooting-network.md

---
sidebar_label: Declarative Networking
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/troubleshooting-network"/>
</head>

import RegistrationWithNetwork from "!!raw-loader!@site/examples/network/machineregistration.yaml"

# Troubleshooting Declarative Networking

Given the following sample registration:  

<CodeBlock language="yaml" title="example MachineRegistration using Declarative Networking" showLineNumbers>{RegistrationWithNetwork}</CodeBlock>

We can expect each Elemental Machine to be configured using the defined `nm-configurator` `_all.yaml` template.  

At the very first boot, the `elemental-register` will try to contact the Rancher API to register a new `MachineInventory`.  
At this stage the machine's network is not configured and will default to DHCP. It is a requirement that the machine is able to contact the Rancher API in this setup, otherwise the registration can not take place.  

Once the `MachineInventory` has first been registered, the `elemental-operator` is going to claim an `IPAddress` for each `IPPool` reference defined in the network configuration.  
On the `MachineInventory`, this will look like the following:  

```yaml
apiVersion: elemental.cattle.io/v1beta1
kind: MachineInventory
metadata:
  finalizers:
  - machineinventory.elemental.cattle.io
  name: m-e5331e3b-1e1b-4ce7-b080-235ed9a6d07c
  namespace: fleet-default
spec:
  ipAddressClaims:
    inventory-ip:
      apiVersion: ipam.cluster.x-k8s.io/v1beta1
      kind: IPAddressClaim
      name: m-e5331e3b-1e1b-4ce7-b080-235ed9a6d07c-inventory-ip
      namespace: fleet-default
      uid: 78f2d07a-7b6d-4b58-b615-c4108b7964b9
  ipAddressPools:
    inventory-ip:
      apiGroup: ipam.cluster.x-k8s.io
      kind: InClusterIPPool
      name: elemental-inventory-pool
  network:
    config:
      dns-resolver:
        config:
          search: []
          server:
          - 192.168.122.1
      interfaces:
      - description: Main-NIC
        ipv4:
          address:
          - ip: "{inventory-ip}"
            prefix-length: 24
          dhcp: false
          enabled: true
        ipv6:
          enabled: false
        name: eth0
        state: up
        type: ethernet
      routes:
        config:
        - destination: 0.0.0.0/0
          metric: 150
          next-hop-address: 192.168.122.1
          next-hop-interface: eth0
          table-id: 254
    ipAddresses:
      inventory-ip: 192.168.122.150
status:
  conditions:
  - lastTransitionTime: "2024-07-30T11:50:47Z"
    message: NetworkConfig is ready
    reason: ReconcilingNetworkConfig
    status: "True"
    type: NetworkConfigReady
```

You will notice that the `MachineInventory` carries the same `network.config` as the `MachineRegistration`, however instead of referencing IPAddressPools, we now have a map of real IPAddresses:  

```yaml
    ipAddresses:
      inventory-ip: 192.168.122.150
```

This `inventory-ip` will then be substituted in the `nm-configurator` config whenever `{inventory-ip}` has been defined.  

Also note that the `MachineInventory` references and owns each `IPAddressClaim` associated with it. Each claim follow the predictable `$MachineInventoryName-$IPPoolRefKey` naming convention: `m-e5331e3b-1e1b-4ce7-b080-235ed9a6d07c-inventory-ip`.  
These claims will follow the lifecycle of the `MachineInventory` object and be deleted on cascade, for example during the [reset workflow](./reset.md).  

If the `IPAddresses` can not be claimed, the `NetworkConfigReady` condition will be `False`, preventing the machine from completing installation. This can be the case if the `IPPool` has no more `IPAddresses` available.  

## On the machine side

During the installation phase, the `elemental-register` process running on the machine will receive the `nm-configurator` `_all.yaml` config template and the list of claimed IPAddresses with their keys. This information will be digested to an applicable `nm-configurator` configuration:

```yaml
config:
  dns-resolver:
    config:
      search: []
      server:
      - 192.168.122.1
  interfaces:
  - description: Main-NIC
    ipv4:
      address:
      - ip: "192.168.122.150"
        prefix-length: 24
      dhcp: false
      enabled: true
    ipv6:
      enabled: false
    name: eth0
    state: up
    type: ethernet
  routes:
    config:
    - destination: 0.0.0.0/0
      metric: 150
      next-hop-address: 192.168.122.1
      next-hop-interface: eth0
      table-id: 254
```

The `elemental-register` will then invoke `nmc generate` and `nmc apply` to apply this configuration into the running system.  
From this moment until reset, the machine will always use the applied configuration.  

Also note that outside of installation and reset, `nm-configurator` is no longer used, since the `elemental-register` will persist the `/etc/NetworkManager/system-connection/*.nmconnection` files generated by `nmc` rather than the `nmc` configuration itself.  

For example on any running system, you will find a [yip](https://github.com/rancher/yip) configuration file (`/oem/elemental-network.yaml`) to apply the desired `nmconnections`, for example:  

```yaml
name: Apply network config
stages:
    initramfs:
        - files:
            - path: /etc/NetworkManager/system-connections/Wired connection 1.nmconnection
              permissions: 384
              owner: 0
              group: 0
              content: |
                [connection]
                id=Wired connection 1
                uuid=d26b4ae4-d525-3cbf-a557-33feb60343c0
                type=ethernet
                autoconnect-priority=-999
                interface-name=eth0
                timestamp=1722340245

                [ethernet]

                [ipv4]
                address1=192.168.122.150/24
                dhcp-timeout=2147483647
                dns=192.168.122.1;
                dns-options=
                dns-priority=40
                method=manual
                route1=0.0.0.0/0,192.168.122.1,150
                route1_options=table=254

                [ipv6]
                addr-gen-mode=eui64
                dhcp-timeout=2147483647
                method=disabled

                [proxy]

                [user]
                nm-configurator.interface.description=Main-NIC
              encoding: ""
              ownerstring: ""
```

### During reset

Whenever [reset](./reset.md) is triggered, the `elemental-register` running on the machine will clear any `/etc/NetworkManager/system-connection/*.nmconnection` file and restart the network stack. The machine should then revert to DHCP and after that confirm to the `elemental-operator` on the management side, that reset has succeeded.  
Note that this only applies whenever the `MachineInventory.spec.network.configurator` value is different than `none`. Otherwise no action will be taken to reset network during the machine reset phase.  

Following network reset, the machine should reboot into recovery mode, perform the actual reset and receive a fresh network configuration to be applied. Potentially this will be the same as before (if the `MachineRegistration` has not been updated), or it may have different IPs since the previous ones may have been claimed by other machines in the meanwhile.  

If reverting to DHCP failed or the machine is anyhow unable to contact the Rancher API back for confirmation, you will notice that the `MachineInventory` will not be deleted, despite having a deletion timestamp.  
Since the machine has now network issue, it won't be possible to remotely fix it.  

You have the option to physically reach the machine or by any means fix the DHCP driven network configuration, or alternatively you can remove the `machineinventory.elemental.cattle.io` finalizer from the `MachineInventory`, to allow deletion, if you intend to decommission the machine.  


---

## Article: troubleshooting-rancher-upgrades.md

---
sidebar_label: Rancher upgrades
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/troubleshooting-rancher-upgrades"/>
</head>

# Troubleshooting Rancher upgrades

:::warning warning
Upgrading to Rancher v2.7.2 will fail if Elemental clusters are defined. The rancher pod gets stuck in a crash loop (see https://github.com/rancher/rancher/issues/41145).
:::

Note that the issue is present only if at least one Elemental cluster is defined.

To workaround the issue create an empty `dynamicschemas.management.cattle.io` resource named `machineinventoryselectortemplate`:

```shell showLineNumbers
kubectl apply -f - <<EOF
kind: DynamicSchema
apiVersion: management.cattle.io/v3
metadata:
  name: machineinventoryselectortemplate
EOF

```

The crash loop will stop and Rancher upgrade will complete successfully.

---

## Article: troubleshooting-reset.md

---
sidebar_label: Reset
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/troubleshooting-reset"/>
</head>

# Troubleshooting reset

Each `MachineInventory` with the `elemental.cattle.io/resettable: "true"` annotation will trigger the execution of a reset plan, upon deletion.  
The `machineinventory.elemental.cattle.io` finalizer is going to be removed only after the plan has been executed successfully by the `elemental-system-agent` running on the machine.  

You can investigate why a `MachineInventory` has not been deleted yet, by examining it:

```yaml
apiVersion: elemental.cattle.io/v1beta1
kind: MachineInventory
metadata:
  # deletionTimestamp has been set. This object has been marked for deletion.
  deletionTimestamp: "2023-08-04T08:41:25Z"
  annotations:
    # `resettable` annotation is enabled.
    # This means the machine has to go through reset, before deletion of this object.
    elemental.cattle.io/resettable: "true"
  # `machineinventory.elemental.cattle.io` finalizer is set.
  #  The `elemental-operator` is going to create a reset plan for the machine to execute.
  #  After successful execution of the reset plan, the finalizer is removed and the object will be deleted.
  finalizers:
  - machineinventory.elemental.cattle.io
status:
  conditions:
  # Most recent condition shows that the MachineInventory is waiting for a plan to be applied.
  - lastTransitionTime: "2023-08-04T08:41:25Z"
    message: waiting for plan to be applied
    reason: WaitingForPlan
    status: "False"
    type: Ready
  # The plan to be executed is referenced. 
  # Normally it has the same name of the MachineInventory and lives within the same namespace.   
  plan:
    checksum: 5aba8b6b3161bc52d8953b2428e54ecda3b59e8e0043b49d761d1e79174eded6
    secretRef:
      name: m-bf1008a1-61d6-4355-b5f5-f7d1c527affe
      namespace: fleet-default
```

You can also examine the referenced plan `Secret`.  
Note that the `elemental-system-agent` running on the machine is watching this secret and it should execute the plan.  
You can also monitor its progress from the machine logs: `journalctl -u elemental-system-agent -f`.

```yaml
apiVersion: v1
kind: Secret
# This is a `elemental.cattle.io/plan` secret plan.
type: elemental.cattle.io/plan
metadata:
  annotations:
    # This is a `reset` plan type.
    elemental.cattle.io/plan.type: reset
  labels:
    elemental.cattle.io/managed: "true"
  name: m-bf1008a1-61d6-4355-b5f5-f7d1c527affe
  namespace: fleet-default
  # It is owned by the `MachineInventory` waiting for deletion.
  ownerReferences:
  - apiVersion: elemental.cattle.io/v1beta1
    controller: true
    kind: MachineInventory
    name: m-bf1008a1-61d6-4355-b5f5-f7d1c527affe
    uid: 5aa3863c-63a5-4cb9-91fd-7a45191d4842
data:
  # The plan has not been applied yet.
  applied-checksum: ""
  # It also hasn't failed.
  failed-checksum: ""
  # The actual plan to be executed, base64 encoded.
  plan: eyJmaWxlcyI6W3siY29udGVudCI6ImJtRnRaVG9nUld4bGJXVnVkR0ZzSUZKbGMyVjBDbk4wWVdkbGN6b0tJQ0FnSUc1bGRIZHZjbXN1WVdaMFpYSTZDaUFnSUNBZ0lDQWdMU0JqYjIxdFlXNWtjem9LSUNBZ0lDQWdJQ0FnSUNBZ0xTQmxiR1Z0Wlc1MFlXd3RjbVZuYVhOMFpYSWdMUzFrWldKMVp5QXRMWEpsYzJWMENpQWdJQ0FnSUNBZ0lDQnBaam9nSjFzZ0xXWWdMM0oxYmk5amIzTXZjbVZqYjNabGNubGZiVzlrWlNCZEp3b2dJQ0FnSUNBZ0lDQWdibUZ0WlRvZ1VuVnVjeUJsYkdWdFpXNTBZV3dnY21WelpYUUsiLCJwYXRoIjoiL29lbS9yZXNldC1jbG91ZC1jb25maWcueWFtbCIsInBlcm1pc3Npb25zIjoiMDYwMCJ9XSwiaW5zdHJ1Y3Rpb25zIjpbeyJuYW1lIjoiY29uZmlndXJlIG5leHQgYm9vdCB0byByZWNvdmVyeSBtb2RlIiwiYXJncyI6WyIvb2VtL2dydWJlbnYiLCJzZXQiLCJuZXh0X2VudHJ5PXJlY292ZXJ5Il0sImNvbW1hbmQiOiJncnViMi1lZGl0ZW52In0seyJuYW1lIjoic2NoZWR1bGUgcmVib290IiwiYXJncyI6WyItciIsIisxIl0sImNvbW1hbmQiOiJzaHV0ZG93biJ9XX0K
```

The plan created by the `elemental-operator` should contain the following instructions:

```json
{
  "files": [
    // A cloud-init config file is created on the default /oem directory.  
    // This config will be executed once in recovery mode.
    {
      "content": "bmFtZTogRWxlbWVudGFsIFJlc2V0CnN0YWdlczoKICAgIG5ldHdvcmsuYWZ0ZXI6CiAgICAgICAgLSBjb21tYW5kczoKICAgICAgICAgICAgLSBlbGVtZW50YWwtcmVnaXN0ZXIgLS1kZWJ1ZyAtLXJlc2V0CiAgICAgICAgICBpZjogJ1sgLWYgL3J1bi9jb3MvcmVjb3ZlcnlfbW9kZSBdJwogICAgICAgICAgbmFtZTogUnVucyBlbGVtZW50YWwgcmVzZXQK",
      "path": "/oem/reset-cloud-config.yaml",
      "permissions": "0600"
    }
  ],
  "instructions": [
    {
      "name": "configure next boot to recovery mode",
      "args": [
        "/oem/grubenv",
        "set",
        "next_entry=recovery"
      ],
      "command": "grub2-editenv"
    },
    {
      "name": "schedule reboot",
      "args": [
        "-r",
        "+1"
      ],
      "command": "shutdown"
    }
  ]
}
```

If the `elemental-system-agent` successfully executed the plan, the `machineinventory.elemental.cattle.io` finalizer on the `MachineInventory` will be removed and the `MachineInventory` will be deleted.  
Note that this is not an indication that the machine has been fully reset yet.  
This is a limitation of the current implementation and it will eventually improve, so that it will be possible to completely track the reset status.  

However, at this stage we do expect the host to undergo reboot, and reboot in recovery mode.  
Once in recovery mode, the `cos-setup-network` should execute the cloud-init config that has been written on `/oem/reset-cloud-config.yaml`.  
You can monitor the status with `journalctl -u cos-setup-network -f`.

The cloud-init instructions should look like the following:

```yaml
name: Elemental Reset
stages:
    network.after:
        - if: '[ -f /run/cos/recovery_mode ]'
          name: Runs elemental reset
          commands:
            - systemctl start elemental-register-reset
```

The `elemental-register` cli will register with the `elemental-operator` as a new machine. This will lead to the creation of a new `MachineInventory` object.  
The remote `MachineRegistration` configuration will also be fetched to apply the reset options, for example `reset-persistent`, `reset-oem`, or the power settings, either `reboot` or `power-off`.  
After reset, depending on the settings, the machine should either shut down or reboot and be ready to be adopted within a new cluster.  

## Forcefully deleting a MachineInventory undergoing reset

If the machine is unable to execute the reset instructions and the related `MachineInventory` is not deleted, there are two equivalent ways you can manually fix the issue.

- Remove the `elemental.cattle.io/resettable: "true"` annotation from the `MachineInventory`.  
- Remove the `machineinventory.elemental.cattle.io` finalizer from the `MachineInventory`.  

Remember to also take care of the machine itself, by fully reprovisioning it or rebooting into recovery mode and using the `elemental reset` command directly.  


---

## Article: troubleshooting-restore.md

---
sidebar_label: Restore
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/troubleshooting-restore"/>
</head>

# Troubleshooting restore

:::warning warning
When a restore is performed, do not restart the `rancher-system-agent` on elemental nodes as it can stale and end with the following error:

`panic: error while connecting to Kubernetes cluster: the server has asked for the client to provide credentials`

If you face this problem, please follow the procedure below.
:::

Before you initiate a restore, you need to copy `/var/lib/rancher/agent/rancher2_connection_info.json` from the elemental node to a place where you have access with Rancher UI.

Once the file is copied, download the `rancher-agent-token-update.sh` script from the [Elemental repository](https://github.com/rancher/elemental):

```shell showLineNumbers
wget -q https://raw.githubusercontent.com/rancher/elemental/main/scripts/rancher-agent-token-update && chmod +x rancher-agent-token-update
```

Execute the script without any additional options:

```shell showLineNumbers
./rancher-agent-token-update
```

After the restore successfully completed, copy `rancher2_connection_info.json` back to the elemental node to the path
`/var/lib/rancher/agent/rancher2_connection_info.json`. Finally, restart the `rancher-system-agent` service:

```shell showLineNumbers
systemctl restart rancher-system-agent
```


---

## Article: troubleshooting-support.md

---
sidebar_label: Support
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/troubleshooting-support"/>
</head>

# Elemental Support

The `elemental-support` is a utility program that collects all information about the machine where the Elemental stack is running.  

When troubleshooting issues with any Elemental node, it is possible to generate a compressed file containing logs and config files.  
Be aware that **the tar.gz generated file may contain sensitive information**, like access tokens to the management cluster and more.  

## Requirements

Optionally, the `elemental-support` utility needs `kubectl` to be available on the target machine, in order to collect node information. This can be very useful to debug Elemental upgrade issues with a node.  
The kubectl config used will be looked up from the `KUBECONFIG` environment variable, or in its absence either the `/etc/rancher/k3s/k3s.yaml` or `/etc/rancher/rke2/rke2.yaml` files will be used instead.  

## Collecting information

Once ran, the `elemental-support` will create the compressed archive within the same directory.  

Example run:  

```bash
m-993369ec-4b3f-4368-86c3-1449ebf57cf2# elemental-support 
I0430 11:33:38.934777    8536 log.go:42] Support version 1.3.4, commit 192dc33d, commit date 20230831
I0430 11:33:38.934865    8536 log.go:46] Getting elemental-support version
I0430 11:33:38.934918    8536 log.go:42] Copying /etc/os-release
I0430 11:33:38.935008    8536 log.go:42] Copying /etc/resolv.conf
I0430 11:33:38.935066    8536 log.go:42] Copying /etc/hostname
I0430 11:33:38.935131    8536 log.go:42] Copying /etc/rancher/agent/config.yaml
I0430 11:33:38.935174    8536 log.go:42] Copying /etc/rancher/agent/config.yaml
I0430 11:33:38.935214    8536 log.go:42] Copying dir /var/lib/elemental/agent/applied/
I0430 11:33:38.935434    8536 log.go:42] Copying dir /var/lib/rancher/agent/applied/
I0430 11:33:38.935572    8536 log.go:42] Copying dir /system/oem
I0430 11:33:38.936584    8536 log.go:42] Copying dir /oem/
I0430 11:33:38.937121    8536 log.go:42] Getting service log elemental-system-agent
I0430 11:33:38.959813    8536 log.go:42] Getting service log rancher-system-agent
I0430 11:33:38.979271    8536 log.go:42] Getting service log k3s
I0430 11:33:39.018177    8536 log.go:42] Getting service log rke2
I0430 11:33:39.036852    8536 log.go:42] Getting service log cos-setup-boot
I0430 11:33:39.053409    8536 log.go:42] Getting service log cos-setup-fs
I0430 11:33:39.068132    8536 log.go:42] Getting service log cos-setup-initramfs
I0430 11:33:39.085935    8536 log.go:42] Getting service log cos-setup-network
I0430 11:33:39.101440    8536 log.go:42] Getting service log cos-setup-reconcile
I0430 11:33:39.119528    8536 log.go:42] Getting service log cos-setup-rootfs
I0430 11:33:39.136151    8536 log.go:42] Getting service log cos-immutable-rootfs
I0430 11:33:39.153850    8536 log.go:42] Getting service log elemental
I0430 11:33:39.170564    8536 log.go:42] Getting service log NetworkManager
I0430 11:33:39.190169    8536 log.go:42] Getting elemental-cli version
I0430 11:33:39.205410    8536 log.go:42] Getting elemental-operator version
I0430 11:33:39.205621    8536 log.go:42] Getting elemental-register version
I0430 11:33:39.266263    8536 log.go:42] Found k3s kubeconfig at /etc/rancher/k3s/k3s.yaml
I0430 11:33:39.266312    8536 log.go:42] Found k3s kubectl at /usr/local/bin/kubectl
I0430 11:33:39.266320    8536 log.go:42] Getting k8s info for pods
I0430 11:33:39.770730    8536 log.go:42] Getting k8s info for secrets
I0430 11:33:39.962037    8536 log.go:42] Getting k8s info for nodes
I0430 11:33:40.118680    8536 log.go:42] Getting k8s info for services
I0430 11:33:40.280797    8536 log.go:42] Getting k8s info for deployments
I0430 11:33:40.458869    8536 log.go:42] Getting k8s info for plans
I0430 11:33:40.612138    8536 log.go:42] Getting k8s info for apps
I0430 11:33:40.787309    8536 log.go:42] Getting k8s info for jobs
I0430 11:33:40.932361    8536 log.go:42] Getting k8s logs for namespace cattle-system
I0430 11:33:41.353556    8536 log.go:42] Getting k8s logs for namespace kube-system
I0430 11:33:41.935885    8536 log.go:42] Getting k8s logs for namespace ingress-nginx
W0430 11:33:41.999345    8536 log.go:62] No pods in namespace ingress-nginx
I0430 11:33:41.999365    8536 log.go:42] Getting k8s logs for namespace calico-system
W0430 11:33:42.061902    8536 log.go:62] No pods in namespace calico-system
I0430 11:33:42.061927    8536 log.go:42] Getting k8s logs for namespace cattle-fleet-system
I0430 11:33:42.204459    8536 log.go:42] Creating final file
I0430 11:33:42.215815    8536 log.go:42] All done. File created at m-993369ec-4b3f-4368-86c3-1449ebf57cf2-2024-04-30T113342Z.tar.gz
```


---

## Article: troubleshooting-upgrade.md

---
sidebar_label: Upgrade
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/troubleshooting-upgrade"/>
</head>

# Troubleshooting upgrade

For a high level overview of the upgrade lifecycle and components, please refer to the [related document](./upgrade-lifecycle).  

## Rancher side

In this example we upgraded the cluster nodes with the following ManagedOSImage definition:

```yaml showLineNumbers
apiVersion: elemental.cattle.io/v1beta1
kind: ManagedOSImage
metadata:
  name: my-upgrade
  namespace: fleet-default
spec:
  # Set to the new Elemental version you would like to upgrade to or track the latest tag
  osImage: "registry.suse.com/rancher/sle-micro/5.5:latest"
  clusterTargets:
    - clusterName: my-cluster
```

Once the `ManagedOSImage` is applied, the `elemental-operator` will verify it and generate a related `Bundle`.  
The `Bundle` name will be prefixed with `mos` and then the `ManagedOSImage` name. In this case `mos-my-upgrade`.  

In the `Bundle` definition, you will find the details about the upgrade plan and the desired target.  
For example:

```shell showLineNumbers
kubectl -n fleet-default get bundle mos-my-upgrade -o yaml
```

<details>
  <summary>Example</summary>

```yaml showLineNumbers
apiVersion: fleet.cattle.io/v1alpha1
kind: Bundle
metadata:
  creationTimestamp: "2023-06-16T09:01:47Z"
  generation: 1
  name: mos-my-upgrade
  namespace: fleet-default
  ownerReferences:
  - apiVersion: elemental.cattle.io/v1beta1
    controller: true
    kind: ManagedOSImage
    name: my-upgrade
    uid: e468ed21-23bb-487a-a022-dbc7ef753720
  resourceVersion: "1038645"
  uid: 35e83fc4-28c8-4b10-8059-cae6cdff2cda
spec:
  resources:
  - content: '{"kind":"ClusterRole","apiVersion":"rbac.authorization.k8s.io/v1","metadata":{"name":"os-upgrader-my-upgrade","creationTimestamp":null},"rules":[{"verbs":["update","get","list","watch","patch"],"apiGroups":[""],"resources":["nodes"]},{"verbs":["list"],"apiGroups":[""],"resources":["pods"]}]}'
    name: ClusterRole--os-upgrader-my-upgrade-296a3abf3451.yaml
  - content: '{"kind":"ClusterRoleBinding","apiVersion":"rbac.authorization.k8s.io/v1","metadata":{"name":"os-upgrader-my-upgrade","creationTimestamp":null},"subjects":[{"kind":"ServiceAccount","name":"os-upgrader-my-upgrade","namespace":"cattle-system"}],"roleRef":{"apiGroup":"rbac.authorization.k8s.io","kind":"ClusterRole","name":"os-upgrader-my-upgrade"}}'
    name: ClusterRoleBinding--os-upgrader-my-upgrade-f63eaecde935.yaml
  - content: '{"kind":"ServiceAccount","apiVersion":"v1","metadata":{"name":"os-upgrader-my-upgrade","namespace":"cattle-system","creationTimestamp":null}}'
    name: ServiceAccount-cattle-system-os-upgrader-my-upgrade-ce93d-01096.yaml
  - content: '{"kind":"Secret","apiVersion":"v1","metadata":{"name":"os-upgrader-my-upgrade","namespace":"cattle-system","creationTimestamp":null},"data":{"cloud-config":""}}'
    name: Secret-cattle-system-os-upgrader-my-upgrade-a997ee6a67ef.yaml
  - content: '{"kind":"Plan","apiVersion":"upgrade.cattle.io/v1","metadata":{"name":"os-upgrader-my-upgrade","namespace":"cattle-system","creationTimestamp":null},"spec":{"concurrency":1,"nodeSelector":{},"serviceAccountName":"os-upgrader-my-upgrade","version":"latest","secrets":[{"name":"os-upgrader-my-upgrade","path":"/run/data"}],"tolerations":[{"operator":"Exists"}],"cordon":true,"upgrade":{"image":"registry.suse.com/suse/sle-micro/5.5","command":["/usr/sbin/suc-upgrade"]}},"status":{}}'
    name: Plan-cattle-system-os-upgrader-my-upgrade-273c2c09afca.yaml
  targets:
  - clusterName: my-cluster
.
.
.
```

</details>

## Elemental Cluster side

Any Elemental node correctly registered and part of the target cluster will fetch the bundle and start applying it.  
This operation is performed by the Rancher's `system-upgrade-controller` running on the Elemental Cluster.  
To monitor the correct operation of this controller, you can read its logs:

```shell showLineNumbers
kubectl -n cattle-system logs deployment/system-upgrade-controller
```

If everything is correct, the `system-upgrade-controller` will create an upgrade Plan on the cluster:

```shell
kubectl -n cattle-system get plans
```

For each Plan, the controller will orchestrate the jobs that will apply it on each targeted node.  
The job names will use the Plan name (`os-upgrader-my-upgrade`) and the target machine hostname (`my-host`) for easy discoverability.  
For example: `apply-os-upgrader-my-upgrade-on-my-host-7a25e`  
You can monitor these jobs with:

```shell showLineNumbers
kubectl -n cattle-system get jobs
```

Each job will use a `privileged: true` container with the SLE Micro image specified in the `ManagedOSImage` definition. This container will try to upgrade the system and perform a reboot.  

If the job fails, you can check its status by examining the logs:

```shell showLineNumbers
kubectl -n cattle-system logs job.batch/apply-os-upgrader-my-upgrade-on-my-host-7a25e
```

:::info Two stages job process

Note that the upgrade process is performed in two stages.  
You will notice that the same job is ran twice and the first one ends with the `Unknown` Status and will not complete.  
**This is to be expected**, as Elemental relies on the job to be ran again after the machine restarts, so that it can verify the new version was installed correctly.  
You will notice a second run of the job, this time completing correctly.

```shell showLineNumbers
kubectl -n cattle-system get jobs 
NAMESPACE       NAME                                            COMPLETIONS   DURATION   AGE
cattle-system   apply-os-upgrader-my-upgrade-on-my-host-0b392   1/1           2m34s      6m23s
cattle-system   apply-os-upgrader-my-upgrade-on-my-host-7a25e   0/1           6m23s      6m23s
```

```shell showLineNumbers
kubectl -n cattle-system get pods 
NAME                                            READY   STATUS      RESTARTS      AGE
apply-os-upgrader-my-upgrade-on-my-host-zbkrh   0/1     Completed   0             9m40s
apply-os-upgrader-my-upgrade-on-my-host-zvrff   0/1     Unknown     0             12m
```

:::

## Recovering from failures

It is possible that the ManagedOSImage upgrade process will fail, leaving one or more nodes in a faulty state.  
For example if the to-be-upgraded image is not found on the registry or otherwise broken, the upgrade job running on the downstream clusters will not succeed.  

When this is the case, the nodes running the failing upgrade job will stay cordoned.  
You can update the ManagedOSImage with a functioning `osImage`, or alternatively delete it to stop any further upgrade attempt.  
In any case, in order to restore them and be able to also schedule following upgrades, the affected nodes need to be uncordoned manually.  


---

## Article: unmanaged-os.md

---
sidebar_label: Register an Unmanaged OS
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/unmanaged-os"/>
</head>

# Register an Unmanaged OS

Normally, the <Vars name="elemental_operator_name" link="elemental_operator_url" /> manages operating systems that are installed and configured by the <Vars name="elemental_toolkit_name" link="elemental_toolkit_url" />.  
For example, in order to automate OS installation, upgrade, and reset, the `elemental-register` relies on the `elemental` CLI to execute these operations.  

However, it is also possible to register and provision Elemental "Toolkitless" systems.  

In this scenario, [elemental-register](architecture-components.md#elemental-register-command-line-tool) needs to be installed on the system.  
Optionally, the [elemental-system-agent](architecture-components.md#elemental-system-agent-daemon) can be installed. Note that without the `elemental-system-agent`, the <Vars name="elemental_operator_name" /> will not be able to provision any k8s cluster on the machine. In this case the <Vars name="elemental_operator_name" /> can only be used for OS Inventory purposes only.  

Finally, on the management cluster, the `MachineRegistration` must enable the `spec.config.elemental.registration.no-toolkit` flag.  

Once `no-toolkit` is enabled on the `MachineRegistration`, and a new registration occurs using `elemental-register --install` on the system, a new `MachineInventory` will be created on the management cluster:  

```bash
kubectl -n fleet-default describe machineinventory my-unmanaged-os-machine
```

The `MachineInventory` will be annotated with the `elemental.cattle.io/os.unmanaged: "true"` annotation, highlighting that this machine is not managed and has limited functionality.  

On the system, upon successful registration, the `/etc/rancher/elemental/agent/config.yaml` and `/var/lib/elemental/agent/elemental_connection.json` files are automatically created to configure the `elemental-system-agent`.
The `elemental-system-agent` component is needed for K8s provisioning and Reset triggers.  

When a [machine reset](reset.md) is triggered, for example by deleting the `MachineInventory` directly, the `elemental-system-agent` will execute a simple reset plan that will create the `/var/lib/elemental/.unmanaged_reset` sentinel file.  
The presence of this file indicates that the machine needs a reset. This may involve stopping services, uninstalling packages, formatting devices, and so on, depending how the machine is customly managed.  

## Supported Features

- Registration of a `MachineInventory`
- K8s provisioning (when the [elemental-system-agent](architecture-components.md#elemental-system-agent-daemon) is installed and running on the machine)
- Reset triggers (when the [elemental-system-agent](architecture-components.md#elemental-system-agent-daemon) is installed and running on the machine)

## Unsupported Features

- Cloud-init driven [configuration](cloud-config-reference.md)
- OS [Upgrades](upgrade.md)
- OS [Reset](reset.md)


---

## Article: upgrade-lifecycle.md

---
sidebar_label: Upgrade Lifecycle
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/upgrade-lifecycle"/>
</head>

# Upgrade Lifecycle

![Upgrade Flow](images/upgrade-lifecycle.png) 

## Components

- [Elemental Operator](https://elemental.docs.rancher.com/upgrade)

The Elemental Operator supports several ways to configure the OS version on a Cluster level.  
Whenever a new `ManagedOSIMage` is defined, a related Fleet `Bundle` is also generated by the Elemental Operator.  

- [Fleet](https://fleet.rancher.io/)

In this context, Fleet is used to apply an `upgrade.cattle.io/v1` `Plan` on each Cluster targeted by the `ManagedOSImage`.  
The `Plan` is then executed by the `System Upgrade Controller` running on the downstream Elemental Cluster.  

- [System Upgrade Controller](https://github.com/rancher/system-upgrade-controller)

The System Upgrade Controller is an upgrade tool used by [k3s](https://docs.k3s.io/upgrades/automated) and [rke2](https://docs.rke2.io/upgrade/automated_upgrade) clusters. It should always be automatically installed on each provisioned Elemental cluster.  
In this context the System Upgrade Controller orchestrates the application of `upgrade.cattle.io/v1` `Plan` to the Elemental nodes.  

## Versioning and components lifecycle

The [Fleet](https://fleet.rancher.io/) and [System Upgrade Controller](https://github.com/rancher/system-upgrade-controller) versions are controlled by Rancher.  

On the Rancher cluster, you can run:  

```bash
kubectl get settings fleet-version system-upgrade-controller-chart-version
```

For more information about chart versions, you can visit the [repository](https://github.com/rancher/charts) or read the [documentation](https://ranchermanager.docs.rancher.com/how-to-guides/new-user-guides/helm-charts-in-rancher).

The charts version is determined by Rancher.  
Depending on how Rancher is installed, the following environment variables can be set, for example when installing the [Rancher Helm chart](https://ranchermanager.docs.rancher.com/getting-started/installation-and-upgrade/installation-references/helm-chart-options#setting-extra-environment-variables):  

```bash
CATTLE_SYSTEM_UPGRADE_CONTROLLER_CHART_VERSION
CATTLE_FLEET_VERSION
```

Manually updating the versions is generally not recommended.  
Also beware that a new `CATTLE_SYSTEM_UPGRADE_CONTROLLER_CHART_VERSION` will trigger a `system-upgrade-controller` chart update on all provisioned downstream clusters.  

## Verifying successful upgrades

Verifying the correct execution of the upgrade jobs is a fundamental part of ensuring the Elemental upgrade process is complete.  

To troubleshoot each step within the entire upgrade process, please consult the [related document](./troubleshooting-upgrade.md).  
This section only describes the last step needed to verify the correct application of the upgrade Plan on one targeted cluster.  

- Versions higher or equal to **v0.13.4** of the `system-upgrade-controller`:  

  ```shell
  kubectl -n cattle-system describe plan os-upgrader-my-upgrade
  ```

  Once the Plan is applied to all nodes, it will have a `Complete` status condition.

- Older versions:

  Each Plan has a `latestHash` status value that should match the Plan label on each node.  
  A simple script can be used to list all nodes where the plan has not be applied yet:  

  ```shell
  #!/bin/bash

  # This script prints a list of all nodes where an upgrade plan was not applied.  
  # To further determine the cause of failures, you can analyze the Plan status and the related jobs.
  # For ex: kubectl -n cattle-system get plans,jobs

  # Edit this variable according to your ManagedOSImage name.
  MANAGED_OS_IMAGE_NAME=my-upgrade

  PLAN_NAME=os-upgrader-$MANAGED_OS_IMAGE_NAME
  LATEST_HASH=$(kubectl -n cattle-system get plan $PLAN_NAME -o=jsonpath='{.status.latestHash}')

  printf "Plan '$PLAN_NAME' latest hash is: '$LATEST_HASH'\n"

  PLAN_LABEL=plan.upgrade.cattle.io/$PLAN_NAME
  printf "Listing all nodes with mismatching or missing label '$PLAN_LABEL':\n"
  kubectl get nodes -l $PLAN_LABEL!=$LATEST_HASH
  ```


---

## Article: upgrade.md

---
sidebar_label: Upgrade
title: ''
---

<head>
  <link rel="canonical" href="https://elemental.docs.rancher.com/upgrade"/>
</head>

import ClusterTarget from "!!raw-loader!@site/examples/upgrade/upgrade-cluster-target.yaml"
import NodeSelector from "!!raw-loader!@site/examples/upgrade/upgrade-node-selector.yaml"
import UpgradeForce from "!!raw-loader!@site/examples/upgrade/upgrade-force.yaml"
import UpgradeRecovery from "!!raw-loader!@site/examples/upgrade/upgrade-recovery.yaml"
import ManagedOSVersion from "!!raw-loader!@site/examples/upgrade/upgrade-managedos-version.yaml"

# Upgrade

All components in Elemental are managed using Kubernetes. Below is how
to use Kubernetes approaches to upgrade the components.

## Elemental node upgrade

Elemental nodes are upgraded with the <Vars name="elemental_operator_name" />. Refer to the
[<Vars name="elemental_operator_name" />](elementaloperatorchart-reference.md) documentation for complete information.

Upgrade can be achieve either with CLI or UI:

- [Command Line Interface](#upgrade-via-command-line-interface)
- [User Interface](#upgrade-via-user-interface)
## Upgrade via command line interface

There are two ways of selecting nodes for upgrading. Via a cluster target, which will match ALL nodes in a cluster that matches our
selector or via node selector, which will match nodes based on the node labels. Node selector allows us to be more targeted with the upgrade
while cluster selector just selects all the nodes in a matched cluster.  

Updating an existing `ManagedOSImage` will trigger a new upgrade cycle, to reconcile the configured image (or image version) to all targeted nodes.  

<Tabs>
<TabItem value="clusterTarget" label="With 'clusterTarget'" default>
You can target nodes for an upgrade via a `clusterTarget` by setting it to the cluster name that you want to upgrade.
All nodes in a cluster that matches that name will match and be upgraded.

<CodeBlock language="yaml" title="upgrade-cluster-target.yaml" showLineNumbers>{ClusterTarget}</CodeBlock>

</TabItem>
<TabItem value="nodeSelector" label="With nodeSelector">
You can target nodes for an upgrade via a `nodeSelector` by setting it to the label and value that you want to match.
Any nodes containing that key with the value will match and be upgraded.

<CodeBlock language="yaml" title="upgrade-node-selector.yaml" showLineNumbers>{NodeSelector}</CodeBlock>

</TabItem>
<TabItem value="forceUpgrade" label="With FORCE flag">
When upgrading to an older version or the same version that is already running the upgrade-procedure will be skipped.

It is possible to force upgrades to older versions by setting the FORCE environment variable as shown below.

<CodeBlock language="yaml" title="upgrade-force.yaml" showLineNumbers>{UpgradeForce}</CodeBlock>

</TabItem>

<TabItem value="recoveryUpgrade" label="With UPGRADE_RECOVERY flag">
You can decide upgrade the Recovery partition when upgrading the system, or alternatively to upgrade the Recovery partition only.

<CodeBlock language="yaml" title="upgrade-recovery.yaml" showLineNumbers>{UpgradeRecovery}</CodeBlock>

</TabItem>
</Tabs>

### Selecting source for upgrade 

<Tabs>
<TabItem value="osImage" label="Via 'osImage'">

Just specify an OCI image on the `osImage` field

<CodeBlock language="yaml" title="upgrade-cluster-target.yaml" showLineNumbers>{ClusterTarget}</CodeBlock>

</TabItem>
<TabItem value="managedOSVersion" label="Via 'ManagedOSVersion'">

In this case we use the auto populated `ManagedOSVersion` resources to set the wanted `managedOSVersionName` field.
See section [Managing available versions](#managing-available-versions) to understand how the `ManagedOSVersion` are managed.

<CodeBlock language="yaml" title="upgrade-managedos-version.yaml" showLineNumbers>{ManagedOSVersion}</CodeBlock>

</TabItem>
</Tabs>

:::warning Warning
If both `osImage` and `ManagedOSVersion` are defined in the same `ManagedOSImage` be aware that `osImage` takes precedence.
:::

### Managing available versions

It is possible to create [ManagedOSVersions](./managedosversion-reference.md) directly, or to subscribe to [ManagedOSVersionChannels](./managedosversionchannel-reference.md) to automatically sync `ManagedOSVersions` from them.  

For more details and a list of available channels, or to even make your own, please read the [documentation](./channels.md).  

## Upgrade via user interface

To upgrade via the UI, you have to go in the Elemental Advanced menu, then click on `Update Groups`.

Choose a name, select clusters to target and choose between the two upgrade ways:

![Elemental Upgrade Menu](images/upgrade-ui-menu.png)


<Tabs>

<TabItem value="managedOSVersion" label="Via Managed OS Version">

In this case, a `OS Version Channels` is used to auto populated `OS Versions` resources.

The channel bellow is provide by us but you can bring your own channel as well.

See section [Managing available versions](#managing-available-versions) to understand how the `ManagedOSVersion` are managed.

![Create OS Version Channel](images/upgrade-ui-version-channel.png)

After a short syncing time, you will see your `OS Versions` appears in the `OS Versions` menu.

![Elemental OS Version menu](images/upgrade-ui-os-version.png)

Finally, you can select the `OS Versions` when you create your `Upgrade Group` according to the following screenshot:

![Select OS Version in Upgrade Group](images/upgrade-ui-upgrade-group-channel.png)

</TabItem>
<TabItem value="imageFromRegistry" label="Via Image from registry">

Just specify an OCI image on the `Image path` field to upgrade to:

![Upgrade via Image Registry](images/upgrade-ui-image-registry.png)

</TabItem>
</Tabs>

Click on the `Create` button to start the upgrade process, if you have multiple nodes, the upgrade will be done sequentially node by node.


---

## Article: partials/_elemental-operator-install.md

## Install Elemental Operator

`elemental-operator` is the management endpoint, running the management
cluster and taking care of creating inventories, registrations for machines and much more.

We will use the Helm package manager to install the elemental-operator chart into our cluster.

```shell showLineNumbers
helm upgrade --create-namespace -n cattle-elemental-system --install elemental-operator-crds oci://registry.suse.com/rancher/elemental-operator-crds-chart
helm upgrade --create-namespace -n cattle-elemental-system --install elemental-operator oci://registry.suse.com/rancher/elemental-operator-chart
```

Now after a few seconds you should see the operator pod appear on the `cattle-elemental-system` namespace:

```shell showLineNumbers
kubectl get pods -n cattle-elemental-system
NAME                                  READY   STATUS    RESTARTS   AGE
elemental-operator-64f88fc695-b8qhn   1/1     Running   0          16s
```

:::info Helm v3.8.0+ required
The Elemental Operator chart is distributed via an OCI registry: Helm correctly supports OCI based registries starting from the v3.8.0 release.
:::

:::warning Swap charts installation order when upgrading from elemental-operator release < 1.2.4
When upgrading from an elemental-operator release embedding the Elemental CRDs (version < 1.2.4) the elemental-operator-crds chart installation will fail.
You will need to upgrade the elemental-operator chart first, and only then install the elemental-operator-crds chart.
:::

### Non-stable installations

Besides the Helm charts listed above, there are two other `non-stable`
versions available.

* **Staging:** refers to the latest tagged release from Github. This is documented in the [Next](/next/quickstart-ui) pages.

* **Development:** refers to the 'tip of HEAD' from Github. This is the ongoing development version and changes constantly.

<Tabs>
<TabItem value="stagingOperator" label="Staging version (x86-64, ARM64 (Raspberry Pi 4))" default>

```shell showLineNumbers
helm upgrade --create-namespace -n cattle-elemental-system --install elemental-operator-crds oci://registry.opensuse.org/isv/rancher/elemental/staging/charts/rancher/elemental-operator-crds-chart
helm upgrade --create-namespace -n cattle-elemental-system --install elemental-operator oci://registry.opensuse.org/isv/rancher/elemental/staging/charts/rancher/elemental-operator-chart
```

</TabItem>
<TabItem value="develOperator" label="Development version (x86-64, ARM64 (Raspberry Pi 4))" default>

:::warning Reminder
The development version is not recommended for production environments. We welcome feedback via Slack or Github issues, but it could be unstable and contain experimental features that can be dropped without notice.
:::

```shell showLineNumbers
helm upgrade --create-namespace -n cattle-elemental-system --install elemental-operator-crds oci://registry.opensuse.org/isv/rancher/elemental/dev/charts/rancher/elemental-operator-crds-chart
helm upgrade --create-namespace -n cattle-elemental-system --install --set image.imagePullPolicy=Always elemental-operator oci://registry.opensuse.org/isv/rancher/elemental/dev/charts/rancher/elemental-operator-chart
```

</TabItem>
</Tabs>

### Installation options

There are a few options that can be set in the chart install but that is out of scope for this document. You can see all the values on the chart [values.yaml](https://github.com/rancher/elemental-operator/blob/main/.obs/chartfile/operator/values.yaml).


---

## Article: partials/_quickstart-prereqs.md

## Prerequisites

* A Rancher server (v2.7.0 or later) configured (server-url set)
  * To configure the Rancher `server-url` please check the [Rancher docs](https://ranchermanager.docs.rancher.com/how-to-guides/new-user-guides/authentication-permissions-and-global-configuration#first-log-in)
* A machine (bare metal or virtualized) with TPM 2.0
  * Hint 1: Libvirt allows setting virtual TPMs for virtual machines [example here](tpm#add-tpm-module-to-virtual-machine)
  * Hint 2: You can enable TPM emulation on bare metal machines missing the TPM 2.0 module [example here](tpm#add-tpm-emulation-to-bare-metal-machine)
  * Hint 3: Make sure you're using UEFI (not BIOS) on x86-64, or the ISO won't boot
  * Hint 4: A minimum volume size of 25 GB is recommended. See the [Elemental partition table](installation#deployed-partition-table) for more details  
  * Hint 5: CPU and RAM requirements depend on the Kubernetes version installed, for example [K3s](https://docs.k3s.io/installation/requirements#hardware) or [RKE2](https://docs.rke2.io/install/requirements#hardware)  
* Helm Package Manager (https://helm.sh/)
* For ARM (aarch64) - One SD-card (32 GB or more, must be **fast** - 40MB/s write speed is acceptable) and a USB-stick for installation
