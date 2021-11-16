dpl-go
======
|codacy|
|code-climate|
|lgtm|
|lgtm-quality|

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

.. |codacy| image:: https://app.codacy.com/project/badge/Grade/74172ff9d3214478a9c33dd4c0339ab9
    :target: https://www.codacy.com/gh/dev-pipeline/dpl-go/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=dev-pipeline/dpl-go&amp;utm_campaign=Badge_Grade
.. |code-climate| image:: https://api.codeclimate.com/v1/badges/8bf6a4d29669138fc13a/maintainability
    :target: https://codeclimate.com/github/dev-pipeline/dpl-go/maintainability
    :alt: Maintainability
.. |lgtm| image:: https://img.shields.io/lgtm/alerts/g/dev-pipeline/dpl-go.svg?logo=lgtm&logoWidth=18
    :target: https://lgtm.com/projects/g/dev-pipeline/dpl-go/alerts/
.. |lgtm-quality| image:: https://img.shields.io/lgtm/grade/go/g/dev-pipeline/dpl-go.svg?logo=lgtm&logoWidth=18
    :target: https://lgtm.com/projects/g/dev-pipeline/dpl-go/context:go
