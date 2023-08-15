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
- `-c` Configurations / additional inputs (optional)
    - Input format and version
    - Output format and version (e.g. CDX 1.5)
    - Fail on non-perfect translation: true-false
    - CDX software identifiers priority (purl\cpe\swid)
    - User-determined fields: list of key-path and values
- `-h` provides the help information
- `-o` output filename (optional)
    - By default: the output filename is the same as the input filename, appended with “-spdx” or “-cyclonedx” depending on the format output.
    - A string that becomes the filename as a .json file
    - ❓ should the input be “output.json” or just “output” (DB prefers the former)


## Outputs

- An SBOM of the desired format in a .json file with the given input filename.
- Written to stdout unless `-o` is specified


### User Interaction

```jsx
convert-sbom input-sbom.json 

convert-sbom input-sbom.xml -o output-sbom.json -f spdx-3.0
```


## Error Handling

We should provide helpful, clear error messages and guidance if any exceptions are found, including:

- Input file isn’t a valid .json file
- Input file isn’t a valid (CDX || SPDX) SBOM
- Execution errors


# Open Questions

| Question | Answer | Date Answered |
| --- | --- | --- |
| What language should this be written in?
- [Go, OPA/Rego] | Start with Golang |  |
| Should we support .spdx? | No |  |
| What’s the expected behavior around versioning for CDX and SPDX? | CDX 1.4+
SPDX 2.2, 2.3+ |  |
|  |  |  |


# Out of Scope

- Support for SWID format
- Support for XML