![impostor](impostor.png)

# Impostor
Impostor is a REST API simulator (Mock server) written in golang. It'll runs based on a configuration provided to simulate various use cases.

## Getting Started
install `impostor`
```
$ go get -u -v github.com/ratanphayade/impostor
```

### Using Impostor
`impostor` has multiple options to make it developer friendly. Custom option are:
```
  -app string
    	name of the app which has to be mocked (default "app")
      
  -config string
    	path with configuration file (default "config.toml")
      
  -host string
    	if you run your server on a different host (default "localhost")
      
  -mock string
    	directory where your mock configs are saved (default "test")
      
  -port int
    	port to run the server (default 9000)
      
  -watch
    	if true, then watch for any change in app mock config and reload
```

`impostor` can be started in two ways

- By providing all the required parameters in command
```
$ impostor -host=127.0.0.1 -port=9000
```
in this case you might need to provide all the options explicitly when we're running this command

- Create the config file and pass it to the binary
```
$ impostor -config=config.toml -app=sample
```
This would make the developer life easy by maintaining multiple config at one place and every time while running this command we just have to choose the application.

Sample config file looks like this: [link](https://github.com/ratanphayade/impostor/blob/master/config.toml)
```
[apps]
    [apps.default]
        port      = 9000
        host      = "localhost"
        mock_path = "test"

    [apps.sample]
        port      = 9001
        host      = "localhost"
```

### Writing Mocks

An application mock details should be present in single directory which should be accessed by
`mock_path`/`app`. Ex: in case of above configuration request mock details should be present in `./test/default`

There are three kind of files:
- `404.json`    - defines the response which has to be sent in case of request not found. Sample [here](https://github.com/ratanphayade/impostor/blob/master/test/app/404.json)
- `cors.json`   - deals with cors request. if not provided, then all the request will be accepted. Sample [here](https://github.com/ratanphayade/impostor/blob/master/test/app/cors.json)
- `<file>.json` - contains request and response details per API. Which means we'll be writing each API mock rules in different files.

 Sample [here](https://github.com/ratanphayade/impostor/blob/master/test/app/users_list.json)

```
{
  "method": "GET",
  "endpoint": "/users/{id}",
  "evaluators": [
    {
      "response": {
        "label": "success",
        "body": "{\"id\":1,\"firstname\":\"Hodor\",\"lastname\":\"\",\"last_seen\":\"The Door\"}",
        "latency": 0,
        "status_code": 200,
        "headers": {
          "Content-Type": "application/json"
        }
      },
      "rules": [
        {
          "target": "resource",
          "modifier": "id",
          "value": "1"
        }
      ]
    }
  ]
}
```

In general every request mock has
* `method`: Request method of the endpoint
* `endpoint`: Request URI with or without response identifier
* `evaluators`: Contains multiple response along with rule. If the evaluation of all the rule results `true`, then respective response will be used.

Evaluator contains:
- `response`: contains details of the response body, status code, header and if any latency to be introduced.
- `rules`: contains list of rules which would help identify which response to use based on the request.

every response associated with some rules, if all the rules evaluates to true then the respective response will be used. If there are no rules specified for a response then it'll be considered as default response and will be run only when no matches found.

Each rule contains:
- `target`: this will identify from where the data has to be fetched. possible values
   - `resource`: URI resource identifier
   - `params`: URL query params
   - `body`: Request body
   - `headers`: Request header
- `modifier`: attribute from the above mentioned targets. This can be plain string or flatten key format with `.` separator.
- `value`: expected value for the attribute
- `is_regex`: if `true` does a regex match on value given

Response contains:
- `label`: response identifier
- `body`: response body structure
- `status_code`: status code which has to be sent in response
- `latency`: if provided, response will be delayed by specified duration in ms.
- `headers`: response headers list. it contains a json object of <string>:<string> type.

#### Generating dynamic response

Response body can also have dynamic fields which can be copy of any `targets` mentioned above or a generated data. Currently, we are supporting
`copy` and `generate` command in response. To use these commands we should also follow some format.

- `generate`: generated the data using given options. Format `{{ generate [string|int] length }}`. Based on the given option generator will generate the data of specified length.
- `copy`: copu the data from any of request field (specified by `targets`). Format `{{ copy <target> <attribute>}}`. Here attribute can be plain string or flatten string based on the nested structure of data
Ex:
```
"body": "{\"id\":{{ generate string 14 }},\"firstname\":\"Hodor\",\"lastname\":\"\",\"status\":\"{{ copy body options.status }}\"}",
```

