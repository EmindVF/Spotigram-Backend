# Spotigram
A backend server which provides messanger and music streaming service functionality.

## Routes

### About
`/about` - Sends an about string.

### Authorization
`/auth/register`- Registers a new user.\
Input:
```json
{
    "name" (8-100 characters)
    "email" (5-100 characters)
    "password" (8-72 characters)
    "password_confirmed" (8-72 characters)
}
```
Output: None

`/auth/login` - logs in.\
Input:
```json
{
    "email" (5-100 characters)
    "password" (8-72 characters)
}
```
Output: 
```json
{
    "id" (id of registered user) (UUID)
}
```
Provides `access_token` and `refresh_token` in the cookies.

`/auth/logout`- logs out.\
Expects `access_token` and `refresh_token`.\
Input: None\
Output: None

`/auth/refresh`- logs out a.\
Expects `refresh_token`.\
Input: None\
Output:
```json
{
    "access_token"
}
```
Provides new  `access_token` in the cookies.

### Me
`/me/friends` - returns a list of friends of current user.\
Expects `access_token`.\
Max amout of friends - 100.\
Input:
```json
{
    "offset" (int)
}
```
Output: 
```json
[
    {
        "user_id1" (UUID)
        "user_id2" (UUID)
        "chat_id" (UUID)
    }
]
```
UserId1 and UserId2 are sorted, so the current user id may be either of them.

`/me/friend-requests-sent` - returns a list of friend requests sent of current user.\
Expects `access_token`.\
Max amout of items - 100.\
Input:
```json
{
    "offset" (int)
}
```
Output: 
```json
[
    {
        "sender_id" (UUID)
        "recipient_id" (UUID)
        "is_ignored" 
    }
]
```

`/me/friend-requests-received` - returns a list of friend requests received of current user.\
Expects `access_token`.\
Max amout of items - 100.\
Input:
```json
{
    "offset" (int)
}
```
Output: 
```json
[
    {
        "sender_id" (UUID)
        "recipient_id" (UUID)
        "is_ignored" 
    }
]
```

`/me/info` - returns info of current user.\
Expects `access_token`.\
Input: None\
Output: 
```json
{
    "id" (UUID)
    "email"
    "name"
    "verified"
}
```

`/me/public-key` - returns the pubic key of current user for end-to-end encryption.\
Expects `access_token`.\
Input: None\
Output: 
```json
{
    "public_key" (base64)
}
```

`/me/picture` - returns the picture of current user.\
Expects `access_token`.\
Input: None\
Output: raw bytes of a webp image.

`/me/change-name` - changes the name of current user.\
Expects `access_token`.\
Input:
```json
{
    "name" (5-100 characters long)
}
```
Output: None 

`/me/change-password` - changes the password of current user.\
Expects `access_token`.\
Input:
```json
{
    "old_password" (8-72 characters long)
    "new_password" (8-72 characters long)
    "new_password_confirmed" (8-72 characters long)
}
```
Output: None 

`/me/change-picture` - changes the picture of current user.\
Expects `access_token`.\
Input: jpg or png image in raw bytes. \
Output: None 

`/me/change-public-key` - changes the public key of current user.\
Expects `access_token`.\
Input:
```json
{
    "public_key" (base64)
}
```
Output: None 

### User
`/user/all` - returns a list of all users.\
Expects `access_token`.\
Max amout of items - 100.\
Input:
```json
{
    "offset" (int)
}
```
Output: 
```json
[
    {
        "id" (UUID)
        "name" 
        "email"
        "verified"
    }
]
```

`/user/info` - returns info of a user.\
Expects `access_token`.\
Input:
```json
{
    "id" (UUID)
}
```
Output: 
```json
{
    "id" (UUID)
    "email"
    "name"
    "verified"
}
```

`/user/picture` - returns picture of a user.\
Expects `access_token`.\
Input:
```json
{
    "id" (UUID)
}
```
Output: raw bytes of a webp image.

`/user/public-key` - returns the pubic key of a user for end-to-end encryption.\
Expects `access_token`.\
Input:
```json
{
    "id" (UUID)
}
```
Output: 
```json
{
    "public_key" (base64)
}
```

### Song
`/song/all` - returns a list of all songs.\
Expects `access_token`.\
Max amout of items - 100.\
Input:
```json
{
    "offset" (int)
}
```
Output: 
```json
[
    {
        "id" (UUID)
        "creator_id" (UUID)
        "name" 
        "length"
    }
]
```

`/song/picture` - returns picture of a song.\
Expects `access_token`.\
Input:
```json
{
    "id" (UUID)
}
```
Output: raw bytes of a webp image.

`/song/download` - downloads a whole song.\
Expects `access_token`.\
Input:
```json
{
    "id" (UUID)
}
```
Output: raw bytes of an mp3 song.

`/song/download` - downloads a whole song.\
Expects `access_token`.\
Input:
```json
{
    "id" (UUID)
}
```
Output: raw bytes of an mp3 song.

`/song/upload/<SONGNAME>` - uploads an mp3 song.\
Expects `access_token`.\
Input: raw bytes of an mp3 song.
Output: None\

`/song/stream` - used for streaming.\
Expects `access_token`.\
Input: None.\
Output: sends files specified in the url \
Example:
`/song/stream/<SONGID>.m3u8` - header for HLS streaming.
`/song/stream/<SONGID>_<CHUNKID>.ts` - a chunk of a song.

### Playlist
`/playlist/all` - returns a list of all user playlists.\
Expects `access_token`.\
Max amout of items - 100.\
Input:
```json
{
    "offset" (int)
}
```
Output: 
```json
[
    {
        "id" (UUID)
        "name" 
        "length"
    }
]
```

`/playlist/songs` - returns a list of all songs of a playlist.\
Expects `access_token`.\
Max amout of items - 100.\
Input:
```json
{
    "id" (UUID)
}
```
Output: 
```json
[
    {
        "id" (UUID)
        "creator_id" (UUID)
        "name" 
        "length"
    }
]
```

`/playlist/create` - creates a playlist.\
Expects `access_token`.\
Input:
```json
{
    "name" (5-100 characters long)
}
```
Output: 
```json
{
    "id" (UUID)
}
```

`/playlist/delete` - deletes a playlist.\
Expects `access_token`.\
Input:
```json
{
    "id" (UUID)
}
```
Output: None\

`/playlist/add-song` - deletes a playlist.\
Expects `access_token`.\
Input:
```json
{
    "id" (UUID)
    "song_id" (UUID)
}
```
Output: None\

`/playlist/delete-song` - deletes a playlist.\
Expects `access_token`.\
Input:
```json
{
    "id" (UUID)
    "song_id" (UUID)
}
```
Output: None\



### Chat
`/chat/messages` - returns a list of messages of a chat.\
Expects `access_token`.\
Max amout of items - 100.\
Input:
```json
{
    "chat_id" (UUID)
    "id" (id of the message based on time)
} 
```
Output: 
```json
[
    {
        "id" (id of the message based on time)
        "user_id" (UUID)
        "chat_id" (UUID)
        "content"
        "date" 
        "is_encrypted"
    }
]
```

`/chat/connect` - connects the user to the websocket for real-time messaging and friend activity.\
Expects `access_token`.\

## How to build
```sh
docker-compose build
docker-compose up
```