%define name obdi
%define version 1
%define php php
%if "%{?dist}" == ".el5"
%define php php53
%endif
# The following line may be required
#%define debug_package %{nil}

Summary: Nagios REST configuration tools.
# Don't replace the next line with %{name}, hudson build script needs it
Name: nagrestconf
Version: %{version}
Release: 1
License: GPL
Group: Application/System
Source: nagrestconf-%{version}.tar.gz
Requires: bash, grep, nagios >= 3, procmail, sed, gawk, grep, %php >= 5.3, httpd, mod_ssl, subversion
# PreReq: sh-utils
BuildArch: noarch
BuildRoot: %{_builddir}/%{name}-%{version}/tmp
Packager: Mark Clarkson
Vendor: Smorg

%description
Configuration tools for Nagios.

This package provides csv2nag, nagctl, the REST interface and the web
configurator GUI.

%package services-tab-plugin
Group: Application/System
Summary: Services Tab plugin for Nagrestconf.
Requires: nagrestconf

%description services-tab-plugin
Configuration tools for Nagios.

This package provides the 'Services Tab' plugin.

%package services-bulktools-plugin
Group: Application/System
Summary: Services Tab Bulk Tools plugin for Nagrestconf.
Requires: nagrestconf, nagrestconf-services-tab-plugin

%description services-bulktools-plugin
Configuration tools for Nagios.

This package provides the 'Bulk Tools' plugin for the Services tab.

%package hosts-bulktools-plugin
Group: Application/System
Summary: Hosts Tab Bulk Tools plugin for Nagrestconf.
Requires: nagrestconf

%description hosts-bulktools-plugin
Configuration tools for Nagios.

This package provides the 'Bulk Tools' plugin for the Hosts tab.

%package backup-plugin
Group: Application/System
Summary: Backup and restore plugin for Nagrestconf.
Requires: nagrestconf

%description backup-plugin
Configuration tools for Nagios.

This package provides the 'Backup & Restore' plugin.

%prep
%setup -q

# Pre Install
%pre

# Post Install
%post

%post services-tab-plugin
if [ "$1" = 1 ]; then
    # New install
    %__ln_s ../plugins/smorg_services_tab.php /usr/share/nagrestconf/htdocs/nagrestconf/plugins-enabled/10_smorg_services_tab.php
fi

%post services-bulktools-plugin
if [ "$1" = 1 ]; then
    # New install
    %__ln_s ../plugins/smorg_services_bulktools_btn.php /usr/share/nagrestconf/htdocs/nagrestconf/plugins-enabled/50_smorg_services_bulktools_btn.php
fi

%post hosts-bulktools-plugin
if [ "$1" = 1 ]; then
    # New install
    %__ln_s ../plugins/smorg_hosts_bulktools_btn.php /usr/share/nagrestconf/htdocs/nagrestconf/plugins-enabled/50_smorg_hosts_bulktools_btn.php
fi

%post backup-plugin
if [ "$1" = 1 ]; then
    # New install
    %__ln_s ../plugins/smorg_backup_btn.php /usr/share/nagrestconf/htdocs/nagrestconf/plugins-enabled/06_smorg_backup_btn.php
fi

# Pre Uninstall
%preun

# Post Uninstall
%postun

%preun services-tab-plugin
if [ "$1" = 0 ]; then
    # uninstall
    %__rm -f /usr/share/nagrestconf/htdocs/nagrestconf/plugins-enabled/10_smorg_services_tab.php
fi

%preun services-bulktools-plugin
if [ "$1" = 0 ]; then
    # uninstall
    %__rm -f /usr/share/nagrestconf/htdocs/nagrestconf/plugins-enabled/50_smorg_services_bulktools_btn.php
fi

%preun hosts-bulktools-plugin
if [ "$1" = 0 ]; then
    # uninstall
    %__rm -f /usr/share/nagrestconf/htdocs/nagrestconf/plugins-enabled/50_smorg_hosts_bulktools_btn.php
fi

%preun backup-plugin
if [ "$1" = 0 ]; then
    # uninstall
    %__rm -f /usr/share/nagrestconf/htdocs/nagrestconf/plugins-enabled/06_smorg_backup_btn.php
fi

%install

echo "php requires = %{php}"

[ "$RPM_BUILD_ROOT" != "/" ] && %{__rm} -rf %{buildroot}

# Config
install -d -m 755 ${RPM_BUILD_ROOT}/%_sysconfdir/
#cp -r etc/httpd ${RPM_BUILD_ROOT}/%_sysconfdir/
install -D -m 640 etc/httpd/conf.d/nagrestconf.conf ${RPM_BUILD_ROOT}/%_sysconfdir/httpd/conf.d/nagrestconf.conf
install -D -m 640 etc/httpd/conf.d/rest.conf ${RPM_BUILD_ROOT}/%_sysconfdir/httpd/conf.d/rest.conf
#cp -r etc/nagrestconf ${RPM_BUILD_ROOT}%_sysconfdir/
install -D -m 640 etc/nagrestconf/csv2nag.conf ${RPM_BUILD_ROOT}/%_sysconfdir/nagrestconf/csv2nag.conf
install -D -m 640 etc/nagrestconf/nagctl.conf ${RPM_BUILD_ROOT}/%_sysconfdir/nagrestconf/nagctl.conf
install -D -m 640 etc/nagrestconf/nagrestconf.ini ${RPM_BUILD_ROOT}/%_sysconfdir/nagrestconf/nagrestconf.ini
install -D -m 640 etc/nagrestconf/restart_nagios.conf ${RPM_BUILD_ROOT}/%_sysconfdir/nagrestconf/restart_nagios.conf

# Scripts
install -D -m 755 scripts/csv2nag ${RPM_BUILD_ROOT}%_bindir/csv2nag
install -D -m 755 scripts/nagctl ${RPM_BUILD_ROOT}%_bindir/nagctl
install -D -m 755 scripts/restart_nagios ${RPM_BUILD_ROOT}%_bindir/restart_nagios
install -D -m 755 scripts/dcc_configure ${RPM_BUILD_ROOT}%_bindir/dcc_configure
install -D -m 755 scripts/slc_configure ${RPM_BUILD_ROOT}%_bindir/slc_configure
install -D -m 755 scripts/upgrade_setup_files.sh ${RPM_BUILD_ROOT}%_bindir/upgrade_setup_files.sh
install -D -m 755 scripts/update_nagios ${RPM_BUILD_ROOT}%_bindir/update_nagios
install -D -m 755 scripts/auto_reschedule_nagios_check ${RPM_BUILD_ROOT}%_bindir/auto_reschedule_nagios_check
install -D -m 755 scripts/nagrestconf_install ${RPM_BUILD_ROOT}%_bindir/nagrestconf_install

# PHP Directories
install -d -m 755 ${RPM_BUILD_ROOT}/usr/share/nagrestconf/htdocs/
cp -r nagrestconf ${RPM_BUILD_ROOT}/usr/share/nagrestconf/htdocs/
cp -r rest ${RPM_BUILD_ROOT}/usr/share/nagrestconf/htdocs/

install -d -m 755 ${RPM_BUILD_ROOT}/usr/share/nagrestconf/htdocs/nagrestconf/scripts/
install -d -m 755 ${RPM_BUILD_ROOT}/usr/share/nagrestconf/htdocs/nagrestconf/upload/
install -d -m 755 ${RPM_BUILD_ROOT}/usr/share/nagrestconf/htdocs/nagrestconf/download/
install -d -m 755 ${RPM_BUILD_ROOT}/usr/share/nagrestconf/htdocs/nagrestconf/plugins/
install -d -m 755 ${RPM_BUILD_ROOT}/usr/share/nagrestconf/htdocs/nagrestconf/plugins-lib/
install -d -m 755 ${RPM_BUILD_ROOT}/usr/share/nagrestconf/htdocs/nagrestconf/plugins-enabled/

# GUI Plugins
install -D -m 755 plugins/smorg_services_tab_impl.php ${RPM_BUILD_ROOT}/usr/share/nagrestconf/htdocs/nagrestconf/plugins-lib/
install -D -m 755 plugins/smorg_services_tab.php ${RPM_BUILD_ROOT}/usr/share/nagrestconf/htdocs/nagrestconf/plugins/
install -D -m 755 plugins/smorg_services_bulktools_btn_impl.php ${RPM_BUILD_ROOT}/usr/share/nagrestconf/htdocs/nagrestconf/plugins-lib/
install -D -m 755 plugins/smorg_services_bulktools_btn.php ${RPM_BUILD_ROOT}/usr/share/nagrestconf/htdocs/nagrestconf/plugins/
install -D -m 755 plugins/smorg_hosts_bulktools_btn_impl.php ${RPM_BUILD_ROOT}/usr/share/nagrestconf/htdocs/nagrestconf/plugins-lib/
install -D -m 755 plugins/smorg_hosts_bulktools_btn.php ${RPM_BUILD_ROOT}/usr/share/nagrestconf/htdocs/nagrestconf/plugins/
install -D -m 755 plugins/smorg_backup_btn_impl.php ${RPM_BUILD_ROOT}/usr/share/nagrestconf/htdocs/nagrestconf/plugins-lib/
install -D -m 755 plugins/smorg_backup_btn.php ${RPM_BUILD_ROOT}/usr/share/nagrestconf/htdocs/nagrestconf/plugins/

%files
%defattr(755,root,root,755)
%_bindir/csv2nag
%_bindir/nagctl
%_bindir/restart_nagios
%_bindir/dcc_configure
%_bindir/slc_configure
%_bindir/upgrade_setup_files.sh
%_bindir/update_nagios
%_bindir/auto_reschedule_nagios_check
%_bindir/nagrestconf_install
/usr/share/nagrestconf/htdocs/nagrestconf/scripts/csv2json.sh
%defattr(644,root,root,755)
%dir /usr/share/nagrestconf/htdocs/nagrestconf/plugins-lib
%dir /usr/share/nagrestconf/htdocs/nagrestconf/plugins-enabled
%dir /usr/share/nagrestconf/htdocs/nagrestconf/plugins
%dir /usr/share/nagrestconf/htdocs/nagrestconf/scripts
%defattr(644,apache,apache,755)
%dir /usr/share/nagrestconf/htdocs/nagrestconf/upload
%dir /usr/share/nagrestconf/htdocs/nagrestconf/download
%defattr(644,root,root,755)
/usr/share/nagrestconf/htdocs/nagrestconf/css
/usr/share/nagrestconf/htdocs/nagrestconf/fonts
/usr/share/nagrestconf/htdocs/nagrestconf/images
/usr/share/nagrestconf/htdocs/nagrestconf/index.php
/usr/share/nagrestconf/htdocs/nagrestconf/upload.php
/usr/share/nagrestconf/htdocs/nagrestconf/js
/usr/share/nagrestconf/htdocs/nagrestconf/main.css
/usr/share/nagrestconf/htdocs/rest
%doc doc/initial-config doc/initial-config.dcc doc/bulk-loading README doc/README.html
%config(noreplace) /etc/httpd/conf.d/rest.conf
%config(noreplace) /etc/httpd/conf.d/nagrestconf.conf
%config(noreplace) /etc/nagrestconf/nagrestconf.ini
%config(noreplace) /etc/nagrestconf/restart_nagios.conf
%config(noreplace) /etc/nagrestconf/csv2nag.conf
%config(noreplace) /etc/nagrestconf/nagctl.conf

%files services-tab-plugin
%defattr(644,root,root,755)
/usr/share/nagrestconf/htdocs/nagrestconf/plugins-lib/smorg_services_tab_impl.php
/usr/share/nagrestconf/htdocs/nagrestconf/plugins/smorg_services_tab.php

%files services-bulktools-plugin
%defattr(644,root,root,755)
/usr/share/nagrestconf/htdocs/nagrestconf/plugins-lib/smorg_services_bulktools_btn_impl.php
/usr/share/nagrestconf/htdocs/nagrestconf/plugins/smorg_services_bulktools_btn.php

%files hosts-bulktools-plugin
%defattr(644,root,root,755)
/usr/share/nagrestconf/htdocs/nagrestconf/plugins-lib/smorg_hosts_bulktools_btn_impl.php
/usr/share/nagrestconf/htdocs/nagrestconf/plugins/smorg_hosts_bulktools_btn.php

%files backup-plugin
%defattr(644,root,root,755)
/usr/share/nagrestconf/htdocs/nagrestconf/plugins-lib/smorg_backup_btn_impl.php
/usr/share/nagrestconf/htdocs/nagrestconf/plugins/smorg_backup_btn.php
/usr/share/nagrestconf/htdocs/nagrestconf/scripts/restore.php

%clean
%{__rm} -rf %{buildroot}

%changelog
* Mon Aug 17 2013 Mark Clarkson <mark.clarkson@smorg.co.uk>
- Added Bulk Tools plugin for the Hosts tab.

* Mon Aug 02 2013 Mark Clarkson <mark.clarkson@smorg.co.uk>
- Added Bulk Tools plugin for the Services tab.

* Mon Jul 29 2013 Mark Clarkson <mark.clarkson@smorg.co.uk>
- Added Services Tab plugin.

%changelog
* Fri Jan 4 2013 Mark Clarkson <mark.clarkson@smorg.co.uk>
- Moved set code out to nagrestconf_install script.

* Wed Nov 7 2012 Mark Clarkson <mark.clarkson@smorg.co.uk>
- Added schedulehostdowntime and delhostdowntime.

* Wed Nov 7 2012 Mark Clarkson <mark.clarkson@smorg.co.uk>
- Added nagios objects: service dependencies, host dependencies,
  service escalations, host escalations, extra service info and
  extra host info.

* Mon Nov 5 2012 Mark Clarkson <mark.clarkson@smorg.co.uk>
- Added more nagios directives to the REST api. Host/Service groups and timeperiods now complete.

* Fri Nov 2 2012 Mark Clarkson <mark.clarkson@smorg.co.uk>
- Added more nagios directives to the REST api. Host templates, services and servicesets now complete.

* Thu Nov 1 2012 Mark Clarkson <mark.clarkson@smorg.co.uk>
- Added more nagios directives to the REST api. Host and Contacts now complete.

* Fri Jul 6 2012 Mark Clarkson <mark.clarkson@smorg.co.uk>
- Lots of additional configuration tabs, fixes and extra checks.

* Wed May 16 2012 Mark Clarkson <mark.clarkson@smorg.co.uk>
- Speech marks fix. Both types can now be used in service command input boxes.

* Wed May 16 2012 Mark Clarkson <mark.clarkson@smorg.co.uk>
- Bug fixes and UI enhancements. Added disable host and services feature.

* Sun May 13 2012 Mark Clarkson <mark.clarkson@smorg.co.uk>
- Added nagrestconf alpha release.

* Tue Sep 6 2011 Mark Clarkson <mark.clarkson@smorg.co.uk>
- Misc fixes and operates as dcc or slc via DCC variable.

* Tue Oct 5 2010 Mark Clarkson <mark.clarkson@smorg.co.uk>
- First packaged version
