# incus_instance

Manages an Incus instance that can be either a container or virtual machine.

An instance can take a number of configuration and device options. A full reference can be found [here](https://linuxcontainers.org/incus/docs/main/reference/instance_options/).

## Basic Example

```hcl
resource "incus_instance" "instance1" {
  name  = "instance1"
  image = "images:ubuntu/22.04"

  config = {
    "boot.autostart" = true
    "limits.cpu"     = 2
  }
}
```

## Example to Attach a Volume

```hcl
resource "incus_storage_pool" "pool1" {
  name   = "mypool"
  driver = "zfs"
}

resource "incus_storage_volume" "volume1" {
  name = "myvolume"
  pool = incus_storage_pool.pool1.name
}

resource "incus_instance" "instance1" {
  name  = "instance1"
  image = "ubuntu"

  device {
    name = "volume1"
    type = "disk"

    properties = {
      path   = "/mount/point/in/instance"
      source = incus_storage_volume.volume1.name
      pool   = incus_storage_pool.pool1.name
    }
  }
}
```

## Example to proxy/forward ports

```hcl
resource "incus_instance" "instance2" {
  name      = "instance2"
  image     = "ubuntu"
  profiles  = ["default"]
  ephemeral = false

  device {
    name = "http"
    type = "proxy"

    properties = {
      # Listen on Incus host's TCP port 80
      listen = "tcp:0.0.0.0:80"
      # And connect to the instance's TCP port 80
      connect = "tcp:127.0.0.1:80"
    }
  }
}
```

## Example to create a new instance from an existing instance

```hcl
resource "incus_instance" "instance1" {
  project = "default"
  name    = "instance1"
  image   = "images:debian/12"
}

resource "incus_instance" "instance2" {
  project = "default"
  name    = "instance2"

  source_instance = {
    project = "default"
    name    = "instance1"
  }
}
```

## Example to create a new instance from an instance backup

```hcl
resource "incus_instance" "instance1" {
  project     = "default"
  name        = "instance1"
  source_file = "/path/to/backup.tar.gz"
}
```

## Example to create a new instance from an instance backup with storage

In order to provide the storage pool name for an instance, which is created
from a backup exactly one `device` configuration of `type = "disk"` might be
provided. The name of the pool is given as the `pool` attribute in
`properties`. Additionally the property `path = "/"` is required.

```hcl
resource "incus_instance" "instance1" {
  project     = "default"
  name        = "instance1"
  source_file = "/path/to/backup.tar.gz"

  device = {
    name = "storage"
    type = "disk"

    properties = {
      path = "/"
      pool = "pool-name"
    }
  }
}
```

## Example of waiting for the Incus agent in a virtual machine

```hcl
resource "incus_instance" "instance1" {
  project = "default"
  name    = "instance1"
  image   = "images:debian/12"
  type    = "virtual-machine"

  wait_for {
    type = "agent"
  }
}
```

## Example of waiting for a certain time period

```hcl
resource "incus_instance" "instance1" {
  project = "default"
  name    = "instance1"
  image   = "images:debian/12"

  wait_for {
    type  = "delay"
    delay = "30s"
  }
}
```

## Example of waiting for the IPv4 network to be ready

```hcl
resource "incus_instance" "instance1" {
  project = "default"
  name    = "instance1"
  image   = "images:debian/12"

  wait_for {
    type = "ipv4"
  }
}
```

## Example of waiting for the IPv6 network to be ready on a specific network interface

```hcl
resource "incus_instance" "instance1" {
  project = "default"
  name    = "instance1"
  image   = "images:debian/12"

  wait_for {
    type = "ipv6"
    nic  = "eth0"
  }
}
```

## Example of waiting for the IPv4 and IPv6 network to be ready on a specific network interface

```hcl
resource "incus_instance" "instance1" {
  project = "default"
  name    = "instance1"
  image   = "images:debian/12"

  wait_for {
    type = "ipv4"
    nic  = "eth0"
  }

  wait_for {
    type = "ipv6"
    nic  = "eth0"
  }
}
```

## Argument Reference

* `name` - **Required** - Name of the instance.

* `image` - *Optional* - Base image from which the instance will be created. Must
  specify [an image accessible from the provider remote](https://linuxcontainers.org/incus/docs/main/reference/image_servers/).

* `source_file` - *Optional* - The souce backup file from which the instance should be restored. For handling of storage pool, see examples.

* `source_instance` - *Optional* - The source instance from which the instance will be created. See reference below.

* `description` - *Optional* - Description of the instance.

* `type` - *Optional* - Instance type. Can be `container`, or `virtual-machine`. Defaults to `container`.

* `ephemeral` - *Optional* - Boolean indicating if this instance is ephemeral. Defaults to `false`.

* `running` - *Optional* - Boolean indicating whether the instance should be started (running). Defaults to `true`.

* `wait_for` - *Optional* - WaitFor definition. See reference below.
  If `running` is set to false or instance is already running (on update), this value has no effect.

* `profiles` - *Optional* - List of Incus config profiles to apply to the new
  instance. Profile `default` will be applied if profiles are not set (are `null`).
  However, if an empty array (`[]`) is set as a value, no profiles will be applied.

* `device` - *Optional* - Device definition. See reference below.

* `file` - *Optional* - File to upload to the instance. See reference below.

* `config` - *Optional* - Map of key/value pairs of
  [instance config settings](https://linuxcontainers.org/incus/docs/main/reference/instance_options/).

* `project` - *Optional* - Name of the project where the instance will be spawned.

* `remote` - *Optional* - The remote in which the resource will be created. If
  not provided, the provider's default remote will be used.

* `target` - *Optional* - Specify a target node in a cluster.

* `architecture` - *Optional* - The instance architecture (e.g. x86_64, aarch64). See [Architectures](https://linuxcontainers.org/incus/docs/main/architectures/) for all possible values.

The `source_instance` block supports:

* `project` - **Required** - Name of the project in which the source instance exists.

* `name` - **Required** - Name of the source instance.

* `snapshot`- *Optional* - Name of the snapshot of the source instance

The `wait_for` block supports:

* `type` - **Required** - Type for what should be waited for. Can be `agent`, `delay`, `ipv4`, `ipv6` or `ready`.

* `delay` - *Optional* - Delay time that should be waited for when type is `delay`, e.g. `30s`.

* `nic` - *Optional* - Network interface that should be waited for when type is `ipv4` or `ipv6`.

The `device` block supports:

* `name` - **Required** - Name of the device.

* `type` - **Required** - Type of the device Must be one of none, disk, nic,
  unix-char, unix-block, usb, gpu, infiniband, proxy, unix-hotplug, tpm, pci.

* `properties`- **Required** - Map of key/value pairs of
  [device properties](https://linuxcontainers.org/incus/docs/main/reference/devices/).

The `file` block supports:

* `content` - **Required** unless source_path is used* - The *contents* of the file.
  Use the `file()` function to read in the content of a file from disk.

* `source_path` - **Required** unless content is used* - The source path to a file to
  copy to the instance.

* `target_path` - **Required** - The absolute path of the file on the instance,
  including the filename.

* `uid` - *Optional* - The UID of the file. Must be an unquoted integer.

* `gid` - *Optional* - The GID of the file. Must be an unquoted integer.

* `mode` - *Optional* - The octal permissions of the file, must be quoted. Defaults to `0755`.

* `create_directories` - *Optional* - Whether to create the directories leading
  to the target if they do not exist.

## Attribute Reference

The following attributes are exported:

* `ipv4_address` - The IPv4 Address of the instance. See Instance Network
  Access for more details.

* `ipv6_address` - The IPv6 Address of the instance. See Instance Network
  Access for more details.

* `mac_address` - The MAC address of the detected NIC. See Instance Network
  Access for more details.

* `status` - The status of the instance.

## Instance Network Access

If your instance has multiple network interfaces, you can specify which one
Terraform should report the IP addresses of. If you do not specify an interface,
Terraform will use the first address from the best interface detected.

To specify an interface, do the following:

```hcl
resource "incus_instance" "instance1" {
  name     = "instance1"
  image    = "images:alpine/edge/amd64"
  profiles = ["default"]

  config = {
    "user.access_interface" = "eth0"
  }
}
```

## Importing

Import ID syntax: `[<remote>:][<project>/]<name>[,image=<image>]`

* `<remote>` - *Optional* - Remote name.
* `<project>` - *Optional* - Project name.
* `<name>` - **Required** - Instance name.
* `image=<image>` - *Optional* - The image used by the instance.

~> **Warning:** Importing the instance without specifying `image` will lead to its replacement
   upon the next apply, rather than an in-place update.

### Import example

Example using terraform import command:

```shell
terraform import incus_instance.myinst proj/c1,image=images:alpine/edge/amd64
```

Example using the import block (only available in Terraform v1.5.0 and later):

```hcl
resource "incus_instance" "myinst" {
  name    = "c1"
  project = "proj"
  image   = "images:alpine/edge/amd64"
}

import {
  to = incus_instance.myinst
  id = "proj/c1,image=images:alpine/edge/amd64"
}
```

## Notes

* The instance resource `config` includes some keys that can be automatically generated by the Incus.
  If these keys are not explicitly defined by the user, they will be omitted from the Terraform
  state and treated as computed values.
  * `image.*`
  * `volatile.*`
