package jvm

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var (
	javaCmd      *exec.Cmd
	javaCmdMutex sync.Mutex
)

// IsRemediated checks if remediation.properties configuration contains the patch
func IsRemediated() bool {
	data, err := ioutil.ReadFile("remediation.properties")
	if err != nil {
		return false
	}
	return strings.Contains(string(data), "remediated=true")
}

// StartJavaApp starts the target Java Spring Boot application locally (for local tests)
func StartJavaApp() error {
	javaCmdMutex.Lock()
	defer javaCmdMutex.Unlock()

	if javaCmd != nil && javaCmd.Process != nil {
		_ = javaCmd.Process.Kill()
		_ = javaCmd.Wait()
		javaCmd = nil
	}

	// Paths to search for the compiled jar
	jarPaths := []string{
		"test/vulnerable-app/target/vulnerable-app-simple-1.0.0.jar",
		"vulnerable-app/target/vulnerable-app-simple-1.0.0.jar",
		"COPY/test/vulnerable-app/target/vulnerable-app-simple-1.0.0.jar",
	}

	var jarPath string
	for _, p := range jarPaths {
		if _, err := os.Stat(p); err == nil {
			jarPath = p
			break
		}
	}

	if jarPath == "" {
		return fmt.Errorf("could not find vulnerable-app jar file in paths: %v", jarPaths)
	}

	args := []string{}
	if IsRemediated() {
		args = append(args, "-Dlog4j2.formatMsgNoLookups=true")
	} else {
		args = append(args, "-Dcom.sun.jndi.ldap.object.trustURLCodebase=true")
		args = append(args, "-Djdk.jndi.object.factoriesFilter=*")
		args = append(args, "-Djdk.jndi.ldap.object.factoriesFilter=*")
	}
	args = append(args, "-jar", jarPath)

	javaPath := "java"
	if _, err := os.Stat("/usr/lib/jvm/java-17-openjdk/bin/java"); err == nil {
		javaPath = "/usr/lib/jvm/java-17-openjdk/bin/java"
	}

	javaCmd = exec.Command(javaPath, args...)

	logFile, err := os.Create("java_target.log")
	if err == nil {
		javaCmd.Stdout = logFile
		javaCmd.Stderr = logFile
	}

	return javaCmd.Start()
}

// StopJavaApp stops the local Java application process
func StopJavaApp() {
	javaCmdMutex.Lock()
	defer javaCmdMutex.Unlock()
	if javaCmd != nil && javaCmd.Process != nil {
		_ = javaCmd.Process.Kill()
		_ = javaCmd.Wait()
		javaCmd = nil
	}
}
