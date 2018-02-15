package smsc

import (
	"errors"
	"net/url"
	"unicode/utf8"
)

var (
	ErrLongText = errors.New("smsc: too long text to send")
	ErrNoPhones = errors.New("smsc: empty phones list")
)

// message controls what and how will sent to API.
type message struct {
	Login    string
	Password string
	Text     string
	Phones   []string
	Charset  charset
	Format   format
	Cost     Cost
	Op       OpOpt
	Err      ErrOpt
	Valid    *valid
	Sender   string
	Translit TranslitOpt
}

const (
	smsMax        = 5
	smsSize       = 160
	smsMaxSize    = smsMax * smsSize
	smsHeaderSize = 7
)

// Validate checks the message integrity and returns an optional error.
func (m *message) Validate() error {
	if n := CountBytes(m.Text); n > smsMaxSize {
		return ErrLongText
	}
	if len(m.Phones) == 0 {
		return ErrNoPhones
	}
	// TODO: Validate options when added.
	return nil
}

// CountBytes returns a number of bytes for text to be send in SMS.
func CountBytes(text string) int {
	if n := len(text); n < smsSize {
		return n
	}

	// In UTF-8, characters from the U+0000..U+10FFFF range (the UTF-16
	// accessible range) are encoded using sequences of 1 to 4 octets.
	// https://tools.ietf.org/html/rfc3629
	buf := make([]byte, 4)

	var size int
	var smsBytes int

	for _, r := range text {
		n := utf8.EncodeRune(buf, r)

		next := smsBytes + n + smsHeaderSize
		if next > smsSize {
			size += smsHeaderSize // count message header
			smsBytes = 0          // ... and start a new sms
		}

		size += n
		smsBytes += n
	}

	return size
}

// Values returns a form for a request to API.
func (m *message) Values() url.Values {
	v := url.Values{
		"login":  []string{m.Login},
		"psw":    []string{m.Password},
		"mes":    []string{m.Text},
		"phones": m.Phones,
	}

	if m.Charset != "" {
		v.Set("charset", formatOpt(m.Charset))
	}
	if m.Format != 0 {
		v.Set("fmt", formatOpt(m.Format))
	}
	if m.Cost != 0 {
		v.Set("cost", formatOpt(m.Cost))
	}
	if m.Op != 0 {
		v.Set("op", formatOpt(m.Op))
	}
	if m.Err != 0 {
		v.Set("err", formatOpt(m.Err))
	}
	if m.Valid != nil {
		v.Set("valid", formatOpt(m.Valid))
	}
	if m.Sender != "" {
		v.Set("sender", formatOpt(m.Sender))
	}
	if m.Translit != 0 {
		v.Set("translit", formatOpt(m.Translit))
	}
	return v
}
