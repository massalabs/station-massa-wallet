package ico

import (
	"fmt"
	"net/http"
)

// ICOQUEST: The all package must be removed when ICO is over
const (
	url         = "http://54.36.174.177:3000/register_quest/"
	projectName = "massastation"
)

// ValidateQuest validates a quest for a given address.
func ValidateQuest(questID string, address string) error {
	url := url + projectName + "/" + questID + "/" + address
	_, err := http.Post(url, "", nil)

	if err != nil {
		fmt.Printf("error validating quest: %v\n", err)
	} else {
		fmt.Printf("quest %s validated\n", questID)
	}

	return err
}
