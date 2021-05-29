#!/bin/ash

apk add curl jq

RUNNER_ARCH=$(arch)
RUNNER_ARCH=${RUNNER_ARCH/x86_/amd}
RUNNER_ARCH=${RUNNER_ARCH/aarch/arm}
BUILDX_VER=$(curl -ks https://api.github.com/repos/docker/buildx/releases/latest | jq -r '.name')
mkdir -p "$HOME/.docker/cli-plugins/"
wget -O "$HOME/.docker/cli-plugins/docker-buildx" "https://github.com/docker/buildx/releases/download/${BUILDX_VER}/buildx-${BUILDX_VER}.linux-${RUNNER_ARCH}"
chmod a+x "$HOME/.docker/cli-plugins/docker-buildx"
echo -e '{\n  "experimental": "enabled"\n}' | tee "$HOME/.docker/config.json"
docker run --rm --privileged multiarch/qemu-user-static --reset -p yes

docker buildx create --use --name builder
docker buildx inspect --bootstrap builder
docker buildx install

echo "${DOCKER_TOKEN}" | docker login -u "${DOCKER_USERNAME}" --password-stdin
docker buildx build --push --cache-from=type=local,src=cache --platform "linux/amd64,linux/arm64,linux/arm/v7" -t "${REGISTRY_IMAGE}:${CI_COMMIT_SHORT_SHA}" -t "${REGISTRY_IMAGE}:${1}" .
