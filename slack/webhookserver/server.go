package webhookserver

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/thatoddmailbox/aislack/config"
	"github.com/thatoddmailbox/aislack/slack"
)

type Server struct {
	port int
	c    *slack.Client
}

func NewServer(port int, c *slack.Client) (*Server, error) {
	return &Server{
		port: port,
		c:    c,
	}, nil
}

func validateRequest(r *http.Request, body []byte) bool {
	timestampString := r.Header.Get("X-Slack-Request-Timestamp")
	if timestampString == "" {
		return false
	}

	slackSignatureString := r.Header.Get("X-Slack-Signature")
	if slackSignatureString == "" {
		return false
	}

	timestamp, err := strconv.ParseInt(timestampString, 10, 64)
	if err != nil {
		return false
	}

	difference := time.Now().Unix() - timestamp
	if difference > timestamp {
		// too old
		log.Println("too old")
		return false
	}

	sigstring := "v0:" + timestampString + ":" + string(body)

	h := hmac.New(sha256.New, []byte(config.Current.Slack.SigningSecret))
	h.Write([]byte(sigstring))
	digestBytes := h.Sum(nil)
	digestHex := "v0=" + hex.EncodeToString(digestBytes)

	valid := hmac.Equal([]byte(digestHex), []byte(slackSignatureString))
	return valid
}

func (s *Server) slashCommandHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		// TODO: handle this??
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !validateRequest(r, body) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	values, err := url.ParseQuery(string(body))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := s.c.HandleSlashCommand(values)
	if err != nil {
		// TODO: handle this??
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		// TODO: handle this??
		log.Println(err)
		return
	}
}

func (s *Server) StartListening() {
	http.HandleFunc("/slack-slash-command", s.slashCommandHandler)

	err := http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
	if err != nil {
		panic(err)
	}
}
