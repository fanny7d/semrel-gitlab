package render

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/blang/semver"
	"github.com/juranki/go-semrel/semrel"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

var (
	releaseNoteTmpl = `# {{ .NextVersion }}
{{ date }}{{ if .Changes.breaking }}

## Breaking changes{{ range .Changes.breaking }}

### {{ if ne "" .Scope }}**{{ .Scope }}:** {{ end}}{{ .Subject }} ({{ .Hash }})

{{ .BreakingMessage }}{{ end }}{{ end }}{{ if .Changes.feature }}

## Features
{{ range .Changes.feature }}
- {{ if ne "" .Scope }}**{{ .Scope }}:** {{ end}}{{ .Subject }} ({{ .Hash }}){{ end }}{{ end }}{{ if .Changes.fix }}

## Fixes
{{ range .Changes.fix }}
- {{ if ne "" .Scope }}**{{ .Scope }}:** {{ end}}{{ .Subject }} ({{ .Hash }}){{ end }}{{ end }}{{ if .Changes.other }}

## Other changes
{{ range .Changes.other }}
- {{ if ne "" .Scope }}**{{ .Scope }}:** {{ end}}{{ .Subject }} ({{ .Hash }}){{ end }}{{ end }}

<!--- downloads here -->`
	changelogTmpl = `## {{ .NextVersion }}
{{ date }}{{ if .Changes.breaking }}

### Breaking changes{{ range .Changes.breaking }}

#### {{ if ne "" .Scope }}**{{ .Scope }}:** {{ end}}{{ .Subject }} ({{ .Hash }})

{{ .BreakingMessage }}{{ end }}{{ end }}{{ if .Changes.feature }}

### Features
{{ range .Changes.feature }}
- {{ if ne "" .Scope }}**{{ .Scope }}:** {{ end}}{{ .Subject }} ({{ .Hash }}){{ end }}{{ end }}{{ if .Changes.fix }}

### Fixes
{{ range .Changes.fix }}
- {{ if ne "" .Scope }}**{{ .Scope }}:** {{ end}}{{ .Subject }} ({{ .Hash }}){{ end }}{{ end }}{{ if .Changes.other }}

### Other changes
{{ range .Changes.other }}
- {{ if ne "" .Scope }}**{{ .Scope }}:** {{ end}}{{ .Subject }} ({{ .Hash }}){{ end }}{{ end }}`
	funcs   = template.FuncMap{"date": func() string { return time.Now().Format("2006-01-02") }}
	preTmpl = []string{
		`{{ (env "CI_COMMIT_REF_SLUG") }}`,
		`{{ seq }}`,
	}
)

type preReleaseInfo struct {
	CommitTS time.Time
}

func render(releaseInfo interface{}, tmpl string, fns template.FuncMap) (string, error) {
	t := template.New("")
	if fns != nil {
		t = t.Funcs(fns)
	}
	t, err := t.Parse(tmpl)
	if err != nil {
		return "", err
	}
	buf := bytes.Buffer{}
	err = t.Execute(&buf, releaseInfo)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// ReleaseNote takes release info object returned by semrel.Release and returns markdown for the release note
func ReleaseNote(releaseInfo interface{}) (string, error) {
	return render(releaseInfo, releaseNoteTmpl, funcs)
}

// ChangelogEntry takes release info object returned by semrel.Release and returns markdown for the changelog entry
func ChangelogEntry(releaseInfo interface{}) (string, error) {
	return render(releaseInfo, changelogTmpl, funcs)
}

// BumpMessage renders tag using tmpl
func BumpMessage(tag string, tmpl string) (string, error) {
	return render(nil, tmpl, template.FuncMap{
		"tag": func() string { return tag },
	})
}

// SetPreReleaseNumber adds pre- and build parts to the next version of releaseData
func SetPreReleaseNumber(releaseData *semrel.ReleaseData, context *cli.Context) error {
	preTs := parseCommaSeparatedList(context.GlobalString("pre-tmpl"))
	buildTs := parseCommaSeparatedList(context.GlobalString("build-tmpl"))

	versions, err := getMatchingVersions(releaseData.NextVersion)
	if err != nil {
		return err
	}
	fns := template.FuncMap{
		"env":      os.Getenv,
		"commitTS": func() time.Time { return releaseData.Time.UTC() },
	}

	if len(preTs) == 0 {
		preTs = preTmpl
	}

	if len(buildTs) > 0 {
		buildIDs := make([]string, len(buildTs))
		for i, buildT := range buildTs {
			buildStr, err := render(nil, buildT, fns)
			if err != nil {
				return errors.Wrapf(err, "formatting %s", buildT)
			}
			buildID, err := semver.NewBuildVersion(buildStr)
			if err != nil {
				return errors.Wrap(err, "build version tmpl")
			}
			buildIDs[i] = buildID
		}
		releaseData.NextVersion.Build = buildIDs
	}

	if len(preTs) > 0 {
		preIDs := make([]semver.PRVersion, 0)
		for _, preT := range preTs {
			fns["seq"] = getVersionsFun(versions, preIDs)
			preStr, err := render(nil, preT, fns)
			if err != nil {
				return errors.Wrap(err, "pre version tmpl")
			}
			preID, err := semver.NewPRVersion(preStr)
			if err != nil {
				return errors.Wrap(err, "pre version tmpl")
			}
			preIDs = append(preIDs, preID)
		}
		releaseData.NextVersion.Pre = preIDs
	}
	return nil
}

func getMatchingVersions(v semver.Version) ([]semver.Version, error) {
	versions := make([]semver.Version, 0)

	r, err := git.PlainOpen(".")
	if err != nil {
		return nil, err
	}

	addIfSemVer := func(sha string, version string) {
		sv, err := semver.ParseTolerant(version)
		if err == nil {
			if v.Major == sv.Major && v.Minor == sv.Minor && v.Patch == sv.Patch {
				versions = append(versions, sv)
			}
		}
	}

	tagRefs, err := r.Tags()
	if err != nil {
		return nil, err
	}
	err = tagRefs.ForEach(func(t *plumbing.Reference) error {
		addIfSemVer(t.Hash().String(), t.Name().Short())
		return nil
	})
	if err != nil {
		return nil, err
	}

	tagObjects, err := r.TagObjects()
	if err != nil {
		return nil, err
	}
	err = tagObjects.ForEach(func(t *object.Tag) error {
		addIfSemVer(t.Target.String(), t.Name)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return versions, nil
}

func getVersionsFun(versions []semver.Version, preIDs []semver.PRVersion) func() (string, error) {
	f := func() (string, error) {
		n := uint64(0)
		if len(preIDs) == 0 {
			return "", errors.New("seq cannot be used in first ID of pre number")
		}
		for _, v := range versions {
			if len(v.Pre) > len(preIDs) {
				if !v.Pre[len(preIDs)].IsNumeric() {
					continue
				}
				if v.Pre[len(preIDs)].VersionNum <= n {
					continue
				}
				match := true
				for i, p := range preIDs {
					if p.Compare(v.Pre[i]) != 0 {
						match = false
						break
					}
				}
				if match {
					n = v.Pre[len(preIDs)].VersionNum
				}
			}
		}
		return fmt.Sprintf("%d", n+1), nil
	}
	return f
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
