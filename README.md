# morse-go

A golang library for morse code encoding/decoding.

## how to get

```bash
$ go get -u github.com/meinside/morse-go
```

## how to use

### sample application

```go
package main

import (
	"log"

	"github.com/meinside/morse-go"
)

const (
	phrase = "Testing morse code..."
)

func main() {
	log.Printf("Will encode: %s", phrase)

	// escape before encoding
	escaped := morse.Escape(phrase)

	log.Printf("Escaped: %s", escaped)

	// encode,
	encoded, _ := morse.Encode(escaped)

	log.Printf("Encoded: %s", encoded)

	// decode,
	decoded, _ := morse.Decode(encoded)

	log.Printf("Decoded: %s", decoded)

	// build codes from durations
	codes := []morse.Code{
		morse.CodeFromDurations(morse.Dit, morse.Dit, morse.Dit),
		morse.Space,
		morse.CodeFromDurations(morse.Dah, morse.Dah, morse.Dah),
		morse.Space,
		morse.CodeFromDurations(morse.Dit, morse.Dit, morse.Dit),
	}

	// decode codes from durations
	decoded, _ = morse.Decode(codes)

	log.Printf("Decoded %s to: %s", codes, decoded)
}

```

Result:

```
2020/03/05 17:22:25 Will encode: Testing morse code...
2020/03/05 17:22:25 Escaped: Testing morse code
2020/03/05 17:22:25 Encoded: [− • ••• − •• −• −−•   −− −−− •−• ••• •   −•−• −−− −•• •]
2020/03/05 17:22:25 Decoded: testing morse code
2020/03/05 17:22:25 Decoded [•••   −−−   •••] to: s o s
```

## how to test/benchmark

```bash
$ go test
$ go test -bench .
```

## License

MIT

