Unmonitor episodes in Sonarr when watching in Jellyfin

Basic Go clone of https://github.com/Shraymonks/unmonitorr except that this only supports Jellyfin and Sonarr. (I don't use Radarr.) This is derived from [my unmonitorr fork](https://github.com/qwerty12/unmonitorr).

I don't like JavaScript, don't like that Node uses 60 MB of private memory for just this and find the Node ecosystem awkward (to deal with on Windows).

Unmonitorr's original README goes into how to set this up with Jellyfin, including configuring shemanaev's Webhook plugin.
However, here you must set the `SONARR_HOST` environment variable yourself.

Thanks to Shraymonks for writing unmonitorr in the first place, and to salvo-github for adding Jellyfin support to it.