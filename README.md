#REST API CRUD operations using GO Language
REST API CRUD operations using go language and mongo database.

######This assignment emphasizes the usage of REST APIs in GO programming language. We will use Mongo NoSQL database for data persistence. We shall be performing 4 major CRUD operations like POST, GET, UPDATE and DELETE.######

The connection string to mongodb on MongoLab is : 
```
mongodb://sagardafle:sagardafle123@ds045454.mongolab.com:45454/cmpe-273-sagardafle
```

The name of the table is : user_details.

Description: Location and Trip Planner (Go) using httprouter handler and Google Maps API
```
Input: Get, Post, Put, Delete requests from Postman
```
```
Response: 200 OK/ 404 not found
```
* External API's: Google maps

Usage:
######Create Location######

* POST Request:
```
localhost:1111/locations 
```
```
{ "Name":"AT&T Park", "Address":"24 Willie Mays Plaza", "City":"San Francisco", "State":"CA", "Zip":"94107" } Response: 201 Created { "_id": "562c525be7024724c440210f", "Name": "AT&T Park", "Address": "24 Willie Mays Plaza", "City": "San Francisco", "State": "CA", "Zip": "94107", "Coordinates": { "Lattitude": "37.7781747", "Longitude": "-122.3907248" } }
```
######Retrieve Location######

* GET Request:
```
localhost:1111/locations/143c52ffe7057723488f2e40
```
Response:
```
200 OK { "_id": "143c52ffe7057723488f2e40", "Name": "AT&T Park", "Address": "24 Willie Mays Plaza", "City": "San Francisco", "State": "CA", "Zip": "94107", "Coordinates": { "Lattitude": "37.7781747", "Longitude": "-122.3907248" } }
```
######Update Location######

PUT Request:
```
localhost:1111/locations/143c52ffe7057723488f2e40
```
```
{ "Address":"900 North Point St #52", "City":"San Francisco", "State":"CA", "Zip":"94109" } Response: 201 Created { "_id": "562c52ffe7024723488f2b30", "Name": "AT&T Park", "Address": "900 North Point St #52", "City": "San Francisco", "State": "CA", "Zip": "94109", "Coordinates": { "Lattitude": "37.8055762", "Longitude": "-122.4229471" } }
```

######Delete Location######

DELETE Request:
```
localhost:1111/locations/143c52ffe7057723488f2e40
```

Response: 200 OK

