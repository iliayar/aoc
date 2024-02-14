package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type NodeId string

type Node struct {
	id   NodeId
	flow int
	to   []NodeId
}

func parseNode(line string) Node {
	line, _ = strings.CutPrefix(line, "Valve ")

	nodeId := NodeId(line[0:2])

	line = line[2:]
	line, _ = strings.CutPrefix(line, " has flow rate=")

	var flowStr string
	flowStr, line, _ = strings.Cut(line, ";")

	flow, _ := strconv.Atoi(flowStr)

	line, _ = strings.CutPrefix(line, " tunnels lead to valves ")
	line, _ = strings.CutPrefix(line, " tunnel leads to valve ")
	toStrs := strings.Split(line, ", ")

	to := []NodeId{}
	for _, toStr := range toStrs {
		to = append(to, NodeId(toStr))
	}

	return Node{
		id:   nodeId,
		flow: flow,
		to:   to,
	}
}

type BitSet int64

func (bs BitSet) get(i int) bool {
	return bs&(1<<i) > 0
}

func (bs BitSet) set(i int, v bool) BitSet {
	if !v {
		return bs & ^(1 << i)
	} else {
		return bs | (1 << i)
	}
}

const MAX_TIME = 26
const PLAYERS = 2

type Graph map[NodeId]Node

func (g Graph) search(nodeId NodeId) int {
	m := []NodeId{}
	for nodeId := range g {
		m = append(m, nodeId)
	}
	mi := map[NodeId]int{}
	for i, n := range m {
		mi[n] = i
	}

	switched := BitSet(0)
	mem := Mem{}
	ans := g.searchImpl(nodeId, mem, mi, nodeId, switched, MAX_TIME, PLAYERS-1)

	return ans
}

type Mem map[MemKey]int

type MemKey struct {
	opH int64
	n   int
	t   int
	p   int
}

func makeMemKey(opH BitSet, n, t, p int) MemKey {
	return MemKey{
		opH: int64(opH),
		n:   n,
		t:   t,
		p:   p,
	}
}

func (g Graph) searchImpl(s NodeId, mm Mem, m map[NodeId]int, n NodeId, op BitSet, t int, p int) int {
	if t == 0 {
		if p == 0 {
			return 0
		} else {
			return g.searchImpl(s, mm, m, s, op, MAX_TIME, p-1)
		}
	}

	mk := makeMemKey(op, m[n], t, p)
	ans, found := mm[mk]
	if found {
		return ans
	}

	ans = math.MinInt

	o := op.get(m[n])
	if !o && g[n].flow > 0 {
		nOp := op.set(m[n], true)

		nAns := g[n].flow*(t-1) + g.searchImpl(s, mm, m, n, nOp, t-1, p)
		if nAns > ans {
			ans = nAns
		}
	}

	for _, to := range g[n].to {
		nAns := g.searchImpl(s, mm, m, to, op, t-1, p)
		if nAns > ans {
			ans = nAns
		}
	}

	mm[mk] = ans
	return ans
}

func (g Graph) searchBfs(nodeId NodeId, time int, opened map[NodeId]bool) (NodeId, int) {
	type E struct {
		nodeId NodeId
		dist   int
	}

	was := map[NodeId]bool{}
	q := []E{{nodeId, 0}}

	mFlow := math.MinInt
	var mE E

	for len(q) > 0 {
		var e E
		e, q = q[0], q[1:]

		_, w := was[e.nodeId]
		if w {
			continue
		}
		was[e.nodeId] = true

		_, w = opened[e.nodeId]
		if !w {
			cFlow := (time - e.dist) * g[e.nodeId].flow
			if cFlow > mFlow {
				mFlow = cFlow
				mE = e
			}
		}

		for _, to := range g[e.nodeId].to {
			q = append(q, E{to, e.dist + 1})
		}
	}

	return mE.nodeId, mE.dist
}

func main() {
	content, _ := os.ReadFile("input.txt")
	contentStr := strings.Trim(string(content), "\n\r\t ")

	nodes := map[NodeId]Node{}

	const startId = NodeId("AA")

	for _, line := range strings.Split(contentStr, "\n") {
		node := parseNode(line)
		nodes[node.id] = node
	}

	fmt.Println(Graph(nodes).search(startId))
}
