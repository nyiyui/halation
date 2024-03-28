<script lang="ts">
  import { getContext, onMount } from 'svelte';
  import { newNode, activateNode } from '$lib/tsapi.ts';

  export let nodes;
  export let nid: string;
  export let reversePromises: Record<string, Array<string>>;
  
  let node = nodes[nid].Node;
  let nodeType = nodes[nid].NodeType;

  const nodeTypeLetters = {
    "nyiyui.ca/halation/node.Manual": "Knob",
    "nyiyui.ca/halation/node.SetState": "State",
    "nyiyui.ca/halation/node.EvalLua": "Lua",
    "nyiyui.ca/halation/node.Timer": "Timer",
  };

  const ctx = getContext('listen');
  ctx.getSource().addEventListener('changed', (e) => {
    const data = JSON.parse(e.data);
    if (data.NodeName != nid) return;
    if (!data.Activated) return;
    bong();
  })

  let nodeElement;
  function bong() {
    nodeElement.style.borderLeftColor = 'aqua';
    setTimeout(() => {
      nodeElement.style.transition = 'border-left-color 1s';
      nodeElement.style.removeProperty('border-left-color');
    }, 3000);
    setTimeout(() => {
      nodeElement.style.removeProperty('transition');
    }, 3000+1000);
  }

  function activate() {
    let response = activateNode(nid);
  }

  function newDownstream() {
    const name = Math.random().toString(36).substr(3);
    // TODO: consider fixing this sus random token code
    let response = newNode(name, {
      NodeType: "nyiyui.ca/halation/node.Manual",
      Node: {
        Description: "Untitled",
        Promises: [
          {FieldName: "dummy", SupplyNodeName: nid},
        ],
      },
    });
    response.then(() => window.location.reload());
    // TODO: implement auto-refresh when a node changes or via svelte event handlers
  }
</script>

<div class="node" bind:this={nodeElement}>
  <div class="node-self">
    <div>
      <!-- <code>{nid}</code> -->
      <small>{nodeTypeLetters[nodeType]}</small>
      {node.Description}
    </div>
    {#if node.Promises}
      <div>
        Promises:
        {#each node.Promises as promise}
          <code>{promise.FieldName}</code>
        {/each}
      </div>
    {/if}
    <div>
      <input type="button" value="Activate" on:click={activate} />
      <input type="button" value="New Downstream" on:click={newDownstream} />
      <a href="/edit?node-name={encodeURIComponent(nid)}">Edit</a>
    </div>
    </div>
  <div>
    {#if reversePromises[nid]}
      {#if reversePromises[nid].length == 1}
        <svelte:self {nodes} nid={reversePromises[nid][0]} {reversePromises} />
      {:else}
      {#each reversePromises[nid] as downstream}
        <div style="display: block;">
          <svelte:self {nodes} nid={downstream} {reversePromises} />
        </div>
      {/each}
      {/if}
    {/if}
  </div>
</div>

<style>
.node {
  display: flex;
  flex-wrap: wrap;
  border-left: solid 2px grey;
  padding-left: 4px;
}
.node > div {
  margin: 4px 6px;
}
.node-self {
  display: flex;
  flex-direction: column;
}

@keyframes activated {
  50% {
    border-left-color: aqua;
    border-left-width: 4px;
    padding-left: 2px;
  }
}
</style>
