====================================
tISM - the Immutable Secrets Manager
====================================

tISM is PGP encryption-as-a-service for secrets management.  Decrypt/Encrypt PGP secrets via API with token authorization.

tISM solves the immutable infrastructure problem of secrets management.

.. contents::
    :local:
    
.. WARNING::

   Use at your own risk!

Features
========

* Does not store any secrets.
* Simple. No databases. The only persistent data is a pgp keyring and configuration file.
* Asymmetric encryption with secure and ubiquitous PGP/GPG.  Allows secrets to be encrypted with distributed public keys.
* Authorization with short lived and revocable JWT tokens.

Security
========

tISM relies on 3 separated components to access secrets.

1.  Access Token.
2.  PGP Encrypted Message
3.  tISM Server

Access Tokens are implemented with JSON Web Token https://tools.ietf.org/html/rfc7519
Message Encryption and Decryption is implemented with OpenPGP https://tools.ietf.org/html/rfc4880

Quick Start
===========

Installation
------------

Install rpm or run container image.

RPM
^^^

sudo dnf install https://github.com/deuscapturus/tism/releases/download/0.0/tism-0.0-1.fc25.x86_64.rpm

docker container
^^^^^^^^^^^^^^^^

docker run tism/tism:0.0

Start tismd
-----------

Initialize
^^^^^^^^^^

First generate a TLS cert and admin token

::

  tism -t -c -n
  2016/10/15 10:22:55 eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6MSwiZXhwIjo5OTk5OTk5OTk5OSwianRpIjoiM3QwOGQxN2VxZHVtcCIsImtleXMiOlsiQUxMIl19.bCBGHR8hCfLT5Pb4iek12T-jawPtX0xINbvhmqG9Jzs
  2016/10/15 10:22:56 written ./cert/cert.crt
  2016/10/15 10:22:56 written ./cert/cert.key

`-t` generates a token, `-c` generates the TLS cert, `-n` tells tism to not start the tism server.

Run
^^^


Now you are ready to run tism

::

   tism

or

::

   systemctl start tism

Web UI  
======

To use the web ui your browser must have es6 module support enabled (a very new feature).

Currently Supported Browers:

- Firefox 54 or greater with `dom.moduleScripts.enabled`
- Safari 10.1 or greater

https://localhost:8080

.. image:: tism-web-ui.png

API
===

[[API.rst]]
