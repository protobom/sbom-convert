FROM alpine:3.20@sha256:0a4eaa0eecf5f8c050e5bba433f58c052be7587ee8af3e8b3910ef9ab5fbe9f5 as build

COPY sbom-convert_*.apk /tmp/
RUN apk add --no-cache --allow-untrusted /tmp/sbom-convert_*.apk

FROM cgr.dev/chainguard/static@sha256:0fa3935a85aa2349cc89d9715d891c318f700ba951f3945610a2b90c6b0d5e76 as runtime

COPY --from=build /usr/bin/sbom-convert /sbom-convert

CMD [ "-h" ]
ENTRYPOINT [ "/sbom-convert" ]
