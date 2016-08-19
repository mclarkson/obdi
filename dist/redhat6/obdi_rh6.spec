%define name obdi
%define version 0.2.2
%define release 1
# Don't strip the Golang binaries
%define __strip /bin/true
# The following line may be required
#%define debug_package %{nil}

Summary: A REST interface and GUI for deploying software
# Don't replace the next line with %{name}, hudson build script needs it
Name: obdi
Version: %{version}
Release: %{release}
License: Apache-2.0
Group: Application/System
Source: obdi-%{version}.tar.gz
Requires: gcc golang
BuildRequires: golang
# PreReq: sh-utils
BuildRoot: %{_builddir}/%{name}-%{version}/tmp
Packager: Mark Clarkson
Vendor: Smorg

%description
A REST interface and GUI for deploying software.

%package worker
Group: Application/System
Summary: Obdi Worker daemon.
#Requires: obdi

%description worker
Obdi worker daemon

%prep
%setup -q

# Pre Install
%pre

# Post Install
%post
if [ "$1" = 1 ]; then
    # New install
    openssl req -new -newkey rsa:4096 -days 3650 -nodes -x509 -sha256 \
        -subj "/C=US/ST=Denial/L=Springfield/O=Dis/CN=snakeoil" \
        -keyout /etc/obdi/certs/key.pem \
        -out /etc/obdi/certs/cert.pem
fi

%post worker
if [ "$1" = 1 ]; then
    # New install
    openssl req -new -newkey rsa:4096 -days 3650 -nodes -x509 -sha256 \
        -subj "/C=US/ST=Denial/L=Springfield/O=Dis/CN=snakeoil" \
        -keyout /etc/obdi-worker/certs/key.pem \
        -out /etc/obdi-worker/certs/cert.pem
fi

# Pre Uninstall
%preun

# Post Uninstall
%postun

%preun worker
if [ "$1" = 0 ]; then
    # uninstall
    :
fi

%build

# Get GOROOT from go
eval `go env | grep GOROOT`

# Include Golang binaries
export PATH=$PATH:/usr/sbin:$GOROOT/bin

# Golang requirements
export GOPATH=$PWD

# Fix include paths
mkdir -p src/github.com/mclarkson/obdi
ln -s ../../../../external src/github.com/mclarkson/obdi/external 

# Build
cd obdi
go build -ldflags "-X main.VERSION %{version}" -o obdi
cd ..
cd obdi-worker
go build -o obdi-worker

%install

[ "$RPM_BUILD_ROOT" != "/" ] && %{__rm} -rf %{buildroot}

# Config
install -d -m 755 ${RPM_BUILD_ROOT}/%{_sysconfdir}/obdi/certs/
install -d -m 755 ${RPM_BUILD_ROOT}/%{_sysconfdir}/obdi-worker/certs/
install -D -m 640 conf/obdi.conf ${RPM_BUILD_ROOT}/%_sysconfdir/obdi/obdi.conf
install -D -m 640 conf/obdi-worker.conf ${RPM_BUILD_ROOT}/%_sysconfdir/obdi-worker/obdi-worker.conf

# Golang single binary (/usr/sbin)
install -D -m 755 obdi/obdi ${RPM_BUILD_ROOT}%{_sbindir}/obdi
install -D -m 755 obdi-worker/obdi-worker ${RPM_BUILD_ROOT}/%{_sbindir}/obdi-worker

# Initrd
install -D -m 755 init/obdi ${RPM_BUILD_ROOT}/%{_initrddir}/obdi
install -D -m 755 init/obdi-worker ${RPM_BUILD_ROOT}/%{_initrddir}/obdi-worker

# Directories

# Cache directory (/var/cache)
install -d -m 755 ${RPM_BUILD_ROOT}/%{_var}/cache/obdi/

# Main database directory (/var/lib)
install -d -m 755 ${RPM_BUILD_ROOT}/%{_sharedstatedir}/obdi/

# Database directory for plugins (/var/lib)
install -d -m 755 ${RPM_BUILD_ROOT}/%{_sharedstatedir}/obdi/plugins/

# Compiled Golang plugin directory (/usr/lib)
install -d -m 755 ${RPM_BUILD_ROOT}/%{_usr}/lib/obdi/plugins/

# Golang plugin source code directory (/usr/lib)
install -d -m 755 ${RPM_BUILD_ROOT}/%{_usr}/lib/obdi/plugins/src/

# Static Content Directories - Web server files (/usr/share)
install -d -m 755 ${RPM_BUILD_ROOT}/%{_datarootdir}/obdi/
cp -r static ${RPM_BUILD_ROOT}/%{_datarootdir}/obdi/

%files
%defattr(755,root,root,755)
%_sbindir/obdi
%_initrddir/obdi
%_sysconfdir/obdi/certs
%defattr(644,root,root,755)
%dir %{_var}/cache/obdi
%dir %{_sharedstatedir}/obdi
%dir %{_sharedstatedir}/obdi/plugins
%dir %{_usr}/lib/obdi/plugins/
%dir %{_usr}/lib/obdi/plugins/src/
%defattr(644,root,root,755)
%{_datarootdir}/obdi/static
##%doc doc/initial-config doc/initial-config.dcc doc/bulk-loading README doc/README.html
%config(noreplace) /etc/obdi/obdi.conf

%files worker
%defattr(755,root,root,755)
%_sbindir/obdi-worker
%_initrddir/obdi-worker
%_sysconfdir/obdi-worker/certs
%config(noreplace) /etc/obdi-worker/obdi-worker.conf

%clean
%{__rm} -rf %{buildroot}

%changelog
* Tue Oct 28 2014 Mark Clarkson <mark.clarkson@smorg.co.uk>
- First packaged version
