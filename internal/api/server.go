package api

import (
	"fmt"
	"github.com/dsrosen6/cw-ticket-bot/internal/webex"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	defaultAddr    = ":8080"
	defaultTimeout = 60
	defaultDbName  = "data.db"

	envVarAddr           = "TICKET_BOT_SERVER_ADDR"
	envVarTimeout        = "TICKET_BOT_TIMEOUT_SECONDS"
	envVarWebexToken     = "TICKET_BOT_WEBEX_TOKEN"
	envVarWebexBotEmail  = "TICKET_BOT_WEBEX_BOT_EMAIL"
	envVarWebhookBaseUrl = "TICKET_BOT_WEBHOOK_BASE_URL"
)

type Server struct {
	httpClient     *http.Client
	webexClient    *webex.Client
	db             *DB
	addr           string
	webhookBaseUrl string
	timeout        int
	botEmail       string
}

func NewServer() (*Server, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("loading .env file: %w", err)
	}

	addr := getAddr(os.Getenv(envVarAddr))
	timeout := getTimeout(os.Getenv(envVarTimeout))
	webexToken := os.Getenv(envVarWebexToken)
	webexBotEmail := os.Getenv(envVarWebexBotEmail)
	webhookBaseUrl := os.Getenv(envVarWebhookBaseUrl)

	if webexToken == "" {
		return nil, fmt.Errorf("webex token is empty")
	}

	if webexBotEmail == "" {
		return nil, fmt.Errorf("webex bot email is empty")
	}

	if webhookBaseUrl == "" {
		return nil, fmt.Errorf("webhook base url is empty")
	}

	httpClient := http.DefaultClient
	db, err := newDB(defaultDbName)
	if err != nil {
		return nil, fmt.Errorf("opening db: %w", err)
	}

	return &Server{
		httpClient:     httpClient,
		webexClient:    webex.NewClient(httpClient, webexToken),
		db:             db,
		addr:           addr,
		webhookBaseUrl: webhookBaseUrl,
		timeout:        timeout,
		botEmail:       webexBotEmail,
	}, nil
}

func (s *Server) Run() error {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Duration(s.timeout) * time.Second))

	if err := s.db.InitSchema(); err != nil {
		return fmt.Errorf("initializing db schema: %w", err)
	}
	log.Println("db initialized")

	r.Mount("/ping", PingRouter())
	r.Mount("/directmsg", s.DirectMessageRouter())

	log.Println("listening at:", s.addr)
	if err := http.ListenAndServe(s.addr, r); err != nil {
		return fmt.Errorf("an error occured running the server: %w", err)
	}

	return nil
}

func getAddr(customAddr string) string {
	addr := defaultAddr
	if customAddr != "" {
		addr = customAddr
	}

	return addr
}

func getTimeout(customTimeout string) int {
	timeout := defaultTimeout
	if customTimeout != "" {
		var err error
		timeout, err = strconv.Atoi(customTimeout)
		if err != nil {
			log.Printf("error converting custom timeout of %s to integer - using default of %d", customTimeout, defaultTimeout)
			timeout = defaultTimeout
		}
	}

	return timeout
}
