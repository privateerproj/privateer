# Roadmap

In this pre-alpha state of Privateer, many of the development tasks are simply being tracked in ad hoc "TODO" entries, which can be reviewed in TODO.md. 

Ad-hoc TODO tracking should be removed before v0.2.

This roadmap should be moved into a GitHub Project before v0.3.

## v0.1

### Feature Additions

- [x] Write plugin logs to independent files
- [x] Write end-to-end summary to independent file
- [x] Create Quickstart guide:
  - can be in readme or elsewhere as appropriate
  - Just fix the readme in general!
- [ ] Install trusted packages if they are not found:
  - finish download logic... unzip

### Feature Improvements

- [x] remove redundancy of plugin name statements
- [x] Create a sample close handler on the raid wireframe

### Bugfixes

- [x] unnecessary error message when no config was supplied (skipping the config file may be intentional, or maybe not)
- [x] config wasn't being used in raid
- [x] Corrected -v usage and removed any instance of `log` in favor of `logger`. Basic `log` should not be used due to persistent unexpected behavior.
- [x] default loglevel is now error

## v0.2

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
- Increase test coverage

### Bugfixes

- Installation is attempting even when package is present
- log.Print usage will result in duplication of timestamp & loglevel
- [ ] Possible Inconsistent RPC error(s)
  - [ERROR] wireframe: plugin: plugin server: accept unix /var/folders/mv/x9vm780x6l755g028989fy500000gn/T/plugin3010505469: use of closed network connection: 

## v0.3

- Remote keystore support (etcd, consul, etc)
- Create website: privateerproj.com
- [ ] Accept all .yml files in a given config directory
  - ref:
    - https://github.com/argoproj/argo-cd/blob/master/cmd/main.go
    - https://github.com/argoproj/argo-cd/blob/master/cmd/argocd/commands/root.go
