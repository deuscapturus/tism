tISM - the Immutable Secrets Manager
====================================

tISM is a secrets management solution similiar to Hashicorp's Vault.  But unlike Vault, tISM does not store any of your secrets.  Instead secrets are stored in your own build artifact, git repository or configuration as encrypted pgp messages which are decrypted by the tISM Server.

tISM solves the immutable infrastructure problem of secrets management.

.. WARNING::
   tISM is currently in the very early stages of development.  It is not yet ready for any real use.

Features
--------

* Does not store any secrets.
* Simple. No databases. The only persistent data is a pgp keyring and configuration file.
* Asymmetric encryption with secure and ubiquitous PGP/GPG.  Allows secrets to be encrypted with decentralized public keys.
* Authentication short lived tokens which are also revocable.

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

Create New Encryption Key
-------------------------

.. code::

  curl -H "Content-Type: application/json" -X POST -d '{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6MSwiZXhwIjoxNTg1MTExNDYwLCJqdGkiOiI3NnA5cWNiMWdtdmw4Iiwia2V5cyI6WyJBTEwiXX0.RtAhG6Uorf5xnSf4Ya_GwJnoHkCsql4r1_hiOeDSLzo", "name":"it-operations", "comment":"Production Environment","email":"it-ops@test.com"}' http://localhost:8080/key/new

Encrypt a Message
-----------------

.. code::

  echo -n "sdf@34s#atrsdfgjo" | gpg --batch --trust-model always --encrypt -r "it-operations (Production Environment) <it-ops@test.com>" | base64 -w 0


List Keys
---------

.. code::

  curl -H "Content-Type: application/json" -X POST -d '{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6MSwiZXhwIjoxNTg1MTExNDYwLCJqdGkiOiI3NnA5cWNiMWdtdmw4Iiwia2V5cyI6WyJBTEwiXX0.RtAhG6Uorf5xnSf4Ya_GwJnoHkCsql4r1_hiOeDSLzo"}' http://localhost:8080/key/list

Get Key by Id
-------------

.. code::

  curl -H "Content-Type: application/json" -X POST -d '{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6MSwiZXhwIjoxNTg1MTExNDYwLCJqdGkiOiI3NnA5cWNiMWdtdmw4Iiwia2V5cyI6WyJBTEwiXX0.RtAhG6Uorf5xnSf4Ya_GwJnoHkCsql4r1_hiOeDSLzo","id":"13ec80c75c697055"}' http://localhost:8080/key/get

Issue a new Token
-----------------

.. code::

  curl -H "Content-Type: application/json" -X POST -d '{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6MSwiZXhwIjoxNTg1MTExNDYwLCJqdGkiOiI3NnA5cWNiMWdtdmw4Iiwia2V5cyI6WyJBTEwiXX0.RtAhG6Uorf5xnSf4Ya_GwJnoHkCsql4r1_hiOeDSLzo","keys":["815f99f8f9d435e3","13ec80c75c697055"]' http://localhost:8080/token/new

Decrypt a Secret
----------------

.. code::

  curl -H "Content-Type: application/json" -X POST -d '{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6MSwiZXhwIjoxNTg1MTExNDYwLCJqdGkiOiI3NnA5cWNiMWdtdmw4Iiwia2V5cyI6WyJBTEwiXX0.RtAhG6Uorf5xnSf4Ya_GwJnoHkCsql4r1_hiOeDSLzo   ","GpgContents":"hQEMAzJ+GfdAB3KqAQf+J/LwHFevlL35lZ5W575/QR9DGbWGZGaukDw9OtPDU0EIUvsTdidJweUV1zCuDCOzfE0AZCBebREwcA7z2N+8h3FP9h6otgnrRjkk1rdzIRBN48n6ojFOafIWNOEVFlkD3R9wA4iYx7Ma/GZoKf7cjJciWT59bW95gvnJUaSOOqSpgHKnz/X8KXkFJNkc5wrlPKir1XeI7YNTGbOPDsMXQ83Jrl9fsHr9/r/oPX33yGq7TOeSaCTH37XxPSwskRhM+wuOcobfxH9MxVGnZZf+gOBxD77KFvTN53pboh6wMoDMeera0ScT79XdrooIaRR0hbSJIDZhrPQ3GTZeftNXn8kqE1qgh7zGD9nMfUEL2Y4VOJVyKwzvsRAWTAJFMVzcolSAYCFF9ASkIk7Q"}"}' http://localhost:8080/decrypt

