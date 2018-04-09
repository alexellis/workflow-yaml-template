# workflow-yaml-template
Rudimentary YAML workflow runner for OpenFaaS

## Example

* Generate a workflow:

```
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

  ## Dependent functions
  gif-maker:
    skip_build: true
    image: functions/gif-maker:latest
    environment:
      write_timeout: 65
      read_timeout: 65

  youtube-dl:
    skip_build: true
    image: alexellis2/faas-youtubedl:latest
    environment:
      write_timeout: 65
      read_timeout: 65
```