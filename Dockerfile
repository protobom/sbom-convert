FROM alpine:3.19@sha256:51b67269f354137895d43f3b3d810bfacd3945438e94dc5ac55fdac340352f48 as build

COPY sbom-convert_*.apk /tmp/
RUN apk add --no-cache --allow-untrusted /tmp/sbom-convert_*.apk

FROM cgr.dev/chainguard/static@sha256:68b8855b2ce85b1c649c0e6c69f93c214f4db75359e4fd07b1df951a4e2b0140 as runtime

COPY --from=build /usr/bin/sbom-convert /sbom-convert

CMD [ "-h" ]
ENTRYPOINT [ "/sbom-convert" ]
