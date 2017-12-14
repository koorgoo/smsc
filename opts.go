package smsc

import (
	"fmt"
)

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

// TODO: Add more options.
// Major - ID, Sender, Translit, Subj.
// Minor - TinyURL, Time, Tz, Period, Freq, Flash, Bin, Push, HLR, Ping, MMS,
// Mail, Viber, FileURL, Call, Voice, List, Valid, MaxSMS, ImgCode, UserIP, PP.

// formatOpt retuns a string for v value of int-like option.
func formatOpt(v interface{}) string {
	var s string
	switch v.(type) {
	case format, Cost, OpOpt, ErrOpt:
		s = fmt.Sprintf("%v", v)
	default:
		panic(fmt.Sprintf("unknown type: %T", v))
	}
	return s
}
