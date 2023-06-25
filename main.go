package main

import (
	"os"
	"log"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"math/rand"

	"time"

	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func generateChatCompletion(model string, message string) (string, error) {
	apiKey := os.Getenv("CHATGPT_KEY")
	url := "https://api.openai.com/v1/chat/completions"

	// Prepare the request payload
	payload := struct {
		Model    string `json:"model"`
		Messages []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"messages"`
		Temperature float64 `json:"temperature"`
	}{
		Model: model,
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{
				Role:    "user",
				Content: message,
			},
		},
		Temperature: 0.7,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Request failed with status: %s, message: %s", resp.Status, respBody)
	}

	// Extract the completion from the response
	var respData struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	err = json.Unmarshal(respBody, &respData)
	if err != nil {
		return "", err
	}

	if len(respData.Choices) == 0 {
		return "", fmt.Errorf("No completion generated")
	}

	return respData.Choices[0].Message.Content, nil
}


// Twitter user-auth requests with an Access Token (token credential)
func makeTweet(tweet_str string) {
	// read credentials from environment variables
	consumerKey := os.Getenv("TWITTER_CONSUMER_KEY")
	consumerSecret := os.Getenv("TWITTER_CONSUMER_SECRET")
	accessToken := os.Getenv("TWITTER_ACCESS_TOKEN")
	accessSecret := os.Getenv("TWITTER_ACCESS_SECRET")
	if consumerKey == "" || consumerSecret == "" || accessToken == "" || accessSecret == "" {
		panic("Missing required environment variable")
	}

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)

	// httpClient will automatically authorize http.Request's
	httpClient := config.Client(oauth1.NoContext, token)
	api := twitter.NewClient(httpClient)

	// Send a Tweet
	tweet, resp, err := api.Statuses.Update(tweet_str, nil)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Response: %+v\n", resp)
	log.Printf("Tweet: %+v\n", tweet)


	// Create a new HTTP Client with the config and token
	// httpClient := config.Client(oauth1.NoContext, token)
	//
	// // Create a new Twitter client
	// client := twitter.NewClient(httpClient)
	//
	// // Open the image file
	// file, err := os.Open("img.jpg")
	// if err != nil {
	// 	log.Fatal("Error opening image file:", err)
	// }
	// defer file.Close()
	//
	// // Upload the image to Twitter
	// media, _, err := client.Media.UploadFile(file)
	// if err != nil {
	// 	log.Fatal("Error uploading image:", err)
	// }
	//
	// // Tweet message with the uploaded image
	// tweet, _, err := client.Statuses.Update("Hello, world! #Golang", &twitter.StatusUpdateParams{
	// 	MediaIds: []int64{media.MediaID},
	// })
	// if err != nil {
	// 	log.Fatal("Error posting tweet:", err)
	// }
	//
	// // Print the created tweet ID
	// log.Println("Tweet ID:", tweet.ID)
}


func main() {
	for {
		model := "gpt-3.5-turbo"
		topics := [32]string{"Web Hacking", "Cybersecurity", "CTF", "Docker", "C programming language", "python programming language", "Golang", "netwroking", "Threat Intelligence", "Small Biz Security", "Think Before You Click", "Threat Hunting", "Financial Security", }

		randomIndex := rand.Intn(12)

		topic := topics[randomIndex]

		fmt.Println(topic, randomIndex, topics[0])

		message := "Write an interesting story with full of true facts for a tweet about " + topic + " within 200 characters. Also use proper hashtags. "

		completion, err := generateChatCompletion(model, message)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		log.Println(message)
		log.Println(completion)

		makeTweet(completion)

		min := 20000
		max := 28800

		// Generate a random number within the range
		randomNum := rand.Intn(max-min+1) + min

		time.Sleep(time.Second * time.Duration(randomNum))
	}
}

