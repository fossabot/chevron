package main

import (
	"encoding/json"
	"fmt"
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron-bootstrap"
	"github.com/quan-to/slog"
)

const (
	messageLog             = "log"
	messageSign            = "sign"
	messageEncrypt         = "encrypt"
	messageDecrypt         = "decrypt"
	messageAddPrivateKey   = "addPrivateKey"
	messageUnlockKey       = "unlockKey"
	messageListPrivateKeys = "listKeys"
)

var electronLog = slog.Scope("Electron")

// handleMessages handles messages
func handleMessages(_ *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	switch m.Name {
	case messageListPrivateKeys:
		return ListPrivateKeys()
	case messageSign:
		var p map[string]interface{}
		if err = json.Unmarshal(m.Payload, &p); err != nil {
			payload = err.Error()
			return
		}

		fingerPrint := p["fingerPrint"]
		data := p["data"]

		if fingerPrint == nil || data == nil {
			err = fmt.Errorf("fingerPrint or data missing")
			return
		}

		return Sign(fingerPrint.(string), data.(string))
	case messageEncrypt:
		return nil, fmt.Errorf("not implemented")
	case messageDecrypt:
		return nil, fmt.Errorf("not implemented")
	case messageAddPrivateKey:
		key := string(m.Payload)
		return AddPrivateKey(key)
	case messageUnlockKey:
		var p map[string]interface{}
		if err = json.Unmarshal(m.Payload, &p); err != nil {
			payload = err.Error()
			return
		}

		fingerPrint := p["fingerPrint"]
		password := p["password"]

		if fingerPrint == nil || password == nil {
			err = fmt.Errorf("fingerPrint or password missing")
			return
		}

		return UnlockKey(fingerPrint.(string), password.(string))
	case messageLog:
		var arguments map[int]interface{}
		if err = json.Unmarshal(m.Payload, &arguments); err != nil {
			payload = err.Error()
			return
		}
		args := make([]interface{}, len(arguments))
		fStr := ""

		i := 0
		for _, v := range arguments {
			args[i] = v
			switch v.(type) {
			case string:
				fStr += "%s "
			case int:
				fStr += "%d "
			case float32:
				fStr += "%f "
			case float64:
				fStr += "%f "
			case map[string]interface{}:
				newV, _ := json.Marshal(v)
				args[i] = string(newV)
				fStr += "%s"
			default:
				fStr += "%+v"
			}
			fStr += " "
			i++
		}
		electronLog.Info(fStr, args...)
	default:
		log.Error("Unknown Message: %+v", m)
	}
	return
}
