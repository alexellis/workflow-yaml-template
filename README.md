# workflow-yaml-template
Rudimentary YAML workflow runner for OpenFaaS

## Example

* Generate a workflow:

```bash
faas template pull https://github.com/alexellis/workflow-yaml-template
faas new --lang workflow youtube2gif --prefix=DOCKER_HUB_NAME
```

* Customise `youtube2gif/workflow.yml`

```yaml
workflow:
  name: youtube2gif
  gateway_url: http://gateway:8080/
  steps:
    - name: download-video
      function: youtube-dl
      method: POST

    - name: convert-to-gif
      function: gif-maker
      method: POST
```

* Add the dependent functions to your `stack.yml`

```yaml
provider:
  name: openfaas

functions:
  ## Workflow runner
  youtube2gif:
    lang: workflow
    handler: ./youtube2gif
    image: YOUR_HUB_USERNAME/youtube-gif
    environment:
      write_timeout: 2m
      read_timeout: 2m
      combine_output: false

  ## Dependent functions
  gif-maker:
    skip_build: true
    image: functions/gif-maker:latest
    environment:
      write_timeout: 65s
      read_timeout: 65s
      combine_output: false

  youtube-dl:
    skip_build: true
    image: rgee0/faas-youtubedl:0.4
    environment:
      write_timeout: 65s
      read_timeout: 65s
      combine_output: false
```

* Build / deploy / test:

```
faas-cli up -f youtube2gif.yml 
```

* Test it out:

```
echo -n "https://www.youtube.com/watch?v=0Bmhjf0rKe8" | faas invoke youtube2gif "https://www.youtube.com/watch?v=0Bmhjf0rKe8" > cat-surprise.gif
```
