windows os:
set COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1
docker build -t jasonsoft/starter -f build/docker/starter.dockerfile .