#!/bin/bash

docker run -v $(pwd):/bitrise/src bitriseio/docker-bitrise-base:latest /bin/bash -c "bitrise update && bitrise run test"