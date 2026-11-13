---
nav_order: 2
---

# Hosting rules repositories in the bazel-contrib org

The rules authors SIG hosts some rulesets in our GitHub org. Benefits of this setup include:

- **Community Governance** eliminates single points of failure and improves trust.
Rather than host under their personal account, Rules authors can make their project appear more stable and mature with the reputation of a group backing it.
Similarly, corporate owners can lose interest in a project when a maintainer leaves the company.
By donating the project, the owners are demonstrating a commitment to the open-source community.
The SIG also has some funding to help keep projects maintained.

- **Intellectual property donations** allow companies a way to contribute to the community in addition to money or engineer's time.

- **Standardization** benefits the project by optionally allowing things like linting and release mechanics to be similar to other projects in the org.
This reduces the overall cognitive burden of working across rulesets and encourages sharing of ideas.

> This policy was developed following discussion in <https://github.com/bazel-contrib/SIG-rules-authors/issues/3>
> and much of the content was copied from <https://github.com/MobileNativeFoundation/foundation/pull/12/files>.
> Thanks to the Mobile Native Foundation for providing a similar resource!

## Adding a repo

### Adding Criteria

We want the bazel-contrib org to be trusted, so it shouldn't fall into disrepair or host "abandonware".

Criteria for accepting rule sets into this repo and avoid them from getting stale:

1. Must use an [OSI-approved OSS license](https://opensource.org/licenses), preferably [Apache-2.0](https://www.apache.org/licenses/LICENSE-2.0).
1. Must have wide applicability in the community.
1. Must have a clear point of contact who answers questions from the SIG.
1. Must be "production quality":
    - clear README or other documentation outlining the goal of these rules, how to use them etc.
    - generated API documentation
    - include examples of use
    - tests that are running continuously
1. Must reply to issues/PRs in 2-3 weeks (exact service level agreement TBD)
1. Must have more than one person who is committed to review/approve PRs
    - We recommend encoding this as a `CODEOWNERS` file.
1. Must publish guidance on the README about the support and maintenance commitment for example:
    - Minimally maintained on a best-effort basis by volunteers in their spare time
    - Maintained and supported by [Some Company]
    - Funded by OpenCollective sponsorship on [link to open collective project]
1. Must include guidance for contributors, see https://docs.github.com/en/communities/setting-up-your-project-for-healthy-contributions/setting-guidelines-for-repository-contributors
1. Must publish semver releases.
    - Optional: follow the same release pattern as the rules-template does.
1. Must work with LTS Bazel version
1. Must publish the rules to the [Bazel Central Registry](https://github.com/bazelbuild/bazel-central-registry), keep that CI green

Where possible, the SIG would also prefer to reduce fragmentation.
New projects will be evaluated for whether they duplicate existing ones, and if so whether that is warranted.

### Adding Procedure

1. Rules authors can reach out to the SIG by filing an issue suggesting the repo be added to our org. The issue should demonstrate that the project meets the criteria above.
1. The SIG will meet and determine whether the criteria are satisfied. This is solely up to discretion of the members attending the meeting.
1. The result will be posted back to the issue.
1. The requestor will work with a bazel-contrib repo admin to perform the actual move. GitHub does a great job of forwarding URLs to the new org.

## Archiving a repo

### Archive Criteria

Software projects don't live forever. Eventually, a project may be archived to indicate that it is no longer maintained and that the rules authors SIG no longer recommends its use. The SIG will always discuss a project’s future with its maintainers to determine the best path forward, which could include finding new maintainers for the project.

Good candidates for archiving include:

1. Rules with no active maintenance: maintainers may decide to abandon a project.
1. Rules with no active usage: projects may fall out of use over time. They may have been replaced by first-party tooling in Bazel, or the community may have adopted a competing solution.

We may periodically audit the repositories in the org.
More often, we expect that this procedure will begin as a reaction to some user report of a repo that appears unmaintained.

>  In the future, we hope that the [rules catalog](https://github.com/bazel-contrib/SIG-rules-authors/issues/2) will make it more obvious when a project should be considered for archiving.

Archived projects will use GitHub’s “archived” status: code and past releases will still be available.

### Audit Metrics

An audit should record the date, data sources, observation periods, and gaps in coverage.
Use the following signals to identify projects that need a discussion with their maintainers.
The thresholds are review triggers, not automatic reasons to archive a repository.

#### Maintenance and responsiveness

| Signal | What to measure | When to investigate |
| --- | --- | --- |
| Unanswered contributions | Count and proportion of human issues and ready-for-review PRs with no maintainer response after 21 days. Include the oldest waiting item. | Repeated misses of the existing 2–3 week response expectation, especially for bug reports and proposed fixes. |
| Maintainer's turn | Count issues and PRs where the next action belongs to a maintainer, and days since that action became pending. | Items waiting at least 90 days for triage, review, a merge decision, or a release. |
| Incoming work and throughput | Over the last 90 days, count new human issues/PRs, closed issues, merged PRs, and PRs closed without merging separately. | A growing queue of actionable work with little maintainer response or resolution. |
| Human maintenance | Dates of the last human code change on the default branch, review, and issue triage, and the number of active maintainers over the last 180 days. | No observed human maintenance for 180 days, or fewer than two people actively reviewing or maintaining the project. |
| Releases | Date of the latest stable release, and merged fixes still awaiting release. For rulesets, check whether fixes reached the BCR. | Useful fixes remain unreleased for 90 days, or no stable release for 12 months despite unresolved maintenance needs. |
| Compatibility and CI | Current CI/BCR results for supported Bazel LTS versions, and unresolved compatibility or security reports with links to proposed fixes. | Reproducible failures or security problems remain unaddressed, particularly when contributors have offered fixes. |

Exclude bots and automated acknowledgments from human response counts; report bot work separately.
A reply must come from someone other than the author; PR reviews and inline review comments count too.
Identify maintainers from repository responsibilities and review permissions, confirmed with the project.
Organization membership alone does not establish responsibility for a particular repository.

Inspect the discussion before labeling an item "maintainer's turn".
Exclude drafts and items waiting for information, author changes, or an external dependency from that queue.
The last commenter alone does not establish whose turn it is.
Report counts with their denominators and distinguish outside contributions from maintainers' own work.
Include closed items when measuring response rates over a period, so the remaining open backlog does not bias the result.

Recent pushes, dependency updates, or automated releases can hide a lack of human maintenance.
Conversely, stable software may need few changes or releases.
Evaluate these metrics against the project's stated support commitment and actual unresolved needs.

#### Usage

Record usage per calendar month, preferably for the last 12 complete months.
Compare the latest three complete months with the preceding three to identify trends.
Use consistent sources and coverage for those comparisons, and flag any changes in collection.
Use the best available source for each project: BCR fetch counts, instrumented artifact downloads,
GitHub release-asset downloads, or repository traffic.
Also check known downstream users and whether a supported replacement exists.

- State exactly what is counted: registry metadata requests, artifact downloads, clones, or views.
  These are different measures; do not add them together or describe them as unique users.
- Record which modules, versions, artifacts, and mirrors are covered, and any known effects from CI, bots, and caches.
  Exclude checksums, signatures, attestations, and documentation assets from software download counts.
- Derive monthly release-asset downloads from differences between cumulative snapshots at month boundaries.
  A single lifetime download count divided by a project's age is not evidence of current monthly usage.
- If only a shorter interval is available, report its dates and estimate a monthly rate as
  `observed count / observed days * 30.44`, clearly labeled as an estimate.
- Mark unavailable data as unknown, not zero. Stars and forks provide context but do not establish active usage.

There is no universal download threshold for archiving: a small tool and a widely shared ruleset have different audiences.
Sustained low or declining measured usage is a reason to ask whether the project is still needed.
Missing telemetry is not evidence that a project has no users.

#### Audit Outcome

Publish a ranked list of candidates with linked evidence, measurement dates, usage coverage, and reasons for each recommendation.
Give priority to projects with multiple maintenance failures or confirmed retirement and little evidence of continued need.
For projects with ongoing demand, seek additional maintainers, a transfer, or a supported migration path.
Record contrary evidence and uncertainty, including known users and potential new maintainers.
Audit findings inform the Archive Procedure below; they do not replace maintainer consultation, community discussion, or the SIG vote.

These signals draw on the [CHAOSS Starter Project Health model](https://www.chaoss.community/kb/metrics-model-starter-project-health/),
including response time, contribution throughput, contributor concentration, and release frequency.

### Archive Procedure

To archive a project:

1. The SIG will first reach out to the project’s maintainers to discuss the current state of the project.
1. The SIG will start a discussion with the broader community by posting both a proposal on the discussion forum and an issue on the project’s repository (see template below). The proposal will be open for at least a month.
1. The SIG must vote in favor of archiving the project.
1. The SIG will create a final issue in the project repository indicating that the project will be archived.
1. As soon as two weeks later, the project will be archived.

Template for Github issue:
"This project appears [to not have much usage] [to lack maintainers/contributors]. If that’s the case, it might be a good candidate for archiving, per policy: (link to this document). If anyone depends on this project and would rather it not be archived, please comment on this issue explaining why archiving isn't appropriate."

