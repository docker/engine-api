package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/calavera/go-jsonschema-generator"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/events"
	"github.com/docker/engine-api/types/registry"
)

const (
	defaultDest = "schema/json"
	defaultTag  = "master"
)

var sources = map[string]interface{}{
	"ps":              []types.Container{},                 // docker ps
	"inspect":         types.ContainerJSON{},               // docker inspect
	"stats":           []types.Stats{},                     // docker stats
	"events":          []events.Message{},                  // docker events
	"commit":          types.ContainerCommitResponse{},     // docker commit
	"wait":            types.ContainerWaitResponse{},       // docker wait
	"create":          types.ContainerCreateResponse{},     // docker create
	"top":             types.ContainerProcessList{},        // docker top
	"cp-stat":         types.ContainerPathStat{},           // docker cp stat
	"diff":            []types.ContainerChange{},           // docker diff
	"exec-create":     types.ContainerExecCreateResponse{}, // docker exec create
	"exec-inspect":    types.ContainerExecInspect{},        // docker exec inspect
	"history":         []types.ImageHistory{},              // docker history
	"image-inspect":   types.ImageInspect{},                // docker inspect --image
	"images":          []types.Image{},                     // docker images
	"image-delete":    []types.ImageDelete{},               // docker rmi
	"search":          []registry.SearchResult{},           // docker search
	"info":            types.Info{},                        // docker info
	"login":           types.AuthResponse{},                // docker login
	"network-create":  types.NetworkCreateResponse{},       // docker network create
	"network-ls":      []types.NetworkResource{},           // docker network ls
	"network-inspect": types.NetworkResource{},             // docker network inspect
	"volume-ls":       types.VolumesListResponse{},         // docker volume ls
	"volume":          types.Volume{},                      // docker volume inspect / docker volume create
}

func main() {
	tag := flag.String("tag", defaultTag, "version tag for the schemas, `master` by default")
	flag.Parse()

	args := flag.Args()

	if len(args) > 2 {
		fmt.Printf("Usage: %s [--tag TAG] [DESTINATION-PATH]\n", args[0])
		os.Exit(1)
	}

	dest := defaultDest
	if len(args) == 2 {
		dest = args[1]
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dest, err = filepath.Abs(filepath.Join(cwd, dest))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Generating json schemas in: %s\n", dest)

	tagDest := filepath.Join(dest, *tag)
	if err := os.MkdirAll(tagDest, 0755); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for f, t := range sources {
		s := &jsonschema.Document{}
		s.Read(t)

		m, err := s.Marshal()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("- writing %s\n", f)
		err = ioutil.WriteFile(filepath.Join(tagDest, f+".json"), m, 0644)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
