FROM alpine:3.18@sha256:82d1e9d7ed48a7523bdebc18cf6290bdb97b82302a8a9c27d4fe885949ea94d1

ARG NAME=sbom-convert
ENV NAME=${NAME}

RUN apk add --no-cache \
	bash \
	docker-cli \
	tini

COPY scripts/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/sbin/tini", "--", "/entrypoint.sh"]
CMD [ "-h" ]

COPY ${NAME}_*.apk /tmp/
RUN apk add --no-cache --allow-untrusted /tmp/${NAME}_*.apk
