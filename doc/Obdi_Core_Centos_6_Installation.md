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
yum install rpmdevtools gcc golang tar
```

Download obdi using Git:

```
git clone https://github.com/mclarkson/obdi.git
```

The following video shows the previous commands being run.

\[VIDEO]

### Run a script to create an RPM

Build the RPM:

```
cd obdi
BUILD_NUMBER=1 ./dist/jenkins-build.sh
```

The following video shows the build in action.

\[VIDEO]

### Install that RPM

Find the RPM files:

```
cd ~/obdi
TheRpm=$(find . -name "*.rpm")
echo $TheRpm
```

Finally, so long as the last echo line, above, managed to
output the location then install the RPM:

```
rpm -ivh $TheRpm
```

And here's the video for this step:

\[VIDEO]

That's it!

