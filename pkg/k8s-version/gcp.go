package k8sversion

import (
	"context"
	"fmt"

	container "cloud.google.com/go/container/apiv1"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
)

type GCP struct {
	projectID string
	location  string
}

func NewGCP(projectID string) (*GCP, error) {
	if projectID == "" {
		return nil, fmt.Errorf("projectID can't be empty")
	}
	return &GCP{
		projectID: projectID,
		location:  "us-central1",
	}, nil
}

func (g *GCP) GetLatestVersion(ctx context.Context) (string, error) {
	c, err := container.NewClusterManagerClient(ctx)
	if err != nil {
		return "", err
	}
	req := &containerpb.GetServerConfigRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s", g.projectID, g.location),
	}
	resp, err := c.GetServerConfig(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.DefaultClusterVersion, nil
}

func (*GCP) GetName() string {
	return "GCP"
}

func (*GCP) GetColor() string {
	return "#db4437"
}
