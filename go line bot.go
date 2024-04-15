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
)

func main() {
    http.HandleFunc("/callback", callbackHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func callbackHandler(w http.ResponseWriter, req *http.Request){
    bot, err := linebot.New(
        "14d603cb06dc40b42a9a2ce5397792cb", //channel secret key
        "8XeU8l1bWwq/svNJi1l3QmGkeLxE7DYzWMGkRzXvd0TIex+Ta/vU3qufbYqMH4jBAUFwcz6CL2O+VT32swmqwtRzMFTo4IHm6oXG8MY2x3hMndTmdgz5m12lToMwa6fwcyEQxuql39kv+wTg03jRzwdB04t89/1O/w1cDnyilFU=", //access token
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
    for _, event := range events {
        if event.Type == linebot.EventTypeMessage {
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
                // 获取图片内容
                content, err := bot.GetMessageContent(message.ID).Do()
                if err != nil {
                    log.Print(err)
                    continue
                }
                defer content.Content.Close()

                // 保存图片到本地文件系统，文件名为消息ID.jpg
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
client, err := genai.NewClient(ctx, option.WithAPIKey("AIzaSyCgn0QdY6p0qFW9c9b_NpTdOGCoI2ifIAY"))
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
client, err := genai.NewClient(ctx, option.WithAPIKey("AIzaSyCgn0QdY6p0qFW9c9b_NpTdOGCoI2ifIAY"))
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