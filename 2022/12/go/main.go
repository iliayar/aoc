package main

import (
	"container/heap"
	"fmt"
	"math"
	"os"
	"strings"
)

type NodeId int

type Node struct {
    height int
}

type Edge struct {
    distance int
}

type Graph struct {
    nodes map[NodeId]*Node
    edges map[NodeId]map[NodeId]*Edge
}

func makeGraph() Graph {
   return Graph {
       nodes: map[NodeId]*Node{},
       edges: map[NodeId]map[NodeId]*Edge{},
   }
}

func (g *Graph) dijkstra(start NodeId) map[NodeId]int {
    pq := make(PQueue, len(g.nodes))

    result := map[NodeId]int{}
    result[start] = 0

    for nodeId := range g.nodes {
        distance := math.MaxInt

        if g.edges[start][nodeId] != nil {
            distance = g.edges[start][nodeId].distance
        }

        pq[int(nodeId)] = &Item{
            nodeId: nodeId,
            distance: distance,
            index: int(nodeId),
        }
    }

    heap.Init(&pq)

    for pq.Len() > 0 {
        item := heap.Pop(&pq).(*Item)
        nodeId := item.nodeId

        if _, ok := result[nodeId]; ok {
            continue
        }
        result[nodeId] = item.distance

        for toId, edge := range g.edges[nodeId] {
            if _, ok := result[toId]; ok {
                continue
            }

            item := Item {
                nodeId: toId,
                distance: item.distance + edge.distance,
            }
            heap.Push(&pq, &item)
        }
    }

    return result
}

func main() {
    content, err := os.ReadFile("input.txt")
    if err != nil {
        panic(err)
    }

    contentStr := strings.Trim(string(content), "\n\r\t ")

    end := NodeId(0)

    surface := [][]int{}

    g := makeGraph()

    linesStrs := strings.Split(contentStr, "\n")
    for i, line := range linesStrs {
        curSurface := []int{}

        for j, c := range line {
            nodeId := NodeId(i * len(linesStrs[0]) + j)
            height := 0

            if c == 'S' {
                height = 0
            } else if c == 'E' {
                end = nodeId
                height = 25
            } else {
                height = int(c - 'a')
            }

            node := Node{
                height: height,
            }
            g.nodes[nodeId] = &node
            g.edges[nodeId] = map[NodeId]*Edge{}


            curSurface = append(curSurface, height)
        }

        surface = append(surface, curSurface)
    }

    for i, line := range surface {
        for j, height := range line {
            nodeId := NodeId(i * len(surface[0]) + j)

            addNode := func(toI, toJ int) {
                toId := NodeId(toI * len(surface[0]) + toJ)
                heightDiff := surface[toI][toJ] - height
                if heightDiff <= 1 {
                    e := Edge {
                        distance: 1,
                    }
                    g.edges[toId][nodeId] = &e
                }
            }
                
            if i > 0 {
                addNode(i - 1, j)
            }
            if j > 0 {
                addNode(i, j - 1)
            }
            if i < len(surface) - 1 {
                addNode(i + 1, j)
            }
            if j < len(surface[0]) - 1 {
                addNode(i, j + 1)
            }
        }
    }

    distances := g.dijkstra(end)

    res := math.MaxInt
    for i, line := range surface {
        for j, height := range line {
            nodeId := NodeId(i * len(surface[0]) + j)

            if nodeId == end {
                continue
            }

            if height == 0 {
                if distances[nodeId] > 0 && distances[nodeId] < res {
                    res = distances[nodeId]
                }
            }
        }
    }

    fmt.Println(res)
}

type Item struct {
    nodeId NodeId
    index int
    distance int
}

type PQueue []*Item

func (pq PQueue) Len() int { return len(pq) }

func (pq PQueue) Less(i, j int) bool {
    return pq[i].distance < pq[j].distance
}

func (pq PQueue) Swap(i, j int) {
    pq[i], pq[j] = pq[j], pq[i]
    pq[i].index = i
    pq[j].index = j
}

func (pq *PQueue) Push(x any) {
    n := len(*pq)
    item := x.(*Item)
    item.index = n
    *pq = append(*pq, item)
}

func (pq *PQueue) Pop() any {
    old := *pq
    n := len(old)
    item := old[n - 1]
    old[n - 1] = nil
    item.index = -1
    *pq = old[0 : n-1]
    return item
}

func (pq *PQueue) update(item *Item, distance int) {
    item.distance = distance
    heap.Fix(pq, item.index)
}
