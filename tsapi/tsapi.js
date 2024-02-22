var __assign = (this && this.__assign) || function () {
    __assign = Object.assign || function(t) {
        for (var s, i = 1, n = arguments.length; i < n; i++) {
            s = arguments[i];
            for (var p in s) if (Object.prototype.hasOwnProperty.call(s, p))
                t[p] = s[p];
        }
        return t;
    };
    return __assign.apply(this, arguments);
};
var ManualTypeName = "nyiyui.ca/halation/node.Manual";
var SetStateTypeName = "nyiyui.ca/halation/node.SetState";
var TimerTypeName = "nyiyui.ca/halation/node.Timer";
var EvalLuaTypeName = "nyiyui.ca/halation/node.EvalLua";
var MPVStateTypeName = "nyiyui.ca/halation/mpv";
var OSCStateTypeName = "nyiyui.ca/halation/osc";
var LinearGradientTypeName = "nyiyui.ca/halation/gradient.LinearGradient";
var baseUrl = "http://localhost:8080";
function doRequest(method, path, params) {
    return fetch(new URL(path, baseUrl).toString(), __assign({ method: method }, params));
}
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
