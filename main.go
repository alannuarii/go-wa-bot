package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)


func eventHandler(evt interface{}, client *whatsmeow.Client) {
    switch v := evt.(type) {
    case *events.Message:
        fmt.Println("Received a message!", v.Message.GetConversation())

        if v.Message.GetConversation() == "Hai" {

			jid := types.NewJID("120363204122320229", "g.us")

			msg := &waProto.Message{
				Conversation: proto.String("This is response")}
            
            _, err := client.SendMessage(context.Background(), types.JID(jid), msg)
            if err != nil {
                fmt.Println("Error sending message:", err)
            } else {
                fmt.Println("Message sent successfully")
            }
        }
    }
}


func main() {
    dbLog := waLog.Stdout("Database", "DEBUG", true)
    container, err := sqlstore.New("sqlite3", "file:gowabot.db?_foreign_keys=on", dbLog)
    if err != nil {
        panic(err)
    }

    deviceStore, err := container.GetFirstDevice()
    if err != nil {
        panic(err)
    }
    clientLog := waLog.Stdout("Client", "DEBUG", true)
    client := whatsmeow.NewClient(deviceStore, clientLog)
    client.AddEventHandler(func(evt interface{}) {
        eventHandler(evt, client)
    })

    if client.Store.ID == nil {
        qrChan, _ := client.GetQRChannel(context.Background())
        err = client.Connect()
        if err != nil {
            panic(err)
        }
        for evt := range qrChan {
            if evt.Event == "code" {

                fmt.Println("QR code generated: ", evt.Code)
            } else {
                fmt.Println("Login event:", evt.Event)
            }
        }
    } else {
        err = client.Connect()
        if err != nil {
            panic(err)
        }
    }

    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    <-c

    client.Disconnect()
}
