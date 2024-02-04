function morphIntoDraggable(elem, handle) {
  let pos1 = 0, pos2 = 0, pos3 = 0, pos4 = 0;
  handle.addEventListener("mousedown", dragMouseDown);

  function dragMouseDown(e) {
    if (e.button != 0) return;
    e.preventDefault();
    // get the mouse cursor position at startup:
    pos3 = e.clientX;
    pos4 = e.clientY;
    document.onmouseup = closeDragElement;
    // call a function whenever the cursor moves:
    document.onmousemove = elementDrag;
  }

  function elementDrag(e) {
    e.preventDefault();
    // calculate the new cursor position:
    pos1 = pos3 - e.clientX;
    pos2 = pos4 - e.clientY;
    pos3 = e.clientX;
    pos4 = e.clientY;
    // set the element's new position:
    elem.style.top = (elem.offsetTop - pos2) + "px";
    elem.style.left = (elem.offsetLeft - pos1) + "px";
  }

  function closeDragElement() {
    // stop moving when mouse button is released:
    document.onmouseup = null;
    document.onmousemove = null;
  }
}

let latest;
// https://stackoverflow.com/a/70691308
// Note: persist popup (even after clicking outside) to allow having multiple popups open at the same time.
function morphIntoPopup(elem) {
  const absX = latest.clientX + window.scrollX;
  const absY = latest.clientY + window.scrollY;
  
  const bcrParent = elem.parentElement.getBoundingClientRect();
  const bcrPopup = elem.getBoundingClientRect();
  
  const maxX = bcrParent.width - bcrPopup.width;
  const maxY = bcrParent.height - bcrPopup.height;
  
  //const x = Math.max(0, Math.min(absX, maxX));
  //const y = Math.max(0, Math.min(absY, maxY));
  const x = absX;
  const y = absY;
  
  elem.classList.remove("popup-backburner");
  Object.assign(elem.style, {
    left: `${x}px`,
    top: `${y}px`,
  });
};
function mousemoveListener(e) {
  latest = e;
}
function popupListener(e) {
  if (!e.target.classList.contains("node-edit")) return;
  let modal = document.getElementById(e.target.dataset.target);
  morphIntoPopup(modal);
}
addEventListener("mousemove", mousemoveListener);
addEventListener("htmx:afterRequest", popupListener);
addEventListener("htmx:afterRequest", (e) => { // cleanup
  if (e.target.id != "content") return;
  removeEventListener("mousemove", mousemoveListener);
  removeEventListener("htmx:afterRequest", popupListener);
}, { once: true });

export { morphIntoDraggable };
