package apk

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math/rand"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	b, err := ioutil.ReadFile("./1.xml")
	if err != nil {
		return
	}
	app := &APK{}
	if err := xml.Unmarshal(b, &app.manifest); err != nil {
		if app.manifest.Package == "" || app.manifest.Application.Label == nil || *app.manifest.Application.Label == "" {
			return
		}
	}
}

func TestNewXMLFile(t *testing.T) {
	app, err := GetApkInfo("D:\\temp\\temp\\laji\\111\\alibaba.apk")
	if err != nil {
		t.Fatal(err)
		return
	}
	fmt.Println(app.manifest.Package)
	fmt.Println(*app.manifest.Application.Label)
	fmt.Println(app.manifest.Activity)
}

/*
随机时间防止检测
3       5%
4	   14%
5 	   36%
6	   23%
7	   9%
8 	   8%
9	   3%
10	   1%
0-19   %1
*/
func getTimeTable(randTime int) int {
	if randTime < 5 {
		return 3
	}
	if randTime < 19 {
		return 4
	}
	if randTime < 55 {
		return 5
	}
	if randTime < 78 {
		return 6
	}
	if randTime < 87 {
		return 7
	}
	if randTime < 95 {
		return 8
	}
	if randTime < 98 {
		return 9
	}
	if randTime < 99 {
		return 10
	}
	return int(rand.Int31n(20))
}

func getRandTime() time.Duration {
	s := getTimeTable(int(rand.Int31n(100)))
	ms := rand.Int31n(1000)
	return time.Duration(s)*1000*time.Millisecond + time.Duration(ms)*time.Millisecond
}

type Node struct {
	Index      string `xml:"index,attr"`
	Text       string `xml:"text,attr"`
	ResourceId string `xml:"resource-id,attr"`
	Clickable  bool   `xml:"clickable,attr"`
	Nodes      []Node `xml:"node"`
}

type NodeTree struct {
	Nodes []Node `xml:"node"`
}

func (t NodeTree) FindValidNodes() []Node {
	nodes := make([]Node, 0)
	for _, node := range t.Nodes {
		if node.Index != "" {
			return t.Nodes
		} else {
			ns := node.FindValidNodes()
			nodes = append(nodes, ns...)
		}
	}
	return nodes
}

func (n Node) FindValidNodes() []Node {
	nodes := make([]Node, 0)
	for _, node := range n.Nodes {
		if node.Index != "" {
			nodes = append(nodes, node)
		} else {
			ns := node.FindValidNodes()
			nodes = append(nodes, ns...)
		}
	}
	return nodes
}

func (n Node) NodeNums() int {
	num := 0
	for _, node := range n.Nodes {
		num += node.NodeNums()
	}
	return num + 1
}

func (n Node) ClickAbleNums() int {
	num := 0
	for _, node := range n.Nodes {
		num += node.ClickAbleNums()
	}
	if n.Clickable {
		num = num + 1
	}
	return num
}

func TestRand(t *testing.T) {
	b, err := ioutil.ReadFile("D:\\temp\\temp\\laji\\1.xml")
	if err != nil {
		panic(err)
		return
	}
	tree := NodeTree{}
	nd := Node{}
	if err = xml.Unmarshal(b, &tree); err != nil {
		t.Fatal(err)
	}
	if err = xml.Unmarshal(b, &nd); err != nil {
		t.Fatal(err)
	}
	n := tree.FindValidNodes()
	t.Log(n)
	//t.Log(tree.Node.NodeNums())
	//t.Log(tree.Node.ClickAbleNums())
	t.Log("asd")
}
