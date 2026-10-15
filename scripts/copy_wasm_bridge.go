package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

func main() {
	goroot := runtime.GOROOT()
	if goroot == "" {
		fmt.Println("Error: runtime unable to detect GOROOT path")
		os.Exit(1)
	}

	// Dynamic fallback sequence to locate the webassembly bootstrap script across various Go layouts
	possibleSrcPaths := []string{
		filepath.Join(goroot, "lib", "wasm", "wasm_exec.js"),  // Modern location
		filepath.Join(goroot, "misc", "wasm", "wasm_exec.js"), // Traditional baseline Go standard location
	}

	var srcPath string
	var srcFile *os.File
	var err error

	// Iterate through paths to find an accessible source bundle
	for _, path := range possibleSrcPaths {
		srcFile, err = os.Open(path)
		if err == nil {
			srcPath = path
			break
		}
	}

	if err != nil {
		fmt.Printf("Critical error: unable to locate 'wasm_exec.js' inside GOROOT layouts: %v\n", possibleSrcPaths)
		os.Exit(1)
	}
	defer srcFile.Close()

	fmt.Printf("[INFO] Found source bridge file at: %s\n", srcPath)

	// Resolve destination path pointing to the static public web layer environment
	dstPath := filepath.Join("web", "app", "public", "wasm_exec.js")

	// Ensure parent directories exist before execution
	err = os.MkdirAll(filepath.Dir(dstPath), 0755)
	if err != nil {
		fmt.Printf("Error creating destination directories: %v\n", err)
		os.Exit(1)
	}

	dstFile, err := os.Create(dstPath)
	if err != nil {
		fmt.Printf("Error creating destination file: %v\n", err)
		os.Exit(1)
	}
	defer dstFile.Close()

	// Perform fast byte copy
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		fmt.Printf("Error copying data stream: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("[SUCCESS] wasm_exec.js successfully copied into public directory.")
}
