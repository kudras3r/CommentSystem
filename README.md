# CommentSystem

Test task in Ozon.

## Description

The project is a system for adding and retrieving posts and comments under them, using GraphQL, with the choice of data storage type: PostgreSQL or in-memory. It includes a pagination system based on the limit-offset method. The project also features operation logging and is deployed as a Docker container.


## Installation

```bash
git clone https://github.com/kudras3r/CommentSystem.git
cd CommentSystem/
```

In root create .env file:

```env
DB_HOST=comment-system-db # for db run with d-compose or 0.0.0.0 for manually run
DB_USER=ozon_keker 
DB_PASS=1234
DB_NAME=comm_sys_db
DB_PORT=5432

LOG_LEVEL=DEBUG # INFO

SERVER_HOST=0.0.0.0
SERVER_PORT=8080
```

### Docker

Check docker/docker-compose.yaml

```yaml
environment:
    - POSTGRES_USER=ozon_keker  # your
    - POSTGRES_PASSWORD=1234    # your
    - POSTGRES_DB=comm_sys_db   # your
```

This values will initialize postgres in container. 
.env file values should be the same.

Run:
```bash
cd docker/
chmod +x docker-up.sh
./docker-up.sh --storage=db
```
Or
```bash
./docker-up.sh --storage=inmemory
```

### Manually

In .env DB_HOST set in localhost / 0.0.0.0 

```bash
cd cmd/
go run main.go --storage=inmemory
```
or
```bash
go run main.go --storage=db
```


## Usage

**Create post:**

```GraphQL
mutation {
  createPost(title:"post_title", content:"post_content", allowComms: true, authorID:"0") {
    id          # returned values
    title
    content
    allowComms
    authorID
    createdAt
  }
}
```

Response:
```json
{
  "data": {
    "createPost": {
      "id": "0",
      "title": "post_title",
      "content": "post_content",
      "allowComms": true,
      "authorID": "0",
      "createdAt": "2025-04-01T14:49:12+03:00"
    }
  }
}
```

**Get post:**

```GraphQL
query {
  post(id:"0") {
    id
    content
    allowComms
    comments {
      id 
      content
    }
  }
}
```

Response:
```json
{
  "data": {
    "post": {
      "id": "0",
      "content": "post_content",
      "allowComms": true,
      "comments": []
    }
  }
}
```

**CreateComment:**

```GraphQL
mutation {
  createComment(postID:"0", content:"comment_at_post_0", authorID:"125") {
    id
    content
    createdAt
  }
}
```

Response:
```json
{
  "data": {
    "createComment": {
      "id": "0",
      "content": "comment_at_post_0",
      "createdAt": "2025-04-01T14:53:42+03:00"
    }
  }
}
```

Lets get post(0) again:
```json
{
  "data": {
    "post": {
      "id": "0",
      "content": "post_content",
      "allowComms": true,
      "comments": [
        {
          "id": "0",
          "content": "comment_at_post_0"
        }
      ]
    }
  }
}
```

Now, create another comment at post(0), another post and comment on comment at post(0):

```GraphQL
mutation {
  createComment(postID:"0", content:"another_comment_at_post_0", authorID:"135") {
    id
    content
    createdAt
  }
}
```

```GraphQL
mutation {
  createPost(title:"another_post_title", content:"another_post_content", allowComms: false, authorID:"1") {
    id
    title
    content
    allowComms
    authorID
    createdAt
  }
}
```

```GraphQL
mutation {
  createComment(postID:"0", content:"comment_on_comment_at_post_0", authorID:"111", parentID:"0") {
    id
    content
    createdAt
  }
}
```

Lets take a look at posts:

```GraphQL
query {
  posts {
    content
    allowComms
    comments(first:10) {
      postID
      parentID
      content
      children(first:10) {
        postID
        parentID
        content
      }
    }
  }
}
```

Response:
```json
{
  "data": {
    "posts": [
      {
        "content": "post_content",
        "allowComms": true,
        "comments": [
          {
            "postID": "0",
            "parentID": null,
            "content": "comment_at_post_0",
            "children": [
              {
                "postID": "0",
                "parentID": "0",
                "content": "comment_on_comment_at_post_0"
              }
            ]
          },
          {
            "postID": "0",
            "parentID": null,
            "content": "another_comment_at_post_0",
            "children": []
          }
        ]
      },
      {
        "content": "another_post_content",
        "allowComms": false,
        "comments": []
      }
    ]
  }
}
```

Lets try to create comment at post that not allow comments:
```GraphQL
mutation {
  createComment(postID:"1", content:"comment_at_post_1", authorID:"123") {
    id
    content
    createdAt
  }
}
```

Response:
```json
{
  "errors": [
    {
      "message": "comments not allow at post with id 1",
      "path": [
        "createComment"
      ]
    }
  ],
  "data": null
}
```

Creating a comment that links to a comment located under another post:
```json
{
  "errors": [
    {
      "message": "parent comment does not belong to the current post",
      "path": [
        "createComment"
      ]
    }
  ],
  "data": null
}
```

Сreating a comment under a post that does not exist:
```json
{
  "errors": [
    {
      "message": "no post with id 100",
      "path": [
        "createComment"
      ]
    }
  ],
  "data": null
}
```

## Roadmap

- [27.03.25] make basic structure | first server run
- [28.03.25] restructurize project | add pg db | add first resolvers | add migrations | add .env load func
- [29.03.25] scheme regenerate | bug fixes
- [30.03.25] add service level | add logger | inmemory rework
- [31.03.25] add logging | add docker | bug fixes
- [01.03.25] docker fix | readme update | bug fixes
- [18.04.25] add rwmutex



## Authors

- [@kudras3r](https://www.github.com/kudras3r)

