Last updated: 8/15/2023

# 🎯 Objective

Collaboratively develop a command line interface (cli) to easily convert between (the most common) SBOM formats, CycloneDX and SPDX.


# 🗒️ Primary User Workflows

1. Users can convert SBOMs from SPDX to CycloneDX.*
2. Users can convert SBOMs from CycloneDX to SPDX.*

*this describes the desired user flow. More nuance will have to be considered when it comes to converting between specific SBOM formats.


# 🗒️ Requirements & Scope

## Supported SBOM Formats

SPDX: 2.2, 2.3, 3.0
CDX: 1.4, 1.5

## Inputs & Parameters

Inputs

- (implied/positional) Input SBOM (required)
    - Files must be valid .json document
    - Files must be valid (CDX || SPDX) SBOM
    - To start, the cli will only accept one input at a time
- `-e`, `--encoding`: (string) The output encoding [spdx: [text, json] cyclonedx: [json] (default "json")]
- `-f`, `--format`: (string) The output format [spdx, spdx-2.2, spdx-2.3, spdx3, spdx-3.0.1, cyclonedx, cyclonedx-1.0 through cyclonedx-1.7]. Plain `spdx` selects the latest SPDX 2.x (2.3); `spdx3` selects SPDX 3.0.1
- `-h`, `--help`:` help for convert
- `-o`, `--output`:  (string) Path to write the converted SBOM. Default: stdout. If just a string is provided, the cli will append ".json" by default. Otherwise, users can specify full filenames+extensions, like myBom.spdx.

Global Flags:
- `-c`, `--config`: (string) Path to config file
- `-v`, `--verbose`: log verbosity level (-v=info, -vv=debug, -vvv=trace)




## Outputs

- An SBOM of the desired format. .json as default, but other formats can be specified.
- Written to stdout unless `-o` is specified


### User Interaction

```jsx
sbom-convert input-sbom.json

sbom-convert -o output-sbom.json -f spdx-3.0
```


## Error Handling

We should provide helpful, clear error messages and guidance if any exceptions are found, including:

- Input file isn’t a valid .json file
- Input file isn’t a valid (CDX || SPDX) SBOM
- Execution errors


# Open Questions

| Question | Answer | Date Answered |
| --- | --- | --- |
| What language should this be written in? [Go, OPA/Rego] | Start with Golang |  |
| Should we support .spdx? | No |  |
| What’s the expected behavior around versioning for CDX and SPDX? | CDX 1.4+
SPDX 2.2, 2.3+ |  |
|  |  |  |


# Out of Scope

- Support for SWID format
- Support for XML
