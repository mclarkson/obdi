# Obdi Plugins

## About

Obdi provides a framework for running scripts. A Linux System Administrator generally works on the command line and types lots of commands to get a job done. When these commands are repeated often, and in much the same order, they are usually grouped together to form scripts. Obdi attempts to provide the final step, putting a User Interface around those scripts to make using them easier. Wrapping the scripts is achieved by extending the Obdi User Interface using Plugins.

## Where are they

Plugins are stored in their own Git repositories. To get a plugin, Obdi needs to be pointed at a list that describes the plugin and also tells Obdi where to get the plugin from. This list of plugins, or plugin repository, is also stored in a git repository containing at least one file, named `repodata.json`.

After installation, Obdi has no plugins installed so doesn't do anything useful. Plugin repositories need to be added to be able to install the plugins. There are some plugin repositories available and it's fairly simple to create a plugin repository of your own, and they can be stored in any public or private Git repository.

## Known Plugins

The following Plugins are free to use and licensed under GPLv3.

### Obdi Core Plugins
The Obdi job viewer is here. When a script is run it is called a Job, and has a Job ID and other information that can be viewed with this plugin.

Plugin repository: https://github.com/mclarkson/obdi-core-repository.git

Click the above link to visit the page and see the plugins.

### Salt Stack Plugins
Plugins that wrap around Salt Stack configuration management scripts.

Plugin repository: https://github.com/mclarkson/obdi-salt-repository.git

Click the above link to visit the page and see the plugins.

Additional post-installation instructions and documentation are kept with the plugins so will not be documented here. For some screen shots see the Blog page: http://blogger.smorg.co.uk/2015/01/obdi-in-pictures.html.

## Installing Plugins

Plugins are installed using the Admin interface. Only the 'admin' user can log into the Admin interface.

Once in the Admin interface, click on the 'Plugins' button, then 'Manage Repositories'. To add the 'Obdi Core Plugins' or the 'Salt Stack Plugins', paste the plugin repository http(s) URL into the 'Repository URL' Text Entry control.

Next, click 'Back', then 'Add Plugin', a list of plugins will be shown. Click the green button to install a plugin. Only installable plugins will be shown, so if all the plugins are installed then the list will be empty.

There is some more information and screen shots on the blog page: http://blogger.smorg.co.uk/2015/03/obdi-plugins-completed.html.

## Anatomy of a plugin

![](../images/PluginsOverview.png?raw=true)

## Next Steps

Once the plugins are installed they may need to be configured. Refer to the documentation supplied with each plugin to find out how to use them.

