import Amplitude from 'amplitudejs'

interface Colours {
  primary: string
  secondary: string
}

function getSchemeColours(): Colours {
  // The scheme colours can be retrieved from the hero and the border colour of the title
  let primary, secondary: string
  primary = getComputedStyle(document.querySelector('.hero.is-scheme-primary')!).backgroundColor
  secondary = getComputedStyle(document.querySelector('.listen-title')!).borderBottomColor!
  return {primary: primary, secondary: secondary}
}

async function loadEpisode(): Promise<void> {
  // Get the episode id
  const dataset = (document.querySelector('#player-container')! as HTMLElement).dataset
  const episodeNum = dataset.episodeNum!
  const storyId = dataset.storyId
  // Get the colours of the scheme
  const colours: Colours = getSchemeColours()
  Amplitude.init({
    bindings: {
      32: 'play_pause',
    },
    songs: [{
      url: `/episodes/${storyId}/${episodeNum}`,
    }]
  })

  // Allow for clicking to move around on the progress bar
  document.getElementById('song-played-progress')!.addEventListener('click', (e: MouseEvent) => {
    const elmnt = e.target as HTMLElement
    const offset = elmnt.getBoundingClientRect()
    const x = e.pageX - offset.left
    Amplitude.setSongPlayedPercentage((x / elmnt.offsetWidth) * 100)
  })
}

loadEpisode()
