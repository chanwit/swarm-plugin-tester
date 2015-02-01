package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/docker/swarm/scheduler/strategy/plugin"
	"github.com/samalba/dockerclient"
)

func test(pluginType string, pluginName string, N uint32) {
	cmd := exec.Command("swarm-" + pluginType + "-" + pluginName)
	go cmd.Run()

	// if pluginType == strategy
	// if pluginType == discovery
	client, err := plugin.NewClient(pluginName)

	r := new(int)
	err = client.Call("Rpc.Initialize", r, r)
	if err != nil {
		panic(err)
	}

	config := &dockerclient.ContainerConfig{
		Hostname: "node1",
	}
	nodes := make([]*plugin.Node, 1000)
	for i := range nodes {
		nodes[i] = &plugin.Node{ID: fmt.Sprintf("AABB:%04d:EEFF", i)}
	}

	start := Start()

	sum := uint32(0)
	for i := uint32(0); i < N; i++ {
		req := &plugin.StrategyPluginRequest{config, nodes}
		reply := &plugin.Node{}
		err = client.Call("Rpc.PlaceContainer", req, reply)
		if reply.ID == "AABB:0000:EEFF" {
			sum++
		}
	}
	w := Stop(start)
	fmt.Printf("Success = %d%%\n", sum*100/N)
	fmt.Printf("total = %v ms\n", w.Milliseconds())
	fmt.Printf("each  = %f ms\n", float64(w.Milliseconds())/float64(N))

	err = cmd.Process.Kill()
	if err != nil {
		fmt.Print(err.Error())
	}
}

func main() {
	i, _ := strconv.Atoi(os.Args[2])
	str := strings.Split(os.Args[1], ":")
	test(str[0], str[1], uint32(i))
}
