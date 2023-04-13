package main

import (
	"fmt"

	"go.starlark.net/starlark"

	pb "github.com/bazel-contrib/SIG-rules-authors/experimental/rules-keeper/proto"
)

// parseModuleFile parses a module file.
//
// The filename and src parameters are as for syntax.Parse: filename is the name
// of the file to execute, and the name that appears in error messages; src is
// an optional source of bytes to use instead of filename.
func parseModuleFile(filename string, src any) (*pb.ModuleFile, error) {
	var mf pb.ModuleFile
	predeclared := starlark.StringDict{
		"module": starlark.NewBuiltin("module", func(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			var (
				ignored                                   starlark.Value
				bazelCompatibility, platforms, toolchains *starlark.List
			)
			if err := starlark.UnpackArgs(b.Name(), args, kwargs,
				"name?", &mf.Name,
				"repo_name", &ignored,
				"version?", &mf.Version,
				"compatibility_level?", &mf.CompatibilityLevel,
				"bazel_compatibility?", &bazelCompatibility,
				"execution_platforms_to_register?", &platforms,
				"toolchains_to_register?", &toolchains,
			); err != nil {
				return nil, err
			}
			if bazelCompatibility != nil {
				mf.BazelCompatibility = appendList(mf.BazelCompatibility, bazelCompatibility)
			}
			if platforms != nil {
				mf.ExecutionPlatformsToRegister = appendList(mf.ExecutionPlatformsToRegister, platforms)
			}
			if toolchains != nil {
				mf.ToolchainsToRegister = appendList(mf.ToolchainsToRegister, toolchains)
			}
			return starlark.None, nil
		}),
		"bazel_dep": starlark.NewBuiltin("bazel_dep", func(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			var (
				ignored starlark.Value
				dep     pb.ModuleFile_Dependency
			)
			if err := starlark.UnpackArgs(b.Name(), args, kwargs,
				"name", &dep.Name,
				"version?", &dep.Version,
				"repo_name?", &ignored,
				"dev_dependency?", &dep.DevDependency,
			); err != nil {
				return nil, err
			}
			mf.Dependencies = append(mf.Dependencies, &dep)
			return starlark.None, nil
		}),
		"register_execution_platforms": starlark.NewBuiltin("register_execution_platforms", func(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			mf.ExecutionPlatformsToRegister = appendList(mf.ExecutionPlatformsToRegister, starlark.NewList(args))
			return starlark.None, nil
		}),
		"register_toolchains": starlark.NewBuiltin("register_toolchains", func(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			mf.ToolchainsToRegister = appendList(mf.ToolchainsToRegister, starlark.NewList(args))
			return starlark.None, nil
		}),
		"use_extension": starlark.NewBuiltin("use_extension", func(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			return noopExtension{}, nil
		}),
		"archive_override":          starlark.NewBuiltin("archive_override", noop),
		"git_override":              starlark.NewBuiltin("git_override", noop),
		"local_path_override":       starlark.NewBuiltin("local_path_override", noop),
		"multiple_version_override": starlark.NewBuiltin("multiple_version_override", noop),
		"single_version_override":   starlark.NewBuiltin("single_version_override", noop),
		"use_repo":                  starlark.NewBuiltin("use_repo", noop),
	}

	thread := &starlark.Thread{
		Name:  "module",
		Print: func(_ *starlark.Thread, msg string) {},
	}
	_, err := starlark.ExecFile(thread, filename, src, predeclared)
	if err != nil {
		return nil, err
	}
	return &mf, nil
}

func appendList(s []string, l *starlark.List) []string {
	for it := l.Iterate(); ; {
		var v starlark.Value
		if !it.Next(&v) {
			break
		}
		t, ok := starlark.AsString(v)
		if ok {
			s = append(s, t)
		}
	}
	return s
}

func noop(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return starlark.None, nil
}

type noopExtension struct{}

var _ starlark.HasAttrs = noopExtension{}

func (noopExtension) Attr(name string) (starlark.Value, error) {
	return starlark.NewBuiltin(name, noop), nil
}
func (noopExtension) AttrNames() []string   { return nil }
func (noopExtension) String() string        { return "{}" }
func (noopExtension) Freeze()               {}
func (noopExtension) Type() string          { return "proxy" }
func (noopExtension) Truth() starlark.Bool  { return starlark.True }
func (noopExtension) Hash() (uint32, error) { return 0, fmt.Errorf("unhashable type: proxy") }
