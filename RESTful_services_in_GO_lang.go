package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//mgoSession will be used to maintain an active session with the mongolab database.
var (
	mgoSession *mgo.Session
)

//Userinput Struct of json string that will collwct data from the POSTMAN/CURL
type Userinput struct {
	ID         bson.ObjectId `json:"id" bson:"_id"`
	Coordinate Geometry      `json:"coordinate"`
	Address    string        `json:"address"`
	City       string        `json:"city"`
	Name       string        `json:"name"`
	State      string        `json:"state"`
	Zip        string        `json:"zip"`
}

var t, new, reentry Userinput //Refrences to the UserInput Struct

//Geometry struct will have only lattitude and longitude values.
type Geometry struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

/**
ResponseBody struct have a nested Geometry struct created above. This struct will be used whenever the response is to be sent back to the user via POSTMAN/CURL.
Usually used in POST,GET,UPDATE methods.
**/
type ResponseBody struct {
	ID         bson.ObjectId `json:"id" bson:"_id"`
	Name       string        `json:"name"`
	Address    string        `json:"address"`
	City       string        `json:"city"`
	State      string        `json:"state"`
	Zip        string        `json:"zip"`
	Coordinate Geometry      `json:"coordinate"`
}

var m ResponseBody //Refrence for the Responsebody struct

/**
Drop the database before the first POST operation
**/
var (
	IsDrop = true
)

/**
 GoogleAPIStruct is used to handle the values(struct) returned by the google mp api page.
We only require the latitude and the longitude from this struct for this assignment.
**/
type GoogleAPIStruct struct {
	Results []struct {
		AddressComponents []struct {
			LongName  string   `json:"long_name"`
			ShortName string   `json:"short_name"`
			Types     []string `json:"types"`
		} `json:"address_components"`
		FormattedAddress string `json:"formatted_address"`
		Geometry         struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
			LocationType string `json:"location_type"`
			Viewport     struct {
				Northeast struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"northeast"`
				Southwest struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"southwest"`
			} `json:"viewport"`
		} `json:"geometry"`
		PlaceID string   `json:"place_id"`
		Types   []string `json:"types"`
	} `json:"results"`
	Status string `json:"status"`
}

var jsonData GoogleAPIStruct //Refrence to the GoogleApiStruct
var oid bson.ObjectId        // oid is used to save the user id as the bson object id.

/**
location function is used to handle the "POST" function.
It reads the data(fields) in the request body in the form of json and then inserts the values in the database via the mongo insert clause.
**/

func location(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {

	if err := json.NewDecoder(req.Body).Decode(&t); err != nil {
	}
	//Getting the address
	garbage := ",+"
	address := t.Address
	address = strings.Replace(address, " ", "+", -1)

	//Getting the City name
	city := t.City
	city = strings.Replace(city, " ", "+", -1)
	city = garbage + city

	//Getting the City name
	state := t.State
	state = strings.Replace(state, " ", "+", -1)
	state = garbage + state

	locationstring := "http://maps.google.com/maps/api/geocode/json?address=" + address + city + state + "&sensor=false"

	response, err := http.Get(locationstring)

	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}

		err = json.Unmarshal([]byte(contents), &jsonData) // here!
		if err != nil {
			panic(err)
		}

	}

	insertdb(rw, t, jsonData)
}

/**
clonemgo function creates a databse connection for us everytime we need it.
It connects to the mongolab database(cmpe-273-sagardafle) and thereby creates a session .
**/
func clonemgo() {
	session, err := mgo.Dial("mongodb://sagardafle:sagardafle123@ds045454.mongolab.com:45454/cmpe-273-sagardafle")
	mgoSession = session
	//defer mgoSession.Close()
	if err != nil {
		panic(err)
	}

}

func insertdb(rw http.ResponseWriter, t Userinput, j GoogleAPIStruct) {
	t.Coordinate.Latitude = jsonData.Results[0].Geometry.Location.Lat
	t.Coordinate.Longitude = jsonData.Results[0].Geometry.Location.Lng

	clonemgo()

	mgoSession.SetMode(mgo.Monotonic, true)

	if IsDrop {
		err := mgoSession.DB("cmpe-273-sagardafle").DropDatabase()
		if err != nil {
			panic(err)
		}
	}

	c := mgoSession.DB("cmpe-273-sagardafle").C("user_details")

	oid := bson.NewObjectId()
	// Insert Datas
	err := c.Insert(&Userinput{ID: oid, Name: t.Name, Address: t.Address, City: t.City, State: t.State, Zip: t.Zip, Coordinate: t.Coordinate})

	if err != nil {
		panic(err)
	}

	m := &ResponseBody{
		ID:         oid,
		Name:       t.Name,
		Address:    t.Address,
		City:       t.City,
		State:      t.State,
		Zip:        t.Zip,
		Coordinate: t.Coordinate,
	}

	js, err := json.Marshal(m)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(js)
}

/**
getUserDetails function will be invoked upon "GET" request.
It fetches the user details and encode the data in a json struct based on the id(resource) given in the URL.
**/

func getUserDetails(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	//fmt.Fprintf(rw, "Entered ID: %s\n", id)
	clonemgo()
	//var results ResponseBody
	//results := ResponseBody{}

	oid := bson.ObjectIdHex(id)
	if err := mgoSession.DB("cmpe-273-sagardafle").C("user_details").FindId(oid).One(&m); err != nil {
		panic(err)
	}

	js, err := json.Marshal(m)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(js)

}

/**
updateUserDetails function will be invoked upon "PUT" request.
It read the request body and updates all the fields in database against the ID entered in the URL as resource.
**/

func updateUserDetails(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {

	if err := json.NewDecoder(req.Body).Decode(&new); err != nil {
	}
	//Getting the address
	garbage := ",+"
	address := new.Address
	address = strings.Replace(address, " ", "+", -1)

	//Getting the City name
	city := new.City
	city = strings.Replace(city, " ", "+", -1)
	city = garbage + city

	//Getting the City name
	state := new.State
	state = strings.Replace(state, " ", "+", -1)
	state = garbage + state

	locationstring := "http://maps.google.com/maps/api/geocode/json?address=" + address + city + state + "&sensor=false"

	response, err := http.Get(locationstring)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}

		err = json.Unmarshal([]byte(contents), &jsonData) // here!
		if err != nil {
			panic(err)
		}
	}

	if new.Name != "" {
		reentry.Name = new.Name
	} else {
		reentry.Name = t.Name
	}

	if new.Address != "" {
		reentry.Address = new.Address
	} else {
		reentry.Address = t.Address
	}

	if new.City != "" {
		reentry.City = new.City
	} else {
		reentry.City = t.City
	}

	if new.State != "" {
		reentry.State = new.State
	} else {
		reentry.State = t.State
	}

	if new.Zip != "" {
		reentry.Zip = new.Zip
	} else {
		reentry.Zip = t.Zip
	}

	new.Coordinate.Latitude = jsonData.Results[0].Geometry.Location.Lat
	new.Coordinate.Longitude = jsonData.Results[0].Geometry.Location.Lng

	if new.Coordinate.Latitude != t.Coordinate.Latitude {
		reentry.Coordinate.Latitude = new.Coordinate.Latitude
	} else {
		reentry.Coordinate.Latitude = t.Coordinate.Latitude
	}

	if new.Coordinate.Longitude != t.Coordinate.Longitude {
		reentry.Coordinate.Longitude = new.Coordinate.Longitude
	} else {
		reentry.Coordinate.Longitude = t.Coordinate.Longitude
	}

	id := ps.ByName("id")
	oid := bson.ObjectIdHex(id)
	colQuerier := bson.M{"_id": oid}
	change := bson.M{"$set": bson.M{"name": reentry.Name, "address": reentry.Address, "city": reentry.City, "state": reentry.State, "zip": reentry.Zip, "coordinate": reentry.Coordinate}}
	err = mgoSession.DB("cmpe-273-sagardafle").C("user_details").Update(colQuerier, change)
	if err != nil {
		panic(err)
	}

	// Query All
	err = mgoSession.DB("cmpe-273-sagardafle").C("user_details").FindId(oid).One(&m)

	js, err := json.Marshal(m)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(js)

}

/**
deleteuserdetails function will be invoked upon "DELETE" request.
It accepts the user unique id and call the databse removeId method to delete the respective entries.
**/
func deleteuserdetails(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {

	id := ps.ByName("id")
	oid := bson.ObjectIdHex(id)
	err := mgoSession.DB("cmpe-273-sagardafle").C("user_details").RemoveId(oid)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	router := httprouter.New()
	router.POST("/location", location)
	router.GET("/location/:id", getUserDetails)
	router.PUT("/location/:id", updateUserDetails)
	router.DELETE("/location/:id", deleteuserdetails)
	log.Fatal(http.ListenAndServe(":8080", router))
}
