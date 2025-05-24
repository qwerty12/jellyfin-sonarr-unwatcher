Unmonitor episodes in Sonarr when watching in Jellyfin

Basic Go clone of unmonitorr except that this only supports Jellyfin and Sonarr. (I don't use Radarr.) This contains the benefits (IMO) of [my unmonitorr fork](https://github.com/qwerty12/unmonitorr).

I don't like JavaScript, don't like that Node uses 60 MB of private memory for just this and find the Node ecosystem awkward (to deal with on Windows). (Though, while still much better by comparison, Go's 14 MB isn't anything to write home about, either...)

Unmonitorr's original README goes into how to set this up with Jellyfin, including configuring shemanaev's Webhook plugin.
The procedure is the same here; however, the `SONARR_HOST` environment variable must be explicitly set by you.

This also rips off prefetcharr - to enable, set the `REMAINING_EPISODES` environment variable to something like `2`, and configure the following webhook:

| Option         | Value                               |
| -------------- | ----------------------------------- |
| URL            | `http://127.0.0.1:9898/prefetcharr` |
| Payload format | Default                             |
| Events         | Scrobble                            |

---

Credits:

* [`Shraymonks/unmonitorr`](https://github.com/Shraymonks/unmonitorr), [`salvo-github/unmonitorr` (Jellyfin support)](https://github.com/Shraymonks/unmonitorr/pull/22)
    <details>
    <summary>LICENSE</summary>
    
    ```
    ISC License
    
    Copyright (c) 2022 Raymond Ha
    
    Permission to use, copy, modify, and/or distribute this software for any
    purpose with or without fee is hereby granted, provided that the above
    copyright notice and this permission notice appear in all copies.
    
    THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH
    REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY
    AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT,
    INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM
    LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR
    OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR
    PERFORMANCE OF THIS SOFTWARE.
    ```
    
    </details>

* [`p-hueber/prefetcharr`](https://github.com/p-hueber/prefetcharr)
    <details>
    <summary>LICENSE</summary>
    
    ```
    Copyright (c) 2024 Paul Hüber
    
    Permission is hereby granted, free of charge, to any person obtaining a copy
    of this software and associated documentation files (the "Software"), to deal
    in the Software without restriction, including without limitation the rights
    to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
    copies of the Software, and to permit persons to whom the Software is
    furnished to do so, subject to the following conditions:
    
    The above copyright notice and this permission notice shall be included in all
    copies or substantial portions of the Software.
    
    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
    IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
    FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
    AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
    LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
    OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
    SOFTWARE.
    ```
    
    </details>
