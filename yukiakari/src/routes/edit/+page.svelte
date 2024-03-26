<script lang="ts">
  import { onMount } from 'svelte';
  import { getNodes, patchNode } from '$lib/tsapi.ts';
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
{#if !!savePromise}
  {#await savePromise}
    saving
  {:then}
    saved
  {:catch err}
    {err}
  {/await}
{:else}
  not saved
{/if}
<hr />
<button on:click={saveNode}>Save</button>
<NodeEdit bind:node />
<button on:click={saveNode}>Save</button>
{/if}
