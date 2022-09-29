package dasvid

import (
  "fmt"
  "go.dedis.ch/kyber/v3"
  "go.dedis.ch/kyber/v3/group/edwards25519"
)

var curve = edwards25519.NewBlakeSHA256Ed25519()
var sha256 = curve.Hash()

type Signature struct {
    R kyber.Point
    S kyber.Scalar
}

func Hash(s string) kyber.Scalar {
    sha256.Reset()
    sha256.Write([]byte(s))

    return curve.Scalar().SetBytes(sha256.Sum(nil))
}

// m: Message
// x: Private key
func Sign(m string, z kyber.Scalar) Signature {
    // Get the base of the curve.
    g := curve.Point().Base()

    // Pick a random k from allowed set.
    k := curve.Scalar().Pick(curve.RandomStream())

    // r = k * G (a.k.a the same operation as r = g^k)
    r := curve.Point().Mul(k, g)

    // h := Hash(publicKey + m + r.String())
    publicKey := curve.Point().Mul(z, g)
    h := Hash(publicKey.String() + m + r.String())
    
    // s = k - e * x
    s := curve.Scalar().Sub(k, curve.Scalar().Mul(h, z))

    return Signature{R: r, S: s}
}

// m: Message
// s: Signature
// y: Public key
func Verify(m string, S Signature, y kyber.Point) bool {
    // Create a generator.
    g := curve.Point().Base()

    // h = Hash(pubkey || m || r)
    h := Hash(y.String() + m + S.R.String())

    // Attempt to reconstruct 's * G' with a provided signature; s * G = r - h * y
    sGv := curve.Point().Sub(S.R, curve.Point().Mul(h, y))

    // Construct the actual 's * G'
    sG := curve.Point().Mul(S.S, g)

    // Equality check; ensure signature and public key outputs to s * G.
    return sG.Equal(sGv)
}

func (S Signature) String() string {
    return fmt.Sprintf("(r=%s, s=%s)", S.R, S.S)
}

// origpubkey = first public key
// setSigR = array containing all Sig.R
// setH = array containing all Hashes
// s1s = last signature.S
func Verifygg(origpubkey kyber.Point, setSigR []kyber.Point, setH []kyber.Scalar, s1s kyber.Scalar) bool {
    // Verify n concatenated signatures using galindo-garcia

    // Important to note that as new assertions are added in the beginning of the token, the content of arrays is in reverse order.
    // e.g. setSigR[0] = last appended signature. Thats why 'i' starts from len(setSigR)-1
    if (len(setSigR)) != len(setH) {
        fmt.Println("Incorrect parameters!")
        return false
    }

    var i = len(setSigR)-1
    var y kyber.Point

    // Create a generator.
    g := curve.Point().Base()

    // calculate all y's from first to last-1 parts
	for (i > 0) {
        if (i == len(setSigR)-1) {
            y = origpubkey
        } else {
            y = curve.Point().Sub(setSigR[i+1], curve.Point().Mul(setH[i+1], y))
        }
        i--
    }

    // calculate last y
    y = curve.Point().Sub(setSigR[i+1], curve.Point().Mul(setH[i+1], y))

    // check if g ^ lastsig.S = lastsig.R - y ^ lastHash
    leftside    := curve.Point().Mul(s1s, g)
    rightside   := curve.Point().Sub(setSigR[i], curve.Point().Mul(setH[i], y))

    return leftside.Equal(rightside)
}

// Given ID, return a keypair 
func IDKeyPair(id string) (kyber.Scalar, kyber.Point){


    privateKey	:= Hash(id)
    publicKey 	:= curve.Point().Mul(privateKey, curve.Point().Base())

    return privateKey, publicKey
}

// Return a new random key pair
func RandomKeyPair() (kyber.Scalar, kyber.Point){

    privateKey	:= curve.Scalar().Pick(curve.RandomStream())
    publicKey 	:= curve.Point().Mul(privateKey, curve.Point().Base())

    return privateKey, publicKey
}