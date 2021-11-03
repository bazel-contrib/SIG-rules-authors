# Bazel Rules Authors SIG

Bazel is a full-stack monorepo build & test tool from Google. <https://bazel.build>

Bazel needs plugins, called rulesets, which teach it new languages and frameworks.
Some rulesets are developed by Google, but most are community-maintained.

> The term "rulesets" is used in this document to refer to all Bazel extensions,
> including Starlark shared libraries and utilities for end users to work with rulesets.

A SIG (Special Interest Group) can be formed under the Bazel project following the
guidance in the [Bazel proposal to host SIGs].

> As of November 2021 this is still an RFC and not yet approved by Google.

The Rules Authors SIG charter is for ruleset authors
to share technical approaches for solving common problems,
to have a single coherent voice for interacting with the core Bazel team, and
to provide a more consistent experience for Bazel end-users.

More details are in the [Community proposal] to form the Rules Authors SIG.

# Participating in the SIG

## As a company

The SIG is funded by companies that rely on the community-maintained rulesets.
We need your support to continue providing the software you depend on!

> As of November 2021, Google is _not_ a funder of the SIG.

Companies can contribute in several ways:

- _Intellectual property_: Upstream fixes and features your organization has made. Donate proprietary rulesets developed in-house.
- _Engineering time_: Give your enthusiastic developers some dedicated "20% time" to make targeted contributions that benefit your use cases.
- _Financial support_: The SIG plans to accept monetary contribution, likely using <https://opencollective.com>.

Participation benefits your company:

- Recruit and retain talented engineers who want to work in open source.
- Attribution of your contributions builds respect for your brand in the community.
- Avoid merge conflicts when upstream changes break your private patches to Bazel rules.

Please get in touch with us if you think your company may be interested. See the contact info below.

## As a contributor

If you maintain a ruleset, you can ask to join the [Members of the SIG].

> We have not yet determined how to admit members, see https://github.com/bazel-contrib/SIG-rules-authors/issues/1

# Resources

Contact the SIG:

- on Slack: `#rules` channel in https://slack.bazel.build
- by Email: bazel-contrib@googlegroups.com
- Email archives: https://groups.google.com/g/bazel-contrib/
- If you need to reach out privately, email the SIG Leads listed below.

Read the [Meeting notes] from prior meetings.

## Members and partners

Full list of [Members of the SIG].

Leads:

- Alex Eagle <alex@aspect.dev>
- Helen Altshuler <helen@engflow.com>
- Keith Smiley <keithbsmiley@gmail.com>

[bazel proposal to host sigs]: https://docs.google.com/document/d/11iOi_J7TxFGJg6q8hKjddtxMuHFO32t1i9g3r9BmU98/edit#heading=h.5mcn15i0e1ch
[community proposal]: https://github.com/bazelbuild/proposals/blob/main/designs/2021-08-10-rules-authors-sig.md
[meeting notes]: https://docs.google.com/document/d/1YGCYAGLzTfqSOgRFVsB8hDz-kEoTgTEKKp9Jd07TJ5c/edit#
[members of the sig]:https://github.com/orgs/bazel-contrib/teams/rules-authors/members
