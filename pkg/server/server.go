package server

import (
	"bytes"
	"context"
	iofs "io/fs"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/go-logr/logr"

	k8sversion "github.com/phillebaba/which-k8s-version/pkg/k8s-version"
)

type FrontendServer struct {
	logger    logr.Logger
	port      string
	styleCss  string
	indexHtml string
}

func NewFrontendServer(logger logr.Logger, fs iofs.FS, port string, vss []k8sversion.VersionSource) (*FrontendServer, error) {
	styleCssData, err := iofs.ReadFile(fs, "static/style.css")
	if err != nil {
		return nil, err
	}
	indexData, err := iofs.ReadFile(fs, "static/index.html.tpl")
	if err != nil {
		return nil, err
	}
	indexTmpl, err := template.New("index").Parse(string(indexData))
	if err != nil {
		return nil, err
	}
	pd := struct {
		VersionSources []k8sversion.VersionSource
	}{
		VersionSources: vss,
	}
	var indexHtml bytes.Buffer
	err = indexTmpl.Execute(&indexHtml, pd)
	if err != nil {
		return nil, err
	}

	return &FrontendServer{
		logger:    logger,
		port:      port,
		styleCss:  string(styleCssData),
		indexHtml: indexHtml.String(),
	}, nil
}

func (s *FrontendServer) ListenAndServe(stopCh <-chan struct{}) {
	mux := http.NewServeMux()
	mux.HandleFunc("/style.css", s.handleStyleCss())
	mux.HandleFunc("/", s.handleIndex())

	srv := &http.Server{
		Addr:    s.port,
		Handler: mux,
	}
	s.logger.Info("Frontend server starting")
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			s.logger.Error(err, "Frontend server crashed")
			os.Exit(1)
		}
	}()

	<-stopCh
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		s.logger.Error(err, "Frontend server graceful shutdown failed")
	} else {
		s.logger.Info("Frontend server stopped")
	}
}

func (s *FrontendServer) handleIndex() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(s.indexHtml))
	}
}

func (s *FrontendServer) handleStyleCss() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		w.Write([]byte(s.styleCss))
	}
}
