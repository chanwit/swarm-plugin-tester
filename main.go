package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/chanwit/tentacle"
	"github.com/samalba/dockerclient"
)

type StopWatch struct {
	start, stop time.Time
}

func Start() time.Time {
	return time.Now()
}

func Stop(start time.Time) *StopWatch {
	watch := StopWatch{start: start, stop: time.Now()}
	return &watch
}

func (self *StopWatch) Milliseconds() uint32 {
	return uint32(self.stop.Sub(self.start) / time.Millisecond)
}

type StrategyPlugin struct {
	Name string
}

func (s *StrategyPlugin) Initialize() error {
	return nil
}

func (s *StrategyPlugin) PlaceContainer(config *dockerclient.ContainerConfig, nodes []*tentacle.Node) (*tentacle.Node, error) {
	exeName := "swarm-strategy-" + s.Name
	cmd := exec.Command(exeName)

	request := &tentacle.Request{config, nodes}
	b, err := json.Marshal(request)
	cmd.Stdin = strings.NewReader(string(b))
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
	}
	response := &tentacle.Response{}
	err = json.Unmarshal(out.Bytes(), &response)
	if err != nil {
		log.Fatal(err)
	}

	if response.Error == "" {
		err = nil
	} else {
		err = errors.New(response.Error)
	}

	for _, node := range nodes {
		if node.ID == response.Node.ID {
			return node, err
		}
	}

	return nil, errors.New("Cannot find a node to place the new container")
}

var count chan uint32
var strategies map[string]tentacle.PlacementStrategy

func test() {
	s := strategies[os.Args[1]]
	s.Initialize()

	config := &dockerclient.ContainerConfig{
		Hostname: "node1",
	}
	nodes := []*tentacle.Node{
		&tentacle.Node{ID: "AABB:CCDD:EEFF"},
		&tentacle.Node{ID: "AABB:1CDD:EEFF"},
		&tentacle.Node{ID: "AABB:2CDD:EEFF"},
		&tentacle.Node{ID: "AABB:3CDD:EEFF"},
		&tentacle.Node{ID: "AABB:4CDD:EEFF"},
		&tentacle.Node{ID: "AABB:5CDD:EEFF"},
		&tentacle.Node{ID: "AABB:6CDD:EEFF"},
		&tentacle.Node{ID: "AABB:6CDD:EEFF"},
		&tentacle.Node{ID: "AABB:6CDD:EEFF"},
		&tentacle.Node{ID: "AABB:6CDD:EEFF"},
		&tentacle.Node{ID: "AABB:6CDD:EEFF"},
		&tentacle.Node{ID: "AABB:6CDD:EEFF"},
		&tentacle.Node{ID: "AABB:6CDD:EEFF"},
		&tentacle.Node{ID: "AABB:7CDD:EEFF"},
		&tentacle.Node{ID: "AABB:7CDD:EEFF"},
		&tentacle.Node{ID: "AABB:7CDD:EEFF"},
		&tentacle.Node{ID: "AABB:7CDD:EEFF"},
		&tentacle.Node{ID: "AABB:7CDD:EEFF"},
	}

	node, _ := s.PlaceContainer(config, nodes)
	if node.ID == "AABB:CCDD:EEFF" {
		count <- 1
	}
}

func main() {
	strategies = make(map[string]tentacle.PlacementStrategy)

	infos, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	for _, info := range infos {
		if strings.HasPrefix(info.Name(), "swarm-strategy-") {
			name := strings.TrimPrefix(info.Name(), "swarm-strategy-")
			if runtime.GOOS == "windows" {
				name = strings.TrimSuffix(name, ".exe")
			}
			strategies[name] = &StrategyPlugin{Name: name}
		}
	}

	count = make(chan uint32)
	start := Start()

	N := uint32(1000)
	for i := uint32(0); i < N; i++ {
		go test()
	}
	sum := uint32(0)
	for ; sum < N; sum = sum + <-count {
	}

	w := Stop(start)
	fmt.Printf("Success = %d%%\n", sum*100/N)
	fmt.Printf("total = %v ms\n", w.Milliseconds())
	fmt.Printf("each  = %v ms\n", w.Milliseconds()/N)
}
