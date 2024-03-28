
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
// Node Types
type Timer = {Description: string,Promises: NodePromise[] | null,Delay: string,};
const TimerTypeName = "nyiyui.ca/halation/node.Timer";
type EvalLua = {Description: string,Promises: NodePromise[] | null,Source: string,};
const EvalLuaTypeName = "nyiyui.ca/halation/node.EvalLua";
type Manual = {Description: string,Promises: NodePromise[] | null,};
const ManualTypeName = "nyiyui.ca/halation/node.Manual";
type SetState = {Description: string,Promises: NodePromise[] | null,SG: SG,};
const SetStateTypeName = "nyiyui.ca/halation/node.SetState";
// State Types
type State = MPVState|OSCState;
type MPVState = {FilePath: string,Paused: boolean,Position: number,Fullscreen: boolean,ExtraProperties: Record<string, number | boolean | string>,};
const MPVStateTypeName = "nyiyui.ca/halation/mpv";
type OSCState = {Blackout: boolean,Channels: Channel[],};
const OSCStateTypeName = "nyiyui.ca/halation/osc";
// Gradient Types
type LinearGradient = {Duration_: string,PreferredResolution_: string,};
const LinearGradientTypeName = "nyiyui.ca/halation/gradient.LinearGradient";
type Gradient = LinearGradient;

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

export { listenChanges, ensureNode, newNode, activateNode, patchNode, deleteNode, getNodes, getNode };
