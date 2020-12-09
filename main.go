package anticapcha

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func Push(base64String string, token string) (string, error){

	uri := "https://rucaptcha.com/in.php"
	params := url.Values{}
	params.Add("key", token)
	params.Add("body", base64String)
	params.Add("method", "base64")

	var data = []byte(params.Encode())

	//fmt.Println(string(data))

	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(data))
	req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.66 Safari/537.36")
	req.Header.Set("content-Type", "multipart/form-data")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func Pull(answer string, token string, try int) (string, error, int)  {
	if answer[0:2] != "OK" {
		return answer, nil, 0
	}

	for v := 0; v < try; v++ {
		uri := "https://rucaptcha.com/res.php?key="+token+"&action=get&id=" + answer[3:len(answer)]
		req, err := http.NewRequest("GET", uri, nil)
		req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.66 Safari/537.36")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return "", err, 0
		}
		defer resp.Body.Close()

		body, e := ioutil.ReadAll(resp.Body)
		if e != nil {
			return "", e, 0
		}
		strResp:=string(body)
		fmt.Println("response Body:", strResp)

		if strResp[0:2] == "OK"{
			return strResp[3:len(strResp)],nil, 1
			break;
		}
		time.Sleep(time.Second)
	}
	return "NOT_RECOGNIZED", nil, 0
}