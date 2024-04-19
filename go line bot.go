package main

import (
    "log"
    "net/http"
    "github.com/line/line-bot-sdk-go/linebot"
    "io/ioutil"
    "context"
    "fmt"
    "github.com/google/generative-ai-go/genai"
    "google.golang.org/api/option"
    "bytes"
)

func main() {
    http.HandleFunc("/callback", callbackHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func callbackHandler(w http.ResponseWriter, req *http.Request){
    secretKey := "channel secret key"
    accessToken := "access token"
    bot, err := linebot.New(
        secretKey, //channel secret key
        accessToken, //access token
    )
    if err != nil {
        log.Fatal(err)
    }
    events, err := bot.ParseRequest(req)
    if err != nil {
        if err == linebot.ErrInvalidSignature {
            w.WriteHeader(400)
        } else {
            w.WriteHeader(500)
        }
        return
    }
    done := make(chan error)
    for _, event := range events {
        if event.Type == linebot.EventTypeMessage {
            var chatId string
            switch event.Source.Type {
            case linebot.EventSourceTypeUser:
                chatId = event.Source.UserID
            case linebot.EventSourceTypeGroup:
                chatId = event.Source.GroupID
            case linebot.EventSourceTypeRoom:
                chatId = event.Source.RoomID
            }
            go sendLoadingAnimation(chatId, accessToken,done)
            go func(){
                if err := <-done; err != nil {
                    log.Fatalf("Failed to send loading animation: %v", err)
                } else {
                    log.Println("Loading animation sent successfully")
                }
            }()
            switch message := event.Message.(type) {
            case *linebot.TextMessage:
                replyMessage ,err:=  replyText(message.Text)
                if err != nil {
                    log.Printf("Error generating reply: %v", err)
                    replyMessage = "Sorry, I encountered an error."
                }
                if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
                    log.Print(err)
                }
            case *linebot.ImageMessage:
                content, err := bot.GetMessageContent(message.ID).Do()
                if err != nil {
                    log.Print(err)
                    continue
                }
                defer content.Content.Close()

                data, err := ioutil.ReadAll(content.Content) // 這個data就可以作為輸入了
                if err != nil {
                    log.Print(err)
                    continue
                }
                replyMessage,err :=  replyImage(data)
                if err != nil {
                    log.Printf("Error generating reply: %v", err)
                    replyMessage = "Sorry, I encountered an error."
                }
                if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
                    log.Print(err)
                }
            }
        }
    }
}

func printResponse(resp *genai.GenerateContentResponse) string {
	var ret string
	for _, cand := range resp.Candidates {
		for _, part := range cand.Content.Parts {
			ret = ret + fmt.Sprintf("%v", part)
			fmt.Println(part)
		}
	}
	return ret 
}

func replyText(s string) (string,error) {
ctx := context.Background()
// Access your API key as an environment variable (see "Set up your API key" above)
client, err := genai.NewClient(ctx, option.WithAPIKey("your API key"))
if err != nil {
  return "" , err
}
defer client.Close()

model := client.GenerativeModel("gemini-pro")
resp, err := model.GenerateContent(ctx, genai.Text(s))
if err != nil {
    return "" , err
}

return printResponse(resp) , nil
}

func replyImage(ImageData []byte) (string,error){
    ctx := context.Background()
// Access your API key as an environment variable (see "Set up your API key" above)
client, err := genai.NewClient(ctx, option.WithAPIKey("your API key"))
if err != nil {
    return "" , err
}
defer client.Close()

model := client.GenerativeModel("gemini-pro-vision")
	prompt := []genai.Part{
		genai.ImageData("png", ImageData),
		genai.Text("Describe this image with scientific detail, reply in zh-TW:"),
	}
	log.Println("Begin processing image...")
	resp, err := model.GenerateContent(ctx, prompt...)
	log.Println("Finished processing image...", resp)
	if err != nil {
		return "" , err
	}
return printResponse(resp),nil
}

func sendLoadingAnimation(chatId string, channelToken string, done chan<- error) {
    url := "https://api.line.me/v2/bot/chat/loading/start"
    requestBody := fmt.Sprintf(`{"chatId":"%s","loadingSeconds":60}`, chatId)
    req, err := http.NewRequest("POST", url, bytes.NewBufferString(requestBody))
    if err != nil {
        done <- err
        return
    }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer " + channelToken)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        done <- err
        return
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        done <- err
        return
    }

    fmt.Println("Response Status:", resp.Status)
    fmt.Println("Response Body:", string(body))
    done <- nil
}

