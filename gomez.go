package main

import (
     "encoding/json"
     "net/http"
     "fmt"
     "html"
     "log"
     "bytes"
     "flag"
     "io/ioutil"
)

var conf Config

// Config is a struct for storing the parsed config.json 
type Config struct {
    FbToken string `json:"fb_token"`
    PageToken string `json:"page_token"`
    CertFile string `json:"cert_file"`
    KeyFile string `json:"key_file"`
}

// MessageJSON is the message that is sent back to the user
type MessageJSON struct {
	AccessToken string `json:"access_token"`
	Recipient struct {
		ID int64 `json:"id"`
	} `json:"recipient"`
	Message struct {
		Text string `json:"text"`
	} `json:"message"`
}
// Response is the message that is sent from from Facebook
type Response struct {
	Object string `json:"object"`
	Entry []struct {
		ID int64 `json:"id"`
		Time int64 `json:"time"`
		Messaging []struct {
			Sender struct {
				ID int64 `json:"id"`
			} `json:"sender"`
			Recipient struct {
				ID int64 `json:"id"`
			} `json:"recipient"`
			Timestamp int64 `json:"timestamp"`
			Message struct {
				Mid string `json:"mid"`
				Seq int `json:"seq"`
				Text string `json:"text"`
			} `json:"message"`
		} `json:"messaging"`
	} `json:"entry"`
}
    
func sendMessage(sender int64, messageData string) {
 m := MessageJSON{}
 m.AccessToken = conf.PageToken
 m.Recipient.ID = sender
 m.Message.Text = messageData
 
 b := new(bytes.Buffer)
 json.NewEncoder(b).Encode(m)
 
  resp, err := http.Post("https://graph.facebook.com/v2.6/me/messages","application/json",b)
  if err != nil {
       log.Println(err)
   } 
     log.Println(resp)
  defer resp.Body.Close()
}    
func webhook(w http.ResponseWriter, r *http.Request) {
        method := r.Method
       
        if method == "GET" {
         q := r.URL.Query()
        
        if q.Get("hub.verify_token") == conf.FbToken {   
            fmt.Fprintf(w,q.Get("hub.challenge")) 
        } else {
            fmt.Fprintf(w, html.EscapeString("Error wrong token"))
        }
        }
        if method == "POST" {
             res := Response{}
             err := json.NewDecoder(r.Body).Decode(&res)
             if err != nil {
                 log.Println(err)
             }
             log.Println(res) 
             if res.Entry[0].Messaging[0].Message.Text != "" {
             
             sendMessage(res.Entry[0].Messaging[0].Sender.ID,"hey") 
             
            }
        }     
}

func config(file string) { 
    content, err := ioutil.ReadFile(file)
    if err != nil {
        log.Println(err)
    }
    err=json.Unmarshal(content, &conf)
    if err != nil {
        log.Println(err)
    }
}

func main() { 
   var confFile = flag.String("config", "", "This is the path to the config file")
   flag.Parse()
   
   config(*confFile)
   http.HandleFunc("/",  webhook)
   http.ListenAndServeTLS(":443",conf.CertFile,conf.KeyFile, nil)
}