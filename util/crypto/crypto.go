package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	ethsecp256k1 "github.com/ethereum/go-ethereum/crypto/secp256k1"
)

// Secp256k1 returns the standard secp256k1 curve implementation.
func Secp256k1() elliptic.Curve {
	return ethsecp256k1.S256()
}

// GenerateECKeyPair generates an ECDSA private key using the specified curve.
func GenerateECKeyPair(curve elliptic.Curve) (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(curve, rand.Reader)
}

// PrivateKeyToPEM converts an ecdsa.PrivateKey to PEM format.
// It uses go-ethereum/crypto for raw key bytes and then PEM encodes.
// OIDs used for secp256k1 PKCS#8 encoding.
var (
	oidPublicKeyECDSA      = asn1.ObjectIdentifier{1, 2, 840, 10045, 2, 1}
	oidNamedCurveSecp256k1 = asn1.ObjectIdentifier{1, 3, 132, 0, 10}
)

type pkcs8AlgorithmIdentifier struct {
	Algorithm  asn1.ObjectIdentifier
	Parameters asn1.RawValue `asn1:"optional"`
}

type pkcs8PrivateKey struct {
	Version    int
	Algo       pkcs8AlgorithmIdentifier
	PrivateKey []byte
}

type ecPrivateKey struct {
	Version       int
	PrivateKey    []byte
	NamedCurveOID asn1.ObjectIdentifier `asn1:"explicit,tag:0,optional"`
	PublicKey     asn1.BitString        `asn1:"explicit,tag:1,optional"`
}

func curveByteSize(curve elliptic.Curve) int {
	bits := curve.Params().BitSize
	return (bits + 7) / 8
}

func marshalECPrivateKey(key *ecdsa.PrivateKey) ([]byte, error) {
	if key == nil {
		return nil, errors.New("private key is nil")
	}

	size := curveByteSize(key.Curve)
	priv := make([]byte, size)
	dBytes := key.D.Bytes()
	copy(priv[size-len(dBytes):], dBytes)

	pubBytes := elliptic.Marshal(key.Curve, key.X, key.Y)

	ecKey := ecPrivateKey{
		Version:       1,
		PrivateKey:    priv,
		NamedCurveOID: oidNamedCurveSecp256k1,
		PublicKey: asn1.BitString{
			Bytes:     pubBytes,
			BitLength: len(pubBytes) * 8,
		},
	}

	return asn1.Marshal(ecKey)
}

// PrivateKeyToPEM converts an ecdsa.PrivateKey to PKCS#8 PEM format so that it matches
// the Python SDK's output.
func PrivateKeyToPEM(privateKey *ecdsa.PrivateKey) ([]byte, error) {
	if privateKey == nil {
		return nil, errors.New("private key is nil")
	}

	if privateKey.Curve != Secp256k1() {
		return nil, fmt.Errorf("unsupported curve for PKCS#8 export: %T", privateKey.Curve)
	}

	ecKey, err := marshalECPrivateKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal EC private key: %w", err)
	}

	params, err := asn1.Marshal(oidNamedCurveSecp256k1)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal curve oid: %w", err)
	}

	pkcs8Key := pkcs8PrivateKey{
		Version: 0,
		Algo: pkcs8AlgorithmIdentifier{
			Algorithm:  oidPublicKeyECDSA,
			Parameters: asn1.RawValue{FullBytes: params},
		},
		PrivateKey: ecKey,
	}

	der, err := asn1.Marshal(pkcs8Key)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal PKCS#8 key: %w", err)
	}

	return pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der}), nil
}

// PrivateKeyFromPEM parses a PEM-encoded private key.
// It uses go-ethereum/crypto to convert raw bytes to ecdsa.PrivateKey.
func parseECPrivateKeyDER(der []byte) (*ecdsa.PrivateKey, error) {
	var ecKey ecPrivateKey
	if _, err := asn1.Unmarshal(der, &ecKey); err != nil {
		return nil, fmt.Errorf("failed to parse EC private key: %w", err)
	}

	if len(ecKey.NamedCurveOID) > 0 && !ecKey.NamedCurveOID.Equal(oidNamedCurveSecp256k1) {
		return nil, fmt.Errorf("unexpected curve OID: %v", ecKey.NamedCurveOID)
	}

	curve := Secp256k1()
	size := curveByteSize(curve)
	if len(ecKey.PrivateKey) != size {
		return nil, fmt.Errorf("invalid private key length: got %d want %d", len(ecKey.PrivateKey), size)
	}

	d := new(big.Int).SetBytes(ecKey.PrivateKey)
	if d.Sign() <= 0 || d.Cmp(curve.Params().N) >= 0 {
		return nil, errors.New("invalid private key scalar")
	}

	x, y := curve.ScalarBaseMult(ecKey.PrivateKey)
	return &ecdsa.PrivateKey{PublicKey: ecdsa.PublicKey{Curve: curve, X: x, Y: y}, D: d}, nil
}

func parsePKCS8PrivateKey(der []byte) (*ecdsa.PrivateKey, error) {
	var pkcs8 pkcs8PrivateKey
	if _, err := asn1.Unmarshal(der, &pkcs8); err != nil {
		return nil, fmt.Errorf("failed to parse PKCS#8 structure: %w", err)
	}

	if !pkcs8.Algo.Algorithm.Equal(oidPublicKeyECDSA) {
		return nil, fmt.Errorf("unexpected algorithm OID: %v", pkcs8.Algo.Algorithm)
	}

	var curveOID asn1.ObjectIdentifier
	if len(pkcs8.Algo.Parameters.FullBytes) > 0 {
		if _, err := asn1.Unmarshal(pkcs8.Algo.Parameters.FullBytes, &curveOID); err != nil {
			return nil, fmt.Errorf("failed to parse curve parameters: %w", err)
		}
	}

	if len(curveOID) == 0 {
		curveOID = oidNamedCurveSecp256k1
	}

	if !curveOID.Equal(oidNamedCurveSecp256k1) {
		return nil, fmt.Errorf("unexpected curve parameters OID: %v", curveOID)
	}

	return parseECPrivateKeyDER(pkcs8.PrivateKey)
}

func PrivateKeyFromPEM(pemBytes []byte) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	switch block.Type {
	case "PRIVATE KEY":
		return parsePKCS8PrivateKey(block.Bytes)
	case "EC PRIVATE KEY":
		// Legacy support: previous versions stored raw 32-byte keys with this label.
		if len(block.Bytes) == 32 {
			privKey, err := ethcrypto.ToECDSA(block.Bytes)
			if err != nil {
				return nil, fmt.Errorf("failed to parse raw EC private key: %w", err)
			}
			return privKey, nil
		}
		return parseECPrivateKeyDER(block.Bytes)
	default:
		return nil, fmt.Errorf("unsupported PEM block type: %s", block.Type)
	}
}
func MakeSignature(privateKey *ecdsa.PrivateKey, content []byte) (string, error) {
	if privateKey == nil {
		return "", errors.New("private key is required")
	}

	digest := sha256.Sum256(content)
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, digest[:])
	if err != nil {
		return "", fmt.Errorf("signing payload: %w", err)
	}

	if privateKey.Curve == nil {
		return "", errors.New("elliptic curve is required")
	}
	params := privateKey.Curve.Params()
	size := (params.BitSize + 7) / 8
	rb := r.Bytes()
	sb := s.Bytes()
	if len(rb) > size || len(sb) > size {
		return "", fmt.Errorf("signature component larger than curve size")
	}

	sig := make([]byte, size*2)
	copy(sig[size-len(rb):size], rb)
	copy(sig[2*size-len(sb):], sb)

	return base64.RawURLEncoding.EncodeToString(sig), nil
}

func VerifySignature(content []byte, signHex string, pubkeyAddress string) bool {
	signBytes, err := base64.RawURLEncoding.DecodeString(signHex)
	if err != nil {
		return false
	}
	curve := Secp256k1()

	size := (curve.Params().BitSize + 7) / 8
	if len(signBytes) != size*2 {
		return false
	}

	r := new(big.Int).SetBytes(signBytes[:size])
	s := new(big.Int).SetBytes(signBytes[size:])

	digest := sha256.Sum256(content)
	sig := make([]byte, 65)
	rBytes := r.Bytes()
	sBytes := s.Bytes()
	copy(sig[32-len(rBytes):32], rBytes)
	copy(sig[64-len(sBytes):64], sBytes)
	var recoveredPubKey *ecdsa.PublicKey
	for recoveryID := byte(0); recoveryID < 2; recoveryID++ {
		sig[64] = recoveryID
		recoveredPubBytes, err := ethcrypto.Ecrecover(digest[:], sig)
		if err != nil {
			continue
		}
		recoveredPubKey, err = ethcrypto.UnmarshalPubkey(recoveredPubBytes)
		if err != nil {
			continue
		}
		pubkeyAddr := ethcrypto.PubkeyToAddress(*recoveredPubKey).Hex()
		// 比较恢复出的公钥是否与预期一致
		if pubkeyAddr == pubkeyAddress {
			return true
		}
	}
	return false
}
