### AI Context

this is the official Longhorn documentation. Longhorn is a cloud native distributed block storage system for Kubernetes.

# Merged Knowledge Base for AI

## Article: _index.md

---
title: The Longhorn Documentation
description: Cloud native distributed block storage for Kubernetes
weight: 1
---

**Longhorn** is a lightweight, reliable, and powerful distributed [block storage](https://cloudacademy.com/blog/object-storage-block-storage/) system for Kubernetes.

Longhorn implements distributed block storage using containers and microservices. Longhorn creates a dedicated storage controller for each block device volume and synchronously replicates the volume across multiple replicas stored on multiple nodes. The storage controller and replicas are themselves orchestrated using Kubernetes.

## Features

* Enterprise-grade distributed block storage with no single point of failure
* Incremental snapshot of block storage
* Backup to secondary storage ([NFS](https://www.extrahop.com/resources/protocols/nfs/) or [S3](https://aws.amazon.com/s3/)-compatible object storage) built on efficient change block detection
* Recurring snapshots and backups
* Automated, non-disruptive upgrades. You can upgrade the entire Longhorn software stack without disrupting running storage volumes.
* An intuitive GUI dashboard

---

## Article: best-practices.md

---
title: Best Practices
weight: 4
---

We recommend the following setup for deploying Longhorn in production.

- [Minimum Recommended Hardware](#minimum-recommended-hardware)
- [Architecture](#architecture)
- [Operating System](#operating-system)
- [Kubernetes](#kubernetes)
  - [Kubernetes Version](#kubernetes-version)
  - [CoreDNS Setup](#coredns-setup)
- [Nodes and Disk Setup](#nodes-and-disk-setup)
  - [Use a Dedicated Disk](#use-a-dedicated-disk)
  - [Minimal Available Storage and Over-provisioning](#minimal-available-storage-and-over-provisioning)
  - [Disk Space Management](#disk-space-management)
  - [Setting up Extra Disks](#setting-up-extra-disks)
- [Configuring Default Disks Before and After Installation](#configuring-default-disks-before-and-after-installation)
- [Volumes Performance Optimization](#volumes-performance-optimization)
  - [IO Performance](#io-performance)
  - [Space Efficiency](#space-efficiency)
  - [Disaster Recovery](#disaster-recovery)
- [Deploying Workloads](#deploying-workloads)
- [Volumes Maintenance](#volumes-maintenance)
- [Guaranteed Instance Manager CPU](#guaranteed-instance-manager-cpu)
  - [V1 Data Engine](#v1-data-engine)
  - [V2 Data Engine](#v2-data-engine)
- [StorageClass](#storageclass)
- [Scheduling Settings](#scheduling-settings)
  - [Replica Node Level Soft Anti-Affinity](#replica-node-level-soft-anti-affinity)
  - [Allow Volumes Creation with Degraded Availability](#allow-volumes-creation-with-degraded-availability)
  - [Replica Auto-Balance](#replica-auto-balance)

## Minimum Recommended Hardware

- 3 nodes
- 4 vCPUs per node
- 4 GiB per node
- SSD/NVMe or similar performance block device on the node for storage (recommended)
- HDD/Spinning Disk or similar performance block device on the node for storage (verified)
  - 500/250 max IOPS per volume (1 MiB I/O)
  - 500/250 max throughput per volume (MiB/s)

> **Warning**: While Longhorn can function with HDDs (spinning disks) as storage, it is important to understand that **latency** plays a much more important role in volume stability than IOPS or throughput. This is because HDDs are mechanical, relying on spinning platters and moving read/write heads to access data. This physical movement introduces inherent delays (seek time and rotational delay), leading to much higher latency compared to the SSDs or NVMe drives, which utilize flash memory and have no moving parts. This can directly cause instability, especially when multiple input-output intensive tasks are running, such as:
>
> - Foreground IOs to the replicas
> - Foreground IOs from the replicas
> - Rebuilding volumes
> - Backups or other workloads
>
> The increased latency due to the use of HDDs, combined with other input-output workloads, can lead to **volume instability**. Therefore, we recommend **SSD or NVMe** drives for better performance and stability, especially for production workloads.
>
> The mentioned IOPS and throughput (500/250 max IOPS per volume and 500/250 max throughput per volume) are intended as general references based on the test setup but **should not be treated as hard requirements**. Latency, not just throughput, is the most important factor in ensuring system stability.

## Architecture

Longhorn supports the following architectures:

1. AMD64
1. ARM64

## Operating System

> **Note:** CentOS Linux has been removed from the verified OS list below, as it has been discontinued in favor of CentOS Stream [[ref](https://www.redhat.com/en/blog/faq-centos-stream-updates#Q5)], a rolling-release Linux distribution. Our focus for verifying RHEL-based downstream open source distributions will be enterprise-grade, such as Rocky and Oracle Linux.

The following Linux OS distributions and versions have been verified during the v{{< current-version >}} release testing. However, this does not imply that Longhorn exclusively supports these distributions. Essentially, Longhorn should function well on any certified Kubernetes cluster running on Linux nodes with a wide range of general-purpose operating systems, as well as verified container-optimized operating systems like SLE Micro.

| No. | OS                           | Versions
|-----|------------------------------| --------
| 1.  | Ubuntu                       | 24.04
| 2.  | SUSE Linux Enterprise Server | 15 SP7
| 3.  | SUSE Linux Enterprise Micro  | 6.1
| 4.  | Red Hat Enterprise Linux     | 10.0
| 5.  | Oracle Linux                 | 9.6
| 6.  | Rocky Linux                  | 10.0
| 7.  | Talos Linux                  | 1.10.6
| 8.  | Container-Optimized OS (GKE) | 121

Longhorn relies heavily on kernel functionality and performs better on some kernel versions. The following activities,
in particular, benefit from usage of specific kernel versions.

- Optimizing or improving the filesystem: Use a kernel with version `v5.8` or later. See [Issue
  #2507](https://github.com/longhorn/longhorn/issues/2507#issuecomment-857195496) for details.
- Enabling the [Freeze Filesystem for Snapshot](../references/settings#freeze-filesystem-for-snapshot) setting: Use a
  kernel with version `5.17` or later to ensure that a volume crash during a filesystem freeze cannot lock up a node.
- Enabling the [V2 Data Engine](../v2-data-engine/prerequisites): Use a kernel with version `5.19` or later to ensure


The list below contains known broken kernel versions that users should avoid using:

| No. | Version          | Distro          | Additional Context
|-----|------------------|-----------------| ------------------
| 1.  | 6.5.6            | Vanilla kernel  | Related to this bug https://longhorn.io/kb/troubleshooting-rwx-volume-fails-to-attached-caused-by-protocol-not-supported/
| 2.  | 5.15.0-94        | Ubuntu          | Related to this bug https://longhorn.io/kb/troubleshooting-rwx-volume-fails-to-attached-caused-by-protocol-not-supported/
| 3.  | 6.5.0-21         | Ubuntu          | Related to this bug https://longhorn.io/kb/troubleshooting-rwx-volume-fails-to-attached-caused-by-protocol-not-supported/
| 4.  | 6.5.0-1014-aws   | Ubuntu          | Related to this bug https://longhorn.io/kb/troubleshooting-rwx-volume-fails-to-attached-caused-by-protocol-not-supported/


## Kubernetes

### Kubernetes Version

Please ensure your Kubernetes cluster is at least v1.21 before upgrading to Longhorn v{{< current-version >}} because this is the minimum version Longhorn v{{< current-version >}} supports.

We recommend running your Kubernetes cluster on one of the following versions. These versions are the active supported versions prior to the Longhorn release, and have been tested with Longhorn v{{< current-version >}}.

| Release | Released     | End-of-life
|---------|--------------| -----------
| 1.33    | 23 Apr 2025  | 28 Jun 2026
| 1.32    | 11 Dec 2024  | 28 Feb 2026
| 1.31    | 13 Aug 2024  | 28 Oct 2025
| 1.30    | 17 Apr 2024  | 28 Jun 2025

Referenced to https://endoflife.date/kubernetes.

### CoreDNS Setup

Ensure that CoreDNS runs with at least 2 replicas to maintain high availability. This setup minimizes interruptions in the DNS resolution if one CoreDNS pod experiences a temporary disruption.

## Nodes and Disk Setup

We recommend the following setup for nodes and disks.

### Use a Dedicated Disk

It's recommended to dedicate a disk for Longhorn storage for production, instead of using the root disk.

### Minimal Available Storage and Over-provisioning

If you need to use the root disk, use the default `minimal available storage percentage` setup which is 25%, and set `overprovisioning percentage` to 100% to minimize the chance of DiskPressure.

If you're using a dedicated disk for Longhorn, you can lower the setting `minimal available storage percentage` to 10%.

For the Over-provisioning percentage, it depends on how much space your volume uses on average. For example, if your workload only uses half of the available volume size, you can set the Over-provisioning percentage to `200`, which means Longhorn will consider the disk to have twice the schedulable size as its full size minus the reserved space.

### Disk Space Management

Since Longhorn doesn't currently support sharding between the different disks, we recommend using [LVM](https://en.wikipedia.org/wiki/Logical_Volume_Manager_(Linux)) to aggregate all the disks for Longhorn into a single partition, so it can be easily extended in the future.

### Setting up Extra Disks

Any extra disks must be written in the `/etc/fstab` file to allow automatic mounting after the machine reboots.

Don't use a symbolic link for the extra disks. Use `mount --bind` instead of `ln -s` and make sure it's in the `fstab` file. For details, see [the section about multiple disk support.](../nodes-and-volumes/nodes/multidisk/#use-an-alternative-path-for-a-disk-on-the-node)

## Configuring Default Disks Before and After Installation

To use a directory other than the default `/var/lib/longhorn` for storage, the `Default Data Path` setting can be changed before installing the system. For details on changing pre-installation settings, refer to [this section.](../advanced-resources/deploy/customizing-default-settings)

The [Default node/disk configuration](../nodes-and-volumes/nodes/default-disk-and-node-config) feature can be used to customize the default disk after installation. Customizing the default configurations for disks and nodes is useful for scaling the cluster because it eliminates the need to configure Longhorn manually for each new node if the node contains more than one disk, or if the disk configuration is different for new nodes. Remember to enable `Create default disk only on labeled node` if applicable.

## Volumes Performance Optimization

Before configuring workloads, ensure that you have set up the following basic requirements for optimal volume performance.

- SATA/NVMe SSDs or disk drives with similar performance
- 10 Gbps network bandwidth between nodes
- Dedicated Priority Class for system-managed and user-deployed Longhorn components. By default, Longhorn installs the default Priority Class `longhorn-critical`.

The following sections outline other recommendations for production environments.

### IO Performance

- **Storage network**: Use a [dedicated storage network](../advanced-resources/deploy/storage-network/#setting-storage-network) to improve IO performance and stability.

- **Longhorn disk**: Use a [dedicated disk](../nodes-and-volumes/nodes/multidisk/#add-a-disk) for Longhorn storage instead of using the root disk.

- **Replica count**: Set the [default replica count](../references/settings/#default-replica-count) to "2" to achieve data availability with better disk space usage or less impact to system performance. This practice is especially beneficial to data-intensive applications.

- **Storage tag**: Use [storage tags](../nodes-and-volumes/nodes/storage-tags) to define storage tiering for data-intensive applications. For example, only high-performance disks can be used for storing performance-sensitive data.

- **Data locality**: Use `best-effort` as the default [data locality](../high-availability/data-locality) of Longhorn StorageClasses.

  For applications that support data replication (for example, a distributed database), you can use the `strict-local` option to ensure that only one replica is created for each volume. This practice prevents the extra disk space usage and IO performance overhead associated with volume replication.

  For data-intensive applications, you can use pod scheduling functions such as node selector or taint toleration. These functions allow you to schedule the workload to a specific storage-tagged node together with one replica.

### Space Efficiency

- **Recurring snapshots**: Periodically clean up system-generated snapshots and retain only the number of snapshots that makes sense for your implementation.

  For applications with replication capability, periodically [delete all types of snapshots](../concepts/#243-deleting-snapshots).

- **Recurring filesystem trim**: Periodically [trim the filesystem](../nodes-and-volumes/volumes/trim-filesystem) inside volumes to reclaim disk space.

- **Snapshot space management**: [Configure global and volume-specific settings](../snapshots-and-backups/snapshot-space-management) to prevent unexpected disk space exhaustion.

### Disaster Recovery

- **Recurring backups**: Create [recurring backup jobs](../snapshots-and-backups/scheduling-backups-and-snapshots/) for mission-critical application volumes.

- **System backup**: Create periodic [system backups](../advanced-resources/system-backup-restore/backup-longhorn-system/#create-longhorn-system-backup).

## Deploying Workloads

If you're using `ext4` as the filesystem of the volume, we recommend adding a liveness check to workloads to help automatically recover from a network-caused interruption, a node reboot, or a Docker restart. See [this section](../high-availability/recover-volume) for details.

## Volumes Maintenance

Using Longhorn's built-in backup feature is highly recommended. You can save backups to an object store such as S3 or to an NFS server. Saving to an object store is preferable because it generally offers better reliability.  Another advantage is that you do not need to mount and unmount the target, which can complicate failover and upgrades.

For each volume, schedule at least one recurring backup. If you must run Longhorn in production without a backupstore, then schedule at least one recurring snapshot for each volume.

Longhorn system will create snapshots automatically when rebuilding a replica. Recurring snapshots or backups can also automatically clean up the system-generated snapshot.

## Guaranteed Instance Manager CPU

We recommend setting the CPU request for Longhorn instance manager pods.

### V1 Data Engine

The `Guaranteed Instance Manager CPU` setting allows you to reserve a percentage of the total allocatable CPU resources on each node for each instance manager pod when the V1 Data Engine is enabled. The default value is 12.

You can also set a specific milli CPU value for instance manager pods on a particular node by updating the node's `Instance Manager CPU Request` field.

> **Note:** This field will overwrite the above setting for the specified node.

Refer to [Guaranteed Instance Manager CPU](../references/settings/#guaranteed-instance-manager-cpu) for more details.

### V2 Data Engine

The `Guaranteed Instance Manager CPU for V2 Data Engine` setting allows you to reserve a specific number of millicpus on each node for each instance manager pod when the V2 Data Engine is enabled. By default, the Storage Performance Development Kit (SPDK) target daemon within each instance manager pod uses 1 CPU core. Configuring a minimum CPU usage value is essential for maintaining engine and replica stability, especially during periods of high node workload. The default value is 1250.

## StorageClass

We don't recommend modifying the default StorageClass named `longhorn`, since the change of parameters might cause issues during an upgrade later. If you want to change the parameters set in the StorageClass, you can create a new StorageClass by referring to the [StorageClass examples](../references/examples/#storageclass).

## Scheduling Settings

### Replica Node Level Soft Anti-Affinity

> Recommend: `false`

This setting should be set to `false` in production environment to ensure the best availability of the volume. Otherwise, one node down event may bring down more than one replicas of a volume.

### Allow Volumes Creation with Degraded Availability

> Recommend: `false`

This setting should be set to `false` in production environment to ensure every volume have the best availability when created. Because with the setting set to `true`, the volume creation won't error out even there is only enough room to schedule one replica. So there is a risk that the cluster is running out of the spaces but the user won't be made aware immediately.

### Replica Auto-Balance

> Recommend: `least-effort`

For production environments, we recommend setting Replica Auto-Balance to `least-effort`. This setting ensures that at least one replica is placed on a different node in each zone, providing extra high availability (HA).

In certain edge cases, you might consider using the `best-effort`, which continuously attempts to evenly distribute replicas across nodes and zones. However, this setting can lead to frequent rebuilds if the cluster is unstable.

For most users, having multiple replicas without Replica Auto-Balance setting is sufficient to achieve basic HA, especially if you prefer to avoid excessive rebuilds and resource usage.


---

## Article: concepts.md

---
title: Architecture and Concepts
weight: 2
---

Longhorn creates a dedicated storage controller for each volume and synchronously replicates the volume across multiple replicas stored on multiple nodes.

The storage controller and replicas are themselves orchestrated using Kubernetes.

For an overview of Longhorn features, refer to [this section.](../what-is-longhorn)

For the installation requirements, go to [this section.](../deploy/install/#installation-requirements)

> This section assumes familiarity with Kubernetes persistent storage concepts. For more information on these concepts, refer to the [appendix.](#appendix-how-persistent-storage-works-in-kubernetes) For help with the terminology used in this page, refer to [this section.](../terminology)

- [1. Design](#1-design)
  - [1.1. The Longhorn Manager and the Longhorn Engine](#11-the-longhorn-manager-and-the-longhorn-engine)
  - [1.2. Advantages of a Microservices Based Design](#12-advantages-of-a-microservices-based-design)
  - [1.3. CSI Driver](#13-csi-driver)
  - [1.4. CSI Plugin](#14-csi-plugin)
  - [1.5. The Longhorn UI](#15-the-longhorn-ui)
- [2. Longhorn Volumes and Primary Storage](#2-longhorn-volumes-and-primary-storage)
    - [2.1. Thin Provisioning and Volume Size](#21-thin-provisioning-and-volume-size)
    - [2.2. Reverting Volumes in Maintenance Mode](#22-reverting-volumes-in-maintenance-mode)
  - [2.3. Replicas](#23-replicas)
    - [2.3.1. How Read and Write Operations Work for Replicas](#231-how-read-and-write-operations-work-for-replicas)
    - [2.3.2 How New Replicas are Added](#232-how-new-replicas-are-added)
    - [2.3.3. How Faulty Replicas are Rebuilt](#233-how-faulty-replicas-are-rebuilt)
  - [2.4. Snapshots](#24-snapshots)
    - [2.4.1. How Snapshots Work](#241-how-snapshots-work)
    - [2.4.2. Recurring Snapshots](#242-recurring-snapshots)
    - [2.4.3. Deleting Snapshots](#243-deleting-snapshots)
    - [2.4.4. Storing Snapshots](#244-storing-snapshots)
    - [2.4.5. Crash Consistency](#245-crash-consistency)
- [3. Backups and Secondary Storage](#3-backups-and-secondary-storage)
  - [3.1. How Backups Work](#31-how-backups-work)
  - [3.2. Recurring Backups](#32-recurring-backups)
  - [3.3. Disaster Recovery Volumes](#33-disaster-recovery-volumes)
  - [3.4. Backupstore Update Intervals, RTO, and RPO](#34-backupstore-update-intervals-rto-and-rpo)
- [Appendix: How Persistent Storage Works in Kubernetes](#appendix-how-persistent-storage-works-in-kubernetes)
  - [How Kubernetes Workloads use New and Existing Persistent Storage](#how-kubernetes-workloads-use-new-and-existing-persistent-storage)
    - [Existing Storage Provisioning](#existing-storage-provisioning)
    - [Dynamic Storage Provisioning](#dynamic-storage-provisioning)
  - [Horizontal Scaling for Kubernetes Workloads with Persistent Storage](#horizontal-scaling-for-kubernetes-workloads-with-persistent-storage)

# 1. Design

The Longhorn design has two layers: the data plane and the control plane. The Longhorn Engine is a storage controller that corresponds to the data plane, and the Longhorn Manager corresponds to the control plane.

## 1.1. The Longhorn Manager and the Longhorn Engine

The Longhorn Manager Pod runs on each node in the Longhorn cluster as a Kubernetes [DaemonSet](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/). It is responsible for creating and managing volumes in the Kubernetes cluster, and handles the API calls from the Longhorn UI or the Longhorn CSI plugin. It follows the Kubernetes controller pattern, which is sometimes called the operator pattern.

The Longhorn Manager communicates with the Kubernetes API server to create a new Longhorn volume [CR](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/). Then the Longhorn Manager watches the API server's response, and when it sees that the Kubernetes API server created a new Longhorn volume CR, the Longhorn Manager creates a new volume.

When the Longhorn Manager is asked to create a volume, it creates a Longhorn Engine instance on the node the volume is attached to, and it creates a replica on each node where a replica will be placed. Replicas should be placed on separate hosts to ensure maximum availability.

The multiple data paths of the replicas ensure high availability of the Longhorn volume. Even if a problem happens with a certain replica or with the Engine, the problem won't affect all the replicas or the Pod's access to the volume. The Pod will still function normally. For a given volume, if the replica count is `N`, the Longhorn volume can tolerate a maximum of `N−1` replica failures. This is because at least one healthy replica is required for the volume to remain operational.

The Longhorn Engine always runs in the same node as the Pod that uses the Longhorn volume. It synchronously replicates the volume across the multiple replicas stored on multiple nodes.

The Engine and replicas are orchestrated using Kubernetes.

In the figure below,

- There are three instances with Longhorn volumes.
- Each volume has a dedicated controller, called the Longhorn Engine. For V1 volumes, the engine runs as a Linux process, while for V2 volumes, it operates as an SPDK RAID block device (bdev).
- Each Longhorn volume has two replicas. In V1, replicas run as Linux processes, while in V2, they are implemented as SPDK logical volume bdevs.
- The arrows in the figure indicate the read/write data flow between the volume, controller instance, replica instances, and disks.
- By creating a separate Longhorn Engine for each volume, if one controller fails, the function of other volumes is not impacted.

**Figure 1. Read/write Data Flow between the Volume, Longhorn Engine, Replica Instances, and Disks**

{{< figure alt="read/write data flow between the volume, controller instance, replica instances, and disks" src="/img/diagrams/architecture/how-longhorn-works-with-kubernetes.svg" >}}

## 1.2. Advantages of a Microservices Based Design

In Longhorn, each Engine only needs to serve one volume, simplifying the design of the storage controllers. Because the failure domain of the controller software is isolated to individual volumes, a controller crash will only impact one volume.

The Longhorn Engine is simple and lightweight, allowing us to create thousands of separate engines. Kubernetes schedules these separate engines, drawing resources from a shared set of disks and working with Longhorn to form a resilient distributed block storage system.

Because each volume has its own controller, the controller and replica instances for each volume can also be upgraded without causing a noticeable disruption in IO operations.

Longhorn can create a long-running job to orchestrate the upgrade of all live volumes without disrupting the on-going operation of the system. To ensure that an upgrade does not cause unforeseen issues, Longhorn can choose to upgrade a small subset of the volumes and roll back to the old version if something goes wrong during the upgrade.

## 1.3. CSI Driver

The Longhorn CSI driver takes the block device, formats it, and mounts it on the node. Then the [kubelet](https://kubernetes.io/docs/reference/command-line-tools-reference/kubelet/) bind-mounts the device inside a Kubernetes Pod. This allows the Pod to access the Longhorn volume.

The required Kubernetes CSI Driver images will be deployed automatically by the longhorn driver deployer.
To install Longhorn in an air gapped environment, refer to [this section](../deploy/install/airgap).

## 1.4. CSI Plugin

Longhorn is managed in Kubernetes via a [CSI Plugin.](https://kubernetes-csi.github.io/docs/) This allows for easy installation of the Longhorn plugin.

The Kubernetes CSI plugin calls Longhorn to create volumes to create persistent data for a Kubernetes workload. The CSI plugin gives you the ability to create, delete, attach, detach, mount the volume, and take snapshots of the volume. All other functionality provided by Longhorn is implemented through the Longhorn UI.

The Kubernetes cluster internally uses the CSI interface to communicate with the Longhorn CSI plugin. And the Longhorn CSI plugin communicates with the Longhorn Manager using the Longhorn API.

For v1 volumes, Longhorn uses iSCSI, which may require additional configuration on your nodes:
- Depending on the Linux distribution, you need to install either `open-iscsi` or `iscsiadm`.

In contrast, v2 volumes come with different prerequisites, depending on the configuration:
- Kernel modules like `vfio_pci` and `uio_pci_generic` are required.
- For the NVMe-TCP frontend, the `nvme_tcp` module is necessary.
- For the UBLK frontend, both the `ublk_drv` module and huge pages support must be enabled.

## 1.5. The Longhorn UI

The Longhorn UI interacts with the Longhorn Manager through the Longhorn API, and acts as a complement of Kubernetes. Through the Longhorn UI, you can manage snapshots, backups, nodes and disks.

Besides, the space usage of the cluster worker nodes is collected and illustrated by the Longhorn UI. See [here](../nodes-and-volumes/nodes/node-space-usage) for details.

# 2. Longhorn Volumes and Primary Storage

When creating a volume, the Longhorn Manager creates the Longhorn Engine microservice and the replicas for each volume as microservices. Together, these microservices form a Longhorn volume. Each replica should be placed on a different node or on different disks.

After the Longhorn Engine is created by the Longhorn Manager, it connects to the replicas. The Engine exposes a block device on the same node where the Pod is running.

A Longhorn volume can be created with kubectl.

### 2.1. Thin Provisioning and Volume Size

Longhorn is a thin-provisioned storage system. That means a Longhorn volume will only take the space it needs at the moment. For example, if you allocated a 20 GB volume but only use 1GB of it, the actual data size on your disk would be 1 GB. You can see the actual data size in the volume details in the UI.

A Longhorn volume itself cannot shrink in size if you’ve removed content from your volume. For example, if you create a volume of 20 GB, used 10 GB, then removed the content of 9 GB, the actual size on the disk would still be 10 GB instead of 1 GB. This happens because Longhorn operates on the block level, not the filesystem level, so Longhorn doesn’t know if the content has been removed by a user or not. That information is mostly kept at the filesystem level.

For more introductions about the volume-size related concepts, see this [doc](../nodes-and-volumes/volumes/volume-size) for more details.

### 2.2. Reverting Volumes in Maintenance Mode

When a volume is attached from the Longhorn UI, there is a checkbox for Maintenance mode. It’s mainly used to revert a volume from a snapshot.

The option will result in attaching the volume without enabling the frontend (block device or iSCSI), to make sure no one can access the volume data when the volume is attached.

After v0.6.0, the snapshot reverting operation required the volume to be in maintenance mode. This is because if the block device’s content is modified while the volume is mounted or being used, it will cause filesystem corruption.

It’s also useful to inspect the volume state without worrying about the data being accessed by accident.

## 2.3. Replicas

Each replica contains a chain of snapshots of a Longhorn volume. A snapshot is like a layer of an image, with the oldest snapshot used as the base layer, and newer snapshots on top. Data is only included in a new snapshot if it overwrites data in an older snapshot. Together, a chain of snapshots shows the current state of the data.

For each Longhorn volume, multiple replicas of the volume should run in the Kubernetes cluster, each on a separate node. All replicas are treated the same, and the Longhorn Engine always runs on the same node as the pod, which is also the consumer of the volume. In that way, we make sure that even if the Pod is down, the Engine can be moved to another Pod and your service will continue undisrupted.

The default replica count can be changed in the [settings.](../references/settings/#default-replica-count) When a volume is attached, the replica count for the volume can be changed in the UI.

If the current healthy replica count is less than specified replica count, Longhorn will start rebuilding new replicas.

If the current healthy replica count is more than the specified replica count, Replica Auto Balance and Data Locality are disabled, Longhorn will do nothing. In this situation, if a replica fails or is deleted, Longhorn won’t start rebuilding new replicas unless the healthy replica count dips below the specified replica count. If Replica Auto Balance or Data Locality are set, Longhorn might delete one of the replicas.

Longhorn replicas are built using Linux [sparse files,](https://en.wikipedia.org/wiki/Sparse_file) which support thin provisioning.

###  2.3.1. How Read and Write Operations Work for Replicas

When data is read from a replica of a volume, if the data can be found in the live data, then that data is used. If not, the newest snapshot will be read. If the data is not found in the newest snapshot, the next-oldest snapshot is read, and so on, until the oldest snapshot is read.

When you take a snapshot, a [differencing](https://en.wikipedia.org/wiki/Data_differencing) disk is created. As the number of snapshots grows, the differencing disk chain (also called a chain of snapshots) could get quite long. To improve read performance, Longhorn therefore maintains a read index that records which differencing disk holds valid data for each 4K block of storage.

In the following figure, the volume has eight blocks. The read index has eight entries and is filled up lazily as read operations take place.

A write operation resets the read index, causing it to point to the live data. The live data consists of data at some indices and empty space in other indices.

Beyond the read index, we currently do not maintain additional metadata to indicate which blocks are used.

**Figure 2. How the Read Index Keeps Track of Which Snapshot Holds the Most Recent Data**

{{< figure alt="how the read index keeps track of which snapshot holds the most recent data" src="/img/diagrams/architecture/read-index.png" >}}

The figure above is color-coded to show which blocks contain the most recent data according to the read index, and the source of the latest data is also listed in the table below:

| Read Index | Source of the latest data |
|---------------|--------------------------------|
| 0                       | Newest snapshot                 |
| 1                       | Live data                                    |
| 2                       | Oldest snapshot                   |
| 3                       | Oldest snapshot                   |
| 4                       | Oldest snapshot                   |
| 5                       | Live data                                   |
| 6                       | Live data                                   |
| 7                       | Live data                                   |

Note that as the green arrow shows in the figure above, Index 5 of the read index previously pointed to the second-oldest snapshot as the source of the most recent data, then it changed to point to the the live data when the 4K block of storage at Index 5 was overwritten by the live data.

The read index is kept in memory and consumes one byte for each 4K block. The byte-sized read index means you can take as many as 254 snapshots for each volume.

The read index consumes a certain amount of in-memory data structure for each replica. A 1 TB volume, for example, consumes 256 MB of in-memory read index.

### 2.3.2 How New Replicas are Added

When a new replica is added, the existing replicas are synced to the new replica. The first replica is created by taking a new snapshot from the live data.

The following steps show a more detailed breakdown of how Longhorn adds new replicas:

1. The Longhorn Engine is paused.
1. Let's say that the chain of snapshots within the replica consists of the live data and a snapshot. When the new replica is created, the live data becomes the newest (second) snapshot and a new, blank version of live data is created.
1. The new replica is created in WO (write-only) mode.
1. The Longhorn Engine is unpaused.
1. All the snapshots are synced.
1. The new replica is set to RW (read-write) mode.

### 2.3.3. How Faulty Replicas are Rebuilt

Longhorn will always try to maintain at least given number of healthy replicas for each volume.

When the controller detects failures in one of its replicas, it marks the replica as being in an error state. The Longhorn Manager is responsible for initiating and coordinating the process of rebuilding the faulty replica.

To rebuild the faulty replica, the Longhorn Manager creates a blank replica and calls the Longhorn Engine to add the blank replica into the volume's replica set.

To add the blank replica, the Engine performs the following operations:
  1. Pauses all read and write operations.
  2. Adds the blank replica in WO (write-only) mode.
  3. Takes a snapshot of all existing replicas, which will now have a blank differencing disk at its head.
  4. Unpauses all read and write operations. Only write operations will be dispatched to the newly added replica.
  5. Starts a background process to sync all but the most recent differencing disk from a good replica to the blank replica.
  6. After the sync completes, all replicas now have consistent data, and the volume manager sets the new replica to RW (read-write) mode.

Finally, the Longhorn Manager calls the Longhorn Engine to remove the faulty replica from its replica set.

## 2.4. Snapshots

The snapshot feature enables a volume to be reverted back to a certain point in history. Backups in secondary storage can also be built from a snapshot.

When a volume is restored from a snapshot, it reflects the state of the volume at the time the snapshot was created.

The snapshot feature is also a part of Longhorn's rebuilding process. Every time Longhorn detects a replica is down, it will automatically take a (system) snapshot and start rebuilding it on another node.

### 2.4.1. How Snapshots Work

A snapshot is like a layer of an image, with the oldest snapshot used as the base layer, and newer snapshots on top. Data is only included in a new snapshot if it overwrites data in an older snapshot. Together, a chain of snapshots shows the current state of the data. For a more detailed breakdown of how data is read from a replica, refer to the section on [read and write operations for replicas.](#231-how-read-and-write-operations-work-for-replicas)

Snapshots cannot change after they are created, unless a snapshot is deleted, in which case its changes are conflated with the next most recent snapshot. New data is always written to the live version. New snapshots are always created from live data.

To create a new snapshot, the live data becomes the newest snapshot. Then a new, blank version of the live data is created, taking the place of the old live data.

### 2.4.2. Recurring Snapshots

To reduce the space taken by snapshots, user can schedule a recurring snapshot or backup with a number of snapshots to retain, which will automatically create a new snapshot/backup on schedule, then clean up for any excessive snapshots/backups.

### 2.4.3. Deleting Snapshots

Unwanted snapshots can be manually deleted through the UI. Any system generated snapshots will be automatically marked for deletion if the deletion of any snapshot was triggered.

In Longhorn, the latest snapshot cannot be deleted. This is because whenever a snapshot is deleted, Longhorn will conflate its content with the next snapshot, so that the next and later snapshot retains the correct content.

But Longhorn cannot do that for the latest snapshot since there is no more recent snapshot to be conflated with the deleted snapshot. The next “snapshot” of the latest snapshot is the live volume (volume-head), which is being read/written by the user at the moment, so the conflation process cannot happen.

Instead, the latest snapshot will be marked as removed, and it will be cleaned up next time when possible.

To clean up the latest snapshot, a new snapshot can be created, then the previous "latest" snapshot can be removed.

### 2.4.4. Storing Snapshots

Snapshots are stored locally, as a part of each replica of a volume. They are stored on the disk of the nodes within the Kubernetes cluster.
Snapshots are stored in the same location as the volume data on the host’s physical disk.

### 2.4.5. Crash Consistency

Longhorn is a crash-consistent block storage solution.

It’s normal for the OS to keep content in the cache before writing into the block layer. This means that if all of the replicas are down, then Longhorn may not contain the changes that occurred immediately before the shutdown, because the content was kept in the OS-level cache and wasn't yet transferred to the Longhorn system.

This problem is similar to problems that could happen if your desktop computer shuts down due to a power outage. After resuming the power, you may find some corrupted files in the hard drive.

To force the data to be written to the block layer at any given moment, the sync command can be manually run on the node, or the disk can be unmounted. The OS would write the content from the cache to the block layer in either situation.

Longhorn runs the sync command automatically before creating a snapshot.

# 3. Backups and Secondary Storage

A backup is an object in the backupstore, which is an NFS or S3 compatible object store external to the Kubernetes cluster. Backups provide a form of secondary storage so that even if your Kubernetes cluster becomes unavailable, your data can still be retrieved.

Because the volume replication is synchronized, and because of network latency, it is hard to do cross-region replication. The backupstore is also used as a medium to address this problem.

When the backup target is configured on the Longhorn UI (**Backup and Restore > Backup Targets**), Longhorn can connect to the backupstore and display a list of existing backups on the **Backup** page.

If Longhorn runs in a second Kubernetes cluster, it can also sync disaster recovery volumes to the backups in secondary storage, so that your data can be recovered more quickly in the second Kubernetes cluster.

## 3.1. How Backups Work

A backup is created using one snapshot as a source, so that it reflects the state of the volume's data at the time that the snapshot was created. A backup is stored remotely outside of the cluster.

By contrast to a snapshot, a backup can be thought of as a flattened version of a chain of snapshots. Similar to the way that information is lost when a layered image is converted to a flat image, data is also lost when a chain of snapshots is converted to a backup. In both conversions, any overwritten data would be lost.

Because backups don't contain snapshots, they don't contain the history of changes to the volume data. After you restore a volume from a backup, the volume initially contains one snapshot. This snapshot is a conflated version of all the snapshots in the original chain, and it reflects the live data of the volume at the time at the time the backup was created.

While snapshots can be hundreds of gigabytes, backups are made of 2 MB files.

Each new backup of the same original volume is incremental, detecting and transmitting the changed blocks between snapshots. This is a relatively easy task because each snapshot is a [differencing](https://en.wikipedia.org/wiki/Data_differencing) file and only stores the changes from the last snapshot. This design also means that if no blocks have changed and a backup is taken, that backup in the backupstore will show as 0 bytes. However if you were to restore from that backup it would still contain the full volume data, since it would restore the necessary blocks already present on the backupstore, that are required for a backup.

To avoid storing a very large number of small blocks of storage, Longhorn performs backup operations using 2 MB blocks. That means that, if any 4K block in a 2MB boundary is changed, Longhorn will back up the entire 2MB block. This offers the right balance between manageability and efficiency.

**Figure 3. The Relationship between Backups in Secondary Storage and Snapshots in Primary Storage**

{{< figure alt="the relationship between backups in secondary storage and snapshots in primary storage" src="/img/diagrams/concepts/longhorn-backup-creation.png" >}}

The above figure describes how backups are created from snapshots in Longhorn:

- The Primary Storage side of the diagram shows one replica of a Longhorn volume in the Kubernetes cluster. The replica consists of a chain of four snapshots. In order from newest to oldest, the snapshots are Live Data, snap3, snap2, and snap1.
- The Secondary Storage side of the diagram shows two backups in an external object storage service such as S3.
- In Secondary Storage, the color coding for backup-from-snap2 shows that it includes both the blue change from snap1 and the green changes from snap2. No changes from snap2 overwrote the data in snap1, therefore the changes from both snap1 and snap2 are both included in backup-from-snap2.
- The backup named backup-from-snap3 reflects the state of the volume's data at the time that snap3 was created. The color coding and arrows indicate that backup-from-snap3 contains all of the dark red changes from snap3, but only one of the green changes from snap2. This is because one of the red changes in snap3 overwrote one of the green changes in snap2. This illustrates how backups don't include the full history of change, because they conflate snapshots with the snapshots that came before them.
- Each backup maintains its own set of 2 MB blocks. Each 2 MB block is backed up only once. The two backups share one green block and one blue block.

When a backup is deleted from the secondary storage, Longhorn does not delete all the blocks that it uses. Instead, it performs a garbage collection periodically to clean up unused blocks from secondary storage.

The 2 MB blocks for all backups belonging to the same volume are stored under a common directory and can therefore be shared across multiple backups.

To save space, the 2 MB blocks that didn't change between backups can be reused for multiple backups that share the same backup volume in secondary storage. Because checksums are used to address the 2 MB blocks, we achieve some degree of deduplication for the 2 MB blocks in the same volume.

Volume-level metadata is stored in volume.cfg. The metadata files for each backup (e.g., snap2.cfg) are relatively small because they only contain the [offsets](https://en.wikipedia.org/wiki/Offset_(computer_science)) and [checksums](https://en.wikipedia.org/wiki/Checksum) of all the 2 MB blocks in the backup.

Each 2 MB block (.blk file) is compressed.

## 3.2. Recurring Backups

Backup operations can be scheduled using the recurring snapshot and backup feature, but they can also be done as needed.

It’s recommended to schedule recurring backups for your volumes. If a backupstore is not available, it’s recommended to have the recurring snapshot scheduled instead.

Backup creation involves copying the data through the network, so it will take time.

## 3.3. Disaster Recovery Volumes

A disaster recovery (DR) volume is a special volume that stores data in a backup cluster in case the whole main cluster goes down. DR volumes are used to increase the resiliency of Longhorn volumes.

Because the main purpose of a DR volume is to restore data from backup, this type of volume doesn’t support the following actions before it is activated:

- Creating, deleting, and reverting snapshots
- Creating backups
- Creating persistent volumes
- Creating persistent volume claims

A DR volume can be created from a volume’s backup in the backupstore. After the DR volume is created, Longhorn will monitor its original backup volume and incrementally restore from the latest backup. A backup volume is an object in the backupstore that contains multiple backups of the same volume.

If the original volume in the main cluster goes down, the DR volume can be immediately activated in the backup cluster, reducing the time needed to restore the data from the backupstore to the volume in the backup cluster.

When a DR volume is activated, Longhorn will check the last backup of the original volume. If that backup has not already been restored, the restoration will be started, and the activate action will fail. Users need to wait for the restoration to complete before retrying.

The Backup Target in the Longhorn settings cannot be updated if any DR volumes exist.

After a DR volume is activated, it becomes a normal Longhorn volume and it cannot be deactivated.

## 3.4. Backupstore Update Intervals, RTO, and RPO

Incremental restoration is usually triggered by the periodic backupstore update. You can set the update interval on the backup target settings page (**Backup and Restore > Backup Targets**).

Notice that this interval can potentially impact Recovery Time Objective (RTO). If it is too long, there may be a large amount of data for the disaster recovery volume to restore, which will take a long time.

As for Recovery Point Objective (RPO), it is determined by recurring backup scheduling of the backup volume. If recurring backup scheduling for normal volume A creates a backup every hour, then the RPO is one hour. You can check here to see how to set recurring backups in Longhorn.

The following analysis assumes that the volume creates a backup every hour, and that incrementally restoring data from one backup takes five minutes:

- If the backupstore Poll Interval is 30 minutes, then there will be at most one backup worth of data since the last restoration. The time for restoring one backup is five minutes, so the RTO would be five minutes.
- If the backupstore Poll Interval is 12 hours, then there will be at most 12 backups worth of data since last restoration. The time for restoring the backups is 5 * 12 = 60 minutes, so the RTO would be 60 minutes.

# Appendix: How Persistent Storage Works in Kubernetes

To understand persistent storage in Kubernetes, it is important to understand Volumes, PersistentVolumes, PersistentVolumeClaims, and StorageClasses, and how they work together.

One important property of a Kubernetes Volume is that it has the same lifecycle as the Pod it belongs to. The Volume is lost if the Pod is gone. In contrast, a PersistentVolume continues to exist in the system until users delete it. Volumes can also be used to share data between containers inside the same Pod, but this isn’t the primary use case because users normally only have one container per Pod.

A [PersistentVolume (PV)](https://kubernetes.io/docs/concepts/storage/persistent-volumes/) is a piece of persistent storage in the Kubernetes cluster, while a [PersistentVolumeClaim (PVC)](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#persistentvolumeclaims) is a request for storage. [StorageClasses](https://kubernetes.io/docs/concepts/storage/storage-classes/) allow new storage to be dynamically provisioned for workloads on demand.

## How Kubernetes Workloads use New and Existing Persistent Storage

Broadly speaking, there are two main ways to use persistent storage in Kubernetes:

- Use an existing persistent volume
- Dynamically provision new persistent volumes

### Existing Storage Provisioning

To use an existing PV, your application will need to use a PVC that is bound to a PV, and the PV should include the minimum resources that the PVC requires.

In other words, a typical workflow for setting up existing storage in Kubernetes is as follows:

1. Set up persistent storage volumes, in the sense of physical or virtual storage that you have access to.
1. Add a PV that refers to the persistent storage.
1. Add a PVC that refers to the PV.
1. Mount the PVC as a volume in your workload.

When a PVC requests a piece of storage, the Kubernetes API server will try to match that PVC with a pre-allocated PV as matching volumes become available. If a match can be found, the PVC will be bound to the PV, and the user will start to use that pre-allocated piece of storage.

if a matching volume does not exist, PersistentVolumeClaims will remain unbound indefinitely. For example, a cluster provisioned with many 50 Gi PVs would not match a PVC requesting 100 Gi. The PVC could be bound after a 100 Gi PV is added to the cluster.

In other words, you can create unlimited PVCs, but they will only be bound to PVs if the Kubernetes master can find a sufficient PV that has at least the amount of disk space required by the PVC.

### Dynamic Storage Provisioning

For dynamic storage provisioning, your application will need to use a PVC that is bound to a StorageClass. The StorageClass contains the authorization to provision new persistent volumes.

The overall workflow for dynamically provisioning new storage in Kubernetes involves a StorageClass resource:

1. Add a StorageClass and configure it to automatically provision new storage from the storage that you have access to.
1. Add a PVC that refers to the StorageClass.
1. Mount the PVC as a volume for your workload.

Kubernetes cluster administrators can use a Kubernetes StorageClass to describe the “classes” of storage they offer. StorageClasses can have different capacity limits, different IOPS, or any other parameters that the provisioner supports. The storage vendor specific provisioner is be used along with the StorageClass to allocate PV automatically, following the parameters set in the StorageClass object. Also, the provisioner now has the ability to enforce the resource quotas and permission requirements for users. In this design, admins are freed from the unnecessary work of predicting the need for PVs and allocating them.

When a StorageClass is used, a Kubernetes administrator is not responsible for allocating every piece of storage. The administrator just needs to give users permission to access a certain storage pool, and decide the quota for the user. Then the user can carve out the needed pieces of the storage from the storage pool.

StorageClasses can also be used without explicitly creating a StorageClass object in Kubernetes. Since the StorageClass is also a field used to match a PVC with a PV, a PV can be created manually with a custom Storage Class name, then a PVC can be created that asks for a PV with that StorageClass name. Kubernetes can then bind your PVC to the PV with the specified StorageClass name, even if the StorageClass object doesn't exist as a Kubernetes resource.

Longhorn introduces a Longhorn StorageClass so that your Kubernetes workloads can carve out pieces of your persistent storage as necessary.

## Horizontal Scaling for Kubernetes Workloads with Persistent Storage

The VolumeClaimTemplate is a StatefulSet spec property, and it provides a way for the block storage solution to scale horizontally for a Kubernetes workload.

This property can be used to create matching PVs and PVCs for Pods that were created by a StatefulSet.

Those PVCs are created using a StorageClass, so they can be set up automatically when the StatefulSet scales up.

When a StatefulSet scales down, the extra PVs/PVCs are kept in the cluster, and they are reused when the StatefulSet scales up again.

The VolumeClaimTemplate is important for block storage solutions like EBS and Longhorn. Because those solutions are inherently [ReadWriteOnce,](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#access-modes) they cannot be shared between the Pods.

Deployments don't work well with persistent storage if you have more than one Pod running with persistent data. For more than one pod, a StatefulSet should be used.


---

## Article: contributing.md

---
title: Contributing
weight: 6
---

Longhorn is open source software, so contributions are greatly welcome. Please read the [Cloud Native Computing Foundation Code of Conduct](https://github.com/cncf/foundation/blob/master/code-of-conduct.md) and [Contributing Guidelines](https://github.com/longhorn/longhorn/blob/master/CONTRIBUTING.md) before contributing.

Contributing code is not the only way of contributing. We value feedback very much and many of the Longhorn features are originated from users' feedback. If you have any feedback, feel free to [file an issue](https://github.com/longhorn/longhorn/issues/new/choose) and talk to the developers at the [CNCF](https://slack.cncf.io/) [#longhorn](https://cloud-native.slack.com/messages/longhorn) slack channel.

Longhorn is a [CNCF Incubating Project.](https://www.cncf.io/projects/longhorn/)

![Longhorn is a CNCF Incubating Project](https://raw.githubusercontent.com/cncf/artwork/master/other/cncf/horizontal/color/cncf-color.svg)

## Source Code

Longhorn is 100% open source software under the auspices of the [Cloud Native Computing Foundation](https://cncf.io). The project's source code is spread across a number of repos:

| Component                      | What it does                                                           | GitHub repo                                                                                 |
| :----------------------------- | :--------------------------------------------------------------------- | :------------------------------------------------------------------------------------------ |
| Longhorn Backing Image Manager | Backing image download, sync, and deletion in a disk                   | [longhorn/backing-image-manager](https://github.com/longhorn/backing-image-manager)         |
| Longhorn Engine                | Core controller/replica logic                                          | [longhorn/longhorn-engine](https://github.com/longhorn/longhorn-engine)                     |
| Longhorn Instance Manager      | Controller/replica instance lifecycle management                       | [longhorn/longhorn-instance-manager](https://github.com/longhorn/longhorn-instance-manager) |
| Longhorn Manager               | Longhorn orchestration, includes CSI driver for Kubernetes             | [longhorn/longhorn-manager](https://github.com/longhorn/longhorn-manager)                   |
| Longhorn Share Manager         | NFS provisioner that exposes Longhorn volumes as ReadWriteMany volumes | [longhorn/longhorn-share-manager](https://github.com/longhorn/longhorn-share-manager)       |
| Longhorn UI                    | The Longhorn dashboard                                                 | [longhorn/longhorn-ui](https://github.com/longhorn/longhorn-ui)                             |

## License

Copyright (c) 2014-2021 The Longhorn Authors.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at [Apache License 2.0](http://www.apache.org/licenses/LICENSE-2.0).

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.


---

## Article: document-conventions.md

---
title: Document Conventions
weight: 5
---

## Release Labels

Each minor release is supported until the next minor release becomes stable.

- **Dev**: The software is under active development and has not been thoroughly tested. Upgrading from earlier releases or to later releases is not supported. "Dev" releases typically occur at the end of a sprint, and may occur when other definitions of done are met. Given the unofficial nature of the release, use the software only for testing purposes.

- **Latest**: The software has passed all stages of verification and testing. All remaining known issues are considered acceptable. Use the software only for testing purposes.

- **Stable**: The software is considered reliable and free of critical issues. This release is recommended for general use.

- **EOL**: The software has reached the end of its useful life and no further code-level maintenance is provided. You may continue to use the software within the terms of the licensing agreement. Upgrading to a stable release is recommended.

> **Note**:
>
> Longhorn provides documentation for new features, enhancements, and bug fixes implemented in major, minor, and patch releases.

## Feature Labels

Features that are labeled **Experimental** and **Technical Preview** provide glimpses into upcoming innovations and offer opportunities to test new technologies within your environment.

- **Experimental**: The feature is incomplete but its essential functionality is operational. Tests in a controlled environment have shown that it can coexist with existing stable features. Because of the missing components and/or evolving functionality, the feature is not recommended for use in production environments.

- **Technical Preview**: The feature is almost complete and its functionality is not expected to change significantly. Tests in a controlled environment have shown that it can coexist with existing stable features. Explore the feature extensively before using in production environments.

These features have the following limitations:

- Still in development and may be functionally incomplete, unstable, or in other ways unsuitable for production use.
- Not supported.
- May only be available for specific hardware architectures. Details and functionality are subject to change. As a result, upgrading to subsequent releases may be impossible and may require a fresh installation.
- Can be removed from the product at any time. This may occur, for example, if we discover that the feature does not meet customer or market needs, or does not comply with enterprise standards.

> **Note**:
> 
> **Experimental** and **Technical Preview** features are usually disabled by default. The documentation provides information about how you can enable and configure such features.

## Admonitions

Admonitions are text blocks that provide additional or important information. Each text block is marked with a special label and icon.

- **Note**: Additional information that is useful but not critical.

- **Important**: Information that requires special attention, such as changes in product behavior.

- **Tip**: Helpful advice and suggestions, such as alternative methods of completing a task.

- **Caution**: Information that must be considered before proceeding with specific actions, such as errors that may occur.

- **Warning**: Critical information about harmful consequences, such as irreversible damage and data loss.

---

## Article: terminology.md

---
title: Terminology
weight: 3
---

- [Attach/Reattach](#attachreattach)
- [Backup](#backup)
- [Backupstore](#backupstore)
- [Backup target](#backup-target)
- [Backup volume](#backup-volume)
- [Block storage](#block-storage)
- [CRD](#crd)
- [CSI Driver](#csi-driver)
- [Disaster Recovery Volumes (DR volume)](#disaster-recovery-volumes-dr-volume)
- [ext4](#ext4)
- [Frontend expansion](#frontend-expansion)
- [Instance Manager](#instance-manager)
- [Longhorn volume](#longhorn-volume)
- [Mount](#mount)
- [NFS](#nfs)
- [Object storage](#object-storage)
- [Offline expansion](#offline-expansion)
- [Overprovisioning](#overprovisioning)
- [PersistentVolume](#persistentvolume)
- [PersistentVolumeClaim](#persistentvolumeclaim)
- [Primary backups](#primary-backups)
- [Remount](#remount)
- [Replica](#replica)
- [S3](#s3)
- [Salvage a volume](#salvage-a-volume)
- [Secondary backups](#secondary-backups)
- [Snapshot](#snapshot)
- [Stable identity](#stable-identity)
- [StatefulSet](#statefulset)
- [StorageClass](#storageclass)
- [System Backup](#system-backup)
- [Thin provisioning](#thin-provisioning)
- [Umount](#umount)
- [Volumes (Kubernetes concept)](#volumes-kubernetes-concept)
- [XFS](#xfs)
- [SMB/CIFS](#smbcifs)

### Attach/Reattach

To attach a block device is to make it appear on the Linux node, e.g. `/dev/longhorn/testvol`

If the volume engine dies unexpectedly, Longhorn will reattach the volume.

### Backup

A backup is an object in the backupstore. The backupstore may contain volume backups and system backups.

### Backupstore

Longhorn backups are saved to the backupstore, which is external to the Kubernetes cluster. The backupstore can be either NFS shares or an S3 compatible server.

Longhorn accesses the backupstore at the endpoint configured in the backuptarget.

### Backup target 

A backup target is the endpoint used to access a backupstore in Longhorn.

### Backup volume

A backup volume is the backup that maps to one original volume, and it is located in the backupstore. Backup volumes can be viewed on the **Backup** page in the Longhorn UI. The backup volume will contain multiple backups for the same volume.

Backups can be created from snapshots. They contain the state of the volume at the time the snapshot was created, but they don't contain snapshots, so they do not contain the history of changes to the volume data. While backups are made of 2 MB files, snapshots can be terabytes.

Backups are made of 2 MB blocks in an object store.

For a longer explanation of how snapshots and backups work, refer to the [conceptual documentation.](../concepts/#241-how-snapshots-work)

### Block storage

An approach to storage in which data stored in fixed-size blocks. Each block is distinguished based on a memory address.

### CRD

A Kubernetes [custom resource definition.](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)

### CSI Driver

The Longhorn CSI Driver is a [container storage interface](https://kubernetes-csi.github.io/docs/drivers.html) that can be used with Kubernetes. The CSI driver for Longhorn volumes is named `driver.longhorn.io`.

### Disaster Recovery Volumes (DR volume)

A DR volume is a special volume that stores data in a backup cluster in case the whole main cluster goes down. DR volumes are used to increase the resiliency of Longhorn volumes.

Each backup volume in the backupstore maps to one original volume in the Kubernetes cluster. Likewise, each DR volume maps to a backup volume in the backupstore.

DR volumes can be created to accurately reflect backups of a Longhorn volume, but they cannot be used as a normal Longhorn volume until they are activated.

### ext4

A file system for Linux. Longhorn supports ext4 for storage.

### Frontend expansion

The frontend here is referring to the block device exposed by the Longhorn volume.

### Instance Manager

The Longhorn component for controller/replica instance lifecycle management.

### Longhorn volume

A Longhorn volume is a Kubernetes volume that is replicated and managed by the Longhorn Manager. For each volume, the Longhorn Manager also creates:

- An instance of the Longhorn Engine
- Replicas of the volume, where each replica consists of a series of snapshots of the volume

Each replica contains a chain of snapshots, which record the changes in the volume's history. Three replicas are created by default, and they are usually stored on separate nodes for high availability.

### Mount

A Linux command to mount the block device to a certain directory on the node, e.g. `mount /dev/longhorn/testvol /mnt`

### NFS

A [distributed file system protocol](https://en.wikipedia.org/wiki/Network_File_System) that allows you to access files over a computer network, similar to the way that local storage is accessed. Longhorn supports using NFS as a backupstore for secondary storage.

### Object storage

Data storage architecture that manages data as objects. Each object typically includes the data itself, a variable amount of metadata, and a globally unique identifier.  Longhorn volumes can be backed up to S3 compatible object storage.

### Offline expansion

In an offline volume expansion, the volume is detached.

### Overprovisioning

Overprovisioning allows a server to view more storage capacity than has been physically reserved. That means we can schedule a total of 750 GiB Longhorn volumes on a 200 GiB disk with 50G reserved for the root file system. The **Storage Over Provisioning Percentage** can be configured in the Longhorn [settings.](../references/settings)

### PersistentVolume

A PersistentVolume (PV) is a Kubernetes resource that represents piece of storage in the cluster that has been provisioned by an administrator or dynamically provisioned using Storage Classes. It is a cluster-level resource, and is required for pods to use persistent storage that is independent of the lifecycle of any individual pod. For more information, see the official [Kubernetes documentation about persistent volumes.](https://kubernetes.io/docs/concepts/storage/persistent-volumes/)

### PersistentVolumeClaim

A PersistentVolumeClaim (PVC) is a request for storage by a user. Pods can request specific levels of resources (CPU and Memory) by using a PVC for storage. Claims can request specific sizes and access modes (e.g., they can be mounted once read/write or many times read-only).

For more information, see the official [Kubernetes documentation.](https://kubernetes.io/docs/concepts/storage/persistent-volumes/)

### Primary backups

The replicas of each Longhorn volume on a Kubernetes cluster can be considered primary backups.

### Remount

In a remount, Longhorn will detect and mount the filesystem for the volume after the reattachment.

### Replica

A replica consists of a chain of snapshots, showing a history of the changes in the data within a volume.

### S3

[Amazon S3](https://aws.amazon.com/s3/) is an object storage service.

### Salvage a volume

The salvage operation is needed when all replicas become faulty, e.g. due to a network disconnection.

When salvaging a volume, Longhorn will try to figure out which replica(s) are usable, then use them to recover the volume.

### Secondary backups

Backups external to the Kubernetes cluster, on S3 or NFS.

### Snapshot

A snapshot in Longhorn captures the state of a volume at the time the snapshot is created. Each snapshot only captures changes that overwrite data from earlier snapshots, so a sequence of snapshots is needed to fully represent the full state of the volume. Volumes can be restored from a snapshot. For a longer explanation of snapshots, refer to the [conceptual documentation.](../concepts)

### Stable identity

[StatefulSets](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/) have a stable identity, which means that Kubernetes won't force delete the Pod for the user.

### StatefulSet

A [Kubernetes resource](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/) used for managing stateful applications.

### StorageClass

A Kubernetes resource that can be used to automatically provision a PersistentVolume for a pod. For more information, refer to the [Kubernetes documentation.](https://kubernetes.io/docs/concepts/storage/storage-classes/#the-storageclass-resource)

### System Backup

Longhorn uploads the system backup to the backupstore. Each system backup contains the system backup resource bundle of the Longhorn system.

See [Longhorn System Backup Bundle](../advanced-resources/system-backup-restore/backup-longhorn-system/#longhorn-system-backup-bundle) for details.

### Thin provisioning

Longhorn is a thin-provisioned storage system. That means a Longhorn volume will only take the space it needs at the moment. For example, if you allocated a 20 GB volume but only use 1 GB of it, the actual data size on your disk would be 1GB. 

### Umount

A [Linux command](https://linux.die.net/man/8/umount) that detaches the file system from the file hierarchy.

### Volumes (Kubernetes concept)

A volume in Kubernetes allows a pod to store files during the lifetime of the pod.

These files will still be available after a container crashes, but they will not be available past the lifetime of a pod. To get storage that is still available after the lifetime of a pod, a Kubernetes [PersistentVolume (PV)](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#persistent-volumes) is required.

For more information, see the Kubernetes documentation on [volumes.](https://kubernetes.io/docs/concepts/storage/volumes/)

### XFS
A [file system](https://en.wikipedia.org/wiki/XFS) supported by most Linux distributions. Longhorn supports XFS for storage.

### SMB/CIFS

A [network file system protocol](https://en.wikipedia.org/wiki/Network_File_System) that allows you to access files over a computer network, similar to the way that local storage is accessed. Longhorn supports using SMB/CIFS as a backupstore for secondary storage.

---

## Article: what-is-longhorn.md

---
title: What is Longhorn?
weight: 1
---
Longhorn is a lightweight, reliable and easy-to-use distributed block storage system for Kubernetes.

Longhorn is free, open source software. Originally developed by Rancher Labs, it is now being developed as a incubating project of the Cloud Native Computing Foundation.

With Longhorn, you can:

- Use Longhorn volumes as persistent storage for the distributed stateful applications in your Kubernetes cluster
- Partition your block storage into Longhorn volumes so that you can use Kubernetes volumes with or without a cloud provider
- Replicate block storage across multiple nodes and data centers to increase availability
- Store backup data in external storage such as NFS or AWS S3
- Create cross-cluster disaster recovery volumes so that data from a primary Kubernetes cluster can be quickly recovered from backup in a second Kubernetes cluster
- Schedule recurring snapshots of a volume, and schedule recurring backups to NFS or S3-compatible secondary storage
- Restore volumes from backup
- Upgrade Longhorn without disrupting persistent volumes

Longhorn comes with a standalone UI, and can be installed using Helm, kubectl, or the Rancher app catalog.

### Simplifying Distributed Block Storage with Microservices

Because modern cloud environments require tens of thousands to millions of distributed block storage volumes, some storage controllers have become highly complex distributed systems. By contrast, Longhorn can simplify the storage system by partitioning a large block storage controller into a number of smaller storage controllers, as long as those volumes can still be built from a common pool of disks. By using one storage controller per volume, Longhorn turns each volume into a microservice. The controller is called the Longhorn Engine.

The Longhorn Manager component orchestrates the Longhorn Engines, so they work together coherently.

### Use Persistent Storage in Kubernetes without Relying on a Cloud Provider

Pods can reference storage directly, but this is not recommended because it doesn't allow the Pod or container to be portable. Instead, the workloads' storage requirements should be defined in Kubernetes Persistent Volumes (PVs) and Persistent Volume Claims (PVCs). With Longhorn, you can specify the size of the volume, the number of synchronous replicas and other volume specific configurations you want across the hosts that supply the storage resource for the volume. Then your Kubernetes resources can use the PVC and corresponding PV for each Longhorn volume, or use a Longhorn storage class to automatically create a PV for a workload.

Replicas are thin-provisioned on the underlying disks or network storage.

### Schedule Multiple Replicas across Multiple Compute or Storage Hosts

To increase availability, Longhorn creates replicas of each volume. Replicas contain a chain of snapshots of the volume, with each snapshot storing the change from a previous snapshot. Each replica of a volume also runs in a container, so a volume with three replicas results in four containers.

The number of replicas for each volume is configurable in Longhorn, as well as the nodes where replicas will be scheduled. Longhorn monitors the health of each replica and performs repairs, rebuilding the replica when necessary.

### Assign Multiple Storage Frontends for Each Volume

Common front-ends include a Linux kernel device (mapped under /dev/longhorn) and an iSCSI target.

### Specify Schedules for Recurring Snapshot and Backup Operations

Specify the frequency of these operations (hourly, daily, weekly, monthly, and yearly), the exact time at which these operations are performed (e.g., 3:00am every Sunday), and how many recurring snapshots and backup sets are kept.


---

## Article: references/_index.md

---
title: References
weight: 9
---


---

## Article: references/examples.md

---
title: Examples
weight: 4
---

For reference, this page provides examples of Kubernetes resources that use Longhorn storage.

- [Block volume](#block-volume)
- [CSI persistent volume](#csi-persistent-volume)
- [Deployment](#deployment)
- [Pod with PersistentVolumeClaim](#pod-with-persistentvolumeclaim)
- [Pod with Generic Ephemeral Volume](#pod-with-generic-ephemeral-volume)
- [Restore to file](#restore-to-file)
- [Simple Pod](#simple-pod)
- [Simple PersistentVolumeClaim](#simple-persistentvolumeclaim)
- [StatefulSet](#statefulset)
- [StorageClass](#storageclass)

### Block Volume

```yaml
    apiVersion: v1
    kind: PersistentVolumeClaim
    metadata:
      name: longhorn-block-vol
    spec:
      accessModes:
        - ReadWriteOnce
      volumeMode: Block
      storageClassName: longhorn
      resources:
        requests:
          storage: 2Gi
    ---
    apiVersion: v1
    kind: Pod
    metadata:
      name: block-volume-test
      namespace: default
    spec:
      containers:
        - name: block-volume-test
          image: nginx:stable-alpine
          imagePullPolicy: IfNotPresent
          volumeDevices:
            - devicePath: /dev/longhorn/testblk
              name: block-vol
          ports:
            - containerPort: 80
      volumes:
        - name: block-vol
          persistentVolumeClaim:
            claimName: longhorn-block-vol
```

### CSI Persistent Volume

```yaml
    apiVersion: v1
    kind: PersistentVolume
    metadata:
      name: longhorn-vol-pv
    spec:
      capacity:
        storage: 2Gi
      volumeMode: Filesystem
      accessModes:
        - ReadWriteOnce
      persistentVolumeReclaimPolicy: Delete
      storageClassName: longhorn
      csi:
        driver: driver.longhorn.io
        fsType: ext4
        volumeAttributes:
          numberOfReplicas: '3'
          staleReplicaTimeout: '2880'
        volumeHandle: existing-longhorn-volume
    ---
    apiVersion: v1
    kind: PersistentVolumeClaim
    metadata:
      name: longhorn-vol-pvc
    spec:
      accessModes:
        - ReadWriteOnce
      resources:
        requests:
          storage: 2Gi
      volumeName: longhorn-vol-pv
      storageClassName: longhorn
    ---
    apiVersion: v1
    kind: Pod
    metadata:
      name: volume-pv-test
      namespace: default
    spec:
      restartPolicy: Always
      containers:
      - name: volume-pv-test
        image: nginx:stable-alpine
        imagePullPolicy: IfNotPresent
        livenessProbe:
          exec:
            command:
              - ls
              - /data/lost+found
          initialDelaySeconds: 5
          periodSeconds: 5
        volumeMounts:
        - name: vol
          mountPath: /data
        ports:
        - containerPort: 80
      volumes:
      - name: vol
        persistentVolumeClaim:
          claimName: longhorn-vol-pvc
```

### Deployment

```yaml
    apiVersion: v1
    kind: Service
    metadata:
      name: mysql
      labels:
        app: mysql
    spec:
      ports:
        - port: 3306
      selector:
        app: mysql
      clusterIP: None
    ---
    apiVersion: v1
    kind: PersistentVolumeClaim
    metadata:
      name: mysql-pvc
    spec:
      accessModes:
        - ReadWriteOnce
      storageClassName: longhorn
      resources:
        requests:
          storage: 2Gi
    ---
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: mysql
      labels:
        app: mysql
    spec:
      selector:
        matchLabels:
          app: mysql # has to match .spec.template.metadata.labels
      strategy:
        type: Recreate
      template:
        metadata:
          labels:
            app: mysql
        spec:
          restartPolicy: Always
          containers:
          - image: mysql:5.6
            name: mysql
            livenessProbe:
              exec:
                command:
                  - ls
                  - /var/lib/mysql/lost+found
              initialDelaySeconds: 5
              periodSeconds: 5
            env:
            - name: MYSQL_ROOT_PASSWORD
              value: changeme
            ports:
            - containerPort: 3306
              name: mysql
            volumeMounts:
            - name: mysql-volume
              mountPath: /var/lib/mysql
          volumes:
          - name: mysql-volume
            persistentVolumeClaim:
              claimName: mysql-pvc
```

### Pod with PersistentVolumeClaim

```yaml
    apiVersion: v1
    kind: PersistentVolumeClaim
    metadata:
      name: longhorn-volv-pvc
    spec:
      accessModes:
        - ReadWriteOnce
      storageClassName: longhorn
      resources:
        requests:
          storage: 2Gi
    ---
    apiVersion: v1
    kind: Pod
    metadata:
      name: volume-test
      namespace: default
    spec:
      restartPolicy: Always
      containers:
      - name: volume-test
        image: nginx:stable-alpine
        imagePullPolicy: IfNotPresent
        livenessProbe:
          exec:
            command:
              - ls
              - /data/lost+found
          initialDelaySeconds: 5
          periodSeconds: 5
        volumeMounts:
        - name: volv
          mountPath: /data
        ports:
        - containerPort: 80
      volumes:
      - name: volv
        persistentVolumeClaim:
          claimName: longhorn-volv-pvc
```

### Pod with Generic Ephemeral Volume

For more information about generic ephemeral volumes, refer to the
[Kubernetes documentation](https://kubernetes.io/docs/concepts/storage/ephemeral-volumes/#generic-ephemeral-volumes).

```yaml
  apiVersion: v1
  kind: Pod
  metadata:
    name: volume-test
    namespace: default
  spec:
    restartPolicy: Always
    containers:
    - name: volume-test
      image: nginx:stable-alpine
      imagePullPolicy: IfNotPresent
      livenessProbe:
        exec:
          command:
            - ls
            - /data/lost+found
        initialDelaySeconds: 5
        periodSeconds: 5
      volumeMounts:
      - name: volv
        mountPath: /data
      ports:
      - containerPort: 80
    volumes:
    - name: volv
      ephemeral:
        volumeClaimTemplate:
          spec:
            accessModes:
              - ReadWriteOnce
            storageClassName: longhorn
            resources:
              requests:
                storage: 2Gi
```

### Restore to File

For more information about restoring to file, refer to [this section.](../../advanced-resources/data-recovery/recover-without-system)

```yaml
    apiVersion: v1
    kind: Pod
    metadata:
      name: restore-to-file
      namespace: longhorn-system
    spec:
      nodeName: <NODE_NAME>
      containers:
      - name: restore-to-file
        command:
        # set restore-to-file arguments here
        - /bin/sh
        - -c
        - longhorn backup restore-to-file
          '<BACKUP_URL>'
          --output-file '/tmp/restore/<OUTPUT_FILE>'
          --output-format <OUTPUT_FORMAT>
        # the version of longhorn engine should be v0.4.1 or higher
        image: longhorn/longhorn-engine:v0.4.1
        imagePullPolicy: IfNotPresent
        securityContext:
          privileged: true
        volumeMounts:
        - name: disk-directory
          mountPath: /tmp/restore  # the argument <output-file> should be in this directory
        env:
        # set Backup Target Credential Secret here.
        - name: AWS_ACCESS_KEY_ID
          valueFrom:
            secretKeyRef:
              name: <S3_SECRET_NAME>
              key: AWS_ACCESS_KEY_ID
        - name: AWS_SECRET_ACCESS_KEY
          valueFrom:
            secretKeyRef:
              name: <S3_SECRET_NAME>
              key: AWS_SECRET_ACCESS_KEY
        - name: AWS_ENDPOINTS
          valueFrom:
            secretKeyRef:
              name: <S3_SECRET_NAME>
              key: AWS_ENDPOINTS
      volumes:
        # the output file can be found on this host path
        - name: disk-directory
          hostPath:
            path: /tmp/restore
      restartPolicy: Never
```

### Simple Pod

```yaml
    apiVersion: v1
    kind: Pod
    metadata:
      name: longhorn-simple-pod
      namespace: default
    spec:
      restartPolicy: Always
      containers:
        - name: volume-test
          image: nginx:stable-alpine
          imagePullPolicy: IfNotPresent
          livenessProbe:
            exec:
              command:
                - ls
                - /data/lost+found
            initialDelaySeconds: 5
            periodSeconds: 5
          volumeMounts:
            - name: volv
              mountPath: /data
          ports:
            - containerPort: 80
      volumes:
        - name: volv
          persistentVolumeClaim:
            claimName: longhorn-simple-pvc
```

### Simple PersistentVolumeClaim

```yaml
    apiVersion: v1
    kind: PersistentVolumeClaim
    metadata:
      name: longhorn-simple-pvc
    spec:
      accessModes:
        - ReadWriteOnce
      storageClassName: longhorn
      resources:
        requests:
          storage: 1Gi
```

### StatefulSet

```yaml
    apiVersion: v1
    kind: Service
    metadata:
      name: nginx
      labels:
        app: nginx
    spec:
      ports:
      - port: 80
        name: web
      selector:
        app: nginx
      type: NodePort
    ---
    apiVersion: apps/v1
    kind: StatefulSet
    metadata:
      name: web
    spec:
      selector:
        matchLabels:
          app: nginx # has to match .spec.template.metadata.labels
      serviceName: "nginx"
      replicas: 2 # by default is 1
      template:
        metadata:
          labels:
            app: nginx # has to match .spec.selector.matchLabels
        spec:
          restartPolicy: Always
          terminationGracePeriodSeconds: 10
          containers:
          - name: nginx
            image: registry.k8s.io/nginx-slim:0.8
            livenessProbe:
              exec:
                command:
                  - ls
                  - /usr/share/nginx/html/lost+found
              initialDelaySeconds: 5
              periodSeconds: 5
            ports:
            - containerPort: 80
              name: web
            volumeMounts:
            - name: www
              mountPath: /usr/share/nginx/html
      volumeClaimTemplates:
      - metadata:
          name: www
        spec:
          accessModes: [ "ReadWriteOnce" ]
          storageClassName: "longhorn"
          resources:
            requests:
              storage: 1Gi
```

### StorageClass

```yaml
    kind: StorageClass
    apiVersion: storage.k8s.io/v1
    metadata:
      name: longhorn
    provisioner: driver.longhorn.io
    allowVolumeExpansion: true
    reclaimPolicy: Delete
    volumeBindingMode: Immediate
    parameters:
      numberOfReplicas: "3"
      staleReplicaTimeout: "2880" # 48 hours in minutes
      fromBackup: ""
      fsType: "ext4"
      #  mkfsParams: "-I 256 -b 4096 -O ^metadata_csum,^64bit"
      #  backingImage: "bi-test"
      #  backingImageDataSourceType: "download"
      #  backingImageDataSourceParameters: '{"url": "https://backing-image-example.s3-region.amazonaws.com/test-backing-image"}'
      #  backingImageChecksum: "SHA512 checksum of the backing image"
      #  diskSelector: "ssd,fast"
      #  nodeSelector: "storage,fast"
      #  recurringJobSelector: '[
      #   {
      #     "name":"snap",
      #     "isGroup":true,
      #   },
      #   {
      #     "name":"backup",
      #     "isGroup":false,
      #   }
      #  ]'

```

Note that Longhorn supports automatic remount only for the workload pod that is managed by a controller (e.g. deployment, statefulset, daemonset, etc...).
See [here](../../high-availability/recover-volume/) for details.


---

## Article: references/helm-values.md

---
title: Helm Values
weight: 5
---

## Values

The `values.yaml` file contains items used to tweak a deployment of this chart.

### Cattle Settings

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| global.cattle.systemDefaultRegistry | string | `""` | Default system registry. |
| global.cattle.windowsCluster.defaultSetting.systemManagedComponentsNodeSelector | string | `"kubernetes.io/os:linux"` | Node selector for system-managed Longhorn components. |
| global.cattle.windowsCluster.defaultSetting.taintToleration | string | `"cattle.io/os=linux:NoSchedule"` | Toleration for system-managed Longhorn components. |
| global.cattle.windowsCluster.enabled | bool | `false` | Setting that allows Longhorn to run on a Rancher Windows cluster. |
| global.cattle.windowsCluster.nodeSelector | object | `{"kubernetes.io/os":"linux"}` | Node selector for Linux nodes that can run user-deployed Longhorn components. |
| global.cattle.windowsCluster.tolerations | list | `[{"effect":"NoSchedule","key":"cattle.io/os","operator":"Equal","value":"linux"}]` | Toleration for Linux nodes that can run user-deployed Longhorn components. |

### Network Policies

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| networkPolicies.enabled | bool | `false` | Setting that allows you to enable network policies that control access to Longhorn pods. |
| networkPolicies.type | string | `"k3s"` | Distribution that determines the policy for allowing access for an ingress. (Options: "k3s", "rke2", "rke1") |

### Image Settings

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| image.csi.attacher.repository | string | `"longhornio/csi-attacher"` | Repository for the CSI attacher image. When unspecified, Longhorn uses the default value. |
| image.csi.attacher.tag | string | `"v4.10.0-20251030"` | Tag for the CSI attacher image. When unspecified, Longhorn uses the default value. |
| image.csi.livenessProbe.repository | string | `"longhornio/livenessprobe"` | Repository for the CSI liveness probe image. When unspecified, Longhorn uses the default value. |
| image.csi.livenessProbe.tag | string | `"v2.17.0-20251030"` | Tag for the CSI liveness probe image. When unspecified, Longhorn uses the default value. |
| image.csi.nodeDriverRegistrar.repository | string | `"longhornio/csi-node-driver-registrar"` | Repository for the CSI Node Driver Registrar image. When unspecified, Longhorn uses the default value. |
| image.csi.nodeDriverRegistrar.tag | string | `"v2.15.0-20251030"` | Tag for the CSI Node Driver Registrar image. When unspecified, Longhorn uses the default value. |
| image.csi.provisioner.repository | string | `"longhornio/csi-provisioner"` | Repository for the CSI Provisioner image. When unspecified, Longhorn uses the default value. |
| image.csi.provisioner.tag | string | `"v5.3.0-20251030"` | Tag for the CSI Provisioner image. When unspecified, Longhorn uses the default value. |
| image.csi.resizer.repository | string | `"longhornio/csi-resizer"` | Repository for the CSI Resizer image. When unspecified, Longhorn uses the default value. |
| image.csi.resizer.tag | string | `"v1.14.0-20251030"` | Tag for the CSI Resizer image. When unspecified, Longhorn uses the default value. |
| image.csi.snapshotter.repository | string | `"longhornio/csi-snapshotter"` | Repository for the CSI Snapshotter image. When unspecified, Longhorn uses the default value. |
| image.csi.snapshotter.tag | string | `"v8.4.0-20251030"` | Tag for the CSI Snapshotter image. When unspecified, Longhorn uses the default value. |
| image.longhorn.backingImageManager.repository | string | `"longhornio/backing-image-manager"` | Repository for the Backing Image Manager image. When unspecified, Longhorn uses the default value. |
| image.longhorn.backingImageManager.tag | string | `"v1.10.1"` | Tag for the Backing Image Manager image. When unspecified, Longhorn uses the default value. |
| image.longhorn.engine.repository | string | `"longhornio/longhorn-engine"` | Repository for the Longhorn Engine image. |
| image.longhorn.engine.tag | string | `"v1.10.1"` | Tag for the Longhorn Engine image. |
| image.longhorn.instanceManager.repository | string | `"longhornio/longhorn-instance-manager"` | Repository for the Longhorn Instance Manager image. |
| image.longhorn.instanceManager.tag | string | `"v1.10.1"` | Tag for the Longhorn Instance Manager image. |
| image.longhorn.manager.repository | string | `"longhornio/longhorn-manager"` | Repository for the Longhorn Manager image. |
| image.longhorn.manager.tag | string | `"v1.10.1"` | Tag for the Longhorn Manager image. |
| image.longhorn.shareManager.repository | string | `"longhornio/longhorn-share-manager"` | Repository for the Longhorn Share Manager image. |
| image.longhorn.shareManager.tag | string | `"v1.10.1"` | Tag for the Longhorn Share Manager image. |
| image.longhorn.supportBundleKit.repository | string | `"longhornio/support-bundle-kit"` | Repository for the Longhorn Support Bundle Manager image. |
| image.longhorn.supportBundleKit.tag | string | `"v0.0.71"` | Tag for the Longhorn Support Bundle Manager image. |
| image.longhorn.ui.repository | string | `"longhornio/longhorn-ui"` | Repository for the Longhorn UI image. |
| image.longhorn.ui.tag | string | `"v1.10.1"` | Tag for the Longhorn UI image. |
| image.openshift.oauthProxy.repository | string | `""` | Repository for the OAuth Proxy image. Specify the upstream image (for example, "quay.io/openshift/origin-oauth-proxy"). This setting applies only to OpenShift users. |
| image.openshift.oauthProxy.tag | float | `""` | Tag for the OAuth Proxy image. Specify OCP/OKD version 4.1 or later (including version 4.15, which is available at quay.io/openshift/origin-oauth-proxy:4.15). This setting applies only to OpenShift users. |
| image.pullPolicy | string | `"IfNotPresent"` | Image pull policy that applies to all user-deployed Longhorn components, such as Longhorn Manager, Longhorn driver, and Longhorn UI. |

### Service Settings

| Key | Description |
|-----|-------------|
| service.manager.nodePort | NodePort port number for Longhorn Manager. When unspecified, Longhorn selects a free port between 30000 and 32767. |
| service.manager.type | Service type for Longhorn Manager. |
| service.ui.nodePort | NodePort port number for Longhorn UI. When unspecified, Longhorn selects a free port between 30000 and 32767. |
| service.ui.type | Service type for Longhorn UI. (Options: "ClusterIP", "NodePort", "LoadBalancer", "Rancher-Proxy") |

### StorageClass Settings

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| persistence.backingImage.dataSourceParameters | string | `nil` | Data source parameters of a backing image used in a Longhorn StorageClass. You can specify a JSON string of a map. (Example: `'{\"url\":\"https://backing-image-example.s3-region.amazonaws.com/test-backing-image\"}'`) |
| persistence.backingImage.dataSourceType | string | `nil` | Data source type of a backing image used in a Longhorn StorageClass. If the backing image exists in the cluster, Longhorn uses this setting to verify the image. If the backing image does not exist, Longhorn creates one using the specified data source type. |
| persistence.backingImage.enable | bool | `false` | Setting that allows you to use a backing image in a Longhorn StorageClass. |
| persistence.backingImage.expectedChecksum | string | `nil` | Expected SHA-512 checksum of a backing image used in a Longhorn StorageClass. |
| persistence.backingImage.name | string | `nil` | Backing image to be used for creating and restoring volumes in a Longhorn StorageClass. When no backing images are available, specify the data source type and parameters that Longhorn can use to create a backing image. |
| persistence.defaultClass | bool | `true` | Setting that allows you to specify the default Longhorn StorageClass. |
| persistence.defaultClassReplicaCount | int | `3` | Replica count of the default Longhorn StorageClass. |
| persistence.defaultDataLocality | string | `"disabled"` | Data locality of the default Longhorn StorageClass. (Options: "disabled", "best-effort") |
| persistence.defaultFsType | string | `"ext4"` | Filesystem type of the default Longhorn StorageClass. |
| persistence.defaultMkfsParams | string | `""` | mkfs parameters of the default Longhorn StorageClass. |
| persistence.defaultNodeSelector.enable | bool | `false` | Setting that allows you to enable the node selector for the default Longhorn StorageClass. |
| persistence.defaultNodeSelector.selector | string | `""` | Node selector for the default Longhorn StorageClass. Longhorn uses only nodes with the specified tags for storing volume data. (Examples: "storage,fast") |
| persistence.disableRevisionCounter | string | `"true"` | Setting that disables the revision counter and thereby prevents Longhorn from tracking all write operations to a volume. When salvaging a volume, Longhorn uses properties of the volume-head-xxx.img file (the last file size and the last time the file was modified) to select the replica to be used for volume recovery. |
| persistence.migratable | bool | `false` | Setting that allows you to enable live migration of a Longhorn volume from one node to another. |
| persistence.nfsOptions | string | `""` | Set NFS mount options for Longhorn StorageClass for RWX volumes |
| persistence.reclaimPolicy | string | `"Delete"` | Reclaim policy that provides instructions for handling of a volume after its claim is released. (Options: "Retain", "Delete") |
| persistence.recurringJobSelector.enable | bool | `false` | Setting that allows you to enable the recurring job selector for a Longhorn StorageClass. |
| persistence.recurringJobSelector.jobList | list | `[]` | Recurring job selector for a Longhorn StorageClass. Ensure that quotes are used correctly when specifying job parameters. (Example: `[{"name":"backup", "isGroup":true}]`) |
| persistence.removeSnapshotsDuringFilesystemTrim | string | `"ignored"` | Setting that allows you to enable automatic snapshot removal during filesystem trim for a Longhorn StorageClass. (Options: "ignored", "enabled", "disabled") |

### CSI Settings

| Key | Description |
|-----|-------------|
| csi.attacherReplicaCount | Replica count of the CSI Attacher. When unspecified, Longhorn uses the default value ("3"). |
| csi.kubeletRootDir | kubelet root directory. When unspecified, Longhorn uses the default value. |
| csi.provisionerReplicaCount | Replica count of the CSI Provisioner. When unspecified, Longhorn uses the default value ("3"). |
| csi.resizerReplicaCount | Replica count of the CSI Resizer. When unspecified, Longhorn uses the default value ("3"). |
| csi.snapshotterReplicaCount | Replica count of the CSI Snapshotter. When unspecified, Longhorn uses the default value ("3"). |

### Longhorn Manager Settings

Longhorn consists of user-deployed components (for example, Longhorn Manager, Longhorn Driver, and Longhorn UI) and system-managed components (for example, Instance Manager, Backing Image Manager, Share Manager, CSI Driver, and Engine Image). The following settings only apply to Longhorn Manager.

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| longhornManager.log.format | string | `"plain"` | Format of Longhorn Manager logs. (Options: "plain", "json") |
| longhornManager.nodeSelector | object | `{}` | Node selector for Longhorn Manager. Specify the nodes allowed to run Longhorn Manager. |
| longhornManager.priorityClass | string | `"longhorn-critical"` | PriorityClass for Longhorn Manager. |
| longhornManager.serviceAnnotations | object | `{}` | Annotation for the Longhorn Manager service. |
| longhornManager.tolerations | list | `[]` | Toleration for Longhorn Manager on nodes allowed to run Longhorn Manager. |

### Longhorn Driver Settings

Longhorn consists of user-deployed components (for example, Longhorn Manager, Longhorn Driver, and Longhorn UI) and system-managed components (for example, Instance Manager, Backing Image Manager, Share Manager, CSI Driver, and Engine Image). The following settings only apply to Longhorn Driver.

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| longhornDriver.nodeSelector | object | `{}` | Node selector for Longhorn Driver. Specify the nodes allowed to run Longhorn Driver. |
| longhornDriver.priorityClass | string | `"longhorn-critical"` | PriorityClass for Longhorn Driver. |
| longhornDriver.tolerations | list | `[]` | Toleration for Longhorn Driver on nodes allowed to run Longhorn components. |

### Longhorn UI Settings

Longhorn consists of user-deployed components (for example, Longhorn Manager, Longhorn Driver, and Longhorn UI) and system-managed components (for example, Instance Manager, Backing Image Manager, Share Manager, CSI Driver, and Engine Image). The following settings only apply to Longhorn UI.

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| longhornUI.nodeSelector | object | `{}` | Node selector for Longhorn UI. Specify the nodes allowed to run Longhorn UI. |
| longhornUI.priorityClass | string | `"longhorn-critical"` | PriorityClass for Longhorn UI. |
| longhornUI.replicas | int | `2` | Replica count for Longhorn UI. |
| longhornUI.tolerations | list | `[]` | Toleration for Longhorn UI on nodes allowed to run Longhorn components. |

### Ingress Settings

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| ingress.annotations | string | `nil` | Ingress annotations in the form of key-value pairs. |
| ingress.enabled | bool | `false` | Setting that allows Longhorn to generate ingress records for the Longhorn UI service. |
| ingress.host | string | `"sslip.io"` | Hostname of the Layer 7 load balancer. |
| ingress.ingressClassName | string | `nil` | IngressClass resource that contains ingress configuration, including the name of the Ingress controller. ingressClassName can replace the kubernetes.io/ingress.class annotation used in earlier Kubernetes releases. |
| ingress.path | string | `"/"` | Default ingress path. You can access the Longhorn UI by following the full ingress path {{host}}+{{path}}. |
| ingress.pathType | string | `"ImplementationSpecific"` | Ingress path type. To maintain backward compatibility, the default value is "ImplementationSpecific". |
| ingress.secrets | string | `nil` | Secret that contains a TLS private key and certificate. Use secrets if you want to use your own certificates to secure ingresses. |
| ingress.secureBackends | bool | `false` | Setting that allows you to enable secure connections to the Longhorn UI service via port 443. |
| ingress.tls | bool | `false` | Setting that allows you to enable TLS on ingress records. |
| ingress.tlsSecret | string | `"longhorn.local-tls"` | TLS secret that contains the private key and certificate to be used for TLS. This setting applies only when TLS is enabled on ingress records. |

### Private Registry Settings

You can install Longhorn in an air-gapped environment with a private registry. For more information, see the **Air Gap Installation** section of the [documentation](https://longhorn.io/docs).

| Key | Description |
|-----|-------------|
| privateRegistry.createSecret | Setting that allows you to create a private registry secret. |
| privateRegistry.registryPasswd | Password for authenticating with a private registry. |
| privateRegistry.registrySecret | Kubernetes secret that allows you to pull images from a private registry. This setting applies only when creation of private registry secrets is enabled. You must include the private registry name in the secret name. |
| privateRegistry.registryUrl | URL of a private registry. When unspecified, Longhorn uses the default system registry. |
| privateRegistry.registryUser | User account used for authenticating with a private registry. |

### Metrics Settings

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| metrics.serviceMonitor.additionalLabels | object | `{}` | Additional labels for the Prometheus ServiceMonitor resource. |
| metrics.serviceMonitor.annotations | object | `{}` | Annotations for the Prometheus ServiceMonitor resource. |
| metrics.serviceMonitor.enabled | bool | `false` | Setting that allows the creation of a Prometheus ServiceMonitor resource for Longhorn Manager components. |
| metrics.serviceMonitor.interval | string | `""` | Interval at which Prometheus scrapes the metrics from the target. |
| metrics.serviceMonitor.metricRelabelings | list | `[]` | Configures the relabeling rules to apply to the samples before ingestion. See the [Prometheus Operator  documentation](https://prometheus-operator.dev/docs/api-reference/api/#monitoring.coreos.com/v1.Endpoint) for formatting details. |
| metrics.serviceMonitor.relabelings | list | `[]` | Configures the relabeling rules to apply the target’s metadata labels. See the [Prometheus Operator  documentation](https://prometheus-operator.dev/docs/api-reference/api/#monitoring.coreos.com/v1.Endpoint) for formatting details. |
| metrics.serviceMonitor.scrapeTimeout | string | `""` | Timeout after which Prometheus considers the scrape to be failed. |

### OS/Kubernetes Distro Settings

#### OpenShift Settings

For more details, see the [ocp-readme](https://github.com/longhorn/longhorn/blob/master/chart/ocp-readme.md).

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| openshift.enabled | bool | `false` | Setting that allows Longhorn to integrate with OpenShift. |
| openshift.ui.port | int | `443` | Port for accessing the OpenShift web console. |
| openshift.ui.proxy | int | `8443` | Port for proxy that provides access to the OpenShift web console. |
| openshift.ui.route | string | `"longhorn-ui"` | Route for connections between Longhorn and the OpenShift web console. |

### Other Settings

| Key | Default | Description |
|-----|---------|-------------|
| annotations | `{}` | Annotation for the Longhorn Manager DaemonSet pods. This setting is optional. |
| defaultBackupStore | `{"backupTarget":null,"backupTargetCredentialSecret":null,"pollInterval":null}` | Setting that allows you to update the default backupstore. |
| defaultBackupStore.backupTarget | `""` | Endpoint used to access the default backupstore. (Options: "NFS", "CIFS", "AWS", "GCP", "AZURE") |
| defaultBackupStore.backupTargetCredentialSecret | `""` | Name of the Kubernetes secret associated with the default backup target. |
| defaultBackupStore.pollInterval | `""` | Number of seconds that Longhorn waits before checking the default backupstore for new backups. The default value is "300". When the value is "0", polling is disabled. |
| enableGoCoverDir | `false` | Setting that allows Longhorn to generate code coverage profiles. |
| enablePSP | `false` | Setting that allows you to enable pod security policies (PSPs) that allow privileged Longhorn pods to start. This setting applies only to clusters running Kubernetes 1.25 and earlier, and with the built-in Pod Security admission controller enabled. |
| namespaceOverride | `""` | Specify override namespace, specifically this is useful for using longhorn as sub-chart and its release namespace is not the `longhorn-system`. |
| preUpgradeChecker.jobEnabled | `true` | Setting that allows Longhorn to perform pre-upgrade checks. Disable this setting when installing Longhorn using Argo CD or other GitOps solutions. |
| preUpgradeChecker.upgradeVersionCheck | `true` | Setting that allows Longhorn to perform upgrade version checks after starting the Longhorn Manager DaemonSet Pods. Disabling this setting also disables `preUpgradeChecker.jobEnabled`. Longhorn recommends keeping this setting enabled. |

### System Default Settings

During installation, you can either allow Longhorn to use the default system settings or use specific flags to modify the default values. After installation, you can modify the settings using the Longhorn UI. For more information, see the **Settings Reference** section of the [documentation](https://longhorn.io/docs).

| Key | Description |
|-----|-------------|
| defaultSettings.allowCollectingLonghornUsageMetrics | Setting that allows Longhorn to periodically collect anonymous usage data for product improvement purposes. Longhorn sends collected data to the [Upgrade Responder](https://github.com/longhorn/upgrade-responder) server, which is the data source of the Longhorn Public Metrics Dashboard (https://metrics.longhorn.io). The Upgrade Responder server does not store data that can be used to identify clients, including IP addresses. |
| defaultSettings.allowEmptyDiskSelectorVolume | Setting that allows scheduling of empty disk selector volumes to any disk. |
| defaultSettings.allowEmptyNodeSelectorVolume | Setting that allows scheduling of empty node selector volumes to any node. |
| defaultSettings.allowRecurringJobWhileVolumeDetached | Setting that allows Longhorn to automatically attach a volume and create snapshots or backups when recurring jobs are run. |
| defaultSettings.allowVolumeCreationWithDegradedAvailability | Setting that allows you to create and attach a volume without having all replicas scheduled at the time of creation. |
| defaultSettings.autoCleanupRecurringJobBackupSnapshot | Setting that allows Longhorn to automatically clean up the snapshot generated by a recurring backup job. |
| defaultSettings.autoCleanupSystemGeneratedSnapshot | Setting that allows Longhorn to automatically clean up the system-generated snapshot after replica rebuilding is completed. |
| defaultSettings.autoDeletePodWhenVolumeDetachedUnexpectedly | Setting that allows Longhorn to automatically delete a workload pod that is managed by a controller (for example, daemonset) whenever a Longhorn volume is detached unexpectedly (for example, during Kubernetes upgrades). After deletion, the controller restarts the pod and then Kubernetes handles volume reattachment and remounting. |
| defaultSettings.autoSalvage | Setting that allows Longhorn to automatically salvage volumes when all replicas become faulty (for example, when the network connection is interrupted). Longhorn determines which replicas are usable and then uses these replicas for the volume. This setting is enabled by default. |
| defaultSettings.backingImageCleanupWaitInterval | Number of minutes that Longhorn waits before cleaning up the backing image file when no replicas in the disk are using it. |
| defaultSettings.backingImageRecoveryWaitInterval | Number of seconds that Longhorn waits before downloading a backing image file again when the status of all image disk files changes to "failed" or "unknown". |
| defaultSettings.backupCompressionMethod | Setting that allows you to specify a backup compression method. |
| defaultSettings.backupConcurrentLimit | Maximum number of worker threads that can concurrently run for each backup. |
| defaultSettings.concurrentAutomaticEngineUpgradePerNodeLimit | Maximum number of engines that are allowed to concurrently upgrade on each node after Longhorn Manager is upgraded. When the value is "0", Longhorn does not automatically upgrade volume engines to the new default engine image version. |
| defaultSettings.concurrentReplicaRebuildPerNodeLimit | Maximum number of replicas that can be concurrently rebuilt on each node. |
| defaultSettings.concurrentVolumeBackupRestorePerNodeLimit | Maximum number of volumes that can be concurrently restored on each node using a backup. When the value is "0", restoration of volumes using a backup is disabled. |
| defaultSettings.createDefaultDiskLabeledNodes | Setting that allows Longhorn to automatically create a default disk only on nodes with the label "node.longhorn.io/create-default-disk=true" (if no other disks exist). When this setting is disabled, Longhorn creates a default disk on each node that is added to the cluster. |
| defaultSettings.defaultDataLocality | Default data locality. A Longhorn volume has data locality if a local replica of the volume exists on the same node as the pod that is using the volume. |
| defaultSettings.defaultDataPath | Default path for storing data on a host. The default value is "/var/lib/longhorn/". |
| defaultSettings.defaultLonghornStaticStorageClass | Default Longhorn StorageClass. "storageClassName" is assigned to PVs and PVCs that are created for an existing Longhorn volume. "storageClassName" can also be used as a label, so it is possible to use a Longhorn StorageClass to bind a workload to an existing PV without creating a Kubernetes StorageClass object. The default value is "longhorn-static". |
| defaultSettings.defaultReplicaCount | Default number of replicas for volumes created using the Longhorn UI. For Kubernetes configuration, modify the `numberOfReplicas` field in the StorageClass. The default value is "3". |
| defaultSettings.deletingConfirmationFlag | Flag that prevents accidental uninstallation of Longhorn. |
| defaultSettings.detachManuallyAttachedVolumesWhenCordoned | Setting that allows automatic detaching of manually-attached volumes when a node is cordoned. |
| defaultSettings.disableRevisionCounter | Setting that disables the revision counter and thereby prevents Longhorn from tracking all write operations to a volume. When salvaging a volume, Longhorn uses properties of the "volume-head-xxx.img" file (the last file size and the last time the file was modified) to select the replica to be used for volume recovery. This setting applies only to volumes created using the Longhorn UI. |
| defaultSettings.disableSchedulingOnCordonedNode | Setting that prevents Longhorn Manager from scheduling replicas on a cordoned Kubernetes node. This setting is enabled by default. |
| defaultSettings.disableSnapshotPurge | Setting that temporarily prevents all attempts to purge volume snapshots. |
| defaultSettings.engineReplicaTimeout | Timeout between the Longhorn Engine and replicas. Specify a value between "8" and "30" seconds. The default value is "8". |
| defaultSettings.failedBackupTTL | Number of minutes that Longhorn keeps a failed backup resource. When the value is "0", automatic deletion is disabled. |
| defaultSettings.fastReplicaRebuildEnabled | Setting that allows fast rebuilding of replicas using the checksum of snapshot disk files. Before enabling this setting, you must set the snapshot-data-integrity value to "enable" or "fast-check". |
| defaultSettings.freezeFilesystemForSnapshot | Setting that freezes the filesystem on the root partition before a snapshot is created. |
| defaultSettings.guaranteedInstanceManagerCPU | Percentage of the total allocatable CPU resources on each node to be reserved for each instance manager pod when the V1 Data Engine is enabled. The default value is "12". |
| defaultSettings.kubernetesClusterAutoscalerEnabled | Setting that notifies Longhorn that the cluster is using the Kubernetes Cluster Autoscaler. |
| defaultSettings.logLevel | Log levels that indicate the type and severity of logs in Longhorn Manager. The default value is "Info". (Options: "Panic", "Fatal", "Error", "Warn", "Info", "Debug", "Trace") |
| defaultSettings.longGRPCTimeOut | Number of seconds that Longhorn allows for the completion of replica rebuilding and snapshot cloning operations. |
| defaultSettings.nodeDownPodDeletionPolicy | Policy that defines the action Longhorn takes when a volume is stuck with a StatefulSet or Deployment pod on a node that failed. |
| defaultSettings.nodeDrainPolicy | Policy that defines the action Longhorn takes when a node with the last healthy replica of a volume is drained. |
| defaultSettings.offlineReplicaRebuilding | Enables automatic rebuilding of degraded replicas while the volume is detached. This setting only takes effect if the individual volume setting is set to `ignored` or `enabled`. |
| defaultSettings.orphanAutoDeletion | Setting that allows Longhorn to automatically delete an orphaned resource and the corresponding data (for example, stale replicas). Orphaned resources on failed or unknown nodes are not automatically cleaned up. |
| defaultSettings.priorityClass | PriorityClass for system-managed Longhorn components. This setting can help prevent Longhorn components from being evicted under Node Pressure. Notice that this will be applied to Longhorn user-deployed components by default if there are no priority class values set yet, such as `longhornManager.priorityClass`. |
| defaultSettings.recurringFailedJobsHistoryLimit | Maximum number of failed recurring backup and snapshot jobs to be retained. When the value is "0", a history of failed recurring jobs is not retained. |
| defaultSettings.recurringJobMaxRetention | Maximum number of snapshots or backups to be retained. |
| defaultSettings.recurringSuccessfulJobsHistoryLimit | Maximum number of successful recurring backup and snapshot jobs to be retained. When the value is "0", a history of successful recurring jobs is not retained. |
| defaultSettings.removeSnapshotsDuringFilesystemTrim | Setting that allows Longhorn to automatically mark the latest snapshot and its parent files as removed during a filesystem trim. Longhorn does not remove snapshots containing multiple child files. |
| defaultSettings.replicaAutoBalance | Setting that automatically rebalances replicas when an available node is discovered. |
| defaultSettings.replicaDiskSoftAntiAffinity | Setting that allows scheduling on disks with existing healthy replicas of the same volume. This setting is enabled by default. |
| defaultSettings.replicaFileSyncHttpClientTimeout | Number of seconds that an HTTP client waits for a response from a File Sync server before considering the connection to have failed. |
| defaultSettings.replicaReplenishmentWaitInterval | Number of seconds that Longhorn waits before reusing existing data on a failed replica instead of creating a new replica of a degraded volume. |
| defaultSettings.replicaSoftAntiAffinity | Setting that allows scheduling on nodes with healthy replicas of the same volume. This setting is disabled by default. |
| defaultSettings.replicaZoneSoftAntiAffinity | Setting that allows Longhorn to schedule new replicas of a volume to nodes in the same zone as existing healthy replicas. Nodes that do not belong to any zone are treated as existing in the zone that contains healthy replicas. When identifying zones, Longhorn relies on the label "topology.kubernetes.io/zone=<Zone name of the node>" in the Kubernetes node object. |
| defaultSettings.restoreConcurrentLimit | Maximum number of worker threads that can concurrently run for each restore operation. |
| defaultSettings.restoreVolumeRecurringJobs | Setting that restores recurring jobs from a backup volume on a backup target and creates recurring jobs if none exist during backup restoration. |
| defaultSettings.snapshotDataIntegrity | Setting that allows you to enable and disable snapshot hashing and data integrity checks. |
| defaultSettings.snapshotDataIntegrityCronjob | Setting that defines when Longhorn checks the integrity of data in snapshot disk files. You must use the Unix cron expression format. |
| defaultSettings.snapshotDataIntegrityImmediateCheckAfterSnapshotCreation | Setting that allows disabling of snapshot hashing after snapshot creation to minimize impact on system performance. |
| defaultSettings.snapshotMaxCount | Maximum snapshot count for a volume. The value should be between 2 to 250 |
| defaultSettings.storageMinimalAvailablePercentage | Percentage of minimum available disk capacity. When the minimum available capacity exceeds the total available capacity, the disk becomes unschedulable until more space is made available for use. The default value is "25". |
| defaultSettings.storageNetwork | Storage network for in-cluster traffic. When unspecified, Longhorn uses the Kubernetes cluster network. |
| defaultSettings.storageOverProvisioningPercentage | Percentage of storage that can be allocated relative to hard drive capacity. The default value is "100". |
| defaultSettings.storageReservedPercentageForDefaultDisk | Percentage of disk space that is not allocated to the default disk on each new Longhorn node. |
| defaultSettings.supportBundleFailedHistoryLimit | Maximum number of failed support bundles that can exist in the cluster. When the value is "0", Longhorn automatically purges all failed support bundles. |
| defaultSettings.systemManagedComponentsNodeSelector | Node selector for system-managed Longhorn components. |
| defaultSettings.systemManagedPodsImagePullPolicy | Image pull policy for system-managed pods, such as Instance Manager, engine images, and CSI Driver. Changes to the image pull policy are applied only after the system-managed pods restart. |
| defaultSettings.taintToleration | Taint or toleration for system-managed Longhorn components. Specify values using a semicolon-separated list in `kubectl taint` syntax (Example: key1=value1:effect; key2=value2:effect). |
| defaultSettings.upgradeChecker | Upgrade Checker that periodically checks for new Longhorn versions. When a new version is available, a notification appears on the Longhorn UI. This setting is enabled by default |
| defaultSettings.v1DataEngine | Setting that allows you to enable the V1 Data Engine. |
| defaultSettings.v2DataEngine | Setting that allows you to enable the V2 Data Engine, which is based on the Storage Performance Development Kit (SPDK). The V2 Data Engine is an experimental feature and should not be used in production environments. |
| defaultSettings.v2DataEngineGuaranteedInstanceManagerCPU | Number of millicpus on each node to be reserved for each Instance Manager pod when the V2 Data Engine is enabled. The default value is "1250". |
| defaultSettings.v2DataEngineHugepageLimit | Setting that allows you to configure maximum huge page size (in MiB) for the V2 Data Engine. |
| defaultSettings.v2DataEngineLogFlags | Setting that allows you to configure the log flags of the SPDK target daemon (spdk_tgt) of the V2 Data Engine. |
| defaultSettings.v2DataEngineLogLevel | Setting that allows you to configure the log level of the SPDK target daemon (spdk_tgt) of the V2 Data Engine. |
| defaultSettings.autoCleanupSnapshotAfterOnDemandBackupCompleted | Setting that automatically cleans up the snapshot after the on-demand backup is completed. |
| defaultSettings.autoCleanupSnapshotWhenDeleteBackup | Setting that automatically cleans up the snapshot when the backup is deleted. |


---

## Article: references/longhorn-client-python.md

---
title: Python Client
weight: 2
---

Currently, you can operate Longhorn using Longhorn UI.
We are planning to build a dedicated Longhorn CLI in the upcoming releases.

In the meantime, you can access Longhorn API using Python binding, as we demonstrated below.

1. Get Longhorn endpoint

   One way to communicate with Longhorn is through `longhorn-frontend` service.

   If you run your automation/scripting tool inside the same cluster in which Longhorn is installed, connect to the endpoint `http://longhorn-frontend.longhorn-system/v1`


   If you run your automation/scripting tool on your local machine,
   use `kubectl port-forward` to forward the `longhorn-frontend` service to localhost:
   ```
   kubectl port-forward services/longhorn-frontend 8080:http -n longhorn-system
   ```
   and connect to endpoint `http://localhost:8080/v1`

2. Using Python Client

    Import file [longhorn.py](https://github.com/longhorn/longhorn-tests/blob/master/manager/integration/tests/longhorn.py) which contains the Python client into your Python script and create a client from the endpoint:
    ```python
    import longhorn

    # If automation/scripting tool is inside the same cluster in which Longhorn is installed
    longhorn_url = 'http://longhorn-frontend.longhorn-system/v1'
    # If forwarding `longhorn-frontend` service to localhost
    longhorn_url = 'http://localhost:8080/v1'

    client = longhorn.Client(url=longhorn_url)

    # Volume operations
    # List all volumes
    volumes = client.list_volume()
    # Get volume by NAME/ID
    testvol1 = client.by_id_volume(id="testvol1")
    # Attach TESTVOL1
    testvol1 = testvol1.attach(hostId="worker-1")
    # Detach TESTVOL1
    testvol1.detach()
    # Create a snapshot of TESTVOL1 with NAME
    snapshot1 = testvol1.snapshotCreate(name="snapshot1")
    # Create a backup from a snapshot NAME
    testvol1.snapshotBackup(name=snapshot1.name)
    # Update the number of replicas of TESTVOL1
    testvol1.updateReplicaCount(replicaCount=2)
    # Find more examples in Longhorn integration tests https://github.com/longhorn/longhorn-tests/tree/master/manager/integration/tests

    # Node operations
    # List all nodes
    nodes = client.list_node()
    # Get node by NAME/ID
    node1 = client.by_id_node(id="worker-1")
    # Disable scheduling for NODE1
    client.update(node1, allowScheduling=False)
    # Enable scheduling for NODE1
    client.update(node1, allowScheduling=True)
    # Find more examples in Longhorn integration tests https://github.com/longhorn/longhorn-tests/tree/master/manager/integration/tests

    # Setting operations
    # List all settings
    settings = client.list_setting()
    # Get setting by NAME/ID
    backupTargetsetting = client.by_id_setting(id="backup-target")
    # Update a setting
    backupTargetsetting = client.update(backupTargetsetting, value="s3://backupbucket@us-east-1/")
    # Find more examples in Longhorn integration tests https://github.com/longhorn/longhorn-tests/tree/master/manager/integration/tests
    ```



---

## Article: references/networking.md

---
title: Longhorn Networking
weight: 3
---

### Overview

This page documents the networking communication between components in the Longhorn system. Using this information, users can write Kubernetes NetworkPolicy
to control the inbound/outbound traffic to/from Longhorn components. This helps to reduce the damage when a malicious pod breaks into the in-cluster network.

The helm chart will install NetworkPolicy objects when the [networkPolicies.enabled value](https://github.com/longhorn/longhorn/blob/v{{< current-version >}}/chart/values.yaml) is set to `true`.
The manifests of these objects can be viewed in the [git repository](https://github.com/longhorn/longhorn/tree/v{{< current-version >}}/chart/templates/network-policies).
Note that depending on the deployed [CNI](https://kubernetes.io/docs/concepts/extend-kubernetes/compute-storage-net/network-plugins/), not all Kubernetes clusters support NetworkPolicy.
See the [Kubernetes documentation](https://kubernetes.io/docs/concepts/services-networking/network-policies/) for details.

> Note: If you are writing network policies, please revisit this page before upgrading Longhorn to make the necessary adjustments to your network policies.
> Note: Depending on your CNI for cluster network, there might be some delay when Kubernetes applying netowk policies to the pod. This delay may fail Longhorn recurring job for taking Snapshot or Backup of the Volume since it cannot access longhorn-manager in the beginning. This is a known issue found in K3s with Traefik and is beyond Longhorn control.

### Longhorn Manager
#### Ingress:
From | Port | Protocol
--- | --- | ---
`Other Longhorn Manager` | 9500 | TCP
`UI` | 9500 | TCP
`Longhorn CSI plugin` | 9500 | TCP
`Backup/Snapshot Recurring Job Pod` | 9500 | TCP
`Longhorn Driver Deployer` | 9500 | TCP
`Conversion Webhook Server` | 9501 | TCP
`Admission Webhook Server` | 9502 | TCP
`Recovery Backend Server` | 9503 | TCP

#### Egress:
To | Port | Protocol
--- | --- | ---
`Other Longhorn Manager` | 9500 | TCP
`Instance Manager` | 8500 (process-manager service); 8501 (proxy service); 8502 (disk service); 8503 (instance service); 8504 (spdk service) | TCP
`Backing Image Manager` | 8000 | TCP
`Backing Image Data Source` | 8000 | TCP
`External Backupstore` | User defined | TCP
`Kubernetes API server` | `Kubernetes API server port` | TCP

### UI
#### ingress:
Users defined
#### egress:
To | Port | Protocol
--- | --- | ---
`Longhorn Manager` | 9500 | TCP

### Instance Manager
#### ingress
From | Port | Protocol
--- | --- | ---
`Longhorn Manager` | 8500 (process-manager service); 8501 (proxy service); 8502 (disk service); 8503 (instance service); 8504 (spdk service) | TCP
`Other Instance Manager` | 10000-30000 | TCP
`Node in the Cluster` | 3260 | TCP
`Backing Image Data Source` | 10000-30000 | TCP

#### egress:
To | Port | Protocol
--- | --- | ---
`Other Instance Manager` | 10000-30000 | TCP
`Backing Image Data Source` |  8002 | TCP
`External Backupstore` | User defined | TCP

### Longhorn CSI plugin
#### ingress
None

#### egress:
To | Port | Protocol
--- | --- | ---
`Longhorn Manager` | 9500 | TCP

#### Additional Info
`Longhorn CSI plugin` pods communitate with `CSI sidecar` pods over the Unix Domain Socket at `<Kuberlet-Directory>/plugins/driver.longhorn.io/csi.sock`


### CSI sidecar (csi-attacher, csi-provisioner, csi-resizer, csi-snapshotter)
#### ingress:
None
#### egress:
To | Port | Protocol
--- | --- | ---
`Kubernetes API server` | `Kubernetes API server port` | TCP

#### Additional Info
`CSI sidecar` pods communitate with `Longhorn CSI plugin` pods over the Unix Domain Socket at `<Kuberlet-Directory>/plugins/driver.longhorn.io/csi.sock`

### Driver deployer
#### ingress:
None
#### egress:
To | Port | Protocol
--- | --- | ---
`Longhorn Manager` | 9500 | TCP
`Kubernetes API server` | `Kubernetes API server port` | TCP

### Engine Image
#### ingress:
None
#### egress:
None

### Backing Image Manager
#### ingress:
From | Port | Protocol
--- | --- | ---
`Longhorn Manager` | 8000 | TCP
`Other Backing Image Manager` | 30001-31000 | TCP

#### egress:
To | Port | Protocol
--- | --- | ---
`Instance Manager` | 10000-30000 | TCP
`Other Backing Image Manager` | 30001-31000 | TCP
`Backing Image Data Source` | 8000 | TCP

### Backing Image Data Source
#### ingress:
From | Port | Protocol
--- | --- | ---
`Longhorn Manager` | 8000 | TCP
`Instance Manager` | 8002 | TCP
`Backing Image Manager` | 8000 | TCP

#### egress:
To | Port | Protocol
--- | --- | ---
`Instance Manager` | 10000-30000 | TCP
`User provided server IP to download the images from` | user defined | TCP

### Share Manager
#### ingress
From | Port | Protocol
--- | --- | ---
`Node in the cluster` | 2049  | TCP

#### egress:
None

### Backup/Snapshot Recurring Job Pod
#### ingress:
None
#### egress:
To | Port | Protocol
--- | --- | ---
`Longhorn Manager` | 9500  | TCP

### Uninstaller
#### ingress:
None
#### egress:
To | Port | Protocol
--- | --- | ---
`Kubernetes API server` | `Kubernetes API server port` | TCP

### Discover Proc Kubelet Cmdline
#### ingress:
None
#### egress:
None

---
Original GitHub issue:
https://github.com/longhorn/longhorn/issues/1805


---

## Article: references/reference-setup-performance-scalability-and-sizing-guidelines.md

---
title: Reference Setup, Performance, Scalability, and Sizing Guidelines
weight: 2
---

You can find the detailed report in the [Longhorn repository](https://github.com/longhorn/longhorn/tree/v{{< current-version >}}/scalability/reference-setup-performance-scalability-and-sizing-guidelines).


---

## Article: references/settings.md

---
title: Settings
weight: 1
---

- [Value Format Types by Supported Data Engines](#value-format-types-by-supported-data-engines)
- [Customizing Default Settings](#customizing-default-settings)
- [System Info](#system-info)
  - [Default Engine Image](#default-engine-image)
  - [Default Instance Manager Image](#default-instance-manager-image)
  - [Default Backing Image Manager Image](#default-backing-image-manager-image)
  - [Support Bundle Manager Image](#support-bundle-manager-image)
  - [Current Longhorn Version](#current-longhorn-version)
  - [Latest Longhorn Version](#latest-longhorn-version)
  - [Stable Longhorn Versions](#stable-longhorn-versions)
- [General](#general)
  - [Node Drain Policy](#node-drain-policy)
  - [Detach Manually Attached Volumes When Cordoned](#detach-manually-attached-volumes-when-cordoned)
  - [Automatically Clean up System Generated Snapshot](#automatically-clean-up-system-generated-snapshot)
  - [Automatically Clean up Outdated Snapshots of Recurring Backup Jobs](#automatically-clean-up-outdated-snapshots-of-recurring-backup-jobs)
  - [Automatically Delete Workload Pod when The Volume Is Detached Unexpectedly](#automatically-delete-workload-pod-when-the-volume-is-detached-unexpectedly)
  - [Blacklist for Automatic Workload Pod Deletion on Unexpected Volume Detachment](#blacklist-for-automatic-workload-pod-deletion-on-unexpected-volume-detachment)
  - [Automatic Salvage](#automatic-salvage)
  - [Concurrent Automatic Engine Upgrade Per Node Limit](#concurrent-automatic-engine-upgrade-per-node-limit)
  - [Concurrent Volume Backup Restore Per Node Limit](#concurrent-volume-backup-restore-per-node-limit)
  - [Create Default Disk on Labeled Nodes](#create-default-disk-on-labeled-nodes)
  - [Custom Resource API Version](#custom-resource-api-version)
  - [Default Data Locality](#default-data-locality)
  - [Default Data Path](#default-data-path)
  - [Default Longhorn Static StorageClass Name](#default-longhorn-static-storageclass-name)
  - [Default Replica Count](#default-replica-count)
  - [Deleting Confirmation Flag](#deleting-confirmation-flag)
  - [Disable Revision Counter](#disable-revision-counter)
  - [Enable Upgrade Checker](#enable-upgrade-checker)
  - [Upgrade Responder URL](#upgrade-responder-url)
  - [Allow Collecting Longhorn Usage Metrics](#allow-collecting-longhorn-usage-metrics)
  - [Pod Deletion Policy When Node is Down](#pod-deletion-policy-when-node-is-down)
  - [Registry Secret](#registry-secret)
  - [Replica Replenishment Wait Interval](#replica-replenishment-wait-interval)
  - [System Managed Pod Image Pull Policy](#system-managed-pod-image-pull-policy)
  - [Backing Image Cleanup Wait Interval](#backing-image-cleanup-wait-interval)
  - [Backing Image Recovery Wait Interval](#backing-image-recovery-wait-interval)
  - [Default Min Number Of Backing Image Copies](#default-min-number-of-backing-image-copies)
  - [Engine Replica Timeout](#engine-replica-timeout)
  - [Support Bundle Failed History Limit](#support-bundle-failed-history-limit)
  - [Support Bundle Node Collection Timeout](#support-bundle-node-collection-timeout)
  - [Fast Replica Rebuild Enabled](#fast-replica-rebuild-enabled)
  - [Timeout of HTTP Client to Replica File Sync Server](#timeout-of-http-client-to-replica-file-sync-server)
  - [Long gRPC Timeout](#long-grpc-timeout)
  - [Offline Replica Rebuilding](#offline-replica-rebuilding)
  - [RWX Volume Fast Failover (Experimental)](#rwx-volume-fast-failover-experimental)
  - [Log Level](#log-level)
  - [Data Engine Log Level](#data-engine-log-level)
  - [Data Engine Log Flags](#data-engine-log-flags)
  - [Replica Rebuilding Bandwidth Limit](#replica-rebuilding-bandwidth-limit)
- [Snapshot](#snapshot)
  - [Snapshot Data Integrity](#snapshot-data-integrity)
  - [Immediate Snapshot Data Integrity Check After Creating a Snapshot](#immediate-snapshot-data-integrity-check-after-creating-a-snapshot)
  - [Snapshot Data Integrity Check CronJob](#snapshot-data-integrity-check-cronjob)
  - [Snapshot Maximum Count](#snapshot-maximum-count)
  - [Freeze Filesystem For Snapshot](#freeze-filesystem-for-snapshot)
- [Orphan](#orphan)
  - [Orphaned Resource Automatic Deletion](#orphaned-resource-automatic-deletion)
  - [Orphaned Resource Automatic Deletion Grace Period](#orphaned-resource-automatic-deletion-grace-period)
- [Backups](#backups)
  - [Allow Recurring Job While Volume Is Detached](#allow-recurring-job-while-volume-is-detached)
  - [Backup Execution Timeout](#backup-execution-timeout)
  - [Failed Backup Time To Live](#failed-backup-time-to-live)
  - [Cronjob Failed Jobs History Limit](#cronjob-failed-jobs-history-limit)
  - [Cronjob Successful Jobs History Limit](#cronjob-successful-jobs-history-limit)
  - [Restore Volume Recurring Jobs](#restore-volume-recurring-jobs)
  - [Backup Compression Method](#backup-compression-method)
  - [Backup Concurrent Limit Per Backup](#backup-concurrent-limit-per-backup)
  - [Restore Concurrent Limit Per Backup](#restore-concurrent-limit-per-backup)
  - [Default Backup Block Size](#default-backup-block-size)
- [Scheduling](#scheduling)
  - [Allow Volume Creation with Degraded Availability](#allow-volume-creation-with-degraded-availability)
  - [Disable Scheduling On Cordoned Node](#disable-scheduling-on-cordoned-node)
  - [Replica Node Level Soft Anti-Affinity](#replica-node-level-soft-anti-affinity)
  - [Replica Zone Level Soft Anti-Affinity](#replica-zone-level-soft-anti-affinity)
  - [Replica Disk Level Soft Anti-Affinity](#replica-disk-level-soft-anti-affinity)
  - [Replica Auto Balance](#replica-auto-balance)
  - [Replica Auto Balance Disk Pressure Threshold (%)](#replica-auto-balance-disk-pressure-threshold-)
  - [Storage Minimal Available Percentage](#storage-minimal-available-percentage)
  - [Storage Over Provisioning Percentage](#storage-over-provisioning-percentage)
  - [Storage Reserved Percentage For Default Disk](#storage-reserved-percentage-for-default-disk)
  - [Allow Empty Node Selector Volume](#allow-empty-node-selector-volume)
  - [Allow Empty Disk Selector Volume](#allow-empty-disk-selector-volume)
- [Danger Zone](#danger-zone)
  - [V1 Data Engine](#v1-data-engine)
  - [V2 Data Engine](#v2-data-engine)
  - [Concurrent Replica Rebuild Per Node Limit](#concurrent-replica-rebuild-per-node-limit)
  - [Concurrent Backing Image Replenish Per Node Limit](#concurrent-backing-image-replenish-per-node-limit)
  - [Kubernetes Taint Toleration](#kubernetes-taint-toleration)
  - [Priority Class](#priority-class)
  - [System Managed Components Node Selector](#system-managed-components-node-selector)
  - [Kubernetes Cluster Autoscaler Enabled (Experimental)](#kubernetes-cluster-autoscaler-enabled-experimental)
  - [Storage Network](#storage-network)
  - [Storage Network For RWX Volume Enabled](#storage-network-for-rwx-volume-enabled)
  - [Remove Snapshots During Filesystem Trim](#remove-snapshots-during-filesystem-trim)
  - [Guaranteed Instance Manager CPU](#guaranteed-instance-manager-cpu)
  - [Disable Snapshot Purge](#disable-snapshot-purge)
  - [Auto Cleanup Snapshot When Delete Backup](#auto-cleanup-snapshot-when-delete-backup)
  - [Auto Cleanup Snapshot After On-Demand Backup Completed](#auto-cleanup-snapshot-after-on-demand-backup-completed)
  - [Instance Manager Pod Liveness Probe Timeout](#instance-manager-pod-liveness-probe-timeout)
  - [Data Engine CPU Mask](#data-engine-cpu-mask)
  - [Data Engine Hugepage Enabled](#data-engine-hugepage-enabled)
  - [Data Engine Memory Size](#data-engine-memory-size)
  - [Data Engine Interrupt Mode Enabled](#data-engine-interrupt-mode-enabled)
  - [Log Path](#log-path)

---

### Value Format Types by Supported Data Engines

Each setting supports only one of the following formats, depending on its definition. The supported format determines which Data Engines can be configured and whether their values can differ.

- Single value for all supported Data Engines
  - Format: Non-JSON string (e.g., `1024`)
  - The value applies to all supported Data Engines and must be the same across them.
  - Data-engine-specific values are not allowed.
- Data-engine-specific values for V1 and V2 Data Engines
  - Format: JSON object (e.g., `{"v1": "value1", "v2": "value2"}`)
  - Allows specifying different values for V1 and V2 Data Engines.
- Data-engine-specific values for V1 Data Engine only
  - Format: JSON object with `v1` key only (e.g., `{"v1": "value1"}`)
  - Only the V1 Data Engine can be configured; the V2 Data Engine is not affected.
- Data-engine-specific values for V2 Data Engine only
  - Format: JSON object with `v2` key only (e.g., `{"v2": "value1"}`)
  - Only the V2 Data Engine can be configured; the V1 Data Engine is not affected.

### Customizing Default Settings

To configure Longhorn before installing it, see [this section](../../advanced-resources/deploy/customizing-default-settings) for details.

### System Info

#### Default Engine Image

The default engine image used by the manager. Can be changed on the manager starting command line only.

Every Longhorn release will ship with a new Longhorn engine image. If the current Longhorn volumes are not using the default engine, a green arrow will show up, indicate this volume needs to be upgraded to use the default engine.

#### Default Instance Manager Image

The default instance manager image used by the manager. Can be changed on the manager starting command line only.

#### Default Backing Image Manager Image

The default backing image manager image used by the manager. Can be changed on the manager starting command line only.

#### Support Bundle Manager Image

Longhorn uses the support bundle manager image to generate the support bundles.

There will be a default image given during installation and upgrade. You can also change it in the settings.

An example of the support bundle manager image:
> Default: `longhornio/support-bundle-kit:v0.0.14`

#### Current Longhorn Version

The current Longhorn version.

#### Latest Longhorn Version

The latest version of Longhorn available. Automatically updated by the Upgrade Checker.

> Only available if `Upgrade Checker` is enabled.

#### Stable Longhorn Versions

The latest stable version of every minor release line. Updated by Upgrade Checker automatically.

### General

#### Node Drain Policy

> Default: `block-if-contains-last-replica`

Define the policy to use when a node with the last healthy replica of a volume is drained. Available options:

- `block-if-contains-last-replica`: Longhorn will block the drain when the node contains the last healthy replica of a
  volume.
- `allow-if-replica-is-stopped`: Longhorn will allow the drain when the node contains the last healthy replica of a
  volume but the replica is stopped.
  WARNING: possible data loss if the node is removed after draining.
- `always-allow`: Longhorn will allow the drain even though the node contains the last healthy replica of a volume.
  WARNING: possible data loss if the node is removed after draining. Also possible data corruption if the last replica
  was running during the draining.
- `block-for-eviction`: Longhorn will automatically evict all replicas and block the drain until eviction is complete.
  WARNING: Can result in slow drains and extra data movement associated with replica rebuilding.
- `block-for-eviction-if-contains-last-replica`: Longhorn will automatically evict any replicas that don't have a
  healthy counterpart and block the drain until eviction is complete.
  WARNING: Can result in slow drains and extra data movement associated with replica rebuilding.

Each option has benefits and drawbacks. See [Node Drain Policy
Recommendations](../../maintenance/maintenance/#node-drain-policy-recommendations) for help deciding which is most
appropriate in your environment.

#### Detach Manually Attached Volumes When Cordoned

> Default: `false`

Longhorn will automatically detach volumes that are manually attached to the nodes which are cordoned.
This prevent the draining process stuck by the PDB of instance-manager which still has running engine on the node.

#### Automatically Clean up System Generated Snapshot

> Default: `true`

Longhorn will generate system snapshot during replica rebuild, and if a user doesn't setup a recurring snapshot schedule, all the system generated snapshots would be left in the replica, and user has to delete them manually, this setting allow Longhorn to automatically cleanup system generated snapshot before and after replica rebuild.

#### Automatically Clean up Outdated Snapshots of Recurring Backup Jobs

> Default: `true`

If enabled, when running a recurring backup job, Longhorn takes a new snapshot before creating the backup. Longhorn retains only the snapshot used by the last backup job even if the value of the retain parameter is not 1.

If disabled, this setting ensures that the retained snapshots directly correspond to the backups on the remote backup target.

#### Automatically Delete Workload Pod when The Volume Is Detached Unexpectedly

> Default: `true`

If enabled, Longhorn will automatically delete the workload pod that is managed by a controller (e.g. deployment, statefulset, daemonset, etc...) when Longhorn volume is detached unexpectedly (e.g. during Kubernetes upgrade, Docker reboot, or network disconnect).
By deleting the pod, its controller restarts the pod and Kubernetes handles volume reattachment and remount.

If disabled, Longhorn will not delete the workload pod that is managed by a controller. You will have to manually restart the pod to reattach and remount the volume.

> **Note:** This setting doesn't apply to below cases.
> - The workload pods don't have a controller; Longhorn never deletes them.
> - Workload pods with *cluster network* RWX volumes. The setting does not apply to such pods because the Longhorn Share Manager, which provides the RWX NFS service, has its own resilience mechanism. This mechanism ensures availability until the volume is reattached without relying on the pod lifecycle to trigger volume reattachment. The setting does apply, however, to workload pods with *storage network* RWX volumes. For more information, see [ReadWriteMany (RWX) Volume](../../nodes-and-volumes/volumes/rwx-volumes) and [Storage Network](../../advanced-resources/deploy/storage-network#limitation).

#### Blacklist for Automatic Workload Pod Deletion on Unexpected Volume Detachment

> Default: `""`

Blacklist of controller api/kind values for the setting [Automatically Delete Workload Pod when the Volume Is Detached Unexpectedly](#automatically-delete-workload-pod-when-the-volume-is-detached-unexpectedly). If a workload pod is managed by a controller whose api/kind is listed in this blacklist, Longhorn will not automatically delete the pod when its volume is unexpectedly detached. Multiple controller api/kind entries can be specified, separated by semicolons. For example: `apps/v1/StatefulSet;apps/v1/DaemonSet`.

> **Note:** The controller api/kind is case sensitive and must exactly match the api/kind in the workload pod's owner reference.

#### Automatic Salvage

> Default: `true`

If enabled, volumes will be automatically salvaged when all the replicas become faulty e.g. due to network disconnection. Longhorn will try to figure out which replica(s) are usable, then use them for the volume.

#### Concurrent Automatic Engine Upgrade Per Node Limit

> Default: `0`

This setting controls how Longhorn automatically upgrades volumes' engines to the new default engine image after upgrading Longhorn manager.
The value of this setting specifies the maximum number of engines per node that are allowed to upgrade to the default engine image at the same time.
If the value is 0, Longhorn will not automatically upgrade volumes' engines to default version.

#### Concurrent Volume Backup Restore Per Node Limit

> Default: `5`

This setting controls how many volumes on a node can restore the backup concurrently.

Longhorn blocks the backup restore once the restoring volume count exceeds the limit.

Set the value to **0** to disable backup restore.

#### Create Default Disk on Labeled Nodes

> Default: `false`

If no other disks exist, create the default disk automatically, only on nodes with the Kubernetes label `node.longhorn.io/create-default-disk=true` .

If disabled, the default disk will be created on all new nodes when the node is detected for the first time.

This option is useful if you want to scale the cluster but don't want to use storage on the new nodes, or if you want to [customize disks for Longhorn nodes](../../nodes-and-volumes/nodes/default-disk-and-node-config).

#### Custom Resource API Version

> Default: `longhorn.io/v1beta2`

The current customer resource's API version, e.g. longhorn.io/v1beta2. Set by manager automatically.

#### Default Data Locality

> Default: `disabled`

We say a Longhorn volume has data locality if there is a local replica of the volume on the same node as the pod which is using the volume.
This setting specifies the default data locality when a volume is created from the Longhorn UI. For Kubernetes configuration, update the dataLocality in the StorageClass

The available modes are:

- `disabled`. This is the default option.
  There may or may not be a replica on the same node as the attached volume (workload).

- `best-effort`. This option instructs Longhorn to try to keep a replica on the same node as the attached volume (workload).
  Longhorn will not stop the volume, even if it cannot keep a replica local to the attached volume (workload) due to environment limitation, e.g. not enough disk space, incompatible disk tags, etc.

- `strict-local`: This option enforces Longhorn keep the **only one replica** on the same node as the attached volume, and therefore, it offers higher IOPS and lower latency performance.

#### Default Data Path

> Default: `/var/lib/longhorn/`

Default path to use for storing data on a host.

Can be used with `Create Default Disk on Labeled Nodes` option, to make Longhorn only use the nodes with specific storage mounted at, for example, `/opt/longhorn` when scaling the cluster.

#### Default Longhorn Static StorageClass Name

> Default: `longhorn-static`

The `storageClassName` is for persistent volumes (PVs) and persistent volume claims (PVCs) when creating PV/PVC for an existing Longhorn volume. Notice that it's unnecessary for users to create the related StorageClass object in Kubernetes since the StorageClass would only be used as matching labels for PVC bounding purposes.  The "storageClassName" needs to be an existing StorageClass. Only the StorageClass named `longhorn-static` will be created if it does not exist. By default 'longhorn-static'.

#### Default Replica Count

> Default: `{"v1":"3","v2":"3"}`

The default number of replicas when creating the volume from Longhorn UI. For Kubernetes, update the `numberOfReplicas` in the StorageClass

The recommended way of choosing the default replica count is: if you have three or more nodes for storage, use 3; otherwise use 2. Using a single replica on a single node cluster is also OK, but the high availability functionality wouldn't be available. You can still take snapshots/backups of the volume.

#### Deleting Confirmation Flag

> Default: `false`

This flag is designed to prevent Longhorn from being accidentally uninstalled which will lead to data loss.

- Set this flag to **true** to allow Longhorn uninstallation.
- If this flag **false**, Longhorn uninstallation job will fail.

#### Disable Revision Counter

> Default: `{"v1":"true"}`

Allows engine controller and engine replica to disable revision counter file update for every data write. This improves the data path performance. See [Revision Counter](../../advanced-resources/deploy/revision_counter) for details.

#### Enable Upgrade Checker

> Default: `true`

Upgrade Checker will check for a new Longhorn version periodically. When there is a new version available, it will notify the user in the Longhorn UI.

#### Upgrade Responder URL

> Default: `https://longhorn-upgrade-responder.rancher.io/v1/checkupgrade`

The Upgrade Responder sends a notification whenever a new Longhorn version that you can upgrade to becomes available.

#### Allow Collecting Longhorn Usage Metrics

> Default: `true`

Enabling this setting will allow Longhorn to provide valuable usage metrics to https://metrics.longhorn.io/.

This information will help us gain insights how Longhorn is being used, which will ultimately contribute to future improvements.

**Node Information collected from all cluster nodes includes:**

- Number of disks of each device type (HDD, SSD, NVMe, unknown).
  > This value may not be accurate for virtual machines.
- Number of disks for each Longhorn disk type (block, filesystem).
- Host system architecture.
- Host kernel release.
- Host operating system (OS) distribution.
- Kubernetes node provider.

**Cluster Information collected from one of the cluster nodes includes:**

- Longhorn namespace UID.
- Number of Longhorn nodes.
- Number of volumes of each access mode (RWO, RWX, unknown).
- Number of volumes of each data engine (v1, v2).
- Number of volumes of each data locality type (disabled, best_effort, strict_local, unknown).
- Number of volumes that are encrypted or unencrypted.
- Number of volumes of each frontend type (blockdev, iscsi).
- Number of replicas.
- Number of snapshots.
- Number of backing images.
- Number of orphans.
- Average volume size in bytes.
- Average volume actual size in bytes.
- Average number of snapshots per volume.
- Average number of replicas per volume.
- Average Longhorn component CPU usage (instance manager, manager) in millicores.
- Average Longhorn component memory usage (instance manager, manager) in bytes.
- Longhorn settings:
  - Partially included:
    - Backup Target Type/Protocol (azblob, cifs, nfs, s3, none, unknown). This is from the Backup Target setting.
  - Included as true or false to indicate if this setting is configured:
    - Priority Class
    - Registry Secret
    - Snapshot Data Integrity CronJob
    - Storage Network
    - System Managed Components Node Selector
    - Taint Toleration
  - Included as it is:
    - Allow Recurring Job While Volume Is Detached
    - Allow Volume Creation With Degraded Availability
    - Automatically Clean up System Generated Snapshot
    - Automatically Clean up Outdated Snapshots of Recurring Backup Jobs
    - Automatically Delete Workload Pod when The Volume Is Detached Unexpectedly
    - Automatic Salvage
    - Backing Image Cleanup Wait Interval
    - Backing Image Recovery Wait Interval
    - Backup Compression Method
    - Backupstore Poll Interval
    - Backup Concurrent Limit
    - Concurrent Automatic Engine Upgrade Per Node Limit
    - Concurrent Backup Restore Per Node Limit
    - Concurrent Replica Rebuild Per Node Limit
    - CRD API Version
    - Create Default Disk Labeled Nodes
    - Default Data Locality
    - Default Replica Count
    - Disable Revision Counter
    - Disable Scheduling On Cordoned Node
    - Engine Replica Timeout
    - Failed Backup TTL
    - Fast Replica Rebuild Enabled
    - Guaranteed Instance Manager CPU
    - Kubernetes Cluster Autoscaler Enabled
    - Node Down Pod Deletion Policy
    - Node Drain Policy
    - Orphan Auto Deletion
    - Recurring Failed Jobs History Limit
    - Recurring Successful Jobs History Limit
    - Remove Snapshots During Filesystem Trim
    - Replica Auto Balance
    - Replica File Sync HTTP Client Timeout
    - Replica Replenishment Wait Interval
    - Replica Soft Anti Affinity
    - Replica Zone Soft Anti Affinity
    - Replica Disk Soft Anti Affinity
    - Restore Concurrent Limit
    - Restore Volume Recurring Jobs
    - Snapshot Data Integrity
    - Snapshot DataIntegrity Immediate Check After Snapshot Creation
    - Storage Minimal Available Percentage
    - Storage Network For RWX Volume Enabled
    - Storage Over Provisioning Percentage
    - Storage Reserved Percentage For Default Disk
    - Support Bundle Failed History Limit
    - Support Bundle Node Collection Timeout
    - System Managed Pods Image Pull Policy

> The `Upgrade Checker` needs to be enabled to periodically send the collected data.

#### Pod Deletion Policy When Node is Down

> Default: `do-nothing`

Defines the Longhorn action when a Volume is stuck with a StatefulSet/Deployment Pod on a node that is down.

- `do-nothing` is the default Kubernetes behavior of never force deleting StatefulSet/Deployment terminating pods. Since the pod on the node that is down isn't removed, Longhorn volumes are stuck on nodes that are down.
- `delete-statefulset-pod` Longhorn will force delete StatefulSet terminating pods on nodes that are down to release Longhorn volumes so that Kubernetes can spin up replacement pods.
- `delete-deployment-pod` Longhorn will force delete Deployment terminating pods on nodes that are down to release Longhorn volumes so that Kubernetes can spin up replacement pods.
- `delete-both-statefulset-and-deployment-pod` Longhorn will force delete StatefulSet/Deployment terminating pods on nodes that are down to release Longhorn volumes so that Kubernetes can spin up replacement pods.

#### Registry Secret

The Kubernetes Secret name.

#### Replica Replenishment Wait Interval

> Default: `600`

When there is at least one failed replica volume in a degraded volume, this interval in seconds determines how long Longhorn will wait at most in order to reuse the existing data of the failed replicas rather than directly creating a new replica for this volume.

Warning: This wait interval works only when there is at least one failed replica in the volume. And this option may block the rebuilding for a while.

#### System Managed Pod Image Pull Policy

> Default: `if-not-present`

This setting defines the Image Pull Policy of Longhorn system managed pods, e.g. instance manager, engine image, CSI driver, etc.

Notice that the new Image Pull Policy will only apply after the system managed pods restart.

This setting definition is exactly the same as that of in Kubernetes. Here are the available options:

- `always`. Every time the kubelet launches a container, the kubelet queries the container image registry to resolve the name to an image digest. If the kubelet has a container image with that exact digest cached locally, the kubelet uses its cached image; otherwise, the kubelet downloads (pulls) the image with the resolved digest, and uses that image to launch the container.

- `if-not-present`. The image is pulled only if it is not already present locally.

- `never`. The image is assumed to exist locally. No attempt is made to pull the image.


#### Backing Image Cleanup Wait Interval
> Default: `60`

This interval in minutes determines how long Longhorn will wait before cleaning up the backing image file when there is no replica in the disk using it.

#### Backing Image Recovery Wait Interval
> Default: `300`

The interval in seconds determines how long Longhorn will wait before re-downloading the backing image file when all disk files of this backing image become `failed` or `unknown`.

> **Note:**
>  - This recovery only works for the backing image of which the creation type is `download`.
>  - File state `unknown` means the related manager pods on the pod is not running or the node itself is down/disconnected.

#### Default Min Number Of Backing Image Copies
> Default: `1`

The default minimum number of backing image copies Longhorn maintains.

#### Engine Replica Timeout

> Default: `{"v1":"8","v2":"8"}`

The time in seconds a v1 engine will wait for a response from a replica before marking it as failed. Values between 8
and 30 are allowed. The engine replica timeout is only in effect while there are I/O requests outstanding.

This setting only applies to additional replicas. A V1 engine marks the last active replica as failed only after twice
the configured number of seconds (timeout value x 2) have passed. This behavior is intended to balance volume
responsiveness with volume availability.

- The engine can quickly (after the configured timeout) ignore individual replicas that become unresponsive in favor of
  other available ones. This ensures future I/O will not be held up.
- The engine waits on the last replica (until twice the configured timeout) to prevent unnecessarily crashing as a
  result of having no available backends.

#### Support Bundle Failed History Limit

> Default: `1`

This setting specifies how many failed support bundles can exist in the cluster.

The retained failed support bundle is for analysis purposes and needs to clean up manually.

Longhorn blocks support bundle creation when reaching the upper bound of the limitation. You can set this value to **0** to have Longhorn automatically purge all failed support bundles.

#### Support Bundle Node Collection Timeout

> Default: `30`

Number of minutes Longhorn allows for collection of node information and node logs for the support bundle.

If the collection process is not completed within the allotted time, Longhorn continues generating the support bundle without the uncollected node data.

#### Fast Replica Rebuild Enabled

> Default: `{"v1":"true","v2":"true"}`

The setting enables fast replica rebuilding feature. It relies on the checksums of snapshot disk files, so setting the snapshot-data-integrity to **enable** or **fast-check** is a prerequisite.

#### Timeout of HTTP Client to Replica File Sync Server

> Default: `30`

The value in seconds specifies the timeout of the HTTP client to the replica's file sync server used for replica rebuilding, volume cloning, snapshot cloning, etc.

#### Long gRPC Timeout

> Default: `86400`

Number of seconds that Longhorn allows for the completion of replica rebuilding and snapshot cloning operations.

#### Offline Replica Rebuilding

> Default: `false`

Controls whether Longhorn automatically rebuilds degraded replicas while the volume is detached. This setting only takes effect if the volume-level setting is set to `ignored` or `enabled`. Available options:

- **true**: Enables offline replica rebuilding for all detached volumes (unless overridden at the volume level).
- **false**: Disables offline replica rebuilding globally (unless overridden at the volume level).

> **Note:** Offline rebuilding occurs only when a volume is detached. Volumes in a faulted state will not trigger offline rebuilding.

This setting allows Longhorn to automatically rebuild replicas for detached volumes when needed.

#### RWX Volume Fast Failover (Experimental)

> Default: `false`

Enable improved ReadWriteMany volume HA by shortening the time it takes to recover from a node failure.

#### Log Level

> Default: `Log Level`

The log level Panic, Fatal, Error, Warn, Info, Debug, Trace used in longhorn manager. By default Info.


#### Data Engine Log Level

> Default: `{"v2":"Notice"}`

Applies only to the V2 Data Engine. Specifies the log level for the Storage Performance Development Kit (SPDK) target daemon. Supported values: `Error`, `Warning`, `Notice`, `Info`, and `Debug`.

#### Data Engine Log Flags

> Default: `{"v2":""}`

Applies only to the V2 Data Engine. Specifies the log flags for the Storage Performance Development Kit (SPDK) target daemon.

#### Replica Rebuilding Bandwidth Limit

> Default: `{"v2":"0"}`

Applies only to the V2 Data Engine. Specifies the default write bandwidth limit, in megabytes per second (MB/s), for volume replica rebuilding.

### Snapshot

#### Snapshot Data Integrity

> Default: `{"v1":"fast-check","v2":"fast-check"}`

This setting allows users to enable or disable snapshot hashing and data integrity checking. Available options are:

- **disabled**: Disable snapshot disk file hashing and data integrity checking.
- **enabled**: Enables periodic snapshot disk file hashing and data integrity checking. To detect the filesystem-unaware corruption caused by bit rot or other issues in snapshot disk files, Longhorn system periodically hashes files and finds corrupted ones. Hence, the system performance will be impacted during the periodical checking.
- **fast-check**: Enable snapshot disk file hashing and fast data integrity checking. Longhorn system only hashes snapshot disk files if their are not hashed or the modification time are changed. In this mode, filesystem-unaware corruption cannot be detected, but the impact on system performance can be minimized.

#### Immediate Snapshot Data Integrity Check After Creating a Snapshot

> Default: `{"v1":"false","v2":"false"}`

Hashing snapshot disk files impacts the performance of the system. The immediate snapshot hashing and checking can be disabled to minimize the impact after creating a snapshot.

#### Snapshot Data Integrity Check CronJob

> Default: `{"v1":"0 0 */7 * *","v2":"0 0 */7 * *"}`

Unix-cron string format. The setting specifies when Longhorn checks the data integrity of snapshot disk files.
> **Warning**
> Hashing snapshot disk files impacts the performance of the system. It is recommended to run data integrity checks during off-peak times and to reduce the frequency of checks.

#### Snapshot Maximum Count

> Default: `250`

Maximum snapshot count for a volume. The value should be between 2 to 250.

#### Freeze Filesystem For Snapshot

> Default: `{"v1":"false"}`

This setting only applies to volumes with the Kubernetes volume mode `Filesystem`. When enabled, Longhorn freezes the
volume's filesystem immediately before creating a user-initiated snapshot. When disabled or when the Kubernetes volume
mode is `Block`, Longhorn instead attempts a system sync before creating a user-initiated snapshot.

Snapshots created when this setting is enabled are more likely to be consistent because the filesystem is in a
consistent state at the moment of creation. However, under very heavy I/O, freezing the filesystem may take a
significant amount of time and may cause workload activity to pause.

When this setting is disabled, all data is flushed to disk just before the snapshot is created, but Longhorn cannot
completely block write attempts during the brief interval between the system sync and snapshot creation. I/O is not
paused during the system sync, so workloads likely do not notice that a snapshot is being created.

The default option for this setting is `false` because kernels with version `v5.17` or earlier may not respond correctly
when a volume crashes while a freeze is ongoing. This is not likely to happen but if it does, an affected kernel will
not allow you to unmount the filesystem or stop processes using the filesystem without rebooting the node. Only enable
this setting if you plan to use kernels with version `5.17` or later, and ext4 or XFS filesystems.

You can override this setting (using the field `freezeFilesystemForSnapshot`) for specific volumes through the Longhorn
UI, a StorageClass, or direct changes to an existing volume. `freezeFilesystemForSnapshot` accepts the following values:

> Default: `ignored`

- `ignored`: Instructs Longhorn to use the global setting. This is the default option.
- `enabled`: Enables freezing before snapshots, regardless of the global setting.
- `disabled`: Disables freezing before snapshots, regardless of the global setting.

### Orphan

#### Orphaned Resource Automatic Deletion

> Example: `replica-data;instance`

This setting allows Longhorn to automatically delete orphan resources and their corresponding orphaned resources. Orphan resources located on nodes that are in down or unknown state will not be cleaned up automatically.

List the enabled resource types in a semicolon-separated list. Available items are:

- `replica-data`: replica data store
- `instance`: engine and replica runtime instance

#### Orphaned Resource Automatic Deletion Grace Period

> Default: `300`

Specifies the wait time, in seconds, before Longhorn automatically deletes an orphaned Custom Resource (CR) and its associated resources.

> **Note:** If a user manually deletes an orphaned CR, the deletion occurs immediately and does not respect this grace period.

### Backups

#### Allow Recurring Job While Volume Is Detached

> Default: `false`

If this setting is enabled, Longhorn automatically attaches the volume and takes snapshot/backup when it is the time to do recurring snapshot/backup.

> **Note:** During the time the volume was attached automatically, the volume is not ready for the workload. The workload will have to wait until the recurring job finishes.

#### Backup Execution Timeout

> Default: `1`

Number of minutes that Longhorn allows for the backup execution.

#### Failed Backup Time To Live

> Default: `1440`

The interval in minutes to keep the backup resource that was failed. Set to 0 to disable the auto-deletion.

Failed backups will be checked and cleaned up during backupstore polling which is controlled by **Backupstore Poll Interval** setting. Hence this value determines the minimal wait interval of the cleanup. And the actual cleanup interval is multiple of **Backupstore Poll Interval**. Disabling **Backupstore Poll Interval** also means to disable failed backup auto-deletion.

#### Cronjob Failed Jobs History Limit

> Default: `1`

This setting specifies how many failed backup or snapshot job histories should be retained.

History will not be retained if the value is 0.

#### Cronjob Successful Jobs History Limit

> Default: `1`

This setting specifies how many successful backup or snapshot job histories should be retained.

History will not be retained if the value is 0.

#### Restore Volume Recurring Jobs

> Default: `false`

This setting allows restoring the recurring jobs of a backup volume from the backup target during a volume restoration if they do not exist on the cluster.
This is also a volume-specific setting with the below options. Users can customize it for each volume to override the global setting.

> Default: `ignored`

- `ignored`: This is the default option that instructs Longhorn to inherit from the global setting.

- `enabled`: This option instructs Longhorn to restore volume recurring jobs/groups from the backup target forcibly.

- `disabled`: This option instructs Longhorn no restoring volume recurring jobs/groups should be done.

#### Backup Compression Method

> Default: `lz4`

This setting allows users to specify backup compression method.

- `none`: Disable the compression method. Suitable for multimedia data such as encoded images and videos.

- `lz4`: Fast compression method. Suitable for flat files.

- `gzip`: A bit of higher compression ratio but relatively slow.

#### Backup Concurrent Limit Per Backup

> Default: `2`

This setting controls how many worker threads per backup concurrently.

#### Restore Concurrent Limit Per Backup

> Default: `2`

This setting controls how many worker threads per restore concurrently.

#### Default Backup Block Size

> Default: `2`

Specifies the default backup block size, in MiB, used when creating a new volume. Supported values are `2` or `16`.

### Scheduling

#### Allow Volume Creation with Degraded Availability

> Default: `true`

This setting allows user to create and attach a volume that doesn't have all the replicas scheduled at the time of creation.

> **Note:** It's recommended to disable this setting when using Longhorn in the production environment. See [Best Practices](../../best-practices/) for details.

#### Disable Scheduling On Cordoned Node

> Default: `true`

When this setting is checked, the Longhorn Manager will not schedule replicas on Kubernetes cordoned nodes.

When this setting is un-checked, the Longhorn Manager will schedule replicas on Kubernetes cordoned nodes.

#### Replica Node Level Soft Anti-Affinity

> Default: `false`

When this setting is checked, the Longhorn Manager will allow scheduling on nodes with existing healthy replicas of the same volume.

When this setting is un-checked, Longhorn Manager will forbid scheduling on nodes with existing healthy replicas of the same volume.

> **Note:**
>   - This setting is superseded if replicas are forbidden to share a zone by the Replica Zone Level Anti-Affinity setting.

#### Replica Zone Level Soft Anti-Affinity

> Default: `true`

When this setting is checked, the Longhorn Manager will allow scheduling new replicas of a volume to the nodes in the same zone as existing healthy replicas.

When this setting is un-checked, Longhorn Manager will forbid scheduling new replicas of a volume to the nodes in the same zone as existing healthy replicas.

> **Note:**
>   - Nodes that don't belong to any zone will be treated as if they belong to the same zone.
>   - Longhorn relies on label `topology.kubernetes.io/zone=<Zone name of the node>` in the Kubernetes node object to identify the zone.

#### Replica Disk Level Soft Anti-Affinity

> Default: `true`

When this setting is checked, the Longhorn Manager will allow scheduling new replicas of a volume to the same disks as existing healthy replicas.

When this setting is un-checked, Longhorn Manager will forbid scheduling new replicas of a volume to the same disks as existing healthy replicas.

> **Note:**
>   - Even if the setting is "true" and disk sharing is allowed, Longhorn will seek to use a different disk if possible, even if on the same node.
>   - This setting is superseded if replicas are forbidden to share a zone or a node by either of the other Soft Anti-Affinity settings.

#### Replica Auto Balance

> Default: `disabled`

Enable this setting automatically rebalances replicas when discovered an available node.

The available global options are:

- `disabled`. This is the default option. No replica auto-balance will be done.

- `least-effort`. This option instructs Longhorn to balance replicas for minimal redundancy.

- `best-effort`. This option instructs Longhorn try to balancing replicas for even redundancy.
  Longhorn does not forcefully re-schedule the replicas to a zone that does not have enough nodes
  to support even balance. Instead, Longhorn will re-schedule to balance at the node level.

Longhorn also supports customizing for individual volume. The setting can be specified in UI or with Kubernetes manifest volume.spec.replicaAutoBalance, this overrules the global setting.
The available volume spec options are:

> Default: `ignored`

- `ignored`. This is the default option that instructs Longhorn to inherit from the global setting.

- `disabled`. This option instructs Longhorn no replica auto-balance should be done."

- `least-effort`. This option instructs Longhorn to balance replicas for minimal redundancy.

- `best-effort`. This option instructs Longhorn to try balancing replicas for even redundancy.
  Longhorn does not forcefully re-schedule the replicas to a zone that does not have enough nodes
  to support even balance. Instead, Longhorn will re-schedule to balance at the node level.

#### Replica Auto Balance Disk Pressure Threshold (%)

> Default: `90`

Percentage of currently used storage that triggers automatic replica rebalancing.

When the threshold is reached, Longhorn automatically rebuilds replicas that are under disk pressure on another disk within the same node.

To disable this setting, set the value to **0**.

This setting takes effect only when the following conditions are met:

- [Replica Auto Balance](#replica-auto-balance) is set to **best-effort**. To disable this setting (disk pressure threshold) when replica auto-balance is set to best-effort, set the value of this setting to **0**.
- At least one other disk on the node has sufficient available space.

This setting is not affected by [Replica Node Level Soft Anti-Affinity](#replica-node-level-soft-anti-affinity), which can prevent Longhorn from rebuilding a replica on the same node. Regardless of that setting's value, this setting still allows Longhorn to attempt replica rebuilding on a different disk on the same node for migration purposes.

#### Storage Minimal Available Percentage

> Default: `25`

This setting controls the minimum free space that must remain on a disk, based on its **Storage Maximum**, before Longhorn can schedule a new replica.

By default, Longhorn ensures that at least **25%** of the disk's total capacity remains free. If adding a replica would reduce the available space below this limit, the disk is temporarily marked as unavailable for scheduling until enough space is freed.

This helps protect your disks from becoming too full, which can cause performance issues or storage failures. Keeping a buffer of free space helps to keep the system stable and ensures there is room for unexpected storage needs.

See [Multiple Disks Support](../../nodes-and-volumes/nodes/multidisk/#configuration) for details.

#### Storage Over Provisioning Percentage

> Default: `100`

The over-provisioning percentage defines the amount of storage that can be allocated relative to the hard drive's capacity.

Adjusting this setting allows Longhorn Manager to schedule new replicas on a disk as long as the combined size of all replicas stays within the permitted over-provisioning percentage of the usable disk space. The usable disk space is calculated as **Storage Maximum** minus **Storage Reserved**.

Note that replicas may consume more space than the volume’s nominal size due to snapshot data. To reclaim disk space, you can delete snapshots that are no longer needed.

> **Example**
>
> Suppose a disk has a Storage Maximum of 100 GiB and Storage Reserved of 10 GiB, resulting in 90 GiB of usable capacity.
>
> If the Storage Over-Provisioning Percentage is set to 200%, the maximum allowed Storage Scheduled is 180 GiB (200% of 90 GiB).
>
> This means Longhorn Manager can continue scheduling replicas to this disk until the total scheduled size reaches 180 GiB, even though the actual usable space is only 90 GiB.

#### Storage Reserved Percentage For Default Disk

> Default: `30`

The reserved percentage specifies the percentage of disk space that will not be allocated to the default disk on each new Longhorn node.

This setting only affects the default disk of a new adding node or nodes when installing Longhorn.

#### Allow Empty Node Selector Volume

> Default: `true`

This setting allows replica of the volume without node selector to be scheduled on node with tags.

#### Allow Empty Disk Selector Volume

> Default: `true`

This setting allows replica of the volume without disk selector to be scheduled on disk with tags.

### Danger Zone

Starting with Longhorn v1.6.0, Longhorn allows you to modify the [Danger Zone settings](https://longhorn.io/docs/1.6.0/references/settings/#danger-zone) without the need to wait for all volumes to become detached. Your preferred settings are immediately applied in the following scenarios:

- No attached volumes: When no volumes are attached before the settings are configured, the setting changes are immediately applied.
- Engine image upgrade (live upgrade): During a live upgrade, which involves creating a new Instance Manager pod, the setting changes are immediately applied to the new pod.

Settings are synchronized hourly. When all volumes are detached, the settings in the following table are immediately applied and the system-managed components (for example, Instance Manager, CSI Driver, and engine images) are restarted.

If you do not detach all volumes before the settings are synchronized, the settings are not applied and you must reconfigure the same settings after detaching the remaining volumes. You can view the list of unapplied settings in the **Danger Zone** section of the Longhorn UI, or run the following CLI command to check the value of the field `APPLIED`.

  ```shell
  ~# kubectl -n longhorn-system get setting priority-class
  NAME             VALUE               APPLIED   AGE
  priority-class   longhorn-critical   true      3h26m
  ```

  | Setting | Additional Information| Affected Components |
  | --- | --- | --- |
  | [Kubernetes Taint Toleration](#kubernetes-taint-toleration)| [Taints and Tolerations](../../advanced-resources/deploy/taint-toleration/) | System-managed components |
  | [Priority Class](#priority-class) | [Priority Class](../../advanced-resources/deploy/priority-class/) | System-managed components |
  | [System Managed Components Node Selector](#system-managed-components-node-selector) | [Node Selector](../../advanced-resources/deploy/node-selector/) | System-managed components |
  | [Storage Network](#storage-network) | [Storage Network](../../advanced-resources/deploy/storage-network/) | Instance Manager and Backing Image components |
  | [V1 Data Engine](#v1-data-engine) || Instance Manager component |
  | [V2 Data Engine](#v2-data-engine) | [V2 Data Engine (Experimental)](../../v2-data-engine/) | Instance Manager component |
  | [Guaranteed Instance Manager CPU](#guaranteed-instance-manager-cpu) || Instance Manager component |

For V1 and V2 Data Engine settings, you can disable the Data Engines only when all associated volumes are detached. For example, you can disable the V2 Data Engine only when all V2 volumes are detached (even when V1 volumes are still attached).

#### V1 Data Engine

> Default: `true`

Setting that allows you to enable the V1 Data Engine.

#### V2 Data Engine

> Default: `false`

Setting that allows you to enable the V2 Data Engine, which is based on the Storage Performance Development Kit (SPDK). The V2 Data Engine is an experimental feature and should not be used in production environments. For more information, see [V2 Data Engine (Experimental)](../../v2-data-engine).

> **Warning**
>
> - DO NOT CHANGE THIS SETTING WITH ATTACHED VOLUMES. Longhorn will block this setting update when there are attached volumes.
>
> - When the V2 Data Engine is enabled, each instance-manager pod utilizes 1 CPU core. This high CPU usage is attributed to the Storage Performance Development Kit (SPDK) target daemon running within each instance-manager pod. The SPDK target daemon is responsible for handling input/output (IO) operations and requires intensive polling. As a result, it consumes 100% of a dedicated CPU core to efficiently manage and process the IO requests, ensuring optimal performance and responsiveness for storage operations.

#### Concurrent Replica Rebuild Per Node Limit

> Default: `5`

This setting controls how many replicas on a node can be rebuilt simultaneously.

Typically, Longhorn can block the replica starting once the current rebuilding count on a node exceeds the limit. But when the value is 0, it means disabling the replica rebuilding.

> **WARNING:**
>  - The old setting "Disable Replica Rebuild" is replaced by this setting.
>  - Different from relying on replica starting delay to limit the concurrent rebuilding, if the rebuilding is disabled, replica object replenishment will be directly skipped.
>  - When the value is 0, the eviction and data locality feature won't work. But this shouldn't have any impact to any current replica rebuild and backup restore.

#### Concurrent Backing Image Replenish Per Node Limit

> Default: `5`

This setting controls how many backing image copies on a node can be replenished simultaneously.

Typically, Longhorn can block the backing image copy starting once the current replenishing count on a node exceeds the limit. But when the value is 0, it means disabling the backing image replenish.

#### Kubernetes Taint Toleration

> Example: `nodetype=storage:NoSchedule`

If you want to dedicate nodes to just store Longhorn replicas and reject other general workloads, you can set tolerations for **all** Longhorn components and add taints to the nodes dedicated for storage.

Longhorn system contains user deployed components (e.g, Longhorn manager, Longhorn driver, Longhorn UI) and system managed components (e.g, instance manager, engine image, CSI driver, etc.)
This setting only sets taint tolerations for system managed components.
Depending on how you deployed Longhorn, you need to set taint tolerations for user deployed components in Helm chart or deployment YAML file.

To apply the modified toleration setting immediately, ensure that all Longhorn volumes are detached. When volumes are in use, Longhorn components are not restarted, and you need to reconfigure the settings after detaching the remaining volumes; otherwise, you can wait for the setting change to be reconciled in an hour.
We recommend setting tolerations during Longhorn deployment because the Longhorn system cannot be operated during the update.

Multiple tolerations can be set here, and these tolerations are separated by semicolon. For example:
* `key1=value1:NoSchedule; key2:NoExecute`
* `:` this toleration tolerates everything because an empty key with operator `Exists` matches all keys, values and effects
* `key1=value1:`  this toleration has empty effect. It matches all effects with key `key1`
  See [Taint Toleration](../../advanced-resources/deploy/taint-toleration) for details.

#### Priority Class

> Default: `longhorn-critical`

By default, Longhorn workloads run with the same priority as other pods in the cluster, meaning in cases of node pressure, such as a node running out of memory, Longhorn workloads will be at the same priority as other Pods for eviction.

The Priority Class setting will specify a Priority Class for the Longhorn workloads to run as. This can be used to set the priority for Longhorn workloads higher so that they will not be the first to be evicted when a node is under pressure.

Longhorn system contains user deployed components (e.g, Longhorn manager, Longhorn driver, Longhorn UI) and system managed components (e.g, instance manager, engine image, CSI driver, etc.).

Note that this setting only sets Priority Class for system managed components.
Depending on how you deployed Longhorn, you need to set Priority Class for user deployed components in Helm chart or deployment YAML file.

> **Warning:** This setting should only be changed after detaching all Longhorn volumes, as the Longhorn system components will be restarted to apply the setting. The Priority Class update will take a while, and users cannot operate Longhorn system during the update. Hence, it's recommended to set the Priority Class during Longhorn deployment.

See [Priority Class](../../advanced-resources/deploy/priority-class) for details.

#### System Managed Components Node Selector

> Example: `label-key1:label-value1;label-key2:label-value2`

If you want to restrict Longhorn components to only run on a particular set of nodes, you can set node selector for all Longhorn components.

Longhorn system contains user deployed components (e.g, Longhorn manager, Longhorn driver, Longhorn UI) and system managed components (e.g, instance manager, engine image, CSI driver, etc.)
You need to set node selector for both of them. This setting only sets node selector for system managed components. Follow the instruction at [Node Selector](../../advanced-resources/deploy/node-selector) to change node selector.

> **Warning:**  Since all Longhorn components will be restarted, the Longhorn system is unavailable temporarily.
To apply a setting immediately, ensure that all Longhorn volumes are detached. When volumes are in use, Longhorn components are not restarted, and you need to reconfigure the settings after detaching the remaining volumes; otherwise, you can wait for the setting change to be reconciled in an hour.
Don't operate the Longhorn system while node selector settings are updated and Longhorn components are being restarted.

#### Kubernetes Cluster Autoscaler Enabled (Experimental)

> Default: `false`

Setting the Kubernetes Cluster Autoscaler Enabled to `true` allows Longhorn to unblock the Kubernetes Cluster Autoscaler scaling.

See [Kubernetes Cluster Autoscaler Support](../../high-availability/k8s-cluster-autoscaler) for details.

> **Warning:** Replica rebuilding could be expensive because nodes with reusable replicas could get removed by the Kubernetes Cluster Autoscaler.

#### Storage Network

> Example: `kube-system/demo-192-168-0-0`

The storage network uses Multus NetworkAttachmentDefinition to segregate the in-cluster data traffic from the default Kubernetes cluster network.

By default, the this setting applies only to RWO (Read-Write-Once) volumes. For RWX (Read-Write-Many) volumes, see [Storage Network for RWX Volume Enabled](#storage-network-for-rwx-volume-enabled) setting.

> **Warning:** This setting should change after all Longhorn volumes are detached because some pods that run Longhorn system components are recreated to apply the setting. When all volumes are detached, Longhorn attempts to restart all Instance Manager and Backing Image Manager pods immediately. When volumes are in use, Longhorn components are not restarted, and you need to reconfigure the settings after detaching the remaining volumes; otherwise, you can wait for the setting change to be reconciled in an hour.

See [Storage Network](../../advanced-resources/deploy/storage-network) for details.

#### Storage Network For RWX Volume Enabled

> Default: `false`

This setting allows Longhorn to use the storage network for RWX volumes.

> **Warning:**
> This setting should change after all Longhorn RWX volumes are detached because some pods that run Longhorn components are recreated to apply the setting. When all RWX volumes are detached, Longhorn attempts to restart all CSI plugin pods immediately. When volumes are in use, pods that run Longhorn components are not restarted, so the settings must be reconfigured after the remaining volumes are detached. If you are unable to manually reconfigure the settings, you can opt to wait because settings are synchronized hourly.
>
> The RWX volumes are mounted with the storage network within the CSI plugin pod container network namespace. As a result, restarting the CSI plugin pod may lead to unresponsive RWX volume mounts. When this occurs, you must restart the workload pod to re-establish the mount connection. Alternatively, you can enable the [Automatically Delete Workload Pod when The Volume Is Detached Unexpectedly](#automatically-delete-workload-pod-when-the-volume-is-detached-unexpectedly) setting.

For more information, see [Storage Network](../../advanced-resources/deploy/storage-network).

#### Remove Snapshots During Filesystem Trim

> Example: `false`

This setting allows Longhorn filesystem trim feature to automatically mark the latest snapshot and its ancestors as removed and stops at the snapshot containing multiple children.

Since Longhorn filesystem trim feature can be applied to the volume head and the followed continuous removed or system snapshots only.

Notice that trying to trim a removed files from a valid snapshot will do nothing but the filesystem will discard this kind of in-memory trimmable file info. Later on if you mark the snapshot as removed and want to retry the trim, you may need to unmount and remount the filesystem so that the filesystem can recollect the trimmable file info.

See [Trim Filesystem](../../nodes-and-volumes/volumes/trim-filesystem) for details.

#### Guaranteed Instance Manager CPU

> Default: `{"v1":"12","v2":"12"}`

Percentage of the total allocatable CPU resources on each node to reserve for each instance manager pod. For example, a value of `10` means 10% of the total CPU on a node will be allocated to each instance manager pod on that node. This helps maintain engine and replica stability during periods of high node workload.

In order to prevent unexpected volume instance (engine/replica) crash as well as guarantee a relative acceptable IO performance, you can use the following formula to calculate a value for this setting:

    Guaranteed Instance Manager CPU = The estimated max Longhorn volume engine and replica count on a node * 0.1 / The total allocatable CPUs on the node * 100.

The result of above calculation doesn't mean that's the maximum CPU resources the Longhorn workloads require. To fully exploit the Longhorn volume I/O performance, you can allocate/guarantee more CPU resources via this setting.

If it's hard to estimate the usage now, you can leave it with the default value, which is 12%. Then you can tune it when there is no running workload using Longhorn volumes.

> **Warning:**
>  - Value 0 means unsetting CPU requests for instance manager pods.
>  - Considering the possible new instance manager pods in the further system upgrade, this float value ranges from 0 to 40.
>  - One more set of instance manager pods may need to be deployed when the Longhorn system is upgraded. If current available CPUs of the nodes are not enough for the new instance manager pods, you need to detach the volumes using the oldest instance manager pods so that Longhorn can clean up the old pods automatically and release the CPU resources. And the new pods with the latest instance manager image will be launched then.
>  - This global setting will be ignored for a node if the field "InstanceManagerCPURequest" on the node is set.
>  - For the v2 Data Engine, the Storage Performance Development Kit (SPDK) target daemon inside each instance manager pod uses one or more dedicated CPU cores. Setting a minimum CPU usage is critical to maintaining stability during periods of high node load.

#### Disable Snapshot Purge

> Default: `false`

When set to true, temporarily prevent all attempts to purge volume snapshots.

Longhorn typically purges snapshots during replica rebuilding and user-initiated snapshot deletion. While purging,
Longhorn coalesces unnecessary snapshots into their newer counterparts, freeing space consumed by historical data.

Allowing snapshot purging during normal operations is ideal, but this process temporarily consumes additional disk
space. If insufficient disk space prevents the process from continuing, consider temporarily disabling purging while
data is moved to other disks.

#### Auto Cleanup Snapshot When Delete Backup

> Default: `false`

When set to true, the snapshot used by the backup will be automatically cleaned up when the backup is deleted.

#### Auto Cleanup Snapshot After On-Demand Backup Completed

> Default: `false`

When set to true, the snapshot used by the backup will be automatically cleaned up after the on-demand backup is completed.

#### Instance Manager Pod Liveness Probe Timeout

> Default: `10`

In seconds. The setting specifies the timeout for the instance manager pod liveness probe. The default value is 10 seconds.

> **Warning**
>
> When applying the setting, Longhorn will try to restart all instance-manager pods if all volumes are detached and eventually restart the instance manager pod without instances running on the instance manager.

#### Data Engine CPU Mask

> Default: `{"v2":"0x1"}`

Applies only to the V2 Data Engine. Specifies the CPU cores on which the Storage Performance Development Kit (SPDK) target daemon runs. The daemon is deployed in each Instance Manager pod. Ensure that the number of assigned cores does not exceed the guaranteed Instance Manager CPUs for the V2 Data Engine.

#### Data Engine Hugepage Enabled

> Default: `{"v2":"true"}`

Applies only to the V2 Data Engine. Enables hugepages for the Storage Performance Development Kit (SPDK) target daemon. If disabled, legacy memory is used. Allocation size is set via the Data Engine Memory Size setting.

#### Data Engine Memory Size

> Default: `{"v2":"2048"}`

Applies only to the V2 Data Engine. Specifies the memory size, in MiB, allocated to the Storage Performance Development Kit (SPDK) target daemon. When hugepage is enabled, this defines the hugepage size; when legacy memory is used, hugepage is disabled.

#### Data Engine Interrupt Mode Enabled

> Default: `{"v2":"false"}`

Applies only to the **V2 Data Engine**.

Controls whether the Storage Performance Development Kit (SPDK) target daemon runs in **interrupt mode** or the default **polling mode**.

- `true`: Enables interrupt mode, reducing CPU usage by handling I/O through interrupts.
- `false`: Keeps polling mode enabled for maximum performance and lowest latency.

> **Warning**
> - DO NOT CHANGE THIS SETTING WITH ATTACHED VOLUMES. Longhorn will block this setting update when there are attached v2 volumes.

#### Log Path

> Default: `/var/lib/longhorn/logs/`

Specifies the directory on the host where Longhorn stores log files for the instance manager pod. Currently, it is only used for instance manager pods in the v2 data engine.


---

## Article: references/storage-class-parameters.md

---
title: Storage Class Parameters
weight: 1
---

## Overview

Storage Class as a resource object has a number of settable parameters.  Here's a sample YAML:
```yaml
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: longhorn-test
provisioner: driver.longhorn.io
allowVolumeExpansion: true
reclaimPolicy: Delete
volumeBindingMode: Immediate
parameters:
  backupTargetName: "default"
  numberOfReplicas: "3"
  staleReplicaTimeout: "2880"
  fromBackup: ""
  fsType: "ext4"
#  mkfsParams: ""
#  migratable: false
#  encrypted: false
#  dataLocality: "disabled"
#  replicaAutoBalance: "ignored"
#  diskSelector: "ssd,fast"
#  nodeSelector: "storage,fast"
#  recurringJobSelector: '[{"name":"snap-group", "isGroup":true},
#                          {"name":"backup", "isGroup":false}]'
#  backingImageName: ""
#  backingImageChecksum: ""
#  backingImageDataSourceType: ""
#  backingImageDataSourceParameters: ""
#  unmapMarkSnapChainRemoved: "ignored"
#  disableRevisionCounter: false
#  replicaSoftAntiAffinity: "ignored"
#  replicaZoneSoftAntiAffinity: "ignored"
#  replicaDiskSoftAntiAffinity: "ignored"
#  nfsOptions: "soft,timeo=150,retrans=3"
#  dataEngine: "v1"
#  freezeFSForSnapshot: "ignored"
```

## Built-in Fields
Some fields are common to all Kubernetes storage classes.
See also [Kubernetes Storage Class](https://kubernetes.io/docs/concepts/storage/storage-classes).

#### Provisioner *(field: `provisioner`)*
Specifies the plugin that will be used for dynamic creation of persistent volumes.  For Longhorn, that is always "driver.longhorn.io".
> See [Kubernetes Storage Class: Provisioner](https://kubernetes.io/docs/concepts/storage/storage-classes/#provisioner).

#### Allow Volume Expansion *(field: `allowVolumeExpansion`)*
> Default: `true`
> See [Kubernetes Storage Class: Allow Volume Expansion](https://kubernetes.io/docs/concepts/storage/storage-classes/#allow-volume-expansion).

#### Reclaim Policy *(field: `reclaimPolicy`)*
> Default: `Delete`
> See [Kubernetes Storage Class: Reclaim Policy](https://kubernetes.io/docs/concepts/storage/storage-classes/#reclaim-policy).

#### Mount Options *(field: `mountOptions`)*
> Default `[]`
> See [Kubernetes Storage Class: Mount Options](https://kubernetes.io/docs/concepts/storage/storage-classes/#mount-options).

#### Volume Binding Mode *(field: `volumeBindingMode`)*
> Default `Immediate`
> See [Kubernetes Storage Class: Volume Binding Mode](https://kubernetes.io/docs/concepts/storage/storage-classes/#volume-binding-mode).

## Longhorn-specific Parameters
Note that some of these parameters also exist and may be specified in global settings.  When a volume is provisioned with Kubernetes against a particular StorageClass, StorageClass parameters override the global settings.
These fields will be applied for new volume creation only.  If a StorageClass is modified, neither Longhorn nor Kubernetes is responsible for propagating changes to its parameters back to volumes previously created with it.

#### Number Of Replicas *(field: `parameters.numberOfReplicas`)*
> Default: `3`

The desired number of copies (replicas) for redundancy.
  - Must be between 1 and 20.
  - Replicas will be placed across the widest possible set of zones, nodes, and disks in a cluster, subject to other constraints, such as NodeSelector.

> Global setting: [Default Replica Count](../settings#default-replica-count).

#### Stale Replica Timeout *(field: `parameters.staleReplicaTimeout`)*
> Default: `30`

Minutes after a replica is marked unhealthy before it is deemed useless for rebuilds and is just deleted.

#### From Backup *(field: `parameters.fromBackup`)*
> Default: `""`
> Example: `"s3://backupbucket@us-east-1?volume=minio-vol01&backup=backup-eeb2782d5b2f42bb"`

URL of a backup to be restored from.

#### FS Type *(field: `parameters.fsType`)*
> Default: `ext4`
> For more details, see [Creating Longhorn Volumes with Kubernetes](../../nodes-and-volumes/volumes/create-volumes#creating-longhorn-volumes-with-kubectl)

#### Mkfs Params *(field: `parameters.mkfsParams`)*
> Default: `""`
> For more details, see [Creating Longhorn Volumes with Kubernetes](../../nodes-and-volumes/volumes/create-volumes#creating-longhorn-volumes-with-kubectl)

#### Migratable *(field: `parameters.migratable`)*
> Default: `false`

Allows for a Longhorn volume to be live migrated from one node to another.  Useful for volumes used by Harvester.

#### Encrypted *(field: `parameters.encrypted`)*
> Default: `false`
> More details in [Encrypted Volumes](../../advanced-resources/security/volume-encryption)

#### Data Locality *(field: `parameters.dataLocality`)*
> Default: `disabled`

If enabled, try to keep the data on the same node as the workload for better performance.
  - For "best-effort", a replica will be co-located if possible, but is permitted to find another node if not.
  - For "strict-local" the Replica count should be 1, or volume creation will fail with a parameter validation error.
  - If "strict-local" is not possible for whatever other reason, volume creation will be failed.  A "strict-local" replica that becomes displaced from its workload will be marked as "Stopped".

>  Global setting: [Default Data Locality](../settings#default-data-locality)
>  More details in [Data Locality](../../high-availability/data-locality).

#### Replica Auto-Balance *(field: `parameters.replicaAutoBalance`)*
> Default: `ignored`

If enabled, move replicas to more lightly-loaded nodes.
  - "ignored" means use the global setting.
  - Other options are "disabled", "least-effort", "best-effort".

> Global setting: [Replica Auto Balance](../settings#replica-auto-balance)
> More details in [Auto Balance Replicas](../../high-availability/auto-balance-replicas).

#### Disk Selector *(field: `parameters.diskSelector`)*
> Default: `""`
> Example: `"ssd,fast"`

A list of tags to select which disks are candidates for replica placement.
> More details in [Storage Tags](../../nodes-and-volumes/nodes/storage-tags)

#### Node Selector *(field: `parameters.nodeSelector`)*
> Default: `""`
> Example: `"storage,fast"`

A list of tags to select which nodes are candidates for replica placement.
> More details in [Storage Tags](../../nodes-and-volumes/nodes/storage-tags)

#### Recurring Job Selector *(field: `parameters.recurringJobSelector`)*
> Default: `""`
> Example:  `[{"name":"backup", "isGroup":true}]`

A list of recurring jobs that are to be run on a volume.
>  More details in [Recurring Snapshots and Backups](../../snapshots-and-backups/scheduling-backups-and-snapshots)

#### Backing Image Name *(field: `parameters.backingImageName`)*
> Default: `""`
> See [Backing Image](../../advanced-resources/backing-image/backing-image#create-and-use-a-backing-image-via-storageclass-and-pvc)

#### Backing Image Checksum *(field: `parameters.backingImageChecksum`)*
> Default: `""`
> See [Backing Image](../../advanced-resources/backing-image/backing-image#create-and-use-a-backing-image-via-storageclass-and-pvc)

#### Backing Image Data Source Type *(field: `parameters.backingImageDataSourceType`)*
> Default: `""`
> See [Backing Image](../../advanced-resources/backing-image/backing-image#create-and-use-a-backing-image-via-storageclass-and-pvc)

#### Backing Image Data Source Parameters *(field: `parameters.backingImageDataSourceParameters`)*
> Default: `""`
> See [Backing Image](../../advanced-resources/backing-image/backing-image#create-and-use-a-backing-image-via-storageclass-and-pvc)

#### Unmap Mark Snap Chain Removed *(field: `parameters.unmapMarkSnapChainRemoved`)*
> Default: `ignored`

  - "ignored" means use the global setting.
  - Other values are "enabled" and "disabled".

> Global setting: [Remove Snapshots During Filesystem Trim](../settings#remove-snapshots-during-filesystem-trim).
> More details in [Trim Filesystem](../../nodes-and-volumes/volumes/trim-filesystem).

#### Disable Revision Counter *(field: `parameters.disableRevisionCounter`)*
> Default: `true`

> Global setting: [Disable Revision Counter](../settings#disable-revision-counter).
> More details in [Revision Counter](../../advanced-resources/deploy/revision_counter).

#### Replica Soft Anti-Affinity *(field: `parameters.replicaSoftAntiAffinity`)*
> Default: `ignored`

  - "ignored" means use the global setting.
  - Other values are "enabled" and "disabled".

> Global setting: [Replica Node Level Soft Anti-Affinity](../settings#replica-node-level-soft-anti-affinity).
> More details in [Scheduling](../../nodes-and-volumes/nodes/scheduling) and [Best Practices](../../best-practices#replica-node-level-soft-anti-affinity).

#### Replica Zone Soft Anti-Affinity *(field: `parameters.replicaZoneSoftAntiAffinity`)*
> Default: `ignored`

  - "ignored" means use the global setting.
  - Other values are "enabled" and "disabled".

> Global setting: [Replica Zone Level Soft Anti-Affinity](../settings#replica-zone-level-soft-anti-affinity).
> More details in [Scheduling](../../nodes-and-volumes/nodes/scheduling).

#### Replica Disk Soft Anti-Affinity *(field: `parameters.replicaDiskSoftAntiAffinity`)*
> Default: `ignored`

  - "ignored" means use the global setting.
  - Other values are "enabled" and "disabled".

> Global setting: [Replica Disk Level Soft Anti-Affinity](../settings#replica-disk-level-soft-anti-affinity).
> More details in [Scheduling](../../nodes-and-volumes/nodes/scheduling).

#### NFS Options *(field: `parameters.nfsOptions`)*
> Default: `""`
> Example: `"hard,sync"`

  - Overrides for NFS mount of RWX volumes to the share-manager.  Use this field with caution.
  - Note:  Built-in options vary by release.  Check your release details before setting this.

> More details in [RWX Workloads](../../nodes-and-volumes/volumes/rwx-volumes#configuring-volume-mount-options)

#### Data Engine *(field: `parameters.dataEngine`)*
> Default: `"v1"`

  - Specify "v2" to enable the V2 Data Engine (experimental feature in v1.6.0). When unspecified, Longhorn uses the default value ("v1").

> Global setting: [V2 Data Engine](../settings#v2-data-engine).
> More details in [V2 Data Engine Quick Start](../../v2-data-engine/quick-start#create-a-storageclass).

#### Freeze Filesystem For Snapshot *(field: `parameters.freezeFilesystemForSnapshot`)*
> Default: `ignored`

  - "ignored" instructs Longhorn to use the global setting.
  - Other values are "enabled" and "disabled".

> Global setting: [Freeze File System For Snapshot](../settings#freeze-filesystem-for-snapshot).

#### Backup Target Name *(field: `parameters.backupTargetName`)*
> Default: `default`
> More details in [default backup target](../../snapshots-and-backups/backup-and-restore/set-backup-target#default-backup-target) and [Create Volumes](../../nodes-and-volumes/volumes/create-volumes).

#### Backup Block Size *(field: `parameters.backupBlockSize`)*
> Default: `""`
> Example: `"2Mi"` or `"16Mi"`

Kubernetes quantity string. Specify the empty string `""` to use the global setting.

> Global setting: [default backup block size](../settings#default-backup-block-size).
> More details in [Configure The Block Size Of Backup](../../snapshots-and-backups/backup-and-restore/configure-backup-block-size)

## Helm Installs

If Longhorn is installed via Helm, values in the default storage class can be set by editing the corresponding item in [values.yaml](https://github.com/longhorn/longhorn/blob/v{{< current-version >}}/chart/values.yaml).  All of the Storage Class parameters have a prefix of "persistence".  For example, `persistence.defaultNodeSelector`.


---

## Article: important-notes/_index.md

---
title: Important Notes
weight: 1
---

This page summarizes the key notes for Longhorn v{{< current-version >}}.
For the full release note, see [here](https://github.com/longhorn/longhorn/releases/tag/v{{< current-version >}}).

- [Warning](#warning)
  - [Upgrade](#upgrade)
    - [Migration Requirement Before Longhorn v1.10 Upgrade](#migration-requirement-before-longhorn-v110-upgrade)
    - [Migration Verification](#migration-verification)
    - [Troubleshooting CRD Upgrade Failures During Upgrade to Longhorn v1.10](#troubleshooting-crd-upgrade-failures-during-upgrade-to-longhorn-v110)
      - [Downgrade Procedure (kubectl Installation)](#downgrade-procedure-kubectl-installation)
      - [Downgrade Procedure (Helm Installation)](#downgrade-procedure-helm-installation)
      - [Post-Downgrade](#post-downgrade)
  - [HotFix](#hotfix)
- [Important Fixes](#important-fixes)
  - [Goroutine Leak in Instance Manager (V2 Data Engine)](#goroutine-leak-in-instance-manager-v2-data-engine)
  - [V2 Volume Attachment Failure in Interrupt Mode](#v2-volume-attachment-failure-in-interrupt-mode)
  - [UI Deployment Failure on IPv4-Only Nodes](#ui-deployment-failure-on-ipv4-only-nodes)
  - [Share Manager Excessive Memory Usage](#share-manager-excessive-memory-usage)
- [Removal](#removal)
  - [`longhorn.io/v1beta1` API](#longhorniov1beta1-api)
  - [`replica.status.evictionRequested` Field](#replicastatusevictionrequested-field)
- [General](#general)
  - [Kubernetes Version Requirement](#kubernetes-version-requirement)
  - [CRD Upgrade Validation](#crd-upgrade-validation)
  - [Upgrade Check Events](#upgrade-check-events)
  - [Manual Checks Before Upgrade](#manual-checks-before-upgrade)
  - [Consolidation of Longhorn Settings](#consolidation-of-longhorn-settings)
  - [System Info Category in Setting](#system-info-category-in-setting)
  - [Volume Attachment Summary](#volume-attachment-summary)
- [Scheduling](#scheduling)
  - [Pod Scheduling with CSIStorageCapacity](#pod-scheduling-with-csistoragecapacity)
- [Performance](#performance)
  - [Configurable Backup Block Size](#configurable-backup-block-size)
  - [Profiling Support for Backup Sync Agent](#profiling-support-for-backup-sync-agent)
- [Resilience](#resilience)
  - [Configurable Liveness Probe for Instance Manager](#configurable-liveness-probe-for-instance-manager)
  - [Backing Image Manager CR Naming](#backing-image-manager-cr-naming)
- [Security](#security)
  - [Refined RBAC Permissions](#refined-rbac-permissions)
- [V1 Data Engine](#v1-data-engine)
  - [IPv6 Support](#ipv6-support)
- [V2 Data Engine](#v2-data-engine)
  - [Longhorn System Upgrade](#longhorn-system-upgrade)
  - [New Functionalities since Longhorn v1.10.0](#new-functionalities-since-longhorn-v1100)
    - [V2 Data Engine Without Hugepage Support](#v2-data-engine-without-hugepage-support)
    - [V2 Data Engine Interrupt Mode Support](#v2-data-engine-interrupt-mode-support)
    - [V2 Data Engine Volume Clone Support](#v2-data-engine-volume-clone-support)
    - [V2 Data Engine Replica Rebuild QoS](#v2-data-engine-replica-rebuild-qos)
    - [V2 Data Engine Volume Expansion](#v2-data-engine-volume-expansion)

## Warning

### Upgrade

If your Longhorn cluster was initially deployed with a version earlier than v1.3.0, the Custom Resources (CRs) were created using the `v1beta1` APIs. While the upgrade from Longhorn v1.8 to v1.9 automatically migrates all CRs to the new `v1beta2` version, **a manual CR migration is strongly advised before upgrading from Longhorn `v1.9` to `v1.10`**.

Certain operations, such as an `etcd` or CRD restore, may leave behind `v1beta1` data. Manually migrating your CRs ensures that all Longhorn data is properly updated to the `v1beta2` API, preventing potential compatibility issues and unexpected behavior with the new Longhorn version.

Following the manual migration, **verify that `v1beta1` has been removed from the CRD stored versions** to ensure completion and a successful upgrade.

For more details, see [Kubernetes official document for CRD storage version](https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definition-versioning/#upgrade-existing-objects-to-a-new-stored-version), and [Issue #11886](https://github.com/longhorn/longhorn/issues/11886).

#### Migration Requirement Before Longhorn v1.10 Upgrade

Before upgrading from Longhorn v1.9 to v1.10, perform the following manual CRD storage version migration.

> **Note**: If your Longhorn installation uses a namespace other than `longhorn-system`, replace `longhorn-system` with your custom namespace throughout the commands.

```bash
# Temporarily disable the CR validation webhook to allow updating read-only settings CRs.
kubectl patch validatingwebhookconfiguration longhorn-webhook-validator \
  --type=merge \
  -p "$(kubectl get validatingwebhookconfiguration longhorn-webhook-validator -o json | \
  jq '.webhooks[0].rules |= map(if .apiGroups == ["longhorn.io"] and .resources == ["settings"] then
    .operations |= map(select(. != "UPDATE")) else . end)')"

# Migrate CRDs that ever stored v1beta1 resources
migration_time="$(date +%Y-%m-%dT%H:%M:%S)"
crds=($(kubectl get crd -l app.kubernetes.io/name=longhorn -o json | jq -r '.items[] | select(.status.storedVersions | index("v1beta1")) | .metadata.name'))
for crd in "${crds[@]}"; do
  echo "Migrating ${crd} ..."
  for name in $(kubectl -n longhorn-system get "$crd" -o jsonpath='{.items[*].metadata.name}'); do
    # Attach additional annotations to trigger v1beta1 resource updating in the latest storage version.
    kubectl patch "${crd}" "${name}" -n longhorn-system --type=merge -p='{"metadata":{"annotations":{"migration-time":"'"${migration_time}"'"}}}'
  done
  # Clean up the stored version in CRD status
  kubectl patch crd "${crd}" --type=merge -p '{"status":{"storedVersions":["v1beta2"]}}' --subresource=status
done

# Re-enable the CR validation webhook.
kubectl patch validatingwebhookconfiguration longhorn-webhook-validator \
  --type=merge \
  -p "$(kubectl get validatingwebhookconfiguration longhorn-webhook-validator -o json | \
  jq '.webhooks[0].rules |= map(if .apiGroups == ["longhorn.io"] and .resources == ["settings"] then
    .operations |= (. + ["UPDATE"] | unique) else . end)')"
```

#### Migration Verification

After running the script, verify the CRD stored versions using this command:

```bash
kubectl get crd -l app.kubernetes.io/name=longhorn -o=jsonpath='{range .items[*]}{.metadata.name}{": "}{.status.storedVersions}{"\n"}{end}'
```

Crucially, all Longhorn CRDs MUST have only `"v1beta2"` listed in `storedVersions` (i.e., `"v1beta1"` must be completely absent) before proceeding to the v1.10 upgrade.

Example of successful output:

```
backingimagedatasources.longhorn.io: ["v1beta2"]
backingimagemanagers.longhorn.io: ["v1beta2"]
backingimages.longhorn.io: ["v1beta2"]
backupbackingimages.longhorn.io: ["v1beta2"]
backups.longhorn.io: ["v1beta2"]
backuptargets.longhorn.io: ["v1beta2"]
backupvolumes.longhorn.io: ["v1beta2"]
engineimages.longhorn.io: ["v1beta2"]
engines.longhorn.io: ["v1beta2"]
instancemanagers.longhorn.io: ["v1beta2"]
nodes.longhorn.io: ["v1beta2"]
orphans.longhorn.io: ["v1beta2"]
recurringjobs.longhorn.io: ["v1beta2"]
replicas.longhorn.io: ["v1beta2"]
settings.longhorn.io: ["v1beta2"]
sharemanagers.longhorn.io: ["v1beta2"]
snapshots.longhorn.io: ["v1beta2"]
supportbundles.longhorn.io: ["v1beta2"]
systembackups.longhorn.io: ["v1beta2"]
systemrestores.longhorn.io: ["v1beta2"]
volumeattachments.longhorn.io: ["v1beta2"]
volumes.longhorn.io: ["v1beta2"]
```

With these steps completed, the Longhorn upgrade to v1.10 should now proceed without issues.

#### Troubleshooting CRD Upgrade Failures During Upgrade to Longhorn v1.10

If you did not apply the required pre-upgrade migration steps and the CRs are not fully migrated to `v1beta2`, the `longhorn-manager` Pods may fail to operate correctly. A common error message for this issue is:

```
Upgrade failed: cannot patch "backingimagedatasources.longhorn.io" with kind CustomResourceDefinition: CustomResourceDefinition.apiextensions.k8s.io "backingimagedatasources.longhorn.io" is invalid: status.storedVersions[0]: Invalid value: "v1beta1": missing from spec.versions; v1beta1 was previously a storage version, and must remain in spec.versions until a storage migration ensures no data remains persisted in v1beta1 and removes v1beta1 from status.storedVersions
```

To fix this issue, you must perform a **forced downgrade** back to the **exact Longhorn v1.9.x version** that was running before the failed upgrade attempt.

##### Downgrade Procedure (kubectl Installation)

If Longhorn was installed using `kubectl`, you must patch the `current-longhorn-version` setting before downgrading. Replace `"v1.9.x" to the original version before upgrade in the following commands.

```bash
# Attaching annotation to allow patching current-longhorn-version.
kubectl patch settings.longhorn.io current-longhorn-version -n longhorn-system --type=merge -p='{"metadata":{"annotations":{"longhorn.io/update-setting-from-longhorn":""}}}'
# Temporarily override current version to allow old version installation
# Replace the value "v1.9.x" with the original version before upgrade.
kubectl patch settings.longhorn.io current-longhorn-version -n longhorn-system --type=merge -p='{"value":"v1.9.x"}'
```

After modifying `current-longhorn-version`, you can proceed to downgrade to the original Longhorn v1.9.x deployment.

##### Downgrade Procedure (Helm Installation)

If Longhorn was installed using Helm, the downgrade is allowed by disabling the [`preUpgradeChecker.upgradeVersionCheck`](https://github.com/longhorn/longhorn/tree/v1.9.x/chart#other-settings) flag.

##### Post-Downgrade

Once the downgrade is complete and the Longhorn system is stable on the v1.9.x version, you must immediately follow the steps outlined in the [Manual CRD Migration Guide](#migration-requirement-before-longhorn-v110-upgrade). This step is crucial to migrate all remaining `v1beta1` CRs to `v1beta2` before attempting the Longhorn v1.10 upgrade again.

### HotFix

The `longhorn-manager:v1.10.1` image is affected by a [regression issue](https://github.com/longhorn/longhorn/issues/12233) that can trigger a nil-pointer dereference under certain conditions, potentially causing unexpected crashes. To mitigate this issue, replace `longhorn-manager:v1.10.1` with the hotfixed image `longhorn-manager:v1.10.1-hotfix-1`.

Follow these steps to apply the update:

1. **Disable the upgrade version check**
   - Helm users: Set `upgradeVersionCheck` to `false` in the `values.yaml` file.
   - Manifest users: Remove the `--upgrade-version-check` flag from the deployment manifest.

2. **Update the `longhorn-manager` image**
   - Change the image tag from `v1.10.1` to `v1.10.1-hotfix-1` in the appropriate file:
     - For Helm: Update `values.yaml`
     - For manifests: Update the deployment manifest directly.

3. **Proceed with the upgrade**
   - Apply the changes using your standard Helm upgrade command or reapply the updated manifest.

## Important Fixes

This release includes several critical stability and performance improvements:

### Goroutine Leak in Instance Manager (V2 Data Engine)

Fixed a goroutine leak in the instance manager when using the V2 data engine. This issue could lead to increased memory usage and potential stability problems over time.

For more details, see [Issue #11962](https://github.com/longhorn/longhorn/issues/11962).

### V2 Volume Attachment Failure in Interrupt Mode

Fixed an issue where V2 volumes using interrupt mode with NVMe disks could fail to complete the attachment process, causing volumes to remain stuck in the attaching state indefinitely.

For more details, see [Issue #11816](https://github.com/longhorn/longhorn/issues/11816).

### UI Deployment Failure on IPv4-Only Nodes

Fixed a bug introduced in v1.10.0 where the Longhorn UI failed to deploy on nodes with only IPv4 enabled. The UI now correctly supports IPv4-only configurations without requiring IPv6.

For more details, see [Issue #11875](https://github.com/longhorn/longhorn/issues/11875).

### Share Manager Excessive Memory Usage

Fixed excessive memory consumption in the share manager for RWX (ReadWriteMany) volumes. The component now maintains stable memory usage under normal operation.

For more details, see [Issue #12043](https://github.com/longhorn/longhorn/issues/12043).

## Removal

### `longhorn.io/v1beta1` API

The `v1beta1` Longhorn API version was removed in v1.10.0.

For more details, see [Issue #10249](https://github.com/longhorn/longhorn/issues/10249).

### `replica.status.evictionRequested` Field

The deprecated `replica.status.evictionRequested` field has been removed.

For more details, see [Issue #7022](https://github.com/longhorn/longhorn/issues/7022)

## General

### Kubernetes Version Requirement

Due to the upgrade of the CSI external snapshotter to v8.2.0, all clusters must be running Kubernetes v1.25 or later before you can upgrade to Longhorn v1.8.0 or a newer version.

### CRD Upgrade Validation

During an upgrade, a new Longhorn manager may start before the Custom Resource Definitions (CRDs) are applied. This sequencing ensures the controller does not process objects containing deprecated data or fields. However, it can cause the Longhorn manager to fail during the initial upgrade phase if the CRD has not yet been applied.

If the Longhorn manager crashes during the upgrade, check the logs to determine if the failure is due to the CRD not being applied. In such cases, the logs may contain error messages similar to the following:

```
time="2025-03-27T06:59:55Z" level=fatal msg="Error starting manager: upgrade resources failed: BackingImage in version \"v1beta2\" cannot be handled as a BackingImage: strict decoding error: unknown field \"spec.diskFileSpecMap\", unknown field \"spec.diskSelector\", unknown field \"spec.minNumberOfCopies\", unknown field \"spec.nodeSelector\", unknown field \"spec.secret\", unknown field \"spec.secretNamespace\"" func=main.main.DaemonCmd.func3 file="daemon.go:94"
```

### Upgrade Check Events

When upgrading via Helm or Rancher App Marketplace, Longhorn performs pre-upgrade checks. If a check fails, the upgrade stops, and the reason for the failure is recorded in an event.

For more detail, see [Upgrading Longhorn Manager](../deploy/upgrade/longhorn-manager).

### Manual Checks Before Upgrade

Automated pre-upgrade checks do not cover all scenarios. Manual checks via kubectl or the UI are recommended:

- Ensure all V2 Data Engine volumes are detached and replicas are stopped. The V2 engine does not support live upgrades.
- Avoid upgrading when volumes are "Faulted", as unusable replicas may be deleted, causing permanent data loss if no backups exist.
- Avoid upgrading if a failed BackingImage exists. See [Backing Image](../advanced-resources/backing-image/backing-image) for details.
- Creating a [Longhorn system backup](../advanced-resources/system-backup-restore/backup-longhorn-system) before upgrading is recommended to ensure recoverability.

### Consolidation of Longhorn Settings

Settings have been consolidated for easier management across V1 and V2 Data Engines. Each setting now uses one of the following formats:

- Single value for all supported Data Engines
  - Format: Non-JSON string (e.g., `1024`)
  - The value applies to all supported Data Engines and must be the same across them.
  - Data-engine-specific values are not allowed.
- Data-engine-specific values for V1 and V2 Data Engines
  - Format: JSON object (e.g., `{"v1": "value1", "v2": "value2"}`)
  - Allows specifying different values for V1 and V2 Data Engines.
- Data-engine-specific values for V1 Data Engine only
  - Format: JSON object with `v1` key only (e.g., `{"v1": "value1"}`)
  - Only the V1 Data Engine can be configured; the V2 Data Engine is not affected.
- Data-engine-specific values for V2 Data Engine only
  - Format: JSON object with `v2` key only (e.g., `{"v2": "value1"}`)
  - Only the V2 Data Engine can be configured; the V1 Data Engine is not affected.

For more information, see [Longhorn Settings](../references/settings).

### System Info Category in Setting

A new **System Info** category has been added to show cluster-level information more clearly.

For more details, see [Issue #11656](https://github.com/longhorn/longhorn/issues/11656)

### Volume Attachment Summary

The UI now display a summary of attachment tickets on each volume overview page for improved visibility into volume state.

For more details, see [Issue #11400](https://github.com/longhorn/longhorn/issues/11400) and [Issue #11401](https://github.com/longhorn/longhorn/issues/11401).

## Scheduling

### Pod Scheduling with CSIStorageCapacity

Longhorn now supports Kubernetes **CSIStorageCapacity**, which enables the scheduler to verify node storage before scheduling pods that use StorageClasses with **WaitForFirstConsumer**.

This reduces scheduling errors and improves reliability.

For more information, see [GitHub Issue #10685](https://github.com/longhorn/longhorn/issues/10685)

## Performance

### Configurable Backup Block Size

Starting in Longhorn v1.10.0,  backup block size can be configured when creating a volume, allowing optimization for performance, efficiency, and cost.

For more information, see [Create Longhorn Volumes](../nodes-and-volumes/volumes/create-volumes).

### Profiling Support for Backup Sync Agent

The backup sync agent exposes a `pprof` server for profiling runtime resource usage during backup sync operations.

For more information, see [Profiling](../troubleshoot/troubleshooting#profiling).

## Resilience

### Configurable Liveness Probe for Instance Manager

You can now configure the instance-manager pod liveness probes. This allows the system to better distinguish between temporary delays and actual failures, which helps reduce unnecessary restarts and improves overall cluster stability.

For more information, see [Longhorn Settings](../references/settings#instance-manager-pod-liveness-probe-timeout).

### Backing Image Manager CR Naming

Backing Image Manager CRs now use a compact, collision-resistant naming format to reduce conflict risk.

For details, see [Issue #11455](https://github.com/longhorn/longhorn/issues/11455)

## Security

### Refined RBAC Permissions

RBAC permissions have been refined to minimize privileges and improve cluster security.

For details, see [Issue #11345](https://github.com/longhorn/longhorn/issues/11345)

## V1 Data Engine

### IPv6 Support

V1 volumes now support single-stack IPv6 Kubernetes clusters.

> **Warning:** Dual-stack Kubernetes clusters and V2 volumes are not supported in this release.

For details, see [Issue #2259](https://github.com/longhorn/longhorn/issues/2259).

## V2 Data Engine

### Longhorn System Upgrade

Live upgrades of V2 volumes are **not supported**. Ensure all V2 volumes are detached before upgrading.

### New Functionalities since Longhorn v1.10.0

#### V2 Data Engine Without Hugepage Support

The V2 Data Engine can run without Hugepage by setting `data-engine-hugepage-enabled` to `{"v2":"false"}`.

This reduces memory pressure on low‑spec nodes and increases deployment flexibility. Performance may be lower compared to running with Hugepage.

#### V2 Data Engine Interrupt Mode Support

Interrupt mode has been added to the V2 Data Engine to help reduce CPU usage. This feature is especially beneficial for clusters with idle or low I/O workloads, where conserving CPU resources is more important than minimizing latency.

While interrupt mode lowers CPU consumption, it may introduce slightly higher I/O latency compared to polling mode. In addition, the current implementation uses a hybrid approach, which still incurs a minimal, constant CPU load even when interrupts are enabled.

For more information, see [Interrupt Mode](../v2-data-engine/features/interrupt-mode) for more information.

> **Note:** In Longhorn v1.10.0, interrupt mode supports only **AIO disks**. Interrupt mode for **NVMe disks** is supported starting in v1.10.1.

#### V2 Data Engine Volume Clone Support

Longhorn now supports volume and snapshot cloning for V2 data engine volumes.
For more information, see [Volume Clone Support](../v2-data-engine/features/volume-clone).

#### V2 Data Engine Replica Rebuild QoS

Provides Quality of Service (QoS) control for V2 volume replica rebuilds. You can configure bandwidth limits globally or per volume to prevent storage throughput overload on source and destination nodes.

For more information, see [Replica Rebuild QoS](../v2-data-engine/features/replica-rebuild-qos).

#### V2 Data Engine Volume Expansion

Longhorn now supports volume expansion for V2 Data Engine volumes. Users can expand the volume through the UI or by modifying the PVC manifest.

For more information, see [V2 Volume Expansion](../v2-data-engine/features/volume-expansion).


---

## Article: monitoring/_index.md

---
title: Monitoring
weight: 7
---

* Setting up Prometheus and Grafana to monitor Longhorn
* Integrating Longhorn metrics into the Rancher monitoring system
* Longhorn Metrics for Monitoring
* Support Kubelet Volume Metrics
* Longhorn Alert Rule Examples


---

## Article: monitoring/alert-rules-example.md

---
title: Longhorn Alert Rule Examples
weight: 5
---

We provide a couple of example Longhorn alert rules below for your references.
See [here](../metrics) for a list of all available Longhorn metrics and build your own alert rules.

```yaml
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  labels:
    prometheus: longhorn
    role: alert-rules
  name: prometheus-longhorn-rules
  namespace: monitoring
spec:
  groups:
  - name: longhorn.rules
    rules:
    - alert: LonghornVolumeActualSpaceUsedWarning
      annotations:
        description: The actual space used by Longhorn volume {{$labels.volume}} on {{$labels.node}} is at {{$value}}% capacity for
          more than 5 minutes.
        summary: The actual used space of Longhorn volume is over 90% of the capacity.
      expr: (longhorn_volume_actual_size_bytes / longhorn_volume_capacity_bytes) * 100 > 90
      for: 5m
      labels:
        issue: The actual used space of Longhorn volume {{$labels.volume}} on {{$labels.node}} is high.
        severity: warning
    - alert: LonghornVolumeStatusCritical
      annotations:
        description: Longhorn volume {{$labels.volume}} on {{$labels.node}} is Fault for
          more than 2 minutes.
        summary: Longhorn volume {{$labels.volume}} is Fault
      expr: longhorn_volume_robustness == 3
      for: 5m
      labels:
        issue: Longhorn volume {{$labels.volume}} is Fault.
        severity: critical
    - alert: LonghornVolumeStatusWarning
      annotations:
        description: Longhorn volume {{$labels.volume}} on {{$labels.node}} is Degraded for
          more than 5 minutes.
        summary: Longhorn volume {{$labels.volume}} is Degraded
      expr: longhorn_volume_robustness == 2
      for: 5m
      labels:
        issue: Longhorn volume {{$labels.volume}} is Degraded.
        severity: warning
    - alert: LonghornNodeStorageWarning
      annotations:
        description: The used storage of node {{$labels.node}} is at {{$value}}% capacity for
          more than 5 minutes.
        summary:  The used storage of node is over 70% of the capacity.
      expr: (longhorn_node_storage_usage_bytes / longhorn_node_storage_capacity_bytes) * 100 > 70
      for: 5m
      labels:
        issue: The used storage of node {{$labels.node}} is high.
        severity: warning
    - alert: LonghornDiskStorageWarning
      annotations:
        description: The used storage of disk {{$labels.disk}} on node {{$labels.node}} is at {{$value}}% capacity for
          more than 5 minutes.
        summary:  The used storage of disk is over 70% of the capacity.
      expr: (longhorn_disk_usage_bytes / longhorn_disk_capacity_bytes) * 100 > 70
      for: 5m
      labels:
        issue: The used storage of disk {{$labels.disk}} on node {{$labels.node}} is high.
        severity: warning
    - alert: LonghornNodeDown
      annotations:
        description: There are {{$value}} Longhorn nodes which have been offline for more than 5 minutes.
        summary: Longhorn nodes is offline
      expr: (avg(longhorn_node_count_total) or on() vector(0)) - (count(longhorn_node_status{condition="ready"} == 1) or on() vector(0)) > 0
      for: 5m
      labels:
        issue: There are {{$value}} Longhorn nodes are offline
        severity: critical
    - alert: LonghornInstanceManagerCPUUsageWarning
      annotations:
        description: Longhorn instance manager {{$labels.instance_manager}} on {{$labels.node}} has CPU Usage / CPU request is {{$value}}% for
          more than 5 minutes.
        summary: Longhorn instance manager {{$labels.instance_manager}} on {{$labels.node}} has CPU Usage / CPU request is over 300%.
      expr: (longhorn_instance_manager_cpu_usage_millicpu/longhorn_instance_manager_cpu_requests_millicpu) * 100 > 300
      for: 5m
      labels:
        issue: Longhorn instance manager {{$labels.instance_manager}} on {{$labels.node}} consumes 3 times the CPU request.
        severity: warning
    - alert: LonghornNodeCPUUsageWarning
      annotations:
        description: Longhorn node {{$labels.node}} has CPU Usage / CPU capacity is {{$value}}% for
          more than 5 minutes.
        summary: Longhorn node {{$labels.node}} experiences high CPU pressure for more than 5m.
      expr: (longhorn_node_cpu_usage_millicpu / longhorn_node_cpu_capacity_millicpu) * 100 > 90
      for: 5m
      labels:
        issue: Longhorn node {{$labels.node}} experiences high CPU pressure.
        severity: warning
```

See more about how to define alert rules at [here](https://prometheus.io/docs/prometheus/latest/configuration/alerting_rules/#alerting-rules).


---

## Article: monitoring/integrating-with-rancher-monitoring.md

---
title: Integrating Longhorn metrics into the Rancher monitoring system
weight: 2
---
## About the Rancher Monitoring System

Using Rancher, you can monitor the state and processes of your cluster nodes, Kubernetes components, and software deployments through integration with [Prometheus](https://prometheus.io/), a leading open-source monitoring solution.

See [here](https://rancher.com/docs/rancher/v2.x/en/monitoring-alerting/) for the instruction about how to deploy/enable the Rancher monitoring system.

## Add Longhorn Metrics to the Rancher Monitoring System

If you are using Rancher to manage your Kubernetes and already enabled Rancher monitoring, you can add Longhorn metrics to Rancher monitoring by simply deploying the following ServiceMonitor:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: longhorn-prometheus-servicemonitor
  namespace: longhorn-system
  labels:
    name: longhorn-prometheus-servicemonitor
spec:
  selector:
    matchLabels:
      app: longhorn-manager
  namespaceSelector:
    matchNames:
    - longhorn-system
  endpoints:
  - port: manager
```

Once the ServiceMonitor is created, Rancher will automatically discover all Longhorn metrics.

You can then set up a Grafana dashboard for visualization.

You can import our prebuilt [Longhorn example dashboard](https://grafana.com/grafana/dashboards/13032) to have an idea.

You can also set up alerts in Rancher UI.


---

## Article: monitoring/kubelet-volume-metrics.md

---
title: Kubelet Volume Metrics Support
weight: 4
---

## About Kubelet Volume Metrics

Kubelet exposes [the following metrics](https://github.com/kubernetes/kubernetes/blob/4b24dca228d61f4d13dcd57b46465b0df74571f6/pkg/kubelet/metrics/collectors/volume_stats.go#L27):

1. kubelet_volume_stats_capacity_bytes
1. kubelet_volume_stats_available_bytes
1. kubelet_volume_stats_used_bytes
1. kubelet_volume_stats_inodes
1. kubelet_volume_stats_inodes_free
1. kubelet_volume_stats_inodes_used

Those metrics measure information related to a PVC's filesystem inside a Longhorn block device.

They are different than [longhorn_volume_*](../metrics) metrics, which measure information specific to a Longhorn block device.

You can set up a monitoring system that scrapes Kubelet metric endpoints to obtains a PVC's status and set up alerts for abnormal events, such as the PVC being about to run out of storage space.

A popular monitoring setup is [prometheus-operator/kube-prometheus-stack,](https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack) which scrapes `kubelet_volume_stats_*` metrics and provides a dashboard and alert rules for them.

## Longhorn CSI Plugin Support

In v1.1.0, Longhorn CSI plugin supports the `NodeGetVolumeStats` RPC according to the [CSI spec](https://github.com/container-storage-interface/spec/blob/master/spec.md#nodegetvolumestats).

This allows the kubelet to query the Longhorn CSI plugin for a PVC's status.

The kubelet then exposes that information in `kubelet_volume_stats_*` metrics.


---

## Article: monitoring/metrics.md

---
title: Longhorn Metrics for Monitoring
weight: 3
---
## Volume

| Name | Description  | Example |
|---|---|---|
| longhorn_volume_actual_size_bytes | Actual space used by each replica of the volume on the corresponding node | longhorn_volume_actual_size_bytes{pvc_namespace="default",node="worker-2",pvc="testvol",volume="testvol"} 1.1917312e+08 |
| longhorn_volume_capacity_bytes | Configured size in bytes for this volume | longhorn_volume_capacity_bytes{pvc_namespace="default",node="worker-2",pvc="testvol",volume="testvol"} 6.442450944e+09 |
| longhorn_volume_state | State of this volume: 1=creating, 2=attached, 3=Detached, 4=Attaching, 5=Detaching, 6=Deleting | longhorn_volume_state{pvc_namespace="default",node="worker-2",pvc="testvol",volume="testvol"} 2 |
| longhorn_volume_robustness | Robustness of this volume: 0=unknown, 1=healthy, 2=degraded, 3=faulted  | longhorn_volume_robustness{pvc_namespace="default",node="worker-2",pvc="testvol",volume="testvol"} 1 |
| longhorn_volume_read_throughput | Read throughput of this volume (Bytes/s) | longhorn_volume_read_throughput{pvc_namespace="default",node="worker-2",pvc="testvol",volume="testvol"} 5120000 |
| longhorn_volume_write_throughput | Write throughput of this volume (Bytes/s) | longhorn_volume_write_throughput{pvc_namespace="default",node="worker-2",pvc="testvol",volume="testvol"} 512000 |
| longhorn_volume_read_iops | Read IOPS of this volume | longhorn_volume_read_iops{pvc_namespace="default",node="worker-2",pvc="testvol",volume="testvol"} 100 |
| longhorn_volume_write_iops | Write IOPS of this volume | longhorn_volume_write_iops{pvc_namespace="default",node="worker-2",pvc="testvol",volume="testvol"} 100 |
| longhorn_volume_read_latency | Read latency of this volume (ns) | longhorn_volume_read_latency{pvc_namespace="default",node="worker-2",pvc="testvol",volume="testvol"} 100000 |
| longhorn_volume_write_latency | Write latency of this volume (ns) | longhorn_volume_write_latency{pvc_namespace="default",node="worker-2",pvc="testvol",volume="testvol"} 100000 |
| longhorn_volume_file_system_read_only | This metric indicates that the volume is now in read-only mode. The metric is either 1 or no record for each volume | longhorn_volume_file_system_read_only{node="worker-2",pvc="testvol",pvc_namespace="default",volume="testvol"} 1

## Node

| Name | Description  | Example |
|---|---|---|
| longhorn_node_status | Status of this node: 1=true, 0=false | longhorn_node_status{condition="ready",condition_reason="",node="worker-2"} 1 |
| longhorn_node_count_total | Total number of nodes in the Longhorn system | longhorn_node_count_total 4 |
| longhorn_node_cpu_capacity_millicpu | The maximum allocatable CPU on this node | longhorn_node_cpu_capacity_millicpu{node="worker-2"} 2000 |
| longhorn_node_cpu_usage_millicpu | The CPU usage on this node | longhorn_node_cpu_usage_millicpu{node="pworker-2"} 186 |
| longhorn_node_memory_capacity_bytes | The maximum allocatable memory on this node | longhorn_node_memory_capacity_bytes{node="worker-2"} 4.031229952e+09 |
| longhorn_node_memory_usage_bytes |  The memory usage on this node | longhorn_node_memory_usage_bytes{node="worker-2"} 1.833582592e+09 |
| longhorn_node_storage_capacity_bytes | The storage capacity of this node | longhorn_node_storage_capacity_bytes{node="worker-3"} 8.3987283968e+10 |
| longhorn_node_storage_usage_bytes | The used storage of this node | longhorn_node_storage_usage_bytes{node="worker-3"} 9.060941824e+09 |
| longhorn_node_storage_reservation_bytes | The reserved storage for other applications and system on this node | longhorn_node_storage_reservation_bytes{node="worker-3"} 2.519618519e+10 |

## Replica

| Name | Description | Example |
|---|---|---|
| longhorn_replica_info | Static metadata for each Replica CR | longhorn_replica_info{replica="testvol-r-abc", volume="testvol", node="node-1", disk_path="/dev/xda", data_engine="v2"} 1 |
| longhorn_replica_state | Current runtime state of the replica: running, stopped, error, starting, stopping, unknown | longhorn_replica_state{replica="testvol-r-abc", volume="testvol", node="node-1", state="running"} 1 |

## Engine

| Name | Description | Example |
|---|---|---|
| longhorn_engine_info | Static metadata for each Engine CR | longhorn_engine_info{engine="testvol-e-0", volume="testvol", node="node-1", data_engine="v2", frontend="blockdev", image="longhorn-instance-manager:latest"} 1 |
| longhorn_engine_state | Runtime state of an engine: running, stopped, error, starting, stopping, unknown | longhorn_engine_state{engine="testvol-e-0", volume="testvol", node="node-1", state="running"} 1 |
| longhorn_engine_replica_mode | The mode reported for each replica by the engine: RW, WO, ERR | longhorn_engine_replica_mode{volume="testvol", engine="testvol-e-0", replica="testvol-r-abc", mode="RW"} 1 |
| longhorn_engine_rebuild_progress | Engine rebuild progress (0–100%). Visible only during replica rebuilding. | longhorn_engine_rebuild_progress{pvc_namespace="default",pvc="testvol",engine="testvol-e-0",rebuild_src="10.42.1.215:20036",rebuild_dst="10.42.0.131:20922"} 42 |

## Disk

| Name | Description  | Example |
|---|---|---|
| longhorn_disk_capacity_bytes | The storage capacity of this disk | longhorn_disk_capacity_bytes{disk="default-disk-8b28ee3134628183",node="worker-3"} 8.3987283968e+10 |
| longhorn_disk_usage_bytes | The used storage of this disk | longhorn_disk_usage_bytes{disk="default-disk-8b28ee3134628183",node="worker-3"} 9.060941824e+09 |
| longhorn_disk_reservation_bytes | The reserved storage for other applications and system on this disk | longhorn_disk_reservation_bytes{disk="default-disk-8b28ee3134628183",node="worker-3"} 2.519618519e+10 |
| longhorn_disk_status | The status of this disk | longhorn_disk_status{condition="ready",condition_reason="",disk="default-disk-ca0300000000",node="worker-3"} |
| longhorn_disk_read_throughput | Read throughput of this disk (Bytes/s) | longhorn_disk_read_throughput{disk="default-disk-8b28ee3134628183",node="worker-3",disk_path="/dev/sda"} 10485760 |
| longhorn_disk_write_throughput | Write throughput of this disk (Bytes/s) | longhorn_disk_write_throughput{disk="default-disk-8b28ee3134628183",node="worker-3",disk_path="/dev/sda"} 2097152 |
| longhorn_disk_read_iops | Read IOPS of this disk | longhorn_disk_read_iops{disk="default-disk-8b28ee3134628183",node="worker-3",disk_path="/dev/sda"} 200 |
| longhorn_disk_write_iops | Write IOPS of this disk | longhorn_disk_write_iops{disk="default-disk-8b28ee3134628183",node="worker-3",disk_path="/dev/sda"} 150 |
| longhorn_disk_read_latency | Read latency of this disk (nanoseconds) | longhorn_disk_read_latency{disk="default-disk-8b28ee3134628183",node="worker-3",disk_path="/dev/sda"} 85000 |
| longhorn_disk_write_latency | Write latency of this disk (nanoseconds) | longhorn_disk_write_latency{disk="default-disk-8b28ee3134628183",node="worker-3",disk_path="/dev/sda"} 95000 |


## Instance Manager

| Name | Description  | Example |
|---|---|---|
| longhorn_instance_manager_cpu_usage_millicpu |  The cpu usage of this longhorn instance manager | longhorn_instance_manager_cpu_usage_millicpu{instance_manager="instance-manager-e-2189ed13",instance_manager_type="engine",node="worker-2"} 80 |
| longhorn_instance_manager_cpu_requests_millicpu | Requested CPU resources in kubernetes of this Longhorn instance manager | longhorn_instance_manager_cpu_requests_millicpu{instance_manager="instance-manager-e-2189ed13",instance_manager_type="engine",node="worker-2"} 250 |
| longhorn_instance_manager_memory_usage_bytes | The memory usage of this longhorn instance manager | longhorn_instance_manager_memory_usage_bytes{instance_manager="instance-manager-e-2189ed13",instance_manager_type="engine",node="worker-2"} 2.4072192e+07 |
| longhorn_instance_manager_memory_requests_bytes | Requested memory in Kubernetes of this longhorn instance manager | longhorn_instance_manager_memory_requests_bytes{instance_manager="instance-manager-e-2189ed13",instance_manager_type="engine",node="worker-2"} 0 |
| longhorn_instance_manager_proxy_grpc_connection | The number of proxy gRPC connection of this longhorn instance manager | longhorn_instance_manager_proxy_grpc_connection{instance_manager="instance-manager-e-814dfd05", instance_manager_type="engine", node="worker-2"} 0

## Manager

| Name | Description  | Example |
|---|---|---|
| longhorn_manager_cpu_usage_millicpu |  The CPU usage of this Longhorn Manager | longhorn_manager_cpu_usage_millicpu{manager="longhorn-manager-5rx2n",node="worker-2"} 27 |
| longhorn_manager_memory_usage_bytes | The memory usage of this Longhorn Manager | longhorn_manager_memory_usage_bytes{manager="longhorn-manager-5rx2n",node="worker-2"} 2.6144768e+07|

## Backup

| Name | Description  | Example |
|---|---|---|
| longhorn_backup_actual_size_bytes | Actual size of this backup | longhorn_backup_actual_size_bytes{backup="backup-4ab66eca0d60473e",volume="testvol", recurring_job="backup"} 6.291456e+07 |
| longhorn_backup_state | State of this backup: 0=New, 1=Pending, 2=InProgress, 3=Completed, 4=Error, 5=Unknown | longhorn_backup_state{backup="backup-4ab66eca0d60473e",volume="testvol", recurring_job=""} 3 |

## Snapshot

| Name | Description  | Example |
|---|---|---|
| longhorn_snapshot_actual_size_bytes | Actual size of this snapshot | longhorn_snapshot_actual_size_bytes{snapshot="f4468111-2efa-45f5-aef6-63109e30d92c",user_created="false",volume="testvol"} 1.048576e+07 |


## BackingImage

| Name | Description  | Example |
|---|---|---|
| longhorn_backing_image_actual_size_bytes | Actual size of this backing image | longhorn_backing_image_actual_size_bytes{backing_image="parrot",disk="ca203ce8-2cad-4cd1-92a7-542851f50518",node="kworker1"} 3.3554432e+07 |
| longhorn_backing_image_state | State of this backing image: 0=Pending, 1=Starting, 2=InProgress, 3=ReadyForTransfer, 4=Ready, 5=Failed, 6=FailedAndCleanUp, 7=Unknown | longhorn_backing_image_state{backing_image="parrot",disk="ca203ce8-2cad-4cd1-92a7-542851f50518",node="kworker1"} 4 |

## BackupBackingImage

| Name | Description  | Example |
|---|---|---|
| longhorn_backup_backing_image_actual_size_bytes | Actual size of this backup backing image | longhorn_backup_backing_image_actual_size_bytes{backup_backing_image="parrot"} 3.3554432e+07 |
| longhorn_backup_backing_image_state | State of this backup backing image: 0=New, 1=Pending, 2=InProgress, 3=Completed, 4=Error, 5=Unknown | longhorn_backup_backing_image_state{backup_backing_image="parrot"} 3 |

## CSI

The CSI sidecar component has built-in metrics for users to get insights into CSI operations. The CSI operations metrics cover total count, error count, and call latency. Longhorn enables the metrics by adding the flag `--http-endpoint` for each CSI sidecar component. You can use [Prometheus's PodMonitor](https://github.com/prometheus-operator/prometheus-operator/blob/main/Documentation/api.md#podmonitor) to collect these metrics. 

| Name | Port |
|---|---|
| longhorn-csi-attacher | 8000 | 
| longhorn-csi-provisioner | 8000 |
| longhorn-csi-resizer | 8000 |
| longhorn-csi-snapshotter | 8000 |

The metrics provided by the CSI sidecar component are provided in a histogram format. For example, you can obtain metrics observing the time it takes to create a Longhorn Volume for the PVC.

```
csi_sidecar_operations_seconds_bucket{driver_name="driver.longhorn.io",grpc_status_code="OK",method_name="/csi.v1.Controller/ControllerPublishVolume",le="0.1"} 0
csi_sidecar_operations_seconds_bucket{driver_name="driver.longhorn.io",grpc_status_code="OK",method_name="/csi.v1.Controller/ControllerPublishVolume",le="0.25"} 0
csi_sidecar_operations_seconds_bucket{driver_name="driver.longhorn.io",grpc_status_code="OK",method_name="/csi.v1.Controller/ControllerPublishVolume",le="0.5"} 0
csi_sidecar_operations_seconds_bucket{driver_name="driver.longhorn.io",grpc_status_code="OK",method_name="/csi.v1.Controller/ControllerPublishVolume",le="1"} 0
csi_sidecar_operations_seconds_bucket{driver_name="driver.longhorn.io",grpc_status_code="OK",method_name="/csi.v1.Controller/ControllerPublishVolume",le="2.5"} 3
csi_sidecar_operations_seconds_bucket{driver_name="driver.longhorn.io",grpc_status_code="OK",method_name="/csi.v1.Controller/ControllerPublishVolume",le="5"} 3
csi_sidecar_operations_seconds_bucket{driver_name="driver.longhorn.io",grpc_status_code="OK",method_name="/csi.v1.Controller/ControllerPublishVolume",le="10"} 3
csi_sidecar_operations_seconds_bucket{driver_name="driver.longhorn.io",grpc_status_code="OK",method_name="/csi.v1.Controller/ControllerPublishVolume",le="15"} 9
csi_sidecar_operations_seconds_bucket{driver_name="driver.longhorn.io",grpc_status_code="OK",method_name="/csi.v1.Controller/ControllerPublishVolume",le="25"} 9
csi_sidecar_operations_seconds_bucket{driver_name="driver.longhorn.io",grpc_status_code="OK",method_name="/csi.v1.Controller/ControllerPublishVolume",le="50"} 9
csi_sidecar_operations_seconds_bucket{driver_name="driver.longhorn.io",grpc_status_code="OK",method_name="/csi.v1.Controller/ControllerPublishVolume",le="120"} 9
csi_sidecar_operations_seconds_bucket{driver_name="driver.longhorn.io",grpc_status_code="OK",method_name="/csi.v1.Controller/ControllerPublishVolume",le="300"} 9
csi_sidecar_operations_seconds_bucket{driver_name="driver.longhorn.io",grpc_status_code="OK",method_name="/csi.v1.Controller/ControllerPublishVolume",le="600"} 9
csi_sidecar_operations_seconds_bucket{driver_name="driver.longhorn.io",grpc_status_code="OK",method_name="/csi.v1.Controller/ControllerPublishVolume",le="+Inf"} 9
csi_sidecar_operations_seconds_sum{driver_name="driver.longhorn.io",grpc_status_code="OK",method_name="/csi.v1.Controller/ControllerPublishVolume"} 66.816478825
csi_sidecar_operations_seconds_count{driver_name="driver.longhorn.io",grpc_status_code="OK",method_name="/csi.v1.Controller/ControllerPublishVolume"} 9
```


---

## Article: monitoring/prometheus-and-grafana-setup.md

---
title: Setting up Prometheus and Grafana to monitor Longhorn
weight: 1
---

This document is a quick guide to setting up the monitor for Longhorn.

Longhorn natively exposes metrics in [Prometheus text format](https://prometheus.io/docs/instrumenting/exposition_formats/#text-based-format) on a REST endpoint `http://LONGHORN_MANAGER_IP:PORT/metrics`.

You can use any collecting tools such as [Prometheus](https://prometheus.io/), [Graphite](https://graphiteapp.org/), [Telegraf](https://www.influxdata.com/time-series-platform/telegraf/) to scrape these metrics then visualize the collected data by tools such as [Grafana](https://grafana.com/).

See [Longhorn Metrics for Monitoring](../metrics) for available metrics.

## High-level Overview

The monitoring system uses `Prometheus` for collecting data and alerting, and `Grafana` for visualizing/dashboarding the collected data.

* Prometheus server which scrapes and stores time-series data from Longhorn metrics endpoints. The Prometheus is also responsible for generating alerts based on configured rules and collected data. Prometheus servers then send alerts to an Alertmanager.
* AlertManager then manages those alerts, including silencing, inhibition, aggregation, and sending out notifications via methods such as email, on-call notification systems, and chat platforms.
* Grafana which queries Prometheus server for data and draws a dashboard for visualization.

The below picture describes the detailed architecture of the monitoring system.

![images](/img/screenshots/monitoring/longhorn-monitoring-system.png)

There are 2 unmentioned components in the above picture:

* Longhorn Backend service is a service pointing to the set of Longhorn manager pods. Longhorn's metrics are exposed in Longhorn manager pods at the endpoint `http://LONGHORN_MANAGER_IP:PORT/metrics`.
* [Prometheus operator](https://github.com/prometheus-operator/prometheus-operator/blob/master/Documentation/user-guides/getting-started.md) makes running Prometheus on top of Kubernetes very easy. The operator watches 3 custom resources: ServiceMonitor, Prometheus ,and AlertManager.
  When you create those custom resources, Prometheus Operator deploys and manages the Prometheus server, AlertManager with the user-specified configurations.

## Installation

This document uses the `default` namespace for the monitoring system. To install on a different namespace, change the field `namespace: <OTHER_NAMESPACE>` in manifests.

### Install Prometheus Operator
Follow instructions in [Prometheus Operator - Quickstart](https://github.com/prometheus-operator/prometheus-operator#quickstart).

> **NOTE:** You may need to choose a release that is compatible with the Kubernetes version of the cluster.

### Install Longhorn ServiceMonitor

#### Install Longhorn ServiceMonitor with Kubectl

Create a ServiceMonitor for Longhorn Manager.

    ```yaml
    apiVersion: monitoring.coreos.com/v1
    kind: ServiceMonitor
    metadata:
      name: longhorn-prometheus-servicemonitor
      namespace: default
      labels:
        name: longhorn-prometheus-servicemonitor
    spec:
      selector:
        matchLabels:
          app: longhorn-manager
      namespaceSelector:
        matchNames:
        - longhorn-system
      endpoints:
      - port: manager
    ```

#### Install Longhorn ServiceMonitor with Helm

1. Modify the YAML file `longhorn/chart/values.yaml`.

    ```yaml
    metrics:
      serviceMonitor:
        # -- Setting that allows the creation of a [Prometheus Operator](https://prometheus-operator.dev/) ServiceMonitor resource for Longhorn Manager components.
        enabled: true
    ```

1. Create a ServiceMonitor for Longhorn Manager using Helm.

    ```bash
      helm upgrade longhorn longhorn/longhorn --namespace longhorn-system -f values.yaml
    ```

Longhorn ServiceMonitor is a [Prometheus Operator](https://prometheus-operator.dev/) custom resource. This setup allows the Prometheus server to discover all Longhorn Manager pods and their respective endpoints.

You can use the label selector `app: longhorn-manager` to select the longhorn-backend service, which points to the set of Longhorn Manager pods.

### Install and configure Prometheus AlertManager

1. Create a highly available Alertmanager deployment with 3 instances.

    ```yaml
    apiVersion: monitoring.coreos.com/v1
    kind: Alertmanager
    metadata:
      name: longhorn
      namespace: default
    spec:
      replicas: 3
    ```

1. The Alertmanager instances will not start unless a valid configuration is given.
See [Prometheus - Configuration](https://prometheus.io/docs/alerting/latest/configuration/) for more explanation.

    ```yaml
    global:
      resolve_timeout: 5m
    route:
      group_by: [alertname]
      receiver: email_and_slack
    receivers:
    - name: email_and_slack
      email_configs:
      - to: <the email address to send notifications to>
        from: <the sender address>
        smarthost: <the SMTP host through which emails are sent>
        # SMTP authentication information.
        auth_username: <the username>
        auth_identity: <the identity>
        auth_password: <the password>
        headers:
          subject: 'Longhorn-Alert'
        text: |-
          {{ range .Alerts }}
            *Alert:* {{ .Annotations.summary }} - `{{ .Labels.severity }}`
            *Description:* {{ .Annotations.description }}
            *Details:*
            {{ range .Labels.SortedPairs }} • *{{ .Name }}:* `{{ .Value }}`
            {{ end }}
          {{ end }}
      slack_configs:
      - api_url: <the Slack webhook URL>
        channel: <the channel or user to send notifications to>
        text: |-
          {{ range .Alerts }}
            *Alert:* {{ .Annotations.summary }} - `{{ .Labels.severity }}`
            *Description:* {{ .Annotations.description }}
            *Details:*
            {{ range .Labels.SortedPairs }} • *{{ .Name }}:* `{{ .Value }}`
            {{ end }}
          {{ end }}
    ```

    Save the above Alertmanager config in a file called `alertmanager.yaml` and create a secret from it using kubectl.

    Alertmanager instances require the secret resource naming to follow the format `alertmanager-<ALERTMANAGER_NAME>`. In the previous step, the name of the Alertmanager is `longhorn`, so the secret name must be `alertmanager-longhorn`

    ```
    $ kubectl create secret generic alertmanager-longhorn --from-file=alertmanager.yaml -n default
    ```

1. To be able to view the web UI of the Alertmanager, expose it through a Service. A simple way to do this is to use a Service of type NodePort.

    ```yaml
    apiVersion: v1
    kind: Service
    metadata:
      name: alertmanager-longhorn
      namespace: default
    spec:
      type: NodePort
      ports:
      - name: web
        nodePort: 30903
        port: 9093
        protocol: TCP
        targetPort: web
      selector:
        alertmanager: longhorn
    ```

    After creating the above service, you can access the web UI of Alertmanager via a Node's IP and the port 30903.

    > Use the above `NodePort` service for quick verification only because it doesn't communicate over the TLS connection. You may want to change the service type to `ClusterIP` and set up an Ingress-controller to expose the web UI of Alertmanager over a TLS connection.

### Install and configure Prometheus server

1. Create PrometheusRule custom resource to define alert conditions. See more examples about Longhorn alert rules at [Longhorn Alert Rule Examples](../alert-rules-example).

    ```yaml
    apiVersion: monitoring.coreos.com/v1
    kind: PrometheusRule
    metadata:
      labels:
        prometheus: longhorn
        role: alert-rules
      name: prometheus-longhorn-rules
      namespace: default
    spec:
      groups:
      - name: longhorn.rules
        rules:
        - alert: LonghornVolumeUsageCritical
          annotations:
            description: Longhorn volume {{$labels.volume}} on {{$labels.node}} is at {{$value}}% used for
              more than 5 minutes.
            summary: Longhorn volume capacity is over 90% used.
          expr: 100 * (longhorn_volume_usage_bytes / longhorn_volume_capacity_bytes) > 90
          for: 5m
          labels:
            issue: Longhorn volume {{$labels.volume}} usage on {{$labels.node}} is critical.
            severity: critical
    ```
   See [Prometheus - Alerting rules](https://prometheus.io/docs/prometheus/latest/configuration/alerting_rules/#alerting-rules) for more information.

1. If [RBAC](https://kubernetes.io/docs/reference/access-authn-authz/authorization/) authorization is activated, Create a ClusterRole and ClusterRoleBinding for the Prometheus Pods.

    ```yaml
    apiVersion: v1
    kind: ServiceAccount
    metadata:
      name: prometheus
      namespace: default
    ```

    ```yaml
    apiVersion: rbac.authorization.k8s.io/v1
    kind: ClusterRole
    metadata:
      name: prometheus
      namespace: default
    rules:
    - apiGroups: [""]
      resources:
      - nodes
      - services
      - endpoints
      - pods
      verbs: ["get", "list", "watch"]
    - apiGroups: [""]
      resources:
      - configmaps
      verbs: ["get"]
    - nonResourceURLs: ["/metrics"]
      verbs: ["get"]
    ```

    ```yaml
    apiVersion: rbac.authorization.k8s.io/v1
    kind: ClusterRoleBinding
    metadata:
      name: prometheus
    roleRef:
      apiGroup: rbac.authorization.k8s.io
      kind: ClusterRole
      name: prometheus
    subjects:
    - kind: ServiceAccount
      name: prometheus
      namespace: default
    ```

1. Create a Prometheus custom resource. Notice that we select the Longhorn service monitor and Longhorn rules in the spec.

    ```yaml
    apiVersion: monitoring.coreos.com/v1
    kind: Prometheus
    metadata:
      name: longhorn
      namespace: default
    spec:
      replicas: 2
      serviceAccountName: prometheus
      alerting:
        alertmanagers:
          - namespace: default
            name: alertmanager-longhorn
            port: web
      serviceMonitorSelector:
        matchLabels:
          name: longhorn-prometheus-servicemonitor
      ruleSelector:
        matchLabels:
          prometheus: longhorn
          role: alert-rules
    ```

1. To be able to view the web UI of the Prometheus server, expose it through a Service. A simple way to do this is to use a Service of type NodePort.

    ```yaml
    apiVersion: v1
    kind: Service
    metadata:
      name: prometheus-longhorn
      namespace: default
    spec:
      type: NodePort
      ports:
      - name: web
        nodePort: 30904
        port: 9090
        protocol: TCP
        targetPort: web
      selector:
        prometheus: longhorn
    ```

    After creating the above service, you can access the web UI of the Prometheus server via a Node's IP and the port 30904.

    > At this point, you should be able to see all Longhorn manager targets as well as Longhorn rules in the targets and rules section of the Prometheus server UI.

    > Use the above NodePort service for quick verification only because it doesn't communicate over the TLS connection. You may want to change the service type to `ClusterIP` and set up an Ingress controller to expose the web UI of the Prometheus server over a TLS connection.

### Setup Grafana

1. Create Grafana datasource ConfigMap.

    ```yaml
    apiVersion: v1
    kind: ConfigMap
    metadata:
      name: grafana-datasources
      namespace: default
    data:
      prometheus.yaml: |-
        {
            "apiVersion": 1,
            "datasources": [
                {
                   "access":"proxy",
                    "editable": true,
                    "name": "prometheus-longhorn",
                    "orgId": 1,
                    "type": "prometheus",
                    "url": "http://prometheus-longhorn.default.svc:9090",
                    "version": 1
                }
            ]
        }
    ```

    > **NOTE:** change field `url` if you are installing the monitoring stack in a different namespace.
    > `http://prometheus-longhorn.<NAMESPACE>.svc:9090"`

1. Create Grafana Deployment.
    ```yaml
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: grafana
      namespace: default
      labels:
        app: grafana
    spec:
      replicas: 1
      selector:
        matchLabels:
          app: grafana
      template:
        metadata:
          name: grafana
          labels:
            app: grafana
        spec:
          containers:
          - name: grafana
            image: grafana/grafana:7.1.5
            ports:
            - name: grafana
              containerPort: 3000
            resources:
              limits:
                memory: "500Mi"
                cpu: "300m"
              requests:
                memory: "500Mi"
                cpu: "200m"
            volumeMounts:
              - mountPath: /var/lib/grafana
                name: grafana-storage
              - mountPath: /etc/grafana/provisioning/datasources
                name: grafana-datasources
                readOnly: false
          volumes:
            - name: grafana-storage
              emptyDir: {}
            - name: grafana-datasources
              configMap:
                  defaultMode: 420
                  name: grafana-datasources
    ```

1. Create Grafana Service.
    ```yaml
    apiVersion: v1
    kind: Service
    metadata:
      name: grafana
      namespace: default
    spec:
      selector:
        app: grafana
      type: ClusterIP
      ports:
        - port: 3000
          targetPort: 3000
    ```

1. Expose Grafana on NodePort `32000`.
    ```yaml
    kubectl -n default patch svc grafana --type='json' -p '[{"op":"replace","path":"/spec/type","value":"NodePort"},{"op":"replace","path":"/spec/ports/0/nodePort","value":32000}]'
    ```

   > Use the above NodePort service for quick verification only because it doesn't communicate over the TLS connection. You may want to change the service type to ClusterIP and set up an Ingress controller to expose Grafana over a TLS connection.

1. Access the Grafana dashboard using any node IP on port `32000`.
    ```
    # Default Credential
    User: admin
    Pass: admin
    ```

1. Setup Longhorn dashboard.

    Once inside Grafana, import the prebuilt [Longhorn example dashboard](https://grafana.com/grafana/dashboards/17626).

    See [Grafana Lab - Export and import](https://grafana.com/docs/grafana/latest/reference/export_import/) for instructions on how to import a Grafana dashboard.

    You should see the following dashboard at successful setup:
    ![images](/img/screenshots/monitoring/longhorn-example-grafana-dashboard.png)


---

## Article: deploy/_index.md

---
title: Installation and Setup
weight: 2
---

---

## Article: deploy/accessing-the-ui/_index.md

---
title: Accessing the UI
weight: 2
---

## Prerequisites for Access and Authentication

These instructions assume that Longhorn is installed.

If you installed Longhorn YAML manifest, you'll need to set up an Ingress controller to allow external traffic into the cluster, and authentication will not be enabled by default. This applies to Helm and kubectl installations. For information on creating an NGINX Ingress controller with basic authentication, refer to [this section.](./longhorn-ingress)

If Longhorn was installed as a Rancher catalog app, Rancher automatically created an Ingress controller for you with access control (the rancher-proxy).

## Accessing the Longhorn UI

Once Longhorn has been installed in your Kubernetes cluster, you can access the UI dashboard.

1. Get the Longhorn's external service IP:

    ```shell
    kubectl -n longhorn-system get svc
    ```

    For Longhorn v0.8.0, the output should look like this, and the `CLUSTER-IP` of the `longhorn-frontend` is used to access the Longhorn UI:

    ```shell
    NAME                TYPE           CLUSTER-IP      EXTERNAL-IP      PORT(S)        AGE
    longhorn-backend    ClusterIP      10.20.248.250   <none>           9500/TCP       58m
    longhorn-frontend   ClusterIP      10.20.245.110   <none>           80/TCP         58m

    ```

    In the example above, the IP is `10.20.245.110`.
    
    > For Longhorn v0.8.0+, UI service type changed from `LoadBalancer` to `ClusterIP.`

2. Navigate to the IP of `longhorn-frontend` in your browser.

    The Longhorn UI looks like this:

    {{< figure src="/img/screenshots/getting-started/v1.10.0/longhorn-ui.png" >}}


---

## Article: deploy/accessing-the-ui/longhorn-ingress.md

---
  title:  Create an Ingress with Basic Authentication (nginx)
  weight: 1
---

If you install Longhorn on a Kubernetes cluster with kubectl or Helm, you will need to create an Ingress to allow external traffic to reach the Longhorn UI.

Authentication is not enabled by default for kubectl and Helm installations. In these steps, you'll learn how to create an Ingress with basic authentication using annotations for the nginx ingress controller.

1. Create a basic auth file `auth`. It's important the file generated is named auth (actually - that the secret has a key `data.auth`), otherwise the Ingress returns a 503.
    ```
    $ USER=<USERNAME_HERE>; PASSWORD=<PASSWORD_HERE>; echo "${USER}:$(openssl passwd -stdin -apr1 <<< ${PASSWORD})" >> auth
    ```
2. Create a secret:
    ```
    $ kubectl -n longhorn-system create secret generic basic-auth --from-file=auth
    ```
3. Create an Ingress manifest `longhorn-ingress.yml` :
    > Since v1.2.0, Longhorn supports uploading backing image from the UI, so please specify `nginx.ingress.kubernetes.io/proxy-body-size: 10000m` as below to ensure uploading images work as expected.

    ```
    apiVersion: networking.k8s.io/v1
    kind: Ingress
    metadata:
      name: longhorn-ingress
      namespace: longhorn-system
      annotations:
        # type of authentication
        nginx.ingress.kubernetes.io/auth-type: basic
        # prevent the controller from redirecting (308) to HTTPS
        nginx.ingress.kubernetes.io/ssl-redirect: 'false'
        # name of the secret that contains the user/password definitions
        nginx.ingress.kubernetes.io/auth-secret: basic-auth
        # message to display with an appropriate context why the authentication is required
        nginx.ingress.kubernetes.io/auth-realm: 'Authentication Required '
        # custom max body size for file uploading like backing image uploading
        nginx.ingress.kubernetes.io/proxy-body-size: 10000m
    spec:
      ingressClassName: nginx
      rules:
      - http:
          paths:
          - pathType: Prefix
            path: "/"
            backend:
              service:
                name: longhorn-frontend
                port:
                  number: 80
    ```
4. Create the Ingress:
    ```
    $ kubectl -n longhorn-system apply -f longhorn-ingress.yml
    ```

e.g.:
```
$ USER=foo; PASSWORD=bar; echo "${USER}:$(openssl passwd -stdin -apr1 <<< ${PASSWORD})" >> auth
$ cat auth
foo:$apr1$FnyKCYKb$6IP2C45fZxMcoLwkOwf7k0

$ kubectl -n longhorn-system create secret generic basic-auth --from-file=auth
secret/basic-auth created
$ kubectl -n longhorn-system get secret basic-auth -o yaml
apiVersion: v1
data:
  auth: Zm9vOiRhcHIxJEZueUtDWUtiJDZJUDJDNDVmWnhNY29Md2tPd2Y3azAK
kind: Secret
metadata:
  creationTimestamp: "2020-05-29T10:10:16Z"
  name: basic-auth
  namespace: longhorn-system
  resourceVersion: "2168509"
  selfLink: /api/v1/namespaces/longhorn-system/secrets/basic-auth
  uid: 9f66233f-b12f-4204-9c9d-5bcaca794bb7
type: Opaque

$ echo "
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: longhorn-ingress
  namespace: longhorn-system
  annotations:
    # type of authentication
    nginx.ingress.kubernetes.io/auth-type: basic
    # prevent the controller from redirecting (308) to HTTPS
    nginx.ingress.kubernetes.io/ssl-redirect: 'false'
    # name of the secret that contains the user/password definitions
    nginx.ingress.kubernetes.io/auth-secret: basic-auth
    # message to display with an appropriate context why the authentication is required
    nginx.ingress.kubernetes.io/auth-realm: 'Authentication Required '
spec:
  rules:
  - http:
      paths:
      - pathType: Prefix
        path: "/"
        backend:
          service:
            name: longhorn-frontend
            port:
              number: 80
" | kubectl -n longhorn-system create -f -
ingress.networking.k8s.io/longhorn-ingress created

$ kubectl -n longhorn-system get ingress
NAME               HOSTS   ADDRESS                                     PORTS   AGE
longhorn-ingress   *       45.79.165.114,66.228.45.37,97.107.142.125   80      2m7s

$ curl -v http://97.107.142.125/
*   Trying 97.107.142.125...
* TCP_NODELAY set
* Connected to 97.107.142.125 (97.107.142.125) port 80 (#0)
> GET / HTTP/1.1
> Host: 97.107.142.125
> User-Agent: curl/7.64.1
> Accept: */*
>
< HTTP/1.1 401 Unauthorized
< Server: openresty/1.15.8.1
< Date: Fri, 29 May 2020 11:47:33 GMT
< Content-Type: text/html
< Content-Length: 185
< Connection: keep-alive
< WWW-Authenticate: Basic realm="Authentication Required"
<
<html>
<head><title>401 Authorization Required</title></head>
<body>
<center><h1>401 Authorization Required</h1></center>
<hr><center>openresty/1.15.8.1</center>
</body>
</html>
* Connection #0 to host 97.107.142.125 left intact
* Closing connection 0

$ curl -v http://97.107.142.125/ -u foo:bar
*   Trying 97.107.142.125...
* TCP_NODELAY set
* Connected to 97.107.142.125 (97.107.142.125) port 80 (#0)
* Server auth using Basic with user 'foo'
> GET / HTTP/1.1
> Host: 97.107.142.125
> Authorization: Basic Zm9vOmJhcg==
> User-Agent: curl/7.64.1
> Accept: */*
>
< HTTP/1.1 200 OK
< Date: Fri, 29 May 2020 11:51:27 GMT
< Content-Type: text/html
< Content-Length: 1118
< Last-Modified: Thu, 28 May 2020 00:39:41 GMT
< ETag: "5ecf084d-3fd"
< Cache-Control: max-age=0
<
<!DOCTYPE html>
<html lang="en">
......
```

## Additional Steps for AWS EKS Kubernetes Clusters

You will need to create an ELB (Elastic Load Balancer) to expose the nginx Ingress controller to the Internet. Additional costs may apply.

1. Create pre-requisite resources according to the [nginx ingress controller documentation.](https://kubernetes.github.io/ingress-nginx/deploy/#prerequisite-generic-deployment-command)

2. Create an ELB by following [these steps.](https://kubernetes.github.io/ingress-nginx/deploy/#aws)

## References
https://kubernetes.github.io/ingress-nginx/


---

## Article: deploy/uninstall/_index.md

---
title: Uninstall Longhorn
weight: 6
---

In this section, you'll learn how to uninstall Longhorn.


- [Prerequisite](#prerequisite)
- [Uninstalling Longhorn from the Rancher UI](#uninstalling-longhorn-from-the-rancher-ui)
- [Uninstalling Longhorn using Helm](#uninstalling-longhorn-using-helm)
- [Uninstalling Longhorn using Helm Controller](#uninstalling-longhorn-using-helm-controller)
- [Uninstalling Longhorn using Fleet](#uninstalling-longhorn-using-fleet)
- [Uninstalling Longhorn using Flux](#uninstalling-longhorn-using-flux)
- [Uninstalling Longhorn using Argo CD](#uninstalling-longhorn-using-argo-cd)
- [Uninstalling Longhorn using kubectl](#uninstalling-longhorn-using-kubectl)
- [Troubleshooting](#troubleshooting)
  - [Uninstalling using Rancher UI or Helm failed, I am not sure why](#uninstalling-using-rancher-ui-or-helm-failed-i-am-not-sure-why)
  - [I deleted the Longhorn App from Rancher UI instead of following the uninstallation procedure](#i-deleted-the-longhorn-app-from-rancher-ui-instead-of-following-the-uninstallation-procedure)
  - [Problems with CRDs](#problems-with-crds)

### Prerequisite
To prevent Longhorn from being accidentally uninstalled (which leads to data lost),
we introduce a new setting, [deleting-confirmation-flag](../../references/settings/#deleting-confirmation-flag).
If this flag is **false**, the Longhorn uninstallation job will fail.
Set this flag to **true** to allow Longhorn uninstallation.
You can set this flag using setting page in Longhorn UI or `kubectl -n longhorn-system patch -p '{"value": "true"}' --type=merge lhs deleting-confirmation-flag`


To prevent damage to the Kubernetes cluster, we recommend deleting all Kubernetes workloads using Longhorn volumes (PersistentVolume, PersistentVolumeClaim, StorageClass, Deployment, StatefulSet, DaemonSet, etc).

### Uninstalling Longhorn from the Rancher UI

From Rancher UI, navigate to `Catalog Apps` tab and delete Longhorn app.

### Uninstalling Longhorn using Helm

Run this command:

```
helm uninstall longhorn -n longhorn-system
```

### Uninstalling Longhorn using Helm Controller

Run this command:

```
kubectl delete helmchart <HelmChart name> -n <HelmChart namespace>
```

### Uninstalling Longhorn Using Fleet

Run the following command:

```
kubectl delete GitRepo longhorn -n fleet-local
```

After the `longhorn-uninstall` job is completed, run the following command:

```
kubectl delete -f https://raw.githubusercontent.com/longhorn/longhorn/v{{< current-version >}}/deploy/longhorn.yaml
```

### Uninstalling Longhorn Using Flux

Run the following command:

```
flux delete helmrelease longhorn-release -n longhorn-system
```

### Uninstalling Longhorn Using Argo CD

Argo CD currently does not support the PreDelete resource hook. Instead of running `argocd app delete longhorn` directly, you must [uninstall Longhorn using kubectl](#uninstalling-longhorn-using-kubectl) to prevent dangling resources from remaining in the `longhorn-system` namespace.

### Uninstalling Longhorn using kubectl

1. Create the uninstallation job to clean up CRDs from the system and wait for success:

    ```
    kubectl create -f https://raw.githubusercontent.com/longhorn/longhorn/v{{< current-version >}}/uninstall/uninstall.yaml
    kubectl get job/longhorn-uninstall -n longhorn-system -w
    ```

    Example output:
    ```
    $ kubectl create -f https://raw.githubusercontent.com/longhorn/longhorn/v{{< current-version >}}/uninstall/uninstall.yaml
    serviceaccount/longhorn-uninstall-service-account created
    clusterrole.rbac.authorization.k8s.io/longhorn-uninstall-role created
    clusterrolebinding.rbac.authorization.k8s.io/longhorn-uninstall-bind created
    job.batch/longhorn-uninstall created

    $ kubectl get job/longhorn-uninstall -n longhorn-system -w
    NAME                 COMPLETIONS   DURATION   AGE
    longhorn-uninstall   0/1           3s         3s
    longhorn-uninstall   1/1           20s        20s
    ```

2. Remove remaining components:
    ```
    kubectl delete -f https://raw.githubusercontent.com/longhorn/longhorn/v{{< current-version >}}/deploy/longhorn.yaml
    kubectl delete -f https://raw.githubusercontent.com/longhorn/longhorn/v{{< current-version >}}/uninstall/uninstall.yaml
    ```

> **Tip:** If you try `kubectl delete -f https://raw.githubusercontent.com/longhorn/longhorn/v{{< current-version >}}/deploy/longhorn.yaml` first and get stuck there,
pressing `Ctrl C` then running `kubectl create -f https://raw.githubusercontent.com/longhorn/longhorn/v{{< current-version >}}/uninstall/uninstall.yaml` can also help you remove Longhorn. Finally, don't forget to cleanup remaining components.




### Troubleshooting
#### Uninstalling using Rancher UI or Helm failed, I am not sure why
You might want to check the logs of the `longhorn-uninstall-xxx` pod inside `longhorn-system` namespace to see why it failed.
One reason can be that [deleting-confirmation-flag](../../references/settings/#deleting-confirmation-flag) is `false`.
You can set it to `true` by using setting page in Longhorn UI or `kubectl -n longhorn-system patch -p '{"value": "true"}' --type=merge lhs deleting-confirmation-flag`
then retry the Helm/Rancher uninstallation.

If the uninstallation was an accident (you don't actually want to uninstall Longhorn),
you can cancel the uninstallation as the following.
1. If you use Rancher UI to deploy Longhorn
   1. Open a kubectl shell on Rancher UI
   1. Find the latest revision of Longhorn release
      ```shell
      > helm list -n longhorn-system -a
      NAME            NAMESPACE       REVISION        UPDATED                                 STATUS          CHART                                   APP VERSION
      longhorn        longhorn-system 2               2022-10-14 01:22:36.929130451 +0000 UTC uninstalling    longhorn-100.2.3+up1.3.2-rc1            v1.3.2-rc1
      longhorn-crd    longhorn-system 3               2022-10-13 22:19:05.976625081 +0000 UTC deployed        longhorn-crd-100.2.3+up1.3.2-rc1        v1.3.2-rc1
      ```
   1. Rollback to the latest revision
      ```shell
      > helm rollback longhorn 2 -n longhorn-system
      checking 22 resources for changes
      ...
      Rollback was a success! Happy Helming!
      ```
1. If you use Helm deploy Longhorn
   1. Open a kubectl terminal
   1. Find the latest revision of Longhorn release
      ```shell
      ➜ helm list --namespace longhorn-system -a
      NAME            NAMESPACE       REVISION        UPDATED                                 STATUS          CHART                   APP VERSION
      longhorn        longhorn-system 1               2022-10-14 13:45:25.341292504 -0700 PDT uninstalling    longhorn-1.4.0-dev      v1.4.0-dev
      ```
   1. Rollback to the latest revision
      ```shell
      ➜  helm rollback longhorn 1 -n longhorn-system
      Rollback was a success! Happy Helming!
      ```


#### I deleted the Longhorn App from Rancher UI instead of following the uninstallation procedure

Redeploy the (same version) Longhorn App. Follow the uninstallation procedure above.

#### Problems with CRDs

If your CRD instances or the CRDs themselves can't be deleted for whatever reason, run the commands below to clean up. Caution: this will wipe all Longhorn state!

```shell
# Delete CRD finalizers, instances and definitions
for crd in $(kubectl get crd -o jsonpath={.items[*].metadata.name} | tr ' ' '\n' | grep longhorn.io); do
  kubectl -n ${NAMESPACE} get $crd -o yaml | sed "s/\- longhorn.io//g" | kubectl apply -f -
  kubectl -n ${NAMESPACE} delete $crd --all
  kubectl delete crd/$crd
done
```

If you encounter the following error, it is possible that an incomplete uninstallation removed the Longhorn validation or modification webhook services, but left the same services registered.  

`for: "STDIN": error when patching "STDIN": Internal error occurred: failed calling webhook "validator.longhorn.io": failed to call webhook: Post "https://longhorn-admission-webhook.longhorn-system.svc:9502/v1/webhook/validation?timeout=10s": service "longhorn-admission-webhook" not found`

You can run the following commands to check the status of the webhook services.

```shell
$ kubectl get ValidatingWebhookConfiguration -A
NAME                               WEBHOOKS   AGE
longhorn-webhook-validator         1          46d
rancher.cattle.io                  7          133d
rke2-ingress-nginx-admission       1          133d
rke2-snapshot-validation-webhook   1          133d

$ kubectl get MutatingWebhookConfiguration -A
NAME                       WEBHOOKS   AGE
longhorn-webhook-mutator   1          46d
rancher.cattle.io          4          133d
```

If either or both are still registered, you can delete the configuration to remove the services from the patch operation call path.

```shell
$ kubectl delete ValidatingWebhookConfiguration longhorn-webhook-validator
validatingwebhookconfiguration.admissionregistration.k8s.io "longhorn-webhook-validator" deleted

$ kubectl delete MutatingWebhookConfiguration longhorn-webhook-mutator
mutatingwebhookconfiguration.admissionregistration.k8s.io "longhorn-webhook-mutator" deleted
```

The script should run successfully after the configuration is deleted.

```shell
Warning: Detected changes to resource pvc-279e8c3e-bfb0-4233-8899-77b5b178c08c which is currently being deleted.
volumeattachment.longhorn.io/pvc-279e8c3e-bfb0-4233-8899-77b5b178c08c configured
No resources found
customresourcedefinition.apiextensions.k8s.io "volumeattachments.longhorn.io" deleted
```

---
Please see [link](https://github.com/longhorn/longhorn) for more information.


---

## Article: deploy/upgrade/_index.md

---
title: Upgrade
weight: 3
---

Here we cover how to upgrade to the latest Longhorn from all previous releases.

# Deprecation & Incompatibility

There are no deprecated or incompatible changes introduced in v{{< current-version >}}.

# Upgrade Path Enforcement and Downgrade Prevention

Starting with v1.5.0, Longhorn only supports upgrades from one minor version to the next. For example, upgrading from 1.5.x to 1.6.x is supported, but skipping versions (e.g., from 1.4.x to 1.6.x) is not. If you attempt to upgrade from an unsupported version or skip a minor version, the operation will fail automatically. However, you can revert to the previously installed version without service interruption or downtime.

Moreover, Longhorn does not support downgrades to earlier versions. This restriction helps prevent unexpected system behavior and issues associated with function incompatibility, deprecation, or removal.

> **Warning**:
> - Once you successfully upgrade to v{{< current-version >}}, you will not be allowed to revert to the previously installed version.
> - The Downgrade Prevention feature was introduced in v1.5.0 so Longhorn is unable to prevent downgrade attempts in older versions.
However, downgrading is completely unsupported and is therefore not recommended.

The following table outlines the supported upgrade paths.

  |  Current version |  Target version |  Supported | Example |
  |    :-:      |    :-:      |   :-:  |    :-:    |
  |  x.y.*      |  x.(y+1).*  |   ✓    |  v1.4.2  to  v1.5.1  |
  |  x.y.*      |  x.y.(*+n)  |   ✓    |  v1.5.0  to  v1.5.1  |
  |  x.y[^lastMinorVersion].*      |  (x+1).y.*  |   ✓    |  v1.30.0 to  v2.0.0  |
  |  x.(y-1).*  |  x.(y+1).*  |   X    |  v1.3.3  to  v1.5.1  |
  |  x.(y-2).*  |  x.(y+1).*  |   X    |  v1.2.6  to  v1.5.1  |

[^lastMinorVersion]: Longhorn only allows upgrades from any patch version of the last minor release before the new major version. For example, if v1.3.0 is the last minor version before v2.0, you can upgrade from any patch version of v1.3.0 to any patch version of v2.0.

## Manual Checks Before Upgrade
Automated checks are only performed on some upgrade paths, and the pre-upgrade checker may not cover some scenarios.  Manual checks, performed using either kubectl or the UI, are recommended for these schenarios.  You can take mitigating actions or defer the upgrade until issues are addressed.
- Ensure that all V2 Data Engine volumes are detached and the replicas are stopped.  The V2 Data Engine currently does not support live upgrades.
- Avoid upgrading when volumes are in the "Faulted" status.  If all the replicas are deemed unusable, they may be deleted and data may be permanently lost (if no usable backups exist).
- Avoid upgrading if a failed BackingImage exists.  For more information, see [Backing Image](../../advanced-resources/backing-image/backing-image).
- It is recommended to create a [Longhorn system backup](../../advanced-resources/system-backup-restore/backup-longhorn-system) before performing the upgrade. This ensures that all critical resources, such as volumes and backing images, are backed up and can be restored in case any issues arise.

# Upgrading Longhorn

There are normally two steps in the upgrade process: first upgrade Longhorn manager to the latest version, then manually upgrade the Longhorn engine to the latest version using the latest Longhorn manager.

## 1. Upgrade Longhorn manager

- To upgrade from v1.9.x, see [this section.](./longhorn-manager)

## 2. Manually Upgrade Longhorn Engine

After Longhorn Manager is upgraded, Longhorn Engine also needs to be upgraded [using the Longhorn UI.](./upgrade-engine)

## 3. Automatically Upgrade Longhorn Engine

Since Longhorn v1.1.1, we provide an option to help you [automatically upgrade engines](./auto-upgrade-engine)

## 4. Automatically Migrate Recurring Jobs

With the introduction of the new label-driven `Recurring Job` feature, Longhorn has removed the `RecurringJobs` field in the Volume Spec and planned to deprecate `RecurringJobs` in the StorageClass.

During the upgrade, Longhorn will automatically:
- Create new recurring job CRs from the `recurringJobs` field in Volume Spec and convert them to the volume labels.
- Create new recurring job CRs from the `recurringJobs` in the StorageClass and convert them to the new `recurringJobSelector` parameter.

Visit [Recurring Snapshots and Backups](../../snapshots-and-backups/scheduling-backups-and-snapshots) for more information about the new `Recurring Job` feature.

# Extended Reading

Visit [Some old instance manager pods are still running after upgrade](https://longhorn.io/kb/troubleshooting-some-old-instance-manager-pods-are-still-running-after-upgrade) for more information about the cleanup strategy of instance manager pods during upgrade.

# Need Help?

If you have any issues, please report it at
https://github.com/longhorn/longhorn/issues and include your backup yaml files
as well as manager logs.


---

## Article: deploy/upgrade/auto-upgrade-engine.md

---
title: Automatically Upgrading Longhorn Engine
weight: 3
---

Since Longhorn v1.1.1, we provide an option to help you automatically upgrade Longhorn volumes to the new default engine version after upgrading Longhorn manager.
This feature reduces the amount of manual work you have to do when upgrading Longhorn.
There are a few concepts related to this feature as listed below:

#### 1. Concurrent Automatic Engine Upgrade Per Node Limit Setting

This is a setting that controls how Longhorn automatically upgrades volumes' engines to the new default engine image after upgrading Longhorn manager.
The value of this setting specifies the maximum number of engines per node that are allowed to upgrade to the default engine image at the same time.
If the value is 0, Longhorn will not automatically upgrade volumes' engines to the default version.
The bigger this value is, the faster the engine upgrade process finishes.

However, giving a bigger value for this setting will consume more CPU and memory of the node during the engine upgrade process.
We recommend setting the value to 3 to leave some room for error but don't overwhelm the system with too many failed upgrades.

#### 2. The behavior of Longhorn with different volume conditions.
In the following cases, assume that the `concurrent automatic engine upgrade per node limit` setting is bigger than 0.

1. Attached Volumes

   If the volume is in attached state and healthy, Longhorn will automatically do a live upgrade for the volume's engine to the new default engine image.

1. Detached Volumes

   Longhorn automatically does an offline upgrade for detached volume.

1. Disaster Recovery Volumes

   Longhorn doesn't automatically upgrade [disaster recovery volumes](../../../snapshots-and-backups/setup-disaster-recovery-volumes/) to the new default engine image because it would trigger a full restoration for the disaster recovery volumes.
The full restoration might affect the performance of other running Longhorn volumes in the system.
So, Longhorn leaves it to you to decide when it is the good time to manually upgrade the engine for disaster recovery volumes (e.g., when the system is idle or during the maintenance time).

   However, when you activate the disaster recovery volume, it will be activated and then detached.
At this time, Longhorn will automatically do offline upgrade for the volume similar to the detached volume case.

#### 3. What Happened If The Upgrade Fails?
If a volume failed to upgrade its engine, the engine image in volume's spec will remain to be different than the engine image in the volume's status.
Longhorn will continuously retry to upgrade until it succeeds.

If there are too many volumes that fail to upgrade per node (i.e., more than the `concurrent automatic engine upgrade per node limit` setting),
Longhorn will stop upgrading volume on that node.


---

## Article: deploy/upgrade/longhorn-manager.md

---
title: Upgrading Longhorn Manager
weight: 1
---

### Upgrading from v1.9.x

We only support upgrading to v{{< current-version >}} from v1.9.x. For other versions, please upgrade to v1.9.x first.

Engine live upgrade is supported from v1.9.x to v{{< current-version >}}.

For airgap upgrades when Longhorn is installed as a Rancher app, you will need to modify the image names and remove the registry URL part.

For example, the image `registry.example.com/longhorn/longhorn-manager:v{{< current-version >}}` is changed to `longhorn/longhorn-manager:v{{< current-version >}}` in Longhorn images section. For more information, see the air gap installation steps [here.](../../install/airgap/#using-a-rancher-app)

#### Preparing for the Upgrade

If Longhorn was installed using a Helm Chart, or if it was installed as Rancher catalog app, check to make sure the parameters in the default StorageClass weren't changed. Changing the default StorageClass's parameter might result in a chart upgrade failure. if you want to reconfigure the parameters in the StorageClass, you can copy the default StorageClass's configuration to create another StorageClass.

    The current default StorageClass has the following parameters:

        parameters:
          numberOfReplicas: <user specified replica count, 3 by default>
          staleReplicaTimeout: "30"
          fromBackup: ""
          baseImage: ""

#### Upgrade

> **Prerequisite:** Always back up volumes before upgrading. If anything goes wrong, you can restore the volume using the backup.

#### Upgrade as a Rancher Catalog App

To upgrade the Longhorn App, make sure which Rancher UI the existing Longhorn App was installed with. There are two Rancher UIs, one is the Cluster Manager (old UI), and the other one is the Cluster Explorer (new UI). The Longhorn App in different UIs considered as two different applications by Rancher. They cannot upgrade to each other. If you installed Longhorn in the Cluster Manager, you need to use the Cluster Manager to upgrade Longhorn to a newer version, and vice versa for the Cluster Explorer.

> Note: Because the Cluster Manager (old UI) is being deprecated, we provided the instruction to migrate the existing Longhorn installation to the Longhorn chart in the Cluster Explorer (new UI) [here](https://longhorn.io/kb/how-to-migrate-longhorn-chart-installed-in-old-rancher-ui-to-the-chart-in-new-rancher-ui/)

Different Rancher UIs screenshots.
- The Cluster Manager (old UI)
{{< figure src="/img/screenshots/install/cluster-manager.png" >}}
- The Cluster Explorer (new UI)
{{< figure src="/img/screenshots/install/cluster-explorer.png" >}}

On Kubernetes clusters managed by Rancher 2.1 or newer, the steps to upgrade the catalog app `longhorn-system` are the similar to the installation steps.

#### Upgrade with Kubectl

To upgrade with kubectl, run this command:

```
kubectl apply -f https://raw.githubusercontent.com/longhorn/longhorn/v{{< current-version >}}/deploy/longhorn.yaml
```

#### Upgrade with Helm

To upgrade with Helm, run this command:

```
helm upgrade longhorn longhorn/longhorn --namespace longhorn-system --version {{< current-version >}}
```

#### Upgrade with Helm Controller

Update the value of `spec.version` in the `HelmChart` YAML file:

```yaml
spec:
  version: v{{< current-version >}} # Replace with the Longhorn version you'd like to upgrade to
  chart: longhorn
  repo: https://charts.longhorn.io
  failurePolicy: abort
```

Alternatively, if using the `spec.chartContent` key, create a patch file with 
```yaml
spec:
  chartContent: <base64 content> # tar cz of longhorn charts directory for release | base64 -w 0
```
and then apply it with 
```
kubectl  patch helmchart longhorn -n <namespace> --type merge --patch-file <name of patch file>
```

> **IMPORTANT!** In both cases, ensure that `spec.failurePolicy` is set to "abort".  The only other value is the default: "reinstall", which performs an uninstall of Longhorn if the pre-upgrade check or the upgrade fails.  With "abort", it retries periodically, giving the user a chance to fix the problem.

#### Upgrade with Fleet

Update the value of `helm.version` in the `fleet` YAML file of your GitOps repository.

```yaml
helm:
  repo: https://charts.longhorn.io
  chart: longhorn
  version: v{{< current-version >}} # Replace with the Longhorn version you'd like to upgrade to
  releaseName: longhorn
```

#### Upgrade with Flux

Update the value of `spec.chart.spec.version` in the `HelmRelease` YAML file of your GitOps repository.

```yaml
spec:
  chart:
    spec:
      chart: longhorn
      reconcileStrategy: ChartVersion
      sourceRef:
        kind: HelmRepository
        name: longhorn
      version: v{{< current-version >}} # Replace with the Longhorn version you'd like to upgrade to
```

#### Upgrade with Argo CD

Update the value of `targetRevision` in the `Application` YAML file of your GitOps repository.

```yaml
spec:
  project: default
  sources:
    - chart: longhorn
      repoURL: https://charts.longhorn.io
      targetRevision: v{{< current-version >}} # Replace with the Longhorn version you'd like to upgrade to
```

Then wait for all the pods to become running and Longhorn UI working. e.g.:

```
$ kubectl -n longhorn-system get pod
NAME                                                  READY   STATUS    RESTARTS      AGE
engine-image-ei-4dbdb778-nw88l                        1/1     Running   0             4m29s
longhorn-ui-b7c844b49-jn5g6                           1/1     Running   0             75s
longhorn-manager-z2p8h                                1/1     Running   0             71s
instance-manager-b34d5db1fe1e2d52bcfb308be3166cfc     1/1     Running   0             65s
longhorn-driver-deployer-6bd59c9f76-jp6pg             1/1     Running   0             75s
engine-image-ei-df38d2e5-zccq5                        1/1     Running   0             65s
csi-snapshotter-588457fcdf-h2lgc                      1/1     Running   0             30s
csi-resizer-6d8cf5f99f-8v4sp                          1/1     Running   1 (30s ago)   37s
csi-snapshotter-588457fcdf-6pgf4                      1/1     Running   0             30s
csi-provisioner-869bdc4b79-7ddwd                      1/1     Running   1 (30s ago)   44s
csi-snapshotter-588457fcdf-p4kkn                      1/1     Running   0             30s
csi-attacher-7bf4b7f996-mfbdn                         1/1     Running   1 (30s ago)   50s
csi-provisioner-869bdc4b79-4dc7n                      1/1     Running   1 (30s ago)   43s
csi-resizer-6d8cf5f99f-vnspd                          1/1     Running   1 (30s ago)   37s
csi-attacher-7bf4b7f996-hrs7w                         1/1     Running   1 (30s ago)   50s
csi-attacher-7bf4b7f996-rt2s9                         1/1     Running   1 (30s ago)   50s
csi-resizer-6d8cf5f99f-7vv89                          1/1     Running   1 (30s ago)   37s
csi-provisioner-869bdc4b79-sn6zr                      1/1     Running   1 (30s ago)   43s
longhorn-csi-plugin-b2zzj                             2/2     Running   0             24s
```

Next, [upgrade Longhorn engine.](../upgrade-engine)

### Upgrading from Unsupported Versions

We only support upgrading to v{{< current-version >}} from v1.9.x. For other versions, please upgrade to v1.9.x first.

If you attempt to upgrade from an unsupported version, the upgrade will fail. When encountering an upgrade failure, please consider the following scenarios to recover the state based on different upgrade methods.

#### Upgrade with Kubectl

When you upgrade with kubectl by running this command:

```shell
kubectl apply -f https://raw.githubusercontent.com/longhorn/longhorn/v{{< current-version >}}/deploy/longhorn.yaml
```

Longhorn will block the upgrade process and provide the failure reason in the logs of the `longhorn-manager` pod.
During the upgrade failure, the user's Longhorn system should remain intact without any impacts except `longhorn-manager` daemon set.

To recover, you need to apply the manifest of the previously installed version using the following command:

```shell
kubectl apply -f https://raw.githubusercontent.com/longhorn/longhorn/[previous installed version]/deploy/longhorn.yaml
```

Besides, users might need to delete new components introduced by the new version manually.

#### Upgrade with Helm or Rancher App Marketplace

To prevent any impact caused by failed upgrades from unsupported versions, Longhorn will automatically initiate a new job (`pre-upgrade`) to verify if the upgrade path is supported before upgrading when upgrading through `Helm` or `Rancher App Marketplace`.

The `pre-upgrade` job will block the upgrade process and provide the failure reason in the logs of the pod.  It will also be recorded in an event, for instance:

```
2m33s               Normal    Created                 Pod/longhorn-pre-upgrade-v5tqq    Created container longhorn-pre-upgrade
2m33s               Warning   FailedUpgradePreCheck   /longhorn-pre-upgrade             failed to upgrade since upgrading from v1.6.2 to v1.8.0 for minor version is not supported
```

During the upgrade failure, the user's Longhorn system should remain intact without any impacts.

To recover, you need to run the below commands to rollback to the previously installed revision:

```shell
# get previous installed Longhorn REVISION
helm history longhorn
helm rollback longhorn [REVISION]

# or
helm upgrade longhorn longhorn/longhorn --namespace longhorn-system --version [previous installed version]
```

To recover, you need to upgrade to the previously installed revision at `Rancher App Marketplace` again.

### TroubleShooting
1. Error: `"longhorn" is invalid: provisioner: Forbidden: updates to provisioner are forbidden.`
- This means there are some modifications applied to the default storageClass and you need to clean up the old one before upgrade.

- To clean up the deprecated StorageClass, run this command:
    ```
    kubectl delete -f https://raw.githubusercontent.com/longhorn/longhorn/v{{< current-version >}}/examples/storageclass.yaml
    ```


---

## Article: deploy/upgrade/upgrade-engine.md

---
title: Manually Upgrading Longhorn Engine
weight: 2
---

In this section, you'll learn how to manually upgrade the Longhorn Engine from the Longhorn UI.

### Prerequisites

Always make backups before upgrading the Longhorn engine images.

Upgrade the Longhorn manager before upgrading the Longhorn engine.

### Offline Upgrade

Follow these steps if the live upgrade is not available, or if the volume is stuck in degraded state:

1. Follow [the detach procedure for relevant workloads](../../../nodes-and-volumes/volumes/detaching-volumes).
2. Select all the volumes using batch selection. Click the batch operation button **Upgrade Engine**, and choose the engine image available in the list. It's the default engine shipped with the manager for this release.
3. Resume all workloads. Any volume not part of a Kubernetes workload must be attached from the Longhorn UI.

### Live upgrade

Live upgrade is supported for upgrading from v1.9.x to v{{< current-version >}}.

The `iSCSI` frontend does not support live upgrades.

Live upgrade should only be done with healthy volumes.

1. Select the volume you want to upgrade.
2. Click `Upgrade Engine` in the drop down.
3. Select the engine image you want to upgrade to.
    1. Normally it's the only engine image in the list, since the UI exclude the current image from the list.
4. Click OK.

During the live upgrade, the user will see double number of the replicas temporarily. After upgrade complete, the user should see the same number of the replicas as before, and the `Engine Image` field of the volume should be updated.

Notice after the live upgrade, Rancher or Kubernetes would still show the old version of image for the engine, and new version for the replicas. It's expected. The upgrade is success if you see the new version of image listed as the volume image in the Volume Detail page.

### Clean up the old image

After you've done upgrade for all the images, select `Settings/Engine Image` from Longhorn UI. Now you should able to remove the non-default image.


---

## Article: deploy/install/_index.md

---
title: Quick Installation
description: Install Longhorn on Kubernetes
weight: 1
---

- [Installation Requirements](#installation-requirements)
  - [OS/Distro Specific Configuration](#osdistro-specific-configuration)
  - [Checking the Kubernetes Version](#checking-the-kubernetes-version)
  - [Pod Security Policy](#pod-security-policy)
  - [Notes on Mount Propagation](#notes-on-mount-propagation)
  - [Root and Privileged Permission](#root-and-privileged-permission)
  - [Installing open-iscsi](#installing-open-iscsi)
  - [Installing NFSv4 client](#installing-nfsv4-client)
  - [Installing Cryptsetup and LUKS](#installing-cryptsetup-and-luks)
  - [Installing Device Mapper Userspace Tool](#installing-device-mapper-userspace-tool)
- [Longhorn Command Line Tool](#longhorn-command-line-tool)
  - [Checking Prerequisites Using Longhorn Command Line Tool](#checking-prerequisites-using-longhorn-command-line-tool)
  - [Installing Prerequisites Using Longhorn Command Line Tool](#installing-prerequisites-using-longhorn-command-line-tool)

> **Note**: This quick installation guide uses some configurations which are not for production usage.
> Please see [Best Practices](../../best-practices/) for how to configure Longhorn for production usage.

Longhorn can be installed on a Kubernetes cluster in several ways:

- [Rancher Catalog App](./install-with-rancher)
- [kubectl](./install-with-kubectl/)
- [Helm](./install-with-helm/)
- [Helm Controller](./install-with-helm-controller/)
- [Fleet](./install-with-fleet/)
- [Flux](./install-with-flux/)
- [ArgoCD](./install-with-argocd/)

To install Longhorn in an air gapped environment, refer to [Air Gap Installation](../install/airgap).

For information on customizing Longhorn's default settings, refer to [Customizing Default Settings](../../advanced-resources/deploy/customizing-default-settings).

For information on deploying Longhorn on specific nodes and rejecting general workloads for those nodes, refer to the section on [Taints and Tolerations](../../advanced-resources/deploy/taint-toleration).

## Installation Requirements

Each node in the Kubernetes cluster where Longhorn is installed must fulfill the following requirements:

-  A container runtime compatible with Kubernetes (Docker v1.13+, containerd v1.3.7+, etc.)
-  Kubernetes >= v1.25
-  `open-iscsi` is installed, and the `iscsid` daemon is running on all the nodes. This is necessary, since Longhorn relies on `iscsiadm` on the host to provide persistent volumes to Kubernetes. For help installing `open-iscsi`, refer to [Installing open-iscsi](#installing-open-iscsi).
-  RWX support requires that each node has a NFSv4 client installed.
    - For installing a NFSv4 client, refer to [Installing NFSv4 Client](#installing-nfsv4-client).
- The host filesystem supports the `file extents` feature to store the data. Currently we support:
    - ext4
    - XFS
- `bash`, `curl`, `findmnt`, `grep`, `awk`, `blkid`, `lsblk` must be installed.
- [Mount propagation](https://kubernetes-csi.github.io/docs/deploying.html#enabling-mount-propagation) must be enabled.

The Longhorn workloads must be able to run as root in order for Longhorn to be deployed and operated properly.

[Longhorn Command Line Tool](../../advanced-resources/longhornctl/) can be used to check the Longhorn environment for potential issues.

For the minimum recommended hardware, refer to the [best practices guide](../../best-practices/#minimum-recommended-hardware).

### OS/Distro Specific Configuration

You must perform additional setups before using Longhorn with certain operating systems and distributions.

- Google Kubernetes Engine (GKE): See [Longhorn CSI on GKE](../../advanced-resources/os-distro-specific/csi-on-gke).
- K3s clusters: See [Longhorn CSI on K3s](../../advanced-resources/os-distro-specific/csi-on-k3s).
- RKE clusters with CoreOS: See [Longhorn CSI on RKE and CoreOS](../../advanced-resources/os-distro-specific/csi-on-rke-and-coreos).
- OCP/OKD clusters: See [OKD Support](../../advanced-resources/os-distro-specific/okd-support).
- Talos Linux clusters: See [Talos Linux Support](../../advanced-resources/os-distro-specific/talos-linux-support).
- Container-Optimized OS: See [Container-Optimized OS Support](../../advanced-resources/os-distro-specific/container-optimized-os-support).

### Checking the Kubernetes Version

Use the following command to check your Kubernetes server version

```shell
kubectl version
```

Result:

```shell
Client Version: version.Info{Major:"1", Minor:"26", GitVersion:"v1.26.10", GitCommit:"b8609d4dd75c5d6fba4a5eaa63a5507cb39a6e99", GitTreeState:"clean", BuildDate:"2023-10-18T11:44:31Z", GoVersion:"go1.20.10", Compiler:"gc", Platform:"linux/amd64"}
Server Version: version.Info{Major:"1", Minor:"26", GitVersion:"v1.26.10+k3s2", GitCommit:"cb5cb5557f34e240e38c68a8c4ca2506c68b1d86", GitTreeState:"clean", BuildDate:"2023-11-08T03:21:46Z", GoVersion:"go1.20.10", Compiler:"gc", Platform:"linux/amd64"}
```

The `Server Version` should be >= v1.25.

### Pod Security Policy

Starting with v1.0.2, Longhorn is shipped with a default Pod Security Policy that will give Longhorn the necessary privileges to be able to run properly.

No special configuration is needed for Longhorn to work properly on clusters with Pod Security Policy enabled.

### Notes on Mount Propagation

If your Kubernetes cluster was provisioned by Rancher v2.0.7+ or later, the MountPropagation feature is enabled by default.

If MountPropagation is disabled, Base Image feature will be disabled.

### Root and Privileged Permission

Longhorn components require root access with privileged permissions to achieve volume operations and management, because Longhorn relies on system resources on the host across different namespaces, for example, Longhorn uses `nsenter` to understand block devices' usage or encrypt/decrypt volumes on the host.

Below are the directories Longhorn components requiring access with root and privileged permissions :

- Longhorn Manager
  - /boot (read only): Get required modules' information from /boot/config-$(uname -r) on the host.
  - /dev: Block devices created by Longhorn are under the `/dev` path.
  - /proc (read only): Find the recognized host process like container runtime, then use `nsenter` to access the mounts on the host to understand disks usage.
  - /etc (read only): Read the necessary system configuration to get node status updated, for example, `nfsmount.conf`.
  - /var/lib/longhorn: The default path for storing volume data on a host.
- Longhorn Engine Image
  - /var/lib/longhorn/engine-binaries: The default path for storing the Longhorn engine binaries.
- Longhorn Instance Manager
  - /: Access any data path on this node and access Longhorn engine binaries.
  - /dev: Block devices created by Longhorn are under the `/dev` path.
  - /proc: Find the recognized host process like container runtime, then use `nsenter` to manage iSCSI targets and initiators, also some file system
- Longhorn Share Manager
  - /dev: Block devices created by Longhorn are under the `/dev` path.
  - /lib/modules: Kernel modules required by `cryptsetup` for volume encryption.
  - /proc: Find the recognized host process like container runtime, then use `nsenter` for volume encryption.
  - /sys: Support volume encryption by `cryptsetup`.
- Longhorn CSI Plugin
  - /: For host checks via the NFS customer mounter (deprecated). Note that, this will be removed in the future release.
  - /dev: Block devices created by Longhorn are under the `/dev` path.
  - /lib/modules: Kernel modules required by Longhorn CSI plugin.
  - /sys: Support volume encryption by `cryptsetup`.
  - /var/lib/kubelet/plugins/kubernetes.io/csi: The path where the Longhorn CSI plugin creates the staging path (via `NodeStageVolume`) of a block device. The staging path will be bind-mounted to the target path `/var/lib/kubelet/pods` (via `NodePublishVolume`) for support single volume could be mounted to multiple Pods.
  - /var/lib/kubelet/plugins_registry: The path where the node-driver-registrar registers the CSI plugin with kubelet.
  - /var/lib/kubelet/plugins/driver.longhorn.io: The path where the socket for the communication between kubelet Longhorn CSI driver.
  - /var/lib/kubelet/pods: The path where the Longhorn CSI driver mounts volume from the target path (via `NodePublishVolume`).
- Longhorn CSI Attacher/Provisioner/Resizer/Snapshotter
  - /var/lib/kubelet/plugins/driver.longhorn.io: The path where the socket for the communication between kubelet Longhorn CSI driver.
- Longhorn Backing Image Manager
  - /var/lib/longhorn: The default path for storing data on the host.
- Longhorn Backing Image Data Source
  - /var/lib/longhorn: The default path for storing data on the host.
- Longhorn System Restore Rollout
  - /var/lib/longhorn/engine-binaries: The default path for storing the Longhorn engine binaries.

### Installing open-iscsi

The command used to install `open-iscsi` differs depending on the Linux distribution.

For GKE, we recommend using Ubuntu as the guest OS image since it contains`open-iscsi` already.

You may need to edit the cluster security group to allow SSH access.

- SUSE and openSUSE: Run the following command:
  ```
  zypper install open-iscsi
  systemctl enable iscsid
  systemctl start iscsid
  ```

- Debian and Ubuntu: Run the following command:
  ```
  apt-get install open-iscsi
  ```

- RHEL, CentOS, and EKS *(EKS Kubernetes Worker AMI with AmazonLinux2 image)*: Run the following commands:
  ```
  yum --setopt=tsflags=noscripts install iscsi-initiator-utils
  echo "InitiatorName=$(/sbin/iscsi-iname)" > /etc/iscsi/initiatorname.iscsi
  systemctl enable iscsid
  systemctl start iscsid
  ```

- Talos Linux: See [Talos Linux Support](../../advanced-resources/os-distro-specific/talos-linux-support).

- Container-Optimized OS: See [Container-Optimized OS Support](../../advanced-resources/os-distro-specific/container-optimized-os-support)

Please ensure iscsi_tcp module has been loaded before iscsid service starts. Generally, it should be automatically loaded along with the package installation.

```
modprobe iscsi_tcp
```

> **Important**: On SUSE and openSUSE, the `iscsi_tcp` module is included only in the `kernel-default` package. If the `kernel-default-base` package is installed on your system, you must replace it with `kernel-default`.

We also provide an `iscsi` installer to make it easier for users to install `open-iscsi` automatically. You can use the [Longhorn CLI](../../advanced-resources/longhornctl/) to install the prerequisites.

You can use the `longhornctl check preflight` command. This command verifies your Kubernetes cluster environment to ensure it meets Longhorn's requirements. It performs a series of checks that can help identify potential issues that may prevent Longhorn from functioning correctly.

Example:
```sh
> longhornctl check preflight
INFO[2025-08-22T12:58:40+08:00] Initializing preflight checker               
INFO[2025-08-22T12:58:40+08:00] Cleaning up preflight checker                
INFO[2025-08-22T12:58:42+08:00] Running preflight checker                    
INFO[2025-08-22T12:58:49+08:00] Retrieved preflight checker result:
ip-10-0-1-132:
  info:
  - Service iscsid is running
  - NFS4 is supported
  - Package nfs-client is installed
  - Package open-iscsi is installed
  - Package cryptsetup is installed
  - Package device-mapper is installed
  - Module dm_crypt is loaded
  warn:
  - Kube DNS "coredns" is set with fewer than 2 replicas; consider increasing replica count for high availability
INFO[2025-08-22T12:58:49+08:00] Cleaning up preflight checker                
INFO[2025-08-22T12:58:50+08:00] Completed preflight checker
```

In rare cases, it may be required to modify the installed SELinux policy to get Longhorn working. If you are running
an up-to-date version of a Fedora downstream distribution (e.g. Fedora, RHEL, Rocky, CentOS, etc.) and plan to leave
SELinux enabled, see [the KB](../../../../kb/troubleshooting-volume-attachment-fails-due-to-selinux-denials) for details.

### Installing NFSv4 client

In Longhorn system, backup feature requires NFSv4, v4.1 or v4.2, and ReadWriteMany (RWX) volume feature requires NFSv4.1. Before installing NFSv4 client userspace daemon and utilities, make sure the client kernel support is enabled on each Longhorn node.

- Check if `NFSv4` support is enabled in the kernel:
  ```
  cat /boot/config-`uname -r`| grep CONFIG_NFS_V4
  ```

- Check if `NFSv4.1` support is enabled in the kernel:
  ```
  cat /boot/config-`uname -r`| grep CONFIG_NFS_V4_1
  ```

- Check if `NFSv4.2` support is enabled in the kernel:
  ```
  cat /boot/config-`uname -r`| grep CONFIG_NFS_V4_2
  ```

The command used to install a NFSv4 client differs depending on the Linux distribution.

- For Debian and Ubuntu, use this command:
  ```
  apt-get install nfs-common
  ```

- For RHEL, CentOS, and EKS with `EKS Kubernetes Worker AMI with AmazonLinux2 image`, use this command:
  ```
  yum install nfs-utils
  ```

- For SUSE/OpenSUSE you can install a NFSv4 client via:
  ```
  zypper install nfs-client
  ```

- For Talos Linux, [the NFS client is part of the `kubelet` image maintained by the Talos team](https://www.talos.dev/v1.6/kubernetes-guides/configuration/storage/#nfs).

- For Container-Optimized OS, [the NFS is supported with the node image](https://cloud.google.com/kubernetes-engine/docs/concepts/node-images#storage_driver_support).

We also provide an `nfs` installer to make it easier for users to install `nfs-client` automatically. You can use this [Longhorn CLI](../../advanced-resources/longhornctl/) to install the prerequisite.

You can use the `longhornctl check preflight` command. This command verifies your Kubernetes cluster environment to ensure it meets Longhorn's requirements. It performs a series of checks that can help identify potential issues that may prevent Longhorn from functioning correctly.

Example:
```sh
> longhornctl check preflight
INFO[2025-08-22T12:58:40+08:00] Initializing preflight checker               
INFO[2025-08-22T12:58:40+08:00] Cleaning up preflight checker                
INFO[2025-08-22T12:58:42+08:00] Running preflight checker                    
INFO[2025-08-22T12:58:49+08:00] Retrieved preflight checker result:
ip-10-0-1-132:
  info:
  - Service iscsid is running
  - NFS4 is supported
  - Package nfs-client is installed
  - Package open-iscsi is installed
  - Package cryptsetup is installed
  - Package device-mapper is installed
  - Module dm_crypt is loaded
  warn:
  - Kube DNS "coredns" is set with fewer than 2 replicas; consider increasing replica count for high availability
INFO[2025-08-22T12:58:49+08:00] Cleaning up preflight checker                
INFO[2025-08-22T12:58:50+08:00] Completed preflight checker
```

> **Notice:**  
> These steps only verify that the kernel supports NFSv4, v4.1, or v4.2.  
> To verify the NFS version in use, run `mount | grep nfs` or `nfsstat -m` to confirm the mounted version. Using the correct NFS version is required for backup and RWX volume features in Longhorn.

### Installing Cryptsetup and LUKS

[Cryptsetup](https://gitlab.com/cryptsetup/cryptsetup) is an open-source utility used to conveniently set up `dm-crypt` based device-mapper targets and Longhorn uses [LUKS2](https://gitlab.com/cryptsetup/cryptsetup#luks-design) (Linux Unified Key Setup) format that is the standard for Linux disk encryption to support volume encryption.

The command used to install the cryptsetup tool differs depending on the Linux distribution.

- For Debian and Ubuntu, use this command:

  ```shell
  apt-get install cryptsetup
  ```

- For RHEL, CentOS, Rocky Linux and EKS with `EKS Kubernetes Worker AMI with AmazonLinux2 image`, use this command:

  ```shell
  yum install cryptsetup
  ```

- For SUSE/OpenSUSE, use this command:

  ```shell
  zypper install cryptsetup
  ```

### Installing Device Mapper Userspace Tool

The device mapper is a framework provided by the Linux kernel for mapping physical block devices onto higher-level virtual block devices. It forms the foundation of the `dm-crypt` disk encryption and provides the linear dm device on the top of v2 volume. The device mapper is typically included by default in many Linux distributions. Some lightweight or highly customized distributions or a minimal installation of a distribution might exclude it to save space or reduce complexity

The command used to install the device mapper differs depending on the Linux distribution.

- For Debian and Ubuntu, use this command:

  ```shell
  apt-get install dmsetup
  ```

- For RHEL, CentOS, Rocky Linux and EKS with `EKS Kubernetes Worker AMI with AmazonLinux2 image`, use this command:

  ```shell
  yum install device-mapper
  ```

- For SUSE/OpenSUSE, use this command:

  ```shell
  zypper install device-mapper
  ```

## Longhorn Command Line Tool

### Checking Prerequisites Using Longhorn Command Line Tool

The `longhornctl` tool is a CLI for Longhorn operations. For more information, see [Command Line Tool (longhornctl)](../../advanced-resources/longhornctl/).

To check the prerequisites and configurations, download the tool and run the `check` sub-command:

```shell
# For AMD64 platform
curl -sSfL -o longhornctl https://github.com/longhorn/cli/releases/download/v{{< current-version >}}/longhornctl-linux-amd64
# For ARM platform
curl -sSfL -o longhornctl https://github.com/longhorn/cli/releases/download/v{{< current-version >}}/longhornctl-linux-arm64

chmod +x longhornctl
./longhornctl check preflight
```

Example of result:

```shell
INFO[2024-01-01T00:00:01Z] Initializing preflight checker
INFO[2024-01-01T00:00:01Z] Cleaning up preflight checker
INFO[2024-01-01T00:00:01Z] Running preflight checker
INFO[2024-01-01T00:00:02Z] Retrieved preflight checker result:
worker1:
  info:
  - Service iscsid is running
  - NFS4 is supported
  - Package nfs-common is installed
  - Package open-iscsi is installed
  warn:
  - multipathd.service is running. Please refer to https://longhorn.io/kb/troubleshooting-volume-with-multipath/ for more information.
worker2:
  info:
  - Service iscsid is running
  - NFS4 is supported
  - Package nfs-common is not installed
  - Package open-iscsi is installed
```

Use the `install` sub-command to install and set up the preflight dependencies before installing Longhorn.

```shell
./longhornctl install preflight
```

> **Note**:
> Some immutable Linux distributions, such as SUSE Linux Enterprise Micro (SLE Micro), require you to reboot worker nodes after running the `install` sub-command. After the reboot, you must run the `install` sub-command again to complete the operation.
>
> The documentation of the Linux distribution you are using should outline such requirements. For example, the [SLE Micro documentation](https://documentation.suse.com/sle-micro/6.0/html/Micro-transactional-updates/index.html#reference-transactional-update-usage) explains how all changes made by the `transactional-update` command become active only after the node is rebooted.

### Installing Prerequisites Using Longhorn Command Line Tool

```shell
longhornctl --kube-config ~/.kube/config --image longhornio/longhorn-cli:v{{< current-version >}} install preflight
```

Example of result:

```shell
INFO[2025-03-11T08:17:57+08:00] Initializing preflight installer
INFO[2025-03-11T08:17:57+08:00] Cleaning up preflight installer
INFO[2025-03-11T08:17:57+08:00] Running preflight installer
INFO[2025-03-11T08:17:57+08:00] Installing dependencies with package manager
INFO[2025-03-11T08:18:28+08:00] Installed dependencies with package manager
INFO[2025-03-11T08:18:28+08:00] Cleaning up preflight installer
INFO[2025-03-11T08:18:28+08:00] Completed preflight installer. Use 'longhornctl check preflight' to check the result.
```

---

## Article: deploy/install/airgap.md

---
title: Air Gap Installation
weight: 100
---

Longhorn can be installed in an air gapped environment by using a manifest file, a Helm chart, or the Rancher UI.

- [Requirements](#requirements)
- [Using a Manifest File](#using-a-manifest-file)
- [Using a Helm chart](#using-a-helm-chart)
- [Using a Rancher app](#using-a-rancher-app)
- [Troubleshooting](#troubleshooting)

## Requirements
  - Deploy Longhorn Components images to your own registry.
  - Deploy Kubernetes CSI driver components images to your own registry.

#### Note:
  - A full list of all needed images is in [longhorn-images.txt](https://raw.githubusercontent.com/longhorn/longhorn/v{{< current-version >}}/deploy/longhorn-images.txt). First, download the images list by running:
    ```shell
    wget https://raw.githubusercontent.com/longhorn/longhorn/v{{< current-version >}}/deploy/longhorn-images.txt
    ```
  - We provide a script, [save-images.sh](https://raw.githubusercontent.com/longhorn/longhorn/v{{< current-version >}}/scripts/save-images.sh), to quickly pull the above `longhorn-images.txt` list. If you specify a `tar.gz` file name for flag `--images`, the script will save all images to the provided filename. In the example below, the script pulls and saves Longhorn images to the file `longhorn-images.tar.gz`. You then can copy the file to your air-gap environment. On the other hand, if you don't specify the file name, the script just pulls the list of images to your computer.
    ```shell
    wget https://raw.githubusercontent.com/longhorn/longhorn/v{{< current-version >}}/scripts/save-images.sh
    chmod +x save-images.sh
    ./save-images.sh --image-list longhorn-images.txt --images longhorn-images.tar.gz
    ```
  - We provide another script, [load-images.sh](https://raw.githubusercontent.com/longhorn/longhorn/v{{< current-version >}}/scripts/load-images.sh), to push Longhorn images to your private registry. If you specify a `tar.gz` file name for flag `--images`, the script loads images from the `tar` file and pushes them. Otherwise, it will find images in your local Docker and push them. In the example below, the script loads images from the file `longhorn-images.tar.gz` and pushes them to `<YOUR-PRIVATE-REGISTRY>`
    ```shell
    wget https://raw.githubusercontent.com/longhorn/longhorn/v{{< current-version >}}/scripts/load-images.sh
    chmod +x load-images.sh
    ./load-images.sh --image-list longhorn-images.txt --images longhorn-images.tar.gz --registry <YOUR-PRIVATE-REGISTRY>
    ```
  - For more options with using the scripts, see flag `--help`:
    ```shell
    ./save-images.sh --help
    ./load-images.sh --help
    ```

## Using a Manifest File

1. Get Longhorn Deployment manifest file

    `wget https://raw.githubusercontent.com/longhorn/longhorn/v{{< current-version >}}/deploy/longhorn.yaml`

2. Create Longhorn namespace

    `kubectl create namespace longhorn-system`


3. If private registry require authentication, Create `docker-registry` secret in `longhorn-system` namespace:

    `kubectl -n longhorn-system create secret docker-registry <SECRET_NAME> --docker-server=<REGISTRY_URL> --docker-username=<REGISTRY_USER> --docker-password=<REGISTRY_PASSWORD>`

    * Add your secret name  `SECRET_NAME` to `imagePullSecrets.name` in the following resources
      * `longhorn-driver-deployer` Deployment
      * `longhorn-manager` DaemonSet
      * `longhorn-ui` Deployment

      Example:
      ```yaml
      apiVersion: apps/v1
      kind: Deployment
      metadata:
        labels:
          app: longhorn-ui
        name: longhorn-ui
        namespace: longhorn-system
      spec:
        replicas: 1
        selector:
          matchLabels:
            app: longhorn-ui
        template:
          metadata:
            labels:
              app: longhorn-ui
          spec:
            containers:
            - name: longhorn-ui
              image: longhornio/longhorn-ui:v0.8.0
              ports:
              - containerPort: 8000
              env:
                - name: LONGHORN_MANAGER_IP
                  value: "http://longhorn-backend:9500"
            imagePullSecrets:
            - name: <SECRET_NAME>                          ## Add SECRET_NAME here
            serviceAccountName: longhorn-service-account
      ```

4. Apply the following modifications to the manifest file

    * Modify Kubernetes CSI driver components environment variables in `longhorn-driver-deployer` Deployment point to your private registry images
      * CSI_ATTACHER_IMAGE
      * CSI_PROVISIONER_IMAGE
      * CSI_NODE_DRIVER_REGISTRAR_IMAGE
      * CSI_RESIZER_IMAGE
      * CSI_SNAPSHOTTER_IMAGE

      ```yaml
      - name: CSI_ATTACHER_IMAGE
        value: <REGISTRY_URL>/csi-attacher:<CSI_ATTACHER_IMAGE_TAG>
      - name: CSI_PROVISIONER_IMAGE
        value: <REGISTRY_URL>/csi-provisioner:<CSI_PROVISIONER_IMAGE_TAG>
      - name: CSI_NODE_DRIVER_REGISTRAR_IMAGE
        value: <REGISTRY_URL>/csi-node-driver-registrar:<CSI_NODE_DRIVER_REGISTRAR_IMAGE_TAG>
      - name: CSI_RESIZER_IMAGE
        value: <REGISTRY_URL>/csi-resizer:<CSI_RESIZER_IMAGE_TAG>
      - name: CSI_SNAPSHOTTER_IMAGE
        value: <REGISTRY_URL>/csi-snapshotter:<CSI_SNAPSHOTTER_IMAGE_TAG>
      ```

    * Modify Longhorn images to point to your private registry images
      * longhornio/longhorn-manager

        `image: <REGISTRY_URL>/longhorn-manager:<LONGHORN_MANAGER_IMAGE_TAG>`

      * longhornio/longhorn-engine

        `image: <REGISTRY_URL>/longhorn-engine:<LONGHORN_ENGINE_IMAGE_TAG>`

      * longhornio/longhorn-instance-manager

        `image: <REGISTRY_URL>/longhorn-instance-manager:<LONGHORN_INSTANCE_MANAGER_IMAGE_TAG>`

      * longhornio/longhorn-share-manager

        `image: <REGISTRY_URL>/longhorn-share-manager:<LONGHORN_SHARE_MANAGER_IMAGE_TAG>`

      * longhornio/longhorn-ui

        `image: <REGISTRY_URL>/longhorn-ui:<LONGHORN_UI_IMAGE_TAG>`

      Example:
      ```yaml
      apiVersion: apps/v1
      kind: Deployment
      metadata:
        labels:
          app: longhorn-ui
        name: longhorn-ui
        namespace: longhorn-system
      spec:
        replicas: 1
        selector:
          matchLabels:
            app: longhorn-ui
        template:
          metadata:
            labels:
              app: longhorn-ui
          spec:
            containers:
            - name: longhorn-ui
              image: <REGISTRY_URL>/longhorn-ui:<LONGHORN_UI_IMAGE_TAG>   ## Add image name and tag here
              ports:
              - containerPort: 8000
              env:
                - name: LONGHORN_MANAGER_IP
                  value: "http://longhorn-backend:9500"
            imagePullSecrets:
            - name: <SECRET_NAME>
            serviceAccountName: longhorn-service-account
      ```

5. Deploy Longhorn using modified manifest file
   `kubectl apply -f longhorn.yaml`

## Using a Helm Chart

In v{{< current-version >}}, Longhorn automatically adds <REGISTRY_URL> prefix to images. You simply need to set the registryUrl parameters to pull images from your private registry.

> **Note:** Once you set registryUrl to your private registry, Longhorn tries to pull images from the registry exclusively. Make sure all Longhorn components' images are in the registry otherwise Longhorn will fail to pull images.

### Use default image name

If you keep the images' names as recommended [here](./#recommendation), you only need to do the following steps:

1. Clone the Longhorn repo:

    `git clone https://github.com/longhorn/longhorn.git`

2. In `chart/values.yaml`

    * Specify `Private registry URL`. If the registry requires authentication, specify `Private registry user`, `Private registry password`, and `Private registry secret`.
    Longhorn will automatically generate a secret with the those information and use it to pull images from your private registry.

      ```yaml
      privateRegistry:
        # -- Setting that allows you to create a private registry secret.
        createSecret: true
        # -- URL of a private registry. When unspecified, Longhorn uses the default system registry.
        registryUrl: <REGISTRY_URL>
        # -- User account used for authenticating with a private registry.
        registryUser: <REGISTRY_USER>
        # -- Password for authenticating with a private registry.
        registryPasswd: <REGISTRY_PASSWORD>
        # -- Kubernetes secret that allows you to pull images from a private registry. This setting applies only when creation of private registry secrets is enabled. You must include the private registry name in the secret name.
        registrySecret: <REGISTRY_SECRET_NAME>
      ```

### Use custom image name

If you want to use custom images' names, you can use the following steps:

1. Clone longhorn repo

    `git clone https://github.com/longhorn/longhorn.git`

2. In `chart/values.yaml`

    > **Note:** Do not include the private registry prefix, it will be added automatically. e.g: if your image is `example.com/username/longhorn-manager`, use `username/longhorn-manager` in the following charts.

    - Specify Longhorn images and tag:

        ```yaml
        image:
          longhorn:
            engine:
              repository: longhornio/longhorn-engine
              tag: <LONGHORN_ENGINE_IMAGE_TAG>
            manager:
              repository: longhornio/longhorn-manager
              tag: <LONGHORN_MANAGER_IMAGE_TAG>
            ui:
              repository: longhornio/longhorn-ui
              tag: <LONGHORN_UI_IMAGE_TAG>
            instanceManager:
              repository: longhornio/longhorn-instance-manager
              tag: <LONGHORN_INSTANCE_MANAGER_IMAGE_TAG>
            shareManager:
              repository: longhornio/longhorn-share-manager
              tag: <LONGHORN_SHARE_MANAGER_IMAGE_TAG>
        ```

    - Specify CSI Driver components images and tag:

        ```yaml
          csi:
            attacher:
              repository: longhornio/csi-attacher
              tag: <CSI_ATTACHER_IMAGE_TAG>
            provisioner:
              repository: longhornio/csi-provisioner
              tag: <CSI_PROVISIONER_IMAGE_TAG>
            nodeDriverRegistrar:
              repository: longhornio/csi-node-driver-registrar
              tag: <CSI_NODE_DRIVER_REGISTRAR_IMAGE_TAG>
            resizer:
              repository: longhornio/csi-resizer
              tag: <CSI_RESIZER_IMAGE_TAG>
            snapshotter:
              repository: longhornio/csi-snapshotter
              tag: <CSI_SNAPSHOTTER_IMAGE_TAG>
        ```

    - Specify `Private registry URL`. If the registry requires authentication, specify `Private registry user`, `Private registry password`, and `Private registry secret`.
    Longhorn will automatically generate a secret with the those information and use it to pull images from your private registry.

      ```yaml
      privateRegistry:
        # -- Setting that allows you to create a private registry secret.
        createSecret: true
        # -- URL of a private registry. When unspecified, Longhorn uses the default system registry.
        registryUrl: <REGISTRY_URL>
        # -- User account used for authenticating with a private registry.
        registryUser: <REGISTRY_USER>
        # -- Password for authenticating with a private registry.
        registryPasswd: <REGISTRY_PASSWORD>
        # -- Kubernetes secret that allows you to pull images from a private registry. This setting applies only when creation of private registry secrets is enabled. You must include the private registry name in the secret name.
        registrySecret: <REGISTRY_SECRET_NAME>
      ```

3. Install Longhorn

  ```shell
  helm install longhorn ./chart --namespace longhorn-system --create-namespace
  ```

# Using a Rancher App

### Use default image name

If you keep the images' names as recommended [here](./#recommendation), you only need to do the following steps:

- In the `Private Registry Settings` section specify:
   - Private registry URL
   - Private registry user
   - Private registry password
   - Private registry secret name

  Longhorn will automatically generate a secret with the those information and use it to pull images from your private registry.

  ![images](/img/screenshots/airgap-deploy/app-default-images.png)

### Use custom image name

- If you want to use custom images' names, you can set `Use Default Images` to `False` and specify images' names.

  > **Note:** Do not include the private registry prefix, it will be added automatically. e.g: if your image is `example.com/username/longhorn-manager`, use `username/longhorn-manager` in the following charts.

  ![images](/img/screenshots/airgap-deploy/app-custom-images.png)

- Specify `Private registry URL`. If the registry requires authentication, specify `Private registry user`, `Private registry password`, and `Private registry secret name`.
  Longhorn will automatically generate a secret with the those information and use it to pull images from your private registry.

  ![images](/img/screenshots/airgap-deploy/app-custom-images-reg.png)

## Troubleshooting

#### For Helm/Rancher installation, if user forgot to submit a secret to authenticate to private registry, `longhorn-manager DaemonSet` will fail to create.


1. Create the Kubernetes secret

    `kubectl -n longhorn-system create secret docker-registry <SECRET_NAME> --docker-server=<REGISTRY_URL> --docker-username=<REGISTRY_USER> --docker-password=<REGISTRY_PASSWORD>`


2. Create `registry-secret` setting object manually.

    ```yaml
    apiVersion: longhorn.io/v1beta2
    kind: Setting
    metadata:
      name: registry-secret
      namespace: longhorn-system
    value: <SECRET_NAME>
    ```

    `kubectl apply -f registry-secret.yml`


3. Delete Longhorn and re-install it again.

    * **Helm2**

      `helm uninstall ./chart --name longhorn --namespace longhorn-system`

      `helm install ./chart --name longhorn --namespace longhorn-system`

    * **Helm3**

      `helm uninstall longhorn ./chart --namespace longhorn-system`

      `helm install longhorn ./chart --namespace longhorn-system`

## Recommendation:
It's highly recommended not to manipulate image tags, especially instance manager image tags such as v1_20200301, because we intentionally use the date to avoid associating it with a Longhorn version.

The images of Longhorn's components are hosted in Dockerhub under the `longhornio` account. For example, `longhornio/longhorn-manager:v{{< current-version >}}`. It's recommended to keep the account name, `longhornio`, the same when you push the images to your private registry. This helps avoid unnecessary configuration issues.


---

## Article: deploy/install/install-with-argocd.md

---
title: Install with ArgoCD
weight: 13
---

## Prerequisites
- Your workstation: Install the [Argo CD CLI](https://argo-cd.readthedocs.io/en/stable/cli_installation/).
- Kubernetes cluster:
  - Ensure that each node fulfills the [installation requirements](../#installation-requirements).
  - Install [Argo CD](https://argo-cd.readthedocs.io/en/stable/).

    ```bash
    kubectl create namespace argocd
    kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/core-install.yaml
    ```
    Allow some time for the deployment of Argo CD components in the `argocd` namespace.

> Use [Longhorn Command Line Tool](../../../advanced-resources/longhornctl/) to check the Longhorn environment for potential issues.

## Installing Longhorn

1. Log in to Argo CD.

    ```bash
    argocd login --core
    ```

1. Set the current namespace to `argocd`.

    ```bash
    kubectl config set-context --current --namespace=argocd
    ```

1. Create the Longhorn Application custom resource.

    ```bash
    cat > longhorn-application.yaml <<EOF
    apiVersion: argoproj.io/v1alpha1
    kind: Application
    metadata:
      name: longhorn
      namespace: argocd
    spec:
      syncPolicy:
        syncOptions:
          - CreateNamespace=true
      project: default
      sources:
        - chart: longhorn
          repoURL: https://charts.longhorn.io/
          targetRevision: v{{< current-version >}} # Replace with the Longhorn version you'd like to install or upgrade to
          helm:
            values: |
              preUpgradeChecker:
                jobEnabled: false
      destination:
        server: https://kubernetes.default.svc
        namespace: longhorn-system
    EOF
    kubectl apply -f longhorn-application.yaml
    ```

1. Deploy Longhorn with the configured settings.

    ```bash
    argocd app sync longhorn
    ```

1. Verify that Longhorn was installed successfully.

    ```bash
    kubectl -n longhorn-system get pod
    ```

    Example of a successful Longhorn installation:

    ```bash
    NAME                                                READY   STATUS    RESTARTS   AGE
    longhorn-ui-b7c844b49-w25g5                         1/1     Running   0          2m41s
    longhorn-manager-pzgsp                              1/1     Running   0          2m41s
    longhorn-driver-deployer-6bd59c9f76-lqczw           1/1     Running   0          2m41s
    longhorn-csi-plugin-mbwqz                           2/2     Running   0          100s
    csi-snapshotter-588457fcdf-22bqp                    1/1     Running   0          100s
    csi-snapshotter-588457fcdf-2wd6g                    1/1     Running   0          100s
    csi-provisioner-869bdc4b79-mzrwf                    1/1     Running   0          101s
    csi-provisioner-869bdc4b79-klgfm                    1/1     Running   0          101s
    csi-resizer-6d8cf5f99f-fd2ck                        1/1     Running   0          101s
    csi-provisioner-869bdc4b79-j46rx                    1/1     Running   0          101s
    csi-snapshotter-588457fcdf-bvjdt                    1/1     Running   0          100s
    csi-resizer-6d8cf5f99f-68cw7                        1/1     Running   0          101s
    csi-attacher-7bf4b7f996-df8v6                       1/1     Running   0          101s
    csi-attacher-7bf4b7f996-g9cwc                       1/1     Running   0          101s
    csi-attacher-7bf4b7f996-8l9sw                       1/1     Running   0          101s
    csi-resizer-6d8cf5f99f-smdjw                        1/1     Running   0          101s
    instance-manager-b34d5db1fe1e2d52bcfb308be3166cfc   1/1     Running   0          114s
    engine-image-ei-df38d2e5-cv6nc                      1/1     Running   0          114s
    ```

1. [Create an NGINX Ingress controller with basic authentication](../../accessing-the-ui/longhorn-ingress) to access the Longhorn UI. Authentication to the Longhorn UI is not enabled by default.

1. [Access the Longhorn UI](../../accessing-the-ui).


---

## Article: deploy/install/install-with-fleet.md

---
title: Install with Fleet
weight: 11
---

## Prerequisites
- Your workstation: Install [Helm](https://helm.sh/docs/) v3.0 or later.
- Kubernetes cluster:
  - Ensure that each node fulfills the [installation requirements](../#installation-requirements).
  - Install [Fleet](https://fleet.rancher.io/) using Helm.

    ```bash
    helm repo add fleet https://rancher.github.io/fleet-helm-charts/
    helm -n cattle-fleet-system install --create-namespace --wait fleet-crd fleet/fleet-crd
    helm -n cattle-fleet-system install --create-namespace --wait fleet fleet/fleet
    ```
    Allow some time for the deployment of Fleet components in the `cattle-fleet-system` namespace.

> Use [Longhorn Command Line Tool](../../../advanced-resources/longhornctl/) to check the Longhorn environment for potential issues.

## Installing Longhorn

1. In your GitOps repository, create a [fleet.yaml](https://fleet.rancher.io/ref-fleet-yaml) file that includes the following:

    - Parameter for installing Longhorn in the `longhorn-system` namespace

    ```yaml
    defaultNamespace: longhorn-system
    ```

    - Parameters for [ignoring modified CRDs](https://fleet.rancher.io/bundle-diffs)

    ```yaml
    diff:
      comparePatches:
      - apiVersion: apiextensions.k8s.io/v1
        kind: CustomResourceDefinition
        name: engineimages.longhorn.io
        operations:
        - {"op": "replace", "path": "/status"}
      - apiVersion: apiextensions.k8s.io/v1
        kind: CustomResourceDefinition
        name: nodes.longhorn.io
        operations:
        - {"op": "replace", "path": "/status"}
      - apiVersion: apiextensions.k8s.io/v1
        kind: CustomResourceDefinition
        name: volumes.longhorn.io
        operations:
        - {"op": "replace", "path": "/status"}
      - apiVersion: apiextensions.k8s.io/v1
        kind: CustomResourceDefinition
        name: engines.longhorn.io
        operations:
        - {"op": "replace", "path": "/status"}
      - apiVersion: apiextensions.k8s.io/v1
        kind: CustomResourceDefinition
        name: instancemanagers.longhorn.io
        operations:
        - {"op": "replace", "path": "/status"}
      - apiVersion: apiextensions.k8s.io/v1
        kind: CustomResourceDefinition
        name: replicas.longhorn.io
        operations:
        - {"op": "replace", "path": "/status"}
      - apiVersion: apiextensions.k8s.io/v1
        kind: CustomResourceDefinition
        name: settings.longhorn.io
        operations:
        - {"op": "replace", "path": "/status"}
    ```

    - Parameters for specifying the version of the Longhorn Helm chart to be installed

    ```yaml
    helm:
      repo: https://charts.longhorn.io
      chart: longhorn
      version: v{{< current-version >}} # Replace with the Longhorn version you'd like to install or upgrade to
      releaseName: longhorn
    ```

    Example of a complete `fleet.yaml` file:

    ```yaml
    defaultNamespace: longhorn-system
    helm:
      repo: https://charts.longhorn.io
      chart: longhorn
      version: v{{< current-version >}}
      releaseName: longhorn
    diff:
      comparePatches:
      - apiVersion: apiextensions.k8s.io/v1
        kind: CustomResourceDefinition
        name: engineimages.longhorn.io
        operations:
        - {"op": "replace", "path": "/status"}
      - apiVersion: apiextensions.k8s.io/v1
        kind: CustomResourceDefinition
        name: nodes.longhorn.io
        operations:
        - {"op": "replace", "path": "/status"}
      - apiVersion: apiextensions.k8s.io/v1
        kind: CustomResourceDefinition
        name: volumes.longhorn.io
        operations:
        - {"op": "replace", "path": "/status"}
    ```

1. Create a GitRepo custom resource (CR) that points to your GitOps repository.

    ```bash
    cat > longhorn-gitrepo.yaml << "EOF"
    apiVersion: fleet.cattle.io/v1alpha1
    kind: GitRepo
    metadata:
      name: longhorn
      namespace: fleet-local
    spec:
      repo: https://github.com/your-username/your-gitops-repo.git
      revision: main
      paths:
      - .
    EOF
    ```

1. Apply the GitRepo CR.

    ```bash
    kubectl apply -f longhorn-gitrepo.yaml
    ```

1. Verify that the GitRepo CR was created and synced successfully.

    ```bash
    kubectl -n fleet-local get gitrepo -w
    ```

1. Verify that Longhorn was installed successfully.

    ```bash
    kubectl -n longhorn-system get pod
    ```

    Example of a successful Longhorn installation:

    ```bash
    NAME                                                READY   STATUS    RESTARTS   AGE
    longhorn-ui-b7c844b49-w25g5                         1/1     Running   0          2m41s
    longhorn-manager-pzgsp                              1/1     Running   0          2m41s
    longhorn-driver-deployer-6bd59c9f76-lqczw           1/1     Running   0          2m41s
    longhorn-csi-plugin-mbwqz                           2/2     Running   0          100s
    csi-snapshotter-588457fcdf-22bqp                    1/1     Running   0          100s
    csi-snapshotter-588457fcdf-2wd6g                    1/1     Running   0          100s
    csi-provisioner-869bdc4b79-mzrwf                    1/1     Running   0          101s
    csi-provisioner-869bdc4b79-klgfm                    1/1     Running   0          101s
    csi-resizer-6d8cf5f99f-fd2ck                        1/1     Running   0          101s
    csi-provisioner-869bdc4b79-j46rx                    1/1     Running   0          101s
    csi-snapshotter-588457fcdf-bvjdt                    1/1     Running   0          100s
    csi-resizer-6d8cf5f99f-68cw7                        1/1     Running   0          101s
    csi-attacher-7bf4b7f996-df8v6                       1/1     Running   0          101s
    csi-attacher-7bf4b7f996-g9cwc                       1/1     Running   0          101s
    csi-attacher-7bf4b7f996-8l9sw                       1/1     Running   0          101s
    csi-resizer-6d8cf5f99f-smdjw                        1/1     Running   0          101s
    instance-manager-b34d5db1fe1e2d52bcfb308be3166cfc   1/1     Running   0          114s
    engine-image-ei-df38d2e5-cv6nc                      1/1     Running   0          114s
    ```

1. [Create an NGINX Ingress controller with basic authentication](../../accessing-the-ui/longhorn-ingress) to access the Longhorn UI. Authentication to the Longhorn UI is not enabled by default.

1. [Access the Longhorn UI](../../accessing-the-ui).


---

## Article: deploy/install/install-with-flux.md

---
title: Install with Flux
weight: 12
---

## Prerequisites
- Your workstation: Install [Helm](https://helm.sh/docs/) v3.0 or later.
- Kubernetes cluster:
  - Ensure that each node fulfills the [installation requirements](../#installation-requirements).
  - [Install the Flux CLI and controllers](https://fluxcd.io/flux/installation/#install-the-flux-cli).
  - [Bootstrap Flux with GitHub](https://fluxcd.io/flux/installation/bootstrap/github/) using the Flux CLI.
    Run the following commands to export your GitHub personal access token (PAT) as an environment variable, deploy the Flux controllers on your cluster, and configure the controllers to sync the cluster state from the specified GitHub repository.

    ```bash
    export GITHUB_TOKEN=<gh-token>
    flux bootstrap github \
      --token-auth \
      --owner=<github_username> \
      --repository=<github_repo_name> \
      --branch=<branch_name> \
      --path=<folder_path_within_repo> \
      --personal
    ```

> Use [Longhorn Command Line Tool](../../../advanced-resources/longhornctl/) to check the Longhorn environment for potential issues.

## Installing Longhorn

1. Create a HelmRepository custom resource (CR) that points to the Longhorn Helm chart URL.

    ```bash
    kubectl create ns longhorn-system
    flux create source helm longhorn-repo \
      --url=https://charts.longhorn.io \
      --namespace=longhorn-system \
      --export > helmrepo.yaml
    kubectl apply -f helmrepo.yaml
    ```

1. Create a HelmRelease CR that references the HelmRepository and specifies the version of the Longhorn Helm chart to be installed.

    ```bash
    flux create helmrelease longhorn-release \
      --chart=longhorn \
      --source=HelmRepository/longhorn-repo \
      --chart-version=v{{< current-version >}} \
      --namespace=longhorn-system \
      --export > helmrelease.yaml
    kubectl apply -f helmrelease.yaml
    ```

1. Verify that the HelmRelease CR was created and synced successfully.

    ```bash
    flux get helmrelease longhorn-release -n longhorn-system
    ```

1. Verify that Longhorn was installed successfully.

    ```bash
    kubectl -n longhorn-system get pod
    ```

    Example of a successful Longhorn installation:

    ```bash
    NAME                                                READY   STATUS    RESTARTS   AGE
    longhorn-ui-b7c844b49-w25g5                         1/1     Running   0          2m41s
    longhorn-manager-pzgsp                              1/1     Running   0          2m41s
    longhorn-driver-deployer-6bd59c9f76-lqczw           1/1     Running   0          2m41s
    longhorn-csi-plugin-mbwqz                           2/2     Running   0          100s
    csi-snapshotter-588457fcdf-22bqp                    1/1     Running   0          100s
    csi-snapshotter-588457fcdf-2wd6g                    1/1     Running   0          100s
    csi-provisioner-869bdc4b79-mzrwf                    1/1     Running   0          101s
    csi-provisioner-869bdc4b79-klgfm                    1/1     Running   0          101s
    csi-resizer-6d8cf5f99f-fd2ck                        1/1     Running   0          101s
    csi-provisioner-869bdc4b79-j46rx                    1/1     Running   0          101s
    csi-snapshotter-588457fcdf-bvjdt                    1/1     Running   0          100s
    csi-resizer-6d8cf5f99f-68cw7                        1/1     Running   0          101s
    csi-attacher-7bf4b7f996-df8v6                       1/1     Running   0          101s
    csi-attacher-7bf4b7f996-g9cwc                       1/1     Running   0          101s
    csi-attacher-7bf4b7f996-8l9sw                       1/1     Running   0          101s
    csi-resizer-6d8cf5f99f-smdjw                        1/1     Running   0          101s
    instance-manager-b34d5db1fe1e2d52bcfb308be3166cfc   1/1     Running   0          114s
    engine-image-ei-df38d2e5-cv6nc                      1/1     Running   0          114s
    ```

1. [Create an NGINX Ingress controller with basic authentication](../../accessing-the-ui/longhorn-ingress) to access the Longhorn UI. Authentication to the Longhorn UI is not enabled by default.

1. [Access the Longhorn UI](../../accessing-the-ui).

## Continuous Operations via GitOps

You can commit and push exported manifests to your GitOps repository.

    ```bash
    git add helmrepo.yaml helmrelease.yaml
    git commit -m "Add HelmRepository and HelmRelease for Longhorn installation"
    git push origin <branch_name>
    ```

Afterwards, you can modify the HelmRelease and HelmRepository CRs by editing the YAML manifests in your GitOps repository. Flux automatically detects and applies the changes without requiring direct access to your Kubernetes cluster.


---

## Article: deploy/install/install-with-helm-controller.md

---
title: Install with Helm Controller
weight: 10
---

In this section, you will learn how to install Longhorn with the HelmChart controller built into RKE2 and K3s.

### Prerequisites

- Kubernetes cluster: Ensure that each node fulfills the [installation requirements](../#installation-requirements).  Cluster should be running RKE2 or K3s.

> [Longhorn Command Line Tool](../../../advanced-resources/longhornctl/) can be used to check the Longhorn environment for potential issues.

### Installing Longhorn


> **Note**:
> * The initial settings for Longhorn can be [customized using Helm options or by editing the deployment configuration file.](../../../advanced-resources/deploy/customizing-default-settings/#using-helm)
> * For Kubernetes < v1.25, if your cluster still enables Pod Security Policy admission controller, set the helm value `enablePSP` to `true` to install `longhorn-psp` PodSecurityPolicy resource which allows privileged Longhorn pods to start.


1. Create a HelmChart yaml file similar to this:

    ```yaml
   apiVersion: helm.cattle.io/v1
    kind: HelmChart
    metadata:
      annotations:
        helmcharts.cattle.io/managed-by: helm-controller
      finalizers:
      - wrangler.cattle.io/on-helm-chart-remove
      generation: 1
      name: longhorn-install
      namespace: default
    spec:
      version: v{{< current-version >}}
      chart: longhorn
      repo: https://charts.longhorn.io
      failurePolicy: abort
      targetNamespace: longhorn-system
      createNamespace: true

    ```

    > **IMPORTANT!** Ensure that `spec.failurePolicy` is set to "abort".  The only other value is the default: "reinstall", which performs an uninstall of Longhorn.  With "abort", it retries periodically, giving the user a chance to fix the problem.

    > **Note:** Rather than specify the repo, version, and chart name, the yaml can also use an image of the charts themselves:
    ```yaml
    spec:
      chartContent:  <tarball of chart directory | base64 -w 0>
    ```
    > For full details see the HelmChart controller docs:  https://docs.rke2.io/helm or https://docs.k3s.io/helm.

2. Apply it to create the HelmChart CR and an installation job:

    ```shell
    $ kubectl apply -f helmchart_repo_install.yaml
    helmchart.helm.cattle.io/longhorn-install created

    ```

    > **Note:** Deleting the helmchart CR will initiate an uninstall of Longhorn.

3. To show the created resources:

    ```shell
    $ kubectl get jobs
    NAME                            COMPLETIONS   DURATION   AGE
    helm-install-longhorn-install   0/1           8s         8s

    $ kubectl get pods
    NAME                                  READY   STATUS      RESTARTS   AGE
    helm-install-longhorn-install-lngm8   0/1     Completed   0          25s

    $ kubectl get helmcharts
    NAME               JOB                     CHART      TARGETNAMESPACE   VERSION   REPO                         HELMVERSION   BOOTSTRAP
    longhorn-install   helm-install-longhorn   longhorn   longhorn-system   v{{< current-version >}}    https://charts.longhorn.io

    ```

4. To confirm that the deployment succeeded, run:

    ```bash
    kubectl -n longhorn-system get pod
    ```

    The result should look like the following:

    ```bash
    NAME                                                READY   STATUS    RESTARTS      AGE
    csi-attacher-85c7684cfd-67kqc                       1/1     Running   0             29m
    csi-attacher-85c7684cfd-jbddj                       1/1     Running   0             29m
    csi-attacher-85c7684cfd-t85bw                       1/1     Running   0             29m
    csi-provisioner-68cdb8b96-46d9q                     1/1     Running   0             29m
    csi-provisioner-68cdb8b96-dgf5f                     1/1     Running   0             29m
    csi-provisioner-68cdb8b96-mh8q7                     1/1     Running   0             29m
    csi-resizer-86dd765b9-d27cs                         1/1     Running   0             29m
    csi-resizer-86dd765b9-scqxm                         1/1     Running   0             29m
    csi-resizer-86dd765b9-zpcv7                         1/1     Running   0             29m
    csi-snapshotter-65b46b8749-dtvh2                    1/1     Running   0             29m
    csi-snapshotter-65b46b8749-g67fn                    1/1     Running   0             29m
    csi-snapshotter-65b46b8749-nfgzm                    1/1     Running   0             29m
    engine-image-ei-221c9c21-gd5d6                      1/1     Running   0             29m
    engine-image-ei-221c9c21-v6clp                      1/1     Running   0             29m
    engine-image-ei-221c9c21-zzdrt                      1/1     Running   0             29m
    instance-manager-77d11dda6091967f9b30011c9876341b   1/1     Running   0             29m
    instance-manager-870c250b69a4fe01382ed46156d33f47   1/1     Running   0             29m
    instance-manager-a4099c5ce28b423c3cc2667906f4b0b4   1/1     Running   0             29m
    longhorn-csi-plugin-jfbh5                           3/3     Running   0             29m
    longhorn-csi-plugin-w768w                           3/3     Running   0             29m
    longhorn-csi-plugin-xcghm                           3/3     Running   0             29m
    longhorn-driver-deployer-586bc86bf9-bkwk6           1/1     Running   0             30m
    longhorn-manager-c4xtv                              1/1     Running   1 (30m ago)   30m
    longhorn-manager-kgqts                              1/1     Running   0             30m
    longhorn-manager-n8xdr                              1/1     Running   0             30m
    longhorn-ui-69667f9678-2lvxn                        1/1     Running   0             30m
    longhorn-ui-69667f9678-2xmc9                        1/1     Running   0             30m

    ```

5. To enable access to the Longhorn UI, you need to set up an Ingress controller. Authentication to the Longhorn UI is not enabled by default. For information on creating an NGINX Ingress controller with basic authentication, refer to [this section.](../../accessing-the-ui/longhorn-ingress)

6. Access the Longhorn UI using [these steps.](../../accessing-the-ui)



---

## Article: deploy/install/install-with-helm.md

---
title: Install with Helm
weight: 9
---

In this section, you will learn how to install Longhorn with Helm.

### Prerequisites

- Kubernetes cluster: Ensure that each node fulfills the [installation requirements](../#installation-requirements).
- Your workstation: Install [Helm](https://helm.sh/docs/) v3.0 or later.

> [Longhorn Command Line Tool](../../../advanced-resources/longhornctl/) can be used to check the Longhorn environment for potential issues.

### Installing Longhorn


> **Note**:
> * The initial settings for Longhorn can be found in [customized using Helm options or by editing the deployment configuration file.](../../../advanced-resources/deploy/customizing-default-settings/#using-helm)
> * For Kubernetes < v1.25, if your cluster still enables Pod Security Policy admission controller, set the helm value `enablePSP` to `true` to install `longhorn-psp` PodSecurityPolicy resource which allows privileged Longhorn pods to start.


1. Add the Longhorn Helm repository:

    ```shell
   helm repo add longhorn https://charts.longhorn.io
    ```

2. Fetch the latest charts from the repository:

    ```shell
   helm repo update
    ```

3. Install Longhorn in the `longhorn-system` namespace.

    ```shell
    helm install longhorn longhorn/longhorn --namespace longhorn-system --create-namespace --version {{< current-version >}}
    ```

4. To confirm that the deployment succeeded, run:

    ```bash
    kubectl -n longhorn-system get pod
    ```

    The result should look like the following:

    ```bash
    NAME                                                READY   STATUS    RESTARTS   AGE
    longhorn-ui-b7c844b49-w25g5                         1/1     Running   0          2m41s
    longhorn-manager-pzgsp                              1/1     Running   0          2m41s
    longhorn-driver-deployer-6bd59c9f76-lqczw           1/1     Running   0          2m41s
    longhorn-csi-plugin-mbwqz                           2/2     Running   0          100s
    csi-snapshotter-588457fcdf-22bqp                    1/1     Running   0          100s
    csi-snapshotter-588457fcdf-2wd6g                    1/1     Running   0          100s
    csi-provisioner-869bdc4b79-mzrwf                    1/1     Running   0          101s
    csi-provisioner-869bdc4b79-klgfm                    1/1     Running   0          101s
    csi-resizer-6d8cf5f99f-fd2ck                        1/1     Running   0          101s
    csi-provisioner-869bdc4b79-j46rx                    1/1     Running   0          101s
    csi-snapshotter-588457fcdf-bvjdt                    1/1     Running   0          100s
    csi-resizer-6d8cf5f99f-68cw7                        1/1     Running   0          101s
    csi-attacher-7bf4b7f996-df8v6                       1/1     Running   0          101s
    csi-attacher-7bf4b7f996-g9cwc                       1/1     Running   0          101s
    csi-attacher-7bf4b7f996-8l9sw                       1/1     Running   0          101s
    csi-resizer-6d8cf5f99f-smdjw                        1/1     Running   0          101s
    instance-manager-b34d5db1fe1e2d52bcfb308be3166cfc   1/1     Running   0          114s
    engine-image-ei-df38d2e5-cv6nc                      1/1     Running   0          114s
    ```

5. To enable access to the Longhorn UI, you will need to set up an Ingress controller. Authentication to the Longhorn UI is not enabled by default. For information on creating an NGINX Ingress controller with basic authentication, refer to [this section.](../../accessing-the-ui/longhorn-ingress)

6. Access the Longhorn UI using [these steps.](../../accessing-the-ui)


---

## Article: deploy/install/install-with-kubectl.md

---
title: Install with Kubectl
description: Install Longhorn with the kubectl client.
weight: 8
---

## Prerequisites

Each node in the Kubernetes cluster where Longhorn will be installed must fulfill [these requirements.](../#installation-requirements)

[Longhorn Command Line Tool](../../../advanced-resources/longhornctl/) can be used to check the Longhorn environment for potential issues.

The initial settings for Longhorn can be customized by [editing the deployment configuration file.](../../../advanced-resources/deploy/customizing-default-settings/#using-the-longhorn-deployment-yaml-file)

## Installing Longhorn

1. Install Longhorn on any Kubernetes cluster using this command:

    ```shell
    kubectl apply -f https://raw.githubusercontent.com/longhorn/longhorn/v{{< current-version >}}/deploy/longhorn.yaml
    ```

    One way to monitor the progress of the installation is to watch pods being created in the `longhorn-system` namespace:

    ```shell
    kubectl get pods \
    --namespace longhorn-system \
    --watch
    ```

2. Check that the deployment was successful:

    ```shell
    $ kubectl -n longhorn-system get pod
    NAME                                                READY   STATUS    RESTARTS   AGE
    longhorn-ui-b7c844b49-w25g5                         1/1     Running   0          2m41s
    longhorn-manager-pzgsp                              1/1     Running   0          2m41s
    longhorn-driver-deployer-6bd59c9f76-lqczw           1/1     Running   0          2m41s
    longhorn-csi-plugin-mbwqz                           2/2     Running   0          100s
    csi-snapshotter-588457fcdf-22bqp                    1/1     Running   0          100s
    csi-snapshotter-588457fcdf-2wd6g                    1/1     Running   0          100s
    csi-provisioner-869bdc4b79-mzrwf                    1/1     Running   0          101s
    csi-provisioner-869bdc4b79-klgfm                    1/1     Running   0          101s
    csi-resizer-6d8cf5f99f-fd2ck                        1/1     Running   0          101s
    csi-provisioner-869bdc4b79-j46rx                    1/1     Running   0          101s
    csi-snapshotter-588457fcdf-bvjdt                    1/1     Running   0          100s
    csi-resizer-6d8cf5f99f-68cw7                        1/1     Running   0          101s
    csi-attacher-7bf4b7f996-df8v6                       1/1     Running   0          101s
    csi-attacher-7bf4b7f996-g9cwc                       1/1     Running   0          101s
    csi-attacher-7bf4b7f996-8l9sw                       1/1     Running   0          101s
    csi-resizer-6d8cf5f99f-smdjw                        1/1     Running   0          101s
    instance-manager-b34d5db1fe1e2d52bcfb308be3166cfc   1/1     Running   0          114s
    engine-image-ei-df38d2e5-cv6nc                      1/1     Running   0          114s
    ```
3. To enable access to the Longhorn UI, you will need to set up an Ingress controller. Authentication to the Longhorn UI is not enabled by default. For information on creating an NGINX Ingress controller with basic authentication, refer to [this section.](../../accessing-the-ui/longhorn-ingress)
4. Access the Longhorn UI using [these steps.](../../accessing-the-ui)

> **Note**:
> For Kubernetes < v1.25, if your cluster still enables Pod Security Policy admission controller, need to apply the [podsecuritypolicy.yaml](https://raw.githubusercontent.com/longhorn/longhorn/master/deploy/podsecuritypolicy.yaml) manifest in addition to applying the `longhorn.yaml` manifests.



### List of Deployed Resources


The following items will be deployed to Kubernetes:

#### Namespace: longhorn-system

All Longhorn bits will be scoped to this namespace.

#### ServiceAccount: longhorn-service-account

Service account is created in the longhorn-system namespace.

#### ClusterRole: longhorn-role

This role will have access to:
  - In apiextension.k8s.io (All verbs)
    - customresourcedefinitions
  - In core (All verbs)
    - pods
      - /logs
    - events
    - persistentVolumes
    - persistentVolumeClaims
      - /status
    - nodes
    - proxy/nodes
    - secrets
    - services
    - endpoints
    - configMaps
  - In core
    - namespaces (get, list)
  - In apps (All Verbs)
    - daemonsets
    - statefulSets
    - deployments
  - In batch (All Verbs)
    - jobs
    - cronjobs
  - In storage.k8s.io (All verbs)
    - storageclasses
    - volumeattachments
    - csinodes
    - csidrivers
  - In coordination.k8s.io
    - leases

#### ClusterRoleBinding: longhorn-bind

This connects the longhorn-role to the longhorn-service-account in the  longhorn-system namespace

#### CustomResourceDefinitions

The following CustomResourceDefinitions will be installed

- In longhorn.io
  - backingimagedatasources
  - backingimagemanagers
  - backingimages
  - backups
  - backuptargets
  - backupvolumes
  - engineimages
  - engines
  - instancemanagers
  - nodes
  - recurringjobs
  - replicas
  - settings
  - sharemanagers
  - volumes

#### Kubernetes API Objects

- A config map with the default settings
- The longhorn-manager DaemonSet
- The longhorn-backend service exposing the longhorn-manager DaemonSet internally to Kubernetes
- The longhorn-ui Deployment
- The longhorn-frontend service exposing the longhorn-ui internally to Kubernetes
- The longhorn-driver-deployer that deploys the CSI driver
- The longhorn StorageClass



---

## Article: deploy/install/install-with-rancher.md

---
title: Install as a Rancher Apps & Marketplace
description: Run Longhorn on Kubernetes with Rancher 2.x
weight: 7
---

One benefit of installing Longhorn through Rancher Apps & Marketplace is that Rancher provides authentication to the Longhorn UI.

If there is a new version of Longhorn available, you will see an `Upgrade Available` sign on the `Apps & Marketplace` page. You can click `Upgrade` button to upgrade Longhorn manager. See more about upgrade [here](../../upgrade).

## Prerequisites

Each node in the Kubernetes cluster where Longhorn is installed must fulfill [these requirements.](../#installation-requirements)

[Longhorn Command Line Tool](../../../advanced-resources/longhornctl/) can be used to check the Longhorn environment for potential issues.

## Installation

> **Note**:
> * For Kubernetes < v1.25, if your cluster still enables Pod Security Policy admission controller, set `Other Settings > Pod Security Policy` to `true` to install `longhorn-psp` PodSecurityPolicy resource which allows privileged Longhorn pods to start.

1. Optional: If Rancher version is 2.5.9 or before, we recommend creating a new project for Longhorn, for example, `Storage`.
2. Navigate to the cluster where you will install Longhorn.
    {{< figure src="/img/screenshots/install/rancher-2.6/select-project.png" >}}
3. Navigate to the `Apps & Marketplace` page.
    {{< figure src="/img/screenshots/install/rancher-2.6/apps-launch.png" >}}
4. Find the Longhorn item in the charts and click it.
    {{< figure src="/img/screenshots/install/rancher-2.6/longhorn.png" >}}
5. Click **Install**.
    {{< figure src="/img/screenshots/install/rancher-2.6/longhorn-chart.png" >}}
6. Optional: Select the project where you want to install Longhorn.
7. Optional: Customize the default settings.
    {{< figure src="/img/screenshots/install/rancher-2.6/launch-longhorn.png" >}}
8. Click Next. Longhorn will be installed in the longhorn-system namespace.
    {{< figure src="/img/screenshots/install/rancher-2.6/installed-longhorn.png" >}}
9. Click the Longhorn App Icon to navigate to the Longhorn dashboard.
    {{< figure src="/img/screenshots/install/rancher-2.6/dashboard.png" >}}

After Longhorn has been successfully installed, you can access the Longhorn UI by navigating to the `Longhorn` option from Rancher left panel.


## Access UI With Network Policy Enabled

Note that when the Network Policy is enabled, access to the UI from Rancher may be restricted.

Rancher interacts with the Longhorn UI via a service called remotedialer, which facilitates connections between Rancher and the downstream clusters it manages. This service allows a user agent to access the cluster through an endpoint on the Rancher server. Remotedialer connects to the Longhorn UI service by using the Kubernetes API Server as a proxy.

However, when the Network Policy is enabled, the Kubernetes API Server may be unable to reach pods on different nodes. This occurs because the Kubernetes API Server operates within the host’s network namespace without a dedicated per-pod IP address. If you're using the Calico CNI plugin, any process in the host’s network namespace (such as the API Server) connecting to a pod triggers Calico to encapsulate the packet in IPIP before forwarding it to the remote host. The tunnel address is chosen as the source to ensure the remote host knows to encapsulate the return packets correctly.

In other words, to allow the proxy to work with the Network Policy, the Tunnel IP of each node must be identified and explicitly permitted in the policy.

You can find the Tunnel IP by:
```
$ kubectl get nodes -oyaml | grep "Tunnel"

      projectcalico.org/IPv4VXLANTunnelAddr: 10.42.197.0
      projectcalico.org/IPv4VXLANTunnelAddr: 10.42.99.0
      projectcalico.org/IPv4VXLANTunnelAddr: 10.42.158.0
      projectcalico.org/IPv4VXLANTunnelAddr: 10.42.80.0
```

Next, permit traffic in the Network Policy using the Tunnel IP. You may need to update the Network Policy whenever new nodes are added to the cluster.
```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: longhorn-ui-frontend
  namespace: longhorn-system
spec:
  podSelector:
    matchLabels:
      app: longhorn-ui
  policyTypes:
  - Ingress
  ingress:
  - from:
    - ipBlock:
        cidr: 10.42.197.0/32
    - ipBlock:
        cidr: 10.42.99.0/32
    - ipBlock:
        cidr: 10.42.158.0/32
    - ipBlock:
        cidr: 10.42.80.0/32
    ports:
      - port: 8000
        protocol: TCP
```

Another way to resolve the issue is by running the server nodes with `egress-selector-mode: cluster`. For more information, see [RKE2 Server Configuration Reference](https://docs.rke2.io/reference/server_config#critical-configuration-values) and [K3s Control-Plane Egress Selector configuration](https://docs.k3s.io/networking/basic-network-options#control-plane-egress-selector-configuration).

---

## Article: troubleshoot/_index.md

---
title: Troubleshoot
weight: 10
---

---

## Article: troubleshoot/support-bundle.md

---
title: Support Bundle
weight: 2
---

Since v1.4.0, Longhorn replaced the in-house support bundle generation with a general-purpose [support bundle kit](https://github.com/rancher/support-bundle-kit). 

You can click the `Generate Support Bundle` at the bottom of Longhorn UI to download a zip file containing cluster manifests and logs.

During support bundle generation, Longhorn will create a Deployment for the support bundle manager.

> **Note:** The support bundle manager will use a dedicated `longhorn-support-bundle` service account and `longhorn-support-bundle` cluster role binding with `cluster-admin` access for bundle collection.

With the support bundle, you can simulate a mocked Kubernetes cluster that is interactable with the `kubectl` command. See [simulator command](https://github.com/rancher/support-bundle-kit#simulator-command) for more details.


## Limitations

Longhorn currently does not support concurrent generation of multiple support bundles. We recommend waiting until the completion of the ongoing support bundle before initiating a new one. If a new support bundle is created while another one is still in progress, Longhorn will overwrite the older support bundle.


## History
[Original Feature Request](https://github.com/longhorn/longhorn/issues/2759)

Available since v1.4.0


---

## Article: troubleshoot/troubleshooting.md

---
title: Troubleshooting Problems
weight: 1
---
- [Troubleshooting Guide](#troubleshooting-guide)
  - [UI](#ui)
  - [Manager and Engines](#manager-and-engines)
  - [CSI driver](#csi-driver)
  - [Flexvolume Driver](#flexvolume-driver)
- [Common issues](#common-issues)
  - [Volume can be attached/detached from UI, but Kubernetes Pod/StatefulSet etc cannot use it](#volume-can-be-attacheddetached-from-ui-but-kubernetes-podstatefulset-etc-cannot-use-it)
    - [Using with Flexvolume Plugin](#using-with-flexvolume-plugin)

## Troubleshooting Guide

> For a more in-depth troubleshooting flow please see https://github.com/longhorn/longhorn/wiki/Troubleshooting

There are a few components in Longhorn: Manager, Engine, Driver and UI. By default, all of those components run as pods in the `longhorn-system` namespace in the Kubernetes cluster.

Most of the logs are included in the Support Bundle. You can click the **Generate Support Bundle** link at the bottom of the UI to download a zip file that contains Longhorn-related configuration and logs.
See [Support Bundle](../support-bundle) for detail.

One exception is the `dmesg`, which needs to be retrieved from each node by the user.

### UI
Make use of the Longhorn UI is a good start for the troubleshooting. For example, if Kubernetes cannot mount one volume correctly, after stop the workload, try to attach and mount that volume manually on one node and access the content to check if volume is intact.

Also, the event logs in the UI dashboard provides some information of probably issues. Check for the event logs in `Warning` level.

### Manager and Engines
You can get the logs from the Longhorn Manager and Engines to help with troubleshooting. The most useful logs are the ones from `longhorn-manager-xxx`, and the logs inside Longhorn instance managers, e.g. `instance-manager-xxxx`, `instance-manager-e-xxxx` and `instance-manager-r-xxxx`.

Since normally there are multiple Longhorn Managers running at the same time, we recommend using [kubetail,](https://github.com/johanhaleby/kubetail) which is a great tool to keep track of the logs of multiple pods. To track the manager logs in real time, you can use:

```
kubetail longhorn-manager -n longhorn-system
```


### CSI driver

For the CSI driver, check the logs for `csi-attacher-0` and `csi-provisioner-0`, as well as containers in `longhorn-csi-plugin-xxx`.

### Flexvolume Driver

The FlexVolume driver is deprecated as of Longhorn v0.8.0 and should no longer be used.

First check where the driver has been installed on the node. Check the log of `longhorn-driver-deployer-xxxx` for that information.

Then check the kubelet logs. The FlexVolume driver itself doesn't run inside the container. It would run along with the kubelet process.

If kubelet is running natively on the node, you can use the following command to get the logs:
```
journalctl -u kubelet
```

Or if kubelet is running as a container (e.g. in RKE), use the following command instead:
```
docker logs kubelet
```

For even more detailed logs of Longhorn FlexVolume, run the following command on the node or inside the container (if kubelet is running as a container, e.g. in RKE):
```
touch /var/log/longhorn_driver.log
```


## Common issues
### Volume can be attached/detached from UI, but Kubernetes Pod/StatefulSet etc cannot use it

#### Using with Flexvolume Plugin
Check if the volume plugin directory has been set correctly. This is automatically detected unless the user explicitly sets it.

By default, Kubernetes uses `/usr/libexec/kubernetes/kubelet-plugins/volume/exec/`, as stated in the [official document](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-storage/flexvolume.md/#prerequisites).

Some vendors choose to change the directory for various reasons. For example, GKE uses `/home/kubernetes/flexvolume` instead.

The correct directory can be found by running `ps aux|grep kubelet` on the host and check the `--volume-plugin-dir` parameter. If there is none, the default `/usr/libexec/kubernetes/kubelet-plugins/volume/exec/` will be used.

## Profiling

### Engine, replica, and sync agent runtime

You can enable the `pprof` server dynamically to perform runtime profiling.
To enable profiling, you can:

1. Shell into the instance manager pod.
2. Identify the runtime process and its port using `ps`:
    ```bash
    $ ps aux | more

    USER         PID %CPU %MEM    VSZ   RSS TTY      STAT START   TIME COMMAND
    ...
    root        1996  0.0  0.6 1990080 20996 ?       Sl   Jul25   0:05 /host/var/lib/longhorn/engine-binaries/longhornio-longhorn-engine-v1.10.0/longhorn --volume-name     vol replica /host/var/lib/longhorn/replicas/vol-3004fc59 --size 1073741824 --disableRevCounter --replica-instance-name vol-r-ec7e35e4 --snapshot-max-count 250     --snapshot-max-size 0 --sync-agent-port-count 7 --listen 0.0.0.0:10000
    root        2004  0.0  0.6 1695152 22708 ?       Sl   Jul25   0:09 /host/var/lib/longhorn/engine-binaries/longhornio-longhorn-engine-v1.10.0/longhorn --volume-name     vol sync-agent --listen 0.0.0.0:10002 --replica 0.0.0.0:10000 --listen-port-range 10003-10009 --replica-instance-name vol-r-ec7e35e4
    root        2031  0.0  0.6 1916348 23760 ?       Sl   Jul25   0:46 /engine-binaries/longhornio-longhorn-engine-v1.10.0/longhorn --engine-instance-name vol-e-0     controller vol --frontend tgt-blockdev --disableRevCounter --size 1073741824 --current-size 0 --engine-replica-timeout 8 --file-sync-http-client-timeout 30     --snapshot-max-count 250 --snapshot-max-size 0 --replica tcp://10.42.2.7:10000 --replica tcp://10.42.0.15:10000 --replica tcp://10.42.1.7:10000 --listen 0.0.0.0:10010
    ```
3. Enable the `pprof` server for the desired runtime (for example, sync-agent):
    > In this example, the sync-agent process listens on port `10002`.
    ```bash
    $ /host/var/lib/longhorn/engine-binaries/longhornio-longhorn-engine-v1.10.0/longhorn --url http://localhost:10002 profiler enable --port 36060
    $ /host/var/lib/longhorn/engine-binaries/longhornio-longhorn-engine-v1.10.0/longhorn --url http://localhost:10002 profiler show

    Profiler enabled at Addr: *:36060
    ```
    > The `pprof` server is now accessible at `http://localhost:36060` *inside the instance manager pod*.
4. Use the `pprof` interface for runtime inspection. For more details, refer to the [official pprof documentation](https://pkg.go.dev/net/http/pprof#hdr-Usage_examples).
5. Disable the profiler after completing your analysis:
    ```bash
    $ /host/var/lib/longhorn/engine-binaries/longhornio-longhorn-engine-v1.10.0/longhorn --url http://localhost:10002 profiler disable

    Profiler is disabled!
    ```


---

## Article: snapshots-and-backups/_index.md

---
title: Backup and Restore
description: Backup and Restore Volume Snapshots in Longhorn 
weight: 6
---

---

## Article: snapshots-and-backups/csi-volume-clone.md

---
title: Volume Clone Support
description: Creating a new volume as a duplicate of an existing volume
weight: 3
---

Longhorn supports [CSI volume cloning](https://kubernetes.io/docs/concepts/storage/volume-pvc-datasource/).


## Volume Cloning

### Clone a Volume Using YAML
Suppose that you have the following `source-pvc`:
```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: source-pvc
spec:
  storageClassName: longhorn
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
```
You can create a new PVC that has the exact same content as the `source-pvc` by applying the following yaml file:
```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: cloned-pvc
spec:
  storageClassName: longhorn
  dataSource:
    name: source-pvc
    kind: PersistentVolumeClaim
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
```

> Note:
> In addition to the requirements listed at [CSI volume cloning](https://kubernetes.io/docs/concepts/storage/volume-pvc-datasource/),
> the `cloned-pvc` must have the same `resources.requests.storage` as the `source-pvc`.


### Clone Volume Using the Longhorn UI

#### Clone a volume

1. Go to the **Volumes** page.
2. Select a volume, and then click **Clone Volume** in the **Operation** menu.
3. (Optional) Configure the settings of the new volume.
4. Click **OK**.

#### Clone a Volume Using a Snapshot

1. Go to the **Volumes** page.
2. Click the name of the volume that you want to clone.
3. In the **Snapshot and Backups** section of the details page, identify the snapshot that you want to use and then click **Clone Volume**.
4. (Optional) Configure the settings of the new volume.
5. Click **OK**.

{{< figure src="/img/screenshots/snapshots-and-backups/clone-volume-modal.png" >}}

#### Clone Multiple Volumes (Bulk Cloning)

1. Go to the **Volumes** page.
2. Select the volume you want to clone.
3. Click **Clone Volume** button on top of the table.
4. (Optional) Configure the settings of the new volumes
5. Click **OK**


**Note**:
> - The Longhorn UI pre-fills certain fields and prevents you from modifying the values to ensure that those match the settings of the source volume.
> - Longhorn automatically attaches the new volume, clones the source volume, and then detaches the new volume.


## Volume Creation

1. Go to the **Volumes** page.
2. Click **Create Volume**.
3. Select the data source (**Volume** or **Volume Snapshot**) that you want to use.
4. If you select **Volume Snapshot**, choose a snapshot.
5. Specify the volume name.
6. Click **OK**.

{{< figure src="/img/screenshots/snapshots-and-backups/create-volume-choose-datasource.png" >}}

## History

- [GitHub Issue](https://github.com/longhorn/longhorn/issues/1815)
- [Longhorn Enhancement Proposal](https://github.com/longhorn/longhorn/pull/2864)

Available since v1.2.0


---

## Article: snapshots-and-backups/scheduling-backups-and-snapshots.md

---
title: Recurring Snapshots and Backups
weight: 3
---

From the Longhorn UI, the volume can refer to recurring snapshots and backups as independent jobs or as recurring job groups.

To create a recurring job, you can go to the `Recurring Jobs` page in Longhorn and `Create Recurring Job` or in the volume detail view in Longhorn.

You can configure,
- Any groups that the job should belong to
- The type of schedule, either `backup`, `backup-force-create`, `snapshot`, `snapshot-force-create`, `snapshot-cleanup`, `snapshot-delete` or `filesystem-trim`
- The time that the backup or snapshot will be created, in the form of a [CRON expression](https://en.wikipedia.org/wiki/Cron#CRON_expression)
- The number of backups or snapshots to retain
- The number of jobs to run concurrently
- Any labels that should be applied to the backup or snapshot
- Parameters that should be applied to the backup
  - `full-backup-interval`: Number of incremental backups that must be completed before Longhorn performs a full backup. This integer parameter is applied only to the backup. Notice that if the value is 0, Longhorn performs a incremental backup every time. For more information, see [Periodic Full Backup](#periodic-full-backup) and [Create a Backup](../backup-and-restore/create-a-backup).

Recurring jobs can be set up using the Longhorn UI, `kubectl`, or by using a Longhorn `RecurringJob`.

To add a recurring job to a volume, you will go to the volume detail view in Longhorn. Then you can set `Recurring Jobs Schedule`.

- Create a new recurring job
- Select from existing recurring jobs
- Select from existing recurring job groups

Then Longhorn will automatically create snapshots or backups for the volume at the recurring job scheduled time, as long as the volume is attached to a node.
If you want to set up recurring snapshots and backups even when the volumes are detached, see the section [Allow Recurring Job While Volume Is Detached](#allow-recurring-job-while-volume-is-detached)

You can set recurring jobs on a Longhorn Volume, Kubernetes Persistent Volume Claim (PVC), or Kubernetes StorageClass.
> Note: When the PVC has recurring job labels, they will override all recurring job labels of the associated Volume.

For more information on how snapshots and backups work, refer to the [concepts](../../concepts) section.

> Note: To avoid the problem that recurring jobs may overwrite the old backups/snapshots with identical backups and empty snapshots when the volume doesn't have new data for a long time, Longhorn does the following:
> 1. Recurring backup job only takes a new backup when the volume has new data since the last backup.
> 1. Recurring snapshot job only takes a new snapshot when the volume has new data in the volume head (the live data).

## Set up Recurring Jobs

### Using the Longhorn UI

Recurring snapshots and backups can be configured from the `Recurring Job` page or the volume detail page.

### Using the manifest

You can also configure the recurring job by directly interacting with the Longhorn RecurringJob custom resource.
```yaml
apiVersion: longhorn.io/v1beta2
kind: RecurringJob
metadata:
  name: snapshot-1
  namespace: longhorn-system
spec:
  cron: "* * * * *"
  task: "snapshot"
  groups:
  - default
  - group1
  retain: 1
  concurrency: 2
  labels:
    label/1: a
    label/2: b
```

The following parameters should be specified for each recurring job selector:

- `name`: Name of the recurring job. Do not use duplicate names. And the length of `name` should be no more than 40 characters.

- `task`: Type of the job. Longhorn supports the following:
  - `backup`: periodically create snapshots then do backups after cleaning up outdated snapshots
  - `backup-force-create`: periodically create snapshots then do backups
  - `snapshot`: periodically create snapshots after cleaning up outdated snapshots
  - `snapshot-force-create`: periodically create snapshots
  - `snapshot-cleanup`: periodically purge removable snapshots and system snapshots
    > **Note:** retain value has no effect for this task, Longhorn automatically mutates the `retain` value to 0.

  - `snapshot-delete`: periodically remove and purge all kinds of snapshots that exceed the retention count.
    > **Note:** The `retain` value is independent of each recurring job.
    >
    > Using a volume with 2 recurring jobs as an example:
    > - `snapshot` with retain value set to 5
    > - `snapshot-delete`: with retain value set to 2
    >
    > Eventually, there will be 2 snapshots retained after a complete `snapshot-delete` task execution.

  - `filesystem-trim`: periodically trim filesystem to reclaim disk space

- `cron`: Cron expression. It tells the execution time of the job.

- `retain`: How many snapshots/backups Longhorn will retain for each volume job. It should be no less than 1.

- `concurrency`: The number of jobs to run concurrently. It should be no less than 1.

Optional parameters can be specified:

- `groups`: Any groups that the job should belong to. Having `default` in groups will automatically schedule this recurring job to any volume with no recurring job.

- `labels`: Any labels that should be applied to the backup or snapshot.

## Add Recurring Jobs to the Default group

Default recurring jobs can be set by tick the checkbox `default` using UI or adding `default` to the recurring job `groups`.

Longhorn will automatically add a volume to the `default` group when the volume has no recurring job.

## Delete Recurring Jobs

Longhorn automatically removes Volume and PVC recurring job labels when a corresponding RecurringJob custom resource is deleted. However, if a recurring job label is added without an existing RecurringJob custom resource, Longhorn does not perform the cleanup process for that label.

## Apply Recurring Job to Longhorn Volume

### Using the Longhorn UI

The recurring job can be assigned on the volume detail page. To navigate to the volume detail page, click **Volume** then click the name of the volume.

## Using the `kubectl` command

Add recurring job group:
```
kubectl -n longhorn-system label volume/<VOLUME-NAME> recurring-job-group.longhorn.io/<RECURRING-JOB-GROUP-NAME>=enabled

# Example:
# kubectl -n longhorn-system label volume/pvc-8b9cd514-4572-4eb2-836a-ed311e804d2f recurring-job-group.longhorn.io/default=enabled
```

Add recurring job:
```
kubectl -n longhorn-system label volume/<VOLUME-NAME> recurring-job.longhorn.io/<RECURRING-JOB-NAME>=enabled

# Example:
# kubectl -n longhorn-system label volume/pvc-8b9cd514-4572-4eb2-836a-ed311e804d2f recurring-job.longhorn.io/backup=enabled
```

Remove recurring job:
```
kubectl -n longhorn-system label volume/<VOLUME-NAME> <RECURRING-JOB-LABEL>-

# Example:
# kubectl -n longhorn-system label volume/pvc-8b9cd514-4572-4eb2-836a-ed311e804d2f recurring-job.longhorn.io/backup-
```

## With PersistentVolumeClaim Using the `kubectl` command

By default, applying a recurring job to a Persistent Volume Claim (PVC) does not have any effect. You can enable or disable this feature using the recurring job source label.

Once the PVC is labeled as the source, any recurring job labels added or removed from the PVC will be periodically synchronized by Longhorn to the associated Volume.
```
kubectl -n <NAMESPACE> label pvc/<PVC-NAME> recurring-job.longhorn.io/source=enabled

# Example:
# kubectl -n default label pvc/sample recurring-job.longhorn.io/source=enabled
```

Add recurring job group:
```
kubectl -n <NAMESPACE> label pvc/<PVC-NAME> recurring-job-group.longhorn.io/<RECURRING-JOB-GROUP-NAME>=enabled

# Example:
# kubectl -n default label pvc/sample recurring-job-group.longhorn.io/default=enabled
```

Add recurring job:
```
kubectl -n <NAMESPACE> label pvc/<PVC-NAME> recurring-job.longhorn.io/<RECURRING-JOB-NAME>=enabled

# Example:
# kubectl -n default label pvc/sample recurring-job.longhorn.io/backup=enabled
```

Remove recurring job:
```
kubectl -n <NAMESPACE> label pvc/<PVC-NAME> <RECURRING-JOB-LABEL>-

# Example:
# kubectl -n default label pvc/sample recurring-job.longhorn.io/backup-
```

## With StorageClass parameters

Recurring job assignment can be configured in the `recurringJobSelector` parameters in a StorageClass.

Any future volumes created using this StorageClass will have those recurring jobs automatically assigned.

The `recurringJobSelector` field should follow JSON format:
```yaml
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: longhorn
provisioner: driver.longhorn.io
parameters:
  numberOfReplicas: "3"
  staleReplicaTimeout: "30"
  fromBackup: ""
  recurringJobSelector: '[
    {
      "name":"snap",
      "isGroup":true
    },
    {
      "name":"backup",
      "isGroup":false
    }
  ]'
```

The following parameters should be specified for each recurring job selector:

1. `name`: Name of an existing recurring job or an existing recurring job group.

2. `isGroup`: is the name that belongs to a recurring job or recurring job group, either `true` or `false`.


## Allow Recurring Job While Volume Is Detached

Longhorn provides the setting `allow-recurring-job-while-volume-detached` that allows you to do recurring backup even when a volume is detached.
You can find the setting in Longhorn UI.

When the setting is enabled, Longhorn will automatically attach the volume and take a snapshot/backup when it is time to do a recurring snapshot/backup.

Note that during the time the volume was attached automatically, the volume is not ready for the workload. Workload will have to wait until the recurring job finishes.

## Periodic Full Backup

Longhorn performs delta backups by default, which means that only data that was changed since the last backup is uploaded. However, when a data block in the backupstore becomes corrupted, Longhorn does not replace that data block with a healthy one during subsequent backup operations. Corrupted data blocks in the backupstore may cause restoration operations to fail.
When a non-zero `full-backup-interval` parameter is set, Longhorn performs a full backup every `full-backup-interval` incremental backups. During a full backup, Longhorn uploads all data blocks in the volume. Data blocks that exist in the backupstore, including corrupted ones, are overwritten.

> **Important**: 
> Performing a full backup might take longer and generate higher network throughput and costs than the default incremental backup.


---

## Article: snapshots-and-backups/setup-a-snapshot.md

---
  title: Create a Snapshot
  weight: 1
---

A [snapshot](../../concepts/#24-snapshots) is the state of a Kubernetes volume at any given moment in time.

## Snapshot Management via UI

To create a snapshot of an existing cluster, follow these steps:

1. In the top navigation bar of the Longhorn UI, click **Volume**.
2. Click the name of the volume for which you want to create a snapshot. This leads to the volume detail page.
3. Click the **Take Snapshot** button.

Once the snapshot is created, you will see it in the list of snapshots for the volume before the Volume Head.

## Snapshot Management with Custom Resources (CRs)

This section demonstrates how to create, list, restore, and delete Longhorn snapshots directly via `kubectl` using Longhorn’s **Custom Resources (CRs)**.

> **Note**: Longhorn uses its own `Snapshot` CRD under the `longhorn.io` API group (for example, `v1beta2`), not the generic Kubernetes `VolumeSnapshot` from `snapshot.storage.k8s.io`.  

### Create a Snapshot

1.  **Prepare the manifest**: Create a file named `longhorn-snapshot.yaml`:

    ```yaml
    apiVersion: longhorn.io/v1beta2
    kind: Snapshot
    metadata:
      name: longhorn-test-snapshot
      namespace: longhorn-system
    spec:
      volume: pvc-840804d8-6f11-49fd-afae-54bc5be639de   # replace with your actual Longhorn volume name
      createSnapshot: true
    ```

2.  **Apply the manifest**:

    ```bash
    kubectl apply -f longhorn-snapshot.yaml
    ```

    Expected output:

    ```bash
    snapshot.longhorn.io/longhorn-test-snapshot created
    ```

    > **Note**: If the volume is detached, you will see a brief warning about the engine not running. Longhorn automatically retries, and the snapshot will complete once the volume is attached.

### List Snapshots

```bash
kubectl get snapshots.longhorn.io -l longhornvolume=pvc-840804d8-6f11-49fd-afae-54bc5be639de -n longhorn-system
```

### Delete a Snapshot

```bash
kubectl delete snapshot.longhorn.io longhorn-test-snapshot -n longhorn-system
```

Expected output:

```
snapshot.longhorn.io "longhorn-test-snapshot" deleted
```

> **Note**: Longhorn automatically handles the cleanup of the underlying data.


---

## Article: snapshots-and-backups/setup-disaster-recovery-volumes.md

---
title: Disaster Recovery (DR) Volumes
weight: 4
---

Ensuring data resilience is important when working with containerized applications.
A **Longhorn Disaster Recovery (DR) volume** is a special type of volume designed to maintain a standby copy of data in a secondary Kubernetes cluster. It is created from backups of a primary volume and kept in sync to enable rapid recovery if the primary cluster becomes unavailable.

The DR volume stores a geographically separated replica of the data. The backup frequency determines how current the DR volume is and, consequently, the potential amount of data loss in the event of a site failure.

## How It Works

The functionality of DR volumes relies on asynchronous replication through a shared backup store.

- **Shared Backup Target**

    Your primary and secondary Kubernetes clusters must be configured to use the exact same external backup target (e.g., an S3-compatible object store or an NFS share).

- **Incremental Backup and Restore**

    A DR volume is created from an existing backup. It continuously polls the backup target for newer backups from the source volume and restores them incrementally. The `Last Backup` field in the UI shows the most recent backup that has been restored.

    To keep the DR volume updated, configure recurring jobs on the source volume to perform regular incremental backups. These recurring backups provide the DR volume with new backups to restore, helping ensure minimal data loss in the event of a disaster.

- **Standby State**

    The DR volume remains in a passive standby state. It is not mounted or accessible by any workloads, which prevents data inconsistencies. The UI indicates the volume’s status with an icon:

    - Gray Icon: The volume is busy restoring data and cannot be activated.
    - Blue Icon: The volume is fully synchronized and ready for activation.

- **Activation**

    In a disaster, you manually activate the DR volume. This process converts it into a standard, writable Longhorn volume that can be attached to your applications in the recovery cluster.

## Creating a DR Volume

You can create a DR volume using either the Longhorn UI or `kubectl`.

> **Prerequisites**: Set up two Kubernetes clusters. These will be called cluster A and cluster B. Install Longhorn on both clusters, and set the same backup target on both clusters. For help setting the backup target, refer to [this page.](../backup-and-restore/set-backup-target)

### Using Longhorn UI

1. In your primary cluster, ensure the source volume has at least one backup.
2. In the Longhorn UI of your secondary (recovery) cluster, navigate to the Backup page.
3. Select the desired backup from the list and choose Create Disaster Recovery Volume. It is recommended to use the same name as the original volume.
4. Longhorn will create the volume, which will appear on the Volume page with a Standby status.

### Using `kubectl` Command

1. **Get the Backup URL**: First, copy the full URL of the source backup from the Backup page in the Longhorn UI. The format of this URL will depend on your configured backup target (e.g., S3, NFS).

2. **Create a YAML Manifest**: Create a file (e.g., dr-volume.yaml) with the following content. Be sure to replace the placeholder URL and adjust the name, size, accessMode, etc. to match your source volume. In this file, the `standby: true` field defines the volume as a DR standby volume.

    ```yaml
    apiVersion: longhorn.io/v1beta2
    kind: Volume
    metadata:
      name: example-dr-volume
      namespace: longhorn-system
    spec:
      size: "2147483648"
      accessMode: rwo
      numberOfReplicas: 3
      fromBackup: "nfs://longhorn-nfs-server.example.com:/opt/backupstore?backup=backup-b69a1249e97f4a27&volume=pvc-33509786-92d7-427c-9b5a-b6d61d56b063"
      # This flag is essential to create a standby volume
      Standby: true
    ```

3. **Apply the Manifest**: Apply the manifest to your **secondary cluster** to create the volume.

## Activating a DR Volume

When a failover is necessary, you must activate the DR volume to make it writable.  
Longhorn supports activation under the following conditions:

- The volume is healthy, indicating that all replicas are in a healthy state.
- The volume is degraded (some replicas have failed), but only if the global setting [`Allow Volume Creation with Degraded Availability`](../../references/settings/#allow-volume-creation-with-degraded-availability) is enabled.

> **Note**: When the setting `Allow Volume Creation with Degraded Availability` is disabled, attempting to activate a degraded DR volume will cause the volume to become stuck in the `Attached` state.  
> After enabling the setting, the DR volume will be activated and converted into a normal volume, remaining in the `Detached` state.

### Using Longhorn UI

1. Go to the `Volumes` page in the Longhorn UI of your secondary cluster.
2. Select the DR volume you want to activate.
3. Click the `Activate Disaster Recovery Volume` button in the **Operation** dropdown menu.
4. The volume will transition to the `Detached` state, and you can attach it with your workloads.

### Using `kubectl` Command

1. Run the following command to activate the DR volume and update the frontend:

    ```bash
    kubectl patch volume example-dr-volume1 -n longhorn-system --type='json' -p='[
      {"op": "replace", "path": "/spec/Standby", "value": false},
      {"op": "replace", "path": "/spec/frontend", "value": "blockdev"}
    ]'
    ```

2. The volume will transition to the `Detached` state, and you can attach it with your workloads.

## Limitations

Since the primary purpose of a DR volume is to restore data from backups, the following actions are not supported until the volume is activated:

- Creating, deleting, or reverting snapshots
- Creating backups
- Creating persistent volumes (PVs)
- Creating persistent volume claims (PVCs)


---

## Article: snapshots-and-backups/snapshot-space-management.md

---
title: Snapshot Space Management
weight: 1
---

Starting with v1.6.0, Longhorn allows you to configure the maximum snapshot count and the maximum aggregate snapshot size for each volume. Both settings do not take into account removed snapshots, backing images, and volume head snapshots. When either of these limits is reached, you must delete snapshots before creating new ones.

In earlier versions, the maximum snapshot count is not configurable (the value is 250) and there is no way to limit snapshot space usage.

## Settings

### Global Settings

**snapshot-max-count**: Maximum number of snapshots that you can create for each volume.

You must specify a value between "2" and "250". Longhorn requires at least two snapshots to function properly, particularly during volume rebuilding. One snapshot is created when the existing snapshots are merged, while the other snapshot is created during the rebuilding process.
The default value is "250".

When you create a volume without changing the default value of `.Spec.SnapshotMaxCount`, Longhorn applies the value of the `snapshot-max-count` setting. Changing the value of `snapshot-max-count` does not affect existing volumes.

### Volume-Specific Settings

**SnapshotMaxCount**: Maximum number of snapshots that you can create for a specific volume.

You can specify "0" or any value between "2" and "250". The default value is "0".

When you create a volume without changing the default value of this setting, Longhorn applies the value of the `snapshot-max-count` setting.

**SnapshotMaxSize**: Maximum aggregate size of snapshots for a specific volume.

You can specify "0" or any value larger than `Volume.Spec.Size` multiplied by 2. You must double the value of `Volume.Spec.Size` because Longhorn requires at least two snapshots to function properly.

The default value is "0", which effectively disables the setting.

When you expand the volume size, Longhorn automatically increases the value of this setting to `Volume.Spec.Size` multiplied by 2 (if the current value is smaller).


---

## Article: snapshots-and-backups/csi-snapshot-support/_index.md

---
title: CSI Snapshot Support
description: Creating and Restoring Longhorn Snapshots/Backups via the kubernetes CSI snapshot mechanism
weight: 3
---

## History
- [GitHub Issue](https://github.com/longhorn/longhorn/issues/304)
- [Longhorn Enhancement Proposal](https://github.com/longhorn/longhorn/blob/master/enhancements/20200904-csi-snapshot-support.md)

Available since v1.1.0


---

## Article: snapshots-and-backups/csi-snapshot-support/csi-volume-snapshot-associated-with-longhorn-backing-image.md

---
title: CSI VolumeSnapshot Associated with Longhorn BackingImage
weight: 2
---

BackingImage in Longhorn is an object that represents a QCOW2 or RAW image which can be set as the backing/base image of a Longhorn volume.

Instead of directly using Longhorn BackingImage resource for BackingImage management. You can also use the generic Kubernetes CSI VolumeSnapshot mechanism. To learn more about the CSI VolumeSnapshot mechanism, click [here](https://kubernetes.io/docs/concepts/storage/volume-snapshots/).

> **Prerequisite:** CSI snapshot support needs to be enabled on your cluster.
> If your kubernetes distribution does not provide the kubernetes snapshot controller
> as well as the snapshot related custom resource definitions, you need to manually deploy them.
> For more information, see [Enable CSI Snapshot Support](../enable-csi-snapshot-support).

## Create A CSI VolumeSnapshot Associated With Longhorn BackingImage

To create a CSI VolumeSnapshot associated with a Longhorn BackingImage, you first need to create a `VolumeSnapshotClass` object
with the parameter `type` set to `bi` as follow:
```yaml
kind: VolumeSnapshotClass
apiVersion: snapshot.storage.k8s.io/v1
metadata:
  name: longhorn-snapshot-vsc
driver: driver.longhorn.io
deletionPolicy: Delete
parameters:
  type: bi
  # export-type default to raw if it is not given
  export-type: qcow2
```
For more information about `VolumeSnapshotClass`, see the kubernetes documentation for [VolumeSnapshotClasses](https://kubernetes.io/docs/concepts/storage/volume-snapshot-classes/).

After that, create a Kubernetes `VolumeSnapshot` object with `volumeSnapshotClassName` points to the name of the `VolumeSnapshotClass` (`longhorn-snapshot-vsc`) and
the `source` points to the PVC of the Longhorn volume for which a Longhorn BackingImage should be exported from.
```yaml
apiVersion: snapshot.storage.k8s.io/v1
kind: VolumeSnapshot
metadata:
  name: test-csi-volume-snapshot-longhorn-backing-image
spec:
  volumeSnapshotClassName: longhorn-snapshot-vsc
  source:
    persistentVolumeClaimName: test-vol
```

**Result:**
A Longhorn BackingImage is created. The `VolumeSnapshot` object creation leads to the creation of a `VolumeSnapshotContent` Kubernetes object.
The `VolumeSnapshotContent` refers to a Longhorn BackingImage in its `VolumeSnapshotContent.snapshotHandle` field with the name `bi://backing?backingImageDataSourceType=export-from-volume&backingImage=${GENERATED_SNAPSHOT_NAME}&volume-name=test-vol&export-type=qcow2`.

### Viewing the Longhorn BackingImage

To see the BackingImage, click **Advanced > Backing Images** in the top navigation bar and click the BackingImage mentioned in the `VolumeSnapshotContent.snapshotHandle`.


### How the CSI Mechanism Works in this Scenario

When the VolumeSnapshot object is created with kubectl, the `VolumeSnapshot.uuid` field is used to identify a Longhorn BackingImage and the associated `VolumeSnapshotContent` object.

This creates a new Longhorn BackingImage named `snapshot-uuid` and the CSI request returns.

Afterwards a `VolumeSnapshotContent` object named `snapcontent-uuid` is created with the `VolumeSnapshotContent.readyToUse` flag is set to **true**.


## Restore PVC from CSI VolumeSnapshot Associated With Longhorn BackingImage
Create a `PersistentVolumeClaim` object where the `dataSource` field points to an existing `VolumeSnapshot` object that is associated with Longhorn BackingImage.

The csi-provisioner will pick this up and instruct the Longhorn CSI driver to provision a new volume using the associated Longhorn BackingImage.

An example `PersistentVolumeClaim` is below. The `dataSource` field needs to point to an existing `VolumeSnapshot` object.

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: test-restore-pvc
spec:
  storageClassName: longhorn
  dataSource:
    name: test-csi-volume-snapshot-longhorn-backing-image
    kind: VolumeSnapshot
    apiGroup: snapshot.storage.k8s.io
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 5Gi
```

### Restore a Longhorn BackingImage that Has No Associated `VolumeSnapshot` (pre-provision)

You can use the CSI mechanism to restore Longhorn BackingImage that has not been created via the CSI mechanism.
To restore Longhorn BackingImage that has not been created via the CSI mechanism, you have to first manually create a `VolumeSnapshot` and `VolumeSnapshotContent` object for the BackingImage.

Create a `VolumeSnapshotContent` object with the `snapshotHandle` field set to `bi://backing?backingImageDataSourceType=${TYPE}&backingImage=${BACKINGIMAGE_NAME}&backingImageChecksum=${backingImageChecksum}&${OTHER_PARAMETERS}` which point to an existing BackingImage.

- Users need to provide following query parameters in `snapshotHandle` for validation purpose:
    - `backingImageDataSourceType`: `sourceType` of existing BackingImage, e.g. `export-from-volume`, `download`
    - `backingImage`: Name of the BackingImage
    - `backingImageChecksum`: Optional. Checksum of the BackingImage.
    - you should also provide the `sourceParameters` of existing BackingImage in the `snapshotHandle` based on the `backingImageDataSourceType`
      - `export-from-volume`:
        - `volume-name`: volume to be expoted from.
        - `export-type`: qcow2 or raw.
      - `download`:
        - `url`: url of the BackingImage.
        - `checksum`: optional.

The parameters can be retrieved from the **Advanced > Backing Images** page in the Longhorn UI.

```yaml
apiVersion: snapshot.storage.k8s.io/v1
kind: VolumeSnapshotContent
metadata:
  name: test-existing-backing
spec:
  volumeSnapshotClassName: longhorn-snapshot-vsc
  driver: driver.longhorn.io
  deletionPolicy: Delete
  source:
    snapshotHandle: bi://backing?backingImageDataSourceType=download&backingImage=test-bi&url=https%3A%2F%2Flonghorn-backing-image.s3-us-west-1.amazonaws.com%2Fparrot.qcow2&backingImageChecksum=bd79ab9e6d45abf4f3f0adf552a868074dd235c4698ce7258d521160e0ad79ffe555b94e7d4007add6e1a25f4526885eb25c53ce38f7d344dd4925b9f2cb5d3b
  volumeSnapshotRef:
    name: test-snapshot-existing-backing
    namespace: default
```

Create the associated `VolumeSnapshot` object with the `name` field set to `test-snapshot-existing-backing`, where the `source` field refers to a `VolumeSnapshotContent` object via the `volumeSnapshotContentName` field.

This differs from the creation of a BackingImage, in which case the `source` field refers to a `PerstistentVolumeClaim` via the `persistentVolumeClaimName` field.

Only one type of reference can be set for a `VolumeSnapshot` object.

```yaml
apiVersion: snapshot.storage.k8s.io/v1beta1
kind: VolumeSnapshot
metadata:
  name: test-snapshot-existing-backing
spec:
  volumeSnapshotClassName: longhorn-snapshot-vsc
  source:
    volumeSnapshotContentName: test-existing-backing
```

Now you can create a `PerstistantVolumeClaim` object that refers to the newly created `VolumeSnapshot` object.
For an example see [Restore PVC from CSI VolumeSnapshot Associated With Longhorn BackingImage](#restore-pvc-from-csi-volumesnapshot-associated-with-longhorn-backingimage) above.


### Restore a Longhorn BackingImage that Has Not Created (on-demand provision)

You can use the CSI mechanism to restore Longhorn BackingImage which has not been created yet. This mechanism only support following 2 kinds of BackingImage data sources.

1. `download`: Download a file from a URL as a BackingImage.
2. `export-from-volume`: Export an existing in-cluster volume as a backing image.

Users need to create the `VolumeSnapshotContent` with an associated `VolumeSnapshot`. The `snapshotHandle` of the `VolumeSnapshotContent` needs to provide the parameters of the data source. Example below for a non-existing BackingImage `test-bi` with two different data sources.

1. `download`: Users need to provide following parameters
    - `backingImageDataSourceType`: `download` for on-demand download.
    - `backingImage`: Name of the BackingImage
    - `url`: Download the file from a URL as a BackingImage.
    - `backingImageChecksum`: Optional. Used for validating the file.
    - example yaml:
        ```yaml
        apiVersion: snapshot.storage.k8s.io/v1
        kind: VolumeSnapshotContent
        metadata:
            name: test-on-demand-backing
        spec:
            volumeSnapshotClassName: longhorn-snapshot-vsc
            driver: driver.longhorn.io
            deletionPolicy: Delete
            source:
              # NOTE: change this to provide the correct parameters
              snapshotHandle: bi://backing?backingImageDataSourceType=download&backingImage=test-bi&url=https%3A%2F%2Flonghorn-backing-image.s3-us-west-1.amazonaws.com%2Fparrot.qcow2&backingImageChecksum=bd79ab9e6d45abf4f3f0adf552a868074dd235c4698ce7258d521160e0ad79ffe555b94e7d4007add6e1a25f4526885eb25c53ce38f7d344dd4925b9f2cb5d3b
        volumeSnapshotRef:
            name: test-snapshot-on-demand-backing
            namespace: default
        ```

2. `export-from-volume`: Users need to provide following parameters
    - `backingImageDataSourceType`: `export-form-volume` for on-demand export.
    - `backingImage`: Name of the BackingImage
    - `volume-name`: Volume to be exported for the BackingImage
    - `export-type`: Currently Longhorn supports `raw` or `qcow2`
    - example yaml:
        ```yaml
        apiVersion: snapshot.storage.k8s.io/v1
        kind: VolumeSnapshotContent
        metadata:
        name: test-on-demand-backing
        spec:
        volumeSnapshotClassName: longhorn-snapshot-vsc
        driver: driver.longhorn.io
        deletionPolicy: Delete
        source: 
          # NOTE: change this to provide the correct parameters
          snapshotHandle: bi://backing?backingImageDataSourceType=export-from-volume&backingImage=test-bi&volume-name=vol-export-src&export-type=qcow2
        volumeSnapshotRef:
            name: test-snapshot-on-demand-backing
            namespace: default
        ```

Create the associated `VolumeSnapshot` object with the `name` field set to `test-snapshot-on-demand-backing`, where the `source` field refers to a `VolumeSnapshotContent` object via the `volumeSnapshotContentName` field.

This differs from the creation of a BackingImage, in which case the `source` field refers to a `PerstistentVolumeClaim` via the `persistentVolumeClaimName` field.

Only one type of reference can be set for a `VolumeSnapshot` object.

```yaml
apiVersion: snapshot.storage.k8s.io/v1beta1
kind: VolumeSnapshot
metadata:
  name: test-snapshot-on-demand-backing
spec:
  volumeSnapshotClassName: longhorn-snapshot-vsc
  source:
    volumeSnapshotContentName: test-on-demand-backing
```

Now you can create a `PerstistantVolumeClaim` object that refers to the newly created `VolumeSnapshot` object.
Longhorn will create the BackingImage with the parameters provide in the `snapshotHandle`.
For an example see [Restore PVC from CSI VolumeSnapshot Associated With Longhorn BackingImage](#restore-pvc-from-csi-volumesnapshot-associated-with-longhorn-backingimage) above.


---

## Article: snapshots-and-backups/csi-snapshot-support/csi-volume-snapshot-associated-with-longhorn-backup.md

---
title: CSI VolumeSnapshot Associated with Longhorn Backup
weight: 3
---

Backups in Longhorn are objects in an off-cluster backupstore, and the endpoint to access the backupstore is the backup target. For more information, see [this section.](../../../concepts/#31-how-backups-work)

To programmatically create backups, you can use the generic Kubernetes CSI VolumeSnapshot mechanism. To learn more about the CSI VolumeSnapshot mechanism, click [here](https://kubernetes.io/docs/concepts/storage/volume-snapshots/).

> **Prerequisite:** CSI snapshot support needs to be enabled on your cluster.
> If your kubernetes distribution does not provide the kubernetes snapshot controller
> as well as the snapshot related custom resource definitions, you need to manually deploy them.
> For more information, see [Enable CSI Snapshot Support](../enable-csi-snapshot-support).

## Create A CSI VolumeSnapshot Associated With Longhorn Backup

To create a CSI VolumeSnapshot associated with a Longhorn backup, you first need to create a `VolumeSnapshotClass` object
with the parameter `type` set to `bak` as follow:
```yaml
kind: VolumeSnapshotClass
apiVersion: snapshot.storage.k8s.io/v1
metadata:
  name: longhorn-backup-vsc
driver: driver.longhorn.io
deletionPolicy: Delete
parameters:
  type: bak
  backupMode: full
```

Longhorn performs incremental backups by default. Add the `backupMode: full` parameter to create a full backup. For more information about `backupMode`, see [Create A Backup](../../backup-and-restore/create-a-backup).

For more information about `VolumeSnapshotClass`, see the kubernetes documentation for [VolumeSnapshotClasses](https://kubernetes.io/docs/concepts/storage/volume-snapshot-classes/).

After that, create a Kubernetes `VolumeSnapshot` object with `volumeSnapshotClassName` points to the name of the `VolumeSnapshotClass` (`longhorn-backup-vsc`) and
the `source` points to the PVC of the Longhorn volume for which a backup should be created.
```yaml
apiVersion: snapshot.storage.k8s.io/v1
kind: VolumeSnapshot
metadata:
  name: test-csi-volume-snapshot-longhorn-backup
spec:
  volumeSnapshotClassName: longhorn-backup-vsc
  source:
    persistentVolumeClaimName: test-vol
```

**Result:**
A backup is created. The `VolumeSnapshot` object creation leads to the creation of a `VolumeSnapshotContent` Kubernetes object.
The `VolumeSnapshotContent` refers to a Longhorn backup in its `VolumeSnapshotContent.snapshotHandle` field with the name `bak://backup-volume/backup-name`.

### Viewing the Backup

To see the backup, click **Backup** in the top navigation bar and navigate to the backup-volume mentioned in the `VolumeSnapshotContent.snapshotHandle`.

For information on how to restore a volume via a `VolumeSnapshot` object, refer to the below sections.

### How the CSI Mechanism Works in this Scenario

When the VolumeSnapshot object is created with kubectl, the `VolumeSnapshot.uuid` field is used to identify a Longhorn snapshot and the associated `VolumeSnapshotContent` object.

This creates a new Longhorn snapshot named `snapshot-uuid`.

Then a backup of that snapshot is initiated, and the CSI request returns.

Afterwards a `VolumeSnapshotContent` object named `snapcontent-uuid` is created.

The CSI snapshotter sidecar periodically queries the Longhorn CSI plugin to evaluate the backup status.

Once the backup is completed, the `VolumeSnapshotContent.readyToUse` flag is set to **true**.


## Restore PVC from CSI VolumeSnapshot Associated With Longhorn Backup
Create a `PersistentVolumeClaim` object where the `dataSource` field points to an existing `VolumeSnapshot` object that is associated with Longhorn backup.

The csi-provisioner will pick this up and instruct the Longhorn CSI driver to provision a new volume with the data from the associated backup.

An example `PersistentVolumeClaim` is below. The `dataSource` field needs to point to an existing `VolumeSnapshot` object.

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: test-restore-pvc
spec:
  storageClassName: longhorn
  dataSource:
    name: test-csi-volume-snapshot-longhorn-backup
    kind: VolumeSnapshot
    apiGroup: snapshot.storage.k8s.io
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 5Gi
```
Note that the `spec.resources.requests.storage` value must be the same as the size of `VolumeSnapshot` object.


#### Restore a Longhorn Backup that Has No Associated `VolumeSnapshot`
You can use the CSI mechanism to restore Longhorn backups that have not been created via the CSI mechanism.
To restore Longhorn backups that have not been created via the CSI mechanism, you have to first manually create a `VolumeSnapshot` and `VolumeSnapshotContent` object for the backup.

Create a `VolumeSnapshotContent` object with the `snapshotHandle` field set to `bak://backup-volume/backup-name`.

The `backup-volume` and `backup-name` values can be retrieved from the **Backup** page in the Longhorn UI.

```yaml
apiVersion: snapshot.storage.k8s.io/v1
kind: VolumeSnapshotContent
metadata:
  name: test-existing-backup
spec:
  volumeSnapshotClassName: longhorn
  driver: driver.longhorn.io
  deletionPolicy: Delete
  source:
    # NOTE: change this to point to an existing backup on the backupstore
    snapshotHandle: bak://test-vol/backup-625159fb469e492e
  volumeSnapshotRef:
    name: test-snapshot-existing-backup
    namespace: default
```

Create the associated `VolumeSnapshot` object with the `name` field set to `test-snapshot-existing-backup`, where the `source` field refers to a `VolumeSnapshotContent` object via the `volumeSnapshotContentName` field.

This differs from the creation of a backup, in which case the `source` field refers to a `PerstistentVolumeClaim` via the `persistentVolumeClaimName` field.

Only one type of reference can be set for a `VolumeSnapshot` object.

```yaml
apiVersion: snapshot.storage.k8s.io/v1
kind: VolumeSnapshot
metadata:
  name: test-snapshot-existing-backup
spec:
  volumeSnapshotClassName: longhorn
  source:
    volumeSnapshotContentName: test-existing-backup
```

Now you can create a `PerstistantVolumeClaim` object that refers to the newly created `VolumeSnapshot` object.
For an example see [Restore PVC from CSI VolumeSnapshot Associated With Longhorn Backup](#restore-pvc-from-csi-volumesnapshot-associated-with-longhorn-backup) above.


---

## Article: snapshots-and-backups/csi-snapshot-support/csi-volume-snapshot-associated-with-longhorn-snapshot.md

---
title: CSI VolumeSnapshot Associated with Longhorn Snapshot
weight: 2
---

Snapshot in Longhorn is an object that represents content of a Longhorn volume at a particular moment. It is stored inside the cluster.

To programmatically create Longhorn snapshots, you can use the generic Kubernetes CSI VolumeSnapshot mechanism. To learn more about the CSI VolumeSnapshot mechanism, click [here](https://kubernetes.io/docs/concepts/storage/volume-snapshots/).

> **Prerequisite:** CSI snapshot support needs to be enabled on your cluster.
> If your kubernetes distribution does not provide the kubernetes snapshot controller
> as well as the snapshot related custom resource definitions, you need to manually deploy them.
> For more information, see [Enable CSI Snapshot Support](../enable-csi-snapshot-support).

## Create A CSI VolumeSnapshot Associated With Longhorn Snapshot

To create a CSI VolumeSnapshot associated with a Longhorn snapshot, you first need to create a `VolumeSnapshotClass` object
with the parameter `type` set to `snap` as follow:
```yaml
kind: VolumeSnapshotClass
apiVersion: snapshot.storage.k8s.io/v1
metadata:
  name: longhorn-snapshot-vsc
driver: driver.longhorn.io
deletionPolicy: Delete
parameters:
  type: snap
```
For more information about `VolumeSnapshotClass`, see the kubernetes documentation for [VolumeSnapshotClasses](https://kubernetes.io/docs/concepts/storage/volume-snapshot-classes/).

After that, create a Kubernetes `VolumeSnapshot` object with `volumeSnapshotClassName` points to the name of the `VolumeSnapshotClass` (`longhorn-snapshot-vsc`) and
the `source` points to the PVC of the Longhorn volume for which a Longhorn snapshot should be created.
```yaml
apiVersion: snapshot.storage.k8s.io/v1
kind: VolumeSnapshot
metadata:
  name: test-csi-volume-snapshot-longhorn-snapshot
spec:
  volumeSnapshotClassName: longhorn-snapshot-vsc
  source:
    persistentVolumeClaimName: test-vol
```

**Result:**
A Longhorn snapshot is created. The `VolumeSnapshot` object creation leads to the creation of a `VolumeSnapshotContent` Kubernetes object.
The `VolumeSnapshotContent` refers to a Longhorn snapshot in its `VolumeSnapshotContent.snapshotHandle` field with the name `snap://volume-name/snapshot-name`.

### Viewing the Longhorn Snapshot

To see the snapshot, click **Volume** in the top navigation bar and click the volume mentioned in the `VolumeSnapshotContent.snapshotHandle`. Scroll down to see the list of all volume snapshots.


### How the CSI Mechanism Works in this Scenario

When the VolumeSnapshot object is created with kubectl, the `VolumeSnapshot.uuid` field is used to identify a Longhorn snapshot and the associated `VolumeSnapshotContent` object.

This creates a new Longhorn snapshot named `snapshot-uuid` and the CSI request returns.

Afterwards a `VolumeSnapshotContent` object named `snapcontent-uuid` is created with the `VolumeSnapshotContent.readyToUse` flag is set to **true**.


## Restore PVC from CSI VolumeSnapshot Associated With Longhorn Snapshot
Create a `PersistentVolumeClaim` object where the `dataSource` field points to an existing `VolumeSnapshot` object that is associated with Longhorn snapshot.

The csi-provisioner will pick this up and instruct the Longhorn CSI driver to provision a new volume with the data from the associated Longhorn snapshot.

An example `PersistentVolumeClaim` is below. The `dataSource` field needs to point to an existing `VolumeSnapshot` object.

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: test-restore-pvc
spec:
  storageClassName: longhorn
  dataSource:
    name: test-csi-volume-snapshot-longhorn-snapshot
    kind: VolumeSnapshot
    apiGroup: snapshot.storage.k8s.io
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 5Gi
```

> **Note**:
> 1. The `spec.resources.requests.storage` value must match the size of the VolumeSnapshot object.
> 2. If you are restoring a volume from a `VolumeSnapshot` associated with a V2 data engine volume, you can use a `StorageClass` to specify the clone mode.
>    - Set the `cloneMode` parameter to either `full-copy` or `linked-clone`.
>    - If `cloneMode` is not specified, the default is `full-copy`.
>
>    For more information, see [the V2 Volume Clone Support documentation](../../../v2-data-engine/features/volume-clone).


---

## Article: snapshots-and-backups/csi-snapshot-support/enable-csi-snapshot-support.md

---
title: Enable CSI Snapshot Support on a Cluster
description: Enable CSI Snapshot Support for Programmatic Creation of Longhorn Snapshots/Backups
weight: 1
---

> **Prerequisite**
>
> It is the responsibility of the Kubernetes distribution to deploy the snapshot controller as well as the related custom resource definitions.
>
> For more information, see [CSI Volume Snapshots](https://kubernetes.io/docs/concepts/storage/volume-snapshots/).

#### If your Kubernetes Distribution Does Not Bundle the Snapshot Controller

You may manually install these components by executing the following steps.


> **Prerequisite**
>
> Please install the same release version of snapshot CRDs and snapshot controller to ensure that the CRD version is compatible with the snapshot controller.
>
> For general use, update the snapshot controller YAMLs with an appropriate **namespace** prior to installing.
>
> For example, on a vanilla Kubernetes cluster, update the namespace from `default` to `kube-system` prior to issuing the `kubectl create` command.

Install the Snapshot CRDs:
1. Download the files from https://github.com/kubernetes-csi/external-snapshotter/tree/v8.4.0/client/config/crd
because Longhorn v{{< current-version >}} uses [CSI external-snapshotter](https://kubernetes-csi.github.io/docs/external-snapshotter.html) v8.4.0
2. Run `kubectl create -k client/config/crd`.
3. Do this once per cluster.

Install the Common Snapshot Controller:
1. Download the files from https://github.com/kubernetes-csi/external-snapshotter/tree/v8.4.0/deploy/kubernetes/snapshot-controller
because Longhorn v{{< current-version >}} uses [CSI external-snapshotter](https://kubernetes-csi.github.io/docs/external-snapshotter.html) v8.4.0
2. Update the namespace to an appropriate value for your environment (e.g. `kube-system`)
3. Run `kubectl create -k deploy/kubernetes/snapshot-controller`.
3. Do this once per cluster.
> **Note:** previously, the snapshot controller YAML files were deployed into the `default` namespace by default.
> The updated YAML files are being deployed into `kube-system` namespace by default.
> Therefore, we suggest deleting the previous snapshot controller in the `default` namespace to avoid having multiple snapshot controllers.

See the [Usage](https://github.com/kubernetes-csi/external-snapshotter#usage) section from the kubernetes
external-snapshotter git repo for additional information.

#### Add a Default `VolumeSnapshotClass`
Ensure the availability of the Snapshot CRDs. Afterwards create a default `VolumeSnapshotClass`.
```yaml
# Use v1 as an example
kind: VolumeSnapshotClass
apiVersion: snapshot.storage.k8s.io/v1
metadata:
  name: longhorn
driver: driver.longhorn.io
deletionPolicy: Delete
```


---

## Article: snapshots-and-backups/backup-and-restore/_index.md

---
title: Backup and Restore
weight: 2
---

> Before v1.2.0, Longhorn used a "blocking way" for communication with the remote backup target. Consequently, there are some involuntary factors impacting the functions relying on remote backup target. For example, network latency, listing backups, or causing further cascading problems after the backup target operation.

> Since v1.2.0, Longhorn started using an asynchronous backup operations to resolve the aforementioned issues in the previous version.
> - To do this, create the backup cluster custom resources first, then perform the following snapshot and backup operations to the remote backup target.
> - Once the backup creation is completed, asynchronously pull the state of backup volumes and backups from the remote backup target. Then, update the status of the corresponding cluster custom resources.
>
> This enhancement is scalable for the backup query to assist with resolving the costly resources caused by the blocking way. This was because all backups are saved as custom resources instead of querying from the remote target directly.
>
> Note: After the Longhorn upgrade, if a volume has not been upgraded to the latest Longhorn engine (>=v1.2.0). When creating a backup, it will have the intermediate transition state of the name of the created backup (due to the different backup name handling in the latest longhorn version >= v1.2.0). Longhorn will then ensure the backup is synced with the remote backup target and the backup will be updated to the final correct state with the remote backup target is the single source of truth. To upgrade the Longhorn engine, refer to [Manually Upgrade Longhorn Engine](../../deploy/upgrade/upgrade-engine) or [Automatically Upgrade Longhorn Engine](../../deploy/upgrade/auto-upgrade-engine).

- [Setting a Backup Target](./set-backup-target)
- [Create a Backup](./create-a-backup)
- [Restore from a Backup](./restore-from-a-backup)
- [Restoring Volumes for Kubernetes StatefulSets](./restore-statefulset)


---

## Article: snapshots-and-backups/backup-and-restore/configure-backup-block-size.md

---
title: Configure The Block Size Of Backup
weight: 2
---

## Backup Blocks In Longhorn

A Longhorn backup is composed of data fragments derived from a snapshot, where each fragment is called a block. Blocks are the fundamental units used for processing, transmission, and storage in the backup target. All blocks within a single backup have the same physical size.

Prior to Longhorn v1.10.0, the backup block size was fixed at **2 MiB**. Starting in Longhorn v1.10.0, users can configure the backup block size during **volume creation**. This value is immutable once the volume is created. The block size used for backups is displayed on the volume detail page in the Longhorn UI, and all backups for a volume will use the size defined at creation.

### Impact of Backup Block Size

Longhorn supports two available backup block size, 2 MiB and 16 MiB. The selected block size affects the efficiency of backup creation and storage:

1. Larger block sizes reduce the total number of blocks, improving transmission efficiency and reducing the number of API requests to the backup target.
2. However, larger block sizes can increase the physical storage footprint due to zero-padding and require more memory during backup creation.

## Global Default Backup Block Size

A global setting allows users to define the default backup block size for new volumes. If a backup block size is not explicitly set during volume creation, Longhorn will apply the default value. To change the default backup block size:

- Using Longhorn UI:
    ```
    Settings > General > Default Backup Block Size
    ```
- Using `kubectl`:
    ```bash
    kubectl -n longhorn-system edit settings.longhorn.io default-backup-block-size
    ```

## Create a Volume And Specify The Backup Block Size

To specify a custom backup block size during volume creation:

1. Navigate to the **Volume** menu.
2. Click **Create Volume**.
3. In the volume creation dialog, inside `Advanced Configurations`, select the desired **Backup Block Size**.

## Specify The Backup Block Size In The Storage Class

For volumes provisioned through a Persistent Volume Claim (PVC), you can set the `backupBlockSize` in the `parameters` section of the `StorageClass`.

Example:

```yaml
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: longhorn-example
provisioner: driver.longhorn.io
parameters:
  backupBlockSize: 16Mi
...
```

## Restoring Volume From a Backup

When restoring a volume from a backup, the restored volume can be configured with a different backup block size than the original.

**Caution**: Longhorn versions prior to v1.10 lack forward compatibility and cannot restore backups created by v1.10 or later. Restoring a backup with a non-default backup block size (anything other than 2 MiB) on Longhorn v1.9.x or older will result in a volume being created with file system corruption.




---

## Article: snapshots-and-backups/backup-and-restore/create-a-backup.md

---
title: Create a Backup
weight: 2
---

## Incremental Backup
Backups in Longhorn are objects in an off-cluster backupstore. A backup of a snapshot is copied to the backupstore, and the endpoint to access the backupstore is the backup target. For more information, see [this section.](../../../concepts/#31-how-backups-work)

> **Prerequisite:** A backup target must be set up. For more information, see [Set the BackupTarget](../set-backup-target). If the BackupTarget has not been set, you'll be presented with an error.

### Create an Incremental Backup Using UI
To create a backup,

1. Navigate to the **Volume** menu.
2. Select the volume you wish to back up.
3. Click **Create Backup.**
4. Add any appropriate labels and click OK.

**Result:** The backup is created. To see it, click **Backup** in the top navigation bar.

For information about restoring a volume from a snapshot, see [Restore from a Backup](../restore-from-a-backup).

### Create an Incremental Using YAML Code

1. Obtain the name of the snapshot that you want to back up (from either the Longhorn UI or the CR).
2. Apply the YAML.

Example:

```yaml
apiVersion: longhorn.io/v1beta2
kind: Backup
metadata:
  name: backup-example
  namespace: longhorn-system
spec:
  backupMode: incremental
  snapshotName: snapshot-name-example
  labels:
    app: test
```

## Full Backup

By default, Longhorn backs up only data that was changed since the last backup. This approach, known as *delta backup*, enhances time efficiency and conserves network throughput. However, when a data block in the backupstore becomes corrupted, Longhorn does not replace that data block with a healthy one during subsequent backup operations.

Starting with v1.7.0, Longhorn can perform full backups that upload all data blocks in the volume and overwrite existing data blocks in the backupstore.

### Create a Full Backup Using the Longhorn UI
1. Go to the **Volume** page.
2. Select the volume that you want to back up.
3. Click **Create Backup**.
4. Add appropriate labels.
5. Select Full Backup.
6. Click **OK**.

### Create a Full Backup Using YAML Code
1. Obtain the name of the snapshot that you want to back up (from either the Longhorn UI or the CR).
2. Apply the YAML.

Example:

```yaml
apiVersion: longhorn.io/v1beta2
kind: Backup
metadata:
  name: backup-example
  namespace: longhorn-system
spec:
  backupMode: full
  snapshotName: snapshot-name-example
  labels:
    app: test
```

## Uploaded Data Size

To facilitate collection of data transfer information for each backup, Longhorn records the information using two metrics in the CR status.

### Newly Uploaded Data Size
`status.newlyUploadDataSize` records the size of data that was uploaded *for the first time* to the backupstore during the latest backup. In other words, it tracks the size of data blocks that did not previously exist in the backupstore.

### Re-Uploaded Data Size
`status.reUploadDataSize` records the size of data that was overwritten during the latest full backup. In other words, it tracks the size of data blocks that previously existed in the backupstore.

---

## Article: snapshots-and-backups/backup-and-restore/restore-from-a-backup.md

---
title: Restore from a Backup
weight: 3
---

Longhorn can easily restore backups to a volume. 

For more information on how backups work, refer to the [concepts](../../../concepts/#3-backups-and-secondary-storage) section.

When you restore a backup, it creates a volume of the same name by default. If a volume with the same name as the backup already exists, the backup will not be restored.

To restore a backup,

1. Navigate to the **Backup.** menu
2. Select the backup(s) you wish to restore and click **Restore Latest Backup.**
3. In the **Name** field, select the volume you wish to restore.
4. Click **OK.**

You can then create the PV/PVC from the volume after restoring a volume from a backup. Here you can specify the `storageClassName` or leave it empty to use the `storageClassName` inherited from the PVC of the backup volume. The `StorageClass` should be already in the cluster to prevent any further issue.

**Result:** The restored volume is available on the **Volume** page.

---

## Article: snapshots-and-backups/backup-and-restore/restore-recurring-jobs-from-a-backup.md

---
title: Restore Volume Recurring Jobs from a Backup
weight: 5
---

Since v1.4.0, Longhorn supports recurring jobs backup and restore along with the volume backup and restore. When restoring a backup volume, if users enable the `Restore Volume Recurring Jobs` setting, the original recurring jobs of the volume will be restored back accordingly.

For more information on the setting `Restore Volume Recurring Jobs`, refer to the [settings](../../../references/settings/#restore-volume-recurring-jobs) section.

For more information on how volume backup works, refer to the [concepts](../../../concepts/#3-backups-and-secondary-storage) section.

When restoring a volume with recurring jobs, Longhorn will restore them together. If the volume name already exists, the volume and the recurring jobs will not be restored.  If the recurring job name already exists but the spec is different, the restoring recurring job will be created with a randomly generated name to avoid conflict. Otherwise, Longhorn will try to reuse existing recurring jobs instead if they are the same as restoring recurring jobs of a backup volume.

By default, Longhorn will not automatically restore volume recurring jobs, users can enable the automatic restoration by Longhorn UI or kubectl.

## Via Longhorn UI

1. Navigate to the **Setting** menu and click **General**
2. Enable the `Restore Volume Recurring Jobs`
3. Navigate to the **Backup** menu
4. Select the backup(s) you wish to restore and click **Restore Latest Backup.**
5. In the **Name** field, select the volume you wish to restore.
6. Click **OK**

## Via Command Line

```bash
# kubectl -n longhorn-system edit settings.longhorn.io restore-volume-recurring-jobs
```

Then, set the value to `true`.

```text
# kubectl -n longhorn-system get setting restore-volume-recurring-jobs
NAME                            VALUE   AGE
restore-volume-recurring-jobs   false   28m
```

### Example of Volume Specific Setting

```yaml
apiVersion: longhorn.io/v1beta2
kind: Volume
metadata:
  labels:
    longhornvolume: vol-01
  name: vol-01
  namespace: longhorn-system
spec:
  restoreVolumeRecurringJob: ignored
  engineImage: longhornio/longhorn-engine:v1.4.0
  fromBackup: "s3://backupbucket@us-east-1?volume=minio-vol01&backup=backup-eeb2782d5b2f42bb"
  frontend: blockdev
```

Users can override the setting `restore-volume-recurring-jobs` by the volume spec property `spec.restoreVolumeRecurringJob`.

- **ignored**. This is the default option that instructs Longhorn to inherit from the global setting.
- **enabled**. This option instructs Longhorn to restore volume recurring jobs from the backup target forcibly.
- **disabled**. This option instructs Longhorn no restoring volume recurring jobs should be done.

**Result:** The restored volume recurring jobs are available on the **RecurringJob** page.


---

## Article: snapshots-and-backups/backup-and-restore/restore-statefulset.md

---
title: Restoring Volumes for Kubernetes StatefulSets
weight: 4
---
Longhorn supports restoring backups, and one of the use cases for this feature is to restore data for use in a Kubernetes StatefulSet, which requires restoring a volume for each replica that was backed up.

To restore, follow the below instructions. The example below uses a StatefulSet with one volume attached to each Pod and two replicas.

1. Connect to the `Longhorn UI` page in your web browser. Under the `Backup` tab, select the name of the StatefulSet volume. Click the dropdown menu of the volume entry and restore it. Name the volume something that can easily be referenced later for the `Persistent Volumes`.
    - Repeat this step for each volume you need restored.
    - For example, if restoring a StatefulSet with two replicas that had volumes named `pvc-01a` and `pvc-02b`, the restore could look like this:

    | Backup Name | Restored Volume   |
    |-------------|-------------------|
    | pvc-01a     | statefulset-vol-0 |
    | pvc-02b     | statefulset-vol-1 |

2. In Kubernetes, create a `Persistent Volume` for each Longhorn volume that was created. Name the volumes something that can easily be referenced later for the `Persistent Volume Claims`. `storage` capacity, `numberOfReplicas`, `storageClassName`, and `volumeHandle` must be replaced below. In the example, we're referencing `statefulset-vol-0` and `statefulset-vol-1` in Longhorn and using `longhorn` as our `storageClassName`.

    ```
    apiVersion: v1
    kind: PersistentVolume
    metadata:
      name: statefulset-vol-0
    spec:
      capacity:
        storage: <size> # must match size of Longhorn volume
      volumeMode: Filesystem
      accessModes:
        - ReadWriteOnce
      persistentVolumeReclaimPolicy: Delete
      csi:
        driver: driver.longhorn.io # driver must match this
        fsType: ext4
        volumeAttributes:
          numberOfReplicas: <replicas> # must match Longhorn volume value
          staleReplicaTimeout: '30' # in minutes
        volumeHandle: statefulset-vol-0 # must match volume name from Longhorn
      storageClassName: longhorn # must be same name that we will use later
    ---
    apiVersion: v1
    kind: PersistentVolume
    metadata:
      name: statefulset-vol-1
    spec:
      capacity:
        storage: <size>  # must match size of Longhorn volume
      volumeMode: Filesystem
      accessModes:
        - ReadWriteOnce
      persistentVolumeReclaimPolicy: Delete
      csi:
        driver: driver.longhorn.io # driver must match this
        fsType: ext4
        volumeAttributes:
          numberOfReplicas: <replicas> # must match Longhorn volume value
          staleReplicaTimeout: '30'
        volumeHandle: statefulset-vol-1 # must match volume name from Longhorn
      storageClassName: longhorn # must be same name that we will use later
    ```

> 2.1 In the case of encrypted volume, make sure you are specifying the `nodePublishSecretRef`, and `nodeStageSecretRef` while creating the `PV`.
>
> ```yaml
> kind: PersistentVolume
> metadata:
>   name: statefulset-encrypted-vol-0
> spec:
>   capacity:
>     storage: <size>
>   volumeMode: Filesystem
>   accessModes:
>     - ReadWriteOnce
>   persistentVolumeReclaimPolicy: Delete
>   csi:
>     driver: driver.longhorn.io
>     fsType: ext4
>     nodePublishSecretRef:
>       name: <secret-name>
>       namespace: <namespace>
>     nodeStageSecretRef:
>       name: <secret-name>
>       namespace: <namespace>
>     volumeAttributes:
>       numberOfReplicas: <replicas>
>       staleReplicaTimeout: "30"
>     volumeHandle: statefulset-encrypted-vol-0
>   storageClassName: longhorn
> ```

3. In the `namespace` the `StatefulSet` will be deployed in, create PersistentVolume Claims **for each** `Persistent Volume`. The name of the `Persistent Volume Claim` must follow this naming scheme:

    ```
    <name of Volume Claim Template>-<name of StatefulSet>-<index>
    ```
    StatefulSet Pods are zero-indexed. In this example, the name of the `Volume Claim
  Template` is `data`, the name of the `StatefulSet` is `webapp`, and there
  are two replicas, which are indexes `0` and `1`.

    ```
    apiVersion: v1
    kind: PersistentVolumeClaim
    metadata:
      name: data-webapp-0
    spec:
      accessModes:
      - ReadWriteOnce
      resources:
        requests:
          storage: 2Gi # must match size from earlier
      storageClassName: longhorn # must match name from earlier
      volumeName: statefulset-vol-0 # must reference Persistent Volume
    ---
    apiVersion: v1
    kind: PersistentVolumeClaim
    metadata:
      name: data-webapp-1
    spec:
      accessModes:
      - ReadWriteOnce
      resources:
        requests:
          storage: 2Gi # must match size from earlier
      storageClassName: longhorn # must match name from earlier
      volumeName: statefulset-vol-1 # must reference Persistent Volume
    ```

4. Create the `StatefulSet`:

    ```
    apiVersion: apps/v1beta2
    kind: StatefulSet
    metadata:
      name: webapp # match this with the PersistentVolumeClaim naming scheme
    spec:
      selector:
        matchLabels:
          app: nginx # has to match .spec.template.metadata.labels
      serviceName: "nginx"
      replicas: 2 # by default is 1
      template:
        metadata:
          labels:
            app: nginx # has to match .spec.selector.matchLabels
        spec:
          terminationGracePeriodSeconds: 10
          containers:
          - name: nginx
            image: registry.k8s.io/nginx-slim:0.8
            ports:
            - containerPort: 80
              name: web
            volumeMounts:
            - name: data
              mountPath: /usr/share/nginx/html
      volumeClaimTemplates:
      - metadata:
          name: data # match this with the PersistentVolumeClaim naming scheme
        spec:
          accessModes: [ "ReadWriteOnce" ]
          storageClassName: longhorn # must match name from earlier
          resources:
            requests:
              storage: 2Gi # must match size from earlier
    ```

**Result:** The restored data should now be accessible from inside the `StatefulSet`
`Pods`.


---

## Article: snapshots-and-backups/backup-and-restore/set-backup-target.md

---
title: Setting a Backup Target
weight: 1
---

A backup target is an endpoint used to access a backupstore. Backup targets can be configured on the Longhorn UI (**Settings > Backup Target**). A backupstore is a server that stores the backups of Longhorn volumes. You can use NFS, SMB/CIFS, Azure Blob Storage, and S3-compatible servers.

{{< figure alt="the backup target UI page" src="/img/screenshots/backup-target/v1.10.0/page.png" >}}

> **Note:**  
> Starting with v1.8.0, Longhorn supports usage of multiple backupstores. Setting the default backup target before creating a new one is recommended.

Saving to an object store such as S3 is preferable because it generally offers better reliability.  Another advantage is that you do not need to mount and unmount the target, which can complicate failover and upgrades.

For more information about how the backupstore works in Longhorn, see the [concepts section.](../../../concepts/#3-backups-and-secondary-storage)

If you don't have access to AWS S3 or want to give the backupstore a try first, we've also provided a way to [setup a local S3 testing backupstore](#set-up-a-local-testing-backupstore) using [MinIO](https://minio.io/).

Longhorn also supports setting up recurring snapshot/backup jobs for volumes, via Longhorn UI or Kubernetes Storage Class. See [here](../../scheduling-backups-and-snapshots) for details.

> **Notice**
>
> - The lifecycle of Longhorn backups within the backupstore is entirely managed by Longhorn. **Any retention policy directly on the backupstore is strictly prohibited**.
>
> - Longhorn attempts to clean up the backup-related custom resources in the following scenarios:
>   - An empty response from the NFS server due to server downtime.
>   - A race condition between related Longhorn backup controllers.
>
>   The backup information is resynchronized during the next polling interval. For more information, see [#9530](https://github.com/longhorn/longhorn/issues/9530).

This page covers the following topics:

- [Default Backup Target](#default-backup-target)
  - [Set the Default Backup Target Using Helm](#set-the-default-backup-target-using-helm)
  - [Set the Default Backup Target Using a Manifest YAML File](#set-the-default-backup-target-using-a-manifest-yaml-file)
- [Set up AWS S3 Backupstore](#set-up-aws-s3-backupstore)
- [Set up GCP Cloud Storage Backupstore](#set-up-gcp-cloud-storage-backupstore)
- [Set up a Local Testing Backupstore](#set-up-a-local-testing-backupstore)
- [Using a self-signed SSL certificate for S3 communication](#using-a-self-signed-ssl-certificate-for-s3-communication)
- [Enable virtual-hosted-style access for S3 compatible Backupstore](#enable-virtual-hosted-style-access-for-s3-compatible-backupstore)
- [Set up NFS Backupstore](#set-up-nfs-backupstore)
- [Set up SMB/CIFS Backupstore](#set-up-smbcifs-backupstore)
- [Set up Azure Blob Storage Backupstore](#set-up-azure-blob-storage-backupstore)

### Default Backup Target

The default backup target (`default`) is automatically created during a fresh installation. You can set the default backup target during or after the installation using either Helm or a [manifest YAML file](https://raw.githubusercontent.com/longhorn/longhorn/v1.8.0/deploy/longhorn.yaml)(`longhorn.yaml`).

#### Set the Default Backup Target Using Helm

In the `values.yaml` file, you can set three parameters to manage the default backup target.

- `defaultBackupStore.backupTarget`: Endpoint used to access the default backupstore.
- `defaultBackupStore.backupTargetCredentialSecret`: Name of the Kubernetes secret associated with the default backup target.
- `defaultBackupStore.pollInterval`: Number of seconds that Longhorn waits before checking the default backupstore for new backups.

```yaml
# -- Setting that allows you to update the default backupstore.
defaultBackupStore:
  # -- Endpoint used to access the default backupstore.
  backupTarget: ~
  # -- Name of the Kubernetes secret associated with the default backup target.
  backupTargetCredentialSecret: ~
  # -- Number of seconds that Longhorn waits before checking the default backupstore for new backups.
  pollInterval: ~
```

#### Set the Default Backup Target Using a Manifest YAML File

Starting with v1.8.0, you can use a new `ConfigMap` resource named `longhorn-default-resource` to manage settings of resources, including the default backup target resource.

- `backup-target`: Endpoint used to access the default backupstore.
- `backup-target-credential-secret`: Name of the Kubernetes secret associated with the default backup target.
- `backupstore-poll-interval`: Number of seconds that Longhorn waits before checking the default backupstore for new backups.

```yaml
# Example
apiVersion: v1
kind: ConfigMap
metadata:
  name: longhorn-default-resource
  namespace: longhorn-system
data:
  default-resource.yaml: |
    "backup-target": "s3://example@us-west-1/"
    "backup-target-credential-secret": "example-secret"
    "backupstore-poll-interval": "180"
```

### Set up AWS S3 Backupstore

1. Create a new bucket in [AWS S3.](https://aws.amazon.com/s3/)

2. Set permissions for Longhorn. There are two options for setting up the credentials. The first is that you can set up a Kubernetes secret with the credentials of an AWS IAM user. The second is that you can use a third-party application to manage temporary AWS IAM permissions for a Pod via annotations rather than operating with AWS credentials.

   - Option 1: Create a Kubernetes secret with IAM user credentials

     1. Follow the [guide](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_users_create.html#id_users_create_console) to create a new AWS IAM user, with the following permissions set. Edit the `Resource` section to use your S3 bucket name:

        ```json
        {
          "Version": "2012-10-17",
          "Statement": [
            {
              "Sid": "GrantLonghornBackupstoreAccess0",
              "Effect": "Allow",
              "Action": [
                "s3:PutObject",
                "s3:GetObject",
                "s3:ListBucket",
                "s3:DeleteObject"
              ],
              "Resource": [
                "arn:aws:s3:::<your-bucket-name>",
                "arn:aws:s3:::<your-bucket-name>/*"
              ]
            }
          ]
        }
        ```

     2. Create a Kubernetes secret with a name such as `aws-secret` in the namespace where Longhorn is placed (`longhorn-system` by default). The secret must be created in the `longhorn-system` namespace for Longhorn to access it:

        ```shell
        kubectl create secret generic <aws-secret> \
            --from-literal=AWS_ACCESS_KEY_ID=<your-aws-access-key-id> \
            --from-literal=AWS_SECRET_ACCESS_KEY=<your-aws-secret-access-key> \
            -n longhorn-system
        ```

   - Option 2: Set permissions with IAM temporary credentials by AWS STS AssumeRole (kube2iam or kiam)

     [kube2iam](https://github.com/jtblin/kube2iam) or [kiam](https://github.com/uswitch/kiam) is a Kubernetes application that allows managing AWS IAM permissions for Pod via annotations rather than operating on AWS credentials. Follow the instructions in the GitHub repository for kube2iam or kiam to install it into the Kubernetes cluster.

     1. Follow the [guide](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_create_for-service.html#roles-creatingrole-service-console) to create a new AWS IAM role for AWS S3 service, with the following permissions set:

         ```json
        {
           "Version": "2012-10-17",
          "Statement": [
             {
               "Sid": "GrantLonghornBackupstoreAccess0",
               "Effect": "Allow",
               "Action": [
                 "s3:PutObject",
                 "s3:GetObject",
                 "s3:ListBucket",
                 "s3:DeleteObject"
               ],
               "Resource": [
                 "arn:aws:s3:::<your-bucket-name>",
                 "arn:aws:s3:::<your-bucket-name>/*"
               ]
             }
           ]
         }
        ```

     2. Edit the AWS IAM role with the following trust relationship:

        ```json
        {
          "Version": "2012-10-17",
          "Statement": [
            {
              "Effect": "Allow",
              "Principal": {
                  "Service": "ec2.amazonaws.com"
              },
              "Action": "sts:AssumeRole"
            },
            {
              "Effect": "Allow",
              "Principal": {
                "AWS": "arn:aws:iam::<AWS_ACCOUNT_ID>:role/<AWS_EC2_NODE_INSTANCE_ROLE>"
              },
              "Action": "sts:AssumeRole"
            }
          ]
        }
        ```

     3. Create a Kubernetes secret with a name such as `aws-secret` in the namespace where Longhorn is placed (`longhorn-system` by default). The secret must be created in the `longhorn-system` namespace for Longhorn to access it:

        ```shell
        kubectl create secret generic <aws-secret> \
            --from-literal=AWS_IAM_ROLE_ARN=<your-aws-iam-role-arn> \
            -n longhorn-system
        ```

3. On the Longhorn UI, go to **Backup and Restore > Backup Targets**, and then create or edit a backup target.

   {{< figure alt="edit a backup target" src="/img/screenshots/backup-target/v1.10.0/edit.png" >}}

   - Set **URL** to:

     ```text
     s3://<your-bucket-name>@<your-aws-region>/
     ```

     Make sure that you have `/` at the end, otherwise you will get an error. A subdirectory (prefix) may be used:

     ```text
     s3://<your-bucket-name>@<your-aws-region>/mypath/
     ```

     Also make sure you've set **`<your-aws-region>` in the URL**.

     For example, For AWS, you can find the region codes [here.](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Concepts.RegionsAndAvailabilityZones.html)
     For Google Cloud Storage, you can find the region codes [here.](https://cloud.google.com/storage/docs/locations)

   - Set **Credential Secret** to:

     ```text
     aws-secret
     ```

     This is the secret name with AWS credentials or AWS IAM role.

**Result:** Longhorn can store backups in S3. To create a backup, see [this section.](../create-a-backup)

**Note:** If you operate Longhorn behind a proxy and you want to use AWS S3 as the backupstore, you must provide Longhorn information about your proxy in the `aws-secret` as below:

```shell
kubectl create secret generic <aws-secret> \
    --from-literal=AWS_ACCESS_KEY_ID=<your-aws-access-key-id> \
    --from-literal=AWS_SECRET_ACCESS_KEY=<your-aws-secret-access-key> \
    --from-literal=HTTP_PROXY=<your-proxy-ip-and-port> \
    --from-literal=HTTPS_PROXY=<your-proxy-ip-and-port> \
    --from-literal=NO_PROXY=<excluded-ip-list> \
    -n longhorn-system
```

Make sure `NO_PROXY` contains the network addresses, network address ranges and domains that should be excluded from using the proxy. In order for Longhorn to operate, the minimum required values for `NO_PROXY` are:

* localhost
* 127.0.0.1
* 0.0.0.0
* 10.0.0.0/8 (K8s components' IPs)
* 192.168.0.0/16 (internal IPs in the cluster)

### Set up GCP Cloud Storage Backupstore

1. Create a new bucket in [Google Cloud Storage](https://console.cloud.google.com/storage/browser?referrer=search&project=elite-protocol-319303)
2. Create a GCP serviceaccount in [IAM & Admin](https://console.cloud.google.com/iam-admin)
3. Give the GCP serviceaccount permissions to read, write, and delete objects in the bucket.

   The serviceaccount will require the `roles/storage.objectAdmin` role to read, write, and delete objects in the bucket.

   Here is a reference to the GCP IAM roles you have available for granting access to a serviceaccount https://cloud.google.com/storage/docs/access-control/iam-roles.

   > **Note:** Consider creating an IAM condition to reduce how many buckets this serviceaccount has object admin access to. On the Google Cloud console, go to **Cloud Storage > Buckets**, and select the target bucket. On the **Bucket details** screen, go to the **Permissions** tab, click **Grant Access**, and grant your service account Storage Object Admin permissions for the target bucket.

4. Navigate to your [buckets in cloud storage](https://console.cloud.google.com/storage/browser) and select your newly created bucket.
5. Go to the cloud storage's settings menu and navigate to the [interoperability tab](https://console.cloud.google.com/storage/settings;tab=interoperability)
6. Scroll down to _Service account HMAC_ and press `+ CREATE A KEY FOR A SERVICE ACCOUNT`
7. Select the GCP serviceaccount you created earlier and press `CREATE KEY`
8. Save the _Access Key_ and _Secret_.

    Also note down the configured _Storage URI_ under the _Request Endpoint_ while you're in the interoperability menu.

   - The Access Key will be mapped to the `AWS_ACCESS_KEY_ID` field in the Kubernetes secret we create later.
   - The Secret will be mapped to the `AWS_SECRET_ACCESS_KEY` field in the Kubernetes secret we create later.
   - The Storage URI will be mapped to the `AWS_ENDPOINTS` field in the Kubernetes secret we create later.

9. Go to the Longhorn UI. In the top navigation bar, click **Backup and Restore/Backup Targets**, and create or edit a backup target.

   - Set **URL** to:

     ```text
     s3://${BUCKET_NAME}@us/
     ```

   - Set **Credential Secret** to:

     ```text
     longhorn-gcp-backups
     ```

10. Create a Kubernetes secret named `longhorn-gcp-backups` in the `longhorn-system` namespace with the following content:

    ```yaml
    apiVersion: v1
    kind: Secret
    metadata:
      name: longhorn-gcp-backups
      namespace: longhorn-system
    type: Opaque
    stringData:
      AWS_ACCESS_KEY_ID: GOOG1EBYHGDE4WIGH2RDYNZWWWDZ5GMQDRMNSAOTVHRAILWAMIZ2O4URPGOOQ
      AWS_ENDPOINTS: https://storage.googleapis.com
      AWS_SECRET_ACCESS_KEY: BKoKpIW021s7vPtraGxDOmsJbkV/0xOVBG73m+8f
    ```

    > **Note:** The secret can be named whatever you like as long as they match what's in longhorn's settings.

Once the secret is created and Longhorn's settings are saved, navigate to the backup tab in Longhorn. If there are any issues, they should pop up as a toast notification.

If you don't get any error messages, try creating a backup and confirm the content is pushed out to your new bucket.

The **Backup Target** screen on the Longhorn UI displays the status of each backup target. If the status is **Error** and no other details are provided, you can use the **Inspect** feature of your browser to view the response data for `/v1/backuptargets`. Errors from GCP are labeled "AWS Error" (for example, "AWS Error: AccessDenied"). For more information, see [Issue #10428](https://github.com/longhorn/longhorn/issues/10428).

### Set up a Local Testing Backupstore

Longhorn provides sample backupstore server setups for testing purposes.  You can find samples for AWS S3 (MinIO), Azure, CIFS and NFS in the `longhorn/deploy/backupstores` folder.

1. Set up a MinIO S3 server for the backupstore in the `longhorn-system` namespace.

   ```shell
   kubectl create -f https://raw.githubusercontent.com/longhorn/longhorn/v{{< current-version >}}/deploy/backupstores/minio-backupstore.yaml
   ```

2. Go to the Longhorn UI. click **Backup and Restore/Backup Targets**, and create or edit a backup target.

   - Set **URL** to:

     ```text
     s3://backupbucket@us-east-1/
     ```

   - Set **Credential Secret** to:

     ```text
     minio-secret
     ```

     The `minio-secret` yaml looks like this:

     ```yaml
     apiVersion: v1
     kind: Secret
     metadata:
       name: minio-secret
       namespace: longhorn-system
     type: Opaque
     data:
       AWS_ACCESS_KEY_ID: bG9uZ2hvcm4tdGVzdC1hY2Nlc3Mta2V5 # longhorn-test-access-key
       AWS_SECRET_ACCESS_KEY: bG9uZ2hvcm4tdGVzdC1zZWNyZXQta2V5 # longhorn-test-secret-key
       AWS_ENDPOINTS: aHR0cHM6Ly9taW5pby1zZXJ2aWNlLmRlZmF1bHQ6OTAwMA== # https://minio-service.default:9000
       AWS_CERT: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURMRENDQWhTZ0F3SUJBZ0lSQU1kbzQycGhUZXlrMTcvYkxyWjVZRHN3RFFZSktvWklodmNOQVFFTEJRQXcKR2pFWU1CWUdBMVVFQ2hNUFRHOXVaMmh2Y200Z0xTQlVaWE4wTUNBWERUSXdNRFF5TnpJek1EQXhNVm9ZRHpJeApNakF3TkRBek1qTXdNREV4V2pBYU1SZ3dGZ1lEVlFRS0V3OU1iMjVuYUc5eWJpQXRJRlJsYzNRd2dnRWlNQTBHCkNTcUdTSWIzRFFFQkFRVUFBNElCRHdBd2dnRUtBb0lCQVFEWHpVdXJnUFpEZ3pUM0RZdWFlYmdld3Fvd2RlQUQKODRWWWF6ZlN1USs3K21Oa2lpUVBvelVVMmZvUWFGL1BxekJiUW1lZ29hT3l5NVhqM1VFeG1GcmV0eDBaRjVOVgpKTi85ZWFJNWRXRk9teHhpMElPUGI2T0RpbE1qcXVEbUVPSXljdjRTaCsvSWo5Zk1nS0tXUDdJZGxDNUJPeThkCncwOVdkckxxaE9WY3BKamNxYjN6K3hISHd5Q05YeGhoRm9tb2xQVnpJbnlUUEJTZkRuSDBuS0lHUXl2bGhCMGsKVHBHSzYxc2prZnFTK3hpNTlJeHVrbHZIRXNQcjFXblRzYU9oaVh6N3lQSlorcTNBMWZoVzBVa1JaRFlnWnNFbQovZ05KM3JwOFhZdURna2kzZ0UrOElXQWRBWHExeWhqRDdSSkI4VFNJYTV0SGpKUUtqZ0NlSG5HekFnTUJBQUdqCmF6QnBNQTRHQTFVZER3RUIvd1FFQXdJQ3BEQVRCZ05WSFNVRUREQUtCZ2dyQmdFRkJRY0RBVEFQQmdOVkhSTUIKQWY4RUJUQURBUUgvTURFR0ExVWRFUVFxTUNpQ0NXeHZZMkZzYUc5emRJSVZiV2x1YVc4dGMyVnlkbWxqWlM1awpaV1poZFd4MGh3Ui9BQUFCTUEwR0NTcUdTSWIzRFFFQkN3VUFBNElCQVFDbUZMMzlNSHVZMzFhMTFEajRwMjVjCnFQRUM0RHZJUWozTk9kU0dWMmQrZjZzZ3pGejFXTDhWcnF2QjFCMVM2cjRKYjJQRXVJQkQ4NFlwVXJIT1JNU2MKd3ViTEppSEtEa0Jmb2U5QWI1cC9VakpyS0tuajM0RGx2c1cvR3AwWTZYc1BWaVdpVWorb1JLbUdWSTI0Q0JIdgpnK0JtVzNDeU5RR1RLajk0eE02czNBV2xHRW95YXFXUGU1eHllVWUzZjFBWkY5N3RDaklKUmVWbENtaENGK0JtCmFUY1RSUWN3cVdvQ3AwYmJZcHlERFlwUmxxOEdQbElFOW8yWjZBc05mTHJVcGFtZ3FYMmtYa2gxa3lzSlEralAKelFadHJSMG1tdHVyM0RuRW0yYmk0TktIQVFIcFc5TXUxNkdRakUxTmJYcVF0VEI4OGpLNzZjdEg5MzRDYWw2VgotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0t
     ```

     For more information on creating a secret, see [the Kubernetes documentation.](https://kubernetes.io/docs/concepts/configuration/secret/#creating-a-secret-manually) The secret must be created in the `longhorn-system` namespace for Longhorn to access it.

     > **Note:** Make sure to use `echo -n` when generating the base64 encoding, otherwise a new line will be added at the end of the string and it will cause error when accessing the S3.

3. Click the **Backup** tab in the UI. It should report an empty list without any errors.

**Result:** Longhorn can store backups in S3. To create a backup, see [this section.](../create-a-backup)

### Using a self-signed SSL certificate for S3 communication

If you want to use a self-signed SSL certificate, you can specify AWS_CERT in the Kubernetes secret you provided to Longhorn. See the example in [Set up a Local Testing Backupstore](#set-up-a-local-testing-backupstore).
It's important to note that the certificate needs to be in PEM format, and must be its own CA. Or one must include a certificate chain that contains the CA certificate.
To include multiple certificates, one can just concatenate the different certificates (PEM files).

### Enable virtual-hosted-style access for S3 compatible Backupstore

**You may need to enable this new addressing approach for your S3 compatible Backupstore when**

1. you want to switch to this new access style right now so that you won't need to worry about [Amazon S3 Path Deprecation Plan](https://aws.amazon.com/blogs/aws/amazon-s3-path-deprecation-plan-the-rest-of-the-story/);
2. the backupstore you are using supports virtual-hosted-style access only, e.g., Alibaba Cloud(Aliyun) OSS;
3. you have configured `MINIO_DOMAIN` environment variable to [enable virtual-host-style requests for the MinIO server](https://min.io/docs/minio/linux/administration/object-management.html#path-vs-virtual-host-bucket-access);
4. the error `...... error: AWS Error: SecondLevelDomainForbidden Please use virtual hosted style to access. .....` is triggered.

**The way to enable virtual-hosted-style access**

1. Add a new field `VIRTUAL_HOSTED_STYLE` with value `true` to your backup target secret. e.g.:

    ```yaml
    apiVersion: v1
    kind: Secret
    metadata:
      name: s3-compatible-backup-target-secret
      namespace: longhorn-system
    type: Opaque
    data:
      AWS_ACCESS_KEY_ID: bG9uZ2hvcm4tdGVzdC1hY2Nlc3Mta2V5
      AWS_SECRET_ACCESS_KEY: bG9uZ2hvcm4tdGVzdC1zZWNyZXQta2V5
      AWS_ENDPOINTS: aHR0cHM6Ly9taW5pby1zZXJ2aWNlLmRlZmF1bHQ6OTAwMA==
      VIRTUAL_HOSTED_STYLE: dHJ1ZQ== # true
    ```

2. Deploy/update the secret.
3. Create correspondence backup target in `Backup and Restore > Backup Targets`.
   1. Name: The target name you want.
   2. URL: `s3://<bucket-name>@<region>/`.
   3. Credential Secret: `s3-compatible-backup-target-secret` in this example.

### Set up NFS Backupstore

Ensure that the NFS server supports NFSv4 and that the target URL points to the service.

Example:

```text
nfs://longhorn-test-nfs-svc.default:/opt/backupstore
```

The default mount options are `actimeo=1,soft,timeo=300,retry=2`.  To use other options, append the keyword "nfsOptions" and the options string to the target URL.  

Example:

```text
nfs://longhorn-test-nfs-svc.default:/opt/backupstore?nfsOptions=soft,timeo=330,retrans=3  
```

Any mount options that you specify will replace, not add to, the default options.

You can find an example NFS backupstore for testing purpose [here](https://github.com/longhorn/longhorn/blob/v{{< current-version >}}/deploy/backupstores/nfs-backupstore.yaml).

**Result:** Longhorn can store backups in NFS. To create a backup, see [this section.](../create-a-backup)

### Set up SMB/CIFS Backupstore

Before configuring a SMB/CIFS backupstore, a credential secret for the backupstore can be created and deployed by

  ```shell
  #!/bin/bash

  USERNAME=${Username of SMB/CIFS Server}
  PASSWORD=${Password of SMB/CIFS Server}

  CIFS_USERNAME=`echo -n ${USERNAME} | base64`
  CIFS_PASSWORD=`echo -n ${PASSWORD} | base64`

  cat <<EOF >>cifs_secret.yml
  apiVersion: v1
  kind: Secret
  metadata:
    name: cifs-secret
    namespace: longhorn-system
  type: Opaque
  data:
    CIFS_USERNAME: ${CIFS_USERNAME}
    CIFS_PASSWORD: ${CIFS_PASSWORD}
  EOF

  kubectl apply -f cifs_secret.yml
  ```

On the Longhorn UI, go to **Backup and Restore > Backup Targets**.

- Create or edit a backup target.
  - Set **URL** to:

    ```text
    cifs://longhorn-test-cifs-svc.default/backupstore
    ```

    The default CIFS mount option is "soft".  To use other options, append the keyword "cifsOptions" and the options string to the target URL.

    Example:

    ```text
    cifs://longhorn-test-cifs-svc.default/backupstore?cifsOptions=rsize=65536,wsize=65536,soft
    ```

    Any mount options that you specify will replace, not add to, the default options.

  - Set **Credential Secret** to:

    ```text
    cifs-secret
    ```

    This is the secret name with CIFS credentials.

You can find an example CIFS backupstore for testing purpose [here](https://github.com/longhorn/longhorn/blob/v{{< current-version >}}/deploy/backupstores/cifs-backupstore.yaml).

**Result:** Longhorn can store backups in CIFS. To create a backup, see [this section.](../create-a-backup)

### Set up Azure Blob Storage Backupstore

1. Verify that a container for the backupstore exists in [Azure Blob Storage](https://portal.azure.com/).
2. Grant the Azure service account permissions to read, write, and delete objects in the container.
  For more information, see [Manage blob containers using the Azure portal](https://learn.microsoft.com/en-us/azure/storage/blobs/blob-containers-portal) in the Microsoft documentation.

3. Go to **Home > `serviceaccount` > Security + networking > Access keys**.
4. Save the following information:

   - `Storage account name`: Maps to the `AZBLOB_ACCOUNT_NAME` field in the Kubernetes secret that you will create.
   - `Key`: Maps to the `AZBLOB_ACCOUNT_KEY` field in the Kubernetes secret that you will create.

5. Go to the Longhorn UI. In the top navigation bar, click **Backup and Restore/Backup Targets**, and create or edit a backup target.

   - Set **URL**. The target URL should look like this:

     ```text
     azblob://[your-container-name]@core.windows.net/
     ```

     Make sure that you have `/` at the end, otherwise you will get an error. A subdirectory (prefix) may be used:

     ```text
     azblob://[your-container-name]@core.windows.net/my-path/
     ```

   - Set **Credential Secret**.

     ```text
     longhorn-azblob-secret
     ```

6. Create a Kubernetes secret named `longhorn-azblob-secret`.
  This secret is used to access the backupstore in the Longhorn namespace (default: `longhorn-system`) with the following content:

   ```shell
   #!/bin/bash
   cat <<EOF >>longhorn-azblob-secret.yml
   apiVersion: v1
   kind: Secret
   metadata:
     name: longhorn-azblob-secret
     namespace: longhorn-system
   type: Opaque
   stringData:
     AZBLOB_ACCOUNT_NAME: "<Storage account name>"
     AZBLOB_ACCOUNT_KEY:  "<Key>"
     ...
     # Parameters below are used for the compatible azure server for instance `Azurite` or 
     # you have a proxy to redirect the requests.
     #AZBLOB_ENDPOINT: ""
     #AZBLOB_CERT: ""
     #HTTP_PROXY: ""
     #HTTPS_PROXY: ""
   EOF

   kubectl apply -f longhorn-azblob-secret.yml
   ```

After configuring the above settings, you can manage backups on Azure Blob storage. See [how to create backup](../create-a-backup) for details.


---

## Article: snapshots-and-backups/backup-and-restore/synchronize_backup_volumes_manually.md

---
title: Synchronize Backup Volumes Manually
weight: 6
---

After creating a backup, Longhorn creates a backup volume that corresponds to the original volume (on which the backup is based). A backup volume is an object in the backupstore that contains multiple backups of the same volume.

Earlier Longhorn versions poll and update all backup volumes at a fixed poll interval. Longhorn v1.6.2 provides a way for you to manually synchronize backup volumes with the backup target.

> **Important:** You must set up a [backup target](../set-backup-target) and verify that a backup volume was created before attempting to synchronize. Longhorn returns an error when no backup target and backup volume exist.

- Synchronize all backup volumes:
  1. On the Longhorn UI, go to **Backup**.
  1. Click **Sync All Backup Volumes**.

- Synchronize a single backup volume:
  1. On the Longhorn UI, go to **Backup**.
  1. Select a backup volume.
  1. Click **Sync Backup Volume**.

To check if synchronization was successful, click the name of the backup volume on the **Backup** page.


---

## Article: nodes-and-volumes/_index.md

---
title: Nodes and Volumes
weight: 4
---

---

## Article: nodes-and-volumes/nodes/_index.md

---
title: Nodes
weight: 1
---

---

## Article: nodes-and-volumes/nodes/default-disk-and-node-config.md

---
title: Configuring Defaults for Nodes and Disks
weight: 3
---

_Available as of v0.8.1_

This feature allows the user to customize the default disks and node configurations in Longhorn for newly added nodes using Kubernetes labels and annotations instead of the Longhorn API or UI.

Customizing the default configurations for disks and nodes is useful for scaling the cluster because it eliminates the need to configure Longhorn manually for each new node if the node contains more than one disk, or if the disk configuration is different for new nodes.

Longhorn will not keep the node labels or annotations in sync with the current Longhorn node disks or tags. Nor will Longhorn keep the node disks or tags in sync with the nodes, labels or annotations after the default disks or tags have been created.

### Adding Node Tags to New Nodes

When a node does not have a tag, you can use a node annotation to set the node tags, as an alternative to using the Longhorn UI or API.

1. Scale up the Kubernetes cluster. The newly added nodes contain no node tags.
2. Add annotations to the new Kubernetes nodes that specify what the default node tags should be. The annotation format is:

    ```
    node.longhorn.io/default-node-tags: <node tag list with JSON string format>
    ```
    For example:

    ```
    node.longhorn.io/default-node-tags: '["fast","storage"]'
    ``` 
3. Wait for Longhorn to sync the node tag automatically.

> **Result:** If the node tag list was originally empty, Longhorn updates the node with the tag list, and you will see the tags for that node updated according to the annotation. If the node already had tags, you will see no change to the tag list.
### Customizing Default Disks for New Nodes

Longhorn uses the **Create Default Disk on Labeled Nodes** setting to enable default disk customization.

If the setting is disabled, Longhorn will create a default disk using `setting.default-data-path` on all new nodes.

If the setting is enabled, Longhorn will decide to create the default disks or not, depending on the node's label value of `node.longhorn.io/create-default-disk`.

- If the node's label value is `true`, Longhorn will create the default disk using `settings.default-data-path` on the node. If the node already has existing disks, Longhorn will not change anything.
- If the node's label value is `config`, Longhorn will check for the `node.longhorn.io/default-disks-config` annotation and create default disks according to it. If there is no annotation, or if the annotation is invalid, or the label value is invalid, Longhorn will not change anything.

The value of the label will be in effect only when the setting is enabled.

If the `create-default-disk` label is not set, the default disk will not be automatically created on the new nodes when the setting is enabled.

The configuration described in the annotation only takes effect when there are no existing disks or tags on the node.

If the label or annotation fails validation, the whole annotation is ignored. 

> **Prerequisite:** The Longhorn setting **Create Default Disk on Labeled Nodes** must be enabled.
1. Add new nodes to the Kubernetes cluster.
2. Add the label to the node. Longhorn relies on the label to decide how to customize default disks:

    ```
    node.longhorn.io/create-default-disk: 'config'
    ```

3. Then add an annotation to the node. The annotation is used to specify the configuration of default disks. The format is:

    ```
    node.longhorn.io/default-disks-config: <disks configuration with JSON string format>
    ```

    For example, the following disk configuration can be specified in the annotation:

    ```
    node.longhorn.io/default-disks-config: 
    '[
        { 
            "path":"/mnt/disk1",
            "allowScheduling":true
        },
        {   
            "name":"fast-ssd-disk", 
            "path":"/mnt/disk2",
            "allowScheduling":false,
            "storageReserved":10485760,
            "tags":[
                "ssd",
                "fast"
            ]
        }
    ]'
    ```

    > **Note:** If the same name is specified for different disks, the configuration will be treated as invalid.
    
4. Wait for Longhorn to create the customized default disks automatically.

> **Result:** The disks will be updated according to the annotation.

### Launch Longhorn with multiple disks
1. Add the label to all nodes before launching Longhorn. 

    ```
    node.longhorn.io/create-default-disk: 'config'
    ```

2. Then add the disk config annotation to all nodes:

    ```
    node.longhorn.io/default-disks-config: '[ { "path":"/var/lib/longhorn", "allowScheduling":true
      }, { "name":"fast-ssd-disk", "path":"/mnt/extra", "allowScheduling":false, "storageReserved":10485760,
      "tags":[ "ssd", "fast" ] }]'
    ```
3. Deploy Longhorn with `create-default-disk-labeled-nodes: true`, check [here](../../../advanced-resources/deploy/customizing-default-settings) for customizing the default settings of Longhorn.


---

## Article: nodes-and-volumes/nodes/disks-or-nodes-eviction.md

---
title: Evicting Replicas on Disabled Disks or Nodes
weight: 6
---

Longhorn supports auto eviction for evicting the replicas on the selected disabled disks or nodes to other suitable disks and nodes. Meanwhile the same level of high availability is maintained during the eviction.

> **Note:** This eviction feature can only be enabled when the selected disks or nodes have scheduling disabled. And during the eviction time, the selected disks or nodes cannot be re-enabled for scheduling.

> **Note:** This eviction feature works for volumes that are `Attached` and `Detached`. If the volume is 'Detached', Longhorn will automatically attach it before the eviction and automatically detach it once eviction is done.

By default, `Eviction Requested` for disks or nodes is `false`. And to keep the same level of high availability during the eviction, Longhorn only evicts a replica per volume after the replica rebuild for this volume is a success.

## Select Disks or Nodes for Eviction

To evict disks for a node,

1. Head to the `Node` tab, select one of the nodes, and select `Edit Node and Disks` in the dropdown menu.
1. Make sure the disk is disabled for scheduling and set `Scheduling` to `Disable`.
2. Set `Eviction Requested` to `true` and save.

To evict a node,

1. Head to the `Node` tab, select one or more nodes, and click `Edit Node`.
1. Make sure the node is disabled for scheduling and set `Scheduling` to `Disable`.
2. Set `Eviction Requested` to `true`, and save.

## Cancel Disks or Nodes Eviction

To cancel the eviction for a disk or a node, set the corresponding `Eviction Requested` setting to `false`.

## Check Eviction Status

The `Replicas` number on the selected disks or nodes should be reduced to 0 once the eviction is a success.

If you click on the `Replicas` number, it will show the replica name on this disk. When you click on the replica name, the Longhorn UI will redirect the webpage to the corresponding volume page, and it will display the volume status. If there is any error, e.g. no space, or couldn't find another schedulable disk (schedule failure), the error will be shown. All of the errors will be logged in the Event log.

If any error happened during the eviction, the eviction will be suspended until new space has been cleared or it will be cancelled. And if the eviction is cancelled, the remaining replicas on the selected disks or nodes will remain on the disks or nodes.


---

## Article: nodes-and-volumes/nodes/multidisk.md

---
title: Multiple Disk Support
weight: 4
---

Longhorn supports using more than one disk on the nodes to store the volume data.

By default, Longhorn stores volume data in the `/var/lib/longhorn` directory on the host. If you want to use a different disk for storage, you can add a new disk and disable scheduling for the default directory. This gives you flexibility to manage storage based on your needs.

## Add a Disk

Before adding a disk to Longhorn, you need to mount it to a directory on the host of the Longhorn node.

1. **Choose a Disk:** Select the physical or virtual disk you want to use for Longhorn storage and format it with an extent-based filesystem (e.g., ext4, xfs).
2. **Mount the Disk:** Mount the disk to a directory on the host, such as `/mnt/example-disk`. Ensure the directory is accessible and properly configured.

After the disk is mounted, you can add it to Longhorn using either the  UI or the `kubectl` command-line tool.

- **Using the Longhorn UI**

  1. Go to the **Nodes** tab, select a node, and choose **Edit Disks** from the dropdown menu.
  2. Add the mount path of the disk to the disk list.

- **Using `kubectl` Command**

  1. Run `kubectl edit node.longhorn.io <node-name>` to edit the Longhorn node resource.
  2. Add the disk path to `spec.disks`. For example:

    ```yaml
    ...
    spec:
      ...
      disks:
        ...
        example-disk:
          allowScheduling: true
          diskDriver: ""
          diskType: filesystem
          evictionRequested: false
          path: /mnt/example-disk
          storageReserved: 0
          tags: []
    ...
    ```

  3. Save and exit the editor.

Once a disk is added:

- Longhorn automatically detects the storage details of the disk, such as maximum and available capacity.
- If the disk is suitable for storing volume data, Longhorn begins scheduling volumes to it.

> **Notice**:
>
> 1. You cannot add a disk path that is already in use by another Longhorn disk.
> 2. Longhorn uses the filesystem ID to detect duplicate mounts. Therefore, you cannot add a disk with the **same filesystem ID** as another disk on the same node.  
>    See: [Issue #2477](https://github.com/longhorn/longhorn/issues/2477)

### Root Disk Reservation

Optionally, you can use the `Space Reserved` field in the UI or `spec.disks.<disk-name>.storageReserved` to reserve a portion of disk space (in bytes) for other purposes. This reserved space will not be used by Longhorn for volume data.

To maintain node stability when compute resources (for example memory or disk) are under pressure, the `kubelet` requires some space to remain free. If these critical resources are exhausted, it can lead to node instability.

Longhorn **reserves 30% of the root disk space (`/var/lib/longhorn`) by default** to prevent issues like `DiskPressure` conditions from the `kubelet`, particularly after scheduling multiple volumes. This behaviour is controlled by the `storage-reserved-percentage-for-default-disk` setting.

### Use an Alternative Path for a Disk on the Node

If you prefer to use a different path for a disk (rather than the original mount point), you can use `mount --bind` to create an alternative path. Do **not** use symbolic link (`lh -s`), as these are not properly resolved inside Longhorn pods.

Make sure the alternative path is remounted after a node reboot, for example, by adding it to `/etc/fstab`.

## Remove a Disk

Nodes and disks can be excluded from future scheduling. Note that any storage already scheduled on the node will not be automatically released when scheduling is disabled for the node.

To remove a disk:

- Disable scheduling for the disk.
- Ensure there are **no replicas or backing images** left on the disk, including any in an error state. For instructions on how to evict replicas from disabled disks, see [Select Disks or Nodes for Eviction](../disks-or-nodes-eviction/#select-disks-or-nodes-for-eviction).

Once the disk is empty and scheduling is disabled, you can safely remove it from the node configuration.

## Configuration

There are two global settings affect the scheduling of the volume.

- `StorageOverProvisioningPercentage` defines the maximum total storage that can be **scheduled** on a disk, relative to its usable capacity. The formula is:

    ```bash
    ScheduledStorage / (MaximumStorage - ReservedStorage)
    ```

    The default is `100` (%).

	On a 200 GiB disk with 50 GiB reserved, Longhorn sees 150 GiB of usable space. With the default setting, it can schedule up to 150 GiB of volume data.

    Since workloads typically don’t consume the entire allocated volume size, and Longhorn uses sparse files to store data, increasing this setting is generally safe and can help optimize disk utilization.

- `StorageMinimalAvailablePercentage` specifies the minimum percentage of free space that must remain on a disk in order to schedule new replicas. The formula is:

    ```bash
    AvailableStorage / MaximumStorage
    ```

    The default is `25` (%).

    For a 200 GiB disk with 50 GiB reserved, Longhorn will stop scheduling new replicas if available space falls below 37.5 GiB  (25% of 150 GiB). A new volume also will not be scheduled if its size would push available space below that limit.

    This setting helps prevent disks from becoming too full, which could lead to scheduling failures or volume operation issues.

> **Warning**:
> 
> Currently, Longhorn cannot fully enforce the `StorageMinimalAvailablePercentage` limit in all scenarios because:
>
> - Longhorn volumes may use more space than their requested size, especially when snapshots are taken.
> - Longhorn allows over-provisioning by default.


---

## Article: nodes-and-volumes/nodes/node-conditions.md

---
title: Node Conditions
weight: 7
---

Node conditions describe the status of all worker nodes and are used to check the environment settings of worker nodes to identify potential issues before any system impact.

Node conditions:

- `Ready`:  
  Indicates that the node is ready for Longhorn operations, including that a `longhorn-manager` pod is running on this node, the Kubernetes node is ready, and there is no physical resources pressure.  

- `Schedulable`:  
  Indicated that the node is not cordoned and workload can be scheduled to this node.

- `MountPropagation`:  
  Indicates that the node supports mount propagation. This is necessary for sharing of volumes mounted by a container with other containers in the same Longhorn pod, or to other Longhorn pods on the same node.  

- `Multipathd`:  
  Confirms if the `multipathd` service is not running on the node, which may affect the pod with the volume startup. See [Troubleshooting: `MountVolume.SetUp failed for volume` due to multipathd on the node](../../../../../kb/troubleshooting-volume-with-multipath).  

- `RequiredPackages`:  
  Checks if all required packages ([NFS client](../../../deploy/install/#installing-nfsv4-client), [iSCSI tool](../../../deploy/install/#installing-open-iscsi), [cryptsetup](../../../deploy/install/#installing-cryptsetup-and-luks), [dmsetup](../../../deploy/install/#installing-device-mapper-userspace-tool)) exist for Longhorn  

- `NFSClientInstalled`:  
  Identifies if any of the following NFS clients are supported: `v4.2`, `v4.1`, or `v4.0`. NFS client is required for RWX volume and backup.  

- `KernelModulesLoaded`:  
  Identifies if the following Kernel modules are loaded:
  - `dm_crypt`: Is required for the volume and backing image encryption.  
  - For engine v2 only:
    - `vfio_pci`: Is required for SPDK and PCI device management
    - `uio_pci_generic`: Is required for SPDK UIO support
    - `nvme_tcp`: Is required for NVMe-over-TCP device usage

- `HugePagesAvailable`:
  Indicates whether the node is properly configured with HugePages (2Mi) as required by the Longhorn v2 data engine. This includes verifying that:
  - HugePages (2Mi) are registered as a Kubernetes resource (`hugepages-2Mi`).
  - The configured HugePages capacity meets or exceeds the value defined in the `v2-data-engine-hugepage-limit` setting.

Node conditions do not block the Longhorn deployment but they result in warnings in the Longhorn `Node` resource.
For more information, see [Longhorn Installation Requirements](../../../deploy/install/#installation-requirements).


---

## Article: nodes-and-volumes/nodes/node-space-usage.md

---
title: Node Space Usage
weight: 1
---

In this section, you'll have a better understanding of the space usage info presented by the Longhorn UI. 


### Whole Cluster Space Usage

In `Dashboard` page, Longhorn will show you the cluster space usage info:

{{< figure src="/img/screenshots/volumes-and-nodes/space-usage-info-dashboard-page.png" >}}

`Schedulable`: The actual space that can be used for Longhorn volume scheduling.

`Reserved`: The space reserved for other applications and system.

`Used`: The actual space that has been used by Longhorn, system, and other applications.

`Disabled`: The total space of the disks/nodes on which Longhorn volumes are not allowed for scheduling.

### Space Usage of Each Node

In `Node` page, Longhorn will show the space allocation, schedule, and usage info for each node:

{{< figure src="/img/screenshots/volumes-and-nodes/space-usage-info-node-page.png" >}}

`Size` column: The **max actual available space** that can be used by Longhorn volumes. It equals the total disk space of the node minus reserved space. 

`Allocated` column: The left number is the size that has been used for **volume scheduling**, and it does not mean the space has been used for the Longhorn volume data store. The right number is the **max** size for volume scheduling, which the result of `Size` multiplying `Storage Over Provisioning Percentage`. (In the above illustration, `Storage Over Provisioning Percentage` is 500.) Hence, the difference between the 2 numbers (let's call it as the allocable space) determines if a volume replica can be scheduled to this node.

`Used` column: The left part indicates the currently used space of this node. The whole bar indicates the total space of the node.

Notice that the allocable space may be greater than the actual available space of the node when setting `Storage Over Provisioning Percentage` to a value greater than 100. If the volumes are heavily used and lots of historical data will be stored in the volume snapshots, please be careful about using a large value for this setting. For more info about the setting, see [here](../../../references/settings/#storage-over-provisioning-percentage) for details. 


---

## Article: nodes-and-volumes/nodes/scheduling.md

---
title: Scheduling
weight: 5
---

In this section, you'll learn how Longhorn schedules replicas based on multiple factors.

## Scheduling Policy

Longhorn's scheduling policy has two stages. The scheduler only goes to the next stage if the previous stage is satisfied. Otherwise, the scheduling will fail.

If any tag has been set in order to be selected for scheduling, the node tag and the disk tag have to match when the node or the disk is selected.

The first stage is the **node and zone selection stage.** Longhorn will filter the node and zone based on the `Replica Node Level Soft Anti-Affinity` and `Replica Zone Level Soft Anti-Affinity` settings.

The second stage is the **disk selection stage.** Longhorn will filter the disks that satisfy the first stage based on the `Replica Disk Level Soft Anti-Affinity`, `Storage Minimal Available Percentage`, `Storage Over Provisioning Percentage`, and other disk-related factors like requested disk space.

### The Node and Zone Selection Stage

Longhorn evaluates which **nodes** are suitable for scheduling a new replica based on a series of criteria. The decision-making process follows a specific order to ensure optimal placement for fault tolerance. 

#### 1. Node Tag Matching

Longhorn first checks for node selector tags on the volume.
- If the volume has node selector tags, only nodes with matching tags are eligible.
- If the volume has **no node selector**, the behavior depends on the setting **Allow Empty Node Selector Volume**:
    - `true` (default): Schedules on nodes **with or without tags**.
    - `false`: Schedules **only** on nodes **without tags**.

#### 2. Cordoned Node Handling

The setting **Disable Scheduling On Cordoned Node** determines whether cordoned nodes are eligible for replica scheduling:
- `true` (default): Cordoned nodes are **excluded**.
- `false`: Cordoned nodes are **eligible**.

#### 3.  Anti-Affinity Rules Across Nodes and Zones

Longhorn prioritizes spreading replicas across different **nodes** and **zones** to improve fault tolerance. A **"new"** node or zone is one that does **not** currently host any replica of the volume, while an **"existing"** node or zone already hosts a replica of the volume. The selection logic proceeds in the following order:

The scheduler attempts to place the new replica in the most "isolated" location possible, following this hierarchy of preference:
1.  **New Node in a New Zone** (most preferred)
2.  **New Node in an Existing Zone**
3.  **Existing Node in an Existing Zone** (least preferred)

The following table details the required settings for a replica to be scheduled in each scenario:

| **Scenario** | **Replica Zone Level Soft Anti-Affinity** | **Replica Node Level Soft Anti-Affinity** | **Scheduler Action** |
| :--- | :--- | :--- | :--- |
| **New Node in a New Zone** | `false` | `false` | **Schedules** the replica. |
| | Any other value | Any other value | **Does not** schedule the replica. |
| **New Node in an Existing Zone** | `true` | `false` | **Schedules** the replica if no new zone is available. |
| | Any other value | Any other value | **Does not** schedule the replica. |
| **Existing Node in an Existing Zone** | `true` | `true` | **Schedules** the replica if no other options are available. |
| | Any other value | Any other value | **Does not** schedule the replica. |

### Disk Selection Stage

Once the node and zone stage is satisfied, Longhorn decides whether it can schedule the replica on any disk of the selected node. It checks the available disks based on matching tags, total disk space, and available disk space. It also considers whether another replica already exists and the anti-affinity settings.

Longhorn checks all available disks on the selected node to ensure they meet the following criteria:

1.  **Disk Tag Matching**:
    - If the volume has disk tags, the disk must match any specified tags required for the replica.
    - If the volume has no disk tags, the behavior depends on the setting **Allow Empty Disk Selector Volume**:
        - `true` (default): Allows scheduling on disks **with or without tags**.
        - `false`: Only allows scheduling on disks **without tags**.
2.  **Available Space Check**:
    - The disk must have sufficient available space based on the configured `Storage Minimal Available Percentage`.
3.  **Anti-Affinity Settings**:
    - **Hard Anti-Affinity**: Prevents scheduling a replica on a disk that already hosts another replica of the same volume.
    - **Soft Anti-Affinity** (when enabled): Prefers scheduling the replica on a disk without an existing replica, even if it’s a less optimal choice in terms of space or other factors.
4.  **Space Conditions**: Two formulas determine if a disk is schedulable:
    - **Actual Space Usage Condition**: Ensures sufficient usable storage remains after accounting for currently used space.
        - Formula: `(Storage Available - Actual Size) > (Storage Maximum × Minimal Available Percentage) / 100`
    - **Scheduling Space Condition**: Ensures the replica’s size (plus any scheduled but unwritten data) fits within the over-provisioning limit.
        - Formula: `(Size + Storage Scheduled) ≤ ((Storage Maximum - Storage Reserved) × Over Provisioning Percentage) / 100`

    > **Note:** During disk evaluation, since no specific replica is being scheduled yet, `Actual Size` and `Size` are temporarily treated as `0` in these formulas.

If any of these conditions fail including disk tag, anti-affinity, or space requirements, the disk is marked unschedulable, and Longhorn will not place the replica on it.

If either condition fails or the disk does not meet tag or anti-affinity requirements, it is marked unschedulable, and Longhorn will not place the replica on that disk.

#### Example Scenario

Consider a node (**Node A**) with two disks:
- **Disk X**: 1 GB available, 4 GB max space
- **Disk Y**: 2 GB available, 8 GB max space

##### Stage 1: Initial Disk Evaluation

During the initial disk selection stage, Longhorn performs a basic check on all available disks. At this point, no specific replica has been selected, so `Actual Size` and `Size` are treated as `0`.

**Disk X Evaluation**
- **Available Space**: 1 GB
- **`Storage Minimal Available Percentage`**: 25% (default)
- **Minimum required available space**: `(4 GB × 25) / 100 = 1 GB`
- **Result**: **Disk X** fails the `Actual Space Usage Condition` because its available space (1 GB) is **not greater than** the minimum required (1 GB). Therefore, Disk X is not schedulable unless the `Storage Minimal Available Percentage` is set to 0.

**Disk Y Evaluation**
- **Available Space**: 2 GB
- **`Storage Minimal Available Percentage`**: 10%
- **Minimum required available space**: `(8 GB × 10) / 100 = 0.8 GB`
- **Result**: **Disk Y** passes the `Actual Space Usage Condition` because its available space (2 GB) is greater than the minimum required (0.8 GB).

Next, we check the **Scheduling Space Condition**:
- **Scheduled Space**: 2 GB
- **`Storage Reserved`**: 1 GB
- **`Over Provisioning Percentage`**: 100% (default)
- **Max Provisionable Storage**: `(8 GB - 1 GB) × 100 / 100 = 7 GB`
- **Result**: **Disk Y** passes the `Scheduling Space Condition` because the currently scheduled space (2 GB) is less than the max provisionable storage (7 GB).

Since Disk Y passes all conditions, it is marked as a schedulable disk candidate.

##### Stage 2: Anti-Affinity Rules

Let's assume both Disk X and Disk Y pass the initial space checks and Disk X already hosts a replica for the same volume.

**Hard Anti-Affinity**
- Longhorn will not schedule the new replica on Disk X. It will instead attempt to schedule it on Disk Y.
* If Disk Y is not suitable (e.g., mismatched disk tags), scheduling for this replica will fail.

**Soft Anti-Affinity**
- If **soft anti-affinity** is enabled, Longhorn **prefers** to schedule the replica on Disk Y to avoid co-locating replicas.
- However, if Disk Y is unsuitable for any reason, Longhorn **may still schedule** the replica on Disk X. This allows for sharing a disk as a fallback option when no other viable candidates are available.

## Settings

For more information on settings that are relevant to scheduling replicas on nodes and disks, refer to the settings reference:

- [Disable Scheduling On Cordoned Node](../../../references/settings/#disable-scheduling-on-cordoned-node)
- [Replica Soft Anti-Affinity](../../../references/settings/#replica-node-level-soft-anti-affinity) (also called Replica Node Level Soft Anti-Affinity)
- [Replica Zone Level Soft Anti-Affinity](../../../references/settings/#replica-zone-level-soft-anti-affinity)
- [Replica Disk Level Soft Anti-Affinity](../../../references/settings/#replica-disk-level-soft-anti-affinity)
- [Storage Minimal Available Percentage](../../../references/settings/#storage-minimal-available-percentage)
- [Storage Over Provisioning Percentage](../../../references/settings/#storage-over-provisioning-percentage)
- [Allow Empty Node Selector Volume](../../../references/settings/#allow-empty-node-selector-volume)
- [Allow Empty Disk Selector Volume](../../../references/settings/#allow-empty-disk-selector-volume)

## Notice
Longhorn relies on label `topology.kubernetes.io/zone=<Zone name of the node>` or `topology.kubernetes.io/region=<Region name of the node>` in the Kubernetes node object to identify the zone/region.

Since these are reserved and used by Kubernetes as [well-known labels](https://kubernetes.io/docs/reference/labels-annotations-taints/#topologykubernetesiozone).


---

## Article: nodes-and-volumes/nodes/storage-tags.md

---
title: Storage Tags
weight: 2
---

## Overview

The storage tag feature enables only certain nodes or disks to be used for storing Longhorn volume data. For example, performance-sensitive data can use only the high-performance disks which can be tagged as `fast`, `ssd` or `nvme`, or only the high-performance node tagged as `baremetal`.

This feature supports both disks and nodes. 

## Setup

The tags can be set up using the Longhorn UI:

1. *Node -> Select one node -> Edit Node and Disks*
2. Click `+New Node Tag` or `+New Disk Tag` to add new tags.

All the existing scheduled replica on the node or disk won't be affected by the new tags.

## Usage

When multiple tags are specified for a volume, the disk and the node (the disk belong to) must have all the specified tags to become usable.

### UI

When creating a volume, specify the disk tag and node tag in the UI.

### Kubernetes

Use Kubernetes StorageClass parameters to specify tags.

You can specify tags in the default Longhorn StorageClass by adding parameter `nodeSelector: "storage,fast"` in the ConfigMap named `longhorn-storageclass`. 
For example:

```yaml
apiVersion: v1
kind: ConfigMap
data:
  storageclass.yaml: |
    kind: StorageClass
    apiVersion: storage.k8s.io/v1
    metadata:
      name: longhorn
      annotations:
        storageclass.kubernetes.io/is-default-class: "true"
    provisioner: driver.longhorn.io
    allowVolumeExpansion: true
    reclaimPolicy: "Delete"
    volumeBindingMode: Immediate
    parameters:
      numberOfReplicas: "3"
      staleReplicaTimeout: "480"
      diskSelector: "ssd"
      nodeSelector: "storage,fast"
```
If Longhorn is installed via Helm, you can achieve that by editing `persistence.defaultNodeSelector` in [values.yaml](https://github.com/longhorn/longhorn/blob/v{{< current-version >}}/chart/values.yaml).

Alternatively, a custom storageClass setting can be used, e.g.:
```yaml
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: longhorn-fast
provisioner: driver.longhorn.io
parameters:
  numberOfReplicas: "3"
  staleReplicaTimeout: "480" # 8 hours in minutes
  diskSelector: "ssd"
  nodeSelector: "storage,fast"
```

## History
* [Original feature request](https://github.com/longhorn/longhorn/issues/311)
* Available since v0.6.0


---

## Article: nodes-and-volumes/volumes/_index.md

---
title: Volumes
weight: 2
---

---

## Article: nodes-and-volumes/volumes/create-volumes.md

---
title: Create Longhorn Volumes
weight: 1
---

In this tutorial, you'll learn how to create Kubernetes persistent storage resources of persistent volumes (PVs) and persistent volume claims (PVCs) that correspond to Longhorn volumes. You will use kubectl to dynamically provision storage for workloads using a Longhorn storage class. For help creating volumes from the Longhorn UI, refer to [this section.](#creating-longhorn-volumes-with-the-longhorn-ui)

> This section assumes that you understand how Kubernetes persistent storage works. For more information, see the [Kubernetes documentation.](https://kubernetes.io/docs/concepts/storage/persistent-volumes/)

### Creating Longhorn Volumes with kubectl

First, you will create a Longhorn StorageClass. The Longhorn StorageClass contains the parameters to provision PVs.

Next, a PersistentVolumeClaim is created that references the StorageClass. Finally, the PersistentVolumeClaim is mounted as a volume within a Pod.

When the Pod is deployed, the Kubernetes master will check the PersistentVolumeClaim to make sure the resource request can be fulfilled. If storage is available, the Kubernetes master will create the Longhorn volume and bind it to the Pod.

1. Use following command to create a StorageClass called `longhorn`:

    ```
    kubectl create -f https://raw.githubusercontent.com/longhorn/longhorn/v{{< current-version >}}/examples/storageclass.yaml
    ```

    The following example StorageClass is created:

    ```
    kind: StorageClass
    apiVersion: storage.k8s.io/v1
    metadata:
      name: longhorn
    provisioner: driver.longhorn.io
    allowVolumeExpansion: true
    parameters:
      numberOfReplicas: "3"
      staleReplicaTimeout: "2880" # 48 hours in minutes
      fromBackup: ""
      fsType: "ext4"
    #  backupTargetName: "default"
    #  mkfsParams: "-I 256 -b 4096 -O ^metadata_csum,^64bit"
    #  diskSelector: "ssd,fast"
    #  nodeSelector: "storage,fast"
    #  recurringJobSelector: '[
    #   {
    #     "name":"snap",
    #     "isGroup":true,
    #   },
    #   {
    #     "name":"backup",
    #     "isGroup":false,
    #   }
    #  ]'
    ```

    In particular, starting with v1.4.0, the parameter `mkfsParams` can be used to specify filesystem format options for each StorageClass.
    Starting with v1.8.0, the parameter `backupTargetName` can be used to specify the backup target. The name of the default backup target (`default`) is used if `backupTargetName` is not specified.
    Parameters may be omitted from the StorageClass specification.  When the storage class is used to create a PV and a volume, parameters that are not specified will generally be set using a default value taken from the global settings.  See [here](../../../references/storage-class-parameters) for the list of storage class parameters, and [here](../../../references/settings) for the full list of global settings.

2. Create a Pod that uses Longhorn volumes by running this command:

    ```
    kubectl create -f https://raw.githubusercontent.com/longhorn/longhorn/v{{< current-version >}}/examples/pod_with_pvc.yaml
    ```

    A Pod named `volume-test` is launched, along with a PersistentVolumeClaim named `longhorn-volv-pvc`. The PersistentVolumeClaim references the Longhorn StorageClass:

    ```
    apiVersion: v1
    kind: PersistentVolumeClaim
    metadata:
      name: longhorn-volv-pvc
    spec:
      accessModes:
        - ReadWriteOnce
      storageClassName: longhorn
      resources:
        requests:
          storage: 2Gi
    ```

    The persistentVolumeClaim is mounted in the Pod as a volume:

    ```
    apiVersion: v1
    kind: Pod
    metadata:
      name: volume-test
      namespace: default
    spec:
      containers:
      - name: volume-test
        image: nginx:stable-alpine
        imagePullPolicy: IfNotPresent
        volumeMounts:
        - name: volv
          mountPath: /data
        ports:
        - containerPort: 80
      volumes:
      - name: volv
        persistentVolumeClaim:
          claimName: longhorn-volv-pvc
    ```
More examples are available [here.](../../../references/examples)

### Binding Workloads to PVs without a Kubernetes StorageClass

It is possible to use a Longhorn StorageClass to bind a workload to a PV without creating a StorageClass object in Kubernetes.

Since the Storage Class is also a field used to match a PVC with a PV, which doesn't have to be created by a Provisioner, you can create a PV manually with a custom StorageClass name, then create a PVC asking for the same StorageClass name.

When a PVC requests a StorageClass that does not exist as a Kubernetes resource, Kubernetes will try to bind your PVC to a PV with the same StorageClass name. The StorageClass will be used like a label to find the matching PV, and only existing PVs labeled with the StorageClass name will be used.

If the PVC names a StorageClass, Kubernetes will:

1. Look for an existing PV that has the label matching the StorageClass
2. Look for an existing StorageClass Kubernetes resource. If the StorageClass exists, it will be used to create a PV.

### Creating Longhorn Volumes with the Longhorn UI

Since the Longhorn volume already exists while creating PV/PVC, a StorageClass is not needed for dynamically provisioning Longhorn volume. However, the field `storageClassName` should be set in PVC/PV, to be used for PVC binding purpose. It's unnecessary for users to create the related StorageClass object.

By default the StorageClass for Longhorn created PV/PVC is `longhorn-static`. Users can modify it in `Setting - General - Default Longhorn Static StorageClass Name` as they need.

Users need to manually delete PVC and PV created by Longhorn.

### PV/PVC Creation for Existing Longhorn Volume

Now users can create PV/PVC via our Longhorn UI for the existing Longhorn volumes.
Only detached volume can be used by a newly created pod.

### The Failure of the Longhorn Volume Creation

Creating a Longhorn volume can fail for different reasons. The issues are categorized into:

- insufficient storage
- disk not found
- disks are unavailable
- tags not fulfilled
- node not found
- nodes are unavailable
- none of the node candidates contains a ready engine image
- hard affinity cannot be satisfied
- replica scheduling failed
- unused failed replica is not supported
- replica already scheduled
- longhorn client operation failed
- incompatible volume size

The failure results in the workload failing to use the provisioned PV and showing a warning message
```
# kubectl describe pod workload-test

Events:
  Type     Reason              Age                From                     Message
  ----     ------              ----               ----                     -------
  Warning  FailedAttachVolume  14s (x8 over 82s)  attachdetach-controller  AttachVolume.Attach
  failed for volume "pvc-e130e369-274d-472d-98d1-f6074d2725e8" : rpc error: code = Aborted
  desc = volume pvc-e130e369-274d-472d-98d1-f6074d2725e8 is not ready for workloads
```

In order to help users understand the error causes, Longhorn summarizes them in the PV annotation, `longhorn.io/volume-scheduling-error`. Failures are combined in this annotation and separated by a semicolon, for example, `longhorn.io/volume-scheduling-error: insufficient storage;disks are unavailable`. The annotation can be checked by using `kubectl describe pv <pvc name>`.
```
# kubectl describe pv pvc-e130e369-274d-472d-98d1-f6074d2725e8
Name:            pvc-e130e369-274d-472d-98d1-f6074d2725e8
Labels:          <none>
Annotations:     longhorn.io/volume-scheduling-error: insufficient storage
                 pv.kubernetes.io/provisioned-by: driver.longhorn.io

...

```


---

## Article: nodes-and-volumes/volumes/delete-volumes.md

---
title: Delete Longhorn Volumes
weight: 2
---
Once you are done utilizing a Longhorn volume for storage, there are a number of ways to delete the volume, depending on how you used the volume.

## Deleting Volumes Through Kubernetes
> **Note:** This method only works if the volume was provisioned by a StorageClass and the PersistentVolume for the Longhorn volume has its Reclaim Policy set to Delete.

You can delete a volume through Kubernetes by deleting the PersistentVolumeClaim that uses the provisioned Longhorn volume. This will cause Kubernetes to clean up the PersistentVolume and then delete the volume in Longhorn.

## Deleting Volumes Through Longhorn
All Longhorn volumes, regardless of how they were created, can be deleted through the Longhorn UI.

To delete a single volume, go to the Volumes page in the UI. Under the Operation dropdown, select Delete. You will be prompted with a confirmation before deleting the volume.

To delete multiple volumes at the same time, you can check multiple volumes on the Volumes page and select Delete at the top.

> **Note:** If Longhorn detects that a volume is tied to a PersistentVolume or PersistentVolumeClaim, then these resources will also be deleted once you delete the volume. You will be warned in the UI about this before proceeding with deletion. Longhorn will also warn you when deleting an attached volume, since it may be in use.


---

## Article: nodes-and-volumes/volumes/detaching-volumes.md

---
title: Detach Longhorn Volumes
weight: 3
---

Shut down all Kubernetes Pods using Longhorn volumes in order to detach the volumes. The easiest way to achieve this is by deleting all workloads and recreate them later after upgrade. If this is not desirable, some workloads may be suspended.

In this section, you'll learn how each workload can be modified to shut down its pods.

#### Deployment
Edit the deployment with `kubectl edit deploy/<name>`.

Set `.spec.replicas` to `0`.

#### StatefulSet
Edit the statefulset with `kubectl edit statefulset/<name>`.

Set `.spec.replicas` to `0`.

#### DaemonSet
Edit the daemonset with `kubectl edit ds/<name>`.

Add a nodeSelector to the pod spec:
```yaml
spec:
  template:
    spec:
      nodeSelector:
        no-schedule: "true"
```

#### Pod
Delete the pod with `kubectl delete pod/<name>`.

There is no way to suspend a pod not managed by a workload controller.

#### CronJob
Edit the cronjob with `kubectl edit cronjob/<name>`.

Set `.spec.suspend` to `true`.

Wait for any currently executing jobs to complete, or terminate them by deleting relevant pods.

#### Job
Consider allowing the single-run job to complete.

Otherwise, delete the job with `kubectl delete job/<name>`.

#### ReplicaSet
Edit the replicaset with `kubectl edit replicaset/<name>`.

Set `.spec.replicas` to `0`.

#### ReplicationController
Edit the replicationcontroller with `kubectl edit rc/<name>`.

Set `.spec.replicas` to `0`.

Wait for the volumes using by the Kubernetes to complete detaching.

Then detach all remaining volumes from Longhorn UI. These volumes were most likely created and attached outside of Kubernetes via Longhorn UI or REST API.

---

## Article: nodes-and-volumes/volumes/expansion.md

---
title: Volume Expansion
weight: 6
---

Volumes are expanded in two stages. First, Longhorn resizes the block device, then it expands the filesystem.

Since v1.4.0, Longhorn supports online expansion. Most of the time Longhorn can directly expand an attached volumes without limitations, no matter if the volume is being R/W or rebuilding.

If the volume was not expanded though the CSI interface (e.g. for Kubernetes older than v1.16), the capacity of the corresponding PVC and PV won't change.

## Prerequisite

- For offline expansion, the Longhorn version must be v0.8.0 or higher.
- For online expansion, the Longhorn version must be v1.4.0 or higher.

## Expand a Longhorn volume

There are two ways to expand a Longhorn volume: with a PersistentVolumeClaim (PVC) and with the Longhorn UI.

#### Via PVC

This method is applied only if:

- The PVC is dynamically provisioned by the Kubernetes with Longhorn StorageClass.
- The field `allowVolumeExpansion` should be `true` in the related StorageClass.

This method is recommended if it's applicable, because the PVC and PV will be updated automatically and everything is kept consistent after expansion.

Usage: Find the corresponding PVC for Longhorn volume, then modify the requested `spec.resources.requests.storage` of the PVC:

```
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"v1","kind":"PersistentVolumeClaim","metadata":{"annotations":{},"name":"longhorn-simple-pvc","namespace":"default"},"spec":{"accessModes":["ReadWriteOnce"],"resources":{"requests":{"storage":"1Gi"}},"storageClassName":"longhorn"}}
    pv.kubernetes.io/bind-completed: "yes"
    pv.kubernetes.io/bound-by-controller: "yes"
    volume.beta.kubernetes.io/storage-provisioner: driver.longhorn.io
  creationTimestamp: "2019-12-21T01:36:16Z"
  finalizers:
  - kubernetes.io/pvc-protection
  name: longhorn-simple-pvc
  namespace: default
  resourceVersion: "162431"
  selfLink: /api/v1/namespaces/default/persistentvolumeclaims/longhorn-simple-pvc
  uid: 0467ae73-22a5-4eba-803e-464cc0b9d975
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
  storageClassName: longhorn
  volumeMode: Filesystem
  volumeName: pvc-0467ae73-22a5-4eba-803e-464cc0b9d975
status:
  accessModes:
  - ReadWriteOnce
  capacity:
    storage: 1Gi
  phase: Bound
```

#### Via Longhorn UI

Usage: On the volume page of Longhorn UI, click `Expand` for the volume.

## Filesystem expansion

Longhorn will try to expand the file system only if:

- The expanded size should be greater than the current size.
- There is a Linux filesystem in the Longhorn volume.
- The filesystem used in the Longhorn volume is one of the following:
    - ext4
    - xfs
- The expanded size must be less than the maximum file size allowed by the file system (for example, 16TiB for `ext4`).
- The Longhorn volume is using the block device frontend.

## Corner cases

#### Handling Volume Revert

If a volume is reverted to a snapshot with smaller size, the frontend of the volume is still holding the expanded size. But the filesystem size will be the same as that of the reverted snapshot. In this case, you will need to handle the filesystem manually:

1. Attach the volume to a random node.
2. Log in to the corresponding node, and expand the filesystem.

    If the filesystem is `ext4`, the volume might need to be [mounted](https://linux.die.net/man/8/mount) and [umounted](https://linux.die.net/man/8/umount) once before resizing the filesystem manually. Otherwise, executing `resize2fs` might result in an error:

    ```
    resize2fs: Superblock checksum does not match superblock while trying to open ......
    Couldn't find valid filesystem superblock.
    ```

    Follow the steps below to resize the filesystem:

    ```
    mount /dev/longhorn/<volume name> <arbitrary mount directory>
    umount /dev/longhorn/<volume name>
    mount /dev/longhorn/<volume name> <arbitrary mount directory>
    resize2fs /dev/longhorn/<volume name>
    umount /dev/longhorn/<volume name>
    ```

3. If the filesystem is `xfs`, you can directly mount, then expand the filesystem.

    ```
    mount /dev/longhorn/<volume name> <arbitrary mount directory>
    xfs_growfs <the mount directory>
    umount /dev/longhorn/<volume name>
    ```

#### Encrypted volume

Longhorn support for online expansion depends on Kubernetes.
- Kubernetes natively supports [authenticated CSI storage resizing](https://kubernetes.io/blog/2023/12/15/csi-node-expand-secret-support-ga/) starting in v1.29.
- In [Kubernetes v1.25 to v1.28](https://kubernetes.io/blog/2022/09/21/kubernetes-1-25-use-secrets-while-expanding-csi-volumes-on-node-alpha/), the feature gate `CSINodeExpandSecret` is required.
  You can enable online expansion for encrypted volumes by specifying the following [encryption parameters in the StorageClass](../../../advanced-resources/security/volume-encryption#setting-up-kubernetes-secrets-and-storageclasses):

- `csi.storage.k8s.io/node-expand-secret-name`
- `csi.storage.k8s.io/node-expand-secret-namespace`

If you cannot enable it but still prefer to do online expansion, you can:
1. Login the node host the encrypted volume is attached to.
2. Execute `cryptsetup resize <volume name>`. The passphrase this command requires is the field `CRYPTO_KEY_VALUE` of the corresponding secret.
3. Expand the filesystem.

#### RWX volume

From v1.8.0, Longhorn supports fully automatic online expansion of the filesystem (NFS) for RWX volumes.  The feature requires the v1.8.0 versions of these components to be running:

- Longhorn-Manager
- CSI plugin
- Share Manager, which manages the NFS export

If you have upgraded from a previous version, the Share Manager pods (one for each RWX volume) are not upgraded automatically, to avoid disruption during the upgrade.

After growing the block device, the CSI layer sends a resize command to the Share Manager to grow the filesystem within the block device.  With a down-rev share-manager, the command fails with an "unimplemented" error code and so no expansion happens.  To get the right image before the expansion, the simplest thing is to force a restart of the pod.  Identify the Share Manager pod of the RWX volume (typically named `share-manager-<volume name>`) and delete it:

```shell
kubectl -n longhorn-system delete pod <the share manager pod>
```

The pod will automatically be recreated using the appropriate version, and the expansion completes.  Further expansions will not require any further intervention.

##### Offline

It's still possible to expand the RWX volume offline using these steps:

1. Detach the RWX volume by scaling down the workload to `replicas=0`. Ensure that the volume is fully detached.

1. After the scale command returns, run the following command and verify that the state is `detached`.
    ```shell
    kubectl -n longhorn-system get volume <volume-name>
    ```
1. Expand the block device using either the PVC or the Longhorn UI.

1. Scale up the workload.

The reattached volume will have the expanded size.  Furthermore, the Share Manager pod will be recreated with the current version.




---

## Article: nodes-and-volumes/volumes/iscsi.md

---
title: Use Longhorn Volume as an iSCSI Target
weight: 4
---

Longhorn supports iSCSI target frontend mode. You can connect to it
through any iSCSI client, including `open-iscsi`, and virtual machine
hypervisor like KVM, as long as it's in the same network as the Longhorn system.

The Longhorn CSI driver doesn't support iSCSI mode.

To start a volume with the iSCSI target frontend mode, select `iSCSI` as the frontend when [creating the volume.](../create-volumes)

After the volume has been attached, you will see something like the following in the `endpoint` field:

```text
iscsi://10.42.0.21:3260/iqn.2014-09.com.rancher:testvolume/1
```

In this example,

- The IP and port is `10.42.0.21:3260`.
- The target name is `iqn.2014-09.com.rancher:testvolume`.
- The volume name is `testvolume`.
- The LUN number is 1. Longhorn always uses LUN 1.

The above information can be used to connect to the iSCSI target provided by Longhorn using an iSCSI client.


---

## Article: nodes-and-volumes/volumes/pvc-ownership-and-permission.md

---
title: Longhorn PVC Ownership and Permission
weight: 1
---

Kubernetes supports the 2 [volume modes](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#volume-mode) for PVC: Filesystem and Block.
When a pod defines the security context and requests a Longhorn PVC, Kubernetes will handle the ownership and permission modification for the PVC differently based on the volume mode.

### Longhorn PVC with Filesystem Volume Mode

Because the Longhorn CSI driver `csiDriver.spec.fsGroupPolicy` is set to `ReadWriteOnceWithFSType`, the Kubelet attempts to change the ownership and permission of a Longhorn PVC in the following manner:
1. Check `pod.spec.securityContext.fsGroup`.
* If non-empty, continue to the next step.
* If empty, the Kubelet doesn't attempt to change the ownership and permission for the volume.
1. Check `fsType` of the PV and `accessModes` of the PVC.
* If the PV's `fsType` is defined and the PVC's `accessModes` list contains `ReadWriteOnce`, continue to the next step.
* Otherwise, the Kubelet doesn't attempt to change the ownership and permission for the volume.
1. Check `pod.spec.securityContext.fsGroupChangePolicy`.
* If the `pod.spec.securityContext.fsGroupChangePolicy` is set to `always` or empty, the kubelet performs the following actions:
  * Ensures that all processes of the containers inside the pod are part of the supplementary group id `pod.spec.securityContext.fsGroup`
  * Ensures that any new files created in the volume will be in group id `pod.spec.securityContext.fsGroup`
  * Recursively changes permission and ownership of the volume to have the same group id as `pod.spec.securityContext.fsGroup` every time the volume is mounted
* If the `pod.spec.securityContext.fsGroupChangePolicy` is set to `OnRootMismatch`:
  * If the root of the volume already has the correct permissions (i.e., belongs to the group id as `pod.spec.securityContext.fsGroup`) , the recursive permission and ownership change will be skipped.
  * Otherwise, Kubelet recursively changes permission and ownership of the volume to have the same group id as `pod.spec.securityContext.fsGroup`

For more information, see:
* https://kubernetes.io/docs/tasks/configure-pod-container/security-context/#configure-volume-permission-and-ownership-change-policy-for-pods
* https://github.com/longhorn/longhorn/issues/2131#issuecomment-778897129

### Longhorn PVC with Block Volume Mode

For PVC with Block volume mode, Kubelet never attempts to change the permission and ownership of the block device when making it available inside the container.
You must set the correct group ID in the `pod.spec.securityContext` for the pod to be able to read and write to the block device or run the container as root.

By default, Longhorn puts the block device into group id 6, which is typically associated with the "disk" group.
Therefore, pods that use Longhorn PVC with Block volume mode must either set the group id 6 in the `pod.spec.securityContext`, or run as root.
For example:
1. Pod that sets the group id 6 in the `pod.spec.securityContext`
    ```yaml
    apiVersion: v1
    kind: PersistentVolumeClaim
    metadata:
      name: longhorn-block-vol
    spec:
      accessModes:
        - ReadWriteOnce
      volumeMode: Block
      storageClassName: longhorn
      resources:
        requests:
          storage: 2Gi
    ---
    apiVersion: v1
    kind: Pod
    metadata:
      name: block-volume-test
      namespace: default
    spec:
      securityContext:
        runAsGroup: 1000
        runAsNonRoot: true
        runAsUser: 1000
        supplementalGroups:
        - 6
      containers:
        - name: block-volume-test
          image: ubuntu:20.04
          command: ["sleep", "360000"]
          imagePullPolicy: IfNotPresent
          volumeDevices:
            - devicePath: /dev/longhorn/testblk
              name: block-vol
      volumes:
        - name: block-vol
          persistentVolumeClaim:
            claimName: longhorn-block-vol
    ```
1. Pod that runs as root
    ```yaml
    apiVersion: v1
    kind: PersistentVolumeClaim
    metadata:
      name: longhorn-block-vol
    spec:
      accessModes:
        - ReadWriteOnce
      volumeMode: Block
      storageClassName: longhorn
      resources:
        requests:
          storage: 2Gi
    ---
    apiVersion: v1
    kind: Pod
    metadata:
      name: block-volume-test
      namespace: default
    spec:
      containers:
        - name: block-volume-test
          image: ubuntu:20.04
          command: ["sleep", "360000"]
          imagePullPolicy: IfNotPresent
          volumeDevices:
            - devicePath: /dev/longhorn/testblk
              name: block-vol
      volumes:
        - name: block-vol
          persistentVolumeClaim:
            claimName: longhorn-block-vol
    ```


---

## Article: nodes-and-volumes/volumes/rwx-volumes.md

---
title: ReadWriteMany (RWX) Volume
weight: 4
---

Longhorn supports ReadWriteMany (RWX) volumes by exposing regular Longhorn volumes via NFSv4 servers that reside in share-manager pods.


# Introduction

Longhorn creates a dedicated `share-manager-<volume-name>` Pod within the `longhorn-system` namespace for each RWX volume that is currently in active use. The Pod facilitates the export of Longhorn volume via an internally hosted NFSv4 server. Additionally, a corresponding Service is created for each RWX volume, serving as the designated endpoint for actual NFSv4 client connections.

{{< figure src="/img/diagrams/rwx/rwx-arch.png" >}}

# Requirements

It is necessary to meet the following requirements in order to use RWX volumes.

1. Each NFS client node needs to have a NFSv4 client installed.

    Please refer to [Installing NFSv4 client](../../../deploy/install/#installing-nfsv4-client) for more installation details.

      > **Troubleshooting:** If the NFSv4 client is not available on the node, when trying to mount the volume the below message will be part of the error:
      > ```
      > for several filesystems (e.g. nfs, cifs) you might need a /sbin/mount.<type> helper program.
      > ```

2. The hostname of each node is unique in the Kubernetes cluster.

    There is a dedicated recovery backend service for NFS servers in Longhorn system. When a client connects to an NFS server, the client's information, including its hostname, will be stored in the recovery backend. When a share-manager Pod or NFS server is abnormally terminated, Longhorn will create a new one. Within the 90-seconds grace period, clients will reclaim locks using the client information stored in the recovery backend.

# Creation and Usage of an RWX Volume

> **Notice**  
> An RWX volume must have the access mode set to `ReadWriteMany` and the "migratable" flag disabled (*parameters.migratable: `false`*).

1. For dynamically provisioned Longhorn volumes, the access mode is based on the PVC's access mode.
2. For manually created Longhorn volumes (restore, DR volume) the access mode can be specified during creation in the Longhorn UI.
3. When creating a PV/PVC for a Longhorn volume via the UI, the access mode of the PV/PVC will be based on the volume's access mode.
4. One can change the Longhorn volume's access mode via the UI as long as the volume is not bound to a PVC.
5. For a Longhorn volume that gets used by an RWX PVC, the volume access mode will be changed to RWX.

## Configuring Volume Locality for RWX Volumes

Longhorn provides new settings that allow you to precisely control the data locality of RWX volumes (through identification of associated Share Manager pods). These granular settings work with related global settings to provide optimal performance, resilience, and adherence to organizational policies or constraints.

### `shareManagerNodeSelector`

You can use the StorageClass parameter `shareManagerNodeSelector` to specify selectors for identifying nodes that RWX volumes can be scheduled on. These selectors are merged with global `system-managed-components-node-selector` settings and then applied to the Share Manager pods of the RWX volumes to provide more control over volume locality.

  Example:
  ```
  kind: StorageClass
  apiVersion: storage.k8s.io/v1
  metadata:
    name: longhorn-rwx
  provisioner: driver.longhorn.io
  parameters:
    shareManagerNodeSelector: label-key1:label-value1;label-key2:label-value2
  ```
  In this example, RWX volumes provisioned with the specified StorageClass will be scheduled on nodes with the labels `label-key1:label-value1` and `label-key2:label-value2`. 

### `allowedTopologies`

Longhorn converts the `storageClass.allowedTopologies` settings into affinity rules for the Share Manager pods of the RWX volumes. This ensures that the pods are scheduled on nodes that meet the specified topological requirements (such as regions and zones) and align with the RWX volume locality.

  Example:
  ```
  kind: StorageClass
  apiVersion: storage.k8s.io/v1
  metadata:
    name: longhorn-rwx
  provisioner: driver.longhorn.io
  allowedTopologies:
  - matchLabelExpressions:
    - key: topology.kubernetes.io/region
      values:
      - us-west-1
  ```
  In this example, the Share Manager pods and RWX volumes will be scheduled in the `us-west-1` region.

### `shareManagerTolerations`

You can also use the StorageClass parameter `shareManagerTolerations` to allow more flexible scheduling based on node taints. The defined tolerations are merged with global `taint-toleration` settings and then applied to the Share Manager pods.

  Example:
  ```
  kind: StorageClass
  apiVersion: storage.k8s.io/v1
  metadata:
    name: longhorn-rwx
  provisioner: driver.longhorn.io
  parameters:
    shareManagerTolerations: nodetype=storage:NoSchedule
  ```
  In this example, the Share Manager pods will tolerate the `nodetype=storage:NoSchedule` taint on nodes, allowing them to be scheduled on those nodes.

## Configuring Volume Mount Options

An RWX volume is accessible only when mounted via NFS. By default Longhorn uses NFS version 4.1 with the `softerr` mount option, a `timeo` value of "600", and a `retrans` value of "5".

If the NFS server becomes inaccessible, requests from NFS clients are retried according to the configured `retrans` value. Longer-duration events such as power outages and factors such as network partitions cause the requests to eventually fail. An NFS error (`ETIMEDOUT` for the `softerr` mount option) is returned to the calling application and data loss may occur. If `softerr` is not supported, Longhorn automatically uses the `soft` mount option instead, which returns an `EIO` as the error.

You can use specific mount options for new volumes. First, create a customized StorageClass with an `nfsOptions` parameter, and then create PVCs for RWX volumes using that specific StorageClass.

Example:

  ```yaml
  kind: StorageClass
  apiVersion: storage.k8s.io/v1
  metadata:
    name: longhorn-test
  provisioner: driver.longhorn.io
  allowVolumeExpansion: true
  reclaimPolicy: Delete
  volumeBindingMode: Immediate
  parameters:
    numberOfReplicas: "3"
    staleReplicaTimeout: "2880"
    fromBackup: ""
    fsType: "ext4"
    nfsOptions: "vers=4.2,noresvport,softerr,timeo=600,retrans=5"
  ```

> **Important:**
> To create PVCs for RWX volumes using the sample StorageClass, replace the `nfsOptions` string with a customized comma-separated list of legal options.

### Notes

1. You must provide the complete set of desired options. Any options not supplied will use the NFS-server side defaults, not Longhorn's own.

2. Longhorn does not validate the `nfsOptions` string, so erroneous values and typographical errors are not flagged. When the string is invalid, the mount is rejected by the NFS server and the volume is not created nor attached.

3. In Longhorn v1.4.0 to 1.4.3 and v1.5.0 to v1.5.1, volumes within a share manager pod (specifically, in the `NodeStageVolume` step) are hard mounted by default by the Longhorn CSI plugin. Hard mounting allows Longhorn to persistently retry sending NFS requests, ensuring that IOs do not fail even when the NFS server becomes inaccessible for some time. IOs resume seamlessly when the server regains connectivity or a replacement server is created.

   This mechanism for guaranteeing data integrity, however, comes with some risk. To maintain stability, the Linux kernel does not allow unmounting of a file system until all pending IOs are completed. This is a concern because the system cannot shut down until all file systems are unmounted. If the NFS server is unable to recover, the client nodes must undergo a forced reboot.

   To mitigate the issue, upgrade to v1.4.4, v1.5.2, or a later version. After upgrading, either `softerr` or `soft` is automatically applied to the `nfsOptions` parameter whenever RWX volumes are reattached (if the default settings are not overridden).

4. You can still use the `hard` mount option (via the `nfsOptions` override mechanism), but hard-mounted volumes are subject to the outlined risks.

For more information, see [#6655](https://github.com/longhorn/longhorn/issues/6655).

# Failure Handling

1. share-manager Pod is abnormally terminated

    Client IO will be blocked until Longhorn creates a new share-manager Pod and the associated volume. Once the Pod is successfully created, the 90-seconds grace period for lock reclamation is started, and users would expect
      - Before the grace period ends, client IO to the RWX volume will still be blocked.
      - The server rejects READ and WRITE operations and non-reclaim locking requests with an error of NFS4ERR_GRACE.
      - The grace period can be terminated early if all locks are successfully reclaimed.

    After exiting the grace period, IOs of the clients successfully reclaiming the locks continue without stale file handle errors or IO errors. If a lock cannot be reclaimed within the grace period, the lock is discarded, and the server returns IO error to the client. The client re-establishes a new lock. The application should handle the IO error. Nevertheless, not all applications can handle IO errors due to their implementation. Thus, it may result in the failure of the IO operation and the data loss. Data consistency may be an issue.

    Here is an example of a DaemonSet using an RWX volume.

    Each Pod of the DaemonSet is writing data to the RWX volume. If the node where the share-manager Pod is running is down, a new share-manager Pod is created on another node. Since one of the clients located on the down node has gone, the lock reclaim process cannot be terminated earlier than the 90-second grace period, even though the remaining clients' locks have been successfully reclaimed. The IOs of these clients continue after the grace period has expired.

2. If the Kubernetes DNS service goes down, share-manager Pods will not be able to communicate with longhorn-nfs-recovery-backend

    The NFS-ganesha server in a share-manager Pod communicates with longhorn-nfs-recovery-backend via the service `longhorn-recovery-backend`'s IP. If the DNS service is out of service, the creation and deletion of RWX volumes as well as the recovery of NFS servers will be inoperable. Thus, the high availability of the DNS service is recommended for avoiding the communication failure.

3. Fast failover feature.

    Longhorn supports a feature that can improve availability by shortening the time it takes to recover from a failure of the node on which the volume's share-manager NFS server pod is running.  The feature uses a direct heartbeat to monitor the server. If the server is unresponsive it acts to create a new one faster than the usual sequence. It also configures the NFS server differently, to shorten the recovery grace period from 90 to 30 seconds.  
    More details are at [RWX Volume Fast Failover](../../../high-availability/rwx-volume-fast-failover).


# Migration from Previous External Provisioner

The below PVC creates a Kubernetes job that can copy data from one volume to another.

- Replace the `data-source-pvc` with the name of the previous NFSv4 RWX PVC that was created by Kubernetes.
- Replace the `data-target-pvc` with the name of the new RWX PVC that you wish to use for your new workloads.

You can manually create a new RWX Longhorn volume + PVC/PV, or just create an RWX PVC and then have Longhorn dynamically provision a volume for you.

Both PVCs need to exist in the same namespace. If you were using a different namespace than the default, change the job's namespace below.

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  namespace: default  # namespace where the PVC's exist
  name: volume-migration
spec:
  completions: 1
  parallelism: 1
  backoffLimit: 3
  template:
    metadata:
      name: volume-migration
      labels:
        name: volume-migration
    spec:
      restartPolicy: Never
      containers:
        - name: volume-migration
          image: ubuntu:xenial
          tty: true
          command: [ "/bin/sh" ]
          args: [ "-c", "cp -r -v /mnt/old /mnt/new" ]
          volumeMounts:
            - name: old-vol
              mountPath: /mnt/old
            - name: new-vol
              mountPath: /mnt/new
      volumes:
        - name: old-vol
          persistentVolumeClaim:
            claimName: data-source-pvc # change to data source PVC
        - name: new-vol
          persistentVolumeClaim:
            claimName: data-target-pvc # change to data target PVC
```


# History
* Available since v1.0.1 [External provisioner](https://github.com/Longhorn/Longhorn/issues/1183)
* Available since v1.1.0 [Native RWX support](https://github.com/Longhorn/Longhorn/issues/1470)


---

## Article: nodes-and-volumes/volumes/trim-filesystem.md

---
title: Trim Filesystem
weight: 7
---

Since v1.4.0, Longhorn supports trimming filesystem inside Longhorn volumes. Trimming will reclaim space wasted by the removed files of the filesystem.

> **Note:**
> - Trimming removed files in snapshots has no effect on the filesystem because valid snapshots are immutable. However,
    the filesystem remembers whenever it has trimmed blocks associated with a snapshot. Because of this, you may need to
    unmount and remount the filesystem before reattempting to trim a snapshot that has been marked as removed.
>
> - If you allow automatic snapshot removal during filesystem trim, use the mount option `discard` with caution.
    `discard` frequently triggers snapshot removal and interrupts operations such as backup creation.

## Prerequisites

- The Longhorn version must be v1.4.0 or higher.
- There is a trimmable filesystem like EXT4 or XFS inside the Longhorn volume.
- The volume is attached and mounted on a mount point before trimming.

## Trim the filesystem in a Longhorn volume

You can trim a Longhorn volume using either the Longhorn UI or the `fstrim` command.

### Via Longhorn UI

You can directly click volume operation `Trim Filesystem` for attached volumes.

Then Longhorn will **try its best** to figure out the mount point and execute `fstrim <the mount point>`.  If something is wrong or the filesystem does not exist, the UI will return an error.

### Via shell command

When using `fstrim`, you must identify the mount point of the volume and then run the command `fstrim <the mount point>`.

- RWO volume: The mount point is either a pod of the workload or the node to which the volume was manually attached.
- RWX volume: The mount point is the share manager pod of the volume. The share manager pod contains the NFS server and is typically named `share-manager-<volume name>`.

To trim an RWX volume, perform the following steps:

1. Identify and then open a shell inside the share manager pod of the volume.
    ```
    kubectl -n longhorn-system exec -it <the share manager pod> -- bash
    ```
1. Identify the work directory of the NFS server (for example, `/export/<volume name>`).
    ```
    mount | grep <volume name>
    /dev/longhorn/<volume name> on /export/<volume name> type ext4 (rw,relatime)
    ```
1. Trim the work directory.
    ```
    fstrim /export/<volume name>
    ```

#### Periodically trim the filesystem

You can set up a [RecurringJob](../../../snapshots-and-backups/scheduling-backups-and-snapshots/#set-up-recurring-jobs) to periodically trim the filesystem.

## Automatically Remove Snapshots During Filesystem Trim

By design, valid snapshots of Longhorn volumes are immutable so you can only use the filesystem trim feature with the
following:

- Volume head
- Preceding continuous chain of snapshots created by the system or marked as removed

If most of the actual space consumed by a volume is associated with valid snapshots, the trim operation is not very
effective.

### Global Setting: "Remove Snapshots During Filesystem Trim"

If you want Longhorn to automatically reclaim the maximum amount of space, you can enable the setting
[_Remove Snapshots During Filesystem Trim_](../../../references/settings/#remove-snapshots-during-filesystem-trim).
When this global setting is enabled, the latest snapshot and the preceding continuous chain of snapshots are
automatically marked as removed, allowing Longhorn to reclaim space for as many snapshots as possible. However, the
setting can cause removal (and eventual purging) of snapshots that you intentionally created.

### The Volume Spec Field "UnmapMarkSnapChainRemoved"

There is a per-volume field `volume.Spec.UnmapMarkSnapChainRemoved` that overwrites the global setting mentioned above.

The options for this volume-specific setting are "disabled", "enabled", and "ignored". When the value is "ignored", the
global setting takes effect.

You can configure this setting in a StorageClass so that the value is applied to all volumes created using that
StorageClass.

## Known Issues & Limitations

### Rebuilding Volumes

By design, Longhorn unmaps blocks in the volume head and in the preceding continuous chain of snapshots marked as
removed. Some of these snapshots may be moved from one replica to another during volume rebuilding, so Longhorn is
unable to trim the filesystem of affected volumes when rebuilding is in progress.

Because rebuilding may take a long time, Longhorn simply does not unmap blocks during a rebuild instead of returning an
I/O error to the filesystem. This behavior benefits VM workloads in particular, which respond poorly when repeated
attempts to complete a trim return errors. See [Issue #7103](https://github.com/longhorn/longhorn/issues/7103) for more
information.

A trim operation that is started during rebuilding has no effect. Future trim operations on the same mounted volume may
also have no effect because the filesystem remembers which blocks it has trimmed. You may need to unmount and remount
the filesystem before attempting to start the trim operation again.

### Expanding Volumes

Longhorn is unable to trim the filesystem during volume expansion. Because expansion is fast, Longhorn returns an I/O
error whenever the issue is encountered. The filesystem recognizes that blocks were not trimmed and can try again
without a remount.

### Encrypted Volumes

- By default, TRIM commands are not enabled by the device-mapper. You can check [this doc](https://wiki.archlinux.org/title/Dm-crypt/Specialties#Discard/TRIM_support_for_solid_state_drives_(SSD)) for details.

- If you still want to trim an encrypted Longhorn volume, you can:
    1. Enter into the node host the volume is attached to.
    2. Enable flag `discards` for the encrypted volume. The passphrase is recorded in the corresponding secret:
    ```shell
    cryptsetup --allow-discards --persistent refresh <Longhorn volume name>
    ```
    3. Directly use Longhorn UI to trim the volume or execute `fstrim` for **the mount point** of `/dev/mapper/<volume name>` manually.


---

## Article: nodes-and-volumes/volumes/volume-conditions.md

---
title: Volume Conditions
weight: 7
---

## Volume Conditions 

Volume conditions describe the current status of a volume and potential issues that may occur.
- `Scheduled`: All replicas were scheduled successfully. 
    If Longhorn was unable to schedule any of the replicas, the condition is set to `false` and error messages are displayed. The condition is set to `true` when the setting [Allow Volume Creation With Degraded Availability](../../../references/settings#allow-volume-creation-with-degraded-availability) is enabled and at least one replica is scheduled.
- `TooManySnapshots`: This specific volume has more than 100 snapshots.
    Longhorn allows you to create a maximum of 250 snapshots for each volume. For more information about configuring the maximum snapshot count, see [Snapshot Space Management](../../../snapshots-and-backups/snapshot-space-management).
- `Restore`: Longhorn is restoring the volume from a backup.
- `WaitForBackingImage`: The replicas have not started because the backing images must first be synced with their disks.

## Engine Conditions

Engine conditions describe the current status of an engine and potential issues that may occur.
`FilesystemReadOnly`: The state of the current volume mount point has changed to read-only. 
This change may prevent workloads from writing data to the volume. For troubleshooting information, see [Volume Recovery](../../../high-availability/recover-volume).


## Replica Conditions

Replica conditions describe the current status of a replica and potential issues that may occur.
- `RebuildFailed`: The replica failed to rebuild.
- `WaitForBackingImage`: The replicas have not started because the backing images must first be synced with their disks.

---

## Article: nodes-and-volumes/volumes/volume-size.md

---
title: Volume Size
weight: 4
---

In this section, you'll have a better understanding of concepts related to volume size.

## Volume `Size`:

{{< figure src="/img/screenshots/volumes-and-nodes/volume-size-nominal-size.png" >}}

This value, which you specified during volume creation, represents the amount of space available to the volume when in use. 

The following are other ways of understanding this concept:

- The volume itself is just a Kubernetes CRD object and volume data is stored in replicas. This value represents the nominal size of each replica. 
- Longhorn replicas use [sparse files](https://wiki.archlinux.org/index.php/Sparse_file) to store data. This value represents the maximum size to which a sparse file may expand.

Replicas are scheduled on nodes with enough allocatable space to cover this nominal size during volume creation. For more information, see [Node Space Usage](../../nodes/node-space-usage).

> **Note**: The maximum volume size is based on the disk's file system (for example, 16383 GiB for `ext4`).

## Volume `Actual Size`

{{< figure src="/img/screenshots/volumes-and-nodes/volume-size-actual-size.png" >}}

This value represents the amount of space used by each replica (including the volume head and snapshots) on the node.

Because all historical data (stored in the snapshots) and active data are included in the calculation, this value can exceed the user-defined nominal size.

The Longhorn UI displays this value only when the volume is running.

## Example

In the example, we will explain how volume `size` and `actual size` get changed after a bunch of IO and snapshot related operations.

> The illustration presents the file organization of **one replica**. The volume head and snapshots are actually sparse files, which we mentioned above.

{{< figure src="/img/screenshots/volumes-and-nodes/volume-size-illustration.png" >}}


1. Create a 12 Gi volume with a single replica, then attach and mount it on a node. See Figure 1 of the illustration.
    - For the empty volume, the nominal `size` is 12 Gi and the `actual size` is almost 0.
    - There is some meta info in the volume hence the `actual size` is 260 Mi and is not exactly 0.

{{< figure src="/img/screenshots/volumes-and-nodes/volume-size-illustration-fig1.png" >}}

2. Write 4 Gi data (data#0) in the volume mount point. The `actual size` is increased by 4 Gi because of the allocated blocks in the replica for the 4 Gi data. Meanwhile, `df` command in the filesystem also shows the 4 Gi used space. See Figure 2 of the illustration.

{{< figure src="/img/screenshots/volumes-and-nodes/volume-size-illustration-fig2.png" >}}

3. Delete the 4 Gi data. Then, `df` command shows that the used space of the filesystem is nearly 0, but the `actual size` is unchanged.

    > Users can see by default the volume `actual size` is not shrunk after deleting the 4 Gi data. Longhorn is a block-level storage system. Therefore, the deletion in the filesystem only marks the blocks that belong to the deleted file as unused. Currently, Longhorn will not apply TRIM/UNMAP operations automatically/periodically. if you want to do filesystem trim, please check [this doc](../trim-filesystem) for details.

{{< figure src="/img/screenshots/volumes-and-nodes/volume-size-illustration-fig2.png" >}}

4. Then, rewrite the 4 Gi data (data#1), and the `df` command in the filesystem shows 4 Gi used space again. However, the `actual size` is increased by 4 Gi and becomes 8.25Gi. See Figure 3(a) of the illustration.

     > After deletion, filesystem may or maynot reuse the recently freed blocks from recently deleted files according to the filesystem design and please refer to [Block allocation strategies of various filesystems](https://www.ogris.de/blkalloc). If the volume nominal `size` is 12 Gi, the `actual size` in the end would range from 4 Gi to 8 Gi since the filesystem may or maynot reuse the freed blocks. On the other hand, if the volume nominal `size` is 6 Gi, the `actual size` at the end would range from 4 Gi to 6 Gi, because the filesystem has to reuse the freed blocks in the 2nd round of writing. See Figure 3(b) of the illustration.
     >
     > Thus, allocating an appropriate nominal `size` for a volume that holds heavy writing tasks according to the IO pattern would make disk space usage more efficient.

{{< figure src="/img/screenshots/volumes-and-nodes/volume-size-illustration-fig3.png" >}}

5. Take a snapshot (snapshot#1). See Figure 4 of the illustration.
    - Now data#1 is stored in snapshot#1.
    - The new volume head size is almost 0.
    - With the volume head and the snapshot included, the `actual size` remains 8.25 Gi.

{{< figure src="/img/screenshots/volumes-and-nodes/volume-size-illustration-fig4.png" >}}

6. Delete data#1 from the mount point.
    - The data#1 filesystem level removal info is stored in current volume head file. For snapshot#1, data#1 is still retained as the historical data.
    - The `actual size` is still 8.25 Gi.

7. Write 8 Gi data (data#2) in the volume mount, then take one more snapshot (snapshot#2). See Figure 5 of the illustration.
    - Now the `actual size` is 16.2 Gi, which is greater than the volume nominal `size`.
    - From a filesystem's perspective, the overlapping part between the two snapshots is considered as the blocks that have to be reused or overwritten. But in terms of Longhorn, these blocks are actually fresh ones held in another snapshot/volume head. See the 2 snapshots in Figure 6.

    > The volume head holds the latest data of the volume only, while each snapshot may store historical data as well as active data, which consumes at most size space. Therefore, the volume `actual size`, which is the size sum of the volume head and all snapshots, is possibly bigger than the size specified by users.
    >
    > Even if users will not take snapshots for volumes, there are operations like rebuilding, expansion, or backing up that would lead to system (hidden) snapshot creation. As a result, volume `actual size` being larger than size is unavoidable under some use cases.

{{< figure src="/img/screenshots/volumes-and-nodes/volume-size-illustration-fig5.png" >}}

8. Delete snapshot#1 and wait for snapshot purge complete. See Figure 7 of the illustration.
    - Here Longhorn actually coalesces the snapshot#1 with the snapshot#2.
    - For the overlapping part during coalescing, the newer data (data#2) will be retained in the blocks. Then some historical data is removed and the volume gets shrunk (from 16.2 Gi to 11.4 Gi in the example).

{{< figure src="/img/screenshots/volumes-and-nodes/volume-size-illustration-fig6.png" >}}

9. Delete all existing data (data#2) and write 11.5 Gi data (data#3) in the volume mount. See Figure 8 of the illustration.
    - this makes the volume head actual size becomes 11.5 Gi and the volume total actual size becomes 22.9 Gi.

{{< figure src="/img/screenshots/volumes-and-nodes/volume-size-illustration-fig7.png" >}}

10. Try to delete the only snapshot (snapshot#2) of the volume. See Figure 9 of the illustration.
    - The snapshot directly behinds the volume head cannot be cleaned up.
      If users try to delete this kind of snapshot, Longhorn will mark the snapshot as Removing, hide it, then try to free the overlapping part between the volume head and the snapshot for the snapshot file.
      The last operation is called snapshot prune in Longhorn and is available since v1.3.0.
    - Since in the example both the snapshot and the volume head use up most of the nominal space, the overlapping part almost equals to the snapshot actual size. After the pruning, the snapshot actual size is down to 259 Mi and the volume gets shrunk from 22.9 Gi to 11.8 Gi.

{{< figure src="/img/screenshots/volumes-and-nodes/volume-size-illustration-fig8.png" >}}


Here we summarize the important things related to disk space usage we have in the example:

- Unused blocks are not released

  Longhorn will not issue TRIM/UNMAP operations automatically. Hence deleting files from filesystems will not lead to volume actual size decreasing/shrinking. You may need to check [the doc](../trim-filesystem) and handle it by yourself if needed.


- Allocated blocks but unused are not reused

  Deleting then writing new files would lead to the actual size keeps increasing. Since the filesystem may not reuse the recently freed blocks from recently deleted files. Thus, allocating an appropriate nominal size for a volume that holds heavy writing tasks according to the IO pattern would make disk space usage more efficient.

- By deleting snapshots, the overlapping part of the used blocks might be eliminated regardless of whether the blocks are recently released blocks by the filesystem or still contain historical data.

## Space Configuration Suggestions for Volumes

1. Reserve enough free space in disks as buffers in case of the actual size of existing volumes keep growing up.
    - A general estimation for the maximum space consumption of a volume is

        ```
        (N + 1) x head/snapshot average actual size
        ```

        - where `N` is the total number of snapshots the volume contains (including the volume head), and the extra `1` is for the temporary space that may be required by snapshot deletion.
        - The average actual size of the snapshots varies and depends on the use cases.
          If snapshots are created periodically for a volume (e.g. by relying on snapshot recurring jobs), the average value would be the average modified data size for the volume in the snapshot creation interval.
          If there are heavy writing tasks for volumes, the head/snapshot average actual size would be volume the nominal size. In this case, it's better to set [`Storage Over Provisioning Percentage`](../../../references/settings/#storage-over-provisioning-percentage) to be smaller than 100% to avoid disk space exhaustion.
        - Some extended cases:
            - There is one snapshot recurring job with retention number is `N`. Then the formula can be extended to:

                ```
                (M + N + 1 + 1 + 1 + 1) x head/snapshot average actual size
                ```

                - The explanation of the formula:
                    - `M` is the snapshots created by users manually. Recurring jobs are not responsible for removing this kind of snapshot. They can be deleted by users only.
                    - `N` is the snapshot recurring job retain number.
                    - The 1st `1` means the volume head.
                    - The 2nd `1` means the extra snapshot created by the recurring job. Since the recurring job always creates a new snapshot then deletes the oldest snapshot when the current snapshots created by itself exceeds the retention number. Before the deletion starts, there is one extra snapshot that can take extra disk space.
                    - The 3rd `1` is the system snapshot. If the rebuilding is triggered or the expansion is issued, Longhorn will create a system snapshot before starting the operations. And this system snapshot may not be able to get cleaned up immediately.
                    - The 4th `1` is for the temporary space that may be required by snapshot deletion/purge.
            - Users don't want snapshot at all. Neither (manually created) snapshot nor recurring job will be launched. Assume [setting _Automatically Cleanup System Generated Snapshot_](../../../references/settings/#automatically-clean-up-system-generated-snapshot) is enabled, then formula would become:

                ```
                (1 + 1 + 1) x head/snapshot average actual size
                ```

                - The worst case that leads to so much space usage:
                    1. At some point the 1st rebuilding/expansion is triggered, which leads to the 1st system snapshot creation.
                        - The purges before and after the 1st rebuilding does nothing.
                    2. There is data written to the new volume head, and the 2nd rebuilding/expansion somehow is triggered.
                        - The snapshot purge before the 2nd rebuilding may lead to the shrink of the 1st system snapshot.
                        - Then the 2nd system snapshot is created and the rebuilding is started.
                        - After the rebuilding done, the subsequent snapshot purge would lead to the coalescing of the 2 system snapshots. This coalescing requires temporary space.
                    3. During the afterward snapshot purging for the 2nd rebuilding, there is more data written to the new volume head.
                - The explanation of the formula:
                    - The 1st `1` means the volume head.
                    - The 2nd `1` is the second system snapshot mentioned in the worst case.
                    - The 3rd `1` is for the temporary space that may be required by the 2 system snapshot purge/coalescing.

2. Do not retain too many snapshots for the volumes.

3. Cleaning up snapshots will help reclaim disk space. There are two ways to clean up snapshots:
    - Delete the snapshots manually via Longhorn UI.
    - Set a snapshot recurring job with retention 1, then the snapshots will be cleaned up automatically.

    Also, notice that the extra space, up to volume nominal `size`, is required during snapshot cleanup and merge.

4. An appropriate the volume nominal `size` according to the workloads.

## Troubleshooting

### Workload is Hitting "no space left on device" Error

To address the "no space left on device" error, you need to understand that the difference between a volume being full and the node disks being full is crucial for proper storage management in Longhorn.

#### Volume Full

A volume is full when the filesystem mounted on it has reached its capacity limit.

- Data writes fail with "no space left on device" errors.
- The host disk for the replica of the volume may still have available space for other volumes.

- **Characteristics:**
  - The `df` command shows 100% usage for the mounted filesystem.
  - Applications cannot write new data to the volume mount point.

- **Example**

    ```bash
    $ df -h /mnt/longhorn-volume-example-dir
    Filesystem                    Size  Used Avail Use% Mounted on
    /dev/longhorn-volume-example   12G   12G    0  100% /mnt/longhorn-volume-example-dir
    ```

#### Node Disk Full

Node disks do not have enough space to accommodate volume operations because volumes are thin-provisioned and the replica sizes on the disks continue to grow.

- Data writes fail with "no space left on device" errors.
- Volume operations, such as volume creation/expansion, snapshot creation/deletion, and replica rebuilding, may be limited.
- Replicas of newly created volumes are prevented from being scheduled to the full disks.

- **Characteristics:**
  - Applications cannot write new data to the volume mount point even though the volume is not full.
  - Longhorn operations for volumes on these node disks, such as replica creation, rebuilding, or snapshot operations, may fail.
  - Volumes with replicas sharing the same node disks are affected.

##### Data Integrity Protection for Disks Out-of-Space Scenarios

When multiple replicas simultaneously encounter "no space left on device" errors, Longhorn implements data integrity protection. If all writable replicas are on disks that are out of space, the system preserves the maximum number of replicas that have written the same amount of data (the highest written bytes count). This ensures data consistency by:

1. **Retaining** the maximum number of replicas that have written the same amount of data.  
This avoids massive replica rebuilding operations, allows for efficient rebuilding when node disks recover from out-of-space issues and ensures that users can still read data consistently.
2. **Marking** other replicas as **ERR** to prevent data corruption.  
This ensures that users can read data consistently from retained replicas, and replicas marked as `ERR` will be rebuilt from retained replicas to maintain data integrity when node disks recover from out-of-space issues.
3. **Maintaining** volume accessibility even when all underlying disks are full.

For example, if replicas A, B, and C encounter space errors after writing 1MB, 2MB, and 2MB respectively, replicas B and C (having written the most data consistently) will remain active while replica A is marked as **ERR**.

### Solutions

When encountering the `no space left on device` error, first check whether the volume’s filesystem is full, and then check whether the disk hosting the volume replica is full.

- **Check volume filesystem usage**

    If the volume's filesystem usage is 100% using the following command:

    ```bash
    df -h /mnt/your-volume
    ```

    Please expand the volume or remove the unnecessary files in the filesystem first.

- **Check Longhorn node disk usage**

    If the disk usage is 100% using the following command or checking the Longhorn UI page `Node > Node Disk`:

    ```bash
    # Check through Longhorn UI: Node > Node Disk
    # Or check the data path mount point
    df -h /var/lib/longhorn # default data path
    ```

    Please follow the steps to recovery the workloads:

    1. Scale down the workload.
    2. Expand the disk size, or remove unnecessary files, for example, [orphaned replica directories](../../../advanced-resources/data-cleanup/orphaned-data-cleanup#orphaned-replica-directories) and [used backing images](../../../advanced-resources/backing-image/backing-image#clean-up-backing-images), on the disk.
    3. Then scale up the workload.


---

## Article: nodes-and-volumes/volumes/workload-identification.md

---
title: Viewing Workloads that Use a Volume
weight: 5
---

Now users can identify current workloads or workload history for existing Longhorn persistent volumes (PVs) and their history of being bound to persistent volume claims (PVCs).

From the Longhorn UI, go to the **Volumes** tab. Each Longhorn volume is listed on the page. The **Attached To** column displays the name of the workload using the volume. If you click the workload name, you will be able to see more details, including the workload type, pod name, and status.

Workload information is also available on the Longhorn volume detail page. To see the details, click the volume name:

```
State: attached
...
Namespace:default
PVC Name:longhorn-volv-pvc
PV Name:pvc-0edf00f3-1d67-4783-bbce-27d4458f6db7
PV Status:Bound
Pod Name:teststatefulset-0
Pod Status:Running
Workload Name:teststatefulset
Workload Type:StatefulSet
```

## History

After the workload is no longer using the Longhorn volume, the volume detail page shows the historical status of the most recent workload that used the volume:

```
Last time used by Pod: a few seconds ago
...
Last Pod Name: teststatefulset-0
Last Workload Name: teststatefulset
Last Workload Type: Statefulset
```

If these fields are set, they indicate that currently no workload is using this volume.

When a PVC is no longer bound to the volume, the following status is shown:

```
Last time bound with PVC:a few seconds ago
Last time used by Pod:32 minutes ago
Last Namespace:default
Last Bounded PVC Name:longhorn-volv-pvc
```

If the `Last time bound with PVC` field is set, it indicates currently there is no bound PVC for this volume. The related fields will show the most recent workload using this volume.


---

## Article: advanced-resources/_index.md

---
title: Advanced Resources
weight: 8
---

---

## Article: advanced-resources/volumeattachment.md

---
title: Longhorn VolumeAttachment
weight: 1
---

**Table of Contents**
- [Kubernetes and Longhorn `VolumeAttachment`](#kubernetes-and-longhorn-volumeattachment)
- [Longhorn VolumeAttachment](#longhorn-volumeattachment)
  - [VolumeAttachment CR](#volumeattachment-cr)
  - [Understanding Attachment Tickets](#understanding-attachment-tickets)
- [The CSI Attachment and Detachment Workflow](#the-csi-attachment-and-detachment-workflow)
  - [Core Components in the Workflow](#core-components-in-the-workflow)
  - [The CSI Volume Attachment Flow](#the-csi-volume-attachment-flow)
  - [The CSI Volume Detachment Flow](#the-csi-volume-detachment-flow)
  - [Summary of the Workflow](#summary-of-the-workflow)
- [Troubleshooting Volume Attachment Issues](#troubleshooting-volume-attachment-issues)
  - [Volume is Stuck in `Attaching` or `Detaching` State](#volume-is-stuck-in-attaching-or-detaching-state)
    - [Possible Causes](#possible-causes)
    - [Resolution Steps](#resolution-steps)
  - [Case Study](#case-study)
    - [Case 1: Failure to Attach Volume Due to Unexpected `longhorn-ui` Attachment Ticket](#case-1-failure-to-attach-volume-due-to-unexpected-longhorn-ui-attachment-ticket)
    - [Case 2: Volume Fails to Attach to New Node Due to Backup Job Stuck in Pending State](#case-2-volume-fails-to-attach-to-new-node-due-to-backup-job-stuck-in-pending-state)

---

This document provides a detailed overview of Longhorn's `VolumeAttachment` custom resource, how it integrates with Kubernetes' native `VolumeAttachment`, its operational flow, and common troubleshooting scenarios.

## Kubernetes and Longhorn `VolumeAttachment`

To understand how volume attachments work in Longhorn, it is important to distinguish between the Kubernetes `VolumeAttachment` and Longhorn's custom `VolumeAttachment`:

- **Kubernetes `VolumeAttachment`**: It is a native Kubernetes API resource that is part of the Container Storage Interface (CSI) specification. Its primary role is to signal to a CSI driver that a volume should be attached to a specific node.

- **Longhorn `VolumeAttachment`**: It is a Custom Resource (CR) defined by Longhorn, with the full name `volumeattachment.longhorn.io`. This internal Longhorn resource is used by the Longhorn Manager to track and manage all attachment requests for a volume.

## Longhorn VolumeAttachment

### VolumeAttachment CR

To retrieve a Longhorn `VolumeAttachment`, use the following command:

```bash
kubectl -n longhorn-system get volumeattachment.longhorn.io <volume-name> -o yaml
```

Example output:

```bash
apiVersion: v1
...
  spec:
    attachmentTickets:
      csi-f0471a334f0b249f964cd1dec461a5eb94c8d268cbbc904c1a8e9a37e2045208:
        generation: 0
        id: csi-f0471a334f0b249f964cd1dec461a5eb94c8d268cbbc904c1a8e9a37e2045208
        nodeID: rancher60-master
        parameters:
          disableFrontend: "false"
          lastAttachedBy: ""
        type: csi-attacher
    volume: pvc-b26e9514-aafd-46e0-b70c-4e3f187c7977
  status:
    attachmentTicketStatuses:
      csi-f0471a334f0b249f964cd1dec461a5eb94c8d268cbbc904c1a8e9a37e2045208:
        conditions:
        - lastProbeTime: ""
          lastTransitionTime: "2025-07-05T09:17:27Z"
          message: ""
          reason: ""
          status: "True"
          type: Satisfied
        generation: 0
        id: csi-f0471a334f0b249f964cd1dec461a5eb94c8d268cbbc904c1a8e9a37e2045208
        satisfied: true
...
```

- `spec.attachmentTickets`: A map containing all active attachment requests (also known as **tickets**). Each ticket includes:
  - `id`: A unique identifier for each attachment ticket.
  - `nodeID`: The ID of the node where the volume should be attached.
  - `parameters`: Optional parameters for the attachment, such as `disableFrontend` and `lastAttachedBy`.
  - `type`: The attacher type, indicating the source of the attachment request.

- `status.attachmentTicketStatuses`: A map containing the current status of each active attachment ticket or request. Each entry includes:
  - `conditions`: The current condition(s) of the ticket, including whether the request is satisfied or not.
  - `satisfied`: A boolean value indicating whether the attachment request has been fulfilled or not.
  - `generation`: The generation number of the ticket, used to track updates.

### Understanding Attachment Tickets

The Longhorn `VolumeAttachment` custom resource (CR) is designed to manage attachment requests from various internal Longhorn system controllers. Each request is represented as an **attachment ticket** within the CR.

All active tickets are stored in the `spec.attachmentTickets` map. The `type` field in each ticket (referred to as the **AttacherType**) identifies the source of the request. The common `AttacherType` values include:

- `csi-attacher`: The most common type which handles standard attachment requests from the Kubernetes CSI plugin, typically when mounting a volume to a pod.
- `longhorn-api`: Represents a manual attachment request initiated by a user, either through the Longhorn UI or the Longhorn API.
- `snapshot-controller`: Used when attaching a volume to create or restore a snapshot.
- `backup-controller`: Used when attaching a volume to perform a backup.
- `volume-restore-controller`: Used when attaching a volume during a restore operation.
- `volume-clone-controller`: Used when attaching a volume for cloning from an existing volume.
- `share-manager-controller`: Manages backend volume attachments for ReadWriteMany (RWX) volumes by attaching them to the share-manager pod.
- `volume-expansion-controller`: Handles attachments needed for online volume expansion.
- `volume-rebuilding-controller`: Used when attaching a volume to rebuild a degraded or missing replica.
- `salvage-controller`: Used during the salvage process when Longhorn attempts to recover and reattach a problematic volume.
- `volume-eviction-controller`: Handles attachments involved in evicting a replica from a node.
- `bim-ds-controller`: Used by the Backing Image Data Source controller when creating a volume from a backing image.

## The CSI Attachment and Detachment Workflow

To understand how Longhorn integrates with Kubernetes, it is important to examine how the native Kubernetes `VolumeAttachment` resource and Longhorn’s custom `VolumeAttachment` CR interact through the CSI interface.

### Core Components in the Workflow

In addition to the Kubernetes and Longhorn `VolumeAttachment` objects, several key components work together to manage volume attachment and detachment:

- `external-attacher`: A CSI sidecar container that monitors Kubernetes `VolumeAttachment` objects and triggers `ControllerPublishVolume` or `ControllerUnpublishVolume` gRPC calls to the Longhorn CSI driver.
- `longhorn-csi-plugin`: The Longhorn CSI driver that implements the required CSI gRPC services.
- `longhorn-manager`: The central controller in Longhorn that manages the full lifecycle of Longhorn volumes. It includes various sub-controllers, including the volume attachment logic.
- `longhorn-volume-attachment-controller`: A sub-controller within `longhorn-manager` that monitors the Longhorn `VolumeAttachment` CR and performs attach or detach operations based on its spec.

### The CSI Volume Attachment Flow

When a pod that uses a Longhorn PersistentVolumeClaim (PVC) is scheduled onto a node, the CSI volume attachment workflow is triggered.

1. **kubelet Request**: The kubelet on the target node detects that a Longhorn volume needs to be mounted and notifies the Kubernetes `attach-detach-controller`.

2. **Kubernetes `VolumeAttachment` Creation**: The `attach-detach-controller` creates a Kubernetes `VolumeAttachment` object, specifying the Longhorn CSI driver (`driver.longhorn.io`), the target node name, and the persistent volume name.

3. **`external-attacher` Triggers CSI Call**: The `external-attacher` sidecar container observes the new Kubernetes `VolumeAttachment` object and issues a `ControllerPublishVolume` gRPC call to the `longhorn-csi-plugin`.

4. **Longhorn `VolumeAttachment` CR Creation**: Rather than attaching the volume directly, the `longhorn-csi-plugin` creates a Longhorn `VolumeAttachment` custom resource (CR). It adds an **attachment ticket** to the CR’s spec to represent the attachment request.

5. **Longhorn Controller Reconciliation**: The `longhorn-volume-attachment-controller`, a sub-controller within `longhorn-manager`, detects the new ticket and begins reconciliation. It verifies that the volume is available and updates the corresponding Volume CR’s `spec.nodeID` with the target node name.

6. **`longhorn-manager` Executes Attachment**: After detecting that `spec.nodeID` is set, `longhorn-manager` starts the Longhorn Engine on the specified node to complete the attachment.

7. **Volume Attachment Completion**:
   - `longhorn-manager` updates the Volume CR’s status to reflect that the volume is attached.
   - The `longhorn-volume-attachment-controller` updates the Longhorn `VolumeAttachment` CR’s status to indicate success.
   - The `longhorn-csi-plugin` receives the successful status and responds to the `external-attacher`.
   - Finally, the `external-attacher` marks the Kubernetes `VolumeAttachment` object’s `status.attached` field as `true`.

8. **kubelet Mounts the Volume**: Once the volume is marked as attached, the kubelet proceeds with the `NodeStageVolume` and `NodePublishVolume` CSI calls to mount the volume into the pod’s container.

### The CSI Volume Detachment Flow

When a pod using a Longhorn volume is deleted or rescheduled, the CSI detachment workflow is triggered.

1. **kubelet Request**: The kubelet signals to the Kubernetes `attach-detach-controller` that the volume is no longer needed on the node.

2. **Kubernetes `VolumeAttachment` Deletion**: The `attach-detach-controller` deletes the corresponding Kubernetes `VolumeAttachment` object.

3. **`external-attacher` Triggers CSI Call**: The `external-attacher` observes the deletion and initiates a `ControllerUnpublishVolume` gRPC call to the `longhorn-csi-plugin`.

4. **Attachment Ticket Removal**: The `longhorn-csi-plugin` processes the request by updating the Longhorn `VolumeAttachment` CR to remove the relevant attachment ticket.

5. **Longhorn Controller Reconciliation**: The `longhorn-volume-attachment-controller` detects that the ticket has been removed. If no other tickets exist for the volume, it clears the `spec.nodeID` field in the Longhorn Volume CR.

6. **`longhorn-manager` Executes Detachment**: With the `spec.nodeID` cleared, `longhorn-manager` initiates the detachment process by stopping the Longhorn Engine on the node.

7. **Volume Detachment Completion**:
   - `longhorn-manager` updates the Volume CR’s status to indicate that the volume is detached.
   - The `longhorn-csi-plugin` receives confirmation and responds with success to the `external-attacher`.
   - The `external-attacher` removes the finalizer from the Kubernetes `VolumeAttachment` object, allowing the API server to fully delete it.

### Summary of the Workflow

Longhorn extends the native volume attachment mechanism of Kubernetes by introducing a custom `VolumeAttachment` CR. This design provides several advantages:

- **Decoupling and Abstraction**: The custom resource encapsulates complex attach or detach logic within Longhorn, reducing the responsibilities of the `longhorn-csi-plugin`.
- **Fine-Grained Control**: The attachment ticket system enables Longhorn to handle requests from multiple sources (for example, pods, snapshots, backups) while ensuring only one valid attachment per volume at any time.
- **Observability and Troubleshooting**: The CR gives clear visibility into attachment state and history of the volume, simplifying monitoring and troubleshooting.

In summary, the Kubernetes `VolumeAttachment` object initiates the attachment or detachment process, while Longhorn’s custom `VolumeAttachment` CR orchestrates and manages the actual operations internally.

## Troubleshooting Volume Attachment Issues

This section outlines common issues related to volume attachment and provides recommended resolution steps. Before making any changes, carefully inspect system logs and the relevant custom resources to avoid disrupting active workloads.

### Volume is Stuck in `Attaching` or `Detaching` State

When a volume remains in the `Attaching` or `Detaching` state for an extended period, the cause is often related to stale or conflicting attachment tickets in the Longhorn `VolumeAttachment` CR.

#### Possible Causes

- **Stale or Orphaned Tickets**: A ticket from a previous workload (for example, a deleted pod or completed backup job) was not properly removed and still exists under `spec.attachmentTickets`.

- **Conflicting Tickets**: An existing ticket (for example, from the CSI attacher) blocks a new request (for example, a manual detach or move to a different node).

#### Resolution Steps

1. **Inspect the Longhorn `VolumeAttachment` CR**: Use the following command to view the attachment tickets:

    ```bash
    kubectl -n longhorn-system get volumeattachment.longhorn.io <volume-name> -o yaml
    ```

2. **Analyze Ticket Sources**: Look under `spec.attachmentTickets` and check the `type` field for each ticket to identify its source (for example, `csi-attacher`, `backup-controller`, etc.).

3. **Remove Invalid Tickets with Caution**: If you confirm a ticket is no longer needed (for example, its corresponding pod has been deleted), you may remove it by editing the CR.

  > **Warning:**
  >
  > Deleting an active ticket can cause serious disruptions. If you remove a ticket still required by a running workload, Longhorn interprets this as a detach request:
  >
  > - The volume engine will stop on the node, causing the pod to lose storage access and encounter I/O errors, likely crashing the pod.
  > - Kubernetes CSI will eventually detect the issue and re-attach the volume, recreating the ticket, but this causes downtime and may require manual pod restart.
  >
  > Always verify that the workload related to the ticket is inactive before removing it.

4. **Verify the State**: After removing invalid tickets, Longhorn should be able to complete the attach or detach operation successfully.

### Case Study

#### Case 1: Failure to Attach Volume Due to Unexpected `longhorn-ui` Attachment Ticket

- **Issue**: Longhorn [#8339](https://github.com/longhorn/longhorn/issues/8339)  
- **Symptom**:  
  - Workloads using the affected volume remain stuck in `Pending` state.  
  - The Longhorn `VolumeAttachment` CR contains an unexpected ticket from `longhorn-ui`.  
- **Workaround**:  
  - Inspect the `VolumeAttachment` CR:  
    ```bash
    kubectl -n longhorn-system get volumeattachment.longhorn.io <volume-name> -o yaml
    ```  
  - If you find a `longhorn-ui` attachment ticket, remove the entire ticket block from the CR.

#### Case 2: Volume Fails to Attach to New Node Due to Backup Job Stuck in Pending State

- **Issue**: Longhorn [#10090](https://github.com/longhorn/longhorn/issues/10090)  
- **Symptom**:  
  - When a workload is rescheduled to a different node, the volume fails to follow.  
  - Backup jobs referencing non-existent snapshots remain stuck in `Pending` state, with `status.message` containing `failed to get the snapshot ... not found`.  
  - These stuck backup jobs hold onto the volume, blocking detach/reattach.  
- **Workaround**:  
  1. Check the Longhorn `VolumeAttachment` CR for any tickets locking the volume:  
      ```bash
      kubectl -n longhorn-system get volumeattachment.longhorn.io <volume-name> -o yaml
      ```  
  2. If you see a ticket from the backup controller, a backup job is locking the volume.  
  3. **Do not delete the `backup-*` attachment ticket directly**, as Longhorn will recreate it.  
  4. Instead, resolve the stuck backup job by removing any `Backup` CRs with:  
       - `status.state = pending`  
       - `status.message` containing `Failed to get the Snapshot...`

      This releases the volume and allows it to be reattached.


---

## Article: advanced-resources/os-distro-specific/_index.md

---
title: OS/Distro Specific Configuration
weight: 2
---



---

## Article: advanced-resources/os-distro-specific/container-optimized-os-support.md

---
title:  Container-Optimized OS (COS) Support
weight: 5
---

## Requirements

> **Note:**
> Longhorn currently supports Container-Optimized OS only when used as the base image for Google Kubernetes Engine (GKE), which includes a pre-configured Kubernetes environment. The following information may not apply to manually created Kubernetes environments, including Kubernetes provisioned with other orchestrators.

The [Container-Optimized OS (COS)](https://cloud.google.com/container-optimized-os/docs) does not include a package manager and does not allow non-containerized applications to run. Additionally, its root filesystem is mounted as read-only, which poses a challenge for IO operations.

In GKE, Kubernetes tackles these constraints by housing necessary dependencies in a chroot environment (`/home/kubernetes/containerized_mounter/rootfs`) and mounting directories within it, enabling the execution of required tasks.

Longhorn provides a GKE COS node agent daemonset, which leverages GKE Kubernetes solutions to configure and run necessary dependencies. This agent is responsible for the following operations:

- Mounting the Longhorn data path.
- Loading the kernel module.
- Installing and running the iSCSI daemon.

## GKE COS Node Agent Installation
1. Configure the Longhorn GKE COS node agent. You can use the default settings, if applicable.
    > **Tip:**
    > You can use a comma-separated list when specifying values for the `node-agent` container's environment variable (`LONGHORN_DATA_PATHS`).
    >
    > Example:
    >
    > ```yaml
    >  containers:
    >    - name: node-agent
    >      env:
    >        - name: LONGHORN_DATA_PATHS
    >          value: /var/lib/longhorn1,/var/lib/longhorn2


1. Install the Longhorn GKE COS node agent.
    ```
    kubectl apply -f https://raw.githubusercontent.com/longhorn/longhorn/v{{< current-version >}}/deploy/prerequisite/longhorn-gke-cos-node-agent.yaml
    ```

1. Check the agent pod's status.
    Example:
    ```
    $ kubectl -n longhorn-system get pod -l app=longhorn-gke-cos-node
    NAME                                READY   STATUS    RESTARTS      AGE
    longhorn-gke-cos-node-agent-222w8   1/1     Running   1 (86m ago)   86m
    longhorn-gke-cos-node-agent-8r26h   1/1     Running   1 (86m ago)   86m
    longhorn-gke-cos-node-agent-nwhsw   1/1     Running   1 (86m ago)   86m
    ```

1. Check the installation result in the agent pod logs.
    ```
    Completed!
    Keep the container running for iscsi daemon
    ```
    > **Note:**
    > The agent installs the iSCSI daemon (iscsid) in a container using a package manager. However, the package manager attempts to initiate iSCSI services through systemd, which the container environment does not fully support. As a result, you will likely see error logs similar to `System has not been booted with systemd as init system (PID 1). Can't operate`. To work around this, the script manually starts the daemon instead of relying on systemd. You can disregard the mentioned errors in this context.

1. Verify that the dependent kernel module is loaded. You must run the command on the host.
    ```
    $ lsmod | grep -q iscsi_tcp && echo "The iSCSI module is loaded" || echo "The iSCSI module is NOT loaded"
    The iSCSI module is loaded
    ```

1. Verify that the iSCSI daemon is running. You must run the command on the host.
    ```
    $ ps aux | grep -q '[i]scsid' && echo "The iSCSI daemon is running" || echo "The iSCSI daemon is NOT running"
    The iSCSI daemon is running
    ```

1. Verify that the Longhorn data path (`/var/lib/longhorn`) is mounted on the host. If you specified multiple Longhorn data paths, run the command for each path on the host.
    ```
    $ findmnt --noheadings "/var/lib/longhorn"
    /var/lib/longhorn /dev/sda1[/var/lib/longhorn] ext4   rw,relatime,commit=30
    ```

## Limitations

- In COS clusters, Longhorn currently supports only V1 data volumes.
- You can use `pbkdf2` for volume encryption if the built-in `cryptsetup` utility in your COS cluster does not support `argon2i` or `argon2id`. For more information, see [Issue #10049](https://github.com/longhorn/longhorn/issues/10049).

## References

- [[FEATURE] Container-Optimized OS support](https://github.com/longhorn/longhorn/issues/6165)


---

## Article: advanced-resources/os-distro-specific/csi-on-gke.md

---
title:  Longhorn CSI on GKE
weight: 3
---

To operate Longhorn on a cluster provisioned with Google Kubernetes Engine, some additional configuration is required.

1. GKE clusters must use the `Ubuntu` OS instead of `Container-Optimized` OS, in order to satisfy Longhorn's `open-iscsi` dependency.

2. GKE requires a user to manually claim themselves as cluster admin to enable role-based access control. Before installing Longhorn, run the following command:

    ```shell
    kubectl create clusterrolebinding cluster-admin-binding --clusterrole=cluster-admin --user=<name@example.com>
    ```

    where `name@example.com` is the user's account name in GCE.  It's case sensitive. See [this document](https://cloud.google.com/kubernetes-engine/docs/how-to/role-based-access-control) for more information.

---

## Article: advanced-resources/os-distro-specific/csi-on-k3s.md

---
  title: Longhorn CSI on K3s
  weight: 1
---

In this section, you'll learn how to install Longhorn on a K3s Kubernetes cluster. [K3s](https://rancher.com/docs/k3s/latest/en/) is a fully compliant Kubernetes distribution that is easy to install, using half the memory, all in a binary of less than 50mb.

## Requirements

  -  Longhorn v0.7.0 or higher.
  -  `open-iscsi` or `iscsiadm` installed on the node.

## Instruction

  Longhorn v0.7.0 and above support k3s v0.10.0 and above only by default. 
  
  If you want to deploy these new Longhorn versions on versions before k3s v0.10.0, you need to set `--kubelet-root-dir` to `<data-dir>/agent/kubelet` for the Deployment `longhorn-driver-deployer` in `longhorn/deploy/longhorn.yaml`. 
  `data-dir` is a `k3s` arg and it can be set when you launch a k3s server. By default it is `/var/lib/rancher/k3s`.

## Troubleshooting

### Common issues

#### Failed to get arg root-dir: Cannot get kubelet root dir, no related proc for root-dir detection ...

This error is due to Longhorn cannot detect where is the root dir setup for Kubelet, so the CSI plugin installation failed.

You can override the root-dir detection by setting environment variable `KUBELET_ROOT_DIR` in https://github.com/longhorn/longhorn/blob/v{{< current-version >}}/deploy/longhorn.yaml.

#### How to find `root-dir`?

**For K3S prior to v0.10.0**

Run `ps aux | grep k3s` and get argument `--data-dir` or `-d` on k3s node.

e.g.
```
$ ps uax | grep k3s
root      4160  0.0  0.0  51420  3948 pts/0    S+   00:55   0:00 sudo /usr/local/bin/k3s server --data-dir /opt/test/kubelet
root      4161 49.0  4.0 259204 164292 pts/0   Sl+  00:55   0:04 /usr/local/bin/k3s server --data-dir /opt/test/kubelet
``` 
You will find `data-dir` in the cmdline of proc `k3s`. By default it is not set and `/var/lib/rancher/k3s` will be used. Then joining `data-dir` with `/agent/kubelet` you will get the `root-dir`. So the default `root-dir` for K3S is `/var/lib/rancher/k3s/agent/kubelet`.

If K3S is using a configuration file, you would need to check the configuration file to locate the `data-dir` parameter.

**For K3S v0.10.0+**

It is always `/var/lib/kubelet`

## Background 
#### Longhorn versions before v0.7.0 don't work on K3S v0.10.0 or above
K3S now sets its kubelet directory to `/var/lib/kubelet`. See [the K3S release comment](https://github.com/rancher/k3s/releases/tag/v0.10.0) for details.

## Reference
https://github.com/kubernetes-csi/driver-registrar


---

## Article: advanced-resources/os-distro-specific/csi-on-rke-and-coreos.md

---
  title: Longhorn CSI on RKE and CoreOS
  weight: 2
---

For minimalist Linux Operating systems, you'll need a little extra configuration to use Longhorn with RKE (Rancher Kubernetes Engine).  This document outlines the requirements for using RKE and CoreOS.

###  Background 

CSI doesn't work with CoreOS + RKE before Longhorn v0.4.1. The reason is that in the case of CoreOS, RKE sets the argument `root-dir=/opt/rke/var/lib/kubelet` for the kubelet , which is different from the default value `/var/lib/kubelet`.
                                                                             
**For k8s v1.12+**, the kubelet will detect the `csi.sock` according to argument `<--kubelet-registration-path>` passed in by Kubernetes CSI driver-registrar, and `<drivername>-reg.sock` (for Longhorn, it's `io.rancher.longhorn-reg.sock`) on kubelet path `<root-dir>/plugins`.
   
  **For k8s v1.11,** the kubelet will find both sockets on kubelet path `/var/lib/kubelet/plugins`.
   
By default, Longhorn CSI driver creates and expose these two sock files on the host path `/var/lib/kubelet/plugins`. Then the kubelet cannot find `<drivername>-reg.sock`, so CSI driver doesn't work.

Furthermore, the kubelet will instruct the CSI plugin to mount the Longhorn volume on `<root-dir>/pods/<pod-name>/volumes/kubernetes.io~csi/<volume-name>/mount`. But this path inside the CSI plugin container won't be bind mounted on the host path. And the mount operation for the Longhorn volume is meaningless.

Therefore, in this case, Kubernetes cannot connect to Longhorn using the CSI driver without additional configuration.

### Requirements

  -  Kubernetes v1.11 or higher.
  -  Longhorn v0.4.1 or higher.

###  1. Add extra binds for the kubelet

> This step is only required for For CoreOS + and Kubernetes v1.11. It is not needed for Kubernetes v1.12+.

Add extra_binds for kubelet in RKE `cluster.yml`:

```

services:
  kubelet:
    extra_binds:
    - "/opt/rke/var/lib/kubelet/plugins:/var/lib/kubelet/plugins" 

```

This makes sure the kubelet plugins directory is exposed for CSI driver installation.

### 2. Start the iSCSI Daemon

If you want to enable iSCSI daemon automatically at boot, you need to enable the systemd service:

```
sudo su
systemctl enable iscsid
reboot
```

Or just start the iSCSI daemon for the current session:

```
sudo su
systemctl start iscsid
```

### Troubleshooting

#### Failed to get arg root-dir: Cannot get kubelet root dir, no related proc for root-dir detection ...

This error happens because Longhorn cannot detect the root dir setup for the kubelet, so the CSI plugin installation failed.

You can override the root-dir detection by setting environment variable `KUBELET_ROOT_DIR` in https://github.com/longhorn/longhorn/blob/v{{< current-version >}}/deploy/longhorn.yaml.

#### How to find `root-dir`?
 
Run `ps aux | grep kubelet` and get the argument `--root-dir` on host node. 

For example,
```

$ ps aux | grep kubelet
root      3755  4.4  2.9 744404 120020 ?       Ssl  00:45   0:02 kubelet --root-dir=/opt/rke/var/lib/kubelet --volume-plugin-dir=/var/lib/kubelet/volumeplugins

```
You will find `root-dir` in the cmdline of proc `kubelet`. If it's not set, the default value `/var/lib/kubelet` would be used. In the case of CoreOS, the root-dir would be `/opt/rke/var/lib/kubelet` as shown above.

If the kubelet is using a configuration file, you need to check the configuration file to locate the `root-dir` parameter.

### References
https://github.com/kubernetes-csi/driver-registrar

https://coreos.com/os/docs/latest/iscsi.html


---

## Article: advanced-resources/os-distro-specific/okd-support.md

---
title:  OCP/OKD Support
weight: 4
---

To deploy Longhorn on a cluster provisioned with OpenShift 4.x, some additional configurations are required.

> **Note**: OKD currently does not support the ARM platform. For more information, see the [OKD website](https://www.okd.io) and [GitHub issue #1165](https://github.com/okd-project/okd/issues/1165) (*OKD in ARM platform*).

- [Install Longhorn](#install-longhorn)
  - [Install With Helm](#install-with-helm)
  - [Install With `oc` Command](#install-with-oc-command)
- [Prepare A Customized Default Longhorn Disk (Optional)](#prepare-a-customized-default-longhorn-disk-optional)
  - [Add An Extra Disk to Longhorn Storage](#add-an-extra-disk-to-longhorn-storage)
    - [Create Filesystem For The Device](#create-filesystem-for-the-device)
    - [Mounting The Device On Boot with MachineConfig CRD](#mounting-the-device-on-boot-with-machineconfig-crd)
    - [Label and Annotate The Node](#label-and-annotate-the-node)
- [Reference](#reference)
- [Main Contributor](#main-contributor)

## Install Longhorn

### Install With Helm

Please refer to this section [Install with Helm](../../../deploy/install/install-with-helm/) first.

Install Longhorn with the following settings:

| Setting | Value | Example |
| --- | --- | --- |
| `openshift.enabled` | `true` | N/A |
| `image.openshift.oauthProxy.repository` | Upstream image | `quay.io/openshift/origin-oauth-proxy` |
| `image.openshift.oauthProxy.tag` | Version 4.1 or later | `4.18` |

```bash
  helm install longhorn longhorn/longhorn \
    --namespace longhorn-system \
    --create-namespace \
    --set openshift.enabled=true \
    --set image.openshift.oauthProxy.repository=quay.io/openshift/origin-oauth-proxy \
    --set image.openshift.oauthProxy.tag=4.18
```

### Install With `oc` Command

Perform the following steps to install Longhorn on [OKD](https://www.okd.io/) clusters.

1. Download the `longhorn-okd.yaml` file.
  ```
  wget https://raw.githubusercontent.com/longhorn/longhorn/v{{< current-version >}}/deploy/longhorn-okd.yaml
  ```
1. Specify the target `oauth-proxy` container image in the `longhorn-okd.yaml` file (for example, `quay.io/openshift/origin-oauth-proxy:4.18`).

1. Run the following command:
  ```shell
  oc apply -f longhorn-okd.yaml
  ```

One way to monitor the progress of the installation is to watch pods being created in the `longhorn-system` namespace:

  ```shell
    oc get pods \
    --namespace longhorn-system \
    --watch
  ```

For more information, see [Install with Kubectl](../../../deploy/install/install-with-kubectl).

## Prepare A Customized Default Longhorn Disk (Optional)

To understand more about configuring the disks for Longhorn, please refer to the section [Configuring Defaults for Nodes and Disks](../../../nodes-and-volumes/nodes/default-disk-and-node-config/#launch-longhorn-with-multiple-disks)

Longhorn will use the directory `/var/lib/longhorn` as default storage mount point and that means Longhorn use the root device as the default storage. If you don't want to use the root device as the Longhorn storage, set ***defaultSettings.createDefaultDiskLabeledNodes*** true when installing Longhorn by helm:

```txt
--set defaultSettings.createDefaultDiskLabeledNodes=true
```

And then add another device formatted to Longhorn storage

### Add An Extra Disk to Longhorn Storage

#### Create Filesystem For The Device

Create the filesystem on the device with the label `longhorn` on the storage node. Get into the node by oc command:

```bash
oc get nodes --no-headers | awk '{print $1}'
oc debug node/${NODE_NAME} -t -- chroot /host bash
```

Check if the device is present and format it with Longhorn label:

```bash
lsblk
sudo mkfs.ext4 -L longhorn /dev/${DEVICE_NAME}
```

#### Mounting The Device On Boot with MachineConfig CRD

The secondary drive needs to be mounted automatically when node boots up by the `MachineConfig` that can be created and deployed by:

```bash
cat <<EOF >>auto-mount-machineconfig.yaml
apiVersion: machineconfiguration.openshift.io/v1
kind: MachineConfig
metadata:
  labels:
    machineconfiguration.openshift.io/role: worker
  name: 71-mount-storage-worker
spec:
  config:
    ignition:
      version: 3.2.0
    systemd:
      units:
        - name: var-mnt-longhorn.mount
          enabled: true
          contents: |
            [Unit]
            Before=local-fs.target
            [Mount]
            # Example mount point, you can change it to where you like for each device.
            Where=/var/mnt/longhorn
            What=/dev/disk/by-label/longhorn
            Options=rw,relatime,discard
            [Install]
            WantedBy=local-fs.target
EOF

oc apply -f auto-mount-machineconfig.yaml
```

#### Label and Annotate The Node

Please refer to the section [Customizing Default Disks for New Nodes](../../../nodes-and-volumes/nodes/default-disk-and-node-config/#customizing-default-disks-for-new-nodes) to label and annotate storage node on where your device is by oc commands:

```bash
oc get nodes --no-headers | awk '{print $1}'

oc annotate node ${NODE_NAME} --overwrite node.longhorn.io/default-disks-config='[{"path":"/var/mnt/longhorn","allowScheduling":true}]'
oc label node ${NODE_NAME} --overwrite node.longhorn.io/create-default-disk=config
```

**Note**: You might need to reboot the node to validate the modified configuration.

## Reference

- [OCP/OKD Documentation and Helm Support](https://github.com/longhorn/longhorn/pull/5004)
- [OKD Official Website](https://www.okd.io/)
- [OKD Official Documentation Website](https://docs.okd.io/latest/welcome/index.html)
- [oauth-proxy](https://github.com/openshift/oauth-proxy/blob/master/contrib/sidecar.yaml)

## Main Contributor

- [@ArthurVardevanyan](https://github.com/ArthurVardevanyan)


---

## Article: advanced-resources/os-distro-specific/talos-linux-support.md

---
title: Talos Linux Support
weight: 5
---

## Requirements

You must meet the following requirements before installing Longhorn on a Talos Linux cluster.

### System Extensions

Some Longhorn-dependent binary executables are not present in the default Talos root filesystem. To have access to these binaries, Talos offers system extension mechanism to extend the installation.

- `siderolabs/iscsi-tools`: this extension enables iscsid daemon and iscsiadm to be available to all nodes for the Kubernetes persistent volumes operations.
- `siderolabs/util-linux-tools`: this extension enables linux tool to be available to all nodes. For example, the `fstrim` binary is used for Longhorn volume trimming.

The most straightforward method is patching the extensions onto existing Talos Linux nodes.

```yaml
customization:
  systemExtensions:
    officialExtensions:
      - siderolabs/iscsi-tools
      - siderolabs/util-linux-tools
```

For detailed instructions, see the Talos documentation on [System Extensions](https://www.talos.dev/v1.6/talos-guides/configuration/system-extensions/) and [Boot Assets](https://www.talos.dev/v1.6/talos-guides/install/boot-assets/).

### Pod Security

Longhorn requires pod security `enforce: "privileged"`.

By default, Talos Linux applies a `baseline` pod security profile across namespaces, except for the kube-system namespace. This default setting restricts Longhorn's ability to manage and access system resources. For more information, see [Root and Privileged Permission](../../../deploy/install/#root-and-privileged-permission).

For detailed instructions, see [Pod Security Policies Disabled & Pod Security Admission Introduction](../../../../archives/1.7.0/important-notes/#pod-security-policies-disabled--pod-security-admission-introduction) and the Talos documentation on [Pod Security](https://www.talos.dev/v1.6/kubernetes-guides/configuration/pod-security/).

## Talos Linux Version Prior to v1.10.x

### Data Path Mounts

You need provide additional data path mounts to be accessible to the Kubernetes Kubelet container.

These mounts are necessary to provide access to the host directories, and attach volumes required by Longhorn components.

```yaml
machine:
  kubelet:
    extraMounts:
      - destination: /var/lib/longhorn
        type: bind
        source: /var/lib/longhorn
        options:
          - bind
          - rshared
          - rw
```

For detailed instructions, see the Talos documentation on [Editing Machine Configuration](https://www.talos.dev/v1.6/talos-guides/configuration/editing-machine-configuration/).

## V2 Data Engine

To use V2 volumes, all nodes must meet the V2 Data Engine [prerequisites](../../../v2-data-engine/prerequisites#prerequisites).

```yaml
machine:
  sysctls:
    vm.nr_hugepages: "1024"
  kernel:
    modules:
      - name: nvme_tcp
      - name: vfio_pci
#     - name: uio_pci_generic
```

> **Note:**
> Talos Linux v1.7.x and earlier versions do not include the `uio_pci_generic` kernel module. If your system device supports `vfio_pci`, which is the preferred kernel module for SPDK application deployment, you are not required to install and enable the `uio_pci_generic` kernel driver. For more information, see [System Configuration User Guide](https://spdk.io/doc/system_configuration.html) in the SPDK documentation.
>
> You can use `uio_pci_generic` if `vfio_pci` is incompatible with your system or specific hardware. Future versions of Talos Linux are expected to include native support for `uio_pci_generic`. For more information, see [Issue #9236](https://github.com/siderolabs/talos/issues/9236).
> Since 1.8.0 `uio_pci_generic` is now supported.

## Talos Linux Upgrades

### Prior to v1.8.x

When [upgrading a Talos Linux node](https://www.talos.dev/v1.7/talos-guides/upgrading-talos/#talosctl-upgrade), always include the `--preserve` option in the command. This option explicitly tells Talos to keep ephemeral data intact.

Example:

```
talosctl upgrade --nodes 10.20.30.40 --image ghcr.io/siderolabs/installer:v1.7.6 --preserve
```

> **Caution:**
> If you do not include the `--preserve` option, Talos wipes `/var/lib/longhorn`, destroying all replicas stored on that node.

### Recovering from an Upgraded Node without Preserving Data

If you were unable to include the `--preserve` option in the upgrade command, perform the following steps:

1. On the Longhorn UI, go to the **Nodes** page.

1. Select the upgraded node, and then select **Edit node and disks** in the **Operation** menu.

1. On the **Edit Node and Disks** page, set **Scheduling** to **Disable**, delete the disk, and then click **Save**.

1. Select the upgraded node again, and then select **Edit node and disks** in the **Operation** menu.

1. On the **Edit Node and Disks** page, add a disk and configure the following settings:

   - **Path**: Specify `/var/lib/longhorn/`.
   - **Storage Reserved**: Specify a value that matches your requirements. By default, it is set to 30% of the disk capacity.
   - **Scheduling**: Select **Enable**.

1. Click **Save**.

Longhorn synchronizes the replicas based on the configured settings.

### After v1.8.x

The `--preserve` is no longer required. The flag is automatically set for `talosctl upgrade` command [here](https://www.talos.dev/v1.8/introduction/what-is-new/#upgrades).

## Talos Linux Version v1.10.x and Later

### Data Path Mounts

Because Talos Linux deprecated `.machine.disks` we recommend using `UserVolumeConfig` to mount a disk for Longhorn. See the [What's new in Talos v1.10](https://www.talos.dev/v1.10/introduction/what-is-new/#user-volumes) for more details.

You can optionally create also a `VolumeConfig` to specify the size of Talos System volumes, which is _recommended_, like this we avoid the set `defaultSettings.storageReservedPercentageForDefaultDisk`.

> More options of disk configuration can be found in the [Talos documentation](https://www.talos.dev/v1.10/talos-guides/configuration/disk-management/#disk-layout).

You need provide additional data path mounts to be accessible to the Kubernetes Kubelet container.

These mounts are necessary to provide access to the host directories, and attach volumes required by Longhorn components.

The [default data path](../../../references/settings#default-data-path) for Longhorn is `/var/lib/longhorn`. In order to use the below configuration in Talos, we must first set our default data path to `/var/mnt/longhorn`. The method to do this will depend on your [deployment method](../../deploy/customizing-default-settings).

```yaml
machine:
  kubelet:
    extraMounts:
      - destination: /var/mnt/longhorn
        type: bind
        source: /var/mnt/longhorn
        options:
          - bind
          - rshared
          - rw
```

You need to create a `UserVolumeConfig` to mount the disk for Longhorn, which will be automatically mounted to `/var/mnt/longhorn` on the configured node.

```
apiVersion: v1alpha1
kind: UserVolumeConfig
name: longhorn # name is used to identify the volume /var/mnt/<name>
provisioning:
  diskSelector:
    match: disk.transport == "nvme"
  grow: false
  maxSize: 1700GB
```

For detailed instructions on `UserVolumeConfig` and `VolumeConfig`, see the Talos documentation on [Block configuration](https://www.talos.dev/v1.10/reference/configuration/block/)

## References

- [[FEATURE] Talos support](https://github.com/longhorn/longhorn/issues/3161)


---

## Article: advanced-resources/data-integrity/_index.md

---
title: Data Integrity
weight: 5
---

---

## Article: advanced-resources/data-integrity/snapshot-data-integrity-check.md

---
title: Snapshot Data Integrity Check
weight: 2
---

Longhorn is capable of hashing snapshot disk files and periodically checking their integrity.

## Introduction

Longhorn system supports volume snapshotting and stores the snapshot disk files on the local disk. However, it is impossible to check the data integrity of snapshots due to the lack of the checksums of the snapshots previously. As a result, when the data is corrupted due to, for example, the bit rot in the underlying storage, there is no way to detect the corruption and repair the replicas. After applying the feature, Longhorn is capable of hashing snapshot disk files and periodically checking their integrity. When a snapshot disk file in one replica is corrupted, Longhorn will automatically start the rebuilding process to fix it.

## Settings

### Global Settings

- **snapshot-data-integrity** <br>

    This setting allows users to enable or disable snapshot hashing and data integrity checking. Available options are:

    - **disabled**: Disable snapshot disk file hashing and data integrity checking.
    - **enabled**: Enables periodic snapshot disk file hashing and data integrity checking. To detect the filesystem-unaware corruption caused by bit rot or other issues in snapshot disk files, Longhorn system periodically hashes files and finds corrupted ones. Hence, the system performance will be impacted during the periodical checking.
    - **fast-check**: Enable snapshot disk file hashing and fast data integrity checking. Longhorn system only hashes snapshot disk files if their are not hashed or the modification time are changed. In this mode, filesystem-unaware corruption cannot be detected, but the impact on system performance can be minimized.

- **snapshot-data-integrity-immediate-check-after-snapshot-creation** <br>

    Hashing snapshot disk files impacts the performance of the system. The immediate snapshot hashing and checking can be disabled to minimize the impact after creating a snapshot.

- **snapshot-data-integrity-cronjob** <br>

    A schedule defined using the unix-cron string format specifies when Longhorn checks the data integrity of snapshot disk files.

    > **Warning**
    > Hashing snapshot disk files impacts the performance of the system. It is recommended to run data integrity checks during off-peak times and to reduce the frequency of checks.

### Per-Volume Settings

Longhorn also supports the per-volume setting by configuring `Volume.Spec.SnapshotDataIntegrity`. The value is `ignored` by default, so data integrity check is determined by the global setting `snapshot-data-integrity`. `Volume.Spec.SnapshotDataIntegrity` supports `ignored`, `disabled`, `enabled` and `fast-check`. Each volume can have its data integrity check setting customized.

## Performance Impact

For detecting data corruption, checksums of snapshot disk files need to be calculated. The calculations consume storage and computation resources. Therefore, the storage performance will be negatively impacted. In order to provide a clear understanding of the impact, we benchmarked storage performance when checksumming disk files. The read IOPS, bandwidth and latency are negatively impacted.

- Environment
    - Host: AWS EC2 c5d.2xlarge
    - CPU: Intel(R) Xeon(R) Platinum 8124M CPU @ 3.00GHz
    - Memory: 16 GB
    - Network: Up to 10Gbps
    - Kubernetes: v1.24.4+rke2r1
- Result
    - Disk: 200 GiB NVMe SSD as the instance store
      - 100 GiB snapshot with full random data
        {{< figure src="/img/diagrams/snapshot/snapshot_hash_ssd_perf.png" >}}

    - Disk: 200 GiB throughput optimized HDD (st1)
      - 30 GiB snapshot with full random data
        {{< figure src="/img/diagrams/snapshot/snapshot_hash_hdd_perf.png" >}}

## Recommendation

The feature helps detect the data corruption in snapshot disk files of volumes. However, the checksum calculation negatively impacts the storage performance. To lower down the impact, the recommendations are
- Checksumming and checking snapshot disk files can be scheduled to off-peak hours by the global setting `snapshot-data-integrity-cronjob`.
- Disable the global setting `snapshot-data-integrity-immediate-check-after-snapshot-creation`.

---

## Article: advanced-resources/data-recovery/_index.md

---
title: Data Recovery
weight: 12
---


---

## Article: advanced-resources/data-recovery/corrupted-replica.md

---
title: Identifying Corrupted Replicas
weight: 3
---

In the case that one of the disks used by Longhorn went bad, you might experience intermittent input/output errors when using a Longhorn volume. 

For example, one file sometimes cannot be read, but later it can. In this scenario, it's likely one of the disks went bad, resulting in one of the replicas returning incorrect data to the user.

To recover the volume, we can identify the corrupted replica and remove it from the volume:

1. Scale down the workload to detach the volume.
2. Find all the replicas' locations by checking the Longhorn UI. The directories used by the replicas will be shown as a tooltip for each replica in the UI.
3. Log in to each node that contains a replica of the volume and get to the directory that contains the replica data.

    For example, the replica might be stored at:
   
        /var/lib/longhorn/replicas/pvc-06b4a8a8-b51d-42c6-a8cc-d8c8d6bc65bc-d890efb2
4. Run a checksum for every file under that directory.

    For example:
    
    ```
    # sha512sum /var/lib/longhorn/replicas/pvc-06b4a8a8-b51d-42c6-a8cc-d8c8d6bc65bc-d890efb2/*
    fcd1b3bb677f63f58a61adcff8df82d0d69b669b36105fc4f39b0baf9aa46ba17bd47a7595336295ef807769a12583d06a8efb6562c093574be7d14ea4d6e5f4  /var/lib/longhorn/replicas/pvc-06b4a8a8-b51d-42c6-a8cc-d8c8d6bc65bc-d890efb2/revision.counter
    c53649bf4ad843dd339d9667b912f51e0a0bb14953ccdc9431f41d46c85301dff4a021a50a0bf431a931a43b16ede5b71057ccadad6cf37a54b2537e696f4780  /var/lib/longhorn/replicas/pvc-06b4a8a8-b51d-42c6-a8cc-d8c8d6bc65bc-d890efb2/volume-head-000.img
    f6cd5e486c88cb66c143913149d55f23e6179701f1b896a1526717402b976ed2ea68fc969caeb120845f016275e0a9a5b319950ae5449837e578665e2ffa82d0  /var/lib/longhorn/replicas/pvc-06b4a8a8-b51d-42c6-a8cc-d8c8d6bc65bc-d890efb2/volume-head-000.img.meta
    e6f6e97a14214aca809a842d42e4319f4623adb8f164f7836e07dc8a3f4816a0389b67c45f7b0d9f833d50a731ae6c4670ba1956833f1feb974d2d12421b03f7  /var/lib/longhorn/replicas/pvc-06b4a8a8-b51d-42c6-a8cc-d8c8d6bc65bc-d890efb2/volume.meta
    ```

5. Compare the output of each replica. One of them should fail or have different results compared to the others. This will be the one replica we need to remove from the volume.
6. Use the Longhorn UI to remove the identified replica from the volume.
7. Scale up the workload to make sure the error is gone.


---

## Article: advanced-resources/data-recovery/data-error.md

---
title: Identifying and Recovering from Data Errors
weight: 1
---

If you've encountered an error message like the following:

    'fsck' found errors on device /dev/longhorn/pvc-6288f5ea-5eea-4524-a84f-afa14b85780d but could not correct them.

Then you have a data corruption situation. This section describes how to address the issue.

## Bad Underlying Disk

To determine if the error is caused because one of the underlying disks went bad, follow [these steps](../corrupted-replica) to identify corrupted replicas.

If most of the replicas on the disk went bad, that means the disk is unreliable now and should be replaced.

If only one replica on the disk went bad, it can be a situation known as `bit rot`. In this case, removing the replica is good enough.

## Recover from a Snapshot

If all the replicas are identical, then the volume needs to be recovered using snapshots.

The reason for this is probably that the bad bit was written from the workload the volume attached to.

To revert to a previous snapshot:

1. In maintenance mode, attach the volume to any node.
2. Revert to a snapshot. You should start with the latest one.
3. Detach the volume from maintenance mode to any node.
4. Re-attach the volume to a node you have access to.
5. Mount the volume from `/dev/longhorn/<volume_name>` and check the volume content.
6. If the volume content is still incorrect, repeat from step 1.
7. Once you find a usable snapshot, make a new snapshot from there and start using the volume as normal.

## Recover from Backup

If all of the methods above failed, use a backup to [recover the volume.](../../../snapshots-and-backups/backup-and-restore/restore-from-a-backup)

---

## Article: advanced-resources/data-recovery/export-from-replica.md

---
title: Exporting a Volume from a Single Replica
weight: 2
---

Each replica of a Longhorn volume contains the full data for the volume.

If the whole Kubernetes cluster or Longhorn system goes offline, the following steps can be used to retrieve the data of the volume.

1. Identify the volume.

    Longhorn uses the disks on the node to store the replica data.

    By default, the data is stored at the directory specified by the setting [`Default Data Path`](https://longhorn.io/docs/0.8.1/references/settings/#default-data-path).

    More disks can be added to a node by either using the Longhorn UI or by using [a node label and annotation](../../../nodes-and-volumes/nodes/default-disk-and-node-config/).

    You can either keep a copy of the path of those disks, or use the following command to find the disks that have been used by Longhorn. For example:

    ```
    # find / -name longhorn-disk.cfg
    /var/lib/longhorn/longhorn-disk.cfg
    ```

    The result above shows that the path `/var/lib/longhorn` has been used by Longhorn to store data.

2. Check the path found in step 1 to see if it contains the data.

    The data will be stored in the `/replicas` directory, for example:

    ```
    # ls /var/lib/longhorn/replicas/
    pvc-06b4a8a8-b51d-42c6-a8cc-d8c8d6bc65bc-d890efb2
    pvc-71a266e0-5db5-44e5-a2a3-e5471b007cc9-fe160a2c
    ```

    The directory naming pattern is:

    ```
    <volume_name>-<8 bytes UUID>
    ```

    So in the example above, there are two volumes stored here, which are `pvc-06b4a8a8-b51d-42c6-a8cc-d8c8d6bc65bc` and `pvc-71a266e0-5db5-44e5-a2a3-e5471b007cc9`.

    The volume name matches the Kubernetes PV name.

3. Use the `lsof` command to make sure no one is currently using the volume, e.g.
   ```
   # lsof pvc-06b4a8a8-b51d-42c6-a8cc-d8c8d6bc65bc-d890efb2/
   COMMAND    PID USER   FD   TYPE DEVICE SIZE/OFF   NODE NAME
   longhorn 14464 root  cwd    DIR    8,0     4096 541456 pvc-06b4a8a8-b51d-42c6-a8cc-d8c8d6bc65bc-d890efb2
   ```
   The above result shows that the data directory is still being used, so don't proceed to the next step. If it's not being used, `lsof` command should return empty result.
4. Check the volume size of the volume you want to restore using the following command inside the directory:
   ```
   # cat pvc-06b4a8a8-b51d-42c6-a8cc-d8c8d6bc65bc-d890efb2/volume.meta
    {"Size":1073741824,"Head":"volume-head-000.img","Dirty":true,"Rebuilding":false,"Parent":"","SectorSize":512,"BackingFileName":""}
   ```
   From the result above, you can see the volume size is `1073741824` (1 GiB). Note the size.
5. To export the content of the volume, follow the instructions below based on your environment:

   Docker (RKE1)

   To export the content of the volume in Docker, use the following command to create a single replica Longhorn volume container:

   ```
   docker run -v /dev:/host/dev -v /proc:/host/proc -v <data_path>:/volume --privileged longhornio/longhorn-engine:v{{< current-version >}} launch-simple-longhorn <volume_name> <volume_size>
   ```

   For example, based on the information above, the command should be:

   ```
   docker run -v /dev:/host/dev -v /proc:/host/proc -v /var/lib/longhorn/replicas/pvc-06b4a8a8-b51d-42c6-a8cc-d8c8d6bc65bc-d890efb2:/volume --privileged longhornio/longhorn-engine:v{{< current-version >}} launch-simple-longhorn pvc-06b4a8a8-b51d-42c6-a8cc-d8c8d6bc65bc 1073741824
   ```

   Containerd (RKE2/k3s)

   To export the content of the volume in RKE2/k3s, you need to create a static pod manifest. This manifest will launch the Longhorn engine and expose the volume.

   Create a file named `longhorn-recovery.yaml` under `/var/lib/rancher/rke2/agent/pod-manifests/` with the following content:

   ```
   apiVersion: v1
   kind: Pod
   metadata:
     name: longhorn-recovery
     namespace: longhorn-system
   spec:
     nodeName: <node-where-the-replica-is-located>
     hostPID: true
     containers:
     - name: engine
       image: longhornio/longhorn-engine:v<current-version>
       securityContext:
         privileged: true
       command: ["launch-simple-longhorn"]
       args: ["<volume-name>", "<volume-size-in-bytes>"]
       volumeMounts:
       - name: dev
         mountPath: /host/dev
       - name: proc
         mountPath: /host/proc
       - name: data
         mountPath: /volume
     volumes:
     - name: dev
       hostPath:
         path: /dev
     - name: proc
       hostPath:
         path: /proc
     - name: data
       hostPath:
         path: <host-path-to-replica>
     restartPolicy: Never
   ```
   Replace the placeholders with the following:
   - `<current-version>`: The version of Longhorn you are using.
   - `<volume-name>`: The name of the volume you want to recover.
   - `<host-path-to-replica>`: The path to the replica directory you found in Step 1.

**Result:** Now you should have a block device created on `/dev/longhorn/<volume_name>` for this device, such as `/dev/longhorn/pvc-06b4a8a8-b51d-42c6-a8cc-d8c8d6bc65bc` for the example above. Now you can mount the block device to get the access to the data.

> To avoid accidental change of the volume content, it's recommended to use `mount -o ro` to mount the directory as `readonly`.

After you are done accessing the volume content, use `docker stop` to stop the container. For RKE2, clean up the resources by removing the static pod manifest file `sudo rm /var/lib/rancher/rke2/agent/pod-manifests/longhorn-recovery.yaml`


---

## Article: advanced-resources/data-recovery/full-disk.md

---
title: Recovering from a Full Disk
weight: 4
---

If one disk used by one of the Longhorn replicas is full, that replica will go to the error state, and Longhorn will rebuild another replica on another node/disk.

To recover from a full disk,

1. Disable the scheduling for the full disk.
    
    Longhorn should have already marked the disk as `unschedulable`.
    
    This step is to make sure the disk will not be scheduled to by accident after more space is freed up.

2. Identify the replicas in the error state on the disk using the Longhorn UI's disk page.
3. Remove the replicas in the error state.

## Recommended after Recovery

We recommend adding more disks or more space to the node that had this situation.

---

## Article: advanced-resources/data-recovery/recover-without-system.md

---
title: Recovering from a Longhorn Backup without System Installed
weight: 5
---

This command gives users the ability to restore a backup to a `raw` image or a `qcow2` image. If the backup is based on a backing file, users should provide the backing file as a `qcow2` image with `--backing file` parameter.

1. Copy the [yaml template](https://github.com/longhorn/longhorn/blob/v{{< current-version >}}/examples/restore_to_file.yaml.template): Make a copy of `examples/restore_to_file.yaml.template` as e.g. `restore.yaml`.
    
2. Set the node which the output file should be placed on by replacing `<NODE_NAME>`, e.g. `node1`.

3. Specify the host path of output file by modifying field `hostpath` of volume `disk-directory`. By default the directory is `/tmp/restore/`.

4. Set the first argument (backup url) by replacing `<BACKUP_URL>`, e.g. `s3://<your-bucket-name>@<your-aws-region>/backupstore?backup=<backup-name>&volume=<volume-name>`.

    - `<backup-name>` and `<volume-name>` can be retrieved from backup.cfg stored in the backup destination folder, e.g. `backup_backup-72bcbdad913546cf.cfg`. The content will be like below: 

        ```json
        {"Name":"backup-72bcbdad913546cf","VolumeName":"volume_1","SnapshotName":"79758033-a670-4724-906f-41921f53c475"}
        ```

5. Set argument `output-file` by replacing `<OUTPUT_FILE>`, e.g. `volume.raw` or `volume.qcow2`.

6. Set argument `output-format` by replacing `<OUTPUT_FORMAT>`. The supported options are `raw` or `qcow2`.

7. Set argument `longhorn-version` by replacing `<LONGHORN_VERSION>`, e.g. `v{{< current-version >}}`

8. Set the S3 Credential Secret by replacing `<S3_SECRET_NAME>`, e.g. `minio-secret`.  

    - The credential secret can be referenced [here](https://longhorn.io/docs/{{< current-version >}}/snapshots-and-backups/backup-and-restore/set-backup-target/#set-up-aws-s3-backupstore) and must be created in the `longhorn-system' namespace.

9. Execute the yaml using e.g.:

        kubectl create -f restore.yaml

10. Watch the result using:

        kubectl -n longhorn-system get pod restore-to-file -w

After the pod status changed to `Completed`, you should able to find `<OUTPUT_FILE>` at e.g. `/tmp/restore` on the `<NODE_NAME>`.

We also provide a script, [restore-backup-to-file.sh](https://raw.githubusercontent.com/longhorn/longhorn/v{{< current-version >}}/scripts/restore-backup-to-file.sh), to restore a backup. The following parameters should be specified:
  - `--backup-url`: Specifies the backups S3/NFS URL. e.g., `s3://backupbucket@us-east-1/backupstore?backup=backup-bd326da2c4414b02&volume=volumeexamplename"`
  
  - `--output-file`: Set the output file name. e.g, `volume.raw`
  
  - `--output-format`: Set the output file format. e.g. `raw` or `qcow2`
  
  - `--version`: Specifies the version of Longhorn to use. e.g., `v{{< current-version >}}`

Optional parameters can be specified:

  - `--aws-access-key`: Specifies AWS credentials access key if backups is s3.
  
  - `--aws-secret-access-key`: Specifies AWS credentials access secret key if backups is s3.
  
  - `--backing-file`: backing image. e.g., `/tmp/backingfile.qcow2`

The output image files can be found in the `/tmp/restore` folder after the script has finished running.

---

## Article: advanced-resources/data-cleanup/_index.md

---
title: Data Cleanup
weight: 6
---

---

## Article: advanced-resources/data-cleanup/orphaned-data-cleanup.md

---
title: Orphaned Data Cleanup
weight: 4
---

Longhorn supports orphaned data cleanup. Currently, Longhorn can identify and clean up the orphaned replica directories on disks.

## Orphaned Replica Directories

When a user introduces a disk into a Longhorn node, it may contain replica directories that are not tracked by the Longhorn system. The untracked replica directories may belong to other Longhorn clusters. Or, the replica CRs associated with the replica directories are removed after the node or the disk is down. When the node or the disk comes back, the corresponding replica data directories are no longer tracked by the Longhorn system. These replica data directories are called orphaned.

Longhorn supports the detection and cleanup of orphaned replica directories. It identifies the directories and gives a list of `orphan` resources that describe those directories. By default, Longhorn does not automatically delete `orphan` resources and their directories. Users can trigger the deletion of orphaned replica directories manually or have it done automatically.

When automatic orphan deletion is enabled, Longhorn automatically deletes orphaned Custom Resources (CRs) and their associated directories after the delay defined by the `orphan-resource-auto-deletion-grace-period` setting. If a user manually deletes an orphaned CR, the deletion occurs immediately and does not respect this grace period.

### Example

In the example, we will explain how to manage orphaned replica directories identified by Longhorn via `kubectl` and Longhorn UI.

#### Manage Orphaned Replica Directories via kubectl

1. Introduce disks containing orphaned replica directories.
   - Orphaned replica directories on Node `worker1` disks
    ```
    # ls /mnt/disk/replicas/
    pvc-19c45b11-28ee-4802-bea4-c0cabfb3b94c-15a210ed
    ```
   - Orphaned replica directories on Node `worker2` disks
    ```
    # ls /var/lib/longhorn/replicas/
    pvc-28255b31-161f-5621-eea3-a1cbafb4a12a-866aa0a5

    # ls /mnt/disk/replicas/
    pvc-19c45b11-28ee-4802-bea4-c0cabfb3b94c-a86771c0
    ```

2. Longhorn detects the orphaned replica directories and creates an `orphan` resources describing the directories.
    ```
    # kubectl -n longhorn-system get orphans -l "longhorn.io/orphan-type=replica"
    NAME                                                                      TYPE      NODE
    orphan-fed8c6c20965c7bdc3e3bbea5813fac52ccd6edcbf31e578f2d8bab93481c272   replica   rancher60-worker1
    orphan-637f6c01660277b5333f9f942e4b10071d89379dbe7b4164d071f4e1861a1247   replica   rancher60-worker2
    orphan-6360f22930d697c74bec4ce4056c05ac516017b908389bff53aca0657ebb3b4a   replica   rancher60-worker2
    ```

3. One can list the `orphan` resources created by Longhorn system by `kubectl -n longhorn-system get orphan`.
    ```
    kubectl -n longhorn-system get orphan
    ```

4. Get the detailed information of one of the orphaned replica directories in `spec.parameters` by `kubcel -n longhorn-system get orphan <name>`.
   ```
    # kubectl -n longhorn-system get orphans orphan-fed8c6c20965c7bdc3e3bbea5813fac52ccd6edcbf31e578f2d8bab93481c272 -o yaml
    apiVersion: longhorn.io/v1beta2
    kind: Orphan
    metadata:
    creationTimestamp: "2022-04-29T10:17:40Z"
    finalizers:
    - longhorn.io
    generation: 1
    labels:
        longhorn.io/component: orphan
        longhorn.io/managed-by: longhorn-manager
        longhorn.io/orphan-type: replica
        longhornnode: rancher60-worker1

    ......

    spec:
    nodeID: rancher60-worker1
    orphanType: replica
    parameters:
        DataName: pvc-19c45b11-28ee-4802-bea4-c0cabfb3b94c-15a210ed
        DiskName: disk-1
        DiskPath: /mnt/disk/
        DiskUUID: 90f00e61-d54e-44b9-a095-35c2b56a0462
    status:
    conditions:
    - lastProbeTime: ""
        lastTransitionTime: "2022-04-29T10:17:40Z"
        message: ""
        reason: ""
        status: "True"
        type: DataCleanable
    - lastProbeTime: ""
        lastTransitionTime: "2022-04-29T10:17:40Z"
        message: ""
        reason: ""
        status: "False"
        type: Error
    ownerID: rancher60-worker1
   ```

5. One can delete the `orphan` resource by `kubectl -n longhorn-system delete orphan <name>` and then the corresponding orphaned replica directory will be deleted.
   ```
    # kubectl -n longhorn-system delete orphan orphan-fed8c6c20965c7bdc3e3bbea5813fac52ccd6edcbf31e578f2d8bab93481c272

    # kubectl -n longhorn-system get orphans -l "longhorn.io/orphan-type=replica"
    NAME                                                                      TYPE      NODE
    orphan-637f6c01660277b5333f9f942e4b10071d89379dbe7b4164d071f4e1861a1247   replica   rancher60-worker2
    orphan-6360f22930d697c74bec4ce4056c05ac516017b908389bff53aca0657ebb3b4a   replica   rancher60-worker2
   ```

    The orphaned replica directory is deleted.
    ```
    # ls /mnt/disk/replicas/

    ```

6. By default, Longhorn will not automatically delete the orphaned replica directory. You can enable automatic deletion by setting the `orphan-resource-auto-deletion` setting.
    ```
    # kubectl -n longhorn-system edit settings.longhorn.io orphan-resource-auto-deletion
    ```
    Then, add `replica-data` to the list by including it as one of the semicolon-separated items.

    ```
    # kubectl -n longhorn-system get settings.longhorn.io orphan-resource-auto-deletion
    NAME                            VALUE          APPLIED   AGE
    orphan-resource-auto-deletion   replica-data   true      26m
    ```

7. After enabling the automatic deletion and wait for a while, the `orphan` resources and directories are deleted automatically.
   ```
    # kubectl -n longhorn-system get orphans.longhorn.io -l "longhorn.io/orphan-type=replica"
    No resources found in longhorn-system namespace.
   ```
   The orphaned replica directories are deleted.
   ```
    # ls /mnt/disk/replicas/

    # ls /var/lib/longhorn/replicas/

   ```

    Additionally, one can delete all orphaned replica directories on the specified node by
    ```
    # kubectl -n longhorn-system delete orphan -l "longhorn.io/orphan-type=replica-instance,longhornnode=<node name>”
    ```

#### Manage Orphaned Replica Directories via Longhorn UI

1. In the top navigation bar, go to `Setting > Orphan Resources > Replica Data`.
2. Review the list of orphaned replica directories grouped by node and disk.
3. To delete a directory, click `Operation > Delete`.

By default, Longhorn does not delete orphaned replica directories automatically. To enable automatic deletion, go to `Settings > Orphan`.

### Exception
Longhorn will not create an `orphan` resource for an orphaned directory when

- The orphaned directory is not an **orphaned replica directory**.
  - The directory name does not follow the replica directory's naming convention.
  - The volume volume.meta file is missing.
- The orphaned replica directory is on an evicted node.
- The orphaned replica directory is in an evicted disk.
- The orphaned data cleanup mechanism does not clean up a stale replica, also known as an error replica. Instead, the stale replica is cleaned up according to the [staleReplicaTimeout](../../../nodes-and-volumes/volumes/create-volumes/#creating-longhorn-volumes-with-kubectl) setting.


---

## Article: advanced-resources/data-cleanup/orphaned-instance-cleanup.md

---
title: Orphaned Instance Cleanup
weight: 4
---

Longhorn can identify and clean up orphaned instances on each node.

## Orphaned Runtime Instance

When a network outage affects a Longhorn node, it may leave behind engine or replica runtime instances that are no longer tracked by the Longhorn system. The corresponding engine and replica CRs might be removed or rescheduled to another node during the outage. When the node comes back, the Longhorn system no longer tracks the corresponding runtime instances. These runtime instances, such as the engine and replica processes for v1 volumes, are called orphaned instances. Orphaned instances continue to consume CPU and memory.

Longhorn supports the detection and cleanup of orphaned instance. It identifies the instances and gives a list of `orphan` resources that describe those orphans. By default, Longhorn does not automatically delete `orphan` instances. Users can trigger the deletion of orphaned instances manually or have it done automatically.

When automatic orphan deletion is enabled, Longhorn automatically deletes orphaned Custom Resources (CRs) and their associated directories after the delay defined by the `orphan-resource-auto-deletion-grace-period` setting. If a user manually deletes an orphaned CR, the deletion occurs immediately and does not respect this grace period.

### Example

The following example shows how to manage orphaned instances using `kubectl`.

#### Manage Orphaned Instances via kubectl

1. Introduce nodes that orphaned instance processes running on
    - Orphaned replica instance on Node `worker1`
     ```
     # kubectl -n longhorn-system describe instancemanager -l "longhorn.io/node=worker1"
     Name:         instance-manager-8ff396d6d3744979b32abafc6346781c
     Namespace:    longhorn-system
     Kind:         InstanceManager
     ...
     Status:
       Instance Replicas:
         pvc-569e44c0-b352-4aca-bf14-2cf7a6cfe86f-r-05660b73:
           Spec:
             Data Engine:  v1
             Name:         pvc-569e44c0-b352-4aca-bf14-2cf7a6cfe86f-r-05660b73
           Status:
             Conditions:         <nil>
             Endpoint:
             Error Msg:
             Listen:
             Port End:           10020
             Port Start:         10011
             Resource Version:   0
             State:              running
             Target Port End:    0
             Target Port Start:  0
             Type:               replica
     ...
     ```
    - Orphaned engine instance on Node `worker2`
     ```
     # kubectl -n longhorn-system describe instancemanager -l "longhorn.io/node=worker2"
     Name:         instance-manager-b87f10b867cec1dca2b814f5e78bcc90
     Namespace:    longhorn-system
     Kind:         InstanceManager
     ...
     Status:
       Instance Engines:
         pvc-569e44c0-b352-4aca-bf14-2cf7a6cfe86f-e-0:
           Spec:
             Data Engine:  v1
             Name:         pvc-569e44c0-b352-4aca-bf14-2cf7a6cfe86f-e-0
           Status:
             Conditions:
               Filesystem Read Only:  false
             Endpoint:
             Error Msg:
             Listen:
             Port End:                10020
             Port Start:              10020
             Resource Version:        0
             State:                   running
             Target Port End:         10020
             Target Port Start:       10020
             Type:                    engine
     ...
     ```

2. Longhorn detects the orphaned instances and creates an `orphan` resources describing the instances.
    ```
    # kubectl -n longhorn-system get orphan -l "longhorn.io/orphan-type in (engine-instance,replica-instance)"
    NAME                                                                      TYPE               NODE
    orphan-1807009489e50534c35c350e22680449c97deca4e5d3b72f4591976145f8bc41   engine-instance    worker2
    orphan-a91aa42ab5eda6b8b9fe1116d5b5f5673e5108d89be3db6fd18a275913463eef   replica-instance   worker1
    ```

3. You can view the list of `orphan` resources created by the Longhorn system by running `kubectl -n longhorn-system get orphan`.
    ```
    kubectl -n longhorn-system get orphan
    ```

4. Get the detailed information of one of the orphaned replica instance in `spec.parameters` by `kubectl -n longhorn-system get orphan <name>`.
    ```
    # kubectl -n longhorn-system get orphans orphan-a91aa42ab5eda6b8b9fe1116d5b5f5673e5108d89be3db6fd18a275913463eef -o yaml
    apiVersion: longhorn.io/v1beta2
    kind: Orphan
    metadata:
    creationTimestamp: "2025-05-02T06:07:32Z"
    finalizers:
    - longhorn.io
    generation: 1
    labels:
        longhorn.io/component: orphan
        longhorn.io/managed-by: longhorn-manager
        longhorn.io/orphan-type: replica-instance
        longhornnode: worker1
        longhornreplica: pvc-569e44c0-b352-4aca-bf14-2cf7a6cfe86f-r-05660b73

    ......

    spec:
      dataEngine: v1
      nodeID: worker1
      orphanType: replica-instance
      parameters:
        InstanceManager: instance-manager-8ff396d6d3744979b32abafc6346781c
        InstanceName: pvc-569e44c0-b352-4aca-bf14-2cf7a6cfe86f-r-05660b73
    status:
    conditions:
    - lastProbeTime: ""
        lastTransitionTime: "2025-05-02T06:06:39Z"
        message: ""
        reason: running
        status: "True"
        type: InstanceExist
    - lastProbeTime: ""
        lastTransitionTime: "2025-05-02T06:06:39Z"
        message: ""
        reason: ""
        status: "False"
        type: Error
    ownerID: worker1
    ```

5. Get the detailed information of one of the orphaned engine instance in `spec.parameters` by `kubectl -n longhorn-system get orphan <name>`.
    ```
    # kubectl -n longhorn-system get orphans orphan-1807009489e50534c35c350e22680449c97deca4e5d3b72f4591976145f8bc41 -o yaml
    apiVersion: longhorn.io/v1beta2
    kind: Orphan
    metadata:
    creationTimestamp: "2025-05-02T06:47:25Z"
    finalizers:
    - longhorn.io
    generation: 1
    labels:
        longhorn.io/component: orphan
        longhorn.io/managed-by: longhorn-manager
        longhorn.io/orphan-type: engine-instance
        longhornengine: pvc-569e44c0-b352-4aca-bf14-2cf7a6cfe86f-e-0
        longhornnode: worker2

    ......

    spec:
      dataEngine: v1
      nodeID: worker2
      orphanType: engine-instance
      parameters:
        InstanceManager: instance-manager-b87f10b867cec1dca2b814f5e78bcc90
        InstanceName: pvc-569e44c0-b352-4aca-bf14-2cf7a6cfe86f-e-0
    status:
    conditions:
    - lastProbeTime: ""
        lastTransitionTime: "2025-05-02T06:47:25Z"
        message: ""
        reason: running
        status: "True"
        type: InstanceExist
    - lastProbeTime: ""
        lastTransitionTime: "2025-05-02T06:47:25Z"
        message: ""
        reason: ""
        status: "False"
        type: Error
    ownerID: worker2
    ```

6. You can delete an `orphan` resource by running `kubectl -n longhorn-system delete orphan <name>`. The corresponding orphaned instance will also be removed.
    ```
    # kubectl -n longhorn-system delete orphan orphan-a91aa42ab5eda6b8b9fe1116d5b5f5673e5108d89be3db6fd18a275913463eef

    # kubectl -n longhorn-system get orphan -l "longhorn.io/orphan-type in (engine-instance,replica-instance)"
    NAME                                                                      TYPE               NODE
    orphan-1807009489e50534c35c350e22680449c97deca4e5d3b72f4591976145f8bc41   engine-instance    worker2
    ```

    The orphaned instance is deleted.
    ```
    # kubectl -n longhorn-system describe instancemanager -l "longhorn.io/node=worker1"
    Name:         instance-manager-8ff396d6d3744979b32abafc6346781c
    Namespace:    longhorn-system
    Kind:         InstanceManager
    ...
    Status:
      Instance Replicas:
    ...
    ```

7. By default, Longhorn does not automatically delete orphaned instances. You can enable automatic deletion by configuring the `orphan-resource-auto-deletion` setting.
    ```
    # kubectl -n longhorn-system edit settings.longhorn.io orphan-resource-auto-deletion
    ```
    Then, add `instance` to the list by including it as one of the semicolon-separated items.

    ```
    # kubectl -n longhorn-system get settings.longhorn.io orphan-resource-auto-deletion
    NAME                            VALUE      APPLIED   AGE
    orphan-resource-auto-deletion   instance   true      45h
    ```

8. After enabling the automatic deletion and wait for a while, the `orphan` resources and processes are deleted automatically.
    ```
    # kubectl -n longhorn-system get orphans.longhorn.io -l "longhorn.io/orphan-type in (engine-instance,replica-instance)"
    No resources found in longhorn-system namespace.
    ```
    The orphaned instances are deleted from instance manager.
    ```
    # kubectl -n longhorn-system describe instancemanager -l "longhorn.io/node=worker1"
    Name:         instance-manager-8ff396d6d3744979b32abafc6346781c
    Namespace:    longhorn-system
    Kind:         InstanceManager
    ...
    Status:
      Instance Replicas:
    ...

    # kubectl -n longhorn-system describe instancemanager -l "longhorn.io/node=worker2"
    Name:         instance-manager-b87f10b867cec1dca2b814f5e78bcc90
    Namespace:    longhorn-system
    Kind:         InstanceManager
    ...
    Status:
      Instance Engines:
    ...

    ```

    Additionally, you can delete all orphaned instances on the specified node by running:
    ```
    # kubectl -n longhorn-system delete orphan -l "longhorn.io/orphan-type in (engine-instance,replica-instance),longhornnode=<node name>"
    ```

#### Manage Orphaned Instances via Longhorn UI

1. In the top navigation bar, go to `Advanced > Orphaned Resources > Instances`.
2. Review the list of orphaned instances, displaying relevant instance information.
3. To delete an orphaned instance, click `Operation > Delete`.

By default, Longhorn does not automatically delete orphaned instances. To enable automatic deletion, go to `Settings > Orphan`.

### Exception

Longhorn does not create an orphan resource in the following scenarios:

- The orphaned engine or replica is rescheduled back to the node.
- The engine or replica is in a migrating, starting, or stopping state.
- The node is evicted.


---

## Article: advanced-resources/cluster-restore/_index.md

---
title: Cluster Restore
weight: 11
---

---

## Article: advanced-resources/cluster-restore/rancher-cluster-restore.md

---
title: Restore cluster with a Rancher snapshot
weight: 4
---

This doc describes what users need to do after restoring the cluster with a Rancher snapshot.

## Assumptions:
- Most of the data and the underlying disks still exist in the cluster before the restore and can be directly reused then.
- There is a backupstore holding all volume data.
- The setting [`Disable Revision Counter`](../../../references/settings/#disable-revision-counter) is false. (It's false by default.) Otherwise, users need to manually check if the data among volume replicas are consistent, or directly restore volumes from backup.

## Expectation:
- All settings and node & disk configs will be restored.
- As long as the valid data still exists, the volumes can be recovered without using a backup. In other words, we will try to avoid restoring backups, which may help reduce Recovery Time Objective (RTO) as well as save bandwidth.
- Detect the invalid or out-of-sync replicas as long as the related volume still contains a valid replica after the restore.

## Behaviors & Requirement of Rancher restore
According to [the Rancher restore article](https://ranchermanager.docs.rancher.com/how-to-guides/new-user-guides/backup-restore-and-disaster-recovery/restore-rancher-launched-kubernetes-clusters-from-backup), you have to restart the Kubernetes components on all nodes. Otherwise, there will be tons of resource update conflicts in Longhorn.

## Actions after the restore
- Restart all Kubernetes components for all nodes. See the above link for more details.

- Kill all longhorn manager pods then Kubernetes will automatically restart them. Wait for conflicts in longhorn manager pods to disappear.

- All volumes may be reattached. If a Longhorn volume is used by a single pod, users need to shut down then recreate it. For Deployments or Statefulsets, Longhorn will automatically kill then restart the related pods.

- If the following happens after the snapshot and before the cluster restore:
    - A volume is unchanged: Users don't need to do anything.
    - The data is updated: Users don't need to do anything typically. Longhorn will automatically fail the replicas that don't contain the latest data.
    - A new volume is created: This volume will disappear after the restore. Users need to recreate a new volume, launch [a single replica volume](../../data-recovery/export-from-replica) based on the replica of the disappeared volume, then transfer the data to the new volume.
    - A volume is deleted: Since the data is cleaned up when the volume is removed, the restored volume contains no data. Users may need to re-delete it.
    - For DR volumes: Users don't need to do anything. Longhorn will redo a full restore.
    - Some operations are applied for a volume:
        - Backup: The backup info of the volume should be resynced automatically.
        - Snapshot: The snapshot info of the volume should be resynced once the volume is attached.
        - Replica rebuilding & replica removal:
            - If there are new replicas rebuilt, those replicas will disappear from the Longhorn system after the restoring. Users need to clean up the replica data manually, or use the data directories of these replicas to export a single replica volume then do data recovery if necessary.
            - If there are some failed/removed replicas and there is at least one replica keeping healthy, those failed/removed replicas will be back after the restoration. Then Longhorn can detect these restored replicas do not contain any data, and copy the latest data from the healthy replica to these replicas.
          - If all replicas are replaced by new replicas after the snapshot, the volume will contain invalid replicas only after the restore. Then users need to export [a single replica volume](../../data-recovery/export-from-replica) for the data recovery.
        -  Engine image upgrade: Users need to redo the upgrade.
        - Expansion: The spec size of the volume will be smaller than the current size. This is like someone requesting volume shrinking but actually Longhorn will refuse to handle it internally. To recover the volume, users need to scale down the workloads and re-do the expansion.

    - **Notice**: If users don't know how to recover a problematic volume, the simplest way is always restoring a new volume from backup.

- If the Longhorn system is upgraded after the snapshot, the new settings and the modifications on the node config will disappear. Users need to re-do the upgrade, then re-modify the settings and node configurations.

- If a node is deleted from Longhorn system after the snapshot, the node won't be back, but the pods on the removed node will be restored. Users need to manually clean up them since these pod may get stuck in state `Terminating`.
- If a node to added to Longhorn system after the snapshot, Longhorn should automatically relaunch all necessary workloads on the node after the cluster restore. But users should be aware that all new replicas or engines on this node will be gone after the restore.


## References
- The related GitHub issue is https://github.com/longhorn/longhorn/issues/2228.
  In this GitHub post, one user is providing a way that restores the Longhorn to a new cluster that doesn't contain any data.


---

## Article: advanced-resources/deploy/_index.md

---
title: Deploy
weight: 1
---

---

## Article: advanced-resources/deploy/customizing-default-settings.md

---
title: Customizing Default Settings
weight: 1
---

You may customize Longhorn's default settings while installing or upgrading. You may specify, for example, `Create Default Disk With Node Labeled` and `Default Data Path` before starting Longhorn.

The default settings can be customized in the following ways:

- [Installation](#installation)
  - [Using the Rancher UI](#using-the-rancher-ui)
  - [Using the Longhorn Deployment YAML File](#using-the-longhorn-deployment-yaml-file)
  - [Using Helm](#using-helm)
  - [Using Helm Controller](#using-helm-controller)
- [Update Settings](#update-settings)
  - [Using the Longhorn UI](#using-the-longhorn-ui)
  - [Using the Rancher UI](#using-the-rancher-ui-1)
  - [Using Kubectl](#using-kubectl)
  - [Using Helm](#using-helm-1)
- [Upgrade](#upgrade)
  - [Using the Rancher UI](#using-the-rancher-ui-2)
  - [Using the Longhorn Deployment YAML File](#using-the-longhorn-deployment-yaml-file-1)
  - [Using Helm](#using-helm-2)
- [History](#history)


> **NOTE:** When using Longhorn Deployment YAML file or Helm for installation, updating or upgrading, if the value of a default setting is an empty string and valid, the default setting will be cleaned up in Longhorn. If not, Longhorn will ignore the invalid values and will not update the default values.

## Installation
### Using the Rancher UI

From the project view in Rancher, go to **Apps && Marketplace > Longhorn > Install > Next > Edit Options > Longhorn Default Settings > Customize Default Settings** and edit the settings before installing the app.

### Using the Longhorn Deployment YAML File

1. Download the longhorn repo:

    ```shell
    git clone https://github.com/longhorn/longhorn.git
    ```

1. Modify the config map named `longhorn-default-setting` in the yaml file `longhorn/deploy/longhorn.yaml`.

    In the below example, users customize the default settings, backup-target, backup-target-credential-secret, and default-data-path.
    When the setting is absent or has a leading `#` symbol, the default setting will use the default value in Longhorn or the customized values previously configured.

    ```yaml
    ---
    apiVersion: v1
    kind: ConfigMap
    metadata:
      name: longhorn-default-setting
      namespace: longhorn-system
    data:
      default-setting.yaml: |-
        backup-target: s3://backupbucket@us-east-1/backupstore
        backup-target-credential-secret: minio-secret
        #allow-recurring-job-while-volume-detached:
        #create-default-disk-labeled-nodes:
        default-data-path: /var/lib/longhorn-example/
        #replica-soft-anti-affinity:
        #replica-auto-balance:
        #storage-over-provisioning-percentage:
        #storage-minimal-available-percentage:
        #upgrade-checker:
        #upgrade-responder-url:
        #default-replica-count:
        #default-data-locality:
        #default-longhorn-static-storage-class:
        #backupstore-poll-interval:
        #taint-toleration:
        #system-managed-components-node-selector:
        #priority-class:
        #auto-salvage:
        #auto-delete-pod-when-volume-detached-unexpectedly:
        #disable-scheduling-on-cordoned-node:
        #replica-zone-soft-anti-affinity:
        #replica-disk-soft-anti-affinity:
        #node-down-pod-deletion-policy:
        #node-drain-policy:
        #replica-replenishment-wait-interval:
        #concurrent-replica-rebuild-per-node-limit:
        #disable-revision-counter:
        #system-managed-pods-image-pull-policy:
        #allow-volume-creation-with-degraded-availability:
        #auto-cleanup-system-generated-snapshot:
        #concurrent-automatic-engine-upgrade-per-node-limit:
        #backing-image-cleanup-wait-interval:
        #backing-image-recovery-wait-interval:
        #guaranteed-instance-manager-cpu:
        #kubernetes-cluster-autoscaler-enabled:
        #orphan-resource-auto-deletion:
        #storage-network:
        #recurring-successful-jobs-history-limit:
        #recurring-failed-jobs-history-limit:
    ---
    ```

### Using Helm

> **NOTE:**
> Use Helm 3 when installing and upgrading Longhorn. Helm 2 is [no longer supported](https://helm.sh/blog/helm-2-becomes-unsupported/).

Use the Helm command with the `--set` flag to modify the default settings. For example:

```shell
helm install longhorn longhorn/longhorn \
  --namespace longhorn-system \
  --create-namespace \
  --set defaultSettings.taintToleration="key1=value1:NoSchedule; key2:NoExecute"
```

You can also provide a copy of the `values.yaml` file with the default settings modified to the `--values` flag when running the Helm command:

1. Obtain a copy of the `values.yaml` file from GitHub:

    ```shell
    curl -Lo values.yaml https://raw.githubusercontent.com/longhorn/charts/master/charts/longhorn/values.yaml
    ```

2. Modify the default settings in the YAML file. The following is an example snippet of `values.yaml`:

   When the setting is absent or has a leading `#` symbol, the default setting will use the default value in Longhorn or the customized values previously configured.

    ```yaml
    defaultSettings:
      backupTarget: s3://backupbucket@us-east-1/backupstore
      backupTargetCredentialSecret: minio-secret
      createDefaultDiskLabeledNodes: true
      defaultDataPath: /var/lib/longhorn-example/
      replicaSoftAntiAffinity: false
      storageOverProvisioningPercentage: 600
      storageMinimalAvailablePercentage: 15
      upgradeChecker: false
      defaultReplicaCount: 2
      defaultDataLocality: disabled
      defaultLonghornStaticStorageClass: longhorn-static-example
      backupstorePollInterval: 500
      taintToleration: key1=value1:NoSchedule; key2:NoExecute
      systemManagedComponentsNodeSelector: "label-key1:label-value1"
      priorityClass: high-priority
      autoSalvage: false
      disableSchedulingOnCordonedNode: false
      replicaZoneSoftAntiAffinity: false
      replicaDiskSoftAntiAffinity: false
      volumeAttachmentRecoveryPolicy: never
      nodeDownPodDeletionPolicy: do-nothing
      guaranteedInstanceManagerCpu: 15
      orphanResourceAutoDeletion: ""
      orphanResourceAutoDeletionGracePeriod: 300
    ```

3. Run Helm with `values.yaml`:

   ```shell
   helm install longhorn longhorn/longhorn \
     --namespace longhorn-system \
     --create-namespace \
     --values values.yaml
   ```

For more info about using helm, see the section about
[installing Longhorn with Helm](../../../deploy/install/install-with-helm)

### Using Helm Controller

In the HelmChart YAML file, add lines to spec.set with the desired settings:
```yaml
spec:
  ...
  set:
    defaultSettings.priorityClass: system-node-critical
    defaultSettings.replicaAutoBalance: least-effort
    defaultSettings.storageOverProvisioningPercentage: "200"
    persistence.defaultClassReplicaCount: "2"

```

## Update Settings

### Using the Longhorn UI

We recommend using the Longhorn UI to change Longhorn setting on the existing cluster. It would make the setting persistent.

### Using the Rancher UI

From the project view in Rancher, go to **Apps && Marketplace > Longhorn > Upgrade > Next > Edit Options > Longhorn Default Settings > Customize Default Settings** and edit the settings before upgrading the app to the current Longhorn version.

### Using Kubectl

If you prefer to update the setting from the command line, use `kubectl`.
To avoid collisions with other CRDs, do not use the simple `settings`. Instead, use `settings.longhorn.io` or `lhs`.
```shell
kubectl edit settings.longhorn.io <SETTING-NAME> -n longhorn-system
```

### Using Helm

Modify the default settings in the YAML file as described in [Fresh Installation > Using Helm](#using-helm) and then update the settings using
```
helm upgrade longhorn longhorn/longhorn --namespace longhorn-system --values ./values.yaml --version `helm list -n longhorn-system -o json | jq -r .'[0].app_version'`
```

## Upgrade

### Using the Rancher UI

From the project view in Rancher, go to **Apps && Marketplace > Longhorn > Upgrade > Next > Edit Options > Longhorn Default Settings > Customize Default Settings** and edit the settings before upgrading the app.
### Using the Longhorn Deployment YAML File

Modify the config map named `longhorn-default-setting` in the yaml file `longhorn/deploy/longhorn.yaml` as described in [Fresh Installation > Using the Longhorn Deployment YAML File](#using-the-longhorn-deployment-yaml-file) and then upgrade the Longhorn system using `kubectl`.

### Using Helm

Modify the default settings in the YAML file as described in [Fresh Installation > Using Helm](#using-helm) and then upgrade the Longhorn system using `helm upgrade`.

## History
Available since v1.3.0 ([Reference](https://github.com/longhorn/longhorn/issues/2570))


---

## Article: advanced-resources/deploy/node-selector.md

---
title: Node Selector
weight: 4
---

If you want to restrict Longhorn components to only run on a particular set of nodes, you can set node selector for all Longhorn components.
For example, you want to install Longhorn in a cluster that has both Linux nodes and Windows nodes but Longhorn cannot run on Windows nodes.
In this case, you can set the node selector to restrict Longhorn to only run on Linux nodes.

For more information about how node selector work, refer to the [official Kubernetes documentation.](https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#nodeselector)

# Setting up Node Selector for Longhorn
Longhorn consists of user-deployed components (for example, Longhorn Manager, Longhorn Driver, and Longhorn UI) and system-managed components (for example, Instance Manager, Backing Image Manager, Share Manager, CSI Driver, and Engine Image).
You need to set node selector for both types of components. See more details below.

### Setting up Node Selector During installing Longhorn
1. Set the node selector for user-deployed components (for example, Longhorn Manager, Longhorn Driver, and Longhorn UI).
   * If you install Longhorn through Rancher, you must copy and paste the following parameters into the YAML on the Rancher UI (click **Edit as YAML** during the installation) to apply the value to all user-deployed components.
      ```yaml
        global:
          nodeSelector:
            label-key1: "label-value1"
      ```
   * You can also specify the node selector for each user-deployed component and it will orverride the global setting.
      ```yaml
        longhornManager:
          nodeSelector:
            label-key1: "label-value1"
        longhornDriver:
          nodeSelector:
            label-key1: "label-value1"
        longhornUI:
          nodeSelector:
            label-key1: "label-value1"
      ```
   * If you install Longhorn by using `kubectl` to apply [the deployment YAML](https://raw.githubusercontent.com/longhorn/longhorn/v{{< current-version >}}/deploy/longhorn.yaml), you need to modify the node selector section for Longhorn Manager, Longhorn UI, and Longhorn Driver Deployer.
    Then apply the YAMl files.
   * If you install Longhorn using Helm, you can change the Helm values for `global.nodeSelector`, `longhornManager.nodeSelector`, `longhornUI.nodeSelector`, `longhornDriver.nodeSelector` in the `values.yaml` file before installing the chart.

2. Set the node selector for system-managed components (for example, Instance Manager, Backing Image Manager, Share Manager, CSI Driver, and Engine Image).

   Follow the [Customize default settings](../customizing-default-settings/) to set node selector by changing the value for the `system-managed-components-node-selector` default setting
   > Note: Because of the limitation of Rancher 2.5.x, if you are using Rancher UI to install Longhorn, you need to click `Edit As Yaml` and add setting `systemManagedComponentsNodeSelector` to `defaultSettings`.
   >
   > For example:
   > ```yaml
   > defaultSettings:
   >   systemManagedComponentsNodeSelector: "label-key1:label-value1"
   >  ```

### Setting up Node Selector After Longhorn has been installed

> **Warning**:
> * Since all Longhorn components will be restarted, the Longhorn system is unavailable temporarily.
> * When all Longhorn volumes are detached, the customized settings are immediately applied to the system-managed components (for example, Instance manager, CSI driver and Engine images).
> * When one or more Longhorn volumes are still attached, the customized setting is applied to the Instance Manager only when no engines and replica instances are running. You are required to reconfigure the setting after detaching the remaining volumes. Alternatively, you can wait for the next setting synchronization, which will occur in an hour.
> * Don't operate the Longhorn system while node selector settings are updated and Longhorn components are being restarted.

1. Prepare
   * To ensure that your preferred settings are immediately applied, stop all workloads and detach all Longhorn volumes before applying it.

2. Set the node selector for user-deployed components (for example, Longhorn Manager, Longhorn Driver, and Longhorn UI).
   * If you install Longhorn through Rancher, you must copy and paste the following parameters into the YAML on the Rancher UI (click **Edit as YAML** during the upgrade) to apply the value to all user-deployed components.
        ```yaml
        global:
          nodeSelector:
            label-key1: "label-value1"
        ```
    * You can also specify the node selector for each user-deployed component and it will override the global setting.
        ```yaml
        longhornManager:
          nodeSelector:
            label-key1: "label-value1"
        longhornDriver:
          nodeSelector:
            label-key1: "label-value1"
        longhornUI:
          nodeSelector:
            label-key1: "label-value1"
        ```
    * If you install Longhorn by using `kubectl` to apply [the deployment YAML](https://raw.githubusercontent.com/longhorn/longhorn/v{{< current-version >}}/deploy/longhorn.yaml), you need to modify the node selector section for Longhorn Manager, Longhorn UI, and Longhorn Driver Deployer.
      Then reapply the YAMl files.
    * If you install Longhorn using Helm, you can change the Helm values for `global.nodeSelector`, `longhornManager.nodeSelector`, `longhornUI.nodeSelector`, `longhornDriverDeployer.nodeSelector` in the `values.yaml` file, and then run `helm upgrade` to upgrade to the new version of the chart.

3. Set the node selector for system-managed components (for example, Instance Manager, Backing Image Manager, Share Manager, CSI Driver, and Engine Image).

   The node selector setting can be found at Longhorn UI under **Settings > System Managed Components Node Selector.**

4. Clean up

   If you are changing node selector in a way so that Longhorn cannot run on some nodes that Longhorn is currently running on,
   those nodes will become `down` state after this process. Verify that there is no replica left on those nodes.
   Disable scheduling for those nodes, and delete them in Longhorn UI

## History
Available since v1.1.1
* [Original feature request](https://github.com/longhorn/longhorn/issues/2199)


---

## Article: advanced-resources/deploy/priority-class.md

---
title: Priority Class
weight: 6
---
The Priority Class setting can be used to set a higher priority on Longhorn workloads in the cluster, preventing them from being the first to be evicted during node pressure situations.

For more information on how pod priority works, refer to the [official Kubernetes documentation](https://kubernetes.io/docs/concepts/configuration/pod-priority-preemption/).

# Setting Priority Class

Longhorn consists of user-deployed components (for example, Longhorn Manager, Longhorn Driver, and Longhorn UI) and system-managed components (for example, Instance Manager, CSI Driver, and Engine images).
You need to set Priority Class for both types of components. See more details below.

### Setting Priority Class During Longhorn Installation

Longhorn creates a Priority Class `longhorn-critical` and sets it as default for its user deployed or system managed components if the following actions are not taken.

1. Set taint Priority Class for system managed components: follow the [Customize default settings](../customizing-default-settings/) to set Priority Class by changing the value for the `priority-class` default setting
1. Set taint Priority Class for user deployed components: modify the Helm chart or deployment YAML file depending on how you deploy Longhorn.

> **Warning:** Longhorn will not start if the Priority Class setting is invalid (such as the Priority Class not existing).
> You can see if this is the case by checking the status of the longhorn-manager DaemonSet with `kubectl -n longhorn-system describe daemonset.apps/longhorn-manager`.
> You will need to uninstall Longhorn and restart the installation if this is the case.

### Setting Priority Class After Longhorn Installation

1. Set taint Priority Class for system managed components: The Priority Class setting can be found in the Longhorn UI by clicking **Settings > Priority Class.**
1. Set taint Priority Class for user deployed components: modify the Helm chart or deployment YAML file depending on how you deploy Longhorn.

Users can update or remove the Priority Class here, but note that this will result in recreation of all the Longhorn system components.
The Priority Class setting will reject values that appear to be invalid Priority Classes.

# Usage

To ensure that your preferred Priority Class settings are immediately applied, stop all workloads and detach all Longhorn volumes before configuring the settings.

Longhorn temporarily becomes unavailable when all components are restarted.
Don't operate the Longhorn system after modifying the Priority Class setting, as the Longhorn components will be restarting.

When all Longhorn volumes are detached, the customized setting is immediately applied to the system-managed components.
When one or more Longhorn volumes are still attached, the customized setting is applied to the Instance Manager only when no engines and replica instances are running. You are required to reconfigure the setting after detaching the remaining volumes. Alternatively, you can wait for the next setting synchronization, which will occur in an hour.

Do not delete the Priority Class in use by Longhorn, as this can cause new Longhorn workloads to fail to come online.

## History

[Original Feature Request](https://github.com/longhorn/longhorn/issues/1487)

Available since v1.0.1


---

## Article: advanced-resources/deploy/rancher_windows_cluster.md

---
title: Rancher Windows Cluster
weight: 5
---

Rancher can provision a Windows cluster with combination of Linux worker nodes and Windows worker nodes.
For more information on the Rancher Windows cluster, see the official [Rancher documentation](https://rancher.com/docs/rancher/v2.x/en/cluster-provisioning/rke-clusters/windows-clusters/).

In a Rancher Windows cluster, all Linux worker nodes are:
- Tainted with the taint `cattle.io/os=linux:NoSchedule`
- Labeled with `kubernetes.io/os:linux`

Follow the below [Deploy Longhorn With Supported Helm Chart](#deploy-longhorn-with-supported-helm-chart) or [Setup Longhorn Components For Existing Longhorn](#setup-longhorn-components-for-existing-longhorn) to know how to deploy or setup Longhorn on a Rancher Windows cluster.

> **Note**: After Longhorn is deployed, you can launch workloads that use Longhorn volumes only on Linux nodes.

## Deploy Longhorn With Supported Helm Chart
You can update the Helm value `global.cattle.windowsCluster.enabled` to allow Longhorn installation on the Rancher Windows cluster.

When this value is set to `true`, Longhorn will recognize the Rancher Windows cluster then deploy Longhorn components with the correct node selector and tolerations so that all Longhorn workloads can be launched on Linux nodes only.

On the Rancher marketplace, the setting can be customized in `customize Helm options` before installation: \
`Edit Options` > `Other Settings` > `Rancher Windows Cluster`

Also in: \
`Edit YAML`
```
global:
  cattle:
    systemDefaultRegistry: ""
    windowsCluster:
      # Enable this to allow Longhorn to run on the Rancher deployed Windows cluster
      enabled: true
```

## Setup Longhorn Components For Existing Longhorn
You can setup the existing Longhorn when its not deployed with the supported Helm chart.

1. Since Longhorn components can only run on Linux nodes,
   you need to set node selector `kubernetes.io/os:linux` for Longhorn to select the Linux nodes.
   Please follow the instruction at [Node Selector](../node-selector) to set node selector for Longhorn.

1. Since all Linux worker nodes in Rancher Windows cluster are tainted with the taint `cattle.io/os=linux:NoSchedule`,
   You need to set the toleration `cattle.io/os=linux:NoSchedule` for Longhorn to be able to run on those nodes.
   Please follow the instruction at [Taint Toleration](../taint-toleration) to set toleration for Longhorn.


---

## Article: advanced-resources/deploy/revision_counter.md

---
title: Revision Counter
weight: 7
---

The revision counter is a mechanism that Longhorn uses to track each replica's updates.

During replica creation, Longhorn will create a 'revision.counter' file with its initial counter set to 0. And for every write to the replica, the counter in 'revision.counter' file will be increased by 1.

The Longhorn engine uses these counters as a heuristic for achieving best-effort consistency among replicas during startup. Note that because the write IOs in Longhorn are parallel, enabling the revision counter does not guarantee data consistency. Longhorn also uses these counters during auto-salvage to identify the replica with the latest update.

Disable Revision Counter is an option in which every write on replicas is not tracked. When this setting is used, performance is improved. This option can be helpful if you prefer higher performance and have a stable network infrastructure (e.g. an internal network) with enough CPU resources. When the revision counter is disabled, the Longhorn Engine skips checking the revision counter for all replicas at startup. However, auto-salvage still functions because Longhorn can use the replica's head file stat to identify the replica to be used for recovery. For more information about how auto-salvage functions without the revision counter, see [Auto-Salvage Support with Revision Counter Disabled](#auto-salvage-support-with-revision-counter-disabled).

By default, the revision counter is disabled.

> **Note:** 'Salvage' is Longhorn trying to recover a volume in a faulted state. A volume is in a faulted state when the Longhorn Engine loses the connection to all the replicas, and all replicas are marked as being in an error state.

# Disable Revision Counter

## Using Longhorn UI

To disable or enable the revision counter from the Longhorn UI, click **Settings > Disable Revision Counter**.

To create individual volumes with settings that are customized against the general settings, go to the **Volumes** page and click **Create Volume**.

## Using a Manifest File

A `StorageClass` can be customized to add a `disableRevisionCounter` parameter.

By default, the `disableRevisionCounter` is false, so the revision counter is enabled.

Set `disableRevisionCounter` to true to disable the revision counter:

```yaml
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: best-effort-longhorn
provisioner: driver.longhorn.io
allowVolumeExpansion: true
parameters:
  numberOfReplicas: "1"
  disableRevisionCounter: "true"
  staleReplicaTimeout: "2880" # 48 hours in minutes
  fromBackup: ""
```

## Auto-Salvage Support with Revision Counter Disabled

The logic for auto-salvage is different when the revision counter is disabled.

When revision counter is enabled and all the replicas in the volume are in the 'ERR' state, the engine controller will be in a faulted state, and for engine to recover the volume, it will get the replica with the largest revision counter as 'Source of Truth' to rebuild the rest replicas.

When the revision counter is disabled in this case, the engine controller will get the `volume-head-xxx.img` last modified time and head file size of all replicas. It will also do the following steps:
1. Identify the replica with the most recent last modified timestamp based on when `volume-head-xxx.img` was last modified
1. Select all replicas with last modified timestamp within 5s of the above replica's last modified timestamp
2. From the replica candidates from the above step, compare the head file size of the candidates, and pick the ones with the largest file size
1. From the replica candidates from the above step, pick the best replica with most recent modified timestamp
3. Change the best replica to 'RW' mode, and the other replicas are marked as 'ERR' mode. The errored replicas are rebuilt based on the best replica


---

## Article: advanced-resources/deploy/storage-network.md

---
title: Storage Network
weight: 8
---

By Default, Longhorn uses the default Kubernetes cluster CNI network that is limited to a single interface and shared with other workloads cluster-wide. In case you have a situation where network segregation is needed, Longhorn supports isolating Longhorn in-cluster data traffic with the Storage Network setting.

The Storage Network setting takes Multus NetworkAttachmentDefinition in `<NAMESPACE>/<NAME>` format.

You can refer to [Comprehensive Document](https://github.com/k8snetworkplumbingwg/multus-cni#comprehensive-documentation) for how to install and set up Multus NetworkAttachmentDefinition.

Applying the setting will add `k8s.v1.cni.cncf.io/networks` annotation and recreate all existing instance-manager, and backing-image-manager pods.
Longhorn will apply the same annotation to any new instance-manager, backing-image-manager, and backing-image-data-source pods.

> **Important**: To ensure that your preferred settings are immediately applied, stop all workloads and detach all Longhorn volumes before configuring the settings.
>
> When all volumes are detached, Longhorn attempts to restart all Instance Manager and Backing Image Manager pods to apply the setting.
> When one or more Longhorn volumes are still attached, the customized setting is applied to the Instance Manager only when no engines and replica instances are running. You are required to reconfigure the setting after detaching the remaining volumes. Alternatively, you can wait for the next setting synchronization, which will occur in an hour.

# Setting Storage Network

## Prerequisite

The Multus NetworkAttachmentDefinition network for the storage network setting must be reachable in pods across different cluster nodes.

You can verify by creating a simple DaemonSet and try ping between pods.

### Setting Storage Network During Longhorn Installation
Follow the [Customize default settings](../customizing-default-settings/) to set Storage Network by changing the value for the `storage-network` default setting

> **Warning:** Longhorn instance-manager will not start if the Storage Network setting is invalid.
>
> You can check the events of the instance-manager Pod to see if it is related to an invalid NetworkAttachmentDefinition with `kubectl -n longhorn-system describe pods -l longhorn.io/component=instance-manager`.
>
> If this is the case, provide a valid `NetworkAttachmentDefinition` and re-run Longhorn install.

### Setting Storage Network After Longhorn Installation

Set the setting [Storage Network](../../../references/settings#storage-network).

> **Warning:** Do not modify the NetworkAttachmentDefinition custom resource after applying it to the setting.
>
> Longhorn is not aware of the updates. Hence this will cause malfunctioning and error. Instead, you can create a new NetworkAttachmentDefinition custom resource and update it to the setting.

### Setting Storage Network For RWX Volumes

Configure the setting [Storage Network For RWX Volume Enabled](../../../references/settings#storage-network-for-rwx-volume-enabled).

# Limitation

When an RWX volume is created with the storage network, the NFS mount point connection must be re-established when the CSI plugin pod restarts. Longhorn provides the [Automatically Delete Workload Pod when The Volume Is Detached Unexpectedly](../../../references/settings#automatically-delete-workload-pod-when-the-volume-is-detached-unexpectedly) setting, which automatically deletes RWX volume workload pods when the CSI plugin pod restarts. However, the workload pod's NFS mount point could become unresponsive when the setting is disabled or the pod is not managed by a controller. In such cases, you must manually restart the CSI plugin pod.

For more information, see [Storage Network Support for Read-Write-Many (RWX) Volume](../../../../archives/1.7.0/important-notes/#storage-network-support-for-read-write-many-rwx-volumes) in Important Notes.

# History
- [Original Feature Request (since v1.3.0)](https://github.com/longhorn/longhorn/issues/2285)
- [[FEATURE] Support storage network for RWX volumes (since v1.7.0)](https://github.com/longhorn/longhorn/issues/8184)
- [[FEATURE] Storage network with V2 data engine (since v1.9.0)](https://github.com/longhorn/longhorn/issues/6450)


---

## Article: advanced-resources/deploy/taint-toleration.md

---
title: Taints and Tolerations
weight: 3
---

If users want to create nodes with large storage spaces and/or CPU resources for Longhorn only (to store replica data) and reject other general workloads, they can taint those nodes and add tolerations for Longhorn components. Then Longhorn can be deployed on those nodes.

Notice that the taint tolerations setting for one workload will not prevent it from being scheduled to the nodes that don't contain the corresponding taints.

For more information about how taints and tolerations work, refer to the [official Kubernetes documentation.](https://kubernetes.io/docs/concepts/configuration/taint-and-toleration/)

# Setting up Taints and Tolerations
Longhorn consists of user-deployed components (for example, Longhorn Manager, Longhorn Driver, and Longhorn UI) and system-managed components (for example, Instance Manager, Backing Image Manager, Share Manager, CSI Driver, and Engine Image).
You need to set tolerations for both types of components. See more details below.

### Setting up Taints and Tolerations During installing Longhorn
1. Set taint tolerations for user-deployed components (for example, Longhorn Manager, Longhorn Driver, and Longhorn UI).
   * If you install Longhorn through Rancher, you must copy and paste the following parameters into the YAML on the Rancher UI (click **Edit as YAML** during the installation) to apply the value to all user-deployed components.
     ```yaml
       global:
         tolerations:
         - key: "key"
           operator: "Equal"
           value: "value"
           effect: "NoSchedule"
    ```
   * You can also specify the tolerations for each user-deployed component and it will override the global setting.
     ```yaml
       longhornManager:
         tolerations:
         - key: "key"
           operator: "Equal"
           value: "value"
           effect: "NoSchedule"
       longhornDriver:
         tolerations:
         - key: "key"
           operator: "Equal"
           value: "value"
           effect: "NoSchedule"
       longhornUI:
         tolerations:
         - key: "key"
           operator: "Equal"
           value: "value"
           effect: "NoSchedule"
     ```
   * If you install Longhorn by using `kubectl` to apply [the deployment YAML](https://raw.githubusercontent.com/longhorn/longhorn/v{{< current-version >}}/deploy/longhorn.yaml), you need to modify the taint tolerations section for Longhorn Manager, Longhorn UI, and Longhorn Driver Deployer.
    Then apply the YAMl files.
   * If you install Longhorn using Helm, you can change the Helm values for `global.tolerations`, `longhornManager.tolerations`, `longhornUI.tolerations`, `longhornDriver.tolerations` in the `values.yaml` file before installing the chart.

2. Set taint tolerations for system-managed components (for example, Instance Manager, CSI Driver, and Engine images)

   Follow the [Customize default settings](../customizing-default-settings/) to set taint tolerations by changing the value for the `taint-toleratio` default setting
   > Note: Because of the limitation of Rancher 2.5.x, if you are using Rancher UI to install Longhorn, you need to click `Edit As Yaml` and add setting `taintToleration` to `defaultSettings`.
   >
   > For example:
   > ```yaml
   > defaultSettings:
   >   taintToleration: "key=value:NoSchedule"
   >  ```

### Setting up Taints and Tolerations After Longhorn has been installed

> **Warning**:
>
> To ensure that your preferred toleration settings are immediately applied, stop all workloads and detach all Longhorn volumes before configuring the settings.
>
> Since all Longhorn components will be restarted, the Longhorn system is unavailable temporarily.
>
> When all Longhorn volumes are detached, the customized setting is immediately applied to the system-managed components.
> When one or more Longhorn volumes are still attached, the customized setting is applied to the Instance Manager only when no engines and replica instances are running. You are required to reconfigure the setting after detaching the remaining volumes. Alternatively, you can wait for the next setting synchronization, which will occur in an hour.
>
> Don't operate the Longhorn system while toleration settings are updated and Longhorn components are being restarted.

1. Prepare

   To ensure that your preferred settings are immediately applied, stop all workloads and detach all Longhorn volumes before configuring the settings.

2. Set taint tolerations for user-deployed components (for example, Longhorn Manager, Longhorn Driver, and Longhorn UI).
   * If you install Longhorn through Rancher, you must copy and paste the following parameters into the YAML on the Rancher UI (click **Edit as YAML** during the upgrade) to apply the value to all user-deployed components.
     ```yaml
       global:
         tolerations:
         - key: "key"
           operator: "Equal"
           value: "value"
           effect: "NoSchedule"
    ```
   * You can also specify the tolerations for each user-deployed component and it will override the global setting.
     ```yaml
       longhornManager:
         tolerations:
         - key: "key"
           operator: "Equal"
           value: "value"
           effect: "NoSchedule"
       longhornDriver:
         tolerations:
         - key: "key"
           operator: "Equal"
           value: "value"
           effect: "NoSchedule"
       longhornUI:
         tolerations:
         - key: "key"
           operator: "Equal"
           value: "value"
           effect: "NoSchedule"
     ```
   * If you install Longhorn by using `kubectl` to apply [the deployment YAML](https://raw.githubusercontent.com/longhorn/longhorn/v{{< current-version >}}/deploy/longhorn.yaml), you need to modify the taint tolerations section for Longhorn Manager, Longhorn UI, and Longhorn Driver Deployer.
  Then reapply the YAMl files.
   * If you install Longhorn using Helm, you can change the Helm values for `global.tolerations`, `longhornManager.tolerations`, `longhornUI.tolerations`, `longhornDriver.tolerations` in the `values.yaml` file, and then run `helm upgrade` to upgrade to the new version of the chart.

3. Set taint tolerations for system-managed components (for example, Instance Manager, Backing Image Manager, Share Manager, CSI Driver, and Engine Image).

   The taint toleration setting can be found at Longhorn UI under **Settings > Kubernetes Taint Toleration.**


## History
Available since v0.6.0
* [Original feature request](https://github.com/longhorn/longhorn/issues/584)
* [Resolve the problem with GitOps](https://github.com/longhorn/longhorn/issues/2120)


---

## Article: advanced-resources/system-backup-restore/_index.md

---
title: Longhorn System Backup And Restore
weight: 10
---

> Before v1.4.0, you can restore Longhorn with third-party tools.

- [Restore to a cluster contains data using Rancher snapshot](./restore-to-a-cluster-contains-data-using-rancher-snapshot)
- [Restore to a new cluster using Velero](./restore-to-a-new-cluster-using-velero)

> Since v1.4.0, Longhorn introduced out-of-the-box Longhorn system backup and restore.
> - Longhorn's custom resources will be backed up and bundled into a single system backup file, then saved to the remote backup target.
> - Later, you can choose a system backup to restore to a new cluster or restore to an existing cluster.

- [Backup Longhorn system](./backup-longhorn-system)
- [Restore Longhorn system](./restore-longhorn-system)


---

## Article: advanced-resources/system-backup-restore/backup-longhorn-system.md

---
title: Backup Longhorn System
weight: 1
---

- [Longhorn System Backup Bundle](#longhorn-system-backup-bundle)
- [Create Longhorn System Backup](#create-longhorn-system-backup)
  - [Prerequisite](#prerequisite)
  - [Configuration](#configuration)
    - [Volume Backup Policy](#volume-backup-policy)
  - [Single Execution](#single-execution)
    - [Create a System Backup Using the Longhorn UI](#create-a-system-backup-using-the-longhorn-ui)
    - [Create a System Backup Using `kubectl`](#create-a-system-backup-using-kubectl)
  - [Recurring Job](#recurring-job)
    - [Create a Recurring Backup Job Using the Longhorn UI](#create-a-recurring-backup-job-using-the-longhorn-ui)
    - [Create a Recurring Backup Job Using `kubectl`](#create-a-recurring-backup-job-using-kubectl)
- [Delete Longhorn System Backup](#delete-longhorn-system-backup)
  - [Delete a System Backup Using the Longhorn UI](#delete-a-system-backup-using-the-longhorn-ui)
  - [Delete a System Backup Using `kubectl`](#delete-a-system-backup-using-kubectl)
- [History](#history)

## Longhorn System Backup Bundle

Longhorn system backup creates a resource bundle and uploads it to the remote backup target.

It includes below resources associating with the Longhorn system:
- BackingImages
- ClusterRoles
- ClusterRoleBindings
- ConfigMaps
- CustomResourceDefinitions
- DaemonSets
- Deployments
- EngineImages
- PersistentVolumes
- PersistentVolumeClaims
- RecurringJobs
- Roles
- RoleBindings
- Settings
- Services
- ServiceAccounts
- StorageClasses
- Volumes

> **Note:**
>
> - The default backup target (`default`) is always used to store system backups.
> - The Longhorn system backup bundle only includes resources operated by Longhorn.
> - Longhorn does not back up the `Nodes` resource. The Longhorn Manager on the target cluster is responsible for creating its own Longhorn `Node` custom resources.
> - Longhorn is unable to back up V2 Data Engine backing images.
>
> Here is an example of a cluster workload with a bare `Pod` workload. The system backup will collect the `PersistentVolumeClaim`, `PersistentVolume`, and `Volume`. The system backup will exclude the `Pod` during system backup resource collection.

## Create Longhorn System Backup

You can create a Longhorn system backup using the Longhorn UI. Or with the `kubectl` command.

### Prerequisite

- [Set the backup target](../../../snapshots-and-backups/backup-and-restore/set-backup-target). Longhorn saves the system backups to the remote backup store. You will see an error during creation when the backup target is unset.

   > **Note:** Unsetting the backup target clears the existing `SystemBackup` custom resource. Longhorn syncs to the remote backup store after setting the backup target. Another cluster can also sync to the same list of system backups when the backup target is the same.

- Create a backup for all volumes (optional).

  > **Note:** Longhorn system restores volume with the latest backup. We recommend updating the last backup for all volumes. By taking volume backups, you ensure that the data is up-to-date with the system backup. For more information, please refer to the [Configuration - Volume Backup Policy](#volume-backup-policy) section.

### Configuration

#### Volume Backup Policy
The Longhorn system backup offers the following volume backup policies:
 - `if-not-present`: Longhorn will create a backup for volumes that either lack an existing backup or have an outdated latest backup.
 - `always`: Longhorn will create a backup for all volumes, regardless of their existing backups.
 - `disabled`: Longhorn will not create any backups for volumes.

### Single Execution

#### Create a System Backup Using the Longhorn UI

1. Go to the `System Backups` page in the `Backup and Restore` drop-down list.

1. Click `Create` under `System Backup`.

1. Give a `Name` for the system backup.

1. Select a `Volume Backup Policy` for the system backup.

1. The system backup will be ready to use when the state changes to `Ready`.

#### Create a System Backup Using `kubectl`

1. Execute `kubectl create` to create a Longhorn `SystemBackup` custom resource.
   ```yaml
   apiVersion: longhorn.io/v1beta2
   kind: SystemBackup
   metadata:
     name: demo
     namespace: longhorn-system
   spec:
     volumeBackupPolicy: if-not-present
   ```

1. The system backup will be ready to use when the state changes to `Ready`.
   ```
   > kubectl -n longhorn-system get systembackup
   NAME   VERSION   STATE   CREATED
   demo   v1.4.0    Ready   2022-11-24T04:23:24Z
   ```

### Recurring Job

#### Create a Recurring Backup Job Using the Longhorn UI

1. Go to the `Recurring Jobs` page.

1. Click on `Create Recurring Job`.

1. Configure the following settings:
   - **Name**: Specify a name for the recurring job.
   - **Task**: Select **System Backup**.
   - **Retain**: Specify the number of system backups that Longhorn must retain.
   - **Cron**: Specify the cron expression (a string consisting of fields separated by whitespace characters) that defines the schedule properties.
   - **Parameters**: Select **volume-backup-policy**.

1. Click **OK**.

Longhorn creates system backups according to the schedule defined in the **Cron** field.

#### Create a Recurring Backup Job Using `kubectl`

Run `kubectl create` to create a Longhorn `RecurringJob` custom resource with the task `system-backup`.

Example:
   ```yaml
   apiVersion: longhorn.io/v1beta2
   kind: RecurringJob
   metadata:
     name: demo
     namespace: longhorn-system
   spec:
     task: system-backup
     cron: '* * * * *'
     retain: 1
     parameters:
       volume-backup-policy: if-not-present
   ```

Longhorn creates system backup according to the schedule defined in the `cron` field.

## Delete Longhorn System Backup

You can delete the Longhorn system backup in the remote backup target using the Longhorn UI. Or with the `kubectl` command.

### Delete a System Backup Using the Longhorn UI

1. Go to the `System Backup` page in the `Setting` drop-down list.

1. Delete a single system backup in the `Operation` drop-down menu next to the system backup. Or delete in batch with the `Delete` button.

   > **Note:** Deleting the system backup will also make a deletion in the backup store.

### Delete a System Backup Using `kubectl`

1. Execute `kubectl delete` to delete a Longhorn `SystemBackup` custom resource.
   ```
   > kubectl -n longhorn-system get systembackup
   NAME   VERSION   STATE   CREATED
   demo   v1.4.0    Ready   2022-11-24T04:23:24Z

   > kubectl -n longhorn-system delete systembackup/demo
   systembackup.longhorn.io "demo" deleted
   ```

## History
[Original Feature Request](https://github.com/longhorn/longhorn/issues/1455)

Available since v1.4.0


---

## Article: advanced-resources/system-backup-restore/restore-longhorn-system.md

---
title: Restore Longhorn System
weight: 2
---

- [What does the Longhorn system restore rollout to the cluster](#longhorn-system-restore-rollouts)
- [What are the limitations](#limitations)
    - [Restore Path](#restore-path)
- [How to restore from Longhorn system backup](#create-longhorn-system-restore)
    - [Prerequisite](#prerequisite)
    - [Using Longhorn UI](#using-longhorn-ui)
    - [Using kubectl command](#using-kubectl-command)
- [How to delete Longhorn system restore](#delete-longhorn-system-restore)
    - [Using Longhorn UI](#using-longhorn-ui-1)
    - [Using kubectl command](#using-kubectl-command-1)
- [How to restart Longhorn System Restore](#restart-longhorn-system-restore)
- [What settings are configurable](#configurable-settings)
- [How to troubleshoot](#troubleshoot)
- [History](#history)

## Longhorn System Restore Rollouts

- Longhorn restores the resource from the [Longhorn System Backup Bundle](../backup-longhorn-system#longhorn-system-backup-bundle).
- Longhorn does not restore existing `Volumes` and their associated `PersistentVolume` and `PersistentVolumeClaim`.
- Longhorn automatically restores a `Volume` from its latest backup.
- To prevent overwriting eligible settings, Longhorn does not restore the `ConfigMap/longhorn-default-setting`.
- Longhorn does not restore [configurable settings](#configurable-settings).
- Since Longhorn does not back up V2 Data Engine backing images, you must ensure that those images are available in the cluster before you restore the Longhorn system. This allows Longhorn to restore volumes that use V2 Data Engine backing images.

## Limitations
### Restore Path

Longhorn does not support cross-major/minor version system restore except for upgrade failures, ex: 1.4.x -> 1.5.
## Create Longhorn System Restore

You can restore the Longhorn system using Longhorn UI. Or with the `kubectl` command.

### Prerequisite

- A running Longhorn cluster for Longhorn to roll out the resources in the system backup bundle.
- Set up the `Nodes` and disk tags for `StorageClass`.
- Have a Longhorn system backup.

  See [Backup Longhorn System - Create Longhorn System Backup](../backup-longhorn-system#create-longhorn-system-backup) for instructions.
- All existing `Volumes` are detached.

### Using Longhorn UI

1. Go to the `System Backups` page in the `Backup and Restore`.
1. Select a system backup to restore.
1. Click `Restore` in the `Operation` drop-down menu.
1. Give a `Name` for the system restore.
1. The system restore starts and show the `Completed` state when done.

## Using `kubectl` Command

1. Find the Longhorn `SystemBackup` to restore.
   ```
   > kubectl -n longhorn-system get systembackup
   NAME     VERSION   STATE   CREATED
   demo     v1.4.0    Ready   2022-11-24T04:23:24Z
   demo-2   v1.4.0    Ready   2022-11-24T05:00:59Z
   ```
1. Execute `kubectl create` to create a Longhorn `SystemRestore` of the `SystemBackup`.
   ```yaml
   apiVersion: longhorn.io/v1beta2
   kind: SystemRestore
   metadata:
     name: restore-demo
     namespace: longhorn-system
   spec:
     systemBackup: demo
   ```
1. The system restore starts.
1. The `SystemRestore` change to state `Completed` when done.
   ```
   > kubectl -n longhorn-system get systemrestore
   NAME           STATE       AGE
   restore-demo   Completed   59s
   ```

## Delete Longhorn System Restore

> **Warning:** Deleting the SystemRestore also deletes the associated job and will abort the remaining resource rollouts. You can [Restart the Longhorn System Restore](#restart-longhorn-system-restore) to roll out the remaining resources.

You can abort or remove a completed Longhorn system restore using Longhorn UI. Or with the `kubectl` command.

### Using Longhorn UI

1. Go to the `System Backups` page in the `Backup and Restore`.
1. Delete a single system restore in the `Operation` drop-down menu next to the system restore. Or delete in batch with the `Delete` button.

### Using `kubectl` Command

1. Execute `kubectl delete` to delete a Longhorn `SystemRestore`.
   ```
   > kubectl -n longhorn-system get systemrestore
   NAME           STATE       AGE
   restore-demo   Completed   2m37s
   
   > kubectl -n longhorn-system delete systemrestore/restore-demo
   systemrestore.longhorn.io "restore-demo" deleted
   ```

## Restart Longhorn System Restore

1. [Delete Longhorn System Restore](#delete-longhorn-system-restore) that is in progress.
1. [Create Longhorn System Restore](#create-longhorn-system-restore).

## Configurable Settings

Some settings are excluded as configurable before the Longhorn system restore.
- [Concurrent volume backup restore per node limit](../../../references/settings/#concurrent-volume-backup-restore-per-node-limit)
- [Concurrent replica rebuild per node limit](../../../references/settings/#concurrent-replica-rebuild-per-node-limit)

## Troubleshoot

### System Restore Hangs

1. Check the longhorn-system-rollout Pod log for any errors.
```
> kubectl -n longhorn-system logs --selector=job-name=longhorn-system-rollout-<SYSTEM-RESTORE-NAME>
```
1. Resolve if the issue is identifiable, ex: remove the problematic restoring resource.
1. [Restart the Longhorn system restore](#restart-longhorn-system-restore).

## History
[Original Feature Request](https://github.com/longhorn/longhorn/issues/1455)

Available since v1.4.0


---

## Article: advanced-resources/system-backup-restore/restore-to-a-cluster-contains-data-using-Rancher-snapshot.md

---
title: Restore to a cluster contains data using Rancher snapshot
weight: 4
---

This doc describes what users need to do after restoring the cluster with a Rancher snapshot.

## Assumptions:
- **Most of the data and the underlying disks still exist** in the cluster before the restore and can be directly reused then.
- There is a backupstore holding all volume data.
- The setting [`Disable Revision Counter`](../../../references/settings/#disable-revision-counter) is false. (It's false by default.) Otherwise, users need to manually check if the data among volume replicas are consistent, or directly restore volumes from backup.

## Expectation:
- All settings and node & disk configs will be restored.
- As long as the valid data still exists, the volumes can be recovered without using a backup. In other words, we will try to avoid restoring backups, which may help reduce Recovery Time Objective (RTO) as well as save bandwidth.
- Detect the invalid or out-of-sync replicas as long as the related volume still contains a valid replica after the restore.

## Behaviors & Requirement of Rancher restore
According to [the Rancher restore article](https://ranchermanager.docs.rancher.com/how-to-guides/new-user-guides/backup-restore-and-disaster-recovery/restore-rancher-launched-kubernetes-clusters-from-backup), you have to restart the Kubernetes components on all nodes. Otherwise, there will be tons of resource update conflicts in Longhorn.

## Actions after the restore
- Restart all Kubernetes components for all nodes. See the above link for more details.

- Kill all longhorn manager pods then Kubernetes will automatically restart them. Wait for conflicts in longhorn manager pods to disappear.

- All volumes may be reattached. If a Longhorn volume is used by a single pod, users need to shut down then recreate it. For Deployments or Statefulsets, Longhorn will automatically kill then restart the related pods.

- If the following happens after the snapshot and before the cluster restore:
    - A volume is unchanged: Users don't need to do anything.
    - The data is updated: Users don't need to do anything typically. Longhorn will automatically fail the replicas that don't contain the latest data.
    - A new volume is created: This volume will disappear after the restore. Users need to recreate a new volume, launch [a single replica volume](../../data-recovery/export-from-replica) based on the replica of the disappeared volume, then transfer the data to the new volume.
    - A volume is deleted: Since the data is cleaned up when the volume is removed, the restored volume contains no data. Users may need to re-delete it.
    - For DR volumes: Users don't need to do anything. Longhorn will redo a full restore.
    - Some operations are applied for a volume:
        - Backup: The backup info of the volume should be resynced automatically.
        - Snapshot: The snapshot info of the volume should be resynced once the volume is attached.
        - Replica rebuilding & replica removal:
            - If there are new replicas rebuilt, those replicas will disappear from the Longhorn system after the restoring. Users need to clean up the replica data manually, or use the data directories of these replicas to export a single replica volume then do data recovery if necessary.
            - If there are some failed/removed replicas and there is at least one replica keeping healthy, those failed/removed replicas will be back after the restoration. Then Longhorn can detect these restored replicas do not contain any data, and copy the latest data from the healthy replica to these replicas.
            - If all replicas are replaced by new replicas after the snapshot, the volume will contain invalid replicas only after the restore. Then users need to export [a single replica volume](../../data-recovery/export-from-replica) for the data recovery.
        -  Engine image upgrade: Users need to redo the upgrade.
        - Expansion: The spec size of the volume will be smaller than the current size. This is like someone requesting volume shrinking but actually Longhorn will refuse to handle it internally. To recover the volume, users need to scale down the workloads and re-do the expansion.

    - **Notice**: If users don't know how to recover a problematic volume, the simplest way is always restoring a new volume from backup.

- If the Longhorn system is upgraded after the snapshot, the new settings and the modifications on the node config will disappear. Users need to re-do the upgrade, then re-modify the settings and node configurations.

- If a node is deleted from Longhorn system after the snapshot, the node won't be back, but the pods on the removed node will be restored. Users need to manually clean up them since these pod may get stuck in state `Terminating`.
- If a node to added to Longhorn system after the snapshot, Longhorn should automatically relaunch all necessary workloads on the node after the cluster restore. But users should be aware that all new replicas or engines on this node will be gone after the restore.


## References
- The related GitHub issue is https://github.com/longhorn/longhorn/issues/2228.
  In this GitHub post, one user is providing a way that restores the Longhorn to a new cluster that doesn't contain any data.


---

## Article: advanced-resources/system-backup-restore/restore-to-a-new-cluster-using-velero.md

---
title: Restore to a new cluster using Velero
weight: 4
---

This doc instructs how users can restore workloads with Longhorn system to a new cluster via Velero.

> **Note:** Need to use [Velero CSI plugin](https://github.com/vmware-tanzu/velero-plugin-for-csi) >= 0.4 to ensure restoring PersistentVolumeClaim successfully. Visit [here](/kb/troubleshooting-restore-pvc-stuck-using-velero-csi-plugin-version-lower-than-0.4) to get more information.


## Assumptions:
- A new cluster means there is **no Longhorn volume data** in it.
- There is a remote backup target holds all Longhorn volume data.
- There is a remote backup server that can store the cluster backups created by Velero.

## Expectation:
- All settings will be restored. But the node & disk configurations won't be applied.
- All workloads using Longhorn volumes will get started after the volumes are restored from the remote backup target.

## Workflow

### Create backup for the old cluster
1. Install Velero into a cluster using Longhorn.
2. Create backups for all Longhorn volumes.
3. Use Velero to create a cluster backup. Here, some Longhorn resources should be excluded from the cluster backup:
    ```bash
    velero backup create lh-cluster --exclude-resources persistentvolumes,persistentvolumeclaims,backuptargets.longhorn.io,backupvolumes.longhorn.io,backups.longhorn.io,nodes.longhorn.io,volumes.longhorn.io,engines.longhorn.io,replicas.longhorn.io,backingimagedatasources.longhorn.io,backingimagemanagers.longhorn.io,backingimages.longhorn.io,sharemanagers.longhorn.io,instancemanagers.longhorn.io,engineimages.longhorn.io
    ```
### Restore Longhorn and workloads to a new cluster
1. Install Velero with the same remote backup sever for the new cluster.
2. Restore the cluster backup. e.g.,
    ```bash
    velero restore create --from-backup lh-cluster
    ```
3. Removing all old instance manager pods and backing image manager pods from namespace `longhorn-system`. These old pods should be created by Longhorn rather than Velero and there should be corresponding CRs for them. The pods are harmless but they would lead to the endless logs printed in longhorn-manager pods. e.g.,:
    ```log
    [longhorn-manager-q6n7x] time="2021-12-20T10:42:49Z" level=warning msg="Can't find instance manager for pod instance-manager-r-1f19ecb0, may be deleted"
    [longhorn-manager-q6n7x] time="2021-12-20T10:42:49Z" level=warning msg="Can't find instance manager for pod instance-manager-e-6c3be222, may be deleted"
    [longhorn-manager-ldlvw] time="2021-12-20T10:42:55Z" level=warning msg="Can't find instance manager for pod instance-manager-e-bbf80f76, may be deleted"
    [longhorn-manager-ldlvw] time="2021-12-20T10:42:55Z" level=warning msg="Can't find instance manager for pod instance-manager-r-3818fdca, may be deleted"
    ```
4. Re-config nodes and disks for the restored Longhorn system if necessary.
5. Re-create backing images if necessary.
6. Restore all Longhorn volumes from the remote backup target.
7. If there are RWX backup volumes, users need to manually update the access mode to `ReadWriteMany` since all restored volumes are mode `ReadWriteOnce` by default.
8. Create PVCs and PVs with previous names for the restored volumes.

Note: We will enhance Longhorn system so that users don't need to apply step3 and step8 in the future.

## References
- The related GitHub issue is https://github.com/longhorn/longhorn/issues/3367


---

## Article: advanced-resources/longhornctl/_index.md

---
title: Command Line Tool (longhornctl)
description: Command line interface (CLI) for Longhorn operations and troubleshooting.
weight: 8
---

The `longhornctl` tool is a CLI interface to Longhorn operations. It interacts with Longhorn by creating Kubernetes Custom Resources (CRs) and executing commands inside a dedicated Pod for in-cluster and host operations.

## Common usage scenarios

* **Installation:**
  * `longhornctl install preflight`: Perform preflight dependencies installation and setup before installing Longhorn.
* **Operations:**
  * `longhornctl export replica`: Extract data from a Longhorn replica data directory to a designated directory on its host machine. This is useful for recovering data when Longhorn is unavailable.
  * `longhornctl trim volume`: Reclaim unused storage space within a Longhorn volume.
* **Troubleshooting:**
  * `longhornctl check preflight`: Identifies potential issues before usage.
  * `longhornctl get replica`: Retrieve details about Longhorn replicas on the host.

## Usage

For more information about the available commands, see [this document](https://github.com/longhorn/cli/tree/v{{< current-version >}}/docs/longhornctl.md) in the GitHub repository, or run `longhornctl help` in your terminal.


---

## Article: advanced-resources/longhornctl/install-longhornctl.md

---
title: Install longhornctl
weight: 1
---

## Use the Prebuilt Binary

1. Download the binary:
   ```bash
   # Choose your architecture (amd64 or arm64).
   ARCH="amd64"

   # Download the release binary.
   curl -LO "https://github.com/longhorn/cli/releases/download/v{{< current-version >}}/longhornctl-linux-${ARCH}"
   ```
1. Validate the binary:
   ```bash
   # Download the checksum for your architecture.
   curl -LO "https://github.com/longhorn/cli/releases/download/v{{< current-version >}}/longhornctl-linux-${ARCH}.sha256"

   # Verify the downloaded binary matches the checksum.
   echo "$(cat longhornctl-linux-${ARCH}.sha256 | awk '{print $1}') longhornctl-linux-${ARCH}" | sha256sum --check
   ```
1. Install the binary:
   ```bash
   sudo install longhornctl-linux-${ARCH} /usr/local/bin/longhornctl
   ```
1. Verify installation:
   ```bash
   longhornctl version
   ```

## Build From Source

See [this document](https://github.com/longhorn/cli/tree/{{< current-version >}}?tab=readme-ov-file#build-from-source) in the GitHub repository.


---

## Article: advanced-resources/security/_index.md

---
title: Security
weight: 7
---


---

## Article: advanced-resources/security/mtls-support.md

---
title: MTLS Support
weight: 6
---

Longhorn supports MTLS to secure and encrypt the grpc communication between the control plane (longhorn-manager) and the data plane (instance-managers).
For Certificate setup we use the Kubernetes secret mechanism in combination with an optional secret mount for the longhorn-manager/instance-manager.


# Requirements
In a default installation mtls is disabled to enable mtls support one needs to create a `longhorn-grpc-tls` secret in the `longhorn-system` namespace before deployment.
The secret is specified as an optional secret mount for the longhorn-manager/instance-managers so if it does not exist when these
components are started, mtls will not be used and a restart of the components will be required to enable tls support.

The longhorn-manager has a non tls client fallback for mixed mode setups where there are old instance-managers that were started without tls support.

# Self Signed Certificate Setup

You should create a `ca.crt` with the CA flag set which is then used to sign the `tls.crt` this will allow you to rotate the `tls.crt` in the future without service interruptions.
You can use [openssl](https://mariadb.com/docs/server/security/securing-mariadb/securing-mariadb-encryption/data-in-transit-encryption/certificate-creation-with-openssl/)
or [cfssl](https://github.com/cloudflare/cfssl) for the `ca.crt` as well as `tls.crt` certificate generation.

The `tls.crt` certificate should use `longhorn-backend` for the common name and the below list of entries for the Subject Alternative Name.
```text
Common Name: longhorn-backend
Subject Alternative Names: longhorn-backend, longhorn-backend.longhorn-system, longhorn-backend.longhorn-system.svc, longhorn-frontend, longhorn-frontend.longhorn-system, longhorn-frontend.longhorn-system.svc, longhorn-engine-manager, longhorn-engine-manager.longhorn-system, longhorn-engine-manager.longhorn-system.svc, longhorn-replica-manager, longhorn-replica-manager.longhorn-system, longhorn-replica-manager.longhorn-system.svc, longhorn-csi, longhorn-csi.longhorn-system, longhorn-csi.longhorn-system.svc, longhorn-backend, IP Address:127.0.0.1
```

# Setting up Kubernetes Secrets

The `ca.crt` is the certificate of the certificate authority that was used to sign
the `tls.crt` which will be used both by the client (longhorn-manager) and the server (instance-manager) for grpc mtls authentication.
The `tls.key` is associated private key for the created `tls.crt`.

The `longhorn-grpc-tls` yaml looks like the below example,
If you are having trouble getting your own certificates to work you can base decode the below certificate
and compare it against your own generated certificates via `openssl x509 -in tls.crt -text -noout`.
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: longhorn-grpc-tls
  namespace: longhorn-system
type: kubernetes.io/tls
data:
  ca.crt: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUREakNDQWZhZ0F3SUJBZ0lVU2tNdUVEOC9XYXphNmpkb1NiTE1qalFqb3JFd0RRWUpLb1pJaHZjTkFRRUwKQlFBd0h6RWRNQnNHQTFVRUF4TVViRzl1WjJodmNtNHRaM0p3WXkxMGJITXRZMkV3SGhjTk1qSXdNVEV4TWpFeApPVEF3V2hjTk1qY3dNVEV3TWpFeE9UQXdXakFmTVIwd0d3WURWUVFERXhSc2IyNW5hRzl5YmkxbmNuQmpMWFJzCmN5MWpZVENDQVNJd0RRWUpLb1pJaHZjTkFRRUJCUUFEZ2dFUEFEQ0NBUW9DZ2dFQkFNY2grbTJhUndnNEtBa0EKT0xzdzdScWlWb1VqL2VPbVhuSE9HVE5nWE4rcFh5bDlCdzVDM1J4UDYzU29qaTVvNEhkU1htVmpwZmhmNjh1YwpvNVJJeUtXM1p6cndteDhXZldEc0dNNEtnYXBvMy84N3pVQ00vdGltOHllTzFUbTZlWVhXcWdlZ2JpM1Q1WnlvCmkzRjdteFg3QlU3Z25uWGthVmJ5UU1xRkEyMDJrK25jaVhaUE9iU0tlc1NvZ20wdWsrYXFvY3N1SjJ6dk9tZG0KMXd0a3ZTUklhL3l6T25JRGlmbFRteXNhZ3oxQy9VM1JxbzJ6TjIwbWJNYUJhMmx5anVZWkdWSnNyNGh4dGhqUApIR2x1UUh2QTlKTE9kc2J0T2xmbjRZNlZpUktCSzZWMVpOeVROMVJpN3ArTXZlaWQ3cE9rNHYweC9qVTc1a0N6Clo1cGJHbGtDQXdFQUFhTkNNRUF3RGdZRFZSMFBBUUgvQkFRREFnRUdNQThHQTFVZEV3RUIvd1FGTUFNQkFmOHcKSFFZRFZSME9CQllFRlBGc0xRbmQxOHFUTVd5djh1STk3Z2hnR2djR01BMEdDU3FHU0liM0RRRUJDd1VBQTRJQgpBUUNMcnk5a2xlSElMdDRwbzd4N0hvSldsMEswYjdwV2Y0Y3ZVeHh1bUdTYUpoQmFHNTVlZFNFSVAzajhsRGg1Cm94ZXJlbjNrRUtzeGZiQVQ0RzU3KzBaeExQSkZQcjFMM3JvcmxUVE1DS1QyY2Z1UDJ3SEIzZndWNDJpSHZSUDgKSUVqU041bFNkWjZnN1NjWFZ2RnpZNzlrbVZDQ2RNYlpGcEFuOElyTkh3L0tTUGZUajNob2VyV3ZGL3huaEo3bQpmSzUrcE5TeWR6QTA1K1Y0ODJhWGlvV2NWcWY2UHpSVndmT0tIalUrbUVDQXZMbitNSzRvN1l2VW1iN2tSUGs5CnBjU1A4N2lpN0hwRVhqZUtRaVJhZElXKzMySXp1UTFiOXRYc3BNTGF0UFA5TXNvWmY0M1EyZWw4bWd1RjRxOUcKVmVUZFZaU2hBNWNucmNRZTEySUs1MzAvCi0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K
  tls.crt: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUVqekNDQTNlZ0F3SUJBZ0lVUjZWcGR5U1Z0MGp6bDcwQnIxMmdZOTB0QVNBd0RRWUpLb1pJaHZjTkFRRUwKQlFBd0h6RWRNQnNHQTFVRUF4TVViRzl1WjJodmNtNHRaM0p3WXkxMGJITXRZMkV3SGhjTk1qSXdNVEV4TWpFeApPVEF3V2hjTk1qTXdNVEV4TWpFeE9UQXdXakFiTVJrd0Z3WURWUVFERXhCc2IyNW5hRzl5YmkxaVlXTnJaVzVrCk1Ga3dFd1lIS29aSXpqMENBUVlJS29aSXpqMERBUWNEUWdBRUxOQVVJZUsvdnppaGc3a1Q5d3E4anU4VU51c24Kc2FzdmlpS1VHQnpkblZndlNSdzNhYzd4RTRSQjlmZytjRnVUenpGaFNHRlVLVUpYaVh5d0FXZ0o4YU9DQXBBdwpnZ0tNTUE0R0ExVWREd0VCL3dRRUF3SUZvREFkQmdOVkhTVUVGakFVQmdnckJnRUZCUWNEQVFZSUt3WUJCUVVICkF3SXdEQVlEVlIwVEFRSC9CQUl3QURBZEJnTlZIUTRFRmdRVXgvaDVCOUFMSExuYWJaNjBzT2dvbnA3YlN0VXcKSHdZRFZSMGpCQmd3Rm9BVThXd3RDZDNYeXBNeGJLL3k0ajN1Q0dBYUJ3WXdnZ0lMQmdOVkhSRUVnZ0lDTUlJQgovb0lRYkc5dVoyaHZjbTR0WW1GamEyVnVaSUlnYkc5dVoyaHZjbTR0WW1GamEyVnVaQzVzYjI1bmFHOXliaTF6CmVYTjBaVzJDSkd4dmJtZG9iM0p1TFdKaFkydGxibVF1Ykc5dVoyaHZjbTR0YzNsemRHVnRMbk4yWTRJUmJHOXUKWjJodmNtNHRabkp2Ym5SbGJtU0NJV3h2Ym1kb2IzSnVMV1p5YjI1MFpXNWtMbXh2Ym1kb2IzSnVMWE41YzNSbApiWUlsYkc5dVoyaHZjbTR0Wm5KdmJuUmxibVF1Ykc5dVoyaHZjbTR0YzNsemRHVnRMbk4yWTRJWGJHOXVaMmh2CmNtNHRaVzVuYVc1bExXMWhibUZuWlhLQ0oyeHZibWRvYjNKdUxXVnVaMmx1WlMxdFlXNWhaMlZ5TG14dmJtZG8KYjNKdUxYTjVjM1JsYllJcmJHOXVaMmh2Y200dFpXNW5hVzVsTFcxaGJtRm5aWEl1Ykc5dVoyaHZjbTR0YzNsegpkR1Z0TG5OMlk0SVliRzl1WjJodmNtNHRjbVZ3YkdsallTMXRZVzVoWjJWeWdpaHNiMjVuYUc5eWJpMXlaWEJzCmFXTmhMVzFoYm1GblpYSXViRzl1WjJodmNtNHRjM2x6ZEdWdGdpeHNiMjVuYUc5eWJpMXlaWEJzYVdOaExXMWgKYm1GblpYSXViRzl1WjJodmNtNHRjM2x6ZEdWdExuTjJZNElNYkc5dVoyaHZjbTR0WTNOcGdoeHNiMjVuYUc5eQpiaTFqYzJrdWJHOXVaMmh2Y200dGMzbHpkR1Z0Z2lCc2IyNW5hRzl5YmkxamMya3ViRzl1WjJodmNtNHRjM2x6CmRHVnRMbk4yWTRJUWJHOXVaMmh2Y200dFltRmphMlZ1WkljRWZ3QUFBVEFOQmdrcWhraUc5dzBCQVFzRkFBT0MKQVFFQWV5UlhCWnI5Z1RmTGlsNGMvZElaSlVYeFh4ckFBQmtJTG55QkdNdkFqaFJoRndLZ09VU0MvMGUyeDYvTQpoTi9SWElYVzdBYUF0a25ZSHFLa3piMDZsbWhxczRHNWVjNkZRZDViSGdGbnFPOHNWNEF6WVFSRWhDZjlrWWhUClVlRnJLdDdOQllHNFNXSnNYK2M0ZzU5RlZGZkIzbTZscStoR3JaY085T2NIQ1NvVDM2SVRPeERDT3lrV002WHcKVW5zYWtaaHRwQ3lxdHlwQXZqaURNM3ZTY2txVTFNSWxLSnA1Z3lGT3k2VHVwQ01tYnRiWlRpSEtaN0ZlcmlmcwoyYng4Z0JmaldFQnEwMEhVWTdyY3RFNzFpVk11WURTczAwYTB2c1ZGQ240akppeWFnM0lHWkdud0FHQk1zR2h3ClFJcndjRHgwdy91NGR1VWRNMzBpaU1WZ0pnPT0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=
  tls.key: LS0tLS1CRUdJTiBFQyBQUklWQVRFIEtFWS0tLS0tCk1IY0NBUUVFSUwzbjZVZzlhZU1Day9XbkZ2L1pmSTlxMkIyakxnbjFRWGQwcjhIL3k2QkhvQW9HQ0NxR1NNNDkKQXdFSG9VUURRZ0FFTE5BVUllSy92emloZzdrVDl3cThqdThVTnVzbnNhc3ZpaUtVR0J6ZG5WZ3ZTUnczYWM3eApFNFJCOWZnK2NGdVR6ekZoU0dGVUtVSlhpWHl3QVdnSjhRPT0KLS0tLS1FTkQgRUMgUFJJVkFURSBLRVktLS0tLQo=
```

For more information on creating a secret, see [the Kubernetes documentation.](https://kubernetes.io/docs/concepts/configuration/secret/#creating-a-secret-manually) The secret must be created in the `longhorn-system` namespace for Longhorn to access it.

> Note: Make sure to use `echo -n` when generating the base64 encoding,
> otherwise a new line will be added at the end of the string
> which will cause an error during loading of the certificates.


# History
Available since v1.3.0 [#3839](https://github.com/longhorn/longhorn/issues/3839)


---

## Article: advanced-resources/security/volume-encryption.md

---
title: Volume Encryption
weight: 2
---

Longhorn supports volume encryption in both `Filesystem` and `Block` modes, providing protection against unauthorized access, data breaches, and compliance violations. Backups created from encrypted volumes are also encrypted.

Volume encryption is made possible by the Linux kernel module `dm_crypt`, the command-line utility `cryptsetup`, and Kubernetes Secrets. `dm_crypt` and `cryptsetup` handle the creation and management of encrypted devices, while Secrets (and related permissions) facilitate secure storage of encryption keys.

# Requirements

To use encrypted volumes, ensure that the `dm_crypt` kernel module is loaded and that `cryptsetup` is installed on your worker nodes.

# Setting up Kubernetes Secrets and StorageClasses

Longhorn uses Kubernetes Secrets for secure storage of encryption keys. Kubernetes allows usage of template parameters that are resolved during volume creation. To use a Secret with an encrypted volume, you must configure the Secret as a StorageClass parameter.

Template parameters allow you to use Secrets with individual volumes or with a collection of volumes. For more information about template parameters, see [StorageClass Secrets](https://kubernetes-csi.github.io/docs/secrets-and-credentials-storage-class.html) in the Kubernetes CSI Developer Documentation.

In the following example, the encryption key is specified as string data in the `CRYPTO_KEY_VALUE` parameter of the Secret. Using string data eliminates the need for Base64 encoding before the Secret is submitted via kubectl create.

Besides `CRYPTO_KEY_VALUE`, parameters `CRYPTO_KEY_CIPHER`, `CRYPTO_KEY_HASH`, `CRYPTO_KEY_SIZE`, and `CRYPTO_PBKDF` provide the customization for volume encryption.
- `CRYPTO_KEY_CIPHER`: Sets the cipher specification algorithm string. The default value is `aes-xts-plain64` for LUKS.
- `CRYPTO_KEY_HASH`: Specifies the passphrase hash for `open`. The default value is `sha256`.
- `CRYPTO_KEY_SIZE`: Sets the key size in bits and it must be a multiple of 8. The default value is `256`.
- `CRYPTO_PBKDF`: Sets Password-Based Key Derivation Function (PBKDF) algorithm for LUKS keyslot. The default value is `argon2i`.

For more information, see [cryptsetup(8)](https://man7.org/linux/man-pages/man8/cryptsetup.8.html) in the Linux man pages.

- Example of a Secret:
  ```yaml
  apiVersion: v1
  kind: Secret
  metadata:
    name: longhorn-crypto
    namespace: longhorn-system
  stringData:
    CRYPTO_KEY_VALUE: "Your encryption passphrase"
    CRYPTO_KEY_PROVIDER: "secret"
    CRYPTO_KEY_CIPHER: "aes-xts-plain64"
    CRYPTO_KEY_HASH: "sha256"
    CRYPTO_KEY_SIZE: "256"
    CRYPTO_PBKDF: "argon2i"
  ```

- Example of a StorageClass with a global Secret:
  ```yaml
  kind: StorageClass
  apiVersion: storage.k8s.io/v1
  metadata:
    name: longhorn-crypto-global
  provisioner: driver.longhorn.io
  allowVolumeExpansion: true
  parameters:
    numberOfReplicas: "3"
    staleReplicaTimeout: "2880" # 48 hours in minutes
    fromBackup: ""
    encrypted: "true"
    # global secret that contains the encryption key that will be used for all volumes
    csi.storage.k8s.io/provisioner-secret-name: "longhorn-crypto"
    csi.storage.k8s.io/provisioner-secret-namespace: "longhorn-system"
    csi.storage.k8s.io/node-publish-secret-name: "longhorn-crypto"
    csi.storage.k8s.io/node-publish-secret-namespace: "longhorn-system"
    csi.storage.k8s.io/node-stage-secret-name: "longhorn-crypto"
    csi.storage.k8s.io/node-stage-secret-namespace: "longhorn-system"
    csi.storage.k8s.io/node-expand-secret-name: "longhorn-crypto"
    csi.storage.k8s.io/node-expand-secret-namespace: "longhorn-system"
  ```

- Example of a StorageClass with a volume-specific Secret:
  ```yaml
  kind: StorageClass
  apiVersion: storage.k8s.io/v1
  metadata:
    name: longhorn-crypto-per-volume
  provisioner: driver.longhorn.io
  allowVolumeExpansion: true
  parameters:
    numberOfReplicas: "3"
    staleReplicaTimeout: "2880" # 48 hours in minutes
    fromBackup: ""
    encrypted: "true"
    # per volume secret which utilizes the `pvc.name` and `pvc.namespace` template parameters
    csi.storage.k8s.io/provisioner-secret-name: ${pvc.name}
    csi.storage.k8s.io/provisioner-secret-namespace: ${pvc.namespace}
    csi.storage.k8s.io/node-publish-secret-name: ${pvc.name}
    csi.storage.k8s.io/node-publish-secret-namespace: ${pvc.namespace}
    csi.storage.k8s.io/node-stage-secret-name: ${pvc.name}
    csi.storage.k8s.io/node-stage-secret-namespace: ${pvc.namespace}
    csi.storage.k8s.io/node-expand-secret-name: ${pvc.name}
    csi.storage.k8s.io/node-expand-secret-namespace: ${pvc.namespace}
  ```

# Using an Encrypted Volume

To create an encrypted volume, you must create a PVC using a StorageClass that has been configured for encryption. The above StorageClass examples can be used as a starting point.

After creation of the PVC it will remain in `Pending` state till the associated secret has been created and can be retrieved
A newly-created PVC remains in the `Pending` state until the associated Secret is created and can be retrieved by the csi `external-provisioner` sidecar. Afterwards, the regular volume creation process continues with encryption taking effect.

# Filesystem Expansion

Longhorn supports [both online and offline expansion](../../../nodes-and-volumes/volumes/expansion/#encrypted-volume) for encrypted volumes.

StorageClass parameters are needed to enable online expansion:

- `csi.storage.k8s.io/node-expand-secret-name`
- `csi.storage.k8s.io/node-expand-secret-namespace`

> **Notice**  
> - Longhorn v1.8.0 does not support expansion of V2 volumes.

# History

- Encryption of volumes in `Filesystem` mode available starting v1.2.0 ([#1859](https://github.com/longhorn/longhorn/issues/1859))
- Encryption of volumes in `Block` mode available starting v1.6.0 ([#4883](https://github.com/longhorn/longhorn/issues/4883))


---

## Article: advanced-resources/support-managed-k8s-service/_index.md

---
title: Support Managed Kubernetes Service
weight: 4
---



---

## Article: advanced-resources/support-managed-k8s-service/manage-node-group-on-aks.md

---
title:  Manage Node-Group on Azure AKS
weight: 2
---

See [Create and manage multiple node pools for a cluster in Azure Kubernetes Service (AKS)](https://docs.microsoft.com/en-us/azure/aks/use-multiple-node-pools) for more information.

Following is an example to replace cluster nodes with a new storage size.


## Storage Expansion

AKS does not support additional disk in its [template](https://docs.microsoft.com/en-us/azure/templates/Microsoft.ContainerService/2022-01-01/managedclusters?tabs=bicep#template-format). It is possible for manual disk attachment. Then raw device needs to be mounted either by manually mounting in VM or during launch with CustomScriptExtension that [is not supported](https://docs.microsoft.com/en-us/azure/aks/support-policies#user-customization-of-agent-nodes) in AKS.

1. In Longhorn, set `replica-replenishment-wait-interval` to `0`.

2. Add a new node-pool. Later Longhorn components will be automatically deployed on the nodes in this pool.

    ```
    AKS_NODEPOOL_NAME_NEW=<new-nodepool-name>
    AKS_RESOURCE_GROUP=<aks-resource-group>
    AKS_CLUSTER_NAME=<aks-cluster-name>
    AKS_DISK_SIZE_NEW=<new-disk-size-in-gb>
    AKS_NODE_NUM=<number-of-nodes>
    AKS_K8S_VERSION=<kubernetes-version>

    az aks nodepool add \
      --resource-group ${AKS_RESOURCE_GROUP} \
      --cluster-name ${AKS_CLUSTER_NAME} \
      --name ${AKS_NODEPOOL_NAME_NEW} \
      --node-count ${AKS_NODE_NUM} \
      --node-osdisk-size ${AKS_DISK_SIZE_NEW} \
      --kubernetes-version ${AKS_K8S_VERSION} \
      --mode System
    ```

3. Using Longhorn UI to disable the disk scheduling and request eviction for nodes in the old node-pool.

4. Cordon and drain Kubernetes nodes in the old node-pool.
    ```
    AKS_NODEPOOL_NAME_OLD=<old-nodepool-name>

    for n in `kubectl get nodes | grep ${AKS_NODEPOOL_NAME_OLD}- | awk '{print $1}'`; do
      kubectl cordon $n && \
      kubectl drain $n --ignore-daemonsets --delete-emptydir-data
    done
    ```

5. Delete old node-pool.
    ```
    az aks nodepool delete \
      --cluster-name ${AKS_CLUSTER_NAME} \
      --name ${AKS_NODEPOOL_NAME_OLD} \
      --resource-group ${AKS_RESOURCE_GROUP}
    ```


---

## Article: advanced-resources/support-managed-k8s-service/manage-node-group-on-eks.md

---
title:  Manage Node-Group on AWS EKS
weight: 1
---

EKS supports configuring the same launch template. The nodes in the node-group will be recycled by new nodes with new configurations when updating the launch template version.

See [Launch template support](https://docs.aws.amazon.com/eks/latest/userguide/launch-templates.html) for more information.

The following is an example to replace cluster nodes with new storage size.


## Storage Expansion

1. In Longhorn, set `replica-replenishment-wait-interval` to `0`.

2. Go to the launch template of the EKS cluster node-group. You can find in the EKS cluster tab `Configuration/Compute/<node-group-name>` and click the launch template.

3. Click `Modify template (Create new version)` in the `Actions` drop-down menu.

4. Choose the `Source template version` in the `Launch template name and version description`.

5. Follow steps to [Expand volume](#expand-volume), or [Create additional volume](#create-additional-volume).
> **Note:** If you choose to expand by [create additional volume](#create-additional-volume), the disks need to be manually added to the disk list of the nodes after the EKS cluster upgrade.


### Expand volume
1. Update the volume size in `Configure storage`.

2. Click `Create template version` to save changes.

3. Go to the EKS cluster node-group and change `Launch template version` in `Node Group configuration`. Track the status in the `Update history` tab.


### Create additional volume
1. Click `Advanced` then `Add new volume` in `Configure storage` and fill in the fields.

2. Adjust the auto-mount script and add to `User data` in `Advanced details`. Make sure the `DEV_PATH` matches the `Device name` of the additional volume.
    ```
    MIME-Version: 1.0
    Content-Type: multipart/mixed; boundary="==MYBOUNDARY=="

    --==MYBOUNDARY==
    Content-Type: text/x-shellscript; charset="us-ascii"

    #!/bin/bash

    # https://docs.aws.amazon.com/eks/latest/userguide/launch-templates.html#launch-template-user-data
    echo "Running custom user data script"

    DEV_PATH="/dev/sdb"
    mkfs -t ext4 ${DEV_PATH}

    MOUNT_PATH="/mnt/longhorn"
    mkdir ${MOUNT_PATH}
    mount ${DEV_PATH} ${MOUNT_PATH}
    ```

3. Click `Create template version` to save changes.

4. Go to the EKS cluster node-group and change `Launch template version` in `Node Group configuration`. Track the status in the `Update history` tab.

5. In Longhorn, add the path of the mounted disk into the disk list of the nodes.


---

## Article: advanced-resources/support-managed-k8s-service/manage-node-group-on-gke.md

---
title:  Manage Node-Group on GCP GKE
weight: 3
---

See [Migrating workloads to different machine types](https://cloud.google.com/kubernetes-engine/docs/tutorials/migrating-node-pool) for more information.

The following is an example to replace cluster nodes with new storage size.


## Storage Expansion

GKE supports adding additional disk with `local-ssd-count`. However, each local SSD is fixed size to 375 GB. We suggest expanding the node size via node pool replacement.

1. In Longhorn, set `replica-replenishment-wait-interval` to `0`.

2. Add a new node-pool. Later Longhorn components will be automatically deployed on the nodes in this pool.

    ```
    GKE_NODEPOOL_NAME_NEW=<new-nodepool-name>
    GKE_REGION=<gke-region>
    GKE_CLUSTER_NAME=<gke-cluster-name>
    GKE_IMAGE_TYPE=Ubuntu
    GKE_MACHINE_TYPE=<gcp-machine-type>
    GKE_DISK_SIZE_NEW=<new-disk-size-in-gb>
    GKE_NODE_NUM=<number-of-nodes>

    gcloud container node-pools create ${GKE_NODEPOOL_NAME_NEW} \
      --region ${GKE_REGION} \
      --cluster ${GKE_CLUSTER_NAME} \
      --image-type ${GKE_IMAGE_TYPE} \
      --machine-type ${GKE_MACHINE_TYPE} \
      --disk-size ${GKE_DISK_SIZE_NEW} \
      --num-nodes ${GKE_NODE_NUM}
  
    gcloud container node-pools list \
      --zone ${GKE_REGION} \
      --cluster ${GKE_CLUSTER_NAME} 
    ```

3. Using Longhorn UI to disable the disk scheduling and request eviction for nodes in the old node-pool.

4. Cordon and drain Kubernetes nodes in the old node-pool.
    ```
    GKE_NODEPOOL_NAME_OLD=<old-nodepool-name>
    for n in `kubectl get nodes | grep ${GKE_CLUSTER_NAME}-${GKE_NODEPOOL_NAME_OLD}- | awk '{print $1}'`; do
      kubectl cordon $n && \
      kubectl drain $n --ignore-daemonsets --delete-emptydir-data
    done
    ```

5. Delete old node-pool.
    ```
    gcloud container node-pools delete ${GKE_NODEPOOL_NAME_OLD}\
      --zone ${GKE_REGION} \
      --cluster ${GKE_CLUSTER_NAME}
    ```


---

## Article: advanced-resources/support-managed-k8s-service/upgrade-k8s-on-aks.md

---
title:  Upgrade Kubernetes on Azure AKS
weight: 5
---

AKS provides `az aks upgrade` for in-places nodes upgrade by node reimaged, but this will cause the original Longhorn disks missing, then there will be no disks allowing replica rebuilding in upgraded nodes anymore.

We suggest using node-pool replacement to upgrade the agent nodes but use `az aks upgrade` for control plane nodes to ensure data safety.

1. In Longhorn, set `replica-replenishment-wait-interval` to `0`.

2. Upgrade AKS control plane.
    ```
    AKS_RESOURCE_GROUP=<aks-resource-group>
    AKS_CLUSTER_NAME=<aks-cluster-name>
    AKS_K8S_VERSION_UPGRADE=<aks-k8s-version>

    az aks upgrade \
        --resource-group ${AKS_RESOURCE_GROUP} \
        --name ${AKS_CLUSTER_NAME} \
        --kubernetes-version ${AKS_K8S_VERSION_UPGRADE} \
        --control-plane-only
    ```

3. Add a new node-pool.

    ```
    AKS_NODEPOOL_NAME_NEW=<new-nodepool-name>
    AKS_DISK_SIZE=<disk-size-in-gb>
    AKS_NODE_NUM=<number-of-nodes>

    az aks nodepool add \
      --resource-group ${AKS_RESOURCE_GROUP} \
      --cluster-name ${AKS_CLUSTER_NAME} \
      --name ${AKS_NODEPOOL_NAME_NEW} \
      --node-count ${AKS_NODE_NUM} \
      --node-osdisk-size ${AKS_DISK_SIZE} \
      --kubernetes-version ${AKS_K8S_VERSION_UPGRADE} \
      --mode System
    ```

4. Using Longhorn UI to disable the disk scheduling and request eviction for nodes in the old node-pool.

5. Cordon and drain Kubernetes nodes in the old node-pool.
    ```
    AKS_NODEPOOL_NAME_OLD=<old-nodepool-name>

    for n in `kubectl get nodes | grep ${AKS_NODEPOOL_NAME_OLD}- | awk '{print $1}'`; do
      kubectl cordon $n && \
      kubectl drain $n --ignore-daemonsets --delete-emptydir-data
    done
    ```

6. Delete old node-pool.
    ```
    az aks nodepool delete \
      --cluster-name ${AKS_CLUSTER_NAME} \
      --name ${AKS_NODEPOOL_NAME_OLD} \
      --resource-group ${AKS_RESOURCE_GROUP}
    ```


---

## Article: advanced-resources/support-managed-k8s-service/upgrade-k8s-on-eks.md

---
title: Upgrade Kubernetes on AWS EKS
weight: 4
---

In Longhorn, set `replica-replenishment-wait-interval` to `0`.

See [Updating a cluster](https://docs.aws.amazon.com/eks/latest/userguide/update-cluster.html) for instructions.

> **Note:** If you have created [addition disks](../manage-node-group-on-eks#create-additional-volume) for Longhorn, you will need to manually add the path of the mounted disk into the disk list of the upgraded nodes.

---

## Article: advanced-resources/support-managed-k8s-service/upgrade-k8s-on-gke.md

---
title: Upgrade Kubernetes on GCP GKE
weight: 6
---

In Longhorn, set `replica-replenishment-wait-interval` to `0`.

See [Upgrading the cluster](https://cloud.google.com/kubernetes-engine/docs/how-to/upgrading-a-cluster#upgrading_the_cluster) and [Upgrading node pools](https://cloud.google.com/kubernetes-engine/docs/how-to/upgrading-a-cluster#upgrading-nodes) for instructions.


---

## Article: advanced-resources/driver-migration/_index.md

---
title: CSI Driver Migration
weight: 100
---

---

## Article: advanced-resources/driver-migration/migrating-flexvolume.md

---
title: Migrating from the Flexvolume Driver to CSI
weight: 5
---

As of Longhorn v0.8.0, the Flexvolume driver is no longer supported. This guide will show you how to migrate from the Flexvolume driver to CSI. CSI is the newest out-of-tree Kubernetes storage interface.

> Note that the volumes created and used through one driver won't be recognized by Kubernetes using the other driver. So please don't switch the driver (e.g. during an upgrade) if you have existing volumes created using the old driver.

Ensure your Longhorn App is up to date. Follow the relevant upgrade procedure before proceeding.

The migration path between drivers requires backing up and restoring each volume and will incur both API and workload downtime. This can be a tedious process. Consider deleting unimportant workloads using the old driver to reduce effort.

1. [Back up existing volumes](../../../snapshots-and-backups/backup-and-restore/create-a-backup).
2. On Rancher UI, navigate to the `Catalog Apps` page, locate the `Longhorn` app and click the `Up to date` button. Under `Kubernetes Driver`, select
`flexvolume`. We recommend leaving `Flexvolume Path` empty. Click `Upgrade`.
3. Restore each volume. This [procedure](../../../snapshots-and-backups/backup-and-restore/restore-statefulset) is tailored to the StatefulSet workload, but the process is approximately the same for all workloads.

---

## Article: advanced-resources/rebuilding/_index.md

---
title: Replica Rebuilding
weight: 6
---
- [Replica Rebuilding Workflow](#replica-rebuilding-workflow)
  - [Full Replica Rebuilding](#full-replica-rebuilding)
  - [Delta Replica Rebuilding](#delta-replica-rebuilding)
  - [Fast Replica Rebuilding](#fast-replica-rebuilding)
- [Factors That Affect Rebuilding Performance](#factors-that-affect-rebuilding-performance)
- [Use Cases](#use-cases)
  - [Node Reboot During Upgrade](#node-reboot-during-upgrade)
  - [Short-Term Node Drain](#short-term-node-drain)
- [Relevant Settings](#relevant-settings)
  - [Settings Trade-Off Analysis](#settings-trade-off-analysis)

When Longhorn detects a failed or deleted replica, it automatically initiates a rebuilding process. This document outlines the replica rebuilding workflow for v1 data engine, including **full**, **delta**, and **fast** rebuilding methods. It also explains the limitations associated with each method.

**Rebuilding will not start in the following scenarios**:

- The volume is migrating to another node.
- The volume is an old restore/DR volume.
- The volume is expanding in size.

## Replica Rebuilding Workflow

Replica rebuilding may be triggered in the following scenarios for v1 data engine:

- A node is rebooted, drained or evicted.
- A replica becomes unhealthy or is deleted.

 {{< figure alt="Replica Rebuilding Flow Diagram" src="/img/diagrams/architecture/replica-rebuilding-flow.png" >}}

1. Mark the target replica with `WO` (Write-Only) mode.
2. Create a new snapshot to serve as the volume head reference point for data integrity checks.
3. Generate the synchronization file list for the volume head and snapshot files.

- For **V1 Data Engine**:  
  4. Launch a receiver server on the target replica for each snapshot, then instruct the source replica to begin data synchronization.  
  5. For each snapshot synchronization:
    - Check whether the snapshot file exists in the target replica's data directory:
        - If **NO**, transfer the entire snapshot data from the source replica to the target replica. See [Full Replica Rebuilding](#full-replica-rebuilding)
        - If **YES**, check whether the snapshot checksum files exist, the modification time and checksums are identical between the target and source replicas:
            - If **YES**, Longhorn skips transferring that snapshot’s data. This optimization reduces CPU usage, disk I/O, network I/O, and overall rebuild time. See [Fast Replica Rebuilding](#fast-replica-rebuilding)
            - If **NO**, Longhorn calculates and compares block-by-block checksums for the snapshot file using the SHA-512 algorithm. If mismatches are found, only the differing data blocks are synchronized. See [Delta Replica Rebuilding](#delta-replica-rebuilding)
- For **V2 Data Engine**:  
  4. Expose the source and target replicas and prepare shallow copy with the SPDK engine.  
  5. For each snapshot:
    - Check if the snapshot timestamp, snapshot actual size, and snapshot checksum match between the source and target snapshots:
      - If **YES**, Longhorn skips transferring that snapshot’s data.
      - If **NO**, check if both the source and target snapshot hold ranged checksums
        - If **YES**, fetching and comparing the entire ranges' checksums of the source and target snapshot. If mismatches are found, only copy mismatched parts. See [Fast Replica Rebuilding](#fast-replica-rebuilding)
        - If **NO**, delete the existing target snapshot. Then start copying the entire snapshot from the source replica to the target replica. See [Full Replica Rebuilding](#full-replica-rebuilding)

### Full Replica Rebuilding

If the replica is unrecoverable or has no existing data, Longhorn synchronizes all data from a healthy replica. It reconstructs the replica by transferring the full snapshot chain.

Full replica rebuilding consumes significant network bandwidth and results in heavy disk write operations on the target node. However, it is required when the target replica has no usable data.

### Delta Replica Rebuilding

Delta replica rebuilding is only for v1 data engine. It starts with a reusable failed replica, and it checks the data integrity for all snapshots' data block by block.

- This is available for failed replica reuse only, and there is an existing snapshot file (with the same name) in the failed replica data directory
- When a snapshot has no checksum, Longhorn performs delta replica rebuilding for this snapshot instead.

- **Pros**:
  - Reduce network bandwidth consumption.
- **Cons**:
  - Increased CPU overhead because Longhorn will compute the checksum of the snapshot data block by block for data integrity check.
  - Rebuilding time is influenced by CPU performance.

### Fast Replica Rebuilding

Fast rebuilding is enabled when:

- The fast replica building setting is enabled:
  - `fast-replica-rebuild-enabled: true`
- Snapshot checksum files are created (the snapshot checksums are pre-computed) via one of the following methods:
  - `snapshot-data-integrity` is set to `enabled`:  scheduled job calculates checksums for all snapshots at a configured interval (default: 7 days), or
  - `snapshot-data-integrity-immediate-check-after-snapshot-creation` is set to `true`: the snapshot checksum is calculated immediately after snapshot creation.

  > **Note**: These checksum calculations consume storage and computing resources. The calculation time is unpredictable and may negatively impact the storage performance.  
  > For more details, see [Snapshot Data Integrity](../data-integrity/snapshot-data-integrity-check)

- **Pros**:
  - Minimize network bandwidth consumption.
  - Minimize disk IO.
- **Cons**:
  - Calculating the snapshot checksum can be time-consuming.
  - The timing of checksum calculation is unpredictable. It can be triggered even if the node is experiencing heavy IO load.

For more details, see [Fast Replica Rebuilding](./fast-replica-rebuilding).

## Factors That Affect Rebuilding Performance

- **Large volume head**
  - **Why it matters**:  
    The volume head is a special file that never has a precomputed checksum. If a replica fails, Longhorn must always synchronize the entire volume head file. The larger the volume head, the longer the rebuild process will take.
  - **How to prevent**:  
    Take snapshots regularly to reduce the amount of data in the volume head. Schedule snapshot creation before planned maintenance to minimize rebuild time.
- **No snapshots exist**
  - **Why it matters**:  
    Without snapshots, Longhorn cannot skip data transfer or reuse existing data. If a snapshot of the volume head was just created, but its checksum has not yet been calculated, Longhorn must perform delta rebuilding for that snapshot. This increases CPU load due to block-by-block checksum comparisons.
  - **How to prevent**:  
    1. Enable `snapshot-data-integrity-immediate-check-after-snapshot-creation` or `snapshot-data-integrity`, so checksums are precomputed.
       Trade-off: Increases CPU, disk I/O, and storage usage during checksum computation.
    2. Use a recurring job to take snapshots regularly.
- **Snapshot Purged**
  - **Why it matters**:  
    When snapshot purging starts, system-generated snapshots are coalesced into the next snapshot. As a result, the checksum of the next snapshot becomes invalid.
  - **How to prevent**:
    1. Enable `snapshot-data-integrity-immediate-check-after-snapshot-creation` to ensure the checksums are generated after purging.
    2. Proactively create a snapshot and allow time for its checksum to be computed before upgrades or rebuilds.
- **Concurrent rebuilds**
  - **Why it matters**:  
    Multiple rebuilds running simultaneously on a node can heavily consume CPU, disk I/O, and network I/O resources, impacting overall performance.
  - **How to prevent**:  
    Tune the number of concurrent rebuilds using the `concurrent-replica-rebuild-per-node-limit` setting.
- **Multiple replica failures**
  - **Why it matters**:  
    Increases rebuild complexity and duration. If the `auto-cleanup-system-generated-snapshot` setting is `true` and no user-created snapshots exist, then when two replicas fail before either has been rebuilt, Longhorn must perform at least one **full data transfer** to restore volume health.  
    For more details, see [Avoid "full data transfer" when rebuilding two failed replicas](https://github.com/longhorn/longhorn/issues/9335)
  - **How to prevent**:
    1. Manually disable `auto-cleanup-system-generated-snapshot before doing maintenance` before performing maintenance.
    2. Take user-created snapshots of all volumes before starting maintenance.
    3. Use a recurring job to take snapshots regularly.

## Use Cases

### Node Reboot During Upgrade

When a worker node with replicas is rebooted as part of a planned upgrade:

1. The replica on that node becomes temporarily unavailable and fails, but read/write operations continue.
2. If the node recovers within the `replica-replenishment-wait-interval`, Longhorn initiates a rebuild using the reusable failed replica.

During the rebuilding process:

  1. Longhorn selects the latest reusable failed replica if multiple reusable failed replicas are available.
  2. Based on the rebuild scenario:
      - **If fast replica rebuilding is enabled and all snapshot checksums exist:**
        [Fast Replica Rebuilding](#fast-replica-rebuilding) is triggered.
        - Only changed blocks in the volume head are synced, avoiding both full and delta rebuilding.
      - **If fast replica rebuilding is enabled but some snapshot checksums are missing:**
        [Delta Replica Rebuilding](#delta-replica-rebuilding) is used.
        - Changed blocks from snapshots without checksums are synced, avoiding full rebuilding.
      - **If fast replica rebuilding is disabled:**
        - Changed blocks of **all snapshots** are synced (delta rebuilding), avoiding full rebuilding.

### Short-Term Node Drain

If a worker node is drained for short-term maintenance and then quickly restored:

1. The replica on the drained node is marked as failed immediately.
2. If the node is uncordoned before `replica-replenishment-wait-interval` expires, Longhorn attempts to reuse the failed replica..
3. Rebuild behavior follows the same logic as described in the previous use case.

## Relevant Settings

| Setting| Default | Description |
| :--- | :---: | :--- |
| `fast-replica-rebuild-enabled` | `true` | Enables fast replica rebuilding. Relies on precomputed snapshot checksums. |
| [snapshot-data-integrity](../data-integrity/snapshot-data-integrity-check) | `fast-check` | Hashes snapshot disk files only if they are unhashed or their modification time has changed. |
| `snapshot-data-integrity-cronjob` | `0 0 */7 * *` | Cron schedule to compute checksums for all snapshots (default: every 7 days). |
| `snapshot-data-integrity-immediate-check-after-snapshot-creation` | `false` | If enabled, checksums are computed immediately after snapshot creation. |
| `replica-replenishment-wait-interval` | `600` |    Time in seconds to wait before creating a new replica, allowing reuse of failed replicas. |
| `concurrent-replica-rebuild-per-node-limit` | `5` | Limits the number of concurrent replica rebuilds per node. |
| `offline-replica-rebuilding` | `false` | Determines if degraded replicas are rebuilt while the volume is detached. |
||||

### Settings Trade-Off Analysis

- **[fast-replica-rebuild-enabled](../../references/settings#fast-replica-rebuild-enabled)**
  - `enabled`  
    Skips snapshot data transfer if checksums are up to date — fast rebuild, but data isn't re-validated.
  - `disabled`  
    Performs delta rebuilding using block comparisons — slower, but ensures snapshot data integrity.
- **[snapshot-data-integrity](../data-integrity/snapshot-data-integrity-check)**
  - `enabled`  
  By default, computes snapshot checksums every 7 days. This consumes resources and increases processing time.
- **[snapshot-data-integrity-cronjob](../../references/settings#snapshot-data-integrity-check-cronjob)**
  - Default: `0 0 */7 * *`  
  If the `snapshot-data-integrity` setting is `enabled`, it defines when snapshot checksums are recalculated. Snapshots created within this interval may lack precomputed checksums.
- **[snapshot-data-integrity-immediate-check-after-snapshot-creation](../../references/settings#immediate-snapshot-data-integrity-check-after-creating-a-snapshot)**
  - `true`  
  Immediately calculates checksums after snapshot creation, increasing CPU and disk I/O usage. Completes at an unpredictable time.
  - `false`  
  Snapshots may not have checksums until the next cron job runs with `snapshot-data-integrity` `enabled`. Delta rebuilding will be used if checksums are missing.
- **[replica-replenishment-wait-interval](../../references/settings#replica-replenishment-wait-interval)**
  - Default: `600` seconds
    - **Short interval**: May skip reusing failed replicas and trigger full rebuilds.
    - **Long interval**: Waits longer to reuse failed replicas but may delay recovery.
- **[concurrent-replica-rebuild-per-node-limit](../../references/settings#concurrent-replica-rebuild-per-node-limit)**
  - Default: `5`
    - **High limit**: May overload node resources and slow down rebuilds and workloads.
    - **Low limit**: Reduces resource strain but may increase total rebuild duration due to queuing.


---

## Article: advanced-resources/rebuilding/fast-replica-rebuilding.md

---
title: Fast Replica Rebuilding
weight: 5
---

Longhorn supports fast replica rebuilding based on the checksums of snapshot disk files.

## Introduction

The legacy replica rebuilding process walks through all snapshot disk files. For each data block, the client (healthy replica) hashes the local data block as well as requests the checksum of the corresponding data block on the remote side (rebuilt replica). Then, the client compares the two checksums to determine if the data block needs to be sent to the remote side and override the data block. Thus, it is an IO- and computing-intensive process, especially if the volume is large or contains a large number of snapshot files.

If users enable the snapshot data integrity check feature by configuring `snapshot-data-integrity` to `enabled` or `fast-check`, the change timestamps and the checksums of snapshot disk files are recorded. As long as the two below conditions are met, we can skip the synchronization of the snapshot disk file. 
- The change timestamps on the snapshot disk file and the value recorded are the same.
- Both the local and remote snapshot disk files have the same checksum.

Then, a reduction in the number of unnecessary computations can speed up the entire process as well as reduce the impact on the system performance.

## Settings
### Global Settings
- fast-replica-rebuild-enabled <br>

    The setting enables fast replica rebuilding feature. It relies on the checksums of snapshot disk files, so setting the snapshot-data-integrity to **enable** or **fast-check** is a prerequisite. Please refer to [Snapshot Data Integrity Check](../../data-integrity/snapshot-data-integrity-check).



---

## Article: advanced-resources/rebuilding/offline-replica-rebuilding.md

---
title: Offline Replica Rebuilding
weight: 5
---

Starting with v1.9.0, Longhorn supports offline replica rebuilding, allowing degraded volumes to automatically rebuild replicas while detached.​

## Global Setting `offline-replica-rebuilding`

- When enabled, Longhorn automatically initiates offline rebuilding for eligible volumes.​
- For more information, see [settings](../../../references/settings#offline-replica-rebuilding).

## Per-Volume Override

- You can override the global `Offline Replica Rebuilding` setting for each volume individually using the Longhorn UI or by running `kubectl -n longhorn-system edit volume [volume-name]` and modifying the `Spec.OfflineRebuilding` field.
- When set to `enabled` or `disabled`, this per-volume setting takes precedence over the global configuration. The default value is `ignored`.

| Global Setting (`offline-replica-rebuilding`) | Per-Volume Setting (`spec.offlineRebuilding`) | Offline Rebuilding Enabled |
| :-------------------------------------------: | :-------------------------------------------: | :------------------------: |
| `true`                                        | `ignored`                                     | Yes                        |
| `false`                                       | `ignored`                                     | No                         |
| `true`                                        | `enabled`                                     | Yes                        |
| `false`                                       | `enabled`                                     | Yes                        |
| `true`                                        | `disabled`                                    | No                         |
| `false`                                       | `disabled`                                    | No                         |

## Rebuilding Process

- When triggered, Longhorn attaches the volume without activating the frontend, rebuilds any missing replicas, and then detaches the volume upon completion.
- This process can be interrupted if the workload scales up.

## Rebuilding Not Started or Canceled

When offline rebuilding starts, degraded volumes can get stuck in the attached state if rebuilding conditions are not met. To prevent this, if the necessary conditions are not satisfied, offline rebuilding will not start or will be canceled.

- **Benefits:**
  - It ensures volumes don't remain stuck in the attached state if rebuilding never finishes.
  - It prevents wasteful rebuilding attempts.
  - It reduces unnecessary volume attachment and detachment cycles.
  - It provides predictable rebuilding behavior based on resource availability.

- **Required conditions:**

  Offline rebuilding automatically starts for degraded volumes once the required conditions are met. These conditions include:

  - A reusable failed replica exists, or
  - A disk candidate exists:
    - The instance manager on the node hosting the disk must be ready.
    - The disk's containing node is schedulable.
    - The disk itself is schedulable.

### Before Offline Rebuilding Starts

When offline rebuilding is enabled, Longhorn determines whether it should start.

1. Longhorn detects a degraded, detached volume.
2. The system validates whether the required conditions are met before starting the rebuild.
3. If the conditions are met, rebuilding proceeds. Otherwise, the volume remains detached.
4. The required conditions are re-evaluated when a node is added, becomes ready, or becomes schedulable.

### During Offline Rebuilding

Longhorn determines if a rebuilding process should be canceled while in progress.

1. Longhorn detects the volume's status when offline rebuilding starts and the volume is attached.
2. If the volume's `Scheduled` condition status becomes `False`, the offline rebuilding is canceled, and the volume is detached.
3. If the required conditions are met again, offline rebuilding restarts; otherwise, the volume remains detached.

### Examples

- Successful offline rebuilding:
  1. A volume is created with 3 replicas in a 3-worker-node cluster.
  2. Offline rebuilding is enabled.
  3. The volume is detached and then a replica of the volume is deleted.
  4. Offline rebuilding begins, and the volume is attached.
  5. After rebuilding finishes, the volume is detached.
- Offline rebuilding does not start even when it is enabled:
  1. A volume is created with 3 replicas in a 3-worker-node (A, B, and C) cluster.
  2. Offline rebuilding is enabled.
  3. Worker node A is unschedulable.
  4. The volume replica on worker node A is deleted.
  5. Because only two schedulable worker nodes exist, offline rebuilding will not start.
- A worker node is drained during offline rebuilding:
  1. A volume is created with 3 replicas in a 3-worker-node (A, B, and C) cluster.
  2. Offline rebuilding is enabled.
  3. The volume is detached, and then the volume replica on worker node A is deleted.
  4. Offline rebuilding begins, and the volume is attached to rebuild a replica on worker node A.
  5. Worker node A is drained making it unschedulable, and the volume replica on worker node A is deleted.
  6. The volume remains attached until the volume's `Scheduled` condition status becomes `False`.
  7. The volume is detached until worker node A is uncordoned or a new schedulable node is added.

## Limitations

Offline rebuilding is not supported for faulted volumes.


---

## Article: advanced-resources/backing-image/_index.md

---
title: Backing Image
weight: 7
---

---

## Article: advanced-resources/backing-image/backing-image-backup.md

---
title: Backing Image Backup
weight: 2
---

As of v1.6.0, Longhorn supports backing up of backing images.

## Prerequisites

You must first [set up a backup target](../../../snapshots-and-backups/backup-and-restore/set-backup-target). If you skip this crucial step, the missing backup target will prevent Longhorn from creating a backup of the backing image.

## Create a Backup of a Backing Image

Because backing images are globally unique within the Longhorn system, the corresponding backups are also globally unique and are identified using the same name.

### Create a Backup Using YAML

Example of backing image:
```yaml
apiVersion: longhorn.io/v1beta2
kind: BackingImage
metadata:
  name: parrot
  namespace: longhorn-system
spec:
  sourceType: download
  sourceParameters:
    url: https://longhorn-backing-image.s3-us-west-1.amazonaws.com/parrot.raw
  checksum: 304f3ed30ca6878e9056ee6f1b02b328239f0d0c2c1272840998212f9734b196371560b3b939037e4f4c2884ce457c2cbc9f0621f4f5d1ca983983c8cdf8cd9a
```

Example of YAML code used to create a backup of the sample backing image:
```yaml
apiVersion: longhorn.io/v1beta2
kind: BackupBackingImage
metadata:
  name: parrot-backup
  namespace: longhorn-system
spec:
  backingImage: parrot
  backupTargetName: default
  userCreated: true
  labels:
    usecase: test
    type: raw
```

> **IMPORTANT:**
> - `name`: If the names are not unique, Longhorn will not be able to create a backup of the backing image.
> - `backingImage`: The backing image for the backup.
> - `backupTargetName`: The backup target that is used to store the backup of the backing image.
> - `userCreated`: Set the value to `true` to indicate that you created the backup custom resource, which enabled the creation of the backup in the backupstore. The value `false` indicates that the backup custom resource was synced from the backupstore.
> - `labels`: You can add labels to the backing image backup.

### Create a Backup Using the Longhorn UI
1. Go to **Setting** > **Backing Image**.
2. Select the backing image that you want to back up, and then click **Back Up** in the **Operation** menu.

Longhorn creates the backup and adds the details to the **Backing Image Backup** list.

{{< figure src="/img/screenshots/backing-image/backup.png" >}}


## Restore a Backing Image from a Backup
You can restore a backing image in another cluster after creating a backup in the backupstore.

Example of YAML code used to restore a backing image:
```yaml
apiVersion: longhorn.io/v1beta2
kind: BackingImage
metadata:
    name: parrot-restore
    namespace: longhorn-system
spec:
    sourceType: restore
    sourceParameters:
        # change to your backup URL
        # backup-url: nfs://longhorn-test-nfs-svc.default:/opt/backupstore?backingImage=parrot
        backup-url: s3://backupbucket@us-east-1/?backingImage=parrot
        concurrent-limit: "2"
    checksum: 304f3ed30ca6878e9056ee6f1b02b328239f0d0c2c1272840998212f9734b196371560b3b939037e4f4c2884ce457c2cbc9f0621f4f5d1ca983983c8cdf8cd9a
```

> **IMPORTANT:**
> - `sourceType`: Set the value to `restore`.
> - `sourceParameters`: Configure the following parameters:
>   - `backup-url`: URL of the backing image resource in the backupstore. You can find this information in the status of the backup custom resource `.Status.URL`.
>   - `concurrent-limit`: Maximum number of worker threads that can concurrently run for each restore operation. When unspecified, Longhorn uses the default value.
> - `checksum`: You can specify the expected SHA-512 checksum of the backing image file, which Longhorn uses to validate the restored file. When unspecified, Longhorn uses the checksum of the restored file as the truth.

### Restore from a Backup Using the Longhorn UI
1. Go to **Setting** > **Backing Image**.
2. Select the backup that you want to use, and then click **Restore** in the **Operation** menu.
3. Click **OK**.

{{< figure src="/img/screenshots/backing-image/1.8.0/restore.png" >}}

## Volume with a Backing Image

When you create a backup of a volume, Longhorn automatically creates a backup of its backing image.

You can restore a volume with a backing image. If the image already exists in the cluster, Longhorn uses the image directly. If the image exists in the backupstore but not in the cluster, Longhorn automatically restores the backing image.


---

## Article: advanced-resources/backing-image/backing-image-encryption.md

---
title: Backing Image Encryption
weight: 2
---

Starting with v1.7.0, Longhorn allows you to encrypt and decrypt a backing image by cloning it. The backing image encryption mechanism utilizes the Linux kernel module `dm_crypt` and the command-line utility `cryptsetup`.

## Clone a Backing Image
You can clone a backing image using YAML code. Notice that, this will create a whole new backing image with the same content as the original one. The backing image also consumes the disk space.

Example of a downloaded backing image:

```yaml
apiVersion: longhorn.io/v1beta2
kind: BackingImage
metadata:
  name: parrot
  namespace: longhorn-system
spec:
  sourceType: download
  sourceParameters:
    url: https://longhorn-backing-image.s3-us-west-1.amazonaws.com/parrot.raw
  checksum: 304f3ed30ca6878e9056ee6f1b02b328239f0d0c2c1272840998212f9734b196371560b3b939037e4f4c2884ce457c2cbc9f0621f4f5d1ca983983c8cdf8cd9a
```

Example of YAML code used to clone the sample backing image:

```yaml
apiVersion: longhorn.io/v1beta2
kind: BackingImage
metadata:
  name: parrot-cloned
  namespace: longhorn-system
spec:
  sourceType: clone
  sourceParameters:
    backing-image: parrot
    encryption: ignore
```

> **Important:**
> - `backing-image`: Specify the name of the backing image to be cloned.
> - `encryption`: Set the value to `ignore` to directly clone the backing image. If the value is not given, Longhorn use `ignore` as default value.

You can also clone a backing image using the Longhorn UI.
1. Go to **Setting** > **Backing Image**.
2. Click **Create Backing Image**.
3. Configure the following settings:
  - **Created From**: Select **Clone From Existing Backing Image**.
  - **Encryption**: Select **Ignore**.
4. Click **OK**.

{{< figure src="/img/screenshots/backing-image/clone.png" >}}

## Encrypt a Backing Image
You can enable encryption during cloning of a backing image so that the image can be used with an encrypted volume.

Example of a downloaded backing image:

```yaml
apiVersion: longhorn.io/v1beta2
kind: BackingImage
metadata:
  name: parrot
  namespace: longhorn-system
spec:
  sourceType: download
  sourceParameters:
    url: https://longhorn-backing-image.s3-us-west-1.amazonaws.com/parrot.raw
  checksum: 304f3ed30ca6878e9056ee6f1b02b328239f0d0c2c1272840998212f9734b196371560b3b939037e4f4c2884ce457c2cbc9f0621f4f5d1ca983983c8cdf8cd9a
```

Example of YAML code used to clone and encrypt the sample backing image:

```yaml
apiVersion: longhorn.io/v1beta2
kind: BackingImage
metadata:
  name: parrot-cloned-encrypted
  namespace: longhorn-system
spec:
  sourceType: clone
  sourceParameters:
    backing-image: parrot
    encryption: encrypt
    secret: longhorn-crypto
    secret-namespace: longhorn-system
```

Example of YAML code used to encrypt the backing image:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: longhorn-crypto
  namespace: longhorn-system
stringData:
  CRYPTO_KEY_VALUE: "Your encryption passphrase"
  CRYPTO_KEY_PROVIDER: "secret"
  CRYPTO_KEY_CIPHER: "aes-xts-plain64"
  CRYPTO_KEY_HASH: "sha256"
  CRYPTO_KEY_SIZE: "256"
  CRYPTO_PBKDF: "argon2i"
```

> **Important:**
> - `backing-image`: Specify the name of the backing image to be cloned.
> - `encryption`: Set the value to `encrypt` to encrypt the backing image during cloning.
> - `secret`: Specify the secret used to encrypt the backing image.
> - `secret-namespace`: Specify the namespace of the secret used to encrypt the backing image.

You can also create an encrypted copy of a backing image using the Longhorn UI.
1. Go to **Setting** > **Backing Image**.
2. Click **Create Backing Image**.
3. Configure the following settings:
  - **Created From**: Select **Clone From Existing Backing Image**.
  - **Encryption**: Select **Encrypt**.
4. Specify the secret and secret namespace to be used for encryption.
5. Click **OK**.

{{< figure src="/img/screenshots/backing-image/encrypt.png" >}}

## Decrypt a Backing Image
You can decrypt an encrypted backing image through cloning.

Example of an encrypted backing image:

```yaml
apiVersion: longhorn.io/v1beta2
kind: BackingImage
metadata:
  name: parrot-cloned-encrypted
  namespace: longhorn-system
spec:
  sourceType: clone
  sourceParameters:
    backing-image: parrot
    encryption: encrypt
    secret: longhorn-crypto
    secret-namespace: longhorn-system
```

Example of YAML code used to encrypt and decrypt the backing image:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: longhorn-crypto
  namespace: longhorn-system
stringData:
  CRYPTO_KEY_VALUE: "Your encryption passphrase"
  CRYPTO_KEY_PROVIDER: "secret"
  CRYPTO_KEY_CIPHER: "aes-xts-plain64"
  CRYPTO_KEY_HASH: "sha256"
  CRYPTO_KEY_SIZE: "256"
  CRYPTO_PBKDF: "argon2i"
```

Example of YAML code used to decrypt the backing image:

```yaml
apiVersion: longhorn.io/v1beta2
kind: BackingImage
metadata:
  name: parrot-cloned-decrypt
  namespace: longhorn-system
spec:
  sourceType: clone
  sourceParameters:
    backing-image: parrot-cloned-encrypted
    encryption: decrypt
    secret: longhorn-crypto
    secret-namespace: longhorn-system
```

> **Important:**
> - `backing-image`: Specify the name of the backing image to be cloned.
> - `encryption`: Set the value to `decrypt` to decrypt the backing image during cloning.
> - `secret`: Specify the secret used to decrypt the backing image.
> - `secret-namespace`: Specify the namespace of the secret used to decrypt the backing image.

You can also decrypt a backing image (through cloning) using the Longhorn UI.
1. Go to **Setting** > **Backing Image**.
2. Click **Create Backing Image**.
3. Configure the following settings:
  - **Created From**: Select **Clone From Existing Backing Image**.
  - **Encryption**: Select **Decrypt**.
4. Specify the secret and secret namespace to be used for decryption.
5. Click **OK**.


{{< figure src="/img/screenshots/backing-image/decrypt.png" >}}


## Use an Encrypted Backing Image with an Encrypted Volume
The secret used to encrypt the backing image and the volume must be identical. Once the encrypted backing image is ready, you can create the StorageClass with the corresponding backing image and the secret to create the volume for the workload.

Example of YAML code for the encryption secret:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: longhorn-crypto
  namespace: longhorn-system
stringData:
  CRYPTO_KEY_VALUE: "Your encryption passphrase"
  CRYPTO_KEY_PROVIDER: "secret"
  CRYPTO_KEY_CIPHER: "aes-xts-plain64"
  CRYPTO_KEY_HASH: "sha256"
  CRYPTO_KEY_SIZE: "256"
  CRYPTO_PBKDF: "argon2i"
```

Example of YAML code for the StorageClass:
```yaml
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: longhorn-crypto-global
provisioner: driver.longhorn.io
allowVolumeExpansion: true
parameters:
  numberOfReplicas: "3"
  staleReplicaTimeout: "2880" # 48 hours in minutes
  fromBackup: ""
  encrypted: "true"
  backingImage: "parrot-cloned-encrypted"
  backingImageDataSourceType: "clone"
  # global secret that contains the encryption key that will be used for all volumes
  csi.storage.k8s.io/provisioner-secret-name: "longhorn-crypto"
  csi.storage.k8s.io/provisioner-secret-namespace: "longhorn-system"
  csi.storage.k8s.io/node-publish-secret-name: "longhorn-crypto"
  csi.storage.k8s.io/node-publish-secret-namespace: "longhorn-system"
  csi.storage.k8s.io/node-stage-secret-name: "longhorn-crypto"
  csi.storage.k8s.io/node-stage-secret-namespace: "longhorn-system"
```

For more information, see [Volume Encryption](../../security/volume-encryption).

## Limitations
- Longhorn is unable to encrypt backing images that are already encrypted, and decrypt backing images that are not encrypted.
- Longhorn does not allow you to change the encryption key of an encrypted backing image.
- When encrypting a qcow2 image, Longhorn first creates a raw image from the qcow2 image and then encrypts it. The resulting encrypted raw image temporarily consumes extra space during cloning. For example,
    1. If we encrypt a 10MiB qcow2 image with a virtual size of 200MiB, we first create the raw image from the qcow2 which will consume 200MiB of the space.
    2. Longhorn then create the encrypted backing image from that 200MiB raw image which will take another 200MiB of the space.
    3. After the encrypted backing image is created, the temporary raw image will be cleaned up and free the 200MiB from the space.
- If the source backing image is a sparse file, the file loses its sparsity after encryption.
- To allow storage of the LUKS metadata during encryption, the image size is increased by 16 MB. For more information, see the [cryptsetup release notes](https://gitlab.com/cryptsetup/cryptsetup/-/blob/master/docs/v2.1.0-ReleaseNotes#L27).


---

## Article: advanced-resources/backing-image/backing-image.md

---
title: Backing Image
weight: 1
---

Longhorn natively supports backing images since v1.1.1.

A QCOW2 or RAW image can be set as the backing/base image of a Longhorn volume, which allows Longhorn to be integrated with a VM like [Harvester](https://github.com/rancher/harvester).

> **Important**:  
> The image size must be **a multiple of 512 bytes**. Longhorn uses direct I/O, which requires alignment of file sizes with the underlying storage block size.

## Create V1 Data Engine Backing Image

### Parameters during creation

#### The data source of a backing image
You can prepare a V1 Data Engine backing image using any of the supported data sources.
1. Download a backing image file (using a URL).
2. Upload a file from your local machine. This option is available to Longhorn UI users.
3. Export an existing in-cluster volume as a backing image.
4. Restore a backing image from the backupstore, For more information, see [Backing Image Backup](../backing-image-backup).
5. Clone a backing image.

#### Volume exporting

A backing image serves as the initial snapshot in the snapshot chain of a Longhorn volume. When you export a volume with an associated backing image, Longhorn merges that image with the delta changes, resulting in a new consolidated backing image.

#### The checksum of a backing image
- The checksum of a backing image is **the SHA512 checksum** of the whole backing image **file** rather than that of the actual content.
  What's the difference? When Longhorn calculates the checksum of a qcow2 file, it will read the file as a raw file instead of using the qcow library to read the correct content. In other words, users always get the correct checksum by executing `shasum -a 512 <the file path>` regardless of the file format.
- It's recommended to provide the expected checksum during backing image creation.
  Otherwise, Longhorn will consider the checksum of the first file as the correct one. Once there is something wrong with the first file preparation, which then leads to an incorrect checksum as the expected value, this backing image is probably unavailable.

#### Scheduling
- Longhorn first prepares and stores the backing image file on a random node and disk, and then duplicates the file to the disks that are storing the replicas.
- For improved space efficiency, you can add `nodeSelector` and `diskSelector` to force storing of backing image files on a specific set of nodes and disks.
- The replicas cannot be scheduled on nodes or disks where the backing image cannot be scheduled.

#### Number of copies
- You can add `minNumberOfCopies` to ensure that multiple backing image files exist in the cluster.
- You can adjust the `minNumberOfCopies` in the global setting to apply the default value to the BackingImage.

### The way of creating a backing image

#### Create a backing image via Longhorn UI
On **Advanced > Backing Images** page, users can create backing images with any kinds of data source.

#### Create a V1 Backing Image Using YAML
You can download a file or export an existing volume as a backing image via YAML.
It's better not to "upload" a file via YAML. Otherwise, you need to manually handle the data upload via HTTP requests.

Here are some examples:
```yaml
apiVersion: longhorn.io/v1beta2
kind: BackingImage
metadata:
  name: bi-download
  namespace: longhorn-system
spec:
  dataEngine: v1
  minNumberOfCopies: 2
  nodeSelector:
    - "node1"
  diskSelector:
    - "disk1"
  sourceType: download
  sourceParameters:
    url: https://longhorn-backing-image.s3-us-west-1.amazonaws.com/parrot.raw
  checksum: 304f3ed30ca6878e9056ee6f1b02b328239f0d0c2c1272840998212f9734b196371560b3b939037e4f4c2884ce457c2cbc9f0621f4f5d1ca983983c8cdf8cd9a
```
```yaml
apiVersion: longhorn.io/v1beta2
kind: BackingImage
metadata:
  name: bi-export
  namespace: longhorn-system
spec:
  dataEngine: v1
  minNumberOfCopies: 2
  nodeSelector:
    - "node1"
  diskSelector:
    - "disk1"
  sourceType: export-from-volume
  sourceParameters:
    volume-name: vol-export-src
    export-type: qcow2
```

#### Create and use a backing image via StorageClass and PVC
1. In a Longhorn StorageClass.
2. Setting parameter `backingImageName` means asking Longhorn to use this backing image during volume creation.
3. If you want to create the backing image as long as it does not exist during the CSI volume creation, parameters `backingImageDataSourceType` and `backingImageDataSourceParameters` should be set as well. Similar to YAML, it's better not to create a backing image via "upload" in StorageClass. Note that if all of these parameters are set and the backing image already exists, Longhorn will validate if the parameters matches the existing one before using it.
    - For `download`:
      ```yaml
      kind: StorageClass
      apiVersion: storage.k8s.io/v1
      metadata:
        name: longhorn-backing-image-example
      provisioner: driver.longhorn.io
      allowVolumeExpansion: true
      reclaimPolicy: Delete
      volumeBindingMode: Immediate
      parameters:
        numberOfReplicas: "3"
        staleReplicaTimeout: "2880"
        backingImage: "bi-download"
        backingImageDataSourceType: "download"
        backingImageDataSourceParameters: '{"url": "https://backing-image-example.s3-region.amazonaws.com/test-backing-image"}'
        backingImageChecksum: "SHA512 checksum of the backing image"
        backingImageMinNumberOfCopies: "2"
        backingImageNodeSelector: "node1"
        backingImageDiskSelector: "disk1"
        dataEngine: "v1"
      ```
    - For `export-from-volume`:
      ```yaml
      kind: StorageClass
      apiVersion: storage.k8s.io/v1
      metadata:
        name: longhorn-backing-image-example
      provisioner: driver.longhorn.io
      allowVolumeExpansion: true
      reclaimPolicy: Delete
      volumeBindingMode: Immediate
      parameters:
        numberOfReplicas: "3"
        staleReplicaTimeout: "2880"
        backingImage: "bi-export-from-volume"
        backingImageDataSourceType: "export-from-volume"
        backingImageDataSourceParameters: '{"volume-name": "vol-export-src", "export-type": "qcow2"}'
        backingImageMinNumberOfCopies: "2"
        backingImageNodeSelector: "node1"
        backingImageDiskSelector: "disk1"
        dataEngine: "v1"
      ```

4. Create a PVC with the StorageClass. Then the backing image will be created (with the Longhorn volume) if it does not exist.
5. Longhorn starts to prepare the backing images to disks for the replicas when a volume using the backing image is attached to a node.

#### Notice:
- Please be careful of the escape character `\` when you input a download URL in a StorageClass.
- A backing image that is created using a StorageClass has the same data engine as the volume.

## Utilize a backing image in a volume

Users can [directly create then immediately use a backing image via StorageClass](./#create-and-use-a-backing-image-via-storageclass-and-pvc),
or utilize an existing backing image as mentioned below.

#### Use an existing backing
##### Use an existing backing Image during volume creation
1. Click **Advanced > Backing Images** in the Longhorn UI.
2. Click **Create Backing Image** to create a backing image with a unique name and a valid URL.
3. Select a backing image from the list. The volume and the backing image must use the same data engine.
4. Longhorn starts to download the backing image to disks for the replicas when a volume using the backing image is attached to a node.

##### Use an existing backing Image during volume restore
1. Click `Backup` and pick up a backup volume for the restore.
2. As long as the backing image is already set for the backup volume, Longhorn will automatically choose the backing image during the restore.
3. Longhorn allows you to re-specify/override the backing image during the restore.

#### Download the backing image file to the local machine
Since v1.3.0, users can download existing backing image files to the local via UI.

#### Notice:
- Users need to make sure the backing image existence when they use UI to create or restore a volume with a backing image specified.
- Before downloading an existing backing image file to the local, users need to guarantee there is a ready file for it.
- Downloading of V2 Data Engine backing images is currently not supported.

## Create a V2 Data Engine Backing Image

Starting v1.8.0, you can create a backing image that is supported by the V2 Data Engine by configuring `Data Engine` in the YAML (through the UI or a StorageClass).

### Parameters During Creation

All parameters are the same as that of the V1 Data Engine backing image, except for `Data Engine`.

#### Backing Image Data Sources

You can prepare a V2 Data Engine backing image using any of the supported data sources.
- Download a backing image file (using a URL).
- Upload a file from your local machine. This option is available to Longhorn UI users.
- Export an existing in-cluster V1 Data Engine volume as a backing image.
- Restore a backing image from the backupstore. For more information, see [Backing Image Backup](../backing-image-backup).
- Clone a V1 backing image.

#### Notice

- The following operations are currently not supported:
  - Exporting from a V2 Data Engine volume
  - Cloning a V2 backing image
  - Backing up a V2 backing image
- Unlike the V1 Data Engine, which is file-based, the V2 Data Engine requires Longhorn to store the backing image data in an SPDK logical volume. As a result, for qcow2 images, Longhorn must first convert the qcow2 image to a raw format before storing the data to the V2 Data Engine backing image, enabling it to read the correct data.

## Clean up backing images

#### Clean up backing images in disks
- Longhorn automatically cleans up the unused backing image files in the disks based on [the setting `Backing Image Cleanup Wait Interval`](../../../references/settings#backing-image-cleanup-wait-interval). But Longhorn will retain at least one file in a disk for each backing image anyway.
- You can manually remove backing images from disks using the Longhorn UI. Go to **Setting** > **Backing Image**, and then click the name of a specific backing image. In the window that opens, select one or more disks and then click **Clean Up**.
- Once there is one replica in a disk using a backing image, no matter what the replica's current state is, the backing image file in this disk cannot be cleaned up.

#### Delete backing images
- The backing image can be deleted only when there is no volume using it.

## Backing image recovery
- If there is still a ready backing image file in one disk, Longhorn will automatically clean up the failed backing image files then re-launch these files from the ready one.
- If somehow all files of a backing image become failed, and the first file is :
  - downloaded from a URL, Longhorn will restart the downloading.
  - exported from an existing volume, Longhorn will (attach the volume if necessary then) restart the export.
  - uploaded from user local env, there is no way to recover it. Users need to delete this backing image then re-create a new one by re-uploading the file.
- When a node is down or the backing image manager pod on the node is unavailable, all backing image files on the node will become `unknown`. Later on if the node is back and the pod is running, Longhorn will detect then reuse the existing files automatically.

## Backing image eviction
- You can manually evict all backing image files from a node or disk by setting `Scheduling` to `Disabled` and `Eviction Requested` to `True` on the Longhorn UI.
- If only one backing image file exists in the cluster, Longhorn first duplicates the file to another disk and then deletes the file.
- If the backing image file cannot be duplicated to other disks, Longhorn does not delete the file. You can update the settings to resolve the issue.

## Backing image Workflow
1. To manage all backing image files in a disk, Longhorn will create one backing image manager pod for each disk. Once the disk has no backing image file requirement, the backing image manager will be removed automatically.
2. Once a backing image file is prepared by the backing image manager for a disk, the file will be shared among all volume replicas in this disk.
3. When a backing image is created, Longhorn will launch a backing image data source pod to prepare the first file. The file data is from the data source users specified (download from remote/upload from local/export from the volume). After the preparation done, the backing image manager pod in the same disk will take over the file then Longhorn will stop the backing image data source pod.
4. Once a new backing image is used by a volume, the backing image manager pods in the disks that the volume replicas reside on will be asked to sync the file from the backing image manager pods that already contain the file.
5. As mentioned in the section [#clean-up-backing-images-in-disks](#clean-up-backing-images-in-disks), the file will be cleaned up automatically if all replicas in one disk do not use one backing image file.

## Concurrent limit of backing image syncing
- `Concurrent Backing Image Replenish Per Node Limit` in the global settings controls how many backing images copies on a node can be replenished simultaneously.
- When set to 0, Longhorn won't replenish the copy automatically event it is less than the `minNumberOfCopies`

## Warning
- The download URL of the backing image should be public. We will improve this part in the future.
- If there is high memory usage of one backing image manager pod after [file download](#download-the-backing-image-file-to-the-local-machine), this is caused by the system cache/buffer. The memory usage will decrease automatically hence you don't need to worry about it. See [the GitHub ticket](https://github.com/longhorn/longhorn/issues/4055) for more details.

## History
* Available since v1.1.1 [Enable backing image feature in Longhorn](https://github.com/Longhorn/Longhorn/issues/2006)
* Support [upload](https://github.com/longhorn/longhorn/issues/2404) and [volume exporting](https://github.com/longhorn/longhorn/issues/2403) since v1.2.0.
* Support [download to local](https://github.com/longhorn/longhorn/issues/2404) and [volume exporting](https://github.com/longhorn/longhorn/issues/3155) since v1.3.0.


---

## Article: maintenance/_index.md

---
title: Maintenance and Upgrade
weight: 3
---

---

## Article: maintenance/maintenance.md

---
title: Node Maintenance and Kubernetes Upgrade Guide
weight: 3
---

This section describes how to handle planned node maintenance or upgrading Kubernetes version for the cluster.

- [Updating the Node OS or Container Runtime](#updating-the-node-os-or-container-runtime)
- [Removing a Disk](#removing-a-disk)
  - [Reusing the Node Name](#reusing-the-node-name)
- [Removing a Node](#removing-a-node)
- [Upgrading Kubernetes](#upgrading-kubernetes)
  - [In-place Upgrade](#in-place-upgrade)
  - [Managed Kubernetes](#managed-kubernetes)
- [Node Drain Policy Recommendations](#node-drain-policy-recommendations)
  - [Important Notes](#important-notes)
  - [Block If Contains Last Replica](#block-if-contains-last-replica)
  - [Allow If Last Replica Is Stopped](#allow-if-last-replica-is-stopped)
  - [Always Allow](#always-allow)
  - [Block For Eviction](#block-for-eviction)
  - [Block For Eviction If Contains Last Replica](#block-for-eviction-if-contains-last-replica)

## Updating the Node OS or Container Runtime

1. Cordon the node. Longhorn will automatically disable the node scheduling when a Kubernetes node is cordoned.

1. Drain the node to move the workload to somewhere else.

   It is necessary to use `--ignore-daemonsets` to drain the node. The `--ignore-daemonsets` is needed because Longhorn
   deployed some daemonsets such as `Longhorn manager`, `Longhorn CSI plugin`, `engine image`.

   While the drain proceeds, engine processes on the node will be migrated with the workload pods to other nodes.

   > **Note:** Volumes that are not attached through the CSI flow on the node (for example, manually attached using
   > UI) will not be automatically attached to new nodes by Kubernetes during the draining. Therefore, Longhorn will
   > prevent the node from completing the drain operation. The user will need to detach these volumes manually to
   > unblock the draining.

   While the drain proceeds, replica processes on the node will either continue to run or eventually be evicted and
   stopped based on the [Node Drain Policy](#node-drain-policy-recommendations).

   > **Note:** By default, if there is one last healthy replica for a volume on the node, Longhorn will prevent the node
   > from completing the drain operation, to protect the last replica and prevent the disruption of the workload. You
   > can control this behavior with the setting [Node Drain Policy](../../references/settings#node-drain-policy), or
   > [evict the replica to other nodes before draining](../../nodes-and-volumes/nodes/disks-or-nodes-eviction). See [Node Drain Policy
   > Recommendations](#node-drain-policy-recommendations) for considerations when selecting a policy.

   After the drain is completed, there should be no engine or replica processes running on the node, as the
   instance-manager pod that was running them will be stopped. Depending on the [Node Drain
   Policy](#node-drain-policy-recommendations), replicas scheduled to the node will either appear as `Failed` or be
   removed in favor of replacements. Workloads using Longhorn volumes will function as expected and enough replicas will
   be running elsewhere to meet the requirements of the policy.

   > **Note:** Normally you don't need to evict the replicas before the drain operation, as long as you have healthy
   > replicas on other nodes. The replicas can be reused later, once the node back online and uncordoned. See [Node
   > Drain Policy](#node-drain-policy-recommendations) for further guidance.

1. Perform the necessary maintenance, including shutting down or rebooting the node.
1. Uncordon the node. Longhorn will automatically re-enable the node scheduling. If there are existing replicas on the
   node, Longhorn might use those replicas to speed up the rebuilding process. You can set the [Replica Replenishment
   Wait Interval](../../references/settings#replica-replenishment-wait-interval) setting to customize how long Longhorn
   should wait for potentially reusable replica to be available.

## Removing a Disk

To remove a disk:

1. Disable the disk scheduling.
1. Evict all the replicas on the disk.
1. Delete the disk.

### Reusing the Node Name

These steps also apply if you've replaced a node using the same node name. Longhorn will recognize that the disks are
different once the new node is up. You will need to remove the original disks first and add them back for the new node
if it uses the same name as the previous node.

## Removing a Node

To remove a node:

1. Disable the disk scheduling.
1. Evict all the replicas on the node.
1. Detach all the volumes on the node.

   If the node has been drained, all the workloads should be migrated to another node already.

   If there are any other volumes remaining attached, detach them before continuing.

1. Remove the node from Longhorn using the `Delete` in the `Node` tab.

   Or, remove the node from Kubernetes, using:

      kubectl delete node <node-name>

1. Longhorn will automatically remove the node from the cluster.

## Upgrading Kubernetes

### In-place Upgrade

In-place upgrade is upgrading method in which nodes are upgraded without being removed from the cluster. Some example
solutions that use this upgrade methods are [k3s automated upgrades](https://docs.k3s.io/upgrades/automated), [Rancher's
Kubernetes upgrade
guide](https://rancher.com/docs/rancher/v2.x/en/cluster-admin/upgrading-kubernetes/#upgrading-the-kubernetes-version),
[Kubeadm upgrade](https://kubernetes.io/docs/tasks/administer-cluster/kubeadm/kubeadm-upgrade/), etc...

With the assumption that node and disks are not being deleted/removed, the recommended upgrading guide is:

1. Cordon and drain a node before upgrading Kubernetes components. Draining instructions are similar to the ones at
   [Updating the Node OS or Container Runtime](#updating-the-node-os-or-container-runtime).
2. The drain `--timeout` should be big enough so that replica rebuildings on healthy nodes can finish between node
   upgrades. The more Longhorn replicas you have on the draining node, the more time it takes for the Longhorn replicas
   to be rebuilt on other healthy nodes. We recommending you to test and select a big enough value or set it to 0 (aka
   never timeout).
3. The number of nodes upgrading at a time should be smaller than the number of Longhorn replicas for each volume.
   This is so that a running Longhorn volume has at least one healthy replica running at a time.
4. Consider setting the setting [Node Drain Policy](../../references/settings#node-drain-policy) to
   `allow-if-replica-is-stopped` so that the drain is not blocked by the last healthy replica of a detached volume. See
   [Node Drain Policy Recommendations](#node-drain-policy-recommendations) for considerations when selecting a policy.

### Managed Kubernetes

See the instruction at [Support Managed Kubernetes Service](../../advanced-resources/support-managed-k8s-service).

## Node Drain Policy Recommendations

There are currently five Node Drain Policies available for selection. Each has its own benefits and drawbacks. This
section provides general guidance on each and suggests situations in which each might be used.

### Important Notes

Node Drain Policy is intended to govern Longhorn behavior when a node is actively being drained. However, there is no
way for Longhorn to determine the difference between the cordoning and draining of a node, so, depending on the policy,
Longhorn may take action any time a node is cordoned, even if it is not being drained.

Node drain policy works to prevent the eviction of an instance-manager pod during a drain until certain conditions are
met. If the instance-manager pod cannot be evicted, the drain cannot complete. This prevents a user (or automated
process) from continuing to shut down or restart a node if it is not safe to do so. It may be tempting to ignore the
drain failure and proceed with maintenance operations if it seems to take too long, but this limits Longhorn's ability
to protect data. Always look at events and/or logs to try to determine WHY the drain is not progressing and take actions
to fix the underlying issue.

### Block If Contains Last Replica

This is the default policy. It is intended to provide a good balance between convenience and data protection. While it
is in effect, Longhorn will prevent the eviction of an instance-manager pod (and the completion of a drain) on a
cordoned node that contains the last healthy replica of a volume.

Benefits:

- Protects data by preventing the drain operation from completing until there is a healthy replica available for each
  volume available on another node.

Drawbacks:

- If there is only one replica for the volume, or if its other replicas are unhealthy, the user may need to manually
  (through the UI) request the eviction of replicas from the disk or node.
- Volumes may be degraded after the drain is complete. If the node is rebooted, redundancy is reduced until it is
  running again. If the node is removed, redundancy is reduced until another replica rebuilds.

### Allow If Last Replica Is Stopped

This policy is similar to `Block If Contains Last Replica`. It is inherently less safe, but can allow drains to complete
more quickly. It only prevents the eviction of an instance-manager pod (and the completion of a drain) on a node that
contains the last RUNNING healthy replica.

Benefits:

- Allows the drain operation to proceed in situations where the node being drained is expected to come back online
  (data will not be lost) and the replicas stored on the node's disks are not actively being used.

Drawbacks:

- Similar drawbacks to `Block If Contains Last Replica`.
- If, for some reason, the node never comes back, data is lost.

### Always Allow

This policy does not protect data in any way, but allows drains to immediately complete. It never prevents the eviction
of an instance-manager pod (and the completion of a drain). Do not use it in a production environment.

Benefits:

- The drain operation completes quickly without Longhorn getting in the way.

Drawbacks:

- There is no opportunity for Longhorn to protect data.

### Block For Eviction

This policy provides the maximum amount of data protection, but can lead to long drain times and unnecessary data
movement. It prevents the eviction of an instance-manager pod (and the completion of a drain) as long as any replicas
remain on a node. In addition, it takes action to automatically evict replicas from the node.

It is not recommended to leave this policy enabled under normal use, as it will trigger replica eviction any time a
node is cordoned. Only enable it during planned maintenance.

A primary use case for this policy is when automatically upgrading clusters in which volumes have no redundancy
(`numberOfReplicas == 1`). Other policies will prevent the drain until such replicas are manually evicted, which is
inconvenient for automation.

Benefits:

- Protects data by preventing the drain operation from completing until all replicas have been relocated.
- Automatically evicts replicas, so the user does not need to do it manually (through the UI).
- Maintains replica redundancy at all times.

Drawbacks:

- The drain operation is significantly slower than for other behaviors. Every replica must be rebuilt on another node
  before it can complete. Drain timeout must be adjusted as appropriate for the amount of data that will move during
  rebuilding.
- The drain operation is data-intensive, especially when replica auto balance is enabled, as evicted replicas may be
  moved back to the drained node when/if it comes back online.
- Like all of these policies, it triggers on cordon, not on drain. If a user regularly cordons nodes without draining
  them, replicas will be rebuilt pointlessly.

### Block For Eviction If Contains Last Replica

This policy provides the data protection of the default `Block If Contains Last Replica` with the added convenience of
automatic eviction. While it is in effect, Longhorn will prevent the eviction of an instance-manager pod (and the
completion of a drain) on a cordoned node that contains the last healthy replica of a volume. In addition, replicas that
meet this condition are automatically evicted from the node.

It is not recommended to leave this policy enabled under normal use, as it may trigger replica eviction any time a
node is cordoned. Only enable it during planned maintenance.

A primary use case for this policy is when automatically upgrading clusters in which volumes have no redundancy
(`numberOfReplicas == 1`). Other policies will prevent the drain until such replicas are manually evicted, which is
inconvenient for automation.

Benefits:

- Protects data by preventing the drain operation from completing until there is a healthy replica available for each
  volume available on another node.
- Automatically evicts replicas, so the user does not need to do it manually (through the UI).
- The drain operation is only as slow and data-intensive as is necessary to protect data.

Drawbacks:

- Volumes may be degraded after the drain is complete. If the node is rebooted, redundancy is reduced until it is
  running again. If the node is removed, redundancy is reduced until another replica rebuilds.
- Like all of these policies, it triggers on cordon, not on drain. If a user regularly cordons nodes without draining
  them, replicas will be rebuilt pointlessly.


---

## Article: v2-data-engine/_index.md

---
title: V2 Data Engine (Experimental)
weight: 10
aliases:
- /spdk/_index.md
---


---

## Article: v2-data-engine/performance.md

---
title: Performance
weight: 3
aliases:
- /spdk/performance.md
---

## Performance Measurement Tools

- [KBench](https://github.com/yasker/kbench): Used to benchmark cluster storage performance
- [Local Path Provisioner](https://github.com/rancher/local-path-provisioner): Used to measure the baseline performance of the data disk

## Equinix (m3.small.x86)

- Machine: Japan/m3.small.x86
- CPU: Intel(R) Xeon(R) E-2378G CPU @ 2.80GHz
- RAM: 64 GiB
- Kubernetes: v1.23.6+rke2r2
- Nodes: 3 (each node is a master and also a worker)
- OS: Ubuntu 22.04 / 5.15.0-33-generic
- Storage: 1 SSD (Micron_5300_MTFD)
- Network throughput between nodes (tested by iperf over 60 seconds): 15.0 Gbits/sec

{{< figure src="/img/diagrams/v2-data-engine/equinix-iops.svg" >}}

{{< figure src="/img/diagrams/v2-data-engine/equinix-bw.svg" >}}

{{< figure src="/img/diagrams/v2-data-engine/equinix-latency.svg" >}}

# AWS EC2 (c5d.xlarge)

- Machine: Tokyo/c5d.xlarge
- CPU: Intel(R) Xeon(R) Platinum 8124M CPU @ 3.00GHz
- RAM: 8 GiB
- Kubernetes: v1.25.10+rke2r1
- Nodes: 3 (each node is a master and also a worker)
- OS: Ubuntu 22.04.2 LTS / 5.19.0-1025-aws
- Storage: 1 SSD (Amazon EC2 NVMe Instance Storage/Local NVMe Storage)
- Network throughput between nodes (tested by iperf over 60 seconds): 7.9 Gbits/sec

{{< figure src="/img/diagrams/v2-data-engine/aws-c5d-xlarge-iops.svg" >}}

{{< figure src="/img/diagrams/v2-data-engine/aws-c5d-xlarge-bw.svg" >}}

{{< figure src="/img/diagrams/v2-data-engine/aws-c5d-xlarge-latency.svg" >}}


---

## Article: v2-data-engine/prerequisites.md

---
title: Prerequisites
weight: 1
aliases:
- /spdk/prerequisites.md
---

## Prerequisites

Longhorn nodes must meet the following requirements:

- AMD64 or ARM64 CPU
  > **NOTICE**
  >
  >  AMD64 CPUs require SSE4.2 instruction support.

- Linux kernel

  5.19 or later is required for NVMe over TCP support
  > **NOTICE**
  >
  > Host machines with Linux kernel 5.15 may unexpectedly reboot when volume-related IO errors occur. Update the Linux kernel on Longhorn nodes to version 5.19 or later to prevent such issues.

  v6.7 or later is recommended for improved system stability
  > **NOTICE**
  >
  > Memory corruption may occur on hosts using versions of the Linux kernel earlier than 6.7, as highlighted by this SPDK upstream issue: https://github.com/spdk/spdk/issues/3116#issuecomment-1890984674. In Longhorn environments the kernel panic can be caused by prevalent IO timeouts in communications between the `nvme-tcp` driver and SPDK. Update the Linux kernel on Longhorn nodes to version 6.7 or later to prevent the issue from occurring.

- Linux kernel modules
  - `vfio_pci`
  - `uio_pci_generic`
  - `nvme-tcp`

- Huge page support
  - 2 GiB of 2 MiB-sized pages

## Notice

### CPU

When the V2 Data Engine is enabled, each instance-manager pod utilizes **1 CPU core**. This high CPU usage is attributed to the `spdk_tgt` process running within each instance-manager pod. The spdk_tgt process is responsible for handling input/output (IO) operations and requires intensive polling. As a result, it consumes 100% of a dedicated CPU core to efficiently manage and process the IO requests, ensuring optimal performance and responsiveness for storage operations.

### Memory

SPDK leverages huge pages for enhancing performance and minimizing memory overhead. You must configure 2 MiB-sized huge pages on each Longhorn node to enable usage of huge pages. Specifically, 1024 pages (equivalent to a total of 2 GiB) must be available on each Longhorn node.


### Disk

SPDK leverages kernel drivers to support every kind of disk that Linux supports. However, SPDK is equipped with a user space NVMe driver that provides zero-copy, highly parallel, direct access to an SSD from a user space application. Because of this, using **local NVMe disks** is highly recommended for enabling V2 volumes to achieve optimal storage performance.

---

## Article: v2-data-engine/quick-start.md

---
title: Quick Start
weight: 2
aliases:
- /spdk/quick-start.md
---

**Table of Contents**
- [Prerequisites](#prerequisites)
  - [Configure Kernel Modules and Huge Pages](#configure-kernel-modules-and-huge-pages)
  - [Load `nvme-tcp` Kernel Module](#load-nvme-tcp-kernel-module)
  - [Load Kernel Modules Automatically on Boot](#load-kernel-modules-automatically-on-boot)
  - [Restart `kubelet`](#restart-kubelet)
  - [Check Environment](#check-environment)
    - [Using the Longhorn Command Line Tool](#using-the-longhorn-command-line-tool)
    - [Use Longhorn Command Line Tool](#use-longhorn-command-line-tool)
- [Installation](#installation)
  - [Install Longhorn System](#install-longhorn-system)
  - [Enable V2 Data Engine](#enable-v2-data-engine)
  - [CPU and Memory Usage](#cpu-and-memory-usage)
  - [Add `block-type` Disks in Longhorn Nodes](#add-block-type-disks-in-longhorn-nodes)
    - [Prepare disks](#prepare-disks)
    - [Add disks to `node.longhorn.io`](#add-disks-to-nodelonghornio)
- [Application Deployment](#application-deployment)
  - [Create a StorageClass](#create-a-storageclass)
  - [Create Longhorn Volumes](#create-longhorn-volumes)

---

Longhorn's V2 Data Engine harnesses the power of the Storage Performance Development Kit (SPDK) to elevate its overall performance. The integration significantly reduces I/O latency while simultaneously boosting IOPS and throughput. The enhancement provides a high-performance storage solution capable of meeting diverse workload demands.

**V2 Data Engine is currently an experimental feature and should NOT be utilized in a production environment.** At present, a volume with V2 Data Engine only supports

- Volume lifecycle (creation, attachment, detachment and deletion)
- Degraded volume
- Block disk management
- Orphaned replica management

In addition to the features mentioned above, additional functionalities such as replica number adjustment, online replica rebuilding, snapshot, backup, restore and so on will be introduced in future versions.

This tutorial will guide you through the process of configuring the environment and create Kubernetes persistent storage resources of persistent volumes (PVs) and persistent volume claims (PVCs) that correspond to Longhorn volumes using V2 Data Engine.

## Prerequisites

### Configure Kernel Modules and Huge Pages

For Debian and Ubuntu, please install Linux kernel extra modules before loading the kernel modules
```
apt install -y linux-modules-extra-`uname -r`
```

To configure the necessary kernel modules and huge pages for SPDK, you can use the [Longhorn CLI](../../advanced-resources/longhornctl/). 

Or, you can install them manually by following these steps.
- Load the kernel modules on the each Longhorn node
  ```
  modprobe vfio_pci
  modprobe uio_pci_generic
  ```

- Configure huge pages
SPDK leverages huge pages for enhancing performance and minimizing memory overhead. You must configure 2 MiB-sized huge pages on each Longhorn node to enable usage of huge pages. Specifically, 1024 pages (equivalent to a total of 2 GiB) must be available on each Longhorn node.

To allocate huge pages, run the following commands on each node.
  ```
  echo 1024 > /sys/kernel/mm/hugepages/hugepages-2048kB/nr_hugepages
  ```

  To make the change permanent, add the following line to the file /etc/sysctl.conf.
  ```
  echo "vm.nr_hugepages=1024" >> /etc/sysctl.conf
  ```

### Load `nvme-tcp` Kernel Module

To ensure the necessary prerequisites for NVMe are met, you can use the [Longhorn CLI](../../advanced-resources/longhornctl/).

Alternatively, you can manually load the `nvme-tcp` kernel module on each Longhorn node by running the following command:
```
modprobe nvme-tcp
```

### Load Kernel Modules Automatically on Boot

Rather than manually loading kernel modules `vfio_pci`, `uio_pci_generic` and `nvme-tcp` each time after reboot, you can streamline the process by configuring automatic module loading during the boot sequence. For detailed instructions, please consult the manual provided by your operating system.

Reference:
- [SUSE/OpenSUSE: Loading kernel modules automatically on boot](https://documentation.suse.com/sles/15-SP4/html/SLES-all/cha-mod.html#sec-mod-modprobe-d)
- [Ubuntu: Configure kernel modules to load at boot](https://manpages.ubuntu.com/manpages/jammy/man5/modules-load.d.5.html)
- [RHEL: Loading kernel modules automatically at system boot time](https://access.redhat.com/documentation/zh-tw/red_hat_enterprise_linux/8/html/managing_monitoring_and_updating_the_kernel/managing-kernel-modules_managing-monitoring-and-updating-the-kernel)

### Restart `kubelet`

After finishing the above steps, restart kubelet on each node.

### Check Environment

#### Using the Longhorn Command Line Tool

The `longhornctl` tool is a CLI for Longhorn operations. For more information, see [Command Line Tool (longhornctl)](../../advanced-resources/longhornctl/).

To check the prerequisites and configurations, download the tool and run the `check` sub-command:

```shell
# For AMD64 platform
curl -sSfL -o longhornctl https://github.com/longhorn/cli/releases/download/v{{< current-version >}}/longhornctl-linux-amd64
# For ARM platform
curl -sSfL -o longhornctl https://github.com/longhorn/cli/releases/download/v{{< current-version >}}/longhornctl-linux-arm64

chmod +x longhornctl
./longhornctl check preflight --enable-spdk
```

Example of result:

```shell
INFO[2024-01-10T00:00:01Z] Initializing preflight checker
INFO[2024-01-01T00:00:01Z] Cleaning up preflight checker
INFO[2024-01-01T00:00:01Z] Running preflight checker
INFO[2024-01-01T00:00:02Z] Retrieved preflight checker result:
worker1:
  error:
  - 'HugePages is insufficient. Required 2MiB HugePages: 1024 pages, Total 2MiB HugePages: 0 pages'
  - 'Module nvme_tcp is not loaded: failed to execute: nsenter [--mount=/host/proc/204896/ns/mnt --net=/host/proc/204896/ns/net grep nvme_tcp /proc/modules], output , stderr : exit status 1'
  - 'Module uio_pci_generic is not loaded: failed to execute: nsenter [--mount=/host/proc/204896/ns/mnt --net=/host/proc/204896/ns/net grep uio_pci_generic /proc/modules], output , stderr : exit status 1'
  info:
  - Service iscsid is running
  - NFS4 is supported
  - Package nfs-common is installed
  - Package open-iscsi is installed
  - CPU instruction set sse4_2 is supported
  warn:
  - multipathd.service is running. Please refer to https://longhorn.io/kb/troubleshooting-volume-with-multipath/ for more information.
```

Use the `install` sub-command to install and set up the preflight dependencies before installing Longhorn.

```shell
master:~# ./longhornctl install preflight --enable-spdk
INFO[2024-01-01T00:00:03Z] Initializing preflight installer
INFO[2024-01-01T00:00:03Z] Cleaning up preflight installer
INFO[2024-01-01T00:00:03Z] Running preflight installer
INFO[2024-01-01T00:00:03Z] Installing dependencies with package manager
INFO[2024-01-01T00:00:10Z] Installed dependencies with package manager
INFO[2024-01-01T00:00:10Z] Cleaning up preflight installer
INFO[2024-01-01T00:00:10Z] Completed preflight installer. Use 'longhornctl check preflight' to check the result.
```

> **Note**:
> Some immutable Linux distributions, such as SUSE Linux Enterprise Micro (SLE Micro), require you to reboot worker nodes after running the `install` sub-command. After the reboot, you must run the `install` sub-command again to complete the operation.
>
> The documentation for your Linux distribution should outline such requirements. For example, the [SLE Micro documentation](https://documentation.suse.com/sle-micro/6.0/html/Micro-transactional-updates/index.html#reference-transactional-update-usage) explains that all changes made by the `transactional-update` command become active only after the node is rebooted.

After installing and setting up the preflight dependencies, you can run the `check` sub-command again to verify that all environment settings are correct.

```shell
master:~# ./longhornctl check preflight --enable-spdk
INFO[2024-01-01T00:00:13Z] Initializing preflight checker
INFO[2024-01-01T00:00:13Z] Cleaning up preflight checker
INFO[2024-01-01T00:00:13Z] Running preflight checker
INFO[2024-01-01T00:00:16Z] Retrieved preflight checker result:
worker1:
  info:
  - Service iscsid is running
  - NFS4 is supported
  - Package nfs-common is installed
  - Package open-iscsi is installed
  - CPU instruction set sse4_2 is supported
  - HugePages is enabled
  - Module nvme_tcp is loaded
  - Module uio_pci_generic is loaded
```

#### Use Longhorn Command Line Tool

Make sure everything is correctly configured and installed by
```
longhornctl --kube-config ~/.kube/config --image longhornio/longhorn-cli:v{{< current-version >}} install preflight --enable-spdk
```

See [Longhorn Command Line Tool](../../advanced-resources/longhornctl/) for more information.

## Installation

### Install Longhorn System

Follow the steps in Quick Installation to install Longhorn system.

### Enable V2 Data Engine

Enable the V2 Data Engine by changing the `v2-data-engine` setting to `true` after installation. Following this, the instance-manager pods will be automatically restarted.

Or, you can enable it in `Settings > V2 Data Engine`.

### CPU and Memory Usage

When the V2 Data Engine is enabled, each Instance Manager pod for the V2 Data Engine uses 1 CPU core. The high CPU usage is caused by `spdk_tgt`, a process running in each Instance Manager pod that handles input/output (IO) operations and requires intensive polling. `spdk_tgt` consumes 100% of a dedicated CPU core to efficiently manage and process the IO requests, ensuring optimal performance and responsiveness for storage operations.

```
NAME                                                CPU(cores)   MEMORY(bytes)
csi-attacher-57c5fd5bdf-jsfs4                       1m           7Mi
csi-attacher-57c5fd5bdf-kb6dv                       1m           9Mi
csi-attacher-57c5fd5bdf-s7fb6                       1m           7Mi
csi-provisioner-7b95bf4b87-8xr6f                    1m           11Mi
csi-provisioner-7b95bf4b87-v4gwb                    1m           9Mi
csi-provisioner-7b95bf4b87-vnt58                    1m           9Mi
csi-resizer-6df9886858-6v2ds                        1m           8Mi
csi-resizer-6df9886858-b6mns                        1m           9Mi
csi-resizer-6df9886858-l4vmj                        1m           8Mi
csi-snapshotter-5d84585dd4-4dwkz                    1m           7Mi
csi-snapshotter-5d84585dd4-km8bc                    1m           9Mi
csi-snapshotter-5d84585dd4-kzh6w                    1m           7Mi
engine-image-ei-b907910b-79k2s                      3m           19Mi
instance-manager-214803c4f23376af5a75418299b12ad6   1015m        133Mi (for V2 Data Engine)
instance-manager-4550bbc4938ff1266584f42943b511ad   4m           15Mi  (for V1 Data Engine)
longhorn-csi-plugin-nz94f                           1m           26Mi
longhorn-driver-deployer-556955d47f-h5672           1m           12Mi
longhorn-manager-2n9hd                              4m           42Mi
longhorn-ui-58db78b68-bzzz8                         0m           2Mi
longhorn-ui-58db78b68-ffbxr                         0m           2Mi
```


You can observe the utilization of allocated huge pages on each node by running the command `kubectl get node <node name> -o yaml`.
```
# kubectl get node sles-pool1-07437316-4jw8f -o yaml
...

status:
  ...
  allocatable:
    cpu: "8"
    ephemeral-storage: "203978054087"
    hugepages-1Gi: "0"
    hugepages-2Mi: 2Gi
    memory: 31813168Ki
    pods: "110"
  capacity:
    cpu: "8"
    ephemeral-storage: 209681388Ki
    hugepages-1Gi: "0"
    hugepages-2Mi: 2Gi
    memory: 32861744Ki
    pods: "110"
...
```

### Add `block-type` Disks in Longhorn Nodes

Unlike `filesystem-type` disks that are designed for legacy volumes, volumes using V2 Data Engine are persistent on `block-type` disks. Therefore, it is necessary to equip Longhorn nodes with `block-type` disks.

#### Prepare disks

If there are no additional disks available on the Longhorn nodes, you can create loop block devices to test the feature. To accomplish this, execute the following command on each Longhorn node to create a 10 GiB block device.
```
dd if=/dev/zero of=blockfile bs=1M count=10240
losetup -f blockfile
```

To display the path of the block device when running the command `losetup -f blockfile`, use the following command.
```
losetup -j blockfile
```

#### Add disks to `node.longhorn.io`

You can add the disk by navigating to the Node UI page and specify the `Disk Type` as `Block`. Next, provide the block device's path in the `Path` field.

Or, edit the `node.longhorn.io` resource.
```
kubectl -n longhorn-system edit node.longhorn.io <NODE NAME>
```

Add the disk to `Spec.Disks`
```
<DISK NAME>:
  allowScheduling: true
  evictionRequested: false
  path: /PATH/TO/BLOCK/DEVICE
  storageReserved: 0
  tags: []
  diskType: block
```

Wait for a while, you will see the disk is displayed in the `Status.DiskStatus`.

## Application Deployment

After the installation and configuration, we can dynamically provision a Persistent Volume using V2 Data Engine as the following steps.

### Create a StorageClass

Run the following command to create a StorageClass named `longhorn-spdk`. Set `parameters.dataEngine` to `v2` to enable the V2 Data Engine.
```
kubectl apply -f https://raw.githubusercontent.com/longhorn/longhorn/v{{< current-version >}}/examples/v2/storageclass.yaml
```

### Create Longhorn Volumes

Create a Pod that uses Longhorn volumes using V2 Data Engine by running this command:
```
kubectl apply -f https://raw.githubusercontent.com/longhorn/longhorn/v{{< current-version >}}/examples/v2/pod_with_pvc.yaml
```

Or, if you are creating a volume on Longhorn UI, please specify the `Data Engine` as `v2`.


---

## Article: v2-data-engine/troubleshooting.md

---
title: Troubleshooting
weight: 4
aliases:
- /spdk/troubleshooting.md
---

- [Installation](#installation)
  - ["Package 'linux-modules-extra-x.x.x-x-generic' Has No Installation Candidate" Error During Installation on Debian Machines](#package-linux-modules-extra-xxx-x-generic-has-no-installation-candidate-error-during-installation-on-debian-machines)
- [Disk](#disk)
  - ["Invalid argument" Error in Disk Status After Adding a Block-Type Disk](#invalid-argument-error-in-disk-status-after-adding-a-block-type-disk)

---

## Installation

### "Package 'linux-modules-extra-x.x.x-x-generic' Has No Installation Candidate" Error During Installation on Debian Machines

For Debian machines, if you encounter errors similar to the below when installing Linux kernel extra modules, you need to find an available version in the pkg collection websites like [this](https://pkgs.org/search/?q=linux-modules-extra) rather than directly relying on `uname -r` instead:
```log
apt install -y linux-modules-extra-`uname -r`
Reading package lists... Done
Building dependency tree... Done
Reading state information... Done
Package linux-modules-extra-5.15.0-67-generic is not available, but is referred to by another package.
This may mean that the package is missing, has been obsoleted, or
is only available from another source

E: Package 'linux-modules-extra-5.15.0-67-generic' has no installation candidate
```

For example, for Ubuntu 22.04, one valid version is `linux-modules-extra-5.15.0-76-generic`:
```shell
apt update -y
apt install -y linux-modules-extra-5.15.0-76-generic
```

## Disk

### "Invalid argument" Error in Disk Status After Adding a Block-Type Disk

After adding a block-type disk, the disk status displays error messages:
```
Disk disk-1(/dev/nvme1n1) on node dereksu-ubuntu-pool1-bf77ed93-2d2p9 is not ready: 
failed to generate disk config: error: rpc error: code = Internal desc = rpc error: code = Internal 
desc = failed to add block device: failed to create AIO bdev: error sending message, id 10441, 
method bdev_aio_create, params {disk-1 /host/dev/nvme1n1 4096}: {"code": -22,"message": "Invalid argument"}
```

Next, inspect the log message of the instance-manager pod on the same node. If the log reveals the following:
```
[2023-06-29 08:51:53.762597] bdev_aio.c: 762:create_aio_bdev: *WARNING*: Specified block size 4096 does not match auto-detected block size 512
[2023-06-29 08:51:53.762640] bdev_aio.c: 788:create_aio_bdev: *ERROR*: Disk size 100000000000 is not a multiple of block size 4096
```
These messages indicate that the size of your disk is not a multiple of the block size 4096 and is not supported by Longhorn system.

To resolve this issue, you can follow the steps
1. Remove the newly added block-type disk from the node.
2. Partition the block-type disk using the `fdisk` utility and ensure that the partition size is a multiple of the block size 4096.
3. Add the partitioned disk to the Longhorn node.



---

## Article: v2-data-engine/features/_index.md

---
title: Features
weight: 5
aliases:
- /spdk/features/_index.md
---

- Support for AMD64 and ARM64 platforms
- Volume lifecycle (creation, attachment, detachment and deletion)
- Degraded volume
- [Block disk management](./node-disk-support)
- Orphaned replica management
- Snapshot creation, deletion and reversion
- Volume backup and restoration
- [Selective V2 Data Engine activation](./selective-v2-data-engine-activation)
- [Filesystem Trim](../../nodes-and-volumes/volumes/trim-filesystem)
- [Backing Image](../../advanced-resources/backing-image/backing-image)
- [Volume Encryption](../../advanced-resources/security/volume-encryption)


In addition to the features mentioned above, additional functionalities such as replica number adjustment, online replica rebuilding and so on will be introduced in future versions.



---

## Article: v2-data-engine/features/configurable-cpu-cores.md

---
title: Configurable CPU Cores
weight: 20
aliases:
- /spdk/features/configurable-cpu-cores.md
---

Longhorn now supports configurable CPU cores for the v2 data engine, offering both global and per-node configuration options.

## Global Configuration

To set CPU cores globally, update the [data-engine-cpu-mask](../../../references/settings#data-engine-cpu-mask) setting using a hexadecimal encoded string. For example:

- Use 0x01 to allocate 1 core
- Use 0x03 to allocate 2 cores
- Use 0x07 to allocate 3 cores

## Per-node Configuration

For node-specific CPU core allocation, update the `spec.dataEngineSpec.v2.cpuMask` field of the instance manager with a hexadecimal encoded string. By default, this value is empty, and the v2 data engine will use the global setting specified by `data-engine-cpu-mask`. When a per-node configuration is set, the v2 data engine will prioritize this value over the global setting for that specific node.


---

## Article: v2-data-engine/features/interrupt-mode.md

---
title: Interrupt Mode Support
weight: 20
aliases:
- /spdk/features/interrupt-support.md
---

Starting with v1.10.0, Longhorn supports **SPDK interrupt mode** for V2 data engine volumes. Interrupt mode provides an alternative to the default **polling mode**, offering improved CPU efficiency in certain environment.

Interrupt mode is particularly suitable for clusters with limited CPU resources and a relatively small number of volumes. While polling mode maximizes performance by keeping CPU utilization close to 100% on allocated cores, interrupt mode reduces CPU usage by allowing the SPDK reactor to adjust its usage dynamically instead of continuously polling.

## Overview

### Polling Mode vs Interrupt Mode

**Polling Mode (Default)**:
- Continuously polls for I/O operations
- Provides the lowest latency
- Consumes ~100% of the allocated CPU core at all times
- Best suited for high-performance workloads with frequent I/O

**Interrupt Mode**:
- Uses interrupt-driven I/O handling
- CPU consumption scales with the number of attached volumes
- Better suited for resource-constrained environments

## Prerequisites

- Longhorn v1.10.0 or later
- V2 data engine enabled
- No attached v2 volumes when changing the setting
- For NVMe disks, IOMMU must be enabled. To verify:
    ```bash
    find /sys/kernel/iommu_groups/ -type l
    ```
    Example output (IOMMU enabled):
    ```
    /sys/kernel/iommu_groups/0/devices/0000:e6:0b.1
    /sys/kernel/iommu_groups/1/devices/0000:34:0a.6
    /sys/kernel/iommu_groups/2/devices/0000:a0:00.0
    ```
    If the command returns no output, IOMMU is not enabled.

    > **Note:** IOMMU support may not be exposed on virtualized instances. If unsure, consider using a bare-metal instance, or consult your cloud provider’s documentation or support team.

    For more information, see the official [SPDK documentation](https://spdk.io/doc/system_configuration.html).

## Configuration

### Global Setting

To enable interrupt mode globally, update the [data-engine-interrupt-mode-enabled](../../../references/settings#data-engine-interrupt-mode-enabled) setting.

### Important Considerations

- **Volume State Requirement**: The setting can only be changed when no V2 volumes are attached. Longhorn blocks updates if any V2 volume is active.
- **Global Effect**: The setting applies to all V2 volumes.

## Performance Characteristics

### Recommended Use Cases

Enable interrupt mode when:
- Running in resource-constrained clusters
- Managing only a small number of volumes
- CPU resources are limited or shared with other workloads
- I/O patterns are sporadic rather than continuous
- Energy efficiency is a priority

## Limitations

### Hybrid Implementation

The current V2 volume interrupt mode uses a hybrid approach for NVMe/TCP transport:

- **Admin Queue Operations**: Still relies on periodic polling for keepalive and controller recovery
- **I/O Queue Completion**: Uses polling for command completion
- **Residual CPU Usage**: Results in a small but constant CPU load, even when attach volumes are idle

### Performance Trade-offs

- **Latency**: Slightly higher than polling mode

### Operational Restrictions

- **Setting Changes**: Cannot be modified while V2 volumes are attached
- **Global Scope**: Applies globally; no per-volume override is available


---

## Article: v2-data-engine/features/node-disk-support.md

---
title: Node Disk Support
weight: 20
aliases:
- /spdk/features/node-disk-support.md
---

Longhorn now supports the addition and management of various disk types (AIO, NVMe, and VirtIO) on nodes, enhancing filesystem operations, storage performance, and compatibility.

- Enhanced Storage Performance

  Utilizing NVMe and VirtIO disks allows for faster disk operations, significantly improving overall performance.

- Filesystem Compatibility

  Disks managed with NVMe or VirtIO drivers offer better filesystem support, including advanced operations like trimming.

- Flexibility

  Users can select the disk type that best fits their environment: AIO for traditional setups, NVMe for high-performance needs, or VirtIO for virtualized environments.

- Ease of Management

  Automatic detection of disk drivers simplifies the addition and management of disks, reducing administrative overhead.

## Configure a Disk on Longhorn Node

Longhorn automatically detects the disk type if `node.spec.disks[i].diskDriver` is set to `auto`, optimizing storage performance. The detection and management is as follows:

- NVMe Disk: managed by spdk_tgt using the nvme bdev driver, and `node.status.diskStatus[i].diskDriver` is set to `nvme`.
- VirtIO Disk: managed by spdk_tgt using the virtio bdev driver, and `node.status.diskStatus[i].diskDriver` is set to `virtio-blk`.
- Other Disks: managed by spdk_tgt using the aio bdev driver, and `node.status.diskStatus[i].diskDriver` is set to `aio`.

Alternatively, users can manually set `node.spec.disks[i].diskDriver` to `aio` to force the use of the aio bdev driver.

To support NVMe and VirtIO disks, you need to find the BDF (Bus, Device, Function) of the disk as a disk path to be added to the Longhorn node. The following examples provide an introduction to configuring NVMe disks, VirtIO disks, and others.

> **Note**
> 
> Once these disks are managed by the NVMe bdev driver or VirtIO bdev driver, instead of the Linux kernel driver, they will no be listed under /dev/nvmeXnY or /dev/vdbX.

### Using NVMe Disks

1. List the disks

   First, identify the NVMe disks available on your system by running the following command:

   ```
   # ls -al /sys/block/
   ```

   Example output:
   ```
   lrwxrwxrwx  1 root root 0  Jul  30 12:20 loop0 -> ../devices/virtual/block/loop0
   lrwxrwxrwx  1 root root 0  Jul  30 12:20 nvme0n1 -> ../devices/pci0000:00/0000:00:01.2/0000:02:00.0/nvme/nvme0/nvme0n1
   lrwxrwxrwx  1 root root 0  Jul  30 12:20 nvme0n1 -> ../devices/pci0000:00/0000:00:01.2/0000:05:00.0/nvme/nvme1/nvme1n1
   ```

1. Get the BDF of the NVMe disk

   Identify the BDF of the NVMe disk `/dev/nvme1n1`. From the example above, the BDF is `0000:05:00.0`.

1. Add the NVMe disk to `spec.disks` of `node.longhorn.io`

   ```
   nvme-disk:
     allowScheduling: true
     diskType: block
     diskDriver: auto
     evictionRequested: false
     path: 0000:05:00.0
     storageReserved: 0
     tags: []
   ```

1. Check the `status.diskStatus`. The disk should be detected without errors, and the diskDriver should be set to `nvme`.

> **Note: Alternative Disk Configuration**
> 
> If you add the disk using a different path, such as:
> 
>  ```
>  nvme-disk:
>    allowScheduling: true
>    diskType: block
>    diskDriver: auto
>    evictionRequested: false
>    path: /dev/nvme1n1
>    storageReserved: 0
>    tags: []
>  ```
> In this case, the disk will be managed by the aio bdev driver, and the `node.status.diskStatus[i].diskDriver` is set to `aio`.

### Using VirtIO Disks

The steps are similar to NVMe disks.

1. List the disks

   First, identify the VirtIO disks available on your system by running the following command:

   ```
   # ls -al /sys/block/
   ```

   Example output:

   ```
   lrwxrwxrwx  1 root root 0  Jul  30 12:20 loop0 -> ../devices/virtual/block/loop0
   lrwxrwxrwx  1 root root 0 Feb 22 14:04 vda -> ../devices/pci0000:00/0000:00:02.3/0000:04:00.0/virtio2/block/vda
   lrwxrwxrwx  1 root root 0 Feb 22 14:24 vdb -> ../devices/pci0000:00/0000:00:02.6/0000:07:00.0/virtio5/block/vdb
   ```

1. Get the BDF of the VirtIO disk

   Identify the BDF of the VirtIO disk `/dev/vdb`. From the example above, the BDF is `0000:07:00.0`.

1. Add the NVMe disk to `spec.disks` of `node.longhorn.io`

   ```
   nvme-disk:
     allowScheduling: true
     diskType: block
     diskDriver: auto
     evictionRequested: false
     path: 0000:07:00.0
     storageReserved: 0
     tags: []
   ```

1. Check the `status.diskStatus`. The disk should be detected without errors, and the `diskDriver` should be set to `virtio-blk`.

> **Note: Alternative Disk Configuration**
> 
> If you add the disk using a different path, such as:
> 
>  ```
>  nvme-disk:
>    allowScheduling: true
>    diskType: block
>    diskDriver: auto
>    evictionRequested: false
>    path: /dev/vdb
>    storageReserved: 0
>    tags: []
>  ```
> In this case, the disk will be managed by the aio bdev driver, and the `node.status.diskStatus[i].diskDriver` is set to `aio`.


### Using AIO Disks

When neither NVMe nor VirtIO drivers can manage a disk, Longhorn will default to using the aio bdev driver. Users can also manually configure this.

1. Add the disk to `spec.disks` of `node.longhorn.io`

    ```
    default-disk-loop:
      allowScheduling: true
      diskDriver: aio
      diskType: block
      evictionRequested: false
      path: /dev/loop12
      storageReserved: 0
      tags: []
    ```

1. Check node.status.diskStatus. The disk should be detected without errors, and the `node.status.diskStatus[i].diskDriver` is set to `aio`.

## History

[Original Feature Request](https://github.com/longhorn/longhorn/issues/7672)


---

## Article: v2-data-engine/features/replica-rebuild-qos.md

---
title: Replica Rebuild QoS
weight: 5
---

Longhorn supports rebuild bandwidth throttling (QoS) for v2 volumes based on SPDK. This feature allows users to apply bandwidth limits to replicas during rebuilding to avoid overloading the source and destination node’s storage throughput.

## Global Setting: `v2-data-engine-rebuilding-mbytes-per-second`

* A cluster-wide setting that defines the maximum write bandwidth (in MB/s) for rebuilding replicas.
* When set to `0`, there is no limit.
* This setting can only be configured via kubectl:

```bash
kubectl -n longhorn-system patch settings.longhorn.io v2-data-engine-rebuilding-mbytes-per-second \
  --type=merge -p '{"value":"100"}'
```

## Per-Volume QoS Override

You can override the global rebuild bandwidth limit per volume by setting `spec.rebuildingMbytesPerSecond` in the `volume` spec:

```yaml
spec:
  rebuildingMbytesPerSecond: 50
```

## Effective QoS Resolution

The effective rebuild bandwidth limit is determined by evaluating both global and volume-specific settings. If the volume-specific value is greater than zero, it overrides the global setting.

| Global Setting | Volume Override | Effective QoS |
| -------------- | --------------- | ------------- |
| 0              | 0               | No limit      |
| 100            | 0               | 100 MB/s      |
| 0              | 200             | 200 MB/s      |
| 100            | 200             | 200 MB/s      |

The applied QoS is recorded in the field `status.rebuildStatus[*].appliedRebuildingMbps` in the `engine` status. 

Example of how the applied bandwidth limit appears in the volume engine status:

```yaml
  Rebuild Status:
    tcp://172.24.1.95:20001:
      Error:
      From Replica Address:  tcp://172.24.8.133:20001
      Is Rebuilding:         true
      Progress:              97
      State:                 in_progress
      appliedRebuildingMbps: 50
```


---

## Article: v2-data-engine/features/selective-v2-data-engine-activation.md

---
title: Selective V2 Data Engine Activation
weight: 20
aliases:
- /spdk/features/selective-v2-data-engine-activation.md
---

Starting with v1.6.0, Longhorn allows you to enable or disable the V2 Data Engine on specific cluster nodes. You can choose to enable the V2 Data Engine only on powerful nodes in a cluster with varied power states. This is not possible in v1.5.0, which enables the V2 Data Engine on all nodes.

## Disabling the V2 Data Engine on Specific Nodes

1. Identify the nodes that should not run the V2 Data Engine.

1. Add the label `node.longhorn.io/disable-v2-data-engine: "true"` to the selected nodes.

1. Enable the global setting `v2-data-engine`.

As a result, the following occur only on *nodes without the label*:
- Instance Manager pods for the V2 Data Engine are spawned.
- V2 Data Engine functionality remains available.

## Notice

V2 volume creation is possible only on nodes where the V2 Data Engine is enabled. You must schedule workloads that use V2 volumes on such nodes.

## Reference

For more information, see [[FEATURE] Selective V2 Data Engine Activation](https://github.com/longhorn/longhorn/issues/7015).


---

## Article: v2-data-engine/features/ublk-frontend-support.md

---
title: UBLK Frontend Support
weight: 20
aliases:
- /spdk/features/ublk-frontend-support.md
---

Starting with v1.9.0, Longhorn supports the UBLK frontend for v2 data engine volumes.
This feature exposes v2 data engine volumes as a block device by using [UBLK SPDK framework](https://spdk.io/doc/ublk.html).
In certain high-specification environments (for example, machines with fast SSDs capable of millions of IOPS and 32 CPU cores), the UBLK frontend may offer better performance compared to the default NVMe-oF frontend for v2 data engine volumes.
See the reference [Longhorn-Performance-Investigation](https://github.com/longhorn/longhorn/wiki/Longhorn-Performance-Investigation).
However, the UBLK frontend is less mature than the default NVMe-oF frontend (see the known limitations below).
UBLK frontend also has more restriction as mentioned below.

## Prerequisites

1. The kernel version on nodes must be v6.0 or above. The UBLK kernel driver is only available starting from kernel v6.0.
2. The kernel module `ublk_drv` must be loaded on each node where UBLK volumes are attached. You can load it manually on each volume for testing using: `modprobe ublk_drv`

## How to use

### When creating the v2 volume from UI
Select `UBLK` as the volume frontend during volume creation.

### When creating the v2 volume from manifest
1. Create a StorageClass specifies UBLK frontend. For example:
    ```yaml
    kind: StorageClass
    apiVersion: storage.k8s.io/v1
    metadata:
      name: my-ublk-frontend-storageclass
    provisioner: driver.longhorn.io
    allowVolumeExpansion: true
    reclaimPolicy: Delete
    volumeBindingMode: Immediate
    parameters:
      numberOfReplicas: "1"
      staleReplicaTimeout: "2880"
      fsType: "ext4"
      dataEngine: "v2"
      frontend: "ublk"
    ```
1. Create a PVC referencing that StorageClass. For example:
    ```yaml
    apiVersion: v1
    kind: PersistentVolumeClaim
    metadata:
      name: my-ublk-frontend-pvc
      namespace: default
    spec:
      accessModes:
        - ReadWriteOnce
      storageClassName: my-ublk-frontend-storageclass
      resources:
        requests:
          storage: 1Gi
    ```
2. Longhorn will automatically provision a v2 volume using the UBLK frontend based on the PVC and StorageClass definition

## Known Limitations

When an instance-manager pod crashes, it may leave orphaned UBLK devices on the node.
Currently, removing these orphan devices manually can be difficult and may sometimes require a node reboot.
We are investigating this issue further in [GitHub Issue #10738](https://github.com/longhorn/longhorn/issues/10738).


## Reference

Original GitHub issue for UBLK frontend support: [GitHub Issue #9456](https://github.com/longhorn/longhorn/issues/9456)


---

## Article: v2-data-engine/features/volume-clone.md

---
title: V2 Volume Clone Support
description: Creating a new volume as a duplicate of an existing volume
weight: 3
---


## Clone Using YAML

### Clone CSI snapshot
To clone a CSI snapshot, refer to the documentation on [Creating a Volume from a Snapshot](../../../snapshots-and-backups/csi-snapshot-support/csi-volume-snapshot-associated-with-longhorn-snapshot).

### Clone Volume with v2 data engine
Assume you have a `StorageClass` named `longhorn-v2`:
```yaml
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: longhorn-v2
provisioner: driver.longhorn.io
allowVolumeExpansion: true
reclaimPolicy: Delete
volumeBindingMode: Immediate
parameters:
  dataEngine: "v2"
  numberOfReplicas: "1"
  staleReplicaTimeout: "2880"
```
And you have a PersistentVolumeClaim (PVC) named `source-pvc-v2` provisioned from it:
```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: source-pvc-v2
spec:
  storageClassName: longhorn-v2
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
```

#### Clone using `full-copy` mode
You can create a new PVC with the exact same content as `source-pvc-v2` by applying the YAML below. Longhorn will copy the data from the source PVC to the new PVC.
```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: cloned-pvc-v2
spec:
  storageClassName: longhorn-v2
  dataSource:
    name: source-pvc-v2
    kind: PersistentVolumeClaim
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
```

#### Clone using `linked-clone` mode

The `full-copy` mode creates a new PVC that is fully independent of the source PVC. However, it requires time and resources to copy the data.
Sometimes, you need to quickly create a temporary PVC with the same content as the source, without copying the data. For example, backup solutions like **Velero** or **Kasten** can use this feature to quickly create a temporary PVC to read data and upload it to an S3 bucket.
In this scenario, you can use the `linked-clone` mode. This mode creates a new PVC that shares the same data blocks as the source PVC. Follow the steps below:

1. Create a StorageClass with `cloneMode` set to `linked-clone`.:
   ```yaml
   kind: StorageClass
   apiVersion: storage.k8s.io/v1
   metadata:
     name: longhorn-v2-linked-clone
   provisioner: driver.longhorn.io
   reclaimPolicy: Delete
   volumeBindingMode: Immediate
   parameters:
     dataEngine: "v2"
     cloneMode: "linked-clone"
     numberOfReplicas: "1"
     staleReplicaTimeout: "2880"
   ```
2. Create a new PVC that uses the above `StorageClass` and references the source PVC in the `dataSource` field:
   ```yaml
   apiVersion: v1
   kind: PersistentVolumeClaim
   metadata:
     name: cloned-pvc-v2-linked-clone
   spec:
     storageClassName: longhorn-v2-linked-clone
     dataSource:
       name: source-pvc-v2
       kind: PersistentVolumeClaim
     accessModes:
       - ReadWriteOnce
     resources:
       requests:
         storage: 10Gi
   ```

> **Note**:
> 1. In addition to the requirements for [CSI Volume Cloning](https://kubernetes.io/docs/concepts/storage/volume-pvc-datasource/), the cloned PVC's (`cloned-pvc`) `resources.requests.storage` must match the source PVC's (`source-pvc`) storage size.
> 2. The `linked-clone` mode is only supported by the v2 data engine.
> 3. A PVC created using `linked-clone` shares data blocks with the source and has the following limitations:
>    - It can have only one replica.
>    - It cannot be snapshotted or backed up.
>    - It cannot be used as the source for another clone operation.
>    - A source PVC can only have one `linked-clone` PVC at a time.
>    - `linked-clone` PVCs are designed to be short-lived.  It is highly recommended to delete them when no longer needed.

For more examples of linked-clone, see the blog post, [Backup Applications with Longhorn V2 Volumes using Velero](https://longhorn.io/blog/20250902-k8s-backup-solutions-and-longhorn/).
### Clone Volume Using the Longhorn UI

You can also clone a v2 data engine volume using the Longhorn UI by one of the following methods:

1. On the **Volumes** page, click **Create Volume** and select the data source (`Volume` or `Volume Snapshot`).
2. From the **Volumes** page, select a volume and click **Clone Volume** in the **Operation** menu.
3. On the **Volumes** page, select a volume, click its name, and in the **Snapshot and Backups** section of the details page, identify the snapshot to use and then click **Clone Volume**.
4. For bulk cloning, on the **Volumes** page, select one or more volumes and click the **Clone Volume** button at the top of the table.

## History

- [GitHub Issue](https://github.com/longhorn/longhorn/issues/7794)
- [Longhorn Enhancement Proposal](https://github.com/longhorn/longhorn/pull/10873)

Available since v1.10.0


---

## Article: v2-data-engine/features/volume-expansion.md

---
title: V2 Volume Expansion
weight: 20
aliases:
- /spdk/features/volume-expansion.md
---

Starting with v1.10.0, Longhorn supports online expansion for v2 data engine volumes that use the NVMe frontend. This feature allows users to expand a volume to the requested size while keeping the workload running.

During the expansion process, Longhorn automatically resizes all replicas to match the user-requested size. This eliminates the need to stop or detach the application from the volume, ensuring a seamless and non-disruptive scaling of storage.

This capability significantly improves flexibility in storage management by enabling volumes to be scaled without any downtime.

## How to use

### When creating the v2 volume from UI

1. Select a volume with `Block Device` or `NVMf` as the frontend.
2. Navigate to the Volumes page in the Longhorn UI.
3. Click **Expand Volume** from the volume operations menu.
4. Enter the new desired size and confirm. The expansion will begin automatically.

### When creating the v2 volume from manifest

1. Create a StorageClass for the v2 data engine. Make sure `allowVolumeExpansion` is set to `true`. For example:
    ```yaml
    kind: StorageClass
    apiVersion: storage.k8s.io/v1
    metadata:
        name: longhorn-v2-data-engine
    provisioner: driver.longhorn.io
    allowVolumeExpansion: true
    reclaimPolicy: Delete
    volumeBindingMode: Immediate
    parameters:
      numberOfReplicas: "3"
      staleReplicaTimeout: "2880"
      fsType: "ext4"
      dataEngine: "v2"
    ```

2. Create a PersistentVolumeClaim (PVC) that references this StorageClass:
    ```yaml
    apiVersion: v1
    kind: PersistentVolumeClaim
    metadata:
      name: longhorn-volv-pvc
      namespace: default
    spec:
      accessModes:
        - ReadWriteOnce
      storageClassName: longhorn-v2-data-engine
      resources:
        requests:
          storage: 2Gi
    ```

3. To expand the volume, edit the PVC manifest to increase the storage request to a larger size, then apply the updated manifest.
    ```yaml
      resources:
        requests:
          storage: 3Gi
    ```

## Known Limitations

The `UBLK` frontend is not supported for online expansion as of v1.10.0. Attempting to expand a volume using the UBLK frontend will not be allowed.

## Reference

For more information, see [[FEATURE] v2 supports volume expansion](https://github.com/longhorn/longhorn/issues/8022).

---

## Article: high-availability/_index.md

---
title: High Availability
weight: 5
---

---

## Article: high-availability/auto-balance-replicas.md

---
  title: Auto Balance Replicas
  weight: 1
---

When replicas are scheduled unevenly on nodes or zones, Longhorn `Replica Auto Balance` setting enables the replicas for automatic balancing when a new node is available to the cluster or when the replica count for a volume is updated.

## Replica Auto Balance Settings

### Global setting

Longhorn supports 3 options for global replica auto-balance setting:

- `disabled`. This is the default option, no replica auto-balance will be done.

- `least-effort`. This option instructs Longhorn to balance replicas for minimal redundancy.
  For example, after adding node-2, a volume with 4 off-balanced replicas will only rebalance 1 replica.
    ```
    node-1
    +-- replica-a
    +-- replica-b
    +-- replica-c
    node-2
    +-- replica-d
    ```

- `best-effort`. This option instructs Longhorn to try balancing replicas for even redundancy.
  For example, after adding node-2, a volume with 4 off-balanced replicas will rebalance 2 replicas.
    ```
    node-1
    +-- replica-a
    +-- replica-b
    node-2
    +-- replica-c
    +-- replica-d
    ```
  Longhorn does not forcefully re-schedule the replicas to a zone that does not have enough nodes
  to support even balance. Instead, Longhorn will re-schedule to balance at the node level.

### Volume specific setting

Longhorn also supports setting individual volume for `Replica Auto Balance`. The setting can be specified in `volume.spec.replicaAutoBalance`, this overrules the global setting.

There are 4 options available for individual volume setting:

- `Ignored`. This is the default option that instructs Longhorn to inherit from the global setting.

- `disabled`. This option instructs Longhorn no replica auto-balance should be done.

- `least-effort`. This option instructs Longhorn to balance replicas for minimal redundancy.
  For example, after adding node-2, a volume with 4 off-balanced replicas will only rebalance 1 replica.
    ```
    node-1
    +-- replica-a
    +-- replica-b
    +-- replica-c
    node-2
    +-- replica-d
    ```

- `best-effort`. This option instructs Longhorn to try balancing replicas for even redundancy.
  For example, after adding node-2, a volume with 4 off-balanced replicas will rebalance 2 replicas.
    ```
    node-1
    +-- replica-a
    +-- replica-b
    node-2
    +-- replica-c
    +-- replica-d
    ```
  Longhorn does not forcefully re-schedule the replicas to a zone that does not have enough nodes
  to support even balance. Instead, Longhorn will re-schedule to balance at the node level.


## How to Set Replica Auto Balance For Volumes

There are 3 ways to set `Replica Auto Balance` for Longhorn volumes:

### Change the global setting

You can change the global default setting for `Replica Auto Balance` inside Longhorn UI settings.
The global setting only functions as a default value, similar to the replica count.
It does not change any existing volume settings.
When a volume is created without specifying `Replica Auto Balance`, Longhorn will automatically set to `ignored` to inherit from the global setting.

### Set individual volumes to auto-balance replicas using the Longhorn UI

You can change the `Replica Auto Balance` setting for individual volume after creation on the volume detail page, or do multiple updates on the listed volume page.

### Set individual volumes to auto-balance replicas using a StorageClass

Longhorn also exposes the `Replica Auto Balance` setting as a parameter in a StorageClass.
You can create a StorageClass with a specified `Replica Auto Balance` setting, then create PVCs using this StorageClass.

For example, the below YAML file defines a StorageClass which tells the Longhorn CSI driver to set the `Replica Auto Balance` to `least-effort`:

```yaml
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: hyper-converged
provisioner: driver.longhorn.io
allowVolumeExpansion: true
parameters:
  numberOfReplicas: "3"
  replicaAutoBalance: "least-effort"
  staleReplicaTimeout: "2880" # 48 hours in minutes
  fromBackup: ""
```

## Replica Auto Balance Disk Pressure Threshold (%)

When `Replica Auto Balance` is enabled with `best-effort`, you can set a `Replica Auto Balance Disk Pressure Threshold (%)`. This threshold defines the disk usage level at which Longhorn will automatically attempt to migrate replicas to another disk on the same node.

For example, if the threshold is set to 75%, Longhorn will try to migrate replicas sequentially when the disk consumption reaches 75% capacity.

Longhorn prioritizes balancing replicas across node and zone first. Once the node and zones are balanced, it will then consider balancing within a single node based on disk pressure.

Since Longhorn v1.7.0, when rebuilding replicas on the same node, Longhorn uses local file data synchronization for more efficient data transfer.

## Limitations

Longhorn's automatic replica balancing feature is only activated for volumes that have a robustness status of `Healthy`.

**Unhealthy volumes** or **detached volumes** will not be automatically rebalanced, even if a node has low available space.

This behavior is a deliberate design choice to ensure system stability and data integrity:
- Moving replicas or triggering automatic rebuilds on an unhealthy volume could further compromise data integrity. This design ensures replica stability and requires manual intervention so an administrator can assess the volume's condition before initiating potentially risky operations.
- Detached volumes are not actively serving I/O. Skipping automatic rebalance for these volumes prevents unnecessary rebuilds and saves cluster resources.

If a volume is unhealthy or detached, moving replicas requires manual intervention, such as:
- Rebuilding the replicas after the volume has been inspected and/or attached.
- Attaching the volume if it is detached (to restore a healthy state, if possible).


---

## Article: high-availability/data-locality.md

---
  title: Data Locality
  weight: 1
---

The data locality setting is intended to be enabled in situations where at least one replica of a Longhorn volume should be scheduled on the same node as the pod that uses the volume, whenever it is possible. We refer to the property of having a local replica as having `data locality`.

For example, data locality can be useful when the cluster's network is bad, because having a local replica increases the availability of the volume.

Data locality can also be useful for distributed applications (e.g. databases), in which high availability is achieved at the application level instead of the volume level. In that case, only one volume is needed for each pod, so each volume should be scheduled on the same node as the pod that uses it.  In addition, the default Longhorn behavior for volume scheduling could cause a problem for distributed applications. The problem is that if there are two replicas of a pod, and each pod replica has one volume each, Longhorn is not aware that those volumes have the same data and should not be scheduled on the same node. Therefore Longhorn could schedule identical replicas on the same node, therefore preventing them from providing high availability for the workload.

When data locality is disabled, a Longhorn volume can be backed by replicas on any nodes in the cluster and accessed by a pod running on any node in the cluster.

## Data Locality Settings

Longhorn currently supports two modes for data locality settings:

- `disabled`: This is the default option. There may or may not be a replica on the same node as the attached volume (workload).

- `best-effort`: This option instructs Longhorn to try to keep a replica on the same node as the attached volume (workload). Longhorn will not stop the volume, even if it cannot keep a replica local to the attached volume (workload) due to an environment limitation, e.g. not enough disk space, incompatible disk tags, etc.

- `strict-local`: This option enforces Longhorn keep the **only one replica** on the same node as the attached volume, and therefore, it offers higher IOPS and lower latency performance. This option is incompatible with [ReadWriteMany (RWX) volume](../../nodes-and-volumes/volumes/rwx-volumes).


## How to Set Data Locality For Volumes

There are three ways to set data locality for Longhorn volumes:

### Change the default global setting

You can change the global default setting for data locality inside Longhorn UI settings.
The global setting only functions as a default value, similar to the replica count.
It doesn't change any existing volume's settings.
When a volume is created without specifying data locality, Longhorn will use the global default setting to determine data locality for the volume.

### Change data locality for an individual volume using the Longhorn UI

You can use Longhorn UI to set data locality for volume upon creation.
You can also change the data locality setting for the volume after creation in the volume detail page.

### Set the data locality for individual volumes using a StorageClass
Longhorn also exposes the data locality setting as a parameter in a StorageClass.
You can create a StorageClass with a specified data locality setting, then create PVCs using the StorageClass.
For example, the below YAML file defines a StorageClass which tells the Longhorn CSI driver to set the data locality to `best-effort`:

```yaml
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: hyper-converged
provisioner: driver.longhorn.io
allowVolumeExpansion: true
parameters:
  numberOfReplicas: "2"
  dataLocality: "best-effort"
  staleReplicaTimeout: "2880" # 48 hours in minutes
  fromBackup: ""
```



---

## Article: high-availability/k8s-cluster-autoscaler.md

---
  title: Kubernetes Cluster Autoscaler Support (Experimental)
  weight: 1
---

By default, Longhorn blocks Kubernetes Cluster Autoscaler from scaling down nodes because:
- Longhorn creates PodDisruptionBudgets for all engine and replica instance-manager pods.
- Longhorn instance manager pods have strict PodDisruptionBudgets.
- Longhorn instance manager pods are not backed by a Kubernetes built-in workload controller .
- Longhorn pods are using local storage volume mounts.

For more information, see [What types of pods can prevent CA from removing a node?](https://github.com/kubernetes/autoscaler/blob/master/cluster-autoscaler/FAQ.md#what-types-of-pods-can-prevent-ca-from-removing-a-node)

If you want to unblock the Kubernetes Cluster Autoscaler scaling, you can set the setting [Kubernetes Cluster Autoscaler Enabled](../../references/settings#kubernetes-cluster-autoscaler-enabled-experimental).

When this setting is enabled, Longhorn will retain the least instance-manager PodDisruptionBudget as possible. Each volume will have at least one replica under the protection of an instance-manager PodDisruptionBudget while no redundant PodDisruptionBudget blocking the Cluster Autoscaler from from scaling down.

When this setting is enabled, Longhorn will also add `cluster-autoscaler.kubernetes.io/safe-to-evict` annotation to Longhorn workloads that are not backed by a Kubernetes built-in workload controller or are using local storage mounts.

> **Warning:** Replica rebuilding could be expensive because nodes with reusable replicas could get removed by the Kubernetes Cluster Autoscaler.


---

## Article: high-availability/node-failure.md

---
title: Node Failure Handling with Longhorn
weight: 2
---

## What to expect when a Kubernetes Node fails

This section is aimed to inform users of what happens during a node failure and what is expected during the recovery.

After **one minute**, `kubectl get nodes` will report `NotReady` for the failure node.

After about **five minutes**, the states of all the pods on the `NotReady` node will change to either `Unknown` or `NodeLost`.

StatefulSets have a stable identity, so Kubernetes won't force delete the pod for the user. See the [official Kubernetes documentation about forcing the deletion of a StatefulSet](https://kubernetes.io/docs/tasks/run-application/force-delete-stateful-set-pod/).

Deployments don't have a stable identity, but for the Read-Write-Once type of storage, since it cannot be attached to two nodes at the same time, the new pod created by Kubernetes won't be able to start due to the RWO volume still attached to the old pod, on the lost node.

In both cases, Kubernetes will automatically evict the pod (set deletion timestamp for the pod) on the lost node, then try to **recreate a new one with old volumes**. Because the evicted pod gets stuck in `Terminating` state and the attached volumes cannot be released/reused, the new pod will get stuck in `ContainerCreating` state, if there is no intervene from admin or storage software.

## Longhorn Pod Deletion Policy When Node is Down

Longhorn provides an option to help users automatically force delete terminating pods of StatefulSet/Deployment on the node that is down. After force deleting, Kubernetes will detach the Longhorn volume and spin up replacement pods on a new node.

You can find more detail about the setting options in the `Pod Deletion Policy When Node is Down` in the **Settings** tab in the Longhorn UI or [Settings reference](../../references/settings/#pod-deletion-policy-when-node-is-down)

## What to expect when a failed Kubernetes Node recovers

If the node is back online within 5 - 6 minutes of the failure, Kubernetes will restart pods, unmount, and re-mount volumes without volume re-attaching and VolumeAttachment cleanup.

Because the volume engines would be down after the node is down, this direct remount won’t work since the device no longer exists on the node.

In this case, Longhorn will detach and re-attach the volumes to recover the volume engines, so that the pods can remount/reuse the volumes safely.

If the node is not back online within 5 - 6 minutes of the failure, Kubernetes will try to delete all unreachable pods based on the pod eviction mechanism and these pods will be in a `Terminating` state. See [pod eviction timeout](https://kubernetes.io/docs/concepts/architecture/nodes/#condition) for details.

Then if the failed node is recovered later, Kubernetes will restart those terminating pods, detach the volumes, wait for the old VolumeAttachment cleanup, and reuse(re-attach & re-mount) the volumes. Typically these steps may take 1 ~ 7 minutes.

In this case, detaching and re-attaching operations are already included in the Kubernetes recovery procedures. Hence no extra operation is needed and the Longhorn volumes will be available after the above steps.

For all above recovery scenarios, Longhorn will handle those steps automatically with the association of Kubernetes.


---

## Article: high-availability/recover-volume.md

---
  title: Volume Recovery
  weight: 1
---

Longhorn provides two mechanisms for maintaining volume functionality in a variety of situations.

## Automatic Workload Pod Deletion

This recovery mechanism is enabled by the setting [*Automatically Delete Workload Pod when The Volume Is Detached Unexpectedly*](../../references/settings#automatically-delete-workload-pod-when-the-volume-is-detached-unexpectedly).

When one of the following situations occurs, Longhorn automatically attempts to delete workload pods that are managed by a controller (for example, Deployment, StatefulSet, or DaemonSet). After deletion, the controller restarts the workload pod and Kubernetes handles volume reattachment and remounting.

1. A volume was unexpectedly detached, possibly because of a [Kubernetes upgrade](https://github.com/longhorn/longhorn/issues/703), [container runtime reboot](https://github.com/longhorn/longhorn/issues/686), network connectivity issue, or volume engine crash.
2. A volume was automatically salvaged after all replicas became faulty, possibly because of a network connectivity issue. Longhorn attempts to identify the usable replicas and uses them for the volume.
3. An error occurred on a Share Manager pod that uses an RWX volume.

If you want to prevent Longhorn from automatically deleting workload pods, disable the setting [*Automatically Delete Workload Pod when The Volume Is Detached Unexpectedly*](../../references/settings#automatically-delete-workload-pod-when-the-volume-is-detached-unexpectedly) on the Longhorn UI.

Longhorn does not delete pods without a controller because such pods cannot be restarted after deletion. To recover volumes that are unexpectedly detached, you must manually delete and restart the pods without a controller.

## Automatic Volume Remounting

This recovery mechanism is not controlled by any specific setting.

The state of a volume can change to read-only when IO errors occur. IO errors can be caused by a variety of issues, including the following:
- Network disconnection: Interrupted connection between the engine and replicas.
- High disk latency: Significant delay in the transfer of data between a replica and the corresponding disk.

Longhorn checks the state of the volume's global mount point every 10 seconds. When the volume's filesystem changes to read-only, Longhorn updates the condition to the volume's data engine. Longhorn then automatically attempts to remount the global mount point on the host to change the state back to read-write. Upon successful remounting, the workload pods continue functioning without disruption. However, if the mount point becomes write-protected and Longhorn fails to remount the mount point, you may still need to manually recreate the workload to force it reattach and remount the volume.

> **Note:**
> This mechanism might not work in some situations. For example, when the volume's data engine crashes, Longhorn automatically detaches and reattaches the volume. The filesystem changes to read-only in this case. Longhorn will detect the read-only mode and update the state, but [Automatic Volume Remounting](#automatic-volume-remounting) cannot change it back to read-write because the device is now write-protected. In this case, you can only rely on the [Automatic Workload Pod Deletion](#automatic-workload-pod-deletion) mechanism, which enables volume remounting after the workload pod is recreated.


## Summary

[Automatic Workload Pod Deletion](#automatic-workload-pod-deletion) is triggered when unexpected failures happen. The controller deletes and then restarts the workload pod, and Kubernetes handles volume reattachment and remounting. The process may cause interruptions to the workload. If you want to prevent Longhorn from automatically deleting workload pods, disable the setting [*Automatically Delete Workload Pod when The Volume Is Detached Unexpectedly*](../../references/settings#automatically-delete-workload-pod-when-the-volume-is-detached-unexpectedly) on the Longhorn UI.

[Automatic Volume Remounting](#automatic-volume-remounting) is triggered when the volume's filesystem changes to read-only. Longhorn remounts the global mount point on the host to change the state back to read-write.

---

## Article: high-availability/rwx-volume-fast-failover.md

---
  title: RWX Volume Fast Failover (Experimental)
  weight: 1
---

Release 1.7.0 adds a feature that minimizes the downtime for ReadWriteMany volumes when a node fails.  When enabled Longhorn uses a lease-based mechanism to monitor the state of the NFS server pod that exports the volume Longhorn reacts quickly to move it to another node if it becomes unresponsive.  See [RWX Volumes](../../nodes-and-volumes/volumes/rwx-volumes) for details on how the NFS server works.

To enable the feature, you set [RWX Volume Fast Failover](../../references/settings#rwx-volume-fast-failover-experimental) to "true".  Existing RWX volumes will need to be restarted to use the feature after the setting is changed.  That is done by scaling the workload down to zero and then back up again.  New volumes will pick up the setting at creation and be configured appropriately.  

With the feature enabled, when a pod is created or re-created, Longhorn also creates an associated lease object in the `longhorn-system` namespace, with the same name as the volume.  The NFS server pod keeps the lease renewed as proof of life.  If the renewal stops happening, Longhorn will take steps to create a new NFS server pod on another node and to re-attach the workload, even before the old node is marked as `Not Ready` by Kubernetes.

Along with adding the monitoring and fast reaction, the feature also changes the NFS server configuration to use a shortened grace period for client re-connection.

If the setting is changed back to "false", the lease check is disabled and pod relocation will use regular Kubernetes rules for node failure, even on existing volumes.  When the server pod is next restarted, it will revert to the normal grace period configuration.

For more information, see https://github.com/longhorn/longhorn/issues/6205.

> **Note:**  In rare circumstances, it is possible for the failover to become deadlocked. This happens if the NFS server pod creation is blocked by a recovery action that is itself blocked by the failover-in-process state.  If that is the case, and failover takes more than a minute or two, the workaround is to delete the associated lease object.  That clears the state, and a new lease is created along with the replacement server pod.  For example, if the stuck volume is named `pvc-2ce4e82e-7ccc-46c0-90a8-a141501fbf93` and the feature is enabled, there will be a lease with the same name.  To delete the associated lease object:
> ```bash
> kubectl -n longhorn-system delete lease pvc-2ce4e82e-7ccc-46c0-90a8-a141501fbf93
> ```
> See, for example, https://github.com/longhorn/longhorn/issues/9093.

### Resource Consumption and System Performance Impact

The Longhorn team has investigated the impact of RWX volumes on resource consumption and system performance. The benchmarking studies, which were completed using 60 RWX volumes, show that enabling the *RWX Volume Fast Failover* feature results in the following: 

- More requests sent to the Kubernetes API server (kube-apiserver)
- More remote procedure calls (RPCs) sent from kube-apiserver to etcd
- Slight increase in CPU and memory usage

#### **Environment:**

- **Setup:** 1 Control Node + 3 Worker Nodes (v1.27.15+rke2r1)
- **Workload:** 60 Deployments with 60 RWX volumes with `soft` mount

#### **Test Results:**

| **Metric**                           | **Fast Failover Disabled** | **Fast Failover Enabled** | **Difference**             |
|--------------------------------------|---------------------------|----------------------------|----------------------------|
| **API request rate (kube-apiserver)**        | 37.5 req/s                  | 59 req/s                 | +57.3%                 |
| **RPC rate (kube-apiserver to etcd)**                      | 37 ops/s                  | 57 ops/s                   | +54.1%                 |
| **Memory usage**                  | Lower Peaks/Minima       | Higher Peaks/Minima         | Increased usage with Fast Failover enabled |
| **Longhorn Manager CPU/RAM usage**      | 405 MB / 0.1 CPU          | 417 MB / 0.13 CPU            | +3% RAM / +30% CPU |
| **Share Manager CPU/RAM usage**         | 2.2 GB / 0.235 CPU         | 2.25 GB / 0.26 CPU          | +2.3% RAM / +10.6% CPU |

For detailed screenshots and further context, please refer to the [related issue discussion](https://github.com/longhorn/longhorn/issues/6205#issuecomment-2262625965).
