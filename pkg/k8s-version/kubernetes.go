package k8sversion

import (
	"context"

	"github.com/google/go-github/v35/github"
)

type Kubernetes struct {
	client *github.Client
}

func NewKubernetes() *Kubernetes {
	client := github.NewClient(nil)
	return &Kubernetes{
		client: client,
	}
}

func (k *Kubernetes) GetLatestVersion(ctx context.Context) (string, error) {
	release, _, err := k.client.Repositories.GetLatestRelease(ctx, "kubernetes", "kubernetes")
	if err != nil {
		return "", err
	}
	return *release.TagName, nil
}

func (k *Kubernetes) GetName() string {
	return "Kubernetes"
}

func (k *Kubernetes) GetColor() string {
	return "#326ce5"
}
