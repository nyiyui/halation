<script lang="ts">
  import { inlineHelp, baseUrl } from '$lib/config.ts';
</script>

<label>
  <input type=checkbox bind:checked={$inlineHelp} />
  Inline Help
  <span>(Show small explanation texts beside buttons etc.)</span>
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
