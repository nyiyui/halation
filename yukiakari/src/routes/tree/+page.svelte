<script lang="ts">
  import { onMount } from 'svelte';
  import { getNodes } from '$lib/tsapi.ts';
  import ManualNode from '$lib/nodes/ManualNode.svelte';
  import SetStateNode from '$lib/nodes/SetStateNode.svelte';
  import EvalLuaNode from '$lib/nodes/EvalLuaNode.svelte';
  import TreeNode from '$lib/TreeNode.svelte';

  import { writable } from 'svelte/store';

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

  let nodes;
  onMount(() => { nodes = getNodes(); });
  let roots = writable([]);
  let reversePromises = writable({});
</script>

{#await nodes}
loading
{:then nodes}
  {$roots.length} roots
  {#each $roots as nid}
    <div style="display: block;">
      <TreeNode nodes={nodes.Nodes} {nid} reversePromises={$reversePromises} />
    </div>
  {/each}
{/await}
