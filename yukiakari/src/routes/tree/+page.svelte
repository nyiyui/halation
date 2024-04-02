<script lang="ts">
  import { setContext, onMount } from 'svelte';
  import { getNodes, listenChanges, newNode } from '$lib/tsapi2.ts';
  import { noEdit } from '$lib/config.ts';
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
      for (const nodeName in data.Nodes) {
        const node = data.Nodes[nodeName].Node;
        if (node.Promises) {
          for (const promise of node.Promises) {
            if (!newReversePromises[promise.SupplyNodeName]) {
              newReversePromises[promise.SupplyNodeName] = [];
            }
            newReversePromises[promise.SupplyNodeName].push(nodeName);
          }
        }
        if (node.Promises && node.Promises.length != 0) continue;
        newRoots.push(nodeName);
      }
      roots.set(newRoots);
      reversePromises.set(newReversePromises);
    })
  })

  let source;
  onMount(() => {
    let reconnectTimeout = 1000;

    function tryConnect() {
        connect();
        reconnectTimeout *= 2;
        if (reconnectTimeout >= 30000) {
            reconnectTimeout = 30000;
        }
    }
    
    
    function connect() {
        source = listenChanges();
        source.addEventListener("connected", (e) => {
          console.log("listenChanges: received connected event");
        });
        source.onmessage = console.log;
        source.addEventListener('open', () => {
          reconnectTimeout = 1000;
          console.log("listenChanges: connected");
        });
        source.addEventListener('error', () => {
          source.close();
          setTimeout(tryConnect, reconnectTimeout)
          console.log(`listenChanges: reconnecting in ${reconnectTimeout} ms`);
        });
    }
    
    connect();
  })
  function getSource() {
    return source;
  }
  setContext('listen', { getSource })

  let nodes;
  onMount(() => { nodes = getNodes(); });
  let roots = writable([]);
  let reversePromises = writable({});
</script>

{#await nodes}
loading
{:then nodes}
  <div>
    <label>
      <input type="checkbox" bind:checked={$noEdit} />
      No Editing
    </label>
    <input type="button" value="New Node" on:click={handleNewNode} />
  </div>
  {#each $roots as nodeName}
    <div style="display: block;">
      <TreeNode nodes={nodes.Nodes} {nodeName} reversePromises={$reversePromises} />
    </div>
  {/each}
{/await}
