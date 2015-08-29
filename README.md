README
======

A chatbot, that queries various SOLR indices.

Install
-------

    $ go get github.com/miku/halbot


Add solr aliases
----------------

Specify SOLR urls in a [`$HOME/.parrotrc`](https://github.com/miku/halbot/blob/master/.parrotrc).


Start server
------------

    $ HAL_ADAPTER=irc HAL_IRC_USER=hal HAL_IRC_NICK=hal HAL_IRC_SERVER=x.y.com HAL_IRC_CHANNELS="#chan" halbot
