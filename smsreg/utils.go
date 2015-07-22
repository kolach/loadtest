package smsreg
import "net/http"

func DropCache(url string) bool {
	log.Notice("Clearing server cache")

	httpReq, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Error("Error creating drop cache request", err)
		return false
	}


	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	client := &http.Client{}
	_, err = client.Do(httpReq)

	if err != nil {
		log.Error("Error on dropping server cache: %s", err.Error())
	}

	return err == nil
}