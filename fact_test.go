package testfact
import "testing"

func TestFact(t *testing.T){
	testCases:=[]struct{
		inp, oup int
	}{
		{0,1},
		{1,1},
		{2,2},
		{3,6},
		{4,24},
		{-1,0},
	}

	for ind := range testCases{
		if v:=Fact(testCases[ind].inp);v!=testCases[ind].oup{
			t.Errorf("FAILED wanted %d got %d for Fact(%d)\n", testCases[ind].oup, v, testCases[ind].inp)
		}
	}

}



