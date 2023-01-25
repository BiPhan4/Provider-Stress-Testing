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
		file, err := os.Create("files3" + strconv.Itoa(fileCounter) + ".txt")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		// Write some text to the file

		fileContents := fmt.Sprintf("This is an files3 text file for file No. %d", fileCounter)

		_, err = file.WriteString(fileContents)
		if err != nil {
			fmt.Println(err)
			return
		}

		for i := 0; i < 20; i++ {
			filename := "files3" + strconv.Itoa(fileCounter) + ".txt"
			go wrapFile(filename, target_url)
			fileCounter++
		}

		time.Sleep(30 * time.Second)

	}
}

func wrapFile(filename string, targetUrl string) {
	cid, err := postFile(filename, targetUrl)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("The CID is %s", cid)
}

func postFile(filename string, targetUrl string) (string, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	// this step is very important
	fileWriter, err := bodyWriter.CreateFormFile("file", filename)
	if err != nil {
		// fmt.Println("error writing to buffer")
		return "", err
	}
	// I need to pass in a sender mang lol
	// open file handle
	fh, err := os.Open(filename)
	if err != nil {
		// fmt.Println("error opening file")
		return "", err
	}

	// iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		// fmt.Println("error copying")
		return "", err
	}

	err = fh.Close()
	if err != nil {
		return "", err
	}
	// Below might need to change - we can look at how HTML does it to replicate it
	err = bodyWriter.WriteField("sender", "jkl10k05lmc88q5ft3lm00q30qkd9x6654h3lejnct")
	if err != nil {
		return "", err
	}
	contentType := bodyWriter.FormDataContentType()
	err = bodyWriter.Close()
	if err != nil {
		return "", err
	}
	resp, error := http.Post(targetUrl, contentType, bodyBuf)
	if error != nil {
		// fmt.Println("Post request from main.go failed.")
		// fmt.Println(err)
		return "", err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	err = resp.Body.Close()

	fmt.Println(resp.Status)
	fmt.Println(string(respBody))

	var data UploadResponse

	err = json.Unmarshal(respBody, &data)
	if err != nil {
		// fmt.Println(err)
		return "", err
	}

	// Access the fields of the struct
	// fmt.Printf("The CID is %s", data.CID)

	return data.CID, nil
}

type UploadResponse struct {
	CID string `json:"cid"`
	FID string `json:"fid"`
}
