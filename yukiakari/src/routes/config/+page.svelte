<script lang="ts">
  import { inlineHelp, baseUrl, noEdit } from '$lib/config.ts';
</script>

<label>
  <input type=checkbox bind:checked={$inlineHelp} />
  Inline Help
  <span>(Show small explanation texts beside buttons etc.)</span>
</label>
<label>
  <input type=checkbox bind:checked={$inlineHelp} />
  Disable Editing
  <span>(maybe useful for during shows)</span>
</label>
<label>
  Halation Server API Base URL
  <input type=url bind:value={$baseUrl} />
</label>
<style>
  label {
    display: block;
  }
</style>
