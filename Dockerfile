FROM alpine:3.24@sha256:294b683cb724975bec92580e1e685676bd4b50bda910ddb8c51d4cabeaec77e6 as build

COPY sbom-convert_*.apk /tmp/
RUN apk add --no-cache --allow-untrusted /tmp/sbom-convert_*.apk

FROM cgr.dev/chainguard/static@sha256:54b5f14edad6f2554774d499539e837588a501bd5bf22ac6cd97b1068c35d44d as runtime

COPY --from=build /usr/bin/sbom-convert /sbom-convert

CMD [ "-h" ]
ENTRYPOINT [ "/sbom-convert" ]
