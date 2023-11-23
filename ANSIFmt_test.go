package ANSIFmt

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestUsage1(t *testing.T) {
	simplyPrint()
}

func TestUsage2(t *testing.T) {
	noExitColorfulLogging()
	//colorfulLogging()
}

func simplyPrint() {
	afmt := New().
		Set(Style.Bold, Fore.Red) // []uint8{1, 35}
	_ = afmt.ToString() // \x1b[1;35m
	afmt.Println("This is BOLD and RED")

	afmt.With(Back.White) // append([]uint8{1, 35}, 47)
	_ = afmt.ToString()   // \x1b[1;35;47m
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

func noExitColorfulLogging() {
	fmt.Print(logLevelBanner[6], "[01/02|15:04:05] ", "Something very low level.\n")
	fmt.Print(logLevelBanner[5], "[01/02|15:04:05] ", "Useful debugging information.\n")
	fmt.Print(logLevelBanner[4], "[01/02|15:04:05] ", "Something noteworthy happened!\n")
	fmt.Print(logLevelBanner[3], "[01/02|15:04:05] ", "You should probably take a look at this.\n")
	fmt.Print(logLevelBanner[2], "[01/02|15:04:05] ", "Something failed but I'm not quitting.\n")
	fmt.Print(logLevelBanner[1], "[01/02|15:04:05] ", "Bye.\n")
	fmt.Print(logLevelBanner[0], "[01/02|15:04:05] ", "I'm bailing.\n")
}

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
		logrus.TraceLevel: New().
			Set(Fore.BrightBlue).
			Sprint("[TRAC]"),
		logrus.DebugLevel: New().
			Set(Fore.BrightGreen).
			Sprint("[DEBU]"),
		logrus.InfoLevel: New().
			Set(Fore.BrightWhite).
			Sprint("[INFO]"),
		logrus.WarnLevel: New().
			Set(Fore.BrightYellow).
			Sprint("[WARN]"),
		logrus.ErrorLevel: New().
			Set(Fore.BrightRed).
			Sprint("[ERRO]"),
		logrus.FatalLevel: New().
			Set(Fore.BrightRed, Style.SlowBlink).
			Sprint("[FATA]"),
		logrus.PanicLevel: New().
			Set(Fore.BrightRed, Style.SlowBlink, Style.Invert).
			Sprint("[PANI]"),
	}
)

func TestOutputCount(t *testing.T) {
	str := "123abc汉字"
	strLen := len(str)
	afmtLen, _ := New().With(1).Print(str)
	fmt.Println()
	fmtLen, _ := fmt.Print(str)
	fmt.Println()
	fmt.Println("afmtLen == strLen:", afmtLen == strLen)
	fmt.Println("fmtLen == strLen:", fmtLen == strLen)
}

func TestXxx(t *testing.T) {
	New().Set(1, 31).Println("1;31")
	New().Set(1).Println("1")
	New().Set(31).Println("31")
	New().Set(91).Println("91")
}

func TestAllANSI(t *testing.T) {
	ifmt := New()
	for i := 0; i < 128; i++ {
		ifmt.Set(Code(i)).Println(
			"                                                                ",
		)
	}
}

func TestTrueColor(t *testing.T) {
	tfmt := New()
	r, g, b := 0, 0, 0
	for r = 0; r < 256; r += 4 {
		for g = 0; g < 256; g += 4 {
			for b = 0; b < 256; b += 4 {
				tfmt.SetBack24bitColor(uint8(r), uint8(g), uint8(b)).Print(" ")
			}
			fmt.Print("\n")
		}
	}
	for i := 0; i < 256; i++ {
		tfmt.SetBack8bitColor(uint8(i)).Print("    ")
		if (i == 7 || i == 15) ||
			(i > 15 && i <= 231 && (i-15)%6 == 0) ||
			(i > 231 && (i-231)%8 == 0) {
			fmt.Print("\n")
		}
	}
}

func TestAllANSIEscapeSequences(t *testing.T) {
	afmt := New()
	fmt.Println("---")
	afmt.Set(0).Println("Reset, Normal, Plain: 0")
	afmt.Set(1).Println("Bold, IncreasedIntensity: 1")
	afmt.Set(2).Println("Faint, Dim, DecreasedIntensity: 2")
	afmt.Set(3).Println("Italic: 3")
	afmt.Set(4).Println("Underline: 4")
	afmt.Set(5).Println("SlowBlink: 5")
	afmt.Set(6).Println("RapidBlink: 6")
	afmt.Set(7).Println("Invert, ReverseVideo: 7")
	afmt.Set(8).Println("Conceal, Hide: 8")
	afmt.Set(9).Println("Strike, CrossedOut, Delete: 9")
	afmt.Set(10).Println("PrimaryFont, DefaultFont: 10")
	afmt.Set(11).Println("AlternativeFont1: 11")
	afmt.Set(12).Println("AlternativeFont2: 12")
	afmt.Set(13).Println("AlternativeFont3: 13")
	afmt.Set(14).Println("AlternativeFont4: 14")
	afmt.Set(15).Println("AlternativeFont5: 15")
	afmt.Set(16).Println("AlternativeFont6: 16")
	afmt.Set(17).Println("AlternativeFont7: 17")
	afmt.Set(18).Println("AlternativeFont8: 18")
	afmt.Set(19).Println("AlternativeFont9: 19")
	afmt.Set(20).Println("Fraktur, Gothic, Blackletter: 20")
	afmt.Set(21).Println("DoublyUnderlined, NoBold: 21")
	afmt.Set(22).Println("NormalIntensity: 22")
	afmt.Set(23).Println("NoItalic, NoFraktur, NoGothic, NoBlackletter: 23")
	afmt.Set(24).Println("NoUnderline: 24")
	afmt.Set(25).Println("NoBlink: 25")
	afmt.Set(26).Println("ProportionalSpacing: 26")
	afmt.Set(27).Println("NoInvert, NoReverse: 27")
	afmt.Set(28).Println("NoConceal, Reveal: 28")
	afmt.Set(29).Println("NoStrike, NoCrossedOut, NoDelete: 29")
	afmt.Set(50).Println("NoProportionalSpacing: 50,")
	afmt.Set(51).Println("Frame: 51,")
	afmt.Set(52).Println("Encircle: 52,")
	afmt.Set(53).Println("OverLine: 53,")
	afmt.Set(54).Println("NoFrame, NoEncircle: 54,")
	afmt.Set(55).Println("NoOverLine: 55,")
	afmt.Set(58).Println("CustomUnderlineColor: 58,")
	afmt.Set(59).Println("DefaultUnderlineColors: 59,")
	afmt.Set(60).Println("IdeogramUnderline, IdeogramRightSideLine: 60,")
	afmt.Set(61).Println("IdeogramDoubleUnderline, IdeogramDoubleRightSideLine: 61,")
	afmt.Set(62).Println("IdeogramOverline, IdeogramLeftSideLine: 62,")
	afmt.Set(63).Println("IdeogramDoubleOverLine, IdeogramDoubleLeftSideLine: 63,")
	afmt.Set(64).Println("IdeogramStressMarking: 64,")
	afmt.Set(65).Println("NoIdeogramAttribute: 65,")
	afmt.Set(73).Println("Superscript: 73,")
	afmt.Set(74).Println("Subscript: 74,")
	afmt.Set(75).Println("NoSuperscript, NoSubscript: 75,")
	fmt.Println("---")
	afmt.Set(30).Println("Black:   30,")
	afmt.Set(31).Println("Red:     31,")
	afmt.Set(32).Println("Green:   32,")
	afmt.Set(33).Println("Yellow:  33,")
	afmt.Set(34).Println("Blue:    34,")
	afmt.Set(35).Println("Magenta: 35,")
	afmt.Set(36).Println("Cyan:    36,")
	afmt.Set(37).Println("White:   37,")
	afmt.Set(38).Println("Custom:  38,")
	afmt.Set(39).Println("Default: 39,")
	afmt.Set(90).Println("BrightBlack:   90,")
	afmt.Set(91).Println("BrightRed:     91,")
	afmt.Set(92).Println("BrightGreen:   92,")
	afmt.Set(93).Println("BrightYellow:  93,")
	afmt.Set(94).Println("BrightBlue:    94,")
	afmt.Set(95).Println("BrightMagenta: 95,")
	afmt.Set(96).Println("BrightCyan:    96,")
	afmt.Set(97).Println("BrightWhite:   97,")
	fmt.Println("---")
	afmt.Set(40).Println("Black:   40,")
	afmt.Set(41).Println("Red:     41,")
	afmt.Set(42).Println("Green:   42,")
	afmt.Set(43).Println("Yellow:  43,")
	afmt.Set(44).Println("Blue:    44,")
	afmt.Set(45).Println("Magenta: 45,")
	afmt.Set(46).Println("Cyan:    46,")
	afmt.Set(47).Println("White:   47,")
	afmt.Set(48).Println("Custom:  48,")
	afmt.Set(49).Println("Default: 49,")
	afmt.Set(100).Println("BrightBlack:   100,")
	afmt.Set(101).Println("BrightRed:     101,")
	afmt.Set(102).Println("BrightGreen:   102,")
	afmt.Set(103).Println("BrightYellow:  103,")
	afmt.Set(104).Println("BrightBlue:    104,")
	afmt.Set(105).Println("BrightMagenta: 105,")
	afmt.Set(106).Println("BrightCyan:    106,")
	afmt.Set(107).Println("BrightWhite:   107,")
	fmt.Println("---")
}
