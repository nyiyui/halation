package main

//go:generate go run ./tsapi.go

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"nyiyui.ca/halation/aiz"
	_ "nyiyui.ca/halation/gradient"
	_ "nyiyui.ca/halation/mpv"
	"nyiyui.ca/halation/node"
	"nyiyui.ca/halation/osc"
	_ "nyiyui.ca/halation/osc"
	"nyiyui.ca/halation/timeutil"
)

func main() {
	f, err := os.Create("./tsapi.ts")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fmt.Fprint(f, `
type SG = { State: State, StateType: string, Gradient: Gradient, GradientType: string }
type Channel = {
  ChannelID: number,
	Level: number,
	Hue: number,
	Saturation: number,
}
type NodePromise = {
  FieldName: string,
	SupplyNodeName: string,
}
type NodeName = {
  Package: string,
	Name: string,
}
type NodeInAPI = {
  NodeType: string,
	Node: Node,
}
`)
	generateNodes(f)
	generateSG(f)

	fmt.Fprint(f, `
let baseUrl = "http://localhost:8080/api/v1/";

function doRequest(method, path, params) {
  return fetch((new URL(path, baseUrl)).toString(), {
	  method,
		...params,
  }).then(async response => {
    if (!response.ok) {
			throw new Error("HTTP error " + response.status + ": " + await response.text());
    }
    return response.json();
  })
}

function listenChanges(name: string, node: NodeInAPI) {
	return new EventSource((new URL("nodes/events", baseUrl)).toString());
}
function ensureNode(name: string, node: NodeInAPI) {
	return doRequest("POST", "node/" + encodeURIComponent(name) + "?override=yes", { body: JSON.stringify(node) })
}
function newNode(name: string, node: NodeInAPI) {
	return doRequest("POST", "node/" + encodeURIComponent(name), { body: JSON.stringify(node) })
}
function activateNode(name: string) {
	return doRequest("POST", "node/" + encodeURIComponent(name) + "/activate", {})
}
function patchNode(name: string, node: NodeInAPI) {
	return doRequest("PATCH", "node/" + encodeURIComponent(name), { body: JSON.stringify(node) })
}
function deleteNode(name: string, node: NodeInAPI) {
	return doRequest("DELETE", "node/" + encodeURIComponent(name), { body: JSON.stringify(node) })
}
type GetNodesResponse = {
  Nodes: Record<string, NodeInAPI>,
};
function getNodes(): GetNodesResponse {
  return doRequest("GET", "nodes", {})
}
function getNode(name: NodeName) {
	return doRequest("GET", "node/" + encodeURIComponent(name.Package + "." + name.Name), {})
}
`)
	fmt.Fprint(f, `
export { listenChanges, ensureNode, newNode, activateNode, patchNode, deleteNode, getNodes, getNode };
`)
}

func generateNodes(w io.Writer) {
	fmt.Fprint(w, "// Node Types\n")
	for tn, newNode := range node.NodeTypes {
		n := newNode()
		t := reflect.TypeOf(n).Elem()
		inner := "Description: string,"
		inner += "Promises: NodePromise[] | null,"
		for _, sf := range reflect.VisibleFields(t) {
			if sf.Name == "BaseNode" || sf.Name == "Description" || sf.Name == "Promises" {
				continue
			}
			inner += fmt.Sprintf("%s: %s,", sf.Name, typeName(sf.Type))
		}
		fmt.Fprintf(w, "type %s = {%s};\n", t.Name(), inner)
		fmt.Fprintf(w, "const %sTypeName = %s;\n", t.Name(), strconv.Quote(tn))
	}
}

func generateSG(w io.Writer) {
	tsNames := map[string]string{
		"nyiyui.ca/halation/mpv": "MPVState",
		"nyiyui.ca/halation/osc": "OSCState",
	}
	fmt.Fprint(w, "// State Types\n")
	names := make([]string, 0, len(tsNames))
	for _, val := range tsNames {
		names = append(names, val)
	}
	fmt.Fprintf(w, "type State = %s;\n", strings.Join(names, "|"))
	for tn, newState := range aiz.StateTypes {
		n := newState()
		t := reflect.TypeOf(n).Elem()
		inner := ""
		for _, sf := range reflect.VisibleFields(t) {
			if tn == "nyiyui.ca/halation/mpv" && sf.Name == "ExtraProperties" {
				inner += "ExtraProperties: Record<string, number | boolean | string>,"
				continue
			}
			inner += fmt.Sprintf("%s: %s,", sf.Name, typeName(sf.Type))
		}
		fmt.Fprintf(w, "type %s = {%s};\n", tsNames[tn], inner)
		fmt.Fprintf(w, "const %sTypeName = %s;\n", tsNames[tn], strconv.Quote(tn))
	}

	fmt.Fprint(w, "// Gradient Types\n")
	names = make([]string, 0, len(tsNames))
	for tn, newGradient := range aiz.GradientTypes {
		n := newGradient()
		t := reflect.TypeOf(n).Elem()
		inner := ""
		for _, sf := range reflect.VisibleFields(t) {
			inner += fmt.Sprintf("%s: %s,", sf.Name, typeName(sf.Type))
		}
		fmt.Fprintf(w, "type %s = {%s};\n", t.Name(), inner)
		fmt.Fprintf(w, "const %sTypeName = %s;\n", t.Name(), strconv.Quote(tn))
		names = append(names, t.Name())
	}
	fmt.Fprintf(w, "type Gradient = %s;\n", strings.Join(names, "|"))
}

func typeName(t reflect.Type) string {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() == reflect.Interface {
		return "unknown"
	}
	if t.Kind() == reflect.Slice {
		return fmt.Sprintf("%s[]", typeName(t.Elem()))
	}
	if t == reflect.TypeOf(*new(time.Duration)) {
		return "string"
	}
	if t == reflect.TypeOf(*new(timeutil.Duration)) {
		return "string"
	}
	if t == reflect.TypeOf(*new(aiz.SG)) {
		return "SG"
	}
	if t == reflect.TypeOf(*new(osc.Channel)) {
		return "Channel"
	}
	if t.Kind() == reflect.Bool {
		return "boolean"
	}
	if t.Kind() == reflect.Int {
		return "number"
	}
	if t.Kind() == reflect.Map {
		return fmt.Sprintf("Record<%s, %s>", typeName(t.Key()), typeName(t.Elem()))
	}
	return fmt.Sprint(t)
}
