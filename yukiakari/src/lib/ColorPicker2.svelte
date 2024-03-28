<script>
  import ColorPicker from 'svelte-awesome-color-picker';

  export let hue;
  export let saturation;
  export let value;

  let hsv = {hue, saturation, value};
  $: {
    hue = hsv.hue;
    console.log('upd');
  }
</script>

<ColorPicker bind:hsv />
