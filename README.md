# workflow-yaml-template
Rudimentary YAML workflow runner for OpenFaaS

## Example

* Generate a workflow:

```bash
faas template pull https://github.com/alexellis/workflow-yaml-template
faas new --lang workflow youtube2gif
```

* Customise:

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

* Custom `stack.yml`

```yaml
provider:
  name: faas
  gateway: http://localhost:8080

functions:
  ## Workflow runner
  youtube-gif:
    lang: workflow
    handler: ./youtube-gif
    image: youtube-gif
    environment:
      write_timeout: 2m
      read_timeout: 2m

  ## Dependent functions
  gif-maker:
    skip_build: true
    image: functions/gif-maker:latest
    environment:
      write_timeout: 65s
      read_timeout: 65s

  youtube-dl:
    skip_build: true
    image: rgee0/faas-youtubedl:0.4
    environment:
      write_timeout: 65s
      read_timeout: 65s
```

* Build / deploy / test:

```
faas build && faas deploy


echo -n "https://www.youtube.com/watch?v=0Bmhjf0rKe8" | faas invoke youtube-gif "https://www.youtube.com/watch?v=0Bmhjf0rKe8" > cat-surprise.gif
```
