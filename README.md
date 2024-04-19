# LINE Bot Go Application  
This repository contains a Go application designed to interact with the LINE Messaging API. It handles incoming messages from users and responds appropriately, utilizing Google's Generative AI models to generate responses and describe images.  
## Features  
Text Response: Responds to user text messages by generating content through Google's Gemini Pro model.  
Image Description: Analyzes incoming images and provides a detailed description using Google's Gemini Pro Vision model.  
Asynchronous Processing: Utilizes goroutines for asynchronous API calls and processing, ensuring non-blocking operations.  
## Prerequisites
Go (at least version 1.15)  
Google Cloud account and access to the Generative AI models  
LINE Developer account and access to Messaging API  
## Installation
Clone the repository:
Copy code  
```git clone https://github.com/Rayui1225/go-line-bot.git```  
```cd go-line-bot```  
Set up environment variables:You need to set the following environment variables:  
GOOGLE_GEMINI_API_KEY: Your API key for Google Cloud services.  
ACCESS_TOKEN: Your LINE Messaging API access token.  
CHANNEL_SECRET: Your LINE Messaging API channel secret.  
You can set these variables in your environment or use a .env file and load it with a package like godotenv.  
##Usage
To run the application, use:
```go run main.go```
This will start the server and listen for incoming webhook events from the LINE platform.

Deploying
You can deploy this application to a cloud provider that supports Go, such as Google Cloud Functions or Heroku. Ensure that you configure the environment variables correctly in your deployment setup.
## doing  
Integrate Firebase to give the Line bot memory capabilities.
