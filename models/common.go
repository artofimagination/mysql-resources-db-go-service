package models

import (
	"encoding/json"
	"fmt"
)

// PrettyPrint logs maps and structs in formatted way in the console.
func PrettyPrint(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
		return
	}
	fmt.Println("Failed to pretty print data")
}
