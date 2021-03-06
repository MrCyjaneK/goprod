package goprod

import (
	"encoding/json"
	"fmt"
	"log"
)

type Intent struct {
	Type string `json:"type"`
	Data struct {
		URI       string `json:"uri"`
		Package   string `json:"package"`
		ExtraUsed bool   `json:"extraused"`
		Extra     []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"extra"`
		CustomComponent bool `json:"customcomponent"`
		Component       struct {
			PKG string `json:"pkg"`
			CLS string `json:"cls"`
		} `json:"component"`
	} `json:"data"`
}

func CallIntent(i Intent) error {
	b, err := json.Marshal(i)
	if err != nil {
		log.Println("goporod/android.go CallIntent()", err)
		return err
	}
	fmt.Println(string(b))
	fmt.Println("goprod:" + string(b))
	return nil
}

type Toast struct {
	Type string `json:"type"`
	Data struct {
		Text string `json:"text"`
	} `json:"data"`
}

func SendToast(t Toast) error {
	b, err := json.Marshal(t)
	if err != nil {
		log.Println("goporod/android.go CallIntent()", err)
		return err
	}
	fmt.Println(string(b))
	fmt.Println("goprod:" + string(b))
	return nil
}
