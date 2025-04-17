package domain

// Release 表示一个发布
type Release struct {
	Version *Version
	Changes map[string][]*Commit
	TagName string
	Message string
	Links   []ReleaseLink
}

// ReleaseLink 表示发布中的下载链接
type ReleaseLink struct {
	Name        string
	URL         string
	Description string
}

// NewRelease 创建一个新的发布对象
func NewRelease(version *Version, tagPrefix string) *Release {
	return &Release{
		Version: version,
		Changes: make(map[string][]*Commit),
		TagName: tagPrefix + version.Next.String(),
	}
}

// AddChange 添加一个变更到发布中
func (r *Release) AddChange(category string, commit *Commit) {
	if _, ok := r.Changes[category]; !ok {
		r.Changes[category] = make([]*Commit, 0)
	}
	r.Changes[category] = append(r.Changes[category], commit)
}

// AddLink 添加一个下载链接到发布中
func (r *Release) AddLink(name, url, description string) {
	r.Links = append(r.Links, ReleaseLink{
		Name:        name,
		URL:         url,
		Description: description,
	})
}

// HasContent 判断发布是否有内容
func (r *Release) HasContent() bool {
	if r.Version.Level == NoBump {
		return false
	}

	if len(r.Version.Next.Pre)+len(r.Version.Next.Build) == 0 {
		return true
	}

	for category, changes := range r.Changes {
		if category == "other" {
			continue
		}
		for _, change := range changes {
			if !change.IsPreReleased() {
				return true
			}
		}
	}

	return false
}
