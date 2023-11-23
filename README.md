# ANSIFmt

"fmt" with ANSI escape sequences.

English | [Chinese](./README_zh.md)

### Reference

[wikipedia/ANSI_escape_code](https://www.wikipedia.org/wiki/ANSI_escape_code)

### Usage

Preview on your terminal:

```
go get github.com/Miuzarte/ANSIFmt
go test -run ^TestUsage github.com/Miuzarte/ANSIFmt -v -timeout 1s

# for all ANSI escape sequences
go test -run ^TestAllANSIEscapeSequences github.com/Miuzarte/ANSIFmt -v -timeout 1s
```

#### Simply Print

**output:**

![Print](https://github.com/Miuzarte/ANSIFmt/assets/66856838/46a5ca92-022b-444c-b654-6a622ba72cfb)

**code:**

```go
package main

import (
	"github.com/Miuzarte/ANSIFmt"
)

func simplyPrint() {
	afmt := ANSIFmt.New().
		Set(ANSIFmt.Style.Bold, ANSIFmt.Fore.Red) // []uint8{1, 35}
	_ = afmt.ToString()                               // \x1b[1;35m
	afmt.Println("This is BOLD and RED")

	afmt.With(ANSIFmt.Back.White) // append([]uint8{1, 35}, 47)
	_ = afmt.ToString()           // \x1b[1;35;47m
	afmt.Println("This is BOLD and RED on WHITE")

	afmt.SetFore24bitColor(63, 127, 191) // []uint8{38, 2, 63, 127, 191} rgb(63, 127, 191)
	_ = afmt.ToString()                  // \x1b[38;2;63;127;191m
	afmt.Println("This is (63,127,191)")

	afmt.WithBack8bitColor(239) // append([]uint8{38, 2, 63, 127, 191}, 48, 5, 239) rgb(78,78,78)
	_ = afmt.ToString()         // \x1b[38;2;63;127;191;48;5;239m
	afmt.Println("This is (63,127,191) on deep gray background")
	afmt.Println("See 256-color lookup table:")
	afmt.Println("https://www.wikipedia.org/wiki/ANSI_escape_code#8-bit")

	afmt.Set()
	afmt.Println("Now reset to normal style")
}
```

`NOTICE: Number of bytes written returned by Printf/Print/Println is no longer reliable.`

#### Colorful Logging

**output:**

![Blinking](https://github.com/Miuzarte/ANSIFmt/assets/66856838/6a2b9475-0be1-4172-bcdc-b03f974f22d0)

*Ignore that there's still Panic after Fatal, just for demonstration purposes.*

**code:**

```go
package main

import (
	"bytes"
	"github.com/Miuzarte/ANSIFmt"
	"github.com/sirupsen/logrus"
)

func colorfulLogging() {
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetFormatter(&logFormat{})

	logrus.Trace("Something very low level.")
	logrus.Debug("Useful debugging information.")
	logrus.Info("Something noteworthy happened!")
	logrus.Warn("You should probably take a look at this.")
	logrus.Error("Something failed but I'm not quitting.")
	logrus.Fatal("Bye.")
	logrus.Panic("I'm bailing.")
}

type logFormat struct{}

func (f *logFormat) Format(entry *logrus.Entry) ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.WriteString(logLevelBanner[entry.Level])
	buf.WriteString(entry.Time.Format("[01/02|15:04:05] "))
	buf.WriteString(entry.Message)
	buf.WriteString("\n")
	return buf.Bytes(), nil
}

var (
	logLevelBanner = map[logrus.Level]string{
		logrus.TraceLevel: ANSIFmt.New().
			Set(ANSIFmt.Fore.BrightBlue).
			Sprint("[TRAC]"),
		logrus.DebugLevel: ANSIFmt.New().
			Set(ANSIFmt.Fore.BrightGreen).
			Sprint("[DEBU]"),
		logrus.InfoLevel: ANSIFmt.New().
			Set(ANSIFmt.Fore.BrightWhite).
			Sprint("[INFO]"),
		logrus.WarnLevel: ANSIFmt.New().
			Set(ANSIFmt.Fore.BrightYellow).
			Sprint("[WARN]"),
		logrus.ErrorLevel: ANSIFmt.New().
			Set(ANSIFmt.Fore.BrightRed).
			Sprint("[ERRO]"),
		logrus.FatalLevel: ANSIFmt.New().
			Set(ANSIFmt.Fore.BrightRed, ANSIFmt.Style.SlowBlink).
			Sprint("[FATA]"),
		logrus.PanicLevel: ANSIFmt.New().
			Set(ANSIFmt.Fore.BrightRed, ANSIFmt.Style.SlowBlink, ANSIFmt.Style.Invert).
			Sprint("[PANI]"),
	}
)
```

### Table

```go
ANSIFmt.Style {
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

ANSIFmt.Fore {
    Black:   30,
    Red:     31,
    Green:   32,
    Yellow:  33,
    Blue:    34,
    Magenta: 35,
    Cyan:    36,
    White:   37,

    Custom:  38, // Next arguments are 5;n or 2;r;g;b
	
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

ANSIFmt.Back {
    Black:   40,
    Red:     41,
    Green:   42,
    Yellow:  43,
    Blue:    44,
    Magenta: 45,
    Cyan:    46,
    White:   47,

    Custom:  48, // Next arguments are 5;n or 2;r;g;b
	
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
```
