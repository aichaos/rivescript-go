package rivescript

// Loading and Parsing Methods

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/shurcooL/httpfs/path/vfspath"
)

/*
LoadFile loads a single RiveScript source file from disk.

Parameters

	path: Path to a RiveScript source file.
*/
func (rs *RiveScript) LoadFile(path string) error {
	rs.say("Load RiveScript file: %s", path)

	fh, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("Failed to open file %s: %s", path, err)
	}

	defer fh.Close()
	scanner := bufio.NewScanner(fh)
	scanner.Split(bufio.ScanLines)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return rs.parse(path, lines)
}

/*
LoadFile loads a single RiveScript source file from disk.

Parameters

	path: Path to a RiveScript source file.
*/
func (rs *RiveScript) LoadHttpFile(httpfs http.FileSystem, path string) error {
	rs.say("Load RiveScript file: %s", path)

	fh, err := httpfs.Open(path)
	if err != nil {
		return fmt.Errorf("Failed to open file %s: %s", path, err)
	}

	defer fh.Close()
	scanner := bufio.NewScanner(fh)
	scanner.Split(bufio.ScanLines)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return rs.parse(path, lines)
}

/*
LoadDirectory loads multiple RiveScript documents from a folder on disk.

Parameters

	path: Path to the directory on disk
	extensions...: List of file extensions to filter on, default is
	               '.rive' and '.rs'
*/
func (rs *RiveScript) LoadDirectory(path string, extensions ...string) error {
	if len(extensions) == 0 {
		extensions = []string{".rive", ".rs"}
	}

	files, err := filepath.Glob(fmt.Sprintf("%s/*", path))
	if err != nil {
		return fmt.Errorf("Failed to open folder %s: %s", path, err)
	}

	// No files matched?
	if len(files) == 0 {
		return fmt.Errorf("No RiveScript source files were found in %s", path)
	}

	var anyValid bool
	for _, f := range files {
		// Restrict file extensions.
		validExtension := false
		for _, exten := range extensions {
			if strings.HasSuffix(f, exten) {
				validExtension = true
				break
			}
		}

		if validExtension {
			anyValid = true
			err := rs.LoadFile(f)
			if err != nil {
				return err
			}
		}
	}

	if !anyValid {
		return fmt.Errorf("No RiveScript source files were found in %s", path)
	}

	return nil
}

/*
LoadHttpFilesystem loads multiple RiveScript documents from a folder on disk.

Parameters

	path: Path to the directory on disk
	extensions...: List of file extensions to filter on, default is
	               '.rive' and '.rs'
*/
func (rs *RiveScript) LoadHttpFilesystem(httpfs http.FileSystem, path string, extensions ...string) error {
	if len(extensions) == 0 {
		extensions = []string{".rive", ".rs"}
	}

	files, err := vfspath.Glob(httpfs, fmt.Sprintf("%s/*", path))
	if err != nil {
		return fmt.Errorf("Failed to open folder %s: %s", path, err)
	}

	// No files matched?
	if len(files) == 0 {
		return fmt.Errorf("No RiveScript source files were found in %s", path)
	}

	var anyValid bool
	for _, f := range files {
		// Restrict file extensions.
		validExtension := false
		for _, exten := range extensions {
			if strings.HasSuffix(f, exten) {
				validExtension = true
				break
			}
		}

		if validExtension {
			anyValid = true
			err := rs.LoadHttpFile(httpfs, f)
			if err != nil {
				return err
			}
		}
	}

	if !anyValid {
		return fmt.Errorf("No RiveScript source files were found in %s", path)
	}

	return nil
}

/*
LoadMultipleDirectory loads multiple RiveScript documents from a comma separated list of folders on disk.

Parameters

	comma separated list of paths: Paths to the directorys on disk
	extensions...: List of file extensions to filter on, default is
	               '.rive' and '.rs'
*/
func (rs *RiveScript) LoadMultipleDirectory(paths []string, extensions ...string) error {
	var files []string
	if len(extensions) == 0 {
		extensions = []string{".rive", ".rs"}
	}
	for _, path := range paths {
		f, err := filepath.Glob(fmt.Sprintf("%s/*", path))
		if err != nil {
			return fmt.Errorf("Failed to open folder %s: %s", path, err)
		}
		for _, line := range f {
			files = append(files, line)
		}
	}

	// No files matched?
	if len(files) == 0 {
		return fmt.Errorf("No RiveScript source files were found in %s", paths)
	}

	var anyValid bool
	for _, f := range files {
		// Restrict file extensions.
		validExtension := false
		for _, exten := range extensions {
			if strings.HasSuffix(f, exten) {
				validExtension = true
				break
			}
		}

		if validExtension {
			anyValid = true
			err := rs.LoadFile(f)
			if err != nil {
				return err
			}
		}
	}

	if !anyValid {
		return fmt.Errorf("No RiveScript source files were found in %s", paths)
	}

	return nil
}

/*
Stream loads RiveScript code from a text buffer.

Parameters

	code: Raw source code of a RiveScript document, with line breaks after
	      each line.
*/
func (rs *RiveScript) Stream(code string) error {
	lines := strings.Split(code, "\n")
	return rs.parse("Stream()", lines)
}
