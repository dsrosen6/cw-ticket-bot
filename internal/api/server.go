package api

import (
	"fmt"
	"github.com/dsrosen6/cw-ticket-bot/internal/connectwise"
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
	envVarCwCompanyId    = "TICKET_BOT_CW_COMPANY_ID"
	envVarCwClientId     = "TICKET_BOT_CW_CLIENT_ID"
	envVarCwPublicKey    = "TICKET_BOT_CW_PUBLIC_KEY"
	envVarCwPrivateKey   = "TICKET_BOT_CW_PRIVATE_KEY"
)

type Server struct {
	httpClient     *http.Client
	webexClient    *webex.Client
	cwClient       *connectwise.Client
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
	cwCompanyId := os.Getenv(envVarCwCompanyId)
	cwClientId := os.Getenv(envVarCwClientId)
	cwPublicKey := os.Getenv(envVarCwPublicKey)
	cwPrivateKey := os.Getenv(envVarCwPrivateKey)

	if webexToken == "" || webexBotEmail == "" {
		return nil, fmt.Errorf("one or more webex environment variable is empty")
	}

	if webhookBaseUrl == "" {
		return nil, fmt.Errorf("webhook base url is empty")
	}

	if cwCompanyId == "" || cwClientId == "" || cwPublicKey == "" || cwPrivateKey == "" {
		return nil, fmt.Errorf("one or more connectwise environment variables are empty")
	}

	httpClient := http.DefaultClient
	db, err := newDB(defaultDbName)
	if err != nil {
		return nil, fmt.Errorf("opening db: %w", err)
	}

	return &Server{
		httpClient:     httpClient,
		webexClient:    webex.NewClient(httpClient, webexToken),
		cwClient:       connectwise.NewClient(httpClient, cwPublicKey, cwPrivateKey, cwClientId, cwCompanyId),
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
