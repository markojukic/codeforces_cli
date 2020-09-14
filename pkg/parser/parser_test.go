package parser

import (
	"os"
	"testing"
)

func TestAll(t *testing.T) {
	file, err := os.Open("parser_test.html")
	if err != nil {
		t.Fatal(err)
	}
	root, err := Parse(file)
	if err != nil {
		t.Fatal(err)
	}

	//(*Node) GetNodeByTag(string) *Node
	if div := root.GetNodeByTag("div"); div == nil {
		t.Fatal()
	} else if div.GetAttr("id") != "div0" {
		t.Fatal()
	}

	//(*Node) GetNodesByTag(string) *Node
	if divs := root.GetNodesByTag("div"); len(divs) != 4 {
		t.Fatal()
	} else {
		if divs[0].GetAttr("id") != "div0" ||
			divs[1].GetAttr("id") != "div1" ||
			divs[2].GetAttr("id") != "div2" ||
			divs[3].GetAttr("id") != "div3" {
			t.Fatal()
		}
	}

	//(*Node) GetNodeByTagAttr(string) *Node
	if div := root.GetNodeByTagAttr("div", "class", "1"); div == nil {
		t.Fatal()
	} else if div.GetAttr("id") != "div1" {
		t.Fatal()
	}

	//(*Node) GetNodesByTagAttr(string) *Node
	if divs := root.GetNodesByTagAttr("div", "class", "1"); len(divs) != 2 {
		t.Fatal()
	} else if divs[0].GetAttr("id") != "div1" ||
		divs[1].GetAttr("id") != "div3" {
		t.Fatal()
	}

	if divs := root.GetNodesByTag("div"); len(divs) != 4 {
		t.Fatal()
	} else {
		//(Nodes) GetNodesByTag(string) *Node
		if spans := divs.GetNodesByTag("span"); len(spans) != 4 {
			t.Fatal()
		} else {
			if spans[0].GetAttr("id") != "span0" ||
				spans[1].GetAttr("id") != "span1" ||
				spans[2].GetAttr("id") != "span2" ||
				spans[3].GetAttr("id") != "span3" {
				t.Fatal()
			}
		}

		//(Nodes) GetNodesByTagAttr(string) *Node
		if spans := divs.GetNodesByTagAttr("span", "class", "0"); len(spans) != 2 {
			t.Fatal()
		} else {
			if spans[0].GetAttr("id") != "span0" ||
				spans[1].GetAttr("id") != "span2" {
				t.Fatal()
			}
		}
	}

	//(*Node) GetText() string
	if span := root.GetNodeByTagAttr("span", "id", "span0"); span == nil {
		t.Fatal()
	} else if span.GetText() != "text0" {
		t.Fatal()
	}
	if p := root.GetNodeByTag("p"); p == nil {
		t.Fatal()
	} else if p.GetText() != "line1\nline2" {
		t.Fatal()
	}
}
