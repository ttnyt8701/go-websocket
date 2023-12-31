package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
	"encoding/json"

)

type ClientMessage struct {
	Contents string `json:"contents"`
	UserId   string `json:"userId"`
}

// ハンドシェイクのカスタマイズ
var upgrader = websocket.Upgrader{
	// クライアントのオリジンの検証
    CheckOrigin: func(r *http.Request) bool { return true },
}

// WebSocket接続のハンドリング
func echo(w http.ResponseWriter, r *http.Request) {
	
	// http→websocketにハンドシェイクして接続を確率
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Upgrade error: %s\n", err)
		return
	}
	defer conn.Close()

	// クライアントからのメッセージを待ち受け、受信するたびにコンソールにメッセージを出力。
	for {
		// メッセージ受信
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("Read message error: %s\n", err)
			break
		}

		var clientMsg ClientMessage
		err = json.Unmarshal(msg, &clientMsg)
		if err != nil {
			fmt.Printf("JSON Unmarshal error: %s\n", err)
			break
		}

		fmt.Printf("Received message: %s from userId: %s\n", clientMsg.Contents, clientMsg.UserId)

		
		// 受信したメッセージに応答してメッセージをJSONに変換して送信。
		responseMessage := "contents:" + clientMsg.Contents + "\n" +"UserId:" + clientMsg.UserId 
		responseJSON, err := json.Marshal(map[string]interface{}{
			"messageId": "",
			"userId": clientMsg.UserId ,
			"contents": responseMessage,
			"senderType": 2,
			"postCategoryType": 0,
			"tarotMasterId": "",
			"createdAt": "",
			"updatedAt": "",
			"isDeleted": "",
		})
		if err != nil {
			fmt.Printf("JSON marshalling error: %s\n", err)
			break
		}
		
		// メッセージを送信
		if err = conn.WriteMessage(msgType, responseJSON); err != nil {
			fmt.Printf("Write message error: %s\n", err)
			break
		}
	}
}

func main() {
	http.HandleFunc("/ws", echo)
	fmt.Println("WebSocket server is listening on ws://localhost:9090/ws")
	err := http.ListenAndServe("localhost:9090", nil)
	if err != nil {
		fmt.Printf("Server error: %s\n", err)
	}
}
