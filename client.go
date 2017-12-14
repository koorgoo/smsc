package smsc

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const DefaultURL = "https://smsc.ru/sys/send.php"

var ErrNoLoginPassword = errors.New("smsc: empty login or password")

// Config is a Client config.
type Config struct {
	URL         string
	Login       string
	Password    string
	PasswordMD5 string
	Opt         Opt
	Client      *http.Client
}

// New initializes a Client.
func New(cfg Config) (*Client, error) {
	if cfg.URL == "" {
		cfg.URL = DefaultURL
	}
	if cfg.Login == "" {
		return nil, ErrNoLoginPassword
	}
	if cfg.Password == "" && cfg.PasswordMD5 == "" {
		return nil, ErrNoLoginPassword
	}
	if cfg.PasswordMD5 == "" {
		cfg.PasswordMD5 = hashPassword(cfg.Password)
	}
	if cfg.Client == nil {
		cfg.Client = http.DefaultClient
	}
	c := &Client{
		url:      cfg.URL,
		login:    cfg.Login,
		password: cfg.PasswordMD5,
		opt:      cfg.Opt,
		http:     cfg.Client,
	}
	return c, nil
}

// hashPassword returns a hash from s.
//
// It's always used to hash a password from Config.
func hashPassword(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

// Client has methods for API calls.
type Client struct {
	url      string
	login    string
	password string
	opt      Opt
	http     *http.Client
}

func (c *Client) Send(text string, phones []string, opts ...Opt) (*Result, error) {
	m := c.prepare(text, phones, opts)
	if err := m.Validate(); err != nil {
		return nil, err
	}

	resp, err := c.http.PostForm(c.url, m.Values())
	if err != nil {
		return nil, wrapErr(err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, wrapErr(err)
	}

	// MetaResult allows to parse Result and Error from response while
	// structures have different fields.
	type MetaResult struct {
		*Result
		*Error
	}

	var mr *MetaResult
	if err := json.Unmarshal(b, &mr); err != nil {
		return nil, wrapErr(err)
	}
	if mr.Error != nil {
		return nil, mr.Error
	}
	return mr.Result, nil
}

// prepare returns a message ready to be sent.
func (c *Client) prepare(text string, phones []string, opts []Opt) *message {
	m := &message{
		Login:    c.login,
		Password: c.password,
		Text:     text,
		Phones:   phones,
		Charset:  charsetUTF8,
		Format:   formatJSON,
	}
	if c.opt != nil {
		c.opt(m)
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// wrapErr add smsc package prefix to error.
func wrapErr(err error) error {
	if err == nil {
		return err
	}
	return fmt.Errorf("smsc: %s", err)
}

type Result struct {
	ID      int     `json:"id"`
	Count   int     `json:"cnt"`
	Cost    *string `json:"cost"`
	Balance *string `json:"balance"`
	Phones  []Phone `json:"phones"`
}

func (r *Result) String() string {
	return fmt.Sprintf("OK - %d SMS, ID - %d", r.Count, r.ID)
}

type Phone struct {
	Phone  string  `json:"phone"`
	Mccmnc string  `json:"mccmnc"`
	Cost   string  `json:"cost"`
	Status *string `json:"status"`
	Error  *string `json:"error"`
}

type Error struct {
	Code int    `json:"error_code"`
	Desc string `json:"error"`
	ID   *int   `json:"id"`
}

func (e *Error) Error() string {
	s := fmt.Sprintf("ERROR = %d (%s)", e.Code, e.Desc)
	if e.ID != nil {
		s += fmt.Sprintf(", ID - %d", *e.ID)
	}
	return s
}
