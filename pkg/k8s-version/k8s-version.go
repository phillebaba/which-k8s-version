package k8sversion

import (
	"context"
)

type VersionGetter interface {
	GetLatestVersion(ctx context.Context) (string, error)
	GetName() string
	GetColor() string
}

type VersionSource struct {
	Name    string
	Version string
	Color   string
}

func GetVersions(ctx context.Context, subscriptionID, projectID string) ([]VersionSource, error) {
	kubernetes := NewKubernetes()
	gcp, err := NewGCP(projectID)
	if err != nil {
		return []VersionSource{}, err
	}
	aws := NewAws()
	azure, err := NewAzure(subscriptionID)
	if err != nil {
		return []VersionSource{}, err
	}

	vgs := []VersionGetter{kubernetes, gcp, aws, azure}
	vss := []VersionSource{}
	for _, vg := range vgs {
		version, err := vg.GetLatestVersion(ctx)
		if err != nil {
			return []VersionSource{}, err
		}
		vs := VersionSource{
			Name:    vg.GetName(),
			Color:   vg.GetColor(),
			Version: version,
		}
		vss = append(vss, vs)
	}
	return vss, nil
}
