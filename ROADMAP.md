# Roadmap

In this pre-alpha state of Privateer, many of the development tasks are simply being tracked in ad hoc "TODO" entries, which can be reviewed in TODO.md. 

Ad-hoc TODO tracking should be removed before v0.2.

This roadmap should be moved into a GitHub Project before v0.3.

## v0.1 - August 2023

### Feature Additions

- Install trusted packages if they are not found:
  - finish download logic... unzip
- Write plugin logs to independent files
- Write end-to-end summary to independent file

### Feature Improvements

- Configuration handling:
  - remove redundancy of plugin name statements
  - keep adhoc plugin calls
- Review CLI options:
  - should package/command be pvt or privateer?
  - a bug seems to exist... redo the CLI/opts entirely: flags.go
  - change from reading config.yml to reading specified input file, or all .yml files in directory
  - ref:
    - https://github.com/argoproj/argo-cd/blob/master/cmd/main.go
    - https://github.com/argoproj/argo-cd/blob/master/cmd/argocd/commands/root.go
- Create Quickstart guide:
  - can be in readme or elsewhere as appropriate
  - Just fix the readme in general!
- Log Handling:
  - change RPC address log to TRACE

### Bugfixes

- unecessary error message when reading config
  - 2023-08-17T14:39:49.949-0500 [ERROR] open /Users/knight/dev/privateerproj/privateer/config.yml: no such file or directory
  - Fixed this by making this a Debug message (skipping the config file may be intentional, or maybe not)
- config isn't being used in raid
  - 
- is trace doing anything inside the raid? or at all? is -v doing anything? is log.fmt breaking the logs? 
  - Corrected -v usage and removed any instance of `log` in favor of `logger`. Basic `log` should not be used.
  - default loglevel is now error
- Inconsistent RPC error(s)
  - [ERROR] wireframe: plugin: plugin server: accept unix /var/folders/mv/x9vm780x6l755g028989fy500000gn/T/plugin3010505469: use of closed network connection: 

## v0.2 - September 2023

### Feature Additions

- Improve version handling
  - a la ArgoCD

- Secret handling
  - plugins should not be able to read configs from other plugins!

### Feature Improvements

- Improve log formatting
  - user friendly by default
- Improve machine output
  - trim the fat
- Handle errors better for unknown sally raids

### Bugfixes

- Installation is attempting even when package is present

## v0.3 - October 2023

- Remote keystore support (etcd, consul, etc)
- Create website: privateerproj.com
