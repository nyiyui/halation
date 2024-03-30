<script>
  import { inlineHelp } from '$lib/config.ts';
  import NodeSmall from '$lib/NodeSmall.svelte';
  import ColorPicker2 from '$lib/ColorPicker2.svelte';

  export let node;
  export let nodeName;

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
    if ('SG' in node.Node) {
      if (channels.length != node.Node.SG.State.channels.length) {
        channels = Array(node.Node.SG.State.channels.length);
      }
      node.Node.SG.State.channels.forEach((channel, i) => {
        channels[i] = {h:channel.hue, s:channel.saturation, v:channel.level};
      })
    }
  }

  function updateColour(i) {
    return (e) => {
      console.log(i, e.detail.hsv);
      // node.Node.SG.State.channels[i].hue = e.detail.hsv.h;
      // node.Node.SG.State.channels[i].saturation = e.detail.hsv.s;
      // node.Node.SG.State.channels[i].level = e.detail.hsv.v;
    }
  }

  let channelProps = {
    0: {name: "invalid"},
    1: {name: "left small"},
    2: {name: "left podium"},
    3: {name: "left front"},
    4: {name: "left side"},
    5: {name: "left flood"},
    6: {name: "left ground"},
    7: {name: "left centre wall"},
    8: {name: "left centre back"},
    10: {name: "right centre wall"},
    11: {name: "right centre back"},
    12: {name: "right centre ground"},
    13: {name: "right centre back2"},
    14: {name: "right wall"},
    15: {name: "right flood"},
    16: {name: "right ground"},
    17: {name: "right front"},
    24: {name: "lx3 left back", colour: false},
    31: {name: "lx4 blue", colour: false},
    33: {name: "lx4 white", colour: false},
    32: {name: "lx4 yellow", colour: false},
    34: {name: "lx4 red", colour: false},
    35: {name: "lx4 blue/green/white", colour: false},
    36: {name: "lx4 pink-purple", colour: false},
  };

  function handleChange(i) {
    return (e) => {
      let ch = parseInt(e.target.value);
      if (node.Node.SG.State.channels[i].channelID != ch) {
        node.Node.SG.State.channels[i].channelID = ch
      }
      console.log(ch);
    }
  }

  function deletePromise(i) {
    node.Node.Promises.splice(i, 1);
    node.Node.Promises = node.Node.Promises;
  }

  function deleteChannel(i) {
    node.Node.SG.State.channels.splice(i, 1);
    node.Node.SG.Statet.channels = node.Node.SG.Statet.channels;
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
  Name
  <input value={nodeName} readonly/>
</label>
<label>
  Description
  <input bind:value={node.Node.Description}/>
</label>
<fieldset>
  <legend>Promises</legend>
  {#if inlineHelp}
    <p>
      Nodes are activated when a promise is fulfilled.
      The field is filled by the return value of the node.
      If the field is <code>dummy</code>, no field is filled.
    </p>
  {/if}
  <button on:click={addPromise}>Add Promise</button>
  {#if node.Node.Promises}
    <table>
      <tr>
        <th>Supply</th>
        <th>Field</th>
        <th>Actions</th>
      </tr>
      {#each node.Node.Promises as promise, i}
        <tr>
          <td>
            <input bind:value={node.Node.Promises[i].SupplyNodeName} />
            <br />
            <NodeSmall nodeName={promise.SupplyNodeName} />
          </td>
          <td>
            <input bind:value={node.Node.Promises[i].FieldName} />
          </td>
          <td>
            <button on:click={() => deletePromise(i)}>Delete</button>
          </td>
        </tr>
      {/each}
    </table>
  {/if}
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
              <th colspan=2>Channel (L/R is audience)</th>
              <th>Colour (HSV) / Level (%)</th>
              <th>Actions</th>
            </tr>
            {#each node.Node.SG.State.channels as channel, i}
              <tr>
                <td><input type=number bind:value={channel.channelID} /></td>
                <td>
                  <select value={channel.channelID.toString()} on:change={handleChange(i)}>
                    {#each Object.entries(channelProps) as [ch, props]}
                      <option value={ch}>{ch} - {props.name}</option>
                    {/each}
                  </select>
                </td>
                {#if channelProps[channel.channelID].colour == false}
                  <td><label>Level <input type=number bind:value={channel.level} /></label></td>
                {:else}
                  <td><ColorPicker2 bind:hue={channel.hue} bind:saturation={channel.saturation} bind:value={channel.level} /></td>
                {/if}
                <td>
                  <button on:click={() => deleteChannel(i)}>Delete</button>
                </td>
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
        {#if inlineHelp}
          <p>
            The duration and resolution strings below are a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h". (See <a href="https://pkg.go.dev/time#ParseDuration">Go time.ParseDuration documentation</a> for details.)
          </p>
        {/if}
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
p, label {
  display: block;
}
</style>
