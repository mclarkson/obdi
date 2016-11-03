# Install Obdi on RedHat/CentOS 6

## Build from Source
 
```
cd
rpm -ivh https://dl.fedoraproject.org/pub/epel/epel-release-latest-6.noarch.rpm
yum install rpmdevtools gcc golang tar
git clone https://github.com/mclarkson/obdi.git
cd obdi
BUILD_NUMBER=1 ./dist/jenkins-build.sh
```

