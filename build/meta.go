package main

import(
	"os"
	"os/exec"
	"fmt"
	"strings"
)

func runGit(args []string) (string, int, error) {
	git, err := exec.LookPath("git")
	if(err != nil) {
		return "", -1, err
	}
	cmd := exec.Command(git, args...)
	out, err := cmd.Output()
	code := cmd.ProcessState.ExitCode()
	return string(out), code, err
}

func gitExtractTag() (string, int, error) {
	return runGit([]string{"describe", "--match",  "'v[0-9]*'",  "--dirty=.m",  "--always", "--tags"})
}

func gitExtractRevision() (string, int, error) {
	return runGit([]string{"rev-parse", "HEAD"})
}

func isDirty() (bool, error) {
	_, code, err := runGit([]string{"diff", "--no-ext-diff", "--quiet",  "--exit-code"})
	if code == 1 {
		return true, nil
	}
	if err != nil {
		return false, err
	}
	return code == 0, err
}

func main() {

	var command string

	if(len(os.Args) > 1) {
		command = os.Args[1]
	}

	tag, _, err := gitExtractTag()
	tag = strings.TrimSpace(tag)
	if(err != nil) {
		fmt.Printf("Error running git %s", err)
		os.Exit(1)
	}

	if(command == "version") {
		fmt.Printf(tag)
		os.Exit(0)
	}

	fmt.Printf("Tag=%s\n", tag)
	revision, _, err := gitExtractRevision()
	if(err != nil) {
		fmt.Printf("Error running git %s", err)
		os.Exit(1)
	}
	revision = strings.TrimSpace(revision)
	fmt.Printf("Revision=%s\n", revision)
	dirty, err := isDirty()
	if(err != nil) {
		fmt.Printf("Error running git %s", err)
		os.Exit(1)
	}
	fmt.Printf("Dirty=%t\n", dirty)
	data := fmt.Sprintf("%s|%s|%t", tag, revision, dirty)
	err = os.WriteFile("pkg/version/version.txt", []byte(data), 0777)
	if(err != nil) {
		fmt.Printf("Error Writing file version.txt %s", err)
		os.Exit(1)
	}
}