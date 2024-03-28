<script>
  import ColorPicker from 'svelte-awesome-color-picker';

  export let hue;
  export let saturation;
  export let value;

  let hsv = {};
  $: {
    if (hue != hsv.h || saturation != hsv.s || value != hsv.v) {
      hsv = {h:hue, s:saturation, v:value};
    }
  }
  function handleChange(e) {
    hue = e.detail.hsv.h;
    saturation = e.detail.hsv.s;
    value = e.detail.hsv.v;
    console.log('change', hue)
  }
</script>

<ColorPicker {hsv} on:input={handleChange}/>
