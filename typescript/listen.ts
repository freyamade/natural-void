import * as WavePlayer from 'waveplayer';

interface Colours {
  primary: string;
  secondary: string;
}

interface TrackMetaData {
  left: number[];
  duration: number;
}

interface TimeData {
  hours: number;
  minutes: number;
  seconds: number;
}

function getSchemeColours(): Colours {
  // The scheme colours can be retrieved from the hero and the border colour of the title
  let primary, secondary: string;
  primary = getComputedStyle(document.querySelector('.hero.is-scheme-primary')!).backgroundColor;
  secondary = getComputedStyle(document.querySelector('.listen-title')!).borderBottomColor!;
  return {primary: primary, secondary: secondary};
}

function convertTime(time: number): TimeData {
  time = Math.ceil(time);
  const hours = Math.floor(time / 60 / 60);
  time = Math.ceil(time % (60 * 60));
  const minutes = Math.floor(time / 60);
  time = Math.ceil(time % 60);
  return {hours, minutes, seconds: time};
}

async function loadEpisode(): Promise<void> {
  // Get the episode id
  const episodeID = (document.querySelector('#player-container')! as HTMLElement).dataset.episodeId!;
  // Get the colours of the scheme
  const colours: Colours = getSchemeColours();
  const player: WavePlayer = new WavePlayer({
    container: '#player-container',
    barWidth: 4,
    barGap: 1,
    height: 128,
    progressColor: "#600",
    waveColour: 'white',
    responsive: true,
  });
  // Load the episode's metadata so that we can also render the time
  let data: TrackMetaData;
  try {
    data = await (await fetch(`/episodes/${episodeID}/episode.json`)).json();
    const duration = convertTime(data.duration);
    // Set the player to load based on the id of the episode
    await player.load(`/episodes/${episodeID}/episode.wav`, {data: data.left});
    console.log('Player should have loaded');
    // player.play();
  }
  catch (e) {
    console.error(e);
  }
}

document.addEventListener('DOMContentLoaded', loadEpisode, false)
