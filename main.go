// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"strconv"
	"time"
	"math"

	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client

func main() {
	var err error
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)

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
				var rtnMsg string = ""				
					
				if message.Text[:1] == "@" {		
					cmd := strings.Split(message.Text, " ")
					switch cmd[0] {
						case "@echo":
							rtnMsg = message.Text[6:len(message.Text)]
						case "@len":
							rtnMsg = strconv.Itoa(len(message.Text) - 5)
						case "@userid":
							rtnMsg = event.Source.UserID
						case "@groupid":
							rtnMsg = event.Source.GroupID
						case "@roomid":
							rtnMsg = event.Source.RoomID
						case "@test":
							rtnMsg = "Test!!"
					}	
				} else {
					/*
					if (os.Getenv("EnableGroup") != event.Source.GroupID) && (len(event.Source.GroupID) > 0) {
						return
					}
					*/
					var msgContent string = strings.ToUpper(message.Text)		
					msgContent = strings.Trim(msgContent, " ")
					msgContent = strings.Trim(msgContent, "　")
					
					today := time.Now().UTC()
					input := time.Now().UTC()
					lmp := time.Now().UTC()
					year := today.Year()

					if (len(msgContent) == 4) {
						msgContent = strconv.Itoa(year) + msgContent
					}

					isTaiwanYear := false
					if (len(msgContent) == 7) {
						msgContent = "0" + msgContent	
						isTaiwanYear = true
					}

					input, err = time.Parse("20060102", msgContent)

					//*****輸入的是最後一次月經，算預產期*****//
					bday := input.AddDate(0, 9, 7)
					// 民國轉西元
					if (isTaiwanYear) {
						lmp = input.AddDate(1911, 0, 0)
					} else {
						lmp = input
					}					
					
					pDiffdays := today.Sub(lmp).Hours() / 24
					if (pDiffdays > 0) {
						pWeek := strconv.Itoa(int(pDiffdays / 7))
						pDays := strconv.Itoa(int(math.Mod(pDiffdays, 7)))

						sBDay := bday.Format("2006/01/02")
						if (isTaiwanYear) {
							sBDay = sBDay[1:len(sBDay)]
						}
						
						rtnMsg = "輸入的是最後一次月經：\n已妊娠 " + pWeek + "週 " + pDays + "天\n預產期為" + sBDay
					} else {
						//*****輸入的是預產期，算週數*****//
						bDiffdays := 280 - (lmp.Sub(today).Hours() / 24)
						if (bDiffdays < 0) {
							bDiffdays = 0 - bDiffdays + 280
						}
						bWeek := strconv.Itoa(int(bDiffdays / 7))
						bDays := strconv.Itoa(int(math.Mod(bDiffdays, 7)))

						rtnMsg = "輸入的是預產期：\n已妊娠 " + bWeek + "週 " + bDays + "天"
					}
				}				
				
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(rtnMsg)).Do(); err != nil {
					log.Print(err)
				}
				
			}
		}
	}
}
