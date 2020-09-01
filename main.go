package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type Participants struct {
	Name  string `json:"name,omitempty" bson:"name,omitempty"`
	Email string `json:"email,omitempty" bson:"email,omitempty"`
	RSVP  string `json:"rsvp,omitempty" bson:"rsvp,omitempty"`
}
type Meeting struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title,omitempty" bson:"title,omitempty"`
	Starttime   string             `json:"starttime,omitempty" bson:"starttime,omitempty"`
	Endtime     string             `json:"endtime,omitempty" bson:"endtime,omitempty"`
	Timenow     string             `json:"timenow,omitempty" bson:"timenow,omitempty"`
	participant Participants       `json:"participant,omitempty" bson:"participant,omitempty"`
}

func CreateMeetings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var meet Meeting
	err := json.NewDecoder(r.Body).Decode(&meet)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	database := client.Database("schedulemeeting")
	meetingdata := database.Collection("meeting")

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	fmt.Println(meet)
	result, _ := meetingdata.InsertOne(ctx, meet)
	json.NewEncoder(w).Encode(result)
}

func GetMeetingWidEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var meet Meeting
	collection := client.Database("schedulemeeting").Collection("meeting")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, Meeting{ID: id}).Decode(&meet)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(meet)
}

func GetMeetingEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	start, _ := primitive.ObjectIDFromHex(params["start"])
	end, _ := primitive.ObjectIDFromHex(params["end"])
	fmt.Println(start)
	fmt.Println(end)
	//	var meet []Meeting
	collection := client.Database("schedulemeeting").Collection("meeting")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cur, err := collection.Find(ctx, bson.D{{"starttime", start}})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)
	for cur.Next(context.TODO()) {
		elem := &bson.D{}
		if err := cur.Decode(elem); err != nil {
			log.Fatal(err)
		}
		// ideally, you would do something with elem....
		// but for now just print it to the console
		fmt.Println(elem)
		json.NewEncoder(response).Encode(elem)
	}
	/*	for cur.Next(ctx) {
			var result bson.M
			err := cur.Decode(&result)
			if err != nil {
				log.Fatal(err)
			}
			// do something with result....
			fmt.Println(result)
			json.NewEncoder(response).Encode(result)
		}
		if err := cur.Err(); err != nil {
			log.Fatal(err)
		} */
	/*
		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		cursor, err := collection.Find(ctx).Decode(&meet)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var meeting Meeting
			cursor.Decode(&meeting)
			meet = append(meet, meeting)
		}
		if err := cursor.Err(); err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		json.NewEncoder(response).Encode(meet) */
}

/*find({"OrderDateTime":{ $gte:ISODate("2019-02-10"), $lt:ISODate("2019-02-21") }
}) */
func GetMeetingOfParticiEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	email, _ := primitive.ObjectIDFromHex(params["email"])
	//	var meeting []Meeting
	collection := client.Database("schedulemeeting").Collection("meeting")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cur, err := collection.Find(ctx, bson.D{{"paticipants", bson.D{{"email", email}}}})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		// do something with result....
		fmt.Println(result)
		json.NewEncoder(response).Encode(result)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	/*	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		cursor, err := collection.Find(ctx, Meeting{Paticipants.Email: email})
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var meets Meeting
			cursor.Decode(&meets)
			meeting = append(meeting, meets)
		}
		if err := cursor.Err(); err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		json.NewEncoder(response).Encode(meeting) */
}

func main() {
	fmt.Println("application started at :8080")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
	router := mux.NewRouter()
	router.HandleFunc("/meetings", CreateMeetings).Methods("POST")
	router.HandleFunc("/meeting/{id}", GetMeetingWidEndpoint).Methods("GET")
	router.HandleFunc("/meeting/{start}/{end}", GetMeetingEndpoint).Methods("GET")
	router.HandleFunc("/articals/{paricipants}", GetMeetingOfParticiEndpoint).Methods("GET")
	http.ListenAndServe(":8080", router)
}
