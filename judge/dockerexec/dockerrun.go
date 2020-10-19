package dockerexec

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"time"
)

func GoDocker(hostPath string, containerPath string, workDir string, runName string, limitTime string, limitMem string, inFile string, outFile string, msgFile string, dockerName string) string {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	if hostPath == "" {
		hostPath = DockerRunPath
	}
	if containerPath == "" {
		containerPath = DockerWorkPath
	}
	if workDir == "" {
		workDir = DockerWorkPath
	}

	workDir = DockerWorkPath
	/*runName = "/b"*/
	if limitTime == "" {
		limitTime = LimitTime
	}
	if limitMem == "" {
		limitMem = LimitMemory
	}
	/*
		inFile = "/b.in"
		outFile = "/b.out"
		msgFile = "/msgb.txt"*/
	fmt.Println(hostPath)
	fmt.Println(containerPath)
	fmt.Println(runName)
	fmt.Println(limitTime)
	fmt.Println(limitMem)
	fmt.Println(inFile)
	fmt.Println(outFile)
	fmt.Println(msgFile)

	resp, err := cli.ContainerCreate(ctx,
		&container.Config{
			Image: "gcc",
			Cmd: []string{
				DockerWorkPath + "/judgepro", DockerWorkPath + "/" + runName, limitTime, limitMem, DockerWorkPath + "/" + inFile, DockerWorkPath + "/" + outFile, DockerWorkPath + "/" + msgFile,
			},
			Tty:         false,
			AttachStdin: true,
			WorkingDir:  workDir,
		},
		&container.HostConfig{
			Mounts: []mount.Mount{
				mount.Mount{
					Type:   mount.TypeBind,
					Source: hostPath,
					Target: containerPath,
				},
			},
		}, nil, dockerName)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.ID)

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
	return resp.ID
}

func DelDocker(id string) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	err = cli.ContainerRemove(ctx, id,
		types.ContainerRemoveOptions{
			RemoveVolumes: false,
			RemoveLinks:   false,
			Force:         false,
		},
	)
	if err != nil {
		panic(err)
	}
}

func StopDocker(id string) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	x := time.Duration(0)
	err = cli.ContainerStop(ctx, id,
		&x,
	)
	if err != nil {
		panic(err)
	}
}
