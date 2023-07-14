# version

<https://github.com/minio/mc>

```bash
git remote add upstream git@github.com:minio/mc.git

git fetch upstream

git merge RELEASE.2023-07-11T23-30-44Z
```

## debug

```bash
# cache
docker run -it \
--rm \
-v $PWD/:/go/src/github.com/minio/mc \
-w /go/src/github.com/minio/mc \
registry.cn-qingdao.aliyuncs.com/wod/golang:1.20 \
rm -rf vendor && go mod vendor

# build
docker run -it \
--rm \
-v $PWD/:/go/src/github.com/minio/mc \
-w /go/src/github.com/minio/mc \
registry.cn-qingdao.aliyuncs.com/wod/golang:1.20 \
bash .beagle/build.sh

# check
file .bin/linux/arm64/mc
```

## cache

```bash
# 构建缓存-->推送缓存至服务器
docker run --rm \
  -e PLUGIN_REBUILD=true \
  -e PLUGIN_ENDPOINT=$PLUGIN_ENDPOINT \
  -e PLUGIN_ACCESS_KEY=$PLUGIN_ACCESS_KEY \
  -e PLUGIN_SECRET_KEY=$PLUGIN_SECRET_KEY \
  -e DRONE_REPO_OWNER="open-beagle" \
  -e DRONE_REPO_NAME="mc" \
  -e DRONE_COMMIT_BRANCH="dev" \
  -e PLUGIN_MOUNT="vendor" \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  registry.cn-qingdao.aliyuncs.com/wod/devops-s3-cache:1.0

# 读取缓存-->将缓存从服务器拉取到本地
docker run --rm \
  -e PLUGIN_RESTORE=true \
  -e PLUGIN_ENDPOINT=$PLUGIN_ENDPOINT \
  -e PLUGIN_ACCESS_KEY=$PLUGIN_ACCESS_KEY \
  -e PLUGIN_SECRET_KEY=$PLUGIN_SECRET_KEY \
  -e DRONE_REPO_OWNER="open-beagle" \
  -e DRONE_REPO_NAME="mc" \
  -e DRONE_COMMIT_BRANCH="dev" \
  -e PLUGIN_MOUNT="vendor" \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  registry.cn-qingdao.aliyuncs.com/wod/devops-s3-cache:1.0
```
