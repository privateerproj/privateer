# Privateer - UNDER CONSTRUCTION

Privateer analyzes the complex behaviours and interactions in your cloud infrastructure resources to enable engineers, developers and operations teams identify and fix security vulnerabilities at different points in the lifecycle.

Designed to comprehensively test aspects of security and compliance, Privateer excels where static code inspection or configuration inspection alone are not enough. 

In an era when trusted cloud providers have been seen to fumble massively with sensitive data, Privateer operates whereever you are to bring security issues to the surface quickly and effectively.

### How it works

Privateer is able to execute one or more raids to test the behaviours of your infrastructure services. Upon completion of the raids, Privateer will return a machine-readable set of structured results that can be integrated into the broader DevSecOps process for decision-making. 

Actions within a raid could be as simple as deploying a Kubernetes Pod and running a command inside of it, to complex control and data plane interactions. If your resource is connected to the internet, then Privateer can raid it.

## Architecture

The architecture consists of _Privateer_ as the core executor and _raids_ containing validation checks for specific services. We have built a number of raids, but you can also build your own raid using the [Privateer SDK](tbd) and following the [Raid Template](tbd).

## Quickstart Guide

### TODO:

- The new quickstart guide will depend heavily on the upcoming development efforts being taken as part of the fork from Probr to Privateer.
- The current `-h` isn't very useful.
- The hugo site has not yet been touched from the old site