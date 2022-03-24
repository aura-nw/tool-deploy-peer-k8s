# Run a full-node to join Public Testnets on Cosmos Ecosystem

## 1. Create a Dockerfile
Write a Dockerfile to install all required library as well as cloning the daemon repository.
```
FROM akcadmin/starport:0.18
RUN apt-get update && apt upgrade -y
RUN apt-get install vim telnet make build-essential gcc git jq chrony -y

WORKDIR /

# Build binary and initialize chain
RUN git clone -b <release-branch> <repository>
```

## 2. Create shell script
Write a `.sh` file to build the Docker image and push it to github registry.
```
#!/bin/sh
set -xe

#Login to registry
echo $GITHUB_PASSWORD | docker login ghcr.io -u $GITHUB_USERNAME --password-stdin
#Build and push image
docker build -t ${CONTAINER_RELEASE_IMAGE} -f Dockerfile/<Dockerfile-name> .
docker push ${CONTAINER_RELEASE_IMAGE}
```

## 3. Edit `yaml` file
Edit the `ci.yml` file in `/.github` folder to change image tag and adding `.sh` script file.
```
name: Continuous integration

on:
  push:
    branches: [ <current-branch> ]

jobs:

  build:

    runs-on: ubuntu-latest

    steps:
    - uses: actions/checkout@v2
    - name: Set environment variable
      run: |
        SHORT_SHA_COMMIT=$(git rev-parse --short HEAD)
        echo CONTAINER_RELEASE_IMAGE=ghcr.io/aura-nw/<image-tag>:${GITHUB_REF_NAME}_${SHORT_SHA_COMMIT} >> $GITHUB_ENV
    - name: Build the Docker image and push it to the registry
      env:
        GITHUB_USERNAME: ${{ secrets.REGISTRY_USERNAME }}
        GITHUB_PASSWORD: ${{ secrets.REGISTRY_PASSWORD }}
      run: |
        chmod 777 -R ./ci
        ./ci/build.sh
```
Replace the `build.sh` file with your shell script file.

## 4. Create a StatefulSet yaml file for deployment on Kubernetes
After the Docker image is built and uploaded to the github registry, create 2 yaml files in your Linux machine:

### a. StatfulSet:
```
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: <your-pod-name>
  namespace: cosmos
spec:
  replicas: 1
  revisionHistoryLimit: 3
  serviceName: <your-pod-name>
  selector:
    matchLabels:
      app: <your-pod-name>
  template:
    metadata:
      labels:
        app: <your-pod-name>
    spec:
      containers:
      - image: ghcr.io/aura-nw/<your-image-tag>
        imagePullPolicy: IfNotPresent
        name: <your-pod-name>
        command: 
        - /bin/sh
        - -c
        - |
          cd <daemon-repo>;
          cd ..;
          FILE=/root/<daemon-binary-folder>/config/genesis.json;
          if [ ! -f "$FILE" ]; then
            junod init <moniker> --chain-id <chain-id>;
            curl <genesis-file-link> > /root/<daemon-binary-folder>/config/genesis.json
            sed -i 's/seeds = ""/seeds = <seeds>"/' /root/<daemon-binary-folder>/config/config.toml
            sed -i 's/minimum-gas-prices = "0stake"/minimum-gas-prices = "<minimum-gas-prices>"/' /root/<daemon-binary-folder>/config/app.toml
            sed -i 's/persistent_peers = ""/persistent_peers = "<persistent-peers>"/' /root/<daemon-binary-folder>//config/config.toml
          fi
          <deamon-binary> start;
        volumeMounts:
        - mountPath: /root
          name: pvc
      imagePullSecrets:
        - name: regcred-github

  volumeClaimTemplates:
  - apiVersion: v1
    kind: PersistentVolumeClaim
    metadata:
      name: pvc
    spec:
      accessModes:
      - ReadWriteOnce
      resources:
        requests:
          storage: 200Gi
```
For the command to init and start the node, please refer to the official document of your desired testnet:

- Gaia (ATOM): https://hub.cosmos.network/main/hub-tutorials/join-testnet.html
- Juno (JUNO): https://docs.junonetwork.io/validators/joining-the-testnets
- Osmosis (OSMO): https://docs.osmosis.zone/developing/network/join-testnet.html

The pod running full-node will need a PersistentVolumeClaims to store the data synced from network so that when the pod died or restarted, the data won't be deleted. <br>
Therefore, the process of initialize the node and start the node need to be set in the yaml file to prevent overriding data inside `/root` folder when creating PersistentVolumeClaims.

### b. Secret
> Ignore this step if already has a Secret.

Your StatefulSet will need a Secret to pull Docker image from github registry.
```
apiVersion: v1
kind: Secret
metadata:
  name: <secret-name>
  namespace: <namespace>
data:
  .dockerconfigjson: <your-secret>
type: kubernetes.io/dockerconfigjson
```

Lastly, apply the yaml files for deployment of full-node on Kubernetes.