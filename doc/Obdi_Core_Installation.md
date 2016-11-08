# How to install Obdi from source

This is currently the recommended way to install Obdi.

## Install Obdi on CentOS 6

Follow these instructions to get Obdi installed on Centos 6. Broadly,
installation involves:

1. Install OS.
2. Download Obdi source using Git.
3. Run a script to create an RPM.
4. Install that RPM.

The RPM could be stored in a yum repository on the network to make
updating the Obdi master and workers easily.

### Build from Source

The following commands should be typed (or copy/pasted) into the linux terminal.

Go to your home directory:

```
cd
```

Install the [EPEL](https://fedoraproject.org/wiki/EPEL) repository for installing extra dependencies:

```
rpm -ivh https://dl.fedoraproject.org/pub/epel/epel-release-latest-6.noarch.rpm
```

Install some dependencies:

```
yum install rpmdevtools gcc golang tar
```

Download obdi using Git:

```
git clone https://github.com/mclarkson/obdi.git
```

Build the RPM:

```
cd obdi
BUILD_NUMBER=1 ./dist/jenkins-build.sh
```

