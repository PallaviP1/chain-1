package ca

import "chain/crypto/ed25519/ecmath"

type AssetRangeProof struct {
	commitments []*AssetCommitment
	signature   *RingSignature
	id          *AssetID // nil means "confidential"
}

// CreateAssetRangeProof creates a confidential asset range proof. The
// caller can decorate the result with an asset ID to make it
// non-confidential.
func CreateAssetRangeProof(msg []byte, ac []*AssetCommitment, acPrime *AssetCommitment, j uint64, c, cPrime ecmath.Scalar) *AssetRangeProof {
	P := arpPubkeys(ac, acPrime)
	var p ecmath.Scalar
	p.Sub(&cPrime, &c)
	rs := CreateRingSignature(msg, []ecmath.Point{G, J}, P, j, p)
	return &AssetRangeProof{
		commitments: ac,
		signature:   rs,
	}
}

func (arp *AssetRangeProof) Validate(msg []byte, acPrime *AssetCommitment) bool {
	// xxx pending: whether/how to hash msg before calling
	// arp.signature.Validate, which also hashes
	P := arpPubkeys(arp.commitments, acPrime)
	if !arp.signature.Validate(msg, []ecmath.Point{G, J}, P) {
		return false
	}
	if arp.id != nil {
		// xxx
	}
	return true
}

func arpMsgHash(msg []byte, ac []*AssetCommitment, acPrime *AssetCommitment) [32]byte {
	hasher := hasher256([]byte("ARP"), acPrime.Bytes())
	for _, aci := range ac {
		hasher.Write(aci.Bytes())
	}
	hasher.Write(msg)
	var result [32]byte
	hasher.Read(result[:])
	return result
}

func arpPubkeys(ac []*AssetCommitment, acPrime *AssetCommitment) [][]ecmath.Point {
	n := len(ac)
	result := make([][]ecmath.Point, n)
	for i := 0; i < n; i++ {
		result[i] = make([]ecmath.Point, 2)
		result[i][0].Sub(&acPrime.Point1, &ac[i].Point1)
		result[i][1].Sub(&acPrime.Point2, &ac[i].Point2)
	}
	return result
}
