# Keywatcher

`keywatcher` is a Go package for watching key combinations using the keylogger library.

## Requirements

[https://github.com/MarinX/keylogger](https://github.com/MarinX/keyloggerhttps)
Visit upper link if any error happens in keylogger package

## Installation

```bash
go get -u github.com/vizvasrj/keywatcher
```

## Here how to use this

```go

import (
	"log"
	"time"

	"github.com/vizvasrj/keywatcher"
)
func main() {
    k1 := keywatcher.Key{KeyString: "L_CTRL"}
    k2 := keywatcher.Key{KeyString: "ENTER"}
    k3 := keywatcher.Key{KeyString: "L_ALT"}

    kc, err := keywatcher.Watch(k1, k2, k3)
    if err != nil {
        log.Println(err)
        // Handle error
    }
    defer kc.Close()
    after := time.After(10 * time.Second)
    for {
        select {
        case <-kc.WatchChan:
            fmt.Println("Key combination pressed!")
            // Handle key combination
        case <-after:
            return
        }
    }
}

```

now press the key combination and see the magic
first press `ctrl` then `enter` and then `alt`

## Limitations

Currently it does not support pressing event of the key combination.
it only supports the release event of the key combination.
