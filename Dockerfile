FROM bitriseio/docker-bitrise-base:latest

WORKDIR $HOME
RUN bitrise update
COPY . $HOME

ENTRYPOINT ["bitrise", "run", "test"]
