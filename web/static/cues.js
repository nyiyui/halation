let btn = document.getElementById("cues-name-column-toggle");
let hidden = false;

btn.addEventListener("click", () => {
  if (hidden) {
    document.getElementById("cues").classList.remove("hidden")
    btn.value = btn.dataset.hide;
  } else {
    document.getElementById("cues").classList.add("hidden")
    btn.value = btn.dataset.show;
  }
  hidden = !hidden;
});
btn.value = btn.dataset.hide;
btn.classList.remove("js-hidden")
