%define name obdi
%define version %{?OBDI_SEMANTIC_VERSION}
%define release %{?BUILD_NUMBER}
# Don't strip the Golang binaries
%define __strip /bin/true
# The following line may be required
#%define debug_package %{nil}

Summary: A REST interface and GUI for deploying software
# Don't replace the next line with %{name}, hudson build script needs it
Name: obdi
Version: %{version}
Release: %{release}
License: GPL
Group: Application/System
Source: obdi-%{version}.tar.gz
Requires: git
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

%post worker
if [ "$1" = 1 ]; then
    # New install
    :
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

# Include Golang binaries
export PATH=$PATH:/usr/sbin:/usr/local/go/bin

# Golang requirements
export GOROOT=/usr/local/go
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
install -d -m 755 ${RPM_BUILD_ROOT}/%{_sysconfdir}/
install -D -m 640 conf/obdi.conf ${RPM_BUILD_ROOT}/%_sysconfdir/obdi/obdi.conf
install -D -m 640 conf/obdi-worker.conf ${RPM_BUILD_ROOT}/%_sysconfdir/obdi-worker/obdi-worker.conf

# Golang single binary (/usr/sbin)
install -D -m 755 obdi/obdi ${RPM_BUILD_ROOT}%{_sbindir}/obdi
install -D -m 755 obdi-worker/obdi-worker ${RPM_BUILD_ROOT}/%{_sbindir}/obdi-worker

# Initrd
install -D -m 755 init/obdi ${RPM_BUILD_ROOT}/%{_initrddir}/obdi
install -D -m 755 init/obdi-worker ${RPM_BUILD_ROOT}/%{_initrddir}/obdi-worker

# Directories

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
%defattr(644,root,root,755)
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
%config(noreplace) /etc/obdi-worker/obdi-worker.conf

%clean
%{__rm} -rf %{buildroot}

%changelog
* Tue Oct 28 2014 Mark Clarkson <mark.clarkson@smorg.co.uk>
- First packaged version
