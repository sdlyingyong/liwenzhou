package snowflake

import (
	sf "github.com/bwmarrin/snowflake"
	"time"
)

var (
	node *sf.Node
)

func Init(startTime string, machineID int64) (err error) {
	var st time.Time
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return
	}

	sf.Epoch = st.UnixNano() / 1000000
	// Create a new Node with a Node number of machineID
	node, err = sf.NewNode(machineID)
	return
}

func GenId() (int64, error) {
	return node.Generate().Int64(), nil
}
