# autoTweet

autoTweet is a Go program that generates automated tweets using the ChatGPT API and posts them as tweets using the Twitter API. It allows you to set the topics for your tweets and control the time frame for posting them. The program also includes proper hashtags to enhance tweet visibility and engagement.

## Prerequisites

Before running autoTweet, ensure you have the following:

- Go installed on your system
- Twitter API credentials (API keys and access tokens)
- ChatGPT API credentials
- Docker and Docker Compose installed

## Setup

1. Clone this Git repository:

   ```shell
   git clone https://github.com/your-username/autoTweet.git
   ```

2. Copy the example.env file and rename it as .env:

   ```shell
   cp example.env .env
   ```

3. Open the .env file and add your Twitter API and ChatGPT API credentials.

4. Create a directory named ```images``` for the DALL.E generated

   ```shell
   mkdir images
   ```

## Usage

1. Start the application using Docker Compose:

   ```shell
   docker-compose up -d
   ```

   This command will build and run the autoTweet container.

2. The program will automatically generate tweets based on the predefined topics and post them as tweets at a specified time interval (every 7-8 hours by default).

3. Customize the topics and time frame:

   - Modify the topics in the Go code (`main.go`) to suit your preferences.
   - Adjust the time frame for posting tweets by modifying the code in `main.go`.

4. Monitor the logs:

   ```shell
   docker-compose logs -f
   ```

   Use this command to view the logs and track the generated tweets and posting activities.

## Contributing

Contributions are welcome! If you have any improvements, bug fixes, or new features to suggest, please feel free to open an issue or submit a pull request.
