# movierss

[![Build Status](https://travis-ci.org/teambrookie/movierss.svg?branch=master)](https://travis-ci.org/teambrookie/movierss)
[![Build Status](https://semaphoreci.com/api/v1/fabienfoerster/movierss/branches/master/shields_badge.svg)](https://semaphoreci.com/fabienfoerster/movierss)
[![wercker status](https://app.wercker.com/status/10f00dc08474fea4bb549a53fd3e47e7/s/master "wercker status")](https://app.wercker.com/project/byKey/10f00dc08474fea4bb549a53fd3e47e7)

## Description

MovieRSS is a small app that let you ask Trakt.tv for your movie watchlist and then find the corresponding torrent using RARBG and expose them as an RSS feed.

### Using it


Then they are the /refresh endpoint, it's role is to refresh the unseen episode and to refresh the torrent. They are use like this:
```
http://yourdomain/refresh?action=movie&slug=xxx
http://yourdomain/refresh?action=torrent
```

And finally what really interest us is the /rss endpoint
```
http://yourdomain/rss?slug=xxx
```

docker run -p 8000:8000 -e TRAKT_KEY=xxx movierss
