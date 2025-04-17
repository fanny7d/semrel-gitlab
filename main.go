/*
go-semrel-gitlab provides tools to automate parts of release process on Gitlab CI

More documentation can be found at https://juhani.gitlab.io/go-semrel-gitlab/
*/
package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"

	"github.com/blang/semver"
	"github.com/fanny7d/semrel-gitlab/pkg/actions"
	"github.com/fanny7d/semrel-gitlab/pkg/gitlabutil"
	"github.com/fanny7d/semrel-gitlab/pkg/render"
	"github.com/fanny7d/semrel-gitlab/pkg/workflow"
	"github.com/juranki/go-semrel/angularcommit"
	"github.com/juranki/go-semrel/inspectgit"
	"github.com/juranki/go-semrel/semrel"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	gitlab "github.com/xanzy/go-gitlab"
)

var (
	version string
)

type fakeFixChange struct {
	Scope   string
	Subject string
	Hash    string
}

func (c *fakeFixChange) Category() string            { return "fix" }
func (c *fakeFixChange) BumpLevel() semrel.BumpLevel { return semrel.BumpPatch }
func (c *fakeFixChange) PreReleased() bool           { return false }

func clientFromContext(c *cli.Context) (*gitlab.Client, error) {
	token := c.GlobalString("token")
	if token == "" {
		return nil, errors.New("token is required")
	}
	api := c.GlobalString("api")
	if api == "" {
		return nil, errors.New("api is required")
	}
	return gitlabutil.NewClient(token, api, c.GlobalBool("skip-ssl-verify"))
}

func parseCommaSeparatedList(s string) []string {
	if len(strings.TrimSpace(s)) == 0 {
		return []string{}
	}
	ss := strings.Split(s, ",")
	rv := []string{}
	for _, t := range ss {
		rv = append(rv, strings.TrimSpace(t))
	}
	return rv
}

func hasReleaseContent(info *semrel.ReleaseData) bool {
	if info.BumpLevel == semrel.NoBump {
		// no changes found
		return false
	}
	if len(info.NextVersion.Pre)+len(info.NextVersion.Build) == 0 {
		// not a pre-release, all changes apply
		return true
	}
	// this is a pre-release, walk changes to find unreleased
	for key, changes := range info.Changes {
		if key == "other" {
			// skip other changes
			continue
		}
		for _, change := range changes {
			if !change.PreReleased() {
				// found unreleased commit
				return true
			}
		}
	}
	// no release content found
	return false
}

func analyzeCommits(c *cli.Context) (*semrel.ReleaseData, error) {
	var patchTypes []string
	var minorTypes []string

	listOtherChanges := c.Bool("list-other-changes")
	patchTypesString := c.GlobalString("patch-commit-types")
	minorTypesString := c.GlobalString("minor-commit-types")
	tagPrefix := strings.TrimSpace(c.GlobalString("tag-prefix"))

	if len(patchTypesString) == 0 {
		patchTypes = parseCommaSeparatedList("fix,refactor,perf,docs,style,test")
	} else {
		patchTypes = parseCommaSeparatedList(strings.ToLower(patchTypesString))
	}
	if len(minorTypesString) == 0 {
		minorTypes = parseCommaSeparatedList("feat")
	} else {
		minorTypes = parseCommaSeparatedList(strings.ToLower(minorTypesString))
	}

	// Get new commits
	vcsData, err := inspectgit.VCSDataWithPrefix(".", tagPrefix)
	if err != nil {
		return nil, err
	}

	// Analyze commit messages
	releaseData, err := semrel.Release(vcsData, angularcommit.NewWithOptions(&angularcommit.Options{
		FixTypes:     patchTypes,
		FeatureTypes: minorTypes,
		BreakingChangeMarkers: []string{
			`BREAKING\s+CHANGE:`,
			`BREAKING\s+CHANGE`,
			`BREAKING:`,
		},
	}))
	if err != nil {
		return nil, err
	}
	if _, hasOther := releaseData.Changes["other"]; !listOtherChanges && hasOther {
		delete(releaseData.Changes, "other")
	}

	if c.GlobalBool("bump-patch") && releaseData.BumpLevel == semrel.NoBump && len(vcsData.UnreleasedCommits) > 0 {
		// bump patch and add all commit headings to release note
		changes := []semrel.Change{}
		for _, commit := range vcsData.UnreleasedCommits {
			change := &fakeFixChange{
				Subject: strings.TrimSpace(strings.Split(commit.Msg, "\n")[0]),
				Hash:    commit.SHA,
			}
			changes = append(changes, change)
		}
		releaseData.Changes["fix"] = changes
		releaseData.NextVersion = releaseData.CurrentVersion
		releaseData.NextVersion.Patch++
		releaseData.BumpLevel = semrel.BumpPatch
		if _, hasOther := releaseData.Changes["other"]; hasOther {
			delete(releaseData.Changes, "other")
		}
	}

	if c.GlobalBoolT("initial-development") == false && releaseData.NextVersion.Major == 0 && releaseData.BumpLevel != semrel.NoBump {
		// bump to 1.0.0
		releaseData.NextVersion = semver.MustParse("1.0.0")
		releaseData.BumpLevel = semrel.BumpMajor
	}

	releaseBranches := parseCommaSeparatedList(c.GlobalString("release-branches"))
	if len(releaseBranches) > 0 {
		// Determine if this is a pre-release
		currBranch := c.GlobalString("ci-commit-ref-name")
		if len(currBranch) == 0 {
			return nil, errors.New("CI environment information not available (ci-commit-ref-name)")
		}
		isReleaseBranch := false
		for _, b := range releaseBranches {
			if currBranch == b {
				isReleaseBranch = true
			}
		}
		if !isReleaseBranch {
			// Turn this release a pre-release
			err = render.SetPreReleaseNumber(releaseData, c)
			if err != nil {
				return nil, err
			}
			// Remove commit entries with PreReleased == true
			newChanges := make(map[string][]semrel.Change)
			for category, changes := range releaseData.Changes {
				for _, change := range changes {
					if !change.PreReleased() {
						newChanges[category] = append(newChanges[category], change)
					}
				}
			}
			releaseData.Changes = newChanges
		}
	}
	return releaseData, nil
}

// release short-version
func shortVersion(c *cli.Context) error {
	fmt.Print(version)
	return nil
}

func releaseAPIAvailable(client *gitlab.Client) (bool, error) {
	// Release api is available from v11.7.0-rcN
	v11_7_0 := semver.MustParse("11.7.0")

	glVer, _, err := client.Version.GetVersion()
	if err != nil {
		return false, errors.Wrap(err, "check release api availability")
	}

	v, err := semver.ParseTolerant(glVer.Version)
	if err != nil {
		return false, errors.Wrapf(err, "gitlab version %s", glVer.Version)
	}

	// reset pre and build to include pre-releases of v11.7.0
	v.Pre = []semver.PRVersion{}
	v.Build = []string{}

	return v.GE(v11_7_0), nil
}

// release test-api
func testAPI(c *cli.Context) error {
	client, err := clientFromContext(c)
	if err != nil {
		return err
	}
	check := actions.NewCheck(client)
	return check.Do()
}

// release add-download
func addDownload(c *cli.Context) error {
	client, err := clientFromContext(c)
	if err != nil {
		return err
	}
	file := c.String("file")
	if file == "" {
		return errors.New("file is required")
	}
	tag := c.String("ci-commit-tag")
	if tag == "" {
		tag = os.Getenv("CI_COMMIT_TAG")
	}
	if tag == "" {
		return errors.New("ci-commit-tag is required")
	}
	projectURL, err := url.Parse(c.GlobalString("project-url"))
	if err != nil {
		return errors.Wrap(err, "parse project url")
	}
	getTag := actions.NewGetTag(client, c.GlobalString("project"), tag)
	if err := getTag.Do(); err != nil {
		return err
	}
	upload := actions.NewUpload(client, c.GlobalString("project"), projectURL, file)
	if err := upload.Do(); err != nil {
		return err
	}
	addLink := actions.NewAddLink(&actions.AddLinkParams{
		Client:               client,
		Project:              c.GlobalString("project"),
		LinkDescription:      "",
		MDLinkFunc:           upload.MDLinkFunc(),
		TagFunc:              getTag.TagFunc(),
		LinkURLFunc:          upload.LinkURLFunc(),
		ReleasesAPIAvailable: c.GlobalBool("releases-api-available"),
	})
	return addLink.Do()
}

// release add-download-link
func addDownloadLink(c *cli.Context) error {
	project := c.GlobalString("ci-project-path")
	tag := c.GlobalString("ci-commit-tag")
	url := c.String("url")
	name := c.String("name")
	description := c.String("description")

	if len(project) == 0 {
		return errors.New("ci-project-path is not set")
	}
	if len(tag) == 0 {
		return errors.New("ci-commit-tag is not set")
	}
	if len(name) == 0 || len(url) == 0 || len(description) == 0 {
		return errors.New("name, url and description must be specified")
	}

	client, err := clientFromContext(c)
	if err != nil {
		return errors.Wrap(err, "Unable to connect")
	}
	relAPI, err := releaseAPIAvailable(client)
	if err != nil {
		return errors.Wrap(err, "detect releases api")
	}
	getTag := actions.NewGetTag(client, project, tag)
	addLink := actions.NewAddLink(
		&actions.AddLinkParams{
			Client:               client,
			Project:              project,
			LinkDescription:      description,
			MDLinkFunc:           func() string { return fmt.Sprintf("[%s](%s)", name, url) },
			LinkURLFunc:          func() string { return url },
			TagFunc:              getTag.TagFunc(),
			ReleasesAPIAvailable: relAPI,
		},
	)
	return workflow.Apply([]workflow.Action{getTag, addLink})
}

// release commit-and-tag
// release tag-and-commit
func commitAndTagBase(c *cli.Context, reversed bool) error {
	client, err := clientFromContext(c)
	if err != nil {
		return err
	}
	tag := c.String("tag")
	if tag == "" {
		return errors.New("tag is required")
	}
	getTag := actions.NewGetTag(client, c.GlobalString("project"), tag)
	if err := getTag.Do(); err != nil {
		return err
	}
	commit := actions.NewCommit(client, c.GlobalString("project"), c.String("branch"), c.String("message"))
	if err := commit.Do(); err != nil {
		return err
	}
	createTag := actions.NewCreateTag(client, c.GlobalString("project"), commit.BranchFunc(), tag, c.String("message"), c.Bool("force"))
	if err := createTag.Do(); err != nil {
		return err
	}
	createPipeline := actions.NewCreatePipeline(client, c.GlobalString("project"), commit.BranchFunc())
	return createPipeline.Do()
}

// release tag
func tag(c *cli.Context) error {
	tagPrefix := strings.TrimSpace(c.GlobalString("tag-prefix"))
	project := c.GlobalString("ci-project-path")
	sha := c.GlobalString("ci-commit-sha")
	client, err := clientFromContext(c)
	if err != nil {
		return errors.Wrap(err, "Unable to connect")
	}
	relAPI, err := releaseAPIAvailable(client)
	if err != nil {
		return errors.Wrap(err, "check releases api")
	}

	if len(project) == 0 {
		return errors.New("ci-project-path is not set")
	}
	if len(sha) == 0 {
		return errors.New("ci-commit-sha is not set")
	}

	info, err := analyzeCommits(c)
	if err != nil {
		return err
	}
	if !hasReleaseContent(info) {
		return errors.New("no changes found in commit messages")
	}

	releaseNote, err := render.ReleaseNote(info)
	if err != nil {
		return err
	}

	createTag := actions.NewCreateTag(
		client,
		project,
		actions.NewFuncOfString(sha),
		tagPrefix+info.NextVersion.String(),
		releaseNote,
		relAPI)

	return workflow.Apply([]workflow.Action{createTag})
}

// release next-version
// release test-git
func inspectAndPrint(c *cli.Context, version bool, releaseNote bool) error {
	bumpPatch := c.Bool("bump-patch")
	allowCurrent := c.Bool("allow-current")

	if bumpPatch && allowCurrent {
		return errors.New("bump-patch and allow-current are mutually exclusive")
	}

	info, err := analyzeCommits(c)
	if err != nil {
		return errors.Wrap(err, "analyze commits")
	}
	if !hasReleaseContent(info) {
		if allowCurrent {
			if len(info.NextVersion.Pre)+len(info.NextVersion.Build) == 0 {
				info.NextVersion = info.CurrentVersion
			}
		} else if bumpPatch {
			if len(info.NextVersion.Pre)+len(info.NextVersion.Build) == 0 {
				info.NextVersion = info.CurrentVersion
				info.NextVersion.Patch++
			}
		} else {
			return errors.New("commit log didn't contain changes that would change the version")
		}
	}
	if version {
		fmt.Println(info.NextVersion.String())
	}
	if releaseNote {
		releaseNoteText, err := render.ReleaseNote(info)
		if err != nil {
			return err
		}
		fmt.Println(releaseNoteText)
	}
	return nil
}

// release changelog
func changelog(c *cli.Context) error {
	filename := c.String("f")

	if len(filename) == 0 {
		filename = "CHANGELOG.md"
	}

	info, err := analyzeCommits(c)
	if err != nil {
		return errors.Wrap(err, "analyze commits")
	}

	if !hasReleaseContent(info) {
		return errors.New("no changes found in commit messages")
	}

	changelogEntry, err := render.ChangelogEntry(info)
	if err != nil {
		return err
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// create a new changelog file
		fileparts := []string{
			"# CHANGELOG",
			"<!--- next entry here -->",
			changelogEntry,
		}
		data := strings.Join(fileparts, "\n\n")
		err = ioutil.WriteFile(filename, []byte(data), 0644)
		if err != nil {
			return errors.Wrap(err, "changelog, write file")
		}
		fmt.Printf("wrote %s\n", filename)
	} else {
		// insert new changelog entry to an existing file
		bytes, err := ioutil.ReadFile(filename)
		if err != nil {
			return errors.Wrap(err, "changelog, read file")
		}

		content := strings.Split(string(bytes), "<!--- next entry here -->")
		if len(content) < 2 {
			return errors.Errorf(`Next entry marker not found.
Please copy following tag to %s to mark where you want new entries inserted.

<!--- next entry here -->
`, filename)
		}
		if len(content) > 2 {
			return errors.Errorf(`Too many markers (<!--- next entry here -->).
Please edit %s so that there's only one.`, filename)
		}
		parts := []string{
			strings.TrimRight(content[0], " \n\r\t"),
			"<!--- next entry here -->",
			changelogEntry,
			strings.TrimLeft(content[1], " \n\r\t"),
		}
		data := strings.Join(parts, "\n\n")

		err = ioutil.WriteFile(filename, []byte(data), 0644)
		if err != nil {
			return errors.Wrap(err, "changelog, write file")
		}
		fmt.Printf("updated %s\n", filename)
	}
	return nil
}

func main() {
	if len(version) == 0 {
		version = "DEV"
	}
	listOtherChanges := cli.BoolFlag{
		Name:   "list-other-changes",
		Usage:  "List changes that don't affect versioning",
		EnvVar: "GSG_LIST_OTHER_CHANGES",
	}
	app := cli.NewApp()
	app.Usage = "semantic release tools for GitLab"
	app.Description = `A collection of utilities to help automate releases

   The recommended way to use global options is to assign value to the
   associated environment variables in Gitlab CI Variables. Please note
   that the ci-* options are automatically populated by Gitlab CI.`
	app.Version = version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "token, t",
			Usage:  "gitlab private `TOKEN`",
			EnvVar: "GITLAB_TOKEN,GL_TOKEN",
		},
		cli.StringFlag{
			Name: "gl-api",
			Usage: `gitlab api URL. If not defined infers from CI_PROJECT_URL.
                            If that fails too, uses https://gitlab.com/api/v4/.`,
			EnvVar: "GITLAB_URL,GL_URL",
		},
		cli.BoolFlag{
			Name:   "skip-ssl-verify",
			Usage:  "don't verify CA certificate on gitlab api",
			Hidden: false,
		},
		cli.StringFlag{
			Name: "patch-commit-types",
			Usage: `comma separated list of commit message types that indicate patch bump
                               (default: "fix,refactor,perf,docs,style,test")`,
			Hidden: false,
			EnvVar: "GSG_PATCH_COMMIT_TYPES",
		},
		cli.StringFlag{
			Name: "minor-commit-types",
			Usage: `comma separated list of commit message types that indicate minor bump
                               (default: "feat")`,
			Hidden: false,
			EnvVar: "GSG_MINOR_COMMIT_TYPES",
		},
		cli.BoolTFlag{
			Name: "initial-development",
			Usage: `set this to false when you're ready to release 1.0.0,
                          ignored if version is already >= 1.0.0 (default: true)`,
			EnvVar: "GSG_INITIAL_DEVELOPMENT",
		},
		cli.BoolFlag{
			Name: "bump-patch",
			Usage: `Force patch bump when when none of the commits would otherwise trigger a bump.
                 First lines of commits will be listed under 'Fixes' in changelog and release note.
                 If there are commits that trigger a bump, this flag is ignored.`,
			EnvVar: "GSG_BUMP_PATCH",
		},
		cli.StringFlag{
			Name: "release-branches",
			Usage: `Comma separated list of branch names. If release-branches is defined, normal
                             releases can be done from listed brances. Other brances will produce pre-release
                             versions. If release-brances is not defined or is empty string, all branches will
                             produce normal releases. WARNING: this is an experimental feature.`,
			EnvVar: "GSG_RELEASE_BRANCHES",
		},
		cli.StringFlag{
			Name:   "tag-prefix",
			Value:  "v",
			Usage:  "`Prefix` to use in version tags.",
			EnvVar: "GSG_TAG_PREFIX",
		},
		cli.StringFlag{
			Name:   "bump-commit-tmpl",
			Value:  "chore: version bump for {{tag}} [skip ci]",
			Usage:  "`template` for bump commit message",
			EnvVar: "GSG_BUMP_COMMIT_TMPL",
		},
		cli.StringFlag{
			Name:   "pre-tmpl",
			Usage:  `Pre-release version template. Comma separated list of ID templates.`,
			EnvVar: "GSG_PRE_TMPL",
		},
		cli.StringFlag{
			Name:   "build-tmpl",
			Usage:  `Build metadata template. Comma separated list of ID templates.`,
			EnvVar: "GSG_BUILD_TMPL",
		},
		cli.StringFlag{
			Name:   "ci-project-path",
			Usage:  "gitlab CI environment variable",
			Hidden: false,
			EnvVar: "CI_PROJECT_PATH",
		},
		cli.StringFlag{
			Name:   "ci-commit-sha",
			Usage:  "gitlab CI environment variable",
			Hidden: false,
			EnvVar: "CI_COMMIT_SHA",
		},
		cli.StringFlag{
			Name:   "ci-commit-tag",
			Usage:  "gitlab CI environment variable",
			Hidden: false,
			EnvVar: "CI_COMMIT_TAG",
		},
		cli.StringFlag{
			Name:   "ci-commit-ref-name",
			Usage:  "gitlab CI environment variable",
			Hidden: false,
			EnvVar: "CI_COMMIT_REF_NAME",
		},
		cli.StringFlag{
			Name:   "ci-project-url",
			Usage:  "gitlab CI environment variable",
			Hidden: false,
			EnvVar: "CI_PROJECT_URL",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:      "next-version",
			Usage:     "Print next version",
			UsageText: "release next-version [command options]",
			Description: `Analyze commits and print the next version. See 'release help tag'
   for more details on how the analysis works.

   The default is to fail, if no commits indicating a version bump are found.
   Use --bump-patch or --allow-current to alter default behaviour.

   --bump-patch and --allow-current are mutually exclusive. On pre-releases they
   don't affect the version calculation, instead the pre and build templates are
   applied to increment version. 
   
   If global --bump-patch is set, it will be applied before next-version options are evaluated.`,
			Action: func(c *cli.Context) error {
				return inspectAndPrint(c, true, false)
			},
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "bump-patch, p",
					Usage: "Bump patch number if no changes are found in log",
				},
				cli.BoolFlag{
					Name:  "allow-current, c",
					Usage: "Print current version number if no changes are found in log",
				},
			},
		},
		{
			Name:  "changelog",
			Usage: "Update changelog",
			Description: `Analyze commits and create or update changelog.
         
   HEAD's parents are traversed and changes are collected from commits
   that haven't been released yet.

   First encountered tag from each branch is compared and semantically
   latest is selected to be the base for calculating the next version.

   GitLab api and environment variables are not needed for this command.

   If commit messages indicating a version bump are not found, the command
   exits with non-zero value.`,
			Action: changelog,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "f",
					Usage: "Write changelog to `FILE` (default is CHANGELOG.md)",
				},
				listOtherChanges,
			},
		},
		{
			Name:  "tag",
			Usage: "Create tag and attach release note to HEAD",
			Description: `Analyze commits and create tag and release note.
         
   HEAD's parents are traversed and changes are collected from commits that
   haven't been released yet.

   First encountered tag from each branch is compared and semantically
   latest is selected to be the base for calculating the next version.

   Tag and release note are created using the collected information.

   If commit messages indicating a version bump are not found, the command
   exits with non-zero value.`,
			UsageText: "release tag",
			Action:    tag,
			Flags:     []cli.Flag{listOtherChanges},
		},
		{
			Name:  "commit-and-tag",
			Usage: "Commit files and tag the new commit",
			Description: `Commit and push the listed files.
   If the files don't contain any changes, command fails.

   The default commit message template contains [skip ci],
   to prevent commit pipeline from running. You can override
   default template using global option --bump-commit-tmpl
   or environment variable GSG_BUMP_COMMIT_TMPL,

   Tag and release note are created for the new commit (check
   'release help tag' for more details).

   If commit messages indicating a version bump are not found,
   the command exits with non-zero value.`,
			UsageText: "release commit-and-tag [files to commit]",
			Action: func(c *cli.Context) error {
				return commitAndTagBase(c, false)
			},
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "create-tag-pipeline",
					Usage: `Needed when tagged commit message has
	[skip ci] and you want to execute tag
	pipeline.`,
				},
				listOtherChanges,
			},
		},
		{
			Name:  "tag-and-commit",
			Usage: "Tag HEAD and commit listed files",
			Description: `Attach tag and release note to HEAD (check
   'release help tag' for more details) and commit the listed files. 

   If the files don't contain any changes or commit messages indicating a
   version bump are not found, the command exits with non-zero value.

   The commit message contains [skip ci], to prevent commit pipeline from
   running.`,
			UsageText: "release tag-and-commit [files to commit]",
			Action: func(c *cli.Context) error {
				return commitAndTagBase(c, true)
			},
			Flags: []cli.Flag{listOtherChanges},
		},
		{
			Name:      "add-download",
			Usage:     "Add download to releasenote",
			UsageText: "release add-download [command options]",
			Description: `Upload file to project uploads and add link to release note.
            
   Requires CI_COMMIT_TAG environment variable or --ci-commit-tag flag.`,
			Action: addDownload,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "file, f",
					Usage: "`FILE` to upload",
				},
				cli.StringFlag{
					Name:  "ci-commit-tag",
					Usage: "`TAG` to add download to",
				},
			},
		},
	}
	app.Run(os.Args)
}
