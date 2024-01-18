function showReloadButtonOnChange(reloadButton) {
  const es = new EventSource("/events/change");
  es.addEventListener("changed", () => {
    reloadButton.style.display = "block";
  });
}

export { showReloadButtonOnChange };
