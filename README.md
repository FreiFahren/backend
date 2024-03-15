# FreiFahren

## Overview

Freifahren is an innovative project designed to map the presence of ticket inspectors across the Berlin public transport network. By offering a live map that tracks inspectors in real-time, the initiative seeks to inform and empower users to navigate the city with added confidence. The project leverages community-driven data from the [Freifahren Telegram group](https://t.me/freifahren_BE), where users report sightings of ticket inspectors. This repository is the backend that powers the Freifahren web application.

## Getting Started

### Prerequisites

- Go version 1.22 or later
- PostgreSQL 13 or later

### Installation

1. Clone the repository
   ```sh
   git clone https://github.com/FreiFahren/backend
    ```

2. Install Go packages
    ```sh
    go mod download
    ```

3. Set up the database

### Running the application

1. Create a `.env` file in the root directory and add the following environment variables
    ```sh
    DB_USER
    DB_PASSWORD
    DB_HOST
    DB_PORT  
    DB_NAME
    ```

2. Run the application
    ```sh
    go run main.go
    ```

## How it works

We have several API endpoints that allow users to interact with the application. The main endpoints are:

### Getting the id of a station

- `/id` - This endpoint is used to get the id of a station given its name. It is case and whitespace insensitive.

The request should be a `GET` request with the following query parameters:
    - `name` - The name of the station

Example:
```sh
curl -X GET "http://localhost:8080/id?name=alexanderplatz"
```

It will return the id as a text response.

Example:
```sh
"SU-A"
```

### Reporting a new inspector sighting

- `/newInspector` - This endpoint is used to add a new inspector sighting to the database.

The request should be a `POST` request with the following JSON body:
    - `line` - The line on which the inspector was sighted (optional)
    - `station` - The station at which the inspector was sighted (optional)
    - `direction` - The direction in which the inspector was headed (optional)

Example:
```sh
curl -X POST http://localhost:8080/newInspector \
     -H "Content-Type: application/json" \
     -d '{"line":"S7","station":"Alexanderplatz","direction":"Ahrensfelde"}'
```

It will return a json response with the content of the inspector sighting.

Example:
```json
{"line":"S7","station":{"id":"SU-A","name":"Alexanderplatz"},"direction":{"id":"S-Ah","name":"Ahrensfelde"}}
```

### Receive the last known stations 15 mins ago

- `/data` - This endpoint is used to GET all sightings from the last 15 minutes

The request should be a `GET` request, with this example:

Example:
```sh
curl -X GET http://localhost:8080/data \
     -H "Content-Type: application/json" 

```

```json
[{
    "Coordinates": [
      52.52179049999999,
      13.4136147
    ],
    "StationID": "SU-A",
    "DirectionID": "SU-WIU",
    "Line": "U8"
  },
]

```

- `/data?names=true` - This will return the names, instead of id

The request should be a `GET` request, with this example:

Example:
```sh
curl -X GET http://localhost:8080/data \
     -H "Content-Type: application/json" 
```

It will return a json response with the content of the inspector sighting:

```json
[{
    "Coordinates": [
      52.52179049999999,
      13.4136147
    ],
    "StationID": "Alexanderplatz",
    "DirectionID": "Wittenau",
    "Line": "U8"
  },
]
```



### Receive the last known stations 15 mins ago

- `/station?id=` - This endpoint is used to GET the station name through the id


The request should be a `GET` request, with this example:

Example:
curl -X GET http://localhost:8081/station?id=SU-S \
     -H "Content-Type: application/json" 

```json
"Spandau"

```