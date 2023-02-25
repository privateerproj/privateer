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
    - [ ] Fix name references- [ ] Figure out what I was wanting to make note of when I got distracted automating this
- [ ] config handling isn't capturing the -v anymore
- [ ] finish download logic... unzip