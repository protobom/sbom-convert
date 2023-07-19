# Configuration 

Configuration search paths:

- .{{.ApplicationName}}.yaml
- .{{.ApplicationName}}/{{.ApplicationName}}.yaml
- ~/.{{.ApplicationName}}.yaml
- \<k\>/{{.ApplicationName}}/{{.ApplicationName}}.yaml

For a custom configuration location use `--config` flag with any command.

Configuration format and default values.

```yaml
{{.DefaultConfig}}```
