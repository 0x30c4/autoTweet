package main;

import (
	"io"
	"os"
	"fmt"
	"log"
	"bytes"
	"strconv"
	"strings"
	"net/url"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"mime/multipart"
	"github.com/google/uuid"
	"github.com/dghubble/oauth1"
)

type MediaUpload struct {
   MediaId int `json:"media_id"`
}

type ImageGenerationRequest struct {
	Prompt string `json:"prompt"`
	N      int    `json:"n"`
	Size   string `json:"size"`
}

type ImageGenerationResponse struct {
	Created int64 `json:"created"`
	Data    []struct {
		URL string `json:"url"`
	} `json:"data"`
}

func downloadFile(url string, filePath string) error {
	// Create the file
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Make the HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check if the response was successful
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: %s", resp.Status)
	}

	// Copy the response body to the file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("File downloaded successfully.")
	return nil
}

func generateImageURLs(prompt string, n int, size string) ([]string, error) {

	apiKey := os.Getenv("CHATGPT_KEY")
	url := "https://api.openai.com/v1/images/generations"

	requestBody := ImageGenerationRequest{
		Prompt: prompt,
		N:      n,
		Size:   size,
	}

	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var imageResponse ImageGenerationResponse
	err = json.Unmarshal(body, &imageResponse)
	if err != nil {
		return nil, err
	}

	var imageURLs []string
	for _, img := range imageResponse.Data {
		imageURLs = append(imageURLs, img.URL)
	}

	return imageURLs, nil
}

type CompletionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Text          string `json:"text"`
		Index         int    `json:"index"`
		Logprobs      interface{} `json:"logprobs"`
		FinishReason  string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens      int `json:"prompt_tokens"`
		CompletionTokens  int `json:"completion_tokens"`
		TotalTokens       int `json:"total_tokens"`
	} `json:"usage"`
}

func generateCompletion(prompt, model string, maxTokens int, temperature float64) (string, error) {

	apiKey := os.Getenv("CHATGPT_KEY")

	url := "https://api.openai.com/v1/completions"

	payload := strings.NewReader(fmt.Sprintf(`{
		"model": "%s",
		"prompt": "%s",
		"max_tokens": %d,
		"temperature": %f
	}`, model, prompt, maxTokens, temperature))

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+apiKey)

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

	var completionResponse CompletionResponse
	err = json.Unmarshal(body, &completionResponse)
	if err != nil {
		return "", err
	}

	if len(completionResponse.Choices) > 0 {
		return completionResponse.Choices[0].Text, nil
	}

	return "", nil
}
func makeTweet(postText, imagePath string) {
   	consumerKey := os.Getenv("TWITTER_CONSUMER_KEY")
	consumerSecret := os.Getenv("TWITTER_CONSUMER_SECRET")
	accessToken := os.Getenv("TWITTER_ACCESS_TOKEN")
	accessSecret := os.Getenv("TWITTER_ACCESS_SECRET")

	if consumerKey == "" || consumerSecret == "" || accessToken == "" || accessSecret == "" {
		panic("Missing required environment variable")
	}

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)

	httpClient := config.Client(oauth1.NoContext, token)

	// create body form
	b := &bytes.Buffer{}
	form := multipart.NewWriter(b)

	// create media paramater
	fw, err := form.CreateFormFile("media", "file.jpg")
	if err != nil {
		panic(err)
	}

	// open file
	opened, err := os.Open(imagePath)
	if err != nil {
		panic(err)
	}

	// copy to form
	_, err = io.Copy(fw, opened)
	if err != nil {
		panic(err)
	}

	// close form
	form.Close()

	// upload media
	resp, err := httpClient.Post("https://upload.twitter.com/1.1/media/upload.json?media_category=tweet_image", form.FormDataContentType(), bytes.NewReader(b.Bytes()))
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	defer resp.Body.Close()

	// decode response and get media id
	m := &MediaUpload{}
	_ = json.NewDecoder(resp.Body).Decode(m)
	mid := strconv.Itoa(m.MediaId)

	// post status with media id
	resp, err = httpClient.PostForm("https://api.twitter.com/1.1/statuses/update.json", url.Values{"status": {postText}, "media_ids": {mid}})
	// parse response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	fmt.Printf("Response: %s\n", body)
}


func main() {

	model := "text-davinci-003"
	prompt := "generate a creative, fictional, futuristic, utopian, dystopian prompt for generating an image"
	maxTokens := 15
	temperature := 0.9

	generatedText, err := generateCompletion(prompt, model, maxTokens, temperature)
	if err != nil {
		log.Fatal(err)
	}

	n := 1
	size := "1024x1024"

	imageURLs, err := generateImageURLs(generatedText, n, size)
	if err != nil {
		log.Fatal(err)
	}

	postText := "prompt: " + generatedText + 
				"\n\nthe prompt for this image was generated with gpt4\n" +
				"\nview source: https://github.com/0x30c4/autoTweet" +
				"\n#chatgpt #chatgpt4 #gpt4 #openai #DALLE #gpt3 #autoTweet"


	fmt.Println(postText, imageURLs)

	for _, url := range imageURLs {
		filePath := "images/" + uuid.New().String() + ".png"
		err := downloadFile(url, filePath)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		makeTweet(postText, filePath)
	}

}
