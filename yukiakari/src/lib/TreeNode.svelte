<script lang="ts">
  import { newNode, activateNode } from '$lib/tsapi.ts';

  export let nodes;
  export let nid: string;
  export let reversePromises: Record<string, Array<string>>;
  
  let node = nodes[nid].Node;
  let nodeType = nodes[nid].NodeType;
  console.log('nodes', nodes);

  const nodeTypeLetters = {
    "nyiyui.ca/halation/node.Manual": "Knob",
    "nyiyui.ca/halation/node.SetState": "State",
    "nyiyui.ca/halation/node.EvalLua": "Lua",
    "nyiyui.ca/halation/node.Timer": "Timer",
  };

  function activate() {
    console.log('nid', nid);
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

<div class="node">
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
</style>
