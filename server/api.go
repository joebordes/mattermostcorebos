package server

import (
	"mattermostcorebos/middleware"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost/server/public/plugin"
)

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	router := mux.NewRouter()

	session, _ := p.API.GetSession(c.SessionId)
	auth := middleware.AuthenticationMiddleware{Session: session}

	// Protected router section
	protected := router.PathPrefix("/").Subrouter()
	protected.HandleFunc("/team/{team-name}/project/{name}/documents", p.GetDocumentsForProject).Methods(http.MethodGet)

	protected.HandleFunc("/team/{team-name}/project/{name}/task", p.CreateTaskForProject).Methods(http.MethodPost)

	protected.HandleFunc("/team/{team-name}/project/{name}/method/{method}/module/{module}/invoke", p.DoInvoke).Methods(http.MethodPost)

	const wikiPath = "/team/{team-name}/project/{name}/wiki"
	protected.HandleFunc(wikiPath, p.CreateWiki).Methods(http.MethodPost)
	protected.HandleFunc(wikiPath, p.UpdateWiki).Methods(http.MethodPut)
	protected.HandleFunc(wikiPath, p.GetWikies).Methods(http.MethodGet)

	// Public router section
	router.Path("/health").HandlerFunc(p.Health).Methods(http.MethodGet)
	router.Path("/key").HandlerFunc(p.DoKeyJob).Methods(http.MethodGet)
	router.Path("/postmessage").HandlerFunc(p.postMessage).Methods(http.MethodPost)
	router.Path("/syncuser").HandlerFunc(p.syncUserWithCoreBOS).Methods(http.MethodPost)

	protected.Use(auth.CheckAuthentication)
	router.ServeHTTP(w, r)
}
