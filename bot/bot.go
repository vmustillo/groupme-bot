package bot

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"context"
	"time"
	"math/rand"
	"bytes"
)

// Bot is the bot that will be used to send all requests to the GroupMe API
type Bot struct {
	GroupID int64 `json:"group_id"`
	AccessToken string `json:"access_token"`
	DBUri string `json:"dburi"`
	SourceGUID string `json:"source_guid"`
	DB *mongo.Database
}

// User holds the sender_id and sender name to store in the database 
type User struct {
	SenderID string `bson:"sender_id" json:"sender_id"`
	Sender string `bson:"name" json:"name"`
}

// Reply type holds response data needed to post a reply to the API 
type Reply struct {
	Message ResponseData `bson:"message" json:"message"`
}

// ResponseData contains source_guid for the message ID and actual text that is sent to the group
type ResponseData struct {
	SourceGUID string `bson:"source_guid" json:"source_guid"`
	Text string `bson:"text" json:"text"`
}

// InitBot stores the group_id and token needed to make requests for a certain group
func InitBot(configFile string) *Bot {
	file, err := ioutil.ReadFile(configFile)
	Handle(err)

	var bot Bot
	json.Unmarshal(file, &bot)

	// Create new client
	client, err := mongo.NewClient(options.Client().ApplyURI(bot.DBUri))
	if err != nil {
		Handle(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		Handle(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		Handle(err)
	} else {
		fmt.Println("Connected to MongoDB!")
	}

	db := client.Database("groupme-bot")

	bot.DB = db

	return &bot
}

// GetMessages makes a get request to the GroupMe API and returns a specified number of messages
func (bot *Bot) GetMessages(numMsgs int) []byte {
	queryString := fmt.Sprintf("https://api.groupme.com/v3/groups/%d/messages?limit=%d&token=%s", bot.GroupID, numMsgs, bot.AccessToken)
	
	res, err := http.Get(queryString)
	if err != nil {
		Handle(err)
	}

	data, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	return data
}

// StoreMessage stores an array of Messages in the database
func (bot *Bot) StoreMessage(msg Message) error {
	collection := bot.DB.Collection("messages")
	insert, err := collection.InsertOne(context.Background(), msg)
	if err != nil {
		Handle(err)
		return err
	}

	log.Print(insert)
	return nil
}

// MessageExists returns true if a message in the db is found with the same id
func (bot *Bot) MessageExists(msg Message) bool {
	var m Message

	filter := bson.D{{"id", msg.MsgID}}
	collection := bot.DB.Collection("messages")
	err := collection.FindOne(context.Background(), filter).Decode(&m)

	if err != nil {
		log.Print("No new unique messages found")
		return false
	}
	return true
}

// StoreUser stores an array of Users in the database 
func (bot *Bot) StoreUser(sender string, senderID string) error {
	collection := bot.DB.Collection("members")
	insert, err := collection.InsertOne(context.Background(), User{senderID, sender})
	if err != nil {
		log.Print("No new users found")
		return err
	}

	log.Print(insert)
	return nil
}

// UserExists returns true if a user in the db is found with the same sender_id
func (bot *Bot) UserExists(senderID string) bool {
	var u User

	filter := bson.D{{"sender_id", senderID}}
	collection := bot.DB.Collection("members")
	err := collection.FindOne(context.Background(), filter).Decode(&u)

	if err != nil {
		return false
	}
	return true
}

// RespondToMessage takes a message as a parameter and chooses a random response to reply with from an array
func (bot *Bot) RespondToMessage(msg Message) error {
	responses := [...]string{"Sheesh that is SCORCHING", "LMAOOOOOOOOO Berg give this man his phone back", "No gas no cap no gascap that's actually a takepiece", "yeah okay and Chase ran a sub 14 second mile....", "That's actually a good take", "the Padres are actually ass", "Moe Harkless hit a buzzer-beater against the Knicks"}
	rand.Seed(time.Now().UnixNano())
	randomResponseIndex := rand.Intn(len(responses)-0) + 0
	fmt.Println(randomResponseIndex)

	replyData := ResponseData{bot.SourceGUID, responses[randomResponseIndex]}
	reply := Reply{replyData}
	jsonReply, err := json.Marshal(reply)
	
	if err != nil {
		Handle(err)
		return err
	}
	
	queryString := fmt.Sprintf("https://api.groupme.com/v3/groups/%d/messages?token=%s", bot.GroupID, bot.AccessToken)
	res, err := http.Post(queryString, "application/json", bytes.NewBuffer(jsonReply))

	res.Body.Close()
	return nil
}

// Handle is a simple error handler to log errors call a panic
func Handle(err error) {
	if err != nil {
		// log.Panic(err)
		fmt.Println(err)
	}
}
