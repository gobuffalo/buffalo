package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_hasYarnScript(t *testing.T) {
	type args struct {
		script string
		data   string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{name: "Find scripts", args: args{script: "buffalo:dev", data: withScripts}, want: true, wantErr: false},
		{name: "Don't find scripts", args: args{script: "buffalo:dev", data: withoutScripts}, want: false, wantErr: false},
		{name: "Parse error", args: args{script: "buffalo:dev", data: invalidJSON}, want: false, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)
			got, err := hasYarnScript(tt.args.script, []byte(tt.args.data))
			if (err != nil) != tt.wantErr {
				r.Fail("hasYarnScript() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			r.Equal(tt.want, got)
		})
	}
}

var withScripts = `{
	"name": "buffalo",
	"scripts": {
	  "buffalo:build": "webpack -p --progress",
	  "buffalo:dev": "webpack --watch",
	  "buffalo:test": ""
	},
	"dependencies": {
	  "jquery-ujs": "~1.2.2"
	},
	"devDependencies": {
	  "webpack-manifest-plugin": "~2.0.0"
	}
}
`

var withoutScripts = `{
	"name": "buffalo",
	"dependencies": {
	  "jquery-ujs": "~1.2.2"
	},
	"devDependencies": {
	  "webpack-manifest-plugin": "~2.0.0"
	}
}
`

var invalidJSON = `{
	"name": "buffalo",
	"scripts": {
	  "buffalo:build": "webpack -p --progress",
	  "buffalo:dev": "webpack --watch"XXXXXX,
	  "buffalo:test": ""
	},
	"dependencies": {
	  "jquery-ujs": "~1.2.2"
	},
	"devDependencies": {
	  "webpack-manifest-plugin": "~2.0.0"
	}
}
`
