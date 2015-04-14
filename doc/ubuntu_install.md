# Install Obdi on Ubuntu Trusty

## Installation
 
The following steps will install Obdi from the Ubuntu Personal Package Archive (PPA).
```
apt-get install software-properties-common
add-apt-repository ppa:mark-clarkson/obdi
apt-get update
apt-get install obdi obdi-worker
```
Obdi should now be installed and running.

## Testing

The Admin and Run interfaces can be tested now that Obdi is installed.

Only the 'admin' user can log into the Admin interface.
Use a Web Browser to connect to the Admin interface at 'https://\<HOST_OR_IP\>/manager/admin'.
The log-in screen should be shown, so log in as the 'admin' user, with the password 'admin'.

Create a user by clicking 'Add User', then filling in the fields, and don't forget to check the 'Enabled' box.

Now log into the Run interface. The URL for this is 'https://\<HOST_OR_IP\>/manager/run'.
There's nothing to do there except log out, since Obdi needs plugins to do anything useful.

## Configuration
There are a number of settings in `/etc/obdi/obdi.conf` and `/etc/obdi-worker/obdi-worker.conf`. The options are documented in those files.

## Next Steps
Plugins will need to be installed to be able to do anything useful with Obdi.
See the [Plugins](https://github.com/mclarkson/obdi/doc/plugins.md) wiki page for information on installing plugins.

