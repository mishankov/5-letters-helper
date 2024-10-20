package httpclient

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type Response struct {
	Status string
	Body   []byte
}

func Post(url string, body any) (Response, error) {
	reqBody, err := json.Marshal(body)
	if err != nil {
		return Response{}, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return Response{}, err
	}

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return Response{}, err
	}

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}

	return Response{Body: resBody, Status: resp.Status}, nil
}
