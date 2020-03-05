package morse

import (
	"strings"
	"testing"
)

const (
	testPhrase = "The Quick & Brown Fox, Jumps Over The Lazy Dog...?!"
)

func TestEncodeAndDecode(t *testing.T) {
	if encodable, _ := Encodable(testPhrase); encodable {
		t.Errorf("only alphabets and numbers are encodable: %s", testPhrase)
	}

	// remove special characters for testing
	escapedPhrase := Escape(testPhrase)

	if encodable, _ := Encodable(escapedPhrase); !encodable {
		t.Errorf("only alphabets and numbers are encodable: %s", escapedPhrase)
	}

	if encoded, err := Encode(escapedPhrase); err != nil {
		t.Errorf("failed to encode: %s", err)
	} else {
		if decoded, err := Decode(encoded); err != nil {
			t.Errorf("failed to decode: %s", err)
		} else {
			// ignore case
			if !strings.EqualFold(decoded, escapedPhrase) {
				t.Errorf("encoded/decoded values do not match: %s / %s", decoded, escapedPhrase)
			}
		}
	}
}

func BenchmarkEncode(b *testing.B) {
	escapedPhrase := Escape(testPhrase)

	for i := 0; i < b.N; i++ {
		Encode(escapedPhrase)
	}
}

func BenchmarkDecode(b *testing.B) {
	escapedPhrase := Escape(testPhrase)

	if encoded, err := Encode(escapedPhrase); err == nil {
		for i := 0; i < b.N; i++ {
			Decode(encoded)
		}
	}
}
