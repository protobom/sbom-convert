FROM alpine:3.19@sha256:51b67269f354137895d43f3b3d810bfacd3945438e94dc5ac55fdac340352f48 as build

COPY sbom-convert_*.apk /tmp/
RUN apk add --no-cache --allow-untrusted /tmp/sbom-convert_*.apk

FROM cgr.dev/chainguard/static@sha256:bda736ff75ec4b76ecd120d6ccf0b9428419cf4568a5d786a7d77ba23a9f6992 as runtime

COPY --from=build /usr/bin/sbom-convert /sbom-convert

CMD [ "-h" ]
ENTRYPOINT [ "/sbom-convert" ]
