package movegen

import (
	bd "perft/board"
)

type tables struct {
	magics  []bd.BitBoard
	shifts  []uint
	masks   []bd.BitBoard
	offsets []int
	attacks []bd.BitBoard
	dirs    []direction
}

var rookTables = tables{
	// Pre-computed using https://github.com/goutham/magic-bits library.
	magics: []bd.BitBoard{
		0x880081080c00020, 0x210020c000308100,
		0x80082001100280, 0x1001000a0050108,
		0x200041029600a00, 0x5100010008220400,
		0x8280120001000d80, 0x1880012100014080,
		0x3040800340008020, 0x400400050026003,
		0x21002000104902, 0x20900200a100100,
		0xd800802840080, 0x2808004000600,
		0x24001002110814, 0x2000800541000480,
		0x8000ee8002400080, 0x24c04010002005,
		0x822002401000c800, 0x2040808010000800,
		0x804080800c000802, 0x2a0080110402004,
		0x201044000810010a, 0x4080020004004483,
		0x4d84400180228000, 0x1406400880200880,
		0x801200402203, 0x1080080280100084,
		0x402140080080080, 0xa880c0080020080,
		0x342000200080405, 0x20004a8200050044,
		0x8280c00020800889, 0x8002201000400940,
		0x44a200101001542, 0x88090021005000,
		0x3008004200c00400, 0x284120080800400,
		0x4462106804000201, 0x1008240382000061,
		0x80400080208002, 0x20100040004020,
		0x4000802042020010, 0x40a002042120008,
		0x12a008820120004, 0x6000408020010,
		0x2008405020008, 0x80100c0040820003,
		0x2800100446100, 0xa0982002400080,
		0x9a0080010014040, 0x380c209200420a00,
		0xc04008108000580, 0xc002008004002280,
		0x2900842a000100, 0x40100008a004300,
		0x10211800020c3, 0xa08412050242,
		0x2001004010200489, 0xa00081000210045,
		0x4512002810204402, 0x8c22000401102802,
		0x485000082005401, 0x100208400ce,
	},
	shifts: []uint{
		12, 11, 11, 11, 11, 11, 11, 12,
		11, 10, 10, 10, 10, 10, 10, 11,
		11, 10, 10, 10, 10, 10, 10, 11,
		11, 10, 10, 10, 10, 10, 10, 11,
		11, 10, 10, 10, 10, 10, 10, 11,
		11, 10, 10, 10, 10, 10, 10, 11,
		11, 10, 10, 10, 10, 10, 10, 11,
		12, 11, 11, 11, 11, 11, 11, 12,
	},
	dirs: []direction{north, south, east, west},
}

var bishopTables = tables{
	// Pre-computed using https://github.com/goutham/magic-bits library.
	magics: []bd.BitBoard{
		0x820420460402a080, 0x20021200451400,
		0x10011200218000, 0x4040888100800,
		0x6211001000400, 0x401042240021400,
		0x884029888090060, 0x24202808080810,
		0x20242038024080, 0x80021081010102,
		0x100004090c030120, 0x210c0420814205,
		0x408311040061010, 0x4900011016100900,
		0x6841020d30461020, 0x220112088080800,
		0x8040000802080628, 0x4a48000408480040,
		0x2010000e00b20060, 0x1004020809409102,
		0x1011090400801, 0x2002000420842000,
		0xa01200443a090402, 0x1010082a4020221,
		0x7118c00204100682, 0x2223440021040c00,
		0xa208018c08020142, 0x4404004010200,
		0x14840004802000, 0x204016024100401,
		0x23021a0005451020, 0x204222022c10410,
		0x122010002002b0, 0x2501000022200,
		0x84002804001800a1, 0x1002080800060a00,
		0x40018020120220, 0x41108881004a0100,
		0x800c041410224502, 0x4001020080006403,
		0x205091140081002, 0x491210901c001808,
		0x400084048001000, 0x8824200910800,
		0xca00400408228102, 0x2042240800221200,
		0x54082081000405, 0x1010202004291,
		0x4040a40920100100, 0x4802060101082c10,
		0x208002623100105, 0x1000e2c084040010,
		0x202302400682008a, 0x20820c50024a0c10,
		0x200c20020c090100, 0x684010822028800,
		0x400e002101482012, 0x800804218044242,
		0x8a0040201008820, 0xc000000024420200,
		0x3404102090c20200, 0x8000840810104981,
		0x80330810d0009101, 0x4011001020084,
	},
	shifts: []uint{
		6, 5, 5, 5, 5, 5, 5, 6,
		5, 5, 5, 5, 5, 5, 5, 5,
		5, 5, 7, 7, 7, 7, 5, 5,
		5, 5, 7, 9, 9, 7, 5, 5,
		5, 5, 7, 9, 9, 7, 5, 5,
		5, 5, 7, 7, 7, 7, 5, 5,
		5, 5, 5, 5, 5, 5, 5, 5,
		6, 5, 5, 5, 5, 5, 5, 6,
	},
	dirs: []direction{northEast, southEast, northWest, southWest},
}

var knightAttacks []bd.BitBoard
var kingAttacks []bd.BitBoard

func init() {
	// Generate masks, offsets and attacks tables for rook and bishop.
	for _, tbl := range []*tables{&rookTables, &bishopTables} {
		tbl.masks = make([]bd.BitBoard, bd.NumSquares)
		tbl.offsets = make([]int, bd.NumSquares)
		tbl.attacks = make([]bd.BitBoard, 0)
		for sq := bd.Square(0); sq < bd.NumSquares; sq++ {
			tbl.masks[sq] = 0
			for _, d := range tbl.dirs {
				tbl.masks[sq] |= d.maskBits(sq)
			}
			tmpAttacks := genAttackTable(tbl.dirs, sq, tbl.shifts[sq], tbl.magics[sq])
			tbl.offsets[sq] = len(tbl.attacks)
			tbl.attacks = append(tbl.attacks, tmpAttacks...)
		}
	}

	knightAttacks = make([]bd.BitBoard, bd.NumSquares)
	for r := bd.Row(0); r < 8; r++ {
		for c := bd.Col(0); c < 8; c++ {
			sq := bd.MakeSquare(r, c)
			knightAttacks[sq] = SetBit(r-1, c-2) |
				SetBit(r+1, c-2) |
				SetBit(r+2, c-1) |
				SetBit(r+2, c+1) |
				SetBit(r+1, c+2) |
				SetBit(r-1, c+2) |
				SetBit(r-2, c+1) |
				SetBit(r-2, c-1)
		}
	}

	kingAttacks = make([]bd.BitBoard, bd.NumSquares)
	for r := bd.Row(0); r < 8; r++ {
		for c := bd.Col(0); c < 8; c++ {
			sq := bd.MakeSquare(r, c)
			kingAttacks[sq] = SetBit(r-1, c-1) |
				SetBit(r, c-1) |
				SetBit(r+1, c-1) |
				SetBit(r+1, c) |
				SetBit(r+1, c+1) |
				SetBit(r, c+1) |
				SetBit(r-1, c+1) |
				SetBit(r-1, c)
		}
	}
}

func SetBit(r bd.Row, c bd.Col) bd.BitBoard {
	if r >= bd.Row(0) && r <= bd.Row(7) && c >= bd.Col(0) && c <= bd.Col(7) {
		return bd.BitBoard(1) << uint(bd.MakeSquare(r, c))
	}
	return bd.BitBoard(0)
}

func combinedOccs(dirs []direction, sq bd.Square) []bd.BitBoard {
	occs := make([]bd.BitBoard, 0)
	for _, d := range dirs {
		dirOcc := d.occs(sq)
		if len(dirOcc) == 0 {
			continue
		}
		if len(occs) == 0 {
			occs = append(occs, dirOcc...)
			continue
		}
		tmp := make([]bd.BitBoard, 0)
		for _, bb := range dirOcc {
			for _, occ := range occs {
				tmp = append(tmp, bb|occ)
			}
		}
		occs = tmp
	}
	return occs
}

func genAttackTable(dirs []direction, sq bd.Square, shift_bits uint, magic bd.BitBoard) []bd.BitBoard {
	occs := combinedOccs(dirs, sq)
	attacks := make([]bd.BitBoard, 0)
	for _, occ := range occs {
		var attack bd.BitBoard
		for _, dir := range dirs {
			attack |= dir.attack(sq, occ)
		}
		attacks = append(attacks, attack)
	}
	invalidAttack := bd.BitBoard(^uint64(0))
	table := make([]bd.BitBoard, 1<<shift_bits)
	for i := 0; i < len(table); i++ {
		table[i] = bd.BitBoard(invalidAttack)
	}
	for i, occ := range occs {
		attack := attacks[i]
		offset := (occ * magic) >> (64 - shift_bits)
		if table[offset] == invalidAttack || table[offset] == attack {
			table[offset] = attack
		} else {
			panic("collision")
		}
	}
	return table
}
