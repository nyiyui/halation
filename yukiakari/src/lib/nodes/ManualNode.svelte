<script lang="ts">
  import { Handle, Position, type NodeProps } from '@xyflow/svelte';
  import type { Writable } from 'svelte/store';
 
  type $$Props = NodeProps;

  export let id: $$Props['id']; id;
  export let dragHandle: $$Props['dragHandle'] = undefined; dragHandle;
  export let type: $$Props['type']  = undefined; type;
  export let selected: $$Props['selected'] = undefined; selected;
  export let isConnectable: $$Props['isConnectable'] = undefined; isConnectable;
  export let zIndex: $$Props['zIndex'] = undefined; zIndex;
  export let width: $$Props['width'] = undefined; width;
  export let height: $$Props['height'] = undefined; height;
  export let dragging: $$Props['dragging']; dragging;
  export let targetPosition: $$Props['targetPosition'] = undefined; targetPosition;
  export let sourcePosition: $$Props['sourcePosition'] = undefined; sourcePosition;
  export let positionAbsoluteX: unknown;
  export let positionAbsoluteY: unknown;
 
  export let data: { color: Writable<string> };
 
  const { Description } = data;
</script>
 
<div class="colorpicker">
  <span>{Description}</span>
  (Manual)
  <Handle type="target" id="dummy" position={Position.Top} />
  <Handle type="source" position={Position.Bottom} />
</div>
