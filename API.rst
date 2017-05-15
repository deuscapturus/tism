========================
tISM - API Documentation
========================

.. contents::
    :local:
    
.. NOTICE::

   This is not a REST API.  If it were, it would be an improper implementaion of REST.

POST /encrypt
=============

Encrypt a secret

Request
-------

+------------+--------------------------------------+----------+
| Name       | Description                          | Required |
+============+======================================+==========+
| token      | Authorization Token.                 | yes      |
+------------+--------------------------------------+----------+
| encoding   | Encoding of output.                  | no       |
|            |                                      |          |
|            | Valid options: "base64" or "armor"   |          |
|            | Default: "base64"                    |          |
+------------+--------------------------------------+----------+
| decsecret  | String that will be encrypted.       | yes      |
+------------+--------------------------------------+----------+
| id         | Id of the encryption key to use for  | yes      |
|            | encryption.                          |          |
+------------+--------------------------------------+----------+

Response
--------

Body: Text

Cipher text.

Example
-------

.. code::

  curl -k -H "Content-Type: application/json" -X POST \
  -d '{
      "token" : "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6MSwiZXhwIjo5OTk5OTk5OTk5OSwianRpIjoibDcxY2NmdDhyaWllIiwia2V5cyI6WyJBTEwiXX0.Rtja9H9SSgAy1oMy9AMgXiflC_nZtKLMWToPoN2H8Io",
      "encoding" : "base64",
      "decsecret" : "Th1s$Secret",
      "id" : "815f99f8f9d435e3"
  }' \
  https://localhost:8080/encrypt


POST /decrypt
=============

Decrypt a secret.

Request
-------

+------------+--------------------------------------+----------+
| Name       |  Description                         | Required |
+============+======================================+==========+
| token      | Authorization Token                  | yes      |
+------------+--------------------------------------+----------+
| encoding   | Encoding of cipher text.             | no       |
|            |                                      |          |
|            | Valid options: "base64" or "armor"   |          |
|            | Default: "base64"                    |          |
+------------+--------------------------------------+----------+
| encsecret  | cipher text.                         | yes      |
+------------+--------------------------------------+----------+

Response
--------

Body: Text

Decrypted text.

Example
-------

.. code::

  curl -k -H "Content-Type: application/json" -X POST \
  -d '{
      "token" : "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6MSwiZXhwIjo5OTk5OTk5OTk5OSwianRpIjoibDcxY2NmdDhyaWllIiwia2V5cyI6WyJBTEwiXX0.Rtja9H9SSgAy1oMy9AMgXiflC_nZtKLMWToPoN2H8Io",
      "encoding" : "base64",
      "encsecret" : "hQEMAzJ+GfdAB3KqAQf9E3cyvrPEWR1sf1tMvH0nrJ0bZa9kDFLPxvtwAOqlRiNp0F7IpiiVRF+h+sW5Mb4ffB1TElMzQ+/G5ptd6CjmgBfBsuGeajWmvLEi4lC6/9v1rYGjjLeOCCcN4Dl5AHlxUUaSrxB8akTDvSAnPvGhtRTZqDlltl5UEHsyYXM8RaeCrBw5Or1yvC9Ctx2saVp3xmALQvyhzkUv5pTb1mH0I9Z7E0ian07ZUOD+pVacDAf1oQcPpqkeNVTQQ15EP0fDuvnW+a0vxeLhkbFLfnwqhqEsvFxVFLHVLcs2ffE5cceeOMtVo7DS9fCtkdZr5hR7a+86n4hdKfwDMFXiBwSIPMkmY980N/H30L/r50+CBkuI/u4M2pXDcMYsvvt4ajCbJn91qaQ7BDI="
  }' \
  https://localhost:8080/decrypt


POST /key/new
=============

Creates a new PGP encryption key

Request
-------

+------------+--------------------------------------+----------+
| Name       | Description                          | Required |
+============+======================================+==========+
| token      | Authorization Token                  | yes      |
+------------+--------------------------------------+----------+
| name       | Name for this key                    | yes      |
+------------+--------------------------------------+----------+
| comment    | Comment for this key                 | no       |
+------------+--------------------------------------+----------+
| email      | Email for this key                   | no       |
+------------+--------------------------------------+----------+

Response
--------

Body: JSON

+------------+--------------------------------------+
| Key        | Value                                |
+============+======================================+
| id         | Id for the key created.              |
+------------+--------------------------------------+
| pubkey     | The public key for the key created.  |
+------------+--------------------------------------+

.. code:: json

  {
    "id": "69b2c77142a7efb4",
    "pubkey": "-----BEGIN PGP PUBLIC KEY BLOCK-----\n\nxsBNBFkYZqkBCADBIlaFoEzWTpz3nDbxgzXAsXMdwYHxMEpuduTnSs1mnQualxmN\ncdJjpcL5JEKQAA1kHGM/IrnNqLNzIAsZICiy9uC74BrT6yKAokMOVxrOxFrtr0dV\nB16QbzvIZBj0IbO1/fEoHdt079CUSrMXHhAfe26KxJndPzWKUXK7aGvdGhrPCWb8\n9PimvhU4B7AlKdTf7xVi40xL5uSGUc1MUQKu3Ywb95TCLiwlck0wmwJkdPIXPB8j\nzNXWdyCm9wm79vR1dwmw2n6KxMi/oMfcD508kx+b73shIRiNz9Tc4yIDL7z9xbui\n5fiQ+mwVGUUC8iNS9CqI/sVgW0DrdcziXti1ABEBAAHNOGl0LW9wZXJhdGlvbnMg\nKFByb2R1Y3Rpb24gRW52aXJvbm1lbnQpIDxpdC1vcHNAdGVzdC5jb20+wsBlBBMB\nCAAZBQJZGGapCRBpssdxQqfvtAIbAwIZAQIVCAAAWVsIAI6SUugG84HZbKw/uWCM\nHDPG1Xyq24+TyK9GUaL9+qc07KPVWS7G568RVcD3Fhu1utUiNCj7aXCVfMLJoY5R\nUvi8QpcTzfMzNLU+xZC+mzVjKIg1QJsJvctGcJgfqXp5SKX6B5Lych2g5B/iSHC4\nDGRGHWrhGGkouzNNrPy53rRK/HqmwAGCTRcI5AjPUQqWxpDFzySB3g5FbXbjuIvr\n+kVB3k3VSwo41XY/jGhcSd4XgRA5O3+qAuym8Hw1IDpYVJkEbLxoAKN7Je4xHICa\nY5hEwIJToYV69u7D84A99LtR7P/ptoMvSeYMF+wPe9e4LDxRttv3XSzxUYviIRWQ\nwQHOwE0EWRhmqQEIALH7G23/wJh4xHjV1ZlbwnBo6k8LNe26oaH860S/8numv1B1\nzAcrfe5LZ8mQrqbgfuNUJZa2vmZn3Cn1YaZjnOLOuo4ya1nzQ5zXdLS3tPtErQbF\nHn6JJIMPF1CldJvhgsq8ebrAnmvAZexRfEBD5XHfdL9EX97lzNQkfsXD4d74sXcB\nrBbZT0A/IdfXE5ZDIeeZD+w21cH4auN9h9I3yjJif97KEHQg4XBAlilDu2n0ULEO\n6xKMN9HOfo2chNKjb+02QqYJYN0Ot58TynbR6nhBic/wy/NDvN8msGl2gMZOIaZt\nqQpuG+iivfnCBJW7/FbNFbRDRmqsP1H0fUWuOZsAEQEAAcLAXwQYAQgAEwUCWRhm\nqQkQabLHcUKn77QCGwwAAIrFCAA8SDiwohNnjxSwGchdHfG6k8HwltMP5KhRhXTW\nb2IaItq10qcVTFUPaYEBm6kBEQecMa+WGRYTmShsFvfRVRuqxinHNhr8/jpmSpys\nMu0JACzZWRp5VR9RFR139MIVYXzjOiI6CvMKRFL266y0We6uJA6WRfDOb6aUwEua\ndfcTWl80kBLQVHMqM1HYAR89knHKROo7uxT2S+9yQ52DJ0rTy2m7rN4+u2xulESY\ntr3PK8+vBDYpP77strapgzeQhxgxgto6J46dvaPXlswYiVXLzfLlYoHtjrvulTy0\nm+M6FGz/svjK/CUnUAgc4a8KWXKpoqfj38gMWvHdc7DwDTT0\n=f/mX\n-----END PGP PUBLIC KEY BLOCK-----"
  }

Example
-------

.. code::

  curl -k -H "Content-Type: application/json" -X POST \
  -d '{
      "token" : "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6MSwiZXhwIjo5OTk5OTk5OTk5OSwianRpIjoibDcxY2NmdDhyaWllIiwia2V5cyI6WyJBTEwiXX0.Rtja9H9SSgAy1oMy9AMgXiflC_nZtKLMWToPoN2H8Io",
      "name" : "it-operations",
      "comment" : "Production Environment",
      "email" : "it-ops@test.com"
    }' \
  https://localhost:8080/key/new


POST /key/list
==============

List all keys that are authorized to a token.

Request
-------

+------------+--------------------------------------+----------+
| Name       | Description                          | Required |
+============+======================================+==========+
| token      | Authorization Token                  | yes      |
+------------+--------------------------------------+----------+

Response
--------

Body: JSON

List of dictionaries.

+--------------+--------------------------------------+
| Key          | Value                                |
+==============+======================================+
| Id           | Id for the key.  If "ALL", this token|
|              | is authorized to all keys, current   |
|              | and future.                          |
+--------------+--------------------------------------+
| CreationTime | Creation time for the key.           |
+--------------+--------------------------------------+
| Name         | The name of the key in GPG format.   |
+--------------+--------------------------------------+

.. code:: json

  [
    {
      "Id": "ALL"
    },
    {
      "CreationTime": "2015-10-20 08:24:59 -0600 MDT",
      "Id": "a14f89ugcsdf4777",
      "Name": "team-dev"
    },
    {
      "CreationTime": "2017-05-14 08:16:05 -0600 MDT",
      "Id": "sd0f93o4jsiojf8b",
      "Name": "it-operations (Production Environment) <it-ops@test.com>"
    },
    {
      "CreationTime": "2017-05-14 08:16:09 -0600 MDT",
      "Id": "69b2c77142a7efb4",
      "Name": "it-operations (Production Environment) <it-ops@test.com>"
    }
  ]

Example
-------

.. code::

  curl -k -s -H "Content-Type: application/json" -X POST \
  -d '{
      "token" : "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6MSwiZXhwIjo5OTk5OTk5OTk5OSwianRpIjoibDcxY2NmdDhyaWllIiwia2V5cyI6WyJBTEwiXX0.Rtja9H9SSgAy1oMy9AMgXiflC_nZtKLMWToPoN2H8Io"
  }' \
  https://localhost:8080/key/list


POST /key/get
=============

Get a encryption key details by key Id.

Request
-------

+------------+--------------------------------------+----------+
| Name       |  Description                          | Required |
+============+======================================+==========+
| token      | Authorization Token                  | yes      |
+------------+--------------------------------------+----------+
| id         | Id of the encryption key.            | yes      |
+------------+--------------------------------------+----------+

Response
--------

Body: JSON

+------------+--------------------------------------+
| Key        | Value                                |
+============+======================================+
| id         | Id of the encryption key.            |
+------------+--------------------------------------+
| pubkey     | The public key for the key requested.|
+------------+--------------------------------------+

.. code:: json

  {
    "id": "69b2c77142a7efb4",
    "pubkey": "-----BEGIN PGP PUBLIC KEY BLOCK-----\n\nxsBNBFkYZqkBCADBIlaFoEzWTpz3nDbxgzXAsXMdwYHxMEpuduTnSs1mnQualxmN\ncdJjpcL5JEKQAA1kHGM/IrnNqLNzIAsZICiy9uC74BrT6yKAokMOVxrOxFrtr0dV\nB16QbzvIZBas34spd9g8$#T9CUSrMXHhAfe26KxJndPzWKUXK7aGvdGhrPCWb8\sdgimvhU4B7AlKdTf7xVi40xL5asdgc1MUQKu3Ywb95TCLiwlck0wmwJkdPIXPB8j\nzNXWdyCm9wm79vR1dwmw2n6KxMi/%^ucD508kx+b73shIRiNz9Tc4yIDL7z9xbuiASGfiQ+mwVGUUC8iNS9CqI/sVgW0DrdcziXti1ABEBAAHNOGl0LW9wZXJhdGlvbnMg\nKFByb2R1YsGASRGgRW52aXJvbm1lbnQpIDxpdC1vcHNAdGVzdC5jb20+wsBlBBMB\nCAAZBQJZGGapCRBpssdxQqfvtAIbAwIZAQIVCAAAWVsIAI6SUugG84HZbKwSRGSRGR3434534524+TyK9GUaL9+qc07KPVWS7G568RVcD3Fhu1usDGASDGASCVfMLJoY5R\sdgi8QpcTzfMzNLU+xZC+dfhjKIg1QJsJvctGcJgfqXp5SKX6B5Lych2g5B/iSHC4\dfhRGHWrhGGkouzNNrPy53rRK/HqmwAGCTRcI5AjPUQqWxpDFzySB3g5FbXbjuIvr\n+DJsdk3VSwo41XY/dJDSBZxbgRA5O3+qAuym8Hw1dDdYVJkEbLxoAKN7Je4xHICa\ndfhEwIJToYV69u7D84A99LtR7P/ptoMvSeYMF+wPe9e4LDxRtsdfhSzxUYviIRWQ\nwQHOwE0EWRhmqQEIALH7G23/wJh4xHjV1ZlbwnBo6k8LNe26oaH860S/8numv1B1\nzAcrfe5LZ8mQrqbgfuNUJZa2vmZn3Cn1YaZjnOLOuo4ya1nzQ5zXdLS3tPtErQbF\nHn6JJIMPF1CldJvhgsq8ebrAnmvAZexRfEBD5XHfdL9EX97lzNQkfsXD4d74sXcB\nrBbZT0A/IdfXE5ZDIeeZD+w21cH4auN9h9I3yjJif97KEHQg4XBAlilDu2n0ULEO\n6xKMN9HOfo2chNKjb+02QqYJYN0Ot58TynbR6nhBic/wy/NDvN8msGl2gMZOIaZt\nqQpuG+iivfnCBJW7/FbNFbRDRmqsP1H0fUWuOZsAEQEAAcLAXwQYAQgAEwUCWRhm\nqQkQabLHcUKn77QCGwwAAIrFCAA8SDiwohNnjxSwGchdHfG6k8HwltMP5KhRhXTW\nb2IaItq10qcVTFUPaYEBm6kBEQecMa+WGRYTmShsFvfRVRuqxinHNhr8/jpmSpys\nMu0JACzZWRp5VR9RFR139MIVYXzjOiI6CvMKRFL266y0We6uJA6WRfDOb6aUwEua\ndfcTWl80kBLQVHMqM1HYAR89knHKROo7uxT2S+9yQ52DJ0rTy2m7rN4+u2xulESY\ntr3PK8+vBDYpP77strapgzeQhxgxgto6J46dvaPXlswYiVXLzfLlYoHtjrvulTy0\nm+M6FGz/svjK/CUnUAgc4a8KWXKpoqfj38gMWvHdc7DwDTT0\n=f/mX\n-----END PGP PUBLIC KEY BLOCK-----"
  }

Example
-------

.. code::

  curl -k -H "Content-Type: application/json" -X POST \
  -d '{
      "token" : "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6MSwiZXhwIjo5OTk5OTk5OTk5OSwianRpIjoibDcxY2NmdDhyaWllIiwia2V5cyI6WyJBTEwiXX0.Rtja9H9SSgAy1oMy9AMgXiflC_nZtKLMWToPoN2H8Io",
      "id" : "13ec80c75c697055"
  }' \
  https://localhost:8080/key/get


POST /key/delete
================

Delete a key by id.

Request
-------

+------------+--------------------------------------+----------+
| Name       | Description                          | Required |
+============+======================================+==========+
| token      | Authorization Token                  | yes      |
+------------+--------------------------------------+----------+
| id         | Id of the encryption key.            | yes      |
+------------+--------------------------------------+----------+

Response
--------

Body: None

Example
-------

.. code::

  curl -k -H "Content-Type: application/json" -X POST \
  -d '{
      "token" : "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6MSwiZXhwIjo5OTk5OTk5OTk5OSwianRpIjoibDcxY2NmdDhyaWllIiwia2V5cyI6WyJBTEwiXX0.Rtja9H9SSgAy1oMy9AMgXiflC_nZtKLMWToPoN2H8Io",
      "id" : "13ec80c75c697055"
  }' \
  https://localhost:8080/key/delete


POST /token/new
===============

Get a new authorization token.

Request
-------

+------------+--------------------------------------+----------+
| Name       | Description                          | Required |
+============+======================================+==========+
| token      | Authorization Token                  | yes      |
+------------+--------------------------------------+----------+
| keys       | List of encryption keys by id that   | yes      |
|            | this token will be authorized to use.|          |
+------------+--------------------------------------+----------+
| admin      | Whether or not to make this token an | yes      |
|            | admin token.  Admin token can create |          |
|            | new tokens and delete keys.          |          |
|            |                                      |          |
|            | Valid options: 0 or 1                |          |
+------------+--------------------------------------+----------+

Response
--------

Body: Text

Token

Example
-------

.. code::

  curl -k -H "Content-Type: application/json" -X POST \
  -d '{
      "token" : "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6MSwiZXhwIjo5OTk5OTk5OTk5OSwianRpIjoibDcxY2NmdDhyaWllIiwia2V5cyI6WyJBTEwiXX0.Rtja9H9SSgAy1oMy9AMgXiflC_nZtKLMWToPoN2H8Io",
      "keys" : ["815f99f8f9d435e3","13ec80c75c697055"]
  }' \
  https://localhost:8080/token/new


POST /token/info
================

Get information for a token.

Request
-------

+------------+--------------------------------------+----------+
| Name       | Description                          | Required |
+============+======================================+==========+
| token      | Authorization Token                  | yes      |
+------------+--------------------------------------+----------+

Response
--------

Body: JSON

+------------+--------------------------------------+
| Key        | Value                                |
+============+======================================+
| keys       | List if all keys by Id this token is |
|            | authorized to.                       |
+------------+--------------------------------------+
| admin      | Token admin status.                  |
|            |                                      |
|            | Return options: 0 or 1               |
+------------+--------------------------------------+

.. code:: json

  {
    "keys": [
      "ALL"
    ],
    "admin": 1
  }

Example
-------

.. code::

  curl -k -H "Content-Type: application/json" -X POST \
  -d '{
      "token" : "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6MSwiZXhwIjo5OTk5OTk5OTk5OSwianRpIjoibDcxY2NmdDhyaWllIiwia2V5cyI6WyJBTEwiXX0.Rtja9H9SSgAy1oMy9AMgXiflC_nZtKLMWToPoN2H8Io"
  }' \
  https://localhost:8080/token/info
