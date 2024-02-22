<script lang="ts">
  import { onMount } from 'svelte';
  import { getNodes } from '$lib/tsapi.ts';
  import ManualNode from '$lib/nodes/ManualNode.svelte';
  import SetStateNode from '$lib/nodes/SetStateNode.svelte';
  import EvalLuaNode from '$lib/nodes/EvalLuaNode.svelte';
  import { writable } from 'svelte/store';
  import {
    SvelteFlow,
    Controls,
    Background,
    BackgroundVariant,
    MiniMap
  } from '@xyflow/svelte';
 
  // 👇 this is important! You need to import the styles for Svelte Flow to work
  import '@xyflow/svelte/dist/style.css';
 
  // We are using writables for the nodes and edges to sync them easily. When a user drags a node for example, Svelte Flow updates its position.
  const nodes = writable([
    {
      id: '1',
      type: 'input',
      data: { label: 'Input Node' },
      position: { x: 0, y: 0 }
    },
    {
      id: '2',
      type: 'default',
      data: { label: 'Node' },
      position: { x: 0, y: 150 }
    }
  ]);
 
  // same for edges
  const edges = writable([
    {
      id: '1-2',
      type: 'default',
      source: '1',
      target: '2',
      label: 'Edge Text'
    }
  ]);
 
  const snapGrid = [25, 25];

  const nodeTypes = {
    'nyiyui.ca/halation/node.Manual': ManualNode,
    'nyiyui.ca/halation/node.SetState': SetStateNode,
    'nyiyui.ca/halation/node.EvalLua': EvalLuaNode,
  };

  onMount(function() {
    getNodes().then((data) => {
      let newNodes = [];
      for (const nodeName in data.Nodes) {
        const entry = data.Nodes[nodeName]
        newNodes.push({
          id: nodeName,
          type: entry.NodeType,
          data: entry.Node,
          position: { x: 0, y: 0 },
        })
      }
      nodes.set(newNodes)
      console.log(newNodes)

      let newEdges = [];
      for (const nodeName in data.Nodes) {
        const entry = data.Nodes[nodeName]
        if (entry.Node.Promises !== null) {
          for (const promise of entry.Node.Promises) {
          console.log(promise)
            newEdges.push({
              id: `${nodeName}-${promise.SupplyNodeName}`,
              source: promise.SupplyNodeName,
              target: nodeName,
              targetHandle: promise.FieldName,
            })
          }
        }
      }
      edges.set(newEdges)
      console.log(newEdges)
    })
  })
</script>
 
<!--
👇 By default, the Svelte Flow container has a height of 100%.
This means that the parent container needs a height to render the flow.
-->
<div style:height="500px">
  <SvelteFlow
    {nodes}
    {edges}
    {snapGrid}
    {nodeTypes}
    fitView
    on:nodeclick={(event) => console.log('on node click', event.detail.node)}
  >
    <Controls />
    <Background variant={BackgroundVariant.Dots} />
    <MiniMap />
  </SvelteFlow>
</div>
