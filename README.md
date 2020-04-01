<p align="center"><img src="https://www.pokepedia.fr/images/b/b6/Luxio-DP.png" width="192px" alt="Luxio © pokepedia.fr" /></p>

# luxio

Luxio is the first evolution form of [S(p)hnix](https://eprint.iacr.org/2018/695).
It's another password manager.

## What is this thing?

Luxio is an implementation of the Sphinx (a password **S**tore that
**P**erfectly **H**ides from **I**tself (**N**o **X**aggeration)) scheme, with
some extra features (while still guarantees the same security).

For an overview of Sphinx, please see
[this talk](https://www.youtube.com/watch?v=px8hiyf81iM)
presented by the Levchin Prize winner 2018 Hugo Krawczyk on Real World Crypto.
His slides are available [here](https://rwc.iacr.org/2017/Slides/hugo.krawczyk.pdf).

Basically, this is different from other password managers which store the passwords in an encrypted database.
It prevents
- offline dictionary attack: whoever has your database cannot try to crack your master password on it (by brute-forcing your master password).
- "on-the-wire" attack, as long as your machine is secure when you are typing your master password.

### Features:
All Sphinx's features including:
- All domain-specific passwords in a user's device or *online*.
- User *memorizes* a single master password.
- All passwords are *random* and *independent* of each other.

And, Luxio supports:
- Generate human-readable passwords. Luxio goes a step further, and convert
Sphinx's password strings to human-readable passphrase.
Each passphrase is a combination of 4 random word, providing `64` bits of entropy, which is strong for most online services.
Luxio uses [`niceware`](https://github.com/diracdeltas/niceware) for this task.
- Change domain-specific passwords.
The new password is also *random* and *independent* of the old one.
This can be done by requiring the device to store a database of random salts for
each account.
However, this database can still be made *online* (*public*).
Even the **metadata** of the user's accounts has been eliminated (that is, no one
can see the username or the domain information).
- **Forward Secrecy**. If some domain-specific passwords are leaked, the attacker
cannot compute the new password from the old ones (even with the database).
However, note that a brute-force attack can be performed if a per-domain password
and the device key are leaked.
This attack is unavoidable for any password-based scheme.

## Installation

Currently, this package provides a binary program to run on your machine,
without needing to separate the device and the client.

You need to have Go version 1.10 or higher installed.
If Go is set up correctly, you can simply run:

```
go install github.com/vqhuy/luxio/cli/luxio
```

## Usage

#### Generate a new device key

```
❯❯❯ luxio keygen
# an ASCII-armored string.
```

#### Create a config file

```
❯❯❯ luxio init
```

This command generates a defaul config file at `~/.luxiorc` and a directory
for storing the database and the device key at `~/.luxio/`.
A random device key will also be generated at `~/.luxio/key.luxio`.

The config file is as follows.

```
DB = "/absolute/path/to/my/password/store"
###
# Turn this flag on if you want to hide your metadata from the database
HideMetadata = true
# ... or turn it off so you can use the `list` functionality.
# HideMetadata = false
###
# Choose either this
Key = "my-hex-encoded-device-key"
# or this
# KeyEval = "run-this-command-to-get-my-ASCII-armored-device-key"
```

#### Retrieve password

```
❯❯❯ luxio request -h
Get password of the given account on the given domain

Usage:
  luxio request [flags]

Examples:
luxio request -d "https://accounts.google.com/" -u "name@gmail.com"

Flags:
  -d, --domain string     Domain name (an URL, a website, etc)
  -h, --help              help for request
      --pin               Print as a PIN code
      --plain             Print as a plain, lower-case passphrase
      --special           Print as a title-case passphrase with a number and a special character
  -u, --username string   Username or Account
❯❯❯ luxio request -d "domain" -u "username"
❯ Enter your Master Password:
```

#### Change password

```
❯❯❯ luxio update -d "domain" -u "username"
❯ Enter your Master Password:
```

#### List all accounts
(This is only supported if `HideMetadata` is `false`.)

```
❯❯❯ luxio list "domain"
o
├──domain.com
|  ├──account1
|  └──account2
└──subdomain.domain.com
   └──account
```
You can also run it as `luxio list "*"` to get all accounts information
from the database.

## Disclaimer

As usual, use at your own risk.
