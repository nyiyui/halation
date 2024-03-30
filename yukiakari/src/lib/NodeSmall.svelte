<script>
  import { getNodes } from '$lib/tsapi2.ts';

  export let nodeName;

  const nodeTypeLetters = {
    "nyiyui.ca/halation/node.Manual": "Knob",
    "nyiyui.ca/halation/node.SetState": "State",
    "nyiyui.ca/halation/node.EvalLua": "Lua",
    "nyiyui.ca/halation/node.Timer": "Timer",
  };

  let nodes = getNodes(nodeName);
  // TODO: getNode not implemented yet, use getNodes instead
</script>

{#await nodes}
loading
{:then nodes}
{#if nodes.Nodes[nodeName]}
<small>{nodeTypeLetters[nodes.Nodes[nodeName].NodeType]}</small>
{nodes.Nodes[nodeName].Node.Description}
{:else}
Nonexistent
{/if}
{/await}
