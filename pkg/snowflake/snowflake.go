package snowflake

import (
	"sync"
	"time"

	"github.com/bwmarrin/snowflake"
)

type GenID struct {
	mu   sync.Mutex
	node *snowflake.Node
}

var G = new(GenID)

func Init(startTime string, machineID int64) (err error) {
	var st time.Time
	var node *snowflake.Node
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return
	}
	snowflake.Epoch = st.UnixNano() / 1000000
	node, err = snowflake.NewNode(machineID)
	if err != nil {
		return
	}
	G.node = node
	return
}

func (g *GenID) GetID() int64 {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.node.Generate().Int64()
}
