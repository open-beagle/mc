# version

<https://github.com/minio/mc>

```bash
git remote add upstream git@github.com:minio/mc.git

git fetch upstream

git merge RELEASE.2024-10-02T08-27-28Z
```

## debug

```bash
# build
docker run -it \
  --rm \
  -v $PWD/:/go/src/github.com/minio/mc \
  -w /go/src/github.com/minio/mc \
  registry.cn-qingdao.aliyuncs.com/wod/golang:1.22-alpine \
  bash .beagle/build.sh

# check
file parts/linux/arm64/mc

# debug
docker run -it \
  --rm \
  -v $PWD/:/go/src/github.com/minio/mc \
  -w /go/src/github.com/minio/mc \
  -e S3_ACCESS_KEY=${PLUGIN_MINIO_ACCESS_KEY_ALIYUN} \
  -e S3_SECRET_KEY=${PLUGIN_MINIO_SECRET_KEY_ALIYUN} \
  registry.cn-qingdao.aliyuncs.com/wod/devops-minio:1.0 \
  ash -c 'mc alias set aliyun --api=S3v4 https://cache.ali.wodcloud.com $S3_ACCESS_KEY $S3_SECRET_KEY
    mkdir -p parts/2024-10-02
    mv parts/linux/amd64/mc parts/2024-10-02/mc_2024-10-02_amd64
    mv parts/linux/arm64/mc parts/2024-10-02/mc_2024-10-02_arm64
    mc cp --recursive parts/2024-10-02/ aliyun/vscode/minio/'
```

## cache

```bash
# 构建缓存-->推送缓存至服务器
docker run --rm \
  -e PLUGIN_REBUILD=true \
  -e PLUGIN_ENDPOINT=${PLUGIN_MINIO_ENDPOINT_ALIYUN} \
  -e PLUGIN_ACCESS_KEY=${PLUGIN_MINIO_ACCESS_KEY_ALIYUN} \
  -e PLUGIN_SECRET_KEY=${PLUGIN_MINIO_SECRET_KEY_ALIYUN} \
  -e DRONE_REPO_OWNER="open-beagle" \
  -e DRONE_REPO_NAME="mc" \
  -e PLUGIN_MOUNT=".git" \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  registry.cn-qingdao.aliyuncs.com/wod/devops-s3-cache:1.0

# 读取缓存-->将缓存从服务器拉取到本地
docker run --rm \
  -e PLUGIN_RESTORE=true \
  -e PLUGIN_ENDPOINT=${PLUGIN_MINIO_ENDPOINT_ALIYUN} \
  -e PLUGIN_ACCESS_KEY=${PLUGIN_MINIO_ACCESS_KEY_ALIYUN} \
  -e PLUGIN_SECRET_KEY=${PLUGIN_MINIO_SECRET_KEY_ALIYUN} \
  -e DRONE_REPO_OWNER="open-beagle" \
  -e DRONE_REPO_NAME="mc" \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  registry.cn-qingdao.aliyuncs.com/wod/devops-s3-cache:1.0
```
