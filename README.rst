dpl-go
======
|codacy|
|code-climate|
|go-report-card|

This is a rewrite of `dev-pipeline`_ in Go.  The python code evolved
from a quick and dirty scripts for my own use into a project that I as
figuring out as I went along, and the code reflects that history.  In
addition, there are several pain points with python (mostly around
testing everything and parallelism); Go either solves or completely
eliminates those pain points.

This project is designed to be a mostly drop-in replacement, but with
smaller short-term goals.  Instead of aiming for maximum flexibility and
keeping everything as a plugin (like the python version), it will be (at
least for now) much more focused on my needs and common use cases.
Eventually I expect to support plugins, but that's not a high priority
at the moment.

Installing
----------
:code:`dpl` is only released via source.  Make sure you have an
up-to-date Go installed and run the following in your terminal:

.. code-block:: bash

    $ go install github.com/dev-pipeline/dpl-go
    $

If that works, you'll have :code:`dpl-go` as an executable in
:code:`$GOROOT/bin`.  If you prefer the name :code:`dpl`, the included
:code:`Makefile` will do the trick.

.. code-block:: bash

    /path/to/dpl/src $ make
    /path/to/dpl/src $ make install

You can specify an optional destination directory by setting
:code:`DESTDIR` during the call to :code:`make install`.

.. code-block:: bash

    /path/to/dpl/src $ DESTDIR=/usr/local/bin make install

.. _dev-pipeline: https://github.com/dev-pipeline/dev-pipeline

.. |codacy| image:: https://app.codacy.com/project/badge/Grade/74172ff9d3214478a9c33dd4c0339ab9
    :target: https://www.codacy.com/gh/dev-pipeline/dpl-go/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=dev-pipeline/dpl-go&amp;utm_campaign=Badge_Grade
.. |code-climate| image:: https://api.codeclimate.com/v1/badges/8bf6a4d29669138fc13a/maintainability
    :target: https://codeclimate.com/github/dev-pipeline/dpl-go/maintainability
    :alt: Maintainability
.. |go-report-card| image:: https://goreportcard.com/badge/github.com/dev-pipeline/dpl-go
    :target: https://goreportcard.com/report/github.com/dev-pipeline/dpl-go
    :alt: Go Report Card
