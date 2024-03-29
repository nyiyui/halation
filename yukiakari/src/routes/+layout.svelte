<nav>
  <img src="/favicon.png" height="16px" />
  Halation
	<a href="/tree">Tree</a>
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
