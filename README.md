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

Once the docker container is up and running, use something like Postman to make a GET request to http://localhost:8080/basic and check you get a response. You can also test a POST request by making a request http://localhost:8080/basic with a JSON body.

### Save a request and response

Make a POST request to http://localhost:8080/save

The body of the request must look like this:

``` json
{   
    "RequestRoute" : "CreateTest",
    "RequestMethod" : "POST",
    "Response" : {
        // Any valid JSON can go in here.
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

## TODO

* Implement some form of header request / response so that you can test if a header is correct
* Clear out request files
* Add better error handling when user provides incorrect JSON