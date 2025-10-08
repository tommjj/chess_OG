// This file contains precomputed attack tables for sliding pieces (rooks and bishops).
//
// References:
//  - https://www.chessprogramming.org/Bitboards#Attack_Tables
//  - https://www.chessprogramming.org/Magic_Bitboards
//  - https://www.chessprogramming.org/Sliding_Piece_Attacks
//  - implementation:  https://github.com/maksimKorzh/bbc
//
// The attack tables are generated using magic bitboards, which provide a fast way to compute sliding piece attacks.
// The code initializes various masks and attack tables for pawns, knights, kings, rooks, and bishops.
// It also includes functions to compute attacks based on the current board occupancy.

package chess_core

import "fmt"

func init() {
	initAttackTables()
}

func initAttackTables() {
	initLeaperAttacks()
	initSliderMasks()
	initSliderAttacks()
	initBetweenLookup()
}

// Attack tables

// file masks
const (
	// not file A board
	//  8  0 1 1 1 1 1 1 1
	//  7  0 1 1 1 1 1 1 1
	//  6  0 1 1 1 1 1 1 1
	//  5  0 1 1 1 1 1 1 1
	//  4  0 1 1 1 1 1 1 1
	//  3  0 1 1 1 1 1 1 1
	//  2  0 1 1 1 1 1 1 1
	//  1  0 1 1 1 1 1 1 1
	//     a b c d e f g h
	notAFile BitBoard = 0xfefefefefefefefe
	// not file H board
	//  8  1 1 1 1 1 1 1 0
	//  7  1 1 1 1 1 1 1 0
	//  6  1 1 1 1 1 1 1 0
	//  5  1 1 1 1 1 1 1 0
	//  4  1 1 1 1 1 1 1 0
	//  3  1 1 1 1 1 1 1 0
	//  2  1 1 1 1 1 1 1 0
	//  1  1 1 1 1 1 1 1 0
	//     a b c d e f g h
	notHFile BitBoard = 0x7f7f7f7f7f7f7f7f
	// not file HG board
	//  8  1 1 1 1 1 1 0 0
	//  7  1 1 1 1 1 1 0 0
	//  6  1 1 1 1 1 1 0 0
	//  5  1 1 1 1 1 1 0 0
	//  4  1 1 1 1 1 1 0 0
	//  3  1 1 1 1 1 1 0 0
	//  2  1 1 1 1 1 1 0 0
	//  1  1 1 1 1 1 1 0 0
	//     a b c d e f g h
	notGHFile BitBoard = 0x3f3f3f3f3f3f3f3f
	// not file AB board
	//  8  0 0 1 1 1 1 1 1
	//  7  0 0 1 1 1 1 1 1
	//  6  0 0 1 1 1 1 1 1
	//  5  0 0 1 1 1 1 1 1
	//  4  0 0 1 1 1 1 1 1
	//  3  0 0 1 1 1 1 1 1
	//  2  0 0 1 1 1 1 1 1
	//  1  0 0 1 1 1 1 1 1
	//     a b c d e f g h
	notABFile BitBoard = 0xfcfcfcfcfcfcfcfc
)

// bishop relevant occupancy bit count for every square on board
var bishopRelevantBits = [64]int{
	6, 5, 5, 5, 5, 5, 5, 6,
	5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 7, 7, 7, 7, 5, 5,
	5, 5, 7, 9, 9, 7, 5, 5,
	5, 5, 7, 9, 9, 7, 5, 5,
	5, 5, 7, 7, 7, 7, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5,
	6, 5, 5, 5, 5, 5, 5, 6,
}

// rook relevant occupancy bit count for every square on board
var rookRelevantBits = [64]int{
	12, 11, 11, 11, 11, 11, 11, 12,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	12, 11, 11, 11, 11, 11, 11, 12,
}

var rookMagicNumbers = [64]BitBoard{
	0x8a80104000800020, 0x140002000100040, 0x2801880a0017001, 0x100081001000420,
	0x200020010080420, 0x3001c0002010008, 0x8480008002000100, 0x2080088004402900,
	0x800098204000, 0x2024401000200040, 0x100802000801000, 0x120800800801000,
	0x208808088000400, 0x2802200800400, 0x2200800100020080, 0x801000060821100,
	0x80044006422000, 0x100808020004000, 0x12108a0010204200, 0x140848010000802,
	0x481828014002800, 0x8094004002004100, 0x4010040010010802, 0x20008806104,
	0x100400080208000, 0x2040002120081000, 0x21200680100081, 0x20100080080080,
	0x2000a00200410, 0x20080800400, 0x80088400100102, 0x80004600042881,
	0x4040008040800020, 0x440003000200801, 0x4200011004500, 0x188020010100100,
	0x14800401802800, 0x2080040080800200, 0x124080204001001, 0x200046502000484,
	0x480400080088020, 0x1000422010034000, 0x30200100110040, 0x100021010009,
	0x2002080100110004, 0x202008004008002, 0x20020004010100, 0x2048440040820001,
	0x101002200408200, 0x40802000401080, 0x4008142004410100, 0x2060820c0120200,
	0x1001004080100, 0x20c020080040080, 0x2935610830022400, 0x44440041009200,
	0x280001040802101, 0x2100190040002085, 0x80c0084100102001, 0x4024081001000421,
	0x20030a0244872, 0x12001008414402, 0x2006104900a0804, 0x1004081002402,
}

var bishopMagicNumbers = [64]BitBoard{
	0x40040844404084, 0x2004208a004208, 0x10190041080202, 0x108060845042010,
	0x581104180800210, 0x2112080446200010, 0x1080820820060210, 0x3c0808410220200,
	0x4050404440404, 0x21001420088, 0x24d0080801082102, 0x1020a0a020400,
	0x40308200402, 0x4011002100800, 0x401484104104005, 0x801010402020200,
	0x400210c3880100, 0x404022024108200, 0x810018200204102, 0x4002801a02003,
	0x85040820080400, 0x810102c808880400, 0xe900410884800, 0x8002020480840102,
	0x220200865090201, 0x2010100a02021202, 0x152048408022401, 0x20080002081110,
	0x4001001021004000, 0x800040400a011002, 0xe4004081011002, 0x1c004001012080,
	0x8004200962a00220, 0x8422100208500202, 0x2000402200300c08, 0x8646020080080080,
	0x80020a0200100808, 0x2010004880111000, 0x623000a080011400, 0x42008c0340209202,
	0x209188240001000, 0x400408a884001800, 0x110400a6080400, 0x1840060a44020800,
	0x90080104000041, 0x201011000808101, 0x1a2208080504f080, 0x801202060021121,
	0x500861011240000, 0x180806108200800, 0x4000020e01040044, 0x300000261044000a,
	0x802241102020002, 0x20906061210001, 0x5a84841004010310, 0x4010801011c04,
	0xa010109502200, 0x4a02012000, 0x500201010098b028, 0x8040002811040900,
	0x28000010020204, 0x6000020202d0240, 0x8918844842082200, 0x4010011029020020,
}

var pawnAttacks [2][64]BitBoard

var knightAttacks [64]BitBoard

var kingAttacks [64]BitBoard

var rookMasks [64]BitBoard

var bishopMasks [64]BitBoard

var rookAttackTable [64][4096]BitBoard

var bishopAttackTable [64][512]BitBoard

var betweenLookup [64][64]BitBoard

// generate pawn attacks for a given square and side
func maskPawnAttacks(side Color, sq Square) BitBoard {
	var attacks BitBoard = 0
	pawn := BitBoard(1) << sq

	if side == Black {
		if (pawn>>7)&notHFile != 0 {
			attacks |= pawn >> 7
		}
		if (pawn>>9)&notAFile != 0 {
			attacks |= pawn >> 9
		}
	} else {
		if (pawn<<7)&notAFile != 0 {
			attacks |= pawn << 7
		}
		if (pawn<<9)&notHFile != 0 {
			attacks |= pawn << 9
		}
	}
	return attacks
}

func maskKnightAttacks(sq Square) BitBoard {
	var attacks BitBoard = 0
	knight := BitBoard(1) << sq

	if (knight>>17)&notHFile != 0 {
		attacks |= knight >> 17
	}
	if (knight>>15)&notAFile != 0 {
		attacks |= knight >> 15
	}
	if (knight>>10)&notGHFile != 0 {
		attacks |= knight >> 10
	}
	if (knight>>6)&notABFile != 0 {
		attacks |= knight >> 6
	}
	if (knight<<17)&notAFile != 0 {
		attacks |= knight << 17
	}
	if (knight<<15)&notHFile != 0 {
		attacks |= knight << 15
	}
	if (knight<<10)&notABFile != 0 {
		attacks |= knight << 10
	}
	if (knight<<6)&notGHFile != 0 {
		attacks |= knight << 6
	}
	return attacks
}

func maskKingAttacks(sq Square) BitBoard {
	var attacks BitBoard = 0
	king := BitBoard(1) << sq

	attacks |= king >> 8
	attacks |= king << 8
	if (king>>1)&notHFile != 0 {
		attacks |= king >> 1
	}
	if (king<<1)&notAFile != 0 {
		attacks |= king << 1
	}
	if (king>>9)&notHFile != 0 {
		attacks |= king >> 9
	}
	if (king>>7)&notAFile != 0 {
		attacks |= king >> 7
	}
	if (king<<9)&notAFile != 0 {
		attacks |= king << 9
	}
	if (king<<7)&notHFile != 0 {
		attacks |= king << 7
	}
	return attacks
}

func maskRookOccupancy(sq Square) BitBoard {
	var attacks BitBoard = 0

	rank := sq / 8
	file := sq % 8

	// north
	for r := rank + 1; r <= 6; r++ {
		attacks |= (BitBoard(1) << (file + r*8))
	}
	// south
	for r := rank - 1; r >= 1; r-- {
		attacks |= (BitBoard(1) << (file + r*8))
	}
	// east
	for f := file + 1; f <= 6; f++ {
		attacks |= (BitBoard(1) << (f + rank*8))
	}
	// west
	for f := file - 1; f >= 1; f-- {
		attacks |= (BitBoard(1) << (f + rank*8))
	}
	return attacks
}

func maskBishopOccupancy(sq Square) BitBoard {
	var attacks BitBoard = 0

	rank := sq / 8
	file := sq % 8

	// north east
	for r, f := rank+1, file+1; r <= 6 && f <= 6; r, f = r+1, f+1 {
		attacks |= (BitBoard(1) << (f + r*8))
	}
	// north west
	for r, f := rank+1, file-1; r <= 6 && f >= 1; r, f = r+1, f-1 {
		attacks |= (BitBoard(1) << (f + r*8))
	}
	// south east
	for r, f := rank-1, file+1; r >= 1 && f <= 6; r, f = r-1, f+1 {
		attacks |= (BitBoard(1) << (f + r*8))
	}
	// south west
	for r, f := rank-1, file-1; r >= 1 && f >= 1; r, f = r-1, f-1 {
		attacks |= (BitBoard(1) << (f + r*8))
	}
	return attacks
}

func initLeaperAttacks() {
	for sq := range Square(64) {
		pawnAttacks[White][sq] = maskPawnAttacks(White, sq)
		pawnAttacks[Black][sq] = maskPawnAttacks(Black, sq)
		kingAttacks[sq] = maskKingAttacks(sq)
		knightAttacks[sq] = maskKnightAttacks(sq)
	}
}

func initSliderMasks() {
	for sq := range Square(64) {
		rookMasks[sq] = maskRookOccupancy(sq)
		bishopMasks[sq] = maskBishopOccupancy(sq)
	}
}

// setOccupancy generates a blocker board based on the index and the attack mask
func setOccupancy(index int, bitsInMask int, attackMask BitBoard) BitBoard {
	var occupancy BitBoard = 0
	for range bitsInMask {
		// 1. Tìm vị trí (square) của bit cản đường thấp nhất
		square := attackMask.LeastSignificantBit()

		// 2. Xóa bit đó khỏi mặt nạ để quét bit tiếp theo trong vòng lặp sau
		attackMask &= attackMask - 1
		// 3. Nếu bit thấp nhất của index hiện tại là 1
		if (index & 1) != 0 {
			// Thì đặt BitBoard cản đường ở vị trí square đó
			occupancy |= (BitBoard(1) << square)
		}
		// 4. Dịch index sang phải để kiểm tra bit tiếp theo
		index >>= 1
	}
	return occupancy
}

func computeRookAttacks(sq Square, block BitBoard) BitBoard {
	var attacks BitBoard = 0

	rank := sq / 8
	file := sq % 8

	// north
	for r := rank + 1; r <= 7; r++ {
		attacks |= (BitBoard(1) << (file + r*8))
		if (BitBoard(1) << (file + r*8) & block) != 0 {
			break
		}
	}
	// south
	for r := rank - 1; r >= 0; r-- {
		attacks |= (BitBoard(1) << (file + r*8))
		if (BitBoard(1) << (file + r*8) & block) != 0 {
			break
		}
	}
	// east
	for f := file + 1; f <= 7; f++ {
		attacks |= (BitBoard(1) << (f + rank*8))
		if (BitBoard(1) << (f + rank*8) & block) != 0 {
			break
		}
	}
	// west
	for f := file - 1; f >= 0; f-- {
		attacks |= (BitBoard(1) << (f + rank*8))
		if (BitBoard(1) << (f + rank*8) & block) != 0 {
			break
		}
	}
	return attacks
}

func computeBishopAttacks(sq Square, block BitBoard) BitBoard {
	var attacks BitBoard = 0

	rank := sq / 8
	file := sq % 8

	// north east
	for r, f := rank+1, file+1; r <= 7 && f <= 7; r, f = r+1, f+1 {
		attacks |= (BitBoard(1) << (f + r*8))
		if (BitBoard(1) << (f + r*8) & block) != 0 {
			break
		}
	}
	// north west
	for r, f := rank+1, file-1; r <= 7 && f >= 0; r, f = r+1, f-1 {
		attacks |= (BitBoard(1) << (f + r*8))
		if (BitBoard(1) << (f + r*8) & block) != 0 {
			break
		}
	}
	// south east
	for r, f := rank-1, file+1; r >= 0 && f <= 7; r, f = r-1, f+1 {
		attacks |= (BitBoard(1) << (f + r*8))
		if (BitBoard(1) << (f + r*8) & block) != 0 {
			break
		}
	}
	// south west
	for r, f := rank-1, file-1; r >= 0 && f >= 0; r, f = r-1, f-1 {
		attacks |= (BitBoard(1) << (f + r*8))
		if (BitBoard(1) << (f + r*8) & block) != 0 {
			break
		}
	}
	return attacks
}

func initSliderAttacks() {
	var occupancyIndices [64]int
	for sq := range 64 {
		occupancyIndices[sq] = 1 << rookRelevantBits[sq]
	}
	var occupancy BitBoard
	for sq := range Square(64) {
		// rook attacks
		for index := 0; index < occupancyIndices[sq]; index++ {
			occupancy = setOccupancy(index, rookRelevantBits[sq], rookMasks[sq])
			magicIndex := (occupancy * rookMagicNumbers[sq]) >> (64 - rookRelevantBits[sq])
			rookAttackTable[sq][magicIndex] = computeRookAttacks(sq, occupancy)
		}
		// bishop attacks
		for index := 0; index < (1 << bishopRelevantBits[sq]); index++ {
			occupancy = setOccupancy(index, bishopRelevantBits[sq], bishopMasks[sq])
			magicIndex := (occupancy * bishopMagicNumbers[sq]) >> (64 - bishopRelevantBits[sq])
			bishopAttackTable[sq][magicIndex] = computeBishopAttacks(sq, occupancy)
		}
	}
}

func computeBetweenBits(from, to Square) BitBoard {
	if from == to {
		return 0
	}

	dir := directionBetween(from, to)
	if dir == 0 {
		return 0 // không cùng hàng, cột, hoặc chéo
	}

	var between BitBoard
	sq := Square(int(from) + dir)
	for sq != to {
		between |= sq.ToBB()
		sq += Square(dir)
	}
	return between
}

// Xác định hướng giữa 2 ô (nếu thẳng hoặc chéo)
func directionBetween(from, to Square) int {
	df := int(to%8) - int(from%8)
	dr := int(to/8) - int(from/8)

	if df == 0 { // dọc
		if dr > 0 {
			return 8
		}
		return -8
	}
	if dr == 0 { // ngang
		if df > 0 {
			return 1
		}
		return -1
	}
	if df == dr { // chéo chính
		if df > 0 {
			return 9
		}
		return -9
	}
	if df == -dr { // chéo phụ
		if df > 0 {
			return -7
		}
		return 7
	}
	return 0
}

func initBetweenLookup() {
	for from := range 64 {
		for to := range 64 {
			betweenLookup[from][to] = computeBetweenBits(Square(from), Square(to))
		}
	}
}

func magicIndexRook(sq Square, occupancy BitBoard) int {
	occ := occupancy & rookMasks[sq]
	return int((occ * rookMagicNumbers[sq]) >> (64 - rookRelevantBits[sq]))
}

func magicIndexBishop(sq Square, occupancy BitBoard) int {
	occ := occupancy & bishopMasks[sq]
	return int((occ * bishopMagicNumbers[sq]) >> (64 - bishopRelevantBits[sq]))
}

// RookAttacks returns the attack bitboard for a rook on square sq given the current occupancy
func RookAttacks(sq Square, occupancy BitBoard) BitBoard {
	magicIndex := magicIndexRook(sq, occupancy)
	return rookAttackTable[sq][magicIndex]
}

// BishopAttacks returns the attack bitboard for a bishop on square sq given the current occupancy
func BishopAttacks(sq Square, occupancy BitBoard) BitBoard {
	magicIndex := magicIndexBishop(sq, occupancy)
	return bishopAttackTable[sq][magicIndex]
}

// QueenAttacks returns the attack bitboard for a queen on square sq given the current occupancy
func QueenAttacks(sq Square, occupancy BitBoard) BitBoard {
	return RookAttacks(sq, occupancy) | BishopAttacks(sq, occupancy)
}

// Leaper attacks
func KingAttacks(sq Square) BitBoard {
	return kingAttacks[sq]
}

// KnightAttacks returns the attack bitboard for a knight on square sq
func KnightAttacks(sq Square) BitBoard {
	return knightAttacks[sq]
}

// PawnAttacks returns the attack bitboard for a pawn on square sq given the side
func PawnAttacks(sq Square, side Color) BitBoard {
	return pawnAttacks[side][sq]
}

// IsKingAttacked checks if the king of the given color is under attack
func IsKingAttacked(side Color, bb *BitBoards) bool {
	kingSquare := Square(bb.PiecesByType(side, King).LeastSignificantBit())
	if kingSquare == -1 {
		// no king found
		panic(fmt.Sprintf("IsKingAttacked: no king found for side %v", side))
	}

	return IsAttacked(kingSquare, side, bb)
}

// IsAttacked checks if the square sq is attacked by any piece of the given side
func IsAttacked(sq Square, defendingColor Color, bb *BitBoards) bool {
	oop := defendingColor.Opposite()

	// check pawn attacks
	if (PawnAttacks(sq, defendingColor) & bb.PiecesByType(oop, Pawn)) != 0 {
		return true
	}

	// check knight attacks
	if (KnightAttacks(sq) & bb.PiecesByType(oop, Knight)) != 0 {
		return true
	}

	// check king attacks
	if (KingAttacks(sq) & bb.PiecesByType(oop, King)) != 0 {
		return true
	}

	// check bishop/queen attacks
	if (BishopAttacks(sq, bb.AllPieces) & (bb.PiecesByType(oop, Bishop) | bb.PiecesByType(oop, Queen))) != 0 {
		return true
	}

	// check rook/queen attacks
	if (RookAttacks(sq, bb.AllPieces) & (bb.PiecesByType(oop, Rook) | bb.PiecesByType(oop, Queen))) != 0 {
		return true
	}

	return false
}

func AttackersTo(sq Square, defendingColor Color, bb *BitBoards) BitBoard {
	var attackers BitBoard
	oop := defendingColor.Opposite()

	attackers |= PawnAttacks(sq, defendingColor) & bb.PiecesByType(oop, Pawn)
	attackers |= KnightAttacks(sq) & bb.PiecesByType(oop, King)
	attackers |= BishopAttacks(sq, bb.AllPieces) & (bb.PiecesByType(oop, Bishop) | bb.PiecesByType(oop, Queen))
	attackers |= RookAttacks(sq, bb.AllPieces) & (bb.PiecesByType(oop, Rook) | bb.PiecesByType(oop, Queen))
	attackers |= KingAttacks(sq) & bb.PiecesByType(oop, King)

	return attackers
}

// BetweenBits
func BetweenBits(from, to Square) BitBoard {
	return betweenLookup[from][to]
}
