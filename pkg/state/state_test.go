// nolint: err113
package state

import (
	"errors"
	"reflect"
	"testing"
)

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SStateError{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
	}
}

func TestNewBoolState(t *testing.T) {
	tests := []struct {
		name string
		want IState
	}{
		{
			name: "new_bool_state",
			want: NewBoolState(),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := NewBoolState(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBoolState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sState_Enable(t *testing.T) {
	type fields struct {
		fEnabled bool
	}
	type args struct {
		f IStateF
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "state_enable_enabled_without_f",
			fields: fields{
				fEnabled: true,
			},
			args:    args{f: nil},
			wantErr: true,
		},
		{
			name: "state_enable_enabled_with_f",
			fields: fields{
				fEnabled: true,
			},
			args:    args{f: func() error { return nil }},
			wantErr: true,
		},
		{
			name: "state_enable_enabled_with_f_error",
			fields: fields{
				fEnabled: true,
			},
			args:    args{f: func() error { return errors.New("some error") }}, //nolint:err113
			wantErr: true,
		},
		{
			name: "state_enable_disabled_without_f",
			fields: fields{
				fEnabled: false,
			},
			args:    args{f: nil},
			wantErr: false,
		},
		{
			name: "state_enable_disabled_with_f",
			fields: fields{
				fEnabled: false,
			},
			args:    args{f: func() error { return nil }},
			wantErr: false,
		},
		{
			name: "state_enable_disabled_with_f_error",
			fields: fields{
				fEnabled: false,
			},
			args:    args{f: func() error { return errors.New("some error") }}, //nolint:err113
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			p := &sState{
				fEnabled: tt.fields.fEnabled,
			}
			if err := p.Enable(tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("sState.Enable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_sState_Disable(t *testing.T) {
	type fields struct {
		fEnabled bool
	}
	type args struct {
		f IStateF
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "state_disable_enabled_without_f",
			fields: fields{
				fEnabled: true,
			},
			args:    args{f: nil},
			wantErr: false,
		},
		{
			name: "state_disable_enabled_with_f",
			fields: fields{
				fEnabled: true,
			},
			args:    args{f: func() error { return nil }},
			wantErr: false,
		},
		{
			name: "state_disable_enabled_with_f_error",
			fields: fields{
				fEnabled: true,
			},
			args:    args{f: func() error { return errors.New("some error") }}, //nolint:err113
			wantErr: true,
		},
		{
			name: "state_disable_disabled_without_f",
			fields: fields{
				fEnabled: false,
			},
			args:    args{f: nil},
			wantErr: true,
		},
		{
			name: "state_disable_disabled_with_f",
			fields: fields{
				fEnabled: false,
			},
			args:    args{f: func() error { return nil }},
			wantErr: true,
		},
		{
			name: "state_disable_disabled_with_f_error",
			fields: fields{
				fEnabled: false,
			},
			args:    args{f: func() error { return errors.New("some error") }}, //nolint:err113
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			p := &sState{
				fEnabled: tt.fields.fEnabled,
			}
			if err := p.Disable(tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("sState.Disable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
