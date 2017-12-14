package smsc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var (
	id      = 1000
	balance = "10.0"
	cost    = "1.5"
)

var ResultStringTests = []struct {
	Result Result
	S      string
}{
	{
		Result{ID: 100, Count: 1, Cost: &cost, Balance: &balance},
		"OK - 1 SMS, ID - 100",
	},
}

func TestResult_String(t *testing.T) {
	for _, tt := range ResultStringTests {
		if s := fmt.Sprintf("%v", &tt.Result); s != tt.S {
			t.Errorf("want %q, got %q", tt.S, s)
		}
	}
}

var ErrorErrorTests = []struct {
	Err Error
	S   string
}{
	{
		Error{Code: 2, Desc: "authorise error"},
		"ERROR = 2 (authorise error)",
	},
	{
		Error{Code: 7, Desc: "invalid number", ID: &id},
		"ERROR = 7 (invalid number), ID - 1000",
	},
}

func TestError_Error(t *testing.T) {
	for _, tt := range ErrorErrorTests {
		var err error = &tt.Err
		if s := fmt.Sprintf("%s", err); s != tt.S {
			t.Errorf("want %q, got %q", tt.S, s)
		}
	}
}

const (
	pass     = "pass"
	passHash = "1a1dc91c907325c69271ddf0c944bc72"
)

var NewTests = []struct {
	Config Config
	Err    error
}{
	{
		Config{Login: "test", Password: "pass"},
		nil,
	},
	{
		Config{Login: "test", PasswordMD5: "pass"},
		nil,
	},
	{
		Config{Login: "test"},
		ErrNoLoginPassword,
	},
	{
		Config{Password: "test"},
		ErrNoLoginPassword,
	},
}

func TestNew(t *testing.T) {
	for _, tt := range NewTests {
		if _, err := New(tt.Config); err != tt.Err {
			t.Errorf("want %v, got %v", tt.Err, err)
		}
	}
}

var ClientPrepareTests = []struct {
	Name    string
	Config  Config
	Text    string
	Phones  []string
	Opts    []Opt
	Message message
}{
	{
		Name: "Use a hash from a password",
		Config: Config{
			Login:    "me",
			Password: pass,
			Opt:      nil,
		},
		Text:   "test",
		Phones: []string{"123"},
		Opts:   nil,
		Message: message{
			Login:    "me",
			Password: passHash,
			Text:     "test",
			Phones:   []string{"123"},
			Charset:  charsetUTF8,
			Format:   formatJSON,
		},
	},
	{
		Name: "Use a hashed password",
		Config: Config{
			Login:       "me",
			PasswordMD5: passHash,
			Opt:         nil,
		},
		Text:   "test",
		Phones: []string{"123"},
		Opts:   nil,
		Message: message{
			Login:    "me",
			Password: passHash,
			Text:     "test",
			Phones:   []string{"123"},
			Charset:  charsetUTF8,
			Format:   formatJSON,
		},
	},
	{
		Name: "Client options are applied to a message",
		Config: Config{
			Login:    "me",
			Password: pass,
			Opt:      With(CostCountBalance),
		},
		Text:   "test",
		Phones: []string{"123"},
		Opts:   nil,
		Message: message{
			Login:    "me",
			Password: passHash,
			Text:     "test",
			Phones:   []string{"123"},
			Charset:  charsetUTF8,
			Format:   formatJSON,
			Cost:     CostCountBalance,
		},
	},
	{
		Name: "Send options overrides Client's one",
		Config: Config{
			Login:    "me",
			Password: pass,
			Opt:      With(CostCountBalance),
		},
		Text:   "test",
		Phones: []string{"123"},
		Opts:   []Opt{With(CostWithoutSend, Valid(0, 1))},
		Message: message{
			Login:    "me",
			Password: passHash,
			Text:     "test",
			Phones:   []string{"123"},
			Charset:  charsetUTF8,
			Format:   formatJSON,
			Cost:     CostWithoutSend,
			Valid:    &valid{0, 1},
		},
	},
}

func TestClient_prepare(t *testing.T) {
	for _, tt := range ClientPrepareTests {
		t.Run(tt.Name, func(t *testing.T) {
			c, _ := New(tt.Config)
			m := c.prepare(tt.Text, tt.Phones, tt.Opts)
			if !reflect.DeepEqual(m, &tt.Message) {
				t.Errorf("want %v, got %v", tt.Message, *m)
			}
		})
	}
}

// func TestClient_Send_fillsRequestWithRequiredParameters(t *testing.T) {
// 	login, password := "test", "pass"
// 	phone := "+71234567890"

// 	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if m := http.MethodPost; r.Method != m {
// 			t.Fatalf("method: want %v, got %v", m, r.Method)
// 		}
// 		if s := r.PostFormValue("login"); s != login {
// 			t.Errorf("login: want %q, got %q", login, s)
// 		}
// 		if s := r.PostFormValue("psw"); s != password {
// 			t.Errorf("psw: want %q, got %q", password, s)
// 		}
// 		if n := len(r.PostForm["phones"]); n != 1 {
// 			t.Fatalf("phones: want 1, got %d", n)
// 		}
// 		if s := r.PostFormValue("phones"); s != phone {
// 			t.Errorf("phones: want %q, got %q", phone, s)
// 		}

// 		n, err := strconv.ParseInt(r.PostFormValue("fmt"), 10, 32)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		if v := format(n); v != formatJSON {
// 			t.Errorf("fmt: want %v, got %v", formatJSON, v)
// 		}

// 		json.NewEncoder(w).Encode(&Result{Count: 1})
// 	}))
// 	defer ts.Close()

// 	c, err := New(Config{
// 		URL:      ts.URL,
// 		Login:    login,
// 		Password: password,
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	r, err := c.Send("A test message.", []string{phone})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if r.Count != 1 {
// 		t.Fatalf("count: want 1, got %d", r.Count)
// 	}
// }

var ClientSendTests = []struct {
	Value  interface{}
	Result *Result
	Err    error
}{
	{
		Value:  &Result{ID: 0, Count: 1},
		Result: &Result{ID: 0, Count: 1},
	},
	{
		Value: &Error{Code: 2},
		Err:   &Error{Code: 2},
	},
}

func TestClient_Send(t *testing.T) {
	minFields := 4 // login, psw, mes, phones

	for _, tt := range ClientSendTests {
		t.Run(fmt.Sprintf("%v", tt.Value), func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if m := http.MethodPost; r.Method != m {
					t.Fatalf("method: want %v, got %v", m, r.Method)
				}
				if err := r.ParseForm(); err != nil {
					t.Fatal(err)
				}
				if n := len(r.PostForm); n < minFields {
					t.Fatalf("form: want %v at least, got %v", minFields, n)
				}
				json.NewEncoder(w).Encode(tt.Value)
			}))
			defer ts.Close()

			c, err := New(Config{URL: ts.URL, Login: "test", Password: "pass"})
			if err != nil {
				t.Fatal(err)
			}

			r, err := c.Send("A test message.", []string{"+71234567890"})
			if !reflect.DeepEqual(tt.Err, err) {
				t.Errorf("error: want %v, got %v", tt.Err, err)
			}
			if !reflect.DeepEqual(tt.Result, r) {
				t.Errorf("result: want %v, got %v", tt.Result, r)
			}
		})
	}
}
