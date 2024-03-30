<script>
  import { inlineHelp } from '$lib/config.ts';
  import { base } from '$app/paths';
</script>

<nav>
  <img src="{base}/favicon.png" height="16px" />
  Halation
	<a href="{base}/config">Config</a>
	<a href="{base}/tree">Tree</a>
</nav>

<div class="slot-container">
  <slot></slot>
</div>

<style>
  nav {
    display: flex;
    background: black;
    color: beige;
  }
  nav > * {
    padding: 4px;
  }
  nav a {
    color: aqua;
  }
  nav a:hover {
    text-decoration: underline;
  }
  nav a:active {
    color: aquamarine;
  }
  .slot-container {
    margin: 8px;
  }
  :global(body) {
    margin: 0;
  }
</style>
