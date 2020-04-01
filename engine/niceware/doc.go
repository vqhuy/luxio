package niceware

// This package implements a variant of the niceware.js package
// (https://github.com/diracdeltas/niceware)
// for generating random-yet-memorable passwords.
// This implementation converts a pseudo-random string (of size 32 bytes)
// to a 4-word phrase (64 bits of randomness).
// The byte array is splitted in to 4 blocks, each of 8 bytes.
// For each block, the first 4 bytes are XOR-ed, and the last 4 bytes
// are XOR-ed together, providing 2 random bytes.
// Then we apply niceware's `BytesToPassphrase` to these 2 random bytes,
// and obtain 4 words in total.
