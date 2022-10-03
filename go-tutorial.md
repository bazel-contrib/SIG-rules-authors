---
nav_order: 3
---

# Go and Bazel

// TODO: is Bazel a proper noun?

## This tutorial covers

- Background of Bazel and Bazel's Go support
- Covering how rules_go and gazelle are used with Bazel and Go
- Creating a basic Go project for the tutorial
- Implementing a WORKSPACE and BUILD.bazel files 
- Using gazelle to generate more WORKSPACE and BUILD.bazel updates
- Utilzing different Bazel commands
- An overview of gazelle and Go dependency management
- Understanding the contents of the WORKSPACE and BUILD.bazel files
- Creating new internal dependencies and using gazelle to update bazel
- Adding an a new external Go dependency and the Go vendoring
- Running and implementing Go tests with bazel
- Learing about other rules in rules_go

## About Bazel and Go

This tutorial is going to cover how Bazel supports the programming language Go.
[Bazel](https://bazel.build) is an open-source build and test application that supports 
the software development lifecycle. Bazel strives to allow a developer to have hermetic and 
deterministic builds, testing, packaging and deployment.  This build tool supports 
multiple languages and cross compilation for different operating systems and hardware architecture.

One of software languages that Bazel supports includes [Go](https://go.dev). Go is an open-source 
programming language that was created by Google. Now that we have some background about Bazel
now lets cover some Bazel concepts.

One of the concepts within bazel is a rule.  A bazel rule defines a series of actions and outputs. Toolchains is another
concept/part of bazel. The toolchain framework, is a way for rule authors to decouple their rule logic from 
platform-based selection of tools. So we need a rule, or usually a set of rules
(ruleset) that provides a toolchain to support a programming language.  
The bazel open-source community maintains
[rules_go](https://github.com/bazelbuild/rules_go).  This ruleset provides the following support:

- Building go libraries, binaries, and tests
- Vendoring and dependency management
- Support for cgo
- The cross-compilation of binaries for different OS and platforms
- Build-time code analysis via nogo
- Support for protocol buffers
- Remote execution
- Code coverage testing
- gopls integration for editor support
- Debugging

The bazel open-source community also provides another tool called [gazelle](https://github.com/bazelbuild/bazel-gazelle).
Gazelle addresses the creation and maintence of [BUILD](https://bazel.build/concepts/build-files) files.  
Every bazel project has BUILD (BUILD.bazel) files that define the various rules that are used within a project.
When you add more code or dependencies to a project you need to update your build files.  When you add a new folder
you need to add another BUILD file. If you have ever worked with bazel you know how much time you spend maintaining
these files, if you maintain the files by hand.  Gazelle was created to reduce the previously mentioned pain points.

> Gazelle is a build file generator for Bazel projects. It can create new BUILD.bazel files for a project that follows language conventions, and it can update existing build files to include new sources, dependencies, and options. Gazelle natively supports Go and protobuf, and it may be extended to support new languages and custom rule sets.
>
>  -- <cite>https://github.com/bazelbuild/bazel-gazelle#gazelle-build-file-generator</cite> 

// TODO list languages
Intially gazelle was created to support Go, and now supports many other languages.

Part of learning Bazel is understanding the configuration language that Bazel uses.
The language is called [StarLark](https://github.com/bazelbuild/starlark).

> Starlark (formerly known as Skylark) is a language intended for use as a configuration language. It was designed for the Bazel build system, but may be useful for other projects as well. This repository is where Starlark features are proposed, discussed, and specified. It contains information about the language, including the specification. There are multiple implementations of Starlark.
>
> Starlark is a dialect of Python. Like Python, it is a dynamically typed language with high-level data types, first-class functions with lexical scope, and garbage collection.
> - <cite>https://github.com/bazelbuild/starlark#overview</cite>

The good news is that Starlark is a dialect of Python, almost a subset of the language.  If you know
Python you have a jump start on learning Starlark.

Before we start going through creating a simple Go project we are going to cover a couple of dependencies for this 
tutorial.

## Dependencies for the tutorial

We use the following dependencies for this tutorial.

- go: https://go.dev/doc/install
- gcc: use your systems package manager
- bazelisk: https://github.com/bazelbuild/bazelisk#installation

Technically we do not need the go binary installed, to use Bazel, but we are going to use
`cobra-cli` to generate some project code.  We did not want to add the 
extra work to run the binary using Bazel. A developer, using go,
does not need to download the go binary.  In order to keep a build deterministic
bazel and rules_go download go. rules_go require that gcc is installed.

We are not installing bazel by hand for this tutorial, but are using Bazelisk.
Bazelisk is a wrapper for Bazel written in Go. It automatically picks the
correct version of Bazel given your current working directory, downloads it from 
the official server (if required) and then transparently passes through all 
command-line arguments to the real Bazel binary.  You can call it just 
like you would call Bazel.

Now how about we actually write some code! We are going to create a
simple Go program and then add Bazel to the project.  We have
structured the tutorial in this manner since at times you migrate
to using Bazel with an existing project, and at other times you
start a new project with Bazel.

## The project

We are going to create a small example project first using go.  As
we mentioned you do not need to use go directly at all, when using bazel.
But to get a "easy" jump start we wanted to quickly generate some code.

// TODO finish the code to do that?

The project is going to consist of a simple cli program that generates a
random number or generates a random word.

## Generate the project framework

First create a git repository to store you work.  For this project we are using
https://github.com/bazel-contrib/Bazel-learning-paths/tutorials/go-code-tutorial, and replace any references
to that repository with your own. You can refer to the above repository for 
the final source code base.

The we are using the [cobra](https://cobra.dev/) CLI framework for this project.
The cobra framework is commonly used by various projects including Kubernetes.
The cobra-cli binary is provided by the project for the intial generation of CLI code.
// TODO include the go call and mentioned the docs if then need them.

Follow the (instructions)[https://github.com/spf13/cobra-cli/blob/main/README.md] and install
cobra-cli.

In the root directory of your project use go mod and init the code vendoring.

// TODO what code formatting is supported.
```
$ go mod init github.com/bazel-contrib/Bazel-learning-paths/tutorials/go-code-tutorial
```

Next use cobra-cli to create go root, rool and word files. Replace 
the NAME variable with your information.

```
$ export NAME="Your Name your@email.com"
$ cobra-cli init -a '${NAME}' --license apache
$ cobra-cli add roll -a '${NAME}' --license apache
$ cobra-cli add word -a '${NAME]' --license apache
```

Run the above commands in the root directory of your project.

You will now have the following files:

```
├── cmd
│   ├── roll.go
│   ├── root.go
│   └── word.go
├── go.mod
├── go.sum
└── main.go
```

Let's add a couple of directories:

```
mkdir -p pkg/{word,roll}
```

Inside of those directories we can add roll_dice.go and generate_word.go files.

In the roll_dice.go file add the following code:

```
package roll

import "fmt"

func Roll() {
        fmt.Println("roll dice")
}
```

In the generate_word.go file add the the following code:

```
package word

import "fmt"

func GenerateWord() {
        fmt.Println("GenerateWord")
}
```
You will end up with the following file structure:

```
├── cmd
│   ├── roll.go
│   ├── root.go
│   └── word.go
├── go.mod
├── go.sum
├── main.go
└── pkg
    ├── roll
    │   └── roll_dice.go
    └── word
        └── generate_word.go
```

Next add a .gitignore file by running the following command.

```
$ tee -a .gitignore << EOF
/bazel-$(basename $(pwd))
/bazel-bazel-gazelle
/bazel-bin
/bazel-out
/bazel-testlogs
EOF
```

Bazel creates various directories in the project root and this file will allow git 
to ignore those directories.

This is a good time to push your files into a remote git repository like GitHub. Now
we cover rules_go and gazelle.

## Go and Bazel

As we mentioned previously bazel provides rules_go and gazelle. You can find more
about them here:

- https://github.com/bazelbuild/rules_go
- https://github.com/bazelbuild/bazel-gazelle

At a high level we use Skylark to define that bazel will use rules from rules_go
to create the Go support within a project. We use gazelle to manage our BUILD.bazel files,
or WORKSPACE files, and other bazel specific files.

If you are not familiar with BUILD.bazel files or WORKSPACE files take a look at:
https://bazel.build/concepts/build-files

Next let's create our WORKSPACE file so that bazel knows it is using rules_go and gazelle.

## Create WORKSPACE file

The [StarLark](https://bazel.build/rules/language) is
used within WORKSPACE and BUILD.bazel files. The definitions within the WORKSPACE files include StarkLark
code for both rules_go and gazelle.

An example WORKSPACE is documented [here](https://github.com/bazelbuild/bazel-gazelle#running-gazelle-with-bazel).

Use your favorite editor and create a file named "WORKSPACE" in the root directory of your project.

Edit the WORKSPACE file and include the following StarLark code.


```
# use http_archive to download bazel rules_go
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "099a9fb96a376ccbbb7d291ed4ecbdfd42f6bc822ab77ae6f1b5cb9e914e94fa",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.35.0/rules_go-v0.35.0.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.35.0/rules_go-v0.35.0.zip",
    ],
)

# use http_archive to download bazel_gazelle dependency
http_archive(
    name = "bazel_gazelle",
    sha256 = "efbbba6ac1a4fd342d5122cbdfdb82aeb2cf2862e35022c752eaddffada7c3f3",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.27.0/bazel-gazelle-v0.27.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.27.0/bazel-gazelle-v0.27.0.tar.gz",
    ],
)

# load bazel and gazelle rules
load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

############################################################
# Define your own dependencies here using go_repository.
# Else, dependencies declared by rules_go/gazelle will be used.
# The first declaration of an external repository "wins".
############################################################

# we are going to store the go dependecy definitions
# in a different file "deps.bzl". We can include those 
# definitions in this file, but it gets quite verbose.
load("//:deps.bzl", "go_dependencies")

# Next we initialize the tool chains

# gazelle:repository_macro deps.bzl%go_dependencies
go_dependencies()

go_rules_dependencies()

# We define the version of go that this project uses
go_register_toolchains(version = "1.19.1")

gazelle_dependencies()
```

The above WORKSPACE file contains specific version numbers for rules_go and gazelle.  Refer to the 
gazelle site to use the latest versions.  Also update the `go_register_toolchains(version = "1.19.1")`
to the version that you would like to use of Go.

Next we need to a BUILD (BUILD.bazel) file in the root project directory.

## Create intial BUILD file

Open your editor and create a file named BUILD.bazel. Write the following contents to the BUILD.bazel
file:

```
# Load the gazelle rule
load("@bazel_gazelle//:def.bzl", "gazelle")

# The following comment defines the import path that corresponds to the repository root directory.
# This is a critical definition, and if you mess this up all of the BUILD file generation 
# will have errors.

# Modify the name to your project name in your git repository.

# gazelle:prefix github.com/bazel-contrib/Bazel-learning-paths/tutorials/go-code-tutorial
gazelle(name = "gazelle")

# Add a rule to call gazelle and pull in new go dependencies.
gazelle(
    name = "gazelle-update-repos",
    args = [
        "-from_file=go.mod",
        "-to_macro=deps.bzl%go_dependencies",
        "-prune",
    ],
    command = "update-repos",
)

```

Again the `gazelle:prefix` is critical.  If the value after the code "prefix:" is not named correctly
gazelle does not update BUILD.bazel file correctly. This value contains the import path
that corresponds to your repository, and drives dependency management. If you
include the incorrect value gazelle will think that a dependency inside of the Go code
lives outside of the repository.

The last rule that we defined is named "gazelle-update-repos".  This is a custom
SkyLark definition that defines a gazelle command, and specific arguments to that command.
Do not run this command yet, but this allows us to run:

```
$ bazelisk run //:gazelle-update-repos
```

Which is the equivilent of running

```
$ bazelist run //:gazelle update-repos -from_file=go.mod -to_macro=deps.bzl%go_dependencies -prune
```

The update-repos command is a very common way of running Gazelle. 
Gazelle scans sources in directories throughout the repository, 
then creates and updates build files. The BUILD.bazel file includes
and alias to run update.

Since we run that command a lot, we create the definition for it.

Now we now have done the intial creation of the WORKSPACE and BUILD.bazel files. 
Next we will use bazel to run the gazelle target.

## Run the gazelle commands

As we previously mentioned we use bazel to run gazelle, and 
gazelle manages the BUILD.bazel files for us.  We are using bazelisk to 
manage and run bazel, but we will typically say "run bazel" 
instead of "run bazelisk".  

Run the following commands to update the root BUILD.bazel, 
the WORKSPACE file, and generate the other BUILD.bazel
files for the project.

```
$ bazelisk run //:gazelle
$ bazelisk run //:gazelle-update-repos
```

You now have the following files:

```
├── BUILD.bazel
├── CREATE.adoc
├── LICENSE
├── WORKSPACE
├── cmd
│   ├── BUILD.bazel
│   ├── roll.go
│   ├── root.go
│   └── word.go
├── deps.bzl
├── go.mod
├── go.sum
├── main.go
└── pkg
    ├── roll
    │   ├── BUILD.bazel
    │   └── roll_dice.go
    └── word
        ├── BUILD.bazel
        └── generate_word.go
```

We now have additional BUILD.bazel files in the cmd and pkg directories.
How about we walk through the StarLark code in the BUILD.bazel and deps.bzl
files?

## The bazel files in the project.

The previous gazelles command updated the BUILD.bazel file in the root directory of the project
and created new BUILD files as well. Here is a layout of the bazel files in the project.

```
├── BUILD.bazel
├── WORKSPACE
├── cmd
│   ├── BUILD.bazel
├── deps.bzl
└── pkg
    ├── roll
    │   ├── BUILD.bazel
    │   └── BUILD.bazel
    └── word
        └── BUILD.bazel
```

The WORKSPACE file was updated as well, and we have a another new file called "deps.bzl". 
We now have a working bazel project, so what commands can we run?

### Basic bazel commands

There are a various bazel [commands](https://bazel.build/run/build#available-commands) that 
are defined.

The main ones that are typically run by developers are [build](https://bazel.build/run/build#bazel-build),
[test](https://bazel.build/docs/user-manual#running-tests) and [run](https://bazel.build/docs/user-manual#running-executables).

The build and test commands are pretty self explanitory.  The build command builds the source code
for your project, and the test command runs any tests that are defined. The run command
execs a rule, for instance executes a go binary.

In the project you can run:

```
$ bazelisk build //...
```

This will build the binary for our example project. We can run the binary that bazel
creates with the following command:

```
$ bazelisk run //:go-code-tutorial
```
You can also pass in the command line option "word" that we defined to the bazel run command.

```
$ bazelisk run //:go-code-tutorial word
```

We will talk about the "test" command later. As we do not have any tests defined
in the project.

So the commands build, run and test are pretty easy to get your head around, but the third part of the
command was a bit confusing for me when I first learned bazel.  The "//..." or "//:something" is 
what is called a target.

You can refer to the documentation [here](https://bazel.build/run/build#bazel-build).  The text "//..."
and "//:go-code-tutorial" are all the targets in a given directory or is the name of a 
specific target.  Some commands like build and test can run multiple targets, 
while a command like run can only execute one target.

The below table provides a great guide for targets:

<table>
<tbody><tr>
  <td><code translate="no" dir="ltr">/<wbr>/<wbr>foo/<wbr>bar:wiz</code></td>
  <td>Just the single target <code translate="no" dir="ltr">/<wbr>/<wbr>foo/<wbr>bar:wiz</code>.</td>
</tr>
<tr>
  <td><code translate="no" dir="ltr">/<wbr>/<wbr>foo/<wbr>bar</code></td>
  <td>Equivalent to <code translate="no" dir="ltr">/<wbr>/<wbr>foo/<wbr>bar:bar</code>.</td>
</tr>
<tr>
  <td><code translate="no" dir="ltr">/<wbr>/<wbr>foo/<wbr>bar:all</code></td>
  <td>All rule targets in the package <code translate="no" dir="ltr">foo/<wbr>bar</code>.</td>
</tr>
<tr>
  <td><code translate="no" dir="ltr">/<wbr>/<wbr>foo/<wbr>.<wbr>.<wbr>.<wbr></code></td>
  <td>All rule targets in all packages beneath the directory <code translate="no" dir="ltr">foo</code>.</td>
</tr>
<tr>
  <td><code translate="no" dir="ltr">/<wbr>/<wbr>foo/<wbr>.<wbr>.<wbr>.<wbr>:all</code></td>
  <td>All rule targets in all packages beneath the directory <code translate="no" dir="ltr">foo</code>.</td>
</tr>
<tr>
  <td><code translate="no" dir="ltr">/<wbr>/<wbr>foo/<wbr>.<wbr>.<wbr>.<wbr>:&#42;</code></td>
  <td>All targets (rules and files) in all packages beneath the directory <code translate="no" dir="ltr">foo</code>.</td>
</tr>
<tr>
  <td><code translate="no" dir="ltr">/<wbr>/<wbr>foo/<wbr>.<wbr>.<wbr>.<wbr>:all-targets</code></td>
  <td>All targets (rules and files) in all packages beneath the directory <code translate="no" dir="ltr">foo</code>.</td>
</tr>
<tr>
  <td><code translate="no" dir="ltr">/<wbr>/<wbr>.<wbr>.<wbr>.<wbr></code></td>
  <td>All targets in packages in the workspace. This does not include targets
  from <a href="/docs/external">external repositories</a>.</td>
</tr>
<tr>
  <td><code translate="no" dir="ltr">/<wbr>/<wbr>:all</code></td>
  <td>All targets in the top-level package, if there is a `BUILD` file at the
  root of the workspace.</td>
</tr>
</tbody></table>

> <cite>https://bazel.build/run/build#specifying-build-targets</cite>

If we look in the BUILD.bazel file in the root directory will will find a go_libary rule
named go-code-tutorial_lib, and this is a target we can build.

```
$ bazelisk build //:go-code-tutorial_lib
```

This "go_libary" target is named by gazelle automatically depending on the name of your project, so
the name may differ.

We can also run the go-code-tutorial binary target using the following command:

```
$ bazelisk run //:go-code-tutorial word
```

Or we can build all of the targets under the pkg directory:

```
$ bazelisk build //pkg/...
```

#### Note about binaries and build

We wanted to include a side note about "bazel build".  You may wonder where the heck is the binary put?
Bazel creates various folders and symlinks in project directory. Within out example we have

- bazel-bazel-gazelle
- bazel-bin
- bazel-out
- bazel-go-code-tutorial
- bazel-testlogs

Binaries from the project are placed under the bazel-bin folder.  Inside of that folder we have another folder
that has the name go-code-tutorial\_ and that folder name is created from the name of the binary that is 
created.  A bazel project can contain multiple binaries, so we have to have that form of naming syntax.  Inside
of the go-code-tutorial\_ folder we have the binary go-code-tutorial\_.

### Where gazelle defines the dependencies

One of the features of gazelle is to "vendor" Go projects.  Within this example we are 
using Go vendoring at the base, but bazel must also have the external dependencies defined.

The gazelle update-repos command takes the go.mod file and creates the StarkLark code that
defines the external vendoring that bazel uses. External dependencies are defined in one 
of two locations; in the WORKSPACE file or in an external file that is references in
the WORKSPACE file. The list of external dependencies can grow very long, so we recommend that
it is defined as a refernce in the WORKSPACE file.

Each of the following lines within the WORKSPACE file defines the location of the deps.bz file:

```
# load bazel and gazelle rules
load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

############################################################
# Define your own dependencies here using go_repository.
# Else, dependencies declared by rules_go/gazelle will be used.
# The first declaration of an external repository "wins".
############################################################

load("//:deps.bzl", "go_dependencies")
```
One challenge you can run into is that you need to manually override a dependency, and  you can
do this listing the http_archive. Below we have an example of overriding the buildtools dependency.

```
http_archive(
    name = "com_github_bazelbuild_buildtools",
    sha256 = "a02ba93b96a8151b5d8d3466580f6c1f7e77212c4eb181cba53eb2cae7752a23",
    strip_prefix = "buildtools-3.5.0",
    urls = [
        "https://github.com/bazelbuild/buildtools/archive/3.5.0.tar.gz",
    ],
)
```

This example is from the cockroach database operator project. You can see
the full definition [here](https://github.com/cockroachdb/cockroach-operator/blob/0ef4d1e1b4c94a8edf1393b0fa72d9de8bc21477/WORKSPACE#L20)

Now lets cover what is inside of the BUILD files. As we mentioned bazel rules are in essence 
StarLark libaries.

### The BUILD files

The rules_go have several "Core rules" defined.  These include:

- go_binary
- go_library
- go_test
- go_source
- go_path

See [here](https://github.com/bazelbuild/rules_go/blob/master/docs/go/core/rules.md) for more details.
And these StarLark rules are used inside of the BUILD files, and often updated automatically by gazelle.

After we ran gazelle the BUILD.bazel file was updated to include two new StarLark definitions:

```
go_library(
    name = "go-code-tutorial_lib",
    srcs = ["main.go"],
    importpath = "github.com/bazel-contrib/Bazel-learning-paths/tutorials/go-code-tutorial",
    visibility = ["//visibility:private"],
    deps = ["//cmd"],
)

go_binary(
    name = "go-code-tutorial",
    embed = [":go-code-tutorial_lib"],
    visibility = ["//visibility:public"],
)
```

Both the go_libary and go_binary rules are defined for our code. The go_libary rule defines the build of a Go library from a set of source files that are all part of the same package. The go_binary rule defines the build of an executable from a set of source files, which must all be in the main package.  The go_rules project includes are great documentation [section](https://github.com/bazelbuild/rules_go/blob/master/docs/go/core/rules.md#introduction) if you want more details.

More BUILD.bazel files where also created. Here is the BUILD.bazel file that was created in 
the cmd folder.

```
load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "cmd",
    srcs = [
        "roll.go",
        "root.go",
        "word.go",
    ],
    importpath = "github.com/bazel-contrib/Bazel-learning-paths/tutorials/go-code-tutorial/cmd",
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_spf13_cobra//:cobra",
    ],
)
```

The first line load the SkyLark definition from the go_rules library. You can then
use "go_libary" which is used directly after.  This go_libary definition also mentions
an external dependency using cobra.

### How these files work together

The WORKSPACE, dep.bzl, and BUILD.bazel files create an object graphs that bazel uses.
The blog [post](https://blog.bazel.build/2015/06/17/visualize-your-build.html) covers
visualizing the object graph.  Take a peak if you want to learn a bit about
"bazel query".

Next we cover more definitions in the WORKSPACE file.  We can start with the following
code:

```
http_archive(
    name = "io_bazel_rules_go",
```

We are not including the full call for the sake of brevity. This http_archive definition tells
bazel to download and use a specific version of rules_go. If you look at the BUILD.bazel file in the
root directory you can see load command for rules_go, which exports go_libary.

```
load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
```

The go_libary definition is then used later in the file.

```
go_library(
    name = "go-code-tutorial_lib",
```

So the WORKSPACE file includes the definition of which rules_go we are using and then the BUILD.bazel
files loads those rules and uses one of the definitions in the rules. 


// TODO a couple of images here for object graphs

The same kind of object graph is used for external dependencies. The WORKSPACE file include the
definition for gazelle (http_archive) and includes an import for the deps.bzl file. The deps.bzl file
includes load definitions for the gazelle "go_repository" rule.   The go_repository rules define various
external go dependencies that are then vendored.  One of those dependencies is cobra, and cobra is used
as a dependecy by all of the go files inside of the cmd directory. Inside of the BUILD.bazel file in the cmd 
directory the a "deps" are a parameter passed in the go_libary rule.

```
    deps = ["@com_github_spf13_cobra//:go_default_library"],
```

So now we have the capabilty for bazel to:

- Build an object graph for the project
- Various rules are defined that impact the object tree
- go_rules and gazelle define various rules
- the bazel object tree includes go_libary rules
- external depencies are defined in go_repository rules
- deps are passed into go_libary rules

All of these definitions create an dependency graph that allows bazel to run:

```
$ bazelisk build //...

```

We the command is executed, bazel will download and cache all dependencies including but not limited to:

- The defined GoLang compiler and libaries
- The defined rules sets
- build the go binary that is defined the in root of the project.

Downloading and caching the above components is part of bazel providing hermetic and
deterministic builds.  All of the downloaded components are checked against a SHA that
verifies the checksum of the downloaded file(s).

Next we will make some code changes and introduce some internal code
dependencies.

## Using the files under pkg

Now we want to modify and use the files under the pkg directory.

Edit roll.go under the cmd folder and add an import to roll_dice.

You will now have:

```
import (
    "fmt"

    "github.com/bazel-contrib/Bazel-learning-paths/tutorials/go-code-tutorial/pkg/roll"
    "github.com/spf13/cobra"
)
```

Then call `roll.Roll()` after the `fmt.Println` statement. This will give you:

```
   Run: func(cmd *cobra.Command, args []string) {
       fmt.Println("roll called")
       roll.Roll()
   },
```

You have edited the following files.

```
├── cmd
│   ├── roll.go
└── pkg
    └── roll
        └── roll_dice.go
```

We now need to update the BUILD.bazel files, and the easiest way to do this is to run gazelle again.

Execute the following command:

```
$ bazelisk run //:gazelle
```
We can now use bazel to run the binary again:

```
$ bazelisk run //:go-code-tutorial roll

```

The above commands build the Go binary and executes it.  The following
is an example of the output from the run command.

```
INFO: Analyzed target //:go-code-tutorial (1 packages loaded, 6 targets configured).
INFO: Found 1 target...
Target //:go-code-tutorial up-to-date:
  bazel-bin/go-code-tutorial\_/go-code-tutorial
INFO: Elapsed time: 0.316s, Critical Path: 0.16s
INFO: 3 processes: 1 internal, 2 linux-sandbox.
INFO: Build completed successfully, 3 total actions
INFO: Build completed successfully, 3 total actions
roll called
roll dice
```

Running the gazelle target modified the Build.bazel file under the cmd directory.  Here is the diff.

```
diff --git a/cmd/BUILD.bazel b/cmd/BUILD.bazel
index ac66183..9033b86 100644
--- a/cmd/BUILD.bazel
+++ b/cmd/BUILD.bazel
@@ -9,5 +9,8 @@ go_library(
     ],
     importpath = "github.com/bazel-contrib/Bazel-learning-paths/tutorials/go-code-tutorial/cmd",
     visibility = ["//visibility:public"],
-    deps = ["@com_github_spf13_cobra//:cobra"],
+    deps = [
+        "//pkg/roll",
+        "@com_github_spf13_cobra//:cobra",
+    ],
 )
```

The line was added inside of the deps stanza that points to the package where roll.go resides.

Next modify cmd/word.go file to include a call to  `word.GenerateWord()` func.

Here is the diff of the code changes:

```
diff --git a/cmd/word.go b/cmd/word.go
index d7d00bb..cddc748 100644
--- a/cmd/word.go
+++ b/cmd/word.go
@@ -1,12 +1,12 @@
 /*
 Copyright © 2022 NAME HERE <EMAIL ADDRESS>
-
 */
 package cmd

 import (
        "fmt"

+       "github.com/bazel-contrib/Bazel-learning-paths/tutorials/go-code-tutorial/pkg/word"
        "github.com/spf13/cobra"
 )

@@ -22,6 +22,7 @@ This application is a tool to generate the needed files
 to quickly create a Cobra application.`,
        Run: func(cmd *cobra.Command, args []string) {
                fmt.Println("word called")
+               word.GenerateWord()
        },
 }
```

We added the import and the call to `word.GenerateWord()`. Again we can run gazelle 
add the new dep to the BUILD.bazel file. 

```
$ bazelisk run //:gazelle
```

Now we have BUILD.bazel updated. Here is the diff:

```
diff --git a/cmd/BUILD.bazel b/cmd/BUILD.bazel
index ac66183..891b0e1 100644
--- a/cmd/BUILD.bazel
+++ b/cmd/BUILD.bazel
@@ -9,5 +9,9 @@ go_library(
     ],
     importpath = "github.com/bazel-contrib/Bazel-learning-paths/tutorials/go-code-tutorial/cmd",
     visibility = ["//visibility:public"],
-    deps = ["@com_github_spf13_cobra//:cobra"],
+    deps = [
+        "//pkg/roll",
+        "//pkg/word",
+        "@com_github_spf13_cobra//:cobra",
+    ],
 )
```

We can use bazel to execute the binary with the new changes.

```
$ bazelisk run //:go-code-tutorial word
```

The above command genertates the following output.

```
INFO: Analyzed target //:go-code-tutorial (0 packages loaded, 0 targets configured).
INFO: Found 1 target...
Target //:go-code-tutorial up-to-date:
  bazel-bin/go-code-tutorial\_/go-code-tutorial
INFO: Elapsed time: 0.107s, Critical Path: 0.00s
INFO: 1 process: 1 internal.
INFO: Build completed successfully, 1 total action
INFO: Build completed successfully, 1 total action
word called
GenerateWord
```

The project is now modified so that the files under the pkg folder are now used.  This is the 
principle of using internal dependencies.  Next we will add a Go project dependency
is hosted out of GitHub, an "external dependency".

## Adding an external dependency

To create our random work generator we are going to use babble, which is located here: 
https://github.com/tjarratt/babble. The babble code On Linux uses "/usr/share/dicts/words" file, and you can use 
the package manager to install wamerican or wbritish. See the babble README for more information
on other operating systems.

Edit generate_word.go to add the call to babble. This is the file
we are refering to:

```
└── pkg
    └── word
        └── generate_word.go
```

We need to add the import to the babble libary and call the babble func. Here is the diff after the updates.

```
diff --git a/pkg/word/generate_word.go b/pkg/word/generate_word.go
index 312a267..37215cf 100644
--- a/pkg/word/generate_word.go
+++ b/pkg/word/generate_word.go
@@ -1,7 +1,12 @@
 package word

-import "fmt"
+import (
+       "fmt"
+
+       "github.com/tjarratt/babble"
+)

 func GenerateWord() {
+       fmt.Println("GenerateWord called")
+       fmt.Println(babble.NewBabbler().Babble())
 }
```

I also cleaned up the Println to add some clarity.

Once that code change is done, we need to run go mod to update the projects 
dependencies. We can use bazel to run the go binary, instead of having
to install go and running that binary.  The Go rules have already downloaded
the Go SDK, so use the following command.

```
$ bazelisk run @go_sdk//:bin/go -- mod tidy
```

Keeping go.mod updated allows us to either use go directly or bazel to build
and run the code.

We now need to update the Bazel import, and the easiest way to do this is to run gazelle again.

```
$ bazelisk run //:gazelle-update-repos
$ bazelisk run //:gazelle
```

The first bazel command updates deps.bzl file. The second command
updates the BUILD.bazel file in pkg/word.  Below is the diff of the 
updates.

```
diff --git a/pkg/word/BUILD.bazel b/pkg/word/BUILD.bazel
index c974b0b..e5c0b28 100644
--- a/pkg/word/BUILD.bazel
+++ b/pkg/word/BUILD.bazel
@@ -5,4 +5,5 @@ go_library(
     srcs = ["generate_word.go"],
     importpath = "github.com/bazel-contrib/Bazel-learning-paths/tutorials/go-code-tutorial/pkg/word",
     visibility = ["//visibility:public"],
+    deps = ["@com_github_tjarratt_babble//:babble"],
 )

```

You can see the deps is now updated and points to the external repo "@com_github_tjarratt_babble//:babble".
The "@" references an external code base that bazel will download, so that the Go SDK can build
the code.

This GitHub repo is defined in deps.bzl file in the following go_repository stanza.

```
go_repository(
    name = "com_github_tjarratt_babble",
    importpath = "github.com/tjarratt/babble",
    sum = "h1:j8whCiEmvLCXI3scVn+YnklCU8mwJ9ZJ4/DGAKqQbRE=",
    version = "v0.0.0-20210505082055-cbca2a4833c1",
)
```

We can now run our Go binary and see the changes.

```
$ bazelisk run //:go-code-tutorial word
INFO: Analyzed target //:go-code-tutorial (0 packages loaded, 0 targets configured).
INFO: Found 1 target...
Target //:go-code-tutorial up-to-date:
  bazel-bin/go-code-tutorial_/go-code-tutorial
INFO: Elapsed time: 0.257s, Critical Path: 0.15s
INFO: 3 processes: 1 internal, 2 linux-sandbox.
INFO: Build completed successfully, 3 total actions
INFO: Build completed successfully, 3 total actions
word called
GenerateWord called
Rheingau-nightclothes
```

One of the things that you may notice is that you do not have to run "bazel build" and then "bazel run".
Bazel will notice that the code is not built, and will run the "build" phase for you automaticallly.

To recap what we have done.  We have modified our code to use the babble Go code which lives on 
GitHub.  We then use bazel to run go mod, which updates go.mod file. Next we ran the targets gazelle-update-repos and gazelle
with bazel. The first bazel alias updated the deps.bzl file with the external dependency, and the gazelle target 
updated the deps section in pkg/word/BUILD.bazel.  Bazel is then able to download the external dependency
and use that dependency when our example Go program is compiled.

How about we add a Go unit test so that we can run "bazel test"?

## Go tests

As we mentioned bazel support running code tests, as defined in bazel rules. One of the rules from go_rules
is go_test.  Now lets add a test.

First refactor the func GenerateWord in the Go file pkg/word/generate.go to return a string, rather than printing it.

Here are the changes:

```
func GenerateWord() string {
    fmt.Println("GenerateWord called")
    return babble.NewBabbler().Babble()
}
```

This will allow use to test that this func does not return a string that is empty.

Now move the Println higher up in the stack so that the random word is still printed.
Edit cmd/word.go file, and add a Println around the call to GenerateWord().

Here are these changes.

```
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("word called")
        fmt.Println(word.GenerateWord())
    },
```

Now let us build the code to make sure that we did not make any mistakes.

```
$ bazelisk build //...
```

Create a new file in the pkg/word directory called generate_word_test.go.
Include the following code:

```
package word

import (
    "testing"
)

func TestGenerateWord(t *testing.T) {
    result := GenerateWord()
    if result == "" {
        t.Error("got an empty string")
    }
}
```

We have a unit test now, but bazel does not know about it.  Again we need 
bazel to have the target in its object graph, and in order to do that we need
to update the BUILD.bazel file.  The easiest way to do that is with gazelle.

Simply run:

```
$ bazelisk run //:gazelle
```

This now updates the the BUILD.bazel file in the pkg/word directory.  Here
is a diff of the update:

```
diff --git a/pkg/word/BUILD.bazel b/pkg/word/BUILD.bazel
index e5c0b28..1b79ce0 100644
--- a/pkg/word/BUILD.bazel
+++ b/pkg/word/BUILD.bazel
@@ -1,4 +1,4 @@
-load("@io_bazel_rules_go//go:def.bzl", "go_library")
+load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

 go_library(
     name = "word",
@@ -7,3 +7,9 @@ go_library(
     visibility = ["//visibility:public"],
     deps = ["@com_github_tjarratt_babble//:babble"],
 )
+
+go_test(
+    name = "word_test",
+    srcs = ["generate_word_test.go"],
+    embed = [":word"],
+)
```

We now have a [go_test](https://github.com/bazelbuild/rules_go/blob/master/docs/go/core/rules.md#go_test) 
rule, which is part of the rules_go ruleset. Now we can run:

```
$ bazelisk test //...
```

The above command should print out results like:

```
$ bazelisk test //...
INFO: Analyzed 6 targets (0 packages loaded, 0 targets configured).
INFO: Found 5 targets and 1 test target...
INFO: Elapsed time: 0.125s, Critical Path: 0.00s
INFO: 1 process: 1 internal.
INFO: Build completed successfully, 1 total action
//pkg/word:word_test                                            (cached) PASSED in 0.0s

Executed 0 out of 1 test: 1 test passes.
INFO: Build completed successfully, 1 total action
```

You may also notice that the command printed out a target named //pkg/word:word_test.
We can also run just the specific target:

```
$ bazelisk test //pkg/word:word_test
```

Lets now see what happens when a test fails, since debugging unit tests are often part of the
development process. In the generate_word_test.go file change the "if" statement as show below.

```
    if result == "" {
```

Now if we run

```
$ bazelisk test //pkg/word:word_test
```

We get an output like

```
$ bazelisk test //...
INFO: Analyzed 6 targets (0 packages loaded, 0 targets configured).
INFO: Found 5 targets and 1 test target...
FAIL: //pkg/word:word_test (see /home/clove/.cache/bazel/_bazel_clove/d1fd07b841c26eda93328e4eeaf2336a/execroot/__main__/bazel-out/k8-fastbuild/testlogs/pkg/word/word_test/test.log)
INFO: Elapsed time: 0.299s, Critical Path: 0.17s
INFO: 6 processes: 1 internal, 5 linux-sandbox.
INFO: Build completed, 1 test FAILED, 6 total actions
//pkg/word:word_test                                                     FAILED in 0.0s
  /home/clove/.cache/bazel/_bazel_clove/d1fd07b841c26eda93328e4eeaf2336a/execroot/__main__/bazel-out/k8-fastbuild/testlogs/pkg/word/word_test/test.log

INFO: Build completed, 1 test FAILED, 6 total actions
```

The line that diplays the path to the test.log file will differ between systems, but it provides output from the unit test.
If we cat the file we see the results:

```
$ cat /home/clove/.cache/bazel/_bazel_clove/d1fd07b841c26eda93328e4eeaf2336a/execroot/__main__/bazel-out/k8-fastbuild/testlogs/pkg/word/word_test/test.log
exec ${PAGER:-/usr/bin/less} "$0" || exit 1
Executing tests from //pkg/word:word_test
-----------------------------------------------------------------------------
GenerateWord called
HERE
grebes-slickness's
--- FAIL: TestGenerateWord (0.00s)
    generate_word_test.go:13: got an empty string
FAIL
```

Adding the "test_ouput" argument to the bazel test command will output the test results to the console.

```
$ bazelisk test --test_output=errors //...
INFO: Analyzed 6 targets (0 packages loaded, 0 targets configured).
INFO: Found 5 targets and 1 test target...
FAIL: //pkg/word:word_test (see /home/clove/.cache/bazel/_bazel_clove/d1fd07b841c26eda93328e4eeaf2336a/execroot/__main__/bazel-out/k8-fastbuild/testlogs/pkg/word/word_test/test.log)
INFO: From Testing //pkg/word:word_test:
==================== Test output for //pkg/word:word_test:
GenerateWord called
HERE
justest-indefensibly
--- FAIL: TestGenerateWord (0.00s)
    generate_word_test.go:13: got an empty string
FAIL
================================================================================
INFO: Elapsed time: 0.191s, Critical Path: 0.02s
INFO: 2 processes: 1 internal, 1 linux-sandbox.
INFO: Build completed, 1 test FAILED, 2 total actions
//pkg/word:word_test                                                     FAILED in 0.0s
  /home/clove/.cache/bazel/_bazel_clove/d1fd07b841c26eda93328e4eeaf2336a/execroot/__main__/bazel-out/k8-fastbuild/testlogs/pkg/word/word_test/test.log

INFO: Build completed, 1 test FAILED, 2 total actions
```

If you like you can change the "if" statement back so that the unit test passes.

So now we know how to include a new unit test, update BUILD.bazel rules with gazelle, and then run the test.

## Other rules in rules_go

The rules_go [documentation](https://github.com/bazelbuild/rules_go#documentation) provides a great reference to the different
rules provided in the ruleset.

We have covered three of the top rules 'go_binary', 'go_library', and 'go_test'.  We also covered a rules that
gazelle uses to manage dependencies called 'go_repository'.

Other go_rules rules include:

- Proto rules that generate Go packages from .proto files. These packages can be imported like regular Go libraries.
- The Go toolchain is a set of rules used to customize the behavior of the core Go rules.  The Go toolchain allows for the configuration
of the Go distribution utilised. The toolchain declare Bazel toolchains for each target platform that Go supports. The context rules all for the writing custom rules
that are compatible with rules_go.
- Also go_rules includes rule for using go mock, and the rule go_embed_data.
The rule go_embed_data generates a .go file that contains data from a file or a list of files. 
- The nogo rule support using nogo during testing. The code analysis tool nogo screens code preventing bugs and code anti-patterns, and can also run vet.

Other capabilities of go_rules include:

- creating pure go binaries
- building go static binaries
- basic race condition detection

And lastly you probably know that Go supports cross-compilation, and this is really nice when we are developing with containers.  Within rules_go they 
have included go_cross_binary, which allows your to define the creation of a binary for a specific operating system and CPU architecture. This
can allow us to develop on a Mac and run the binary on that Mac, while also building a binary for Linux.  We then would use a set of bazel
rules that support the building of containers, and bazel can put the Linux binary in the container.

## Summary

- Bazel supports the building and testing of the Go programming language using the rules_go ruleset.
- Intially you need to create a basic WORKSPACE and BUILD.bazel file in the root directory of your project.
- You can use gazelle to create and maintain various bazel files.
- Gazelle can update various bazel files when you add new go file or go tests.
- Bazel supports many commands, and we covered the build, run and test commands.
- Bazel uses an object graph that is based on WORKSPACE, BUILD.bazel and other bazel files.
- The ruleset rules_go provides various rules like go_binary, go_libary and go_test.  They are used
to build binaries, libaries and supporting unit testing.
- Gazelle can update BUILD.bazel and dep.bzl files with either internal or external Go dependencies.
- The go_test rule is used to defined Go unit tests.
- Various other rules are defined by go_rules.  These rules include managing protocol buffers, grpc, cross compilation, and controlling various
aspects of how the Go SDK is downloaded and configured.
