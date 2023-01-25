package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"
)

// https://storagep1.chainstrategies.cloud/upload
// https://node.marstonconnell.xyz/upload

func main() {
	target_url := "http://192.168.2.165:3333/upload"

	var fileCounter int

	for {
		// Create a new text file
		file, err := os.Create("Examples" + strconv.Itoa(fileCounter) + ".txt")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		// Write some text to the file

		fileContents := fmt.Sprintf("This is an Examples text file for file No. %d", fileCounter)

		_, err = file.WriteString(fileContents)
		if err != nil {
			fmt.Println(err)
			return
		}

		filename := "Examples" + strconv.Itoa(fileCounter) + ".txt"
		postFile(filename, target_url)

		fileCounter++

		time.Sleep(4 * time.Second)

	}
}

func postFile(filename string, targetUrl string) error {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	// this step is very important
	fileWriter, err := bodyWriter.CreateFormFile("file", filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		return err
	}
	// I need to pass in a sender mang lol
	// open file handle
	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		return err
	}
	defer fh.Close()

	// iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		fmt.Println("error copying")
		return err
	}
	// Below might need to change - we can look at how HTML does it to replicate it
	bodyWriter.WriteField("sender", "jkl10k05lmc88q5ft3lm00q30qkd9x6654h3lejnct")
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(targetUrl, contentType, bodyBuf)
	if err != nil {
		fmt.Println("Post request from main.go failed.")
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(resp.Status)
	fmt.Println(string(resp_body))

	var data UploadResponse

	err = json.Unmarshal(resp_body, &data)
	if err != nil {
		fmt.Println(err)
	}

	// Access the fields of the struct
	fmt.Printf("The CID is %s", data.CID)

	return nil
}

type UploadResponse struct {
	CID string `json:"cid"`
	FID string `json:"fid"`
}
