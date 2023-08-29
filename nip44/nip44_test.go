package nip44

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/studiokaiji/go-nostr"
)

func TestValidSec(t *testing.T) {
	vectors := []struct {
		Sec1       string `json:"sec1"`
		Sec2       string `json:"sec2"`
		Shared     string `json:"shared"`
		Nonce      string `json:"nonce"`
		Plaintext  string `json:"plaintext"`
		Ciphertext string `json:"ciphertext"`
		Note       string `json:"note"`
	}{}
	if err := json.Unmarshal([]byte(VECTORS_VALID_SEC), &vectors); err != nil {
		t.Errorf("failed to read test vectors json")
	}

	for _, vec := range vectors {
		pub1, _ := nostr.GetPublicKey(vec.Sec1)
		pub2, _ := nostr.GetPublicKey(vec.Sec2)
		ss1, _ := ComputeSharedSecret(pub1, vec.Sec2)
		ss2, _ := ComputeSharedSecret(pub2, vec.Sec1)
		if hex.EncodeToString(ss1) != hex.EncodeToString(ss2) || hex.EncodeToString(ss1) != vec.Shared {
			t.Errorf("invalid shared secret: %x != %x or %x != %s", ss1, ss2, ss1, vec.Shared)
		}

		nonce, _ := hex.DecodeString(vec.Nonce)
		ciphertext, _ := encryptWithNonce(vec.Plaintext, ss1, nonce)
		if ciphertext != vec.Ciphertext {
			t.Errorf("invalid ciphertext: %s != %s", ciphertext, vec.Ciphertext)
		}

		plaintext, _ := Decrypt(ciphertext, ss1)
		if plaintext != vec.Plaintext {
			t.Errorf("invalid plaintext")
		}
	}
}

func TestValidPub(t *testing.T) {
	vectors := []struct {
		Sec1       string `json:"sec1"`
		Pub2       string `json:"pub2"`
		Shared     string `json:"shared"`
		Nonce      string `json:"nonce"`
		Plaintext  string `json:"plaintext"`
		Ciphertext string `json:"ciphertext"`
		Note       string `json:"note"`
	}{}

	for _, vec := range vectors {
		ss, _ := ComputeSharedSecret(vec.Pub2, vec.Sec1)
		if hex.EncodeToString(ss) != vec.Shared {
			t.Errorf("invalid shared secret")
		}

		nonce, _ := hex.DecodeString(vec.Nonce)
		ciphertext, _ := encryptWithNonce(vec.Plaintext, ss, nonce)
		if ciphertext != vec.Ciphertext {
			t.Errorf("invalid ciphertext")
		}

		plaintext, _ := Decrypt(ciphertext, ss)
		if plaintext != vec.Plaintext {
			t.Errorf("invalid plaintext")
		}
	}
}

func TestInvalid(t *testing.T) {
	vectors := []struct {
		Sec1      string `json:"sec1"`
		Pub2      string `json:"pub2"`
		Plaintext string `json:"plaintext"`
		Note      string `json:"note"`
	}{}

	for _, vec := range vectors {
		_, err := ComputeSharedSecret(vec.Pub2, vec.Sec1)
		if err == nil {
			t.Errorf("should have failed, but didn't: %s", vec.Note)
		}
	}
}

const (
	VECTORS_VALID_SEC = `
[
  {
    "sec1": "0000000000000000000000000000000000000000000000000000000000000001",
    "sec2": "0000000000000000000000000000000000000000000000000000000000000002",
    "shared": "0135da2f8acf7b9e3090939432e47684eb888ea38c2173054d4eedffdf152ca5",
    "nonce": "121f9d60726777642fd82286791ab4d7461c9502ebcbb6e6",
    "plaintext": "a",
    "ciphertext": "ARIfnWByZ3dkL9gihnkatNdGHJUC68u25qM=",
    "note": "sk1 = 1, sk2 = random, 0x02"
  },
  {
    "sec1": "0000000000000000000000000000000000000000000000000000000000000002",
    "sec2": "0000000000000000000000000000000000000000000000000000000000000001",
    "shared": "0135da2f8acf7b9e3090939432e47684eb888ea38c2173054d4eedffdf152ca5",
    "plaintext": "a",
    "ciphertext": "AeCt7jJ8L+WBOTiCSfeXEGXB/C/wgsrSRek=",
    "nonce": "e0adee327c2fe58139388249f7971065c1fc2ff082cad245",
    "note": "sk1 = 1, sk2 = random, 0x02"
  }
]
`
	VECTORS_VALID_PUB = `
[
  {
    "sec1": "fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364139",
    "pub2": "0000000000000000000000000000000000000000000000000000000000000002",
    "shared": "a6d6a2f7011cdd1aeef325948f48c6efa40f0ec723ae7f5ac7e3889c43481500",
    "nonce": "f481750e13dfa90b722b7cce0db39d80b0db2e895cc3001a",
    "plaintext": "a",
    "ciphertext": "AfSBdQ4T36kLcit8zg2znYCw2y6JXMMAGjM=",
    "note": "sec1 = n-2, pub2: random, 0x02"
  },
  {
    "sec1": "0000000000000000000000000000000000000000000000000000000000000002",
    "pub2": "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdeb",
    "shared": "4908464f77dd74e11a9b4e4a3bc2467445bd794e8abcbfafb65a6874f9e25a8f",
    "nonce": "45c484ba2c0397853183adba6922156e09a2ad4e3e6914f2",
    "plaintext": "A Peer-to-Peer Electronic Cash System",
    "ciphertext": "AUXEhLosA5eFMYOtumkiFW4Joq1OPmkU8k/25+3+VDFvOU39qkUDl1aiy8Q+0ozTwbhD57VJoIYayYS++hE=",
    "note": "sec1 = 2, pub2: "
  },
  {
    "sec1": "0000000000000000000000000000000000000000000000000000000000000001",
    "pub2": "79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798",
    "shared": "132f39a98c31baaddba6525f5d43f2954472097fa15265f45130bfdb70e51def",
    "nonce": "d60de08405cf9bde508147e82224ac6af409c12b9e5492e1",
    "plaintext": "A purely peer-to-peer version of electronic cash would allow online payments to be sent directly from one party to another without going through a financial institution. Digital signatures provide part of the solution, but the main benefits are lost if a trusted third party is still required to prevent double-spending.",
    "ciphertext": "AdYN4IQFz5veUIFH6CIkrGr0CcErnlSS4VdvoQaP2DCB1dIFL72HSriG1aFABcTlu86hrsG0MdOO9rPdVXc3jptMMzqvIN6tJlHPC8GdwFD5Y8BT76xIIOTJR2W0IdrM7++WC/9harEJAdeWHDAC9zNJX81CpCz4fnV1FZ8GxGLC0nUF7NLeUiNYu5WFXQuO9uWMK0pC7tk3XVogk90X6rwq0MQG9ihT7e1elatDy2YGat+VgQlDrz8ZLRw/lvU+QqeXMQgjqn42sMTrimG6NdKfHJSVWkT6SKZYVsuTyU1Iu5Nk0twEV8d11/MPfsMx4i36arzTC9qxE6jftpOoG8f/jwPTSCEpHdZzrb/CHJcpc+zyOW9BZE2ZOmSxYHAE0ustC9zRNbMT3m6LqxIoHq8j+8Ysu+Cwqr4nUNLYq/Q31UMdDg1oamYS17mWIAS7uf2yF5uT5IlG",
    "note": "sec1 == pub2"
  }
]
`
	VECTORS_INVALID = `
[
  {
    "sec1": "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
    "pub2": "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
    "plaintext": "a",
    "note": "sec1 higher than curve.n"
  },
  {
    "sec1": "0000000000000000000000000000000000000000000000000000000000000000",
    "pub2": "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
    "plaintext": "a",
    "note": "sec1 is 0"
  },
  {
    "sec1": "fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364139",
    "pub2": "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
    "plaintext": "a",
    "note": "pub2 is invalid, no sqrt, all-ff"
  },
  {
    "sec1": "fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141",
    "pub2": "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
    "plaintext": "a",
    "note": "sec1 == curve.n"
  },
  {
    "sec1": "0000000000000000000000000000000000000000000000000000000000000002",
    "pub2": "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
    "plaintext": "a",
    "note": "pub2 is invalid, no sqrt"
  }
]
`
)
