package smsc

import (
	"errors"
	"net/url"
	"strconv"
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
	Format   format
	Cost     Cost
}

const (
	smsMax        = 5
	smsSize       = 160
	smsMaxSize    = smsMax * smsSize
	smsHeaderSize = 7
)

// Validate checks the message integrity and returns an optional error.
func (m *message) Validate() error {
	if n := m.countBytes(); n > smsMaxSize {
		return ErrLongText
	}
	if len(m.Phones) == 0 {
		return ErrNoPhones
	}
	// TODO: Validate options when added.
	return nil
}

// countBytes returns a rough number of bytes to send a text of m taking into
// account that a multi-sms message will contain headers.
func (m *message) countBytes() int {
	if n := len(m.Text); n < smsSize {
		return n
	}

	// In UTF-8, characters from the U+0000..U+10FFFF range (the UTF-16
	// accessible range) are encoded using sequences of 1 to 4 octets.
	// https://tools.ietf.org/html/rfc3629
	buf := make([]byte, 4)

	var size int
	var smsBytes int

	for _, r := range m.Text {
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
	if m.Format != 0 {
		v.Set("fmt", strconv.FormatInt(int64(m.Format), 10))
	}
	if m.Cost != 0 {
		v.Set("cost", strconv.FormatInt(int64(m.Cost), 10))
	}
	return v
}
