# appointy
meeting project

# urls 
router.HandleFunc("/meetings", CreateMeetings).Methods("POST")
router.HandleFunc("/meeting/{id}", GetMeetingWidEndpoint).Methods("GET")
router.HandleFunc("/meeting/{start}/{end}", GetMeetingEndpoint).Methods("GET")
router.HandleFunc("/articals/{paricipants}", GetMeetingOfParticiEndpoint).Methods("GET")
http.ListenAndServe(":8080", router)
  
 # json Structure
  
  {
	"title":"batt krne ke liye",
    "start_time":"2013-10-21T13:28:06.419Z",
    "end_time":"2013-10-21T13:28:06.419Z",
    "time_now":"2013-10-21T13:28:06.419Z",
      "paticipants":[{
        "name":"ayush",
        "email":"xyz@gmail.com",
        "rsvp":"no"
    },
    {
        "name":"piyush",
        "email":"xyz@gmail.com",
        "rsvp":"yes"
    },
    {
        "name":"sachin",
        "email":"xyz@gmail.com",
        "rsvp":"maybe"
    }],
}
