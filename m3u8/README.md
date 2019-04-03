# Live stream record system

Target to write a litte software to record
ORF live streams by using m3u8 library to
get the informaion of the .ts files to
download.


# m3u8 playlist

First read the masterplaylist. The masterplaylist contains
the variant of all stream with different qualities.
Each variant has a URL where you get the real mediaplaylist.
In a live stream there are only x segments inside to download.


# Setup

Install m3u8 library

`> go get github.com/grafov/m3u8`

# m3u8 Playlist file

## Master Playlist

Master playlist shows all available qualities and url to playlist file

    #EXTM3U
    #EXT-X-VERSION:3
    #EXT-X-STREAM-INF:BANDWIDTH=2194075,RESOLUTION=960x540,CODECS="avc1.4D401F,mp4a.40.2"
    manifest_4.m3u8?m=1552488594


## Playlist Live Stream (Sliding Window)

    #EXTM3U
    #EXT-X-VERSION:3
    #EXT-X-TARGETDURATION:10
    #EXT-X-MEDIA-SEQUENCE:200754
    #EXTINF:10.000,
    manifest_4_200754.ts?m=1552488594&f=5
    #EXTINF:10.000,
    manifest_4_200755.ts?m=1552488594&f=5
    #EXTINF:10.000,
    manifest_4_200756.ts?m=1552488594&f=5



