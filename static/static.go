package static

import (
	"fmt"
	"llmProxy/gomie/config"
	"llmProxy/gomie/logger"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

var log *zap.SugaredLogger

const msgPad int = 20

func AddToRouter(router *mux.Router) {
	cfg := config.GetConfig()
	log = logger.GetLogger()
	action := "add route static"

	root_prefix := cfg.GetString("root_prefix")
	static_path := cfg.GetString("static_path")
	absPath, _ := filepath.Abs(static_path)

	msg := cfg.GetString("app_name") + fmt.Sprintf(": %-*v", msgPad, action) + root_prefix
	log.Infow(msg, "action", action, "route", root_prefix, "staticFolder", absPath)

	// static path is relative to the domain root
	addPrefixRoutes(router, static_path, root_prefix)
}

func addPrefixRoutes(router *mux.Router, static_path string, mountPrefix string) {

	// Choose the folder to serve
	staticPrefix := "/"
	if len(mountPrefix) > 0 {
		staticPrefix = mountPrefix
	}

	fileHandler := http.FileServer(http.Dir(static_path))
	wrappedHandler := loggingFileServer(fileHandler)

	//need to strip the staticPrefix for golang to parse the subtree correctly
	// router.Handle(staticPrefix, http.StripPrefix(staticPrefix, wrappedHandler))
	// router.Handle(staticPrefix, http.FileServer(http.Dir(static_path)))
	// router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("..\\dlt-af-product-definition\\dalet-production-folder"))))
	router.PathPrefix(staticPrefix).Handler(http.StripPrefix(staticPrefix, wrappedHandler))
}

// create a custom logging statcc file server
func loggingFileServer(h http.Handler) http.Handler {
	cfg := config.GetConfig()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Infow(cfg.GetString("app_name")+": static", "url", r.RequestURI)
		h.ServeHTTP(w, r)
	})
}
