# cliphist-wofi-img

Quick hack to display images in wofi for cliphist binary image history entries.

## Install

```shell
go install github.com/pdf/cliphist-wofi-img@latest
````

## Usage

- Ensure that both this application and cliphist are on the path.
- Set `pre_display_exec=true` in your wofi config (available as of wofi hg commit `846b4aeae44c`)
- Run wofi as follows:

```shell
cliphist list | wofi --dmenu --allow-images --pre-display-cmd "cliphist-wofi-img %s" | cliphist decode | wl-copy
```

