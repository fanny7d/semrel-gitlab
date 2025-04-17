package domain

import (
	"time"

	"github.com/blang/semver"
)

// Version 表示语义化版本
type Version struct {
	Current semver.Version
	Next    semver.Version
	Level   BumpLevel
}

// BumpLevel 表示版本升级级别
type BumpLevel int

const (
	NoBump BumpLevel = iota
	BumpPatch
	BumpMinor
	BumpMajor
)

// NewVersion 创建一个新的版本对象
func NewVersion(t time.Time) *Version {
	return &Version{
		Current: semver.Version{},
		Next:    semver.Version{},
		Level:   NoBump,
	}
}

// Bump 根据指定的级别升级版本
func (v *Version) Bump(level BumpLevel) {
	if level > v.Level {
		v.Level = level
	}

	switch v.Level {
	case BumpMajor:
		v.Next.Major++
		v.Next.Minor = 0
		v.Next.Patch = 0
	case BumpMinor:
		v.Next.Minor++
		v.Next.Patch = 0
	case BumpPatch:
		v.Next.Patch++
	}
}

// SetPreRelease 设置预发布版本
func (v *Version) SetPreRelease(pre string) error {
	preRelease, err := semver.NewPRVersion(pre)
	if err != nil {
		return err
	}
	v.Next.Pre = []semver.PRVersion{preRelease}
	return nil
}

// SetBuild 设置构建元数据
func (v *Version) SetBuild(build string) {
	v.Next.Build = []string{build}
}
