// +build mage

// Usage:
//	Develop / Test Commands
package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
)

const (
	packageName  = "github.com/babeloff/yq"
	PROJECT      = "yq"
	GITHUB_TOKEN = ""
	SPONSOR      = "babeloff"
)

var ROOT, root_err = os.Getwd()
var DEV_IMAGE = fmt.Sprintf("%s_dev", PROJECT)
var IMPORT_PATH = fmt.Sprintf("github.com/%s/%s", SPONSOR, PROJECT)

var GIT_COMMIT, git_commit_err = sh.Output("git", "rev-parse", "--short", "HEAD")
var _, git_dirty_err = sh.Output("git", "status", "--porcelain")

//if git_dirty_err == nil {
var GIT_DIRTY = "+CHANGES"

//} else {
//	GIT_DIRTY = ""
//}
var GIT_DESCRIBE, git_describe_err = sh.Output("git", "describe", "--tags", "--always")
var LD_FLAGS = fmt.Sprintf("-X main.GitCommit=%s%s -X main.GitDescribe=%s",
	GIT_COMMIT, GIT_DIRTY, GIT_DESCRIBE)

var dockerRun = sh.RunCmd("docker",
	"run", "--rm",
	"-e", fmt.Sprintf("LDFLAGS=%s", LD_FLAGS),
	"-e", fmt.Sprintf("GITHUB_TOKEN=\"%s\"", GITHUB_TOKEN),
	"-v", fmt.Sprintf("%s/vendor:/go/src", ROOT),
	"-v", fmt.Sprintf("%s:/%s/src/%s", ROOT, PROJECT, IMPORT_PATH),
	"-w", fmt.Sprintf("/%s/src/%s", PROJECT, IMPORT_PATH),
	DEV_IMAGE)

var Default = Install

// Uninstall project.
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

	files, files_err := filepath.Glob("*.out")
	if files_err != nil {
		panic(files_err)
	}
	for _, f := range files {
		if rm_err := os.Remove(f); rm_err != nil {
			panic(rm_err)
		}
	}
	//clean_err := sh.Run("go", "clean", "-i", packageName)
	//if clean_err != nil {
	//	return clean_err
	//}

	return nil
}

// Prefix before other make targets to run in your local dev environment.
func Local() error {
	dockerRun()
	os.Mkdir("tmp", 0755)
	ioutil.WriteFile("tmp/dev_image-id", []byte(""), 0644)
	return nil
}

// Configure the docker environment.
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

// Construct the yq binary.
func Build() error {
	delta, _ := target.Dir("build/dev",
		"test", "*.go")
	if !delta {
		return nil
	}

	os.Mkdir("bin", 0755)

	dockerRun("go", "build", "--ldflags", LD_FLAGS)
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

// Construct and install yq.
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

// Publish the product.
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
