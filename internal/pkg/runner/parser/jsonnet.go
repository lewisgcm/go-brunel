package parser

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"
	"github.com/pkg/errors"
	"go-brunel/internal/pkg/runner/environment"
	"go-brunel/internal/pkg/runner/vcs"
	"go-brunel/internal/pkg/shared"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	tempDirectory = "./.brunel-tmp/"
)

type sharedLibrary struct {
	Repository string
	Branch     string
	File       string
}

type JsonnetParser struct {
	WorkingDirectory    string
	EnvironmentProvider environment.Provider
	VCS                 vcs.VCS
}

var library = `
local brunel = {
    shared(config):: std.parseJson(std.native('shared')(std.toString(config))),
    secret(name):: std.native('secret')(name),
    environment(name):: std.native('environment')(name)
};
`

// Parse will load and parse the supplied jsonnet file, progress io.Writer is used to write the progress of any vcs output.
// Function calls i.e brunel.shared etc are resolved at time of parsing
func (parser *JsonnetParser) Parse(file string, progress io.Writer) (*shared.Spec, error) {
	s, err := parser.parseToJSON(parser.WorkingDirectory, sharedLibrary{File: file}, []sharedLibrary{}, progress)
	_ = os.RemoveAll(parser.WorkingDirectory + tempDirectory) // Clean up our temp directory used for fetching shared libraries etc
	if err != nil {
		return nil, errors.Wrap(err, "error parsing jsonnet file")
	}

	var spec shared.Spec
	err = json.Unmarshal([]byte(s), &spec)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing json")
	}
	return &spec, nil
}

// hasBeenLoaded checks to make sure a library hasnt already been loaded
func hasBeenLoaded(s sharedLibrary, stack []sharedLibrary) bool {
	for _, l := range stack {
		if l.File == s.File && l.Branch == s.Branch && l.Repository == s.Repository {
			return true
		}
	}
	return false
}

// String converts a shared library to a string, this is used for printing our library stack if there are circular dependencies
func (l *sharedLibrary) String() string {
	if l.Repository == "" && l.Branch == "" {
		return l.File
	} else {
		return fmt.Sprintf("%s@%s/%s", l.Repository, l.Branch, l.File)
	}
}

// stackString will return the current stack + last library as a human readable string
func stackString(last sharedLibrary, stack []sharedLibrary) string {
	s := ""
	for _, l := range stack {
		s += l.String()
	}
	s += last.String()
	return s
}

func (parser *JsonnetParser) loadSharedLibraries(workingDir string, stack []sharedLibrary, progress io.Writer, args []interface{}) (interface{}, error) {
	// Load shared library
	var rep sharedLibrary
	err := json.Unmarshal([]byte(args[0].(string)), &rep)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing shared library")
	}

	// We only need to clone if a repository has been specified
	if rep.Repository != "" && rep.Branch != "" {

		// If we are still in the working directory, then we need to switch to our temp directory.
		// This is where we store any external libraries, this simplifies cleanup when we are done
		if workingDir == parser.WorkingDirectory {
			workingDir = workingDir + tempDirectory
			if err := os.MkdirAll(workingDir, os.ModePerm); err != nil {
				return nil, errors.Wrap(err, "error creating temp working directory")
			}
		}

		// Generate an md5 based off of repo and branch so we dont need to worry about illegal directory path characters
		// We also use this as a sort of 'caching' mechanism to know if we need to re-clone or not
		workingDir = workingDir + fmt.Sprintf("%x/", md5.Sum([]byte(rep.Repository+rep.Branch)))
		if _, err := os.Stat(workingDir); os.IsNotExist(err) {
			if err = parser.VCS.Clone(vcs.Options{
				Directory:     workingDir,
				RepositoryURL: rep.Repository,
				Branch:        rep.Branch,
				Progress:      progress,
			}); err != nil {
				return nil, errors.Wrap(err, "error cloning shared library")
			}
		}
	}

	// We check if we have already loaded the current library by looking it up in our stack,
	// if it has been loaded then error out.
	if hasBeenLoaded(rep, stack) {
		return nil, fmt.Errorf("circular dependency in imports: %s", stackString(rep, stack))
	}

	// Load in the shared library based on our sharedLibrary object
	js, err := parser.parseToJSON(workingDir, rep, stack, progress)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing file")
	}
	return js, nil
}

func (parser *JsonnetParser) parseToJSON(workingDir string, lib sharedLibrary, stack []sharedLibrary, progress io.Writer) (string, error) {
	rel, err := filepath.Rel(parser.WorkingDirectory, workingDir+lib.File)
	if err != nil {
		return "", errors.Wrap(err, "error getting relative file path")
	}

	if strings.Contains(rel, "../") {
		return "", errors.New(fmt.Sprintf("request path %s is outside of working directory %s", workingDir+lib.File, parser.WorkingDirectory))
	}

	snippet, err := ioutil.ReadFile(filepath.Clean(workingDir + lib.File))
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("error reading file '%s' in workspace '%s'", lib.File, workingDir))
	}

	// We append our current file to the stack to prevent any circular dependencies
	stack = append(stack, lib)

	vm := jsonnet.MakeVM()

	// Register our method for loading shared libraries
	vm.NativeFunction(&jsonnet.NativeFunction{
		Func: func(args []interface{}) (interface{}, error) {
			return parser.loadSharedLibraries(workingDir, stack, progress, args)
		},
		Name:   "shared",
		Params: ast.Identifiers{"string"},
	})

	vm.NativeFunction(&jsonnet.NativeFunction{
		Func: func(args []interface{}) (interface{}, error) {
			return parser.EnvironmentProvider.GetValue(args[0].(string))
		},
		Name:   "environment",
		Params: ast.Identifiers{"string"},
	})

	vm.NativeFunction(&jsonnet.NativeFunction{
		Func: func(args []interface{}) (interface{}, error) {
			return parser.EnvironmentProvider.GetSecret(args[0].(string))
		},
		Name:   "secret",
		Params: ast.Identifiers{"string"},
	})

	js, err := vm.EvaluateSnippet(lib.File, library+string(snippet))
	if err != nil {
		return "", errors.Wrap(err, "error evaluating snippet")
	}
	return js, nil
}
