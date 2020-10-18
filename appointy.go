package main

import(
	"context"
	"fmt"
	"time"
	"net/http"
	"log"
	"testing"
	
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type apiHandler struct{}

func (apiHandler) ServeHTTP(http.ResponseWriter, *http.Request) {}

var collection *mongo.Collection

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://mongodb0.example.com:27017"))
    if err != nil {
        log.Fatal(err)
    }
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
    err = client.Connect(ctx)
    if err != nil {
            log.Fatal(err)
    }
	defer client.Disconnect(ctx)
	
}


type meeting struct {
	Id          		string `json:"id"`
	Title				string `json:"Title"`
	Participants		string `json:"Participants"`
	Start_Time			string `json:"StartTime"`
	End_Time			string `json:"EndTime"`
	Creation_Timestamp	time.Time `json:"CreationTimestamp"`
}

type participant struct {
	Name	string `json:"Name"`
	Email	string `json:"Email"`
	RSVP	string `json:"RSVP"`
}

func createMeeting(){
	mux := http.NewServeMux()
	mux.Handle("/meetings", apiHandler{})
	mux.HandleFunc("/meetings", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/meetings" {
			http.NotFound(w, req)
			return
		}
		fmt.Fprintf(w, "Welcome to the home page! Now you can create meetings")
	})
}

func searchMeeting(id meeting, c *mongo.Database){
	collection := c.Collection("meeting")
	filter := id
	var meet meeting
	mux := http.NewServeMux()
	meeturl:="/meetings/"+meet.Id
	mux.Handle(meeturl, apiHandler{})
	mux.HandleFunc(meeturl, func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != meeturl {
			http.NotFound(w, req)
			return
		}
		fmt.Fprintf(w, "Welcome to the home page! Now you can create meetings")
	})
	err := collection.FindOne(context.TODO(), filter).Decode(&meet)
	if err != nil {
	log.Fatal(err)
	}
	fmt.Println("Found meeting with Id: ",meet.Id)
}

func meetingInATime(time meeting, c *mongo.Database){
	collection := c.Collection("meeting")
	filter := time
	var meet meeting
	mux := http.NewServeMux()
	meeturl:="/meetings/start="+meet.Start_Time+"&end="+meet.End_Time
	mux.Handle(meeturl, apiHandler{})
	mux.HandleFunc(meeturl, func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != meeturl {
			http.NotFound(w, req)
			return
		}
	})
	err := collection.FindOne(context.TODO(), filter).Decode(&meet)
	if err != nil {
	log.Fatal(err)
	}
	fmt.Println("Found meeting during time: ",meet.Start_Time)
}

func myMeetings(email participant, c *mongo.Database){
	collection := c.Collection("meeting")
	filter := email
	var p participant
	mux := http.NewServeMux()
	meeturl:="/meetings/participant="+p.Email
	mux.Handle(meeturl, apiHandler{})
	mux.HandleFunc(meeturl, func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != meeturl {
			http.NotFound(w, req)
			return
		}
	})
	err := collection.FindOne(context.TODO(), filter).Decode(&p)
	if err != nil {
	log.Fatal(err)
	}
	fmt.Println("Found meetings mailed to you: ",p.Email)
}


func MeetingToDB(c *mongo.Database){
	m:= meeting{}
	collection := c.Collection("meeting")
	insertResult, err := collection.InsertOne(context.TODO(), m)
	if err != nil {
	log.Fatal(err)
	}
	fmt.Println("Inserted meeting successfully at id: ",insertResult.InsertedID)
}

func ParticipantToDB(c *mongo.Database){
	p:= participant{}
	collection:= c.Collection("participant")
	insertResult, err := collection.InsertOne(context.TODO(), p)
	if err != nil {
	log.Fatal(err)
	}
	fmt.Println("Added participant successfully! Paticipant id: !",insertResult.InsertedID)
}
