package internal

import (
	"fmt"
	"os"
	"path/filepath"
)

var (
	Deno_Dockerfile = `FROM denoland/deno

RUN mkdir -p /app
WORKDIR /app

USER deno
COPY . /app

RUN deno cache src/mod.ts

CMD ["run", "--unstable", "--allow-all", "src/mod.ts"]`
	Node_Dockerfile = `FROM node:lts-alpine

RUN mkdir -p /app
WORKDIR /app

COPY package.json /app/
RUN npm install

COPY . /app

RUN npm run build

CMD ["node", "dist/index.js"]`
	Node_DockerIgnore = `node_modules
build
dist
dest`
)

func AddDockerFiles(platform, dir string) {
	switch platform {
	case "Node":
		err1 := os.WriteFile(filepath.Join(dir, "Dockerfile"), []byte(Node_Dockerfile), 0777)
		err2 := os.WriteFile(filepath.Join(dir, ".dockerignore"), []byte(Node_DockerIgnore), 0777)
		if err1 != nil || err2 != nil {
			fmt.Println(" - Failed to add Docker files.")
		}
		break
	case "Deno":
		err := os.WriteFile(filepath.Join(dir, "Dockerfile"), []byte(Deno_Dockerfile), 0777)
		if err != nil {
			fmt.Println(" - Failed to add Docker files.")
		}
		break
	}
}
