# subexecd - trigger commands from a redis pubsub channel

subexecd subscribes to a redis pubsub channel and executes pre-defined commands when a matching message is received.

## Usage
Run the go binary from your local path
```
# subexecd -f subexecd.json > /dev/null 2>&1 &
```

## Configuration

subexecd reads its configuration from the config file `subexecd.json` (default location, can be changed with the `-f` flag).

Sample config:
```
{
  "RedisServer": "127.0.0.1:6379",
  "RedisChannel": "subexecd",
  "TriggerList": [
    {
      "Message": "mymessage",
      "Exec": [
        {
          "Cmd": "/usr/bin/somecmd",
          "Args": [
            "--some",
            "parameter"
          ],
          "Timeout": "5s"
        }
      ]
    }
  ]
}
```
