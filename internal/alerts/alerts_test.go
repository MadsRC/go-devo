package alerts

import “testing”
 
func TestNewClient(t *testing.T) {
	client := NewClient()
	if NewClient == nil {
		t.Fail()
	}
}