FROM alpine:3.19@sha256:51b67269f354137895d43f3b3d810bfacd3945438e94dc5ac55fdac340352f48 as build

COPY sbom-convert_*.apk /tmp/
RUN apk add --no-cache --allow-untrusted /tmp/sbom-convert_*.apk

FROM cgr.dev/chainguard/static@sha256:8665c8a9fcdab0f8afc09533ee23287c7870de26064d464a10e3baa52f337734 as runtime

COPY --from=build /usr/bin/sbom-convert /sbom-convert

CMD [ "-h" ]
ENTRYPOINT [ "/sbom-convert" ]
