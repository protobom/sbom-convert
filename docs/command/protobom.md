## protobom

translate an SBOM into another format

### Synopsis

translate an SBOM into another format

```
protobom [path] [flags]
```

### Optional flags 
Flags for `protobom`


| Short | Long | Description | Default |
| --- | --- | --- | --- |
| -c | --config | Configuration file path | |
| -E | --encoding | Select encoding, options=map[cyclonedx:[json] spdx:[text json]] | "json" |
| -f | --format | Select Formats, options=[cyclonedx spdx] | |
| -h | --help | help for protobom | |
| -D | --level | Log depth level, options=[panic fatal error warning info debug trace] | |
| -o | --output | Output path | |
| -q | --quiet | Suppress all logging output | |
| | --structured | Enable structured logger | |
| -V | --ver | Select Specific version, options=map[cyclonedx:[1.4 1.5] spdx:[2.2 2.3]] (default map[cyclonedx:1.4 spdx:2.3]) | |
| -v | --verbose | Log verbosity level [-v,--verbose=1] = info, [-vv,--verbose=2] = debug | |


### Examples for running `protobom`

```

	protobom  sbom.spdx.json                        translate SPDX to CycloneDX (default)
	protobom  sbom.spdx.json -f cyclonedx           translate SPDX to CycloneDX
	protobom  sbom.cdx.json  -f spdx                translate CycloneDX to SPDX
	protobom  sbom.cdx.json  -o sbom.spdx.json      output sbom to file

```

