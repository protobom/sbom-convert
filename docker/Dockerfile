FROM alpine:latest AS build
RUN apk --no-cache add ca-certificates coreutils

FROM scratch

ARG DISCRIPTION
ARG BUILD_DATE
ARG BUILD_VERSION
ARG VCS_REF
ARG VCS_URL
ARG PROJECT_NAME
ENV PROJECT_NAME=$PROJECT_NAME

LABEL org.opencontainers.image.created=$BUILD_DATE
LABEL org.opencontainers.image.title=$PROJECT_NAME
LABEL org.opencontainers.image.description=$DISCRIPTION
LABEL org.opencontainers.image.source=$VCS_URL
LABEL org.opencontainers.image.revision=$VCS_REF
# LABEL org.opencontainers.image.vendor=""
LABEL org.opencontainers.image.version=$BUILD_VERSION

LABEL io.artifacthub.package.readme-url="https://raw.githubusercontent.com/bom-squad/$PROJECT_NAME/README.md"
LABEL io.artifacthub.package.name=$PROJECT_NAME
LABEL io.artifacthub.package.license=Apache-2.0 
# LABEL io.artifacthub.package.vendor=""


COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build / /
COPY . /home/
WORKDIR /home/

ENV PATH="/home/:${PATH}"
ENTRYPOINT [ "/home/docker/entrypoint.sh", "/home/${PROJECT_NAME}" ] 