package server

import (
	"awsx-api/config"
	"awsx-api/log"
	"awsx-api/routing"
	"awsx-api/util"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// func HandleRequests() {
// 	// http.HandleFunc("/grafana-dashboard", getDs)
// 	// log.Fatal(http.ListenAndServe(":10000", nil))

// 	// creates a new instance of a mux router
// 	myRouter := mux.NewRouter().StrictSlash(true)
// 	// replace http.HandleFunc with myRouter.HandleFunc
// 	myRouter.HandleFunc("/grafana-ds", handlers.GetDs)
// 	myRouter.HandleFunc("/grafana-ds-query", handlers.GrafanaQueryHandler)
// 	myRouter.HandleFunc("/grafana-ds/query-range", handlers.GrafanaQueryRangeHandler)

// 	// finally, instead of passing in nil, we want
// 	// to pass in our newly created router as the second
// 	// argument
// 	log.Fatal(http.ListenAndServe(":10000", myRouter))
// }

type Server struct {
	httpServer *http.Server
	router     *mux.Router
	// tracer     *sdktrace.TracerProvider
}

func NewServer() *Server {
	conf := config.Get()

	// create a router that will route all incoming API server requests to different handlers
	router := routing.NewRouter()
	// var tracingProvider *sdktrace.TracerProvider
	// if conf.Server.Observability.Tracing.Enabled {
	// 	// log.Infof("Tracing Enabled. Initializing tracer with collector url: %s", conf.Server.Observability.Tracing.CollectorURL)
	// 	tracingProvider = observability.InitTracer(conf.Server.Observability.Tracing.CollectorURL)
	// }

	middlewares := []mux.MiddlewareFunc{}
	if conf.Server.CORSAllowAll {
		middlewares = append(middlewares, corsAllowed)
	}
	// if conf.Server.Observability.Tracing.Enabled {
	// 	middlewares = append(middlewares, otelmux.Middleware(observability.TracingService))
	// }

	router.Use(middlewares...)

	handler := http.Handler(router)
	// if conf.Server.GzipEnabled {
	// 	handler = configureGzipHandler(router)
	// }

	// The Kiali server has only a single http server ever during its lifetime. But to support
	// testing that wants to start multiple servers over the lifetime of the process,
	// we need to override the default server mux with a new one everytime.
	mux := http.NewServeMux()
	http.DefaultServeMux = mux
	http.Handle("/", handler)
	http.Handle("/management/prometheus", promhttp.Handler())

	// Clients must use TLS 1.2 or higher
	// tlsConfig := &tls.Config{
	// 	MinVersion: tls.VersionTLS12,
	// }

	// create the server definition that will handle both console and api server traffic
	httpServer := &http.Server{
		// Addr:         fmt.Sprintf("%v:%v", conf.Server.Address, conf.Server.Port),
		Addr: fmt.Sprintf("%v:%v", conf.Server.Address, conf.Server.Port),
		// TLSConfig:    tlsConfig,
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
	}

	// return our new Server
	s := &Server{
		httpServer: httpServer,
		router:     router,
	}
	// if conf.Server.Observability.Tracing.Enabled && tracingProvider != nil {
	// 	s.tracer = tracingProvider
	// }
	return s
}

// Start HTTP server asynchronously. TLS may be active depending on the global configuration.
func (s *Server) Start() {
	// Start the business to initialize cache dependencies.
	// The business cache should start before the server endpoint to ensure
	// that the cache is ready before it's used by one of the server handlers.
	// business.Start()

	// conf := config.Get()
	//log.Infof("Server endpoint will start at [%v%v]", s.httpServer.Addr, conf.Server.WebRoot)
	log.Infof("Server endpoint will start at [%v%v]", s.httpServer.Addr, "/")
	// log.Infof("Server endpoint will serve static content from [%v]", conf.Server.StaticContentRootDirectory)
	// secure := conf.Identity.CertFile != "" && conf.Identity.PrivateKeyFile != ""
	go func() {
		var err error
		// if secure {
		// 	// log.Infof("Server endpoint will require https")
		// 	s.router.Use(secureHttpsMiddleware)
		// 	err = s.httpServer.ListenAndServeTLS(conf.Identity.CertFile, conf.Identity.PrivateKeyFile)
		// } else {
		s.router.Use(plainHttpMiddleware)
		err = s.httpServer.ListenAndServe()
		util.CommonError(err)
		// }
		// log.Warning(err)
	}()

	// Start the Metrics Server
	// if conf.Server.Observability.Metrics.Enabled {
	// 	StartMetricsServer()
	// }
}

// Stop the HTTP server
func (s *Server) Stop() {
	// StopMetricsServer()
	// business.Stop()
	// log.Infof("Server endpoint will stop at [%v]", s.httpServer.Addr)
	s.httpServer.Close()
	// observability.StopTracer(s.tracer)
}

func corsAllowed(next http.Handler) http.Handler {
	//conf := config.Get()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//if conf.Server.WhiteListUrls != "" {
		//	w.Header().Set("Access-Control-Allow-Origin", conf.Server.WhiteListUrls)
		//} else {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		//}
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
		next.ServeHTTP(w, r)

	})
}

// func configureGzipHandler(handler http.Handler) http.Handler {
// 	contentTypeOption := gziphandler.ContentTypes([]string{
// 		"application/javascript",
// 		"application/json",
// 		"image/svg+xml",
// 		"text/css",
// 		"text/html",
// 	})
// 	if handlerFunc, err := gziphandler.GzipHandlerWithOpts(contentTypeOption); err == nil {
// 		return handlerFunc(handler)
// 	} else {
// 		// This could happen by a wrong configuration being sent to GzipHandlerWithOpts
// 		panic(err)
// 	}
// }

func plainHttpMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Scheme = "http"
		next.ServeHTTP(w, r)
	})
}

// func secureHttpsMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		r.URL.Scheme = "https"
// 		next.ServeHTTP(w, r)
// 	})
// }
