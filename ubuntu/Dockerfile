FROM bitriseio/docker-bitrise-base:latest

RUN apt update
RUN apt install -y cmake wget

WORKDIR /sources
COPY build-libgit2-deb.sh .

VOLUME /products

CMD ["/bin/bash", "build-libgit2-deb.sh"]
