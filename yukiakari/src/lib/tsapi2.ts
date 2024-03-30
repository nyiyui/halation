import { get } from 'svelte/store';
import * as api from './tsapi.ts';
import type { NodeInAPI, GetNodesResponse } from './tsapi.ts';
import { baseUrl } from './config.ts';

function listenChanges(name: string, node: NodeInAPI) {
  return api.listenChanges(name, node, get(baseUrl))
}
function ensureNode(name: string, node: NodeInAPI) {
  return api.ensureNode(name, node, get(baseUrl))
}
function newNode(name: string, node: NodeInAPI) {
  return api.newNode(name, node, get(baseUrl))
}
function activateNode(name: string) {
  return api.activateNode(name, get(baseUrl))
}
function patchNode(name: string, node: NodeInAPI) {
  return api.patchNode(name, node, get(baseUrl))
}
function deleteNode(name: string, node: NodeInAPI) {
  return api.deleteNode(name, get(baseUrl))
}
function getNodes(): GetNodesResponse {
  return api.getNodes(get(baseUrl))
}
function getNode(name: NodeName, baseUrl: string = baseUrl) {
  return api.getNode(name, get(baseUrl))
}

export { listenChanges, ensureNode, newNode, activateNode, patchNode, deleteNode, getNodes, getNode };
