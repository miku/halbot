README
======

A chatbot, that queries various SOLR indices.

Install
-------

    $ go get github.com/miku/halbot


Add solr aliases
----------------

Specify SOLR urls in a [`$HOME/.halrc`](https://github.com/miku/halbot/blob/master/.halrc).


Start server
------------

    $ HAL_ADAPTER=irc HAL_IRC_USER=hal HAL_IRC_NICK=hal \
      HAL_IRC_SERVER=x.y.com HAL_IRC_CHANNELS="#chan" halbot

Query
-----

Query syntax is `hal <alias> q <query>`, e.g.

    [16:01] <        human> | hal ai q "Roboterarmee"
    [16:01] < hal> 7 in ai for "Roboterarmee"

The first titles can be queries with:

    [16:02] <        human> | hal ai qq "Roboterarmee"
    [16:02] < hal> 7 in ai for "Roboterarmee" -- (1) Amazon startet Roboter-Armee [48],
                   (2) Ballzauberer gegen das Böse [48], (3) TERMINATOR 3
