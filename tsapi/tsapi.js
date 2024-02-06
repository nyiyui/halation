var EvalLuaTypeName = "nyiyui.ca/halation/node.EvalLua";
var ManualTypeName = "nyiyui.ca/halation/node.Manual";
var SetStateTypeName = "nyiyui.ca/halation/node.SetState";
var TimerTypeName = "nyiyui.ca/halation/node.Timer";
var MPVStateTypeName = "nyiyui.ca/halation/mpv";
var OSCStateTypeName = "nyiyui.ca/halation/osc";
var LinearGradientTypeName = "nyiyui.ca/halation/gradient.LinearGradient";
function newNode(node) {
    return doRequest("POST", "/node/new", { body: JSON.stringify(node) });
}
function activateNode(name) {
    return doRequest("POST", "/node/" + encodeURIComponent(name.Package + "." + name.Name) + "/activate", {});
}
function patchNode(name, node) {
    return doRequest("PATCH", "/node/" + encodeURIComponent(name.Package + "." + name.Name), { body: JSON.stringify(node) });
}
function deleteNode(name, node) {
    return doRequest("DELETE", "/node/" + encodeURIComponent(name.Package + "." + name.Name), { body: JSON.stringify(node) });
}
function getNodes() {
    return doRequest("GET", "/nodes", {});
}
function getNode(name) {
    return doRequest("GET", "/node/" + encodeURIComponent(name.Package + "." + name.Name), {});
}
