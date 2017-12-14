package smsc

import (
	"fmt"
	"strconv"
)

// Opt configures a send message and a Result.
type Opt func(*message)

// TODO: Add generic With that will set a field in message based on option type.
//
//     func With(opts ...interface{}) Opt
//

func WithCost(c Cost) Opt {
	return func(m *message) {
		m.Cost = c
	}
}

func WithOp() Opt {
	return func(m *message) {
		m.Op = Op
	}
}

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
	CostOmit Cost = iota
	CostWithoutSend
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

// OpOpt controls whether Response must contain information about sent SMS.
type OpOpt int

const (
	OpOmit OpOpt = iota
	Op
)

// TODO: Add more options.
// Major - ID, Sender, Translit, Subj.
// Minor - TinyURL, Time, Tz, Period, Freq, Flash, Bin, Push, HLR, Ping, MMS,
// Mail, Viber, FileURL, Call, Voice, List, Valid, MaxSMS, ImgCode, UserIP, Err,
// PP.

// formatOpt retuns a string for v value of int-like option.
func formatOpt(v interface{}) string {
	var n int
	switch v := v.(type) {
	case format:
		n = int(v)
	case Cost:
		n = int(v)
	case OpOpt:
		n = int(v)
	default:
		panic(fmt.Sprintf("unknown type: %T", v))
	}
	return strconv.FormatInt(int64(n), 10)
}
