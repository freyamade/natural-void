interface Colours {
  primary: string;
  secondary: string;
}

function getSchemeColours(): Colours {
  // The scheme colours can be retrieved from the hero and the border colour of the title
  let primary, secondary: string;
  primary = getComputedStyle(document.querySelector('.hero.is-scheme-primary')!).backgroundColor;
  secondary = getComputedStyle(document.querySelector('.listen-title')!).borderBottomColor!;
  return {primary: primary, secondary: secondary};
}

async function loadEpisode(): Promise<void> {
  // Get the episode id
  const episodeID = (document.querySelector('#player-container')! as HTMLElement).dataset.episodeId!;
  // Get the colours of the scheme
  const colours: Colours = getSchemeColours();
}

document.addEventListener('DOMContentLoaded', loadEpisode, false)
