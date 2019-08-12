[![Go Report Card](https://goreportcard.com/badge/github.com/willdot/NotARealServer)](https://goreportcard.com/report/github.com/willdot/NotARealServer)


# Not a real server
This tool was created because sometimes when I'm creating front end applications that talk to an API that hasn't been build yet, it's frustrating to not be able to make real API calls.

This tool will run a fake http server and allow you to create fake requests and responses so that when you make a call to the server, it will return the fake response you created.

## Installation

Clone the repository.

```
docker-compose up
```

## Usage

Once the docker container is up and running, use something like Postman to make a GET request to http://localhost:8080/basic and check you get a response. You can also test a POST request by making a request http://localhost:8080/basicwithbody with a JSON body.

### Save a request and response

Make a POST request to http://localhost:8080/save

The body of the request must look like this:

``` json
{   
    "RequestRoute" : "CreateTest",
    "RequestMethod" : "POST",
    "Response" : {
        "Thing" : "Some thing",
        "Number" : "123",
        "Something": {
            "Text" : "Some text"
        }
    }
}
```

The "RequestRoute" will be the url route you will want to call. In the above example, it would be http://localhost:8080/createtest 

The "RequestMethod" is the http method you want the request to be. In the example above, to get your response you would need to make a http POST request to http://localhost:8080/createtest 

In "Response" you can put any valid JSON. This is what will be returned to you when you make your request to this route with this method. For the example above, if you made a POST request to http://localhost:8080/createtest you would get the following JSON response:
``` json
{
    "Thing" : "Some thing",
    "Number" : "123",
    "Something": {
        "Text" : "Some text"
    }
}
```

### Using headers
If you want to add in some bad header error responses, you can add in a headers section to the request body. 

To use headers you need to supply an array of Header requests, which contain a Header (the key of the header) and then an array of values. You then need to supply a bad response, which is what the user will see if the headers they supply, don't match the headers you want. For example:

``` json
"Headers": [
        {
            "Header": {
                "Content-Type": [
                    "application/json"
                ]
            },
            "BadResponse": {
                "Message": "Content type not allowed",
                "ErrorCode": 400
            }
        },
        {
            "Header": {
                "UserKeys": [
                    "1234",
                    "5678"
                ]
            },
            "BadResponse": {
                "Message": "Api key not valid",
                "ErrorCode": 401
            }
        }
    ]
```

### Remove some or all requests

You can send a request to remove one, many or all saved requests files.

#### Remove one or many
Make a POST request to http://localhost:8080/remove 

The body of the request must look like this, where the requests is an array of each request you wish to delete:

``` json
{
    "Requests": [
        {
            "RequestRoute": "CreateTest",
            "RequestMethod": "POST"
        },
        {
            "RequestRoute": "GetTest",
            "RequestMethod": "GET"
        },
    ]
}
```
#### Remove all
Make a POST request to http://localhost:8080/removeall

Be careful though, as this will remove all requests!!
