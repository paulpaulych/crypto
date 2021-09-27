package core

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"
)

type keys struct{ c, d *big.Int }
type testCase = struct {
	p     *big.Int
	alice keys
	bob   keys
	msg   *big.Int
}

func TestShamir(t *testing.T) {
	tests := append(
		generateCasesForAandB(
			37,
			keys{big.NewInt(5), big.NewInt(29)},
			keys{big.NewInt(7), big.NewInt(31)},
		),
		generateCasesForAandB(
			101,
			keys{big.NewInt(31), big.NewInt(71)},
			keys{big.NewInt(53), big.NewInt(17)},
		)...,
	)
	for _, tt := range tests {
		alice := Alice{p: tt.p, c: tt.alice.c, d: tt.alice.d}
		bob := Bob{p: tt.p, c: tt.bob.c, d: tt.bob.d}
		testName := fmt.Sprintf("testCase for A={p=%v,c=%v,d=%v}, B={p=%v,c=%v,d=%v}, msg=%d",
			alice.p, alice.c, alice.d, bob.p, bob.c, bob.d, tt.msg,
		)
		t.Run(testName, func(t *testing.T) {
			step1Out := alice.Step1(tt.msg)
			t.Logf("step 1 output: %v", step1Out)
			step2Out := bob.Step2(step1Out)
			t.Logf("step 2 output: %v", step2Out)
			step3Out := alice.Step3(step2Out)
			t.Logf("step 3 output: %v", step3Out)
			decoded := bob.Decode(step3Out)

			if !reflect.DeepEqual(decoded, tt.msg) {
				t.Errorf("decoded = %v, want %v", decoded, tt.msg)
			}
		})
	}
}

func generateCasesForAandB(p int64, alice, bob keys) []testCase {
	cases := make([]testCase, p-2)
	for i := int64(0); i < p-2; i++ {
		cases[i] = testCase{
			p:     big.NewInt(p),
			alice: alice,
			bob:   bob,
			msg:   big.NewInt(i + 1),
		}
	}
	return cases
}
