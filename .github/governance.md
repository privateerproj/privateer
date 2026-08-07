# Privateer Governance

This document defines the roles, responsibilities, and decision-making processes for the Privateer project (`privateer`, `privateer-sdk`, and related repositories under the [privateerproj](https://github.com/privateerproj) organization).

Privateer does not have a formal collegiate body in charge of steering. Decisions are guided by the consensus of community members who have achieved maintainer status.

While maintainer consensus shall be the process for decision making, all issues and proposals shall be governed by the project's Guiding Governance Principles.

This governance explains how the project is run.

- [Values](#values)
- [Community Roles](#community-roles)
- [Maintainers](#maintainers)
  - [Becoming a Maintainer](#becoming-a-maintainer)
  - [Removing a Maintainer](#removing-a-maintainer)
  - [Emeritus Maintainers](#emeritus-maintainers)
- [Meetings](#meetings)
- [Community Resources](#community-resources)
- [Code of Conduct](#code-of-conduct)
- [Security Response Team](#security-response-team)
- [Voting](#voting)
- [Modifying this Governance](#modifying-this-governance)

## Values

The Privateer project and its leadership embrace the following values:

* Openness: Communication and decision-making happens in the open and is discoverable for future
  reference. As much as possible, all discussions and work take place in public
  forums and open repositories.

* Fairness: All stakeholders have the opportunity to provide feedback and submit
  contributions, which will be considered on their merits.

* Community over Product or Company: Sustaining and growing our community takes
  priority over shipping code or sponsors' organizational goals.  Each
  contributor participates in the project as an individual.

* Vendor Neutrality: The project direction and decisions are not controlled by
  any single organization. Maintainer selection, roadmap prioritization, and
  release decisions are made based on project merit, not employer affiliation.

* Inclusivity: We innovate through different perspectives and skill sets, which
  can only be accomplished in a welcoming and respectful environment.

* Participation: Responsibilities within the project are earned through
  participation, and there is a clear path up the contributor ladder into leadership
  positions.
  
 # Community Roles

Everyone is welcome to contribute through discussion, issues, and pull requests.

The following are roles and additional responsibilities that a person may recieve in the community.


| Role | Responsibilities | Requirements                                                        | Defined by |
|:---|:---|:--------------------------------------------------------------------|:---|
| Member | Active contributor, participates in discussions and reviews | Multiple contributions over time, sponsored by 1 maintainer | GitHub `privateerproj` Organization Member |
| Approver | Review and approve PRs within a specific scope | Member with history of quality reviews | [CODEOWNERS] entry for specific files or directories |
| Core Maintainer | Org-wide oversight, spec authority, binding governance votes | Approver with cross-project contributions | [MAINTAINERS.md] entry and GitHub `privateer-maintainers` Team |
| Community Manager | Outreach, moderation, documentation maintenance (lateral role) | Active community engagement | [MAINTAINERS.md] entry |

For a complete description of all roles, requirements, and promotion processes, see the [Contributor Ladder]. Changes to the [Contributor Ladder] require approval from at least 66% of active maintainers.

## Maintainers

Privateer Maintainers have write access to the project's GitHub repositories.
They may merge their own contributions as well as contributions from other
community members. The current Maintainers are listed in
[MAINTAINERS.md](./MAINTAINERS.md). Together, Maintainers are responsible for
stewarding the project's codebase, community, and technical direction.

Maintainer status carries both privilege and responsibility. Maintainers are
trusted contributors who have consistently demonstrated technical expertise,
collaboration, and commitment to the long-term success of Privateer. They are
expected to review contributions thoughtfully, mentor community members, uphold
project standards, and help resolve issues in both code and documentation.

Active Maintainers collectively govern the Privateer project through open
discussion and consensus. Project decisions are made by the active Maintainers in accordance
with this governance document.

### Becoming a Maintainer

To become a Maintainer, a contributor should demonstrate:

- sustained commitment to the project by:
  - actively participating in discussions, code reviews, and documentation reviews for a sustained period;
  - contributing multiple meaningful pull requests that have been accepted;
  - regularly reviewing and providing constructive feedback on contributions from others;
- the ability to produce high-quality code and/or documentation;
- effective collaboration and communication within the community;
- an understanding of the project's architecture, development workflows, coding standards, and review processes; and
- a commitment to the long-term health of the project and community.

A new Maintainer must be nominated by an existing Maintainer. Following public
community discussion, the active Maintainers may appoint the nominee through a
simple majority vote. Appointments are based solely on merit and demonstrated
contributions, regardless of employer or organizational affiliation.

New Maintainers will be granted the appropriate GitHub permissions and access to
maintainer-only communication channels where necessary.

### Removing a Maintainer

Maintainers may step down at any time if they are no longer able to fulfill
their responsibilities.

A Maintainer may also be removed for reasons including prolonged inactivity,
failure to fulfill maintainer responsibilities, repeated violations of project
policies, or violations of the Code of Conduct.

Unless otherwise communicated, inactivity is generally considered to be little
or no meaningful participation in the project for six consecutive months.

Removal of a Maintainer requires approval by a two-thirds majority vote of the
remaining Maintainers.

### Emeritus Maintainers

Maintainers who step down after making significant contributions may be granted
**Emeritus Maintainer** status.

Emeritus Maintainers are recognized for their past service and may continue to
provide advice and historical context to the project. They do not retain voting
rights or repository write access. Emeritus Maintainers are listed separately in
[EMERITUS.md](./EMERITUS.md).

An Emeritus Maintainer may be reinstated through a simple majority vote of the
active Maintainers.

## Meetings

Maintainers are encouraged to participate in regular public community meetings
where project direction, roadmap planning, releases, and technical proposals are
discussed.

Private meetings may be held only when necessary to discuss sensitive matters,
including security vulnerabilities or Code of Conduct reports. All active
Maintainers should be invited to such meetings except any Maintainer directly
involved in the matter under discussion.

## Community Resources

Maintainers may propose requests for community infrastructure, tooling,
administrative resources, or other project support. Such requests may be
discussed publicly during community meetings or through the project's official
communication channels and approved by a simple majority of the active
Maintainers.

The active Maintainers may delegate responsibility for managing specific
resources or community initiatives to trusted contributors when appropriate.

## Code of Conduct

All members of the Privateer community are expected to follow the project's
[Code of Conduct](./CODE_OF_CONDUCT.md).

Reports of Code of Conduct violations will be handled confidentially by the
active Maintainers. Any Maintainer directly involved in a report must recuse
themselves from the discussion and resolution process.

## Security Response Team

The active Maintainers are responsible for ensuring that security reports are
handled promptly and responsibly.

The Maintainers may designate a Security Response Team consisting of two or
more trusted contributors to coordinate vulnerability reports and security
releases. If no separate team is appointed, the active Maintainers will serve
as the Security Response Team.

Security issues should be handled according to the project's
[Security Policy](./SECURITY.md).

## Voting

Privateer operates primarily through
[lazy consensus](https://community.apache.org/committers/lazyConsensus.html),
where proposals are discussed publicly and proceed unless substantial objections
are raised.

When consensus cannot be reached, any Maintainer may request a formal vote.

Unless otherwise specified:

- A simple majority of active Maintainers is required for ordinary project decisions.
- A two-thirds majority of active Maintainers is required for governance changes, Maintainer removal, or other decisions explicitly requiring a supermajority.

## Modifying this Governance

Changes to this Governance document or its supporting policies may be proposed by
any community member.

After public discussion, modifications require approval by a two-thirds
majority vote of the active Maintainers.
