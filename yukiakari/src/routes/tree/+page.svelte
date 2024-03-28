<script lang="ts">
  import { setContext, onMount } from 'svelte';
  import { getNodes, listenChanges, newNode } from '$lib/tsapi.ts';
  import ManualNode from '$lib/nodes/ManualNode.svelte';
  import SetStateNode from '$lib/nodes/SetStateNode.svelte';
  import EvalLuaNode from '$lib/nodes/EvalLuaNode.svelte';
  import TreeNode from '$lib/TreeNode.svelte';

  import { writable } from 'svelte/store';

  function handleNewNode() {
    const name = Math.random().toString(36).substr(3);
    // TODO: consider fixing this sus random token code
    let response = newNode(name, {
      NodeType: "nyiyui.ca/halation/node.Manual",
      Node: {
        Description: "Untitled",
      },
    });
    response.then(() => window.location.reload());
    // TODO: implement auto-refresh when a node changes or via svelte event handlers
  }

  onMount(function() {
    getNodes().then((data) => {
      let newRoots = [];
      let newReversePromises = {};
      for (const nid in data.Nodes) {
        const node = data.Nodes[nid].Node;
        if (node.Promises) {
          for (const promise of node.Promises) {
            if (!newReversePromises[promise.SupplyNodeName]) {
              newReversePromises[promise.SupplyNodeName] = [];
            }
            newReversePromises[promise.SupplyNodeName].push(nid);
          }
        }
        if (node.Promises && node.Promises.length != 0) continue;
        newRoots.push(nid);
      }
      roots.set(newRoots);
      reversePromises.set(newReversePromises);
    })
  })

  let source;
  onMount(() => {
    console.log('listen');
    source = listenChanges();
    source.onmessage = console.log;
    console.log('listen2', source);
    source.addEventListener("message", (e) => {
      console.log("listenChanges: "+e);
    });
    source.addEventListener("connected", (e) => {
      console.log("listenChanges: connected");
    });
    source.addEventListener("changed", handleChange)
  })
  function getSource() {
    return source;
  }
  setContext('listen', { getSource })

  function handleChange(e) {
    console.log(e)
    let data = JSON.parse(e.data);
    console.log(data.NodeName, data.Activated);
  }

  let nodes;
  onMount(() => { nodes = getNodes(); });
  let roots = writable([]);
  let reversePromises = writable({});
</script>

{#await nodes}
loading
{:then nodes}
  {$roots.length} roots
  <input type="button" value="New Node" on:click={handleNewNode} />
  {#each $roots as nid}
    <div style="display: block;">
      <TreeNode nodes={nodes.Nodes} {nid} reversePromises={$reversePromises} />
    </div>
  {/each}
{/await}
