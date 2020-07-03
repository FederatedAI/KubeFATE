package orm

import "testing"

func TestMysql_Setup(t *testing.T) {
	InitConfigForTest()
	tests := []struct {
		name    string
		e       *Mysql
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "The test database can be connected normally",
			e:       new(Mysql),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Mysql{}
			if err := e.Setup(); (err != nil) != tt.wantErr {
				t.Errorf("Mysql.Setup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
