package morse

import (
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
)

// https://en.wikipedia.org/wiki/Morse_code

// Duration for morse code
type Duration string

// Morse code signal durations
const (
	Dit Duration = "•" // short
	Dah Duration = "−" // long
)

// constants for beep sounds
const (
	hz  = 800
	wpm = 10

	durationShort = 1200 / wpm
	durationLong  = durationShort * 3
	durationGap   = durationShort * 2
)

// Code for morse code strings
type Code string

// International morse codes (ITU)
const (
	// alphabets
	A Code = Code(Dit + Dah)
	B Code = Code(Dah + Dit + Dit + Dit)
	C Code = Code(Dah + Dit + Dah + Dit)
	D Code = Code(Dah + Dit + Dit)
	E Code = Code(Dit)
	F Code = Code(Dit + Dit + Dah + Dit)
	G Code = Code(Dah + Dah + Dit)
	H Code = Code(Dit + Dit + Dit + Dit)
	I Code = Code(Dit + Dit)
	J Code = Code(Dit + Dah + Dah + Dah)
	K Code = Code(Dah + Dit + Dah)
	L Code = Code(Dit + Dah + Dit + Dit)
	M Code = Code(Dah + Dah)
	N Code = Code(Dah + Dit)
	O Code = Code(Dah + Dah + Dah)
	P Code = Code(Dit + Dah + Dah + Dit)
	Q Code = Code(Dah + Dah + Dit + Dah)
	R Code = Code(Dit + Dah + Dit)
	S Code = Code(Dit + Dit + Dit)
	T Code = Code(Dah)
	U Code = Code(Dit + Dit + Dah)
	V Code = Code(Dit + Dit + Dit + Dah)
	W Code = Code(Dit + Dah + Dah)
	X Code = Code(Dah + Dit + Dit + Dah)
	Y Code = Code(Dah + Dit + Dah + Dah)
	Z Code = Code(Dah + Dah + Dit + Dit)

	// numbers
	One   Code = Code(Dit + Dah + Dah + Dah + Dah)
	Two   Code = Code(Dit + Dit + Dah + Dah + Dah)
	Three Code = Code(Dit + Dit + Dit + Dah + Dah)
	Four  Code = Code(Dit + Dit + Dit + Dit + Dah)
	Five  Code = Code(Dit + Dit + Dit + Dit + Dit)
	Six   Code = Code(Dah + Dit + Dit + Dit + Dit)
	Seven Code = Code(Dah + Dah + Dit + Dit + Dit)
	Eight Code = Code(Dah + Dah + Dah + Dit + Dit)
	Nine  Code = Code(Dah + Dah + Dah + Dah + Dit)
	Zero  Code = Code(Dah + Dah + Dah + Dah + Dah)

	Space Code = " "
	None  Code = ""
)

// CodeFromDurations returns a Code from given `durations`.
func CodeFromDurations(durations ...Duration) Code {
	strs := []string{}
	for _, d := range durations {
		strs = append(strs, string(d))
	}

	return Code(strings.Join(strs, ""))
}

// map for codes and characters
var codesMap map[rune]Code
var charsMap map[Code]rune

// regular expression for non-encodable strings
var regexToEscape *regexp.Regexp
var regexRedundantSpaces *regexp.Regexp

// initialize maps and other values
func init() {
	// codes' map
	codesMap = map[rune]Code{
		'a': A,
		'b': B,
		'c': C,
		'd': D,
		'e': E,
		'f': F,
		'g': G,
		'h': H,
		'i': I,
		'j': J,
		'k': K,
		'l': L,
		'm': M,
		'n': N,
		'o': O,
		'p': P,
		'q': Q,
		'r': R,
		's': S,
		't': T,
		'u': U,
		'v': V,
		'w': W,
		'x': X,
		'y': Y,
		'z': Z,

		'1': One,
		'2': Two,
		'3': Three,
		'4': Four,
		'5': Five,
		'6': Six,
		'7': Seven,
		'8': Eight,
		'9': Nine,
		'0': Zero,

		' ': Space,
	}

	// characters' map
	charsMap = make(map[Code]rune)
	for k, v := range codesMap {
		charsMap[v] = k
	}

	regexToEscape = regexp.MustCompile("[^a-zA-Z0-9\\s]+")
	regexRedundantSpaces = regexp.MustCompile("\\s{2,}")
}

// Encode encodes morse codes from given `text`.
//
// Will return an error when given `text` includes non-encodable characters.
func Encode(text string) (codes []Code, err error) {
	codes = []Code{}

	if _, err = Encodable(text); err == nil {
		for _, chr := range strings.ToLowerSpecial(unicode.TurkishCase, text) {
			if code, err := charToCode(chr); err == nil {
				codes = append(codes, code)
			}
		}
	} else {
		err = fmt.Errorf("'%s' is not encodable: %s", text, err)
	}

	return codes, err
}

// Decode decodes given morse `codes` to a string.
func Decode(codes []Code) (decoded string, err error) {
	chars := []rune{}

	if _, err = Decodable(codes); err == nil {
		for _, code := range codes {
			if chr, err := codeToChar(code); err == nil {
				chars = append(chars, chr)
			}
		}
	} else {
		err = fmt.Errorf("'%v' are not decodable: %s", codes, err)
	}

	return string(chars), err
}

// Encodable returns whether given `text` is encodable or not.
func Encodable(text string) (encodable bool, err error) {
	for _, chr := range strings.ToLowerSpecial(unicode.TurkishCase, text) {
		if _, err = charToCode(chr); err != nil {
			return false, err
		}
	}

	return true, nil
}

// Decodable returns whether given `codes` are decodable or not.
func Decodable(codes []Code) (decodable bool, err error) {
	for _, code := range codes {
		if _, err = codeToChar(code); err != nil {
			return false, err
		}
	}

	return true, nil
}

// Escape returns `text` with non-encodable characters and redundant spaces removed/replaced.
func Escape(text string) string {
	return regexRedundantSpaces.ReplaceAllString(regexToEscape.ReplaceAllString(text, ""), " ")
}

// Beep plays sounds for given `codes` synchronously.
func Beep(codes []Code) {
	sr := beep.SampleRate(44100)
	speaker.Init(sr, sr.N(time.Second/100))

	done := make(chan bool)
	for i, code := range codes {
		if i > 0 {
			time.Sleep(durationGap * time.Millisecond)
		}

		for _, chr := range code {
			var duration int
			switch Duration(chr) {
			case Dit:
				duration = durationShort
			case Dah:
				duration = durationLong
			}

			speaker.Play(beep.Seq(beep.Take(sr.N(time.Duration(duration)*time.Millisecond), beeper()), beep.Callback(func() {
				done <- true
			})))
			<-done
		}
	}
}

// beep sound stream
func beeper() beep.Streamer {
	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		for i := range samples {
			samples[i][0] = math.Sin(float64(i) * math.Pi * 2 * hz / 44100)
			samples[i][1] = math.Sin(float64(i) * math.Pi * 2 * hz / 44100)
		}
		return len(samples), true
	})
}

// converts given character to a morse code.
func charToCode(chr rune) (code Code, err error) {
	var found bool
	if code, found = codesMap[chr]; !found {
		err = fmt.Errorf("no matching character in the codes map: '%c'", chr)
	}

	return code, err
}

// converts given morse code to a character.
func codeToChar(code Code) (chr rune, err error) {
	var found bool
	if chr, found = charsMap[code]; !found {
		err = fmt.Errorf("no matching code in the chars map: '%s'", code)
	}

	return chr, err
}
