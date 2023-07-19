# Configuration 

Configuration search paths:

- .protobom.yaml
- .protobom/protobom.yaml
- ~/.protobom.yaml
- \<k\>/protobom/protobom.yaml

For a custom configuration location use `--config` flag with any command.

Configuration format and default values.

```yaml
logger:
    verbose: 2
translate:
    format: cyclonedx
    output: /dev/null
    encoding: json
```
