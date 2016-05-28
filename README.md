# gomez

Gomez is a minimal implementation of a Facebook Messenger chat bot. 

To run the bot you'll need to go through the process of setting up a Messenger app at Facebook and then
place the credentials in a JSON configuration file:
```json
{
"fb_token" : "somethingsomething",
"page_token" : "somethinsomething",
"cert_file" : "/path/to/cert.pem",
"key_file" : "/path/to/keyfile.pem"
}
```

Facebook requires your bot to use https so you'll need to include the paths to certificates in the JSON file.

To run your your bot:
```sh
gomez -config="/path/to/config.json"
```
