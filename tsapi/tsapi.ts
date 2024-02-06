type SG = {
  State: State;
  StateType: string;
  Gradient: Gradient;
  GradientType: string;
};
type Channel = {
  ChannelID: number;
  Level: number;
  Hue: number;
  Saturation: number;
};
type NodePromise = {
  FieldName: string;
  SupplyNodeName: NodeName;
};
type NodeName = {
  Package: string;
  Name: string;
};
type NodeInAPI = {
  NodeType: string;
  Node: Node;
};
// Node Types
type EvalLua = { Description: string; Promises: NodePromise[]; Source: string };
const EvalLuaTypeName = "nyiyui.ca/halation/node.EvalLua";
type Manual = { Description: string; Promises: NodePromise[] };
const ManualTypeName = "nyiyui.ca/halation/node.Manual";
type SetState = { Description: string; Promises: NodePromise[]; SG: SG };
const SetStateTypeName = "nyiyui.ca/halation/node.SetState";
type Timer = { Description: string; Promises: NodePromise[]; Delay: string };
const TimerTypeName = "nyiyui.ca/halation/node.Timer";
// State Types
type State = MPVState | OSCState;
type MPVState = {
  FilePath: string;
  Paused: boolean;
  Position: number;
  Fullscreen: boolean;
  ExtraProperties: Record<string, number | boolean | string>;
};
const MPVStateTypeName = "nyiyui.ca/halation/mpv";
type OSCState = { Blackout: boolean; Channels: Channel[] };
const OSCStateTypeName = "nyiyui.ca/halation/osc";
// Gradient Types
type LinearGradient = { Duration_: string; PreferredResolution_: string };
const LinearGradientTypeName = "nyiyui.ca/halation/gradient.LinearGradient";
type Gradient = LinearGradient;

function newNode(node: NodeInAPI) {
  return doRequest("POST", "/node/new", { body: JSON.stringify(node) });
}
function activateNode(name: NodeName) {
  return doRequest(
    "POST",
    "/node/" + encodeURIComponent(name.Package + "." + name.Name) + "/activate",
    {},
  );
}
function patchNode(name: NodeName, node: NodeInAPI) {
  return doRequest(
    "PATCH",
    "/node/" + encodeURIComponent(name.Package + "." + name.Name),
    { body: JSON.stringify(node) },
  );
}
function deleteNode(name: NodeName, node: NodeInAPI) {
  return doRequest(
    "DELETE",
    "/node/" + encodeURIComponent(name.Package + "." + name.Name),
    { body: JSON.stringify(node) },
  );
}
function getNodes() {
  return doRequest("GET", "/nodes", {});
}
function getNode(name: NodeName) {
  return doRequest(
    "GET",
    "/node/" + encodeURIComponent(name.Package + "." + name.Name),
    {},
  );
}
