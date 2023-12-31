package keywatcher

import (
	"fmt"
	"log"
	"sync"

	"github.com/MarinX/keylogger"
)

type Key struct {
	KeyString string
}

type KeyCombination struct {
	keys      []Key
	wg        *sync.WaitGroup
	doneChan  chan any
	kl        *keylogger.KeyLogger
	WatchChan chan any
	LastKey   chan string
}

func Watch(k ...Key) (*KeyCombination, error) {
	kc := KeyCombination{}
	for _, this_key := range k {
		if this_key.KeyString == "" {
			return nil, fmt.Errorf("key string cannot be empty")
		}
	}

	kc.keys = k

	wg := sync.WaitGroup{}
	// Create a keylogger
	devs := keylogger.FindAllKeyboardDevices()

	kl, err := keylogger.New(devs[0])
	if err != nil {
		// fmt.Println("Error creating keylogger:", err)
		return nil, err
	}
	// defer kl.Close()
	doneChan := make(chan any)
	kc.wg = &wg
	kc.doneChan = doneChan
	kc.kl = kl

	watchChan := make(chan any)
	kc.WatchChan = watchChan
	last_key := make(chan string)
	kc.LastKey = last_key
	kc.wg.Add(1)
	go func() {
		defer kc.wg.Done()

		kc.wg.Add(1)

		go func() {
			defer kc.wg.Done()
			// defer color.Red("Close first part go routine")

			for {
				select {
				case _, ok := <-kc.doneChan:
					if !ok {
						return
					}
				case e := <-kc.kl.Read():

					// for _, this_key := range kc.keys {
					// 	if this_key.KeyString == e.KeyString() {
					if e.KeyString() != "" && e.KeyString() != "3" {
						// fmt.Printf("key: %s code: %#v, type: %#v, value: %#v, press: %#v, release: %#v\n", e.KeyString(), e.Code, e.Type, e.Value, e.KeyPress(), e.KeyRelease())
						kc.LastKey <- e.KeyString()
					}
					// 	}
					// }

				}
			}

		}()

		// *second part
		lKey := make([]string, len(kc.keys))

		kc.wg.Add(1)
		go func() {
			defer kc.wg.Done()
			// defer color.Red("Close second part go routine")
			for {
				select {
				case k := <-kc.LastKey:
					// color.Green("recieved k %s", k)
					lKey = append(lKey[1:], k)
					if checkKeyCombination(lKey, kc.keys) {
						// color.Red("Key combination pressed!")
						kc.WatchChan <- struct{}{}
					}

				case _, ok := <-kc.doneChan:
					if !ok {
						return
					}
				}
			}
		}()
		kc.wg.Wait()

	}()
	return &kc, nil
}

func checkKeyCombination(currentKeys []string, expectedKeys []Key) bool {
	if len(currentKeys) != len(expectedKeys) {
		return false
	}

	for i := range currentKeys {
		if currentKeys[i] != expectedKeys[i].KeyString {
			return false
		}
	}

	return true
}

func (kc *KeyCombination) Close() {
	log.Println("closing this?")
	if kc.doneChan != nil {
		close(kc.doneChan)

	}

	err := kc.kl.Close()
	if err != nil {
		log.Println(err)
	}

	if kc.WatchChan != nil {
		close(kc.WatchChan)
	}

	if kc.LastKey != nil {
		close(kc.LastKey)
	}

}
