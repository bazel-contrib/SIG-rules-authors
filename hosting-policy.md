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

## Criteria for adding a repo

We want the bazel-contrib org to be trusted, so it shouldn't fall into disrepair or host "abandonware".

Criteria for accepting rule sets into this repo and avoid them from getting stale:

1. Must use an open-source license, preferably [Apache-2.0](https://www.apache.org/licenses/LICENSE-2.0).
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
1. Must publish semver releases.
    - Optional: follow the same release pattern as the rules-template does.
1. Must work with LTS Bazel version

Where possible, the SIG would also prefer to reduce fragmentation.
New projects will be evaluated for whether they duplicate existing ones, and if so whether that is warranted.

When the [bzlmod](https://docs.bazel.build/versions/5.0.0/bzlmod.html) feature graduates from experimental,
we'll also add the following criteria:

1. Include the rules in bazel-central-registry, keep that CI green

## Procedure for archiving a repo

Software projects don't live forever. Eventually, a project may be archived to indicate that it is no longer maintained and that the rules authors SIG no longer recommends its use. The SIG will always discuss a project’s future with its maintainers to determine the best path forward, which could include finding new maintainers for the project.

Good candidates for archiving include:

1. Rules with no active maintenance: maintainers may decide to abandon a project.
1. Rules with no active usage: projects may fall out of use over time. They may have been replaced by first-party tooling in Bazel, or the community may have adopted a competing solution.

We may periodically audit the repositories in the org.
More often, we expect that this procedure will begin as a reaction to some user report of a repo that appears unmaintained.

>  In the future, we hope that the [rules catalog](https://github.com/bazel-contrib/SIG-rules-authors/issues/2) will make it more obvious when a project should be considered for archiving.

Archived projects will use GitHub’s “archived” status: code and past releases will still be available.

To archive a project:

1. The SIG will first reach out to the project’s maintainers to discuss the current state of the project.
1. The SIG will start a discussion with the broader community by posting both a proposal on the discussion forum and an issue on the project’s repository (see template below). The proposal will be open for at least a month.
1. The SIG must vote in favor of archiving the project.
1. The SIG will create a final issue in the project repository indicating that the project will be archived.
1. As soon as two weeks later, the project will be archived.

Template for Github issue:
"This project appears [to not have much usage] [to lack maintainers/contributors]. If that’s the case, it might be a good candidate for archiving, per policy: (link to this document). If anyone depends on this project and would rather it not be archived, please comment on this issue explaining why archiving isn't appropriate."

