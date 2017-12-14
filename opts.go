package smsc

import (
	"errors"
	"fmt"
)

var ErrBadValid = errors.New("smsc: invalid period for valid")

// Opt configures a send message and a Result.
type Opt func(*message)

// With returns an Opt from options values v.
func With(v ...interface{}) Opt {
	var opts []Opt
	for _, t := range v {
		switch t := t.(type) {
		case Cost:
			opts = append(opts, withCost(t))
		case OpOpt:
			opts = append(opts, withOp(t))
		case ErrOpt:
			opts = append(opts, withErr(t))
		case Opt:
			opts = append(opts, t)
		default:
			panic(fmt.Sprintf("%T", t))
		}
	}
	return func(m *message) {
		for _, opt := range opts {
			opt(m)
		}
	}
}

func withCost(o Cost) Opt  { return func(m *message) { m.Cost = o } }
func withOp(o OpOpt) Opt   { return func(m *message) { m.Op = o } }
func withErr(o ErrOpt) Opt { return func(m *message) { m.Err = o } }

// format defines API output format. JSON format is used because it is simpler
// to parse response from.
type format int

const (
	formatInlineVerbose format = iota
	formatInline
	formatXML
	formatJSON
)

// Cost defines whether API should send a Result with cost information.
type Cost int

const (
	CostWithoutSend Cost = iota + 1
	CostCount
	CostCountBalance
)

// charset defines message text encoding. utf-8 is used always.
type charset string

const (
	charsetWindows1251 charset = "windows-1251"
	charsetUTF8                = "utf-8"
	charsetKOI8R               = "koi8-r"
)

// OpOpt controls whether Response must contain information about all phone
// numbers.
type OpOpt int

const Op OpOpt = 1

// ErrOpt controls whether Response must include information about failed phone
// numbers.
type ErrOpt int

const Err ErrOpt = 1

func Valid(h, m int) Opt {
	if (h < 0 || m < 0) ||
		(h > 24 || m > 59) ||
		(h == 0 && m < 1) ||
		(h == 24 && m > 0) {
		panic(ErrBadValid)
	}
	return (&valid{h, m}).Apply
}

// valid defines how long an operator must try to send a message.
type valid struct {
	Hours   int
	Minutes int
}

func (v *valid) Apply(m *message) {
	m.Valid = v
}

func (v *valid) String() string {
	var h, m string
	if v.Hours < 9 {
		h = "0"
	}
	if v.Minutes < 9 {
		m = "0"
	}
	return fmt.Sprintf("%s%d:%s%d", h, v.Hours, m, v.Minutes)
}

// TODO: Add more options.
// ID, Sender, Translit, Subj, TinyURL, Time, Tz, Period, Freq, Flash, Bin,
// Push, HLR, Ping, MMS, Mail, Viber, FileURL, Call, Voice, List, MaxSMS,
// ImgCode, UserIP, PP.

// formatOpt retuns a string for v value of option.
func formatOpt(v interface{}) string {
	return fmt.Sprintf("%v", v)
}
