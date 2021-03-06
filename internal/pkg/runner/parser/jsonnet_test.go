package parser

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"go-brunel/internal/pkg/runner/environment"
	"go-brunel/internal/pkg/runner/trigger"
	vcs2 "go-brunel/internal/pkg/runner/vcs"
	"go-brunel/internal/pkg/shared"
	"go-brunel/test"
	"go-brunel/test/mocks/go-brunel/pkg/runner/vcs"
	"io/ioutil"
	"os"
	"testing"
)

func TestJsonnetParser_Parse(t *testing.T) {
	var suites = []struct {
		env        map[string]string                           // environment variables to create for testing
		files      map[string]string                           // files to create for testing
		cloneFiles map[string]string                           // files to create when vcs.Clone is called
		cloneTimes int                                         // How many times should vcs.Clone be called?
		expect     func(t *testing.T, s *shared.Spec, e error) // Function gor handling expects
	}{
		// Tests that a broken file wont be parsed
		{
			files: map[string]string{
				".brunel.jsonnet": `sdsd`,
			},
			expect: func(t *testing.T, spec *shared.Spec, err error) {
				test.ExpectErrorLike(t, errors.New("error parsing"), err)
			},
		},

		// Tests that stage names should be unique
		{
			files: map[string]string{
				".brunel.jsonnet": `
{
    stages: [
		{
			name: 'test'
		},
		{
			name: 'test'
		}
    ]
}`,
			},
			expect: func(t *testing.T, spec *shared.Spec, err error) {
				test.ExpectErrorLike(t, errors.New("stage names should be unique"), err)
			},
		},

		// Tests that stage names cannot be empty
		{
			files: map[string]string{
				".brunel.jsonnet": `
{
    stages: [
		{
			name: '   '
		},
    ]
}`,
			},
			expect: func(t *testing.T, spec *shared.Spec, err error) {
				test.ExpectErrorLike(t, errors.New("stage names must be specified"), err)
			},
		},

		// Tests that build information is available
		{
			files: map[string]string{
				".brunel.jsonnet": `
{
    stages: [
		{
			name: brunel.build.revision,
			when: brunel.match('^revision$', brunel.build.revision)
		},
		{
			name: brunel.build.branch
		},
    ]
}`,
			},
			expect: func(t *testing.T, spec *shared.Spec, err error) {
				if err != nil {
					t.Fatal(err)
				}
				if spec.Stages[0].When == nil || *spec.Stages[0].When != true {
					t.Fail()
					t.Log("when not correctly calculated")
				}
				test.ExpectString(t, "revision", string(spec.Stages[0].ID))
				test.ExpectString(t, "branch", string(spec.Stages[1].ID))
			},
		},

		// Tests that we can read a file, load local shared library and read values from our environment in both the library and local file
		{
			env: map[string]string{
				"MY_ENV": "my env!!",
			},
			files: map[string]string{
				".brunel.jsonnet": `
local shared = brunel.shared({
    file: "file.jsonnet"
});

{
    version: "v1",
    description: shared.SOME_VAL,
    stages: [
        {
			name: 'test',
            services: [
                {
                    image: "nginx:latest",
                    hostname: brunel.environment.variable("MY_ENV")
                }
            ],
            steps: [
                {
                    image: "byrnedo/alpine-curl",
                    entrypoint: brunel.environment.variable("MY_ENV"),
                    args: [ "-c", "--", "curl http://nginx" ]
                },

            ],
        },
    ]
}`,
				"file.jsonnet": `
{
	SOME_VAL: "value " + brunel.environment.variable("MY_ENV")
}`,
			},
			expect: func(t *testing.T, spec *shared.Spec, err error) {
				test.ExpectString(t, spec.Version, "v1")
				test.ExpectString(t, spec.Description, "value my env!!")
				test.ExpectString(t, spec.Stages[0].Services[0].Image, "nginx:latest")
				test.ExpectString(t, spec.Stages[0].Services[0].Hostname, "my env!!")
				test.ExpectString(t, spec.Stages[0].Steps[0].EntryPoint, "my env!!")
			},
		},

		// Tests that we dont allow circular dependencies with local files
		{
			files: map[string]string{
				".brunel.jsonnet": `
local shared = brunel.shared({
    file: "file.jsonnet"
});

{
    version: "v1",
    description: shared.SOME_VAL,
    stages: []
}`,
				"file.jsonnet": `
local shared = brunel.shared({
    file: "file.jsonnet"
});
{
	SOME_VAL: shared.SOME_VAL
}`,
			},
			expect: func(t *testing.T, spec *shared.Spec, err error) {
				test.ExpectErrorLike(t, errors.New("circular dependency in imports"), err)
			},
		},

		// Tests that we will throw an error if we load a non-existent local library
		{
			files: map[string]string{
				".brunel.jsonnet": `
local shared = brunel.shared({
    file: "non-existent-file.jsonnet"
});

{
    version: "v1",
    description: shared.SOME_VAL,
    stages: []
}`,
			},
			expect: func(t *testing.T, spec *shared.Spec, err error) {
				test.ExpectErrorLike(t, errors.New("error reading file"), err)
			},
		},

		// Tests that we cant break out of our current working directory (relative path)
		{
			files: map[string]string{
				".brunel.jsonnet": `
local shared = brunel.shared({
    file: "../../../../file.jsonnet"
});

{
    version: "v1",
    description: shared.SOME_VAL,
    stages: []
}`,
			},
			expect: func(t *testing.T, spec *shared.Spec, err error) {
				test.ExpectErrorLike(t, errors.New("is outside of working directory"), err)
			},
		},

		// Tests that we can load a shared library from git and that we will cache the clone such that other
		// references to the same library will not be re-downloaded.
		{
			files: map[string]string{
				".brunel.jsonnet": `
local shared = brunel.shared({
	repository: "http://some-repo.com",
	branch: "some-branch",
    file: "file.jsonnet"
});

local sharedTwo = brunel.shared({
	repository: "http://some-repo.com",
	branch: "some-branch",
    file: "file.jsonnet"
});

{
    version: "v1",
	description: shared.description + sharedTwo.description,
    stages: []
}`,
			},
			cloneFiles: map[string]string{
				"file.jsonnet": `
{
	description: "description"
}`,
			},
			cloneTimes: 1,
			expect: func(t *testing.T, spec *shared.Spec, err error) {
				test.ExpectError(t, nil, err)
				test.ExpectString(t, "descriptiondescription", spec.Description)
			},
		},

		// Tests that when we load vcs libraries, that also contain references to other vcs libraries we can clone and
		// handle them correctly
		{
			files: map[string]string{
				".brunel.jsonnet": `
local shared = brunel.shared({
	repository: "http://some-repo.com",
	branch: "some-branch",
    file: "file.jsonnet"
});

{
    version: "v1",
	description: shared.description,
    stages: []
}`,
			},
			cloneFiles: map[string]string{
				"file.jsonnet": `
local shared = brunel.shared({
	repository: "http://some-repo.com",
	branch: "some-branch",
    file: "file2.jsonnet"
});

{
	description: shared.description
}`,
				"file2.jsonnet": `
{
	description: "lolz"
}`,
			},
			cloneTimes: 2,
			expect: func(t *testing.T, spec *shared.Spec, err error) {
				test.ExpectError(t, nil, err)
				test.ExpectString(t, spec.Description, "lolz")
			},
		},

		// Tests that we can detect circular dependencies between shared vcs libraries
		{
			files: map[string]string{
				".brunel.jsonnet": `
local shared = brunel.shared({
	repository: "http://some-repo.com",
	branch: "some-branch",
    file: "file.jsonnet"
});

{
    version: "v1",
	description: shared.description,
    stages: []
}`,
			},
			cloneFiles: map[string]string{
				"file.jsonnet": `
local shared = brunel.shared({
    file: "file2.jsonnet"
});

{
	description: shared.description
}`,
				"file2.jsonnet": `
local shared = brunel.shared({
	repository: "http://some-repo.com",
	branch: "some-branch",
    file: "file.jsonnet"
});

{
	description: shared.description
}`,
			},
			cloneTimes: 2,
			expect: func(t *testing.T, spec *shared.Spec, err error) {
				test.ExpectErrorLike(t, errors.New("circular dependency in imports"), err)
			},
		},
	}
	for i, suite := range suites {
		t.Run(
			fmt.Sprintf("suites[%d]", i),
			func(t *testing.T) {
				testWorkSpaceDir := ".test-workspace/"
				_ = os.RemoveAll(testWorkSpaceDir)

				if err := os.Mkdir(testWorkSpaceDir, os.ModePerm); err != nil {
					t.Fatal(err)
				}

				for e, v := range suite.env {
					if err := os.Setenv(e, v); err != nil {
						t.Fatal(err)
					}
				}

				for f, c := range suite.files {
					if err := ioutil.WriteFile(testWorkSpaceDir+f, []byte(c), os.ModePerm); err != nil {
						t.Fatal(err)
					}
				}

				factory := environment.LocalEnvironmentFactory{}

				controller := gomock.NewController(t)
				mockVCS := vcs.NewMockVCS(controller)

				mockVCS.
					EXPECT().
					Clone(gomock.Any()).
					Return(nil).
					Times(suite.cloneTimes).
					Do(func(options vcs2.Options) {
						_ = os.MkdirAll(options.Directory, os.ModePerm)
						for f, c := range suite.cloneFiles {
							_ = ioutil.WriteFile(options.Directory+f, []byte(c), os.ModePerm)
						}
					})

				p := JsonnetParser{
					Event: trigger.Event{
						Job: shared.Job{
							Commit: shared.Commit{
								Branch:   "branch",
								Revision: "revision",
							},
						},
						JobState: nil,
						WorkDir:  testWorkSpaceDir,
						Context:  nil,
					},
					EnvironmentProvider: factory.Create(nil),
					VCS:                 mockVCS,
				}
				spec, err := p.Parse(".brunel.jsonnet", ioutil.Discard)
				suite.expect(t, spec, err)

				for e := range suite.env {
					if err := os.Unsetenv(e); err != nil {
						t.Fatal(err)
					}
				}

				if err := os.RemoveAll(testWorkSpaceDir); err != nil {
					t.Fatal(err)
				}

				controller.Finish()
			},
		)
	}
}
