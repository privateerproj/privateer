# TODO

In-line todo list of things to tackle during the Probr->Privateer overhaul.

NOW
- [ ] Organize this list by priority

HIGH PRIORITY
- Fix all "Probr" to "Privateer"
    - [x] go files
    - [x] build files
    - [x] docs
- [x] Fix all "service pack" to "raid"

EXTERNAL TASKS
- [x] Rename & Adjust pack naming convention
    - ie, servicePack.RunProbes() to raid.Start()
- [x] Create & test "HelloWorld" Raid
- [ ] Remove all BDD logic from raid

SECONDARY PRIORITY
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



LOW PRIORITY
- [ ] Create Quickstart guide
- Website rework
    - [ ] Fix name references
