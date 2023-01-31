# Social Media Backend

This is a small project to help me learn Go and how to write a RESTful API that handles HTTP Requests. 

Included additionally that does not add functionality to the go project, is a Dockerfile for dockerizing the setup.<br>

Example sequence of docker commands (cwd is at root of project):
```
$ docker build . -t goserver:latest
$ docker run -p PORT:PORT goserver
```

## How to use (because I will also forget)
> Server: Compile and set an env variable for the PORT, run the binary.<br>
> Client: Do whatever you want (GET, POST, PUT, DELETE)

Use a REST client to send http requests, will pass JSON and info in the path

### <u>GET (retrieve user / posts)</u>
To get a specific user, you will need to know their email and make a request like:

GET with url: `localhost:port/users/$EMAIL` <br>
You will get a JSON with the user's info.

You can get all the posts of a specific user by making a request like:

GET with url: `localhost:port/posts/$EMAIL` <br>
You will get a JSON of all the posts of the user specified via their email


### <u>POST (create user / post)</u>

To create a user you make a request like:
POST with url: `localhost:port/users` with a JSON body containing the user info
```
{
    "email": "test@example.com",
    "password": "12345",
    "name": "john doe",
    "age": 18
}
```

To create a post, you make a request like:

POST with url: `localhost:port/posts/`, with a JSON body containing the user's email and post content:
```
{
    "email": "test@example.com",
    "text": "omg my cat is about to give birth"
}
```

### <u> PUT (update user / post) </u>
To update a user's information (everything except email and creation time can be changed), make a request like:

PUT with url: `localhost:port/users/$EMAIL`, with a JSON body:
```
{
  "password": "new_secure_password",
  "name": "Uncle Sam Jr",
  "age": 99
}
```

(WIP) No post editing feature at this time. If there will be, we will implement a timestamping feature of when the last edit was made. Should also store the original post creation date as an immutable property.

### <u> DELETE (delete user /post) </u>

To delete a specific user, you make a request like:

DELETE with url: `localhost:port/users/$EMAIL`<br>
(need to the know the user's email)

To delete a specific post, you make a request like:

DELETE with url: `localhost:port/posts/$UUID`<br>
(need to know the post's UUID)


## With guidance from
[Boot.dev](https://boot.dev)