<script>
  import NodeSmall from '$lib/NodeSmall.svelte';
  import ColorPicker2 from '$lib/ColorPicker2.svelte';

  export let node;

  function addPromise() {
    if (!node.Node.Promises) node.Node.Promises = [];
    node.Node.Promises.push({SupplyNodeName: "", FieldName: ""});
    node = node;
  }

  function addChannel() {
    node.Node.SG.State.channels.push({channelID:0, level:0, hue:0, saturation:0});
    node = node;
  }

  function updateNode() {
    if (node.NodeType == "nyiyui.ca/halation/node.Timer") {
      if (!("delay" in node.Node)) {
        node.Node.delay = "";
      }
    }
    if (node.NodeType == "nyiyui.ca/halation/node.SetState") {
      if (!node.Node.SG) {
        node.Node.SG = {
          StateType: "",
          State: {},
          GradientType: "",
          Gradient: {},
        };
        node = node;
      }
      if (node.Node.SG.StateType == "nyiyui.ca/halation/osc" && !("channels" in node.Node.SG.State)) {
        node.Node.SG.State = {
          blackout: true,
          channels: [],
        };
        node = node;
      }
      if (node.Node.SG.GradientType == "nyiyui.ca/halation/gradient.LinearGradient" && !("duration" in node.Node.SG.Gradient)) {
        node.Node.SG.Gradient = {
          duration: "",
          preferredResolution: "",
        };
        node = node;
      }
    }
  }

  let channels = [];
  $: {
    if (channels.length != node.Node.SG.State.channels.length) {
      channels = Array(node.Node.SG.State.channels.length);
    }
    node.Node.SG.State.channels.forEach((channel, i) => {
      channels[i] = {h:channel.hue, s:channel.saturation, v:channel.level};
    })
  }

  function updateColour(i) {
    return (e) => {
      console.log(i, e.detail.hsv);
      // node.Node.SG.State.channels[i].hue = e.detail.hsv.h;
      // node.Node.SG.State.channels[i].saturation = e.detail.hsv.s;
      // node.Node.SG.State.channels[i].level = e.detail.hsv.v;
    }
  }
</script>

<label>
  Type
  <select bind:value={node.NodeType} on:change={updateNode}>
    <option value="nyiyui.ca/halation/node.EvalLua">Lua</option>
    <option value="nyiyui.ca/halation/node.Manual">Knob</option>
    <option value="nyiyui.ca/halation/node.SetState">State</option>
    <option value="nyiyui.ca/halation/node.Timer">Timer</option>
  </select>
</label>
<label>
  Description
  <input bind:value={node.Node.Description}/>
</label>
<fieldset>
  <legend>Promises</legend>
  <button on:click={addPromise}>Add Promise</button>
  <ul>
    {#if node.Node.Promises}
      {#each node.Node.Promises as promise, i}
        <li>
          <label>
            Supply
            <input bind:value={node.Node.Promises[i].SupplyNodeName} />
          </label>
          <label>
            Field
            <input bind:value={node.Node.Promises[i].FieldName} />
          </label>
          <NodeSmall nodeName={promise.SupplyNodeName} />
        </li>
      {/each}
    {/if}
  </ul>
  <button on:click={addPromise}>Add Promise</button>
</fieldset>
<fieldset>
  <legend>Type-Specific</legend>
  {#if node.NodeType == "nyiyui.ca/halation/node.Manual"}
  {:else if node.NodeType == "nyiyui.ca/halation/node.Timer" && "delay" in node.Node}
    <label>
      Duration
      <input bind:value={node.Node.delay} />
    </label>
  {:else if node.NodeType == "nyiyui.ca/halation/node.SetState" && node.Node.SG}
    <fieldset>
      <legend>State</legend>
      <label>
        State Type
        <select bind:value={node.Node.SG.StateType} on:change={updateNode}>
          <option value="nyiyui.ca/halation/osc">Lightboard</option>
          <option value="nyiyui.ca/halation/mpv">MPV</option>
        </select>
      </label>
      {#if node.Node.SG.StateType == "nyiyui.ca/halation/osc"}
        <label>
          <input type=checkbox bind:checked={node.Node.SG.State.blackout} />
          Blackout
        </label>
        {#if !node.Node.SG.State.blackout}
          <button on:click={addChannel}>Add Row</button>
          <table>
            <tr>
              <th>Channel</th>
              <th>Colour</th>
            </tr>
            {#each node.Node.SG.State.channels as channel, i}
              <tr>
                <td><input type=number bind:value={channel.channelID} /></td>
                <td><ColorPicker2 bind:hue={channel.hue} bind:saturation={channel.saturation} bind:value={channel.level} /></td>
              </tr>
            {/each}
          </table>
          <button on:click={addChannel}>Add Row</button>
        {/if}
      {/if}
    </fieldset>
    <fieldset>
      <legend>Gradient</legend>
      <label>
        Gradient Type
        <select bind:value={node.Node.SG.GradientType} on:change={updateNode}>
          <option value="nyiyui.ca/halation/gradient.LinearGradient">Linear</option>
        </select>
      </label>
      {#if node.Node.SG.GradientType == "nyiyui.ca/halation/gradient.LinearGradient" && "duration" in node.Node.SG.Gradient}
        <label>
          Duration
          <input bind:value={node.Node.SG.Gradient.duration} />
        </label>
        <label>
          Preferred Resolution
          <input bind:value={node.Node.SG.Gradient.preferredResolution} />
        </label>
      {/if}
    </fieldset>
  {/if}
</fieldset>
<details>
  <summary>JSON</summary>
  <code>{JSON.stringify(node)}</code>
</details>

<style>
label {
  display: block;
}
</style>
