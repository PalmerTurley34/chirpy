# chirpy

Chirpy is a simple web server, built as a guided project for learning 
web servers in Go on [boot.dev](https://www.boot.dev). The Chirpy API is for a simple
social media website where users can create "chirps" to share with other users.

## Documentation

Documentation for the API endpoints are below:

### GET /api/users

Lists all of Chirpy's users. Shows ID, email address, and Chirpy Red Membership status

Response Example:

    [
        {
            "id": 1,
            "email": "alice@example.com",
            "is_chirpy_red": false
        }
    ]

### POST /api/users

Creates a new user resource and returns that resource in response.

Body Example:

    {
        "email": "alice@example.com",
        "password": "pass123"
    }

Response Example:

    {
        "id": 1,
        "email": "alice@example.com",
        "is_chirpy_red": false
    }

### PUT /api/users

Update a user resource and respond with the updated resource. Must be authenticated as the given user to update the resource.

Headers Required:

```Autorization: Bearer ...```

Body Example:

    {
        "email": "alice@example.com",
        "password": "newPassword!"
    }

Response Example:

    {
        "id": 1,
        "email": "alice@example.com",
        "is_chirpy_red": false
    }

### GET /api/users/{ID}

Get single user by their ID. ID values are intergers. Returns a 404 response if the User ID does not exist.

```GET /api/users/1```

    {
        "id": 1,
        "email": "alice@example.com",
        "is_chirpy_red": false
    }

### POST /api/login

Logs in a user. Returns a JWT access token and a JWT refresh token for authentication. The body of the request should contain "email" and "password" fields:

    {
        "email": "alice@example.com",
        "password": "newPassword!"
    }
If password is incorrect, a 401 response is made. Otherwise, a 200 response, and the tokens will be returned:

    {
        "id": 1,
        "email": "alice@example.com",
        "is_chirpy_red": false,
        "token": "<access-token>",
        "refresh_token": "<refresh-token>"
    }

### POST /api/refresh

Expects a `refresh_token` in the header:
`Authorization: Bearer <refresh-token>`
If an access token is provided a 401 response is returned.
Will return a 200 and new access token on success:

    {
        "token": "<access-token>"
    }

### POST /api/revoke

Revokes the refresh token provided in the request header: `Authorization: Bearer <refresh-token>`

Returns a 401 response if token is invalid.
Returns a 200 on success

### GET /api/chirps

Returns all the chirps that have been created, sorted by ID in ascending order:

    [
        {
            "id": 1,
            "body": "example text",
            "author_id": 1
        }
    ]
`author_id` is the ID of the user that created the tweet.

This endpoint accepts two optional query parameters:

`?author_id=1` Returns only chrips belonging to the given author_id.

`?sort=desc` Sort chirps in ascending or descending order. Defaults to ascending.

### POST /api/chirps

Creates a new chirp. User must be authenticated.

Example Body:

    {
        "body": "example text"
    }

Example Response:

    {
        "id": 1,
        "body": "example text",
        "author_id": 1
    }
Returns 401 if `Authorization: Bearer <token>` is not present or invalid

Returns 200 on success

### GET /api/chirps/{ID}

Get single chirp by ID

Returns 404 is ID is invalid

Returns 200 on success

`GET /api/chirps/1`

    {
        "id": 1,
        "body": "example text"
        "author_id": 1
    }

### DELETE /api/chirps/{ID}

Delete chirp with the given ID.

`Authorization: Bearer <token>` header must be present and must be the token of the corresponding `author_id`.

Returns 401 is user is not authorized to delete.

Returns 200 on success.
