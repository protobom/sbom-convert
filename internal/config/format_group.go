package config

import (
	translate "github.com/bom-squad/go-cli/pkg"
	"github.com/bom-squad/protobom/pkg/formats"
)

// Helper Object
type FormatGroup struct {
	Format           string
	URI              string
	Versions         map[string]interface{}
	Encodings        map[string]interface{}
	DefaultVersions  map[string]interface{}
	DefaultEncodings map[string]interface{}
}

type FormatGroups map[string]FormatGroup

func GetVersionGroups() FormatGroups {
	versionGroups := make(FormatGroups)

	for _, f := range formats.List {
		group := versionGroups[f.Type()]
		if group.URI != "" {
			group.URI = f.URI()
		}
		if group.Versions == nil {
			group.Versions = make(map[string]interface{})
		}
		if group.Encodings == nil {
			group.Encodings = make(map[string]interface{})
		}
		group.Versions[f.Version()] = ""
		group.Encodings[f.Encoding()] = ""
		versionGroups[f.Type()] = group
	}

	return versionGroups
}

func (v FormatGroups) Versions(format string) []string {
	var l []string
	for k, _ := range v[format].Versions {
		l = append(l, k)
	}
	return l
}

func (v FormatGroups) Encodings(format string) []string {
	var l []string
	for k, _ := range v[format].Encodings {
		l = append(l, k)
	}
	return l
}

func (v FormatGroups) EncodingMap() map[string][]string {
	m := make(map[string][]string)
	for k, _ := range v {
		m[k] = v.Encodings(k)
	}

	return m
}

func (v FormatGroups) URI(Type string) string {
	return v[Type].URI
}

func (v FormatGroups) VersionMap() map[string][]string {
	m := make(map[string][]string)
	for k, _ := range v {
		m[k] = v.Versions(k)
	}

	return m
}

func (v FormatGroups) DefaultVersions() map[string]string {
	m := make(map[string]string)
	for _, f := range translate.DefaultsFormatsList {
		m[f.Type()] = f.Version()
	}

	return m
}

func (v FormatGroups) Formats() []string {
	var keys []string
	for k, _ := range v {
		keys = append(keys, k)
	}

	return keys
}
