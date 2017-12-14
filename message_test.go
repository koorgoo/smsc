package smsc

import (
	"net/url"
	"reflect"
	"testing"
)

const (
	abc    = "abcdefghijklmnopqrstuvwxyz"
	abcRus = "абвгдеёжзийклмнопрстуфхцчшщъыьэюя"
)

// Generate returns a string of n length using an abc alphabet.
// XXX: Not a production implementation - it doesn't use bytes.Buffer & sync.Pool.
func Generate(abc string, n int) string {
	runes := []rune(abc)
	if n < len(runes) {
		return string(runes[:n])
	}
	s := string(runes)
	s += Generate(abc, n-len(runes))
	return s
}

var GenerateTests = []struct {
	ABC string
	N   int
	S   string
}{
	{"abc", 0, ""},
	{"abc", 2, "ab"},
	{"abc", 5, "abcab"},
}

func TestGenerate(t *testing.T) {
	for _, tt := range GenerateTests {
		if s := Generate(tt.ABC, tt.N); s != tt.S {
			t.Errorf("(%q, %d): want %q, got %q", tt.ABC, tt.N, tt.S, s)
		}
	}
}

var somePhone = []string{"+71234567890"}

var MessageValidateTests = []struct {
	Name    string
	Message message
	Err     error
}{
	{
		"Single-sms-length text in latin is ok",
		message{Text: Generate(abc, 160), Phones: somePhone},
		nil,
	},
	{
		"Single-sms-length text in cyrillic is ok",
		message{Text: Generate(abcRus, 70), Phones: somePhone},
		nil,
	},
	{
		"765 chars in latin is maximum",
		message{Text: Generate(abc, smsMaxSize-smsMax*smsHeaderSize), Phones: somePhone},
		nil,
	},
	{
		"766 chars in latin is too long already",
		message{Text: Generate(abc, smsMaxSize-smsMax*smsHeaderSize+1), Phones: somePhone},
		ErrLongText,
	},
	{
		"335 chars in cyrillic is maximum",
		message{Text: Generate(abcRus, 335), Phones: somePhone},
		nil,
	},
	{
		"336 chars in cyrillic is too long already",
		message{Text: Generate(abcRus, 336), Phones: somePhone},
		nil,
	},
	{
		"At least one phone is required",
		message{Phones: []string{}, Text: "test"},
		ErrNoPhones,
	},
}

func TestMessage_Validate(t *testing.T) {
	for _, tt := range MessageValidateTests {
		t.Run(tt.Name, func(t *testing.T) {
			if err := tt.Message.Validate(); err != tt.Err {
				t.Errorf("want %v, got %v", tt.Err, err)
			}
		})
	}
}

var MessageValuesTests = []struct {
	Message message
	Values  url.Values
}{
	{
		message{
			Login:    "me",
			Password: "pass",
			Text:     "test",
			Phones:   []string{"1234"},
		},
		url.Values{
			"login":  []string{"me"},
			"psw":    []string{"pass"},
			"mes":    []string{"test"},
			"phones": []string{"1234"},
		},
	},
	{
		message{
			Charset: charsetUTF8,
			Format:  formatJSON,
			Cost:    CostCountBalance,
			Op:      Op,
			Err:     Err,
		},
		url.Values{
			"login":   []string{""},
			"psw":     []string{""},
			"mes":     []string{""},
			"phones":  []string{},
			"charset": []string{charsetUTF8},
			"fmt":     []string{formatOpt(formatJSON)},
			"cost":    []string{formatOpt(CostCountBalance)},
			"op":      []string{formatOpt(Op)},
			"err":     []string{formatOpt(Err)},
		},
	},
}

func TestMessage_Values(t *testing.T) {
	for _, tt := range MessageValuesTests {
		if v := tt.Message.Values(); !reflect.DeepEqual(v, tt.Values) {
			t.Errorf("want %v, got %v", tt.Values, v)
		}
	}
}
