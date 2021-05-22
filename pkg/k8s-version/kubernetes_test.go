package k8sversion

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBasic(t *testing.T) {
	k := NewKubernetes()
	v, err := k.GetLatestVersion(context.TODO())
	require.NoError(t, err)
	require.Equal(t, v, "v1.21.1")
}
