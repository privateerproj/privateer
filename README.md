[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/privateerproj/privateer/badge)](https://securityscorecards.dev/viewer/?uri=github.com/privateerproj/privateer)

# Privateer - UNDER CONSTRUCTION

This interface enables the quick execution of Privateer Raids,
with a shared input and output if multiple are executed.

Several Privateer commands use unconventional terms
to encourage users to act carefully when using this CLI.
This is due to the fact that your Privateer config is likely
to contain secrets that can be destructive if misused.

The "sally" command will start all requested raids.
Raids are intended to directly interact with running services
and only should be used with caution and proper planning.
Never use a custom-built raid from an unknown source.

You may also streamline the creation of
a new Raid using the generate-raid command, or
the creation of Strikes for a Raid using generate-strike.
Review the help documentation for each command to learn more.

## This whole doc needs refined... x.x

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