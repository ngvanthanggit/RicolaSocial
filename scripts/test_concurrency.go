package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type UpdatePostPayload struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func updatePostTest(postID int64, p UpdatePostPayload, wg *sync.WaitGroup) {
	defer wg.Done()

	url := fmt.Sprintf("http://localhost:3030/v1/posts/%d", postID)

	// create the JSON payload
	b, _ := json.Marshal(p)

	// create new request
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(b))

	if err != nil {
		fmt.Println("Error creating new request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	// send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	defer resp.Body.Close()
	fmt.Println("Update response status:", resp.Status)
}

func main() {
	var wg sync.WaitGroup

	postID := 13

	wg.Add(2)
	title := "New title from A"
	content := "New content from B"

	go updatePostTest(int64(postID), UpdatePostPayload{Title: title}, &wg)
	go updatePostTest(int64(postID), UpdatePostPayload{Content: content}, &wg)

	wg.Wait()
}
