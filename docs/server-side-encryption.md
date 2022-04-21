# Server side encryption
[This PR](https://github.com/soraro/tmpnotes/pull/29) added a way to transparently encrypt a note on the server side so an administrator can not read notes stored in Redis. 

## How it works:
The note is encrypted using AES256 with a key derived from the generated uuid in the following way. A uuid has 32 characters, so the first 8 characters are taken, and that is used as the Redis key. The remaining 24 characters are hashed using SHA512. Using the result of that hash we take the first 32 characters as the encryption key. The note string is encrypted, and hex encoded. This hex encoded text is then what is stored in Redis.

Example on a locally running tmpnotes instance:
```
$ curl -X POST localhost:5000/new -d '{"message": "test note", "ttl": 1}'
a66b0dcabc604ab39c4dfc4884eb428e

# on the redis server:
127.0.0.1:6379> get a66b0dca
"1ee49f9de0ebfcfeacaad9ff699a15f64c3cf1608ea847a33c4a02f2b88ce1d92a159df1a7"
```
Logs from this interaction:
```
{"level":"info","msg":"Server listening at :5000","time":"2022-04-15T10:26:30-06:00"}
{"level":"info","msg":"/new","time":"2022-04-15T10:32:47-06:00"}
{"level":"info","msg":"a66b0dca","time":"2022-04-15T10:50:34-06:00"}
```

## Considerations:
* This code ONLY logs the first 8 characters of the uuid, making it very difficult to derive the full key from incoming requests. The code will never store the full uuid.
* Using the first 8 characters of the generated uuid gives us more opportunities for collisions, but since notes are only kept in Redis for a maximum of 24 hours, it should be *unique enough*:  `16^8 = 4,294,967,296` unique combinations

# Other thoughts:
It is still highly recommended to use the optional encryption key feature so you are the sole owner of the key material and the server side only sees the encrypted data. Additionally, you can use tmpnotes like an API and encrypt your string however you see fit and send the characters to us that way via curl or any other HTTP client.