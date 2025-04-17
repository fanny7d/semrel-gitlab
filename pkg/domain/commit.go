package domain

// CommitType 表示提交类型
type CommitType string

const (
	TypeFix      CommitType = "fix"
	TypeFeat     CommitType = "feat"
	TypeRefactor CommitType = "refactor"
	TypePerf     CommitType = "perf"
	TypeDocs     CommitType = "docs"
	TypeStyle    CommitType = "style"
	TypeTest     CommitType = "test"
	TypeChore    CommitType = "chore"
)

// Commit 表示一个提交
type Commit struct {
	Hash       string
	Type       CommitType
	Scope      string
	Subject    string
	Body       string
	Breaking   bool
	PreRelease bool
}

// NewCommit 创建一个新的提交对象
func NewCommit(hash string, commitType CommitType, scope, subject, body string, breaking bool) *Commit {
	return &Commit{
		Hash:     hash,
		Type:     commitType,
		Scope:    scope,
		Subject:  subject,
		Body:     body,
		Breaking: breaking,
	}
}

// DetermineLevel 根据提交类型和是否破坏性变更确定版本升级级别
func (c *Commit) DetermineLevel(patchTypes, minorTypes []string) BumpLevel {
	if c.Breaking {
		return BumpMajor
	}

	for _, t := range minorTypes {
		if string(c.Type) == t {
			return BumpMinor
		}
	}

	for _, t := range patchTypes {
		if string(c.Type) == t {
			return BumpPatch
		}
	}

	return NoBump
}

// IsPreReleased 判断是否为预发布版本的提交
func (c *Commit) IsPreReleased() bool {
	return c.PreRelease
}

// SetPreReleased 设置提交为预发布版本
func (c *Commit) SetPreReleased(preRelease bool) {
	c.PreRelease = preRelease
}
