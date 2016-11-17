# How to install Obdi from source

This is currently the recommended way to install Obdi.

## Install Obdi on CentOS 7

Follow these instructions to get Obdi installed on Centos 7. Broadly,
installation involves:

1. Install OS.
2. Download Obdi source using Git.
3. Run a script to create an RPM.
4. Install that RPM.

The RPM could be stored in a yum repository on the network to make
updating the Obdi master and worker(s) easy.

### Install OS

This can be done in many ways. The video gives an example of how to accomplish
this task using the Amazon AWS management console.

<video src="/videos/centos7install_installos.webm" style="width: 100%" controls preload></video>

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
extra dependencies can be installed:

```
rpm -ivh https://dl.fedoraproject.org/pub/epel/epel-release-latest-7.noarch.rpm
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

<video src="/videos/centos7install_downloadsource.webm" style="width: 100%" controls preload></video>

### Run a script to create an RPM

Build the RPM:

```
cd obdi
./dist/jenkins-build-redhat7.sh
```

The following video shows the build in action.

<video src="/videos/centos7install_runscript.webm" style="width: 100%" controls preload></video>

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

<video src="/videos/centos7install_installrpm.webm" style="width: 100%" controls preload></video>

To start obdi:

```
systemctl start obdi
systemctl start obdi-worker
```

If you installed Centos 7 in Amazon, as shown at the top of this page, then you will probably need to open the HTTPS port 443 in the AWS EC2 security group using the Amazon AWS management console, and also in the Centos 7 virtual machine like so:

```
# -- IF THIS IS AN AMAZON INSTANCE USING OFFICIAL CENTOS 7 --

# -- NOTE THAT I DID NOT NEED TO DO THIS! MAYBE YOU WILL?    --
# -- CHECK 'iptables -L' AND IF EMPTY DON'T DO THE FOLLOWING --

# Get the line number of the last item in the INPUT chain, minus 1
num=$(iptables -L INPUT --line-numbers | tail -1 | awk '{print $1-1}')

# Then insert before the last item
iptables -I INPUT $num -p tcp -m tcp --dport 443 -j ACCEPT

# And make it survive reboots
service iptables save
```

Now you can test the Obdi user interfaces by opening your Web Browser and connecting to 'https://MACHINE_ADDRESS/manager/run' and 'https//MACHINE_ADDRESS/manager/admin'.

That's it!

