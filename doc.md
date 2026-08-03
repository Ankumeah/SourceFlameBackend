# Base URL: /api/(version number, default v1)/

## /ping
- GET("") -> `"PONG"`

## /login
> requires -H X-Authorization: Basic base64(username:password)
- POST("") ->
```json
{
  "refresh_token": refresh_token,"
  JWT": JWT
}
```

## /pat
> requires -H X-Authorization: Basic base64(username:password)
- POST("/:pat_name") -> `{ "PAT": string }`
- DELETE("/:pat_name") -> `null`
- GET("/:pat_name") ->
```json
{
  "name": string,
  "creation": uint64,
  "last_used": uint64 | null
}
```
- GET("/all") -> `{ "PATs": []string }`

## /repo
- GET("/:repo_owner/:repo_name/meta") ->
```json
{
  "creation": uint64,
  "stars": uint64,
  "owner": string,
  "name": string
}
```
- GET("/:repo_owner/:repo_name/blob/*path?hash=hash") -> -H "Content-Type: type" file_contents
- GET("/:repo_owner/:repo_name/list/*path?hash=hash") ->
```json
{
  "files": []{
    "file_name": string,
    "dir": bool
  }
}
```
- GET("/:repo_owner/:repo_name/commits/:branch") ->
```json
{
  "commits": []{
    "author" {
      "name": string,
      "email": string
    },
    "message": string,
    "hash": string,
    "timestamp": uint64
  }
}
```
- GET("/:repo_owner/:repo_name/branches") -> `{ "branches": []string }`
- GET("/:repo_owner/:repo_name/blame/*path?hash=hash") ->
```json
{
  "blame": []{
    "author" {
      "name": string,
      "email": string
    },
    "text": string,
    "timestamp": uint64,
    "hash": string
  }
}
```
- GET("/:repo_owner/:repo_name/members") ->
```json
{
  "members": []{
    "username": string,
    "addition": uint64
  }
}
```
> Git client use
  - GET("/:repo_owner/:repo_name/info/refs?service=(git-receive-pack|git-upload-pack)")
  - POST("/:repo_owner/:repo_name/git-upload-pack")
  > requires -H "Authorization: Basic"
  - POST("/:repo_owner/:repo_name/git-receive-pack")

## /repos
> requires -H "X-Authorization: Bearer JWT"
- POST("/:repo_name") -> `null`
- DELETE("/:repo_name") -> `null`
- POST("/:repo_name/transfer/:new_owner") -> `null`
- POST("/:repo_name/members/:member") -> `null`
- DELETE("/:repo_name/members/:member") -> `null`

## /session
> requires -H "X-Authorization: Bearer refresh_token" -H "X-Username: username"
- POST("/renew_jwt") -> `{ "JWT": string }`
- POST("/renew_session") -> `{ "refresh_token": string }`
- DELETE("/") -> `null`

## /user
- GET("/:username/meta") ->
```json
{
  "creation": uint64,
  "username": string
}
```
- GET("/:username/repos?limit=uint8(default 10)&offset=uint64(default 0)") ->
```json
{
  "repos": []{
    "commits": []{
      "author" {
        "name": string,
        "email": string
      },
      "message": string,
      "hash": string,
      "timestamp": uint64
    }
  }
}
```

## Errors
> Any of the APIs may return the following if they encounter an error
- `{ "error": string }`
