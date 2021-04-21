FROM ubuntu:18.04

ENV VEGETA_VERSION 12.8.4
ENV VEGETA_RATE 100
ENV VEGETA_DURATION 60s

ENV PAUSE_WAIT 10
ENV MAX_LOOP 5000

LABEL \
    maintainer="Benoit Moussaud <mail@moussaud.org>" \
    org.opencontainers.image.title="vegeta_runner" \
    org.opencontainers.image.description="Docker image for the vegeta_runner HTTP load testing tool in a *loop*. based on the work of https://github.com/peter-evans/vegeta-docker" \
    org.opencontainers.image.authors="Benoit Moussaud <mail@moussaud.org>" \
    org.opencontainers.image.url="https://github.com/bmoussaud/micropets-app" \
    org.opencontainers.image.vendor="https://moussaud.org" \
    org.opencontainers.image.licenses="MIT" \
    app.tag="vegeta$VEGETA_VERSION"

RUN apt-get update && apt-get install jq ca-certificates supervisor openssl wget -y

RUN set -ex \
    && wget -q "https://github.com/tsenart/vegeta/releases/download/v${VEGETA_VERSION}/vegeta_${VEGETA_VERSION}_linux_amd64.tar.gz" -O /tmp/vegeta.tar.gz \
    && cd bin \
    && tar xzf /tmp/vegeta.tar.gz \
    && rm /tmp/vegeta.tar.gz 

RUN mkdir -p /var/run/vegeta /var/run/vegeta /var/log/supervisor
COPY vegeta-runner.sh /bin
RUN chmod +x /bin/vegeta-runner.sh
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf
CMD ["/usr/bin/supervisord"]
