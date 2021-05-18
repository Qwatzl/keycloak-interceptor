package main

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"keycloak-interceptor/oidc"
	"log"
	"net/http"
	"testing"
)

func serviceMock(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	_, _ = w.Write([]byte("OK"))
}

func start() {
	mux := http.NewServeMux()

	finalHandler := http.HandlerFunc(serviceMock)
	kInterceptor := oidc.Init("keycloak.json")
	mux.Handle("/users", kInterceptor.Intercept(finalHandler))

	log.Fatal(http.ListenAndServe(":8080", mux))
}

func startContainer(cli *client.Client, ctx context.Context, name string, port string) string {
	config := container.Config{
		Image:        name,
		ExposedPorts: nat.PortSet{"8080": struct{}{}},
	}

	body, _ := cli.ImagePull(ctx, name, types.ImagePullOptions{})
	body.Close()

	hostConfig := container.HostConfig{
		PortBindings: map[nat.Port][]nat.PortBinding{nat.Port("8080"): {{HostIP: "127.0.0.1", HostPort: port}}},
	}
	networkConfig := network.NetworkingConfig{}
	c, err := cli.ContainerCreate(ctx, &config, &hostConfig, &networkConfig, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, c.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	return c.ID
}

func stopContainer(cli *client.Client, ctx context.Context, containerID string) {
	if err := cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{RemoveVolumes: true, Force: true}); err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	log.Println("Start environment...")
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	id := startContainer(cli, ctx, "jboss/keycloak", "8000")

	m.Run()

	log.Println("Stop environment...")
	stopContainer(cli, ctx, id)
}

func TestConfidentialClient(t *testing.T) {
	log.Println("Start server")
}
