# TODO

In-line todo list of things to tackle during the Probr->Privateer overhaul.

## WORKING NOTES (`make todo`)
- [ ] Improve config handling - remove redundancy, keep adhoc plugin calls
- [ ] Improve log formatting - user friendly by default
- [ ] Improve machine output - trim the fat
- [ ] Improve version handling, a la ArgoCD
- [ ] Review CLI options
    - should package/command be pvt or privateer?
    - a bug seems to exist... redo the CLI/opts entirely: flags.go
    - change from reading config.yml to reading specified input file, or all .yml files in directory
    - https://github.com/argoproj/argo-cd/blob/master/cmd/main.go
    - https://github.com/argoproj/argo-cd/blob/master/cmd/argocd/commands/root.go
- [ ] Create Quickstart guide
- Website rework
    - [ ] Fix name references
- [ ] finish download logic... unzip
- [ ] handle errors better for unknown sally raids
- [ ] default loglevel should be error
- [?] for some reason the raid is logging info by default
- [ ] installation is attempting even when package is present
- [ ] change RPC address log to TRACE
- [ ] config isn't being used in raid
- [ ] is trace doing anything inside the raid? or at all? is -v doing anything? is log.fmt breaking the logs? 
- [ ] running from git bash just now, available raids used the full path
- [ ] 2023-08-17T14:39:49.949-0500 [ERROR] open /Users/knight/dev/privateerproj/privateer/config.yml: no such file or directory (unecessary print)
- [ ] InitializeConfig logs the wrong loglevel (trace appears even if level is higher)
