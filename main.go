package main

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context) error {
	buf := bytes.Buffer{}
	mw := multipart.NewWriter(&buf)
	mw.WriteField("message", os.Getenv("MESSAGE"))
	mw.Close()

	req, err := http.NewRequest("POST", "https://notify-api.line.me/api/notify", &buf)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", mw.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+os.Getenv("ACCESS_TOKEN"))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("line notify failed: %d", resp.StatusCode)
	}
	return resp.Body.Close()
}

func main() {
	lambda.Start(HandleRequest)
}
