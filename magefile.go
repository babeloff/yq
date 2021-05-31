// +build mage

// Usage:
//	Develop / Test Commands
package main

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	//"github.com/phreed/yq/codegen"
	//"github.com/phreed/yq/resources/page/page_generate"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
)

const (
	packageName  = "github.com/babeloff/yq"
	noGitLdflags = "-X $PACKAGE/common/yq.buildDate=$BUILD_DATE"

	PROJECT     = "yq"
	IMPORT_PATH = "github.com/mikefarah/" + PROJECT

	GITHUB_TOKEN = ""
	DEV_IMAGE    = PROJECT + "_dev"
)

var ROOT, root_err = os.Getwd()

var GIT_COMMIT, git_commit_err = sh.Output("git", "rev-parse", "--short", "HEAD")
var GIT_DIRTY, git_dirty_err = sh.Output("git", "status", "--porcelain")
var GIT_DESCRIBE, git_describe_err = sh.Output("git", "describe", "--tags", "--always")
var LDFLAGS = "-X main.GitCommit=" + GIT_COMMIT + GIT_DIRTY + " " +
	"-X main.GitDescribe=" + GIT_DESCRIBE

var dockerRun = sh.RunCmd("docker",
	"run", "--rm",
	"-e", "LDFLAGS=\""+LDFLAGS+"\"",
	"-e", "GITHUB_TOKEN=\""+GITHUB_TOKEN+"\"",
	"-v", ROOT+"/vendor:/go/src",
	"-v", ROOT+":/"+PROJECT+"/src/"+IMPORT_PATH,
	"-w", "/"+PROJECT+"/src/"+IMPORT_PATH,
	DEV_IMAGE)

// allow user to override go executable by running as GOEXE=xxx make ... on unix-like systems
var goexe = "go"

// Clean project
func Clean() error {
	rm_bin_err := os.RemoveAll("bin")
	if rm_bin_err != nil {
		return rm_bin_err
	}
	rm_build_err := os.RemoveAll("build")
	if rm_build_err != nil {
		return rm_build_err
	}
	rm_cover_err := os.RemoveAll("cover")
	if rm_cover_err != nil {
		return rm_cover_err
	}
	rm_out_err := os.RemoveAll("*.out")
	if rm_out_err != nil {
		return rm_out_err
	}

	clean_err := sh.Run(goexe, "clean", "-i", packageName)
	if clean_err != nil {
		return clean_err
	}

	return nil
}

func Local() error {
	dockerRun()
	os.Mkdir("tmp", 0755)
	ioutil.WriteFile("tmp/dev_image-id", []byte(""), 0644)
	return nil
}

func Prepare() error {
	delta, _ := target.Dir("tmp/dev_image_id",
		"Dockerfile.dev", "scripts/devtools.sh")
	if !delta {
		return nil
	}

	os.Mkdir("tmp", 0755)
	sh.Run("docker", "rmi", "-f", DEV_IMAGE)
	sh.Run("docker", "build",
		"-t", DEV_IMAGE,
		"-f", "Dockerfile.dev",
		".")
	docker_id, _ := sh.Output("docker", "inspect",
		"-t", "{{ .ID }}",
		DEV_IMAGE,
		"-f", "Dockerfile.dev")
	ioutil.WriteFile("tmp/dev_image-id", []byte(docker_id), 0644)
	return nil
}

// Build yq binary.
func Build() error {
	delta, _ := target.Dir("build/dev",
		"test", "*.go")
	if !delta {
		return nil
	}

	mk_err := os.Mkdir("bin", 0755)
	if mk_err != nil {
		return mk_err
	}

	dockerRun("go", "build", "--ldflags", LDFLAGS)
	dockerRun("bash", "./scripts/acceptance.sh")
	return nil
}

// Build cross-compiled binaries of yq.
func Xcompile() error {
	mg.Deps(Check)

	os.RemoveAll("build")
	os.Mkdir("build", 0755)
	dockerRun("bash", "./scripts/xcompile.sh")
	filepath.Walk("build",
		func(path string, info fs.FileInfo, err error) error {
			return os.Chmod(path, 0755)
		})
	return nil
}

// Install yq.
func Install() error {
	mg.Deps(Build)
	dockerRun("go", "install")
	return nil
}

// Install dependencies to vendor directory.
func Vendor() error {
	mg.Deps(Prepare)
	dockerRun("go", "mod", "vendor")
	return nil
}

// Run code formatter.
func Format() error {
	mg.Deps(Vendor)
	dockerRun("bash", "./scripts/format.sh")
	return nil
}

// Run static code analysis (lint).
func Check() error {
	mg.Deps(Format)
	dockerRun("bash", "./scripts/check.sh")
	return nil
}

// Run gosec.
func Secure() error {
	dockerRun("bash", "./scripts/secure.sh")
	return nil
}

// Run tests on project.
func Test() error {
	dockerRun("bash", "./scripts/test.sh")
	return nil
}

// Run tests and capture code coverage metrics on project.
func Cover() error {
	mg.Deps(Check)

	os.RemoveAll("cover")
	os.Mkdir("cover", 0755)
	dockerRun("bash", "./scripts/coverage.sh")
	filepath.Walk("cover",
		func(path string, info fs.FileInfo, err error) error {
			return os.Chmod(path, 0755)
		})
	return nil
}

// Clean the directory tree of produced artifacts.
func Release() error {
	mg.Deps(Xcompile)
	dockerRun("bash", "./scripts/publish.sh")
	return nil
}

// utility: Configures Minishfit/Docker directory mounts.
func Setup() error {
	dockerRun("bash", "./scripts/setup.sh")
	return nil
}
