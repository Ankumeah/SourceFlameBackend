## /login
> requires { "username": string, "password": string }
- POST("") -> { "refresh_token": refresh_token,"JWT": JWT }

## /pat
> requires { "username": string, "password": string }
- POST("/:pat_name") -> { "PAT": string }
- DELETE("/:pat_name") -> null
- GET("/:pat_name") -> { "name": string, "creation": uint64, "last_used": uint64 | null }
- GET("/all") -> { "PATs": []string }

## /repo
> requires -H "X-Authorization: JWT"
- GET("/:repo_owner/:repo_name/meta") -> { "creation": uint64, "stars": uint64, "private": bool, "owner": string }
- GET("/:repo_owner/:repo_name/blob/*path?hash=hash") -H "Content-Type: type" file_contents
- GET("/:repo_owner/:repo_name/list/*path?hash=hash") -> { "files": []{ "file_name": string, "dir": bool } }
> requires -H "Authorization: PAT" (git client use)
- GET("/:repo_owner/:repo_name/info/refs?service=(git-receive-pack|git-upload-pack)")
- POST("/:repo_owner/:repo_name/git-upload-pack")
- POST("/:repo_owner/:repo_name/git-receive-pack")

## /repos
> requires -H "X-Authorization: JWT"
- POST("/:repo_name?private=bool") -> null
- DELETE("/:repo_name") -> null
- GET("/all?limit=uin8(delfault 10)&offset=uint64(delfault 0)") -> { "repos": []string }

## /session
> requires -H "X-Authorization: refresh_token" -H "X-Username: username"
- POST("/renew_jwt") -> { "JWT": string }
- POST("/renew_session") -> { "refresh_token": string }
- DELETE("/") -> null

## /user
- GET("/:username/meta") -> { "creation": uint64 }
- GET("/:username/repos?limit=uin8(delfault 10)&offset=uint64(delfault 0)") -> { "repos": []string }

## Errors
> Any of the api may return the following if they encounter an error
- { "error": string }
