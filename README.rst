tISM - the Immutable Secrets Manager
====================================

tISM is PGP encryption-as-a-service for secrets management.  Decrypt/Encrypt PGP secrets via REST with token based authorization or web based UI.

tISM solves the immutable infrastructure problem of secrets management.

.. WARNING::

   Use at your own risk!

.. contents::
    :local:

Features
--------

* Does not store any secrets.
* Simple. No databases. The only persistent data is a pgp keyring and configuration file.
* Asymmetric encryption with secure and ubiquitous PGP/GPG.  Allows secrets to be encrypted with distributed public keys.
* Authorization with short lived and revocable JWT tokens.

Security
--------

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
------

To use the web ui your browser must have es6 module support enabled (a very new feature).

Currently Supported Browers:

- Firefox 54 or greater with `dom.moduleScripts.enabled`
- Safari 10.1 or greater

https://localhost:8080

.. image:: tism-web-ui.png

REST API
--------

.. code::

  # go run tismd.go -c -t
  2016/10/15 10:22:55 eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6MSwiZXhwIjo5OTk5OTk5OTk5OSwianRpIjoiM3QwOGQxN2VxZHVtcCIsImtleXMiOlsiQUxMIl19.bCBGHR8hCfLT5Pb4iek12T-jawPtX0xINbvhmqG9Jzs
  2016/10/15 10:22:56 written ./cert/cert.crt
  2016/10/15 10:22:56 written ./cert/cert.key

Create New Encryption Key
^^^^^^^^^^^^^^^^^^^^^^^^^

.. code::

  curl -k -H "Content-Type: application/json" -X POST \
  -d '{
      "token" : "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6MSwiZXhwIjoxNTg1MTExNDYwLCJqdGkiOiI3NnA5cWNiMWdtdmw4Iiwia2V5cyI6WyJBTEwiXX0.RtAhG6Uorf5xnSf4Ya_GwJnoHkCsql4r1_hiOeDSLzo",
      "name" : "it-operations",
      "comment" : "Production Environment",
      "email" : "it-ops@test.com"
    }' \
  https://localhost:8080/key/new

Encrypt a Message
^^^^^^^^^^^^^^^^^

.. code::

  echo -n "sdf@34s#atrsdfgjo" | gpg --batch --trust-model always --encrypt -r "it-operations (Production Environment) <it-ops@test.com>" | base64 -w 0

List Keys
^^^^^^^^^

.. code::

  curl -k -s -H "Content-Type: application/json" -X POST \
  -d '{
      "token" : "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6MSwiZXhwIjoxNTg1MTExNDYwLCJqdGkiOiI3NnA5cWNiMWdtdmw4Iiwia2V5cyI6WyJBTEwiXX0.RtAhG6Uorf5xnSf4Ya_GwJnoHkCsql4r1_hiOeDSLzo"
  }' \
  https://localhost:8080/key/list

Get Key by Id
^^^^^^^^^^^^^

.. code::

  curl -k -H "Content-Type: application/json" -X POST \
  -d '{
      "token" : "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6MSwiZXhwIjoxNTg1MTExNDYwLCJqdGkiOiI3NnA5cWNiMWdtdmw4Iiwia2V5cyI6WyJBTEwiXX0.RtAhG6Uorf5xnSf4Ya_GwJnoHkCsql4r1_hiOeDSLzo",
      "id" : "13ec80c75c697055"
  }' \
  https://localhost:8080/key/get

Delete Key by Id
^^^^^^^^^^^^^^^^

.. code::

  curl -k -H "Content-Type: application/json" -X POST \
  -d '{
      "token" : "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6MSwiZXhwIjoxNTg1MTExNDYwLCJqdGkiOiI3NnA5cWNiMWdtdmw4Iiwia2V5cyI6WyJBTEwiXX0.RtAhG6Uorf5xnSf4Ya_GwJnoHkCsql4r1_hiOeDSLzo",
      "id" : "13ec80c75c697055"
  }' \
  https://localhost:8080/key/delete

Issue a new Token
^^^^^^^^^^^^^^^^^

.. code::

  curl -k -H "Content-Type: application/json" -X POST \
  -d '{
      "token" : "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6MSwiZXhwIjoxNTg1MTExNDYwLCJqdGkiOiI3NnA5cWNiMWdtdmw4Iiwia2V5cyI6WyJBTEwiXX0.RtAhG6Uorf5xnSf4Ya_GwJnoHkCsql4r1_hiOeDSLzo",
      "keys" : ["815f99f8f9d435e3","13ec80c75c697055"]
  }' \
  https://localhost:8080/token/new

Get Token Info
^^^^^^^^^^^^^^

.. code::

  curl -k -H "Content-Type: application/json" -X POST \
  -d '{
      "token" : "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6MSwiZXhwIjoxNTg1MTExNDYwLCJqdGkiOiI3NnA5cWNiMWdtdmw4Iiwia2V5cyI6WyJBTEwiXX0.RtAhG6Uorf5xnSf4Ya_GwJnoHkCsql4r1_hiOeDSLzo"
  }' \
  https://localhost:8080/token/info

Encrypt a Secret
^^^^^^^^^^^^^^^^

.. code::

  curl -k -H "Content-Type: application/json" -X POST \
  -d '{
      "token" : "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6MSwiZXhwIjoxNTg1MTExNDYwLCJqdGkiOiI3NnA5cWNiMWdtdmw4Iiwia2V5cyI6WyJBTEwiXX0.RtAhG6Uorf5xnSf4Ya_GwJnoHkCsql4r1_hiOeDSLzo",
      "encoding" : "base64"
      "decsecret" : "Th1s$Secret",
      "id" : "815f99f8f9d435e3"
  }' \
  https://localhost:8080/encrypt

Decrypt a Secret
^^^^^^^^^^^^^^^^

.. code::

  curl -k -H "Content-Type: application/json" -X POST \
  -d '{
      "token" : "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6MSwiZXhwIjoxNTg1MTExNDYwLCJqdGkiOiI3NnA5cWNiMWdtdmw4Iiwia2V5cyI6WyJBTEwiXX0.RtAhG6Uorf5xnSf4Ya_GwJnoHkCsql4r1_hiOeDSLzo",
      "encoding" : "base64"
      "encsecret" : "hQEMAzJ+GfdAB3KqAQf9E3cyvrPEWR1sf1tMvH0nrJ0bZa9kDFLPxvtwAOqlRiNp0F7IpiiVRF+h+sW5Mb4ffB1TElMzQ+/G5ptd6CjmgBfBsuGeajWmvLEi4lC6/9v1rYGjjLeOCCcN4Dl5AHlxUUaSrxB8akTDvSAnPvGhtRTZqDlltl5UEHsyYXM8RaeCrBw5Or1yvC9Ctx2saVp3xmALQvyhzkUv5pTb1mH0I9Z7E0ian07ZUOD+pVacDAf1oQcPpqkeNVTQQ15EP0fDuvnW+a0vxeLhkbFLfnwqhqEsvFxVFLHVLcs2ffE5cceeOMtVo7DS9fCtkdZr5hR7a+86n4hdKfwDMFXiBwSIPMkmY980N/H30L/r50+CBkuI/u4M2pXDcMYsvvt4ajCbJn91qaQ7BDI="
  }' \
  https://localhost:8080/decrypt
