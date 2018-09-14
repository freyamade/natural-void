# natural-void
A podcast site made specifically for the ExVo D&amp;D games.

That being said, I'll be making it as easy as possible for others to use so feel free :)

Made with Go, Bulma, TypeScript and &lt;3

# Wavelength for player
The wavelength player tool is provided by https://www.npmjs.com/package/waveplayer. This will also be styled using the main colours of the (chosen) theme.

## wav2json
[wav2json](https://github.com/beschulz/wav2json) needs to be installed on the machine you are running the server from (it's installed in the Dockerfile anyway).

This is because waveplayer requires JSON metadata about the file to draw the chart without doing analysis.

After installing the script, the following command should be run to generate the json metadata;

`wav2json episode.wav --channels left right -o episode.json`

That should be run in the same directory as the episode was installed into.
