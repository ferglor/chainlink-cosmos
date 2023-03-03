package cosmwasm

import (
	"strings"
	"testing"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/stretchr/testify/assert"
)

var testConfig = types.ContractConfig{
	ConfigCount: 1,
	Signers: []types.OnchainPublicKey{
		[]byte{28, 69, 220, 255, 145, 161, 41, 242, 208, 125, 181, 65, 174, 4, 255, 77, 61, 37, 134, 54, 130, 230, 11, 172, 175, 166, 100, 99, 69, 122, 138, 128},
		[]byte{245, 123, 138, 66, 133, 4, 54, 37, 129, 106, 119, 250, 131, 43, 174, 81, 139, 147, 232, 202, 3, 177, 159, 111, 170, 76, 143, 137, 250, 67, 69, 125},
		[]byte{146, 9, 73, 38, 35, 203, 190, 72, 88, 255, 219, 63, 192, 95, 118, 108, 236, 15, 144, 179, 62, 29, 223, 222, 245, 61, 164, 73, 208, 76, 72, 59},
		[]byte{243, 96, 118, 131, 178, 167, 101, 157, 94, 246, 73, 127, 240, 101, 36, 36, 102, 191, 168, 19, 47, 217, 47, 45, 245, 233, 119, 230, 53, 102, 153, 74},
	},
	Transmitters: []types.Account{
		"wasm1cd65xyq076dm9cw3xxqtdh4d6ypzug0edd9958",
		"wasm19ctxyyc49cf42cfx3vvj3kmkrgzflw72h4afvv",
		"wasm1ysjdehnf3a3kpndx74yyg6ry90258y4z5vawjz",
		"wasm1fucynrfkrt684pm8jrt8la5h2csvs5cnldcgqc",
	},
	F:                     1,
	OnchainConfig:         []byte{},
	OffchainConfigVersion: 2,
	OffchainConfig:        []byte{8, 128, 168, 214, 185, 7, 16, 128, 148, 235, 220, 3, 24, 128, 148, 235, 220, 3, 32, 128, 202, 181, 238, 1, 40, 128, 168, 214, 185, 7, 48, 3, 58, 4, 1, 1, 1, 1, 66, 32, 136, 99, 127, 179, 112, 251, 210, 5, 179, 14, 165, 40, 178, 72, 177, 95, 153, 70, 125, 163, 116, 227, 213, 217, 77, 208, 194, 7, 151, 116, 212, 160, 66, 32, 88, 247, 149, 158, 51, 177, 58, 136, 11, 0, 206, 196, 97, 202, 194, 189, 249, 27, 54, 211, 54, 208, 184, 216, 15, 61, 233, 177, 39, 97, 213, 69, 66, 32, 56, 122, 236, 208, 44, 127, 77, 118, 178, 31, 172, 160, 227, 177, 171, 61, 137, 247, 136, 89, 211, 54, 157, 119, 235, 17, 213, 190, 36, 80, 68, 233, 66, 32, 89, 136, 237, 203, 198, 53, 101, 102, 194, 23, 36, 136, 10, 131, 164, 242, 82, 56, 135, 70, 71, 252, 228, 74, 22, 145, 234, 199, 176, 124, 240, 110, 74, 52, 49, 50, 68, 51, 75, 111, 111, 87, 69, 105, 57, 68, 107, 85, 68, 56, 66, 122, 110, 109, 89, 70, 119, 109, 70, 110, 107, 115, 52, 69, 111, 65, 53, 70, 74, 110, 102, 114, 67, 77, 88, 57, 54, 121, 112, 53, 74, 120, 104, 71, 99, 89, 74, 52, 49, 50, 68, 51, 75, 111, 111, 87, 81, 107, 52, 54, 68, 113, 52, 103, 122, 104, 70, 65, 78, 68, 77, 100, 85, 51, 88, 102, 113, 90, 57, 69, 117, 97, 78, 110, 109, 53, 85, 120, 122, 116, 75, 72, 100, 118, 107, 75, 55, 102, 86, 117, 74, 52, 49, 50, 68, 51, 75, 111, 111, 87, 67, 90, 113, 106, 121, 115, 87, 57, 81, 104, 118, 49, 68, 72, 82, 76, 69, 117, 83, 81, 85, 88, 98, 109, 49, 80, 113, 88, 88, 65, 74, 102, 80, 57, 71, 97, 112, 99, 99, 107, 74, 105, 119, 109, 74, 52, 49, 50, 68, 51, 75, 111, 111, 87, 66, 97, 106, 119, 111, 106, 72, 121, 109, 102, 86, 66, 68, 82, 56, 120, 67, 104, 118, 84, 55, 76, 52, 56, 50, 50, 74, 53, 107, 89, 70, 107, 65, 113, 115, 80, 104, 112, 81, 78, 121, 81, 74, 114, 82, 16, 0, 0, 0, 0, 0, 45, 198, 192, 0, 0, 0, 0, 0, 0, 0, 0, 88, 128, 225, 235, 23, 96, 128, 225, 235, 23, 104, 128, 225, 235, 23, 112, 128, 225, 235, 23, 120, 128, 225, 235, 23, 130, 1, 140, 1, 10, 32, 181, 88, 224, 203, 224, 168, 227, 94, 172, 196, 63, 164, 146, 136, 91, 186, 91, 111, 129, 189, 44, 156, 168, 196, 184, 11, 66, 188, 195, 115, 137, 20, 18, 32, 26, 27, 247, 83, 194, 80, 93, 171, 82, 4, 178, 132, 241, 254, 253, 125, 196, 235, 131, 246, 48, 70, 100, 201, 194, 244, 60, 191, 66, 86, 102, 96, 26, 16, 198, 216, 18, 243, 142, 77, 172, 75, 15, 162, 86, 225, 152, 154, 30, 61, 26, 16, 178, 15, 108, 85, 146, 35, 250, 89, 217, 237, 111, 20, 87, 215, 173, 142, 26, 16, 43, 40, 56, 124, 16, 177, 93, 71, 182, 227, 75, 241, 146, 174, 93, 244, 26, 16, 92, 222, 61, 101, 97, 47, 37, 21, 39, 37, 172, 18, 140, 109, 236, 20},
}

func TestConfigDigester(t *testing.T) {
	d := NewOffchainConfigDigester(
		"ibiza-808",
		MustAccAddress("wasm1cd65xyq076dm9cw3xxqtdh4d6ypzug0edd9958"),
	)

	digest, err := d.ConfigDigest(testConfig)
	assert.NoError(t, err)
	assert.Equal(t, "000289b55121341b1ff99cc8e15659fb8de14fca52a695b2b269a7fb94059b9f", digest.Hex())
}

func TestConfigDigester_InvalidChainID(t *testing.T) {
	d := NewOffchainConfigDigester(
		strings.Repeat("a", 256), // chain ID is too long
		MustAccAddress("wasm1cd65xyq076dm9cw3xxqtdh4d6ypzug0edd9958"),
	)

	_, err := d.ConfigDigest(testConfig)
	assert.Error(t, err)
}
