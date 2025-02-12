package util

import (
	"encoding/json"
	"fmt"
)

func PrintJSON(v any) {
	o, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(o))
}
