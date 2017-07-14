# movierss

[![Build Status](https://travis-ci.org/teambrookie/movierss.svg?branch=master)](https://travis-ci.org/teambrookie/movierss)
[![Build Status](https://semaphoreci.com/api/v1/fabienfoerster/movierss/branches/master/shields_badge.svg)](https://semaphoreci.com/fabienfoerster/movierss)
[![wercker status](https://app.wercker.com/status/10f00dc08474fea4bb549a53fd3e47e7/s/master "wercker status")](https://app.wercker.com/project/byKey/10f00dc08474fea4bb549a53fd3e47e7)

## Description

MovieRSS is a small app that let you ask Trakt.tv for your movie watchlist and then find the corresponding torrent using RARBG and expose them as an RSS feed.

## Installation

```
docker run -p 8000:8000 -e TRAKT_KEY=xxx teambrookie/movierss
```

## Using it

First of all, what you want is to set when you want the application to refetch the content of your watchlist and when to search for new torrents.

Personally I use IFTTT with a cron job.

First I call the endpoint responsable for refetching the content of the watchlist
```
http://localhost/refresh?action=movie&slug=trakt_username
```

Then I set a cron job for searching the corresponding torrent
```
http://localhost/refresh?action=torrent
```

And finally what really interest us is the /rss endpoint
```
http://localhost/rss?slug=xxx
```


