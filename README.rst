dpl-go
======
This is a rewrite of `dev-pipeline`_ in golang.  The python code evolved
from a quick and dirty scripts for my own use into a project that I as
figuring out as I went along, and the code reflects that history.  In
addition, there are several pain points with python.  Golang solves a
number of my pain points, hence that choice.

This project is designed to be a mostly drop-in replacement, but with
smaller short-term goals.  Instead of aiming for maximum flexibility and
keeping everything as a plugin (like the python version), it will be (at
least for now) much more focused on my needs and common use cases.
Eventually I expect to support plugins, but that's not a high priority
at the moment.

.. _dev-pipeline: https://github.com/dev-pipeline/dev-pipeline
