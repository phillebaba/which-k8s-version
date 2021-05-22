package k8sversion

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

type orchestratorVersionProfileListResult struct {
	Orchestrators []orchestratorVersionProfile `json:"properties.orchestrators"`
}

type orchestratorVersionProfile struct {
	IsPreview           bool   `j́son:"isPreview"`
	Default             bool   `j́son:"default"`
	OrchestratorType    string `j́son:"orchestratorType"`
	OrchestratorVersion string `j́son:"orchestratorVersion"`
}

type Azure struct {
	subscriptionID string
	location       string
}

func NewAzure(subscriptionID string) (*Azure, error) {
	if subscriptionID == "" {
		return nil, fmt.Errorf("subscriptionID can't be empty")
	}
	return &Azure{
		subscriptionID: subscriptionID,
		// TODO (Philip): Should query multiple regions and take the latest
		location: "West%20Europe",
	}, nil
}

func (a *Azure) GetLatestVersion(ctx context.Context) (string, error) {
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("https://management.azure.com/subscriptions/%v/providers/Microsoft.ContainerService/locations/%v/orchestrators?api-version=2019-08-01&resource-type=kubernetes", a.subscriptionID, a.location)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req, err = autorest.Prepare(req, authorizer.WithAuthorization())
	if err != nil {
		return "", err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return "", fmt.Errorf("invalid status code: %v", res.StatusCode)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	listResult := &orchestratorVersionProfileListResult{}
	err = json.Unmarshal(body, listResult)
	if err != nil {
		return "", err
	}
	if len(listResult.Orchestrators) == 0 {
		return "", fmt.Errorf("empty orchestrator list response")
	}
	latest := listResult.Orchestrators[len(listResult.Orchestrators)-1]
	return latest.OrchestratorVersion, nil
}

func (*Azure) GetName() string {
	return "Azure"
}

func (*Azure) GetColor() string {
	return "#008ad7"
}
