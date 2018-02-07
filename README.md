# fedviz

Generates an image of the Fediverse.

```
go get github.com/cjslep/fedviz
go install github.com/cjslep/fedviz
```

The data can be found [here](https://gist.github.com/Gargron/48e67b1b14723cd178c951fe7f373a38).

Then you can simply run `./fedviz` with the data.

The program is structured so different algorithms can easily be swapped in and
out of the main program depending on if you want to hack on layout or color
palette or image format.

Each instance will be assigned a different color. The number of pixels of each
color corresponds to how well connected each instance is. The larger patches of
color federate with more peers than the smaller ones.

The default result is:

![Fediverse](https://github.com/cjslep/fedviz/blob/master/result.png "Fediverse")

Please consider donating to the [Mastodon](https://liberapay.com/Mastodon/)
project, for helping grow the Fediverse and providing this data.
