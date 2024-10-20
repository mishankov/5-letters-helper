package httpclient

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Response struct {
	Status string
	Body   []byte
}

func Post(url string, body any) (Response, error) {
	log.Println("POST url:", url)
	reqBody, err := json.Marshal(body)
	if err != nil {
		return Response{}, err
	}

	log.Println("POST body:", string(reqBody))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(string(reqBody))))
	if err != nil {
		return Response{}, err
	}

	req.Header.Add("Content-Type", "application/json")

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}

	return Response{Body: resBody, Status: resp.Status}, nil
}
