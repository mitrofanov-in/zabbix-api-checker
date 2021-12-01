package auth

import (
	"bytes"
	//	"encoding/json"
	"fmt"
	"net/http"
)

// Get cookie
var BodyStatus_auth int = 0

var cookie_auth []*http.Cookie

func HttpQueryPost(url string, jstr []byte) []*http.Cookie {

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}

	req_lgn, _ := http.NewRequest("POST", url, bytes.NewBuffer(jstr))
	req_lgn.Header.Set("Content-Type", "application/json")
	resp_lgn, err := client.Do(req_lgn)
	if err != nil {
		panic(err)
	}
	cookie_auth = resp_lgn.Cookies()

	for _, c := range cookie_auth {
		fmt.Println(c.Name, c.Value)
	}
	defer resp_lgn.Body.Close()

	BodyStatus_auth = resp_lgn.StatusCode
	fmt.Println(BodyStatus_auth)
	return cookie_auth

}
