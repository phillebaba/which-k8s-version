package k8sversion

import (
	"context"
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type Aws struct {
}

func NewAws() *Aws {
	return &Aws{}
}

func (a *Aws) GetLatestVersion(ctx context.Context) (string, error) {
	url := "https://docs.aws.amazon.com/eks/latest/userguide/kubernetes-versions.html"
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return "", fmt.Errorf("invalid status code: %v", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", err
	}

	version := doc.Find("div.itemizedlist:nth-child(7) > ul:nth-child(1) > li:nth-child(1) > p:nth-child(1)").Text()
	if version == "" {
		return "", fmt.Errorf("version not found in page")
	}
	return version, nil
}

func (*Aws) GetName() string {
	return "AWS"
}

func (*Aws) GetColor() string {
	return "#ff9900"
}
