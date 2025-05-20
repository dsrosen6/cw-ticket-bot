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

	envVarAddr          = "TICKET_BOT_SERVER_ADDR"
	envVarTimeout       = "TICKET_BOT_TIMEOUT_SECONDS"
	envVarWebexToken    = "TICKET_BOT_WEBEX_TOKEN"
	envVarWebexBotEmail = "TICKET_BOT_WEBEX_BOT_EMAIL"
)

type Server struct {
	httpClient  *http.Client
	webexClient *webex.Client
	addr        string
	timeout     int
	botEmail    string
}

func NewServer() (*Server, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("loading .env file: %w", err)
	}

	addr := getAddr(os.Getenv(envVarAddr))
	timeout := getTimeout(os.Getenv(envVarTimeout))
	webexToken := os.Getenv(envVarWebexToken)
	webexBotEmail := os.Getenv(envVarWebexBotEmail)

	if webexToken == "" {
		return nil, fmt.Errorf("webex token is empty")
	}

	if webexBotEmail == "" {
		return nil, fmt.Errorf("webex bot email is empty")
	}

	httpClient := http.DefaultClient

	return &Server{
		httpClient:  httpClient,
		webexClient: webex.NewClient(httpClient, webexToken),
		addr:        addr,
		timeout:     timeout,
		botEmail:    webexBotEmail,
	}, nil
}

func (s *Server) Run() error {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Duration(s.timeout) * time.Second))

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
