package ANSIFmt

import (
	"fmt"
	"io"
	"strconv"
)

type Code = uint8

type Sequences = []Code

type ANSIFmt struct {
	Sequences    Sequences
	FormatSwitch bool
}

type style struct {
	Reset, Normal, Plain           Code //0
	Bold, IncreasedIntensity       Code //1
	Faint, Dim, DecreasedIntensity Code //2
	Italic                         Code //3
	Underline                      Code //4
	SlowBlink                      Code //5
	RapidBlink                     Code //6
	Invert, ReverseVideo           Code //7
	Conceal, Hide                  Code //8
	Strike, CrossedOut, Delete     Code //9
	PrimaryFont, DefaultFont       Code //10
	AlternativeFont1, AlternativeFont2, AlternativeFont3,
	AlternativeFont4, AlternativeFont5, AlternativeFont6,
	AlternativeFont7, AlternativeFont8, AlternativeFont9 Code //11-19
	Fraktur, Gothic, Blackletter Code //20

	DoublyUnderlined, NoBold                     Code //21
	NormalIntensity                              Code //22
	NoItalic, NoFraktur, NoGothic, NoBlackletter Code //23
	NoUnderline                                  Code //24
	NoBlink                                      Code //25
	ProportionalSpacing                          Code //26 (not known to be used on terminals)
	NoInvert, NoReverse                          Code //27
	NoConceal, Reveal                            Code //28
	NoStrike, NoCrossedOut, NoDelete             Code //29

	NoProportionalSpacing Code //50
	Frame                 Code //51
	Encircle              Code //52
	OverLine              Code //53
	NoFrame, NoEncircle   Code //54
	NoOverLine            Code //55

	CustomUnderlineColor   Code //58
	DefaultUnderlineColors Code //59

	IdeogramUnderline, IdeogramRightSideLine             Code //60
	IdeogramDoubleUnderline, IdeogramDoubleRightSideLine Code //61
	IdeogramOverline, IdeogramLeftSideLine               Code //62
	IdeogramDoubleOverLine, IdeogramDoubleLeftSideLine   Code //63
	IdeogramStressMarking                                Code //64
	NoIdeogramAttribute                                  Code //65

	Superscript                Code //73
	Subscript                  Code //74
	NoSuperscript, NoSubscript Code //75
}

type color struct {
	Black   Code //F30, B40
	Red     Code //F31, B41
	Green   Code //F32, B42
	Yellow  Code //F33, B43
	Blue    Code //F34, B44
	Magenta Code //F35, B45
	Cyan    Code //F36, B46
	White   Code //F37, B47

	// https://www.wikipedia.org/wiki/ANSI_escape_code#8-bit
	// https://www.wikipedia.org/wiki/ANSI_escape_code#24-bit
	Custom Code //F38, B48

	Default Code //F39, B49

	BrightBlack   Code //F90, B100
	BrightRed     Code //F91, B101
	BrightGreen   Code //F92, B102
	BrightYellow  Code //F93, B103
	BrightBlue    Code //F94, B104
	BrightMagenta Code //F95, B105
	BrightCyan    Code //F96, B106
	BrightWhite   Code //F97, B107
}

const (
	pref  = "\x1b["
	suff  = "m"
	reset = pref + suff
)

var (
	Style = style{
		Reset: 0, Normal: 0, Plain: 0,
		Bold: 1, IncreasedIntensity: 1,
		Faint: 2, Dim: 2, DecreasedIntensity: 2,
		Italic:     3,
		Underline:  4,
		SlowBlink:  5,
		RapidBlink: 6,
		Invert:     7, ReverseVideo: 7,
		Conceal: 8, Hide: 8,
		Strike: 9, CrossedOut: 9, Delete: 9,
		PrimaryFont: 10, DefaultFont: 10,
		AlternativeFont1: 11,
		AlternativeFont2: 12,
		AlternativeFont3: 13,
		AlternativeFont4: 14,
		AlternativeFont5: 15,
		AlternativeFont6: 16,
		AlternativeFont7: 17,
		AlternativeFont8: 18,
		AlternativeFont9: 19,
		Fraktur:          20, Gothic: 20, Blackletter: 20,
		DoublyUnderlined: 21, NoBold: 21,
		NormalIntensity: 22,
		NoItalic:        23, NoFraktur: 23, NoGothic: 23, NoBlackletter: 23,
		NoUnderline:         24,
		NoBlink:             25,
		ProportionalSpacing: 26,
		NoInvert:            27, NoReverse: 27,
		NoConceal: 28, Reveal: 28,
		NoStrike: 29, NoCrossedOut: 29, NoDelete: 29,

		NoProportionalSpacing: 50,
		Frame:                 51,
		Encircle:              52,
		OverLine:              53,
		NoFrame:               54, NoEncircle: 54,
		NoOverLine: 55,

		CustomUnderlineColor:   58,
		DefaultUnderlineColors: 59,

		IdeogramUnderline: 60, IdeogramRightSideLine: 60,
		IdeogramDoubleUnderline: 61, IdeogramDoubleRightSideLine: 61,
		IdeogramOverline: 62, IdeogramLeftSideLine: 62,
		IdeogramDoubleOverLine: 63, IdeogramDoubleLeftSideLine: 63,
		IdeogramStressMarking: 64,
		NoIdeogramAttribute:   65,

		Superscript:   73,
		Subscript:     74,
		NoSuperscript: 75, NoSubscript: 75,
	}
	Fore = color{
		Black:   30,
		Red:     31,
		Green:   32,
		Yellow:  33,
		Blue:    34,
		Magenta: 35,
		Cyan:    36,
		White:   37,

		Custom: 38, // Next arguments are 5;n or 2;r;g;b

		Default: 39,

		BrightBlack:   90,
		BrightRed:     91,
		BrightGreen:   92,
		BrightYellow:  93,
		BrightBlue:    94,
		BrightMagenta: 95,
		BrightCyan:    96,
		BrightWhite:   97,
	}
	Back = color{
		Black:   40,
		Red:     41,
		Green:   42,
		Yellow:  43,
		Blue:    44,
		Magenta: 45,
		Cyan:    46,
		White:   47,

		Custom: 48, // Next arguments are 5;n or 2;r;g;b

		Default: 49,

		BrightBlack:   100,
		BrightRed:     101,
		BrightGreen:   102,
		BrightYellow:  103,
		BrightBlue:    104,
		BrightMagenta: 105,
		BrightCyan:    106,
		BrightWhite:   107,
	}
)

// New an ANSIFmt
func New() *ANSIFmt {
	return &ANSIFmt{
		Sequences:    nil,
		FormatSwitch: true,
	}
}

// Set new ANSI escape sequences
func (af *ANSIFmt) Set(Formats ...Code) *ANSIFmt {
	af.Sequences = nil
	return af.With(Formats...)
}

// With additional ANSI escape sequences
func (af *ANSIFmt) With(Formats ...Code) *ANSIFmt {
	af.Sequences = append(af.Sequences, Formats...)
	return af
}

// SetFore8bitColor
// from 256-color lookup tables
// https://www.wikipedia.org/wiki/ANSI_escape_code#8-bit
func (af *ANSIFmt) SetFore8bitColor(n uint8) *ANSIFmt {
	return af.Set(38, 5, n)
}

// WithFore8bitColor
// from 256-color lookup tables
// https://www.wikipedia.org/wiki/ANSI_escape_code#8-bit
func (af *ANSIFmt) WithFore8bitColor(n uint8) *ANSIFmt {
	return af.With(38, 5, n)
}

// SetBack8bitColor
// from 256-color lookup tables
// https://www.wikipedia.org/wiki/ANSI_escape_code#8-bit
func (af *ANSIFmt) SetBack8bitColor(n uint8) *ANSIFmt {
	return af.Set(48, 5, n)
}

// WithBack8bitColor
// from 256-color lookup tables
// https://www.wikipedia.org/wiki/ANSI_escape_code#8-bit
func (af *ANSIFmt) WithBack8bitColor(n uint8) *ANSIFmt {
	return af.With(48, 5, n)
}

// SetFore24bitColor in R, G, B
func (af *ANSIFmt) SetFore24bitColor(r, g, b uint8) *ANSIFmt {
	return af.Set(38, 2, r, g, b)
}

// WithFore24bitColor in R, G, B
func (af *ANSIFmt) WithFore24bitColor(r, g, b uint8) *ANSIFmt {
	return af.With(38, 2, r, g, b)
}

// SetBack24bitColor in R, G, B
func (af *ANSIFmt) SetBack24bitColor(r, g, b uint8) *ANSIFmt {
	return af.Set(48, 2, r, g, b)
}

// WithBack24bitColor in R, G, B
func (af *ANSIFmt) WithBack24bitColor(r, g, b uint8) *ANSIFmt {
	return af.With(48, 2, r, g, b)
}

// DisableFmt temporarily if u want
func (af *ANSIFmt) DisableFmt() *ANSIFmt {
	af.FormatSwitch = false
	return af
}

// EnableFmt if u DisableFmt() before
func (af *ANSIFmt) EnableFmt() *ANSIFmt {
	af.FormatSwitch = true
	return af
}

// ToString return full formatted ANSI escape sequences
func (af *ANSIFmt) ToString() string {
	if !af.FormatSwitch {
		return "\x1b[m"
	}
	s := "\x1b["
	l := len(af.Sequences)
	for i := 0; i < l; i++ {
		s += strconv.Itoa(int(af.Sequences[i]))
		if i < l-1 {
			s += ";"
		}
	}
	s += "m"
	return s
}

// ---Print

// Printf just like fmt.Printf
func (af *ANSIFmt) Printf(format string, a ...any) (n int, err error) {
	return fmt.Print(af.ToString(), fmt.Sprintf(format, a...), reset)
}

// Print just like fmt.Print
func (af *ANSIFmt) Print(a ...any) (n int, err error) {
	return fmt.Print(af.ToString(), fmt.Sprint(a...), reset)
}

// Println just like fmt.Println
func (af *ANSIFmt) Println(a ...any) (n int, err error) {
	return fmt.Print(af.ToString(), fmt.Sprintln(a...), reset)
}

// ---Sprint

// Sprintf just like fmt.Sprintf
func (af *ANSIFmt) Sprintf(format string, a ...any) string {
	return fmt.Sprint(af.ToString(), fmt.Sprintf(format, a...), reset)
}

// Sprint just like fmt.Sprint
func (af *ANSIFmt) Sprint(a ...any) string {
	return fmt.Sprint(af.ToString(), fmt.Sprint(a...), reset)
}

// Sprintln just like fmt.Sprintln
func (af *ANSIFmt) Sprintln(a ...any) string {
	return fmt.Sprint(af.ToString(), fmt.Sprintln(a...), reset)
}

// ---Fprint

// Fprintf just like fmt.Fprintf
func (af *ANSIFmt) Fprintf(w io.Writer, format string, a ...any) (n int, err error) {
	return fmt.Fprint(w, fmt.Sprint(af.ToString(), fmt.Sprintf(format, a...), reset))
}

// Fprint just like fmt.Fprint
func (af *ANSIFmt) Fprint(w io.Writer, a ...any) (n int, err error) {
	return fmt.Fprint(w, fmt.Sprint(af.ToString(), fmt.Sprint(a...), reset))
}

// Fprintln just like fmt.Fprintln
func (af *ANSIFmt) Fprintln(w io.Writer, a ...any) (n int, err error) {
	return fmt.Fprint(w, fmt.Sprint(af.ToString(), fmt.Sprintln(a...), reset))
}

// ---Append

// Appendf just like fmt.Appendf
func (af *ANSIFmt) Appendf(b []byte, format string, a ...any) []byte {
	return fmt.Append(b, fmt.Sprint(af.ToString(), fmt.Sprintf(format, a...), reset))
}

// Append just like fmt.Append
func (af *ANSIFmt) Append(b []byte, a ...any) []byte {
	return fmt.Append(b, fmt.Sprint(af.ToString(), fmt.Sprint(a...), reset))
}

// Appendln just like fmt.Appendln
func (af *ANSIFmt) Appendln(b []byte, a ...any) []byte {
	return fmt.Append(b, fmt.Sprint(af.ToString(), fmt.Sprintln(a...), reset))
}
