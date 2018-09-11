function toggleMenu() {
  document.querySelector('.navbar-burger')!.classList.toggle('is-active');
  document.querySelector('.navbar-menu')!.classList.toggle('is-active');
}

document.querySelector('.navbar-burger')!.addEventListener('click', toggleMenu, false);

function tabChange(e: Event) {
  const tab = e.target as HTMLElement;
  // Switch the active tab and story
  const active = document.querySelector('.hero-foot .tabs li.is-active a')! as HTMLElement;
  active.parentElement!.classList.remove('is-active');
  document.getElementById(active.dataset.target!)!.classList.add('is-removed');
  tab.parentElement!.classList.add('is-active');
  document.getElementById(tab.dataset.target!)!.classList.remove('is-removed');
}

const tabs = document.querySelectorAll('.hero-foot .tabs a');
for (let i = 0; i < tabs.length; i ++) {
  tabs[i].addEventListener('click', tabChange, false);
}
