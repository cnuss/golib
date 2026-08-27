// The tests for lib.go, internal and external alike.
//
// This file is `package golib`, not `package golib_test`, because the Version
// resolution tests reach the unexported version stamp and vcsVersion. The
// façade tests below would read slightly better from outside the package, but
// splitting them across two files to gain that would mint a second test file
// for one source file — see CONTRIBUTING's test-file naming rule. The
// genuine outside-the-package view is covered by v1/example_test.go (package
// v1_test, importing the root) and by e2e, which builds and runs the example
// binaries as real consumers.
package golib

import (
	"runtime/debug"
	"testing"

	v1 "github.com/cnuss/golib/v1"
)

// TestAliasesAreAliases pins that the root re-exports are type aliases (`=`)
// rather than defined types, so the two spellings of each type are one type
// and values cross freely between code that names them through the root and
// code that names them through v1.
//
// The Result assignments are the load-bearing half: structs are assignable
// only when at most one side is a defined type, so declaring Result without
// the `=` fails to compile right here. The Builder assignments can't catch the
// same slip — interface assignability is structural, so a defined interface
// type with the same method set stays assignable — and are here to document
// the intended usage.
func TestAliasesAreAliases(t *testing.T) {
	var fromRoot BuilderV1[string] = New[string]()
	var fromV1 v1.Builder[string] = fromRoot
	fromRoot = fromV1

	var resultRoot Result[string] = fromRoot.WithValue("hello").Build()
	var resultV1 v1.Result[string] = resultRoot
	resultRoot = resultV1

	if resultRoot.Value != "hello" {
		t.Errorf("Value = %q, want %q", resultRoot.Value, "hello")
	}
}

// TestFacadeCoversTheCommonPath walks the whole façade path a caller uses —
// build a value, name the types it holds, ask which release is linked — so a
// break anywhere along it fails here rather than only in an example binary.
func TestFacadeCoversTheCommonPath(t *testing.T) {
	var b BuilderV1[int] = New[int]()

	res := b.WithName("count").WithValue(7).Build()
	if res.Name != "count" || res.Value != 7 {
		t.Errorf("Build() = %+v, want {Name:count Value:7}", res)
	}

	if Version() == "" {
		t.Error("Version() = empty, want a derived identifier")
	}
}

// TestVersionStampWins pins the release-stamp precedence: when the build
// stamps the version (via -ldflags -X), Version returns it verbatim, ahead of
// any build-info resolution.
func TestVersionStampWins(t *testing.T) {
	old := version
	t.Cleanup(func() { version = old })

	version = "v9.9.9"
	if got := Version(); got != "v9.9.9" {
		t.Errorf("Version() = %q, want the stamp v9.9.9", got)
	}
}

// TestVersionFallbackNonEmpty pins that Version always self-identifies: with
// no stamp it derives an identifier from build info (module version, main
// version, or VCS revision) and never returns "".
func TestVersionFallbackNonEmpty(t *testing.T) {
	old := version
	t.Cleanup(func() { version = old })

	version = ""
	if got := Version(); got == "" {
		t.Error("Version() = empty with no stamp, want a derived identifier")
	}
}

// TestVCSVersion covers the local-build fallback against synthetic build info,
// which is the only way to exercise every branch — a real `go test` binary
// carries whichever stamp the toolchain chose to embed.
func TestVCSVersion(t *testing.T) {
	const long = "0123456789abcdef0123456789abcdef01234567"

	cases := []struct {
		name     string
		settings []debug.BuildSetting
		want     string
	}{
		{
			name: "revision truncated to 12",
			settings: []debug.BuildSetting{
				{Key: "vcs.revision", Value: long},
			},
			want: "0123456789ab",
		},
		{
			name: "dirty tree suffixed",
			settings: []debug.BuildSetting{
				{Key: "vcs.revision", Value: long},
				{Key: "vcs.modified", Value: "true"},
			},
			want: "0123456789ab-dirty",
		},
		{
			name: "clean tree unsuffixed",
			settings: []debug.BuildSetting{
				{Key: "vcs.revision", Value: "abc123"},
				{Key: "vcs.modified", Value: "false"},
			},
			want: "abc123",
		},
		{
			name:     "no vcs stamp at all",
			settings: nil,
			want:     "unknown",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := vcsVersion(&debug.BuildInfo{Settings: tc.settings})
			if got != tc.want {
				t.Errorf("vcsVersion() = %q, want %q", got, tc.want)
			}
		})
	}
}
