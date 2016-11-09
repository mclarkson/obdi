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

### Install OS

This can be done in many ways. The video gives an example of how to accomplish
this task using the Amazon AWS management console.

<video src="/videos/centos6install_installos.webm" style="width: 100%" controls preload></video>

### Download Obdi source using Git

Now log into the new Linux box and type, or copy/paste, the next commands
into the linux terminal.

Switch user to the super user:

```
sudo su
```

Go to your home directory:

```
cd
```

Install the [EPEL](https://fedoraproject.org/wiki/EPEL) repository so that
extra extra dependencies can be installed:

```
rpm -ivh https://dl.fedoraproject.org/pub/epel/epel-release-latest-6.noarch.rpm
```

Install some dependencies:

```
yum install rpmdevtools gcc golang tar git
```

Download obdi using Git:

```
git clone https://github.com/mclarkson/obdi.git
```

The following video shows the previous commands being run.

<video src="/videos/centos6install_downloadsource.webm" style="width: 100%" controls preload></video>

### Run a script to create an RPM

Build the RPM:

```
cd obdi
BUILD_NUMBER=1 ./dist/jenkins-build.sh
```

The following video shows the build in action.

<video src="/videos/centos6install_runscript.webm" style="width: 100%" controls preload></video>

### Install that RPM

Find the RPM files:

```
cd ~/obdi
TheRpm=$(find rpmbuild/RPMS -name "*.rpm")
echo $TheRpm
```

Finally, so long as the last echo line, above, managed to
output the location then install the RPM:

```
rpm -ivh $TheRpm
```

And here's the video for this step:

<video src="/videos/centos6install_installrpm.webm" style="width: 100%" controls preload></video>

To start obdi:

```
service obdi start
service obdi-worker start
```

If you installed Centos 6 in Amazon, as shown at the top of this page, then you will probably need to open the HTTPS port 443 in the AWS EC2 security group using the Amazon AWS management console, and also in the Centos 6 virtual machine like so:

```
iptables -A INPUT -p tcp -m tcp --dport 443 -j ACCEPT
/etc/init.d/iptables save
```

Now you can test the Obdi user interfaces by connecting to 'https://MACHINE_ADDRESS/manager/run' and 'https//MACHINE_ADDRESS/manager/admin'.

That's it!

