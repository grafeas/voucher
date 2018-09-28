# Voucher Tutorial

This tutorial should give you an overview on how to set up Voucher. These instructions will work with both Voucher Server (`voucher_server`) and Voucher Standalone (`voucher_cli`).

## Installing Voucher

First, install Voucher using the instructions in the [README](README.md).

## Create configuration

Create a new configuration file. You may want to refer to the [example configuration](config/config.toml).

Make sure that you update the `ejson` specific blocks to point to the `ejson` file and key you will be making in the next step.

## Create ejson configuration

If you plan on creating attestations (rather than just running checks against your images), or if you plan on using Clair as your Vulnerability Scanner, you will need to create an `ejson` file to store the OpenPGP keys and/or Clair login information respectively.

Note: this step is unnecessary if you are a Shopify employee and are running Voucher in Shopify's cloud platform.

First, create a public and private key pair (this is unnecessary if your platform automatically generates a key for you):

```shell
$ ejson keygen
```

You will see output similar to the following:

```
Public Key:
45960c1576c3a5caa13fea5630b56ba2c48dd67b9701bddf9f24666123122306
Private Key:
0de9401520d770fe9cd4bc985dd949a3f537d338f3954ba12c65be07c5e4637f
```

You can then create an `ejson` file using the public key included. For example:

```json
{
    "_public_key": "<public key>",
    "openpgpkeys": {},
    "clair": {}
}
```

The private key should then be stored in a file, where the filename is the public key, and the body of the file is the private key.

For example:

```shell
$ echo 0de9401520d770fe9cd4bc985dd949a3f537d338f3954ba12c65be07c5e4637f > 45960c1576c3a5caa13fea5630b56ba2c48dd67b9701bddf9f24666123122306
```

You can now decrypt your `ejson` file using:

```shell
$ ejson --keydir=<path to the directory containing the keyfile> decrypt <filename of your ejson file>
```

## Generating Keys for Attestation

Attestation uses GPG signing keys to sign the image. It's suggested that you use a primary GPG key instead of a subkey.

To generate a signing key, use the following command. This will ensure that you're generating a new key that only can
be used for signing (as the other attributes are unnecessary).

```
$ gpg --full-generate-key
```

You will first be asked what type of key you want to create. Select an RSA signing key.

```
Please select what kind of key you want:
   (1) RSA and RSA (default)
   (2) DSA and Elgamal
   (3) DSA (sign only)
   (4) RSA (sign only)
Your selection? 4
```

Next you will be prompted for the key length. We'll use the largest possible value, 4096 bits.

```
RSA keys may be between 1024 and 4096 bits long.
What keysize do you want? (2048) 4096
Requested keysize is 4096 bits       
```

When prompted for how long the key should be valid, you may want to specify that the key does not expire, especially if the team maintaining your Binary Authorization configuration is the same team managing your Voucher install.

```
Please specify how long the key should be valid.
         0 = key does not expire
      <n>  = key expires in n days
      <n>w = key expires in n weeks
      <n>m = key expires in n months
      <n>y = key expires in n years
Key is valid for? (0) 0
Key does not expire at all
Is this correct? (y/N) y
```

You will next be asked to provide an ID for the GPG key. You will want to add a comment to claify which Check this key is for.

```                        
GnuPG needs to construct a user ID to identify your key.

Real name: Cloud Security 
Email address: cloudsecurityteam@example.com
Comment: DIY                       
You selected this USER-ID:
    "Cloud Security (DIY) <cloudsecurityteam@example.com>"

Change (N)ame, (C)omment, (E)mail or (O)kay/(Q)uit? o
```

The system will generate the private key. Once that has completed, you'll get a message similar to the following:

```
Note that this key cannot be used for encryption.  You may want to use
the command "--edit-key" to generate a subkey for this purpose.
pub   rsa4096/0x2A468DCA15B582C7 2018-08-17 [SC]
      Key fingerprint = 2032 24C4 5F50 3F4E 4D2F  534D 2A46 8DCA 15B5 82C7
uid                   Cloud Security (DIY) <cloudsecurityteam@example.com>
```

You can then export that key for use in Voucher, by running:

```
$ gpg -a --export-secret-key 0x2A468DCA15B582C7 > diy.gpg
```

If you look in `diy.gpg`, you will see something similar to the following:

```
-----BEGIN PGP PRIVATE KEY BLOCK-----

lQcYBFt23t8BEADuZqi....
```

This key will need to be put into the ejson file, so you will need to replace all of the newlines with "\n".

Our example key from before would then look like this.

```
-----BEGIN PGP PRIVATE KEY BLOCK-----\nlQcYBFt23t8BEADuZqi....
```

Next, create a new value in the `openpgpkeys` block in your `ejson` file. Make sure the key name is the same as it's name in the source code (eg, for the "DIY" test, use "diy"):

```json
{
    "_public_key": "<public key>",
    "openpgpkeys": {
        "diy": "-----BEGIN PGP PRIVATE KEY BLOCK-----\nlQcYBFt23t8BEADuZqi...."
    },
    "clair": {}
}
```

and call `ejson encrypt` to encrypt it:

```shell
$ ejson encrypt secrets.ejson
```
