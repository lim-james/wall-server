# Wall Server


# Wall Posts API

## Overview

The Wall Posts API allows you to retrieve existing posts and create new posts on the wall.

## Base URL

`https://ec2-52-77-251-29.ap-southeast-1.compute.amazonaws.com`

## Endpoints

### 1. Get All Posts

#### Endpoint

`GET /`

#### Description

Retrieve all posts on the wall.

#### Response

```json
[
    {
        "id": 1,
        "title": "title",
        "Body": "body"
    },
    {
        "id": 2,
        "title": "mystical",
        "Body": "body"
    }
]
```

#### Example Usage

```bash
curl https://ec2-52-77-251-29.ap-southeast-1.compute.amazonaws.com
```

### 2. Create a New Post

#### Endpoint

`POST /`

#### Description

Create a new post on the wall.

#### Request

```json
{
    "title": "new title",
    "Body": "new body"
}
```

#### Response

```json
{
    "id": 3,
    "title": "new title",
    "Body": "new body"
}
```

#### Example Usage

```bash
curl -X POST -H "Content-Type: application/json" -d '{"title":"new title", "Body":"new body"}' https://ec2-52-77-251-29.ap-southeast-1.compute.amazonaws.com
```


