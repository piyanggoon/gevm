package types

import (
	"math/big"

	"github.com/holiman/uint256"
)

var (
	U256Zero = uint256.Int{}
	U256One  = uint256.Int{1, 0, 0, 0}
	U256Max  = uint256.Int{
		0xffffffffffffffff,
		0xffffffffffffffff,
		0xffffffffffffffff,
		0xffffffffffffffff,
	}
	U256MaxPositiveI256 = uint256.Int{
		0xffffffffffffffff,
		0xffffffffffffffff,
		0xffffffffffffffff,
		0x7fffffffffffffff,
	}
	U256MinNegativeI256 = uint256.Int{
		0,
		0,
		0,
		0x8000000000000000,
	}
)

func U256From(v uint64) uint256.Int {
	var z uint256.Int
	z.SetUint64(v)
	return z
}

func U256FromLimbs(a, b, c, d uint64) uint256.Int {
	return uint256.Int{a, b, c, d}
}

func U256FromBytes32(b [32]byte) uint256.Int {
	var z uint256.Int
	z.SetBytes32(b[:])
	return z
}

func U256FromBytes(b []byte) uint256.Int {
	var z uint256.Int
	z.SetBytes(b)
	return z
}

func U256FromBig(v *big.Int) uint256.Int {
	var z uint256.Int
	z.SetFromBig(v)
	return z
}

func U256ToAddress(u *uint256.Int) Address {
	b := u.Bytes32()
	var addr Address
	copy(addr[:], b[12:32])
	return addr
}

func U256AsUsize(u *uint256.Int) uint64 {
	v, overflow := u.Uint64WithOverflow()
	if overflow {
		return ^uint64(0)
	}
	return v
}

func U256AsUsizeSaturated(u *uint256.Int) uint64 { return U256AsUsize(u) }
func U256LowU64(u *uint256.Int) uint64           { return u.Uint64() }

func U256ByteBE(u *uint256.Int, n uint) byte {
	if n >= 32 {
		return 0
	}
	var idx uint256.Int
	idx.SetUint64(uint64(n))
	z := *u
	z.Byte(&idx)
	return byte(z.Uint64())
}

func U256ByteLE(u *uint256.Int, n uint) byte {
	if n >= 32 {
		return 0
	}
	b := u.Bytes32()
	return b[31-n]
}

func U256LeadingZeros(u *uint256.Int) uint { return uint(256 - u.BitLen()) }

func U256Bit(u *uint256.Int, n uint) bool {
	if n >= 256 {
		return false
	}
	return u.ToBig().Bit(int(n)) == 1
}

func U256PutBytes32(u *uint256.Int, dst []byte) {
	b := u.Bytes32()
	copy(dst, b[:])
}

func AddTo(result, a, b *uint256.Int) { result.Add(a, b) }
func SubTo(result, a, b *uint256.Int) { result.Sub(a, b) }
func AndTo(result, a, b *uint256.Int) { result.And(a, b) }
func OrTo(result, a, b *uint256.Int)  { result.Or(a, b) }
func XorTo(result, a, b *uint256.Int) { result.Xor(a, b) }
func NotTo(result, a *uint256.Int)    { result.Not(a) }

func EqPtr(a, b *uint256.Int) bool      { return a.Eq(b) }
func LtPtr(a, b *uint256.Int) bool      { return a.Lt(b) }
func GtPtr(a, b *uint256.Int) bool      { return a.Gt(b) }
func IsZeroPtr(a *uint256.Int) bool     { return a.IsZero() }
func IsOnePtr(a *uint256.Int) bool      { return a.Eq(&U256One) }
func Cmp(a, b *uint256.Int) int         { return a.Cmp(b) }
func Add(a, b *uint256.Int) uint256.Int { var z uint256.Int; z.Add(a, b); return z }
func Sub(a, b *uint256.Int) uint256.Int { var z uint256.Int; z.Sub(a, b); return z }
func Mul(a, b *uint256.Int) uint256.Int { var z uint256.Int; z.Mul(a, b); return z }
func Div(a, b *uint256.Int) uint256.Int { var z uint256.Int; z.Div(a, b); return z }
func Mod(a, b *uint256.Int) uint256.Int { var z uint256.Int; z.Mod(a, b); return z }
func AddMod(a, b, m *uint256.Int) uint256.Int {
	var z uint256.Int
	z.AddMod(a, b, m)
	return z
}
func MulMod(a, b, m *uint256.Int) uint256.Int {
	var z uint256.Int
	z.MulMod(a, b, m)
	return z
}
func Exp(base, exponent *uint256.Int) uint256.Int {
	var z uint256.Int
	z.Exp(base, exponent)
	return z
}
func Neg(a *uint256.Int) uint256.Int         { var z uint256.Int; z.Neg(a); return z }
func Not(a *uint256.Int) uint256.Int         { var z uint256.Int; z.Not(a); return z }
func Shl(a *uint256.Int, n uint) uint256.Int { var z uint256.Int; z.Lsh(a, n); return z }
func Shr(a *uint256.Int, n uint) uint256.Int { var z uint256.Int; z.Rsh(a, n); return z }
func Sar(a *uint256.Int, n uint) uint256.Int { var z uint256.Int; z.SRsh(a, n); return z }

func OverflowingAdd(a, b *uint256.Int) (uint256.Int, bool) {
	var z uint256.Int
	_, overflow := z.AddOverflow(a, b)
	return z, overflow
}

func OverflowingSub(a, b *uint256.Int) (uint256.Int, bool) {
	var z uint256.Int
	_, underflow := z.SubOverflow(a, b)
	return z, underflow
}

func SignExtend(b uint256.Int, x uint256.Int) uint256.Int {
	var z uint256.Int
	z.ExtendSign(&x, &b)
	return z
}

type Sign int8

const (
	SignMinus Sign = -1
	SignZero  Sign = 0
	SignPlus  Sign = 1
)

func I256Sign(u uint256.Int) Sign {
	if u[3]>>63 != 0 {
		return SignMinus
	}
	if u.IsZero() {
		return SignZero
	}
	return SignPlus
}

func I256SignCompl(val *uint256.Int) Sign {
	if val[3]>>63 != 0 {
		val.Neg(val)
		return SignMinus
	}
	if val.IsZero() {
		return SignZero
	}
	return SignPlus
}

func U256RemoveSign(val *uint256.Int) {
	val[3] &= 0x7fffffffffffffff
}

func I256Cmp(a, b uint256.Int) int {
	if a.Slt(&b) {
		return -1
	}
	if a.Sgt(&b) {
		return 1
	}
	return 0
}

func I256Lt(a, b uint256.Int) bool { return a.Slt(&b) }
func I256Gt(a, b uint256.Int) bool { return a.Sgt(&b) }

func SDiv(a, b uint256.Int) uint256.Int {
	var z uint256.Int
	z.SDiv(&a, &b)
	return z
}

func I256ToU256(v int64) uint256.Int {
	if v >= 0 {
		return U256From(uint64(v))
	}
	return uint256.Int{uint64(v), ^uint64(0), ^uint64(0), ^uint64(0)}
}

func SMod(a, b uint256.Int) uint256.Int {
	var z uint256.Int
	z.SMod(&a, &b)
	return z
}

func U256ByteLen(u *uint256.Int) uint { return uint(u.ByteLen()) }
