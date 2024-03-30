<script lang="ts">
  import { onMount } from 'svelte';
  import { getNodes, listenChanges, activateNode, patchNode, ensureNode } from '$lib/tsapi2.ts';
  import ManualNode from '$lib/nodes/ManualNode.svelte';
  import SetStateNode from '$lib/nodes/SetStateNode.svelte';
  import EvalLuaNode from '$lib/nodes/EvalLuaNode.svelte';
  import NodeEdit from '$lib/NodeEdit.svelte';

  import { writable } from 'svelte/store';

  let nodeName: string;
  onMount(() => {
    const urlParams = new URLSearchParams(window.location.search);
    nodeName = urlParams.get("node-name");
  });

  let livePromise;
  let live = false;
  $: {
    let _ = node;
    renderLive();
  }
  function renderLive() {
    if (!live) return;
    console.log('pass');
    let node2 = structuredClone(node);
    node2.Node.Promises = null;
    livePromise = ensureNode("__live", node2).then(async () => await activateNode("__live"));
  }

  let node;
  let promise;
  let loading = true;
  onMount(() => {
    promise = getNodes().then((data) => {
      if (!(nodeName in data.Nodes)) {
        throw new Error("not found");
      }
      node = data.Nodes[nodeName]
      loading = false;
    });
  });

  let savePromise;
  function saveNode() {
    savePromise = patchNode(nodeName, node);
  }

  $: { node; savePromise = null; }

  onMount(() => {
    let source = listenChanges();
    source.addEventListener("changed", (e) => {
      const data = JSON.parse(e.data);
      if (data.NodeName != nodeName) return;
      if (!data.Activated) return;
      bong();
    })
  })
  let activated = false;
  function bong() {
    activated = true;
    setTimeout(() => {
      activated = false;
    }, 3000);
  }

  function activate() {
    let response = activateNode(nodeName);
  }
</script>


{#if loading}
{#if promise}
{#await promise}
loading
{:catch err}
{err}
{/await}
{/if}
{:else}
<button on:click={activate}>Activate</button>
<button on:click={saveNode}>Save</button>
<label>
  <input type=checkbox bind:checked={live} />
  Live
</label>
{#if live}
  {#if !!livePromise}
    {#await livePromise}
      (wait…)
    {:then}
      (live)
    {:catch err})
      ({err})
    {/await}
  {:else}
    (not activated)
  {/if}
{/if}
{#if !!savePromise}
  {#await savePromise}
    (saving)
  {:then}
    (saved)
  {:catch err})
    ({err})
  {/await}
{:else}
  (not saved)
{/if}
{#if activated}
(activated)
{/if}
<NodeEdit bind:node {nodeName} />
{/if}
